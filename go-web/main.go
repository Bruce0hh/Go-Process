package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	log.Fatal(http.ListenAndServe(":9999", new(Engine)))
}

type Engine struct{}

func (e *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	switch req.URL.Path {
	case "/":
		_, _ = fmt.Fprintf(w, "URL.Path = %v\n", req.URL.Path)
	case "/hello":
		for k, v := range req.Header {
			_, _ = fmt.Fprintf(w, "Header[%v] = %v", k, v)
		}
	default:
		_, _ = fmt.Fprintf(w, "404 NOT FOUND: %v", req.URL)
	}
}
