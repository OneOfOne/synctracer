package synctracer

import (
	"sync"
	"testing"
	"time"
)

func TestRWMutex(t *testing.T) {
	var (
		mu RWMutex
		wg sync.WaitGroup
	)
	PrintAfter = time.Microsecond * 100
	wg.Add(3)
	mu.Lock()
	go func() {
		defer wg.Done()
		time.Sleep(time.Second * 2)
		mu.Unlock()
	}()
	go doOther(&wg, &mu)
	go doOther2(&wg, &mu)

	wg.Wait()
}

func TestSlowCall(t *testing.T) {
	SlowCall(func() {
		time.Sleep(time.Second)
		t.Log("done")
	}, 1, "blah", t)
}

func doOther(wg *sync.WaitGroup, mu *RWMutex) {
	defer wg.Done()
	mu.Lock()
	time.Sleep(PrintAfter * 10)
	mu.Unlock()
}

func doOther2(wg *sync.WaitGroup, mu *RWMutex) {
	defer wg.Done()
	mu.Lock()
	time.Sleep(PrintAfter * 10)
	mu.Unlock()
}
