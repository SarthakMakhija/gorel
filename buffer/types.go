package buffer

import (
	"gorel"
	"unsafe"
)

var ReservedSizeForAType = int(unsafe.Sizeof(uint8(0)))

type TypeDescription uint8

const (
	TypeUint8     TypeDescription = 1
	TypeUint16    TypeDescription = 2
	TypeUint32    TypeDescription = 3
	TypeUint64    TypeDescription = 4
	TypeString    TypeDescription = 5
	TypeByteSlice TypeDescription = 6
)

func (typeDescription TypeDescription) AsString() string {
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

func (typeDescription TypeDescription) EndOffsetPostDecode(source []byte, fromOffset uint16) gorel.EndOffset {
	var endOffset gorel.EndOffset
	switch typeDescription {
	case TypeUint8:
		_, endOffset = gorel.DecodeUint8(source, fromOffset)
	case TypeUint16:
		_, endOffset = gorel.DecodeUint16(source, fromOffset)
	case TypeUint32:
		_, endOffset = gorel.DecodeUint32(source, fromOffset)
	case TypeUint64:
		_, endOffset = gorel.DecodeUint64(source, fromOffset)
	case TypeString:
		_, endOffset = gorel.DecodeByteSlice(source, fromOffset)
	case TypeByteSlice:
		_, endOffset = gorel.DecodeByteSlice(source, fromOffset)
	}
	return endOffset
}

func (typeDescription TypeDescription) Equals(other TypeDescription) bool {
	return uint8(typeDescription) == uint8(other)
}

type Types struct {
	description []TypeDescription
}

func NewTypes() *Types {
	return &Types{}
}

func DecodeTypesFrom(buffer []byte) *Types {
	types := NewTypes()
	for _, description := range buffer {
		types.description = append(types.description, TypeDescription(description))
	}
	return types
}

func (types *Types) AddTypeDescription(description TypeDescription) {
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

func (types *Types) GetTypeAt(index int) TypeDescription {
	return types.description[index]
}

func (types *Types) SizeUsedInBytes() int {
	return ReservedSizeForAType * len(types.description)
}

func SizeUsedInBytes(numberOfDescriptions uint16) int {
	return ReservedSizeForAType * int(numberOfDescriptions)
}

func (types *Types) Length() int {
	return len(types.description)
}
