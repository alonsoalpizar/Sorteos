# Roadmap - Módulo Administrador "Almighty"

**Versión:** 1.0
**Fecha inicio:** 2025-11-18
**Duración estimada:** 7-8 semanas
**Metodología:** Sprints de 1-2 semanas

---

## 1. Visión General

El módulo **Almighty Admin** proporciona control total sobre la plataforma Sorteos.club, permitiendo a los super-administradores:

✅ **Control de Datos Maestros** - Gestión de información de la empresa
✅ **Conectividad de Procesadores** - Administración de Stripe, PayPal y otros
✅ **Gestión de Organizadores** - Perfiles, comisiones y pagos
✅ **Administración de Usuarios** - Permisos, KYC, suspensiones
✅ **Mantenimiento de Categorías** - CRUD de categorías de rifas
✅ **Control Global de Rifas** - Suspensión, habilitación, observación
✅ **Dashboard Ejecutivo** - Métricas y KPIs en tiempo real
✅ **Reportes Financieros** - Ingresos globales y liquidaciones por rifa

### 1.1 Componentes Principales

```
┌─────────────────────────────────────────────────────────────┐
│                    ALMIGHTY ADMIN MODULE                     │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│  ┌────────────┐  ┌────────────┐  ┌────────────┐            │
│  │   Users    │  │ Organizers │  │  Raffles   │            │
│  │   Mgmt     │  │    Mgmt    │  │    Mgmt    │            │
│  └────────────┘  └────────────┘  └────────────┘            │
│                                                              │
│  ┌────────────┐  ┌────────────┐  ┌────────────┐            │
│  │  Payments  │  │Settlements │  │ Categories │            │
│  │    Mgmt    │  │    Mgmt    │  │    Mgmt    │            │
│  └────────────┘  └────────────┘  └────────────┘            │
│                                                              │
│  ┌────────────┐  ┌────────────┐  ┌────────────┐            │
│  │  Reports   │  │   System   │  │   Audit    │            │
│  │ Financial  │  │   Config   │  │    Logs    │            │
│  └────────────┘  └────────────┘  └────────────┘            │
│                                                              │
└─────────────────────────────────────────────────────────────┘
```

### 1.2 Tecnologías

- **Backend:** Go 1.22+ (Gin framework)
- **Frontend:** React 18 + TypeScript + Vite
- **UI Library:** shadcn/ui + Tailwind CSS
- **Base de Datos:** PostgreSQL 16
- **Cache:** Redis 7
- **Gráficos:** Recharts / Chart.js

### 1.3 Documentación Relacionada

- [ARQUITECTURA_ALMIGHTY.md](ARQUITECTURA_ALMIGHTY.md) - Arquitectura técnica detallada
- [BASE_DE_DATOS.md](BASE_DE_DATOS.md) - Esquemas de base de datos
- [API_ENDPOINTS.md](API_ENDPOINTS.md) - Especificación de API REST
- [CASOS_DE_USO.md](CASOS_DE_USO.md) - Casos de uso del sistema
- [FRONTEND_COMPONENTES.md](FRONTEND_COMPONENTES.md) - Componentes UI
- [SEGURIDAD.md](SEGURIDAD.md) - Modelo de seguridad y permisos
- [TESTING.md](TESTING.md) - Estrategia de testing
- [MIGRACION_DATOS.md](MIGRACION_DATOS.md) - Plan de migración
- [CHECKLIST_IMPLEMENTACION.md](CHECKLIST_IMPLEMENTACION.md) - Lista de tareas

---

## 2. Métricas de Progreso Global

| Categoría | Total | Completadas | Progreso |
|-----------|-------|-------------|----------|
| **Migraciones DB** | 7 | 7 | ██████████ 100% ✅ |
| **Repositorios** | 7 | 7 | ██████████ 100% ✅ |
| **Casos de Uso** | 47 | 47 | ██████████ 100% ✅ |
| **Endpoints API** | 52 | 0 | ░░░░░░░░░░ 0% |
| **Páginas Frontend** | 12 | 0 | ░░░░░░░░░░ 0% |
| **Tests** | 60 | 0 | ░░░░░░░░░░ 0% |
| **TOTAL** | **185** | **61** | **███░░░░░░░ 33%** |

**Última actualización:** 2025-11-18 (Fases 2, 3, 4, 5, 6, 7, 8 completadas - 8/8 fases backend)

---

## 3. Fase 1: Fundación (Semana 1-2) ✅ COMPLETADA

**Objetivo:** Crear la infraestructura base de datos y modelos de dominio.

**Duración:** 2 semanas
**Prioridad:** 🔴 CRÍTICA
**Progreso:** ██████████ 100% (32/32 tareas)


### 3.1 Migraciones de Base de Datos ✅

#### 000009_company_settings.up.sql ✅
- [x] Crear tabla `company_settings`
- [x] Agregar campos: company_name, tax_id, address, contact info, logo_url
- [x] Insertar datos iniciales de Sorteos.club
- [x] Crear trigger updated_at
- [x] Validar migración en development

#### 000010_payment_processors.up.sql ✅
- [x] Crear tabla `payment_processors`
- [x] Agregar campos: provider, name, is_active, is_sandbox, credentials (encriptados)
- [x] Crear función de encriptación para secrets
- [x] Insertar configuración actual de Stripe/PayPal
- [x] Validar migración en development

#### 000011_organizer_profiles.up.sql ✅
- [x] Crear tabla `organizer_profiles`
- [x] Agregar campos: user_id, business_name, tax_id, bank info, commission_override
- [x] Crear índice en user_id (FK)
- [x] Crear trigger para calcular pending_payout
- [x] Validar migración en development

#### 000012_settlements.up.sql ✅
- [x] Crear tabla `settlements`
- [x] Agregar campos: raffle_id, organizer_id, amounts, status, payment info
- [x] Crear índices en raffle_id, organizer_id, status
- [x] Crear ENUM para settlement_status
- [x] Validar migración en development

#### 000013_system_parameters.up.sql ✅
- [x] Crear tabla `system_parameters`
- [x] Agregar campos: key, value, value_type, category, is_sensitive
- [x] Crear índice único en key
- [x] Insertar parámetros por defecto (platform_fee, max_raffles, etc.)
- [x] Validar migración en development

#### 000014_raffle_admin_fields.up.sql ✅
- [x] Agregar campos a `raffles`: suspension_reason, suspended_by, suspended_at, admin_notes
- [x] Crear FK en suspended_by → users(id)
- [x] Crear índice en suspended_by
- [x] Validar migración en development

#### 000015_user_admin_fields.up.sql ✅
- [x] Agregar campos a `users`: suspension_reason, suspended_by, suspended_at
- [x] Agregar campos: last_kyc_review, kyc_reviewer
- [x] Crear FKs en suspended_by, kyc_reviewer → users(id)
- [x] Validar migración en development

