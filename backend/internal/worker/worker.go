package worker

import (
	"context"
	"sync"

	"github.com/Tabhi109/investwise/internal/logger"
)

// Job represents a background execution unit
type Job func(ctx context.Context)

// WorkerPool manages concurrent execution of tasks in background goroutines
type WorkerPool struct {
	jobQueue   chan Job
	numWorkers int
	ctx        context.Context
	cancel     context.CancelFunc
	wg         sync.WaitGroup
}

// NewWorkerPool instantiates a WorkerPool with specific depth and worker counts
func NewWorkerPool(numWorkers int) *WorkerPool {
	ctx, cancel := context.WithCancel(context.Background())
	return &WorkerPool{
		jobQueue:   make(chan Job, 1024), // Buffered queue
		numWorkers: numWorkers,
		ctx:        ctx,
		cancel:     cancel,
	}
}

// Start spawns the background worker routines
func (p *WorkerPool) Start() {
	for i := 0; i < p.numWorkers; i++ {
		p.wg.Add(1)
		go p.worker(i)
	}
	logger.Info("Background concurrency worker pool started", "count", p.numWorkers)
}

// Submit enqueues a job for background processing
func (p *WorkerPool) Submit(job Job) {
	select {
	case p.jobQueue <- job:
	default:
		// Queue full: log and drop or handle backpressure (e.g. run synchronously, block, etc.)
		logger.Warn("Worker pool job queue full, task dropped")
	}
}

// worker represents a single consumer routine loop
func (p *WorkerPool) worker(id int) {
	defer p.wg.Done()
	logger.Debug("Worker routine initialized", "worker_id", id)

	for {
		select {
		case <-p.ctx.Done():
			logger.Debug("Worker routine terminating", "worker_id", id)
			return
		case job, ok := <-p.jobQueue:
			if !ok {
				return
			}
			p.runJob(id, job)
		}
	}
}

// runJob wrapper handles execution and safe panic recovery
func (p *WorkerPool) runJob(workerID int, job Job) {
	defer func() {
		if r := recover(); r != nil {
			logger.Error("Panic recovered inside background worker job", nil, "worker_id", workerID, "panic", r)
		}
	}()
	job(p.ctx)
}

// Stop cancels context, waits for workers to complete their active jobs, and shuts down
func (p *WorkerPool) Stop() {
	p.cancel()
	p.wg.Wait()
	// Drain channel
	select {
	case job, ok := <-p.jobQueue:
		if ok {
			p.runJob(-1, job)
		}
	default:
	}
	logger.Info("Worker pool gracefully stopped")
}
