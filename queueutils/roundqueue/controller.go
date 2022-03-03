package roundqueue

import (
	"errors"

	"github.com/akayna/Go-dreamBridgeUtils/queueutils"
)

// NewRoundQueue - Return a new queue with max size setted. Use maxSize = 0 to unlimited size.
func NewRoundQueue(maxSize int) *RoundQueue {
	newRoundQueue := new(RoundQueue)

	newRoundQueue.pointer = nil
	newRoundQueue.actualSize = 0
	newRoundQueue.maxSize = maxSize

	return newRoundQueue
}

// IsFull - Return true if the queue is full
func (roundQueue *RoundQueue) IsFull() bool {
	if roundQueue.maxSize == 0 {
		return false
	}
	return roundQueue.actualSize >= roundQueue.maxSize
}

// IsEmpty - Return true if the queue is empty
func (roundQueue *RoundQueue) IsEmpty() bool {
	return roundQueue.actualSize == 0
}

// ActualSize - Return the queue actual size
func (roundQueue *RoundQueue) ActualSize() int {
	return roundQueue.actualSize
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
		roundQueue.pointer.NextBlock = newBlock
		roundQueue.pointer.PreviousBlock = newBlock
	} else {
		newBlock.PreviousBlock = roundQueue.pointer
		newBlock.NextBlock = roundQueue.pointer.NextBlock

		roundQueue.pointer.NextBlock.PreviousBlock = newBlock
		roundQueue.pointer.NextBlock = newBlock
	}

	roundQueue.actualSize++

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
		roundQueue.pointer.NextBlock = newBlock
		roundQueue.pointer.PreviousBlock = newBlock
	} else {
		newBlock.PreviousBlock = roundQueue.pointer.PreviousBlock
		newBlock.NextBlock = roundQueue.pointer

		roundQueue.pointer.PreviousBlock.NextBlock = newBlock
		roundQueue.pointer.PreviousBlock = newBlock
	}

	roundQueue.actualSize++

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

	roundQueue.pointer.PreviousBlock.NextBlock = roundQueue.pointer.NextBlock
	roundQueue.pointer.NextBlock.PreviousBlock = roundQueue.pointer.PreviousBlock

	roundQueue.pointer = roundQueue.pointer.NextBlock

	removedBlock.NextBlock = nil
	removedBlock.PreviousBlock = nil

	roundQueue.actualSize--

	return removedBlock
}

// MovePointerToNext - Move the pointer to the next block
func (roundQueue *RoundQueue) MovePointerToNext() {
	if roundQueue.IsEmpty() {
		return
	}

	roundQueue.pointer = roundQueue.pointer.NextBlock
}

// MovePointerToPrevious - Move the pointer to the previous block
func (roundQueue *RoundQueue) MovePointerToPrevious() {
	if roundQueue.IsEmpty() {
		return
	}

	roundQueue.pointer = roundQueue.pointer.PreviousBlock
}

// GetPointerData - Returns the data stored in the block pointed by the pointer
func (roundQueue *RoundQueue) GetPointerData() interface{} {
	if roundQueue.pointer == nil {
		return nil
	}
	return roundQueue.pointer.Data
}

// SetPointerData - Sets the data of the block pointed by the pointer
func (roundQueue *RoundQueue) SetPointerData(data interface{}) error {
	roundQueue.mu.Lock()
	defer roundQueue.mu.Unlock()

	if roundQueue.pointer == nil {
		return errors.New("roundqueue.SetPointerData - The queue is empty")
	}

	roundQueue.pointer.Data = data

	return nil
}
