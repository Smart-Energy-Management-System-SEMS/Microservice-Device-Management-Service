// Package transform converts domain objects into "resources" (also called DTOs:
// Data Transfer Objects) for the API responses. This is the outgoing twin of
// the persistence mappers: just like we keep DB models separate from the domain,
// we keep the API shape separate too. That way we decide exactly what the
// outside world sees, and an internal refactor of the domain does not silently
// change our public JSON contract.
package transform

import (
	"device-management-service/device-management/domain/model/aggregates"
	"device-management-service/device-management/interfaces/rest/resources"
)

// ToDeviceResource maps a single domain Device to its API resource. UUIDs and
// value objects are turned into plain strings (.String(), string(...)) because
// that is the friendliest form for JSON clients to consume.
func ToDeviceResource(device *aggregates.Device) resources.DeviceResource {
	return resources.DeviceResource{
		DeviceID:           device.DeviceID.String(),
		ExternalDeviceCode: device.ExternalDeviceCode,
		UserID:             device.UserID.String(),
		DeviceName:         device.DeviceName,
		DeviceType:         device.DeviceType,
		Brand:              device.Brand,
		Model:              device.Model,
		ConnectionProtocol: string(device.ConnectionProtocol),
		Status:             string(device.Status),
		RegisteredAt:       device.RegisteredAt,
		UpdatedAt:          device.UpdatedAt,
	}
}

// ToDeviceResources maps a whole slice by reusing the single-item function.
// The "item := device" line is important: it makes a fresh copy each iteration
// so that &item points at this element, not at the loop variable being reused.
// Taking the address of the range variable directly would be a classic Go bug.
func ToDeviceResources(devices []aggregates.Device) []resources.DeviceResource {
	result := make([]resources.DeviceResource, 0, len(devices))
	for _, device := range devices {
		item := device
		result = append(result, ToDeviceResource(&item))
	}
	return result
}
