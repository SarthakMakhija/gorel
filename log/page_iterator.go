package log

type BackwardRecordIterator struct {
	page        *Page
	offsetIndex int
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
