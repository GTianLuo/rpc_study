package example_1

import (
	"context"
	"fmt"
	"github.com/smallnest/rpcx/client"
	"github.com/smallnest/rpcx/log"
)

func myClient() {
	//定义点对点的服务发现，客户端直连服务器
	d, err := client.NewPeer2PeerDiscovery("tcp@"+"localhost"+":8972", "")
	if err != nil {
		log.Error(err, "Failed to discover server")
	}

	//servicePath string,  FailMode,  SelectMode,  ServiceDiscovery, Option
	//servicePath 服务名称
	//FailMode 告诉客户端如何处理调用失败：重试、快速返回，或者 尝试另一台服务器。
	//SelectMode 告诉客户端如何在有多台服务器提供了同一服务的情况下选择服务器。
	xClient := client.NewXClient("Arith", client.Failtry, client.RandomSelect, d, client.DefaultOption)
	defer xClient.Close()

	args := &Args{A: 10, B: 10}
	reply := &Reply{}
	err = xClient.Call(context.Background(), "Mul", args, reply)
	if err != nil {
		log.Error(err, "Failed to call")
	}

	fmt.Printf("%v * %v = %v", args.A, args.B, reply.C)

}
