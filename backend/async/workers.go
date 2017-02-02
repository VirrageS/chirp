package async

import (
	"sync"
)

const (
	// TODO: should maxWorkers be depended on number of CPU's or current
	// number of connections to database or cache?

	// Number of workers spawned for each pool.
	maxWorkers = 10

	// Size of task and result queues. It allows to buffer tasks and results
	// so no blocking occurs (until buffer is not full).
	bufferSize = 100
)

// Task represents abstract type to hold any value which will be delieverd to
// workers as a value to process.
type Task interface{}

// Result is struct to represent result of Task, which is either some kind
// of value or error if any occured.
type Result struct {
	Value interface{}
	Error error
}

// WorkerPool is struct to represent pool of workers to which one can send task
// and receive results.
type WorkerPool struct {
	wg        *sync.WaitGroup
	inChan    chan Task
	outChan   chan *Result
	closeChan chan bool
}

// NewWorkerPool creates new instance of `WorkerPool`. As parameter it receives
// `handler` which is function that is responsible for processing task
// and returning the result.
func NewWorkerPool(handler func(value Task) *Result) *WorkerPool {
	var (
		wg        = &sync.WaitGroup{}
		inChan    = make(chan Task, bufferSize)
		outChan   = make(chan *Result, bufferSize)
		closeChan = make(chan bool, maxWorkers)
	)

	for i := 0; i < maxWorkers; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			for {
				select {
				case value := <-inChan:
					outChan <- handler(value)
				case <-closeChan:
					return
				}
			}
		}()
	}

	return &WorkerPool{
		wg:        wg,
		inChan:    inChan,
		outChan:   outChan,
		closeChan: closeChan,
	}
}

// PostTask adds new task to task queue.
// Adding tasks after `Close` results in panic.
func (wp *WorkerPool) PostTask(task Task) {
	wp.inChan <- task
}

// GetResult returns first result from result queue or blocks until new one appears.
func (wp *WorkerPool) GetResult() *Result {
	return <-wp.outChan
}

// Close finishes all workers and closes all channels.
func (wp *WorkerPool) Close() {
	for i := 0; i < maxWorkers; i++ {
		wp.closeChan <- true
	}

	wp.wg.Wait()
	close(wp.inChan)
	close(wp.outChan)
	close(wp.closeChan)
}
