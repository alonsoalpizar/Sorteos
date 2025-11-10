# Roadmap de Desarrollo - Plataforma de Sorteos

**Versión:** 1.0
**Fecha:** 2025-11-10
**Metodología:** Sprints de 2 semanas (Scrum adaptado)

---

## 1. Visión General

Este roadmap define las **3 fases principales** del proyecto, desde el MVP hasta la plataforma completa con aplicaciones móviles nativas. Cada fase incluye hitos medibles, criterios de aceptación y estimaciones realistas.

**Horizonte temporal:**
- **Fase 1 (MVP):** 8-10 semanas
- **Fase 2 (Escalamiento):** 10-12 semanas
- **Fase 3 (Expansión):** 12-16 semanas

---

## 2. Fase 1 - MVP (Producto Mínimo Viable)

**Objetivo:** Lanzar plataforma funcional con un único proveedor de pagos y funcionalidades core.

**Duración estimada:** 8-10 semanas (4-5 sprints)

---

### Sprint 1-2: Infraestructura y Autenticación ✅ COMPLETADO

**Fecha inicio:** 2025-11-10
**Fecha finalización:** 2025-11-10
**Estado Backend:** 100% completado ✅
**Estado Frontend:** 100% completado ✅
**Última actualización:** 2025-11-10 21:30

#### Tareas Backend
- [x] ✅ Setup proyecto Go con estructura hexagonal (2025-11-10)
  - go.mod con 40+ dependencias
  - Estructura de carpetas hexagonal (cmd, internal, pkg)
- [x] ✅ Configuración Docker Compose (Postgres, Redis) (2025-11-10)
  - PostgreSQL 15-alpine con health checks
  - Redis 7-alpine con persistencia
  - Adminer y Redis Commander (debug profile)
- [x] ✅ Migraciones base (users, user_consents, audit_logs) (2025-11-10)
  - 001_create_users_table: tabla users con ENUMs (role, kyc_level, status)
  - 002_create_user_consents_table: consentimientos GDPR
  - 003_create_audit_logs_table: auditoría con índices optimizados
- [x] ✅ Logging estructurado con Zap (2025-11-10)
  - pkg/logger/logger.go con diferentes niveles
- [x] ✅ Configuración Viper con .env (2025-11-10)
  - pkg/config/config.go con validaciones
  - .env.example con todas las variables
- [x] ✅ Entry point main.go (2025-11-10)
  - Servidor Gin con middlewares (CORS, logging, recovery, request ID)
  - Health checks (/health, /ready)
  - Conexión a PostgreSQL y Redis con pools
  - Graceful shutdown
- [x] ✅ Sistema de errores personalizado (2025-11-10)
  - pkg/errors/errors.go con códigos HTTP
- [x] ✅ Dockerfile multi-stage (2025-11-10)
- [x] ✅ Makefile con comandos de desarrollo (2025-11-10)
- [x] ✅ README.md completo con guías (2025-11-10)
- [x] ✅ Domain entities (2025-11-10 19:00)
  - internal/domain/user.go con validaciones (email, phone, password)
  - internal/domain/user_consent.go para GDPR
  - internal/domain/audit_log.go con builder pattern
- [x] ✅ User repository con GORM (2025-11-10 19:00)
  - internal/adapters/db/user_repository.go
  - CRUD completo con soft delete
  - Búsquedas optimizadas (email, phone, cedula)
  - Listado paginado con filtros
- [x] ✅ JWT Token Manager con Redis (2025-11-10 19:00)
  - internal/adapters/redis/token_manager.go
  - Generación de access/refresh tokens
  - Validación y rotación de tokens
  - Blacklist de tokens
  - Códigos de verificación
- [x] ✅ Crypto utilities (2025-11-10 19:00)
  - pkg/crypto/password.go (bcrypt cost 12)
  - pkg/crypto/code.go (códigos de 6 dígitos)
- [x] ✅ Use cases de autenticación (2025-11-10 19:00)
  - internal/usecase/auth/register.go
  - internal/usecase/auth/login.go
  - internal/usecase/auth/refresh_token.go
  - internal/usecase/auth/verify_email.go
- [x] ✅ HTTP handlers para autenticación (2025-11-10 20:00)
  - internal/adapters/http/handler/auth/register_handler.go
  - internal/adapters/http/handler/auth/login_handler.go
  - internal/adapters/http/handler/auth/refresh_token_handler.go
  - internal/adapters/http/handler/auth/verify_email_handler.go
- [x] ✅ Middlewares (2025-11-10 20:00)
  - internal/adapters/http/middleware/auth.go (JWT + Roles + KYC)
  - internal/adapters/http/middleware/rate_limit.go (Redis sliding window)
- [x] ✅ Repositorios adicionales (2025-11-10 20:00)
  - internal/adapters/db/user_consent_repository.go
  - internal/adapters/db/audit_log_repository.go
- [x] ✅ Integración SendGrid (2025-11-10 20:00)
  - internal/adapters/notifier/sendgrid.go
  - Templates HTML para emails de verificación
- [x] ✅ Rutas conectadas en main.go (2025-11-10 20:00)
  - cmd/api/routes.go con todas las rutas de auth
  - Endpoints: POST /api/v1/auth/{register,login,refresh,verify-email}

