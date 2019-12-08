package work

import (
	"log"
	"sync"
)

type Worker interface {
	Task()
}

type Pool struct {
	work chan Worker
	wg   sync.WaitGroup
}

// New create the pool and limit to max goroutine number.
func New(maxGoroutines int) *Pool {
	p := Pool{
		work: make(chan Worker),
	}

	p.wg.Add(maxGoroutines)
	for i := 0; i < maxGoroutines; i++ {
		go func(i int) {
			for w := range p.work {
				log.Printf("goroutine: %d", i)
				w.Task()
			}
			p.wg.Done()
		}(i)
	}

	return &p
}

// Run post new worker to the pool.
func (p *Pool) Run(w Worker) {
	p.work <- w
}

// Shutdown wait to stop for all goroutine.
func (p *Pool) Shutdown() {
	close(p.work)
	p.wg.Wait()
}
