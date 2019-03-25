package multithreading

import (
	"fmt"
	"log"
	"sync"
)

// ThreadTracker allows you to easily manage goroutines
type ThreadTracker struct {
	activeThreadsCounter *sync.WaitGroup
	stopMode             *StopMode
	count                uint64
	mutexCount           *sync.Mutex
}

// MakeThreadTracker is a factory method for ThreadTracker
func MakeThreadTracker() *ThreadTracker {
	return &ThreadTracker{
		activeThreadsCounter: &sync.WaitGroup{},
		stopMode:             nil,
		count:                0,
		mutexCount:           &sync.Mutex{},
	}
}

// StopMode represents the behavior of the threadTracker once it has been stopped
type StopMode uint8

// types of StopMode
const (
	StopModeError StopMode = iota
	StopModeNoop
)

// TriggerGoroutine initiates a new goroutine while tracking it
// typical usage 1:
//     threadTracker.TriggerGoroutine(func(inputs []interface{}) {
//         fmt.Printf("Hello %s\n", "World")
//     })
// typical usage 2:
//     threadTracker.TriggerGoroutine(func(inputs []interface{}) {
//         fmt.Printf("Hello %s\n", inputs[0])
//     }, "World")
func (t *ThreadTracker) TriggerGoroutine(fn func(inputs []interface{}), inputs []interface{}) error {
	return t.TriggerGoroutineWithDefers(nil, fn, inputs)
}

// TriggerGoroutineWithDefers initiates a new goroutine while tracking it
// typical usage:
//     threadTracker.TriggerGoroutineWithDefers([]func(){
// 	       func() { fmt.Printf("this should appear third\n") },
// 	       func() { fmt.Printf("this should appear second\n") },
//     }, func(inputs []interface{}) {
//         fmt.Printf("Hello %s -- this should appear first\n", "World")
//     })
func (t *ThreadTracker) TriggerGoroutineWithDefers(deferredFns []func(), fn func(inputs []interface{}), inputs []interface{}) error {
	if t.stopMode != nil && *t.stopMode == StopModeError {
		return fmt.Errorf("cannot add more threads since this threadTracker has been stopped")
	} else if t.stopMode != nil && *t.stopMode == StopModeNoop {
		log.Printf("cannot add more threads since this threadTracker has been stopped")
		return nil
	}

	t.incThreadCount()

	go func() {
		// defer this func so it is called even if the code panics
		defer t.decThreadCount()
		if deferredFns != nil {
			for _, dFn := range deferredFns {
				defer dFn()
			}
		}

		fn(inputs)
	}()
	return nil
}

// Wait blocks until all goroutines finish
func (t *ThreadTracker) Wait() {
	t.activeThreadsCounter.Wait()
}

// Stop prevents this threadTracker from accepting any more threads
func (t *ThreadTracker) Stop(stopMode StopMode) {
	t.stopMode = &stopMode
}

// NumActiveThreads fetches the number of active threads managed by this threadTracker
func (t *ThreadTracker) NumActiveThreads() uint64 {
	t.mutexCount.Lock()
	defer t.mutexCount.Unlock()

	return t.count
}

func (t *ThreadTracker) incThreadCount() {
	t.mutexCount.Lock()
	defer t.mutexCount.Unlock()

	t.activeThreadsCounter.Add(1)
	t.count++
}

func (t *ThreadTracker) decThreadCount() {
	t.mutexCount.Lock()
	defer t.mutexCount.Unlock()

	t.activeThreadsCounter.Done()
	t.count--
}