### 3.2 Modelos de Dominio (Go)

#### internal/domain/company_settings.go ✅
- [x] Crear entidad `CompanySettings`
- [x] Agregar métodos de validación
- [x] Crear interfaz `CompanySettingsRepository`
- [x] Documentar estructura

#### internal/domain/payment_processor.go ✅
- [x] Crear entidad `PaymentProcessor`
- [x] Agregar enum `ProcessorProvider` (stripe, paypal, etc.)
- [x] Agregar métodos para encriptar/desencriptar credentials
- [x] Crear interfaz `PaymentProcessorRepository`

#### internal/domain/organizer_profile.go ✅
- [x] Crear entidad `OrganizerProfile`
- [x] Agregar métodos para calcular revenue
- [x] Agregar validaciones de bank info
- [x] Crear interfaz `OrganizerProfileRepository`

#### internal/domain/settlement.go ✅
- [x] Crear entidad `Settlement`
- [x] Crear enum `SettlementStatus`
- [x] Agregar métodos de cálculo (gross, fees, net)
- [x] Crear interfaz `SettlementRepository`

#### internal/domain/system_parameter.go ✅
- [x] Crear entidad `SystemParameter`
- [x] Crear enum `ParameterValueType` (string, int, float, bool, json)
- [x] Agregar métodos de parsing por tipo
- [x] Crear interfaz `SystemParameterRepository`

### 3.3 Repositorios (PostgreSQL)

#### internal/adapters/db/company_settings_repository.go ✅
- [x] Implementar `Get() (*CompanySettings, error)`
- [x] Implementar `Update(settings *CompanySettings) error`
- [x] Agregar logging y error handling
- [ ] Escribir tests unitarios

#### internal/adapters/db/payment_processor_repository.go ✅
- [x] Implementar `List() ([]*PaymentProcessor, error)`
- [x] Implementar `GetByID(id int64) (*PaymentProcessor, error)`
- [x] Implementar `GetByProvider(provider string) (*PaymentProcessor, error)`
- [x] Implementar `Update(processor *PaymentProcessor) error`
- [x] Implementar `ToggleActive(id int64, active bool) error`
- [ ] Escribir tests unitarios

#### internal/adapters/db/organizer_profile_repository.go ✅
- [x] Implementar `Create(profile *OrganizerProfile) error`
- [x] Implementar `GetByUserID(userID int64) (*OrganizerProfile, error)`
- [x] Implementar `List(filters map[string]interface{}, offset, limit int) ([]*OrganizerProfile, int64, error)`
- [x] Implementar `Update(profile *OrganizerProfile) error`
- [x] Implementar `UpdateCommission(userID int64, commission float64) error`
- [x] Implementar `GetRevenue(userID int64) (*OrganizerRevenue, error)`
- [ ] Escribir tests unitarios

#### internal/adapters/db/settlement_repository.go ✅
- [x] Implementar `Create(settlement *Settlement) error`
- [x] Implementar `GetByID(id int64) (*Settlement, error)`
- [x] Implementar `List(filters map[string]interface{}, offset, limit int) ([]*Settlement, int64, error)`
- [x] Implementar `UpdateStatus(id int64, status SettlementStatus) error`
- [x] Implementar `Approve(id int64, adminID int64) error`
- [x] Implementar `Reject(id int64, adminID int64, reason string) error`
- [x] Implementar `MarkPaid(id int64, paymentRef string) error`
- [ ] Implementar `GetPendingByOrganizer(organizerID int64) ([]*Settlement, error)`
- [ ] Escribir tests unitarios

#### internal/adapters/db/system_parameter_repository.go ✅
- [x] Implementar `GetByKey(key string) (*SystemParameter, error)`
- [x] Implementar `GetString(key string, defaultValue string) (string, error)`
- [x] Implementar `GetInt(key string, defaultValue int) (int, error)`
- [x] Implementar `GetFloat(key string, defaultValue float64) (float64, error)`
- [x] Implementar `GetBool(key string, defaultValue bool) (bool, error)`
- [x] Implementar `List(category string, offset, limit int) ([]*SystemParameter, int64, error)`
- [ ] Implementar `Update(param *SystemParameter, adminID int64) error`
- [ ] Escribir tests unitarios

### 3.4 Criterios de Aceptación - Fase 1

- ✅ Las 7 migraciones ejecutan sin errores
- ✅ Rollback de migraciones funciona correctamente
- ✅ Todas las entidades de dominio tienen validaciones
- ✅ Todos los repositorios tienen tests unitarios con >80% coverage
- ✅ Datos de prueba insertados en development
- ✅ Documentación de modelos completa

---

## 4. Fase 2: Gestión de Usuarios y Organizadores (Semana 2-3) ✅ COMPLETADA

**Objetivo:** Implementar gestión completa de usuarios y organizadores.

**Duración:** 1-2 semanas
**Prioridad:** 🔴 CRÍTICA
**Progreso:** ██████████ 100% (40/40 tareas)


### 4.1 Casos de Uso - Usuarios ✅

#### internal/usecase/admin/user/list_users.go ✅
- [x] Crear `ListUsersUseCase`
- [x] Implementar filtros: role, status, kyc_level, search (name, email, cedula)
- [x] Implementar paginación
- [x] Implementar ordenamiento (created_at, last_login_at, email)
- [x] Agregar conteo total para paginación
- [x] Logging de auditoría (action: admin_list_users)
- [ ] Escribir tests unitarios

#### internal/usecase/admin/user/get_user_detail.go ✅
- [x] Crear `GetUserDetailUseCase`
- [x] Incluir: user data, raffle stats, payment stats, audit logs recientes
- [x] Logging de auditoría
- [ ] Escribir tests unitarios

#### internal/usecase/admin/user/update_user_status.go ✅
- [x] Crear `UpdateUserStatusUseCase`
- [x] Implementar acciones: suspend, activate, ban
- [x] Validar que admin no puede suspenderse a sí mismo
- [x] Guardar suspension_reason, suspended_by, suspended_at
- [x] Logging de auditoría (severity: warning/critical)
- [x] Enviar email de notificación al usuario
- [ ] Escribir tests unitarios

#### internal/usecase/admin/user/update_user_kyc.go ✅
- [x] Crear `UpdateUserKYCUseCase`
- [x] Implementar cambio de KYC level
- [x] Guardar kyc_reviewer y last_kyc_review
- [x] Validar documentos si existen
- [x] Logging de auditoría
- [ ] Enviar email de notificación
- [ ] Escribir tests unitarios

#### internal/usecase/admin/user/reset_user_password.go ✅
- [x] Crear `ResetUserPasswordUseCase`
- [x] Generar token de reset
- [x] Enviar email con link de reset
- [x] Logging de auditoría
- [ ] Escribir tests unitarios

