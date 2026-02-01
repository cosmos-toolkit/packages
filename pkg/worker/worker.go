// Package worker provides worker pool with configurable concurrency and retry.
// Base for SQS, cron, internal queue.
package worker

import (
	"context"
	"sync"

	"github.com/cosmos-toolkit/pkgs/pkg/retry"
)

// Job represents a unit of work.
type Job interface {
	Run(ctx context.Context) error
}

// JobFunc adapts a function to Job.
type JobFunc func(ctx context.Context) error

func (f JobFunc) Run(ctx context.Context) error { return f(ctx) }

// Config configures the pool.
type Config struct {
	Concurrency int
	Retry       retry.Config
}

// DefaultConfig returns default configuration (4 workers, default retry).
func DefaultConfig() Config {
	return Config{
		Concurrency: 4,
		Retry:       retry.DefaultConfig(),
	}
}

// Pool processes jobs with a fixed number of workers.
type Pool struct {
	cfg   Config
	jobs  <-chan Job
	done  chan struct{}
	wg    sync.WaitGroup
	start sync.Once
	stop  sync.Once
}

// NewPool creates a pool that reads jobs from jobsCh.
func NewPool(cfg Config, jobsCh <-chan Job) *Pool {
	if cfg.Concurrency < 1 {
		cfg.Concurrency = 1
	}
	return &Pool{cfg: cfg, jobs: jobsCh, done: make(chan struct{})}
}

// Start starts the workers. Returns immediately.
func (p *Pool) Start(ctx context.Context) {
	p.start.Do(func() {
		for i := 0; i < p.cfg.Concurrency; i++ {
			p.wg.Add(1)
			go p.worker(ctx)
		}
	})
}

func (p *Pool) worker(ctx context.Context) {
	defer p.wg.Done()
	for {
		select {
		case <-ctx.Done():
			return
		case <-p.done:
			return
		case job, ok := <-p.jobs:
			if !ok {
				return
			}
			_ = retry.Do(ctx, p.cfg.Retry, func() error { return job.Run(ctx) })
		}
	}
}

// Stop signals workers to stop and waits for completion.
func (p *Pool) Stop() {
	p.stop.Do(func() { close(p.done) })
	p.wg.Wait()
}
