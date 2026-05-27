package valueobjects

import dmerrors "device-management-service/device-management/domain"

type BindingStatus string

const (
	BindingStatusLinked   BindingStatus = "LINKED"
	BindingStatusUnlinked BindingStatus = "UNLINKED"
	BindingStatusPending  BindingStatus = "PENDING"
)

func NewBindingStatus(value string) (BindingStatus, error) {
	status := BindingStatus(value)
	if !status.IsValid() {
		return "", dmerrors.NewValidationError("invalid binding status")
	}
	return status, nil
}

func (s BindingStatus) IsValid() bool {
	return s == BindingStatusLinked || s == BindingStatusUnlinked || s == BindingStatusPending
}
