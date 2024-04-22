package server

import (
	"Relatdb/utils"
	"bytes"
	"encoding/binary"
)

type DataPacket interface {
	GetPacketId() byte
	GetPacketBytes() []byte
}

type AbstractDataPacket struct {
	PacketSize uint32
	PacketId   byte
}

func (self *AbstractDataPacket) GetPacketId() byte {
	return self.PacketId
}

func getDataLengthBytes(dataLength uint32) []byte {
	bytes := utils.Uint32ToBytes(dataLength, true)[1:]
	return []byte{bytes[2], bytes[1], bytes[0]}
}

/*
握手包
*/
type HandshakePacket struct {
	AbstractDataPacket
	ProtocolVersion     byte   //协议版本（1个字节）
	ServerVersion       []byte //数据库版本（n个字节,结束补0）
	ConnectionId        uint32 //连接ID（4个字节）
	AuthPluginDataPart1 []byte //认证插件随机数1（8个字节）
	ServerCapabilities  uint32 //数据库支持的功能（2个字节 + 填充字节）
	ServerCharsetIndex  byte   //使用的字符集（1个字节）
	ServerStatus        uint16 //数据库状态（2个字节）
	AuthPluginDataPart2 []byte //认证插件随机数2（12位）
	AuthPluginName      []byte //日志插件名称
}

func (self *HandshakePacket) GetPacketBytes() []byte {
	serverCapabilitiesFiller := []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	var buf bytes.Buffer
	buf.WriteByte(0)
	buf.WriteByte(self.ProtocolVersion)
	buf.Write(self.ServerVersion)
	buf.WriteByte(0)
	binary.Write(&buf, binary.LittleEndian, self.ConnectionId)
	buf.Write(self.AuthPluginDataPart1)
	buf.WriteByte(0)
	binary.Write(&buf, binary.LittleEndian, uint16(self.ServerCapabilities))
	buf.WriteByte(self.ServerCharsetIndex)
	binary.Write(&buf, binary.LittleEndian, self.ServerStatus)
	binary.Write(&buf, binary.LittleEndian, uint16(self.ServerCapabilities>>16))
	buf.WriteByte(byte(len(self.AuthPluginDataPart1) + len(self.AuthPluginDataPart2) + 1))
	buf.Write(serverCapabilitiesFiller)
	buf.Write(self.AuthPluginDataPart2)
	buf.WriteByte(0)
	buf.Write(self.AuthPluginName)
	buf.WriteByte(0)
	bytes := buf.Bytes()
	return append(getDataLengthBytes(uint32(len(bytes))-1), bytes...)
}

/*
二进制数据包
*/
type BinaryPacket struct {
	AbstractDataPacket
	Data []byte //数据
}

func (self *BinaryPacket) GetPacketBytes() []byte {
	return self.Data
}

type AuthPacket struct {
	AbstractDataPacket
	ClientFlags    uint32 //客户端功能
	MaxPacketSize  uint32 //最大包大小
	CharsetIndex   byte   //字符集
	Extra          []byte //控制信息
	UserName       string //用户名
	Password       []byte //密码
	DataBase       string //数据库
	AuthPluginName string //认证插件
}

func (self *AuthPacket) GetPacketBytes() []byte {
	return nil
}

func LoadAuthPacket(packet *BinaryPacket) *AuthPacket {
	authPacket := &AuthPacket{}
	bytesReader := utils.NewBytesReader(packet.Data)
	authPacket.PacketSize = packet.PacketSize
	authPacket.PacketId = packet.PacketId
	authPacket.ClientFlags = bytesReader.ReadLittleEndianUint32()
	authPacket.MaxPacketSize = bytesReader.ReadLittleEndianUint32()
	authPacket.CharsetIndex = bytesReader.ReadByte()
	offset := bytesReader.Offset
	length := readLength(bytesReader)
	fillerLength := uint64(23)
	if length > 0 && length < fillerLength {
		authPacket.Extra = bytesReader.ReadBytes(length)
	}
	bytesReader.Offset = offset + fillerLength
	authPacket.UserName = string(bytesReader.ReadToZero())
	length = readLength(bytesReader)
	if length > 0 {
		authPacket.Password = bytesReader.ReadBytes(length)
	} else {
		authPacket.Password = []byte{}
	}
	if authPacket.ClientFlags&CLIENT_CONNECT_WITH_DB != 0 {
		authPacket.DataBase = string(bytesReader.ReadToZero())
	}
	if authPacket.ClientFlags&CLIENT_PLUGIN_AUTH != 0 {
		authPacket.AuthPluginName = string(bytesReader.ReadToZero())
	}
	return authPacket
}

func readLength(bytesReader *utils.BytesReader) uint64 {
	length := bytesReader.ReadByte()
	switch length {
	case 251:
		return 0
	case 252:
		return uint64(bytesReader.ReadLittleEndianUint16())
	case 253:
		return uint64(utils.Uint32(bytesReader.ReadBytes(3), false))
	case 254:
		return bytesReader.ReadLittleEndianUint64()
	default:
		return uint64(length)
	}
}

type ErrorPacket struct {
	AbstractDataPacket
	FieldCount     byte   //包中的字段个数
	ErrorCode      uint16 //错误代码
	SqlStateMarker byte   //SQL状态标识符
	SqlState       []byte //SQL状态
	Message        []byte //错误消息内容
}

func (self *ErrorPacket) GetPacketBytes() []byte {
	var buf bytes.Buffer
	buf.WriteByte(self.PacketId)
	buf.WriteByte(self.FieldCount)
	binary.Write(&buf, binary.LittleEndian, self.ErrorCode)
	buf.WriteByte(self.SqlStateMarker)
	buf.Write(self.SqlState)
	buf.Write(self.Message)
	bytes := buf.Bytes()
	return append(getDataLengthBytes(uint32(len(bytes))-1), bytes...)
}
