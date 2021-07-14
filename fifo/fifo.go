package fifo

import "sync"

const TAMANHO_MAX_FILA = 50

type FIFO struct {
	mu       sync.Mutex
	qtdItems int64
	inicio   int64
	fim      int64
	fila     [TAMANHO_MAX_FILA]int64
}

func NovaFIFO() *FIFO {
	novaFIFO := new(FIFO)

	return novaFIFO
}

func InsereDadoFIFO(fifo *FIFO, dado int64) bool {
	fifo.mu.Lock()
	defer fifo.mu.Unlock()

	// Testa se tem espaço na fila
	if fifo.qtdItems >= TAMANHO_MAX_FILA {
		return false
	}
	// Acrescenta o dado na fila
	fifo.fila[fifo.fim] = dado
	fifo.fim++
	fifo.qtdItems++

	if fifo.fim >= TAMANHO_MAX_FILA {
		fifo.fim = 0
	}

	return true
}

func RetiraDadoFIFO(fifo *FIFO, dado *int64) bool {
	fifo.mu.Lock()
	defer fifo.mu.Unlock()

	// Verifica se tem dados na fila
	if fifo.qtdItems == 0 {
		return false
	}

	// Remove o dado da fila
	*dado = fifo.fila[fifo.inicio]
	fifo.inicio++
	fifo.qtdItems--

	if fifo.inicio >= TAMANHO_MAX_FILA {
		fifo.inicio = 0
	}

	return true
}

func QtdItensFIFO(fifo *FIFO) int64 {
	fifo.mu.Lock()
	x := fifo.qtdItems
	fifo.mu.Unlock()
	return x
}
