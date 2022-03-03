package roundqueue

import (
	"sync"

	"github.com/akayna/Go-dreamBridgeUtils/queueutils"
)

type RoundQueue struct {
	mu         sync.Mutex
	pointer    *queueutils.Block
	actualSize int
	maxSize    int
}
