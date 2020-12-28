package server

/**
* @DateTime   : 2020/12/28
* @Author     : xumamba
* @Description: 请求消息对象
**/

type Message struct {
	Length uint32
	ID uint32
	Data []byte
}

func (m *Message) GetMsgLen() uint32 {
	return m.Length
}

func (m *Message) GetMsgID() uint32 {
	return m.ID
}

func (m *Message) GetRealData() []byte {
	return m.Data
}

func (m *Message) SetMsgLen(len uint32) {
	m.Length = len
}

func (m *Message) SetMsgID(id uint32) {
	m.ID = id
}

func (m *Message) SetRealData(data []byte) {
	m.Data = data
}

func NewMessage(msgID uint32, data []byte) *Message {
	return &Message{
		Length: uint32(len(data)),
		ID:     msgID,
		Data:   data,
	}
}