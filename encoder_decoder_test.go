package gorel

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBytesNeededForEncodingAByteSlice(t *testing.T) {
	assert.Equal(t, 6, BytesNeededForEncodingAByteSlice([]byte("raft")))
}

func TestEncodeAndDecodeAByteSlice(t *testing.T) {
	destination := make([]byte, 100)
	numberOfBytesForEncoding := EncodeByteSlice([]byte("LSM stands for log-structured merge tree"), destination, 0)

	sourceBytes, _ := DecodeByteSlice(destination[:numberOfBytesForEncoding], 0)
	assert.Equal(t, "LSM stands for log-structured merge tree", string(sourceBytes))
}

func TestEncodeAndDecodeAnUint16(t *testing.T) {
	value := uint16(100)
	destination := make([]byte, 2)

	EncodeUint16(value, destination, 0)

	decoded, endOffset := DecodeUint16(destination, 0)
	assert.Equal(t, uint16(100), decoded)
	assert.Equal(t, EndOffset(uint16Size), endOffset)
}

func TestEncodeAndDecodeAnUint32(t *testing.T) {
	value := uint32(5400)
	destination := make([]byte, 4)

	EncodeUint32(value, destination, 0)

	decoded, endOffset := DecodeUint32(destination, 0)
	assert.Equal(t, uint32(5400), decoded)
	assert.Equal(t, EndOffset(uint32Size), endOffset)
}

func TestEncodeAndDecodeAnUint64(t *testing.T) {
	value := uint64(10000)
	destination := make([]byte, 8)

	EncodeUint64(value, destination, 0)

	decoded, endOffset := DecodeUint64(destination, 0)
	assert.Equal(t, uint64(10000), decoded)
	assert.Equal(t, EndOffset(uint64Size), endOffset)
}
