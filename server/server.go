package server

import (
	"fmt"
	"net"
)

type Options struct {
	BindIp   string
	BindPort uint
}

type Server struct {
	version string
	options *Options
	ln      net.Listener

	autoConnId uint64
	connMap    map[uint64]*Connection
}

func CreateServer(options *Options) *Server {
	server := &Server{
		version: SERVER_VERSION,
		options: options,

		autoConnId: 0,
		connMap:    map[uint64]*Connection{},
	}
	return server
}

func (self *Server) Start() {
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
		self.handlingConn(conn)
	}
}

func (self *Server) Stop() {
	err := self.ln.Close()
	if err != nil {
		fmt.Println("Listening stop error:", err.Error())
	}
}

func (self *Server) getServerCapabilities() uint32 {
	flag := uint32(0)
	flag |= CLIENT_LONG_PASSWORD
	flag |= CLIENT_FOUND_ROWS
	flag |= CLIENT_LONG_FLAG
	flag |= CLIENT_CONNECT_WITH_DB
	// flag |=  common.CLIENT_NO_SCHEMA;
	// flag |=  common.CLIENT_COMPRESS;
	flag |= CLIENT_ODBC
	// flag |=  common.CLIENT_LOCAL_FILES;
	flag |= CLIENT_IGNORE_SPACE
	flag |= CLIENT_PROTOCOL_41
	flag |= CLIENT_INTERACTIVE
	// flag |=  common.CLIENT_SSL;
	flag |= CLIENT_IGNORE_SIGPIPE
	flag |= CLIENT_TRANSACTIONS
	flag |= CLIENT_SECURE_CONNECTION
	flag |= CLIENT_PLUGIN_AUTH
	return flag
}

func (self *Server) nextConnId() uint64 {
	self.autoConnId++
	return self.autoConnId
}

func (self *Server) handlingConn(c net.Conn) {
	conn := NewConnection(self, c)

	go self.onConn(conn)
}

func (self *Server) registerConn(conn *Connection) {
	self.connMap[conn.connId] = conn
}

func (self *Server) onConn(conn *Connection) {
	conn.sendHandshakePacket()
	if !conn.authentication() {
		return
	}
	self.registerConn(conn)
	conn.receiveCommandHandler()
}
