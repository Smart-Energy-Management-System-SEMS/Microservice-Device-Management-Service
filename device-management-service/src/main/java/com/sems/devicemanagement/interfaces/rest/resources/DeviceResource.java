package com.sems.devicemanagement.interfaces.rest.resources;
public record DeviceResource(Long id, String userId, String name, String category, String type, String status, String lastActivity, String location, boolean active) {}
