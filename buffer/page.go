package buffer

import (
	"encoding/binary"
	"gorel"
	"gorel/file"
	"unsafe"
)

var reservedSizeForNumberOfOffsets = int(unsafe.Sizeof(uint16(0)))

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

	offsetAtWhichEncodedStartingOffsetsAreWritten := len(buffer) - reservedSizeForNumberOfOffsets - file.SizeUsedInBytesFor(numberOfOffsets)
	startingOffsets := file.DecodeStartingOffsetsFrom(
		buffer[offsetAtWhichEncodedStartingOffsetsAreWritten : offsetAtWhichEncodedStartingOffsetsAreWritten+reservedSizeForNumberOfOffsets*int(numberOfOffsets)],
	)

	offsetAtWhichEncodedTypeDescriptionsAreWritten := offsetAtWhichEncodedStartingOffsetsAreWritten - file.SizeUsedInBytes(numberOfTypeDescriptions)
	types := file.DecodeTypesFrom(buffer[offsetAtWhichEncodedTypeDescriptionsAreWritten : offsetAtWhichEncodedTypeDescriptionsAreWritten+file.ReservedSizeForAType*int(numberOfTypeDescriptions)])

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
	page.assertIndexInBounds(index)
	page.assertTypeDescriptionMatch(file.TypeUint8, page.types.GetTypeAt(index))
	decoded, _ := gorel.DecodeUint8(page.buffer, page.startingOffsets.OffsetAtIndex(index))
	return decoded
}

func (page *Page) GetUint16(index int) uint16 {
	page.assertIndexInBounds(index)
	page.assertTypeDescriptionMatch(file.TypeUint16, page.types.GetTypeAt(index))
	decoded, _ := gorel.DecodeUint16(page.buffer, page.startingOffsets.OffsetAtIndex(index))
	return decoded
}

func (page *Page) GetUint32(index int) uint32 {
	page.assertIndexInBounds(index)
	page.assertTypeDescriptionMatch(file.TypeUint32, page.types.GetTypeAt(index))
	decoded, _ := gorel.DecodeUint32(page.buffer, page.startingOffsets.OffsetAtIndex(index))
	return decoded
}

func (page *Page) GetUint64(index int) uint64 {
	page.assertIndexInBounds(index)
	page.assertTypeDescriptionMatch(file.TypeUint64, page.types.GetTypeAt(index))
	decoded, _ := gorel.DecodeUint64(page.buffer, page.startingOffsets.OffsetAtIndex(index))
	return decoded
}

func (page *Page) GetString(index int) string {
	page.assertIndexInBounds(index)
	page.assertTypeDescriptionMatch(file.TypeString, page.types.GetTypeAt(index))
	decoded, _ := gorel.DecodeByteSlice(page.buffer, page.startingOffsets.OffsetAtIndex(index))
	return string(decoded)
}

func (page *Page) GetBytes(index int) []byte {
	page.assertIndexInBounds(index)
	page.assertTypeDescriptionMatch(file.TypeByteSlice, page.types.GetTypeAt(index))
	decoded, _ := gorel.DecodeByteSlice(page.buffer, page.startingOffsets.OffsetAtIndex(index))
	return decoded
}

func (page *Page) assertIndexInBounds(index int) {
	gorel.Assert(
		index < page.startingOffsets.Length(),
		"index out of bounds, index = %d, startingOffsets = %d",
		index,
		page.startingOffsets.Length(),
	)
}

func (page *Page) assertTypeDescriptionMatch(expectedTypeDescription, actualTypeDescription file.TypeDescription) {
	gorel.Assert(
		expectedTypeDescription.Equals(actualTypeDescription),
		"type description mismatch, expected type %s actual type %s",
		expectedTypeDescription.AsString(),
		actualTypeDescription.AsString(),
	)
}

func (page *Page) addField(encodeFn func() gorel.BytesNeededForEncoding, typeDescription file.TypeDescription) {
	bytesNeededForEncoding := encodeFn()
	page.startingOffsets.Append(uint16(page.currentWriteOffset))
	page.types.AddTypeDescription(typeDescription)
	page.moveCurrentWriteOffsetBy(bytesNeededForEncoding)
}

func (page *Page) moveCurrentWriteOffsetBy(offset uint) {
	page.currentWriteOffset += offset
}

func (page *Page) updateCurrentWriteOffset() {
	lastTypeDescription := page.types.GetTypeAt(page.types.Length() - 1)
	endOffset := lastTypeDescription.EndOffsetPostDecode(page.buffer, page.startingOffsets.OffsetAtIndex(page.startingOffsets.Length()-1))
	page.currentWriteOffset = uint(endOffset)
}