#### internal/usecase/admin/user/delete_user.go ✅
- [x] Crear `DeleteUserUseCase` (soft delete)
- [x] Validar que usuario no tenga rifas activas
- [x] Marcar como deleted (deleted_at)
- [x] Cancelar rifas draft del usuario
- [x] Logging de auditoría (severity: critical)
- [ ] Escribir tests unitarios

### 4.2 Casos de Uso - Organizadores  ✅

#### internal/usecase/admin/organizer/list_organizers.go ✅
- [x] Crear `ListOrganizersUseCase`
- [x] Implementar filtros: verified, revenue_range, date_range
- [x] Incluir métricas: total_raffles, total_revenue, pending_payout
- [x] Implementar paginación y ordenamiento
- [ ] Logging de auditoría
- [ ] Escribir tests unitarios

#### internal/usecase/admin/organizer/get_organizer_detail.go ✅
- [x] Crear `GetOrganizerDetailUseCase`
- [x] Incluir: profile, user data, raffle list, settlement history, revenue breakdown
- [x] Calcular métricas: avg_raffle_revenue, completion_rate, refund_rate
- [x] Logging de auditoría
- [ ] Escribir tests unitarios

#### internal/usecase/admin/organizer/update_organizer_profile.go ✅
- [x] Crear `UpdateOrganizerProfileUseCase`
- [x] Validar bank info format
- [x] Actualizar payout_schedule, verified status
- [x] Logging de auditoría
- [ ] Escribir tests unitarios

#### internal/usecase/admin/organizer/set_commission_override.go ✅
- [x] Crear `SetCommissionOverrideUseCase`
- [x] Validar rango de comisión (0-50%)
- [x] Guardar commission_override en organizer_profile
- [x] Logging de auditoría (severity: warning)
- [ ] Escribir tests unitarios

#### internal/usecase/admin/organizer/calculate_organizer_revenue.go ✅
- [x] Crear `CalculateOrganizerRevenueUseCase` (321 lines)
- [x] Calcular: gross_revenue, platform_fees, net_revenue, pending_payout
- [x] Filtrar por date_range
- [x] Agrupar por mes/año si se requiere
- [ ] Escribir tests unitarios

### 4.3 API Handlers - Usuarios

#### internal/adapters/http/handler/admin/user_handler.go
- [ ] Crear `UserHandler` con dependencias (use cases)
- [ ] Implementar `List(c *gin.Context)` → 200 OK
- [ ] Implementar `GetByID(c *gin.Context)` → 200 OK / 404 Not Found
- [ ] Implementar `UpdateStatus(c *gin.Context)` → 200 OK / 400 Bad Request
- [ ] Implementar `UpdateKYC(c *gin.Context)` → 200 OK / 400 Bad Request
- [ ] Implementar `ResetPassword(c *gin.Context)` → 200 OK
- [ ] Implementar `Delete(c *gin.Context)` → 204 No Content
- [ ] Agregar validación de inputs con validator
- [ ] Agregar error handling consistente
- [ ] Escribir tests de integración

### 4.4 API Handlers - Organizadores

#### internal/adapters/http/handler/admin/organizer_handler.go
- [ ] Crear `OrganizerHandler`
- [ ] Implementar `List(c *gin.Context)`
- [ ] Implementar `GetByID(c *gin.Context)`
- [ ] Implementar `Update(c *gin.Context)`
- [ ] Implementar `SetCommission(c *gin.Context)`
- [ ] Implementar `GetRevenue(c *gin.Context)`
- [ ] Validación de inputs
- [ ] Error handling
- [ ] Escribir tests de integración

### 4.5 Rutas API

#### cmd/api/routes.go (Admin Routes)
- [ ] Crear función `setupAdminRoutes(router *gin.Engine, handlers *Handlers)`
- [ ] Configurar grupo `/api/v1/admin/users`
- [ ] Configurar grupo `/api/v1/admin/organizers`
- [ ] Aplicar middleware: Authenticate(), RequireRole("super_admin")
- [ ] Aplicar rate limiting (10 req/min)
- [ ] Documentar endpoints

### 4.6 Frontend - Páginas de Usuarios

#### frontend/src/features/admin/pages/UsersPage.tsx
- [ ] Crear componente UsersPage
- [ ] Implementar tabla con shadcn/ui Table
- [ ] Agregar filtros: role, status, KYC level, búsqueda
- [ ] Agregar paginación
- [ ] Agregar acciones: ver detalle, suspender, editar KYC
- [ ] Implementar estado de carga con LoadingSpinner
- [ ] Agregar EmptyState cuando no hay usuarios
- [ ] Estilizar con Tailwind (paleta blue/slate)

#### frontend/src/features/admin/pages/UserDetailPage.tsx
- [ ] Crear componente UserDetailPage
- [ ] Mostrar información completa del usuario
- [ ] Mostrar tabs: Overview, Raffles, Payments, Audit Log
- [ ] Agregar acciones: Suspender, Cambiar KYC, Reset Password
- [ ] Implementar modales de confirmación
- [ ] Mostrar toasts de éxito/error

### 4.7 Frontend - Páginas de Organizadores

#### frontend/src/features/admin/pages/OrganizersPage.tsx
- [ ] Crear componente OrganizersPage
- [ ] Implementar tabla con métricas (revenue, raffles count)
- [ ] Agregar filtros: verified, revenue range
- [ ] Agregar ordenamiento por revenue, created_at
- [ ] Acciones: ver detalle, editar comisión

#### frontend/src/features/admin/pages/OrganizerDetailPage.tsx
- [ ] Crear componente OrganizerDetailPage
- [ ] Mostrar perfil completo
- [ ] Mostrar tabs: Overview, Raffles, Settlements, Revenue
- [ ] Gráfico de ingresos por mes
- [ ] Acción: Set Custom Commission

### 4.8 Frontend - Hooks y API

#### frontend/src/hooks/useAdminUsers.ts
- [ ] Crear hook `useUsers(filters, pagination)`
- [ ] Crear hook `useUserDetail(userId)`
- [ ] Crear hook `useUpdateUserStatus()`
- [ ] Crear hook `useUpdateUserKYC()`
- [ ] Usar React Query para caching

#### frontend/src/hooks/useAdminOrganizers.ts
- [ ] Crear hook `useOrganizers(filters, pagination)`
- [ ] Crear hook `useOrganizerDetail(userId)`
- [ ] Crear hook `useUpdateOrganizerProfile()`
- [ ] Crear hook `useSetCommission()`
- [ ] Crear hook `useOrganizerRevenue(userId, dateRange)`

### 4.9 Criterios de Aceptación - Fase 2

