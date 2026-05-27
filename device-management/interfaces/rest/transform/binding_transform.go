package transform

import (
	"device-management-service/device-management/domain/model/entities"
	"device-management-service/device-management/interfaces/rest/resources"
)

func ToBindingResource(binding *entities.DeviceBinding) resources.DeviceBindingResource {
	var homeID *string
	if binding.HomeID != nil {
		value := binding.HomeID.String()
		homeID = &value
	}
	return resources.DeviceBindingResource{
		BindingID:     binding.BindingID.String(),
		DeviceID:      binding.DeviceID.String(),
		UserID:        binding.UserID.String(),
		HomeID:        homeID,
		BindingStatus: string(binding.BindingStatus),
		LinkedAt:      binding.LinkedAt,
		UnlinkedAt:    binding.UnlinkedAt,
		UpdatedAt:     binding.UpdatedAt,
	}
}

func ToBindingResources(bindings []entities.DeviceBinding) []resources.DeviceBindingResource {
	result := make([]resources.DeviceBindingResource, 0, len(bindings))
	for _, binding := range bindings {
		item := binding
		result = append(result, ToBindingResource(&item))
	}
	return result
}
