package com.sems.devicemanagement.config;

import io.swagger.v3.oas.models.Components;
import io.swagger.v3.oas.models.OpenAPI;
import io.swagger.v3.oas.models.info.Contact;
import io.swagger.v3.oas.models.info.Info;
import io.swagger.v3.oas.models.security.SecurityRequirement;
import io.swagger.v3.oas.models.security.SecurityScheme;
import io.swagger.v3.oas.models.servers.Server;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;

import java.util.List;

@Configuration
public class OpenApiConfig {

    @Bean
    public OpenAPI deviceManagementOpenAPI() {
        final String securitySchemeName = "bearerAuth";

        return new OpenAPI()
                .info(new Info()
                        .title("SEMS - Device Management Service")
                        .description("""
                                ## Microservicio de Gestión de Dispositivos IoT
                                
                                **Bounded Context:** Device Management
                                
                                Este microservicio es responsable de:
                                - Vincular y desvincular dispositivos IoT al hogar del usuario
                                - Consultar el estado de los dispositivos registrados
                                - Gestionar las preferencias de monitoreo energético por usuario
                                
                                ### Autenticación
                                Utiliza el mismo token JWT generado por el **IAM Service**.  
                                Incluir en el header: `Authorization: Bearer <token>`
                                """)
                        .version("v1.0.0")
                        .contact(new Contact()
                                .name("Energix - SEMS Team")
                                .email("sems@energix.pe")))
                .servers(List.of(
                        new Server().url("/").description("Current server")
                ))
                .addSecurityItem(new SecurityRequirement().addList(securitySchemeName))
                .components(new Components()
                        .addSecuritySchemes(securitySchemeName, new SecurityScheme()
                                .name(securitySchemeName)
                                .type(SecurityScheme.Type.HTTP)
                                .scheme("bearer")
                                .bearerFormat("JWT")
                                .description("Ingresa el JWT token obtenido del IAM Service")));
    }
}
