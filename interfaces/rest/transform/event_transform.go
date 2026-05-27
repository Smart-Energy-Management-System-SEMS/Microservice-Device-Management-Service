package transform

import (
	"device-management-service/domain/model/entities"
	"device-management-service/interfaces/rest/resources"
)

func ToEventResource(event *entities.DeviceEvent) resources.DeviceEventResource {
	return resources.DeviceEventResource{
		EventID:     event.EventID.String(),
		DeviceID:    event.DeviceID.String(),
		EventType:   event.EventType,
		Description: event.Description,
		OccurredAt:  event.OccurredAt,
	}
}

func ToEventResources(events []entities.DeviceEvent) []resources.DeviceEventResource {
	result := make([]resources.DeviceEventResource, 0, len(events))
	for _, event := range events {
		item := event
		result = append(result, ToEventResource(&item))
	}
	return result
}