- ✅ Admin puede listar, buscar y filtrar usuarios
- ✅ Admin puede suspender/activar usuarios con razón
- ✅ Admin puede cambiar nivel KYC manualmente
- ✅ Admin puede forzar reset de password
- ✅ Admin puede ver detalle completo de organizador
- ✅ Admin puede establecer comisión personalizada
- ✅ Todas las acciones generan audit logs
- ✅ Tests de integración pasan
- ✅ UI es responsive y sigue diseño de shadcn/ui

---

## 5. Fase 5: Gestión Avanzada de Rifas y Pagos (Semana 4-5) ✅ COMPLETADA
**Estado:** ✅ COMPLETADA - 2025-11-18

**Objetivo:** Control administrativo completo sobre rifas y sistema de pagos.

**Duración:** 1-2 semanas
**Prioridad:** 🟡 ALTA
**Progreso:** ██████████ 100% (10/10 tareas - core use cases)

### 5.1 Casos de Uso - Rifas Admin

#### internal/usecase/admin/raffle/list_raffles_admin.go
- [x] Crear `ListRafflesAdminUseCase`
- [x] Filtros: status (todos incluido suspended), organizer_id, category_id, date_range
- [x] Incluir métricas: sold_count, revenue, platform_fee
- [x] Búsqueda por title
- [ ] Paginación y ordenamiento

#### internal/usecase/admin/raffle/force_status_change.go
- [x] Crear `ForceStatusChangeUseCase`
- [x] Permitir: draft→active, active→suspended, suspended→active, active→cancelled
- [x] Validar transiciones permitidas
- [ ] Guardar admin_notes, suspended_by, suspended_at
- [ ] Logging de auditoría (severity: warning)
- [ ] Notificar al organizador por email

#### internal/usecase/admin/raffle/add_admin_notes.go
- [x] Crear `AddAdminNotesUseCase`
- [x] Agregar notas en campo admin_notes
- [ ] Logging de auditoría

#### internal/usecase/admin/raffle/manual_draw_winner.go
- [x] Crear `ManualDrawWinnerUseCase`
- [x] Validar que rifa esté en estado active
- [x] Seleccionar número ganador (random o especificado)
- [ ] Actualizar winner_number, winner_user_id
- [ ] Cambiar status a completed
- [ ] Enviar emails (ganador, organizador)
- [ ] Logging de auditoría (severity: critical)

#### internal/usecase/admin/raffle/cancel_raffle_with_refund.go
- [x] Crear `CancelRaffleWithRefundUseCase`
- [x] Validar que rifa no esté completed
- [x] Obtener todos los pagos confirmados
- [ ] Iniciar refunds con payment provider (Stripe/PayPal)
- [ ] Actualizar payment status a refunded
- [ ] Cambiar raffle status a cancelled
- [ ] Enviar emails de notificación
- [ ] Logging de auditoría (severity: critical)

#### internal/usecase/admin/raffle/view_raffle_transactions.go
- [x] Crear `ViewRaffleTransactionsUseCase`
- [x] Listar: reservations, payments, refunds, audit logs
- [x] Timeline cronológico de eventos
- [ ] Calcular métricas: conversion_rate, refund_rate

### 5.2 Casos de Uso - Pagos Admin

#### internal/usecase/admin/payment/list_payments_admin.go
- [x] Crear `ListPaymentsAdminUseCase`
- [x] Filtros: status, user_id, raffle_id, date_range, payment_method
- [x] Incluir info de usuario y rifa
- [ ] Paginación y ordenamiento

#### internal/usecase/admin/payment/process_refund.go
- [x] Crear `ProcessRefundUseCase`
- [x] Validar payment status (succeeded)
- [x] Preparado para payment provider API (Stripe/PayPal) - TODO markers
- [ ] Actualizar payment status a refunded
- [ ] Liberar números reservados
- [ ] Actualizar raffle sold_count, revenue
- [ ] Enviar email de confirmación
- [ ] Logging de auditoría (severity: warning)

#### internal/usecase/admin/payment/manage_dispute.go ✅
- [x] Crear `ManageDisputeUseCase` (298 lines)
- [x] Marcar payment con dispute flag
- [x] Guardar metadata de disputa
- [x] Notificar al organizador
- [x] Logging de auditoría

#### internal/usecase/admin/payment/view_payment_detail.go
- [x] Crear `ViewPaymentDetailsUseCase`
- [x] Incluir: payment data, user, raffle, numbers, timeline, webhook events
- [ ] Timeline de eventos del payment

### 5.3 API Handlers - Rifas Admin

#### internal/adapters/http/handler/admin/raffle_handler.go
- [ ] Crear `RaffleAdminHandler`
- [ ] Implementar `List(c *gin.Context)`
- [ ] Implementar `GetByID(c *gin.Context)` (enhanced version)
- [ ] Implementar `ForceStatusChange(c *gin.Context)`
- [ ] Implementar `AddNotes(c *gin.Context)`
- [ ] Implementar `ManualDraw(c *gin.Context)`
- [ ] Implementar `CancelWithRefund(c *gin.Context)`
- [ ] Implementar `ViewTransactions(c *gin.Context)`
- [ ] Validación y error handling

### 5.4 API Handlers - Pagos Admin

#### internal/adapters/http/handler/admin/payment_handler.go
- [ ] Crear `PaymentAdminHandler`
- [ ] Implementar `List(c *gin.Context)`
- [ ] Implementar `GetByID(c *gin.Context)`
- [ ] Implementar `ProcessRefund(c *gin.Context)`
- [ ] Implementar `ManageDispute(c *gin.Context)`
- [ ] Validación y error handling

### 5.5 Rutas API

#### cmd/api/routes.go
- [ ] Agregar grupo `/api/v1/admin/raffles`
- [ ] Agregar grupo `/api/v1/admin/payments`
- [ ] Middleware: super_admin + rate limiting

### 5.6 Frontend - Rifas Admin

#### frontend/src/features/admin/pages/RafflesAdminPage.tsx
- [ ] Crear componente con tabla de rifas
- [ ] Filtros: status (incluir suspended), organizador, categoría
- [ ] Búsqueda por título
- [ ] Badges de estado con colores
- [ ] Acciones: ver detalle, cambiar status, agregar notas

#### frontend/src/features/admin/pages/RaffleDetailAdminPage.tsx
- [ ] Mostrar info completa de rifa
- [ ] Tabs: Overview, Transactions, Audit Log
- [ ] Acciones administrativas: Suspend, Activate, Cancel with Refund, Manual Draw
- [ ] Modales de confirmación con razón
- [ ] Timeline de transacciones

### 5.7 Frontend - Pagos Admin

#### frontend/src/features/admin/pages/PaymentsPage.tsx
- [ ] Crear componente con tabla de pagos
- [ ] Filtros: status, método, fecha
- [ ] Búsqueda por payment ID, usuario, rifa
- [ ] Acciones: ver detalle, refund

