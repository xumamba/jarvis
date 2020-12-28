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
	// Start 启动连接
	Start()
	// Stop 关闭连接
	Stop()
	// GetTCPConn 获取原生连接
	GetTCPConn() *net.TCPConn
	// GetConnID 获取当前连接ID
	GetConnID() uint32
	// GetRemoteAddr 获取当前连接对端地址
	GetRemoteAddr() net.Addr
	// SendMsg 向客户端发送数据
	SendMsg(msgID uint32, data []byte) error
}
