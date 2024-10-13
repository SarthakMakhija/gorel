package file

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

const blockSize = 4096

func TestPageWithAnInt8Value(t *testing.T) {
	table := []struct {
		value int8
	}{
		{value: 10},
		{value: -10},
	}

	for _, entry := range table {
		page := NewPage(blockSize)
		page.setInt8(entry.value)

		assert.Equal(t, entry.value, page.getInt8())
	}
}

func TestPageWithAnInt16Value(t *testing.T) {
	table := []struct {
		value int16
	}{
		{value: 10},
		{value: -10},
	}

	for _, entry := range table {
		page := NewPage(blockSize)
		page.setInt16(entry.value)

		assert.Equal(t, entry.value, page.getInt16())
	}
}

func TestPageWithAnInt32Value(t *testing.T) {
	table := []struct {
		value int32
	}{
		{value: 10},
		{value: -10},
	}

	for _, entry := range table {
		page := NewPage(blockSize)
		page.setInt32(entry.value)

		assert.Equal(t, entry.value, page.getInt32())
	}
}

func TestPageWithAnInt64Value(t *testing.T) {
	table := []struct {
		value int64
	}{
		{value: 10},
		{value: -10},
	}

	for _, entry := range table {
		page := NewPage(blockSize)
		page.setInt64(entry.value)

		assert.Equal(t, entry.value, page.getInt64())
	}
}

func TestPageWithAnIntValue(t *testing.T) {
	table := []struct {
		value int
	}{
		{value: 10},
		{value: -10},
	}

	for _, entry := range table {
		page := NewPage(blockSize)
		page.setInt(entry.value)

		assert.Equal(t, entry.value, page.getInt())
	}
}

func TestPageWithAnUInt8Value(t *testing.T) {
	page := NewPage(blockSize)
	page.setUint8(100)

	assert.Equal(t, uint8(100), page.getUint8())
}

func TestPageWithAnUInt16Value(t *testing.T) {
	page := NewPage(blockSize)
	page.setUint16(100)

	assert.Equal(t, uint16(100), page.getUint16())
}

func TestPageWithAnUInt32Value(t *testing.T) {
	page := NewPage(blockSize)
	page.setUint32(100)

	assert.Equal(t, uint32(100), page.getUint32())
}

func TestPageWithAnUInt64Value(t *testing.T) {
	page := NewPage(blockSize)
	page.setUint64(100)

	assert.Equal(t, uint64(100), page.getUint64())
}

func TestPageWithAnUintValue(t *testing.T) {
	page := NewPage(blockSize)
	page.setUint(100)

	assert.Equal(t, uint(100), page.getUint())
}

func TestPageWithAFloat32Value(t *testing.T) {
	table := []struct {
		value float32
	}{
		{value: float32(3.19)},
		{value: float32(-3.19)},
	}

	for _, entry := range table {
		page := NewPage(blockSize)
		page.setFloat32(entry.value)

		assert.Equal(t, entry.value, page.getFloat32())
	}
}

func TestPageWithAFloat64Value(t *testing.T) {
	table := []struct {
		value float64
	}{
		{value: 3.19},
		{value: -3.19},
	}

	for _, entry := range table {
		page := NewPage(blockSize)
		page.setFloat64(entry.value)

		assert.Equal(t, entry.value, page.getFloat64())
	}
}

func TestPageWithAStringValue(t *testing.T) {
	page := NewPage(blockSize)
	page.setBytes([]byte("PebbleDB is also an LSM-based storage engine"))

	assert.Equal(t, "PebbleDB is also an LSM-based storage engine", page.getString())
}

func TestPageWithAByteSliceValue(t *testing.T) {
	page := NewPage(blockSize)
	page.setBytes([]byte("RockDB is an LSM-based storage engine"))

	assert.Equal(t, []byte("RockDB is an LSM-based storage engine"), page.getBytes())
}

func TestPageWithABoolValue(t *testing.T) {
	table := []struct {
		value bool
	}{
		{value: true},
		{value: false},
	}

	for _, entry := range table {
		page := NewPage(blockSize)
		page.setBool(entry.value)

		assert.Equal(t, entry.value, page.getBool())
	}
}

func TestPageWithMultipleValueTypes(t *testing.T) {
	page := NewPage(blockSize)

	page.setUint8(127)
	page.setInt(500)
	page.setInt(-500)
	page.setString("SlateDB is an object storage built on top of LSM")
	page.setFloat32(float32(12.67))
	page.setUint16(3000)
	page.setFloat64(-12.67)

	assert.Equal(t, uint8(127), page.getUint8())
	assert.Equal(t, 500, page.getInt())
	assert.Equal(t, -500, page.getInt())
	assert.Equal(t, "SlateDB is an object storage built on top of LSM", page.getString())
	assert.Equal(t, float32(12.67), page.getFloat32())
	assert.Equal(t, uint16(3000), page.getUint16())
	assert.Equal(t, -12.67, page.getFloat64())
}
