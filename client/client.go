package client

import (
	"errors"
	"gorpc"
	"gorpc/codec"
	"io"
	"sync"
)

// Call 一次RPC调用所需要的信息
type Call struct {
	Seq           uint64
	ServiceMethod string
	Args          interface{}
	Reply         interface{}
	Error         error
	Done          chan *Call // 当调用结束时，会调用 call.done() 通知调用方
}

func (c *Call) done() {
	c.Done <- c
}

type Client struct {
	cc       codec.Codec // 消息解码器
	opt      *gorpc.Option
	sending  sync.Mutex   // 保证请求有序，防止多个报文混淆
	header   codec.Header // 客户端每次只处理一个请求，所以header可以复用
	mu       sync.Mutex
	seq      uint64
	pending  map[uint64]*Call // 存储未处理完的请求
	closing  bool             // 主动关闭
	shutdown bool             // 错误发生
}

var ErrShutdown = errors.New("connection is shutdown")

// Close 关闭连接
func (c *Client) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.closing {
		return ErrShutdown
	}
	c.closing = true
	return c.cc.Close()
}

var _ io.Closer = (*Client)(nil)

func (c *Client) IsAvailable() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return !c.shutdown && !c.closing
}

/*
Call相关的三个方法
*/

func (c *Client) registerCall(call *Call) (uint64, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.closing || c.shutdown {
		return 0, ErrShutdown
	}
	call.Seq = c.seq
	c.pending[call.Seq] = call
	c.seq++
	return call.Seq, nil
}

func (c *Client) removeCall(seq uint64) *Call {
	c.mu.Lock()
	defer c.mu.Unlock()
	call := c.pending[seq]
	delete(c.pending, seq)
	return call
}

func (c *Client) terminalCalls(err error) {
	c.sending.Lock()
	defer c.sending.Unlock()
	c.mu.Lock()
	defer c.mu.Unlock()
	c.shutdown = true
	for _, call := range c.pending {
		call.Error = err
		call.done()
	}
}
