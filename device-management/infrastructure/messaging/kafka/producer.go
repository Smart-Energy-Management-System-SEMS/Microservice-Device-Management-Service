package kafka

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"device-management-service/device-management/application/outboundservices"
	"github.com/segmentio/kafka-go"
)

type Producer struct {
	brokers      []string
	clientID     string
	writeTimeout time.Duration
	writersMu    sync.Mutex
	writers      map[string]*kafka.Writer
}

func NewProducer(brokers []string, clientID string, writeTimeout time.Duration) *Producer {
	return &Producer{
		brokers:      brokers,
		clientID:     clientID,
		writeTimeout: writeTimeout,
		writers:      make(map[string]*kafka.Writer),
	}
}

func (p *Producer) Publish(ctx context.Context, topic string, event outboundservices.IntegrationEvent) error {
	payload, err := json.Marshal(event)
	if err != nil {
		return err
	}
	writer := p.writerForTopic(topic)
	writeCtx, cancel := context.WithTimeout(ctx, p.writeTimeout)
	defer cancel()

	return writer.WriteMessages(writeCtx, kafka.Message{
		Key:   []byte(event.DeviceID.String()),
		Value: payload,
		Headers: []kafka.Header{
			{Key: "client_id", Value: []byte(p.clientID)},
			{Key: "event_type", Value: []byte(event.EventType)},
		},
		Time: time.Now().UTC(),
	})
}

func (p *Producer) Close() error {
	p.writersMu.Lock()
	defer p.writersMu.Unlock()

	var closeErr error
	for topic, writer := range p.writers {
		if writer == nil {
			continue
		}
		if err := writer.Close(); err != nil && closeErr == nil {
			closeErr = err
		}
		delete(p.writers, topic)
	}
	return closeErr
}

func (p *Producer) writerForTopic(topic string) *kafka.Writer {
	p.writersMu.Lock()
	defer p.writersMu.Unlock()

	if writer, ok := p.writers[topic]; ok {
		return writer
	}

	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers:      p.brokers,
		Topic:        topic,
		Balancer:     &kafka.LeastBytes{},
		BatchTimeout: 10 * time.Millisecond,
	})
	p.writers[topic] = writer
	return writer
}
