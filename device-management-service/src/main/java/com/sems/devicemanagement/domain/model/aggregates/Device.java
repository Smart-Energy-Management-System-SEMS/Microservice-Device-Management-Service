package com.sems.devicemanagement.domain.model.aggregates;

import com.sems.devicemanagement.domain.model.commands.AddDeviceCommand;
import com.sems.devicemanagement.domain.model.valueobjects.*;
import jakarta.persistence.*;
import lombok.Getter;

@Getter
@Entity
@Table(name = "devices")
public class Device extends AuditableAbstractAggregateRoot {

    @Embedded
    @AttributeOverride(name = "id", column = @Column(name = "user_id"))
    private UserId userId;

    @Embedded
    @AttributeOverride(name = "name", column = @Column(name = "device_name"))
    private DeviceName name;

    @Embedded
    @AttributeOverride(name = "category", column = @Column(name = "device_category"))
    private DeviceCategory category;

    @Embedded
    @AttributeOverride(name = "type", column = @Column(name = "device_type"))
    private DeviceType type;

    @Embedded
    @AttributeOverride(name = "status", column = @Column(name = "device_status"))
    private DeviceStatus status;

    @Embedded
    @AttributeOverride(name = "activity", column = @Column(name = "device_activity"))
    private DeviceActivity activity;

    @Embedded
    @AttributeOverride(name = "location", column = @Column(name = "device_location"))
    private DeviceLocation location;

    @Column(nullable = false)
    private boolean active;

    public Device() {}

    public Device(UserId userId, DeviceName name, DeviceCategory category,
                  DeviceType type, DeviceStatus status, DeviceActivity activity,
                  DeviceLocation location, boolean active) {
        this.userId = userId;
        this.name = name;
        this.category = category;
        this.type = type;
        this.status = status;
        this.activity = activity;
        this.location = location;
        this.active = active;
    }

    public static Device create(AddDeviceCommand command, UserId userId) {
        return new Device(
                userId,
                new DeviceName(command.name()),
                new DeviceCategory(command.category()),
                new DeviceType(command.type()),
                new DeviceStatus(command.status()),
                new DeviceActivity(command.lastActivity()),
                new DeviceLocation(command.location()),
                command.active()
        );
    }

    public void updateStatus(DeviceStatus status) {
        this.status = status;
    }

    public void updateActivity(DeviceActivity activity) {
        this.activity = activity;
    }
}
