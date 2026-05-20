package infra

import (
	"context"
	"fmt"
	"strings"
	"sync"

	rmq "github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"

	"github.com/huodaoshi/harness/backend/conf"
)

// RocketMQQueue implements MessageQueue using Apache RocketMQ (sync publish + push consume).
type RocketMQQueue struct {
	ns              primitive.NamesrvAddr
	consumerGroup   string
	producerGroup   string
	mu              sync.Mutex
	producer        rmq.Producer
	pushConsumer    rmq.PushConsumer
	consumerStarted bool
}

func newRocketMQQueue(cfg *conf.MQConfig) (*RocketMQQueue, error) {
	raw := strings.TrimSpace(cfg.NameServer)
	if raw == "" {
		return nil, fmt.Errorf("infra: mq: rocketmq requires mq.name_server or ROCKETMQ_NAME_SERVER")
	}
	ns, err := primitive.NewNamesrvAddr(raw)
	if err != nil {
		return nil, fmt.Errorf("infra: mq: rocketmq name_server: %w", err)
	}

	prodGroup := strings.TrimSpace(cfg.ProducerGroup)
	if prodGroup == "" {
		prodGroup = "einoagent-producer"
	}
	consGroup := strings.TrimSpace(cfg.Group)

	return &RocketMQQueue{
		ns:            ns,
		consumerGroup: consGroup,
		producerGroup: prodGroup,
	}, nil
}

// Publish sends a single message to the topic (lazy-starts producer).
func (r *RocketMQQueue) Publish(ctx context.Context, topic string, msg []byte) error {
	if err := r.ensureProducer(); err != nil {
		return err
	}
	r.mu.Lock()
	p := r.producer
	r.mu.Unlock()

	m := primitive.NewMessage(topic, msg)
	_, err := p.SendSync(ctx, m)
	if err != nil {
		return fmt.Errorf("infra: mq: rocketmq publish topic=%q: %w", topic, err)
	}
	return nil
}

func (r *RocketMQQueue) ensureProducer() error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.producer != nil {
		return nil
	}
	p, err := rmq.NewProducer(
		producer.WithNameServer(r.ns),
		producer.WithGroupName(r.producerGroup),
	)
	if err != nil {
		return fmt.Errorf("infra: mq: rocketmq new producer: %w", err)
	}
	if err := p.Start(); err != nil {
		return fmt.Errorf("infra: mq: rocketmq producer start: %w", err)
	}
	r.producer = p
	return nil
}

// Subscribe registers a push consumer for topic and invokes handler for each message batch.
// consumer group comes from mq.group (must differ per consumer app when running multiple workers).
func (r *RocketMQQueue) Subscribe(ctx context.Context, topic string, handler func([]byte)) error {
	if strings.TrimSpace(r.consumerGroup) == "" {
		return fmt.Errorf("infra: mq: rocketmq Subscribe requires mq.group (consumer group), e.g. knowledgeindexing")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if r.pushConsumer == nil {
		pc, err := rmq.NewPushConsumer(
			consumer.WithNameServer(r.ns),
			consumer.WithGroupName(r.consumerGroup),
		)
		if err != nil {
			return fmt.Errorf("infra: mq: rocketmq new consumer: %w", err)
		}
		r.pushConsumer = pc
	}

	if err := r.pushConsumer.Subscribe(topic, consumer.MessageSelector{},
		func(_ context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
			for _, m := range msgs {
				body := append([]byte(nil), m.Body...)
				handler(body)
			}
			return consumer.ConsumeSuccess, nil
		}); err != nil {
		return fmt.Errorf("infra: mq: rocketmq subscribe %q: %w", topic, err)
	}

	if !r.consumerStarted {
		if err := r.pushConsumer.Start(); err != nil {
			return fmt.Errorf("infra: mq: rocketmq consumer start: %w", err)
		}
		r.consumerStarted = true
	}

	_ = ctx // graceful stop: caller should Close() the queue (consumer Shutdown).
	return nil
}

// Close shuts down producer and consumer.
func (r *RocketMQQueue) Close() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	var firstErr error
	set := func(err error) {
		if err != nil && firstErr == nil {
			firstErr = err
		}
	}

	if r.producer != nil {
		set(r.producer.Shutdown())
		r.producer = nil
	}
	if r.pushConsumer != nil {
		set(r.pushConsumer.Shutdown())
		r.pushConsumer = nil
		r.consumerStarted = false
	}
	return firstErr
}
