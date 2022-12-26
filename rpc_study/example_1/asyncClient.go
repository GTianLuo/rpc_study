package example_1

import (
	"context"
	"fmt"
	"github.com/smallnest/rpcx/client"
	"github.com/smallnest/rpcx/log"
	"time"
)

func asyncClient() {
	//定义点对点的服务发现，客户端直连服务器
	d, err := client.NewPeer2PeerDiscovery("tcp@"+"localhost"+":8972", "")
	if err != nil {
		log.Error(err, "Failed to discover server")
	}

	args := &Args{
		A: 10,
		B: 10,
	}
	reply := &Reply{}
	xClient := client.NewXClient("Arith", client.Failtry, client.RandomSelect, d, client.DefaultOption)
	call, err := xClient.Go(context.Background(), "Mul", args, reply, nil)

	for {
		select {
		case <-call.Done:
			fmt.Printf("%v * %v = %v", args.A, args.B, reply.C)
			return
		default:
			fmt.Println("服务器计算未完成！")
			time.Sleep(time.Second)
		}
	}
}
