package file

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStartingOffsetsLength(t *testing.T) {
	startingOffsets := NewStartingOffsets()
	startingOffsets.Append(100)
	startingOffsets.Append(400)
	startingOffsets.Append(530)

	encoded := startingOffsets.Encode()

	decodedStartingOffsets := DecodeStartingOffsetsFrom(encoded)
	assert.Equal(t, 3, decodedStartingOffsets.Length())
}

func TestEncodeAndDecodeASingleStartingOffset(t *testing.T) {
	startingOffsets := NewStartingOffsets()
	startingOffsets.Append(20)

	encoded := startingOffsets.Encode()

	decodedStartingOffsets := DecodeStartingOffsetsFrom(encoded)
	assert.Equal(t, 1, decodedStartingOffsets.Length())
	assert.Equal(t, uint16(20), decodedStartingOffsets.offsets[0])
}

func TestEncodeAndDecodeAFewStartingOffsets(t *testing.T) {
	startingOffsets := NewStartingOffsets()
	startingOffsets.Append(20)
	startingOffsets.Append(400)
	startingOffsets.Append(520)

	encoded := startingOffsets.Encode()

	decodedStartingOffsets := DecodeStartingOffsetsFrom(encoded)
	assert.Equal(t, 3, decodedStartingOffsets.Length())
	assert.Equal(t, uint16(20), decodedStartingOffsets.offsets[0])
	assert.Equal(t, uint16(400), decodedStartingOffsets.offsets[1])
	assert.Equal(t, uint16(520), decodedStartingOffsets.offsets[2])
}
