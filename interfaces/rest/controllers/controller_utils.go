package controllers

import (
	"device-management-service/shared/domain"
	"github.com/google/uuid"
)

func parseUUID(value string, field string) (uuid.UUID, error) {
	id, err := uuid.Parse(value)
	if err != nil {
		return uuid.Nil, domain.NewValidationError(field + " must be a valid UUID")
	}
	return id, nil
}
