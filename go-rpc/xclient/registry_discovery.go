package xclient

import (
	"log"
	"net/http"
	"strings"
	"time"
)

type RegistryDiscovery struct {
	*MultiServersDiscovery
	registry   string        // 注册中心地址
	timeout    time.Duration // 服务列表的过期时间
	lastUpdate time.Time     // 最后从注册中心更新服务列表的时间，默认10s过期；过期后，需要从注册中心更新新的列表
}

const defaultUpdateTimeout = time.Second * 10

func NewRegistryDiscovery(registryAddr string, timeout time.Duration) *RegistryDiscovery {
	if timeout == 0 {
		timeout = defaultUpdateTimeout
	}
	return &RegistryDiscovery{
		MultiServersDiscovery: NewMultiServerDiscovery(make([]string, 0)),
		registry:              registryAddr,
		timeout:               timeout,
	}
}

func (rd *RegistryDiscovery) Update(servers []string) error {
	rd.mu.Lock()
	defer rd.mu.Unlock()
	rd.servers = servers
	rd.lastUpdate = time.Now()
	return nil
}

func (rd *RegistryDiscovery) Refresh() error {
	rd.mu.Lock()
	defer rd.mu.Unlock()

	if rd.lastUpdate.Add(rd.timeout).After(time.Now()) {
		return nil
	}
	log.Printf("rpc registry: refresh servers from registry: %v\n", rd.registry)
	resp, err := http.Get(rd.registry)
	if err != nil {
		log.Printf("rpc registry refresh err: %v\n", err)
		return err
	}
	servers := strings.Split(resp.Header.Get("X-Gorpc-Servers"), ",")
	rd.servers = make([]string, 0, len(servers))
	for _, server := range servers {
		if strings.TrimSpace(server) != "" {
			rd.servers = append(rd.servers, strings.TrimSpace(server))
		}
	}
	rd.lastUpdate = time.Now()

	return nil
}

func (rd *RegistryDiscovery) Get(mode SelectMode) (string, error) {
	if err := rd.Refresh(); err != nil {
		return "", err
	}
	return rd.MultiServersDiscovery.Get(mode)
}

func (rd *RegistryDiscovery) GetAll() ([]string, error) {
	if err := rd.Refresh(); err != nil {
		return nil, err
	}
	return rd.MultiServersDiscovery.GetAll()
}
