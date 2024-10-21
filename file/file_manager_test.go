package file

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

const blockSize = 4096

func TestWriteAPageUsingBlockFileManager(t *testing.T) {
	fileManager, err := NewBlockFileManager(".", blockSize)
	assert.Nil(t, err)

	defer func() {
		fileManager.Close()
		_ = os.Remove(t.Name())
	}()

	pageBuilder := NewPageBuilder(blockSize)
	pageBuilder.addBytes([]byte("RockDB is an LSM-based storage engine"))

	fileName := t.Name()
	blockId := NewBlockId(fileName, 3)

	err = fileManager.Write(blockId, pageBuilder.build())
	assert.Nil(t, err)
}

func TestWriteAPageAtBlockZeroAndThenReadItUsingBlockFileManager(t *testing.T) {
	fileManager, err := NewBlockFileManager(".", blockSize)
	assert.Nil(t, err)

	defer func() {
		fileManager.Close()
		_ = os.Remove(t.Name())
	}()

	pageBuilder := NewPageBuilder(blockSize)
	pageBuilder.addUint32(500)
	pageBuilder.addString("RockDB is an LSM-based storage engine")

	fileName := t.Name()
	blockId := NewBlockId(fileName, 0)

	err = fileManager.Write(blockId, pageBuilder.build())
	assert.Nil(t, err)

	readPage, err := fileManager.Read(blockId)
	assert.Nil(t, err)

	assert.Equal(t, uint32(500), readPage.getUint32(0))
	assert.Equal(t, "RockDB is an LSM-based storage engine", readPage.getString(1))
}

func TestWriteAPageAtBlockHigherThanZeroAndThenReadItUsingBlockFileManager(t *testing.T) {
	fileManager, err := NewBlockFileManager(".", blockSize)
	assert.Nil(t, err)

	defer func() {
		fileManager.Close()
		_ = os.Remove(t.Name())
	}()

	pageBuilder := NewPageBuilder(blockSize)
	pageBuilder.addUint32(500)
	pageBuilder.addString("PebbleDB is an LSM-based storage engine")

	fileName := t.Name()
	blockId := NewBlockId(fileName, 10)

	err = fileManager.Write(blockId, pageBuilder.build())
	assert.Nil(t, err)

	readPage, err := fileManager.Read(blockId)
	assert.Nil(t, err)

	assert.Equal(t, uint32(500), readPage.getUint32(0))
	assert.Equal(t, "PebbleDB is an LSM-based storage engine", readPage.getString(1))
}
