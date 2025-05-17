package worker

import (
	"errors"
	"sync"
)

// Pool manages a set of workers to process Operations concurrently.
// It stops all workers on the first error, but collects all errors.
type Pool struct {
	workers []*Worker
	jobs    chan Operation
	errors  chan error
	stop    chan struct{}

	wg       sync.WaitGroup
	stopOnce sync.Once

	errMu sync.Mutex
	errs  []error
}

// NewPool creates a new Pool with the specified number of workers and job queue size.
func NewPool(workerCount, queueSize int) *Pool {
	stop := make(chan struct{})
	jobs := make(chan Operation, queueSize)
	errors := make(chan error, workerCount*2)

	workers := make([]*Worker, workerCount)
	for i := 0; i < workerCount; i++ {
		workers[i] = NewWorker(i, jobs, errors, stop)
	}

	return &Pool{
		workers: workers,
		jobs:    jobs,
		errors:  errors,
		stop:    stop,
	}
}

// Start launches all workers and an error watcher goroutine.
func (p *Pool) Start() {
	for _, w := range p.workers {
		w.Start(&p.wg)
	}

	go func() {
		var once sync.Once
		for err := range p.errors {
			if err == nil {
				continue
			}
			p.addError(err)
			once.Do(func() {
				// On first error, broadcast stop signal
				p.Stop()
			})
		}
	}()
}

// Submit adds a job to the pool's job queue.
func (p *Pool) Submit(job Operation) error {
	select {
	case <-p.stop:
		return errors.New("pool stopped, cannot submit new jobs")
	default:
	}

	select {
	case p.jobs <- job:
		return nil
	case <-p.stop:
		return errors.New("pool stopped, cannot submit new jobs")
	}
}

// Stop closes the stop and jobs channels once.
func (p *Pool) Stop() {
	p.stopOnce.Do(func() {
		close(p.stop)
		close(p.jobs)
	})
}

// Wait blocks until all workers have exited, then closes the errors channel.
// Returns all errors collected during processing.
func (p *Pool) Wait() []error {
	p.wg.Wait()
	close(p.errors)

	p.errMu.Lock()
	defer p.errMu.Unlock()
	return p.errs
}

// addError appends a new error to the internal slice in a thread-safe manner.
func (p *Pool) addError(err error) {
	p.errMu.Lock()
	defer p.errMu.Unlock()
	p.errs = append(p.errs, err)
}
