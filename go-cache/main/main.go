package main

import (
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

func main() {
	go_cache.NewGroup("scores", 2<<10, go_cache.GetterFunc(
		func(key string) ([]byte, error) {
			log.Println("searching key...", key)
			if v, ok := db[key]; ok {
				return []byte(v), nil
			}
			return nil, fmt.Errorf("%s not exist", key)
		},
	))

	addr := "localhost:9999"
	peers := go_cache.NewHTTPPool(addr)
	log.Printf("gocache is running at %s", addr)
	log.Fatal(http.ListenAndServe(addr, peers))
}
