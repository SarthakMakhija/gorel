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
