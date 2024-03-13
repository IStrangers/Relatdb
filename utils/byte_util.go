package utils

import (
	"encoding/binary"
)

func Uint16(bytes []byte) uint16 {
	return binary.LittleEndian.Uint16(bytes)
}

func Uint16ToBytes(u uint16) []byte {
	bytes := make([]byte, 2)
	binary.LittleEndian.PutUint16(bytes, u)
	return bytes
}

func Uint32(bytes []byte) uint32 {
	return binary.LittleEndian.Uint32(bytes)
}

func Uint32ToBytes(u uint32) []byte {
	bytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(bytes, u)
	return bytes
}

func Uint64(bytes []byte) uint64 {
	return binary.LittleEndian.Uint64(bytes)
}

func Uint64ToBytes(u uint64) []byte {
	bytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(bytes, u)
	return bytes
}
