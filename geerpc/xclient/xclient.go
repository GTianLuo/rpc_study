package xclient

import (
	"context"
	"geerpc"
	"reflect"
	"sync"
)

type XClient struct {
	d       Discovery
	mode    SelectMode
	opt     *geerpc.Option
	mu      sync.Mutex
	clients map[string]*geerpc.Client
}

func NewXClient(d Discovery, mode SelectMode, opt *geerpc.Option) *XClient {
	return &XClient{
		d:       d,
		mode:    mode,
		opt:     opt,
		clients: make(map[string]*geerpc.Client),
	}
}

func (xc *XClient) Close() error {
	xc.mu.Lock()
	defer xc.mu.Unlock()
	for key, client := range xc.clients {
		_ = client.Close()
		delete(xc.clients, key)
	}
	return nil
}

func (xc *XClient) dail(rpcAddr string) (*geerpc.Client, error) {
	c, ok := xc.clients[rpcAddr]
	if ok && !c.IsAvailable() {
		delete(xc.clients, rpcAddr)
		c = nil
	}
	if c == nil {
		var err error
		if c, err = geerpc.XDail(rpcAddr, xc.opt); err != nil {
			return nil, err
		}
		xc.clients[rpcAddr] = c
	}

	return c, nil
}

func (xc *XClient) call(rpcAddr string, context context.Context, serviceMethod string, args interface{}, reply interface{}) error {
	client, err := xc.dail(rpcAddr)
	if err != nil {
		return err
	}
	return client.Call(context, serviceMethod, args, reply)
}

func (xc *XClient) Call(context context.Context, serviceMethod string, args interface{}, reply interface{}) error {
	rpcAddr, err := xc.d.Get(xc.mode)
	if err != nil {
		return err
	}
	return xc.call(rpcAddr, context, serviceMethod, args, reply)
}

func (xc *XClient) Broadcast(ctx context.Context, serviceMethod string, args interface{}, reply interface{}) error {
	rpcAddrs, err := xc.d.GetAll()
	if err != nil {
		return err
	}
	var wg sync.WaitGroup
	var mu sync.Mutex
	var e error
	ctx, cancel := context.WithCancel(ctx)
	replyDone := reply == nil
	for _, rpcAddr := range rpcAddrs {
		wg.Add(1)
		go func(rpcAddr string) {
			defer wg.Done()
			var cloneReply interface{}
			if cloneReply != nil {
				cloneReply = reflect.New(reflect.TypeOf(reply).Elem()).Interface()
			}
			err := xc.Call(ctx, serviceMethod, args, cloneReply)
			mu.Lock()
			if err != nil && e == nil {
				cancel()
				e = err
			} else if err == nil && !replyDone {
				reflect.ValueOf(reply).Set(reflect.ValueOf(cloneReply).Elem())
				replyDone = true
				cancel()
			}
			mu.Unlock()
		}(rpcAddr)
	}
	wg.Wait()
	return e
}
