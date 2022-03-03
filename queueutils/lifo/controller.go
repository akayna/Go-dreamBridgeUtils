package lifo

import (
	"errors"

	"github.com/akayna/Go-dreamBridgeUtils/queueutils"
	"github.com/akayna/Go-dreamBridgeUtils/queueutils/roundqueue"
)

// NewLifo - Returns a new lifo with max size setted. Use maxSize = 0 to unlimited size.
func NewLifo(maxSize int) *Lifo {
	lifo := new(Lifo)
	lifo.queue = roundqueue.NewRoundQueue(maxSize)

	return lifo
}

// AddBlock - Add a new block in the end of the lifo
func (lifo *Lifo) AddBlock(newBlock *queueutils.Block) error {
	err := lifo.queue.AddPreviousBlock(newBlock)

	if err != nil {
		return errors.New("lifo.AddBlock - The queue is full")
	}

	return nil
}

// RemoveBlock - Remove the last block of the lifo
func (lifo *Lifo) RemoveBlock() *queueutils.Block {
	lifo.queue.MovePointerToPrevious()
	return lifo.queue.RemoveBlock()
}

// IsEmpty - return true if the lifo is empty
func (lifo *Lifo) IsEmpty() bool {
	return lifo.queue.IsEmpty()
}
