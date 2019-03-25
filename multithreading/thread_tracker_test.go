package multithreading

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestThreadTracker_TriggerGoroutine(t *testing.T) {
	var counter int8
	testCases := []struct {
		fns  []func(inputs []interface{})
		want int8
	}{
		{
			fns: []func(inputs []interface{}){
				func(inputs []interface{}) {
					counter = 1
				},
			},
			want: 1,
		}, {
			fns: []func(inputs []interface{}){
				func(inputs []interface{}) {
					counter = 2
				},
				func(inputs []interface{}) {
					// this will execute last because of the sleep
					time.Sleep(time.Duration(250) * time.Millisecond)
					counter = 1
				},
			},
			want: 1,
		}, {
			fns: []func(inputs []interface{}){
				func(inputs []interface{}) {
					time.Sleep(time.Duration(250) * time.Millisecond)
					// this will execute last because of the sleep
					counter = 2
				},
				func(inputs []interface{}) {
					counter = 1
				},
			},
			want: 2,
		},
	}

	for i, kase := range testCases {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			counter = -1
			threadTracker := MakeThreadTracker()

			for _, fn := range kase.fns {
				e := threadTracker.TriggerGoroutine(fn, nil)
				if !assert.NoError(t, e) {
					return
				}
			}
			threadTracker.Wait()
			assert.Equal(t, kase.want, counter)
		})
	}
}

func TestThreadTracker_TriggerGoroutine_Values(t *testing.T) {
	values := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	m := map[int]bool{}
	threadTracker := MakeThreadTracker()
	mutex := &sync.Mutex{}

	for _, v := range values {
		e := threadTracker.TriggerGoroutine(func(inputs []interface{}) {
			v := inputs[0].(int)

			mutex.Lock()
			m[v] = true
			mutex.Unlock()
		}, []interface{}{v})
		if !assert.NoError(t, e) {
			return
		}
	}

	threadTracker.Wait()
	assert.Equal(t, 10, len(m))
}

func TestThreadTracker_TriggerGoroutineWithDefers(t *testing.T) {
	var counter int8
	testCases := []struct {
		defers []func()
		want   int8
	}{
		{
			defers: nil,
			want:   10,
		}, {
			defers: []func(){},
			want:   10,
		}, {
			defers: []func(){
				func() {
					counter = 1
				},
			},
			want: 1,
		}, {
			defers: []func(){
				func() {
					// this defer will execute second
					counter = 2
				},
				func() {
					counter = 1
				},
			},
			want: 2,
		},
	}

	for i, kase := range testCases {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			counter = -1
			threadTracker := MakeThreadTracker()

			e := threadTracker.TriggerGoroutineWithDefers(
				kase.defers,
				func(inputs []interface{}) {
					counter = 10
				},
				nil,
			)
			if !assert.NoError(t, e) {
				return
			}

			threadTracker.Wait()
			assert.Equal(t, kase.want, counter)
		})
	}
}

func TestThreadTracker_panic(t *testing.T) {
	var counter int8
	testCases := []struct {
		defers []func()
		want   int8
	}{
		{
			defers: []func(){
				func() {
					// this will execute because we catch the panic here
					if r := recover(); r != nil {
						counter = 2
					}
				},
			},
			want: 2,
		}, {
			defers: []func(){
				func() {
					// this will execute after handling the panic
					counter = 3
				},
				func() {
					if r := recover(); r != nil {
						counter = 2
					}
				},
			},
			want: 3,
		}, {
			defers: []func(){
				func() {
					if r := recover(); r != nil {
						// this will not get executed because the panic has been handled in the last defer (first to execute)
						counter = 4
					}
				},
				func() {
					counter = 3
				},
				func() {
					if r := recover(); r != nil {
						counter = 2
					}
				},
			},
			want: 3,
		},
	}

	for i, kase := range testCases {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			counter = -1
			threadTracker := MakeThreadTracker()

			e := threadTracker.TriggerGoroutineWithDefers(
				kase.defers,
				func(inputs []interface{}) {
					panic("some error")
				},
				nil,
			)
			if !assert.NoError(t, e) {
				return
			}

			threadTracker.Wait()
			assert.Equal(t, kase.want, counter)
		})
	}
}

func TestThreadTracker_NumActiveThreads(t *testing.T) {
	testCases := []struct {
		numFns int
		want   uint64
	}{
		{
			numFns: 0,
			want:   0,
		}, {
			numFns: 1,
			want:   1,
		}, {
			numFns: 2,
			want:   2,
		},
	}

	for _, kase := range testCases {
		t.Run(fmt.Sprintf("%d", kase.numFns), func(t *testing.T) {
			threadTracker := MakeThreadTracker()

			mutex := &sync.Mutex{}
			mutex.Lock()
			for i := 0; i < kase.numFns; i++ {
				e := threadTracker.TriggerGoroutine(func(inputs []interface{}) {
					mutex.Lock() // blocking call
					mutex.Unlock()
				}, nil)
				if !assert.NoError(t, e) {
					mutex.Unlock()
					return
				}
			}
			if !assert.Equal(t, kase.want, threadTracker.NumActiveThreads()) {
				mutex.Unlock()
				return
			}

			mutex.Unlock()
			threadTracker.Wait()
			if !assert.Equal(t, uint64(0), threadTracker.NumActiveThreads()) {
				return
			}
		})
	}
}

func TestThreadTracker_Stop(t *testing.T) {
	testCases := []struct {
		stopMode  StopMode
		nThreads  uint64
		wantError bool
	}{
		{
			stopMode:  StopModeNoop,
			nThreads:  0,
			wantError: false,
		}, {
			stopMode:  StopModeNoop,
			nThreads:  1,
			wantError: false,
		}, {
			stopMode:  StopModeNoop,
			nThreads:  2,
			wantError: false,
		}, {
			stopMode:  StopModeError,
			nThreads:  0,
			wantError: true,
		}, {
			stopMode:  StopModeError,
			nThreads:  1,
			wantError: true,
		}, {
			stopMode:  StopModeError,
			nThreads:  2,
			wantError: true,
		},
	}

	for _, kase := range testCases {
		t.Run(fmt.Sprintf("%v", kase.stopMode), func(t *testing.T) {
			threadTracker := MakeThreadTracker()
			mutex := &sync.Mutex{}
			mutex.Lock()
			blockingFn := func(inputs []interface{}) {
				mutex.Lock() // blocking call
				mutex.Unlock()
			}

			// spin up number of prerequisite threads
			for i := 0; uint64(i) < kase.nThreads; i++ {
				e := threadTracker.TriggerGoroutine(blockingFn, nil)
				if !assert.NoError(t, e) {
					mutex.Unlock()
					return
				}
			}
			// sanity check
			if !assert.Equal(t, kase.nThreads, threadTracker.NumActiveThreads()) {
				mutex.Unlock()
				return
			}

			// run stop command and validate if there's an error or not
			threadTracker.Stop(kase.stopMode)
			e := threadTracker.TriggerGoroutine(blockingFn, nil)
			if kase.wantError && !assert.Error(t, e) {
				mutex.Unlock()
				return
			} else if !kase.wantError && !assert.NoError(t, e) {
				mutex.Unlock()
				return
			}

			// cleanup
			mutex.Unlock()
			threadTracker.Wait()
		})
	}
}
