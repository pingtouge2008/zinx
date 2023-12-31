package znet

type Message struct {
	Id      uint32
	DataLen uint32
	Data    []byte
}

func NewMsgPacket(id uint32, data []byte) *Message {
	return &Message{
		Id:      id,
		DataLen: uint32(len(data)),
		Data:    data,
	}
}
func (m *Message) GetDataLen() uint32 {
	return m.DataLen
}

func (m *Message) GetMsgId() uint32 {
	return m.Id
}

func (m *Message) GetData() []byte {
	return m.Data
}

func (m *Message) SetDataLen(len uint32) {
	m.DataLen = len
}

func (m *Message) SetMsgId(msgId uint32) {
	m.Id = msgId
}

func (m *Message) SetData(data []byte) {
	m.Data = data
}
