package file

import (
	"encoding/binary"
	"gorel"
)

type Page struct {
	buffer          []byte
	startingOffsets *StartingOffsets
	types           *Types
}

func DecodePageFrom(buffer []byte) Page {
	numberOfOffsets := binary.LittleEndian.Uint16(buffer[len(buffer)-reservedSizeForNumberOfOffsets:])
	numberOfTypeDescriptions := numberOfOffsets

	offsetAtWhichEncodedStartingOffsetsAreWritten := binary.LittleEndian.Uint16(buffer[len(buffer)-reservedSizeForNumberOfOffsets-reservedSizeForNumberOfOffsets:])
	startingOffsets := DecodeStartingOffsetsFrom(
		buffer[offsetAtWhichEncodedStartingOffsetsAreWritten : int(offsetAtWhichEncodedStartingOffsetsAreWritten)+reservedSizeForNumberOfOffsets*int(numberOfOffsets)],
	)

	offsetAtWhichEncodedTypeDescriptionsAreWritten := int(offsetAtWhichEncodedStartingOffsetsAreWritten) + reservedSizeForNumberOfOffsets*int(numberOfOffsets)
	types := decodeTypesFrom(buffer[offsetAtWhichEncodedTypeDescriptionsAreWritten : offsetAtWhichEncodedTypeDescriptionsAreWritten+reservedSizeForAType*int(numberOfTypeDescriptions)])

	return Page{
		buffer:          buffer,
		startingOffsets: startingOffsets,
		types:           types,
	}
}

func (page Page) getUint8(index int) uint8 {
	page.assertIndexInBounds(index)
	page.assertTypeDescriptionMatch(typeUint8, page.types.getTypeAt(index))
	return page.buffer[page.startingOffsets.OffsetAtIndex(index)]
}

func (page Page) getUint16(index int) uint16 {
	page.assertIndexInBounds(index)
	page.assertTypeDescriptionMatch(typeUint16, page.types.getTypeAt(index))
	return gorel.DecodeUint16(page.buffer, page.startingOffsets.OffsetAtIndex(index))
}

func (page Page) getUint32(index int) uint32 {
	page.assertIndexInBounds(index)
	page.assertTypeDescriptionMatch(typeUint32, page.types.getTypeAt(index))
	return gorel.DecodeUint32(page.buffer, page.startingOffsets.OffsetAtIndex(index))
}

func (page Page) getUint64(index int) uint64 {
	page.assertIndexInBounds(index)
	page.assertTypeDescriptionMatch(typeUint64, page.types.getTypeAt(index))
	return gorel.DecodeUint64(page.buffer, page.startingOffsets.OffsetAtIndex(index))
}

func (page Page) getString(index int) string {
	page.assertIndexInBounds(index)
	page.assertTypeDescriptionMatch(typeString, page.types.getTypeAt(index))
	return string(gorel.DecodeByteSlice(page.buffer, page.startingOffsets.OffsetAtIndex(index)))
}

func (page Page) getBytes(index int) []byte {
	page.assertIndexInBounds(index)
	page.assertTypeDescriptionMatch(typeByteSlice, page.types.getTypeAt(index))
	return gorel.DecodeByteSlice(page.buffer, page.startingOffsets.OffsetAtIndex(index))
}

func (page Page) assertIndexInBounds(index int) {
	gorel.Assert(
		index < page.startingOffsets.Length(),
		"index out of bounds, index = %d, startingOffsets = %d",
		index,
		page.startingOffsets.Length(),
	)
}

func (page Page) assertTypeDescriptionMatch(expectedTypeDescription, actualTypeDescription Type) {
	gorel.Assert(
		expectedTypeDescription.equals(actualTypeDescription),
		"type description mismatch, expected type %s actual type %s",
		expectedTypeDescription.asString(),
		actualTypeDescription.asString(),
	)
}