#### Tareas Frontend
- [x] ✅ Setup proyecto Vite + React + TypeScript (2025-11-10 21:00)
- [x] ✅ Configuración Tailwind + shadcn/ui (2025-11-10 21:00)
- [x] ✅ Componentes base (Button, Input, Card, Label, Alert, Badge) (2025-11-10 21:15)
- [x] ✅ Páginas: Register, Login, VerifyEmail, Dashboard (2025-11-10 21:30)
- [x] ✅ React Query setup con Axios (2025-11-10 21:10)
- [x] ✅ Zustand store para autenticación (2025-11-10 21:10)
- [x] ✅ Protected routes (2025-11-10 21:20)

#### Entregables
- ✅ Usuario puede registrarse, verificar cuenta y hacer login
- ✅ Tokens JWT funcionales con refresh automático
- ✅ Dark mode funcional

#### Archivos Creados (2025-11-10) - SISTEMA DE AUTENTICACIÓN COMPLETO
```
backend/
├── cmd/api/
│   ├── main.go                                        ✅ (actualizado)
│   └── routes.go                                      ✅ NEW
├── internal/
│   ├── domain/
│   │   ├── user.go                                    ✅ NEW
│   │   ├── user_consent.go                            ✅ NEW
│   │   └── audit_log.go                               ✅ NEW
│   ├── usecase/auth/
│   │   ├── register.go                                ✅ NEW
│   │   ├── login.go                                   ✅ NEW
│   │   ├── refresh_token.go                           ✅ NEW
│   │   └── verify_email.go                            ✅ NEW
│   └── adapters/
│       ├── db/
│       │   ├── user_repository.go                     ✅ NEW
│       │   ├── user_consent_repository.go             ✅ NEW
│       │   └── audit_log_repository.go                ✅ NEW
│       ├── redis/
│       │   └── token_manager.go                       ✅ NEW
│       ├── http/
│       │   ├── handler/auth/
│       │   │   ├── register_handler.go                ✅ NEW
│       │   │   ├── login_handler.go                   ✅ NEW
│       │   │   ├── refresh_token_handler.go           ✅ NEW
│       │   │   └── verify_email_handler.go            ✅ NEW
│       │   └── middleware/
│       │       ├── auth.go                            ✅ NEW
│       │       └── rate_limit.go                      ✅ NEW
│       └── notifier/
│           └── sendgrid.go                            ✅ NEW
├── pkg/
│   ├── config/config.go                               ✅
│   ├── logger/logger.go                               ✅ (actualizado)
│   ├── errors/errors.go                               ✅
│   └── crypto/
│       ├── password.go                                ✅ NEW
│       └── code.go                                    ✅ NEW
├── migrations/
│   ├── 001_create_users_table.up.sql                  ✅
│   ├── 001_create_users_table.down.sql                ✅
│   ├── 002_create_user_consents_table.up.sql          ✅
│   ├── 002_create_user_consents_table.down.sql        ✅
│   ├── 003_create_audit_logs_table.up.sql             ✅
│   └── 003_create_audit_logs_table.down.sql           ✅
├── go.mod                                             ✅
├── .env.example                                       ✅
├── .env                                               ✅
├── .gitignore                                         ✅
├── Dockerfile                                         ✅
├── .dockerignore                                      ✅
├── Makefile                                           ✅
└── README.md                                          ✅
docker-compose.yml                                     ✅ (actualizado)

frontend/
├── src/
│   ├── components/ui/
│   │   ├── Button.tsx                                     ✅ NEW
│   │   ├── Input.tsx                                      ✅ NEW
│   │   ├── Label.tsx                                      ✅ NEW
│   │   ├── Card.tsx                                       ✅ NEW
│   │   ├── Alert.tsx                                      ✅ NEW
│   │   └── Badge.tsx                                      ✅ NEW
│   ├── features/
│   │   ├── auth/
│   │   │   ├── pages/
│   │   │   │   ├── LoginPage.tsx                          ✅ NEW
│   │   │   │   ├── RegisterPage.tsx                       ✅ NEW
│   │   │   │   └── VerifyEmailPage.tsx                    ✅ NEW
│   │   │   └── components/
│   │   │       └── ProtectedRoute.tsx                     ✅ NEW
│   │   └── dashboard/
│   │       └── pages/
│   │           └── DashboardPage.tsx                      ✅ NEW
│   ├── lib/
│   │   ├── utils.ts                                       ✅ NEW
│   │   ├── api.ts                                         ✅ NEW
│   │   └── queryClient.ts                                 ✅ NEW
│   ├── store/
│   │   └── authStore.ts                                   ✅ NEW
│   ├── types/
│   │   └── auth.ts                                        ✅ NEW
│   ├── api/
│   │   └── auth.ts                                        ✅ NEW
│   ├── hooks/
│   │   └── useAuth.ts                                     ✅ NEW
│   ├── App.tsx                                            ✅ NEW
│   ├── main.tsx                                           ✅ NEW
│   ├── index.css                                          ✅ NEW
│   └── vite-env.d.ts                                      ✅ NEW
├── package.json                                           ✅ NEW
├── tsconfig.json                                          ✅ NEW
├── tsconfig.node.json                                     ✅ NEW
├── vite.config.ts                                         ✅ NEW
├── tailwind.config.js                                     ✅ NEW (COLORES APROBADOS)
├── postcss.config.js                                      ✅ NEW
└── index.html                                             ✅ NEW
```

