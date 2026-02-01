// Package outbox implements the outbox pattern: persist + publish.
// Base for event-driven (persist event to DB and publish later).
package outbox

import (
	"context"
	"database/sql"
	"encoding/json"
	"sync"
	"time"
)

// Store persists and reads outbox events.
type Store interface {
	Save(ctx context.Context, topic string, payload []byte, metadata map[string]string) error
	Pending(ctx context.Context, limit int) ([]*Event, error)
	MarkPublished(ctx context.Context, id int64) error
}

// Publisher publishes a message (e.g. queue, broker).
type Publisher interface {
	Publish(ctx context.Context, topic string, body []byte, headers map[string]string) error
}

// Event represents an outbox event.
type Event struct {
	ID        int64
	Topic     string
	Payload   []byte
	Metadata  map[string]string
	CreatedAt time.Time
}

// DefaultTTL is the default TTL for processing (avoid reprocessing indefinitely).
const DefaultTTL = 24 * time.Hour

// Processor processes pending events and publishes them.
type Processor struct {
	store     Store
	publisher Publisher
	interval  time.Duration
	stop      chan struct{}
	wg        sync.WaitGroup
}

// NewProcessor creates an outbox processor.
func NewProcessor(store Store, publisher Publisher, interval time.Duration) *Processor {
	if interval <= 0 {
		interval = 10 * time.Second
	}
	return &Processor{store: store, publisher: publisher, interval: interval, stop: make(chan struct{})}
}

// Start starts processing in the background. Use Stop() to stop.
func (p *Processor) Start(ctx context.Context) {
	p.wg.Add(1)
	go func() {
		defer p.wg.Done()
		tick := time.NewTicker(p.interval)
		defer tick.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-p.stop:
				return
			case <-tick.C:
				_ = p.process(ctx)
			}
		}
	}()
}

// Stop signals stop and waits for completion.
func (p *Processor) Stop() {
	close(p.stop)
	p.wg.Wait()
}

func (p *Processor) process(ctx context.Context) error {
	events, err := p.store.Pending(ctx, 100)
	if err != nil || len(events) == 0 {
		return err
	}
	for _, e := range events {
		if err := p.publisher.Publish(ctx, e.Topic, e.Payload, e.Metadata); err != nil {
			continue
		}
		_ = p.store.MarkPublished(ctx, e.ID)
	}
	return nil
}

// SQLStore implements Store with outbox table (id, topic, payload, metadata_json, created_at, published_at).
type SQLStore struct {
	DB    *sql.DB
	Table string
}

// Save inserts an event into the outbox table.
func (s *SQLStore) Save(ctx context.Context, topic string, payload []byte, metadata map[string]string) error {
	table := s.Table
	if table == "" {
		table = "outbox"
	}
	metaJSON, _ := json.Marshal(metadata)
	_, err := s.DB.ExecContext(ctx,
		`INSERT INTO `+table+` (topic, payload, metadata_json, created_at) VALUES (?, ?, ?, ?)`,
		topic, payload, metaJSON, time.Now())
	return err
}

// Pending returns events not yet published (published_at IS NULL).
func (s *SQLStore) Pending(ctx context.Context, limit int) ([]*Event, error) {
	table := s.Table
	if table == "" {
		table = "outbox"
	}
	rows, err := s.DB.QueryContext(ctx,
		`SELECT id, topic, payload, metadata_json, created_at FROM `+table+` WHERE published_at IS NULL ORDER BY id LIMIT ?`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []*Event
	for rows.Next() {
		var e Event
		var metaJSON []byte
		if err := rows.Scan(&e.ID, &e.Topic, &e.Payload, &metaJSON, &e.CreatedAt); err != nil {
			return nil, err
		}
		_ = json.Unmarshal(metaJSON, &e.Metadata)
		out = append(out, &e)
	}
	return out, rows.Err()
}

// MarkPublished marks the event as published.
func (s *SQLStore) MarkPublished(ctx context.Context, id int64) error {
	table := s.Table
	if table == "" {
		table = "outbox"
	}
	_, err := s.DB.ExecContext(ctx, `UPDATE `+table+` SET published_at = ? WHERE id = ?`, time.Now(), id)
	return err
}
