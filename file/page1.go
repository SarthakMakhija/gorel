package file

import (
	"encoding/binary"
	"gorel"
)

type Page1 struct {
	buffer          []byte
	startingOffsets *StartingOffsets
	types           *Types
}

func DecodePageFrom(buffer []byte) Page1 {
	numberOfOffsets := binary.LittleEndian.Uint16(buffer[len(buffer)-reservedSizeForNumberOfOffsets:])
	numberOfTypeDescriptions := numberOfOffsets

	offsetAtWhichEncodedStartingOffsetsAreWritten := binary.LittleEndian.Uint16(buffer[len(buffer)-reservedSizeForNumberOfOffsets-reservedSizeForNumberOfOffsets:])
	startingOffsets := DecodeStartingOffsetsFrom(
		buffer[offsetAtWhichEncodedStartingOffsetsAreWritten : int(offsetAtWhichEncodedStartingOffsetsAreWritten)+reservedSizeForNumberOfOffsets*int(numberOfOffsets)],
	)

	offsetAtWhichEncodedTypeDescriptionsAreWritten := int(offsetAtWhichEncodedStartingOffsetsAreWritten) + reservedSizeForNumberOfOffsets*int(numberOfOffsets)
	types := decodeTypesFrom(buffer[offsetAtWhichEncodedTypeDescriptionsAreWritten : offsetAtWhichEncodedTypeDescriptionsAreWritten+reservedSizeForAType*int(numberOfTypeDescriptions)])

	return Page1{
		buffer:          buffer,
		startingOffsets: startingOffsets,
		types:           types,
	}
}

func (page Page1) getUint8(index int) uint8 {
	return page.buffer[page.startingOffsets.OffsetAtIndex(index)]
}

func (page Page1) getUint16(index int) uint16 {
	return gorel.DecodeUint16(page.buffer, page.startingOffsets.OffsetAtIndex(index))
}

func (page Page1) getUint32(index int) uint32 {
	return gorel.DecodeUint32(page.buffer, page.startingOffsets.OffsetAtIndex(index))
}

func (page Page1) getUint64(index int) uint64 {
	return gorel.DecodeUint64(page.buffer, page.startingOffsets.OffsetAtIndex(index))
}

func (page Page1) getBytes(index int) []byte {
	return gorel.DecodeByteSlice(page.buffer, page.startingOffsets.OffsetAtIndex(index))
}

func (page Page1) getString(index int) string {
	return string(page.getBytes(index))
}
