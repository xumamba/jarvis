package log

/**
 * @Time       : 2020/12/27
 * @Author     : xumamba
 * @Description: 日志记录模块，抽离出来为了方便使用者自定义该模块。
 */

import (
	sysLog "log"
)

type ILog interface {
	Info(info interface{})
	Debug(debugInfo interface{})
	Error(errInfo interface{})
}

type logger struct{}

func (l *logger) Info(info interface{}) {
	sysLog.Println("[JARVIS Info]: ", info)
}

func (l *logger) Debug(debugInfo interface{}) {
	sysLog.Println("[JARVIS Debug]: ", debugInfo)
}

func (l *logger) Error(errInfo interface{}) {
	sysLog.Println("[JARVIS Error]: ", errInfo)
}

func NewLogger() ILog{
	return &logger{}
}

var Logger = NewLogger()