package server

import (
	"Relatdb/common"
	"Relatdb/protocol"
	"bufio"
)

type Connection struct {
	Reader *bufio.Reader
	Writer *bufio.Writer

	AuthPluginDataPart []byte
}

func (self *Connection) Read(bytes []byte) ([]byte, error) {
	_, err := self.Reader.Read(bytes)
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
	_, err := self.Writer.Write(bytes)
	if err == nil {
		err = self.Writer.Flush()
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

func (self *Connection) InitDB(packet *protocol.BinaryPacket) {

}

func (self *Connection) Query(packet *protocol.BinaryPacket) {

}

func (self *Connection) Ping() {
	self.Write(common.SERVER_OK)
}

func (self *Connection) Close() {

}

func (self *Connection) Kill(packet *protocol.BinaryPacket) {

}

func (self *Connection) StmtPrepare(packet *protocol.BinaryPacket) {

}

func (self *Connection) StmtExecute(packet *protocol.BinaryPacket) {

}

func (self *Connection) StmtClose(packet *protocol.BinaryPacket) {

}

func (self *Connection) Heartbeat(packet *protocol.BinaryPacket) {
	self.Ping()
}
