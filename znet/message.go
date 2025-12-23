package znet

type Message struct {
	MsgId  uint32
	MsgLen uint32
	Data   []byte
}

func (m *Message) GetMsgId() uint32 {
	return m.MsgId
}

func (m *Message) GetMsgLen() uint32 {
	return m.MsgLen
}

func (m *Message) GetData() []byte {
	return m.Data
}

func (m *Message) SetMsgId(msgId uint32) {
	m.MsgId = msgId
}

func (m *Message) SetMsgLen(msgLen uint32) {
	m.MsgLen = msgLen
}

func (m *Message) SetData(data []byte) {
	m.Data = data
}
