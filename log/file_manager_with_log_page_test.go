package log

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
	pageBuilder.Add([]byte("RockDB is an LSM-based storage engine"))

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
	pageBuilder.Add([]byte("RockDB is an LSM-based storage engine"))

	fileName := t.Name()
	blockId := file.NewBlockId(fileName, 0)

	err = fileManager.Write(blockId, pageBuilder.Build())
	assert.Nil(t, err)

	readPage := &Page{}
	err = fileManager.ReadInto(blockId, readPage)
	assert.Nil(t, err)

	iterator := readPage.BackwardIterator()

	assert.True(t, iterator.IsValid())
	assert.Equal(t, "RockDB is an LSM-based storage engine", string(iterator.Record()))

	iterator.Previous()
	assert.False(t, iterator.IsValid())
}

func TestWriteAPageAtBlockHigherThanZeroAndThenReadItUsingBlockFileManager(t *testing.T) {
	fileManager, err := file.NewBlockFileManager(".", blockSize)
	assert.Nil(t, err)

	defer func() {
		fileManager.Close()
		_ = os.Remove(t.Name())
	}()

	pageBuilder := NewPageBuilder(blockSize)
	pageBuilder.Add([]byte("PebbleDB is an LSM-based storage engine"))

	fileName := t.Name()
	blockId := file.NewBlockId(fileName, 10)

	err = fileManager.Write(blockId, pageBuilder.Build())
	assert.Nil(t, err)

	readPage := &Page{}
	err = fileManager.ReadInto(blockId, readPage)
	assert.Nil(t, err)

	iterator := readPage.BackwardIterator()

	assert.True(t, iterator.IsValid())
	assert.Equal(t, "PebbleDB is an LSM-based storage engine", string(iterator.Record()))

	iterator.Previous()
	assert.False(t, iterator.IsValid())
}
