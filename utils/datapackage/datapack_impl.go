package datapackage

/**
 * @DateTime   : 2020/12/24
 * @Author     : xumamba
 * @Description:
 **/

import (
	"bytes"
	"encoding/binary"
	"errors"
)

// DefaultMaxPackByte 最大数据包大小 10MB
const DefaultMaxPackByte = 10 * 1024 * 1024

type DataPack struct {
	MaxPackByte uint32
}

// GetHeadLen 获取数据包头部大小
func (d *DataPack) GetHeadLen() uint32 {
	// ID uint32(4字节) + DataLen uint32(4字节)
	return 8
}

// PackageMsg 数据封包
func (d *DataPack) PackageMsg(msg IMessage) ([]byte, error) {
	// 判断数据包大小是否超过最大限制
	if msg.GetMsgLen() > d.MaxPackByte {
		return nil, errors.New("package msg length too long")
	}

	// 创建一个存放bytes字节的缓冲
	buffer := bytes.NewBuffer([]byte{})
	// 写入数据长度
	if err := binary.Write(buffer, binary.BigEndian, msg.GetMsgLen()); err != nil {
		return nil, err
	}
	// 写入数据标识
	if err := binary.Write(buffer, binary.BigEndian, msg.GetMsgID()); err != nil{
		return nil, err
	}
	// 写入数据体
	if err := binary.Write(buffer, binary.BigEndian, msg.GetRealData()); err != nil{
		return nil, err
	}
	return buffer.Bytes(), nil
}

// UnPackageMsg 数据解包，只把包头信息读取出来，调用者根据head长度标识从conn读取数据体。
func (d *DataPack) UnPackageMsg(data []byte) (IMessage, error) {
	// 创建一个io reader
	reader := bytes.NewReader(data)
	msg := &Message{}

	// 读取数据包 包头
	if err := binary.Read(reader, binary.BigEndian, &msg.Length); err != nil{
		return nil, err
	}
	if err := binary.Read(reader, binary.BigEndian, &msg.ID); err != nil{
		return nil, err
	}

	// 判断数据包大小是否超过最大限制
	if msg.GetMsgLen() > d.MaxPackByte {
		return nil, errors.New("package msg length too long")
	}

	return msg, nil

}

func NewDataPack(maxPackByte uint32) *DataPack {
	return &DataPack{
		MaxPackByte: maxPackByte,
	}
}

var DPHelper = NewDataPack(DefaultMaxPackByte)