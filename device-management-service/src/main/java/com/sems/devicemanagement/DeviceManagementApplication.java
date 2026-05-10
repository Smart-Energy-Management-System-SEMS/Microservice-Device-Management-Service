package com.sems.devicemanagement;

import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.data.jpa.repository.config.EnableJpaAuditing;

@SpringBootApplication
@EnableJpaAuditing
public class DeviceManagementApplication {

    public static void main(String[] args) {
        SpringApplication.run(DeviceManagementApplication.class, args);
        System.out.println("===========================================");
        System.out.println("  SEMS - Device Management Service");
        System.out.println("  Swagger UI: /swagger-ui.html");
        System.out.println("  API Docs:   /v3/api-docs");
        System.out.println("===========================================");
    }
}
