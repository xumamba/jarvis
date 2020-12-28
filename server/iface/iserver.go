package iface

/**
 * @Time       : 2020/12/27 11:54 上午
 * @Author     : xumamba
 * @Description: 定义服务器功能方法
 */


// IServer 服务器接口
type IServer interface {
	// 启动服务器
	Start()
	// 停止服务器
	Stop()
	// 开启业务服务方法
	Serve()
	// 添加业务处理路由
	AddRouter(router IRouter)
}
