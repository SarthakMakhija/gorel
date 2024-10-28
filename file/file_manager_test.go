package file

import (
	"github.com/stretchr/testify/assert"
	"gorel"
	"os"
	"testing"
)

type testPage struct {
	buffer             []byte
	currentWriteOffset uint
	blockSize          uint
}

func newTestPage(blockSize uint) *testPage {
	return &testPage{
		buffer:             make([]byte, blockSize),
		currentWriteOffset: 0,
		blockSize:          blockSize,
	}
}

func (page *testPage) DecodeFrom(buffer []byte) {
	page.buffer = buffer
}

func (page *testPage) add(buffer []byte) {
	numberOfBytesForEncoding := gorel.EncodeByteSlice(buffer, page.buffer, page.currentWriteOffset)
	page.currentWriteOffset += numberOfBytesForEncoding
}

func (page *testPage) getBytes(offset uint16) []byte {
	return gorel.DecodeByteSlice(page.buffer, offset)
}

func (page *testPage) Content() []byte {
	return page.buffer
}

const blockSize = 4096

func TestWriteAPageUsingBlockFileManager(t *testing.T) {
	fileManager, err := NewBlockFileManager(".", blockSize)
	assert.Nil(t, err)

	defer func() {
		fileManager.Close()
		_ = os.Remove(t.Name())
	}()

	page := newTestPage(blockSize)
	page.add([]byte("RockDB is an LSM-based storage engine"))

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

	page := newTestPage(blockSize)
	page.add([]byte("RockDB is an LSM-based storage engine"))

	fileName := t.Name()
	blockId := NewBlockId(fileName, 0)

	err = fileManager.Write(blockId, page)
	assert.Nil(t, err)

	readPage := &testPage{}
	err = fileManager.ReadInto(blockId, readPage)
	assert.Nil(t, err)

	assert.Equal(t, "RockDB is an LSM-based storage engine", string(readPage.getBytes(0)))
}

func TestWriteAPageAtBlockHigherThanZeroAndThenReadItUsingBlockFileManager(t *testing.T) {
	fileManager, err := NewBlockFileManager(".", blockSize)
	assert.Nil(t, err)

	defer func() {
		fileManager.Close()
		_ = os.Remove(t.Name())
	}()

	page := newTestPage(blockSize)
	page.add([]byte("PebbleDB is an LSM-based storage engine"))

	fileName := t.Name()
	blockId := NewBlockId(fileName, 10)

	err = fileManager.Write(blockId, page)
	assert.Nil(t, err)

	readPage := &testPage{}
	err = fileManager.ReadInto(blockId, readPage)
	assert.Nil(t, err)

	assert.Equal(t, "PebbleDB is an LSM-based storage engine", string(readPage.getBytes(0)))
}

func TestAppendAnEmptyBlockUsingBlockFileManager(t *testing.T) {
	fileManager, err := NewBlockFileManager(".", blockSize)
	assert.Nil(t, err)

	defer func() {
		fileManager.Close()
		_ = os.Remove(t.Name())
	}()

	fileName := t.Name()
	blockId, err := fileManager.AppendEmptyBlock(fileName)
	assert.Nil(t, err)

	assert.Equal(t, uint(0), blockId.blockNumber)
}

func TestAppendACoupleOfEmptyBlocksUsingBlockFileManager(t *testing.T) {
	fileManager, err := NewBlockFileManager(".", blockSize)
	assert.Nil(t, err)

	defer func() {
		fileManager.Close()
		_ = os.Remove(t.Name())
	}()

	fileName := t.Name()
	blockId, err := fileManager.AppendEmptyBlock(fileName)
	assert.Nil(t, err)

	page := newTestPage(blockSize)
	page.add([]byte("PebbleDB is an LSM-based storage engine"))

	assert.Nil(t, fileManager.Write(blockId, page))

	blockId, err = fileManager.AppendEmptyBlock(fileName)
	assert.Nil(t, err)

	assert.Equal(t, uint(1), blockId.blockNumber)
}
