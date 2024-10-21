package file

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateALogPageWithASingleField(t *testing.T) {
	pageBuilder := NewPageBuilder(blockSize)
	pageBuilder.addUint8(10)

	page := pageBuilder.build()
	assert.Equal(t, uint8(10), page.getUint8(0))
}

func TestCreateALogPageWithFewFields(t *testing.T) {
	pageBuilder := NewPageBuilder(blockSize)
	pageBuilder.addUint16(16)
	pageBuilder.addUint8(8)
	pageBuilder.addUint32(32)
	pageBuilder.addUint16(100)

	page := pageBuilder.build()
	assert.Equal(t, uint16(16), page.getUint16(0))
	assert.Equal(t, uint8(8), page.getUint8(1))
	assert.Equal(t, uint32(32), page.getUint32(2))
	assert.Equal(t, uint16(100), page.getUint16(3))
}

func TestCreateALogPageByteSliceAndString(t *testing.T) {
	pageBuilder := NewPageBuilder(blockSize)
	pageBuilder.addBytes([]byte("RocksDB is an LSM-based key/value storage engine"))
	pageBuilder.addString("PebbleDB is an LSM-based key/value storage engine")

	page := pageBuilder.build()
	assert.Equal(t, []byte("RocksDB is an LSM-based key/value storage engine"), page.getBytes(0))
	assert.Equal(t, "PebbleDB is an LSM-based key/value storage engine", page.getString(1))
}
