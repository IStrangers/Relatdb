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

func Uint16(bytes []byte) uint16 {
	return binary.LittleEndian.Uint16(FillZero(bytes, 2))
}

func Uint16ToBytes(u uint16) []byte {
	bytes := make([]byte, 2)
	binary.LittleEndian.PutUint16(bytes, u)
	return bytes
}

func Uint32(bytes []byte) uint32 {
	return binary.LittleEndian.Uint32(FillZero(bytes, 4))
}

func Uint32ToBytes(u uint32) []byte {
	bytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(bytes, u)
	return bytes
}

func Uint64(bytes []byte) uint64 {
	return binary.LittleEndian.Uint64(FillZero(bytes, 8))
}

func Uint64ToBytes(u uint64) []byte {
	bytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(bytes, u)
	return bytes
}
