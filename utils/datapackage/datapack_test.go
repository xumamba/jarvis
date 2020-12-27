package datapackage

/**
 * @DateTime   : 2020/12/24
 * @Author     : xumamba
 * @Description:
 **/

import (
	"io"
	"net"
	"testing"
	"time"
)

func TestDataPack(t *testing.T) {
	testAddr := "127.0.0.1:8888"
	listener, err := net.Listen("tcp", testAddr)
	if err != nil {
		t.Fatal(err)
	}

	// 服务端
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				t.Log(err)
				continue
			}

			go func(conn net.Conn) {
				for {
					// 读取第一个数据包的包头
					headData := make([]byte, DPHelper.GetHeadLen())
					_, err := io.ReadFull(conn, headData)
					if err != nil {
						t.Log("read head failed: ", err)
						return
					}
					msg, err := DPHelper.UnPackageMsg(headData)
					if err != nil {
						t.Log("unpackage error: ", err)
						return
					}
					if msg.GetMsgLen() > 0 {
						message := msg.(*Message)
						message.Data = make([]byte, msg.GetMsgLen())
						if _, err := io.ReadFull(conn, message.Data); err != nil {
							t.Log("read body failed: ", err)
							return
						}
						t.Logf("[Receive Message]: len: %d, id: %d, data: %s", message.Length, message.ID, message.Data )
					}
				}
			}(conn)
		}
	}()

	// 客户端
	go func() {
		conn, err := net.Dial("tcp", testAddr)
		if err != nil{
			t.Log(err)
		}
		msg1 := &Message{
			Length: 5,
			ID:     1,
			Data:   []byte{'h', 'e', 'l', 'l', 'o'},
		}
		packageMsg1, err := DPHelper.PackageMsg(msg1)
		if err != nil{
			t.Log(err)
			return
		}
		msg2 := &Message{
			Length: 6,
			ID:     2,
			Data:   []byte{'j', 'a', 'r', 'v', 'i', 's'},
		}
		packageMsg2, err := DPHelper.PackageMsg(msg2)
		if err != nil{
			t.Log(err)
			return
		}
		// 模拟tcp粘包
		packageMsg1 = append(packageMsg1, packageMsg2...)

		conn.Write(packageMsg1)
	}()

	select {
	case <-time.After(time.Second):
		return
	}
}
