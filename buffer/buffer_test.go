package buffer

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"gorel/file"
	"gorel/log"
	"os"
	"testing"
)

func TestBufferIsPinned(t *testing.T) {
	fileManager, err := file.NewBlockFileManager(".", blockSize)
	assert.Nil(t, err)

	defer func() {
		fileManager.Close()
		_ = os.Remove(t.Name())
	}()

	fileName := t.Name()

	logManager, err := log.NewBlockLogManager(fileManager, fileName)
	assert.Nil(t, err)

	buffer := NewBuffer(fileManager, logManager)
	buffer.pin()

	assert.True(t, buffer.IsPinned())
}

func TestBufferIsNotPinned(t *testing.T) {
	fileManager, err := file.NewBlockFileManager(".", blockSize)
	assert.Nil(t, err)

	defer func() {
		fileManager.Close()
		_ = os.Remove(t.Name())
	}()

	fileName := t.Name()

	logManager, err := log.NewBlockLogManager(fileManager, fileName)
	assert.Nil(t, err)

	buffer := NewBuffer(fileManager, logManager)
	assert.False(t, buffer.IsPinned())
}

func TestTotalPinsForABuffer(t *testing.T) {
	fileManager, err := file.NewBlockFileManager(".", blockSize)
	assert.Nil(t, err)

	defer func() {
		fileManager.Close()
		_ = os.Remove(t.Name())
	}()

	fileName := t.Name()

	logManager, err := log.NewBlockLogManager(fileManager, fileName)
	assert.Nil(t, err)

	buffer := NewBuffer(fileManager, logManager)
	buffer.pin()
	buffer.pin()

	assert.Equal(t, 2, buffer.pins)
}

func TestTotalPinsForABufferAfterUnpinningIt(t *testing.T) {
	fileManager, err := file.NewBlockFileManager(".", blockSize)
	assert.Nil(t, err)

	defer func() {
		fileManager.Close()
		_ = os.Remove(t.Name())
	}()

	fileName := t.Name()

	logManager, err := log.NewBlockLogManager(fileManager, fileName)
	assert.Nil(t, err)

	buffer := NewBuffer(fileManager, logManager)
	buffer.pin()
	buffer.pin()
	buffer.unpin()

	assert.Equal(t, 1, buffer.pins)
}

func TestAssignBlockToABuffer(t *testing.T) {
	fileManager, err := file.NewBlockFileManager(".", blockSize)
	assert.Nil(t, err)

	bufferFileName := t.Name()
	logFileName := fmt.Sprintf("%v_%v", t.Name(), "log")

	defer func() {
		fileManager.Close()
		_ = os.Remove(bufferFileName)
		_ = os.Remove(logFileName)
	}()

	page := NewPage(fileManager.BlockSize())
	page.AddUint32(32)
	page.AddString("BoltDB is a B+Tree based storage engine")
	page.finish()

	assert.Nil(t, fileManager.Write(file.NewBlockId(bufferFileName, 0), page))

	logManager, err := log.NewBlockLogManager(fileManager, logFileName)
	assert.Nil(t, err)

	buffer := NewBuffer(fileManager, logManager)
	assert.Nil(t, buffer.AssignToBlock(file.NewBlockId(bufferFileName, 0)))

	bufferPage := buffer.page
	assert.Equal(t, uint32(32), bufferPage.GetUint32(0))
	assert.Equal(t, "BoltDB is a B+Tree based storage engine", bufferPage.GetString(1))
}

func TestFlushABuffer(t *testing.T) {
	fileManager, err := file.NewBlockFileManager(".", blockSize)
	assert.Nil(t, err)

	bufferFileName := t.Name()
	logFileName := fmt.Sprintf("%v_%v", t.Name(), "log")

	defer func() {
		fileManager.Close()
		_ = os.Remove(bufferFileName)
		_ = os.Remove(logFileName)
	}()

	logManager, err := log.NewBlockLogManager(fileManager, logFileName)
	assert.Nil(t, err)

	_, err = fileManager.AppendEmptyBlock(bufferFileName)
	assert.Nil(t, err)

	buffer := NewBuffer(fileManager, logManager)
	assert.Nil(t, buffer.AssignToBlock(file.NewBlockId(bufferFileName, 0)))

	bufferPage := buffer.page
	bufferPage.AddUint32(32)
	bufferPage.AddString("BoltDB is a B+Tree based storage engine")

	anyTransactionNumber := 10
	anyLogSequenceNumber := uint(10)
	buffer.SetModified(anyTransactionNumber, anyLogSequenceNumber)
	assert.Nil(t, buffer.flush())

	assert.Nil(t, buffer.AssignToBlock(file.NewBlockId(bufferFileName, 0)))
	reAssignedBufferPage := buffer.page

	assert.Equal(t, uint32(32), reAssignedBufferPage.GetUint32(0))
	assert.Equal(t, "BoltDB is a B+Tree based storage engine", reAssignedBufferPage.GetString(1))
}