**Total archivos nuevos en Sprint 1-2:**
- Backend: 22 archivos
- Frontend: 31 archivos
- **TOTAL: 53 archivos**

**Backend:**
- Domain: 3 archivos (User, UserConsent, AuditLog)
- Use Cases: 4 archivos (Register, Login, RefreshToken, VerifyEmail)
- Repositories: 3 archivos (User, UserConsent, AuditLog)
- Handlers: 4 archivos (Register, Login, Refresh, VerifyEmail)
- Middlewares: 2 archivos (Auth, RateLimit)
- Adapters: 2 archivos (TokenManager, SendGrid)
- Crypto: 2 archivos (Password, Code)
- Routes: 1 archivo (routes.go)
- Actualizados: 2 archivos (main.go, logger.go)

**Frontend:**
- Componentes UI: 6 archivos (Button, Input, Label, Card, Alert, Badge)
- Páginas: 4 archivos (Login, Register, VerifyEmail, Dashboard)
- Hooks: 1 archivo (useAuth con 8 hooks)
- Store: 1 archivo (authStore con Zustand)
- API Client: 2 archivos (api.ts, auth.ts)
- Types: 1 archivo (auth.ts con tipos completos)
- Utils: 2 archivos (utils.ts, queryClient.ts)
- Routing: 2 archivos (App.tsx, ProtectedRoute)
- Config: 7 archivos (package.json, tsconfig, vite, tailwind, postcss, html, css)

**Características Implementadas:**
- ✅ Sistema de autenticación completo (register, login, verify, logout)
- ✅ Gestión de tokens JWT con refresh automático
- ✅ Rate limiting por IP y usuario
- ✅ Validación de formularios con Zod
- ✅ Manejo de errores con UI feedback
- ✅ Dark mode support
- ✅ Protected routes
- ✅ Email templates con SendGrid
- ✅ Audit logging completo
- ✅ GDPR compliance (user consents)
- ✅ Responsive design con Tailwind
- ✅ COLORES APROBADOS: Blue #3B82F6 / Slate #64748B (NO purple/pink)

---

### Sprint 3-4: Gestión de Sorteos (CRUD Básico) ✅ COMPLETADO

**Fecha inicio:** 2025-11-10
**Fecha finalización:** 2025-11-10
**Estado Backend:** 100% completado ✅
**Estado Frontend:** 100% completado ✅
**Última actualización:** 2025-11-10 08:35

#### Tareas Backend
- [x] ✅ Migraciones: raffles, raffle_numbers, raffle_images (2025-11-10 07:25)
  - 004_create_raffles_table: tabla raffles con ENUMs (status, draw_method, settlement_status)
  - 005_create_raffle_numbers_table: tabla raffle_numbers con ENUM (status: available/reserved/sold)
  - 006_create_raffle_images_table: tabla raffle_images con validaciones MIME y tamaño
  - Triggers automáticos para updated_at, revenue calculation
  - Función para liberar reservas expiradas
- [x] ✅ Domain entities (2025-11-10 06:15)
  - internal/domain/raffle.go: 15+ métodos de negocio (Publish, Suspend, Complete, etc.)
  - internal/domain/raffle_number.go: gestión de reservas con TTL
  - internal/domain/raffle_image.go: validación de archivos (MIME types, size limits)
- [x] ✅ Repositorios GORM para sorteos (2025-11-10 06:20)
  - internal/adapters/db/raffle_repository.go: 16 métodos (CRUD, búsquedas, filtros)
  - internal/adapters/db/raffle_number_repository.go: 14 métodos (batch creation, reservations)
  - internal/adapters/db/raffle_image_repository.go: 10 métodos (primary image logic)
- [x] ✅ Casos de uso (2025-11-10 07:30)
  - CreateRaffle (con validaciones, generación de números, audit log) ✅
  - ListRaffles (paginación, filtros por estado) ✅
  - GetRaffleDetail (con números disponibles) ✅
  - PublishRaffle (validaciones completas de publicación) ✅
  - UpdateRaffle (solo owner o admin) ✅
  - SuspendRaffle (admin only) ✅
  - DeleteRaffle (soft delete, owner o admin) ✅
- [x] ✅ HTTP Handlers (2025-11-10 07:40)
  - CreateRaffleHandler: POST /api/v1/raffles ✅
  - ListRafflesHandler: GET /api/v1/raffles ✅
  - GetRaffleDetailHandler: GET /api/v1/raffles/:id ✅
  - PublishRaffleHandler: POST /api/v1/raffles/:id/publish ✅
  - UpdateRaffleHandler: PUT /api/v1/raffles/:id ✅
  - SuspendRaffleHandler: POST /api/v1/raffles/:id/suspend (admin) ✅
  - DeleteRaffleHandler: DELETE /api/v1/raffles/:id ✅
- [x] ✅ Generación automática de rango de números (2025-11-10 06:25)
  - Números formateados (00-99, 000-999 según cantidad)
  - Creación en batch (100 números por lote)
- [x] ✅ Rutas conectadas en main.go (2025-11-10 07:40)
  - cmd/api/routes.go: función setupRaffleRoutes() con 7 endpoints
  - Rutas públicas (GET raffles list y detail)
  - Rutas protegidas con autenticación + KYC (POST, PUT, DELETE)
  - Rutas admin (POST suspend)
  - Rate limiting en creación de sorteos (10/hora)
