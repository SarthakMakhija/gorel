package log

import "gorel/file"

type BackwardLogIterator struct {
	fileManager     *file.BlockFileManager
	logPageIterator *BackwardRecordIterator
	currentBlockId  file.BlockId
}

func NewBackwardLogIterator(fileManager *file.BlockFileManager, currentBlockId file.BlockId) (*BackwardLogIterator, error) {
	iterator := &BackwardLogIterator{
		fileManager:    fileManager,
		currentBlockId: currentBlockId,
	}
	page := NewPage(fileManager.BlockSize())
	if err := iterator.readBlockInto(currentBlockId, page); err != nil {
		return nil, err
	}
	iterator.logPageIterator = page.BackwardIterator()
	return iterator, nil
}

func (iterator *BackwardLogIterator) IsValid() bool {
	return iterator.logPageIterator.IsValid()
}

func (iterator *BackwardLogIterator) Previous() error {
	if iterator.logPageIterator.IsValid() {
		iterator.logPageIterator.Previous()
		return nil
	}
	if iterator.currentBlockId.BlockNumber() > 0 {
		iterator.currentBlockId = iterator.currentBlockId.Previous()
		page := NewPage(iterator.fileManager.BlockSize())
		if err := iterator.readBlockInto(iterator.currentBlockId, page); err != nil {
			return err
		}
		iterator.logPageIterator = page.BackwardIterator()
	}
	return nil
}

func (iterator *BackwardLogIterator) Record() []byte {
	return iterator.logPageIterator.Record()
}

func (iterator *BackwardLogIterator) readBlockInto(blockId file.BlockId, page *Page) error {
	return iterator.fileManager.ReadInto(blockId, page)
}
