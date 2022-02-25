package lifoutils

import "github.com/akayna/Go-dreamBridgeUtils/queueutils"

const maxBucketSize = 100

var freeBlocksBucket *Lifo

type Lifo struct {
	maxSize int
	queue   *queueutils.Queue
}
