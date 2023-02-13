package consistenthash

import (
	"hash/crc32"
	"sort"
	"strconv"
)

type Hash func(data []byte) uint32

// Map 存储所有hash key
type Map struct {
	hash     Hash           // Hash函数
	replicas int            // 虚拟节点倍数
	keys     []int          // 哈希环
	hashMap  map[int]string // key是虚拟节点的hash值，value是真实节点的名称
}

// New 实例初始化，允许自定义虚拟节点倍数和Hash函数
func New(replicas int, fn Hash) *Map {
	m := &Map{
		hash:     fn,
		replicas: replicas,
		hashMap:  make(map[int]string),
	}
	if m.hash == nil {
		m.hash = crc32.ChecksumIEEE
	}
	return m
}

// Add 可传入0或多个真实节点的名称
// 对每一个真实节点key，创建m.replicas个虚拟节点，名称为strconv.Iota(i)+key
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

// Get 选择节点
func (m *Map) Get(key string) string {
	if len(m.keys) == 0 {
		return ""
	}
	hash := int(m.hash([]byte(key)))                   // 计算hash
	idx := sort.Search(len(m.keys), func(i int) bool { // 找到最近匹配的虚拟节点下标idx
		return m.keys[i] >= hash
	})
	return m.hashMap[m.keys[idx%len(m.keys)]] // m.keys是环状，所以要取余
}
