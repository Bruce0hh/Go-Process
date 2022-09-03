package gorpc

import (
	"context"
	"net"
	"strings"
	"testing"
	"time"
)

// 测试处理超时
type Bar int

func (b Bar) Timeout(argv int, reply *int) error {
	time.Sleep(time.Second * 2)
	return nil
}

func startServer(addr chan string) {
	var b Bar
	_ = Register(&b)
	l, _ := net.Listen("tcp", ":0")
	addr <- l.Addr().String()
	Accept(l)
}

func TestClient_Call(t *testing.T) {
	t.Parallel()
	var reply *int
	addrCh := make(chan string)
	go startServer(addrCh)
	addr := <-addrCh
	time.Sleep(time.Second)
	// 客户端设置超时时间为1s，服务端无限制
	client, _ := Dial("tcp", addr)
	ctx, _ := context.WithTimeout(context.Background(), time.Second)
	// 客户端无限制，服务端设置超时时间为1s
	client2, _ := Dial("tcp", addr, &Option{HandleTimeout: time.Second})
	ctx2 := context.Background()

	type args struct {
		ctx           context.Context
		serviceMethod string
		args          interface{}
		reply         interface{}
	}
	tests := []struct {
		name    string
		client  *Client
		args    args
		wantErr bool
	}{
		{"client timeout", client,
			args{
				ctx:           ctx,
				serviceMethod: "Bar.Timeout",
				args:          1,
				reply:         &reply,
			}, true},
		{"server handle timeout", client2,
			args{
				ctx:           ctx2,
				serviceMethod: "Bar.Timeout",
				args:          1,
				reply:         &reply,
			}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "client timeout" {
				err := tt.client.Call(tt.args.ctx, tt.args.serviceMethod, tt.args.args, tt.args.reply)
				if (err != nil && strings.Contains(err.Error(), ctx.Err().Error())) != tt.wantErr {
					t.Errorf("Call() error = %v, context error = %v, wantErr %v",
						err.Error(), ctx.Err().Error(), tt.wantErr)
				}
			}
			if tt.name == "server handle timeout" {
				err := tt.client.Call(tt.args.ctx, tt.args.serviceMethod, tt.args.args, tt.args.reply)
				if (err != nil && strings.Contains(err.Error(), "handle timeout")) != tt.wantErr {
					t.Errorf("Call() error = %v, wantErr %v", err.Error(), tt.wantErr)
				}
			}
		})
	}
}

/*
	测试连接超时，NewClient耗时2s，ConnectionTimeout分别设置1s和0s
*/
func Test_dialTimeout(t *testing.T) {
	t.Parallel()
	l, _ := net.Listen("tcp", ":0")
	f := func(conn net.Conn, opt *Option) (client *Client, err error) {
		_ = conn.Close()
		time.Sleep(time.Second * 2)
		return nil, nil
	}

	type args struct {
		f       newClientFunc
		network string
		address string
		opts    []*Option
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"test-timeout", args{
			f:       f,
			network: "tcp",
			address: l.Addr().String(),
			opts:    []*Option{{ConnectTimeout: time.Second}},
		}, true},
		{
			"test-no-timeout", args{
				f:       f,
				network: "tcp",
				address: l.Addr().String(),
				opts:    []*Option{{ConnectTimeout: 0}},
			}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := dialTimeout(tt.args.f, tt.args.network, tt.args.address, tt.args.opts...)
			if tt.name == "test-timeout" {
				if (err != nil && strings.Contains(err.Error(), "connect timeout")) != tt.wantErr {
					t.Errorf("dialTimeout() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
			}
			if tt.name == "test-no-timeout" {
				if (err == nil) != tt.wantErr {
					return
				}
			}
		})
	}
}
