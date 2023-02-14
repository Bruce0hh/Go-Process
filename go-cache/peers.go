package go_cache

import (
	"fmt"
	pb "gocache/gocachepb"
	"google.golang.org/protobuf/proto"
	"io/ioutil"
	"net/http"
	"net/url"
)

type PeerGetter interface {
	// Get 根据group和key查找缓存值
	Get(in *pb.Request, out *pb.Response) error
}

type PeerPicker interface {
	// PickPeer 用于传入的key选择响应的PeerGetter
	PickPeer(key string) (peer PeerGetter, ok bool)
}

// HTTP客户端类，实现PeerGetter接口
type httpGetter struct {
	baseURL string
}

func (h *httpGetter) Get(in *pb.Request, out *pb.Response) error {
	u := fmt.Sprintf(
		"%v%v/%v",
		h.baseURL,
		url.QueryEscape(in.GetGroup()),
		url.QueryEscape(in.GetKey()),
	)
	res, err := http.Get(u)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("server return: %v", res.Status)
	}

	bytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("reading response body: %v", err)
	}
	if err = proto.Unmarshal(bytes, out); err != nil {
		return fmt.Errorf("decoding response body: %v", err)
	}

	return nil
}

var _ PeerGetter = (*httpGetter)(nil)
