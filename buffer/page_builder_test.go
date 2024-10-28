package buffer

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

const blockSize = 4096

func TestCreateALogPageWithASingleField(t *testing.T) {
	pageBuilder := NewPageBuilder(blockSize)
	pageBuilder.AddUint8(10)

	page := pageBuilder.Build()
	assert.Equal(t, uint8(10), page.GetUint8(0))
}

func TestCreateALogPageWithFewFields(t *testing.T) {
	pageBuilder := NewPageBuilder(blockSize)
	pageBuilder.AddUint16(16)
	pageBuilder.AddUint8(8)
	pageBuilder.AddUint32(32)
	pageBuilder.AddUint16(100)

	page := pageBuilder.Build()
	assert.Equal(t, uint16(16), page.GetUint16(0))
	assert.Equal(t, uint8(8), page.GetUint8(1))
	assert.Equal(t, uint32(32), page.GetUint32(2))
	assert.Equal(t, uint16(100), page.GetUint16(3))
}

func TestCreateALogPageByteSliceAndString(t *testing.T) {
	pageBuilder := NewPageBuilder(blockSize)
	pageBuilder.AddBytes([]byte("RocksDB is an LSM-based key/value storage engine"))
	pageBuilder.AddString("PebbleDB is an LSM-based key/value storage engine")

	page := pageBuilder.Build()
	assert.Equal(t, []byte("RocksDB is an LSM-based key/value storage engine"), page.GetBytes(0))
	assert.Equal(t, "PebbleDB is an LSM-based key/value storage engine", page.GetString(1))
}
