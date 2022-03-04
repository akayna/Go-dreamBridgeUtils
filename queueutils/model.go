package queueutils

import "sync"

type Block struct {
	mu            sync.Mutex
	nextBlock     *Block
	previousBlock *Block
	data          interface{}
}
