package com.sems.devicemanagement.domain.model.valueobjects;
import jakarta.persistence.Embeddable;
@Embeddable
public record DeviceLocation(String location) {
    public DeviceLocation { if (location == null || location.isBlank()) throw new IllegalArgumentException("Device location cannot be null or blank"); }
}
