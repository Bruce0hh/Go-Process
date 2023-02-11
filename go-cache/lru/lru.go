package lru

import "container/list"

// Cache LRU数据结构：字典+双向链表
type Cache struct {
	maxBytes  int64 // 允许使用的最大内存
	nbytes    int64 // 当前已使用的内存
	ll        *list.List
	cache     map[string]*list.Element
	OnEvicted func(key string, value Value) // 某条记录被移除时的回调函数
}

type entry struct {
	key   string
	value Value
}

// Value 可以接收任意类型
type Value interface {
	Len() int
}

// New 初始化
func New(maxBytes int64, onEvicted func(string2 string, value Value)) *Cache {
	return &Cache{
		maxBytes:  maxBytes,
		ll:        list.New(),
		cache:     make(map[string]*list.Element),
		OnEvicted: onEvicted,
	}
}

// Get 查找功能
// 1.从字典中找到对应的双向链表的节点
// 2.将该节点移动到队尾
func (c *Cache) Get(key string) (value Value, ok bool) {
	if ele, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ele) // 约定Front为队尾
		kv := ele.Value.(*entry)
		return kv.value, true
	}
	return
}

// RemoveOldest 移除最近最少访问的节点（队首）
func (c *Cache) RemoveOldest() {
	ele := c.ll.Back()
	if ele != nil {
		c.ll.Remove(ele) // 取到队首节点，从链表中删除
		kv := ele.Value.(*entry)
		delete(c.cache, kv.key)                                // 从字典中删除该节点的映射关系
		c.nbytes -= int64(len(kv.key)) + int64(kv.value.Len()) // 更新内存
		if c.OnEvicted != nil {
			c.OnEvicted(kv.key, kv.value) // 调用回调函数
		}
	}
}

// Add 新增/修改
func (c *Cache) Add(key string, value Value) {
	if ele, ok := c.cache[key]; ok { // 如果键存在，则更新对应节点的值，并将该节点移到队尾
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		c.nbytes += int64(value.Len()) - int64(kv.value.Len())
		kv.value = value
	} else { // 如果键不存在，则添加新节点，并在字典中添加key和value
		ele := c.ll.PushFront(&entry{key, value})
		c.cache[key] = ele
		c.nbytes += int64(len(key)) + int64(value.Len())
	}
	// 更新c.nbytes，如果超过了设定的最大值，则移除最少访问的节点
	for c.maxBytes != 0 && c.maxBytes < c.nbytes {
		c.RemoveOldest()
	}
}

// Len 获取缓存数
func (c *Cache) Len() int {
	return c.ll.Len()
}
