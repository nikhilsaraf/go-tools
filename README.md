# go-tools
A collection of tools for Golang, focusing on concurrency and goroutines

[![Godoc](https://godoc.org/github.com/nikhilsaraf/go-tools?status.svg)](https://godoc.org/github.com/nikhilsaraf/go-tools)
[![Go Report Card](https://goreportcard.com/badge/github.com/nikhilsaraf/go-tools)](https://goreportcard.com/report/github.com/nikhilsaraf/go-tools)
[![Build Status](https://travis-ci.org/nikhilsaraf/go-tools.svg?branch=master)](https://travis-ci.org/nikhilsaraf/go-tools)
[![License](https://img.shields.io/github/license/mashape/apistatus.svg?maxAge=2592000)](https://raw.githubusercontent.com/nikhilsaraf/go-tools/master/LICENSE)

The [multithreading](multithreading) library currently supports a `ThreadTracker` struct that allows you to easily manage goroutines.

- Create new goroutines.
- Wait for all goroutines to finish.
- Set deferred functions to be executed after goroutines finish.
- Easily handle panics inside goroutines with a panic handler.
- Stop the `threadTracker` from receiving new functions.
- Fetch the number of currently active goroutines.

## Install

Install the package with:

```bash
go get github.com/nikhilsaraf/go-tools/multithreading
```

Import it with:

```go
import "github.com/nikhilsaraf/go-tools/multithreading"
```

and use `multithreading` as the package name inside the code.

## Example

```go
package main

import (
  "fmt"
  "github.com/nikhilsaraf/go-tools/multithreading"
)

func main() {
  // create thread tracker instance
  threadTracker := multithreading.MakeThreadTracker()

  // start thread functions
  for i := 0; i < 10; i++ {
    err := threadTracker.TriggerGoroutine(func(inputs []interface{}) {
      // pass `i` as a value to the goroutine and read from `inputs`.
      // this is needed to "bind" the variable to this goroutine.
      value := inputs[0].(int)
      
      fmt.Printf("Goroutine #%d\n", value)
    }, []interface{}{i})

    if err != nil {
      panic(err)
    }
  }

  // wait for all threads to finish
  threadTracker.Wait()
  fmt.Printf("done\n")
}

```

Sample Output:
```
Goroutine #1
Goroutine #2
Goroutine #9
Goroutine #0
Goroutine #3
Goroutine #7
Goroutine #6
Goroutine #4
Goroutine #8
Goroutine #5
done
```

## Test Examples
[thread_tracker_test.go](/multithreading/thread_tracker_test.go)
