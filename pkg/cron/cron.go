// Package cron provides a wrapper for robfig/cron for job scheduling.
// Distributed lock and per-job metrics can be added later.
package cron

import (
	"context"
	"sync"

	"github.com/robfig/cron/v3"
)

// Job is a scheduled task.
type Job interface {
	Run(ctx context.Context) error
}

// JobFunc adapts a function to Job.
type JobFunc func(ctx context.Context) error

func (f JobFunc) Run(ctx context.Context) error { return f(ctx) }

// Scheduler schedules and runs jobs.
type Scheduler struct {
	cron    *cron.Cron
	entries map[string]cron.EntryID
	mu      sync.Mutex
}

// New creates a scheduler. Use Start() to begin.
func New() *Scheduler {
	return &Scheduler{
		cron:    cron.New(),
		entries: make(map[string]cron.EntryID),
	}
}

// Add registers a job to run at spec (cron format: "0 * * * *" = every minute).
// name identifies the job (for removal later).
func (s *Scheduler) Add(name, spec string, job Job) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	id, err := s.cron.AddFunc(spec, func() {
		_ = job.Run(context.Background())
	})
	if err != nil {
		return err
	}
	s.entries[name] = id
	return nil
}

// Remove removes a job by name.
func (s *Scheduler) Remove(name string) {
	s.mu.Lock()
	id, ok := s.entries[name]
	delete(s.entries, name)
	s.mu.Unlock()
	if ok {
		s.cron.Remove(id)
	}
}

// Start starts the scheduler. Blocks until ctx is cancelled.
func (s *Scheduler) Start(ctx context.Context) error {
	s.cron.Start()
	<-ctx.Done()
	s.cron.Stop()
	return ctx.Err()
}

// Stop stops the scheduler.
func (s *Scheduler) Stop() {
	s.cron.Stop()
}
