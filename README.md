# go-tools
A collection of tools for Golang

[![Godoc](https://godoc.org/github.com/nikhilsaraf/go-tools?status.svg)](https://godoc.org/github.com/nikhilsaraf/go-tools)
[![Go Report Card](https://goreportcard.com/badge/github.com/nikhilsaraf/go-tools)](https://goreportcard.com/report/github.com/nikhilsaraf/go-tools)
[![Build Status](https://travis-ci.org/nikhilsaraf/go-tools.svg?branch=master)](https://travis-ci.org/nikhilsaraf/go-tools)

The [multithreading](multithreading) library currently supports a ThreadTracker struct that allows you to easily manage goroutines.

- It allows you to create new goroutines and Wait for all of them to finish.
- It allows you to set deferred functions inside the new goroutine it creates.
- You can easily handle panics inside these goroutines by passing a panic handler as a deferred function.
- It allows you to stop the threadTracker from receiving new functions with the option to return an error.
- It allows you to fetch the number of currently active goroutines that are managed by this threadTracker.
