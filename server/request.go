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
	data []byte
}

func (r *Request) GetConn() iface.IConnection {
	return r.conn
}

func (r *Request) GetData() []byte {
	return r.data
}


