package log

import (
	"encoding/binary"
	"gorel"
	"gorel/file"
	"unsafe"
)

var reservedSizeForNumberOfOffsets = int(unsafe.Sizeof(uint16(0)))

type Page struct {
	buffer          []byte
	startingOffsets *file.StartingOffsets
}

type BackwardRecordIterator struct {
	page        Page
	offsetIndex int
}

func NewPage(buffer []byte, startingOffsets *file.StartingOffsets) Page {
	return Page{
		buffer:          buffer,
		startingOffsets: startingOffsets,
	}
}

func DecodePageFrom(buffer []byte) Page {
	numberOfOffsets := binary.LittleEndian.Uint16(buffer[len(buffer)-reservedSizeForNumberOfOffsets:])
	offsetAtWhichEncodedStartingOffsetsAreWritten := binary.LittleEndian.Uint16(buffer[len(buffer)-reservedSizeForNumberOfOffsets-reservedSizeForNumberOfOffsets:])
	startingOffsets := file.DecodeStartingOffsetsFrom(
		buffer[offsetAtWhichEncodedStartingOffsetsAreWritten : int(offsetAtWhichEncodedStartingOffsetsAreWritten)+reservedSizeForNumberOfOffsets*int(numberOfOffsets)],
	)
	return NewPage(buffer, startingOffsets)
}

func (page Page) BackwardIterator() *BackwardRecordIterator {
	return &BackwardRecordIterator{
		page:        page,
		offsetIndex: page.startingOffsets.Length() - 1,
	}
}

func (page Page) getBytesAt(offset uint16) []byte {
	return gorel.DecodeByteSlice(page.buffer, offset)
}

func (iterator *BackwardRecordIterator) IsValid() bool {
	return iterator.offsetIndex >= 0
}

func (iterator *BackwardRecordIterator) Previous() {
	iterator.offsetIndex = iterator.offsetIndex - 1
}

func (iterator *BackwardRecordIterator) Record() []byte {
	recordStartingOffset := iterator.page.startingOffsets.OffsetAtIndex(iterator.offsetIndex)
	return iterator.page.getBytesAt(recordStartingOffset)
}
