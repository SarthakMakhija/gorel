package file

import (
	"encoding/binary"
	"unsafe"
)

var reservedSizeForAnOffset = int(unsafe.Sizeof(uint16(0)))

type StartingOffsets struct {
	offsets []uint16
}

func NewStartingOffsets() *StartingOffsets {
	return &StartingOffsets{}
}

func DecodeStartingOffsetsFrom(buffer []byte) *StartingOffsets {
	startingOffsets := NewStartingOffsets()
	for len(buffer) > 0 {
		startingOffsets.offsets = append(startingOffsets.offsets, binary.LittleEndian.Uint16(buffer[:]))
		buffer = buffer[reservedSizeForAnOffset:]
	}
	return startingOffsets
}

func (startingOffsets *StartingOffsets) Append(offset uint16) {
	startingOffsets.offsets = append(startingOffsets.offsets, offset)
}

func (startingOffsets *StartingOffsets) Encode() []byte {
	buffer := make([]byte, len(startingOffsets.offsets)*reservedSizeForAnOffset)
	offsetIndex := 0
	for _, offset := range startingOffsets.offsets {
		binary.LittleEndian.PutUint16(buffer[offsetIndex:], offset)
		offsetIndex += reservedSizeForAnOffset
	}
	return buffer
}

func (startingOffsets *StartingOffsets) Length() int {
	return len(startingOffsets.offsets)
}

func (startingOffsets *StartingOffsets) OffsetAtIndex(index int) uint16 {
	return startingOffsets.offsets[index]
}

func (startingOffsets *StartingOffsets) SizeInBytesForAnOffset() int {
	return reservedSizeForAnOffset
}

func (startingOffsets *StartingOffsets) SizeUsedInBytes() int {
	return reservedSizeForAnOffset * len(startingOffsets.offsets)
}

func SizeUsedInBytesFor(numberOfOffsets uint16) int {
	return reservedSizeForAnOffset * int(numberOfOffsets)
}
