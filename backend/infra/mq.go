package infra

import (
	"context"
	"fmt"
	"sync"

	"github.com/huodaoshi/harness/backend/conf"
)

// MessageQueue is the interface for publishing and subscribing to messages.
type MessageQueue interface {
	Publish(ctx context.Context, topic string, msg []byte) error
	Subscribe(ctx context.Context, topic string, handler func([]byte)) error
	Close() error
}

// LocalQueue is an in-process message queue backed by buffered channels.
type LocalQueue struct {
	mu       sync.RWMutex
	channels map[string]chan []byte
}

// Publish writes a message to the topic's buffered channel (non-blocking).
// If the channel does not exist, it is created with a buffer of 100.
func (q *LocalQueue) Publish(_ context.Context, topic string, msg []byte) error {
	q.mu.Lock()
	ch, ok := q.channels[topic]
	if !ok {
		ch = make(chan []byte, 100)
		q.channels[topic] = ch
	}
	q.mu.Unlock()

	select {
	case ch <- msg:
	default:
		return fmt.Errorf("infra: localqueue: topic %q channel full, message dropped", topic)
	}
	return nil
}

// Subscribe starts a goroutine that calls handler for every message on topic.
// The goroutine exits when ctx is cancelled or the channel is closed.
// If the channel does not exist, it is created with a buffer of 100.
func (q *LocalQueue) Subscribe(ctx context.Context, topic string, handler func([]byte)) error {
	q.mu.Lock()
	ch, ok := q.channels[topic]
	if !ok {
		ch = make(chan []byte, 100)
		q.channels[topic] = ch
	}
	q.mu.Unlock()

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case msg, ok := <-ch:
				if !ok {
					return
				}
				handler(msg)
			}
		}
	}()
	return nil
}

// Close closes all topic channels.
func (q *LocalQueue) Close() error {
	q.mu.Lock()
	defer q.mu.Unlock()

	for topic, ch := range q.channels {
		close(ch)
		delete(q.channels, topic)
	}
	return nil
}

// NewMessageQueue creates a MessageQueue implementation based on cfg.Provider.
// "local" returns an in-process LocalQueue.
// "rocketmq" returns a RocketMQ-backed queue (Apache rocketmq-client-go/v2).
func NewMessageQueue(cfg conf.MQConfig) (MessageQueue, error) {
	switch cfg.Provider {
	case "local":
		return &LocalQueue{channels: make(map[string]chan []byte)}, nil
	case "rocketmq":
		return newRocketMQQueue(&cfg)
	default:
		return nil, fmt.Errorf("infra: mq: unknown provider %q", cfg.Provider)
	}
}
