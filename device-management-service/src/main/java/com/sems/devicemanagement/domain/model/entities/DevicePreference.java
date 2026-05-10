package com.sems.devicemanagement.domain.model.entities;

import jakarta.persistence.*;
import jakarta.validation.constraints.NotNull;
import lombok.Getter;
import org.springframework.data.annotation.CreatedDate;
import org.springframework.data.annotation.LastModifiedDate;
import org.springframework.data.jpa.domain.support.AuditingEntityListener;

import java.util.Date;

@Getter
@Entity
@Table(name = "device_preferences")
@EntityListeners(AuditingEntityListener.class)
public class DevicePreference {

    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long id;

    @Column(name = "user_id", unique = true, nullable = false)
    @NotNull
    private Long userId;

    @Column(nullable = false)
    private Double threshold;

    @Column(nullable = false)
    private boolean notificationEnabled;

    @Column(nullable = false)
    private boolean habilitarMonitoreoEnergia;

    @Column(nullable = false)
    private boolean recibirAlertasAltoConsumo;

    @Column(nullable = false)
    private boolean monitorearCalefaccionRefrigeracion;

    @Column(nullable = false)
    private boolean monitorearElectrodomesticosPrincipales;

    @Column(nullable = false)
    private boolean monitorearElectronicos;

    @Column(nullable = false)
    private boolean monitorearDispositivosCocina;

    @Column(nullable = false)
    private boolean incluirIluminacionExterior;

    @Column(nullable = false)
    private boolean rastrearEnergiaEspera;

    @Column(nullable = false)
    private boolean emailsResumenDiario;

    @Column(nullable = false)
    private boolean reportesProgresoSemanal;

    @Column(nullable = false)
    private boolean sugerirAutomizacionesAhorro;

    @Column(nullable = false)
    private boolean alertasDispositivosDesconectados;

    @CreatedDate
    @Column(nullable = false, updatable = false)
    private Date createdAt;

    @LastModifiedDate
    @Column(nullable = false)
    private Date updatedAt;

    public DevicePreference() {}

    public DevicePreference(Long userId, Double threshold, boolean notificationEnabled,
                            boolean habilitarMonitoreoEnergia, boolean recibirAlertasAltoConsumo,
                            boolean monitorearCalefaccionRefrigeracion, boolean monitorearElectrodomesticosPrincipales,
                            boolean monitorearElectronicos, boolean monitorearDispositivosCocina,
                            boolean incluirIluminacionExterior, boolean rastrearEnergiaEspera,
                            boolean emailsResumenDiario, boolean reportesProgresoSemanal,
                            boolean sugerirAutomizacionesAhorro, boolean alertasDispositivosDesconectados) {
        this.userId = userId;
        this.threshold = threshold;
        this.notificationEnabled = notificationEnabled;
        this.habilitarMonitoreoEnergia = habilitarMonitoreoEnergia;
        this.recibirAlertasAltoConsumo = recibirAlertasAltoConsumo;
        this.monitorearCalefaccionRefrigeracion = monitorearCalefaccionRefrigeracion;
        this.monitorearElectrodomesticosPrincipales = monitorearElectrodomesticosPrincipales;
        this.monitorearElectronicos = monitorearElectronicos;
        this.monitorearDispositivosCocina = monitorearDispositivosCocina;
        this.incluirIluminacionExterior = incluirIluminacionExterior;
        this.rastrearEnergiaEspera = rastrearEnergiaEspera;
        this.emailsResumenDiario = emailsResumenDiario;
        this.reportesProgresoSemanal = reportesProgresoSemanal;
        this.sugerirAutomizacionesAhorro = sugerirAutomizacionesAhorro;
        this.alertasDispositivosDesconectados = alertasDispositivosDesconectados;
    }

    // Update methods
    public void updateHabilitarMonitoreoEnergia(boolean v) { this.habilitarMonitoreoEnergia = v; }
    public void updateRecibirAlertasAltoConsumo(boolean v) { this.recibirAlertasAltoConsumo = v; }
    public void updateMonitorearCalefaccionRefrigeracion(boolean v) { this.monitorearCalefaccionRefrigeracion = v; }
    public void updateMonitorearElectrodomesticosPrincipales(boolean v) { this.monitorearElectrodomesticosPrincipales = v; }
    public void updateMonitorearElectronicos(boolean v) { this.monitorearElectronicos = v; }
    public void updateMonitorearDispositivosCocina(boolean v) { this.monitorearDispositivosCocina = v; }
    public void updateIncluirIluminacionExterior(boolean v) { this.incluirIluminacionExterior = v; }
    public void updateRastrearEnergiaEspera(boolean v) { this.rastrearEnergiaEspera = v; }
    public void updateEmailsResumenDiario(boolean v) { this.emailsResumenDiario = v; }
    public void updateReportesProgresoSemanal(boolean v) { this.reportesProgresoSemanal = v; }
    public void updateSugerirAutomizacionesAhorro(boolean v) { this.sugerirAutomizacionesAhorro = v; }
    public void updateAlertasDispositivosDesconectados(boolean v) { this.alertasDispositivosDesconectados = v; }
}
