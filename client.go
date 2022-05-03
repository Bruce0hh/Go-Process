package Go_RPC

import (
	"Go-RPC/codec"
	"errors"
	"fmt"
	"io"
	"sync"
)

// 封装结构体Call来承载一次RPC的调用信息
type Call struct {
	Seq           uint64
	ServiceMethod string
	Args          interface{}
	Reply         interface{}
	Error         error
	Done          chan *Call
}

func (c *Call) done() {
	c.Done <- c
}

type Client struct {
	cc       codec.Codec // 消息编解码器
	opt      *Option
	sending  sync.Mutex   // 互斥锁，保证请求有序
	header   codec.Header // 消息头，只有在请求发送时才需要
	mu       sync.Mutex
	seq      uint64           // 发送请求的编号
	pending  map[uint64]*Call // 存储未处理完的请求
	closing  bool             // 主动关闭
	shutdown bool             // 有错误发生导致关闭
}

var _ io.Closer = (*Client)(nil)

var ErrShutdown = errors.New("connetction is shutdown")

// 关闭连接
func (client *Client) Close() error {
	client.mu.Lock()
	defer client.mu.Unlock()
	if client.closing {
		return ErrShutdown
	}
	client.closing = true
	return client.cc.Close()
}

// 判断客户端是否在工作
func (client *Client) IsAvailable() bool {
	client.mu.Lock()
	defer client.mu.Unlock()
	return !client.shutdown && !client.closing
}

/**
 * @description: 将参数call添加到client.pending中，更新client.seq
 * @param {*Call} call
 * @return {*}
 */
func (client *Client) registerCall(call *Call) (uint64, error) {
	client.mu.Lock()
	defer client.mu.Unlock()
	if client.closing || client.shutdown {
		return 0, ErrShutdown
	}
	call.Seq = client.seq
	client.pending[call.Seq] = call
	client.seq++
	return call.Seq, nil
}

/**
 * @description: 根据 seq，从 client.pending 中移除对应的 call，并返回
 * @param {uint64} seq
 * @return {*}
 */
func (client *Client) removeCall(seq uint64) *Call {
	client.mu.Lock()
	defer client.mu.Unlock()
	call := client.pending[seq]
	delete(client.pending, seq)
	return call
}

/**
 * @description: 服务端或客户端发生错误时调用，将 shutdown 设置为 true，且将错误信息通知所有 pending 状态的 call
 * @param {error} err
 * @return {*}
 */
func (client *Client) termianate(err error) {
	client.sending.Lock()
	defer client.sending.Unlock()
	client.mu.Lock()
	defer client.mu.Unlock()
	client.shutdown = true
	for _, call := range client.pending {
		call.Error = err
		call.done()
	}
}

/**
 * @description: 接收
 */
func (client *Client) receive() {
	var err error
	for err == nil {
		var h codec.Header
		if err = client.cc.ReaderHeader(&h); err != nil {
			break
		}
		call := client.removeCall(h.Seq)
		switch {
		case call == nil:
			err = client.cc.ReadBody(nil)
		case h.Error != "":
			call.Error = fmt.Errorf(h.Error)
		}
	}
}
