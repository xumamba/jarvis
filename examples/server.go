package main

/**
 * @DateTime   : 2020/12/28
 * @Author     : xumamba
 * @Description:
 **/

import (
	"log"

	"jarvis/server"
	"jarvis/server/iface"
)

type PingRouter struct {
	server.BaseRouter
}

func (pr *PingRouter) PreHandle(request iface.IRequest) {
	if err := request.GetConn().SendMsg(1, []byte("[PingRouter] PreHandle...")); err != nil {
		log.Fatal("[PingRouter PreHandle] error", err)
	}
}

func (pr *PingRouter) Handle(request iface.IRequest) {
	log.Println("[PingRouter Handle] receive request: data=", string(request.GetData()))
}

func (pr *PingRouter) PostHandle(request iface.IRequest) {
	if err := request.GetConn().SendBuffMsg(1, []byte("[PingRouter] PostHandle...")); err != nil{
		log.Println("post handle send msg to client failed: ", err)
	}
}

type HelloRouter struct {
	server.BaseRouter
}

func (hr *HelloRouter) Handle(request iface.IRequest) {
	if err := request.GetConn().SendMsg(1, []byte("[HelloRouter]: Handle...")); err != nil{
		log.Println("post handle send msg to client failed: ", err)
	}
}

func main() {
	s := server.NewServer()
	s.AddRouter(0, &PingRouter{})
	s.AddRouter(1, &HelloRouter{})
	s.Serve()
}
