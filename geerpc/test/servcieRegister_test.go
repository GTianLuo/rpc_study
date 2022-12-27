package test

import (
	"context"
	"fmt"
	"geerpc"
	"geerpc/codec"
	"geerpc/example"
	"log"
	"net"
	"testing"
)

func TestServer2(t *testing.T) {
	err := geerpc.Register(&example.Arith{})
	if err != nil {
		log.Println(err)
		return
	}
	conn, err := net.Listen("tcp", "localhost:9997")
	if err != nil {
		fmt.Println(err)
	}
	geerpc.Accept(conn)
}

func TestClient4(t *testing.T) {
	client, _ := geerpc.Dail("tcp", "localhost:9997", &geerpc.Option{MagicNumber: geerpc.MagicNumber, CodecType: codec.GobType, HandleTimeOut: 1})
	defer client.Close()
	arg := &example.Args{
		A: 10,
		B: 10,
	}
	reply := &example.Reply{}

	err := client.Call(context.Background(), "Arith.Add", arg, reply)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(context.Background(), "Arith.Add:", reply.C)
	err = client.Call(context.Background(), "Arith.Mul", arg, reply)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(context.Background(), "Arith.Mul:", reply.C)
	err = client.Call(context.Background(), "Arith.Sub", arg, reply)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(context.Background(), "Arith.Sub:", reply.C)
}
