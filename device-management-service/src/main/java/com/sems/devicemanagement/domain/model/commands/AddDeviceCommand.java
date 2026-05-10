package com.sems.devicemanagement.domain.model.commands;

public record AddDeviceCommand(
        String name,
        String category,
        String type,
        String status,
        String lastActivity,
        String location,
        boolean active
) {}
