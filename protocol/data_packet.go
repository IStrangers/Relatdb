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

func (self *AbstractDataPacket) GetDataLength() uint32 {
	panic("Unrealized method")
}

func (self *AbstractDataPacket) GetPacketBytes() []byte {
	panic("Unrealized method")
}

type HandshakePacket struct {
	AbstractDataPacket
	ProtocolVersion     byte
	ServerVersion       []byte
	ConnectionId        uint32
	AuthPluginDataPart1 []byte
	ServerCapabilities  uint16
	ServerCharsetIndex  byte
	ServerStatus        uint16
	AuthPluginDataPart2 []byte
}

func (self *HandshakePacket) GetDataLength() uint32 {
	return uint32(1 + 1 + len(self.ServerVersion) + 1 + 4 + len(self.AuthPluginDataPart1) + 1 + 2 + 1 + 2 + len(self.AuthPluginDataPart2) + 13)
}

func (self *HandshakePacket) GetPacketBytes() []byte {
	serverCapabilitiesFiller := []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}

	var buf bytes.Buffer
	binary.Write(&buf, binary.LittleEndian, self.GetDataLength())
	buf.WriteByte(0)
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
	buf.WriteByte(0)
	return buf.Bytes()
}

type BinaryPacket struct {
	AbstractDataPacket
	Data []byte
}

func (self *BinaryPacket) GetDataLength() uint {
	return uint(len(self.Data))
}

func (self *BinaryPacket) GetPacketBytes() []byte {
	return nil
}
