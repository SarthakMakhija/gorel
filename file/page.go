package file

import (
	"encoding/binary"
	"gorel/convert"
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

func (page *Page) setInt(offset uint, value int) {
	copy(page.buffer[offset:], convert.IntToBytes(value))
}

func (page *Page) getInt(offset uint) int {
	return convert.BytesToInt(page.buffer[offset : offset+intSize])
}

func (page *Page) setUint32(offset uint, value uint32) {
	binary.LittleEndian.PutUint32(page.buffer[offset:], value)
}

func (page *Page) getUint32(offset uint) uint32 {
	return binary.LittleEndian.Uint32(page.buffer[offset:])
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
