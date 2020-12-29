package conf

/**
 * @DateTime   : 2020/12/28
 * @Author     : xumamba
 * @Description: 全局配置文件
 **/

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

// Config jarvis框架配置
type Config struct {
	ConfFilePath string // 配置文件路径

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

// GlobalConfObj 全局配置对象
var GlobalConfObj *Config

// PathExists 判断路径是否存在
func PathExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	return false
}

// Reload 根据配置文件重载全局配置对象
func (c *Config) Reload() {
	if exists := PathExists(c.ConfFilePath); !exists {
		return
	}
	fileData, err := ioutil.ReadFile(c.ConfFilePath)
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
		ConfFilePath:    "./conf/jarvis.json",
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
	GlobalConfObj.Reload()
}
