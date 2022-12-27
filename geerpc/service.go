package geerpc

import (
	"geerpc/log"
	"go/ast"
	"reflect"
	"sync/atomic"
)

type methodType struct {
	method    reflect.Method
	ArgType   reflect.Type
	ReplyType reflect.Type
	numCalls  uint64
}

func (m *methodType) NumCalls() uint64 {
	return atomic.LoadUint64(&m.numCalls)
}

func (m *methodType) newArgVal() reflect.Value {
	//这里的ArgType可能是一个指针类型
	//如果是指针类型，我们不能单纯的返回该类型的value(空指针)
	if m.ArgType.Kind() == reflect.Ptr {
		return reflect.New(m.ArgType.Elem())
	} else {
		return reflect.New(m.ArgType).Elem()
	}
}

func (m *methodType) newReplyVal() reflect.Value {
	//这里的返回值一定是一个指针类型，但是我们需要注意，如果是map或者slice类型，需要完成初始化
	replyVal := reflect.New(m.ReplyType.Elem())
	switch m.ReplyType.Elem().Kind() {
	case reflect.Slice:
		replyVal.Set(reflect.MakeSlice(m.ReplyType.Elem(), 0, 0))
	case reflect.Map:
		replyVal.Set(reflect.MakeMap(m.ReplyType.Elem()))
	}
	return replyVal
}

type service struct {
	name   string
	typ    reflect.Type
	rcvr   reflect.Value
	method map[string]*methodType
}

func newService(rcvr interface{}) *service {
	//rcvr可能是一个类型的指针
	s := new(service)
	s.name = reflect.Indirect(reflect.ValueOf(rcvr)).Type().Name()
	s.rcvr = reflect.ValueOf(rcvr)
	s.typ = reflect.TypeOf(rcvr)
	if !ast.IsExported(s.name) {
		log.Errorf("rpc server: %s is a invalid service name\n", s.name)
		return nil
	}
	s.registerMethod()
	return s
}

func (s *service) registerMethod() {
	s.method = make(map[string]*methodType)
	for i := 0; i < s.rcvr.NumMethod(); i++ {
		method := s.typ.Method(i)
		mType := method.Type
		if mType.NumIn() != 3 || mType.NumOut() != 1 {
			continue
		}
		if mType.Out(0) != reflect.TypeOf((*error)(nil)).Elem() {
			//返回类型不为error
			continue
		}
		argType, replyType := mType.In(1), mType.In(2)
		if !isExportOrBuiltin(argType) || !isExportOrBuiltin(replyType) {
			continue
		}
		if replyType.Kind() != reflect.Pointer {
			continue
		}
		s.method[method.Name] = &methodType{
			method:    method,
			ArgType:   argType,
			ReplyType: replyType,
		}
		log.Infof("rpc server: register %s.%s", s.name, method.Name)
	}
}

func isExportOrBuiltin(typ reflect.Type) bool {
	return ast.IsExported(typ.Name()) || typ.PkgPath() == ""
}

func (s *service) call(m *methodType, argv, reply reflect.Value) error {
	atomic.AddUint64(&m.numCalls, 1)
	f := m.method.Func
	rntValues := f.Call([]reflect.Value{s.rcvr, argv, reply})
	if errInter := rntValues[0].Interface(); errInter != nil {
		return errInter.(error)
	}
	return nil
}
