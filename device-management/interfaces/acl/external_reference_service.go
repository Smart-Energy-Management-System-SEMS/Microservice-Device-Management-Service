package acl

import (
	"context"

	shared "device-management-service/shared/domain"
	"github.com/google/uuid"
)

type ExternalReferenceService interface {
	ValidateUserReference(ctx context.Context, userID uuid.UUID) error
	ValidateHomeReference(ctx context.Context, homeID *uuid.UUID) error
}

type LocalExternalReferenceService struct{}

func NewLocalExternalReferenceService() ExternalReferenceService {
	return &LocalExternalReferenceService{}
}

func (s *LocalExternalReferenceService) ValidateUserReference(_ context.Context, userID uuid.UUID) error {
	if userID == uuid.Nil {
		return shared.NewValidationError("user_id external reference is required")
	}
	return nil
}

func (s *LocalExternalReferenceService) ValidateHomeReference(_ context.Context, homeID *uuid.UUID) error {
	if homeID != nil && *homeID == uuid.Nil {
		return shared.NewValidationError("home_id external reference must be a valid UUID")
	}
	return nil
}