- [ ] Upload de imágenes (S3 o local storage) ⏳ PENDIENTE (Sprint 5-6)
- [ ] Cache Redis de sorteos activos ⏳ PENDIENTE (Sprint 5-6)

#### Tareas Frontend
- [x] ✅ Tipos TypeScript (2025-11-10 08:25)
  - src/types/raffle.ts: tipos completos para sorteos, números, imágenes
- [x] ✅ API Client (2025-11-10 08:26)
  - src/api/raffles.ts: cliente HTTP con 7 endpoints
- [x] ✅ Custom Hooks con React Query (2025-11-10 08:27)
  - useRafflesList, useRaffleDetail, useCreateRaffle
  - useUpdateRaffle, usePublishRaffle, useDeleteRaffle, useSuspendRaffle
- [x] ✅ Componentes (2025-11-10 08:30)
  - RaffleCard: card con preview, barra de progreso, stats
  - NumberGrid: grid de números 00-99 con estados visuales
- [x] ✅ Páginas (2025-11-10 08:33)
  - RafflesListPage: listado con filtros y paginación
  - RaffleDetailPage: detalle completo con acciones
  - CreateRafflePage: formulario de creación con validaciones
- [x] ✅ Rutas configuradas en App.tsx (2025-11-10 08:34)
  - Rutas públicas: /raffles, /raffles/:id
  - Rutas protegidas: /raffles/create
- [x] ✅ Utilidades y componentes actualizados (2025-11-10 08:34)
  - Badge: variantes info, error agregadas
  - Alert: variantes info, error agregadas
  - utils.ts: getStatusColor, getStatusLabel, getDrawMethodLabel
  - useAuth: hook agregado
- [ ] ImageUploader ⏳ PENDIENTE (Sprint 5-6)
- [ ] Página de editar sorteo ⏳ PENDIENTE (futuro)

#### Entregables Completados
- ✅ Usuario puede crear sorteo con detalles completos (title, description, price, numbers, draw date/method)
- ✅ Usuario puede listar sorteos públicos con paginación y filtros
- ✅ Usuario puede ver detalle de sorteo con números disponibles/reservados/vendidos
- ✅ Usuario puede publicar sorteo (con validaciones: imágenes, números, fecha futura)
- ✅ Usuario puede actualizar sorteo (title, description, draw date) si no tiene ventas
- ✅ Administrador puede suspender sorteos
- ✅ Usuario puede eliminar sorteos (soft delete) si no tienen ventas
- ✅ Vista pública de sorteos activos con grid responsive
- ✅ UI para crear sorteos con formulario completo y validaciones
- ✅ Vista de detalle con grid de números y acciones para owner/admin

#### Archivos Creados Sprint 3-4 (2025-11-10) - GESTIÓN DE SORTEOS (Fullstack 100% ✅)
```
backend/
├── migrations/
│   ├── 004_create_raffles_table.up.sql                ✅ NEW
│   ├── 004_create_raffles_table.down.sql              ✅ NEW
│   ├── 005_create_raffle_numbers_table.up.sql         ✅ NEW
│   ├── 005_create_raffle_numbers_table.down.sql       ✅ NEW
│   ├── 006_create_raffle_images_table.up.sql          ✅ NEW
│   └── 006_create_raffle_images_table.down.sql        ✅ NEW
├── internal/
│   ├── domain/
│   │   ├── raffle.go                                  ✅ NEW (actualizado: Metadata datatypes.JSON)
│   │   ├── raffle_number.go                           ✅ NEW
│   │   └── raffle_image.go                            ✅ NEW
│   ├── usecase/raffle/
│   │   ├── create_raffle.go                           ✅ NEW
│   │   ├── list_raffles.go                            ✅ NEW
│   │   ├── get_raffle_detail.go                       ✅ NEW
│   │   ├── publish_raffle.go                          ✅ NEW
│   │   └── update_raffle.go                           ✅ NEW (3 use cases: Update, Suspend, Delete)
│   ├── adapters/
│   │   ├── db/
│   │   │   ├── raffle_repository.go                   ✅ NEW
│   │   │   ├── raffle_number_repository.go            ✅ NEW
│   │   │   └── raffle_image_repository.go             ✅ NEW
│   │   └── http/handler/raffle/
│   │       ├── create_raffle_handler.go               ✅ NEW
│   │       ├── list_raffles_handler.go                ✅ NEW
│   │       ├── get_raffle_detail_handler.go           ✅ NEW
│   │       ├── publish_raffle_handler.go              ✅ NEW
│   │       ├── update_raffle_handler.go               ✅ NEW (3 handlers)
│   │       └── common.go                              ✅ NEW (DTOs y error handling)
├── cmd/api/
│   ├── main.go                                        ✅ (actualizado: +setupRaffleRoutes)
│   └── routes.go                                      ✅ (actualizado: +setupRaffleRoutes func)
├── go.mod                                             ✅ (actualizado: +shopspring/decimal +datatypes)
└── go.sum                                             ✅ (actualizado)

frontend/
├── src/
│   ├── types/
│   │   └── raffle.ts                                  ✅ NEW
│   ├── api/
│   │   └── raffles.ts                                 ✅ NEW
│   ├── hooks/
│   │   ├── useRaffles.ts                              ✅ NEW
│   │   └── useAuth.ts                                 ✅ (actualizado: +useAuth)
│   ├── features/raffles/
│   │   ├── components/
│   │   │   ├── RaffleCard.tsx                         ✅ NEW
│   │   │   └── NumberGrid.tsx                         ✅ NEW
│   │   └── pages/
│   │       ├── RafflesListPage.tsx                    ✅ NEW
│   │       ├── RaffleDetailPage.tsx                   ✅ NEW
│   │       └── CreateRafflePage.tsx                   ✅ NEW
│   ├── components/ui/
│   │   ├── Badge.tsx                                  ✅ (actualizado: +info +error)
│   │   └── Alert.tsx                                  ✅ (actualizado: +info +error)
│   ├── lib/
│   │   ├── utils.ts                                   ✅ (actualizado: +3 funciones)
│   │   └── api.ts                                     ✅ (usado como apiClient)
│   └── App.tsx                                        ✅ (actualizado: +rutas raffles)
```

