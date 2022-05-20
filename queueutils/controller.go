package queueutils

var blockIdCounter = uint(1)

// NewBlock - Create or recover a new block
func NewBlock(data interface{}) *Block {
	newBlock := new(Block)

	newBlock.data = data
	newBlock.id = blockIdCounter
	blockIdCounter++

	return newBlock
}

// CleanBlock - Erase all pointers of a block returning it's data
func (block *Block) Clean() interface{} {
	block.mu.Lock()
	defer block.mu.Unlock()

	data := block.data

	block.previousBlock = nil
	block.nextBlock = nil
	block.data = nil

	return data
}

// SetData - Set the block data
func (block *Block) SetData(data interface{}) {
	block.mu.Lock()
	defer block.mu.Unlock()

	block.data = data
}

// GetData - Return the data stored in the block
func (block *Block) GetData() interface{} {
	block.mu.Lock()
	defer block.mu.Unlock()

	return block.data
}

// SetPreviousBlock - Set the previous block
func (block *Block) SetPreviousBlock(previousBlock *Block) {
	block.mu.Lock()
	defer block.mu.Unlock()

	block.previousBlock = previousBlock
}

// GetPreviousBlock - Return the previous block
func (block *Block) GetPreviousBlock() *Block {
	block.mu.Lock()
	defer block.mu.Unlock()

	return block.previousBlock
}

// SetNextBlock - Set the next block
func (block *Block) SetNextBlock(nextBlock *Block) {
	block.mu.Lock()
	defer block.mu.Unlock()

	block.nextBlock = nextBlock
}

// GetNextBlock - Return the next block
func (block *Block) GetNextBlock() *Block {
	block.mu.Lock()
	defer block.mu.Unlock()

	return block.nextBlock
}

// GetId - Returns the block ID
func (block *Block) GetId() uint {
	block.mu.Lock()
	defer block.mu.Unlock()

	return block.id
}
