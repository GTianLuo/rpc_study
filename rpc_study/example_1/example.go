package example_1

import (
	"context"
	"time"
)

type Args struct {
	A int
	B int
}

type Reply struct {
	C int
}

type Arith int

// Mul 提供服务的方法，要满足以下三点：
// 1. 方法必须是可导出的
// 2. 三个参数，其中第一个上下文参数不能是指针，第三个参数必须是指针类型
// 3. 有一个error类型的返回值
func (a *Arith) Mul(c context.Context, args *Args, reply *Reply) error {
	time.Sleep(time.Second * 5)
	reply.C = args.A * args.B
	return nil
}
