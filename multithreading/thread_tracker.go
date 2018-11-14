package multithreading

import (
	"sync"
)

// ThreadTracker allows you to easily manage goroutines
type ThreadTracker struct {
	activeThreadsCounter *sync.WaitGroup
}

// MakeThreadTracker is a factory method for ThreadTracker
func MakeThreadTracker() *ThreadTracker {
	return &ThreadTracker{
		activeThreadsCounter: &sync.WaitGroup{},
	}
}

// TriggerGoroutine initiates a new goroutine while tracking it
// typical usage:
//     threadTracker.TriggerGoroutine(func() {
//         fmt.Printf("Hello %s\n", "World")
//     })
func (t *ThreadTracker) TriggerGoroutine(fn func()) {
	t.TriggerGoroutineWithDefers(fn, nil)
}

// TriggerGoroutineWithDefers initiates a new goroutine while tracking it
// typical usage:
//     threadTracker.TriggerGoroutine(func() {
//         fmt.Printf("Hello %s -- this should appear first\n", "World")
//     }, []func(){
// 	       func() { fmt.Printf("this should appear third\n") },
// 	       func() { fmt.Printf("this should appear second\n") },
//     })
func (t *ThreadTracker) TriggerGoroutineWithDefers(fn func(), deferredFns []func()) {
	t.activeThreadsCounter.Add(1)

	go func() {
		// defer this func so it is called even if the code panics
		defer func() {
			t.activeThreadsCounter.Done()
		}()
		if deferredFns != nil {
			for _, dFn := range deferredFns {
				defer dFn()
			}
		}

		fn()
	}()
}

// Wait blocks until all goroutines finish
func (t *ThreadTracker) Wait() {
	t.activeThreadsCounter.Wait()
}
