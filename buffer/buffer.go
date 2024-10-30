package buffer

import (
	"gorel/file"
	"gorel/log"
)

type Buffer struct {
	fileManager       *file.BlockFileManager
	logManager        *log.BlockLogManager
	page              *Page
	blockId           file.BlockId
	pins              int
	transactionNumber int
	logSequenceNumber uint
}

func NewBuffer(fileManager *file.BlockFileManager, logManager *log.BlockLogManager) *Buffer {
	return &Buffer{
		fileManager:       fileManager,
		logManager:        logManager,
		page:              NewPage(fileManager.BlockSize()),
		blockId:           file.MissingBlockId,
		pins:              0,
		transactionNumber: -1,
	}
}

func (buffer *Buffer) Page() *Page {
	return buffer.page
}

func (buffer *Buffer) IsPinned() bool {
	return buffer.pins > 0
}

func (buffer *Buffer) SetModified(transactionNumber int, logSequenceNumber uint) {
	buffer.transactionNumber = transactionNumber
	if logSequenceNumber >= 0 {
		buffer.logSequenceNumber = logSequenceNumber
	}
}

func (buffer *Buffer) AssignToBlock(blockId file.BlockId) error {
	if err := buffer.flush(); err != nil {
		return err
	}
	if err := buffer.fileManager.ReadInto(blockId, buffer.page); err != nil {
		return err
	}
	buffer.blockId = blockId
	buffer.pins = 0
	return nil
}

func (buffer *Buffer) flush() error {
	if buffer.transactionNumber >= 0 {
		if err := buffer.logManager.Flush(buffer.logSequenceNumber); err != nil {
			return err
		}
		buffer.page.finish()
		if err := buffer.fileManager.Write(buffer.blockId, buffer.page); err != nil {
			return err
		}
		buffer.transactionNumber = -1
	}
	return nil
}

func (buffer *Buffer) pin() {
	buffer.pins += 1
}

func (buffer *Buffer) unpin() {
	buffer.pins -= 1
}
