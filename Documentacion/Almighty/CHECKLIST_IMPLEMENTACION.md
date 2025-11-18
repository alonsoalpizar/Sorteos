# Checklist de Implementación - Módulo Almighty

**Versión:** 1.0
**Fecha inicio:** Pendiente
**Duración estimada:** 7-8 semanas

---

## Instrucciones de Uso

- Marcar con `[x]` cuando la tarea esté completada
- Actualizar fecha de completado en columna derecha
- Cada semana, actualizar el porcentaje de progreso
- Reportar bloqueadores en sección de Notas al final

---

## Semana 1: Fundación - Base de Datos

### Migraciones (7 archivos)

- [ ] 012_create_company_settings.up.sql | Fecha: ____
- [ ] 012_create_company_settings.down.sql | Fecha: ____
- [ ] 013_create_payment_processors.up.sql | Fecha: ____
- [ ] 013_create_payment_processors.down.sql | Fecha: ____
- [ ] 014_create_organizer_profiles.up.sql | Fecha: ____
- [ ] 014_create_organizer_profiles.down.sql | Fecha: ____
- [ ] 015_create_settlements.up.sql | Fecha: ____
- [ ] 015_create_settlements.down.sql | Fecha: ____
- [ ] 016_create_system_parameters.up.sql | Fecha: ____
- [ ] 016_create_system_parameters.down.sql | Fecha: ____
- [ ] 017_add_raffle_admin_fields.up.sql | Fecha: ____
- [ ] 017_add_raffle_admin_fields.down.sql | Fecha: ____
- [ ] 018_add_user_admin_fields.up.sql | Fecha: ____
- [ ] 018_add_user_admin_fields.down.sql | Fecha: ____

### Testing de Migraciones

- [ ] Ejecutar migraciones en development | Fecha: ____
- [ ] Validar integridad de datos | Fecha: ____
- [ ] Probar rollback de cada migración | Fecha: ____
- [ ] Insertar datos de prueba | Fecha: ____

### Modelos de Dominio (Go)

- [ ] internal/domain/company_settings.go | Fecha: ____
- [ ] internal/domain/payment_processor.go | Fecha: ____
- [ ] internal/domain/organizer_profile.go | Fecha: ____
- [ ] internal/domain/settlement.go | Fecha: ____
- [ ] internal/domain/system_parameter.go | Fecha: ____

### Repositorios (PostgreSQL)

- [ ] internal/adapters/db/company_settings_repository.go | Fecha: ____
- [ ] internal/adapters/db/payment_processor_repository.go | Fecha: ____
- [ ] internal/adapters/db/organizer_profile_repository.go | Fecha: ____
- [ ] internal/adapters/db/settlement_repository.go | Fecha: ____
- [ ] internal/adapters/db/system_parameter_repository.go | Fecha: ____

### Tests Unitarios de Repositorios

- [ ] Tests de company_settings_repository | Fecha: ____
- [ ] Tests de payment_processor_repository | Fecha: ____
- [ ] Tests de organizer_profile_repository | Fecha: ____
- [ ] Tests de settlement_repository | Fecha: ____
- [ ] Tests de system_parameter_repository | Fecha: ____

**Progreso Semana 1:** ░░░░░░░░░░ 0% (0/38 tareas)

---

## Semana 2: Gestión de Usuarios

### Casos de Uso - Usuarios

- [ ] internal/usecase/admin/user/list_users.go | Fecha: ____
- [ ] internal/usecase/admin/user/get_user_detail.go | Fecha: ____
- [ ] internal/usecase/admin/user/update_user_status.go | Fecha: ____
- [ ] internal/usecase/admin/user/update_user_kyc.go | Fecha: ____
- [ ] internal/usecase/admin/user/reset_user_password.go | Fecha: ____
- [ ] internal/usecase/admin/user/delete_user.go | Fecha: ____

### API Handlers - Usuarios

