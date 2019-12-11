package work

import (
	"log"
	"sync"
)

// 定义Worker接口，只能实现了该接口的方法，都可以使用协程池执行。
type Worker interface {
	Task()
}

// 定义一个协程池结构体
type Pool struct {
	work chan Worker
	wg   sync.WaitGroup
}

// New 创建一个协程池并设置协程最大运行数量
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

// Run 发送一个Worker到协程池
func (p *Pool) Run(w Worker) {
	p.work <- w
}

// Shutdown 等待和停止所有的协程.
func (p *Pool) Shutdown() {
	close(p.work)
	p.wg.Wait()
}
