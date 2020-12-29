package server

/**
 * @Time       : 2020/12/28
 * @Author     : xumamba
 * @Description: 路由管理模块
 */

import (
	"fmt"
	"strconv"

	"jarvis/conf"
	"jarvis/server/iface"
	"jarvis/utils/log"
)

type MsgHandler struct {
	APIs           map[uint32]iface.IRouter // 不同的消息ID，调用不同的路由处理
	WorkerPoolSize uint32                   // 业务处理Worker池大小
	TaskQueue      []chan iface.IRequest    // 消息队列，与Worker绑定，用于向worker分发任务
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

// StartOneWorker 启动一个worker处理请求
func (m *MsgHandler) StartOneWorker(workerID int, taskQueue chan iface.IRequest)  {
	log.Logger.Info("Worker ID = " + strconv.Itoa(workerID) + " is started.")
	for {
		select {
		case request := <- taskQueue:
			m.Do(request)
		}
	}
}

// StartWorkerPool 启动worker工作池
func (m *MsgHandler) StartWorkerPool() {
	for i := 0; i < int(m.WorkerPoolSize); i++{
		m.TaskQueue[i] = make(chan iface.IRequest, conf.GlobalConfObj.MaxTaskQueueLen)
		go m.StartOneWorker(i, m.TaskQueue[i])
	}
}

// SendMsgToTaskQueue 发送请求消息至任务队列
func (m *MsgHandler) SendMsgToTaskQueue(request iface.IRequest) {
	// 各worker之间负载均衡
	// 轮询
	workerID := request.GetConn().GetConnID() % m.WorkerPoolSize
	m.TaskQueue[workerID] <- request
}


func NewMsgHandler() *MsgHandler {
	return &MsgHandler{
		APIs: make(map[uint32]iface.IRouter),
		WorkerPoolSize: conf.GlobalConfObj.WorkerPoolSize,
		TaskQueue: make([]chan iface.IRequest, conf.GlobalConfObj.WorkerPoolSize),
	}
}
