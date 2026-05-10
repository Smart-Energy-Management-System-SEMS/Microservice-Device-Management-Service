package com.sems.devicemanagement.domain.services;

import com.sems.devicemanagement.domain.model.aggregates.Device;
import com.sems.devicemanagement.domain.model.entities.DevicePreference;
import com.sems.devicemanagement.domain.model.queries.*;

import java.util.List;
import java.util.Optional;

public interface DeviceQueryService {
    List<Device> handle(GetAllDevicesQuery query);
    Optional<Device> handle(GetDeviceByIdQuery query);
    List<Device> handle(GetDevicesByUserIdQuery query);
    Optional<DevicePreference> handle(GetPreferencesByUserIdQuery query);
}
