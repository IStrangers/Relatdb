package server

import (
	"Relatdb/common"
	"Relatdb/protocol"
	"Relatdb/utils"
	"bufio"
	"bytes"
	"crypto/sha1"
	"fmt"
	"net"
)

type Connection struct {
	server *Server
	connId uint64
	conn   net.Conn
	reader *bufio.Reader
	writer *bufio.Writer

	authPluginDataPart []byte
}

func NewConnection(server *Server, conn net.Conn) *Connection {
	return &Connection{
		server: server,
		connId: server.nextConnId(),
		conn:   conn,
		reader: bufio.NewReader(conn),
		writer: bufio.NewWriter(conn),
	}
}

func (self *Connection) read(bytes []byte) ([]byte, error) {
	_, err := self.reader.Read(bytes)
	return bytes, err
}

func (self *Connection) ReadBySize(size uint32) ([]byte, error) {
	return self.read(make([]byte, size))
}

func (self *Connection) readByte() (byte, error) {
	bytes, err := self.ReadBySize(1)
	if err != nil {
		return 0, err
	}
	return bytes[0], nil
}

func (self *Connection) write(bytes []byte) error {
	_, err := self.writer.Write(bytes)
	if err == nil {
		err = self.writer.Flush()
	}
	return err
}

func (self *Connection) writeErrorMessage(packetId byte, errorCode uint16, message string) error {
	errorPacket := &protocol.ErrorPacket{}
	errorPacket.PacketId = packetId
	errorPacket.FieldCount = 0xff
	errorPacket.ErrorCode = errorCode
	errorPacket.SqlStateMarker = '#'
	errorPacket.SqlState = []byte("HY000")
	errorPacket.Message = []byte(message)
	packetBytes := errorPacket.GetPacketBytes()
	return self.write(packetBytes)
}

func (self *Connection) sendHandshakePacket() {
	handshakePacket := &protocol.HandshakePacket{
		ProtocolVersion:     common.PROTOCOL_VERSION,
		ServerVersion:       []byte(self.server.version),
		ConnectionId:        1,
		AuthPluginDataPart1: utils.RandomBytes(8),
		ServerCapabilities:  self.server.getServerCapabilities(),
		ServerCharsetIndex:  33,
		ServerStatus:        2,
		AuthPluginDataPart2: utils.RandomBytes(12),
		AuthPluginName:      []byte("mysql_native_password"),
	}
	self.authPluginDataPart = append(handshakePacket.AuthPluginDataPart1, handshakePacket.AuthPluginDataPart2...)
	self.sendDataPacket(handshakePacket)
}

func (self *Connection) sendDataPacket(dataPacket protocol.DataPacket) {
	packetBytes := dataPacket.GetPacketBytes()
	err := self.write(packetBytes)
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
	packetId, _ := self.readByte()
	data, _ := self.ReadBySize(packetSize)
	binaryPacket := &protocol.BinaryPacket{}
	binaryPacket.PacketSize = packetSize
	binaryPacket.PacketId = packetId
	binaryPacket.Data = data
	return binaryPacket
}

func (self *Connection) authentication() bool {
	binaryPacket := self.receiveBinaryPacket()
	if binaryPacket == nil {
		return false
	}
	authPacket := protocol.LoadAuthPacket(binaryPacket)
	if !checkUserNamePassword(authPacket.UserName, authPacket.Password, self.authPluginDataPart) {
		self.writeErrorMessage(2, common.ER_ACCESS_DENIED_ERROR, fmt.Sprintf("Access denied for user '%s'", authPacket.UserName))
		return false
	}
	self.authOK()
	return true
}

func checkUserNamePassword(userName string, password []byte, authPluginDataPart []byte) bool {
	if userName != common.SERVER_ROOT_USERNAME || password == nil || len(password) == 0 {
		return false
	}
	rootPassword := scramble411([]byte(common.SERVER_ROOT_PASSWORD), authPluginDataPart)
	return bytes.Equal(rootPassword, password)
}

func scramble411(data []byte, seed []byte) []byte {
	crypt := sha1.New()

	crypt.Write(data)
	stage1 := crypt.Sum(nil)

	crypt.Reset()
	crypt.Write(stage1)
	stage2 := crypt.Sum(nil)

	crypt.Reset()
	crypt.Write(seed)
	crypt.Write(stage2)
	stage3 := crypt.Sum(nil)
	for i := range stage3 {
		stage3[i] ^= stage1[i]
	}

	return stage3
}

func (self *Connection) receiveCommandHandler() {
	for {
		binaryPacket := self.receiveBinaryPacket()
		if binaryPacket == nil {
			continue
		}
		bytesReader := utils.NewBytesReader(binaryPacket.Data)
		switch bytesReader.ReadByte() {
		case common.COM_INIT_DB:
			self.initDB(bytesReader)
			break
		case common.COM_QUERY:
			self.query(bytesReader)
			break
		case common.COM_PING:
			self.ping()
			break
		case common.COM_QUIT:
			self.close()
			break
		case common.COM_PROCESS_KILL:
			self.kill(bytesReader)
			break
		case common.COM_STMT_PREPARE:
			self.stmtPrepare(bytesReader)
			break
		case common.COM_STMT_EXECUTE:
			self.stmtExecute(bytesReader)
			break
		case common.COM_STMT_CLOSE:
			self.stmtClose(bytesReader)
			break
		case common.COM_HEARTBEAT:
			self.heartbeat(bytesReader)
			break
		default:
			self.writeErrorMessage(1, common.ER_UNKNOWN_COM_ERROR, "Unknown command")
			break
		}
	}
}

func (self *Connection) authOK() {
	self.write(common.SERVER_AUTH_OK)
}

func (self *Connection) initDB(bytesReader *utils.BytesReader) {

}

func (self *Connection) query(bytesReader *utils.BytesReader) {
	bytes := bytesReader.ReadRemainingBytes()
	querySql := string(bytes)
	println(querySql)
}

func (self *Connection) ping() {
	self.write(common.SERVER_OK)
}

func (self *Connection) close() {
	self.conn.Close()
}

func (self *Connection) kill(bytesReader *utils.BytesReader) {
}

func (self *Connection) stmtPrepare(bytesReader *utils.BytesReader) {

}

func (self *Connection) stmtExecute(bytesReader *utils.BytesReader) {

}

func (self *Connection) stmtClose(bytesReader *utils.BytesReader) {

}

func (self *Connection) heartbeat(bytesReader *utils.BytesReader) {
	self.ping()
}
