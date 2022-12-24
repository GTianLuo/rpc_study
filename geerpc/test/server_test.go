package test

import (
	"encoding/json"
	"fmt"
	"geerpc"
	"geerpc/codec"
	"net"
	"testing"
)

func TestServer(t *testing.T) {
	conn, err := net.Listen("tcp", "localhost:9999")
	if err != nil {
		fmt.Println(err)
	}
	geerpc.Accept(conn)
}

func TestClient(t *testing.T) {
	conn, err := net.Dial("tcp", "localhost:9999")
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
