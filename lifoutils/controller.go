package lifoutils

import (
	"errors"

	"github.com/akayna/Go-dreamBridgeUtils/queueutils"
)

// NewFifo - Create, initialize and return a new queue
func NewLifo(maxSize int) *Lifo {

	newLifo := new(Lifo)

	newLifo.queue = queueutils.NewQueue()
	newLifo.maxSize = maxSize

	return newLifo
}

// ActualSize - Returns the actual size of the queue
func (lifo *Lifo) ActualSize() int {
	return lifo.queue.ActualSize
}

// IsFull - Return if the lifo is full or not
func (lifo *Lifo) IsFull() bool {
	return (lifo.ActualSize() >= lifo.maxSize)
}

// AddData - Add data in the top of the lifo
func (lifo *Lifo) AddData(data interface{}) error {

	newBlock := NewBlock(data)

	return lifo.addBlock(newBlock)
}

// addBlock - Add a block in the top of the lifo
func (lifo *Lifo) addBlock(newBlock *queueutils.Block) error {
	lifo.queue.Mu.Lock()
	defer lifo.queue.Mu.Unlock()

	if lifo.IsFull() {
		return errors.New("lifoutils.AddData - The queue is full. Can add more data")
	}

	if lifo.queue.FirstBlock == nil {
		lifo.queue.FirstBlock = newBlock
		lifo.queue.LastBlock = newBlock
	} else {
		newBlock.PreviousBlock = lifo.queue.LastBlock

		lifo.queue.LastBlock.NextBlock = newBlock
		lifo.queue.LastBlock = newBlock
	}

	lifo.queue.ActualSize++

	return nil
}

// RemoveData - Remove the last entered data from the lifo, return the data and erased the block
func (lifo *Lifo) RemoveData() interface{} {

	removedBlock := lifo.removeBlock()

	return ReleaseBlock(removedBlock)
}

// removeBlock - Remove and return the last entered block from the lifo
func (lifo *Lifo) removeBlock() *queueutils.Block {
	lifo.queue.Mu.Lock()
	defer lifo.queue.Mu.Unlock()

	if lifo.queue.ActualSize <= 0 {
		return nil
	} else if lifo.queue.ActualSize == 1 {
		removedBlock := lifo.queue.LastBlock

		lifo.queue.FirstBlock = nil
		lifo.queue.LastBlock = nil
		lifo.queue.ActualSize = 0

		return removedBlock
	}

	removedBlock := lifo.queue.LastBlock

	lifo.queue.LastBlock = lifo.queue.LastBlock.PreviousBlock
	lifo.queue.LastBlock.NextBlock = nil

	lifo.queue.ActualSize--

	return removedBlock
}

// startNewBucket - Initialize a new free blocks bucket
func startNewBucket() {
	freeBlocksBucket = NewLifo(maxBucketSize)
}

// ReleaseBlock - Release the block, returning it's data and store it in the bucket for future use.
func ReleaseBlock(block *queueutils.Block) interface{} {

	if freeBlocksBucket == nil {
		startNewBucket()
	}

	data := block.Clean()

	if !freeBlocksBucket.IsFull() {
		freeBlocksBucket.addBlock(block)
	}

	return data
}

// NewBlock - Return and empty block from the bucket or generate a new one it the bucket is empty
func NewBlock(data interface{}) *queueutils.Block {

	if freeBlocksBucket == nil {
		startNewBucket()
	}

	newBlock := freeBlocksBucket.removeBlock()

	if newBlock == nil {
		newBlock = queueutils.NewBlock(data)
	} else {
		newBlock.Data = data
	}

	return newBlock
}
