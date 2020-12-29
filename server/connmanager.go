package server

/**
 * @DateTime   : 2020/12/29
 * @Author     : xumamba
 * @Description: 连接管理模块
 **/

import (
	"errors"
	"fmt"
	"sync"

	"jarvis/server/iface"
	"jarvis/utils/log"
)

type ConnManager struct {
	locker sync.RWMutex

	connections map[uint32]iface.IConnection // 连接容器
}

func NewConnManager() iface.IConnManager {
	return &ConnManager{
		connections: make(map[uint32]iface.IConnection),
	}
}

// Add 添加连接
func (c *ConnManager) Add(conn iface.IConnection) {
	c.locker.Lock()
	defer c.locker.Unlock()

	c.connections[conn.GetConnID()] = conn
	log.Logger.Info(fmt.Sprintf("connection add to ConnManager successful: connID = %d, conn nums = %d",
		conn.GetConnID(), c.Len()))
}

// Remove 移除连接
func (c *ConnManager) Remove(conn iface.IConnection) {
	c.locker.Lock()
	defer c.locker.Unlock()

	delete(c.connections, conn.GetConnID())
	log.Logger.Info(fmt.Sprintf("connection add to ConnManager successful: connID = %d, conn nums = %d",
		conn.GetConnID(), c.Len()))
}

// Get 获取指定连接
func (c *ConnManager) Get(connID uint32) (iface.IConnection, error) {
	c.locker.RLock()
	defer c.locker.RUnlock()

	if conn, ok := c.connections[connID]; ok {
		return conn, nil
	}
	return nil, errors.New("connection not found")
}

// Len 获取当前总连接数
func (c *ConnManager) Len() int {
	return len(c.connections)
}

// ClearAll 关闭并清除所有连接
func (c *ConnManager) ClearAll() {
	c.locker.Lock()
	defer c.locker.Unlock()

	for connID, conn := range c.connections {
		// 停止连接
		conn.Stop()
		// 清除连接
		delete(c.connections, connID)
	}
	log.Logger.Info(fmt.Sprintf("Clear All Connections successfully: conn num = %d", c.Len()))
}
