package buffer

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"gorel/file"
	"gorel/log"
	"os"
	"testing"
)

func TestFailsToPinABuffer(t *testing.T) {
	fileManager, err := file.NewBlockFileManager(".", blockSize)
	assert.Nil(t, err)

	fileName := t.Name()
	logFileName := fmt.Sprintf("%v_%v", t.Name(), "log")

	defer func() {
		fileManager.Close()
		_ = os.Remove(fileName)
		_ = os.Remove(logFileName)
	}()

	logManager, err := log.NewBlockLogManager(fileManager, logFileName)
	assert.Nil(t, err)

	bufferManager := NewBufferManager(1, fileManager, logManager)
	bufferManager.bufferPool[0].pin()

	_, err = bufferManager.Pin(file.NewBlockId(fileName, 0))
	assert.EqualError(t, err, NoBufferAvailableForPinningError.Error())
}

func TestAvailableBuffers(t *testing.T) {
	fileManager, err := file.NewBlockFileManager(".", blockSize)
	assert.Nil(t, err)

	fileName := t.Name()
	logFileName := fmt.Sprintf("%v_%v", t.Name(), "log")

	defer func() {
		fileManager.Close()
		_ = os.Remove(fileName)
		_ = os.Remove(logFileName)
	}()

	logManager, err := log.NewBlockLogManager(fileManager, logFileName)
	assert.Nil(t, err)

	bufferManager := NewBufferManager(1, fileManager, logManager)
	assert.Equal(t, 1, bufferManager.Available())
}

func TestPinABuffer(t *testing.T) {
	fileManager, err := file.NewBlockFileManager(".", blockSize)
	assert.Nil(t, err)

	fileName := t.Name()
	logFileName := fmt.Sprintf("%v_%v", t.Name(), "log")

	defer func() {
		fileManager.Close()
		_ = os.Remove(fileName)
		_ = os.Remove(logFileName)
	}()

	blockId, err := fileManager.AppendEmptyBlock(fileName)
	assert.Nil(t, err)

	logManager, err := log.NewBlockLogManager(fileManager, logFileName)
	assert.Nil(t, err)

	bufferManager := NewBufferManager(1, fileManager, logManager)

	buffer, err := bufferManager.Pin(blockId)
	assert.Nil(t, err)

	page := buffer.Page()
	page.AddString("RocksDB is an LSM based storage engine")
	page.AddUint32(32)

	anyTransactionNumber := 10
	anyLogSequenceNumber := uint(10)
	buffer.SetModified(anyTransactionNumber, anyLogSequenceNumber)
	assert.Nil(t, buffer.flush())

	buffer, err = bufferManager.Pin(blockId)
	assert.Nil(t, err)

	assert.Equal(t, "RocksDB is an LSM based storage engine", buffer.page.GetString(0))
	assert.Equal(t, uint32(32), buffer.page.GetUint32(1))
}
