package file

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestOffsetWithBlockZero(t *testing.T) {
	blockId := NewBlockId("lsm.log", 0)

	blockSize := uint(400)
	assert.Equal(t, int64(0), blockId.offset(blockSize))
}

func TestOffset(t *testing.T) {
	blockId := NewBlockId("lsm.log", 3)
	blockSize := uint(400)

	assert.Equal(t, int64(1200), blockId.offset(blockSize))
}
