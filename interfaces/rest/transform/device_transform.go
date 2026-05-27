package transform

import (
	"device-management-service/domain/model/aggregates"
	"device-management-service/interfaces/rest/resources"
)

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

func ToDeviceResources(devices []aggregates.Device) []resources.DeviceResource {
	result := make([]resources.DeviceResource, 0, len(devices))
	for _, device := range devices {
		item := device
		result = append(result, ToDeviceResource(&item))
	}
	return result
}
