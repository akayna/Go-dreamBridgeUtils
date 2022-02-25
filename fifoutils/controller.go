package fifoutils

import (
	"errors"

	"github.com/akayna/Go-dreamBridgeUtils/lifoutils"
	"github.com/akayna/Go-dreamBridgeUtils/queueutils"
)

// NewFifo - Create, initialize and return a new queue
func NewFifo(maxSize int) *Fifo {
	newFifo := new(Fifo)

	newFifo.queue = queueutils.NewQueue()
	newFifo.maxSize = maxSize

	return newFifo
}

// AddData - Add data in the end of the fifo
func (fifo *Fifo) AddData(data interface{}) error {
	fifo.queue.Mu.Lock()
	defer fifo.queue.Mu.Unlock()

	if fifo.queue.ActualSize >= fifo.maxSize {
		return errors.New("fifoutils.AddData - The queue is full. Can add more data")
	}

	newBlock := lifoutils.NewBlock(data)

	if fifo.queue.FirstBlock == nil {
		fifo.queue.FirstBlock = newBlock
		fifo.queue.LastBlock = newBlock
	} else {
		newBlock.PreviousBlock = fifo.queue.LastBlock

		fifo.queue.LastBlock.NextBlock = newBlock
		fifo.queue.LastBlock = newBlock
	}

	fifo.queue.ActualSize++
	return nil
}

// RemoveData - Remove the block from the fifo, return the data and erased the block
func (fifo *Fifo) RemoveData() interface{} {
	fifo.queue.Mu.Lock()
	defer fifo.queue.Mu.Unlock()

	if fifo.queue.ActualSize <= 0 {
		return nil
	} else if fifo.queue.ActualSize == 1 {
		removedBlock := fifo.queue.FirstBlock

		fifo.queue.FirstBlock = nil
		fifo.queue.LastBlock = nil
		fifo.queue.ActualSize = 0

		return lifoutils.ReleaseBlock(removedBlock)
	}

	removedBlock := fifo.queue.FirstBlock

	fifo.queue.FirstBlock = fifo.queue.FirstBlock.NextBlock
	fifo.queue.FirstBlock.PreviousBlock = nil
	fifo.queue.ActualSize--

	return lifoutils.ReleaseBlock(removedBlock)
}

// ActualSize - Returns the actual size of the queue
func (fifo *Fifo) ActualSize() int {
	return fifo.queue.ActualSize
}

// IsFull - Return if the fifo is full or not
func (fifo *Fifo) IsFull() bool {
	return fifo.ActualSize() >= fifo.maxSize
}
