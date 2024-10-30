package buffer

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

const blockSize = 4096

func TestCreateAPageWithASingleField(t *testing.T) {
	page := NewPage(blockSize)
	page.AddUint8(10)
	page.finish()

	assert.Equal(t, uint8(10), page.GetUint8(0))
}

func TestCreateAPageWithFewFields(t *testing.T) {
	page := NewPage(blockSize)
	page.AddUint16(16)
	page.AddUint8(8)
	page.AddUint32(32)
	page.AddUint16(100)
	page.finish()

	assert.Equal(t, uint16(16), page.GetUint16(0))
	assert.Equal(t, uint8(8), page.GetUint8(1))
	assert.Equal(t, uint32(32), page.GetUint32(2))
	assert.Equal(t, uint16(100), page.GetUint16(3))
}

func TestCreateAPageWithByteSliceAndString(t *testing.T) {
	page := NewPage(blockSize)
	page.AddBytes([]byte("RocksDB is an LSM-based key/value storage engine"))
	page.AddString("PebbleDB is an LSM-based key/value storage engine")
	page.finish()

	assert.Equal(t, []byte("RocksDB is an LSM-based key/value storage engine"), page.GetBytes(0))
	assert.Equal(t, "PebbleDB is an LSM-based key/value storage engine", page.GetString(1))
}

func TestAttemptToGetTheValueAtAnIndexGreaterThanTheNumberOfAvailableFields(t *testing.T) {
	page := NewPage(blockSize)
	page.AddString("PebbleDB is an LSM-based key/value storage engine")
	page.finish()

	decodedPage := &Page{}
	decodedPage.DecodeFrom(page.buffer)

	assert.Panics(t, func() {
		decodedPage.GetString(1)
	})
}

func TestAttemptToGetTheValueWithMismatchedTypeDescription(t *testing.T) {
	page := NewPage(blockSize)
	page.AddString("PebbleDB is an LSM-based key/value storage engine")
	page.finish()

	decodedPage := &Page{}
	decodedPage.DecodeFrom(page.buffer)

	assert.Panics(t, func() {
		decodedPage.GetUint8(0)
	})
}

func TestDecodeAPageWithASingleField(t *testing.T) {
	page := NewPage(blockSize)
	page.AddString("PebbleDB is an LSM-based key/value storage engine")
	page.finish()

	decodedPage := &Page{}
	decodedPage.DecodeFrom(page.buffer)

	assert.Equal(t, "PebbleDB is an LSM-based key/value storage engine", decodedPage.GetString(0))
}

func TestDecodeAPageWithFewFields(t *testing.T) {
	page := NewPage(blockSize)
	page.AddString("PebbleDB is an LSM-based key/value storage engine")
	page.AddUint16(16)
	page.AddUint16(160)
	page.AddUint64(64)
	page.finish()

	decodedPage := &Page{}
	decodedPage.DecodeFrom(page.buffer)

	assert.Equal(t, uint64(64), decodedPage.GetUint64(3))
	assert.Equal(t, uint16(160), decodedPage.GetUint16(2))
	assert.Equal(t, uint16(16), decodedPage.GetUint16(1))
	assert.Equal(t, "PebbleDB is an LSM-based key/value storage engine", decodedPage.GetString(0))
}

func TestAddFewFieldsInPageDecodeThePageAndAddFieldsInTheDecodedPageToSimulateLoadingPageFromDiskAndChangingIt(t *testing.T) {
	page := NewPage(blockSize)
	page.AddString("PebbleDB is an LSM-based key/value storage engine")
	page.AddUint32(32)
	page.AddUint16(16)
	page.AddUint64(64)
	page.finish()

	decodedPage := &Page{}
	decodedPage.DecodeFrom(page.buffer)
	decodedPage.AddString("RocksDB is an LSM-based key/value storage engine")

	assert.Equal(t, "PebbleDB is an LSM-based key/value storage engine", decodedPage.GetString(0))
	assert.Equal(t, uint32(32), decodedPage.GetUint32(1))
	assert.Equal(t, uint16(16), decodedPage.GetUint16(2))
	assert.Equal(t, uint64(64), decodedPage.GetUint64(3))
	assert.Equal(t, "RocksDB is an LSM-based key/value storage engine", decodedPage.GetString(4))
}

func TestMutateAnUint8InPage(t *testing.T) {
	page := NewPage(blockSize)
	page.AddUint8(10)
	page.MutateUint8(0, 20)

	assert.Equal(t, uint8(20), page.GetUint8(0))
}

func TestMutateAnUint16InPage(t *testing.T) {
	page := NewPage(blockSize)
	page.AddUint16(16)
	page.MutateUint16(0, 20)

	assert.Equal(t, uint16(20), page.GetUint16(0))
}

func TestMutateAnUint32InPage(t *testing.T) {
	page := NewPage(blockSize)
	page.AddUint32(32)
	page.MutateUint32(0, 64)

	assert.Equal(t, uint32(64), page.GetUint32(0))
}

func TestMutateAnUint64InPage(t *testing.T) {
	page := NewPage(blockSize)
	page.AddUint64(64)
	page.MutateUint64(0, 640)

	assert.Equal(t, uint64(640), page.GetUint64(0))
}

func TestMutateAByteSliceInPage(t *testing.T) {
	page := NewPage(blockSize)
	page.AddBytes([]byte("Bolt-DB"))
	page.MutateBytes(0, []byte("RocksDB"))

	assert.Equal(t, []byte("RocksDB"), page.GetBytes(0))
}

func TestMutateAStringInPage(t *testing.T) {
	page := NewPage(blockSize)
	page.AddString("Bolt-DB")
	page.MutateString(0, "RocksDB")

	assert.Equal(t, "RocksDB", page.GetString(0))
}

func TestAddFewFieldsInPageDecodeThePageAndMutateFieldsInTheDecodedPageToSimulateLoadingPageFromDiskAndChangingIt(t *testing.T) {
	page := NewPage(blockSize)
	page.AddString("PebbleDB is an LSM-based key/value storage engine")
	page.AddUint32(32)
	page.AddUint16(16)
	page.AddUint64(64)
	page.finish()

	decodedPage := &Page{}
	decodedPage.DecodeFrom(page.buffer)
	decodedPage.MutateString(0, "Rocks-DB is an LSM-based key/value storage engine")
	decodedPage.MutateUint16(2, 160)
	decodedPage.MutateUint64(3, 640)

	assert.Equal(t, "Rocks-DB is an LSM-based key/value storage engine", decodedPage.GetString(0))
	assert.Equal(t, uint32(32), decodedPage.GetUint32(1))
	assert.Equal(t, uint16(160), decodedPage.GetUint16(2))
	assert.Equal(t, uint64(640), decodedPage.GetUint64(3))
}
