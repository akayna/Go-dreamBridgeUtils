package fifo

import (
	"errors"

	"github.com/akayna/Go-dreamBridgeUtils/queueutils"
	"github.com/akayna/Go-dreamBridgeUtils/queueutils/roundqueue"
)

// NewFifo - Returns a new fifo with max size setted. Use maxSize = 0 to unlimited size.
func NewFifo(maxSize int) *Fifo {
	fifo := new(Fifo)
	fifo.queue = roundqueue.NewRoundQueue(maxSize)

	return fifo
}

// AddBlock - Add a new block in the end of the queue
func (fifo *Fifo) AddBlock(newBlock *queueutils.Block) error {
	err := fifo.queue.AddPreviousBlock(newBlock)

	if err != nil {
		return errors.New("fifo.AddBlock - The queue is full")
	}

	return nil
}

// RemoveBlock - Remove the first block of the queue
func (fifo *Fifo) RemoveBlock() *queueutils.Block {
	return fifo.queue.RemoveBlock()
}

// IsEmpty - return true if the fifo is empty
func (fifo *Fifo) IsEmpty() bool {
	return fifo.queue.IsEmpty()
}
