package recmutex

import (
	"bytes"
	"runtime"
	"strconv"
	"sync"
	"time"
)

func getGID() uint64 {
	b := make([]byte, 64)
	b = b[:runtime.Stack(b, false)]
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	b = b[:bytes.IndexByte(b, ' ')]
	n, _ := strconv.ParseUint(string(b), 10, 64)
	return n
}

type RecursiveMutex struct {
	mutex            sync.Mutex
	internalMutex    sync.Mutex
	currentGoRoutine uint64
	lockCount        uint64
}

func (rm *RecursiveMutex) Lock() {
	goRoutineID := getGID()

	for {
		rm.internalMutex.Lock()
		if rm.currentGoRoutine == 0 {
			rm.currentGoRoutine = goRoutineID
			break
		} else if rm.currentGoRoutine == goRoutineID {
			break
		} else {
			rm.internalMutex.Unlock()
			time.Sleep(time.Millisecond)
			continue
		}
	}
	rm.lockCount++
	rm.internalMutex.Unlock()
}

func (rm *RecursiveMutex) Unlock() {
	rm.internalMutex.Lock()
	rm.lockCount--
	if rm.lockCount == 0 {
		rm.currentGoRoutine = 0
	}
	rm.internalMutex.Unlock()
}
