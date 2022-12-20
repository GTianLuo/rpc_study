package example_1

import (
	"github.com/smallnest/rpcx/log"
	"github.com/smallnest/rpcx/server"
)

func myServer() {
	//注册服务
	s := server.NewServer()
	err := s.RegisterName("Arith", new(Arith), "")
	if err != nil {
		log.Error(err, "Failed to register server")
	}

	err = s.Serve("tcp", ":8972")
	if err != nil {
		log.Error(err, "Failed to start serve")
	}
}
