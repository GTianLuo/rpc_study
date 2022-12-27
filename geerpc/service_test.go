package geerpc

import (
	"fmt"
	"geerpc/example"
	"reflect"
	"testing"
)

func TestService(t *testing.T) {
	service := newService(&example.Arith{})
	args := example.Args{
		A: 10,
		B: 10,
	}
	reply := &example.Reply{}
	_ = service.call(service.method["Add"], reflect.ValueOf(args), reflect.ValueOf(reply))
	fmt.Println(reply.C)
	_ = service.call(service.method["Mul"], reflect.ValueOf(args), reflect.ValueOf(reply))
	fmt.Println(reply.C)
	_ = service.call(service.method["Sub"], reflect.ValueOf(args), reflect.ValueOf(reply))
	fmt.Println(reply.C)
	//println(reflect.ValueOf(&example.Arith{}).Method(1).Type())
}
