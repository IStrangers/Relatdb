package utils

import (
	"encoding/binary"
)

func FillZero(bytes []byte, end uint32) []byte {
	for i := uint32(len(bytes)); i < end; i++ {
		bytes = append(bytes, 0)
	}
	return bytes
}

func Uint16(bytes []byte, bigEndian bool) uint16 {
	if bigEndian {
		return binary.BigEndian.Uint16(FillZero(bytes, 2))
	} else {
		return binary.LittleEndian.Uint16(FillZero(bytes, 2))
	}
}

func Uint16ToBytes(u uint16, bigEndian bool) []byte {
	bytes := make([]byte, 2)
	if bigEndian {
		binary.BigEndian.PutUint16(bytes, u)
	} else {
		binary.LittleEndian.PutUint16(bytes, u)
	}
	return bytes
}

func Uint32(bytes []byte, bigEndian bool) uint32 {
	if bigEndian {
		return binary.BigEndian.Uint32(FillZero(bytes, 4))
	} else {
		return binary.LittleEndian.Uint32(FillZero(bytes, 4))
	}
}

func Uint32ToBytes(u uint32, bigEndian bool) []byte {
	bytes := make([]byte, 4)
	if bigEndian {
		binary.BigEndian.PutUint32(bytes, u)
	} else {
		binary.LittleEndian.PutUint32(bytes, u)
	}
	return bytes
}

func Uint64(bytes []byte, bigEndian bool) uint64 {
	if bigEndian {
		return binary.BigEndian.Uint64(FillZero(bytes, 8))
	} else {
		return binary.LittleEndian.Uint64(FillZero(bytes, 8))
	}
}

func Uint64ToBytes(u uint64, bigEndian bool) []byte {
	bytes := make([]byte, 8)
	if bigEndian {
		binary.BigEndian.PutUint64(bytes, u)
	} else {
		binary.LittleEndian.PutUint64(bytes, u)
	}
	return bytes
}

type BytesReader struct {
	Offset uint64
	Data   []byte
	Length uint64
}

func NewBytesReader(data []byte) *BytesReader {
	return &BytesReader{
		Offset: 0,
		Data:   data,
		Length: uint64(len(data)),
	}
}

func (self *BytesReader) ReadByte() byte {
	defer func() {
		self.Offset += 1
	}()
	return self.Data[self.Offset]
}

func (self *BytesReader) ReadBytes(size uint64) []byte {
	defer func() {
		self.Offset += size
	}()
	return self.Data[self.Offset : self.Offset+size]
}

func (self *BytesReader) ReadLittleEndianUint16() uint16 {
	return self.ReadUint16(false)
}

func (self *BytesReader) ReadBigEndianUint16() uint16 {
	return self.ReadUint16(false)
}

func (self *BytesReader) ReadUint16(bigEndian bool) uint16 {
	defer func() {
		self.Offset += 2
	}()
	return Uint16(self.Data[self.Offset:self.Offset+2], bigEndian)
}

func (self *BytesReader) ReadLittleEndianUint32() uint32 {
	return self.ReadUint32(false)
}

func (self *BytesReader) ReadBigEndianUint32() uint32 {
	return self.ReadUint32(false)
}

func (self *BytesReader) ReadUint32(bigEndian bool) uint32 {
	defer func() {
		self.Offset += 4
	}()
	return Uint32(self.Data[self.Offset:self.Offset+4], bigEndian)
}

func (self *BytesReader) ReadLittleEndianUint64() uint64 {
	return self.ReadUint64(false)
}

func (self *BytesReader) ReadBigEndianUint64() uint64 {
	return self.ReadUint64(false)
}

func (self *BytesReader) ReadUint64(bigEndian bool) uint64 {
	defer func() {
		self.Offset += 8
	}()
	return Uint64(self.Data[self.Offset:self.Offset+8], bigEndian)
}

func (self *BytesReader) ReadToZero() []byte {
	offset := self.Offset
	for {
		offset++
		if offset >= self.Length {
			return []byte{}
		}
		if self.Data[offset] == 0 {
			data := self.Data[self.Offset:offset]
			self.Offset = offset + 1
			return data
		}
	}
}

func (self *BytesReader) ReadRemainingBytes() []byte {
	defer func() {
		self.Offset = self.Length
	}()
	return self.Data[self.Offset:]
}
