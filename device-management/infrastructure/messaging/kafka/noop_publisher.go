package kafka

import (
	"context"

	"device-management-service/device-management/application/outboundservices"
)

type NoopPublisher struct{}

func NewNoopPublisher() *NoopPublisher {
	return &NoopPublisher{}
}

func (n *NoopPublisher) Publish(_ context.Context, _ string, _ outboundservices.IntegrationEvent) error {
	return nil
}
