package example

import "time"

type Arith struct {
}
type Args struct {
	A int
	B int
}

type Reply struct {
	C int
}

func (a *Arith) Add(args Args, reply *Reply) error {
	time.Sleep(time.Second * 2)
	reply.C = args.A + args.B
	return nil
}

func (a *Arith) Sub(args Args, reply *Reply) error {
	reply.C = args.A * args.B
	return nil
}

func (a *Arith) Mul(args Args, reply *Reply) error {
	reply.C = args.A / args.B
	return nil
}

func (a *Arith) div(args Args, reply *Reply) error {
	reply.C = args.A / args.B
	return nil
}

func (a *Arith) Add2(args Args, reply *Reply) {
	reply.C = args.A + args.B
}

func (a *Arith) Ad3(args Args, reply *Reply, args2 Args) error {
	reply.C = args.A + args.B
	return nil
}