- [ ] internal/adapters/http/handler/admin/user_handler.go | Fecha: ____
- [ ] Implementar List() endpoint | Fecha: ____
- [ ] Implementar GetByID() endpoint | Fecha: ____
- [ ] Implementar UpdateStatus() endpoint | Fecha: ____
- [ ] Implementar UpdateKYC() endpoint | Fecha: ____
- [ ] Implementar ResetPassword() endpoint | Fecha: ____
- [ ] Implementar Delete() endpoint | Fecha: ____

### Rutas API

- [ ] cmd/api/routes.go - Crear setupAdminRoutes() | Fecha: ____
- [ ] Configurar grupo /api/v1/admin/users | Fecha: ____
- [ ] Aplicar middleware Authenticate() | Fecha: ____
- [ ] Aplicar middleware RequireRole(super_admin) | Fecha: ____
- [ ] Aplicar rate limiting (10 req/min) | Fecha: ____

### Tests de Integración - Usuarios

- [ ] Tests de GET /admin/users | Fecha: ____
- [ ] Tests de GET /admin/users/:id | Fecha: ____
- [ ] Tests de PATCH /admin/users/:id/status | Fecha: ____
- [ ] Tests de PATCH /admin/users/:id/kyc | Fecha: ____
- [ ] Tests de POST /admin/users/:id/reset-password | Fecha: ____
- [ ] Tests de DELETE /admin/users/:id | Fecha: ____

**Progreso Semana 2:** ░░░░░░░░░░ 0% (0/24 tareas)

---

## Semana 3: Gestión de Organizadores

### Casos de Uso - Organizadores

- [ ] internal/usecase/admin/organizer/list_organizers.go | Fecha: ____
- [ ] internal/usecase/admin/organizer/get_organizer_detail.go | Fecha: ____
- [ ] internal/usecase/admin/organizer/update_organizer_profile.go | Fecha: ____
- [ ] internal/usecase/admin/organizer/set_commission_override.go | Fecha: ____
- [ ] internal/usecase/admin/organizer/calculate_organizer_revenue.go | Fecha: ____

### API Handlers - Organizadores

- [ ] internal/adapters/http/handler/admin/organizer_handler.go | Fecha: ____
- [ ] Implementar List() endpoint | Fecha: ____
- [ ] Implementar GetByID() endpoint | Fecha: ____
- [ ] Implementar Update() endpoint | Fecha: ____
- [ ] Implementar SetCommission() endpoint | Fecha: ____
- [ ] Implementar GetRevenue() endpoint | Fecha: ____

### Rutas API

- [ ] Configurar grupo /api/v1/admin/organizers | Fecha: ____

### Tests de Integración - Organizadores

- [ ] Tests de GET /admin/organizers | Fecha: ____
- [ ] Tests de GET /admin/organizers/:id | Fecha: ____
- [ ] Tests de PUT /admin/organizers/:id | Fecha: ____
- [ ] Tests de PATCH /admin/organizers/:id/commission | Fecha: ____
- [ ] Tests de GET /admin/organizers/:id/revenue | Fecha: ____

**Progreso Semana 3:** ░░░░░░░░░░ 0% (0/17 tareas)

---

## Semana 4: Gestión Avanzada de Rifas

### Casos de Uso - Rifas Admin

- [ ] internal/usecase/admin/raffle/list_raffles_admin.go | Fecha: ____
- [ ] internal/usecase/admin/raffle/force_status_change.go | Fecha: ____
- [ ] internal/usecase/admin/raffle/add_admin_notes.go | Fecha: ____
- [ ] internal/usecase/admin/raffle/manual_draw_winner.go | Fecha: ____
- [ ] internal/usecase/admin/raffle/cancel_raffle_with_refund.go | Fecha: ____
- [ ] internal/usecase/admin/raffle/view_raffle_transactions.go | Fecha: ____

### Casos de Uso - Pagos Admin

