package test

import (
	"context"
	"fmt"
	"geerpc"
	"geerpc/example"
	"geerpc/xclient"
	"log"
	"net"
	"testing"
)

func TestBalanceServer1(t *testing.T) {
	err := geerpc.Register(&example.Arith{})
	if err != nil {
		log.Println(err)
		return
	}
	conn, err := net.Listen("tcp", ":9996")
	if err != nil {
		fmt.Println(err)
	}
	geerpc.Accept(conn)
}

func TestBalanceServer2(t *testing.T) {
	err := geerpc.Register(&example.Arith{})
	if err != nil {
		log.Println(err)
		return
	}
	conn, err := net.Listen("tcp", ":9997")
	if err != nil {
		fmt.Println(err)
	}
	geerpc.Accept(conn)
}

func TestBalanceServer3(t *testing.T) {
	err := geerpc.Register(&example.Arith{})
	if err != nil {
		log.Println(err)
		return
	}
	conn, err := net.Listen("tcp", ":9998")
	if err != nil {
		fmt.Println(err)
	}
	geerpc.Accept(conn)
}

func TestBalanceServer4(t *testing.T) {
	err := geerpc.Register(&example.Arith{})
	if err != nil {
		log.Println(err)
		return
	}
	conn, err := net.Listen("tcp", ":9999")
	if err != nil {
		fmt.Println(err)
	}
	geerpc.Accept(conn)
}

func TestBalanceClient(t *testing.T) {
	//"tcp@localhost:9998", "tcp@localhost:9997", "tcp@localhost:9996"
	d := xclient.NewMultiServerDiscovery([]string{"tcp@localhost:9999", "tcp@localhost:9998", "tcp@localhost:9997", "tcp@localhost:9996"})
	xc := xclient.NewXClient(d, xclient.RoundRobinSelect, geerpc.DefaultOption)
	args := &example.Args{9, 8}
	reply := &example.Reply{}
	err := xc.Call(context.Background(), "Arith.Add", args, reply)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(*reply)
	xc.Call(context.Background(), "Arith.Mul", args, reply)
	fmt.Println(*reply)
	//xc.Call(context.Background(), "Method.Div", args, reply)
	xc.Call(context.Background(), "Arith.Sub", args, reply)
	fmt.Println(*reply)
	xc.Call(context.Background(), "Arith.Sub", args, reply)
	fmt.Println(*reply)
}
