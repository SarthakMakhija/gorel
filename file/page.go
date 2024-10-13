package file

import (
	"encoding/binary"
	"gorel/convert"
	"math"
	"unsafe"
)

var (
	reservedBufferSize = uint(unsafe.Sizeof(uint32(0)))
	uint8Size          = uint(unsafe.Sizeof(uint8(0)))
	uint16Size         = uint(unsafe.Sizeof(uint16(0)))
	uint32Size         = uint(unsafe.Sizeof(uint32(0)))
	uint64Size         = uint(unsafe.Sizeof(uint64(0)))
	uintSize           = uint(unsafe.Sizeof(uint(0)))
	intSize            = uint(unsafe.Sizeof(0))
	float32Size        = uint(unsafe.Sizeof(float32(0)))
	float64Size        = uint(unsafe.Sizeof(float64(0)))
)

type Page struct {
	buffer             []byte
	currentWriteOffset uint
	currentReadOffset  uint
}

// NewPage
// TODO: Check valid currentWriteOffset in all the set methods, currentWriteOffset < len(buffer)
// TODO: Check valid currentWriteOffset in all the get method,  currentReadOffset < currentWriteOffset
// TODO: Check that the page has enough space to accommodate the incoming value
// TODO: Add support for Date
func NewPage(blockSize uint) *Page {
	return &Page{
		buffer:             make([]byte, blockSize),
		currentWriteOffset: 0,
		currentReadOffset:  0,
	}
}

func (page *Page) setInt8(value int8) {
	page.setUint8(uint8(value))
}

func (page *Page) getInt8() int8 {
	return int8(page.getUint8())
}

func (page *Page) setInt16(value int16) {
	page.setUint16(uint16(value))
}

func (page *Page) getInt16() int16 {
	return int16(page.getUint16())
}

func (page *Page) setInt32(value int32) {
	page.setUint32(uint32(value))
}

func (page *Page) getInt32() int32 {
	return int32(page.getUint32())
}

func (page *Page) setInt64(value int64) {
	page.setUint64(uint64(value))
}

func (page *Page) getInt64() int64 {
	return int64(page.getUint64())
}

func (page *Page) setUint8(value uint8) {
	page.buffer[page.currentWriteOffset] = value
	page.moveCurrentWriteOffsetBy(uint8Size)
}

func (page *Page) getUint8() uint8 {
	value := page.buffer[page.currentReadOffset]
	page.moveCurrentReadOffsetBy(uint8Size)
	return value
}

func (page *Page) setUint16(value uint16) {
	binary.LittleEndian.PutUint16(page.buffer[page.currentWriteOffset:], value)
	page.moveCurrentWriteOffsetBy(uint16Size)
}

func (page *Page) getUint16() uint16 {
	value := binary.LittleEndian.Uint16(page.buffer[page.currentReadOffset:])
	page.moveCurrentReadOffsetBy(uint16Size)
	return value
}

func (page *Page) setUint32(value uint32) {
	binary.LittleEndian.PutUint32(page.buffer[page.currentWriteOffset:], value)
	page.moveCurrentWriteOffsetBy(uint32Size)
}

func (page *Page) getUint32() uint32 {
	value := binary.LittleEndian.Uint32(page.buffer[page.currentReadOffset:])
	page.moveCurrentReadOffsetBy(uint32Size)
	return value
}

func (page *Page) setUint64(value uint64) {
	binary.LittleEndian.PutUint64(page.buffer[page.currentWriteOffset:], value)
	page.moveCurrentWriteOffsetBy(uint64Size)
}

func (page *Page) getUint64() uint64 {
	value := binary.LittleEndian.Uint64(page.buffer[page.currentReadOffset:])
	page.moveCurrentReadOffsetBy(uint64Size)
	return value
}

func (page *Page) setInt(value int) {
	copy(page.buffer[page.currentWriteOffset:], convert.IntToBytes[int](value))
	page.moveCurrentWriteOffsetBy(intSize)
}

func (page *Page) getInt() int {
	value := convert.BytesToInt[int](page.buffer[page.currentReadOffset : page.currentReadOffset+intSize])
	page.moveCurrentReadOffsetBy(intSize)
	return value
}

func (page *Page) setUint(value uint) {
	copy(page.buffer[page.currentWriteOffset:], convert.IntToBytes[uint](value))
	page.moveCurrentWriteOffsetBy(uintSize)
}

func (page *Page) getUint() uint {
	value := convert.BytesToInt[uint](page.buffer[page.currentReadOffset : page.currentReadOffset+intSize])
	page.moveCurrentReadOffsetBy(uintSize)
	return value
}

func (page *Page) setFloat32(value float32) {
	binary.LittleEndian.PutUint32(page.buffer[page.currentWriteOffset:], math.Float32bits(value))
	page.moveCurrentWriteOffsetBy(float32Size)
}

func (page *Page) getFloat32() float32 {
	return math.Float32frombits(page.getUint32())
}

func (page *Page) setFloat64(value float64) {
	binary.LittleEndian.PutUint64(page.buffer[page.currentWriteOffset:], math.Float64bits(value))
	page.moveCurrentWriteOffsetBy(float64Size)
}

func (page *Page) getFloat64() float64 {
	return math.Float64frombits(page.getUint64())
}

func (page *Page) setString(str string) {
	//TODO: str size should be less than 2^32-1 (close to 4G)
	page.setBytes([]byte(str))
}

func (page *Page) getString() string {
	return string(page.getBytes())
}

func (page *Page) setBytes(buffer []byte) {
	//TODO: Buffer size should be less than 2^32-1 (close to 4G)
	binary.LittleEndian.PutUint32(page.buffer[page.currentWriteOffset:], uint32(len(buffer)))
	copy(page.buffer[page.currentWriteOffset+reservedBufferSize:], buffer)
	page.moveCurrentWriteOffsetBy(uint(len(buffer)) + reservedBufferSize)
}

func (page *Page) getBytes() []byte {
	bufferLength := binary.LittleEndian.Uint32(page.buffer[page.currentReadOffset:])
	endOffset := page.currentReadOffset + reservedBufferSize + uint(bufferLength)
	value := page.buffer[page.currentReadOffset+reservedBufferSize : endOffset]
	page.moveCurrentReadOffsetBy(uint(len(value)) + reservedBufferSize)
	return value
}

func (page *Page) setBool(value bool) {
	if value {
		page.buffer[page.currentWriteOffset] = 1
	} else {
		page.buffer[page.currentWriteOffset] = 0
	}
	page.moveCurrentWriteOffsetBy(1)
}

func (page *Page) getBool() bool {
	value := page.buffer[page.currentReadOffset] != 0
	page.moveCurrentReadOffsetBy(1)
	return value
}

func (page *Page) moveCurrentWriteOffsetBy(offset uint) {
	page.currentWriteOffset += offset
}

func (page *Page) moveCurrentReadOffsetBy(offset uint) {
	page.currentReadOffset += offset
}
