package file

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestWriteAPageUsingBlockFileManager(t *testing.T) {
	fileManager, err := NewBlockFileManager(".", blockSize)
	assert.Nil(t, err)

	defer func() {
		fileManager.Close()
		_ = os.Remove(t.Name())
	}()

	page := NewPage(blockSize)
	page.setBytes([]byte("RockDB is an LSM-based storage engine"))

	fileName := t.Name()
	blockId := NewBlockId(fileName, 3)

	err = fileManager.Write(blockId, page)
	assert.Nil(t, err)
}

func TestWriteAPageAtBlockZeroAndThenReadItUsingBlockFileManager(t *testing.T) {
	fileManager, err := NewBlockFileManager(".", blockSize)
	assert.Nil(t, err)

	defer func() {
		fileManager.Close()
		_ = os.Remove(t.Name())
	}()

	page := NewPage(blockSize)
	page.setUint32(uint32(500))
	page.setString("RockDB is an LSM-based storage engine")

	fileName := t.Name()
	blockId := NewBlockId(fileName, 0)

	err = fileManager.Write(blockId, page)
	assert.Nil(t, err)

	readPage := NewPage(blockSize)

	err = fileManager.ReadInto(blockId, readPage)
	assert.Nil(t, err)

	assert.Equal(t, uint32(500), readPage.getUint32())
	assert.Equal(t, "RockDB is an LSM-based storage engine", readPage.getString())
}

func TestWriteAPageAtBlockHigherThanZeroAndThenReadItUsingBlockFileManager(t *testing.T) {
	fileManager, err := NewBlockFileManager(".", blockSize)
	assert.Nil(t, err)

	defer func() {
		fileManager.Close()
		_ = os.Remove(t.Name())
	}()

	page := NewPage(blockSize)
	page.setUint32(uint32(500))
	page.setString("PebbleDB is an LSM-based storage engine")

	fileName := t.Name()
	blockId := NewBlockId(fileName, 10)

	err = fileManager.Write(blockId, page)
	assert.Nil(t, err)

	readPage := NewPage(blockSize)

	err = fileManager.ReadInto(blockId, readPage)
	assert.Nil(t, err)

	assert.Equal(t, uint32(500), readPage.getUint32())
	assert.Equal(t, "PebbleDB is an LSM-based storage engine", readPage.getString())
}
