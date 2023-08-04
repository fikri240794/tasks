package tasks

import (
	"math"
	"math/rand"
	"sync"
	"testing"
	"time"
)

func TestTask_Concurrent(t *testing.T) {
	var (
		maxConcurrentTask int
		task              Task
		taskLen           int
		taskCount         int
		taskActive        int
		syncMutex         sync.Mutex
	)

	rand.Seed(time.Now().UnixNano())
	maxConcurrentTask = rand.Intn(math.MaxInt8-2) + 2
	taskLen = rand.Intn(math.MaxInt8-2) + 2
	task = NewTask(maxConcurrentTask)

	t.Logf("maximum concurrent task: %d", maxConcurrentTask)
	t.Logf("task length: %d", taskLen)

	for i := 0; i < taskLen; i++ {
		task.Go(func() {
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
				t.Errorf("Expected maximum task active: %d, Got: %d", maxConcurrentTask, taskActive)
			}
		})
	}

	task.Wait()

	if taskCount != taskLen {
		t.Errorf("Expected task count: %d, Got: %d", taskLen, taskCount)
	}
}
