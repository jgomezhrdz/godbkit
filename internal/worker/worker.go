package worker

import (
	"fmt"
	"sync"
)

type Operation func(...any) error

type Worker struct {
	id     int
	jobs   <-chan Operation
	errors chan<- error
	stop   <-chan struct{}
}

func NewWorker(id int, jobs <-chan Operation, errChan chan<- error, stopSignal <-chan struct{}) *Worker {
	return &Worker{
		id:     id,
		jobs:   jobs,
		errors: errChan,
		stop:   stopSignal,
	}
}

func (w *Worker) Start(wg *sync.WaitGroup) {
	wg.Add(1)
	go w.process(wg)
}

func (w *Worker) process(wg *sync.WaitGroup) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Worker %d panicked: %v\n", w.id, r)
		}
		wg.Done()
	}()

	for {
		select {
		case <-w.stop:
			fmt.Printf("Worker %d stopped\n", w.id)
			return
		case op, ok := <-w.jobs:
			if !ok {
				return
			}
			err := op()

			if err != nil {
				w.errors <- fmt.Errorf("an error occurred in worker %s", err.Error())
				return
			}
		}
	}
}
