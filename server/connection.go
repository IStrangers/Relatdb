package server

import (
	"Relatdb/common"
	"Relatdb/parser"
	"Relatdb/protocol"
	"Relatdb/utils"
	"bufio"
	"bytes"
	"crypto/sha1"
	"fmt"
	"log"
	"net"
)

type Connection struct {
	server *Server
	connId uint64
	conn   net.Conn
	reader *bufio.Reader
	writer *bufio.Writer

	authPluginDataPart []byte
	clientCapabilities uint32
	userName           string
	database           string
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

func (self *Connection) read(bytes []byte) []byte {
	_, err := self.reader.Read(bytes)
	if err != nil {
		log.Println("conn read error:", err)
	}
	return bytes
}

func (self *Connection) ReadBySize(size uint32) []byte {
	return self.read(make([]byte, size))
}

func (self *Connection) readByte() byte {
	return self.ReadBySize(1)[0]
}

func (self *Connection) write(bytes []byte) {
	_, err := self.writer.Write(bytes)
	if err != nil {
		log.Println("conn write error:", err)
	}
	err = self.writer.Flush()
	if err != nil {
		log.Println("conn write flush error:", err)
	}
}

func (self *Connection) writeErrorMessage(packetId byte, errorCode uint16, message string) {
	errorPacket := &protocol.ErrorPacket{}
	errorPacket.PacketId = packetId
	errorPacket.FieldCount = 0xff
	errorPacket.ErrorCode = errorCode
	errorPacket.SqlStateMarker = '#'
	errorPacket.SqlState = []byte("HY000")
	errorPacket.Message = []byte(message)
	packetBytes := errorPacket.GetPacketBytes()
	self.write(packetBytes)
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
	self.write(packetBytes)
}

func (self *Connection) receiveBinaryPacket() *protocol.BinaryPacket {
	bytes := self.ReadBySize(3)
	packetSize := utils.Uint32(bytes, false)
	if packetSize <= 0 || packetSize > common.MAX_PACKET_SIZE {
		fmt.Println("Received packet size error:", packetSize)
		return nil
	}
	packetId := self.readByte()
	data := self.ReadBySize(packetSize)
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
	self.clientCapabilities = authPacket.ClientFlags
	self.userName = authPacket.UserName
	self.writeAuthOK()
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

func (self *Connection) writeAuthOK() {
	self.write(common.SERVER_AUTH_OK)
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
			self.handlingInitDB(string(bytesReader.ReadRemainingBytes()))
			break
		case common.COM_QUERY:
			self.handlingQuery(string(bytesReader.ReadRemainingBytes()))
			break
		case common.COM_PING:
			self.ping()
			break
		case common.COM_QUIT:
			self.close()
			break
		case common.COM_PROCESS_KILL:
			self.kill()
			break
		case common.COM_STMT_PREPARE:
			self.handlingStmtPrepare()
			break
		case common.COM_STMT_EXECUTE:
			self.handlingStmtExecute()
			break
		case common.COM_STMT_CLOSE:
			self.handlingStmtClose()
			break
		case common.COM_HEARTBEAT:
			self.heartbeat()
			break
		default:
			self.writeErrorMessage(1, common.ER_UNKNOWN_COM_ERROR, "Unknown command")
			break
		}
	}
}

func (self *Connection) writeOk() {
	self.write(common.SERVER_OK)
}

func (self *Connection) ping() {
	self.writeOk()
}

func (self *Connection) heartbeat() {
	self.ping()
}

func (self *Connection) close() {
	err := self.conn.Close()
	if err != nil {
		log.Println("conn close error:", err)
	}
}

func (self *Connection) kill() {
}

func (self *Connection) handlingInitDB(database string) {
	log.Println("use database:", database)
	self.database = database
	self.writeOk()
}

func (self *Connection) handlingQuery(querySql string) {
	defer func() {
		err := recover()
		if err != nil {
			log.Printf("parser query sql error: sql=%s, err=%v\n", querySql, err)
		}
	}()
	parser := parser.CreateParser(1, querySql, true, true)
	stmts := parser.Parse()
	for _, stmt := range stmts {
		println(stmt)
	}
}

func (self *Connection) handlingStmtPrepare() {

}

func (self *Connection) handlingStmtExecute() {

}

func (self *Connection) handlingStmtClose() {

}
