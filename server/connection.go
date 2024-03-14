package server

import (
	"Relatdb/protocol"
	"net"
)

type Connection struct {
	conn net.Conn
}

func (self *Connection) Read(bytes []byte) ([]byte, error) {
	_, err := self.conn.Read(bytes)
	return bytes, err
}

func (self *Connection) ReadBySize(size uint32) ([]byte, error) {
	return self.Read(make([]byte, size))
}

func (self *Connection) ReadByte() (byte, error) {
	bytes, err := self.ReadBySize(1)
	if err != nil {
		return 0, err
	}
	return bytes[0], nil
}

func (self *Connection) Write(bytes []byte) error {
	_, err := self.conn.Write(bytes)
	return err
}

func (self *Connection) WriteErrorMessage(packetId byte, errorCode uint16, message string) error {
	errorPacket := &protocol.ErrorPacket{}
	errorPacket.PacketId = packetId
	errorPacket.ErrorCode = errorCode
	errorPacket.Message = []byte(message)
	packetBytes := errorPacket.GetPacketBytes()
	return self.Write(packetBytes)
}
