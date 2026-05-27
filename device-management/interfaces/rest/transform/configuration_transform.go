package transform

import (
	"device-management-service/device-management/domain/model/entities"
	"device-management-service/device-management/interfaces/rest/resources"
)

func ToConfigurationResource(configuration *entities.DeviceConfiguration) resources.DeviceConfigurationResource {
	return resources.DeviceConfigurationResource{
		ConfigurationID: configuration.ConfigurationID.String(),
		DeviceID:        configuration.DeviceID.String(),
		ConfigKey:       configuration.ConfigKey,
		ConfigValue:     configuration.ConfigValue,
		UpdatedAt:       configuration.UpdatedAt,
	}
}

func ToConfigurationResources(configurations []entities.DeviceConfiguration) []resources.DeviceConfigurationResource {
	result := make([]resources.DeviceConfigurationResource, 0, len(configurations))
	for _, configuration := range configurations {
		item := configuration
		result = append(result, ToConfigurationResource(&item))
	}
	return result
}
