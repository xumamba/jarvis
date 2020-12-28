package iface

/**
* @DateTime   : 2020/12/28
* @Author     : xumamba
 * @Description: 直接面向TCP连接中的数据流。LTV(长度、类型、数据)结构，用于处理TCP粘包问题。
 **/

// IMessage 对一个请求消息的封装
type IMessage interface {
	GetMsgLen() uint32  // 获取消息数据部分长度
	GetMsgID() uint32  // 获取消息唯一标识
	GetRealData() []byte  // 获取消息原始数据

	SetMsgLen(len uint32)  // 设置消息数据部分长度
	SetMsgID(id uint32)  // 设置消息唯一标识
	SetRealData(data []byte)  // 设置消息原始数据
}

// IDataPackage 对数据进行封包和拆包
type IDataPackage interface {
	GetHeadLen() uint32  // 获取数据包头部长度
	PackageMsg(msg IMessage)([]byte, error)  // 数据封包
	UnPackageMsg(data []byte)(IMessage, error)  // 数据拆包
}