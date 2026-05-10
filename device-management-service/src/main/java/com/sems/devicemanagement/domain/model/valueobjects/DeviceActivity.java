package com.sems.devicemanagement.domain.model.valueobjects;
import jakarta.persistence.Embeddable;
@Embeddable
public record DeviceActivity(String activity) {
    public DeviceActivity { if (activity == null || activity.isBlank()) throw new IllegalArgumentException("Device activity cannot be null or blank"); }
}
