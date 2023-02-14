package main

import (
	"flag"
	"fmt"
	go_cache "gocache"
	"log"
	"net/http"
)

var db = map[string]string{
	"Tom":  "666",
	"Jack": "555",
	"Sam":  "777",
}

// 启动缓存服务器
func startCacheServer(addr string, addrs []string, g *go_cache.Group) {
	peers := go_cache.NewHTTPPool(addr) // 创建HTTPPool
	peers.Set(addrs...)                 // 添加节点信息
	g.RegisterPeers(peers)              // 注册到g中
	log.Println("gocache is running at", addr)
	log.Fatal(http.ListenAndServe(addr[7:], peers))
}

// 启动一个端口为9999的API服务
func startAPIServer(apiAddr string, g *go_cache.Group) {
	http.Handle("/api", http.HandlerFunc(
		func(writer http.ResponseWriter, request *http.Request) {
			key := request.URL.Query().Get("key")
			view, err := g.Get(key)
			if err != nil {
				http.Error(writer, err.Error(), http.StatusInternalServerError)
				return
			}
			writer.Header().Set("Content-Type", "application/octet-stream")
			writer.Write(view.ByteSlice())
		},
	))
	log.Println("end server is running at", apiAddr)
	log.Fatal(http.ListenAndServe(apiAddr[7:], nil))
}

func createGroup() *go_cache.Group {
	return go_cache.NewGroup("scores", 2<<10, go_cache.GetterFunc(
		func(key string) ([]byte, error) {
			log.Println("searching key...", key)
			if v, ok := db[key]; ok {
				return []byte(v), nil
			}
			return nil, fmt.Errorf("%s not exist", key)
		},
	))
}

func main() {

	var port int
	var api bool
	flag.IntVar(&port, "port", 8001, "gocache server port")
	flag.BoolVar(&api, "api", false, "Start a api server?")

	apiAddr := "http://localhost:9999"
	addrMap := map[int]string{
		8001: "http://localhost:8001",
		8002: "http://localhost:8002",
		8003: "http://localhost:8003",
	}

	var addrs []string
	for _, v := range addrMap {
		addrs = append(addrs, v)
	}

	g := createGroup()
	if api {
		go startAPIServer(apiAddr, g)
	}
	startCacheServer(addrMap[port], addrs, g)
}
