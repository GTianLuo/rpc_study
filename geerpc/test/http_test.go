package test

import (
	"context"
	"fmt"
	"geerpc"
	"geerpc/example"
	"geerpc/log"
	"net"
	"net/http"
	"os"
	"testing"
)

func TestServer5(t *testing.T) {
	err := geerpc.Register(&example.Arith{})
	if err != nil {
		log.Error(err)
		return
	}
	l, err := net.Listen("tcp", ":9997")
	if err != nil {
		fmt.Println(err)
	}
	geerpc.HandleHTTP()
	err = http.Serve(l, nil)
}

func TestClient5(t *testing.T) {
	client, err := geerpc.DialHTTP("tcp", ":9997", geerpc.DefaultOption)
	if err != nil {
		log.Error(err)
	}
	for i := 0; i < 10; i++ {
		args := &example.Args{A: 10, B: i}
		reply := &example.Reply{}
		err := client.Call(context.Background(), "Arith.Add", args, reply)
		if err != nil {
			log.Error(err)
		}
		fmt.Println(reply)
	}
}

func TestXDail(t *testing.T) {
	ch := make(chan struct{})
	addr := "/tmp/geerpc.tcp"
	go func() {
		_ = os.Remove(addr)
		l, err := net.Listen("unix", addr)
		if err != nil {
			log.Error(err)
		}
		ch <- struct{}{}
		geerpc.Accept(l)
	}()
	<-ch
	_, err := geerpc.XDail("unix@"+addr, geerpc.DefaultOption)
	if err != nil {
		log.Error(err)
	}
}
