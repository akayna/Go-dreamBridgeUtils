package fifo

import (
	"github.com/akayna/Go-dreamBridgeUtils/queueutils/roundqueue"
)

type Fifo struct {
	queue *roundqueue.RoundQueue
}
