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
	if _, err := request.GetConn().GetTCPConn().Write([]byte("PingRouter PreHandle...")); err != nil {
		log.Fatal("[PingRouter PreHandle] error", err)
	}
}

func (pr *PingRouter) Handle(request iface.IRequest) {
	log.Println("[PingRouter Handle] receive request: data=", string(request.GetData()))
}

func (pr *PingRouter) PostHandle(request iface.IRequest) {
	if err := request.GetConn().SendMsg(1, []byte("this is server response")); err != nil{
		log.Println("post handle send msg to client failed: ", err)
	}
}

func main() {
	s := server.NewServer()
	s.AddRouter(&PingRouter{})
	s.Serve()
}
