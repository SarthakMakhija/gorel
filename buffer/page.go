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

type Page struct {
	buffer             []byte
	startingOffsets    *file.StartingOffsets
	types              *file.Types
	currentWriteOffset uint
}

func NewPage(blockSize uint) *Page {
	return &Page{
		buffer:             make([]byte, blockSize),
		startingOffsets:    file.NewStartingOffsets(),
		types:              file.NewTypes(),
		currentWriteOffset: 0,
	}
}

func (page *Page) DecodeFrom(buffer []byte) {
	numberOfOffsets := binary.LittleEndian.Uint16(buffer[len(buffer)-reservedSizeForNumberOfOffsets:])
	numberOfTypeDescriptions := numberOfOffsets

	offsetAtWhichEncodedStartingOffsetsAreWritten := binary.LittleEndian.Uint16(buffer[len(buffer)-reservedSizeForNumberOfOffsets-reservedSizeForNumberOfOffsets:])
	startingOffsets := file.DecodeStartingOffsetsFrom(
		buffer[offsetAtWhichEncodedStartingOffsetsAreWritten : int(offsetAtWhichEncodedStartingOffsetsAreWritten)+reservedSizeForNumberOfOffsets*int(numberOfOffsets)],
	)

	offsetAtWhichEncodedTypeDescriptionsAreWritten := int(offsetAtWhichEncodedStartingOffsetsAreWritten) + reservedSizeForNumberOfOffsets*int(numberOfOffsets)
	types := file.DecodeTypesFrom(buffer[offsetAtWhichEncodedTypeDescriptionsAreWritten : offsetAtWhichEncodedTypeDescriptionsAreWritten+file.ReservedSizeForAType*int(numberOfTypeDescriptions)])

	page.buffer = buffer
	page.startingOffsets = startingOffsets
	page.types = types
}

// AddUint8 TODO: validate capacity before adding, for all the methods.
func (page *Page) AddUint8(value uint8) {
	page.addField(
		func() gorel.BytesNeededForEncoding {
			page.buffer[page.currentWriteOffset] = value
			return uint8Size
		},
		file.TypeUint8,
	)
}

func (page *Page) AddUint16(value uint16) {
	page.addField(
		func() gorel.BytesNeededForEncoding {
			return gorel.EncodeUint16(value, page.buffer, page.currentWriteOffset)
		},
		file.TypeUint16,
	)
}

func (page *Page) AddUint32(value uint32) {
	page.addField(
		func() gorel.BytesNeededForEncoding {
			return gorel.EncodeUint32(value, page.buffer, page.currentWriteOffset)
		},
		file.TypeUint32,
	)
}

func (page *Page) AddUint64(value uint64) {
	page.addField(
		func() gorel.BytesNeededForEncoding {
			return gorel.EncodeUint64(value, page.buffer, page.currentWriteOffset)
		},
		file.TypeUint64,
	)
}

func (page *Page) AddBytes(buffer []byte) {
	page.addField(
		func() gorel.BytesNeededForEncoding {
			return gorel.EncodeByteSlice(buffer, page.buffer, page.currentWriteOffset)
		},
		file.TypeByteSlice,
	)
}

func (page *Page) AddString(str string) {
	page.addField(
		func() gorel.BytesNeededForEncoding {
			return gorel.EncodeByteSlice([]byte(str), page.buffer, page.currentWriteOffset)
		},
		file.TypeString,
	)
}

func (page *Page) Finish() {
	resultingBuffer := page.buffer

	encodedStartingOffsets := page.startingOffsets.Encode()
	encodedTypeDescription := page.types.Encode()

	offsetToWriteTheEncodedStartingOffsets := page.currentWriteOffset
	copy(resultingBuffer[offsetToWriteTheEncodedStartingOffsets:], encodedStartingOffsets)

	offsetToWriteTypeDescription := offsetToWriteTheEncodedStartingOffsets + uint(len(encodedStartingOffsets))
	copy(resultingBuffer[offsetToWriteTypeDescription:], encodedTypeDescription)

	binary.LittleEndian.PutUint16(resultingBuffer[len(resultingBuffer)-reservedSizeForNumberOfOffsets:], uint16(page.startingOffsets.Length()))
	binary.LittleEndian.PutUint16(resultingBuffer[len(resultingBuffer)-reservedSizeForNumberOfOffsets-reservedSizeForNumberOfOffsets:], uint16(offsetToWriteTheEncodedStartingOffsets))
}

func (page *Page) Content() []byte {
	return page.buffer
}

func (page *Page) GetUint8(index int) uint8 {
	page.assertIndexInBounds(index)
	page.assertTypeDescriptionMatch(file.TypeUint8, page.types.GetTypeAt(index))
	return page.buffer[page.startingOffsets.OffsetAtIndex(index)]
}

func (page *Page) GetUint16(index int) uint16 {
	page.assertIndexInBounds(index)
	page.assertTypeDescriptionMatch(file.TypeUint16, page.types.GetTypeAt(index))
	return gorel.DecodeUint16(page.buffer, page.startingOffsets.OffsetAtIndex(index))
}

func (page *Page) GetUint32(index int) uint32 {
	page.assertIndexInBounds(index)
	page.assertTypeDescriptionMatch(file.TypeUint32, page.types.GetTypeAt(index))
	return gorel.DecodeUint32(page.buffer, page.startingOffsets.OffsetAtIndex(index))
}

func (page *Page) GetUint64(index int) uint64 {
	page.assertIndexInBounds(index)
	page.assertTypeDescriptionMatch(file.TypeUint64, page.types.GetTypeAt(index))
	return gorel.DecodeUint64(page.buffer, page.startingOffsets.OffsetAtIndex(index))
}

func (page *Page) GetString(index int) string {
	page.assertIndexInBounds(index)
	page.assertTypeDescriptionMatch(file.TypeString, page.types.GetTypeAt(index))
	return string(gorel.DecodeByteSlice(page.buffer, page.startingOffsets.OffsetAtIndex(index)))
}

func (page *Page) GetBytes(index int) []byte {
	page.assertIndexInBounds(index)
	page.assertTypeDescriptionMatch(file.TypeByteSlice, page.types.GetTypeAt(index))
	return gorel.DecodeByteSlice(page.buffer, page.startingOffsets.OffsetAtIndex(index))
}

func (page *Page) assertIndexInBounds(index int) {
	gorel.Assert(
		index < page.startingOffsets.Length(),
		"index out of bounds, index = %d, startingOffsets = %d",
		index,
		page.startingOffsets.Length(),
	)
}

func (page *Page) assertTypeDescriptionMatch(expectedTypeDescription, actualTypeDescription file.Type) {
	gorel.Assert(
		expectedTypeDescription.Equals(actualTypeDescription),
		"type description mismatch, expected type %s actual type %s",
		expectedTypeDescription.AsString(),
		actualTypeDescription.AsString(),
	)
}

func (page *Page) addField(encodeFn func() gorel.BytesNeededForEncoding, typeDescription file.Type) {
	bytesNeededForEncoding := encodeFn()
	page.startingOffsets.Append(uint16(page.currentWriteOffset))
	page.types.AddTypeDescription(typeDescription)
	page.moveCurrentWriteOffsetBy(bytesNeededForEncoding)
}

func (page *Page) moveCurrentWriteOffsetBy(offset uint) {
	page.currentWriteOffset += offset
}
