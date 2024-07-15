package worker_pool

import (
	"context"
	"log"
	"time"
)

type Worker struct {
	interval time.Duration
	task     Task
}

func NewWorker(task Task, interval time.Duration) Worker {
	return Worker{
		task:     task,
		interval: interval,
	}
}

func (w Worker) Start(ctx context.Context) {
LOOP:
	for {
		select {
		case <-ctx.Done():
			break LOOP
		case <-time.After(w.interval):
			if err := w.task.Process(ctx); err != nil {
				log.Printf("task process error: %v", err) // todo: log
			}
		}
	}
}

type WorkerPool struct {
	workers []Worker
	limit   int
}

func NewWorkerPool(workers []Worker) WorkerPool {
	return WorkerPool{
		workers: workers,
	}
}

func (wp WorkerPool) Run(ctx context.Context) {
	for _, w := range wp.workers {
		w.Start(ctx)
	}
}
