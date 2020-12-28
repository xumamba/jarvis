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
	log.Println("[PingRouter PreHandle]")
}

func (pr *PingRouter)Handle(request iface.IRequest)  {
	log.Println("[PingRouter Handle] receive request: data=", string(request.GetData()))
}

func main() {
	s := server.NewServer()
	s.AddRouter(&PingRouter{})
	s.Serve()
}
