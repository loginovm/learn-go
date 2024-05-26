//go:build !race

package hw05parallelexecution

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestConcurrency(t *testing.T) {
	var runTasksCount int
	go func() {
		for {
			runTasksCount = 0
			tasks := []Task{
				func() error { runTasksCount++; return nil },
				func() error { runTasksCount++; return nil },
			}
			Run(tasks, len(tasks), -1)
			if runTasksCount == 1 {
				break
			}
		}
	}()

	require.Eventually(t, func() bool {
		return runTasksCount == 1
	}, 5*time.Second, 1*time.Second, "tasks run not concurrently")
}
