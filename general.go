package synctracer

import (
	"fmt"
	"path/filepath"
	"runtime"
	"time"
)

var (
	PrintAfter = time.Millisecond * 200
	PrintEvery = time.Second * 5
)

func now() int64 {
	return time.Now().UnixNano()
}

func callerPath(skip int) string {
	fn, file, line, ok := caller(3 + skip)
	if !ok {
		return "unknown"
	}

	return fmt.Sprintf("%s (%s:%d)", fn, filepath.Base(file), line)
}

func caller(skip int) (fn, file string, line int, ok bool) {
	rpc := make([]uintptr, 1)
	n := runtime.Callers(skip+1, rpc[:])
	if n < 1 {
		return
	}
	frame, _ := runtime.CallersFrames(rpc).Next()
	return frame.Function, frame.File, frame.Line, frame.PC != 0
}
