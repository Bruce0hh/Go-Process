package registry

import (
	"log"
	"net/http"
	"sort"
	"strings"
	"sync"
	"time"
)

type RegisterCenter struct {
	timeout time.Duration
	mu      sync.Mutex
	servers map[string]*ServerItem
}

type ServerItem struct {
	Addr  string
	start time.Time
}

const (
	defaultPath    = "/_gorpc_/register"
	defaultTimeout = time.Minute * 5
)

// NewRegister 初始化实例
func NewRegister(timeout time.Duration) *RegisterCenter {
	return &RegisterCenter{servers: make(map[string]*ServerItem), timeout: timeout}
}

var DefaultRegister = NewRegister(defaultTimeout)

// 添加服务实例，如果服务已经存在，则更新start
func (r *RegisterCenter) putServer(addr string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	s := r.servers[addr]
	if s == nil {
		r.servers[addr] = &ServerItem{Addr: addr, start: time.Now()}
	} else {
		s.start = time.Now()
	}
}

// 返回可用的服务列表，如果存在超时的服务，则删除
func (r *RegisterCenter) aliveServers() []string {
	r.mu.Lock()
	defer r.mu.Unlock()

	var alive []string
	for addr, serverItem := range r.servers {
		if r.timeout == 0 || serverItem.start.Add(r.timeout).After(time.Now()) {
			alive = append(alive, addr)
		} else {
			delete(r.servers, addr)
		}
	}
	sort.Strings(alive)
	return alive
}

func (r *RegisterCenter) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "GET": // 返回所有可用的服务列表
		w.Header().Set("X-Gorpc-Servers", strings.Join(r.aliveServers(), ","))
	case "POST": // 添加服务实例或者发送心跳
		addr := req.Header.Get("X-Gorpc-Servers")
		if addr == "" {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		r.putServer(addr)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (r *RegisterCenter) HandleHTTP(registryPath string) {
	http.Handle(registryPath, r)
	log.Printf("rpc registry path: %v\n", registryPath)
}

func HandleHTTP() {
	DefaultRegister.HandleHTTP(defaultPath)
}

func HeartBeat(registry, addr string, duration time.Duration) {
	if duration == 0 {
		duration = defaultTimeout - time.Duration(1)*time.Minute
	}

	var err error
	err = sendHeartbeat(registry, addr)
	go func() {
		t := time.NewTicker(duration)
		for err == nil {
			<-t.C
			err = sendHeartbeat(registry, addr)
		}
	}()
}

func sendHeartbeat(registry, addr string) error {
	log.Printf("%v send heart beat to registry %v", addr, registry)
	httpClient := &http.Client{}
	req, _ := http.NewRequest("POST", registry, nil)
	req.Header.Set("X-Gorpc-Server", addr)
	if _, err := httpClient.Do(req); err != nil {
		return err
	}
	return nil
}
