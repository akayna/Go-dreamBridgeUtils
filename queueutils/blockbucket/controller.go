package blockbucket

import (
	"github.com/akayna/Go-dreamBridgeUtils/queueutils"
	"github.com/akayna/Go-dreamBridgeUtils/queueutils/roundqueue"
)

// GetNewBlock - Get a new block from the bucket or generate a new one if the bucket is empty.
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

// StoreBlock - Store the new block in the bucket for future use.
func StoreBlock(block *queueutils.Block) interface{} {
	if blockBucket == nil {
		blockBucket = roundqueue.NewRoundQueue(0)
	}

	data := block.Clean()

	blockBucket.AddPreviousBlock(block)

	return data
}
