package com.sems.devicemanagement.domain.model.valueobjects;
import jakarta.persistence.Embeddable;
@Embeddable
public record DeviceType(String type) {
    public DeviceType { if (type == null || type.isBlank()) throw new IllegalArgumentException("Device type cannot be null or blank"); }
}
