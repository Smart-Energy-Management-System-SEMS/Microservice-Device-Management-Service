package resources

import "time"

type CreateDeviceBindingResource struct {
	UserID string  `json:"userId" binding:"required"`
	HomeID *string `json:"homeId"`
}

type DeviceBindingResource struct {
	BindingID     string     `json:"bindingId"`
	DeviceID      string     `json:"deviceId"`
	UserID        string     `json:"userId"`
	HomeID        *string    `json:"homeId,omitempty"`
	BindingStatus string     `json:"bindingStatus"`
	LinkedAt      time.Time  `json:"linkedAt"`
	UnlinkedAt    *time.Time `json:"unlinkedAt,omitempty"`
	UpdatedAt     time.Time  `json:"updatedAt"`
}
