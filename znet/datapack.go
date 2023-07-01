package znet

import (
	"bytes"
	"encoding/binary"
	"errors"

	"github.com/pingtouge2008/zinx/utils"
	"github.com/pingtouge2008/zinx/ziface"
)

type DataPack struct{}

func NewDataPack() *DataPack {
	return &DataPack{}
}

func (d *DataPack) GetHeadLen() uint32 {
	return 4 + 4
}

func (d *DataPack) Pack(msg ziface.IMessage) ([]byte, error) {
	buffer := bytes.NewBuffer([]byte{})

	if err := binary.Write(buffer, binary.LittleEndian, msg.GetDataLen()); err != nil {
		return nil, err
	}

	if err := binary.Write(buffer, binary.LittleEndian, msg.GetMsgId()); err != nil {
		return nil, err
	}

	if err := binary.Write(buffer, binary.LittleEndian, msg.GetData()); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func (d *DataPack) Unpack(binaryData []byte) (ziface.IMessage, error) {
	binaryBuffer := bytes.NewReader(binaryData)
	msg := &Message{}

	if err := binary.Read(binaryBuffer, binary.LittleEndian, &msg.DataLen); err != nil {
		return nil, err
	}

	if err := binary.Read(binaryBuffer, binary.LittleEndian, &msg.Id); err != nil {
		return nil, err
	}

	max := utils.GlobalObject.MaxPacketSize
	if max > 0 && msg.DataLen > max {
		return nil, errors.New("too large msg received")
	}

	return msg, nil
}
