package log

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
	currentWriteOffset uint
	blockSize          uint
}

func NewPage(blockSize uint) *Page {
	return &Page{
		buffer:             make([]byte, blockSize),
		startingOffsets:    file.NewStartingOffsets(),
		currentWriteOffset: 0,
		blockSize:          blockSize,
	}
}

func (page *Page) DecodeFrom(buffer []byte) {
	numberOfOffsets := binary.LittleEndian.Uint16(buffer[len(buffer)-reservedSizeForNumberOfOffsets:])
	offsetAtWhichEncodedStartingOffsetsAreWritten := binary.LittleEndian.Uint16(buffer[len(buffer)-reservedSizeForNumberOfOffsets-reservedSizeForNumberOfOffsets:])
	startingOffsets := file.DecodeStartingOffsetsFrom(
		buffer[offsetAtWhichEncodedStartingOffsetsAreWritten : int(offsetAtWhichEncodedStartingOffsetsAreWritten)+reservedSizeForNumberOfOffsets*int(numberOfOffsets)],
	)
	page.buffer = buffer
	page.startingOffsets = startingOffsets
	page.updateCurrentWriteOffset()
}

func (page *Page) Add(buffer []byte) bool {
	if page.hasCapacityFor(buffer) {
		numberOfBytesForEncoding := gorel.EncodeByteSlice(buffer, page.buffer, page.currentWriteOffset)
		page.startingOffsets.Append(uint16(page.currentWriteOffset))
		page.moveCurrentWriteOffsetBy(numberOfBytesForEncoding)
		return true
	}
	return false
}

func (page *Page) Finish() {
	resultingBuffer := page.buffer

	encodedStartingOffsets := page.startingOffsets.Encode()
	offsetToWriteTheEncodedStartingOffsets := page.currentWriteOffset
	copy(resultingBuffer[offsetToWriteTheEncodedStartingOffsets:], encodedStartingOffsets)

	binary.LittleEndian.PutUint16(resultingBuffer[len(resultingBuffer)-reservedSizeForNumberOfOffsets:], uint16(page.startingOffsets.Length()))
	binary.LittleEndian.PutUint16(resultingBuffer[len(resultingBuffer)-reservedSizeForNumberOfOffsets-reservedSizeForNumberOfOffsets:], uint16(offsetToWriteTheEncodedStartingOffsets))
}

func (page *Page) Content() []byte {
	return page.buffer
}

func (page *Page) BackwardIterator() *BackwardRecordIterator {
	return &BackwardRecordIterator{
		page:        page,
		offsetIndex: page.startingOffsets.Length() - 1,
	}
}

func (page *Page) getBytesAt(offset uint16) []byte {
	return gorel.DecodeByteSlice(page.buffer, offset)
}

func (page *Page) moveCurrentWriteOffsetBy(offset uint) {
	page.currentWriteOffset += offset
}

func (page *Page) hasCapacityFor(buffer []byte) bool {
	bytesAvailable :=
		len(page.buffer) -
			int(page.currentWriteOffset) -
			page.startingOffsets.SizeUsedInBytes() -
			2*reservedSizeForNumberOfOffsets

	bytesNeeded := gorel.BytesNeededForEncodingAByteSlice(buffer) + page.startingOffsets.SizeInBytesForAnOffset()
	return bytesAvailable >= bytesNeeded
}

func (page *Page) updateCurrentWriteOffset() {
	lastStartingOffset := page.startingOffsets.OffsetAtIndex(page.startingOffsets.Length() - 1)
	_, endOffset := gorel.DecodeByteSliceWithEndOffset(page.buffer, lastStartingOffset)

	page.currentWriteOffset = uint(endOffset)
}
