package convert

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConvertUIntToBytesAndViceVersa(t *testing.T) {
	bytes := IntToBytes[uint](5000)
	assert.Equal(t, uint(5000), BytesToInt[uint](bytes))
}

func TestConvertNegativeIntToBytesAndViceVersa(t *testing.T) {
	bytes := IntToBytes[int](-5000)
	assert.Equal(t, -5000, BytesToInt[int](bytes))
}

func TestConvertPositiveIntToBytesAndViceVersa(t *testing.T) {
	bytes := IntToBytes[int](5000)
	assert.Equal(t, 5000, BytesToInt[int](bytes))
}
