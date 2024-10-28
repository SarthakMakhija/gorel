package buffer

import (
	"encoding/binary"
	"gorel"
	"gorel/file"
	"unsafe"
)

var (
	reservedSizeForNumberOfOffsets = int(unsafe.Sizeof(uint16(0)))
	uint8Size                      = uint(unsafe.Sizeof(uint8(0)))
)

type PageBuilder struct {
	buffer             []byte
	startingOffsets    *file.StartingOffsets
	types              *file.Types
	currentWriteOffset uint
	blockSize          uint
}

func NewPageBuilder(blockSize uint) *PageBuilder {
	return &PageBuilder{
		buffer:             make([]byte, blockSize),
		startingOffsets:    file.NewStartingOffsets(),
		types:              file.NewTypes(),
		currentWriteOffset: 0,
		blockSize:          blockSize,
	}
}

// AddUint8 TODO: validate capacity before adding, for all the methods.
func (builder *PageBuilder) AddUint8(value uint8) {
	builder.addField(
		func() gorel.BytesNeededForEncoding {
			builder.buffer[builder.currentWriteOffset] = value
			return uint8Size
		},
		file.TypeUint8,
	)
}

func (builder *PageBuilder) AddUint16(value uint16) {
	builder.addField(
		func() gorel.BytesNeededForEncoding {
			return gorel.EncodeUint16(value, builder.buffer, builder.currentWriteOffset)
		},
		file.TypeUint16,
	)
}

func (builder *PageBuilder) AddUint32(value uint32) {
	builder.addField(
		func() gorel.BytesNeededForEncoding {
			return gorel.EncodeUint32(value, builder.buffer, builder.currentWriteOffset)
		},
		file.TypeUint32,
	)
}

func (builder *PageBuilder) AddUint64(value uint64) {
	builder.addField(
		func() gorel.BytesNeededForEncoding {
			return gorel.EncodeUint64(value, builder.buffer, builder.currentWriteOffset)
		},
		file.TypeUint64,
	)
}

func (builder *PageBuilder) AddBytes(buffer []byte) {
	builder.addField(
		func() gorel.BytesNeededForEncoding {
			return gorel.EncodeByteSlice(buffer, builder.buffer, builder.currentWriteOffset)
		},
		file.TypeByteSlice,
	)
}

func (builder *PageBuilder) AddString(str string) {
	builder.addField(
		func() gorel.BytesNeededForEncoding {
			return gorel.EncodeByteSlice([]byte(str), builder.buffer, builder.currentWriteOffset)
		},
		file.TypeString,
	)
}

func (builder *PageBuilder) Build() *Page {
	resultingBuffer := builder.buffer

	encodedStartingOffsets := builder.startingOffsets.Encode()
	encodedTypeDescription := builder.types.Encode()

	offsetToWriteTheEncodedStartingOffsets := builder.currentWriteOffset
	copy(resultingBuffer[offsetToWriteTheEncodedStartingOffsets:], encodedStartingOffsets)

	offsetToWriteTypeDescription := offsetToWriteTheEncodedStartingOffsets + uint(len(encodedStartingOffsets))
	copy(resultingBuffer[offsetToWriteTypeDescription:], encodedTypeDescription)

	binary.LittleEndian.PutUint16(resultingBuffer[len(resultingBuffer)-reservedSizeForNumberOfOffsets:], uint16(builder.startingOffsets.Length()))
	binary.LittleEndian.PutUint16(resultingBuffer[len(resultingBuffer)-reservedSizeForNumberOfOffsets-reservedSizeForNumberOfOffsets:], uint16(offsetToWriteTheEncodedStartingOffsets))

	return &Page{
		buffer:          resultingBuffer,
		startingOffsets: builder.startingOffsets,
		types:           builder.types,
	}
}

func (builder *PageBuilder) addField(encodeFn func() gorel.BytesNeededForEncoding, typeDescription file.Type) {
	bytesNeededForEncoding := encodeFn()
	builder.startingOffsets.Append(uint16(builder.currentWriteOffset))
	builder.types.AddTypeDescription(typeDescription)
	builder.moveCurrentWriteOffsetBy(bytesNeededForEncoding)
}

func (builder *PageBuilder) moveCurrentWriteOffsetBy(offset uint) {
	builder.currentWriteOffset += offset
}