**Total archivos nuevos en Sprint 3-4:**

Backend:
- Migraciones: 6 archivos (3 up + 3 down)
- Domain: 3 archivos (Raffle, RaffleNumber, RaffleImage)
- Use Cases: 5 archivos (Create, List, GetDetail, Publish, Update/Suspend/Delete)
- Repositories: 3 archivos (Raffle, RaffleNumber, RaffleImage)
- Handlers: 6 archivos (Create, List, GetDetail, Publish, Update, Common)
- Config: 2 archivos actualizados (main.go, routes.go)
- **Subtotal Backend: 23 archivos creados + 2 actualizados**

Frontend:
- Types: 1 archivo (raffle.ts)
- API Client: 1 archivo (raffles.ts)
- Hooks: 1 archivo nuevo (useRaffles.ts) + 1 actualizado (useAuth.ts)
- Componentes: 2 archivos (RaffleCard, NumberGrid)
- Páginas: 3 archivos (List, Detail, Create)
- UI Components: 2 actualizados (Badge, Alert)
- Lib: 1 actualizado (utils.ts)
- Config: 1 actualizado (App.tsx)
- **Subtotal Frontend: 8 archivos creados + 5 actualizados**

**TOTAL SPRINT 3-4: 31 archivos creados + 7 actualizados**

**Dependencias añadidas:**
- github.com/shopspring/decimal v1.3.1 (aritmética decimal precisa para dinero)
- gorm.io/datatypes v1.2.0 (soporte para campos JSON en PostgreSQL)

**Características Implementadas:**
- ✅ Sistema de sorteos con ENUMs (draft, active, suspended, completed, cancelled)
- ✅ Sistema de reserva de números con TTL (Time To Live)
- ✅ Cálculo automático de revenue vía triggers de base de datos
- ✅ Gestión de imágenes con validaciones (MIME type, file size)
- ✅ Creación de sorteos con generación automática de números
- ✅ Soft delete en todas las tablas
- ✅ Audit logging integrado con builder pattern
- ✅ Soporte para múltiples métodos de sorteo (loteria_nacional_cr, manual, random)
- ✅ Settlement tracking (pending, processing, completed, failed)
- ✅ Platform fee configurable (default 10%)
- ✅ Función PostgreSQL para liberar reservas expiradas (preparado para cron job)
- ✅ Listado paginado con filtros (status, search, user_id)
- ✅ Detalle de sorteo con conteo de números (disponibles/reservados/vendidos)
- ✅ Validaciones de publicación (imágenes, números, fecha futura)
- ✅ Restricciones de edición para sorteos con ventas
- ✅ Sistema de permisos (owner o admin para ciertas acciones)
- ✅ 7 endpoints HTTP REST funcionales con rate limiting

**Endpoints Backend Implementados:**
- GET /api/v1/raffles - Listar sorteos (público)
- GET /api/v1/raffles/:id - Detalle de sorteo (público)
- POST /api/v1/raffles - Crear sorteo (autenticado + KYC + rate limit 10/hora)
- PUT /api/v1/raffles/:id - Actualizar sorteo (autenticado + KYC + owner/admin)
- POST /api/v1/raffles/:id/publish - Publicar sorteo (autenticado + KYC + owner)
- DELETE /api/v1/raffles/:id - Eliminar sorteo (autenticado + KYC + owner/admin)
- POST /api/v1/raffles/:id/suspend - Suspender sorteo (admin only)

**Rutas Frontend Implementadas:**
- GET /raffles - Listado de sorteos (público)
- GET /raffles/:id - Detalle de sorteo (público)
- GET /raffles/create - Crear sorteo (protegido: auth + KYC)
- GET / - Redirige a /raffles

**Características Frontend:**
- ✅ Grid responsive de sorteos con cards
- ✅ Filtros por estado (todos, activos, borradores, completados, cancelados)
- ✅ Búsqueda por título o descripción
- ✅ Paginación funcional
- ✅ Barra de progreso de ventas en cada card
- ✅ Grid de números 00-99 con estados visuales (disponible, reservado, vendido)
- ✅ Leyenda de colores para números
- ✅ Formulario de creación con validaciones en tiempo real
- ✅ Resumen con cálculo automático de recaudación
- ✅ Acciones para owner/admin (publicar, editar, eliminar)
- ✅ Dark mode support completo
- ✅ Loading states y error handling
- ✅ Badges con colores según estado del sorteo
- ✅ Alertas informativas (success, warning, error, info)
- ✅ React Query para cache y sincronización
- ✅ Zustand para estado global de autenticación

