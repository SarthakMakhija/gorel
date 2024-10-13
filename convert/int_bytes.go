package convert

import "unsafe"

type MachineDependentNumericValueType interface {
	int | uint
}

func IntToBytes[V MachineDependentNumericValueType](value V) []byte {
	size := int(unsafe.Sizeof(value))
	bytes := make([]byte, size)
	for index := 0; index < size; index++ {
		aByte := *(*uint8)(unsafe.Pointer(uintptr(unsafe.Pointer(&value)) + uintptr(index)))
		bytes[index] = aByte
	}
	return bytes
}

func BytesToInt[V MachineDependentNumericValueType](bytes []byte) V {
	var value V
	size := len(bytes)
	for index := 0; index < size; index++ {
		*(*uint8)(unsafe.Pointer(uintptr(unsafe.Pointer(&value)) + uintptr(index))) = bytes[index]
	}
	return value
}
