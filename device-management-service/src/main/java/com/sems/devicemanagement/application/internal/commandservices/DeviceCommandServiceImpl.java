package com.sems.devicemanagement.application.internal.commandservices;

import com.sems.devicemanagement.domain.model.aggregates.Device;
import com.sems.devicemanagement.domain.model.commands.AddDeviceCommand;
import com.sems.devicemanagement.domain.model.commands.DeleteDeviceCommand;
import com.sems.devicemanagement.domain.model.commands.UpdatePreferencesCommand;
import com.sems.devicemanagement.domain.model.entities.DevicePreference;
import com.sems.devicemanagement.domain.model.valueobjects.UserId;
import com.sems.devicemanagement.domain.services.DeviceCommandService;
import com.sems.devicemanagement.infrastructure.persistence.jpa.repositories.DeviceRepository;
import com.sems.devicemanagement.infrastructure.persistence.jpa.repositories.PreferencesRepository;
import jakarta.persistence.EntityNotFoundException;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

@Service
public class DeviceCommandServiceImpl implements DeviceCommandService {

    private final DeviceRepository deviceRepository;
    private final PreferencesRepository preferencesRepository;

    public DeviceCommandServiceImpl(DeviceRepository deviceRepository,
                                    PreferencesRepository preferencesRepository) {
        this.deviceRepository = deviceRepository;
        this.preferencesRepository = preferencesRepository;
    }

    @Override
    public Long handle(AddDeviceCommand command, UserId userId) {
        Device device = Device.create(command, userId);
        device = deviceRepository.save(device);

        // Create default preferences if they don't exist for this user
        if (!preferencesRepository.existsByUserId(userId.id())) {
            DevicePreference preferences = new DevicePreference(
                    userId.id(), 100.0, true,
                    false, false, false, false, false, false, false, false, false, false, false, false
            );
            preferencesRepository.save(preferences);
        }

        return device.getId();
    }

    @Override
    public void handle(UpdatePreferencesCommand command) {
        DevicePreference preferences = preferencesRepository
                .findByUserId(command.userId())
                .orElseGet(() -> new DevicePreference(
                        command.userId(), 100.0, true,
                        command.habilitarMonitoreoEnergia(),
                        command.recibirAlertasAltoConsumo(),
                        command.monitorearCalefaccionRefrigeracion(),
                        command.monitorearElectrodomesticosPrincipales(),
                        command.monitorearElectronicos(),
                        command.monitorearDispositivosCocina(),
                        command.incluirIluminacionExterior(),
                        command.rastrearEnergiaEspera(),
                        command.emailsResumenDiario(),
                        command.reportesProgresoSemanal(),
                        command.sugerirAutomizacionesAhorro(),
                        command.alertasDispositivosDesconectados()
                ));

        if (preferences.getId() != null) {
            preferences.updateHabilitarMonitoreoEnergia(command.habilitarMonitoreoEnergia());
            preferences.updateRecibirAlertasAltoConsumo(command.recibirAlertasAltoConsumo());
            preferences.updateMonitorearCalefaccionRefrigeracion(command.monitorearCalefaccionRefrigeracion());
            preferences.updateMonitorearElectrodomesticosPrincipales(command.monitorearElectrodomesticosPrincipales());
            preferences.updateMonitorearElectronicos(command.monitorearElectronicos());
            preferences.updateMonitorearDispositivosCocina(command.monitorearDispositivosCocina());
            preferences.updateIncluirIluminacionExterior(command.incluirIluminacionExterior());
            preferences.updateRastrearEnergiaEspera(command.rastrearEnergiaEspera());
            preferences.updateEmailsResumenDiario(command.emailsResumenDiario());
            preferences.updateReportesProgresoSemanal(command.reportesProgresoSemanal());
            preferences.updateSugerirAutomizacionesAhorro(command.sugerirAutomizacionesAhorro());
            preferences.updateAlertasDispositivosDesconectados(command.alertasDispositivosDesconectados());
        }

        preferencesRepository.save(preferences);
    }

    @Override
    @Transactional
    public void handle(DeleteDeviceCommand command) {
        if (!deviceRepository.existsById(command.deviceId())) {
            throw new EntityNotFoundException("Device with id " + command.deviceId() + " not found");
        }
        deviceRepository.deleteById(command.deviceId());
    }
}
