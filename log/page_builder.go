package log

import (
	"encoding/binary"
	"gorel"
	"gorel/file"
)

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

func (builder *PageBuilder) Add(buffer []byte) bool {
	if builder.hasCapacityFor(buffer) {
		numberOfBytesForEncoding := gorel.EncodeByteSlice(buffer, builder.buffer, builder.currentWriteOffset)
		builder.startingOffsets.Append(uint16(builder.currentWriteOffset))
		builder.moveCurrentWriteOffsetBy(numberOfBytesForEncoding)
		return true
	}
	return false
}

func (builder *PageBuilder) Build() Page {
	resultingBuffer := builder.buffer

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

func (builder *PageBuilder) hasCapacityFor(buffer []byte) bool {
	bytesAvailable :=
		len(builder.buffer) -
			int(builder.currentWriteOffset) -
			builder.startingOffsets.SizeUsedInBytes() -
			2*reservedSizeForNumberOfOffsets

	bytesNeeded := gorel.BytesNeededForEncodingAByteSlice(buffer) + builder.startingOffsets.SizeInBytesForAnOffset()
	return bytesAvailable >= bytesNeeded
}
