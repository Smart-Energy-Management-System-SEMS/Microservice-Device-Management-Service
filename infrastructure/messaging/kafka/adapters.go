package kafka

import "device-management-service/application/outboundservices"

var _ outboundservices.DeviceEventPublisher = (*Producer)(nil)
