package datapackage

/**
* @DateTime   : 2020/12/24
* @Author     : xumamba
* @Description:
**/

type Message struct {
	Length uint32 // 消息长度
	ID     uint32 // 消息唯一标识
	Data   []byte // 消息内容
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

