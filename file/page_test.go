package file

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

const blockSize = 4096

func TestPageWithAPositiveIntValue(t *testing.T) {
	page := NewPage(blockSize)
	page.setInt(5, 1000)

	assert.Equal(t, 1000, page.getInt(5))
}

func TestPageWithANegativeIntValue(t *testing.T) {
	page := NewPage(blockSize)
	page.setInt(5, -1000)

	assert.Equal(t, -1000, page.getInt(5))
}

func TestPageWithAnUint32Value(t *testing.T) {
	page := NewPage(blockSize)
	page.setUint32(5, uint32(100))

	assert.Equal(t, uint32(100), page.getUint32(5))
}

func TestPageWithAByteSliceValue(t *testing.T) {
	page := NewPage(blockSize)
	page.setBytes(5, []byte("RockDB is an LSM-based storage engine"))

	assert.Equal(t, []byte("RockDB is an LSM-based storage engine"), page.getBytes(5))
}

func TestPageWithAStringValue(t *testing.T) {
	page := NewPage(blockSize)
	page.setBytes(5, []byte("PebbleDB is also an LSM-based storage engine"))

	assert.Equal(t, "PebbleDB is also an LSM-based storage engine", page.getString(5))
}
