package server

import (
	"Relatdb/common"
	"Relatdb/meta"
	"Relatdb/utils"
	"bytes"
	"encoding/binary"
)

var (
	SERVER_OK           = []byte{7, 0, 0, 1, 0, 0, 0, 2, 0, 0, 0}
	SERVER_AUTH_OK      = []byte{7, 0, 0, 2, 0, 0, 0, 2, 0, 0, 0}
	OK_HEADER      byte = 0x00
	ERR_HEADER     byte = 0xff
	EOFHeader      byte = 0xfe
	CATALOG_VAL         = "def"
)

type DataPacket interface {
	GetPacketId() byte
	GetPacketBytes() []byte
}

type AbstractDataPacket struct {
	PacketId byte
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

func lengthEncodedInt(buf *bytes.Buffer, n uint64) {
	switch {
	case n < 251:
		buf.WriteByte(byte(n))
	case n <= 0xffff:
		buf.Write([]byte{252, byte(n), byte(n >> 8)})
	case n <= 0xffffff:
		buf.Write([]byte{253, byte(n), byte(n >> 8), byte(n >> 16)})
	case n <= 0xffffffffffffffff:
		buf.Write([]byte{254, byte(n), byte(n >> 8), byte(n >> 16), byte(n >> 24), byte(n >> 32), byte(n >> 40), byte(n >> 48), byte(n >> 56)})
	}
}

func lengthEncodedBytes(buf *bytes.Buffer, bytes []byte) {
	n := uint64(len(bytes))
	switch {
	case n < 251:
		buf.WriteByte(byte(n))
	case n <= 0xffff:
		buf.WriteByte(252)
	case n <= 0xffffff:
		buf.WriteByte(253)
	case n <= 0xffffffffffffffff:
		buf.WriteByte(254)
	}
	buf.Write(bytes)
}

func lengthEncodedString(buf *bytes.Buffer, s string) {
	lengthEncodedBytes(buf, []byte(s))
}

func getIntLength(n uint64) uint {
	switch {
	case n < 251:
		return 1
	case n <= 0xffff:
		return 3
	case n <= 0xffffff:
		return 4
	default:
		return 9
	}
}

type OkPacket struct {
	AbstractDataPacket
	OkHeader     byte
	AffectedRows uint64 //影响行数
	InsertId     uint64 //插入ID
	ServerStatus uint16
	WarningCount uint16
}

func (self *OkPacket) GetPacketBytes() []byte {
	var buf bytes.Buffer
	buf.WriteByte(self.PacketId)
	buf.WriteByte(self.OkHeader)
	lengthEncodedInt(&buf, self.AffectedRows)
	lengthEncodedInt(&buf, self.AffectedRows)
	binary.Write(&buf, binary.LittleEndian, self.ServerStatus)
	binary.Write(&buf, binary.LittleEndian, self.WarningCount)
	bytes := buf.Bytes()
	return append(getDataLengthBytes(uint32(len(bytes))-1), bytes...)
}

type ErrorPacket struct {
	AbstractDataPacket
	ErrorHeader    byte
	ErrorCode      uint16 //错误代码
	SqlStateMarker byte   //SQL状态标识符
	SqlState       []byte //SQL状态
	Message        []byte //错误消息内容
}

func (self *ErrorPacket) GetPacketBytes() []byte {
	var buf bytes.Buffer
	buf.WriteByte(self.PacketId)
	buf.WriteByte(self.ErrorHeader)
	binary.Write(&buf, binary.LittleEndian, self.ErrorCode)
	buf.WriteByte(self.SqlStateMarker)
	buf.Write(self.SqlState)
	buf.Write(self.Message)
	bytes := buf.Bytes()
	return append(getDataLengthBytes(uint32(len(bytes))-1), bytes...)
}

type EofPacket struct {
	AbstractDataPacket
	EofHeader    byte
	WarningCount uint16
	ServerStatus uint16
}

func NewEofPacket(packetId byte, warningCount uint16, serverStatus uint16) *EofPacket {
	eofPacket := &EofPacket{
		EofHeader:    EOFHeader,
		WarningCount: warningCount,
		ServerStatus: serverStatus,
	}
	eofPacket.PacketId = packetId
	return eofPacket
}

func (self *EofPacket) GetPacketBytes() []byte {
	var buf bytes.Buffer
	buf.WriteByte(self.PacketId)
	buf.WriteByte(self.EofHeader)
	binary.Write(&buf, binary.LittleEndian, self.WarningCount)
	binary.Write(&buf, binary.LittleEndian, self.ServerStatus)
	bytes := buf.Bytes()
	return append(getDataLengthBytes(uint32(len(bytes))-1), bytes...)
}

type ColumnPacket struct {
	AbstractDataPacket
	Catalog    string
	Database   string
	Table      string
	OrgTable   string
	Name       string
	OrgName    string
	Charset    uint
	Length     uint64
	Type       uint64
	Flag       uint
	Decimals   byte
	Definition []byte
}

func (self *ColumnPacket) GetPacketBytes() []byte {
	var buf bytes.Buffer
	buf.WriteByte(self.PacketId)
	lengthEncodedString(&buf, self.Catalog)
	if self.Database == "" {
		buf.WriteByte(0)
	} else {
		lengthEncodedString(&buf, self.Database)
	}
	if self.Table == "" {
		buf.WriteByte(0)
	} else {
		lengthEncodedString(&buf, self.Table)
	}
	if self.OrgTable == "" {
		buf.WriteByte(0)
	} else {
		lengthEncodedString(&buf, self.OrgTable)
	}
	if self.Name == "" {
		buf.WriteByte(0)
	} else {
		lengthEncodedString(&buf, self.Name)
	}
	if self.OrgName == "" {
		buf.WriteByte(0)
	} else {
		lengthEncodedString(&buf, self.OrgName)
	}
	buf.WriteByte(0x0c)
	binary.Write(&buf, binary.LittleEndian, self.Charset)
	binary.Write(&buf, binary.LittleEndian, self.Length)
	buf.WriteByte(byte(self.Type))
	binary.Write(&buf, binary.LittleEndian, self.Flag)
	buf.WriteByte(self.Decimals)
	filler := []byte{0, 0}
	buf.Write(filler)
	if len(self.Definition) > 0 {
		lengthEncodedBytes(&buf, self.Definition)
	}
	bytes := buf.Bytes()
	return append(getDataLengthBytes(uint32(len(bytes))-1), bytes...)
}

type RowPacket struct {
	AbstractDataPacket
	Values [][]byte
}

func (self *RowPacket) GetPacketBytes() []byte {
	var buf bytes.Buffer
	buf.WriteByte(self.PacketId)
	nullMark := byte(251)
	for _, value := range self.Values {
		length := uint64(len(value))
		if value == nil || length == 0 {
			buf.WriteByte(nullMark)
		} else {
			lengthEncodedInt(&buf, length)
			buf.Write(value)
		}
	}
	bytes := buf.Bytes()
	return append(getDataLengthBytes(uint32(len(bytes))-1), bytes...)
}

type TablePacket struct {
	AbstractDataPacket
	ColumnPackets    []*ColumnPacket
	ColumnsEofPacket *EofPacket
	RowPackets       []*RowPacket
	RowsEofPacket    *EofPacket
}

func NewTablePacket(columns []meta.Value, rows [][]meta.Value) *TablePacket {
	packetId := byte(2)
	columnPackets := make([]*ColumnPacket, len(columns))
	for i, column := range columns {
		columnPacket := &ColumnPacket{
			Catalog: CATALOG_VAL,
			Name:    column.ToString(),
			Type:    common.FIELD_TYPE_VAR_STRING,
		}
		columnPacket.PacketId = packetId
		columnPackets[i] = columnPacket
		packetId++
	}
	columnsEofPacket := NewEofPacket(packetId, 0, 2)
	packetId++
	rowPackets := make([]*RowPacket, len(rows))
	for i, row := range rows {
		values := make([][]byte, len(row))
		for j, value := range row {
			values[j] = value.ToValueBytes()
		}
		rowPacket := &RowPacket{
			Values: values,
		}
		rowPacket.PacketId = packetId
		rowPackets[i] = rowPacket
		packetId++
	}
	rowsEofPacket := NewEofPacket(packetId, 0, 2)
	selectPacket := &TablePacket{
		ColumnPackets:    columnPackets,
		ColumnsEofPacket: columnsEofPacket,
		RowPackets:       rowPackets,
		RowsEofPacket:    rowsEofPacket,
	}
	selectPacket.PacketId = 1
	return selectPacket
}

func (self *TablePacket) GetPacketBytes() []byte {
	var buf bytes.Buffer
	columnCount := uint64(len(self.ColumnPackets))
	buf.Write(getDataLengthBytes(uint32(getIntLength(columnCount))))
	buf.WriteByte(self.PacketId)
	lengthEncodedInt(&buf, columnCount)

	for _, columnPacket := range self.ColumnPackets {
		bytes := columnPacket.GetPacketBytes()
		buf.Write(bytes)
	}
	buf.Write(self.ColumnsEofPacket.GetPacketBytes())

	for _, rowPacket := range self.RowPackets {
		bytes := rowPacket.GetPacketBytes()
		buf.Write(bytes)
	}
	buf.Write(self.RowsEofPacket.GetPacketBytes())

	bytes := buf.Bytes()
	return bytes
}
