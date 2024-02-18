package test

import (
	"OrdDeFi-Virtual-Machine/concurrent"
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestConcurrentPool(t *testing.T) {
	tasks := make(chan concurrent.Task, 50)
	var wg sync.WaitGroup
	storage := concurrent.NewResultStorage()

	// Start the worker pool with 12 concurrent worker goroutines.
	concurrent.StartWorkerPool(tasks, 12, &wg, storage)

	// Assign tasks to the worker pool.
	for i := 0; i < 100; i++ {
		i := i // Capture the loop variable.
		tasks <- func() (string, string) {
			println("Current task:", i)
			time.Sleep(time.Second)
			// The task returns a key-value pair.
			return fmt.Sprintf("task%d", i), fmt.Sprintf("result of task %d", i)
		}
	}

	fmt.Printf("Closing task.\n")
	close(tasks) // All tasks have been assigned, close the tasks channel.
	wg.Wait()    // Wait for all worker goroutines to complete.
	fmt.Printf("Finished executing task.\n")

	// Print the results.
	for key, value := range storage.Results {
		fmt.Printf("%s: %s\n", key, value)
	}
}
