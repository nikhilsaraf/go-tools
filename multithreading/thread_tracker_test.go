package multithreading

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestThreadTracker_TriggerGoroutine(t *testing.T) {
	var counter int8
	testCases := []struct {
		fns  []func()
		want int8
	}{
		{
			fns: []func(){
				func() {
					counter = 1
				},
			},
			want: 1,
		}, {
			fns: []func(){
				func() {
					counter = 2
				},
				func() {
					// this will execute last because of the sleep
					time.Sleep(time.Duration(250) * time.Millisecond)
					counter = 1
				},
			},
			want: 1,
		}, {
			fns: []func(){
				func() {
					time.Sleep(time.Duration(250) * time.Millisecond)
					// this will execute last because of the sleep
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

			for _, fn := range kase.fns {
				threadTracker.TriggerGoroutine(fn)
			}
			threadTracker.Wait()
			assert.Equal(t, kase.want, counter)
		})
	}
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

			threadTracker.TriggerGoroutineWithDefers(func() {
				counter = 10
			}, kase.defers)
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

			threadTracker.TriggerGoroutineWithDefers(func() {
				panic("some error")
			}, kase.defers)
			threadTracker.Wait()
			assert.Equal(t, kase.want, counter)
		})
	}
}
