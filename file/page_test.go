package file

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAttemptToGetTheValueAtAnIndexGreaterThanTheNumberOfAvailableFields(t *testing.T) {
	pageBuilder := NewPageBuilder(blockSize)
	pageBuilder.addString("PebbleDB is an LSM-based key/value storage engine")

	page := pageBuilder.build()
	decodedPage := DecodePageFrom(page.buffer)

	assert.Panics(t, func() {
		decodedPage.getString(1)
	})
}

func TestAttemptToGetTheValueWithMismatchedTypeDescription(t *testing.T) {
	pageBuilder := NewPageBuilder(blockSize)
	pageBuilder.addString("PebbleDB is an LSM-based key/value storage engine")

	page := pageBuilder.build()
	decodedPage := DecodePageFrom(page.buffer)

	assert.Panics(t, func() {
		decodedPage.getUint8(0)
	})
}

func TestDecodeAPageWithASingleField(t *testing.T) {
	pageBuilder := NewPageBuilder(blockSize)
	pageBuilder.addString("PebbleDB is an LSM-based key/value storage engine")

	page := pageBuilder.build()
	decodedPage := DecodePageFrom(page.buffer)

	assert.Equal(t, "PebbleDB is an LSM-based key/value storage engine", decodedPage.getString(0))
}

func TestDecodeAPageWithFewFields(t *testing.T) {
	pageBuilder := NewPageBuilder(blockSize)
	pageBuilder.addString("PebbleDB is an LSM-based key/value storage engine")
	pageBuilder.addUint16(16)
	pageBuilder.addUint16(160)
	pageBuilder.addUint64(64)

	page := pageBuilder.build()
	decodedPage := DecodePageFrom(page.buffer)

	assert.Equal(t, uint64(64), decodedPage.getUint64(3))
	assert.Equal(t, uint16(160), decodedPage.getUint16(2))
	assert.Equal(t, uint16(16), decodedPage.getUint16(1))
	assert.Equal(t, "PebbleDB is an LSM-based key/value storage engine", decodedPage.getString(0))
}
