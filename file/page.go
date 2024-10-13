package file

import (
	"encoding/binary"
	"gorel/convert"
	"math"
	"unsafe"
)

var (
	reservedBufferSize = uint(unsafe.Sizeof(uint32(0)))
	intSize            = uint(unsafe.Sizeof(0))
)

type Page struct {
	buffer []byte
}

// NewPage
// TODO: Check valid offset in all the methods, offset < len(buffer)
// TODO: Check that the page has enough space to accommodate the incoming value
func NewPage(blockSize uint) *Page {
	return &Page{
		buffer: make([]byte, blockSize),
	}
}

func (page *Page) setInt8(offset uint, value int8) {
	page.setUint8(offset, uint8(value))
}

func (page *Page) getInt8(offset uint) int8 {
	return int8(page.getUint8(offset))
}

func (page *Page) setInt16(offset uint, value int16) {
	page.setUint16(offset, uint16(value))
}

func (page *Page) getInt16(offset uint) int16 {
	return int16(page.getUint16(offset))
}

func (page *Page) setInt32(offset uint, value int32) {
	page.setUint32(offset, uint32(value))
}

func (page *Page) getInt32(offset uint) int32 {
	return int32(page.getUint32(offset))
}

func (page *Page) setInt64(offset uint, value int64) {
	page.setUint64(offset, uint64(value))
}

func (page *Page) getInt64(offset uint) int64 {
	return int64(page.getUint64(offset))
}

func (page *Page) setUint8(offset uint, value uint8) {
	page.buffer[offset] = value
}

func (page *Page) getUint8(offset uint) uint8 {
	return page.buffer[offset]
}

func (page *Page) setUint16(offset uint, value uint16) {
	binary.LittleEndian.PutUint16(page.buffer[offset:], value)
}

func (page *Page) getUint16(offset uint) uint16 {
	return binary.LittleEndian.Uint16(page.buffer[offset:])
}

func (page *Page) setUint32(offset uint, value uint32) {
	binary.LittleEndian.PutUint32(page.buffer[offset:], value)
}

func (page *Page) getUint32(offset uint) uint32 {
	return binary.LittleEndian.Uint32(page.buffer[offset:])
}

func (page *Page) setUint64(offset uint, value uint64) {
	binary.LittleEndian.PutUint64(page.buffer[offset:], value)
}

func (page *Page) getUint64(offset uint) uint64 {
	return binary.LittleEndian.Uint64(page.buffer[offset:])
}

func (page *Page) setInt(offset uint, value int) {
	copy(page.buffer[offset:], convert.IntToBytes[int](value))
}

func (page *Page) getInt(offset uint) int {
	return convert.BytesToInt[int](page.buffer[offset : offset+intSize])
}

func (page *Page) setUint(offset uint, value uint) {
	copy(page.buffer[offset:], convert.IntToBytes[uint](value))
}

func (page *Page) getUint(offset uint) uint {
	return convert.BytesToInt[uint](page.buffer[offset : offset+intSize])
}

func (page *Page) setFloat32(offset uint, value float32) {
	binary.LittleEndian.PutUint32(page.buffer[offset:], math.Float32bits(value))
}

func (page *Page) getFloat32(offset uint) float32 {
	return math.Float32frombits(page.getUint32(offset))
}

func (page *Page) setFloat64(offset uint, value float64) {
	binary.LittleEndian.PutUint64(page.buffer[offset:], math.Float64bits(value))
}

func (page *Page) getFloat64(offset uint) float64 {
	return math.Float64frombits(page.getUint64(offset))
}

func (page *Page) setBytes(offset uint, buffer []byte) {
	//TODO: Buffer size should be less than 2^32-1 (close to 4G)
	binary.LittleEndian.PutUint32(page.buffer[offset:], uint32(len(buffer)))
	copy(page.buffer[offset+reservedBufferSize:], buffer)
}

func (page *Page) getBytes(offset uint) []byte {
	bufferLength := binary.LittleEndian.Uint32(page.buffer[offset:])
	endOffset := offset + reservedBufferSize + uint(bufferLength)
	return page.buffer[offset+reservedBufferSize : endOffset]
}

func (page *Page) setString(offset uint, str string) {
	//TODO: str size should be less than 2^32-1 (close to 4G)
	page.setBytes(offset, []byte(str))
}

func (page *Page) getString(offset uint) string {
	return string(page.getBytes(offset))
}
