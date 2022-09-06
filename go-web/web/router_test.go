package web

import (
	"reflect"
	"testing"
)

func Test_parsePattern(t *testing.T) {

	tests := []struct {
		name string
		args string
		want []string
	}{
		{"test1", "/p/:name", []string{"p", ":name"}},
		{"test1", "/p/*", []string{"p", "*"}},
		{"test1", "/p/*name/*", []string{"p", "*name"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parsePattern(tt.args); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parsePattern() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_router_addRouter(t *testing.T) {
	type fields struct {
		roots    map[string]*node
		handlers map[string]HandlerFunc
	}
	type args struct {
		method  string
		pattern string
		handler HandlerFunc
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{"test1", fields{newRouter().roots, newRouter().handlers},
			args{"GET", "/", nil}},
		{"test2", fields{newRouter().roots, newRouter().handlers},
			args{"GET", "/hello/:name", nil}},
		{"test3", fields{newRouter().roots, newRouter().handlers},
			args{"GET", "/hi/:name", nil}},
		{"test4", fields{newRouter().roots, newRouter().handlers},
			args{"GET", "/assets/*filepath", nil}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &router{
				roots:    tt.fields.roots,
				handlers: tt.fields.handlers,
			}
			r.addRouter(tt.args.method, tt.args.pattern, tt.args.handler)
		})
	}
}

func Test_router_getRoute(t *testing.T) {
	r := newRouter()
	r.addRouter("GET", "hello/:name", nil)
	type fields struct {
		roots    map[string]*node
		handlers map[string]HandlerFunc
	}
	type args struct {
		method string
		path   string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
		want1  string
	}{
		{"test", fields{r.roots, r.handlers}, args{"GET", "/hello/goweb"},
			"hello/:name", "goweb"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &router{
				roots:    tt.fields.roots,
				handlers: tt.fields.handlers,
			}
			got, got1 := r.getRoute(tt.args.method, tt.args.path)
			if !reflect.DeepEqual(got.pattern, tt.want) {
				t.Errorf("getRoute() got = %v, want %v", got.pattern, tt.want)
			}
			if !reflect.DeepEqual(got1["name"], tt.want1) {
				t.Errorf("getRoute() got1 = %v, want1 %v", got1["name"], tt.want1)
			}
		})
	}
}
