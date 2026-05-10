package com.sems.devicemanagement.infrastructure.persistence.jpa.repositories;

import com.sems.devicemanagement.domain.model.entities.DevicePreference;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;

import java.util.Optional;

@Repository
public interface PreferencesRepository extends JpaRepository<DevicePreference, Long> {
    Optional<DevicePreference> findByUserId(Long userId);
    boolean existsByUserId(Long userId);
}