- [ ] internal/usecase/admin/payment/list_payments_admin.go | Fecha: ____
- [ ] internal/usecase/admin/payment/process_refund.go | Fecha: ____
- [ ] internal/usecase/admin/payment/manage_dispute.go | Fecha: ____
- [ ] internal/usecase/admin/payment/view_payment_detail.go | Fecha: ____

### API Handlers

- [ ] internal/adapters/http/handler/admin/raffle_handler.go | Fecha: ____
- [ ] internal/adapters/http/handler/admin/payment_handler.go | Fecha: ____

### Rutas API

- [ ] Configurar grupo /api/v1/admin/raffles | Fecha: ____
- [ ] Configurar grupo /api/v1/admin/payments | Fecha: ____

### Tests de Integración

- [ ] Tests de raffles admin endpoints (6 tests) | Fecha: ____
- [ ] Tests de payments admin endpoints (4 tests) | Fecha: ____

**Progreso Semana 4:** ░░░░░░░░░░ 0% (0/16 tareas)

---

## Semana 5: Liquidaciones (Settlements)

### Casos de Uso - Settlements

- [ ] internal/usecase/admin/settlement/create_settlement.go | Fecha: ____
- [ ] internal/usecase/admin/settlement/approve_settlement.go | Fecha: ____
- [ ] internal/usecase/admin/settlement/reject_settlement.go | Fecha: ____
- [ ] internal/usecase/admin/settlement/mark_settlement_paid.go | Fecha: ____
- [ ] internal/usecase/admin/settlement/list_settlements.go | Fecha: ____
- [ ] internal/usecase/admin/settlement/auto_create_settlements.go | Fecha: ____

### API Handlers

- [ ] internal/adapters/http/handler/admin/settlement_handler.go | Fecha: ____

### Rutas API

- [ ] Configurar grupo /api/v1/admin/settlements | Fecha: ____

### Scheduled Jobs

- [ ] internal/infrastructure/scheduler/settlement_job.go | Fecha: ____
- [ ] Configurar cron job diario | Fecha: ____

### Tests

- [ ] Tests de settlement use cases (6 tests) | Fecha: ____
- [ ] Tests de settlement endpoints (6 tests) | Fecha: ____

**Progreso Semana 5:** ░░░░░░░░░░ 0% (0/14 tareas)

---

## Semana 6: Reportes Financieros

### Casos de Uso - Reports

- [ ] internal/usecase/admin/reports/global_dashboard.go | Fecha: ____
- [ ] internal/usecase/admin/reports/revenue_report.go | Fecha: ____
- [ ] internal/usecase/admin/reports/raffle_liquidations_report.go | Fecha: ____
- [ ] internal/usecase/admin/reports/organizer_payouts_report.go | Fecha: ____
- [ ] internal/usecase/admin/reports/commission_breakdown.go | Fecha: ____
- [ ] internal/usecase/admin/reports/export_report.go | Fecha: ____

### API Handlers

- [ ] internal/adapters/http/handler/admin/reports_handler.go | Fecha: ____

### Rutas API

- [ ] Configurar grupo /api/v1/admin/reports | Fecha: ____

### Tests

- [ ] Tests de reports use cases (6 tests) | Fecha: ____
- [ ] Tests de reports endpoints (6 tests) | Fecha: ____

**Progreso Semana 6:** ░░░░░░░░░░ 0% (0/12 tareas)

---

## Semana 7: Configuración del Sistema y Categorías

### Casos de Uso - Categorías

- [ ] internal/usecase/admin/category/create_category.go | Fecha: ____
- [ ] internal/usecase/admin/category/update_category.go | Fecha: ____
- [ ] internal/usecase/admin/category/delete_category.go | Fecha: ____
- [ ] internal/usecase/admin/category/reorder_categories.go | Fecha: ____

### Casos de Uso - System Config

