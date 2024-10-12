package file

type BlockId struct {
	fileName    string
	blockNumber uint
}

func NewBlockId(fileName string, blockNumber uint) BlockId {
	return BlockId{
		fileName:    fileName,
		blockNumber: blockNumber,
	}
}

func (blockId BlockId) Offset(blockSize uint) int64 {
	return int64(blockId.blockNumber * blockSize)
}

func (blockId BlockId) FileName() string {
	return blockId.fileName
}

func (blockId BlockId) BlockNumber() uint {
	return blockId.blockNumber
}
