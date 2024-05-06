package common

type Buffer struct {
	Data       []byte
	ReadIndex  uint
	WriteIndex uint
	Length     uint
}

func NewBuffer(data []byte) *Buffer {
	return &Buffer{
		Data:       data,
		ReadIndex:  0,
		WriteIndex: 0,
		Length:     uint(len(data)),
	}
}

func NewBufferBySize(size uint) *Buffer {
	return NewBuffer(make([]byte, size))
}

func (self *Buffer) Remaining() uint {
	return self.Length - self.ReadIndex
}

func (self *Buffer) WriteByte(b byte) {
	self.Data[self.WriteIndex] = b
	self.WriteIndex++
}

func (self *Buffer) WriteZero() {
	self.WriteByte(0)
}

func (self *Buffer) WriteIntByPos(pos uint, i int) {
	backUp := self.WriteIndex
	self.WriteIndex = pos
	self.WriteInt(i)
	self.WriteIndex = backUp
}

func (self *Buffer) WriteInt(i int) {
	self.WriteByte(byte(i))
	self.WriteByte(byte((i) >> 8))
	self.WriteByte(byte((i) >> 16))
	self.WriteByte(byte((i) >> 24))
}

func (self *Buffer) WriteInt64(i int64) {
	self.WriteByte(byte(i))
	self.WriteByte(byte(i >> 8))
	self.WriteByte(byte(i >> 16))
	self.WriteByte(byte(i >> 24))
	self.WriteByte(byte(i >> 32))
	self.WriteByte(byte(i >> 40))
	self.WriteByte(byte(i >> 48))
	self.WriteByte(byte(i >> 56))
}

func (self *Buffer) WriteBytesByPos(pos uint, bytes []byte) {
	backUp := self.WriteIndex
	self.WriteIndex = pos
	self.WriteBytes(bytes)
	self.WriteIndex = backUp
}

func (self *Buffer) WriteBytes(bytes []byte) {
	for _, b := range bytes {
		self.Data[self.WriteIndex] = b
		self.WriteIndex++
	}
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
	i := int(self.ReadByte())
	i |= int(self.ReadByte()) << 8
	i |= int(self.ReadByte()) << 16
	i |= int(self.ReadByte()) << 24
	return i
}

func (self *Buffer) ReadInt64() int64 {
	i := int64(self.ReadByte())
	i |= int64(self.ReadByte()) << 8
	i |= int64(self.ReadByte()) << 16
	i |= int64(self.ReadByte()) << 24
	i |= int64(self.ReadByte()) << 32
	i |= int64(self.ReadByte()) << 40
	i |= int64(self.ReadByte()) << 18
	i |= int64(self.ReadByte()) << 56
	return i
}

func (self *Buffer) ReadBytesByOffset(offset int, length int) []byte {
	bytes := self.Data[offset : offset+length]
	return bytes
}

func (self *Buffer) ReadBytes(length uint) []byte {
	bytes := self.Data[self.ReadIndex : self.ReadIndex+length]
	self.ReadIndex += length
	return bytes
}

func (self *Buffer) ReadBytesWithZero() []byte {
	length := uint(len(self.Data))
	for i := self.ReadIndex; i < length; i++ {
		if self.Data[i] == 0 {
			length = i
			break
		}
	}
	bytes := self.ReadBytes(length - self.ReadIndex)
	self.ReadIndex += 1
	return bytes
}

func (self *Buffer) ReadStringWithZero() string {
	return string(self.ReadBytesWithZero())
}
