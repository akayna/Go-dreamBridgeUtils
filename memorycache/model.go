package memorycache

import (
	"github.com/akayna/Go-dreamBridgeUtils/queueutils/roundqueue"
)

// MemoryCache - Struct to manage memory cache
type MemoryCache struct {
	DataType         map[string]*roundqueue.RoundQueue
	maxTypeCacheSize int
}
