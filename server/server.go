package server

/**
* @DateTime   : 2020/9/18 16:14
* @Author     : xumamba
* @Description: 服务器实现
**/

import (
	"fmt"
	"io"
	"net"
	"reflect"
	"sync"
	"time"
	"unicode"
	"unicode/utf8"

	"jarvis/conf"
	"jarvis/server/iface"
	"jarvis/utils/log"
)

// Server 服务器
type Server struct {
	ServerName string // 服务器名称
	Network    string // tcp or tcp4 or tcp6
	IP         string // 服务器绑定的IP地址
	Port       int    // 服务器绑定的端口

	serviceMap sync.Map      // 已注册的服务
	Handlers   HandlersChain // 服务器中间件

	MsgHandler iface.IMsgHandler  // 服务路由方法
	ConnMgr    iface.IConnManager // 连接管理
}

// Start 服务器启动
func (s *Server) Start() {
	log.Logger.Info(fmt.Sprintf("Server Lintener at IP: %s, Port: %d is starting", s.IP, s.Port))

	// 开启服务器监听
	go func() {
		// 启动业务处理worker池
		s.MsgHandler.StartWorkerPool()

		tcpAddr, err := net.ResolveTCPAddr(s.Network, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			log.Logger.Error("resolve tcp address error: " + err.Error())
			return
		}
		listener, err := net.ListenTCP(s.Network, tcpAddr)
		if err != nil {
			log.Logger.Error(fmt.Sprintf("listen: %s, error: %s", s.Network, err))
			return
		}
		log.Logger.Info(fmt.Sprintf("start JARVIS server %s success, now listenning...", s.ServerName))

		// todo 连接ID生成方法
		var cid uint32
		cid = 0

		for {
			// 阻塞等待客户端建立连接
			conn, err := listener.AcceptTCP()
			if err != nil {
				log.Logger.Error("Accept error: " + err.Error())
				continue
			}
			// 判断最大连接数
			if s.ConnMgr.Len() >= conf.GlobalConfObj.MaxConnNum{
				conn.Close()
				continue
			}
			// 创建连接实体，处理连接绑定的业务方法
			dealConn := NewConn(s, conn, cid, s.Handlers, s.MsgHandler)
			cid++
			go dealConn.Start()
		}
	}()
}

// Stop 关闭服务器
func (s *Server) Stop() {
	log.Logger.Info("stop server, name: " + s.ServerName)
	// 关闭并清理当前建立的连接
	s.ConnMgr.ClearAll()
}

func (s *Server) Serve() {
	s.Start()

	// 阻塞主Goroutine退出
	for {
		time.Sleep(10 * time.Second)
	}
}

// GetConnMgr
func (s *Server) GetConnMgr() iface.IConnManager {
	return s.ConnMgr
}

// AddRouter 向服务器添加路由
func (s *Server) AddRouter(msgID uint32, router iface.IRouter) {
	s.MsgHandler.AddHandler(msgID, router)
}

// NewServer 服务器初始化
func NewServer() iface.IServer {
	return &Server{
		ServerName: conf.GlobalConfObj.Name,
		Network:    "tcp4",
		IP:         conf.GlobalConfObj.IP,
		Port:       conf.GlobalConfObj.Port,
		serviceMap: sync.Map{},
		Handlers:   make(HandlersChain, 0),
		MsgHandler: NewMsgHandler(),
		ConnMgr:    NewConnManager(),
	}
}

// Use 为服务添加中间件
func (s *Server) Use(middleware ...HandlerFunc) *Server {
	s.Handlers = append(s.Handlers, middleware...)
	return s
}

// GetContext 获取请求上下文
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

// HandleContext 将服务器绑定的中间件，添加置请求上下文并调用。
func (s *Server) HandleContext(c *Context) {
	c.handlers = s.Handlers
	c.Next()
}

// MustHandle 请求的入口处理函数
func (s *Server) MustHandle(name string, w io.Writer, r io.Reader) {
	ctx, err := s.GetContext(name, w, r)
	if err != nil {
		panic(err)
	}
	s.HandleContext(ctx)
}

// ForEachService 遍历服务器已注册的服务
func (s *Server) ForEachService(fun func(name string, svc *Service)) {
	s.serviceMap.Range(func(key, value interface{}) bool {
		svcName := key.(string)
		s := value.(*Service)
		fun(svcName, s)
		return true
	})
}

// RegisterService 向服务器中注册服务
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
