package buffer

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAttemptToGetTheValueAtAnIndexGreaterThanTheNumberOfAvailableFields(t *testing.T) {
	pageBuilder := NewPageBuilder(blockSize)
	pageBuilder.AddString("PebbleDB is an LSM-based key/value storage engine")

	page := pageBuilder.Build()
	decodedPage := &Page{}
	decodedPage.DecodePageFrom(page.buffer)

	assert.Panics(t, func() {
		decodedPage.GetString(1)
	})
}

func TestAttemptToGetTheValueWithMismatchedTypeDescription(t *testing.T) {
	pageBuilder := NewPageBuilder(blockSize)
	pageBuilder.AddString("PebbleDB is an LSM-based key/value storage engine")

	page := pageBuilder.Build()
	decodedPage := &Page{}
	decodedPage.DecodePageFrom(page.buffer)

	assert.Panics(t, func() {
		decodedPage.GetUint8(0)
	})
}

func TestDecodeAPageWithASingleField(t *testing.T) {
	pageBuilder := NewPageBuilder(blockSize)
	pageBuilder.AddString("PebbleDB is an LSM-based key/value storage engine")

	page := pageBuilder.Build()
	decodedPage := &Page{}
	decodedPage.DecodePageFrom(page.buffer)

	assert.Equal(t, "PebbleDB is an LSM-based key/value storage engine", decodedPage.GetString(0))
}

func TestDecodeAPageWithFewFields(t *testing.T) {
	pageBuilder := NewPageBuilder(blockSize)
	pageBuilder.AddString("PebbleDB is an LSM-based key/value storage engine")
	pageBuilder.AddUint16(16)
	pageBuilder.AddUint16(160)
	pageBuilder.AddUint64(64)

	page := pageBuilder.Build()
	decodedPage := &Page{}
	decodedPage.DecodePageFrom(page.buffer)

	assert.Equal(t, uint64(64), decodedPage.GetUint64(3))
	assert.Equal(t, uint16(160), decodedPage.GetUint16(2))
	assert.Equal(t, uint16(16), decodedPage.GetUint16(1))
	assert.Equal(t, "PebbleDB is an LSM-based key/value storage engine", decodedPage.GetString(0))
}
