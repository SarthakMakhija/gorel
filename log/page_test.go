package log

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

const blockSize = 4096

func TestAttemptToAddARecordInAPageWithInsufficientSize(t *testing.T) {
	page := NewPage(30)
	assert.False(t, page.Add([]byte("RocksDB is an LSM-based key/value storage engine")))
}

func TestAttemptToAddACoupleOfRecordsInAPageWithSizeSufficientForOnlyOneRecord(t *testing.T) {
	page := NewPage(60)
	assert.True(t, page.Add([]byte("RocksDB is an LSM-based key/value storage engine")))
	assert.False(t, page.Add([]byte("RocksDB is an LSM-based key/value storage engine")))
}

func TestAttemptToAddACoupleOfRecordsSuccessfullyInAPageWithJustEnoughSize(t *testing.T) {
	page := NewPage(108)
	assert.True(t, page.Add([]byte("RocksDB is an LSM-based key/value storage engine")))
	assert.True(t, page.Add([]byte("RocksDB is an LSM-based key/value storage engine")))
}

func TestCreateALogPageWithNoRecords(t *testing.T) {
	page := NewPage(blockSize)

	page.finish()
	iterator := page.BackwardIterator()

	assert.False(t, iterator.IsValid())
}

func TestCreateALogPageWithASingleRecord(t *testing.T) {
	page := NewPage(blockSize)
	page.Add([]byte("RocksDB is an LSM-based key/value storage engine"))

	page.finish()
	iterator := page.BackwardIterator()

	assert.True(t, iterator.IsValid())
	assert.Equal(t, "RocksDB is an LSM-based key/value storage engine", string(iterator.Record()))

	iterator.Previous()
	assert.False(t, iterator.IsValid())
}

func TestCreateALogPageWithCoupleOfRecords(t *testing.T) {
	page := NewPage(blockSize)
	page.Add([]byte("RocksDB is an LSM-based key/value storage engine"))
	page.Add([]byte("PebbleDB is an LSM-based key/value storage engine"))

	page.finish()
	iterator := page.BackwardIterator()

	assert.True(t, iterator.IsValid())
	assert.Equal(t, "PebbleDB is an LSM-based key/value storage engine", string(iterator.Record()))

	iterator.Previous()

	assert.True(t, iterator.IsValid())
	assert.Equal(t, "RocksDB is an LSM-based key/value storage engine", string(iterator.Record()))

	iterator.Previous()
	assert.False(t, iterator.IsValid())
}

func TestCreateALogPageWithFewRecords(t *testing.T) {
	page := NewPage(blockSize)

	const records = 100
	for record := 1; record <= records; record++ {
		page.Add([]byte(fmt.Sprintf("Record %d", record)))
	}

	page.finish()
	iterator := page.BackwardIterator()

	for record := 100; record >= 1; record-- {
		assert.True(t, iterator.IsValid())
		assert.Equal(t, fmt.Sprintf("Record %d", record), string(iterator.Record()))
		iterator.Previous()
	}
	assert.False(t, iterator.IsValid())
}

func TestDecodeAPageWithASingleRecord(t *testing.T) {
	page := NewPage(blockSize)
	page.Add([]byte("RocksDB is an LSM-based key/value storage engine"))
	page.finish()

	decodedPage := &Page{}
	decodedPage.DecodeFrom(page.buffer)

	iterator := decodedPage.BackwardIterator()

	assert.True(t, iterator.IsValid())
	assert.Equal(t, "RocksDB is an LSM-based key/value storage engine", string(iterator.Record()))

	iterator.Previous()
	assert.False(t, iterator.IsValid())
}

func TestDecodeAPageWithCoupleOfRecords(t *testing.T) {
	page := NewPage(blockSize)
	page.Add([]byte("RocksDB is an LSM-based key/value storage engine"))
	page.Add([]byte("PebbleDB is an LSM-based key/value storage engine"))
	page.finish()

	decodedPage := &Page{}
	decodedPage.DecodeFrom(page.buffer)

	iterator := decodedPage.BackwardIterator()

	assert.True(t, iterator.IsValid())
	assert.Equal(t, "PebbleDB is an LSM-based key/value storage engine", string(iterator.Record()))

	iterator.Previous()

	assert.True(t, iterator.IsValid())
	assert.Equal(t, "RocksDB is an LSM-based key/value storage engine", string(iterator.Record()))

	iterator.Previous()
	assert.False(t, iterator.IsValid())
}

func TestDecodeAPageWithFewRecords(t *testing.T) {
	page := NewPage(blockSize)

	const records = 100
	for record := 1; record <= records; record++ {
		page.Add([]byte(fmt.Sprintf("Record %d", record)))
	}
	page.finish()

	decodedPage := &Page{}
	decodedPage.DecodeFrom(page.buffer)

	iterator := decodedPage.BackwardIterator()

	for record := 100; record >= 1; record-- {
		assert.True(t, iterator.IsValid())
		assert.Equal(t, fmt.Sprintf("Record %d", record), string(iterator.Record()))
		iterator.Previous()
	}
	assert.False(t, iterator.IsValid())
}
