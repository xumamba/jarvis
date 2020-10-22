/**
* @DateTime   : 2020/9/18 16:14
* @Author     : xumamba
* @Description:
**/
package server

import (
	"fmt"
	"io"
	"reflect"
	"sync"
	"unicode"
	"unicode/utf8"
)

type Server struct {
	serviceMap sync.Map
	Handlers   HandlersChain
}

func NewServer() *Server {
	return &Server{
		Handlers: make(HandlersChain, 0),
	}
}

func (s *Server) Use(middleware ...HandlerFunc) *Server {
	s.Handlers = append(s.Handlers, middleware...)
	return s
}

func (s *Server) GetContext(svcName string, w io.Writer, r io.Reader) (*Context, error) {
	svc, ok := s.serviceMap.Load(svcName)
	if !ok {
		return nil, fmt.Errorf("service not found:%s", svcName)
	}
	return &Context{
		server:    s,
		w:         w,
		r:         r,
		svc:       svc.(*Service),
		container: make(map[string]interface{}),
	}, nil
}

func (s *Server) ForEachService(fun func(name string, svc *Service)) {
	s.serviceMap.Range(func(key, value interface{}) bool {
		svcName := key.(string)
		s := value.(*Service)
		fun(svcName, s)
		return true
	})
}

func (s *Server) HandleContext(c *Context) {
	c.handlers = s.Handlers
	c.Next()
}

func (s *Server) MustHandle(name string, w io.Writer, r io.Reader) {
	ctx, err := s.GetContext(name, w, r)
	if err != nil {
		panic(err)
	}
	s.HandleContext(ctx)
}

func (s *Server) RegisterService(svcName string, fun interface{}) error {
	return s.register(svcName, fun)
}

func (s *Server) RegisterGRPCService(svc interface{}) {
	svcType := reflect.TypeOf(svc)
	svcValue := reflect.ValueOf(svc)
	for i := 0; i < svcValue.NumMethod(); i++ {
		method := svcValue.Method(i)
		methodType := svcType.Method(i)
		err := s.registerGRPC(svcValue, methodType.Name, method.Interface())
		if err != nil {
			panic(err)
		}
	}
}

func (s *Server) register(svcName string, fun interface{}) error {
	funType := reflect.TypeOf(fun)
	if funType.Kind() != reflect.Func {
		return fmt.Errorf("unknown function type:%s", funType)
	}
	if funType.NumIn() != 3 || funType.NumOut() != 1 {
		return fmt.Errorf("register function parameters must be:func Func(ctx context.Context, argv *Arg, reply *Reply) error")
	}

	reqType := funType.In(1)
	if !isAvailableType(reqType) {
		return fmt.Errorf("register function arg type must be ptr and export:%s", reqType)
	}
	respType := funType.In(2)
	if !isAvailableType(respType) {
		return fmt.Errorf("register function reply type must be ptr and export:%s", respType)
	}
	errType := funType.Out(0)
	if errType != reflect.TypeOf((*error)(nil)).Elem() {
		return fmt.Errorf("register function return type must be error:%s", errType)
	}

	funVal := reflect.ValueOf(fun)
	svc := &Service{
		name:     svcName,
		fun:      funVal,
		reqType:  reqType,
		respType: respType,
		IsGRPC:   false,
	}

	if _, loaded := s.serviceMap.LoadOrStore(svcName, svc); loaded {
		return fmt.Errorf("service already registered:%s", svcName)
	}
	return nil

}

func (s *Server) registerGRPC(svcValue reflect.Value, svcName string, fun interface{}) error {
	funType := reflect.TypeOf(fun)
	if funType.Kind() != reflect.Func {
		return fmt.Errorf("unknown function type:%s", funType)
	}
	if funType.NumIn() != 2 || funType.NumOut() != 2 {
		return fmt.Errorf("register grpc function parameters must be:func Func(ctx context.Context, argv *Arg)(replyv *Reply, err error)")
	}
	reqType := funType.In(1)
	if !isAvailableType(reqType) {
		return fmt.Errorf("register function arg type must be ptr and export:%s", reqType)
	}
	respType := funType.Out(0)
	if !isAvailableType(respType) {
		return fmt.Errorf("register function reply type must be ptr and export:%s", respType)
	}
	if errType := funType.Out(1); errType != reflect.TypeOf((*error)(nil)).Elem() {
		return fmt.Errorf("register function return type must be error")
	}

	funValue := reflect.ValueOf(fun)
	svc := &Service{
		name:     svcName,
		fun:      funValue,
		reqType:  reqType,
		respType: respType,
		IsGRPC:   true,
	}

	if _, loaded := s.serviceMap.LoadOrStore(svcName, svc); loaded {
		return fmt.Errorf("service already registered:%s", svcName)
	}
	return nil
}

func isAvailableType(t reflect.Type) bool {
	for t.Kind() != reflect.Ptr {
		return false
	}
	t = t.Elem()
	return isExported(t.Name()) || t.PkgPath() == ""
}

func isExported(name string) bool {
	r, _ := utf8.DecodeRuneInString(name)
	return unicode.IsUpper(r)
}
