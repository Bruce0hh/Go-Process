package go_cache

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

type PeerGetter interface {
	// Get 根据group和key查找缓存值
	Get(group string, key string) ([]byte, error)
}

type PeerPicker interface {
	// PickPeer 用于传入的key选择响应的PeerGetter
	PickPeer(key string) (peer PeerGetter, ok bool)
}

// HTTP客户端类，实现PeerGetter接口
type httpGetter struct {
	baseURL string
}

func (h *httpGetter) Get(group string, key string) ([]byte, error) {
	u := fmt.Sprintf("%v%v/%v", h.baseURL, url.QueryEscape(group), url.QueryEscape(key))
	res, err := http.Get(u)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server return: %v", res.Status)
	}

	bytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response body: %v", err)
	}
	return bytes, nil
}

var _ PeerGetter = (*httpGetter)(nil)