---

### Sprint 5-6: Reservas y Pagos

#### Tareas Backend
- [ ] Migraciones: reservations, payments, idempotency_keys
- [ ] Sistema de reserva temporal:
  - Lock distribuido Redis por número
  - Crear reserva (status=pending, expires_at=now+5min)
  - Cron job para liberar reservas expiradas
- [ ] Integración con PSP (Stripe como primera opción):
  - Interfaz PaymentProvider
  - Implementación StripeProvider
  - Manejo de webhooks (payment.succeeded, payment.failed)
  - Idempotencia con Idempotency-Key
- [ ] Flujo completo:
  1. POST /raffles/{id}/reservations → crea reserva + lock
  2. POST /payments → intenta cargo con Stripe
  3. Webhook confirma → marca números como sold
  4. Si falla/expira → libera números
- [ ] Tests de concurrencia (vegeta/k6)

#### Tareas Frontend
- [ ] Página de checkout:
  - Selección de números (click en NumberGrid)
  - Carrito temporal (Zustand)
  - Formulario de pago (Stripe Elements)
  - Pantalla de confirmación
- [ ] Componentes:
  - NumberSelector (multi-selección)
  - PaymentForm (iframe Stripe o tarjeta directa)
  - OrderSummary (precio, fees, total)
- [ ] Manejo de estados:
  - Reserva pendiente (timer 5 min)
  - Pago procesando (spinner)
  - Pago exitoso (confetti + redirect)
  - Pago fallido (reintentar)

#### Entregables
- Usuario puede reservar números y pagar con tarjeta
- Números no se duplican (prueba con 500 req concurrentes)
- Reservas expiradas se liberan automáticamente
- Webhooks procesan pagos correctamente

---

### Sprint 7-8: Selección de Ganador y Backoffice Mínimo

#### Tareas Backend
- [ ] Sistema de selección de ganador:
  - Integración con API Lotería Nacional (o mock)
  - Cron job que consulta resultados en draw_date
  - Marca ganadores en raffle_numbers
  - Notificación por email/SMS al ganador
- [ ] Endpoints backoffice:
  - GET /admin/raffles (listado completo con filtros)
  - PATCH /admin/raffles/{id} (suspender/activar)
  - GET /admin/users (con filtros KYC)
  - POST /admin/settlements (crear liquidación manual)
- [ ] Audit log para todas las acciones de admin

#### Tareas Frontend
- [ ] Panel de usuario (dashboard):
  - Mis sorteos publicados (estados, % vendido)
  - Sorteos en los que participé
  - Sorteos ganados
  - Historial de pagos
- [ ] Panel de admin (backoffice básico):
  - Listado de sorteos con acciones (suspender/activar)
  - Listado de usuarios (verificar/suspender)
  - Vista de liquidaciones pendientes
- [ ] Componentes:
  - DataTable reutilizable (con sorting, paginación)
  - StatusBadge (draft/active/suspended/completed)
  - ActionMenu (suspender, editar, ver detalles)

#### Entregables
- Ganadores se determinan automáticamente según lotería
- Usuario recibe notificación al ganar
- Admin puede gestionar sorteos y usuarios desde backoffice
- Todas las acciones de admin quedan registradas (audit log)

---

### Sprint 9-10: Testing, Optimización y Lanzamiento MVP

#### Tareas
- [ ] Tests de aceptación:
  - Flujo completo end-to-end (Playwright/Cypress)
  - Pruebas de carga (k6): 1000 usuarios concurrentes
  - Pruebas de seguridad (OWASP ZAP)
- [ ] Optimizaciones:
  - Índices de base de datos (EXPLAIN ANALYZE)
  - Lazy loading de imágenes
  - Code splitting en React
  - CDN para assets estáticos
- [ ] Documentación:
  - README con setup instructions
  - API docs (Swagger/OpenAPI)
  - Guía de usuario (screenshots)
- [ ] Deploy a staging:
  - CI/CD pipeline completo
  - Health checks y rollback automático
  - Monitoreo con Prometheus + Grafana
- [ ] Beta testing con 50 usuarios reales
- [ ] Corrección de bugs críticos

#### Entregables
- MVP en producción con dominio custom
- Métricas de rendimiento (p95 < 500ms)
- Documentación completa para usuarios y desarrolladores

---

## 3. Fase 2 - Escalamiento y Funcionalidades Avanzadas

**Objetivo:** Expandir capacidades de la plataforma y preparar para crecimiento.

**Duración estimada:** 10-12 semanas (5-6 sprints)

---

### Sprint 11-12: Múltiples PSPs y Modo "Sin Cobro"

#### Backend
- [ ] Implementar providers adicionales:
  - PayPalProvider
  - LocalCRProvider (procesador de CR por definir)
- [ ] Sistema de routing de pagos:
  - Feature flags por sorteo (Stripe/PayPal/Local)
  - Fallback automático si PSP falla
- [ ] Modo "sin cobro en plataforma":
  - Sorteos gratuitos (owner coordina pago fuera)
  - Solo cobro de suscripción mensual al owner
  - Modelo de suscripción (Stripe Billing)

