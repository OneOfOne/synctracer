package synctracer

import (
	"sync"
	"testing"
	"time"
)

func TestSlow(t *testing.T) {
	var (
		mu RWMutex
		wg sync.WaitGroup
	)
	timeout = time.Microsecond * 100
	wg.Add(2)
	mu.Lock()
	go func() {
		defer wg.Done()
		time.Sleep(timeout * 20)
		mu.Unlock()
	}()
	go doOther(&wg, &mu)
	go doOther2(&wg, &mu)

	wg.Wait()
}

func doOther(wg *sync.WaitGroup, mu *RWMutex) {
	defer wg.Done()
	mu.Lock()
	time.Sleep(timeout * 10)
	mu.Unlock()
}

func doOther2(wg *sync.WaitGroup, mu *RWMutex) {
	defer wg.Done()
	mu.Lock()
	time.Sleep(timeout * 10)
	mu.Unlock()
}
