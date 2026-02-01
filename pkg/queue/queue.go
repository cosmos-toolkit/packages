// Package queue provides queue interface (Publish, Consume) and in-memory implementation.
// Other implementations: SQS, Rabbit, etc.
package queue

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
)

// Message represents a queue message.
type Message struct {
	ID      string
	Body    []byte
	Headers map[string]string
}

// Publisher publishes messages.
type Publisher interface {
	Publish(ctx context.Context, topic string, body []byte, headers map[string]string) error
}

// Consumer consumes messages (handler receives the message; Ack confirms processing).
type Consumer interface {
	Consume(ctx context.Context, topic string, handler func(ctx context.Context, m *Message) error) error
}

// Queue combines Publisher and Consumer.
type Queue interface {
	Publisher
	Consumer
}

// InMemory implements Queue in memory (useful for tests and internal queues).
type InMemory struct {
	mu      sync.Mutex
	topics  map[string][]*Message
	waiters map[string][]chan *Message
}

// NewInMemory creates an in-memory queue.
func NewInMemory() *InMemory {
	return &InMemory{
		topics:  make(map[string][]*Message),
		waiters: make(map[string][]chan *Message),
	}
}

// Publish adds a message to the topic and delivers to a blocked consumer if any.
func (q *InMemory) Publish(ctx context.Context, topic string, body []byte, headers map[string]string) error {
	m := &Message{ID: newID(), Body: body, Headers: make(map[string]string)}
	if headers != nil {
		m.Headers = headers
	}
	q.mu.Lock()
	if len(q.waiters[topic]) > 0 {
		ch := q.waiters[topic][0]
		q.waiters[topic] = q.waiters[topic][1:]
		q.mu.Unlock()
		select {
		case ch <- m:
		case <-ctx.Done():
			return ctx.Err()
		}
		return nil
	}
	q.topics[topic] = append(q.topics[topic], m)
	q.mu.Unlock()
	return nil
}

// Consume processes messages from the topic; blocks until ctx is cancelled.
func (q *InMemory) Consume(ctx context.Context, topic string, handler func(ctx context.Context, m *Message) error) error {
	ch := make(chan *Message, 64)
	q.mu.Lock()
	q.waiters[topic] = append(q.waiters[topic], ch)
	// process already queued messages (without sending to ch to avoid deadlock)
	for len(q.topics[topic]) > 0 {
		m := q.topics[topic][0]
		q.topics[topic] = q.topics[topic][1:]
		q.mu.Unlock()
		if err := handler(ctx, m); err != nil {
			_ = err
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		q.mu.Lock()
	}
	q.mu.Unlock()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case m := <-ch:
			if err := handler(ctx, m); err != nil {
				_ = err
			}
		}
	}
}

var idCounter atomic.Uint64

func newID() string {
	return fmt.Sprintf("%d", idCounter.Add(1))
}
