package synctracer

import (
	"fmt"
	"log"
	"sync/atomic"
	"time"
)

type wrap struct {
	start int64
	args  string
	ch    chan struct{}
}

func SlowCall(fn func(), args ...interface{}) {
	start, cp, done := now(), callerPath(0), int64(0)
	argsStr := fmt.Sprintf("%#+v", args)[len("[]interface {}"):]
	defer atomic.StoreInt64(&done, 1)
	go func() {
		for atomic.LoadInt64(&done) == 0 {
			if t := time.Duration(now() - start); t > PrintAfter {
				log.Printf("FUNC STUCK %v @ %s, args: %s", t, cp, argsStr)
				time.Sleep(PrintEvery)
			} else {
				time.Sleep(PrintAfter)
			}
		}
	}()
	fn()
}
