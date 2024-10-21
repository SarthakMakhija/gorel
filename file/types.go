package file

import (
	"unsafe"
)

var reservedSizeForAType = int(unsafe.Sizeof(uint8(0)))

type Type uint8

const (
	typeUint8     Type = 1
	typeUint16    Type = 2
	typeUint32    Type = 3
	typeUint64    Type = 4
	typeString    Type = 5
	typeByteSlice Type = 6
)

func (typeDescription Type) asString() string {
	switch typeDescription {
	case typeUint8:
		return "uint8"
	case typeUint16:
		return "uint16"
	case typeUint32:
		return "uint32"
	case typeUint64:
		return "uint64"
	case typeString:
		return "string"
	case typeByteSlice:
		return "[]byte"
	}
	return ""
}

func (typeDescription Type) equals(other Type) bool {
	return uint8(typeDescription) == uint8(other)
}

type Types struct {
	description []Type
}

func newTypes() *Types {
	return &Types{}
}

func decodeTypesFrom(buffer []byte) *Types {
	types := newTypes()
	for _, description := range buffer {
		types.description = append(types.description, Type(description))
	}
	return types
}

func (types *Types) addTypeDescription(description Type) {
	types.description = append(types.description, description)
}

func (types *Types) encode() []byte {
	buffer := make([]byte, len(types.description)*reservedSizeForAType)
	offsetIndex := 0
	for _, definition := range types.description {
		buffer[offsetIndex] = uint8(definition)
		offsetIndex += reservedSizeForAType
	}
	return buffer
}

func (types *Types) getTypeAt(index int) Type {
	return types.description[index]
}
