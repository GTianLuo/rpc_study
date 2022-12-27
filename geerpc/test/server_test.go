package test

import (
	"encoding/json"
	"fmt"
	"geerpc"
	"geerpc/codec"
	"net"
	"reflect"
	"testing"
	"time"
)

func TestServer(t *testing.T) {
	conn, err := net.Listen("tcp", "localhost:9997")
	if err != nil {
		fmt.Println(err)
	}
	geerpc.Accept(conn)
}

func TestClient(t *testing.T) {
	conn, err := net.Dial("tcp", "localhost:9997")
	defer conn.Close()
	if err != nil {
		fmt.Println(err)
	}

	//time.Sleep(time.Second)
	e := json.NewEncoder(conn)
	_ = e.Encode(geerpc.DefaultOption)

	cc := codec.NewGobCodec(conn)
	for i := 1; i <= 5; i++ {
		h := &codec.Header{
			ServiceMethod: "Test.Add",
			Seq:           uint64(i),
		}
		_ = cc.Write(h, fmt.Sprintf("geerpc request %d", h.Seq))
		_ = cc.ReadHeader(h)
		var reply string
		_ = cc.ReadBody(&reply)
		fmt.Println(reply)
	}
}

func TestClient2(t *testing.T) {
	conn, err := net.Dial("tcp", "localhost:9999")
	defer conn.Close()
	if err != nil {
		fmt.Println(err)
	}

	time.Sleep(time.Second)
	e := json.NewEncoder(conn)
	_ = e.Encode(geerpc.DefaultOption)

	cc := codec.NewGobCodec(conn)
	for i := 1; i <= 1000; i++ {
		h := &codec.Header{
			ServiceMethod: "Test.Add",
			Seq:           uint64(i),
		}
		_ = cc.Write(h, fmt.Sprintf("geerpc request %d", h.Seq))
		_ = cc.ReadHeader(h)
		var reply string
		_ = cc.ReadBody(&reply)
		fmt.Println(reply)
	}
}

func TestClient3(t *testing.T) {
	client, _ := geerpc.Dail("tcp", "localhost:9997", geerpc.DefaultOption)
	defer client.Close()
	for i := 0; i < 10; i++ {
		arg := "rpc request"
		var reply string
		fmt.Println("")
		err := client.Call("Method.Add", arg, &reply)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(reply)
	}
}

func TestReflect(t *testing.T) {
	s := make([]int, 0)
	sTyp := reflect.TypeOf(&s)
	fmt.Println(sTyp)
	fmt.Println(sTyp.Elem())

}
