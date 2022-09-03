package xclient

import (
	"context"
	"gorpc"
	"io"
	"sync"
)

/*
 支持负载均衡的客户端XCLIENT
*/

type XClient struct {
	discovery Discovery
	mode      SelectMode
	opt       *gorpc.Option
	mu        sync.Mutex
	clients   map[string]*gorpc.Client
}

func (x *XClient) Close() error {
	x.mu.Lock()
	defer x.mu.Unlock()
	for key, client := range x.clients {
		_ = client.Close()
		delete(x.clients, key)
	}
	return nil
}

var _ io.Closer = (*XClient)(nil)

func (x *XClient) dial(rpcAddr string) (*gorpc.Client, error) {
	x.mu.Lock()
	defer x.mu.Unlock()

	// 如果clients有缓存的Client，检查是否可用；可用则返回Client，不可用则删除
	client, ok := x.clients[rpcAddr]
	if ok && !client.IsAvailable() {
		_ = client.Close()
		delete(x.clients, rpcAddr)
		client = nil
	}
	// 没有一个client的话，创建新的Client
	if client == nil {
		var err error
		client, err = gorpc.XDial(rpcAddr, x.opt)
		if err != nil {
			return nil, err
		}
		x.clients[rpcAddr] = client
	}
	return client, nil
}

func (x *XClient) call(rpcAddr string, ctx context.Context, serviceMethod string, args, reply interface{}) error {
	client, err := x.dial(rpcAddr)
	if err != nil {
		return err
	}
	return client.Call(ctx, serviceMethod, args, reply)
}

func (x *XClient) Call(ctx context.Context, serviceMethod string, args, reply interface{}) error {
	rpcAddr, err := x.discovery.Get(x.mode)
	if err != nil {
		return err
	}
	return x.call(rpcAddr, ctx, serviceMethod, args, reply)
}
