package lazyWorker

import (
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
	defer w.Stop()

	w.WithQueueSize(50).WithLazyInterval(2 * time.Second).WithLazySize(10).Run()

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
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
	defer w.Stop()

	w.WithQueueSize(50).WithLazyInterval(2 * time.Second).WithLazySize(10).Debug().Run()

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		j := i
		func() {
			wg.Add(1)
			defer wg.Done()

			w.Push(j)
		}()
	}
	wg.Wait()
}

func TestRestart(t *testing.T) {
	handler := func(jobs []interface{}) {
		t.Log(jobs)
	}

	var wg sync.WaitGroup
	w := New("test", handler)

	w.WithQueueSize(50).WithLazyInterval(2 * time.Second).WithLazySize(10).Run()
	for i := 0; i < 10; i++ {
		func() {
			wg.Add(1)
			defer wg.Done()

			w.Push(i)
			time.Sleep(time.Duration(rand.Intn(2)) * time.Second)
		}()
	}
	wg.Wait()
	w.Stop()

	t.Log("restart")
	w.WithQueueSize(50).WithLazyInterval(2 * time.Second).WithLazySize(10).Run()
	for i := 0; i < 100; i++ {
		j := i
		func() {
			wg.Add(1)
			defer wg.Done()

			w.Push(j)
		}()
	}
	wg.Wait()
	w.Stop()
}
