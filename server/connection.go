package server

import (
	"Relatdb/common"
	"Relatdb/protocol"
	"Relatdb/utils"
	"bufio"
	"fmt"
	"net"
)

type Connection struct {
	Server *Relatdb
	Conn   net.Conn
	Reader *bufio.Reader
	Writer *bufio.Writer

	AuthPluginDataPart []byte
}

func NewConnection(server *Relatdb, conn net.Conn) *Connection {
	return &Connection{
		Server: server,
		Conn:   conn,
		Reader: bufio.NewReader(conn),
		Writer: bufio.NewWriter(conn),
	}
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
	self.SendDataPacket(connection, handshakePacket)
}

func (self *Connection) SendDataPacket(connection *Connection, dataPacket protocol.DataPacket) {
	packetBytes := dataPacket.GetPacketBytes()
	err := connection.Write(packetBytes)
	if err != nil {
		fmt.Println("Send data packet error:", err.Error())
	}
}

func (self *Connection) ReceiveBinaryPacket() *protocol.BinaryPacket {
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

func (self *Connection) InitDB(bytesReader *utils.BytesReader) {

}

func (self *Connection) Query(bytesReader *utils.BytesReader) {
	bytes := bytesReader.ReadRemainingBytes()
	querySql := string(bytes)
	println(querySql)
}

func (self *Connection) Ping() {
	self.Write(common.SERVER_OK)
}

func (self *Connection) Close() {
	self.Conn.Close()
}

func (self *Connection) Kill(bytesReader *utils.BytesReader) {
}

func (self *Connection) StmtPrepare(bytesReader *utils.BytesReader) {

}

func (self *Connection) StmtExecute(bytesReader *utils.BytesReader) {

}

func (self *Connection) StmtClose(bytesReader *utils.BytesReader) {

}

func (self *Connection) Heartbeat(bytesReader *utils.BytesReader) {
	self.Ping()
}
