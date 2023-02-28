package main

import (
	"fmt"
	"hash/crc32"
	"sort"
	"strconv"
)

// Hash maps bytes to uint32
type Hash func(data []byte) uint32

// Map contains all hashed keys
type Map struct {
	hash     Hash           // hash函数
	replicas int            // 虚拟节点倍数
	keys     []int          // 哈希环
	hashMap  map[int]string // 节点与哈希值的映射表
}

// NewConsistentHash creates a Map instance
func NewConsistentHash(replicas int, fn Hash) *Map {
	m := &Map{
		replicas: replicas,
		hash:     fn,
		hashMap:  make(map[int]string),
	}
	if m.hash == nil {
		m.hash = crc32.ChecksumIEEE
	}
	return m
}

// Add adds some keys to the hash
func (m *Map) Add(keys ...string) {
	for _, key := range keys {
		for i := 0; i < m.replicas; i++ {
			hash := int(m.hash([]byte(strconv.Itoa(i) + key)))
			m.keys = append(m.keys, hash)
			m.hashMap[hash] = key
		}
	}
	sort.Ints(m.keys)
}

// Get gets the closest item in the hash to the provided key
func (m *Map) Get(key string) string {
	if len(m.keys) == 0 {
		return ""
	}

	hash := int(m.hash([]byte(key)))
	// Binary search for appropriate replica
	idx := sort.Search(len(m.keys), func(i int) bool { return m.keys[i] >= hash })
	// If we have gone through full circle
	if idx == len(m.keys) {
		idx = 0
	}
	return m.hashMap[m.keys[idx]]
}

// Remove removes some keys from the hash
func (m *Map) Remove(keys ...string) {
	for _, key := range keys {
		for i := 0; i < m.replicas; i++ {
			hash := int(m.hash([]byte(strconv.Itoa(i) + key)))
			j := sort.Search(len(m.keys), func(j int) bool { return m.keys[j] >= hash })
			if j < len(m.keys) && m.keys[j] == hash {
				m.keys = append(m.keys[:j], m.keys[j+1:]...)
				delete(m.hashMap, hash)
			}
		}
	}
}
func main() {
	// 新建一个一致性哈希对象
	ch := NewConsistentHash(3, nil)

	// 添加节点
	for i := 0; i < 3; i++ {
		ch.Add("node" + strconv.Itoa(i))
	}
	m := make(map[string]int, 3)
	// 将100个数据分配到节点上
	for i := 0; i < 100; i++ {
		node := ch.Get(strconv.Itoa(i))
		m[node]++
		fmt.Printf("key=%d, node=%s\n", i, node)
	}
	fmt.Printf("%v\n", m)
}
