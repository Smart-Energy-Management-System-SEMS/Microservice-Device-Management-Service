package kafka

import (
	"context"
	"log"

	"github.com/segmentio/kafka-go"
)

type MessageHandler func(ctx context.Context, message kafka.Message) error

type Consumer struct {
	reader  *kafka.Reader
	handler MessageHandler
}

func NewConsumer(brokers []string, groupID string, topic string, handler MessageHandler) *Consumer {
	return &Consumer{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers: brokers,
			GroupID: groupID,
			Topic:   topic,
		}),
		handler: handler,
	}
}

func (c *Consumer) Start(ctx context.Context) {
	if c == nil || c.reader == nil || c.handler == nil {
		return
	}
	go func() {
		defer c.reader.Close()
		for {
			message, err := c.reader.ReadMessage(ctx)
			if err != nil {
				if ctx.Err() != nil {
					return
				}
				log.Printf("kafka read failed: %v", err)
				continue
			}
			if err := c.handler(ctx, message); err != nil {
				log.Printf("kafka message handler failed: %v", err)
			}
		}
	}()
}
