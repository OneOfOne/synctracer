package synctracer // import "go.oneofone.dev/synctracer"

import (
	"fmt"
	"log"
	"sync"
	"sync/atomic"
	"time"
)

var timeout = time.Millisecond * 500

func SetTimeout(to time.Duration) time.Duration {
	old := atomic.SwapInt64((*int64)(&timeout), int64(to))
	return time.Duration(old)
}

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
	printedLock, printedRLock := false, false
	for !printedLock && !printedRLock {
		if ts := atomic.LoadInt64(&m.tsLock); ts != 0 && !printedLock {
			if diff := time.Duration(now() - ts); diff > timeout {
				if fn, _ := m.lastLock.Load().(string); fn != "" {
					log.Output(1, fmt.Sprintf("[%s] LOCK STUCK (%v) @ %s", callerPath(1), diff, fn))
					printedLock = true
				}
			}
		} else if ts := atomic.LoadInt64(&m.tsRLock); ts != 0 && !printedRLock {
			if diff := time.Duration(now() - ts); diff > timeout {
				if fn, _ := m.lastLock.Load().(string); fn != "" {
					log.Output(1, fmt.Sprintf("[%s] RLOCK STUCK (%v) @ %s", callerPath(1), diff, fn))
					printedRLock = true
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
	for {
		if ts := atomic.LoadInt64(&m.tsLock); ts != 0 {
			if diff := time.Duration(now() - ts); diff > timeout {
				if fn, _ := m.lastLock.Load().(string); fn != "" {
					log.Output(1, fmt.Sprintf("[%s] LOCK STUCK (%v) @ %s", callerPath(1), diff, fn))
					break
				}
			}
		} else {
			break
		}
		time.Sleep(time.Millisecond)
	}
}

type WaitGroup = sync.WaitGroup