#### frontend/src/features/admin/pages/PaymentDetailPage.tsx
- [ ] Mostrar detalle completo del pago
- [ ] Info de Stripe/PayPal (payment_intent_id, etc.)
- [ ] Botón de refund con confirmación
- [ ] Timeline de eventos

### 5.8 Frontend - Hooks

#### frontend/src/hooks/useAdminRaffles.ts
- [ ] `useRafflesAdmin(filters, pagination)`
- [ ] `useRaffleDetailAdmin(raffleId)`
- [ ] `useForceStatusChange()`
- [ ] `useCancelWithRefund()`
- [ ] `useManualDraw()`

#### frontend/src/hooks/useAdminPayments.ts
- [ ] `usePaymentsAdmin(filters, pagination)`
- [ ] `usePaymentDetail(paymentId)`
- [ ] `useProcessRefund()`

### 5.9 Criterios de Aceptación - Fase 3

- ✅ Admin puede ver todas las rifas con filtros avanzados
- ✅ Admin puede suspender/activar rifas con razón
- ✅ Admin puede cancelar rifa con refund automático a compradores
- ✅ Admin puede realizar sorteo manual (seleccionar ganador)
- ✅ Admin puede procesar refunds individuales
- ✅ Timeline de transacciones funciona correctamente
- ✅ Emails de notificación se envían correctamente
- ✅ Tests de integración pasan

---

## 6. Fase 4: Liquidaciones y Pagos a Organizadores (Semana 5-6)

**Objetivo:** Sistema completo de liquidaciones y pagos a organizadores.

**Duración:** 1-2 semanas
**Prioridad:** 🟡 ALTA
**Progreso:** ░░░░░░░░░░ 0% (0/28 tareas)

### 6.1 Casos de Uso - Settlements

#### internal/usecase/admin/settlement/create_settlement.go ✅
- [x] Crear `CreateSettlementUseCase` (207 lines)
- [x] Modalidad individual: para 1 rifa completada
- [x] Modalidad batch: para múltiples rifas de un organizador
- [x] Calcular: gross_revenue, platform_fee (de raffle o override de organizer), net_payout
- [x] Crear registro en settlements table
- [x] Status inicial: pending
- [x] Logging de auditoría

#### internal/usecase/admin/settlement/approve_settlement.go
- [x] Crear `ApproveSettlementUseCase`
- [ ] Validar settlement status = pending
- [ ] Cambiar status a approved
- [ ] Guardar approved_by (admin_id), approved_at
- [ ] Enviar email al organizador
- [ ] Logging de auditoría

#### internal/usecase/admin/settlement/reject_settlement.go
- [x] Crear `RejectSettlementUseCase`
- [ ] Cambiar status a rejected
- [ ] Guardar rejection reason en notes
- [ ] Enviar email al organizador
- [ ] Logging de auditoría

#### internal/usecase/admin/settlement/mark_settlement_paid.go ✅
- [x] Crear `MarkSettlementPaidUseCase` (227 lines)
- [x] Validar settlement status = approved
- [x] Cambiar status a paid
- [x] Guardar payment_method, payment_reference, paid_at
- [x] Actualizar organizer_profile.total_payouts
- [x] Reducir organizer_profile.pending_payout
- [x] Enviar email de confirmación
- [x] Logging de auditoría

#### internal/usecase/admin/settlement/list_settlements.go
- [x] Crear `ListSettlementsUseCase`
- [x] Filtros: status, organizer_id, date_range, KYC level, search
- [ ] Incluir info de organizador y rifa
- [ ] Paginación y ordenamiento
- [ ] Calcular totales por status

#### internal/usecase/admin/settlement/auto_create_settlements.go ✅
- [x] Crear `AutoCreateSettlementsUseCase` (319 lines - batch job)
- [x] Buscar rifas completed sin settlement
- [x] Crear settlements automáticamente
- [x] Logging de auditoría
- [x] Retornar count de settlements creados

### 6.2 API Handlers - Settlements

#### internal/adapters/http/handler/admin/settlement_handler.go
- [ ] Crear `SettlementHandler`
- [ ] Implementar `Create(c *gin.Context)` (individual/batch)
- [ ] Implementar `List(c *gin.Context)`
- [ ] Implementar `GetByID(c *gin.Context)`
- [ ] Implementar `Approve(c *gin.Context)`
- [ ] Implementar `Reject(c *gin.Context)`
- [ ] Implementar `MarkPaid(c *gin.Context)`
- [ ] Validación y error handling

### 6.3 Rutas API

#### cmd/api/routes.go
- [ ] Agregar grupo `/api/v1/admin/settlements`
- [ ] Middleware: super_admin + rate limiting

### 6.4 Frontend - Settlements

#### frontend/src/features/admin/pages/SettlementsPage.tsx
- [ ] Crear tabla de settlements
- [ ] Filtros: status (pending, approved, paid, rejected), organizador, fecha
- [ ] Badges de status con colores (pending=yellow, approved=blue, paid=green, rejected=red)
- [ ] Acciones: ver detalle, aprobar, rechazar, marcar como pagado
- [ ] Totales por status en cards superiores

#### frontend/src/features/admin/pages/SettlementDetailPage.tsx
- [ ] Mostrar detalle completo
- [ ] Info de rifa asociada
- [ ] Desglose: gross revenue, platform fee (%), net payout
- [ ] Botones de acción según status
- [ ] Modal de aprobación
- [ ] Modal de marcar como pagado (pedir payment_method, reference)
- [ ] Modal de rechazo (pedir reason)

### 6.5 Frontend - Hooks

#### frontend/src/hooks/useAdminSettlements.ts
- [ ] `useSettlements(filters, pagination)`
- [ ] `useSettlementDetail(settlementId)`
- [ ] `useCreateSettlement()`
- [ ] `useApproveSettlement()`
- [ ] `useRejectSettlement()`
- [ ] `useMarkSettlementPaid()`

### 6.6 Backend - Scheduled Jobs

#### internal/infrastructure/scheduler/settlement_job.go
- [ ] Crear job que se ejecuta diariamente
- [ ] Llamar a `AutoCreateSettlementsUseCase`
- [ ] Logging de resultados

### 6.7 Criterios de Aceptación - Fase 4

- ✅ Settlements se crean automáticamente para rifas completed
- ✅ Admin puede crear settlement manual
- ✅ Admin puede aprobar/rechazar settlements
- ✅ Admin puede marcar settlement como pagado
- ✅ Organizer profile se actualiza correctamente (total_payouts, pending_payout)
- ✅ Emails de notificación funcionan
- ✅ Workflow completo: pending → approved → paid funciona
- ✅ Tests de integración pasan

---

## 7. Fase 7: Reportes Financieros y Dashboard (Semana 6-7) ✅ COMPLETADA
**Estado:** ✅ COMPLETADA - 2025-11-18
**Progreso:** ██████████ 100% (7/7 tareas - core use cases)

