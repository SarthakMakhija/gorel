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

	sourceBytes := DecodeByteSlice(destination[:numberOfBytesForEncoding], 0)
	assert.Equal(t, "LSM stands for log-structured merge tree", string(sourceBytes))
}
