package iface

/**
* @Time       : 2020/12/28
* @Author     : xumamba
* @Description: imsghandler 消息处理模块，根据不同的消息ID处理不同的路由模块
 */

type IMsgHandler interface {
	// 非阻塞方式处理请求
	Do(request IRequest)
	// 添加消息处理路由
	AddHandler(msgID uint32, router IRouter)
	// 发送消息至任务队列
	SendMsgToTaskQueue(request IRequest)
	// 启动worker池
	StartWorkerPool()
}
