package file

import (
	"unsafe"
)

var reservedSizeForAType = int(unsafe.Sizeof(uint8(0)))

const (
	typeUint8     uint8 = 1
	typeUint16    uint8 = 2
	typeUint32    uint8 = 3
	typeUint64    uint8 = 4
	typeString    uint8 = 5
	typeByteSlice uint8 = 6
)

type Types struct {
	description []uint8
}

func NewTypes() *Types {
	return &Types{}
}

func DecodeTypesFrom(buffer []byte) *Types {
	types := NewTypes()
	for _, description := range buffer {
		types.description = append(types.description, description)
	}
	return types
}

func (types *Types) addTypeDescriptionUint8() {
	types.description = append(types.description, typeUint8)
}

func (types *Types) addTypeDescriptionUint16() {
	types.description = append(types.description, typeUint16)
}

func (types *Types) addTypeDescriptionUint32() {
	types.description = append(types.description, typeUint32)
}

func (types *Types) addTypeDescriptionUint64() {
	types.description = append(types.description, typeUint64)
}

func (types *Types) addTypeDescriptionString() {
	types.description = append(types.description, typeString)
}

func (types *Types) addTypeDescriptionByteSlice() {
	types.description = append(types.description, typeByteSlice)
}

func (types *Types) encode() []byte {
	buffer := make([]byte, len(types.description)*reservedSizeForAType)
	offsetIndex := 0
	for _, definition := range types.description {
		buffer[offsetIndex] = definition
		offsetIndex += reservedSizeForAType
	}
	return buffer
}
