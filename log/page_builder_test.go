package log

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

const blockSize = 4096

func TestCreateALogPageWithASingleRecord(t *testing.T) {
	pageBuilder := NewPageBuilder(blockSize)
	pageBuilder.Add([]byte("RocksDB is an LSM-based key/value storage engine"))

	page := pageBuilder.Build()
	iterator := page.BackwardIterator()

	assert.True(t, iterator.IsValid())
	assert.Equal(t, "RocksDB is an LSM-based key/value storage engine", string(iterator.Record()))

	iterator.Previous()
	assert.False(t, iterator.IsValid())
}

func TestCreateALogPageWithCoupleOfRecords(t *testing.T) {
	pageBuilder := NewPageBuilder(blockSize)
	pageBuilder.Add([]byte("RocksDB is an LSM-based key/value storage engine"))
	pageBuilder.Add([]byte("PebbleDB is an LSM-based key/value storage engine"))

	page := pageBuilder.Build()
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
	pageBuilder := NewPageBuilder(blockSize)

	const records = 100
	for record := 1; record <= records; record++ {
		pageBuilder.Add([]byte(fmt.Sprintf("Record %d", record)))
	}

	page := pageBuilder.Build()
	iterator := page.BackwardIterator()

	for record := 100; record >= 1; record-- {
		assert.True(t, iterator.IsValid())
		assert.Equal(t, fmt.Sprintf("Record %d", record), string(iterator.Record()))
		iterator.Previous()
	}
	assert.False(t, iterator.IsValid())
}
