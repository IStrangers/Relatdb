package server

import (
	"Relatdb/common"
	"Relatdb/protocol"
	"Relatdb/utils"
	"fmt"
	"net"
)

type Options struct {
	BindIp   string
	BindPort uint
}

type Relatdb struct {
	version string
	options *Options
	ln      net.Listener
}

func CreateRelatdb(options *Options) *Relatdb {
	server := &Relatdb{
		version: common.RELATDB_VERSION,
		options: options,
	}
	return server
}

func (self *Relatdb) Start() {
	bindAddress := fmt.Sprintf("%s:%d", self.options.BindIp, self.options.BindPort)
	ln, err := net.Listen("tcp", bindAddress)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		return
	}
	self.ln = ln
	fmt.Println("Listening on ", bindAddress)

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			continue
		}
		go self.handlingConnection(conn)
	}
}

func (self *Relatdb) Stop() {
	err := self.ln.Close()
	if err != nil {
		fmt.Println("Listening stop error:", err.Error())
	}
}

func (self *Relatdb) getServerCapabilities() uint16 {
	flag := uint16(0)
	flag |= common.CLIENT_LONG_PASSWORD
	flag |= common.CLIENT_FOUND_ROWS
	flag |= common.CLIENT_LONG_FLAG
	flag |= common.CLIENT_CONNECT_WITH_DB
	// flag |=  common.CLIENT_NO_SCHEMA;
	// flag |=  common.CLIENT_COMPRESS;
	flag |= common.CLIENT_ODBC
	// flag |=  common.CLIENT_LOCAL_FILES;
	flag |= common.CLIENT_IGNORE_SPACE
	flag |= common.CLIENT_PROTOCOL_41
	flag |= common.CLIENT_INTERACTIVE
	// flag |=  common.CLIENT_SSL;
	flag |= common.CLIENT_IGNORE_SIGPIPE
	flag |= common.CLIENT_TRANSACTIONS
	flag |= common.CLIENT_SECURE_CONNECTION
	return flag
}

func (self *Relatdb) handlingConnection(conn net.Conn) {
	connection := &Connection{conn: conn}
	self.sendHandshakePacket(connection)
	if !self.authentication(connection) {
		return
	}
	self.receiveDataHandler(connection)
}

func (self *Relatdb) sendHandshakePacket(connection *Connection) {
	handshakePacket := &protocol.HandshakePacket{
		ProtocolVersion:     common.PROTOCOL_VERSION,
		ServerVersion:       []byte(self.version),
		ConnectionId:        1,
		AuthPluginDataPart1: utils.RandomBytes(8),
		ServerCapabilities:  self.getServerCapabilities(),
		ServerCharsetIndex:  33,
		ServerStatus:        2,
		AuthPluginDataPart2: utils.RandomBytes(12),
	}
	self.sendDataPacket(connection, handshakePacket)
}

func (self *Relatdb) sendDataPacket(connection *Connection, dataPacket protocol.DataPacket) {
	packetBytes := dataPacket.GetPacketBytes()
	err := connection.Write(packetBytes)
	if err != nil {
		fmt.Println("Send data packet error:", err.Error())
	}
}

func (self *Relatdb) receiveBinaryPacket(connection *Connection) *protocol.BinaryPacket {
	bytes, _ := connection.ReadBySize(3)
	packetSize := utils.Uint32(bytes)
	if packetSize <= 0 || packetSize > common.MAX_PACKET_SIZE {
		fmt.Println("Received packet size error:", packetSize)
		return nil
	}
	packetId, _ := connection.ReadByte()
	data, _ := connection.ReadBySize(packetSize)
	binaryPacket := &protocol.BinaryPacket{}
	binaryPacket.PacketSize = packetSize
	binaryPacket.PacketId = packetId
	binaryPacket.Data = data
	return binaryPacket
}

func (self *Relatdb) authentication(connection *Connection) bool {
	binaryPacket := self.receiveBinaryPacket(connection)
	if binaryPacket == nil {
		return false
	}
	if true {
		connection.WriteErrorMessage(2, common.ER_ACCESS_DENIED_ERROR, fmt.Sprintf("Access denied for user '%s'", ""))
		return false
	}
	connection.Write(common.SERVER_AUTH_OK)
	return true
}

func (self *Relatdb) receiveDataHandler(connection *Connection) {
	for {
		binaryPacket := self.receiveBinaryPacket(connection)
		if binaryPacket == nil {
			continue
		}

	}
}
