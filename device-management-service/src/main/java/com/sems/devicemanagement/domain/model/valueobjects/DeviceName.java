package com.sems.devicemanagement.domain.model.valueobjects;

import jakarta.persistence.Embeddable;

@Embeddable
public record DeviceName(String name) {
    public DeviceName {
        if (name == null || name.isBlank())
            throw new IllegalArgumentException("Device name cannot be null or blank");
    }
}
