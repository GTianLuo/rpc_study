package geerpc

import (
	"encoding/json"
	"errors"
	"fmt"
	"geerpc/codec"
	"geerpc/log"
	"io"
	"net"
	"reflect"
	"strings"
	"sync"
	"time"
)

const MagicNumber = 0x3bef5c //16 767 534

type Option struct {
	MagicNumber   int
	CodecType     codec.Type
	ConnTimeOut   time.Duration //0 表示不超时
	HandleTimeOut time.Duration //0 表示不超时
}

func NewOption(codecType codec.Type) *Option {
	return &Option{
		MagicNumber: MagicNumber,
		CodecType:   codecType,
		ConnTimeOut: 10,
	}
}

type Server struct {
	serviceMap sync.Map
}

func NewServer() *Server {
	return &Server{}
}

var DefaultServer = NewServer()

var DefaultOption = &Option{
	MagicNumber: MagicNumber,
	CodecType:   codec.GobType,
}

func (server *Server) Register(rcvr interface{}) error {
	service := newService(rcvr)
	if _, loaded := server.serviceMap.LoadOrStore(service.name, service); loaded {
		return errors.New("rpc: service already defined:" + service.name)
	}
	return nil
}

func Register(rcvr interface{}) error {
	return DefaultServer.Register(rcvr)
}

func (server *Server) findService(serverMethod string) (*service, *methodType, error) {
	dot := strings.LastIndex(serverMethod, ".")
	if dot < 0 {
		err := errors.New("rpc server :service/method request ill-form: " + serverMethod)
		return nil, nil, err
	}
	svci, ok := server.serviceMap.Load(serverMethod[:dot])
	if !ok {
		err := errors.New("rpc server : can't find service: " + serverMethod[:dot])
		return nil, nil, err
	}
	svc := svci.(*service)
	mTyp := svc.method[serverMethod[dot+1:]]
	if mTyp == nil {
		err := errors.New("rpc server : can't find method: " + serverMethod[dot+1:])
		return svc, nil, err
	}
	return svc, mTyp, nil
}

func (server *Server) Accept(lis net.Listener) {
	for {
		conn, err := lis.Accept()
		if err != nil {
			log.Error("rpc server : accept error:", err)
			return
		}
		log.Infof("server Connect client : %s\n", conn.RemoteAddr())
		go server.ServeConn(conn)
	}
}

func Accept(lis net.Listener) {
	DefaultServer.Accept(lis)
}

type request struct {
	h            *codec.Header
	argv, replyv reflect.Value
	service      *service
	mTyp         *methodType
}

func (server *Server) ServeConn(conn io.ReadWriteCloser) {
	defer func() {
		_ = conn.Close()
	}()
	var opt Option
	if err := json.NewDecoder(conn).Decode(&opt); err != nil {
		log.Error("rpc server :option error", err)
		return
	}
	if opt.MagicNumber != MagicNumber {
		log.Errorf("rpc server: invalid magic number %d \n", opt.MagicNumber)
		return
	}
	f := codec.NewCodecFuncMap[opt.CodecType]
	if f == nil {
		log.Error("rpc server: invalid codec type")
		return
	}
	server.ServeCodec(f(conn), opt.HandleTimeOut)
}

var invalidRequest = struct{}{}

func (server *Server) ServeCodec(cc codec.Codec, handleTimeOut time.Duration) {
	sending := new(sync.Mutex)
	wg := new(sync.WaitGroup)
	for {
		request, err := server.readRequest(cc)
		if err != nil {
			if request == nil {
				//连接关闭
				break
			}
			request.h.Error = err.Error()
			server.sendResponse(cc, *request.h, invalidRequest, sending)
			continue
		}
		wg.Add(1)
		go server.handleRequest(cc, request, sending, wg, handleTimeOut)
	}
	wg.Wait()
	_ = cc.Close()
}

func (server *Server) readRequest(cc codec.Codec) (*request, error) {
	h, err := server.readRequestHeader(cc)
	if err != nil {
		return nil, err
	}
	req := &request{h: h}
	if req.service, req.mTyp, err = server.findService(h.ServiceMethod); err != nil {
		return req, err
	}
	req.argv = req.mTyp.newArgVal()
	req.replyv = req.mTyp.newReplyVal()
	argvi := req.argv.Interface()
	if req.argv.Kind() != reflect.Pointer {
		argvi = req.argv.Addr().Interface()
	}
	err = cc.ReadBody(argvi)
	if err != nil {
		log.Error("rpc server : read body error:", err)
	}
	return req, nil
}

func (server *Server) readRequestHeader(cc codec.Codec) (*codec.Header, error) {
	var h codec.Header
	if err := cc.ReadHeader(&h); err != nil {
		if err != io.EOF && err != io.ErrUnexpectedEOF {
			log.Error("rpc server : read header error :", err)
		}
		return nil, err
	}
	return &h, nil
}

func (server *Server) sendResponse(cc codec.Codec, h codec.Header, body interface{}, sending *sync.Mutex) {
	sending.Lock()
	defer sending.Unlock()
	if err := cc.Write(h, body); err != nil {
		log.Error("rpc err: write response error", err)
	}
}

func (server *Server) handleRequest(cc codec.Codec, req *request, sending *sync.Mutex, wg *sync.WaitGroup, handleTimeOut time.Duration) {
	defer wg.Done()
	called := make(chan bool)
	go func() {
		log.Info(req.h, req.argv)
		err := req.service.call(req.mTyp, req.argv, req.replyv)
		select {
		case called <- true:
		default: //超时
			return
		}
		if err != nil {
			req.h.Error = err.Error()
			server.sendResponse(cc, *req.h, invalidRequest, sending)
			return
		}
		server.sendResponse(cc, *req.h, req.replyv.Interface(), sending)
	}()
	if handleTimeOut == 0 {
		<-called
		return
	}
	select {
	case <-time.After(handleTimeOut):
		req.h.Error = fmt.Sprintf("rpc server: request handle timeout: expect within %s", handleTimeOut)
		server.sendResponse(cc, *req.h, invalidRequest, sending)
	case <-called:
	}

}
