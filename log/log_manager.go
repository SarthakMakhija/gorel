package log

import (
	"gorel"
	"gorel/file"
)

// BlockLogManager TODO: concurrency + persistence of latestLogSequenceNumber
type BlockLogManager struct {
	fileManager                *file.BlockFileManager
	logFile                    string
	logPage                    *Page
	currentBlockId             file.BlockId
	latestLogSequenceNumber    uint
	lastSavedLogSequenceNumber uint
}

func NewBlockLogManager(fileManager *file.BlockFileManager, logFile string) (*BlockLogManager, error) {
	numberOfBlocks, err := fileManager.NumberOfBlocks(logFile)
	if err != nil {
		return nil, err
	}
	logManager := &BlockLogManager{
		fileManager: fileManager,
		logFile:     logFile,
		logPage:     NewPage(fileManager.BlockSize()),
	}

	var blockId file.BlockId
	if numberOfBlocks == 0 {
		blockId, err = logManager.appendNewBlock()
		if err != nil {
			return nil, err
		}
	} else {
		blockId = file.NewBlockId(logFile, uint(numberOfBlocks-1))
		if err := fileManager.ReadInto(blockId, logManager.logPage); err != nil {
			return nil, err
		}
	}
	logManager.currentBlockId = blockId
	return logManager, nil
}

func (logManager *BlockLogManager) Append(buffer []byte) error {
	couldAdd := logManager.logPage.Add(buffer)
	if !couldAdd {
		if err := logManager.forceFlush(); err != nil {
			return err
		}
		blockId, err := logManager.appendNewBlock()
		if err != nil {
			return err
		}
		logManager.currentBlockId = blockId
		logManager.logPage = NewPage(logManager.fileManager.BlockSize())
		gorel.Assert(logManager.logPage.Add(buffer), "could not add the bytes to the new log page")
	}
	logManager.latestLogSequenceNumber += 1
	return nil
}

func (logManager *BlockLogManager) Flush(logSequenceNumber uint) error {
	if logSequenceNumber >= logManager.lastSavedLogSequenceNumber {
		return logManager.forceFlush()
	}
	return nil
}

func (logManager *BlockLogManager) BackwardIterator() (*BackwardLogIterator, error) {
	if err := logManager.forceFlush(); err != nil {
		return nil, err
	}
	return NewBackwardLogIterator(logManager.fileManager, logManager.currentBlockId)
}

func (logManager *BlockLogManager) appendNewBlock() (file.BlockId, error) {
	return logManager.fileManager.AppendEmptyBlock(logManager.logFile)
}

func (logManager *BlockLogManager) forceFlush() error {
	logManager.logPage.finish()
	if err := logManager.fileManager.Write(logManager.currentBlockId, logManager.logPage); err != nil {
		return err
	}
	logManager.lastSavedLogSequenceNumber = logManager.latestLogSequenceNumber
	return nil
}
