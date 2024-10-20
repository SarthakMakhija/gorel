package file

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEncodeAndDecodeTypesWithASingleTypeDescription(t *testing.T) {
	types := NewTypes()
	types.addTypeDescriptionUint8()

	encoded := types.encode()
	decodedTypes := DecodeTypesFrom(encoded)

	assert.Equal(t, 1, len(decodedTypes.description))
	assert.Equal(t, typeUint8, decodedTypes.description[0])
}

func TestEncodeAndDecodeTypesWithACoupleOfTypeDescriptions(t *testing.T) {
	types := NewTypes()
	types.addTypeDescriptionUint16()
	types.addTypeDescriptionUint8()

	encoded := types.encode()
	decodedTypes := DecodeTypesFrom(encoded)

	assert.Equal(t, 2, len(decodedTypes.description))
	assert.Equal(t, typeUint16, decodedTypes.description[0])
	assert.Equal(t, typeUint8, decodedTypes.description[1])
}

func TestEncodeAndDecodeTypesWithFewTypeDescriptions(t *testing.T) {
	types := NewTypes()
	types.addTypeDescriptionUint16()
	types.addTypeDescriptionUint8()
	types.addTypeDescriptionUint64()
	types.addTypeDescriptionUint32()

	encoded := types.encode()
	decodedTypes := DecodeTypesFrom(encoded)

	assert.Equal(t, 4, len(decodedTypes.description))
	assert.Equal(t, typeUint16, decodedTypes.description[0])
	assert.Equal(t, typeUint8, decodedTypes.description[1])
	assert.Equal(t, typeUint64, decodedTypes.description[2])
	assert.Equal(t, typeUint32, decodedTypes.description[3])
}
