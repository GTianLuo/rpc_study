package test

import (
	"fmt"
	"geerpc"
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
	client, _ := geerpc.Dail("tcp", "localhost:9997", geerpc.DefaultOption)
	defer client.Close()
	arg := &example.Args{
		A: 10,
		B: 10,
	}
	reply := &example.Reply{}

	client.Call("Arith.Add", arg, reply)
	fmt.Println("Arith.Add:", reply.C)
	client.Call("Arith.Mul", arg, reply)
	fmt.Println("Arith.Mul:", reply.C)
	client.Call("Arith.Sub", arg, reply)
	fmt.Println("Arith.Sub:", reply.C)
}

