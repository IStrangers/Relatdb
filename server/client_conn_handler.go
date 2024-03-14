package server

import (
	"Relatdb/common"
	"Relatdb/protocol"
	"Relatdb/utils"
	"bufio"
	"fmt"
	"net"
)

func buildClientConnHandler(server *Relatdb) ClientConnHandler {
	dataReceiveHandler := &DataReceiveHandler{
		AbstractConnHandler{
			server:      server,
			nextHandler: nil,
		},
	}
	return &HandshakeHandler{
		AbstractConnHandler{
			server:      server,
			nextHandler: dataReceiveHandler,
		},
	}
}

type ClientConnHandler interface {
	handling(conn net.Conn)
}

type AbstractConnHandler struct {
	server      *Relatdb
	nextHandler ClientConnHandler
}

func (self *AbstractConnHandler) handling(conn net.Conn) {
	panic("Unrealized method")
}

type HandshakeHandler struct {
	AbstractConnHandler
}

func (self *HandshakeHandler) handling(conn net.Conn) {
	handshakePacket := &protocol.HandshakePacket{
		ProtocolVersion:     common.Protocol_Version,
		ServerVersion:       []byte(common.Relatdb_Version),
		ConnectionId:        1,
		AuthPluginDataPart1: utils.RandomBytes(8),
		ServerCapabilities:  self.server.getServerCapabilities(),
		ServerCharsetIndex:  33,
		ServerStatus:        2,
		AuthPluginDataPart2: utils.RandomBytes(12),
	}
	handshakePacket.SendDataPacket(conn)
	self.nextHandler.handling(conn)
}

type DataReceiveHandler struct {
	AbstractConnHandler
}

func (self *DataReceiveHandler) handling(conn net.Conn) {
	reader := bufio.NewReader(conn)
	for {
		lengthBytes := make([]byte, 3)
		_, _ = reader.Read(lengthBytes)
		length := utils.Uint32(lengthBytes)
		if length <= 0 {
			continue
		}
		data := make([]byte, length)
		_, err := reader.Read(data)
		if err != nil {
			fmt.Println("Received data fail")
			continue
		}
		println(string(data))
	}
}
