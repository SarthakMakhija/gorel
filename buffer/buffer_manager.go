package buffer

import (
	"errors"
	"gorel/file"
	"gorel/log"
)

var NoBufferAvailableForPinningError = errors.New("no buffer available for pinning")

// BufferManager TODO: concurrency + notify other goroutines on unpin
type BufferManager struct {
	bufferPool []*Buffer
	available  uint
}

func NewBufferManager(
	capacity uint,
	fileManager *file.BlockFileManager,
	logManager *log.BlockLogManager,
) *BufferManager {
	bufferPool := make([]*Buffer, capacity)
	for index := uint(0); index < capacity; index++ {
		bufferPool[index] = NewBuffer(fileManager, logManager)
	}
	return &BufferManager{
		bufferPool: bufferPool,
		available:  capacity,
	}
}

func (bufferManager *BufferManager) Pin(blockId file.BlockId) (*Buffer, error) {
	buffer, err := bufferManager.tryPin(blockId)
	if err != nil {
		return nil, err
	}
	if buffer == nil {
		return nil, NoBufferAvailableForPinningError
	}
	return buffer, nil
}

func (bufferManager *BufferManager) Unpin(buffer *Buffer) {
	buffer.unpin()
	if !buffer.isPinned() {
		bufferManager.available += 1
	}
}

func (bufferManager *BufferManager) Available() int {
	return int(bufferManager.available)
}

func (bufferManager *BufferManager) tryPin(blockId file.BlockId) (*Buffer, error) {
	buffer := bufferManager.findAnExistingBuffer(blockId)
	if buffer == nil {
		buffer = bufferManager.chooseUnpinnedBuffer()
		if buffer == nil {
			return nil, nil
		}
		if err := buffer.AssignToBlock(blockId); err != nil {
			return nil, err
		}
	}
	if !buffer.isPinned() {
		bufferManager.available -= 1
	}
	buffer.pin()
	return buffer, nil
}

func (bufferManager *BufferManager) findAnExistingBuffer(blockId file.BlockId) *Buffer {
	for _, buffer := range bufferManager.bufferPool {
		if buffer.blockId == blockId {
			return buffer
		}
	}
	return nil
}

func (bufferManager *BufferManager) chooseUnpinnedBuffer() *Buffer {
	for _, buffer := range bufferManager.bufferPool {
		if !buffer.isPinned() {
			return buffer
		}
	}
	return nil
}
