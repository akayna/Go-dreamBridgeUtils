package fifoutils

import "github.com/akayna/Go-dreamBridgeUtils/queueutils"

type Fifo struct {
	maxSize int
	queue   *queueutils.Queue
}
