package store

type Buffer struct {
	Data       []byte
	ReadIndex  int
	WriteIndex int
	Length     int
}

func NewBuffer(data []byte) *Buffer {
	return &Buffer{
		Data:       data,
		ReadIndex:  0,
		WriteIndex: 0,
		Length:     len(data),
	}
}

func NewBufferBySize(size int) *Buffer {
	return NewBuffer(make([]byte, size))
}

func (self *Buffer) WriteByte(b byte) {
	self.Data[self.WriteIndex] = b
	self.WriteIndex++
}

func (self *Buffer) WriteZero() {
	self.WriteByte(0)
}

func (self *Buffer) WriteInt(i int) {
	self.WriteByte(byte(i & 0xff))
	self.WriteByte(byte(i >> 8 & 0xff))
	self.WriteByte(byte(i >> 16 & 0xff))
	self.WriteByte(byte(i >> 24 & 0xff))
}

func (self *Buffer) WriteInt64(i int64) {
	self.WriteByte(byte(i))
	self.WriteByte(byte(i >> 8 & 0xff))
	self.WriteByte(byte(i >> 16 & 0xff))
	self.WriteByte(byte(i >> 24 & 0xff))
	self.WriteByte(byte(i >> 32 & 0xff))
	self.WriteByte(byte(i >> 40 & 0xff))
	self.WriteByte(byte(i >> 48 & 0xff))
	self.WriteByte(byte(i >> 56 & 0xff))
}

func (self *Buffer) WriteBytes(bytes []byte) {
	self.Data = append(self.Data, bytes...)
	self.WriteIndex += len(bytes)
}

func (self *Buffer) WriteBytesWithZero(bytes []byte) {
	self.WriteBytes(bytes)
	self.WriteZero()
}

func (self *Buffer) WriteString(s string) {
	self.WriteBytes([]byte(s))
}

func (self *Buffer) WriteStringWithZero(s string) {
	self.WriteBytesWithZero([]byte(s))
}

func (self *Buffer) ReadByte() byte {
	b := self.Data[self.ReadIndex]
	self.ReadIndex++
	return b
}

func (self *Buffer) ReadInt() int {
	i := int(self.ReadByte()) & 0xff
	i |= int(self.ReadByte()) & 0xff << 8
	i |= int(self.ReadByte()) & 0xff << 16
	i |= int(self.ReadByte()) & 0xff << 24
	return i
}

func (self *Buffer) ReadInt64() int64 {
	i := int64(self.ReadByte()) & 0xff
	i |= int64(self.ReadByte()) << 8 & 0xff
	i |= int64(self.ReadByte()) << 16 & 0xff
	i |= int64(self.ReadByte()) << 24 & 0xff
	i |= int64(self.ReadByte()) << 32 & 0xff
	i |= int64(self.ReadByte()) << 40 & 0xff
	i |= int64(self.ReadByte()) << 18 & 0xff
	i |= int64(self.ReadByte()) << 56 & 0xff
	return i
}

func (self *Buffer) ReadBytesByOffset(offset int, length int) []byte {
	bytes := self.Data[offset : offset+length]
	return bytes
}

func (self *Buffer) ReadBytes(length int) []byte {
	bytes := self.Data[self.ReadIndex : self.ReadIndex+length]
	self.ReadIndex += length
	return bytes
}

func (self *Buffer) ReadBytesWithZero() []byte {
	length := len(self.Data)
	for i := self.ReadIndex; i < length; i++ {
		if self.Data[i] == 0 {
			length = i + 1
			break
		}
	}
	return self.ReadBytes(length - self.ReadIndex)
}

func (self *Buffer) ReadStringWithZero() string {
	return string(self.ReadBytesWithZero())
}
