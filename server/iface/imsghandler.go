package iface

/**
* @Time       : 2020/12/28
* @Author     : xumamba
* @Description: imsghandler 消息处理模块，根据不同的消息ID处理不同的路由模块
 */

type IMsgHandler interface {
	Do(request IRequest)
	AddHandler(msgID uint32, router IRouter)
}