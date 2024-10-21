package file

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEncodeAndDecodeTypesWithASingleTypeDescription(t *testing.T) {
	types := newTypes()
	types.addTypeDescription(typeUint8)

	encoded := types.encode()
	decodedTypes := decodeTypesFrom(encoded)

	assert.Equal(t, 1, len(decodedTypes.description))
	assert.Equal(t, typeUint8, decodedTypes.description[0])
}

func TestEncodeAndDecodeTypesWithACoupleOfTypeDescriptions(t *testing.T) {
	types := newTypes()
	types.addTypeDescription(typeUint16)
	types.addTypeDescription(typeUint8)

	encoded := types.encode()
	decodedTypes := decodeTypesFrom(encoded)

	assert.Equal(t, 2, len(decodedTypes.description))
	assert.Equal(t, typeUint16, decodedTypes.description[0])
	assert.Equal(t, typeUint8, decodedTypes.description[1])
}

func TestEncodeAndDecodeTypesWithFewTypeDescriptions(t *testing.T) {
	types := newTypes()
	types.addTypeDescription(typeUint16)
	types.addTypeDescription(typeUint8)
	types.addTypeDescription(typeUint64)
	types.addTypeDescription(typeUint32)

	encoded := types.encode()
	decodedTypes := decodeTypesFrom(encoded)

	assert.Equal(t, 4, len(decodedTypes.description))
	assert.Equal(t, typeUint16, decodedTypes.description[0])
	assert.Equal(t, typeUint8, decodedTypes.description[1])
	assert.Equal(t, typeUint64, decodedTypes.description[2])
	assert.Equal(t, typeUint32, decodedTypes.description[3])
}

func TestEncodeAndDecodeTypesWithAStringTypeDescription(t *testing.T) {
	types := newTypes()
	types.addTypeDescription(typeString)

	encoded := types.encode()
	decodedTypes := decodeTypesFrom(encoded)

	assert.Equal(t, 1, len(decodedTypes.description))
	assert.Equal(t, typeString, decodedTypes.description[0])
}

func TestEncodeAndDecodeTypesWithAByteSliceTypeDescription(t *testing.T) {
	types := newTypes()
	types.addTypeDescription(typeByteSlice)

	encoded := types.encode()
	decodedTypes := decodeTypesFrom(encoded)

	assert.Equal(t, 1, len(decodedTypes.description))
	assert.Equal(t, typeByteSlice, decodedTypes.description[0])
}
