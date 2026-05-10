package com.sems.devicemanagement.interfaces.rest.resources;
public record PreferencesResource(
        Long id, String userId, Double threshold, boolean notificationEnabled,
        boolean habilitarMonitoreoEnergia, boolean recibirAlertasAltoConsumo,
        boolean monitorearCalefaccionRefrigeracion, boolean monitorearElectrodomesticosPrincipales,
        boolean monitorearElectronicos, boolean monitorearDispositivosCocina,
        boolean incluirIluminacionExterior, boolean rastrearEnergiaEspera,
        boolean emailsResumenDiario, boolean reportesProgresoSemanal,
        boolean sugerirAutomizacionesAhorro, boolean alertasDispositivosDesconectados
) {}