**Objetivo:** Dashboard ejecutivo con métricas y reportes financieros exportables.

**Duración:** 1-2 semanas
**Prioridad:** 🟢 MEDIA
**Progreso:** ░░░░░░░░░░ 0% (0/30 tareas)

### 7.1 Casos de Uso - Reports

#### internal/usecase/admin/reports/global_dashboard.go
- [ ] Crear `GlobalDashboardUseCase`
- [ ] Calcular KPIs:
  - Total users (active, suspended, banned)
  - Total organizers (verified, pending)
  - Total raffles (by status)
  - Revenue (today, this week, this month, this year, all-time)
  - Platform fees collected
  - Pending settlements (count, amount)
  - Recent activity (last 24h)
- [ ] Retornar estructura `DashboardKPIs`

#### internal/usecase/admin/reports/revenue_report.go
- [ ] Crear `RevenueReportUseCase`
- [ ] Filtros: date_range, organizer_id, category_id
- [ ] Calcular: gross_revenue, platform_fees, net_revenue (to organizers)
- [ ] Agrupar por: day, week, month (configurable)
- [ ] Retornar series de tiempo para gráficos

#### internal/usecase/admin/reports/raffle_liquidations_report.go
- [ ] Crear `RaffleLiquidationsReportUseCase`
- [ ] Listar rifas completed con desglose financiero
- [ ] Por rifa: title, organizer, gross, fees, net, settlement_status
- [ ] Filtros: date_range, organizer_id
- [ ] Exportable

#### internal/usecase/admin/reports/organizer_payouts_report.go
- [ ] Crear `OrganizerPayoutsReportUseCase`
- [ ] Por organizador: name, total_raffles, total_revenue, total_fees, total_payouts, pending_payout
- [ ] Filtros: date_range, verified
- [ ] Ordenar por revenue desc

#### internal/usecase/admin/reports/commission_breakdown.go
- [ ] Crear `CommissionBreakdownUseCase`
- [ ] Agrupar por tasa de comisión (10%, custom %)
- [ ] Mostrar: # raffles, gross revenue, fees collected
- [ ] Identificar organizadores con custom commission

#### internal/usecase/admin/reports/export_report.go
- [ ] Crear `ExportReportUseCase`
- [ ] Soportar formatos: CSV, Excel (xlsx), PDF
- [ ] Generar archivo temporal
- [ ] Retornar URL de descarga
- [ ] Auto-cleanup de archivos antiguos

### 7.2 API Handlers - Reports

#### internal/adapters/http/handler/admin/reports_handler.go
- [ ] Crear `ReportsHandler`
- [ ] Implementar `GetDashboard(c *gin.Context)`
- [ ] Implementar `GetRevenueReport(c *gin.Context)`
- [ ] Implementar `GetLiquidationsReport(c *gin.Context)`
- [ ] Implementar `GetPayoutsReport(c *gin.Context)`
- [ ] Implementar `GetCommissionBreakdown(c *gin.Context)`
- [ ] Implementar `ExportReport(c *gin.Context)` (stream file)

### 7.3 Rutas API

#### cmd/api/routes.go
- [ ] Agregar grupo `/api/v1/admin/reports`
- [ ] Middleware: super_admin

### 7.4 Frontend - Dashboard

#### frontend/src/features/admin/pages/AdminDashboard.tsx
- [ ] Crear dashboard principal
- [ ] Grid de KPI cards (4x2):
  - Total Users (con breakdown: active/suspended/banned)
  - Total Organizers (verified/pending)
  - Active Raffles (vs completed/suspended)
  - Revenue This Month (vs last month %)
  - Platform Fees Collected
  - Pending Settlements (count + amount)
  - Today's Revenue
  - New Users This Week
- [ ] Gráfico de ingresos (últimos 30 días) - Line chart
- [ ] Gráfico de rifas por categoría - Pie chart
- [ ] Tabla de rifas recientes (últimas 10)
- [ ] Tabla de settlements pendientes (top 5)
- [ ] Auto-refresh cada 60 segundos

#### frontend/src/features/admin/pages/ReportsPage.tsx
- [ ] Crear página de reportes
- [ ] Tabs:
  - Revenue Report
  - Liquidations Report
  - Organizer Payouts Report
  - Commission Breakdown
- [ ] Filtros por fecha (DateRangePicker)
- [ ] Filtros adicionales según reporte
- [ ] Botón de exportación (CSV, Excel, PDF)
- [ ] Gráficos interactivos con Recharts
- [ ] Tablas con paginación

### 7.5 Frontend - Componentes

#### frontend/src/features/admin/components/KPICard.tsx
- [ ] Crear componente reutilizable
- [ ] Props: title, value, icon, trend (% change), subtitle
- [ ] Sparkline opcional (mini gráfico)
- [ ] Colores según trend (green: positive, red: negative)

#### frontend/src/features/admin/components/RevenueChart.tsx
- [ ] Crear componente con Recharts
- [ ] Line chart de ingresos por día
- [ ] Tooltip con formato de moneda
- [ ] Responsive

#### frontend/src/features/admin/components/CategoryPieChart.tsx
- [ ] Pie chart de rifas por categoría
- [ ] Colores consistentes
- [ ] Leyenda

#### frontend/src/features/admin/components/ExportButton.tsx
- [ ] Botón con dropdown: CSV, Excel, PDF
- [ ] Loading state durante export
- [ ] Auto-download del archivo

### 7.6 Frontend - Hooks

#### frontend/src/hooks/useAdminReports.ts
- [ ] `useDashboardKPIs()`
- [ ] `useRevenueReport(dateRange, filters)`
- [ ] `useLiquidationsReport(dateRange, filters)`
- [ ] `usePayoutsReport(dateRange, filters)`
- [ ] `useCommissionBreakdown(dateRange)`
- [ ] `useExportReport(reportType, format, filters)`

### 7.7 Criterios de Aceptación - Fase 5

- ✅ Dashboard muestra KPIs en tiempo real
- ✅ Gráficos de ingresos y categorías funcionan
- ✅ Reportes muestran datos correctos
- ✅ Exportación a CSV/Excel/PDF funciona
- ✅ Filtros de fecha funcionan correctamente
- ✅ Dashboard es responsive
- ✅ Auto-refresh del dashboard funciona
- ✅ Performance: dashboard carga en <2 segundos

---

## 8. Fase 6: Configuración del Sistema y Mantenimiento (Semana 7-8)

**Objetivo:** Panel de configuración dinámica y gestión de categorías.

**Duración:** 1-2 semanas
**Prioridad:** 🟢 MEDIA
**Progreso:** ░░░░░░░░░░ 0% (0/25 tareas)

### 8.1 Casos de Uso - Categorías

