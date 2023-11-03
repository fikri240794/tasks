package tasks

import (
	"runtime"
	"sync"
)

type ErrorTask interface {
	Go(task func() error)
	Wait() error
}

type errorTask struct {
	wg       sync.WaitGroup
	c        chan struct{}
	errsChan chan error
}

func NewErrorTask(maxConcurrentTask int) ErrorTask {
	if maxConcurrentTask < 1 {
		maxConcurrentTask = runtime.NumCPU()
	}

	return &errorTask{
		c:        make(chan struct{}, maxConcurrentTask),
		errsChan: make(chan error, 1),
	}
}

func (et *errorTask) Go(task func() error) {
	et.wg.Add(1)

	go func(taskToDo func() error) {
		var errRoutine error

		defer func() {
			et.wg.Done()
			<-et.c
		}()

		et.c <- struct{}{}

		if len(et.errsChan) > 0 {
			return
		}

		errRoutine = taskToDo()

		if errRoutine != nil {
			et.errsChan <- errRoutine
		}
	}(task)
}

func (et *errorTask) Wait() error {
	var err error

	et.wg.Wait()

	if len(et.errsChan) > 0 {
		err = <-et.errsChan

		return err
	}

	return nil
}
