package controllers

import (
	dmerrors "device-management-service/device-management/domain"
	"github.com/google/uuid"
)

func parseUUID(value string, field string) (uuid.UUID, error) {
	id, err := uuid.Parse(value)
	if err != nil {
		return uuid.Nil, dmerrors.NewValidationError(field + " must be a valid UUID")
	}
	return id, nil
}
