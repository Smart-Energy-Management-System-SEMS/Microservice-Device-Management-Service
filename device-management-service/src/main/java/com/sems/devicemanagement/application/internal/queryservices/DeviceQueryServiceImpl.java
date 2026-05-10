package com.sems.devicemanagement.application.internal.queryservices;

import com.sems.devicemanagement.domain.model.aggregates.Device;
import com.sems.devicemanagement.domain.model.entities.DevicePreference;
import com.sems.devicemanagement.domain.model.queries.*;
import com.sems.devicemanagement.domain.services.DeviceQueryService;
import com.sems.devicemanagement.infrastructure.persistence.jpa.repositories.DeviceRepository;
import com.sems.devicemanagement.infrastructure.persistence.jpa.repositories.PreferencesRepository;
import org.springframework.stereotype.Service;

import java.util.List;
import java.util.Optional;

@Service
public class DeviceQueryServiceImpl implements DeviceQueryService {

    private final DeviceRepository deviceRepository;
    private final PreferencesRepository preferencesRepository;

    public DeviceQueryServiceImpl(DeviceRepository deviceRepository,
                                  PreferencesRepository preferencesRepository) {
        this.deviceRepository = deviceRepository;
        this.preferencesRepository = preferencesRepository;
    }

    @Override
    public List<Device> handle(GetAllDevicesQuery query) {
        return deviceRepository.findAll();
    }

    @Override
    public Optional<Device> handle(GetDeviceByIdQuery query) {
        return deviceRepository.findById(query.deviceId());
    }

    @Override
    public List<Device> handle(GetDevicesByUserIdQuery query) {
        return deviceRepository.findByUserId(query.userId());
    }

    @Override
    public Optional<DevicePreference> handle(GetPreferencesByUserIdQuery query) {
        return preferencesRepository.findByUserId(query.userId());
    }
}
