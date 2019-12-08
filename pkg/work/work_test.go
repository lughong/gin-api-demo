package work

import (
	"log"
	"testing"
	"time"
)

type testWork struct {
	n int
}

func (tw testWork) Task() {
	log.Println(tw.n)
}

func TestRun(t *testing.T) {
	var tw testWork

	maxGoroutines := 2
	pool := New(maxGoroutines)
	t.Log(pool)
	for i := 0; i < 10; i++ {
		tw.n = i
		go pool.Run(tw)
	}

	time.Sleep(time.Second * 2)
	pool.Shutdown()
}
