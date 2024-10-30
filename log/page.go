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
	offsetAtWhichEncodedStartingOffsetsAreWritten := len(buffer) - reservedSizeForNumberOfOffsets - file.SizeUsedInBytesFor(numberOfOffsets)
	startingOffsets := file.DecodeStartingOffsetsFrom(
		buffer[offsetAtWhichEncodedStartingOffsetsAreWritten : offsetAtWhichEncodedStartingOffsetsAreWritten+file.SizeUsedInBytesFor(numberOfOffsets)],
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

func (page *Page) finish() {
	resultingBuffer := page.buffer

	encodedStartingOffsets := page.startingOffsets.Encode()
	offsetToWriteTheEncodedStartingOffsets := len(resultingBuffer) - reservedSizeForNumberOfOffsets - page.startingOffsets.SizeUsedInBytes()

	copy(resultingBuffer[offsetToWriteTheEncodedStartingOffsets:], encodedStartingOffsets)
	binary.LittleEndian.PutUint16(resultingBuffer[len(resultingBuffer)-reservedSizeForNumberOfOffsets:], uint16(page.startingOffsets.Length()))
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
	decoded, _ := gorel.DecodeByteSlice(page.buffer, offset)
	return decoded
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
	_, endOffset := gorel.DecodeByteSlice(page.buffer, lastStartingOffset)

	page.currentWriteOffset = uint(endOffset)
}
