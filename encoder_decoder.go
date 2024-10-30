package gorel

import (
	"encoding/binary"
	"unsafe"
)

var (
	reservedSizeForByteSlice = uint(unsafe.Sizeof(uint16(0)))
	uint8Size                = uint(unsafe.Sizeof(uint8(0)))
	uint16Size               = uint(unsafe.Sizeof(uint16(0)))
	uint32Size               = uint(unsafe.Sizeof(uint32(0)))
	uint64Size               = uint(unsafe.Sizeof(uint64(0)))
)

type BytesNeededForEncoding = uint
type EndOffset = uint16

func BytesNeededForEncodingAByteSlice(buffer []byte) int {
	return int(reservedSizeForByteSlice) + len(buffer)
}

func EncodeByteSlice(source []byte, destination []byte, destinationStartingOffset uint) BytesNeededForEncoding {
	binary.LittleEndian.PutUint16(destination[destinationStartingOffset:], uint16(len(source)))
	copy(destination[destinationStartingOffset+reservedSizeForByteSlice:], source)

	return reservedSizeForByteSlice + uint(len(source))
}

func DecodeByteSlice(source []byte, fromOffset uint16) ([]byte, EndOffset) {
	byteSliceLength := binary.LittleEndian.Uint16(source[fromOffset:])
	endOffset := fromOffset + uint16(reservedSizeForByteSlice) + byteSliceLength
	return source[fromOffset+uint16(reservedSizeForByteSlice) : endOffset], endOffset
}

func EncodeUint8(source uint8, destination []byte, destinationStartingOffset uint) BytesNeededForEncoding {
	destination[destinationStartingOffset] = source
	return uint8Size
}

func DecodeUint8(source []byte, fromOffset uint16) (uint8, EndOffset) {
	return source[fromOffset], fromOffset + EndOffset(uint8Size)
}

func EncodeUint16(source uint16, destination []byte, destinationStartingOffset uint) BytesNeededForEncoding {
	binary.LittleEndian.PutUint16(destination[destinationStartingOffset:], source)
	return uint16Size
}

func DecodeUint16(source []byte, fromOffset uint16) (uint16, EndOffset) {
	return binary.LittleEndian.Uint16(source[fromOffset:]), fromOffset + EndOffset(uint16Size)
}

func EncodeUint32(source uint32, destination []byte, destinationStartingOffset uint) BytesNeededForEncoding {
	binary.LittleEndian.PutUint32(destination[destinationStartingOffset:], source)
	return uint32Size
}

func DecodeUint32(source []byte, fromOffset uint16) (uint32, EndOffset) {
	return binary.LittleEndian.Uint32(source[fromOffset:]), fromOffset + EndOffset(uint32Size)
}

func EncodeUint64(source uint64, destination []byte, destinationStartingOffset uint) BytesNeededForEncoding {
	binary.LittleEndian.PutUint64(destination[destinationStartingOffset:], source)
	return uint64Size
}

func DecodeUint64(source []byte, fromOffset uint16) (uint64, EndOffset) {
	return binary.LittleEndian.Uint64(source[fromOffset:]), fromOffset + EndOffset(uint64Size)
}
