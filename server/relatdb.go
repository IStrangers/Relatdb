package server

import (
	"Relatdb/common"
	"Relatdb/protocol"
	"Relatdb/utils"
	"bufio"
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
	return &Relatdb{
		version: common.Relatdb_Version,
		options: options,
	}
}

func (self *Relatdb) Start() {
	bindAddress := fmt.Sprintf("%s:%d", self.options.BindIp, self.options.BindPort)
	ln, err := net.Listen("tcp", bindAddress)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		return
	}
	fmt.Println("Listening on ", bindAddress)
	self.ln = ln

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			continue
		}
		go self.handleClient(conn)
	}
}

func (self *Relatdb) Stop() {
	err := self.ln.Close()
	if err != nil {
		fmt.Println("Listening stop error:", err.Error())
	}
}

func (self *Relatdb) handleClient(conn net.Conn) {
	self.sendHandshakePacket(conn)
	reader := bufio.NewReader(conn)
	for {
		lengthBytes := make([]byte, 3)
		_, err := reader.Read(lengthBytes)
		if err != nil {
			fmt.Println("Received data fail")
			continue
		}
		length := utils.Uint32(lengthBytes)
		if length <= 0 {
			continue
		}
		data := make([]byte, length)
		_, err = reader.Read(data)
		if err != nil {
			fmt.Println("Received data fail")
			continue
		}
		println(string(data))
	}
}

func (self *Relatdb) sendDataPacket(conn net.Conn, dataPacket protocol.DataPacket) {
	packetBytes := dataPacket.GetPacketBytes()
	_, err := conn.Write(packetBytes)
	if err != nil {
		fmt.Println("Send data packet error:", err.Error())
	}
}

func (self *Relatdb) sendHandshakePacket(conn net.Conn) {
	handshakePacket := &protocol.HandshakePacket{
		ProtocolVersion:     common.Protocol_Version,
		ServerVersion:       []byte(common.Relatdb_Version),
		ConnectionId:        1,
		AuthPluginDataPart1: utils.RandomBytes(8),
		ServerCapabilities:  self.getServerCapabilities(),
		ServerCharsetIndex:  33,
		ServerStatus:        2,
		AuthPluginDataPart2: utils.RandomBytes(12),
	}
	self.sendDataPacket(conn, handshakePacket)
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