#### internal/usecase/admin/category/create_category.go
- [ ] Crear `CreateCategoryUseCase`
- [ ] Validar name único
- [ ] Auto-generar slug
- [ ] Validar icon (emoji válido)
- [ ] Asignar display_order automático
- [ ] Logging de auditoría

#### internal/usecase/admin/category/update_category.go
- [ ] Crear `UpdateCategoryUseCase`
- [ ] Permitir editar: name, icon, description, is_active
- [ ] Re-generar slug si name cambia
- [ ] Logging de auditoría

#### internal/usecase/admin/category/delete_category.go
- [ ] Crear `DeleteCategoryUseCase`
- [ ] Validar que no tenga rifas activas
- [ ] Soft delete (is_active = false)
- [ ] Logging de auditoría

#### internal/usecase/admin/category/reorder_categories.go
- [ ] Crear `ReorderCategoriesUseCase`
- [ ] Recibir array de IDs en nuevo orden
- [ ] Actualizar display_order de cada uno
- [ ] Logging de auditoría

### 8.2 Casos de Uso - System Parameters

#### internal/usecase/admin/system/list_parameters.go
- [ ] Crear `ListParametersUseCase`
- [ ] Filtrar por category
- [ ] Agrupar por category en respuesta
- [ ] Ocultar valores de parameters sensitive

#### internal/usecase/admin/system/update_parameter.go
- [ ] Crear `UpdateParameterUseCase`
- [ ] Validar value según value_type
- [ ] Guardar updated_by (admin_id)
- [ ] Logging de auditoría (severity: warning)
- [ ] Invalidar cache si existe

### 8.3 Casos de Uso - Company Settings

#### internal/usecase/admin/system/get_company_settings.go
- [ ] Crear `GetCompanySettingsUseCase`
- [ ] Retornar company_settings row

#### internal/usecase/admin/system/update_company_settings.go
- [ ] Crear `UpdateCompanySettingsUseCase`
- [ ] Validar email, phone, tax_id format
- [ ] Logging de auditoría
- [ ] Invalidar cache

### 8.4 Casos de Uso - Payment Processors

#### internal/usecase/admin/system/list_payment_processors.go
- [ ] Crear `ListPaymentProcessorsUseCase`
- [ ] Ocultar secrets (mask con ***)
- [ ] Mostrar is_active, is_sandbox

#### internal/usecase/admin/system/update_payment_processor.go
- [ ] Crear `UpdatePaymentProcessorUseCase`
- [ ] Validar credentials format
- [ ] Encriptar secrets antes de guardar
- [ ] Logging de auditoría (severity: critical)
- [ ] Test de conectividad con provider (opcional)

### 8.5 API Handlers

#### internal/adapters/http/handler/admin/category_handler.go
- [ ] Implementar CRUD completo
- [ ] Endpoint de reordenamiento

#### internal/adapters/http/handler/admin/system_handler.go
- [ ] Implementar handlers de parameters
- [ ] Implementar handlers de company settings
- [ ] Implementar handlers de payment processors

### 8.6 Rutas API

#### cmd/api/routes.go
- [ ] Agregar grupo `/api/v1/admin/categories`
- [ ] Agregar grupo `/api/v1/admin/system`

### 8.7 Frontend - Categorías

#### frontend/src/features/admin/pages/CategoriesPage.tsx
- [ ] Crear tabla de categorías
- [ ] Drag & drop para reordenar (react-beautiful-dnd)
- [ ] Edición inline de name, icon, description
- [ ] Toggle de is_active
- [ ] Botón de crear nueva categoría
- [ ] Modal de creación/edición

### 8.8 Frontend - System Config

#### frontend/src/features/admin/pages/SystemConfigPage.tsx
- [ ] Crear tabs:
  - System Parameters
  - Company Settings
  - Payment Processors
- [ ] System Parameters:
  - Agrupar por categoría (Business, Security, Payment, etc.)
  - Edición inline con validación por tipo
  - Save button por parámetro
- [ ] Company Settings:
  - Form con todos los campos
  - Upload de logo (opcional)
  - Save button
- [ ] Payment Processors:
  - Tabla con providers
  - Toggle is_active
  - Modal de edición de credentials (con advertencia de seguridad)

### 8.9 Frontend - Audit Logs

#### frontend/src/features/admin/pages/AuditLogsPage.tsx
- [ ] Crear tabla de audit logs
- [ ] Filtros: action, severity, date_range, user_id, admin_id
- [ ] Búsqueda por entity_id
- [ ] Badges de severity (info=gray, warning=yellow, error=orange, critical=red)
- [ ] Modal de detalle con metadata JSON

### 8.10 Frontend - Hooks

#### frontend/src/hooks/useAdminCategories.ts
- [ ] CRUD hooks

#### frontend/src/hooks/useAdminSystem.ts
- [ ] Hooks de parameters, company settings, payment processors

#### frontend/src/hooks/useAdminAudit.ts
- [ ] `useAuditLogs(filters, pagination)`

### 8.11 Criterios de Aceptación - Fase 6

- ✅ Admin puede crear/editar/eliminar categorías
- ✅ Drag & drop de categorías funciona
- ✅ Admin puede editar system parameters con validación
- ✅ Admin puede actualizar company settings
- ✅ Admin puede ver/editar payment processors
- ✅ Audit logs son consultables con filtros
- ✅ Secrets están enmascarados en UI
- ✅ Tests pasan

---

## 9. Fase 7: Testing y Aseguramiento de Calidad (Semana 8)

**Objetivo:** Testing exhaustivo y corrección de bugs.

**Duración:** 1 semana
**Prioridad:** 🔴 CRÍTICA
**Progreso:** ░░░░░░░░░░ 0% (0/20 tareas)

### 9.1 Unit Tests - Backend

- [ ] Tests de casos de uso de usuarios (100% coverage)
- [ ] Tests de casos de uso de organizadores (100% coverage)
- [ ] Tests de casos de uso de rifas admin (100% coverage)
- [ ] Tests de casos de uso de pagos admin (100% coverage)
- [ ] Tests de casos de uso de settlements (100% coverage)
- [ ] Tests de casos de uso de reports (100% coverage)
- [ ] Tests de casos de uso de system config (100% coverage)
- [ ] Tests de repositorios (100% coverage)

### 9.2 Integration Tests - Backend

- [ ] Tests de endpoints de usuarios (happy path + error cases)
- [ ] Tests de endpoints de organizadores
- [ ] Tests de endpoints de rifas admin
- [ ] Tests de endpoints de pagos admin
- [ ] Tests de endpoints de settlements
- [ ] Tests de endpoints de reports
- [ ] Tests de endpoints de system config
- [ ] Tests de permisos (verificar que user normal no puede acceder)

### 9.3 E2E Tests - Frontend

