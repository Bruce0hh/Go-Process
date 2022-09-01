package main

import (
	"fmt"
	"gorpc"
	"gorpc/client"
	"log"
	"net"
	"sync"
	"time"
)

func startServer(addr chan string) {
	l, err := net.Listen("tcp", ":0")
	if err != nil {
		log.Fatalf("network error: %+v", err)
	}
	log.Printf("start rpc server on: %+v", l.Addr())
	addr <- l.Addr().String()
	gorpc.Accept(l)
}

func main() {

	log.SetFlags(0) // 不输出任何的日志信息头
	addr := make(chan string)
	go startServer(addr)               // 开启服务端
	c, _ := client.Dial("tcp", <-addr) // 开启客户端通信
	defer func() { c.Close() }()
	time.Sleep(time.Second)

	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			args := fmt.Sprintf("gorpc req %d", i)
			var reply string
			if err := c.Call("Foo.Sum", args, &reply); err != nil {
				log.Fatal("call Foo.Sum error: ", err)
			}
			log.Println("reply: ", reply)
		}(i)
	}
	wg.Wait()
}
