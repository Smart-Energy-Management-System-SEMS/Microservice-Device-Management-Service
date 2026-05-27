package resources

import "time"

type CreateDeviceResource struct {
	ExternalDeviceCode string  `json:"externalDeviceCode" binding:"required"`
	UserID             string  `json:"userId" binding:"required"`
	DeviceName         string  `json:"deviceName" binding:"required"`
	DeviceType         string  `json:"deviceType" binding:"required"`
	Brand              *string `json:"brand"`
	Model              *string `json:"model"`
	ConnectionProtocol string  `json:"connectionProtocol" binding:"required"`
}

type UpdateDeviceResource struct {
	DeviceName         string  `json:"deviceName" binding:"required"`
	DeviceType         string  `json:"deviceType" binding:"required"`
	Brand              *string `json:"brand"`
	Model              *string `json:"model"`
	ConnectionProtocol string  `json:"connectionProtocol" binding:"required"`
}

type UpdateDeviceStatusResource struct {
	Status string `json:"status" binding:"required"`
}

type DeviceResource struct {
	DeviceID           string    `json:"deviceId"`
	ExternalDeviceCode string    `json:"externalDeviceCode"`
	UserID             string    `json:"userId"`
	DeviceName         string    `json:"deviceName"`
	DeviceType         string    `json:"deviceType"`
	Brand              *string   `json:"brand,omitempty"`
	Model              *string   `json:"model,omitempty"`
	ConnectionProtocol string    `json:"connectionProtocol"`
	Status             string    `json:"status"`
	RegisteredAt       time.Time `json:"registeredAt"`
	UpdatedAt          time.Time `json:"updatedAt"`
}
