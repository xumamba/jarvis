package server

/**
 * @DateTime   : 2020/12/28
 * @Author     : xumamba
 * @Description:
 **/

import (
	"jarvis/server/iface"
)

// BaseRouter 路由基类，所有路由实现可以继承该基类，目的是方便自定义路由有选择性的添加路由函数。
type BaseRouter struct{}

func (b *BaseRouter) PreHandle(request iface.IRequest) {}

func (b *BaseRouter) Handle(request iface.IRequest) {}

func (b *BaseRouter) PostHandle(request iface.IRequest) {}
