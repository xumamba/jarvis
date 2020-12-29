package main

/**
 * @DateTime   : 2020/12/28
 * @Author     : xumamba
 * @Description:
 **/

import (
	"fmt"
	"io"
	"net"
	"time"

	"jarvis/server"
)

func main() {
	conn, err := net.Dial("tcp4", "0.0.0.0:9999")
	if err != nil {
		fmt.Println("[client] net dial error: ", err.Error())
		return
	}
	var flag = 0
	for {
		msg, _ := server.DPHelper.PackageMsg(server.NewMessage(0, []byte("hello jarvis v0.7")))
		if _, err := conn.Write(msg); err != nil {
			fmt.Println("[client] write error: ", err.Error())
			return
		}

		msgHead := make([]byte, server.DPHelper.GetHeadLen())
		if _, err := io.ReadFull(conn, msgHead); err != nil {
			fmt.Println("receive buf error: " + err.Error())
			break
		}
		res, err := server.DPHelper.UnPackageMsg(msgHead)
		if err != nil {
			fmt.Println("unpack error: " + err.Error())

			break
		}
		data := make([]byte, res.GetMsgLen())
		if _, err := io.ReadFull(conn, data); err != nil {
			fmt.Println("read msg body error")
			break
		}

		res.SetRealData(data)

		fmt.Printf("[client] receive server call back : msgID=%d, msgLen=%d, msgData=%s\n", res.GetMsgID(),
			res.GetMsgLen(), res.GetRealData())

		flag++
		time.Sleep(1 * time.Second)
		if flag == 20 {
			// conn.Close()
		}
	}
}
