package go_cache

import (
	"fmt"
	"log"
	"sync"
)

type Getter interface {
	Get(key string) ([]byte, error)
}

// GetterFunc
// 方便使用者在调用时既能够传入函数作为参数，也能够传入实现该接口的结构体作为参数
type GetterFunc func(key string) ([]byte, error)

func (f GetterFunc) Get(key string) ([]byte, error) {
	return f(key)
}

// Group 缓存的命名空间,类似redis的db
type Group struct {
	name      string
	getter    Getter // 缓存未命中时获取源数据的回调
	mainCache cache  // 一开始实现的并发缓存
}

var (
	mu     sync.RWMutex              // 只读锁，不涉及写操作
	groups = make(map[string]*Group) // 存储group的全局变量
)

// NewGroup 创建实例
func NewGroup(name string, cacheBytes int64, getter Getter) *Group {
	if getter == nil {
		panic("nil Getter")
	}
	mu.Lock()
	defer mu.Unlock()

	g := &Group{
		name:      name,
		getter:    getter,
		mainCache: cache{cacheBytes: cacheBytes},
	}
	groups[name] = g
	return g
}

// GetGroup 通过name获取group
func GetGroup(name string) *Group {
	mu.RLock()
	defer mu.RUnlock()
	g := groups[name]
	return g
}

// Get 从mainCache中查找缓存，如果存在则返回
// 不存在则调用load，load调用getLocally
// getLocally则调用回调函数获取源数据，并通过populateCache添加到缓存mainCache中
func (g *Group) Get(key string) (ByteView, error) {
	if key == "" {
		return ByteView{}, fmt.Errorf("key is required")
	}
	if v, ok := g.mainCache.get(key); ok {
		log.Println("Cache hit")
		return v, nil
	}
	return g.load(key)
}

func (g *Group) load(key string) (value ByteView, err error) {
	return g.getLocally(key)
}

func (g *Group) getLocally(key string) (ByteView, error) {
	bytes, err := g.getter.Get(key)
	if err != nil {
		return ByteView{}, err
	}
	value := ByteView{b: cloneBytes(bytes)}
	g.populateCache(key, value)
	return value, nil
}

func (g *Group) populateCache(key string, value ByteView) {
	g.mainCache.add(key, value)
}
