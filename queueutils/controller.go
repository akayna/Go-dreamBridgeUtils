package queueutils

// NewBlock - Create or recover a new block
func NewBlock(data interface{}) *Block {
	newBlock := new(Block)

	newBlock.Data = data

	return newBlock
}

// CleanBlock - Erase all pointers of a block returning it's data
func (block *Block) Clean() interface{} {

	data := block.Data

	block.PreviousBlock = nil
	block.NextBlock = nil
	block.Data = nil

	return data
}
