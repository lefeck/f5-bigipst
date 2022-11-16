package pkg

import (
	"github.com/e-XpertSolutions/f5-rest-client/f5"
	"log"
)

//Job interface
type Job interface {
	Create(client *f5.Client) (err error)
}

// The job processor, where the real business logic is handled, is responsible for receiving and processing tasks,
// and it needs to tell the scheduler if it is ready to receive more tasks.
type worker struct {
	// Multiple workers share a worker queue WorkerQueue, which is used to register their own work channel (chan Job) to WorkerQueue
	// when the worker is idle, so that the worker can receive task requests
	workerQueue chan chan Job
	// jobQueue is an unbuffered job channel that receives Job
	jobQueue chan Job
	quit     chan bool
}

// Initialize a worker thread
func NewWorkers(workPool chan chan Job) *worker {
	return &worker{
		workerQueue: workPool,
		jobQueue:    make(chan Job),
		quit:        make(chan bool),
	}
}

// Define a start method for the thread to indicate that it is listening for a task to begin processing
func (w *worker) Start(client *f5.Client, ch chan struct{}) {
	go func() {
		for {
			w.workerQueue <- w.jobQueue //Register the worker channel to the thread pool
			select {
			case task := <-w.jobQueue:
				if err := task.Create(client); err != nil {
					log.Fatalf("create configure failed :%s", err)
				}
				ch <- struct{}{} // The goroutine ends, then send to signals
			case <-w.quit:
				return
			}
		}
	}()
}

// The thread stops working
func (w *worker) Stop() {
	go func() {
		w.quit <- true
	}()
}

//The task distributor can distribute the tasks in the task queue to the threads in the thread pool one by one for processing
type Dispatcher struct {
	WorkerQueue chan chan Job
	MaxNum      int
	JobQueue    chan Job
}

// Instantiate a task dispatcher
func NewDispatcher(maxWorkerNum int) *Dispatcher {
	return &Dispatcher{
		WorkerQueue: make(chan chan Job, maxWorkerNum),
		MaxNum:      maxWorkerNum,
		JobQueue:    make(chan Job),
	}
}

// Assign Tasks
func (d *Dispatcher) Dispatch() {
	for {
		select {
		// Remove a task from the task queue
		case jobObj := <-d.JobQueue:
			go func(job Job) {
				// Take a thread out of the thread pool
				workChan := <-d.WorkerQueue
				// A task from the task queue is processed by the worker thread
				workChan <- job
			}(jobObj)
		}
	}
}

// Start Task allocator starts running and distributing tasks
func (d *Dispatcher) Run(client *f5.Client, ch chan struct{}) {
	//Creating a new worker Thread
	for i := 0; i < d.MaxNum; i++ {
		workerObj := NewWorkers(d.WorkerQueue)
		//Start the thread
		workerObj.Start(client, ch)
	}
	// Distribute Tasks
	go d.Dispatch()
}
