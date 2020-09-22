/**
* @DateTime   : 2020/9/18 15:18
* @Author     : xumamba
* @Description:
**/
package server

var debugMode = false

func IsDebugMode() bool {
	return debugMode
}

func EnableDebug() {
	debugMode = true
}

func DisableDebug() {
	debugMode = false
}
