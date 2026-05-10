package com.sems.devicemanagement.domain.model.valueobjects;
import jakarta.persistence.Embeddable;
@Embeddable
public record DeviceCategory(String category) {
    public DeviceCategory { if (category == null || category.isBlank()) throw new IllegalArgumentException("Device category cannot be null or blank"); }
}
