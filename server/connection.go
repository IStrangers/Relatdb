package server

import (
	"Relatdb/common"
	"Relatdb/protocol"
	"Relatdb/utils"
	"bufio"
	"fmt"
)

type Connection struct {
	Server *Relatdb
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

func (self *Connection) SendHandshakePacket(connection *Connection) {
	handshakePacket := &protocol.HandshakePacket{
		ProtocolVersion:     common.PROTOCOL_VERSION,
		ServerVersion:       []byte(self.Server.version),
		ConnectionId:        1,
		AuthPluginDataPart1: utils.RandomBytes(8),
		ServerCapabilities:  self.Server.getServerCapabilities(),
		ServerCharsetIndex:  33,
		ServerStatus:        2,
		AuthPluginDataPart2: utils.RandomBytes(12),
	}
	connection.AuthPluginDataPart = append(handshakePacket.AuthPluginDataPart1, handshakePacket.AuthPluginDataPart2...)
	self.sendDataPacket(connection, handshakePacket)
}

func (self *Connection) sendDataPacket(connection *Connection, dataPacket protocol.DataPacket) {
	packetBytes := dataPacket.GetPacketBytes()
	err := connection.Write(packetBytes)
	if err != nil {
		fmt.Println("Send data packet error:", err.Error())
	}
}

func (self *Connection) receiveBinaryPacket() *protocol.BinaryPacket {
	bytes, _ := self.ReadBySize(3)
	packetSize := utils.Uint32(bytes, false)
	if packetSize <= 0 || packetSize > common.MAX_PACKET_SIZE {
		fmt.Println("Received packet size error:", packetSize)
		return nil
	}
	packetId, _ := self.ReadByte()
	data, _ := self.ReadBySize(packetSize)
	binaryPacket := &protocol.BinaryPacket{}
	binaryPacket.PacketSize = packetSize
	binaryPacket.PacketId = packetId
	binaryPacket.Data = data
	return binaryPacket
}

func (self *Connection) AuthOK() {
	self.Write(common.SERVER_AUTH_OK)
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
