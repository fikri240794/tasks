package tasks

import (
	"errors"
	"fmt"
	"math"
	"math/rand"
	"sync"
	"testing"
	"time"
)

func TestErrorTask_Error(t *testing.T) {
	var (
		maxConcurrentTask int
		errTask           ErrorTask
		taskLen           int
		err               error
	)

	rand.Seed(time.Now().UnixNano())
	maxConcurrentTask = rand.Intn(math.MaxInt8-2) + 2
	taskLen = rand.Intn(math.MaxInt8-2) + 2
	errTask = NewErrorTask(maxConcurrentTask)

	for i := 0; i < taskLen; i++ {
		errTask.Go(func() error {
			var someRandInt int = rand.Intn(math.MaxInt8-2) + 2
			var someRandInt2 int = someRandInt

			if someRandInt%1 == 0 && someRandInt%someRandInt2 == 0 {
				return errors.New("some error")
			}

			return nil
		})
	}

	err = errTask.Wait()

	if err == nil {
		t.Fatal("Expected error: some error, Got: nil")
	}
}

func TestErrorTask_Concurrent(t *testing.T) {
	var (
		maxConcurrentTask int
		errTask           ErrorTask
		taskLen           int
		taskCount         int
		taskActive        int
		syncMutex         sync.Mutex
		err               error
	)

	rand.Seed(time.Now().UnixNano())
	maxConcurrentTask = rand.Intn(math.MaxInt8-2) + 2
	taskLen = rand.Intn(math.MaxInt8-2) + 2
	errTask = NewErrorTask(maxConcurrentTask)

	t.Logf("maximum concurrent task: %d", maxConcurrentTask)
	t.Logf("task length: %d", taskLen)

	for i := 0; i < taskLen; i++ {
		errTask.Go(func() error {
			defer func() {
				syncMutex.Lock()
				taskActive--
				syncMutex.Unlock()
			}()

			syncMutex.Lock()
			taskCount++
			syncMutex.Unlock()

			syncMutex.Lock()
			taskActive++
			syncMutex.Unlock()
			if taskActive > maxConcurrentTask {
				return fmt.Errorf("Expected maximum task active: %d, Got: %d", maxConcurrentTask, taskActive)
			}

			return nil
		})
	}

	err = errTask.Wait()

	if err != nil {
		t.Fatalf("Expected error: nil, Got: %s", err.Error())
	}

	if taskCount != taskLen {
		t.Errorf("Expected task count: %d, Got: %d", taskLen, taskCount)
	}
}
