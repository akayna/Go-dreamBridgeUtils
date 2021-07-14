package synchcounter

import "sync"

type SynchCounter struct {
	mu   sync.Mutex
	cont int64
}

func Soma(c *SynchCounter) {
	c.mu.Lock()
	c.cont++
	c.mu.Unlock()
}

func Subtrai(c *SynchCounter) {
	c.mu.Lock()
	c.cont--
	c.mu.Unlock()
}

func Value(c *SynchCounter) int64 {
	c.mu.Lock()
	x := c.cont
	c.mu.Unlock()
	return x
}
