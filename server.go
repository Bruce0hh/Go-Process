package gorpc

import (
	"encoding/json"
	"errors"
	"gorpc/codec"
	"io"
	"log"
	"net"
	"reflect"
	"strings"
	"sync"
)

// 服务端报文设计
// 在单次连接中，报文可能是下面这样的：
// | Option | Header1 | Body1 | Header2 | Body2 | ...

const MagicNumber = 0x3bef5c

type Option struct {
	MagicNumber int
	CodecType   codec.Type
}

var DefaultOption = &Option{
	MagicNumber: MagicNumber,
	CodecType:   codec.GobType,
}

// Server 服务端实现
type Server struct {
	serviceMap sync.Map
}

func (s *Server) Register(rec interface{}) error {
	svc := newService(rec)
	if _, dup := s.serviceMap.LoadOrStore(svc.name, svc); dup {
		return errors.New("rpc: server already registered: " + svc.name)
	}
	return nil
}

func Register(rec interface{}) error {
	return DefaultServer.Register(rec)
}

func (s *Server) findService(serviceMethod string) (svc *service, mType *methodType, err error) {

	dot := strings.LastIndex(serviceMethod, ".")
	if dot < 0 {
		err = errors.New("rpc server: service/method: " + serviceMethod)
		return
	}
	serviceName, methodName := serviceMethod[:dot], serviceMethod[dot+1:]
	svci, ok := s.serviceMap.Load(serviceName)
	if !ok {
		err = errors.New("rpc server: can't not find service: " + serviceName)
		return
	}
	svc = svci.(*service)
	mType = svc.method[methodName]
	if mType == nil {
		err = errors.New("rpc server: can't find method: " + methodName)
	}

	return
}

func NewServer() *Server {
	return &Server{}
}

// DefaultServer 提供一个默认实例
var DefaultServer = NewServer()

// Accept 循环等待socket连接，开启子协程进行处理
func (s *Server) Accept(lis net.Listener) {
	for {
		conn, err := lis.Accept()
		if err != nil {
			log.Printf("rpc server: accept error: %+v", err)
			return
		}
		go s.ServeConn(conn)
	}
}

func Accept(lis net.Listener) {
	DefaultServer.Accept(lis)
}

// ServeConn 处理报文Option部分
func (s *Server) ServeConn(conn io.ReadWriteCloser) {
	defer func() {
		_ = conn.Close()
	}()
	var opt Option
	if err := json.NewDecoder(conn).Decode(&opt); err != nil {
		log.Printf("rpc server: option error: %+v", err)
		return
	}
	if opt.MagicNumber != MagicNumber {
		log.Printf("rpc server: invalid magic number: %+v", opt.MagicNumber)
		return
	}
	codeType := codec.NewCodecFuncMap[opt.CodecType]
	if codeType == nil {
		log.Printf("rpc server: invalid codec type: %+v", opt.CodecType)
		return
	}
	s.serveCodec(codeType(conn))
}

var invalidRequest = struct{}{}

func (s *Server) serveCodec(cc codec.Codec) {
	sending := new(sync.Mutex)
	wg := new(sync.WaitGroup)

	for {
		req, err := s.readRequest(cc)
		if err != nil {
			if req == nil {
				break
			}
			req.header.Error = err.Error()
			s.sendResponse(cc, req.header, invalidRequest, sending)
			continue
		}
		wg.Add(1)
		go s.handleRequest(cc, req, sending, wg)
	}
	wg.Wait()
	_ = cc.Close()
}

type request struct {
	header     *codec.Header
	arg, reply reflect.Value
	mType      *methodType
	svc        *service
}

// 读取请求中的Header
func (s *Server) readRequestHeader(cc codec.Codec) (*codec.Header, error) {
	var h codec.Header
	if err := cc.ReadHeader(&h); err != nil {
		if err != io.EOF && err != io.ErrUnexpectedEOF {
			log.Printf("rpc server: read header error: %+v", err)
		}
		return nil, err
	}
	return &h, nil
}

// 读取请求
func (s *Server) readRequest(cc codec.Codec) (*request, error) {
	h, err := s.readRequestHeader(cc)
	if err != nil {
		return nil, err
	}
	req := &request{
		header: h,
	}

	req.svc, req.mType, err = s.findService(h.ServiceMethod)
	if err != nil {
		return req, nil
	}
	req.arg = req.mType.newArgv()
	req.reply = req.mType.newReply()

	argv := req.arg.Interface()
	if req.arg.Type().Kind() != reflect.Ptr {
		argv = req.arg.Addr().Interface()
	}
	if err = cc.ReadBody(argv); err != nil {
		log.Printf("rpc server: read arg err: %+v", err)
		return req, nil
	}
	return req, nil
}

// 回复请求
func (s *Server) sendResponse(cc codec.Codec, header *codec.Header, body interface{}, sending *sync.Mutex) {
	sending.Lock()
	defer sending.Unlock()
	if err := cc.Write(header, body); err != nil {
		log.Printf("rpc server: write response error: %+v", err)
	}
}

// 处理请求
func (s *Server) handleRequest(cc codec.Codec, req *request, sending *sync.Mutex, wg *sync.WaitGroup) {
	defer wg.Done()
	if err := req.svc.call(req.mType, req.arg, req.reply); err != nil {
		req.header.Error = err.Error()
		s.sendResponse(cc, req.header, invalidRequest, sending)
		return
	}
	s.sendResponse(cc, req.header, req.reply.Interface(), sending)
}
