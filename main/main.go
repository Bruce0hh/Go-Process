package main

import (
	"context"
	"gorpc"
	"log"
	"net"
	"net/http"
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

func startServer(addr chan string) {

	var foo Foo
	if err := gorpc.Register(&foo); err != nil {
		log.Fatal("register error: ", err)
	}

	l, err := net.Listen("tcp", ":9999")
	if err != nil {
		log.Fatalf("network error: %+v", err)
	}
	log.Printf("start rpc server on: %+v", l.Addr())
	gorpc.HandleHTTP()
	addr <- l.Addr().String()
	//gorpc.Accept(l)
	_ = http.Serve(l, nil)
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
	addr := make(chan string)
	go call(addr)
	startServer(addr) // 开启服务端

}
