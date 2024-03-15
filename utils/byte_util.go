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
