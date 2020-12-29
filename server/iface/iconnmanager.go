package iface

/**
* @DateTime   : 2020/12/29
* @Author     : xumamba
* @Description: 连接管理模块
**/

type IConnManager interface {
	// Add 添加连接
	Add(conn IConnection)
	// Remove 停止并清除指定连接
	Remove(conn IConnection)
	// Get 获取指定连接
	Get(connID uint32) (IConnection, error)
	// Len 获取当前连接数量
	Len() int
	// ClearAll 停止并清除所有连接
	ClearAll()
}
