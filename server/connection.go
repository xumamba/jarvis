package server

/**
 * @Time       : 2020/12/27
 * @Author     : xumamba
 * @Description: connection.go
 */

import "C"
import (
	"errors"
	"io"
	"net"
	"strconv"
	"sync"

	"jarvis/conf"
	"jarvis/server/iface"
	"jarvis/utils/log"
)

type Connection struct {
	sync.RWMutex

	Server iface.IServer

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
	MsgHandler iface.IMsgHandler
	// 无缓冲 消息传递管道，用于读写分离
	msgChan chan []byte
	// 有缓冲 消息传递管道，用于读写分离
	msgBuffChan chan []byte
}

func (c *Connection) Start() {
	go c.StartReader()
	go c.StartWriter()

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
	// 从连接管理器中移除连接
	c.Server.GetConnMgr().Remove(c)
	// 关闭连接信道
	close(c.ExitChan)
	close(c.msgChan)
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
	if c.isClosed == true {
		return errors.New("Connection closed when send msg: connID= " + strconv.Itoa(int(c.GetConnID())))
	}
	msg := NewMessage(msgID, data)
	packageMsg, err := DPHelper.PackageMsg(msg)
	if err != nil {
		log.Logger.Error("PackageMsg error: " + err.Error())
		return err
	}
	c.msgChan <- packageMsg
	return nil
}

func (c *Connection) SendBuffMsg(msgID uint32, data []byte) error {
	if c.isClosed == true {
		return errors.New("Connection closed when send msg: connID= " + strconv.Itoa(int(c.GetConnID())))
	}
	packageMsg, err := DPHelper.PackageMsg(NewMessage(msgID, data))
	if err != nil {
		log.Logger.Error("PackageMsg error: " + err.Error())
		return err
	}
	// 回写消息给客户端
	c.msgBuffChan <- packageMsg
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
		if conf.GlobalConfObj.WorkerPoolSize > 0 {
			// 启动worker工作池机制
			c.MsgHandler.SendMsgToTaskQueue(request)
		} else {
			// 直接开协程处理请求
			go c.MsgHandler.Do(request)
		}
	}

}

// StartWriter 写数据给客户端
func (c *Connection) StartWriter() {
	log.Logger.Info("Writer Goroutine is running...")
	defer c.Stop()

	for {
		select {
		case data := <-c.msgChan:
			// 有数据要写给客户端
			if _, err := c.Conn.Write(data); err != nil {
				log.Logger.Error("Send data error: " + err.Error())
				return
			}
		case data, ok := <-c.msgBuffChan:
			if ok {
				if _, err := c.Conn.Write(data); err != nil {
					log.Logger.Error("Send data error: " + err.Error())
					return
				}
			} else {
				log.Logger.Info("msgBuffChan is closed")
				break
			}
		case <-c.ExitChan:
			// conn 已关闭
			return
		}
	}
}

// NewConn 创建连接
func NewConn(ser iface.IServer, conn *net.TCPConn, connID uint32, handlers HandlersChain, msgHandler iface.IMsgHandler) iface.IConnection {
	c := &Connection{
		Server:      ser,
		ConnID:      connID,
		Conn:        conn,
		isClosed:    false,
		ExitChan:    make(chan bool, 1),
		Handlers:    handlers,
		MsgHandler:  msgHandler,
		msgChan:     make(chan []byte),
		msgBuffChan: make(chan []byte, conf.GlobalConfObj.MaxMsgChanLen),
	}
	// 将新建连接交由连接管理模块
	c.Server.GetConnMgr().Add(c)
	return c
}
