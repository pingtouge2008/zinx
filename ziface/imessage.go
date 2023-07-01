package ziface

type IMessage interface {
	GetDataLen() uint32
	GetMsgId() uint32
	GetData() []byte

	SetDataLen(uint32)
	SetMsgId(uint32)
	SetData([]byte)
}
