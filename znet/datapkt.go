package znet

import (
	"bytes"
	"encoding/binary"
	"errors"
	"github.com/ashin9/zinx/utils"
	"github.com/ashin9/zinx/ziface"
)

type DataPack struct {
}

func NewDataPack() *DataPack {
	return &DataPack{}
}

func (dp *DataPack) GetHeadLen() uint32 {
	return 8
}

func (dp *DataPack) Pack(msg ziface.IMessage) ([]byte, error) {
	dataBuff := bytes.NewBuffer([]byte{})
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetMsgLen()); err != nil {
		return nil, err
	}
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetMsgId()); err != nil {
		return nil, err
	}
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetData()); err != nil {
		return nil, err
	}
	return dataBuff.Bytes(), nil
}

func (dp *DataPack) UnPack(binData []byte) (ziface.IMessage, error) {
	dataBuff := bytes.NewReader(binData)
	// 只解压 Head 信息, 得到 dataLen 和 MsgID
	msg := &Message{}
	// 读 dataLen 和 id
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.MsgLen); err != nil {
		return nil, err
	}
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.MsgId); err != nil {
		return nil, err
	}
	// 是否超出长度
	if utils.GlobalObj.MaxPackageSize > 0 && msg.MsgLen > utils.GlobalObj.MaxPackageSize {
		return nil, errors.New("msg too long")
	}

	return msg, nil
}
