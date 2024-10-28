package file

import (
	"gorel"
	"os"
	"path/filepath"
)

// BlockFileManager
// TODO: append method
// TODO: synchronization
type BlockFileManager struct {
	dbDirectory string
	blockSize   uint
	openFiles   map[string]*os.File
}

func NewBlockFileManager(dbDirectory string, blockSize uint) (*BlockFileManager, error) {
	if _, err := os.Stat(dbDirectory); os.IsNotExist(err) {
		if err := os.MkdirAll(dbDirectory, os.ModePerm); err != nil {
			return nil, err
		}
	}
	//TODO: remove temp files
	return &BlockFileManager{
		dbDirectory: dbDirectory,
		blockSize:   blockSize,
		openFiles:   make(map[string]*os.File),
	}, nil
}

func (fileManager *BlockFileManager) ReadInto(blockId BlockId, page gorel.Page) error {
	buffer := make([]byte, fileManager.blockSize)
	err := fileManager.seekWithinFileAndRun(blockId, func(file *os.File) error {
		if _, err := file.Read(buffer); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	page.DecodePageFrom(buffer)
	return nil
}

func (fileManager *BlockFileManager) Write(blockId BlockId, page gorel.Page) error {
	return fileManager.seekWithinFileAndRun(blockId, func(file *os.File) error {
		if _, err := file.Write(page.Content()); err != nil {
			return err
		}
		return nil
	})
}

func (fileManager *BlockFileManager) Close() {
	for _, file := range fileManager.openFiles {
		if file != nil {
			_ = file.Close()
		}
	}
}

func (fileManager *BlockFileManager) BlockSize() uint {
	return fileManager.blockSize
}

func (fileManager *BlockFileManager) seekWithinFileAndRun(blockId BlockId, block func(*os.File) error) error {
	file, err := fileManager.getOrCreateFile(blockId.fileName)
	if err != nil {
		return err
	}
	if _, err := file.Seek(blockId.offset(fileManager.blockSize), 0); err != nil {
		return err
	}
	return block(file)
}

func (fileManager *BlockFileManager) getOrCreateFile(fileName string) (*os.File, error) {
	file, ok := fileManager.openFiles[fileName]
	if ok {
		return file, nil
	}
	file, err := os.OpenFile(filepath.Join(fileManager.dbDirectory, fileName), os.O_RDWR|os.O_SYNC|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}
	fileManager.openFiles[fileName] = file
	return file, nil
}
