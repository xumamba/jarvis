package main

/**
 * @DateTime   : 2020/12/28
 * @Author     : xumamba
 * @Description:
 **/

import (
	"fmt"
	"net"
	"time"
)

func main() {
	conn, err := net.Dial("tcp4", "0.0.0.0:9999")
	if err != nil {
		fmt.Println("[client] net dial error: ", err.Error())
		return
	}
	var flag = 0
	for {
		if _, err := conn.Write([]byte("hello jarvis v0.3")); err != nil {
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
		if flag == 20 {
			// conn.Close()
		}
	}
}
