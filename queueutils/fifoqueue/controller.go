package fifoqueue

import (
	"errors"
	"log"

	"github.com/akayna/Go-dreamBridgeUtils/queueutils"
	"github.com/akayna/Go-dreamBridgeUtils/queueutils/blockbucket"
)

// NewFifo - Returns a new fifo with max size setted. Use maxSize = 0 to unlimited size.
func NewFifo(maxSize uint) *FifoQueue {
	fifo := new(FifoQueue)
	fifo.maxSize = maxSize

	return fifo
}

// AddData - Add a new data in the end of the queue using block bucket
func (fifo *FifoQueue) AddData(data interface{}) error {
	if fifo.IsFull() {
		return errors.New("fifoqueue.AddData - Fifo is full")
	}

	return fifo.AddBlock(blockbucket.GetNewBlock(data))
}

// AddBlock - Add a new block in the end of the queue
func (fifo *FifoQueue) AddBlock(newBlock *queueutils.Block) error {
	fifo.mu.Lock()
	defer fifo.mu.Unlock()

	if fifo.IsFull() {
		return errors.New("fifoqueue.AddBlock - Fifo is full")
	}

	newBlock.SetPreviousBlock(fifo.lastBlock)
	newBlock.SetNextBlock(nil)

	if fifo.IsEmpty() {
		fifo.firstBlock = newBlock
		fifo.freePointer = newBlock
	} else {
		fifo.lastBlock.SetNextBlock(newBlock)
	}

	fifo.lastBlock = newBlock
	fifo.size++

	return nil
}

// Removes the first block of the queue and returns its data. Stores the block for future use using block bucket.
func (fifo *FifoQueue) RemoveData() interface{} {
	if fifo.IsEmpty() {
		return nil
	}

	return blockbucket.StoreBlock(fifo.RemoveBlock())
}

// RemoveBlock - Remove the first block of the queue
func (fifo *FifoQueue) RemoveBlock() *queueutils.Block {
	fifo.mu.Lock()
	defer fifo.mu.Unlock()

	if fifo.IsEmpty() {
		return nil
	}

	removedBlock := fifo.firstBlock

	if fifo.size == 1 {
		fifo.lastBlock = nil
		fifo.freePointer = nil
	}

	fifo.firstBlock = removedBlock.GetNextBlock()
	fifo.size--

	removedBlock.SetNextBlock(nil)
	removedBlock.SetPreviousBlock(nil)

	return removedBlock
}

// IsEmpty - return true if the fifo is empty
func (fifo *FifoQueue) IsEmpty() bool {
	return fifo.size == 0
}

// IsFull - return true if the fifo is full
func (fifo *FifoQueue) IsFull() bool {
	return fifo.maxSize <= fifo.size
}

// Return the fifo actual size
func (fifo *FifoQueue) Size() uint {
	return fifo.size
}

// Return the fifo actual size
func (fifo *FifoQueue) HasData() bool {
	return fifo.size > 0
}

// Set the free pointer to the requested position and returns its block
func (fifo *FifoQueue) SetFreePointer(position uint) (*queueutils.Block, error) {
	if fifo.IsEmpty() {
		return nil, errors.New("fifoqueue.SetSearchPointer - Fifo is empty")
	}

	if position > (fifo.Size() - 1) {
		return nil, errors.New("fifoqueue.SetSearchPointer - Position greater them fifo")
	}

	fifo.freePointer = fifo.firstBlock

	for cont := uint(0); cont < position; cont++ {
		fifo.freePointer = fifo.freePointer.GetNextBlock()
	}

	return fifo.freePointer, nil
}

// Set the free pointer to the next position and returns its block. If the pointer is nil, reset it to the begining of the queue.
func (fifo *FifoQueue) SetFreePointerNext() (*queueutils.Block, error) {
	if fifo.IsEmpty() {
		return nil, errors.New("fifoqueue.SetFreePointerNext - Fifo is empty")
	}

	if fifo.freePointer != nil {
		fifo.freePointer = fifo.freePointer.GetNextBlock()
	}

	return fifo.freePointer, nil
}

// Set the free pointer to the previous position and returns its block.
func (fifo *FifoQueue) SetFreePointerPrevious() (*queueutils.Block, error) {
	if fifo.IsEmpty() {
		return nil, errors.New("fifoqueue.SetFreePointerPrevious - Fifo is empty")
	}

	if fifo.freePointer != nil {
		fifo.freePointer = fifo.freePointer.GetPreviousBlock()
	}

	return fifo.freePointer, nil
}

// Search for a block in the queue with the specified id.
func (fifo *FifoQueue) GetBlock(id uint) (*queueutils.Block, error) {
	block, err := fifo.SetFreePointer(0)

	if err != nil {
		log.Println("fifoqueue.GetBlock - Error reseting free pointer position.")
		return nil, err
	}

	for block != nil {
		if block.GetId() == id {
			return block, nil
		} else {
			block, _ = fifo.SetFreePointerNext()
		}
	}

	return nil, nil
}
