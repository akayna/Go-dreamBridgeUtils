package queueutils

import (
	"sync"
)

type Block struct {
	NextBlock     *Block
	PreviousBlock *Block
	Data          interface{}
}

type Queue struct {
	Mu          sync.Mutex
	ActualSize  int
	FirstBlock  *Block
	LastBlock   *Block
	FreePointer *Block
}
