package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"gorpc"
	"gorpc/codec"
	"io"
	"log"
	"net"
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
func (client *Client) Close() error {
	client.mu.Lock()
	defer client.mu.Unlock()
	if client.closing {
		return ErrShutdown
	}
	client.closing = true
	return client.cc.Close()
}

var _ io.Closer = (*Client)(nil)

func (client *Client) IsAvailable() bool {
	client.mu.Lock()
	defer client.mu.Unlock()
	return !client.shutdown && !client.closing
}

/*
  Call相关的三个方法
*/

// 将参数call添加到client.pending中，并更新client.seq
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

// 根据seq，将call从pending中移除
func (client *Client) removeCall(seq uint64) *Call {
	client.mu.Lock()
	defer client.mu.Unlock()
	call := client.pending[seq]
	delete(client.pending, seq)
	return call
}

// 发生错误时，将错误信息通知到所有pending的call
func (client *Client) terminalCalls(err error) {
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

// 接收响应
func (client *Client) receive() {
	var err error
	for err == nil {
		var h codec.Header
		if err = client.cc.ReadHeader(&h); err != nil {
			break
		}
		call := client.removeCall(h.Seq)
		switch {
		case call == nil: // call不存在
			err = client.cc.ReadBody(nil)
		case h.Error != "": // 服务端处理出错，即h.Error不为空
			call.Error = fmt.Errorf(h.Error)
			err = client.cc.ReadBody(nil)
			call.done()
		default: // 正常处理
			err = client.cc.ReadBody(call.Reply)
			if err != nil {
				call.Error = errors.New("reading body " + err.Error())
			}
			call.done()
		}
	}
	client.terminalCalls(err)
}

// NewClient Client实例处理
func NewClient(conn net.Conn, opt *gorpc.Option) (*Client, error) {
	f := codec.NewCodecFuncMap[opt.CodecType]
	if f == nil {
		err := fmt.Errorf("invalid codec type %s", opt.CodecType)
		log.Println("rpc client: codec error: ", err)
		return nil, err
	}

	// 发送options
	if err := json.NewEncoder(conn).Encode(opt); err != nil {
		log.Println("rpc client: options error: ", err)
		_ = conn.Close()
		return nil, err
	}
	return newClientCodec(f(conn), opt), nil
}

// 根据编解码将Options信息发送给Server，调用receive()
func newClientCodec(cc codec.Codec, opt *gorpc.Option) *Client {
	client := &Client{
		cc:      cc,
		opt:     opt,
		seq:     1,
		pending: make(map[uint64]*Call),
	}
	go client.receive()
	return client
}

// 解析options
func parseOptions(opts ...*gorpc.Option) (*gorpc.Option, error) {

	if len(opts) == 0 || opts[0] == nil {
		return gorpc.DefaultOption, nil
	}
	if len(opts) != 1 {
		return nil, errors.New("number of options is more than 1")
	}
	opt := opts[0]
	opt.MagicNumber = gorpc.DefaultOption.MagicNumber
	if opt.CodecType == "" {
		opt.CodecType = gorpc.DefaultOption.CodecType
	}
	return opt, nil
}

// Dial 通信
func Dial(network, address string, opts ...*gorpc.Option) (client *Client, err error) {

	opt, err := parseOptions(opts...)
	if err != nil {
		return nil, err
	}
	conn, err := net.Dial(network, address)
	if err != nil {
		return nil, err
	}
	defer func() {
		if client == nil {
			_ = conn.Close()
		}
	}()
	return NewClient(conn, opt)
}

func (client *Client) send(call *Call) {

	client.sending.Lock()
	defer client.sending.Unlock()
	// 注册call
	seq, err := client.registerCall(call)
	if err != nil {
		call.Error = err
		call.done()
		return
	}
	// 准备request header
	client.header.ServiceMethod = call.ServiceMethod
	client.header.Seq = seq
	client.header.Error = ""

	// 解码并发送请求
	if err := client.cc.Write(&client.header, call.Args); err != nil {
		call := client.removeCall(seq)
		if call != nil {
			call.Error = err
			call.done()
		}
	}
}

// Go 异步Call实例
func (client *Client) Go(serviceMethod string, args, reply interface{}, done chan *Call) *Call {
	if done != nil {
		done = make(chan *Call, 10)
	} else if cap(done) == 0 {
		log.Panic("rpc client: done channel is unbuffered")
	}

	call := &Call{
		ServiceMethod: serviceMethod,
		Args:          args,
		Reply:         reply,
		Done:          done,
	}
	client.send(call)
	return call
}

// Call 阻塞，同步接口
func (client *Client) Call(serviceMethod string, args, reply interface{}) error {
	call := <-client.Go(serviceMethod, args, reply, make(chan *Call, 1)).Done
	return call.Error
}
