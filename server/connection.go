package server

import (
	"Relatdb/protocol"
	"bufio"
)

type Connection struct {
	reader *bufio.Reader
	writer *bufio.Writer
}

func (self *Connection) Read(bytes []byte) ([]byte, error) {
	_, err := self.reader.Read(bytes)
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
	_, err := self.writer.Write(bytes)
	if err == nil {
		err = self.writer.Flush()
	}
	return err
}

func (self *Connection) WriteErrorMessage(packetId byte, errorCode uint16, message string) error {
	errorPacket := &protocol.ErrorPacket{}
	errorPacket.PacketId = packetId
	errorPacket.FieldCount = 0xff
	errorPacket.ErrorCode = errorCode
	errorPacket.SqlStateMarker = '#'
	errorPacket.SqlState = []byte("HY000")
	errorPacket.Message = []byte(message)
	packetBytes := errorPacket.GetPacketBytes()
	return self.Write(packetBytes)
}
