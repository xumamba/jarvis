package iface

/**
 * @Time       : 2020/12/27
 * @Author     : xumamba
 * @Description: 提供连接管理功能方法
 */

import (
	"net"
)

// IConnection 连接接口
type IConnection interface {
	// 启动连接
	Start()
	// 关闭连接
	Stop()
	// 获取原生连接
	GetTCPConn() *net.TCPConn
	// 获取当前连接ID
	GetConnID() uint32
	// 获取当前连接对端地址
	GetRemoteAddr() net.Addr
}
