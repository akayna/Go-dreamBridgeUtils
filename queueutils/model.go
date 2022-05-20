package queueutils

import "sync"

type Block struct {
	mu            sync.Mutex
	id            uint
	nextBlock     *Block
	previousBlock *Block
	data          interface{}
}
