package conf

/**
 * @DateTime   : 2020/12/28
 * @Author     : xumamba
 * @Description: 全局配置文件
 **/

import (
	"encoding/json"
	"io/ioutil"
)

// Config jarvis框架配置
type Config struct {
	Name    string // 服务器名称
	IP      string // 服务器IP
	Port    int    // 服务器监听端口
	Version string // 服务器版本

	MaxPacketSize   uint32 // 最大数据包大小
	MaxConnNum      int    // 最大连接数
	WorkerPoolSize  uint32 // 业务处理工作池的worker数量
	MaxTaskQueueLen uint32 // 与worker绑定的任务队列最大任务存储数量
	MaxMsgChanLen   uint32 // 读写分离管道最大缓冲数量
}

var GlobalConfObj *Config

func (c *Config) Reload() {
	fileData, err := ioutil.ReadFile("./conf/jarvis.json")
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(fileData, &GlobalConfObj)
	if err != nil {
		panic(err)
	}
}

func init() {
	GlobalConfObj = &Config{
		Name:            "JarvisServer",
		IP:              "0.0.0.0",
		Port:            9999,
		Version:         "v1.1",
		MaxPacketSize:   4096,
		MaxConnNum:      12000,
		WorkerPoolSize:  10,
		MaxTaskQueueLen: 1024,
		MaxMsgChanLen:   1024,
	}
	// GlobalConfObj.Reload()
}
