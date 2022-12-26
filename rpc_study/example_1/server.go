package example_1

import (
	"github.com/smallnest/rpcx/log"
	"github.com/smallnest/rpcx/server"
)

func myServer() {
	//注册服务
	s := server.NewServer()
	//RegisterName 进行注册时可以指定服务名
	//Register   默认会使用 rcvr 的类型名进行注册
	//这两个方法都有metadata参数，可以在注册中心添加一些元数据提供客户端或服务端使用
	err := s.RegisterName("Arith", new(Arith), "")
	if err != nil {
		log.Error(err, "Failed to register server")
	}
	//启动tcp服务监听请求
	err = s.Serve("tcp", ":8972")
	if err != nil {
		log.Error(err, "Failed to start serve")
	}
}
