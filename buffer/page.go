package buffer

import (
	"encoding/binary"
	"gorel"
	"gorel/file"
	"unsafe"
)

var reservedSizeForNumberOfOffsets = int(unsafe.Sizeof(uint16(0)))

// Page TODO: Deletion of value(s) in page? handling holes in page?
type Page struct {
	buffer             []byte
	startingOffsets    *file.StartingOffsets
	types              *Types
	currentWriteOffset uint
}

func NewPage(blockSize uint) *Page {
	return &Page{
		buffer:             make([]byte, blockSize),
		startingOffsets:    file.NewStartingOffsets(),
		types:              NewTypes(),
		currentWriteOffset: 0,
	}
}

func (page *Page) DecodeFrom(buffer []byte) {
	numberOfOffsets := binary.LittleEndian.Uint16(buffer[len(buffer)-reservedSizeForNumberOfOffsets:])
	numberOfTypeDescriptions := numberOfOffsets

	if numberOfOffsets == 0 {
		page.buffer = buffer
		page.startingOffsets = file.NewStartingOffsets()
		page.types = NewTypes()
		page.currentWriteOffset = 0
		return
	}
	offsetAtWhichEncodedStartingOffsetsAreWritten := len(buffer) - reservedSizeForNumberOfOffsets - file.SizeUsedInBytesFor(numberOfOffsets)
	startingOffsets := file.DecodeStartingOffsetsFrom(
		buffer[offsetAtWhichEncodedStartingOffsetsAreWritten : offsetAtWhichEncodedStartingOffsetsAreWritten+reservedSizeForNumberOfOffsets*int(numberOfOffsets)],
	)

	offsetAtWhichEncodedTypeDescriptionsAreWritten := offsetAtWhichEncodedStartingOffsetsAreWritten - SizeUsedInBytes(numberOfTypeDescriptions)
	types := DecodeTypesFrom(buffer[offsetAtWhichEncodedTypeDescriptionsAreWritten : offsetAtWhichEncodedTypeDescriptionsAreWritten+ReservedSizeForAType*int(numberOfTypeDescriptions)])

	page.buffer = buffer
	page.startingOffsets = startingOffsets
	page.types = types
	page.updateCurrentWriteOffset()
}

// AddUint8 TODO: validate capacity before adding, for all the methods.
func (page *Page) AddUint8(value uint8) {
	page.addField(
		func() gorel.BytesNeededForEncoding {
			return gorel.EncodeUint8(value, page.buffer, page.currentWriteOffset)
		},
		TypeUint8,
	)
}

func (page *Page) MutateUint8(index int, value uint8) {
	page.mutateField(index, TypeUint8, func(destinationOffset uint) gorel.BytesNeededForEncoding {
		return gorel.EncodeUint8(value, page.buffer, destinationOffset)
	})
}

func (page *Page) AddUint16(value uint16) {
	page.addField(
		func() gorel.BytesNeededForEncoding {
			return gorel.EncodeUint16(value, page.buffer, page.currentWriteOffset)
		},
		TypeUint16,
	)
}

func (page *Page) MutateUint16(index int, value uint16) {
	page.mutateField(index, TypeUint16, func(destinationOffset uint) gorel.BytesNeededForEncoding {
		return gorel.EncodeUint16(value, page.buffer, destinationOffset)
	})
}

func (page *Page) AddUint32(value uint32) {
	page.addField(
		func() gorel.BytesNeededForEncoding {
			return gorel.EncodeUint32(value, page.buffer, page.currentWriteOffset)
		},
		TypeUint32,
	)
}

func (page *Page) MutateUint32(index int, value uint32) {
	page.mutateField(index, TypeUint32, func(destinationOffset uint) gorel.BytesNeededForEncoding {
		return gorel.EncodeUint32(value, page.buffer, destinationOffset)
	})
}

func (page *Page) AddUint64(value uint64) {
	page.addField(
		func() gorel.BytesNeededForEncoding {
			return gorel.EncodeUint64(value, page.buffer, page.currentWriteOffset)
		},
		TypeUint64,
	)
}

func (page *Page) MutateUint64(index int, value uint64) {
	page.mutateField(index, TypeUint64, func(destinationOffset uint) gorel.BytesNeededForEncoding {
		return gorel.EncodeUint64(value, page.buffer, destinationOffset)
	})
}

func (page *Page) AddBytes(buffer []byte) {
	page.addField(
		func() gorel.BytesNeededForEncoding {
			return gorel.EncodeByteSlice(buffer, page.buffer, page.currentWriteOffset)
		},
		TypeByteSlice,
	)
}

// MutateBytes TODO: what if the value does not fit?
func (page *Page) MutateBytes(index int, value []byte) {
	page.mutateField(index, TypeByteSlice, func(destinationOffset uint) gorel.BytesNeededForEncoding {
		return gorel.EncodeByteSlice(value, page.buffer, destinationOffset)
	})
}

