package main

import (
	"fmt"
	"goweb/web"
	"net/http"
)

func main() {
	r := web.New()
	r.GET("/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, "URL.Path = %v\n", r.URL.Path)
	})

	r.GET("/hello", func(w http.ResponseWriter, r *http.Request) {
		for k, v := range r.Header {
			fmt.Fprintf(w, "Header[%v] = %v\n", k, v)
		}
	})
	_ = r.Run(":9999")
}
