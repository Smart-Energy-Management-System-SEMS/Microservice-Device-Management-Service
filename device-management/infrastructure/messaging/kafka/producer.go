package kafka

import (
	"context"
	"encoding/json"
	"time"

	"device-management-service/device-management/application/outboundservices"
	"github.com/segmentio/kafka-go"
)

type Producer struct {
	brokers  []string
	clientID string
}

func NewProducer(brokers []string, clientID string) *Producer {
	return &Producer{brokers: brokers, clientID: clientID}
}

func (p *Producer) Publish(ctx context.Context, topic string, event outboundservices.IntegrationEvent) error {
	payload, err := json.Marshal(event)
	if err != nil {
		return err
	}
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers:      p.brokers,
		Topic:        topic,
		Balancer:     &kafka.LeastBytes{},
		BatchTimeout: 10 * time.Millisecond,
	})
	defer writer.Close()

	return writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte(event.DeviceID.String()),
		Value: payload,
		Headers: []kafka.Header{
			{Key: "client_id", Value: []byte(p.clientID)},
			{Key: "event_type", Value: []byte(event.EventType)},
		},
		Time: time.Now().UTC(),
	})
}
