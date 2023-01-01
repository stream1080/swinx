package znet

import (
	"bytes"
	"encoding/binary"
	"errors"

	"github.com/stream1080/swinx/conf"
	"github.com/stream1080/swinx/face"
)

// 封包拆包类实例
type DataPack struct{}

func NewDataPack() *DataPack {
	return &DataPack{}
}

// 获取包头长度, Id uint32(4字节) +  DataLen uint32(4字节)
func (d *DataPack) GetHeadLen() uint32 {
	return 8
}

// 封包方法, |dataLen|msgId|data|
func (d *DataPack) Pack(msg face.Message) ([]byte, error) {
	// 创建一个存放bytes字节的缓冲
	dataBuff := bytes.NewBuffer([]byte{})

	// 写 dataLen
	err := binary.Write(dataBuff, binary.LittleEndian, msg.GetDataLen())
	if err != nil {
		return nil, err
	}

	// 写 msgId
	err = binary.Write(dataBuff, binary.LittleEndian, msg.GetMsgId())
	if err != nil {
		return nil, err
	}

	// 写 data
	err = binary.Write(dataBuff, binary.LittleEndian, msg.GetData())
	if err != nil {
		return nil, err
	}

	return dataBuff.Bytes(), nil
}

// 拆包方法
func (d *DataPack) UnPack(byteData []byte) (face.Message, error) {
	// 创建一个从输入二进制数据的ioReader
	dataBuff := bytes.NewReader(byteData)

	// 新建 msg
	msg := &Message{}

	// 解析 Head 信息，获取 dataLen
	err := binary.Read(dataBuff, binary.LittleEndian, &msg.DataLen)
	if err != nil {
		return nil, err
	}

	// 判断是否超过最大长度
	if conf.ServerConfig.MaxPackageSize < msg.DataLen {
		return nil, errors.New("too large msg data recieved")
	}

	// 解析 Head 信息，获取 msgID
	err = binary.Read(dataBuff, binary.LittleEndian, &msg.Id)
	if err != nil {
		return nil, err
	}

	// 这里只需要把 head 的数据拆包出来就可以了，然后再通过 head 的长度，再从 conn 读取一次数据
	return msg, nil
}
