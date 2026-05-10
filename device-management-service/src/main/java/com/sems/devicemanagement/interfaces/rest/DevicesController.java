package com.sems.devicemanagement.interfaces.rest;

import com.sems.devicemanagement.domain.model.commands.AddDeviceCommand;
import com.sems.devicemanagement.domain.model.commands.DeleteDeviceCommand;
import com.sems.devicemanagement.domain.model.commands.UpdatePreferencesCommand;
import com.sems.devicemanagement.domain.model.queries.*;
import com.sems.devicemanagement.domain.model.valueobjects.UserId;
import com.sems.devicemanagement.domain.services.DeviceCommandService;
import com.sems.devicemanagement.domain.services.DeviceQueryService;
import com.sems.devicemanagement.interfaces.rest.resources.*;
import io.swagger.v3.oas.annotations.Operation;
import io.swagger.v3.oas.annotations.Parameter;
import io.swagger.v3.oas.annotations.responses.ApiResponse;
import io.swagger.v3.oas.annotations.responses.ApiResponses;
import io.swagger.v3.oas.annotations.security.SecurityRequirement;
import io.swagger.v3.oas.annotations.tags.Tag;
import jakarta.validation.Valid;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.security.core.Authentication;
import org.springframework.web.bind.annotation.*;

import java.util.List;
import java.util.Map;

@RestController
@RequestMapping(value = "/api/v1/devices", produces = "application/json")
@Tag(name = "Devices", description = "Gestión de dispositivos IoT del hogar")
@SecurityRequirement(name = "bearerAuth")
public class DevicesController {

    private final DeviceCommandService deviceCommandService;
    private final DeviceQueryService deviceQueryService;

    public DevicesController(DeviceCommandService deviceCommandService,
                             DeviceQueryService deviceQueryService) {
        this.deviceCommandService = deviceCommandService;
        this.deviceQueryService = deviceQueryService;
    }

    @PostMapping
    @Operation(summary = "Vincular nuevo dispositivo",
               description = "Registra un dispositivo IoT y lo vincula al usuario autenticado")
    @ApiResponses({
        @ApiResponse(responseCode = "201", description = "Dispositivo creado exitosamente"),
        @ApiResponse(responseCode = "400", description = "Datos de entrada inválidos"),
        @ApiResponse(responseCode = "401", description = "No autenticado")
    })
    public ResponseEntity<DeviceResource> createDevice(
            @Valid @RequestBody CreateDeviceResource resource,
            Authentication authentication) {

        Long userId = extractUserId(authentication);
        if (userId == null) return ResponseEntity.status(HttpStatus.UNAUTHORIZED).build();

        var command = new AddDeviceCommand(
                resource.name(), resource.category(), resource.type(),
                resource.status(), resource.lastActivity(), resource.location(), resource.active()
        );

        Long deviceId = deviceCommandService.handle(command, new UserId(userId));
        if (deviceId == null) return ResponseEntity.badRequest().build();

        return deviceQueryService.handle(new GetDeviceByIdQuery(deviceId))
                .map(d -> ResponseEntity.status(HttpStatus.CREATED).body(toResource(d)))
                .orElse(ResponseEntity.notFound().build());
    }

    @GetMapping
    @Operation(summary = "Listar todos los dispositivos",
               description = "Retorna todos los dispositivos del sistema (admin)")
    public ResponseEntity<List<DeviceResource>> getAllDevices() {
        var devices = deviceQueryService.handle(new GetAllDevicesQuery())
                .stream().map(this::toResource).toList();
        return ResponseEntity.ok(devices);
    }

    @GetMapping("/{deviceId}")
    @Operation(summary = "Obtener dispositivo por ID")
    @ApiResponses({
        @ApiResponse(responseCode = "200", description = "Dispositivo encontrado"),
        @ApiResponse(responseCode = "404", description = "Dispositivo no encontrado")
    })
    public ResponseEntity<DeviceResource> getDeviceById(
            @Parameter(description = "ID del dispositivo") @PathVariable Long deviceId) {
        return deviceQueryService.handle(new GetDeviceByIdQuery(deviceId))
                .map(d -> ResponseEntity.ok(toResource(d)))
                .orElse(ResponseEntity.notFound().build());
    }

    @GetMapping("/my-devices")
    @Operation(summary = "Mis dispositivos",
               description = "Retorna los dispositivos del usuario autenticado")
    public ResponseEntity<List<DeviceResource>> getMyDevices(Authentication authentication) {
        Long userId = extractUserId(authentication);
        if (userId == null) return ResponseEntity.status(HttpStatus.UNAUTHORIZED).build();

        var devices = deviceQueryService.handle(new GetDevicesByUserIdQuery(new UserId(userId)))
                .stream().map(this::toResource).toList();
        return ResponseEntity.ok(devices);
    }

