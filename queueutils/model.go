package queueutils

type Block struct {
	NextBlock     *Block
	PreviousBlock *Block
	Data          interface{}
}
