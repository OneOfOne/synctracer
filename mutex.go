package synctracer // import "go.oneofone.dev/synctracer"

import (
	"fmt"
	"log"
	"sync"
	"sync/atomic"
	"time"
)

var (
	PrintAfter = time.Millisecond * 500
	PrintEvery = time.Second * 2
)

type RWMutex struct {
	mu sync.RWMutex

	tsLock    int64
	tsRLock   int64
	lastLock  atomic.Value
	lastRLock atomic.Value
}

func (m *RWMutex) Lock() {
	m.waitForLock()
	m.mu.Lock()
	m.lastLock.Store(callerPath(0))
	atomic.StoreInt64(&m.tsLock, now())
}

func (m *RWMutex) Unlock() {
	m.mu.Unlock()
	m.lastLock.Store("")
	atomic.StoreInt64(&m.tsLock, 0)
}

func (m *RWMutex) RLock() {
	m.waitForLock()
	m.mu.RLock()
	m.lastRLock.Store(callerPath(0))
	atomic.StoreInt64(&m.tsRLock, now())
}

func (m *RWMutex) RUnlock() {
	m.mu.RUnlock()
	m.lastRLock.Store("")
	atomic.StoreInt64(&m.tsRLock, 0)
}

func (m *RWMutex) waitForLock() {
	lastPrint := int64(0)
	for {
		n := now()
		if ts := atomic.LoadInt64(&m.tsLock); ts != 0 {
			if diff := time.Duration(n - ts); diff > PrintAfter {
				if fn, _ := m.lastLock.Load().(string); fn != "" && (n-lastPrint) >= int64(PrintEvery) {
					log.Output(1, fmt.Sprintf("LOCK STUCK [%s] %v @ %s", callerPath(1), diff, fn))
					lastPrint = n
				}
			}
		} else if ts := atomic.LoadInt64(&m.tsRLock); ts != 0 && (n-lastPrint) >= int64(PrintEvery) {
			if diff := time.Duration(n - ts); diff > PrintAfter {
				if fn, _ := m.lastLock.Load().(string); fn != "" {
					log.Output(1, fmt.Sprintf("RLOCK STUCK [%s] %v @ %s", callerPath(1), diff, fn))
					lastPrint = n
				}
			}
		} else {
			break
		}
		time.Sleep(time.Millisecond)
	}
}

type Mutex struct {
	mu sync.Mutex

	tsLock   int64
	lastLock atomic.Value
}

func (m *Mutex) Lock() {
	m.waitForLock()
	m.mu.Lock()
	m.lastLock.Store(callerPath(0))
	atomic.StoreInt64(&m.tsLock, now())
}

func (m *Mutex) Unlock() {
	m.mu.Unlock()
	m.lastLock.Store("")
	atomic.StoreInt64(&m.tsLock, 0)
}

func (m *Mutex) waitForLock() {
	lastPrint := int64(0)
	for {
		n := now()
		if ts := atomic.LoadInt64(&m.tsLock); ts != 0 {
			if diff := time.Duration(n - ts); diff > PrintAfter {
				if fn, _ := m.lastLock.Load().(string); fn != "" && (n-lastPrint) >= int64(PrintEvery) {
					log.Output(1, fmt.Sprintf("LOCK STUCK [%s] %v @ %s", callerPath(1), diff, fn))
					lastPrint = n
				}
			}
		} else {
			break
		}
		time.Sleep(time.Millisecond)
	}
}

type WaitGroup = sync.WaitGroup
