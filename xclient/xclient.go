package xclient

import (
	"context"
	"gorpc"
	"io"
	"reflect"
	"sync"
)

/*
 支持负载均衡的客户端XCLIENT，重新包装client
*/

type XClient struct {
	discovery Discovery     // 服务发现实例
	mode      SelectMode    // 负载均衡模式
	opt       *gorpc.Option // 协议选项
	mu        sync.Mutex
	clients   map[string]*gorpc.Client
}

func NewXClient(d Discovery, mode SelectMode, opt *gorpc.Option) *XClient {
	return &XClient{
		discovery: d,
		mode:      mode,
		opt:       opt,
		clients:   make(map[string]*gorpc.Client),
	}
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

// Broadcast 广播：如果任意一个服务实例发生错误，则返回其中一个错误；如果调用成功，则返回其中一个结果
func (x *XClient) Broadcast(ctx context.Context, serviceMethod string, args, reply interface{}) error {
	servers, err := x.discovery.GetAll()
	if err != nil {
		return err
	}

	var (
		wg sync.WaitGroup
		mu sync.Mutex
		e  error
	)

	replyDone := reply == nil
	ctx, cancel := context.WithCancel(ctx)
	// 为了提升性能，请求是并发的
	for _, rpcAddr := range servers {
		wg.Add(1)
		go func(rpcAddr string) {
			defer wg.Done()

			var clonedReply interface{}
			if reply != nil {
				clonedReply = reflect.New(reflect.ValueOf(reply).Elem().Type()).Interface()
			}
			err := x.call(rpcAddr, ctx, serviceMethod, args, clonedReply)
			// 使用互斥锁波阿虎增error和reply能被正确赋值
			mu.Lock()
			if err != nil && e == nil {
				e = err
				cancel() // 借助cancel，确保有错误发生时，能快速失败
			}
			if err == nil && !replyDone {
				reflect.ValueOf(reply).Elem().Set(reflect.ValueOf(clonedReply).Elem())
				replyDone = true
			}
			mu.Unlock()
		}(rpcAddr)
	}
	wg.Wait()
	return e
}
