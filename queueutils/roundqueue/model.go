package roundqueue

import (
	"sync"

	"github.com/akayna/Go-dreamBridgeUtils/queueutils"
)

type RoundQueue struct {
	mu      sync.Mutex
	pointer *queueutils.Block
	size    int
	maxSize int
}
