package buffer

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

const blockSize = 4096

func TestCreateAPageWithASingleField(t *testing.T) {
	page := NewPage(blockSize)
	page.AddUint8(10)
	page.Finish()

	assert.Equal(t, uint8(10), page.GetUint8(0))
}

func TestCreateAPageWithFewFields(t *testing.T) {
	page := NewPage(blockSize)
	page.AddUint16(16)
	page.AddUint8(8)
	page.AddUint32(32)
	page.AddUint16(100)
	page.Finish()

	assert.Equal(t, uint16(16), page.GetUint16(0))
	assert.Equal(t, uint8(8), page.GetUint8(1))
	assert.Equal(t, uint32(32), page.GetUint32(2))
	assert.Equal(t, uint16(100), page.GetUint16(3))
}

func TestCreateAPageWithByteSliceAndString(t *testing.T) {
	page := NewPage(blockSize)
	page.AddBytes([]byte("RocksDB is an LSM-based key/value storage engine"))
	page.AddString("PebbleDB is an LSM-based key/value storage engine")
	page.Finish()

	assert.Equal(t, []byte("RocksDB is an LSM-based key/value storage engine"), page.GetBytes(0))
	assert.Equal(t, "PebbleDB is an LSM-based key/value storage engine", page.GetString(1))
}

func TestAttemptToGetTheValueAtAnIndexGreaterThanTheNumberOfAvailableFields(t *testing.T) {
	page := NewPage(blockSize)
	page.AddString("PebbleDB is an LSM-based key/value storage engine")
	page.Finish()

	decodedPage := &Page{}
	decodedPage.DecodeFrom(page.buffer)

	assert.Panics(t, func() {
		decodedPage.GetString(1)
	})
}

func TestAttemptToGetTheValueWithMismatchedTypeDescription(t *testing.T) {
	page := NewPage(blockSize)
	page.AddString("PebbleDB is an LSM-based key/value storage engine")
	page.Finish()

	decodedPage := &Page{}
	decodedPage.DecodeFrom(page.buffer)

	assert.Panics(t, func() {
		decodedPage.GetUint8(0)
	})
}

func TestDecodeAPageWithASingleField(t *testing.T) {
	page := NewPage(blockSize)
	page.AddString("PebbleDB is an LSM-based key/value storage engine")
	page.Finish()

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
	page.Finish()

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
	page.Finish()

	decodedPage := &Page{}
	decodedPage.DecodeFrom(page.buffer)
	decodedPage.AddString("RocksDB is an LSM-based key/value storage engine")

	assert.Equal(t, "PebbleDB is an LSM-based key/value storage engine", decodedPage.GetString(0))
	assert.Equal(t, uint32(32), decodedPage.GetUint32(1))
	assert.Equal(t, uint16(16), decodedPage.GetUint16(2))
	assert.Equal(t, uint64(64), decodedPage.GetUint64(3))
	assert.Equal(t, "RocksDB is an LSM-based key/value storage engine", decodedPage.GetString(4))
}
