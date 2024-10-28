package buffer

import (
	"github.com/stretchr/testify/assert"
	"gorel/file"
	"os"
	"testing"
)

func TestWriteAPageUsingBlockFileManager(t *testing.T) {
	fileManager, err := file.NewBlockFileManager(".", blockSize)
	assert.Nil(t, err)

	defer func() {
		fileManager.Close()
		_ = os.Remove(t.Name())
	}()

	pageBuilder := NewPageBuilder(blockSize)
	pageBuilder.AddBytes([]byte("RockDB is an LSM-based storage engine"))

	fileName := t.Name()
	blockId := file.NewBlockId(fileName, 3)

	err = fileManager.Write(blockId, pageBuilder.Build())
	assert.Nil(t, err)
}

func TestWriteAPageAtBlockZeroAndThenReadItUsingBlockFileManager(t *testing.T) {
	fileManager, err := file.NewBlockFileManager(".", blockSize)
	assert.Nil(t, err)

	defer func() {
		fileManager.Close()
		_ = os.Remove(t.Name())
	}()

	pageBuilder := NewPageBuilder(blockSize)
	pageBuilder.AddUint32(500)
	pageBuilder.AddString("RockDB is an LSM-based storage engine")

	fileName := t.Name()
	blockId := file.NewBlockId(fileName, 0)

	err = fileManager.Write(blockId, pageBuilder.Build())
	assert.Nil(t, err)

	readPage := &Page{}
	err = fileManager.ReadInto(blockId, readPage)
	assert.Nil(t, err)

	assert.Equal(t, uint32(500), readPage.GetUint32(0))
	assert.Equal(t, "RockDB is an LSM-based storage engine", readPage.GetString(1))
}

func TestWriteAPageAtBlockHigherThanZeroAndThenReadItUsingBlockFileManager(t *testing.T) {
	fileManager, err := file.NewBlockFileManager(".", blockSize)
	assert.Nil(t, err)

	defer func() {
		fileManager.Close()
		_ = os.Remove(t.Name())
	}()

	pageBuilder := NewPageBuilder(blockSize)
	pageBuilder.AddUint32(500)
	pageBuilder.AddString("PebbleDB is an LSM-based storage engine")

	fileName := t.Name()
	blockId := file.NewBlockId(fileName, 10)

	err = fileManager.Write(blockId, pageBuilder.Build())
	assert.Nil(t, err)

	readPage := &Page{}
	err = fileManager.ReadInto(blockId, readPage)
	assert.Nil(t, err)

	assert.Equal(t, uint32(500), readPage.GetUint32(0))
	assert.Equal(t, "PebbleDB is an LSM-based storage engine", readPage.GetString(1))
}
