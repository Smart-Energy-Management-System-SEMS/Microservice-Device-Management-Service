package com.sems.devicemanagement.domain.model.commands;
public record UpdatePreferencesCommand(
        Long userId,
        boolean habilitarMonitoreoEnergia,
        boolean recibirAlertasAltoConsumo,
        boolean monitorearCalefaccionRefrigeracion,
        boolean monitorearElectrodomesticosPrincipales,
        boolean monitorearElectronicos,
        boolean monitorearDispositivosCocina,
        boolean incluirIluminacionExterior,
        boolean rastrearEnergiaEspera,
        boolean emailsResumenDiario,
        boolean reportesProgresoSemanal,
        boolean sugerirAutomizacionesAhorro,
        boolean alertasDispositivosDesconectados
) {}
