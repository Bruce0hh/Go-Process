package web

import (
	"log"
	"net/http"
)

type router struct {
	handlers map[string]HandlerFunc
}

func newRouter() *router {
	return &router{handlers: make(map[string]HandlerFunc)}
}

func (r *router) addRouter(method, pattern string, handle HandlerFunc) {
	log.Printf("Route %v - %v\n", method, pattern)
	key := method + "-" + pattern
	r.handlers[key] = handle
}

func (r *router) handle(c *Context) {
	key := c.Method + "-" + c.Path
	if handler, ok := r.handlers[key]; ok {
		handler(c)
	} else {
		c.String(http.StatusNotFound, "404 NOT FOUND: %v\n", c.Path)
	}
}
