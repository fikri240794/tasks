# Tasks
Limit your asynchronous goroutine processes in simple way. Inspired from https://pkg.go.dev/golang.org/x/sync/errgroup.

## Installation
```bash
go get github.com/fikri240794/tasks
```

## Usage
Example how to use Task:
```go
package main

import "github.com/fikri240794/tasks"

func main() {
	var (
		maxConcurrentTask int
		task              tasks.Task
	)

	maxConcurrentTask = 2 // set limit your async gorotine processes
	task = tasks.NewTask(maxConcurrentTask)

	task.Go(func() {
		// task 1
	})
	task.Go(func() {
		// task 2
	})
	task.Go(func() {
		// task 3
	})
	task.Go(func() {
		// task n...
	})

	task.Wait()
}
```

Example how to use ErrorTask:
```go
package main

import (
	"context"

	"github.com/fikri240794/tasks"
)

func main() {
	var (
		ctx               context.Context
		maxConcurrentTask int
		errTask           tasks.ErrorTask
		errTaskCtx        context.Context
		err               error
	)

	ctx = context.Background() // any context from (from param, request, etc...)
	maxConcurrentTask = 2      // set limit your async gorotine processes
	errTask, errTaskCtx = tasks.NewErrorTask(maxConcurrentTask, ctx) // always create new context for errTask

	errTask.Go(func() error {
		// task 1
		someFunc(errTaskCtx, args...) // use errTaskCtx context for all errTask goroutine
	})
	errTask.Go(func() error {
		// task 2
	})
	errTask.Go(func() error {
		// task 3
	})
	errTask.Go(func() error {
		// task n...
	})

	err = errTask.Wait()
	if err != nil {
		log.Fatal(err)
	}
}
```