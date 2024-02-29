package concurrent

import "sync"

// Task is now defined as a function that returns a key-value pair (string, string).
type Task func() (string, string)

// ResultStorage is a struct that encapsulates a map and a mutex.
// The map stores task results, and the mutex ensures safe concurrent access.
type ResultStorage struct {
	sync.Mutex
	Results map[string]string
}

// NewResultStorage creates a new resultStorage instance.
func NewResultStorage() *ResultStorage {
	return &ResultStorage{
		Results: make(map[string]string),
	}
}

// store safely stores a key-value pair in the results map.
func (rs *ResultStorage) store(key, value string) {
	rs.Lock()
	defer rs.Unlock()
	rs.Results[key] = value
}

// StartWorkerPool starts a pool with a fixed number of worker goroutines and waits for all tasks to complete.
// It takes a resultStorage pointer to store task results.
func StartWorkerPool(tasks chan Task, numberOfWorkers int, wg *sync.WaitGroup, storage *ResultStorage) {
	for i := 0; i < numberOfWorkers; i++ {
		wg.Add(1)
		go worker(tasks, wg, storage)
	}
}

// worker is a worker goroutine in the pool, responsible for executing tasks and storing their results.
func worker(tasks chan Task, wg *sync.WaitGroup, storage *ResultStorage) {
	defer wg.Done()
	for task := range tasks {
		key, value := task()      // Execute the task and get the result.
		storage.store(key, value) // Store the result safely.
	}
}
