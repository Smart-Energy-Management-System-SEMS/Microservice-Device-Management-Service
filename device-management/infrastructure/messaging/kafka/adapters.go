package kafka

import "device-management-service/device-management/application/outboundservices"

var _ outboundservices.DeviceEventPublisher = (*Producer)(nil)
