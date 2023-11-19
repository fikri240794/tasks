package tasks

import (
	"context"
	"runtime"
	"sync"
)

type ErrorTask interface {
	Go(task func() error)
	Wait() error
}

type errorTask struct {
	wg        sync.WaitGroup
	c         chan struct{}
	err       error
	once      sync.Once
	ctx       context.Context
	cancelFun context.CancelFunc
}

func NewErrorTask(maxConcurrentTask int, ctx context.Context) (ErrorTask, context.Context) {
	var errTask *errorTask

	if maxConcurrentTask < 1 {
		maxConcurrentTask = runtime.NumCPU()
	}

	errTask = &errorTask{
		c: make(chan struct{}, maxConcurrentTask),
	}

	errTask.ctx, errTask.cancelFun = context.WithCancel(ctx)

	return errTask, errTask.ctx
}

func (et *errorTask) Go(task func() error) {
	if et.ctx.Err() != nil {
		return
	}

	et.c <- struct{}{}
	et.wg.Add(1)

	go func(taskToDo func() error) {
		var errRoutine error

		defer func() {
			<-et.c
			et.wg.Done()
		}()

		if et.ctx.Err() != nil {
			return
		}

		errRoutine = taskToDo()
		if errRoutine != nil {
			et.once.Do(func() {
				et.err = errRoutine
				et.cancelFun()
			})
		}
	}(task)
}

func (et *errorTask) Wait() error {
	var err error

	et.wg.Wait()
	et.cancelFun()
	err = et.err

	return err
}
