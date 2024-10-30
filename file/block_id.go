package file

type BlockId struct {
	fileName    string
	blockNumber uint
}

const missingFileName = "FILE@NONE@"

var MissingBlockId = BlockId{fileName: missingFileName}

func NewBlockId(fileName string, blockNumber uint) BlockId {
	return BlockId{
		fileName:    fileName,
		blockNumber: blockNumber,
	}
}

func (blockId BlockId) offset(blockSize uint) int64 {
	return int64(blockId.blockNumber * blockSize)
}

func (blockId BlockId) BlockNumber() uint {
	return blockId.blockNumber
}

func (blockId BlockId) Previous() BlockId {
	return NewBlockId(blockId.fileName, blockId.blockNumber-1)
}

func (blockId BlockId) IsMissing() bool {
	return blockId.fileName == missingFileName
}
