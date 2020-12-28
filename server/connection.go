package server

/**
 * @Time       : 2020/12/27
 * @Author     : xumamba
 * @Description: connection.go
 */

import "C"
import (
	"errors"
	"fmt"
	"io"
	"net"
	"strconv"
	"sync"

	"jarvis/server/iface"
	"jarvis/utils/log"
)

type Connection struct {
	sync.RWMutex

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
	// 连接处理路由函数
	Router iface.IRouter
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

func (c *Connection) SendMsg(msgID uint32, data []byte) error {
	c.RLock()
	defer c.RUnlock()
	if c.isClosed == true{
		return errors.New("Connection closed when send msg: connID= " + strconv.Itoa(int(c.GetConnID())))
	}
	msg := NewMessage(msgID, data)
	packageMsg, err := DPHelper.PackageMsg(msg)
	if err != nil{
		log.Logger.Error("PackageMsg error: " + err.Error())
		return err
	}
	fmt.Println(packageMsg)
	if _, err := c.GetTCPConn().Write(packageMsg); err != nil{
		log.Logger.Error("write message to client error: " + err.Error())
		c.ExitChan <- true
		return err
	}
	return nil
}

// StartReader 处理读取连接数据
func (c *Connection) StartReader() {
	log.Logger.Info("Reader Goroutine is running...")
	defer c.Stop()

	for {
		msgHead := make([]byte, DPHelper.GetHeadLen())
		if _, err := io.ReadFull(c.GetTCPConn(), msgHead); err != nil {
			log.Logger.Error("receive buf error: " + err.Error())
			c.ExitChan <- true
			continue
		}
		msg, err := DPHelper.UnPackageMsg(msgHead)
		if err != nil {
			log.Logger.Error("unpack error: " + err.Error())
			c.ExitChan <- true
			continue
		}

		data := make([]byte, msg.GetMsgLen())
		if _, err := io.ReadFull(c.GetTCPConn(), data); err != nil {
			log.Logger.Error("read msg data error: " + err.Error())
			c.ExitChan <- true
			continue
		}

		msg.SetRealData(data)

		request := &Request{
			conn: c,
			msg:  msg,
		}

		go func(request iface.IRequest) {
			c.Router.PreHandle(request)
			c.Router.Handle(request)
			c.Router.PostHandle(request)
		}(request)

	}

}

// NewConn 创建连接
func NewConn(conn *net.TCPConn, connID uint32, handlers HandlersChain, router iface.IRouter) iface.IConnection {
	return &Connection{
		ConnID:   connID,
		Conn:     conn,
		isClosed: false,
		ExitChan: make(chan bool, 1),
		Handlers: handlers,
		Router:   router,
	}
}
