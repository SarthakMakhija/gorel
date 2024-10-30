package file

import (
	"unsafe"
)

var ReservedSizeForAType = int(unsafe.Sizeof(uint8(0)))

type Type uint8

const (
	TypeUint8     Type = 1
	TypeUint16    Type = 2
	TypeUint32    Type = 3
	TypeUint64    Type = 4
	TypeString    Type = 5
	TypeByteSlice Type = 6
)

func (typeDescription Type) AsString() string {
	switch typeDescription {
	case TypeUint8:
		return "uint8"
	case TypeUint16:
		return "uint16"
	case TypeUint32:
		return "uint32"
	case TypeUint64:
		return "uint64"
	case TypeString:
		return "string"
	case TypeByteSlice:
		return "[]byte"
	}
	return ""
}

func (typeDescription Type) Equals(other Type) bool {
	return uint8(typeDescription) == uint8(other)
}

type Types struct {
	description []Type
}

func NewTypes() *Types {
	return &Types{}
}

func DecodeTypesFrom(buffer []byte) *Types {
	types := NewTypes()
	for _, description := range buffer {
		types.description = append(types.description, Type(description))
	}
	return types
}

func (types *Types) AddTypeDescription(description Type) {
	types.description = append(types.description, description)
}

func (types *Types) Encode() []byte {
	buffer := make([]byte, len(types.description)*ReservedSizeForAType)
	offsetIndex := 0
	for _, definition := range types.description {
		buffer[offsetIndex] = uint8(definition)
		offsetIndex += ReservedSizeForAType
	}
	return buffer
}

func (types *Types) GetTypeAt(index int) Type {
	return types.description[index]
}

func (types *Types) SizeUsedInBytes() int {
	return ReservedSizeForAType * len(types.description)
}

func SizeUsedInBytes(numberOfDescriptions uint16) int {
	return ReservedSizeForAType * int(numberOfDescriptions)
}
