package main

import (
	"bufio"
	"fmt"
	"net"
)

func main() {
	ln, err := net.Listen("tcp", ":3306")
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		return
	}
	defer ln.Close()
	fmt.Println("Listening on 127.0.0.1:3306")

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			continue
		}
		go handleClient(conn)
	}
}

func handleClient(conn net.Conn) {
	defer conn.Close()
	handshake := []byte{
		59, 0, 0, 0, 10, 53, 46, 49, 46, 49, 45, 102, 114, 101, 101, 100, 111, 109, 0, 0, 0, 0, 0,
		111, 112, 115, 57, 114, 103, 105, 73, 0, 79, 183, 28, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 106, 81, 69, 116, 50, 115, 51, 77, 77, 84, 103, 117, 0,
	}
	conn.Write(handshake)
	reader := bufio.NewReader(conn)
	for {
		length := readUB3(reader)
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

func readUB3(reader *bufio.Reader) uint32 {
	d1, _ := reader.ReadByte()
	d2, _ := reader.ReadByte()
	d3, _ := reader.ReadByte()
	i := uint32(d1 & 0xff)
	i |= uint32(d2&0xff) << 8
	i |= uint32(d3&0xff) << 16
	return i
}
