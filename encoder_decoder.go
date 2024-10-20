package gorel

import (
	"encoding/binary"
	"unsafe"
)

var reservedSizeForByteSlice = uint(unsafe.Sizeof(uint16(0)))

func BytesNeededForEncodingAByteSlice(buffer []byte) int {
	return int(reservedSizeForByteSlice) + len(buffer)
}

func EncodeByteSlice(source []byte, destination []byte, destinationStartingOffset uint) uint {
	binary.LittleEndian.PutUint16(destination[destinationStartingOffset:], uint16(len(source)))
	copy(destination[destinationStartingOffset+reservedSizeForByteSlice:], source)

	return reservedSizeForByteSlice + uint(len(source))
}

func DecodeByteSlice(source []byte, fromOffset uint16) []byte {
	byteSliceLength := binary.LittleEndian.Uint16(source[fromOffset:])
	endOffset := fromOffset + uint16(reservedSizeForByteSlice) + byteSliceLength
	return source[fromOffset+uint16(reservedSizeForByteSlice) : endOffset]
}

func EncodeUint16(source uint16, destination []byte, destinationStartingOffset uint) {
	binary.LittleEndian.PutUint16(destination[destinationStartingOffset:], source)
}

func DecodeUint16(source []byte, fromOffset uint16) uint16 {
	return binary.LittleEndian.Uint16(source[fromOffset:])
}

func EncodeUint32(source uint32, destination []byte, destinationStartingOffset uint) {
	binary.LittleEndian.PutUint32(destination[destinationStartingOffset:], source)
}

func DecodeUint32(source []byte, fromOffset uint16) uint32 {
	return binary.LittleEndian.Uint32(source[fromOffset:])
}

func EncodeUint64(source uint64, destination []byte, destinationStartingOffset uint) {
	binary.LittleEndian.PutUint64(destination[destinationStartingOffset:], source)
}

func DecodeUint64(source []byte, fromOffset uint16) uint64 {
	return binary.LittleEndian.Uint64(source[fromOffset:])
}
