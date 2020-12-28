package server

/**
 * @Time       : 2020/12/27
 * @Author     : xumamba
 * @Description: request.go
 */

import (
	"jarvis/server/iface"
)


type Request struct {
	conn iface.IConnection
	msg iface.IMessage
}

func (r *Request) GetConn() iface.IConnection {
	return r.conn
}

func (r *Request) GetData() []byte {
	return r.msg.GetRealData()
}

func (r *Request) GetMsgID() uint32 {
	return r.msg.GetMsgID()
}