- [ ] internal/usecase/admin/system/list_parameters.go | Fecha: ____
- [ ] internal/usecase/admin/system/update_parameter.go | Fecha: ____
- [ ] internal/usecase/admin/system/get_company_settings.go | Fecha: ____
- [ ] internal/usecase/admin/system/update_company_settings.go | Fecha: ____
- [ ] internal/usecase/admin/system/list_payment_processors.go | Fecha: ____
- [ ] internal/usecase/admin/system/update_payment_processor.go | Fecha: ____

### API Handlers

- [ ] internal/adapters/http/handler/admin/category_handler.go | Fecha: ____
- [ ] internal/adapters/http/handler/admin/system_handler.go | Fecha: ____

### Rutas API

- [ ] Configurar grupo /api/v1/admin/categories | Fecha: ____
- [ ] Configurar grupo /api/v1/admin/system | Fecha: ____
- [ ] Configurar grupo /api/v1/admin/audit | Fecha: ____

### Tests

- [ ] Tests de category endpoints (4 tests) | Fecha: ____
- [ ] Tests de system endpoints (6 tests) | Fecha: ____

**Progreso Semana 7:** ░░░░░░░░░░ 0% (0/17 tareas)

---

## Frontend - Estructura Base (Paralelo a Backend)

### Layout y Rutas

- [ ] frontend/src/features/admin/layout/AdminLayout.tsx | Fecha: ____
- [ ] frontend/src/features/admin/layout/AdminSidebar.tsx | Fecha: ____
- [ ] frontend/src/app/routes.tsx - Agregar rutas admin | Fecha: ____
- [ ] frontend/src/components/ProtectedRoute.tsx - Validar super_admin | Fecha: ____

### Páginas Principales

- [ ] frontend/src/features/admin/pages/AdminDashboard.tsx | Fecha: ____
- [ ] frontend/src/features/admin/pages/UsersPage.tsx | Fecha: ____
- [ ] frontend/src/features/admin/pages/UserDetailPage.tsx | Fecha: ____
- [ ] frontend/src/features/admin/pages/OrganizersPage.tsx | Fecha: ____
- [ ] frontend/src/features/admin/pages/OrganizerDetailPage.tsx | Fecha: ____
- [ ] frontend/src/features/admin/pages/RafflesAdminPage.tsx | Fecha: ____
- [ ] frontend/src/features/admin/pages/RaffleDetailAdminPage.tsx | Fecha: ____
- [ ] frontend/src/features/admin/pages/PaymentsPage.tsx | Fecha: ____
- [ ] frontend/src/features/admin/pages/SettlementsPage.tsx | Fecha: ____
- [ ] frontend/src/features/admin/pages/CategoriesPage.tsx | Fecha: ____
- [ ] frontend/src/features/admin/pages/ReportsPage.tsx | Fecha: ____
- [ ] frontend/src/features/admin/pages/SystemConfigPage.tsx | Fecha: ____
- [ ] frontend/src/features/admin/pages/AuditLogsPage.tsx | Fecha: ____

### Componentes Reutilizables

- [ ] frontend/src/features/admin/components/KPICard.tsx | Fecha: ____
- [ ] frontend/src/features/admin/components/DataTable.tsx | Fecha: ____
- [ ] frontend/src/features/admin/components/StatusBadge.tsx | Fecha: ____
- [ ] frontend/src/features/admin/components/RevenueChart.tsx | Fecha: ____
- [ ] frontend/src/features/admin/components/CategoryPieChart.tsx | Fecha: ____
- [ ] frontend/src/features/admin/components/ExportButton.tsx | Fecha: ____

### Hooks Personalizados

