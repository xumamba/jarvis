package server

/**
 * @Time       : 2020/12/27
 * @Author     : xumamba
 * @Description: connection.go
 */

import "C"
import (
	"net"

	"jarvis/server/iface"
	"jarvis/utils/log"
)

type Connection struct {
	// 当前连接的唯一标识
	ConnID uint32
	// 当前连接的socket TCP套接字
	Conn *net.TCPConn
	// 当前连接的关闭状态
	isClosed bool
	// 告知该连接退出/停止的信道
	ExitChan chan bool
	// 连接处理函数
	Handlers HandlersChain
}

func (c *Connection) Start() {
	go c.StartReader()

	for {
		select {
		case <-c.ExitChan:
			return
		}
	}
}

func (c *Connection) Stop() {
	log.Logger.Info(c.GetRemoteAddr().String() + " connection exit.")
	if c.isClosed == true {
		return
	}
	c.isClosed = true

	// todo 执行用户注册的关闭连接后的回调业务方法

	// 关闭socket连接
	c.Conn.Close()
	// 通知缓冲队列读数据业务，该连接已关闭
	c.ExitChan <- true
	// 关闭连接信道
	close(c.ExitChan)
}

func (c *Connection) GetTCPConn() *net.TCPConn {
	return c.Conn
}

func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

func (c *Connection) GetRemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

// StartReader 处理读取连接数据
func (c *Connection) StartReader() {
	log.Logger.Info("Reader Goroutine is running...")
	defer c.Stop()

	for {
		buf := make([]byte, 512)

		count, err := c.Conn.Read(buf)
		if err != nil {
			log.Logger.Error("receive buf error: " + err.Error())
			c.ExitChan <- true
			continue
		}

		if _, err := c.Conn.Write(buf[:count]); err != nil {
			log.Logger.Error("write buf error: " + err.Error())
			c.ExitChan <- true
			return
		}

	}

}

// NewConn 创建连接
func NewConn(conn *net.TCPConn, connID uint32, handlers HandlersChain) iface.IConnection {
	return &Connection{
		ConnID:   connID,
		Conn:     conn,
		isClosed: false,
		ExitChan: make(chan bool, 1),
		Handlers: handlers,
	}
}
