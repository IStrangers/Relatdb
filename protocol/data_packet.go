package protocol

import (
	"bytes"
	"encoding/binary"
)

type DataPacket interface {
	GetPacketId() byte
	GetDataLength() uint32
	GetPacketBytes() []byte
}

type AbstractDataPacket struct {
	PacketId byte
}

func (self *AbstractDataPacket) GetPacketId() byte {
	return self.PacketId
}

/*
握手包
*/
type HandshakePacket struct {
	AbstractDataPacket
	ProtocolVersion     byte   //协议版本（1个字节）
	ServerVersion       []byte //数据库版本（n个字节,结束补0）
	ConnectionId        uint32 //连接ID（4个字节）
	AuthPluginDataPart1 []byte //认证插件随机数1（8个字节）
	ServerCapabilities  uint16 //数据库支持的功能（2个字节）
	ServerCharsetIndex  byte   //使用的字符集（1个字节）
	ServerStatus        uint16 //数据库状态（2个字节）
	AuthPluginDataPart2 []byte //认证插件随机数2（12位）
}

func (self *HandshakePacket) GetDataLength() uint32 {
	return uint32(1 + len(self.ServerVersion) + 1 + 4 + len(self.AuthPluginDataPart1) + 1 + 2 + 1 + 2 + len(self.AuthPluginDataPart2) + 13 + 1)
}

func (self *HandshakePacket) GetPacketBytes() []byte {
	serverCapabilitiesFiller := []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	var buf bytes.Buffer
	binary.Write(&buf, binary.LittleEndian, self.GetDataLength())
	buf.WriteByte(self.ProtocolVersion)
	buf.Write(self.ServerVersion)
	buf.WriteByte(0)
	binary.Write(&buf, binary.LittleEndian, self.ConnectionId)
	buf.Write(self.AuthPluginDataPart1)
	buf.WriteByte(0)
	binary.Write(&buf, binary.LittleEndian, self.ServerCapabilities)
	buf.WriteByte(self.ServerCharsetIndex)
	binary.Write(&buf, binary.LittleEndian, self.ServerStatus)
	buf.Write(serverCapabilitiesFiller)
	buf.Write(self.AuthPluginDataPart2)
	buf.WriteByte(0)
	return buf.Bytes()
}

/*
二进制数据包
*/
type BinaryPacket struct {
	AbstractDataPacket
	PacketSize uint32
	Data       []byte //数据
}

func (self *BinaryPacket) GetDataLength() uint32 {
	return uint32(len(self.Data))
}

func (self *BinaryPacket) GetPacketBytes() []byte {
	return nil
}

type AuthPacket struct {
	AbstractDataPacket
	UserName []byte //用户名
	Password []byte //密码
}

func (self *AuthPacket) GetDataLength() uint32 {
	return 0
}

func (self *AuthPacket) GetPacketBytes() []byte {
	return nil
}

type ErrorPacket struct {
	AbstractDataPacket
	ErrorCode uint16
	Message   []byte
}

func (self *ErrorPacket) GetDataLength() uint32 {
	return uint32(1 + 1 + 2 + 1 + len([]byte("HY000")) + len(self.Message))
}

func (self *ErrorPacket) GetPacketBytes() []byte {
	var buf bytes.Buffer
	binary.Write(&buf, binary.LittleEndian, self.GetDataLength())
	buf.WriteByte(self.PacketId)
	buf.WriteByte(0xff)
	binary.Write(&buf, binary.LittleEndian, self.ErrorCode)
	buf.WriteByte('#')
	buf.WriteString("HY000")
	buf.Write(self.Message)
	return buf.Bytes()
}
