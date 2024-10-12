package convert

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConvertPositiveIntToBytesAndViceVersa(t *testing.T) {
	bytes := IntToBytes(5000)
	assert.Equal(t, 5000, BytesToInt(bytes))
}

func TestConvertNegativeIntToBytesAndViceVersa(t *testing.T) {
	bytes := IntToBytes(-5000)
	assert.Equal(t, -5000, BytesToInt(bytes))
}
