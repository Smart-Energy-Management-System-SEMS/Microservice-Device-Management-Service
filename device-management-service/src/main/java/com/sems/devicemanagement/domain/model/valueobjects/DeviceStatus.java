package com.sems.devicemanagement.domain.model.valueobjects;
import jakarta.persistence.Embeddable;
@Embeddable
public record DeviceStatus(String status) {
    public DeviceStatus { if (status == null || status.isBlank()) throw new IllegalArgumentException("Device status cannot be null or blank"); }
}
