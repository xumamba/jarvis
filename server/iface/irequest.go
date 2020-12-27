package iface

/**
* @Time       : 2020/12/27
* @Author     : xumamba
* @Description: 客户端连接信息和请求的数据信息 的封装
 */

type IRequest interface {
	// GetConn 获取连接信息
	GetConn() IConnection
	// GetData 获取请求消息
	GetData() []byte
}