#### Frontend
- [ ] Selector de método de pago en checkout
- [ ] Modal de suscripción (planes Basic/Pro)
- [ ] Dashboard de owner con estado de suscripción

#### Entregables
- Usuario puede pagar con Stripe, PayPal o método local
- Owners pueden publicar sorteos sin cobro + pagar suscripción

---

### Sprint 13-14: Búsqueda Avanzada y Sistema de Afiliados

#### Backend
- [ ] Full-text search con PostgreSQL (pg_trgm):
  - Búsqueda por título, descripción, categoría
  - Filtros combinados (precio, fecha, % vendido)
  - Ordenamiento por relevancia
- [ ] Sistema de afiliados:
  - Tabla affiliate_links (user_id, code, clicks, conversions)
  - Endpoint para generar link único
  - Tracking de registros por afiliado
  - Cálculo de comisiones

#### Frontend
- [ ] Barra de búsqueda con autocomplete
- [ ] Filtros avanzados (sidebar)
- [ ] Panel de afiliados (generar link, estadísticas)

#### Entregables
- Búsqueda rápida y precisa de sorteos
- Usuarios pueden generar links de afiliado y ganar comisiones

---

### Sprint 15-16: Multilenguaje y Comunicación entre Usuarios

#### Backend
- [ ] i18n en backend (mensajes de error, emails)
- [ ] Sistema de mensajería privada:
  - Tabla messages (sender_id, receiver_id, content, read_at)
  - Notificaciones en tiempo real (WebSockets)

#### Frontend
- [ ] Selector de idioma (Español/Inglés)
- [ ] Inbox de mensajes (estilo chat)
- [ ] Notificaciones en tiempo real (toast)

#### Entregables
- Plataforma disponible en ES/EN
- Usuarios pueden comunicarse vía mensajes privados

---

### Sprint 17-18: Comentarios, Valoraciones e Integración con Redes Sociales

#### Backend
- [ ] Sistema de reviews:
  - Tabla reviews (raffle_id, user_id, rating, comment)
  - Moderación (admin puede ocultar reviews)
- [ ] Open Graph tags dinámicos (meta tags para compartir)

#### Frontend
- [ ] Sección de comentarios en detalle de sorteo
- [ ] Botones de compartir (Facebook, Twitter, WhatsApp)
- [ ] Modal de valoración post-sorteo

#### Entregables
- Usuarios pueden comentar y valorar sorteos
- Compartir en redes sociales genera preview atractivo

---

### Sprint 19-20: Notificaciones en Tiempo Real y Dashboards Avanzados

#### Backend
- [ ] WebSockets para eventos en vivo:
  - Nuevo sorteo publicado
  - Sorteo próximo a cerrarse
  - Ganador anunciado
- [ ] Vistas materializadas para KPIs:
  - Total vendido por sorteo/usuario/período
  - Tasa de conversión reserva → pago
  - Top sorteos por ingresos

#### Frontend
- [ ] Dashboard de owner con gráficos (Chart.js):
  - Ingresos por mes
  - % de vendido por sorteo
  - Tasa de conversión
- [ ] Notificaciones push (PWA)

#### Entregables
- Notificaciones en tiempo real funcionales
- Dashboards con métricas accionables para owners

---

### Sprint 21-22: Optimización y Preparación para Escala

#### Tareas
- [ ] Caching agresivo:
  - CDN para imágenes (CloudFront/Cloudflare)
  - Cache de listados en Redis (invalidación inteligente)
- [ ] Database tuning:
  - Índices compuestos optimizados
  - Particionamiento de tablas grandes (audit_logs)
- [ ] Horizontal scaling:
  - Balanceador de carga (Nginx/HAProxy)
  - Réplicas de lectura en Postgres
- [ ] Pruebas de carga: 10k usuarios concurrentes

#### Entregables
- Plataforma soporta 10k usuarios simultáneos
- Latencia p95 < 300ms en operaciones críticas

---

## 4. Fase 3 - Expansión y Aplicaciones Móviles

**Objetivo:** Alcance global y experiencia móvil nativa.

**Duración estimada:** 12-16 semanas (6-8 sprints)

---

### Sprint 23-26: Aplicación Móvil (React Native)

#### Tareas
- [ ] Setup React Native con TypeScript
- [ ] Compartir lógica con web (custom hooks)
- [ ] Pantallas principales:
  - Login/Register
  - Listado y detalle de sorteos
  - Checkout con Apple Pay / Google Pay
  - Dashboard de usuario
- [ ] Push notifications (FCM)
- [ ] Deep links (abrir sorteo desde notificación)
- [ ] Beta en TestFlight / Google Play Beta

#### Entregables
- Apps nativas iOS + Android en beta pública
- Notificaciones push funcionales

---

### Sprint 27-30: Sorteos Temáticos y Campañas Automatizadas

#### Backend
- [ ] Taxonomía de categorías (Viajes, Tecnología, Moda, etc.)
- [ ] Sistema de tags y recomendaciones
- [ ] Integración con herramienta de marketing automation (HubSpot/Mailchimp):
  - Campañas por email basadas en comportamiento
  - Segmentación de usuarios

#### Frontend
- [ ] Landing pages por categoría
- [ ] Recomendaciones personalizadas
- [ ] Builder de campañas (admin)

