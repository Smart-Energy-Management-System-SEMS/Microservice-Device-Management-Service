package com.sems.devicemanagement.infrastructure.persistence.jpa.repositories;

import com.sems.devicemanagement.domain.model.aggregates.Device;
import com.sems.devicemanagement.domain.model.valueobjects.UserId;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;

import java.util.List;

@Repository
public interface DeviceRepository extends JpaRepository<Device, Long> {
    List<Device> findByUserId(UserId userId);
}