    @GetMapping("/users/{userId}")
    @Operation(summary = "Dispositivos por usuario", description = "Retorna dispositivos de un usuario específico")
    public ResponseEntity<List<DeviceResource>> getDevicesByUserId(@PathVariable Long userId) {
        var devices = deviceQueryService.handle(new GetDevicesByUserIdQuery(new UserId(userId)))
                .stream().map(this::toResource).toList();
        return ResponseEntity.ok(devices);
    }

    @DeleteMapping("/{deviceId}")
    @Operation(summary = "Desvincular dispositivo", description = "Elimina un dispositivo del sistema")
    @ApiResponses({
        @ApiResponse(responseCode = "200", description = "Dispositivo eliminado"),
        @ApiResponse(responseCode = "404", description = "Dispositivo no encontrado")
    })
    public ResponseEntity<Map<String, String>> deleteDevice(@PathVariable Long deviceId) {
        deviceCommandService.handle(new DeleteDeviceCommand(deviceId));
        return ResponseEntity.ok(Map.of("message", "Dispositivo eliminado correctamente"));
    }

    @GetMapping("/users/{userId}/preferences")
    @Operation(summary = "Obtener preferencias del usuario")
    public ResponseEntity<PreferencesResource> getPreferences(@PathVariable Long userId) {
        return deviceQueryService.handle(new GetPreferencesByUserIdQuery(userId))
                .map(p -> ResponseEntity.ok(new PreferencesResource(
                        p.getId(), String.valueOf(p.getUserId()), p.getThreshold(),
                        p.isNotificationEnabled(), p.isHabilitarMonitoreoEnergia(),
                        p.isRecibirAlertasAltoConsumo(), p.isMonitorearCalefaccionRefrigeracion(),
                        p.isMonitorearElectrodomesticosPrincipales(), p.isMonitorearElectronicos(),
                        p.isMonitorearDispositivosCocina(), p.isIncluirIluminacionExterior(),
                        p.isRastrearEnergiaEspera(), p.isEmailsResumenDiario(),
                        p.isReportesProgresoSemanal(), p.isSugerirAutomizacionesAhorro(),
                        p.isAlertasDispositivosDesconectados()
                )))
                .orElse(ResponseEntity.notFound().build());
    }

    @PutMapping("/users/{userId}/preferences")
    @Operation(summary = "Actualizar preferencias del usuario")
    public ResponseEntity<Map<String, String>> updatePreferences(
            @PathVariable Long userId,
            @RequestBody UpdatePreferencesResource resource) {

        var command = new UpdatePreferencesCommand(
                userId,
                resource.habilitarMonitoreoEnergia(), resource.recibirAlertasAltoConsumo(),
                resource.monitorearCalefaccionRefrigeracion(), resource.monitorearElectrodomesticosPrincipales(),
                resource.monitorearElectronicos(), resource.monitorearDispositivosCocina(),
                resource.incluirIluminacionExterior(), resource.rastrearEnergiaEspera(),
                resource.emailsResumenDiario(), resource.reportesProgresoSemanal(),
                resource.sugerirAutomizacionesAhorro(), resource.alertasDispositivosDesconectados()
        );
        deviceCommandService.handle(command);
        return ResponseEntity.ok(Map.of("message", "Preferencias actualizadas correctamente"));
    }

    // ---- Helpers ----

    private DeviceResource toResource(com.sems.devicemanagement.domain.model.aggregates.Device d) {
        return new DeviceResource(
                d.getId(),
                String.valueOf(d.getUserId().id()),
                d.getName().name(),
                d.getCategory().category(),
                d.getType().type(),
                d.getStatus().status(),
                d.getActivity().activity(),
                d.getLocation().location(),
                d.isActive()
        );
    }

    private Long extractUserId(Authentication authentication) {
        if (authentication == null || authentication.getPrincipal() == null) return null;
        // The principal name is the email; userId comes from JWT claims attribute
        Object principal = authentication.getPrincipal();
        if (principal instanceof org.springframework.security.core.userdetails.UserDetails ud) {
            // In this microservice, userId is stored as a claim in the JWT - see JwtAuthFilter
            var details = authentication.getDetails();
            if (details instanceof Long id) return id;
        }
        return null;
    }
}
