package protocol

import (
	"Relatdb/common"
	"Relatdb/utils"
	"bytes"
	"encoding/binary"
)

type DataPacket interface {
	GetPacketId() byte
	GetDataLength() uint32
	GetDataLengthBytes(uint32) []byte
	GetPacketBytes() []byte
}

type AbstractDataPacket struct {
	PacketSize uint32
	PacketId   byte
}

func (self *AbstractDataPacket) GetPacketId() byte {
	return self.PacketId
}

func (self *AbstractDataPacket) GetDataLengthBytes(dataLength uint32) []byte {
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
}

func (self *HandshakePacket) GetDataLength() uint32 {
	return uint32(1 + len(self.ServerVersion) + 1 + 4 + len(self.AuthPluginDataPart1) + 1 + 2 + 1 + 2 + len(self.AuthPluginDataPart2) + 13 + 1)
}

func (self *HandshakePacket) GetPacketBytes() []byte {
	serverCapabilitiesFiller := []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	var buf bytes.Buffer
	buf.Write(self.GetDataLengthBytes(self.GetDataLength()))
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
	buf.Write(serverCapabilitiesFiller)
	buf.Write(self.AuthPluginDataPart2)
	buf.WriteByte(0)
	return buf.Bytes()
}

/*
二进制数据包
*/
type BinaryPacket struct {
	AbstractDataPacket
	Data []byte //数据
}

func (self *BinaryPacket) GetDataLength() uint32 {
	return uint32(len(self.Data))
}

func (self *BinaryPacket) GetPacketBytes() []byte {
	return nil
}

type AuthPacket struct {
	AbstractDataPacket
	ClientFlags   uint32 //客户端功能
	MaxPacketSize uint32 //最大包大小
	CharsetIndex  byte   //字符集
	Extra         []byte //控制信息
	UserName      string //用户名
	Password      []byte //密码
	DataBase      string //数据库
}

func (self *AuthPacket) GetDataLength() uint32 {
	return 0
}

func (self *AuthPacket) GetPacketBytes() []byte {
	return nil
}

func (self *AuthPacket) Load(packet *BinaryPacket) {
	bytesReader := utils.NewBytesReader(packet.Data)
	self.PacketSize = packet.PacketSize
	self.PacketId = packet.PacketId
	self.ClientFlags = bytesReader.ReadLittleEndianUint32()
	self.MaxPacketSize = bytesReader.ReadLittleEndianUint32()
	self.CharsetIndex = bytesReader.ReadByte()
	offset := bytesReader.Offset
	length := readLength(bytesReader)
	fillerLength := uint64(23)
	if length > 0 && length < fillerLength {
		self.Extra = bytesReader.ReadBytes(length)
	}
	bytesReader.Offset = offset + fillerLength
	self.UserName = string(bytesReader.ReadToZero())
	length = readLength(bytesReader)
	if length > 0 {
		self.Password = bytesReader.ReadBytes(length)
	} else {
		self.Password = []byte{}
	}
	if self.ClientFlags&common.CLIENT_CONNECT_WITH_DB != 0 {
		self.DataBase = string(bytesReader.ReadToZero())
	}
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

func (self *ErrorPacket) GetDataLength() uint32 {
	return uint32(1 + 2 + 1 + len(self.SqlState) + len(self.Message))
}

func (self *ErrorPacket) GetPacketBytes() []byte {
	var buf bytes.Buffer
	buf.Write(self.GetDataLengthBytes(self.GetDataLength()))
	buf.WriteByte(self.PacketId)
	buf.WriteByte(self.FieldCount)
	binary.Write(&buf, binary.LittleEndian, self.ErrorCode)
	buf.WriteByte(self.SqlStateMarker)
	buf.Write(self.SqlState)
	buf.Write(self.Message)
	return buf.Bytes()
}
