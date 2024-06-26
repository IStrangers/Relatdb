package server

import (
	"Relatdb/parser"
	"Relatdb/parser/ast"
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
	closed bool

	authPluginDataPart []byte
	clientCapabilities uint32
	userName           string
	database           string
	session            *Session
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

func (self *Connection) GetDatabase() string {
	return self.database
}

func (self *Connection) SetDatabase(database string) {
	self.database = database
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
	errorPacket := &ErrorPacket{}
	errorPacket.PacketId = packetId
	errorPacket.ErrorHeader = ERR_HEADER
	errorPacket.ErrorCode = errorCode
	errorPacket.SqlStateMarker = '#'
	errorPacket.SqlState = []byte("HY000")
	errorPacket.Message = []byte(message)
	packetBytes := errorPacket.GetPacketBytes()
	self.write(packetBytes)
}

func (self *Connection) sendHandshakePacket() {
	handshakePacket := &HandshakePacket{
		ProtocolVersion:     PROTOCOL_VERSION,
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

func (self *Connection) sendOkPacket(packetId byte, affectedRows uint64, insertId uint64) {
	okPacket := &OkPacket{}
	okPacket.PacketId = packetId
	okPacket.OkHeader = OK_HEADER
	okPacket.AffectedRows = affectedRows
	okPacket.InsertId = insertId
	okPacket.ServerStatus = 2
	okPacket.WarningCount = 0
	self.sendDataPacket(okPacket)
}

func (self *Connection) sendDataPacket(dataPacket DataPacket) {
	packetBytes := dataPacket.GetPacketBytes()
	self.write(packetBytes)
}

func (self *Connection) receiveBinaryPacket() *BinaryPacket {
	bytes := self.ReadBySize(3)
	packetSize := utils.Uint32(bytes, false)
	if packetSize <= 0 || packetSize > MAX_PACKET_SIZE {
		fmt.Println("Received packet size error:", packetSize)
		return nil
	}
	packetId := self.readByte()
	data := self.ReadBySize(packetSize)
	binaryPacket := &BinaryPacket{}
	binaryPacket.PacketId = packetId
	binaryPacket.Data = data
	return binaryPacket
}

func (self *Connection) authentication() bool {
	binaryPacket := self.receiveBinaryPacket()
	if binaryPacket == nil {
		return false
	}
	authPacket := LoadAuthPacket(binaryPacket)
	if !checkUserNamePassword(authPacket.UserName, authPacket.Password, self.authPluginDataPart) {
		self.writeErrorMessage(2, ER_ACCESS_DENIED_ERROR, fmt.Sprintf("Access denied for user '%s'", authPacket.UserName))
		return false
	}
	self.clientCapabilities = authPacket.ClientFlags
	self.userName = authPacket.UserName
	self.database = authPacket.DataBase
	if self.database == "" {
		self.database = "default"
	}
	self.session = NewSession()
	self.writeAuthOK()
	return true
}

func checkUserNamePassword(userName string, password []byte, authPluginDataPart []byte) bool {
	if userName != SERVER_ROOT_USERNAME || password == nil || len(password) == 0 {
		return false
	}
	rootPassword := scramble411([]byte(SERVER_ROOT_PASSWORD), authPluginDataPart)
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
	self.write(SERVER_AUTH_OK)
}

func (self *Connection) receiveCommandHandler() {
	for !self.closed {
		binaryPacket := self.receiveBinaryPacket()
		if binaryPacket == nil {
			continue
		}
		bytesReader := utils.NewBytesReader(binaryPacket.Data)
		switch bytesReader.ReadByte() {
		case COM_INIT_DB:
			self.handlingInitDB(string(bytesReader.ReadRemainingBytes()))
			break
		case COM_QUERY:
			self.handlingQuery(string(bytesReader.ReadRemainingBytes()))
			break
		case COM_PING:
			self.ping()
			break
		case COM_QUIT:
			self.close()
			break
		case COM_PROCESS_KILL:
			self.kill()
			break
		case COM_STMT_PREPARE:
			self.handlingStmtPrepare()
			break
		case COM_STMT_EXECUTE:
			self.handlingStmtExecute()
			break
		case COM_STMT_CLOSE:
			self.handlingStmtClose()
			break
		case COM_HEARTBEAT:
			self.heartbeat()
			break
		default:
			self.writeErrorMessage(1, ER_UNKNOWN_COM_ERROR, "Unknown command")
			break
		}
	}
}

func (self *Connection) writeOk() {
	self.write(SERVER_OK)
}

func (self *Connection) ping() {
	self.writeOk()
}

func (self *Connection) heartbeat() {
	self.ping()
}

func (self *Connection) close() {
	self.server.removeConn(self.connId)
	err := self.conn.Close()
	if err != nil {
		log.Println("conn close error:", err)
	}
	self.closed = true
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
			errorMsg := fmt.Sprintf("handling sql error: sql=%s, err=%v\n", querySql, err)
			log.Printf(errorMsg)
		}
	}()
	log.Printf("handling query: sql=%s", querySql)
	parser := parser.CreateParser(1, querySql, true, true)
	stmts := parser.Parse()
	stmtLength := len(stmts)
	if stmtLength > 1 && self.clientCapabilities&CLIENT_MULTI_STATEMENTS == 0 {
		//return
	}
	ctx := NewContext(self)
	for i, stmt := range stmts {
		self.handlingStmt(ctx, stmt, stmtLength-1 == i)
	}
}

func (self *Connection) handlingStmt(ctx *Context, stmt ast.Statement, isLastStmt bool) {
	recordSet := ctx.executeStmt(stmt)
	columns := recordSet.GetColumns()
	rows := recordSet.GetRows()
	if len(columns) != 0 || len(rows) != 0 {
		selectPacket := NewTablePacket(columns, rows)
		self.sendDataPacket(selectPacket)
		return
	}
	self.sendOkPacket(0, recordSet.GetAffectedRows(), recordSet.GetInsertId())
}

func (self *Connection) handlingStmtPrepare() {

}

func (self *Connection) handlingStmtExecute() {

}

func (self *Connection) handlingStmtClose() {

}
