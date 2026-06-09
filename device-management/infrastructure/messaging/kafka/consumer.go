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

func NewConsumer(options ConnectionOptions, topic string, handler MessageHandler) (*Consumer, error) {
	dialer, err := options.Dialer()
	if err != nil {
		return nil, err
	}
	return &Consumer{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers: options.Brokers,
			GroupID: options.ConsumerGroup,
			Topic:   topic,
			Dialer:  dialer,
		}),
		handler: handler,
	}, nil
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