- [ ] frontend/src/hooks/useAdminUsers.ts | Fecha: ____
- [ ] frontend/src/hooks/useAdminOrganizers.ts | Fecha: ____
- [ ] frontend/src/hooks/useAdminRaffles.ts | Fecha: ____
- [ ] frontend/src/hooks/useAdminPayments.ts | Fecha: ____
- [ ] frontend/src/hooks/useAdminSettlements.ts | Fecha: ____
- [ ] frontend/src/hooks/useAdminCategories.ts | Fecha: ____
- [ ] frontend/src/hooks/useAdminReports.ts | Fecha: ____
- [ ] frontend/src/hooks/useAdminSystem.ts | Fecha: ____
- [ ] frontend/src/hooks/useAdminAudit.ts | Fecha: ____

**Progreso Frontend:** ░░░░░░░░░░ 0% (0/34 tareas)

---

## Testing (Semana 8)

### Unit Tests Backend

- [ ] Tests de user use cases (>80% coverage) | Fecha: ____
- [ ] Tests de organizer use cases (>80% coverage) | Fecha: ____
- [ ] Tests de raffle admin use cases (>80% coverage) | Fecha: ____
- [ ] Tests de payment admin use cases (>80% coverage) | Fecha: ____
- [ ] Tests de settlement use cases (>80% coverage) | Fecha: ____
- [ ] Tests de reports use cases (>80% coverage) | Fecha: ____
- [ ] Tests de system config use cases (>80% coverage) | Fecha: ____

### Integration Tests Backend

- [ ] Tests de endpoints de usuarios | Fecha: ____
- [ ] Tests de endpoints de organizadores | Fecha: ____
- [ ] Tests de endpoints de rifas admin | Fecha: ____
- [ ] Tests de endpoints de pagos admin | Fecha: ____
- [ ] Tests de endpoints de settlements | Fecha: ____
- [ ] Tests de endpoints de reports | Fecha: ____
- [ ] Tests de endpoints de system config | Fecha: ____
- [ ] Tests de permisos (user normal no puede acceder) | Fecha: ____

### E2E Tests Frontend

- [ ] Test: Login como super_admin y acceder a /admin | Fecha: ____
- [ ] Test: Suspender usuario y verificar audit log | Fecha: ____
- [ ] Test: Cambiar KYC de usuario | Fecha: ____
- [ ] Test: Aprobar settlement y marcar como pagado | Fecha: ____
- [ ] Test: Cancelar rifa con refund | Fecha: ____
- [ ] Test: Crear categoría y reordenar | Fecha: ____
- [ ] Test: Editar system parameter | Fecha: ____
- [ ] Test: Exportar reporte a CSV | Fecha: ____

### Security Tests

- [ ] Penetration testing de permisos | Fecha: ____
- [ ] Test de rate limiting | Fecha: ____
- [ ] Test de validación de inputs (SQL injection, XSS) | Fecha: ____
- [ ] Audit de dependencias (npm audit, go mod) | Fecha: ____

**Progreso Testing:** ░░░░░░░░░░ 0% (0/27 tareas)

---

## Despliegue a Producción (Semana 8)

### Preparación

- [ ] Backup completo de base de datos de producción | Fecha: ____
- [ ] Ejecutar migraciones en staging | Fecha: ____
- [ ] Validar migraciones en staging | Fecha: ____
- [ ] Plan de rollback documentado | Fecha: ____

### Despliegue Backend

- [ ] Build de binario Go | Fecha: ____
- [ ] Deploy a servidor de producción | Fecha: ____
- [ ] Ejecutar migraciones en producción (ventana de mantenimiento) | Fecha: ____
- [ ] Restart de servicio sorteos-backend | Fecha: ____
- [ ] Verificar health checks | Fecha: ____

### Despliegue Frontend

- [ ] Build de producción (npm run build) | Fecha: ____
- [ ] Deploy a /var/www/sorteos.club | Fecha: ____
- [ ] Clear cache de Nginx | Fecha: ____
- [ ] Verificar que /admin carga correctamente | Fecha: ____

### Smoke Testing en Producción

