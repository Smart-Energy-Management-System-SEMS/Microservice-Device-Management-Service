// Package kafka is the infrastructure "adapter" that actually talks to Apache
// Kafka. It implements the DeviceEventPublisher port defined in the application
// layer, so the rest of the code never imports kafka-go directly. If we ever
// switched message brokers, only this package would change.
package kafka

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"device-management-service/device-management/application/outboundservices"
	"github.com/segmentio/kafka-go"
)

// Producer sends messages to Kafka. It keeps one writer per topic in a map and
// reuses them, because creating a writer is relatively expensive. Because the
// service handles many requests at the same time (concurrently), the map is
// protected by a mutex to avoid two goroutines corrupting it at once.
type Producer struct {
	brokers      []string                 // addresses of the Kafka brokers
	clientID     string                   // identifies this service to Kafka
	writeTimeout time.Duration            // max time allowed for a single write
	writersMu    sync.Mutex               // guards the writers map below
	writers      map[string]*kafka.Writer // one cached writer per topic
}

// NewProducer creates a Producer with an empty writer cache. The map MUST be
// initialised with make here; writing to a nil map would panic at runtime.
func NewProducer(brokers []string, clientID string, writeTimeout time.Duration) *Producer {
	return &Producer{
		brokers:      brokers,
		clientID:     clientID,
		writeTimeout: writeTimeout,
		writers:      make(map[string]*kafka.Writer),
	}
}

// Publish implements the DeviceEventPublisher interface. It converts the event
// to JSON and writes it to the given topic.
func (p *Producer) Publish(ctx context.Context, topic string, event outboundservices.IntegrationEvent) error {
	// Marshal turns the Go struct into the JSON bytes that travel on the wire.
	payload, err := json.Marshal(event)
	if err != nil {
		return err
	}
	writer := p.writerForTopic(topic)
	// We derive a child context with a timeout so a stuck broker cannot make
	// this call hang forever. cancel() must be deferred to free its resources.
	writeCtx, cancel := context.WithTimeout(ctx, p.writeTimeout)
	defer cancel()

	// The message Key is the device ID. Kafka uses the key to decide the
	// partition, so all events for the same device land on the same partition
	// and therefore keep their order. Headers carry small metadata that
	// consumers can read without parsing the whole JSON body.
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

// Close shuts down every cached writer, typically during graceful shutdown of
// the service. It tries to close them all and remembers the FIRST error rather
// than stopping early, so one bad writer does not leave the others open.
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

// writerForTopic returns the cached writer for a topic, creating it on first
// use ("lazy initialisation"). The lock makes the whole check-then-create
// sequence atomic, so two concurrent calls cannot both build a writer for the
// same topic. The comma-ok form (writer, ok := map[key]) tells us whether the
// key was already present.
func (p *Producer) writerForTopic(topic string) *kafka.Writer {
	p.writersMu.Lock()
	defer p.writersMu.Unlock()

	if writer, ok := p.writers[topic]; ok {
		return writer
	}

	// LeastBytes balancer spreads messages toward the partition with the least
	// data; BatchTimeout lets the client group messages for ~10ms for
	// efficiency before sending.
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers:      p.brokers,
		Topic:        topic,
		Balancer:     &kafka.LeastBytes{},
		BatchTimeout: 10 * time.Millisecond,
	})
	p.writers[topic] = writer
	return writer
}
