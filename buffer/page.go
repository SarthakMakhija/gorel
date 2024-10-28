package buffer

import (
	"encoding/binary"
	"gorel"
	"gorel/file"
)

type Page struct {
	buffer          []byte
	startingOffsets *file.StartingOffsets
	types           *file.Types
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
