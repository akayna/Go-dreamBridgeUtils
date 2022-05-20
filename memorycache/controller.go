package memorycache

import (
	"log"

	"github.com/akayna/Go-dreamBridgeUtils/queueutils/blockbucket"
	"github.com/akayna/Go-dreamBridgeUtils/queueutils/roundqueue"
)

// NewMemoryCache - Returns a new memory cache manager
func NewMemoryCache(maxSize int) *MemoryCache {
	var newMemoryCache MemoryCache

	newMemoryCache.maxTypeCacheSize = maxSize

	newMemoryCache.DataType = make(map[string]*roundqueue.RoundQueue)

	return &newMemoryCache
}

// Store data in memory, associated with a data type and returns the data id.
func (cache *MemoryCache) StoreData(dataType string, data interface{}) (uint, error) {
	// Verifies if the queue of the type already exists
	if cache.DataType[dataType] == nil {
		cache.DataType[dataType] = roundqueue.NewRoundQueue(cache.maxTypeCacheSize)
	}

	// Create a new blobk and stores it in the queue
	newBlock := blockbucket.GetNewBlock(data)

	err := cache.DataType[dataType].AddNextBlock(newBlock)
	if err != nil {
		log.Println("memorycache.StoreData - Error adding new block to queue.")
		return 0, err
	}

	return newBlock.GetId(), nil
}

// Get data from memory cache by its type and id. This function removes the data from memory.
func (cache *MemoryCache) GetData(dataType string, id uint) interface{} {
	// No data of this type.
	if cache.DataType[dataType] == nil {
		return nil
	}

	// Search for the data with specified id
	block := cache.DataType[dataType].GetBlockId(id)

	if block == nil {
		return nil
	}

	cache.DataType[dataType].RemoveBlock()

	return blockbucket.StoreBlock(block)
}

// Read data from memory cache by its type and id. This function doesn´t remove the data from memory.
func (cache *MemoryCache) ReadData(dataType string, id uint) interface{} {
	// No data of this type.
	if cache.DataType[dataType] == nil {
		return nil
	}

	// Search for the data with specified id
	block := cache.DataType[dataType].GetBlockId(id)

	if block == nil {
		return nil
	}

	return block.GetData()
}
