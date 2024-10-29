package log

import (
	"github.com/stretchr/testify/assert"
	"gorel/file"
	"os"
	"testing"
)

func TestAppendARecordInLogManager(t *testing.T) {
	fileManager, err := file.NewBlockFileManager(".", blockSize)
	assert.Nil(t, err)

	defer func() {
		fileManager.Close()
		_ = os.Remove(t.Name())
	}()

	fileName := t.Name()
	logManager, err := NewBlockLogManager(fileManager, fileName)

	assert.Nil(t, err)
	assert.Nil(t, logManager.Append([]byte("RocksDB is an LSM-based storage engine")))
}

func TestAppendARecordInLogManagerAndIterateOverIt(t *testing.T) {
	fileManager, err := file.NewBlockFileManager(".", blockSize)
	assert.Nil(t, err)

	defer func() {
		fileManager.Close()
		_ = os.Remove(t.Name())
	}()

	fileName := t.Name()
	logManager, err := NewBlockLogManager(fileManager, fileName)

	assert.Nil(t, err)
	assert.Nil(t, logManager.Append([]byte("RocksDB is an LSM-based storage engine")))

	iterator, err := logManager.BackwardIterator()
	assert.Nil(t, err)

	assert.True(t, iterator.IsValid())
	assert.Equal(t, "RocksDB is an LSM-based storage engine", string(iterator.Record()))

	assert.Nil(t, iterator.Previous())
	assert.False(t, iterator.IsValid())
}

func TestAppendAFewRecordsInLogManagerAndIterateOverThem(t *testing.T) {
	fileManager, err := file.NewBlockFileManager(".", blockSize)
	assert.Nil(t, err)

	defer func() {
		fileManager.Close()
		_ = os.Remove(t.Name())
	}()

	fileName := t.Name()
	logManager, err := NewBlockLogManager(fileManager, fileName)

	assert.Nil(t, err)
	assert.Nil(t, logManager.Append([]byte("RocksDB is an LSM-based storage engine")))
	assert.Nil(t, logManager.Append([]byte("PebbleDB is an LSM-based storage engine")))
	assert.Nil(t, logManager.Append([]byte("BoltDB is a B+Tree storage engine")))

	iterator, err := logManager.BackwardIterator()
	assert.Nil(t, err)

	assert.True(t, iterator.IsValid())
	assert.Equal(t, "BoltDB is a B+Tree storage engine", string(iterator.Record()))

	assert.Nil(t, iterator.Previous())

	assert.True(t, iterator.IsValid())
	assert.Equal(t, "PebbleDB is an LSM-based storage engine", string(iterator.Record()))

	assert.Nil(t, iterator.Previous())

	assert.True(t, iterator.IsValid())
	assert.Equal(t, "RocksDB is an LSM-based storage engine", string(iterator.Record()))

	assert.Nil(t, iterator.Previous())
	assert.False(t, iterator.IsValid())
}

func TestAppendAFewRecordsInLogManagerWithSmallerBlockSizeAndIterateOverThem(t *testing.T) {
	const blockSizeInBytes = 150

	fileManager, err := file.NewBlockFileManager(".", blockSizeInBytes)
	assert.Nil(t, err)

	defer func() {
		fileManager.Close()
		_ = os.Remove(t.Name())
	}()

	fileName := t.Name()
	logManager, err := NewBlockLogManager(fileManager, fileName)

	assert.Nil(t, err)
	assert.Nil(t, logManager.Append([]byte("RocksDB is an LSM-based storage engine")))
	assert.Nil(t, logManager.Append([]byte("PebbleDB is an LSM-based storage engine")))
	assert.Nil(t, logManager.Append([]byte("BoltDB is a B+Tree storage engine")))

	iterator, err := logManager.BackwardIterator()
	assert.Nil(t, err)

	assert.True(t, iterator.IsValid())
	assert.Equal(t, "BoltDB is a B+Tree storage engine", string(iterator.Record()))

	assert.Nil(t, iterator.Previous())

	assert.True(t, iterator.IsValid())
	assert.Equal(t, "PebbleDB is an LSM-based storage engine", string(iterator.Record()))

	assert.Nil(t, iterator.Previous())

	assert.True(t, iterator.IsValid())
	assert.Equal(t, "RocksDB is an LSM-based storage engine", string(iterator.Record()))

	assert.Nil(t, iterator.Previous())
	assert.False(t, iterator.IsValid())
}

func TestAppendAFewRecordsInLogManagerAndRecreatesLogManagerInstanceToSimulateRestart(t *testing.T) {
	fileManager, err := file.NewBlockFileManager(".", blockSize)
	assert.Nil(t, err)

	defer func() {
		fileManager.Close()
		_ = os.Remove(t.Name())
	}()

	fileName := t.Name()
	logManager, err := NewBlockLogManager(fileManager, fileName)

	assert.Nil(t, err)
	assert.Nil(t, logManager.Append([]byte("RocksDB is an LSM-based storage engine")))
	assert.Nil(t, logManager.Append([]byte("PebbleDB is an LSM-based storage engine")))
	assert.Nil(t, logManager.forceFlush())

	reloadedLogManager, err := NewBlockLogManager(fileManager, fileName)
	assert.Nil(t, err)

	assert.Nil(t, reloadedLogManager.Append([]byte("BoltDB is a B+Tree storage engine")))

	iterator, err := reloadedLogManager.BackwardIterator()
	assert.Nil(t, err)

	assert.True(t, iterator.IsValid())
	assert.Equal(t, "BoltDB is a B+Tree storage engine", string(iterator.Record()))

	assert.Nil(t, iterator.Previous())

	assert.True(t, iterator.IsValid())
	assert.Equal(t, "PebbleDB is an LSM-based storage engine", string(iterator.Record()))

	assert.Nil(t, iterator.Previous())

	assert.True(t, iterator.IsValid())
	assert.Equal(t, "RocksDB is an LSM-based storage engine", string(iterator.Record()))

	assert.Nil(t, iterator.Previous())
	assert.False(t, iterator.IsValid())
}
