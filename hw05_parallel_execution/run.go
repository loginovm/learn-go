package hw05parallelexecution

import (
	"errors"
	"sync"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	maxErrors := int32(m)
	var errNum int32
	wg := sync.WaitGroup{}
	taskChan := make(chan Task)

	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(taskChan <-chan Task) {
			for task := range taskChan {
				err := task()
				if err != nil {
					atomic.AddInt32(&errNum, 1)
				}
			}
			defer wg.Done()
		}(taskChan)
	}

	var err error
	for _, t := range tasks {
		if maxErrors > 0 && atomic.LoadInt32(&errNum) >= maxErrors {
			err = ErrErrorsLimitExceeded
			break
		}
		taskChan <- t
	}

	close(taskChan)
	wg.Wait()

	return err
}
