package com.sems.devicemanagement.interfaces.rest.resources;
import io.swagger.v3.oas.annotations.media.Schema;
import jakarta.validation.constraints.NotBlank;
@Schema(description = "Payload para crear un nuevo dispositivo IoT")
public record CreateDeviceResource(
        @NotBlank @Schema(example = "Smart TV Samsung") String name,
        @NotBlank @Schema(example = "Electrónica") String category,
        @NotBlank @Schema(example = "Televisor") String type,
        @NotBlank @Schema(example = "connected") String status,
        @NotBlank @Schema(example = "high") String lastActivity,
        @NotBlank @Schema(example = "Sala") String location,
        @Schema(example = "true") boolean active
) {}