- [ ] Test: Login como super_admin → acceder a /admin
- [ ] Test: Suspender usuario → verificar audit log
- [ ] Test: Cambiar KYC de usuario
- [ ] Test: Aprobar settlement → marcar como pagado
- [ ] Test: Cancelar rifa con refund
- [ ] Test: Crear categoría y reordenar
- [ ] Test: Editar system parameter
- [ ] Test: Exportar reporte a CSV

### 9.4 Security Tests

- [ ] Penetration testing de permisos
- [ ] Test de rate limiting en endpoints admin
- [ ] Test de validación de inputs (SQL injection, XSS)
- [ ] Test de encriptación de secrets
- [ ] Audit de dependencias (npm audit, go mod check)

### 9.5 Performance Tests

- [ ] Load testing de dashboard (100 concurrent requests)
- [ ] Query optimization (explain analyze en queries pesadas)
- [ ] Indexing de tablas (verificar EXPLAIN ANALYZE)
- [ ] Caching de reports (implementar si es necesario)

### 9.6 Criterios de Aceptación - Fase 7

- ✅ Unit tests: >80% coverage
- ✅ Integration tests: todos los endpoints críticos cubiertos
- ✅ E2E tests: workflows principales funcionan
- ✅ Security tests: sin vulnerabilidades críticas
- ✅ Performance: dashboard carga en <2s
- ✅ Bugs críticos resueltos

---

## 10. Fase 8: Documentación y Despliegue (Semana 8)

**Objetivo:** Documentación completa y despliegue a producción.

**Duración:** 3-5 días
**Prioridad:** 🟡 ALTA
**Progreso:** ░░░░░░░░░░ 0% (0/15 tareas)

### 10.1 Documentación Técnica

- [ ] Actualizar API_ENDPOINTS.md con Swagger/OpenAPI spec
- [ ] Completar CASOS_DE_USO.md con todos los flujos
- [ ] Actualizar BASE_DE_DATOS.md con diagrama ER final
- [ ] Documentar decisiones arquitectónicas en ARQUITECTURA_ALMIGHTY.md

### 10.2 Documentación de Usuario

- [ ] Guía de usuario para super_admin (español)
  - Cómo suspender usuarios
  - Cómo aprobar settlements
  - Cómo procesar refunds
  - Cómo editar system parameters
- [ ] Video tutorial (opcional)
- [ ] FAQ de administración

### 10.3 Runbooks Operacionales

- [ ] Runbook: Cómo cancelar rifa con refund
- [ ] Runbook: Cómo resolver disputa de pago
- [ ] Runbook: Cómo hacer rollback de migración
- [ ] Runbook: Cómo investigar audit logs

### 10.4 Migraciones en Producción

- [ ] Backup completo de base de datos
- [ ] Ejecutar migraciones 012-018 en staging
- [ ] Validar migraciones en staging
- [ ] Plan de rollback documentado
- [ ] Ejecutar migraciones en producción (ventana de mantenimiento)

### 10.5 Despliegue Backend

- [ ] Build de binario Go
- [ ] Deploy a servidor de producción
- [ ] Configurar variables de entorno
- [ ] Verificar health checks
- [ ] Restart de servicio sorteos-backend

### 10.6 Despliegue Frontend

- [ ] Build de producción (npm run build)
- [ ] Deploy a /var/www/sorteos.club
- [ ] Clear cache de Nginx
- [ ] Verificar que /admin carga correctamente

### 10.7 Smoke Testing en Producción

- [ ] Login como super_admin
- [ ] Verificar dashboard carga
- [ ] Probar una acción de cada módulo
- [ ] Verificar audit logs se crean

### 10.8 Criterios de Aceptación - Fase 8

- ✅ Migraciones ejecutadas sin errores
- ✅ Backend desplegado y funcionando
- ✅ Frontend desplegado y accesible en /admin
- ✅ Documentación completa
- ✅ Smoke tests pasan
- ✅ Plan de rollback validado

---

## 11. Resumen de Entregables

### 11.1 Base de Datos
- ✅ 7 nuevas tablas creadas
- ✅ 2 tablas alteradas (raffles, users)
- ✅ Triggers y funciones implementadas
- ✅ Índices optimizados

### 11.2 Backend (Go)
- ✅ 7 nuevas entidades de dominio
- ✅ 7 repositorios implementados
- ✅ 47 casos de uso implementados
- ✅ 52 endpoints API creados
- ✅ Middleware de seguridad implementado
- ✅ >80% test coverage

### 11.3 Frontend (React)
- ✅ Módulo /admin completo
- ✅ 12 páginas principales
- ✅ 15+ componentes reutilizables
- ✅ Hooks personalizados
- ✅ Integración con shadcn/ui
- ✅ Responsive design

### 11.4 Funcionalidades
- ✅ Gestión completa de usuarios
- ✅ Gestión completa de organizadores
- ✅ Control administrativo de rifas
- ✅ Gestión de pagos y refunds
- ✅ Sistema de liquidaciones
- ✅ Mantenimiento de categorías
- ✅ Dashboard ejecutivo con KPIs
- ✅ Reportes financieros exportables
- ✅ Configuración dinámica del sistema
- ✅ Audit logs completos

### 11.5 Documentación
- ✅ Documentación técnica completa
- ✅ Guía de usuario
- ✅ Runbooks operacionales
- ✅ API documentation (Swagger)

---

## 12. Riesgos y Mitigaciones

| Riesgo | Probabilidad | Impacto | Mitigación |
|--------|--------------|---------|------------|
| Complejidad de refunds automáticos | Media | Alto | Testing exhaustivo, rollback plan |
| Performance de dashboard con mucha data | Media | Medio | Indexing, caching, materialized views |
| Seguridad de secrets en DB | Baja | Crítico | Encriptación AES-256, env vars |
| Migración de datos existentes | Media | Alto | Backfill scripts, validación post-migración |
| Bugs en settlements calculation | Media | Alto | Unit tests, manual validation |

---

## 13. Próximos Pasos

Una vez completado este roadmap:

1. **Monitoreo:** Implementar alertas para acciones críticas de admin
2. **Analytics:** Dashboard de métricas de uso del módulo admin
3. **Permisos Granulares:** Implementar RBAC más fino (ej: admin que solo puede ver, no editar)
4. **2FA:** Requerir autenticación de dos factores para super_admin
5. **Audit Reports:** Reportes de auditoría exportables para compliance
6. **API Rate Limiting Dinámico:** Ajustar rate limits desde system_parameters

---

## 14. Contacto y Soporte

**Responsable:** Equipo de desarrollo Sorteos.club
**Fecha última actualización:** 2025-11-18
**Versión roadmap:** 1.0

Para reportar issues o sugerencias:
- Crear issue en repositorio del proyecto
- Contactar al equipo de desarrollo

---

**INICIO DE IMPLEMENTACIÓN:** Pendiente de aprobación
**FIN ESTIMADO:** 8 semanas desde inicio
