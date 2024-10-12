package convert

import "unsafe"

func IntToBytes(value int) []byte {
	size := int(unsafe.Sizeof(value))
	bytes := make([]byte, size)
	for index := 0; index < size; index++ {
		aByte := *(*uint8)(unsafe.Pointer(uintptr(unsafe.Pointer(&value)) + uintptr(index)))
		bytes[index] = aByte
	}
	return bytes
}

func BytesToInt(bytes []byte) int {
	value := 0
	size := len(bytes)
	for index := 0; index < size; index++ {
		*(*uint8)(unsafe.Pointer(uintptr(unsafe.Pointer(&value)) + uintptr(index))) = bytes[index]
	}
	return value
}
