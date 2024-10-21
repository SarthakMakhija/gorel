package file

import (
	"encoding/binary"
	"gorel"
	"unsafe"
)

var reservedSizeForNumberOfOffsets = int(unsafe.Sizeof(uint16(0)))

type PageBuilder struct {
	buffer             []byte
	startingOffsets    *StartingOffsets
	types              *Types
	currentWriteOffset uint
	blockSize          uint
}

func NewPageBuilder(blockSize uint) *PageBuilder {
	return &PageBuilder{
		buffer:             make([]byte, blockSize),
		startingOffsets:    NewStartingOffsets(),
		types:              newTypes(),
		currentWriteOffset: 0,
		blockSize:          blockSize,
	}
}

func (builder *PageBuilder) addUint8(value uint8) {
	builder.addField(
		func() gorel.BytesNeededForEncoding {
			builder.buffer[builder.currentWriteOffset] = value
			return uint8Size
		},
		typeUint8,
	)
}

func (builder *PageBuilder) addUint16(value uint16) {
	builder.addField(
		func() gorel.BytesNeededForEncoding {
			return gorel.EncodeUint16(value, builder.buffer, builder.currentWriteOffset)
		},
		typeUint16,
	)
}

func (builder *PageBuilder) addUint32(value uint32) {
	builder.addField(
		func() gorel.BytesNeededForEncoding {
			return gorel.EncodeUint32(value, builder.buffer, builder.currentWriteOffset)
		},
		typeUint32,
	)
}

func (builder *PageBuilder) addUint64(value uint64) {
	builder.addField(
		func() gorel.BytesNeededForEncoding {
			return gorel.EncodeUint64(value, builder.buffer, builder.currentWriteOffset)
		},
		typeUint64,
	)
}

func (builder *PageBuilder) addBytes(buffer []byte) {
	builder.addField(
		func() gorel.BytesNeededForEncoding {
			return gorel.EncodeByteSlice(buffer, builder.buffer, builder.currentWriteOffset)
		},
		typeByteSlice,
	)
}

func (builder *PageBuilder) addString(str string) {
	builder.addField(
		func() gorel.BytesNeededForEncoding {
			return gorel.EncodeByteSlice([]byte(str), builder.buffer, builder.currentWriteOffset)
		},
		typeString,
	)
}

func (builder *PageBuilder) build() Page1 {
	resultingBuffer := builder.buffer

	encodedStartingOffsets := builder.startingOffsets.Encode()
	encodedTypeDescription := builder.types.encode()

	offsetToWriteTheEncodedStartingOffsets := builder.currentWriteOffset
	copy(resultingBuffer[offsetToWriteTheEncodedStartingOffsets:], encodedStartingOffsets)

	offsetToWriteTypeDescription := offsetToWriteTheEncodedStartingOffsets + uint(len(encodedStartingOffsets))
	copy(resultingBuffer[offsetToWriteTypeDescription:], encodedTypeDescription)

	binary.LittleEndian.PutUint16(resultingBuffer[len(resultingBuffer)-reservedSizeForNumberOfOffsets:], uint16(builder.startingOffsets.Length()))
	binary.LittleEndian.PutUint16(resultingBuffer[len(resultingBuffer)-reservedSizeForNumberOfOffsets-reservedSizeForNumberOfOffsets:], uint16(offsetToWriteTheEncodedStartingOffsets))

	return Page1{
		buffer:          resultingBuffer,
		startingOffsets: builder.startingOffsets,
		types:           builder.types,
	}
}

func (builder *PageBuilder) addField(encodeFn func() gorel.BytesNeededForEncoding, typeDescription uint8) {
	bytesNeededForEncoding := encodeFn()
	builder.startingOffsets.Append(uint16(builder.currentWriteOffset))
	builder.types.addTypeDescription(typeDescription)
	builder.moveCurrentWriteOffsetBy(bytesNeededForEncoding)
}

func (builder *PageBuilder) moveCurrentWriteOffsetBy(offset uint) {
	builder.currentWriteOffset += offset
}
