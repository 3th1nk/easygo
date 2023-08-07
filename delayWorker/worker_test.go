package delayWorker

import (
	"github.com/stretchr/testify/assert"
	"math/rand"
	"sync"
	"testing"
	"time"
)

func TestSerialWorker(t *testing.T) {
	handler := func(jobs []interface{}) {
		t.Log(jobs)
	}

	w := New("test", handler)
	defer w.Close()

	err := w.WithQueueSize(50).WithDelaySec(2).WithDelayCnt(10).Run()
	assert.NoError(t, err)

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		func() {
			wg.Add(1)
			defer wg.Done()

			w.Push(i)
			time.Sleep(time.Duration(rand.Intn(2)) * time.Second)
		}()
	}
	wg.Wait()
}

func TestParallelWorker(t *testing.T) {
	handler := func(jobs []interface{}) {
		t.Log(jobs)
	}

	w := New("test", handler)
	defer w.Close()

	err := w.WithQueueSize(50).WithDelaySec(2).WithDelayCnt(10).WithParallel(true).Run()
	assert.NoError(t, err)

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		j := i
		go func() {
			wg.Add(1)
			defer wg.Done()

			w.Push(j)
		}()
	}
	wg.Wait()
}
