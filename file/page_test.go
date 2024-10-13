package file

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

const blockSize = 4096

func TestPageWithAnInt8Value(t *testing.T) {
	table := []struct {
		offset uint
		value  int8
	}{
		{offset: 5, value: 10},
		{offset: 5, value: -10},
	}

	for _, entry := range table {
		page := NewPage(blockSize)
		page.setInt8(entry.offset, entry.value)

		assert.Equal(t, entry.value, page.getInt8(entry.offset))
	}
}

func TestPageWithAnInt16Value(t *testing.T) {
	table := []struct {
		offset uint
		value  int16
	}{
		{offset: 5, value: 10},
		{offset: 5, value: -10},
	}

	for _, entry := range table {
		page := NewPage(blockSize)
		page.setInt16(entry.offset, entry.value)

		assert.Equal(t, entry.value, page.getInt16(entry.offset))
	}
}

func TestPageWithAnInt32Value(t *testing.T) {
	table := []struct {
		offset uint
		value  int32
	}{
		{offset: 5, value: 10},
		{offset: 5, value: -10},
	}

	for _, entry := range table {
		page := NewPage(blockSize)
		page.setInt32(entry.offset, entry.value)

		assert.Equal(t, entry.value, page.getInt32(entry.offset))
	}
}

func TestPageWithAnInt64Value(t *testing.T) {
	table := []struct {
		offset uint
		value  int64
	}{
		{offset: 5, value: 10},
		{offset: 5, value: -10},
	}

	for _, entry := range table {
		page := NewPage(blockSize)
		page.setInt64(entry.offset, entry.value)

		assert.Equal(t, entry.value, page.getInt64(entry.offset))
	}
}

func TestPageWithAnIntValue(t *testing.T) {
	table := []struct {
		offset uint
		value  int
	}{
		{offset: 5, value: 10},
		{offset: 5, value: -10},
	}

	for _, entry := range table {
		page := NewPage(blockSize)
		page.setInt(entry.offset, entry.value)

		assert.Equal(t, entry.value, page.getInt(entry.offset))
	}
}

func TestPageWithAnUInt8Value(t *testing.T) {
	page := NewPage(blockSize)
	page.setUint8(5, 100)

	assert.Equal(t, uint8(100), page.getUint8(5))
}

func TestPageWithAnUInt16Value(t *testing.T) {
	page := NewPage(blockSize)
	page.setUint16(5, 100)

	assert.Equal(t, uint16(100), page.getUint16(5))
}

func TestPageWithAnUInt32Value(t *testing.T) {
	page := NewPage(blockSize)
	page.setUint32(5, 100)

	assert.Equal(t, uint32(100), page.getUint32(5))
}

func TestPageWithAnUInt64Value(t *testing.T) {
	page := NewPage(blockSize)
	page.setUint64(5, 100)

	assert.Equal(t, uint64(100), page.getUint64(5))
}

func TestPageWithAnUintValue(t *testing.T) {
	page := NewPage(blockSize)
	page.setUint(5, 100)

	assert.Equal(t, uint(100), page.getUint(5))
}

func TestPageWithAFloat32Value(t *testing.T) {
	table := []struct {
		offset uint
		value  float32
	}{
		{offset: 5, value: float32(3.19)},
		{offset: 5, value: float32(-3.19)},
	}

	for _, entry := range table {
		page := NewPage(blockSize)
		page.setFloat32(entry.offset, entry.value)

		assert.Equal(t, entry.value, page.getFloat32(entry.offset))
	}
}

func TestPageWithAFloat64Value(t *testing.T) {
	table := []struct {
		offset uint
		value  float64
	}{
		{offset: 5, value: 3.19},
		{offset: 5, value: -3.19},
	}

	for _, entry := range table {
		page := NewPage(blockSize)
		page.setFloat64(entry.offset, entry.value)

		assert.Equal(t, entry.value, page.getFloat64(entry.offset))
	}
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
