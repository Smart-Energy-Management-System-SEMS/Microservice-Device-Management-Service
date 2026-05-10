package com.sems.devicemanagement.domain.services;

import com.sems.devicemanagement.domain.model.commands.AddDeviceCommand;
import com.sems.devicemanagement.domain.model.commands.DeleteDeviceCommand;
import com.sems.devicemanagement.domain.model.commands.UpdatePreferencesCommand;
import com.sems.devicemanagement.domain.model.valueobjects.UserId;

public interface DeviceCommandService {
    Long handle(AddDeviceCommand command, UserId userId);
    void handle(UpdatePreferencesCommand command);
    void handle(DeleteDeviceCommand command);
}
