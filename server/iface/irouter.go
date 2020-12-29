package iface

/**
* @DateTime   : 2020/12/28
* @Author     : xumamba
* @Description: 基础路由模块
**/

type IRouter interface {
	// 业务函数前置处理函数
	PreHandle(request IRequest)
	// 业务处理函数
	Handle(request IRequest)
	// 业务处理后置处理函数
	PostHandle(request IRequest)  
}
