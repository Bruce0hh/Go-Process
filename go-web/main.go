package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/hello", helloHandler)
	log.Fatal(http.ListenAndServe(":9999", nil))
}

func helloHandler(w http.ResponseWriter, request *http.Request) {
	for k, v := range request.Header {
		_, _ = fmt.Fprintf(w, "Header[%v] = %v\n", k, v)
	}
}

func indexHandler(w http.ResponseWriter, request *http.Request) {
	_, _ = fmt.Fprintf(w, "URL.Path = %v\n", request.URL.Path)
}
