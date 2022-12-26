package geerpc

import (
	"encoding/json"
	"errors"
	"fmt"
	"geerpc/codec"
	"geerpc/log"
	"io"
	"net"
	"sync"
)

// Call 一次RPC调用的全部信息
type Call struct {
	Seq           uint64
	ServiceMethod string
	Args          interface{}
	Reply         interface{}
	Error         error
	Done          chan *Call
}

// RPC调用结束
func (c *Call) done() *Call {
	c.Done <- c
	return c
}

// Client 处理RPC调用的客户端
type Client struct {
	cc       codec.Codec
	opt      *Option
	sending  sync.Mutex //保证消息的有序发送
	mu       sync.Mutex //保证对Client自身操作的安全性
	conn     io.ReadWriteCloser
	header   codec.Header
	seq      uint64
	pending  map[uint64]*Call
	closing  bool //主动关闭
	shutdown bool //因为异常关闭
}

var ErrShutdown = errors.New("connecting is shut down")

func (client *Client) Close() error {
	client.mu.Lock()
	defer client.mu.Unlock()
	if client.closing || client.shutdown {
		return ErrShutdown
	}
	client.closing = true
	_ = client.conn.Close()
	return nil
}

func (client *Client) IsAvailable() bool {
	client.mu.Lock()
	defer client.mu.Unlock()
	return !client.closing && !client.shutdown
}

func (client *Client) registerCall(c *Call) (uint64, error) {
	client.mu.Lock()
	defer client.mu.Unlock()
	if client.closing || client.shutdown {
		return 0, ErrShutdown
	}
	c.Seq = client.seq
	client.seq++
	client.pending[c.Seq] = c
	return c.Seq, nil
}

func (client *Client) removeCall(seq uint64) *Call {
	client.mu.Lock()
	defer client.mu.Unlock()
	c := client.pending[seq]
	if c != nil {
		delete(client.pending, seq)
	}
	return c
}

func (client *Client) terminateCalls(err error) {
	client.sending.Lock()
	defer client.sending.Unlock()
	client.mu.Lock()
	defer client.mu.Unlock()
	client.shutdown = true
	for _, call := range client.pending {
		call.Error = err
		call.done()
	}
}

func (client *Client) receive() {
	var err error
	for err == nil {
		var h = &codec.Header{}
		err = client.cc.ReadHeader(h)
		if err != nil {
			//与服务端的连接出现异常或关闭
			break
		}
		call := client.removeCall(h.Seq)
		switch {
		case call == nil:
			//TODO:没搞懂这一步有什么用
			err = client.cc.ReadBody(nil)
		case h.Error != "":
			//服务端在处理该请求时出现异常
			call.Error = errors.New(h.Error)
			call.done()
			err = client.cc.ReadBody(nil)
		default:
			err = client.cc.ReadBody(call.Reply)
			if err != nil {
				call.Error = errors.New("reading body " + err.Error())
			}
			call.done()
		}
	}
	client.terminateCalls(err)
}

func NewClient(conn io.ReadWriteCloser, opt *Option) (*Client, error) {
	f := codec.NewCodecFuncMap[opt.CodecType]
	if f == nil {
		err := fmt.Errorf("invaild codec type %v", opt.CodecType)
		log.Error("rpc client: option error: ", err)
		return nil, err
	}
	if err := json.NewEncoder(conn).Encode(opt); err != nil {
		log.Error("rpc client : encode error: ", err)
		_ = conn.Close()
		return nil, err
	}
	return newClient(f(conn), conn, opt), nil
}

func newClient(cc codec.Codec, conn io.ReadWriteCloser, opt *Option) *Client {
	client := &Client{
		cc:      cc,
		opt:     opt,
		seq:     1,
		conn:    conn,
		pending: make(map[uint64]*Call),
	}
	go client.receive()
	return client
}

func Dail(network, address string, opt *Option) (*Client, error) {
	conn, err := net.Dial(network, address)
	if err != nil {
		return nil, err
	}
	client, err := NewClient(conn, opt)
	if err != nil {
		_ = conn.Close()
	}
	return client, nil
}

func (client *Client) send(c *Call) error {
	client.sending.Lock()
	defer client.sending.Unlock()

	seq, err := client.registerCall(c)
	if err != nil {
		return err
	}
	//封装消息头
	var h codec.Header
	h.Seq = seq
	h.ServiceMethod = c.ServiceMethod

	//发送消息
	if err := client.cc.Write(&h, c.Args); err != nil {
		call := client.removeCall(c.Seq)
		if call != nil {
			call.Error = err
			call.done()
		}
		return err
	}
	return nil
}

func (client *Client) Call(serviceMethod string, args interface{}, reply interface{}) error {
	done, err := client.Go(serviceMethod, args, reply)
	if err != nil {
		return err
	}
	return (<-done).Error

}

func (client *Client) Go(serviceMethod string, args interface{}, reply interface{}) (chan *Call, error) {
	//封装Call
	call := &Call{
		ServiceMethod: serviceMethod,
		Args:          args,
		Reply:         reply,
		Done:          make(chan *Call),
	}

	if err := client.send(call); err != nil {
		return nil, err
	}

	return call.Done, nil
}