func (page *Page) AddString(str string) {
	page.addField(
		func() gorel.BytesNeededForEncoding {
			return gorel.EncodeByteSlice([]byte(str), page.buffer, page.currentWriteOffset)
		},
		TypeString,
	)
}

func (page *Page) MutateString(index int, value string) {
	page.mutateField(index, TypeString, func(destinationOffset uint) gorel.BytesNeededForEncoding {
		return gorel.EncodeByteSlice([]byte(value), page.buffer, destinationOffset)
	})
}

func (page *Page) finish() {
	resultingBuffer := page.buffer

	encodedStartingOffsets := page.startingOffsets.Encode()
	encodedTypeDescription := page.types.Encode()

	offsetToWriteTheEncodedStartingOffsets := len(resultingBuffer) - reservedSizeForNumberOfOffsets - page.startingOffsets.SizeUsedInBytes()
	copy(resultingBuffer[offsetToWriteTheEncodedStartingOffsets:], encodedStartingOffsets)

	offsetToWriteTypeDescription := offsetToWriteTheEncodedStartingOffsets - page.types.SizeUsedInBytes()
	copy(resultingBuffer[offsetToWriteTypeDescription:], encodedTypeDescription)

	binary.LittleEndian.PutUint16(resultingBuffer[len(resultingBuffer)-reservedSizeForNumberOfOffsets:], uint16(page.startingOffsets.Length()))
}

func (page *Page) Content() []byte {
	return page.buffer
}

func (page *Page) GetUint8(index int) uint8 {
	page.assertFieldAt(index, TypeUint8)
	decoded, _ := gorel.DecodeUint8(page.buffer, page.startingOffsets.OffsetAtIndex(index))
	return decoded
}

func (page *Page) GetUint16(index int) uint16 {
	page.assertFieldAt(index, TypeUint16)
	decoded, _ := gorel.DecodeUint16(page.buffer, page.startingOffsets.OffsetAtIndex(index))
	return decoded
}

func (page *Page) GetUint32(index int) uint32 {
	page.assertFieldAt(index, TypeUint32)
	decoded, _ := gorel.DecodeUint32(page.buffer, page.startingOffsets.OffsetAtIndex(index))
	return decoded
}

func (page *Page) GetUint64(index int) uint64 {
	page.assertFieldAt(index, TypeUint64)
	decoded, _ := gorel.DecodeUint64(page.buffer, page.startingOffsets.OffsetAtIndex(index))
	return decoded
}

func (page *Page) GetString(index int) string {
	page.assertFieldAt(index, TypeString)
	decoded, _ := gorel.DecodeByteSlice(page.buffer, page.startingOffsets.OffsetAtIndex(index))
	return string(decoded)
}

func (page *Page) GetBytes(index int) []byte {
	page.assertFieldAt(index, TypeByteSlice)
	decoded, _ := gorel.DecodeByteSlice(page.buffer, page.startingOffsets.OffsetAtIndex(index))
	return decoded
}

func (page *Page) assertIndexInBounds(index int) {
	gorel.Assert(
		index < page.startingOffsets.Length(),
		"index out of bounds, index = %d, available startingOffsets = %d",
		index,
		page.startingOffsets.Length(),
	)
}

func (page *Page) assertTypeDescriptionMatch(expectedTypeDescription, actualTypeDescription TypeDescription) {
	gorel.Assert(
		expectedTypeDescription.Equals(actualTypeDescription),
		"type description mismatch, expected type %s actual type %s",
		expectedTypeDescription.AsString(),
		actualTypeDescription.AsString(),
	)
}

func (page *Page) addField(encodeFn func() gorel.BytesNeededForEncoding, typeDescription TypeDescription) {
	bytesNeededForEncoding := encodeFn()
	page.startingOffsets.Append(uint16(page.currentWriteOffset))
	page.types.AddTypeDescription(typeDescription)
	page.moveCurrentWriteOffsetBy(bytesNeededForEncoding)
}

func (page *Page) mutateField(index int, typeDescription TypeDescription, encodeFn func(destinationOffset uint) gorel.BytesNeededForEncoding) {
	page.assertFieldAt(index, typeDescription)
	encodeFn(uint(page.startingOffsets.OffsetAtIndex(index)))
}

func (page *Page) assertFieldAt(index int, typeDescription TypeDescription) {
	page.assertIndexInBounds(index)
	page.assertTypeDescriptionMatch(typeDescription, page.types.GetTypeAt(index))
}

func (page *Page) moveCurrentWriteOffsetBy(offset uint) {
	page.currentWriteOffset += offset
}

func (page *Page) updateCurrentWriteOffset() {
	lastTypeDescription := page.types.GetTypeAt(page.types.Length() - 1)
	endOffset := lastTypeDescription.EndOffsetPostDecode(page.buffer, page.startingOffsets.OffsetAtIndex(page.startingOffsets.Length()-1))
	page.currentWriteOffset = uint(endOffset)
}
