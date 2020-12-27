package test

import (
	"fmt"
	"net"
	"testing"
	"time"

	"jarvis/server"
)

/**
* @Time       : 2020/12/27
* @Author     : xumamba
* @Description: jarvis_test.go
 */

func ClientTest() {
	time.Sleep(1 * time.Second)
	conn, err := net.Dial("tcp4", "0.0.0.0:9999")
	if err != nil {
		fmt.Println("[client] net dial error: ", err.Error())
		return
	}
	var flag = 0
	for {
		if _, err := conn.Write([]byte("hello jarvis v0.2")); err != nil {
			fmt.Println("[client] write error: ", err.Error())
			return
		}

		bytes := make([]byte, 512)
		if _, err = conn.Read(bytes); err != nil {
			fmt.Println("[client] read error: ", err.Error())
			return
		}
		fmt.Println("[client] receive server call back : ", string(bytes))

		flag++
		time.Sleep(1 * time.Second)
		if flag == 20{
			// conn.Close()
		}
	}
}

func TestJarvis(t *testing.T) {
	ser := server.NewServer("jarvis")

	go ClientTest()

	ser.Serve()
}