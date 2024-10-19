package log

import (
	"encoding/binary"
	"gorel/file"
	"unsafe"
)

var reservedSizeForByteSlice = uint(unsafe.Sizeof(uint16(0)))

type PageBuilder struct {
	buffer             []byte
	startingOffsets    *file.StartingOffsets
	currentWriteOffset uint
	blockSize          uint
}

func NewPageBuilder(blockSize uint) *PageBuilder {
	return &PageBuilder{
		buffer:             make([]byte, blockSize),
		startingOffsets:    file.NewStartingOffsets(),
		currentWriteOffset: 0,
		blockSize:          blockSize,
	}
}

func (builder *PageBuilder) Add(buffer []byte) {
	//TODO: check if the page has the capacity to fit the incoming record.
	binary.LittleEndian.PutUint16(builder.buffer[builder.currentWriteOffset:], uint16(len(buffer)))
	copy(builder.buffer[builder.currentWriteOffset+reservedSizeForByteSlice:], buffer)

	builder.startingOffsets.Append(uint16(builder.currentWriteOffset))
	builder.moveCurrentWriteOffsetBy(uint(len(buffer)) + reservedSizeForByteSlice)
}

func (builder *PageBuilder) Build() Page {
	resultingBuffer := make([]byte, builder.blockSize)
	copy(resultingBuffer, builder.buffer[:builder.currentWriteOffset])

	encodedStartingOffsets := builder.startingOffsets.Encode()
	offsetToWriteTheEncodedStartingOffsets := builder.currentWriteOffset
	copy(resultingBuffer[offsetToWriteTheEncodedStartingOffsets:], encodedStartingOffsets)

	binary.LittleEndian.PutUint16(resultingBuffer[len(resultingBuffer)-reservedSizeForNumberOfOffsets:], uint16(builder.startingOffsets.Length()))
	binary.LittleEndian.PutUint16(resultingBuffer[len(resultingBuffer)-reservedSizeForNumberOfOffsets-reservedSizeForNumberOfOffsets:], uint16(offsetToWriteTheEncodedStartingOffsets))
	return NewPage(resultingBuffer, builder.startingOffsets)
}

func (builder *PageBuilder) moveCurrentWriteOffsetBy(offset uint) {
	builder.currentWriteOffset += offset
}
