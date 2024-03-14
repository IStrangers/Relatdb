package protocol

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
)

type DataPacket interface {
	GetPacketId() byte
	GetDataLength() uint32
	GetPacketBytes() []byte
	SendDataPacket(conn net.Conn)
}

type AbstractDataPacket struct {
	PacketId byte
}

func (self *AbstractDataPacket) GetPacketId() byte {
	return self.PacketId
}

func (self *AbstractDataPacket) GetDataLength() uint32 {
	panic("Unrealized method")
}

func (self *AbstractDataPacket) GetPacketBytes() []byte {
	panic("Unrealized method")
}

func (self *AbstractDataPacket) SendDataPacket(conn net.Conn) {
	packetBytes := self.GetPacketBytes()
	_, err := conn.Write(packetBytes)
	if err != nil {
		fmt.Println("Send data packet error:", err.Error())
	}
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
	Data []byte //数据
}

func (self *BinaryPacket) GetDataLength() uint {
	return uint(len(self.Data))
}

func (self *BinaryPacket) GetPacketBytes() []byte {
	return nil
}
