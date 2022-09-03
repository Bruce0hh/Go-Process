package xclient

import (
	"errors"
	"math"
	"math/rand"
	"sync"
	"time"
)

type SelectMode int

const (
	RandomSelect     SelectMode = iota // 随机算法
	RoundRobinSelect                   // Robin 轮询算法
)

type Discovery interface {
	Refresh() error                      // 从注册中心更新服务列表
	Update(servers []string) error       // 手动更新服务列表
	Get(mode SelectMode) (string, error) // 根据负载均衡策略，选择一个服务实例
	GetAll() ([]string, error)           // 返回所有服务实例
}

// MultiServersDiscovery 不使用注册中心维护的服务列表
type MultiServersDiscovery struct {
	r       *rand.Rand // r是一个产生随机数的实例
	mu      sync.Mutex
	servers []string
	index   int // index记录robin已经轮询的位置
}

// Refresh MultiServersDiscovery不使用注册中心，无需Refresh
func (m *MultiServersDiscovery) Refresh() error {
	return nil
}

// Update 动态更新
func (m *MultiServersDiscovery) Update(servers []string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.servers = servers
	return nil
}

func (m *MultiServersDiscovery) Get(mode SelectMode) (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	n := len(m.servers)
	if n == 0 {
		return "", errors.New("rpc discovery: no available servers")
	}
	switch mode {
	case RandomSelect:
		return m.servers[m.r.Intn(n)], nil
	case RoundRobinSelect:
		s := m.servers[m.index%n]
		m.index = (m.index + 1) % n
		return s, nil
	default:
		return "", errors.New("rpc discovery: not supported select mode")
	}
}

// GetAll 返回所有服务列表
func (m *MultiServersDiscovery) GetAll() ([]string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	servers := make([]string, len(m.servers), len(m.servers))
	copy(servers, m.servers)
	return servers, nil
}

// NewMultiServerDiscovery 创建实例
func NewMultiServerDiscovery(servers []string) *MultiServersDiscovery {
	d := &MultiServersDiscovery{
		r:       rand.New(rand.NewSource(time.Now().UnixNano())), // 初始化使用时间戳设定随机数种子，避免每次产生相同的随机数序列
		servers: servers,
	}
	d.index = d.r.Intn(math.MaxInt32 - 1) // 为了避免每次从0开始，初始化时随机设定一个值
	return d
}

var _ Discovery = (*MultiServersDiscovery)(nil)