#### Entregables
- Sorteos organizados por temas
- Campañas automatizadas de email marketing

---

### Sprint 31-34: Analytics Avanzado y A/B Testing

#### Backend
- [ ] Integración con Google Analytics 4
- [ ] Events tracking personalizado
- [ ] Sistema de feature flags (LaunchDarkly/Unleash)

#### Frontend
- [ ] Dashboards de analytics para owners
- [ ] A/B testing en páginas clave (checkout, landing)

#### Entregables
- Análisis detallado de comportamiento de usuarios
- Optimización basada en datos (A/B tests)

---

### Sprint 35-38: Programa de Fidelización y Gamificación

#### Backend
- [ ] Sistema de puntos y niveles:
  - Puntos por compra, referido, compartir
  - Niveles (Bronce, Plata, Oro)
  - Recompensas (descuentos, boletos gratis)
- [ ] Tabla de logros (achievements)

#### Frontend
- [ ] Perfil con badges y nivel actual
- [ ] Marketplace de recompensas
- [ ] Animaciones de logros desbloqueados

#### Entregables
- Sistema de fidelización activo
- Incremento en retención de usuarios (meta: +20%)

---

## 5. Hitos Críticos

| Hito | Fecha Estimada | Criterio de Éxito |
|------|----------------|-------------------|
| MVP Lanzado | Semana 10 | 100 sorteos publicados, 500 usuarios registrados |
| 1 PSP Adicional | Semana 14 | 30% de pagos con PSP alternativo |
| App Móvil Beta | Semana 26 | 1000 descargas en beta |
| 10k Usuarios Activos | Semana 32 | 10k MAU con < 300ms p95 latency |

---

## 6. Riesgos y Mitigaciones

| Riesgo | Probabilidad | Impacto | Mitigación |
|--------|--------------|---------|------------|
| Integración PSP falla | Media | Alto | Mock provider para tests, fallback automático |
| Doble venta de números | Baja | Crítico | Tests de concurrencia en CI, locks distribuidos |
| Escalado de DB | Media | Alto | Réplicas de lectura, caché agresivo |
| Retraso en app móvil | Alta | Medio | Priorizar web, liberar móvil en Fase 3.5 si necesario |

---

## 7. Recursos Necesarios

### Equipo Mínimo (Fase 1)
- 1 Backend Developer (Go)
- 1 Frontend Developer (React)
- 1 Full-Stack Developer (Go + React)
- 1 DevOps (part-time)
- 1 QA (part-time)

### Equipo Fase 2-3
- +1 Backend Developer
- +1 Mobile Developer (React Native)
- +1 UX/UI Designer
- DevOps full-time

---

## 8. Presupuesto Estimado (Infraestructura)

**Fase 1 (MVP):**
- AWS/DigitalOcean: $100-200/mes
- Stripe fees: 2.9% + $0.30 por transacción
- SendGrid: $15/mes (40k emails)
- Twilio: ~$0.01/SMS

**Fase 2:**
- Infra: $300-500/mes (réplicas, CDN)
- Multiple PSPs: fees variables

**Fase 3:**
- Infra: $800-1200/mes (app móvil, analytics)

---

## 9. Métricas de Éxito por Fase

**Fase 1 (MVP):**
- 500 usuarios registrados
- 100 sorteos publicados
- 70% tasa de conversión reserva → pago
- 0 incidentes de doble venta

**Fase 2:**
- 5000 usuarios activos mensuales (MAU)
- 3 PSPs integrados
- NPS > 40

**Fase 3:**
- 20k MAU
- Apps móviles con 4.5+ estrellas
- 80% retención mensual

---

## 10. Dependencias Externas

- **API Lotería Nacional de Costa Rica:** Confirmación de disponibilidad y documentación
- **PSP Local (CR):** Identificar y firmar contrato antes de Sprint 11
- **Revisión legal:** Términos, privacidad, compliance con regulaciones de sorteos

---

## 11. Próximos Pasos Inmediatos

1. **Definir stack de desarrollo:** ✅ Completado (ver [stack_tecnico.md](./stack_tecnico.md))
2. **Crear estructura de carpetas:** ✅ Completado (2025-11-10)
3. **Setup repositorio Git:** ⏳ Pendiente
4. **Diseño de base de datos:** ✅ Migraciones iniciales completadas (users, user_consents, audit_logs)
5. **Sprint 1-2 (Infraestructura):** ⏳ 60% completado (2025-11-10)

### Próximas Tareas (Sprint 1-2 continuación)

**Backend:**
1. Implementar domain entities (`internal/domain/user.go`)
2. Implementar user repository (`internal/adapters/db/user_repository.go`)
3. Implementar JWT token manager (`internal/adapters/redis/token_manager.go`)
4. Implementar use cases de autenticación (`internal/usecase/auth/`)
5. Implementar HTTP handlers (`internal/adapters/http/handler/auth/`)
6. Implementar rate limiting middleware
7. Integrar SendGrid para emails

**Frontend:**
1. Setup Vite + React + TypeScript
2. Configurar Tailwind CSS + shadcn/ui
3. Crear componentes base
4. Implementar páginas de autenticación
5. Configurar React Query y Zustand

---

**Actualizado:** 2025-11-10 18:30
**Próxima revisión:** Después de completar Sprint 1-2
**Última modificación:** Actualizado progreso de infraestructura backend (60% completado)
