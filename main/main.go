package main

import (
	"context"
	"gorpc"
	"gorpc/xclient"
	"log"
	"net"
	"sync"
	"time"
)

type Foo int

type Args struct {
	Num1, Num2 int
}

func (f *Foo) Sum(args Args, reply *int) error {
	*reply = args.Num1 + args.Num2
	return nil
}
func (f *Foo) Sleep(args Args, reply *int) error {
	time.Sleep(time.Second * time.Duration(args.Num1)) // 测试XClient的超时机制
	*reply = args.Num1 + args.Num2
	return nil
}

func foo(x *xclient.XClient, ctx context.Context, typ, servicedMethod string, args *Args) {
	var (
		reply int
		err   error
	)
	switch typ {
	case "call":
		err = x.Call(ctx, servicedMethod, args, &reply)
	case "broadcast":
		err = x.Broadcast(ctx, servicedMethod, args, &reply)
	}
	if err != nil {
		log.Printf("%v, %v error: %v", typ, servicedMethod, err)
	} else {
		log.Printf("%v, %v success: %v + %v = %v", typ, servicedMethod, args.Num1, args.Num2, reply)
	}
}

func startServer(addr chan string) {

	var foo Foo
	server := gorpc.NewServer()
	if err := server.Register(&foo); err != nil {
		log.Fatal("register error: ", err)
	}

	l, err := net.Listen("tcp", ":0")
	if err != nil {
		log.Fatalf("network error: %+v", err)
	}
	log.Printf("start rpc server on: %+v", l.Addr())
	addr <- l.Addr().String()
	server.Accept(l)
}

// 调用单服务实例
func callSingle(addr1, addr2 string) {
	d := xclient.NewMultiServerDiscovery([]string{"tcp@" + addr1, "tcp@" + addr2})
	x := xclient.NewXClient(d, xclient.RandomSelect, nil)
	defer func() { _ = x.Close() }()

	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			foo(x, context.Background(), "call", "Foo.Sum", &Args{
				Num1: i,
				Num2: i * i,
			})
		}(i)
	}
	wg.Wait()
}

// 调用所有服务实例
func broadcast(addr1, addr2 string) {
	d := xclient.NewMultiServerDiscovery([]string{"tcp@" + addr1, "tcp@" + addr2})
	x := xclient.NewXClient(d, xclient.RandomSelect, nil)
	defer func() { _ = x.Close() }()

	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			foo(x, context.Background(), "broadcast", "Foo.Sum", &Args{Num1: i, Num2: i * i})
			ctx, _ := context.WithTimeout(context.Background(), time.Second*2)
			foo(x, ctx, "broadcast", "Foo.Sleep", &Args{Num1: i, Num2: i * i})
		}(i)
	}
	wg.Wait()
}

func call(addrCh chan string) {
	c, _ := gorpc.DialHTTP("tcp", <-addrCh)
	defer func() { _ = c.Close() }()
	time.Sleep(time.Second)

	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			args := &Args{i, i * i}
			var reply int
			if err := c.Call(context.Background(), "Foo.Sum", args, &reply); err != nil {
				log.Fatal("call Foo.Sum error: ", err)
			}
			log.Printf("%d + %d = %d", args.Num1, args.Num2, reply)
		}(i)
	}
	wg.Wait()
}

func main() {

	log.SetFlags(0) // 不输出任何的日志信息头

	ch1 := make(chan string)
	ch2 := make(chan string)

	go startServer(ch1)
	go startServer(ch2)

	addr1 := <-ch1
	addr2 := <-ch2

	time.Sleep(time.Second)
	callSingle(addr1, addr2)
	broadcast(addr1, addr2)
}
