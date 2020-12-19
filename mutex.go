package synctracer // import "go.oneofone.dev/synctracer"

import (
	"fmt"
	"log"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

type RWMutex struct {
	mu sync.RWMutex

	tsLock    int64
	tsRLock   int64
	lastLock  atomic.Value
	lastRLock atomic.Value
}

func (m *RWMutex) Lock() {
	if PrintAfter == 0 {
		m.mu.Lock()
		return
	}
	m.waitForLock()
	m.mu.Lock()
	m.lastLock.Store(callerPath(0))
	atomic.StoreInt64(&m.tsLock, now())
}

func (m *RWMutex) Unlock() {
	if PrintAfter == 0 {
		m.mu.Unlock()
		return
	}
	m.mu.Unlock()
	m.lastLock.Store("")
	atomic.StoreInt64(&m.tsLock, 0)
}

func (m *RWMutex) RLock() {
	if PrintAfter == 0 {
		m.mu.RLock()
		return
	}
	m.waitForLock()
	m.mu.RLock()
	m.lastRLock.Store(callerPath(0))
	atomic.StoreInt64(&m.tsRLock, now())
}

func (m *RWMutex) RUnlock() {
	if PrintAfter == 0 {
		m.mu.RUnlock()
		return
	}

	m.mu.RUnlock()
	m.lastRLock.Store("")
	atomic.StoreInt64(&m.tsRLock, 0)
}

func (m *RWMutex) waitForLock() {
	lastPrint := int64(0)
	for PrintAfter > 0 {
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
		runtime.Gosched()
	}
}

type Mutex struct {
	mu sync.Mutex

	tsLock   int64
	lastLock atomic.Value
}

func (m *Mutex) Lock() {
	if PrintAfter == 0 {
		m.mu.Lock()
		return
	}

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
	for PrintAfter > 0 {
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
		runtime.Gosched()
	}
}

type WaitGroup = sync.WaitGroup
