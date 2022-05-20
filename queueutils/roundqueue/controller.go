package roundqueue

import (
	"errors"

	"github.com/akayna/Go-dreamBridgeUtils/queueutils"
)

// NewRoundQueue - Return a new queue with max size setted. Use maxSize = 0 to unlimited size.
func NewRoundQueue(maxSize int) *RoundQueue {
	newRoundQueue := new(RoundQueue)

	newRoundQueue.pointer = nil
	newRoundQueue.size = 0
	newRoundQueue.maxSize = maxSize

	return newRoundQueue
}

// IsFull - Return true if the queue is full
func (roundQueue *RoundQueue) IsFull() bool {
	if roundQueue.maxSize == 0 {
		return false
	}
	return roundQueue.size >= roundQueue.maxSize
}

// IsEmpty - Return true if the queue is empty
func (roundQueue *RoundQueue) IsEmpty() bool {
	return roundQueue.size == 0
}

// Return the queue actual size
func (roundQueue *RoundQueue) Size() int {
	return roundQueue.size
}

// Returns true if the queue is not empty.
func (roundQueue *RoundQueue) HasData() bool {
	return roundQueue.size > 0
}

// AddNextBlock - Add a new block in the next position related to the pointer.
func (roundQueue *RoundQueue) AddNextBlock(newBlock *queueutils.Block) error {
	roundQueue.mu.Lock()
	defer roundQueue.mu.Unlock()

	if roundQueue.IsFull() {
		return errors.New("roundqueue.AddNextBlock - The queue is full")
	}

	if roundQueue.pointer == nil {
		roundQueue.pointer = newBlock
		roundQueue.pointer.SetNextBlock(newBlock)
		roundQueue.pointer.SetPreviousBlock(newBlock)
	} else {
		newBlock.SetPreviousBlock(roundQueue.pointer)
		newBlock.SetNextBlock(roundQueue.pointer.GetNextBlock())

		roundQueue.pointer.GetNextBlock().SetPreviousBlock(newBlock)
		roundQueue.pointer.SetNextBlock(newBlock)
	}

	roundQueue.size++

	return nil
}

// AddPreviousBlock - Add a new block in the previus position related to the pointer.
func (roundQueue *RoundQueue) AddPreviousBlock(newBlock *queueutils.Block) error {
	roundQueue.mu.Lock()
	defer roundQueue.mu.Unlock()

	if roundQueue.IsFull() {
		return errors.New("roundqueue.AddPreviousBlock - The queue is full")
	}

	if roundQueue.pointer == nil {
		roundQueue.pointer = newBlock
		roundQueue.pointer.SetNextBlock(newBlock)
		roundQueue.pointer.SetPreviousBlock(newBlock)
	} else {
		newBlock.SetPreviousBlock(roundQueue.pointer.GetPreviousBlock())
		newBlock.SetNextBlock(roundQueue.pointer)

		roundQueue.pointer.GetPreviousBlock().SetNextBlock(newBlock)
		roundQueue.pointer.SetPreviousBlock(newBlock)
	}

	roundQueue.size++

	return nil
}

// RemoveBlock - Remove and return the block pointed by the pointer and set pointer to the next block.
func (roundQueue *RoundQueue) RemoveBlock() *queueutils.Block {
	roundQueue.mu.Lock()
	defer roundQueue.mu.Unlock()

	if roundQueue.IsEmpty() {
		return nil
	}

	removedBlock := roundQueue.pointer

	roundQueue.pointer.GetPreviousBlock().SetNextBlock(roundQueue.pointer.GetNextBlock())
	roundQueue.pointer.GetNextBlock().SetPreviousBlock(roundQueue.pointer.GetPreviousBlock())

	roundQueue.pointer = roundQueue.pointer.GetNextBlock()

	removedBlock.SetNextBlock(nil)
	removedBlock.SetPreviousBlock(nil)

	roundQueue.size--

	if roundQueue.IsEmpty() {
		roundQueue.pointer = nil
	}

	return removedBlock
}

// MovePointerToNext - Move the pointer to the next block
func (roundQueue *RoundQueue) MovePointerToNext() {
	if roundQueue.IsEmpty() {
		return
	}

	roundQueue.pointer = roundQueue.pointer.GetNextBlock()
}

// MovePointerToPrevious - Move the pointer to the previous block
func (roundQueue *RoundQueue) MovePointerToPrevious() {
	if roundQueue.IsEmpty() {
		return
	}

	roundQueue.pointer = roundQueue.pointer.GetPreviousBlock()
}

// ReadData - Returns the data stored in the actual pointed block
func (roundQueue *RoundQueue) ReadData() interface{} {
	if roundQueue.pointer == nil {
		return nil
	}
	return roundQueue.pointer.GetData()
}

// GetBlock - Returns the actual pointed block pointed by the pointer
func (roundQueue *RoundQueue) GetBlock() *queueutils.Block {
	return roundQueue.pointer
}

// SetData - Sets the data of the actual pointed block
func (roundQueue *RoundQueue) SetData(data interface{}) error {
	roundQueue.mu.Lock()
	defer roundQueue.mu.Unlock()

	if roundQueue.pointer == nil {
		return errors.New("roundqueue.SetPointerData - The queue is empty")
	}

	roundQueue.pointer.SetData(data)

	return nil
}

// Position the pointer to the block with the specified id and returns it.
func (roundQueue *RoundQueue) GetBlockId(id uint) *queueutils.Block {
	if roundQueue.IsEmpty() {
		return nil
	}

	initialBlock := roundQueue.pointer

	if roundQueue.pointer.GetId() == id {
		return roundQueue.pointer
	}

	roundQueue.MovePointerToNext()

	for initialBlock != roundQueue.pointer {
		if roundQueue.pointer.GetId() == id {
			return roundQueue.pointer
		}
		roundQueue.MovePointerToNext()
	}

	return nil
}