- [ ] Login como super_admin | Fecha: ____
- [ ] Verificar dashboard carga | Fecha: ____
- [ ] Probar suspender usuario de prueba | Fecha: ____
- [ ] Probar aprobar settlement de prueba | Fecha: ____
- [ ] Verificar audit logs se crean | Fecha: ____

### Documentación

- [ ] Actualizar API_ENDPOINTS.md con Swagger | Fecha: ____
- [ ] Completar guía de usuario para super_admin | Fecha: ____
- [ ] Escribir runbooks operacionales | Fecha: ____
- [ ] Video tutorial (opcional) | Fecha: ____

**Progreso Despliegue:** ░░░░░░░░░░ 0% (0/18 tareas)

---

## Resumen de Progreso Global

| Categoría | Total | Completadas | Progreso |
|-----------|-------|-------------|----------|
| **Semana 1 - Fundación** | 38 | 0 | ░░░░░░░░░░ 0% |
| **Semana 2 - Usuarios** | 24 | 0 | ░░░░░░░░░░ 0% |
| **Semana 3 - Organizadores** | 17 | 0 | ░░░░░░░░░░ 0% |
| **Semana 4 - Rifas/Pagos** | 16 | 0 | ░░░░░░░░░░ 0% |
| **Semana 5 - Settlements** | 14 | 0 | ░░░░░░░░░░ 0% |
| **Semana 6 - Reportes** | 12 | 0 | ░░░░░░░░░░ 0% |
| **Semana 7 - Config/Cats** | 17 | 0 | ░░░░░░░░░░ 0% |
| **Frontend** | 34 | 0 | ░░░░░░░░░░ 0% |
| **Testing** | 27 | 0 | ░░░░░░░░░░ 0% |
| **Despliegue** | 18 | 0 | ░░░░░░░░░░ 0% |
| **TOTAL** | **217** | **0** | **░░░░░░░░░░ 0%** |

---

## Notas y Bloqueadores

### Semana 1
- Bloqueadores: _____
- Decisiones tomadas: _____
- Cambios al plan: _____

### Semana 2
- Bloqueadores: _____
- Decisiones tomadas: _____
- Cambios al plan: _____

### Semana 3
- Bloqueadores: _____
- Decisiones tomadas: _____
- Cambios al plan: _____

### Semana 4
- Bloqueadores: _____
- Decisiones tomadas: _____
- Cambios al plan: _____

### Semana 5
- Bloqueadores: _____
- Decisiones tomadas: _____
- Cambios al plan: _____

### Semana 6
- Bloqueadores: _____
- Decisiones tomadas: _____
- Cambios al plan: _____

### Semana 7
- Bloqueadores: _____
- Decisiones tomadas: _____
- Cambios al plan: _____

### Semana 8
- Bloqueadores: _____
- Decisiones tomadas: _____
- Cambios al plan: _____

---

## Hitos Principales

- [ ] **Hito 1:** Base de datos completa (Semana 1) | Fecha objetivo: ____ | Completado: ____
- [ ] **Hito 2:** Backend gestión de usuarios funcional (Semana 2) | Fecha objetivo: ____ | Completado: ____
- [ ] **Hito 3:** Backend gestión de organizadores funcional (Semana 3) | Fecha objetivo: ____ | Completado: ____
- [ ] **Hito 4:** Backend gestión de rifas y pagos funcional (Semana 4) | Fecha objetivo: ____ | Completado: ____
- [ ] **Hito 5:** Sistema de settlements completo (Semana 5) | Fecha objetivo: ____ | Completado: ____
- [ ] **Hito 6:** Dashboard y reportes funcionando (Semana 6) | Fecha objetivo: ____ | Completado: ____
- [ ] **Hito 7:** Frontend completo (Semana 7) | Fecha objetivo: ____ | Completado: ____
- [ ] **Hito 8:** Testing y despliegue a producción (Semana 8) | Fecha objetivo: ____ | Completado: ____

---

**Última actualización:** ____
**Actualizado por:** ____
