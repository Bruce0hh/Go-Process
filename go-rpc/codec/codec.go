package codec

import "io"

// 客户端发送请求：服务名-方法名-参数
// 服务端响应：返回值-错误
// 参数和返回值可抽象为 Header+Body

type Header struct {
	ServiceMethod string // 服务名和方法名，与Go中结构体和方法名映射
	Seq           uint64 // 请求序号，可以认为是某个请求的ID
	Error         string // 错误信息
}

// Codec 编解码接口
type Codec interface {
	io.Closer
	ReadHeader(*Header) error
	ReadBody(interface{}) error
	Write(*Header, interface{}) error
}

// NewCodecFunc Codec构造函数
type NewCodecFunc func(io.ReadWriteCloser) Codec

type Type string

const (
	GobType  Type = "application/gob"
	JsonType Type = "application/json"
)

var NewCodecFuncMap map[Type]NewCodecFunc

func init() {
	NewCodecFuncMap = make(map[Type]NewCodecFunc)
	NewCodecFuncMap[GobType] = NewGobCodec
}
