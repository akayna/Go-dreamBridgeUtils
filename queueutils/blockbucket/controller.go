package blockbucket

import (
	"github.com/akayna/Go-dreamBridgeUtils/queueutils"
	"github.com/akayna/Go-dreamBridgeUtils/queueutils/roundqueue"
)

// Get a new block from the bucket or generate a new one if the bucket is empty.
func GetNewBlock(data interface{}) *queueutils.Block {
	var newBlock *queueutils.Block

	if (blockBucket == nil) || (blockBucket.IsEmpty()) {
		newBlock = queueutils.NewBlock(data)
	} else {
		newBlock = blockBucket.RemoveBlock()
		newBlock.SetData(data)
	}

	return newBlock
}

// Store the new block in the bucket for future use returning its data.
func StoreBlock(block *queueutils.Block) interface{} {
	if block == nil {
		return nil
	}

	if blockBucket == nil {
		blockBucket = roundqueue.NewRoundQueue(0)
	}

	data := block.Clean()

	blockBucket.AddPreviousBlock(block)

	return data
}

// Returns the actual bucket size
func BucketSize() uint {
	if blockBucket == nil {
		return 0
	}

	return uint(blockBucket.Size())
}
