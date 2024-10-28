package file

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEncodeAndDecodeTypesWithASingleTypeDescription(t *testing.T) {
	types := NewTypes()
	types.AddTypeDescription(TypeUint8)

	encoded := types.Encode()
	decodedTypes := DecodeTypesFrom(encoded)

	assert.Equal(t, 1, len(decodedTypes.description))
	assert.Equal(t, TypeUint8, decodedTypes.description[0])
}

func TestEncodeAndDecodeTypesWithACoupleOfTypeDescriptions(t *testing.T) {
	types := NewTypes()
	types.AddTypeDescription(TypeUint16)
	types.AddTypeDescription(TypeUint8)

	encoded := types.Encode()
	decodedTypes := DecodeTypesFrom(encoded)

	assert.Equal(t, 2, len(decodedTypes.description))
	assert.Equal(t, TypeUint16, decodedTypes.description[0])
	assert.Equal(t, TypeUint8, decodedTypes.description[1])
}

func TestEncodeAndDecodeTypesWithFewTypeDescriptions(t *testing.T) {
	types := NewTypes()
	types.AddTypeDescription(TypeUint16)
	types.AddTypeDescription(TypeUint8)
	types.AddTypeDescription(TypeUint64)
	types.AddTypeDescription(TypeUint32)

	encoded := types.Encode()
	decodedTypes := DecodeTypesFrom(encoded)

	assert.Equal(t, 4, len(decodedTypes.description))
	assert.Equal(t, TypeUint16, decodedTypes.description[0])
	assert.Equal(t, TypeUint8, decodedTypes.description[1])
	assert.Equal(t, TypeUint64, decodedTypes.description[2])
	assert.Equal(t, TypeUint32, decodedTypes.description[3])
}

func TestEncodeAndDecodeTypesWithAStringTypeDescription(t *testing.T) {
	types := NewTypes()
	types.AddTypeDescription(TypeString)

	encoded := types.Encode()
	decodedTypes := DecodeTypesFrom(encoded)

	assert.Equal(t, 1, len(decodedTypes.description))
	assert.Equal(t, TypeString, decodedTypes.description[0])
}

func TestEncodeAndDecodeTypesWithAByteSliceTypeDescription(t *testing.T) {
	types := NewTypes()
	types.AddTypeDescription(TypeByteSlice)

	encoded := types.Encode()
	decodedTypes := DecodeTypesFrom(encoded)

	assert.Equal(t, 1, len(decodedTypes.description))
	assert.Equal(t, TypeByteSlice, decodedTypes.description[0])
}
