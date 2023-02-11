package lru

import (
	"container/list"
	"reflect"
	"testing"
)

type String string

func (d String) Len() int {
	return len(d)
}

func TestCache_RemoveOldest(t *testing.T) {
	type fields struct {
		maxBytes  int64
		nbytes    int64
		ll        *list.List
		cache     map[string]*list.Element
		OnEvicted func(key string, value Value)
	}

	k1, k2, k3 := "key1", "key2", "k3"
	v1, v2, v3 := "value1", "value2", "v3"
	cap := len(k1 + k2 + v1 + v2)
	lru := New(int64(cap), nil)
	lru.Add(k1, String(v1))
	lru.Add(k2, String(v2))
	lru.Add(k3, String(v3))

	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
		{"rt1", fields(*lru)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Cache{
				maxBytes:  tt.fields.maxBytes,
				nbytes:    tt.fields.nbytes,
				ll:        tt.fields.ll,
				cache:     tt.fields.cache,
				OnEvicted: tt.fields.OnEvicted,
			}
			c.RemoveOldest()
		})
	}
}

func TestCache_Add(t *testing.T) {
	type fields struct {
		maxBytes  int64
		nbytes    int64
		ll        *list.List
		cache     map[string]*list.Element
		OnEvicted func(key string, value Value)
	}
	type args struct {
		key   string
		value Value
	}

	c := New(0, nil)

	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
		{"at", fields(*c), args{"key1", String("v1")}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Cache{
				maxBytes:  tt.fields.maxBytes,
				nbytes:    tt.fields.nbytes,
				ll:        tt.fields.ll,
				cache:     tt.fields.cache,
				OnEvicted: tt.fields.OnEvicted,
			}
			c.Add(tt.args.key, tt.args.value)
		})
	}
}

func TestCache_Get(t *testing.T) {
	type cache struct {
		maxBytes  int64
		nbytes    int64
		ll        *list.List
		cache     map[string]*list.Element
		OnEvicted func(key string, value Value)
	}
	type args struct {
		key string
	}
	c := New(0, nil)
	c.Add("key1", String("v"))
	c.Add("key1", String("s"))
	tests := []struct {
		name      string
		fields    cache
		args      args
		wantValue Value
		wantOk    bool
	}{
		// TODO: Add test cases.
		{"test2", cache(*c), args{key: "key1"}, String("s"), true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Cache{
				maxBytes:  tt.fields.maxBytes,
				nbytes:    tt.fields.nbytes,
				ll:        tt.fields.ll,
				cache:     tt.fields.cache,
				OnEvicted: tt.fields.OnEvicted,
			}
			gotValue, gotOk := c.Get(tt.args.key)
			if !reflect.DeepEqual(gotValue, tt.wantValue) {
				t.Errorf("Get() gotValue = %v, want %v", gotValue, tt.wantValue)
			}
			if gotOk != tt.wantOk {
				t.Errorf("Get() gotOk = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}
