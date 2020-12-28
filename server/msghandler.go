package server

/**
 * @Time       : 2020/12/28
 * @Author     : xumamba
 * @Description: 路由管理模块
 */

import (
	"fmt"
	"strconv"

	"jarvis/server/iface"
	"jarvis/utils/log"
)

type MsgHandler struct {
	APIs map[uint32]iface.IRouter // 不同的消息ID，调用不同的路由处理
}

func (m *MsgHandler) Do(request iface.IRequest) {
	router, ok := m.APIs[request.GetMsgID()]
	if !ok {
		log.Logger.Error(fmt.Sprintf("api not found, msgID = %d", request.GetMsgID()))
		return
	}
	router.PreHandle(request)
	router.Handle(request)
	router.PostHandle(request)
}

func (m *MsgHandler) AddHandler(msgID uint32, router iface.IRouter) {
	if _, ok := m.APIs[msgID]; ok {
		log.Logger.Error(fmt.Sprintf("api repeated register, msgID = %d", msgID))
		panic(fmt.Sprintf("api repeated register, msgID = %d", msgID))
	}
	m.APIs[msgID] = router
	log.Logger.Info("api register success, msgID = " + strconv.Itoa(int(msgID)))
}

func NewMsgHandler() *MsgHandler {
	return &MsgHandler{APIs: make(map[uint32]iface.IRouter)}
}
