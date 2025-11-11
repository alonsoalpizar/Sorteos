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
**Estado Deployment:** 100% completado ✅
**URL Producción:** https://sorteos.club
**Última actualización:** 2025-11-10 08:50

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

#### Tareas Deployment
- [x] ✅ Integración Frontend en contenedor Docker (2025-11-10 08:42)
  - Multi-stage build: Node (frontend) + Go (backend) + Alpine (runtime)
  - Frontend servido por backend en /assets y / (SPA)
- [x] ✅ Configuración Nginx como reverse proxy (2025-11-10 08:47)
  - SSL/TLS con Let's Encrypt (https://sorteos.club)
  - HTTP → HTTPS redirect
  - www → non-www redirect
  - Compresión gzip
  - Headers de seguridad (HSTS, X-Frame-Options, etc.)
- [x] ✅ Dominio sorteos.club configurado (2025-11-10 08:47)
  - DNS apuntando a 62.171.188.255
  - Certificado SSL válido
- [x] ✅ Fix rutas API frontend (2025-11-10 08:50)
  - Actualizado baseURL: /api → /api/v1
  - Corregidas rutas auth: /v1/auth → /auth

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
- ⚠️ Validaciones de publicación (imágenes, números, fecha futura) - **Validación de imágenes temporalmente deshabilitada**
  - **NOTA (2025-11-11):** Upload de imágenes no implementado aún → validaciones comentadas en `publish_raffle.go`
  - **TODO:** Re-habilitar cuando Sprint 4 (Image Upload) esté completo
  - Ver: Issues Resueltos en Sprint 5-6 para más detalles
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

### Sprint 3.5: Mejora UX/UI - Navegación y Experiencia de Usuario ✅ COMPLETADO

**Fecha inicio:** 2025-11-10 18:00
**Fecha finalización:** 2025-11-10 18:30
**Estado Frontend:** 100% completado ✅
**Estado Deployment:** 100% completado ✅
**URL Producción:** https://sorteos.club
**Última actualización:** 2025-11-10 18:30

#### Contexto
Después del Sprint 3-4, identificamos que la interfaz estaba muy limitada:
- Dashboard sin enlaces útiles (crear sorteo, ver sorteos disponibles)
- Falta de distinción clara entre experiencia de comprador vs vendedor
- Navegación confusa sin menús persistentes
- Páginas faltantes (Mis Sorteos, Mis Compras)

Se decidió hacer una revisión completa de UX/UI antes de continuar con Sprint 5-6 (Pagos).

#### Tareas Completadas

**Estructura de Navegación:**
- [x] ✅ Navbar persistente con logo, search y user menu (2025-11-10 18:05)
  - Logo con link a home
  - Barra de búsqueda para usuarios autenticados
  - Enlaces de navegación (Explorar Sorteos, Crear Sorteo)
  - Menú de usuario con dropdown
  - Responsive con menú mobile

- [x] ✅ UserMenu dropdown component (2025-11-10 18:07)
  - Avatar con iniciales del usuario
  - Información del usuario (nombre, email)
  - Links rápidos: Dashboard, Mis Sorteos, Mis Compras
  - Botón de logout
  - Login/Register para no autenticados

- [x] ✅ MainLayout wrapper component (2025-11-10 18:10)
  - Navbar persistente
  - Footer con links útiles
  - Aplicado a todas las rutas protegidas

**Componentes Reutilizables:**
- [x] ✅ StatsCard - Card para mostrar estadísticas con icono (2025-11-10 18:12)
- [x] ✅ EmptyState - Placeholder con acción para estados vacíos (2025-11-10 18:13)
- [x] ✅ LoadingSpinner - Indicador de carga con texto opcional (2025-11-10 18:14)

**Páginas Mejoradas:**

- [x] ✅ DashboardPage rediseñado completamente (2025-11-10 18:16)
  - Welcome section personalizado
  - Quick actions: Crear Sorteo, Explorar, Mis Sorteos
  - Stats overview: Sorteos Activos, Ventas Totales, Compras Pendientes, Participaciones
  - Recent activity section (preparado para datos reales)
  - Account information section

- [x] ✅ MyRafflesPage - Vista de vendedor (2025-11-10 18:20)
  - Filtros por estado (Todos, Borrador, Activo, Suspendido, Completado, Cancelado)
  - Tabla con: título, estado, progreso de ventas, ingresos, fecha sorteo, acciones
  - Progress bars visuales
  - Paginación
  - Empty state con CTA
  - Stats: números vendidos, recaudación, días restantes

- [x] ✅ MyPurchasesPage - Vista de comprador (2025-11-10 18:22)
  - Lista de compras con números adquiridos
  - Status visual (Pendiente, Completado, Cancelado)
  - Resumen de inversión total
  - Empty state para nuevos usuarios
  - Preparado para datos reales cuando se implemente Sprint 5-6

- [x] ✅ RafflesListPage mejorado (2025-11-10 18:24)
  - Search bar prominente con clear button
  - Filtros por estado mejorados (Todos, Activos, Completados)
  - Contador de resultados
  - Paginación mejorada con números de página
  - URL-based search parameters
  - Botón flotante mobile para "Crear Sorteo"
  - EmptyState con CTA

- [x] ✅ RaffleDetailPage con hero section (2025-11-10 18:27)
  - Hero gradient con título, descripción y precio destacado
  - CTA prominente "Comprar Números" (preparado para pagos)
  - Progress bar de ventas
  - Countdown de días restantes
  - Stats grid mejorado (Disponibles, Vendidos, Reservados, Recaudación)
  - Sección de información del sorteo
  - Grid de números visualizado

**Routing y Estructura:**
- [x] ✅ App.tsx actualizado con MainLayout (2025-11-10 18:11)
  - Landing page sin layout (pública)
  - Auth pages sin layout
  - Todas las páginas protegidas con MainLayout
  - Nuevas rutas: /my-raffles, /my-purchases

**Correcciones Técnicas:**
- [x] ✅ Fixed TypeScript errors (2025-11-10 18:28)
  - Corregido import path: @/stores → @/store
  - Añadidos type annotations (raffle: Raffle, n: string)
  - Fixed User type usage: name → first_name + last_name
  - Removed unused variables (isCancelled)
  - Fixed hook import: useRaffles → useRafflesList

**Build y Deployment:**
- [x] ✅ Version bump v1.1.0 en main.tsx (2025-11-10 18:29)
- [x] ✅ Clean build sin errores (2025-11-10 18:29)
  - Bundle: 441.86 kB JS (gzipped: 129.33 kB)
  - TypeScript compilation: 0 errors
- [x] ✅ Docker multi-stage build exitoso (2025-11-10 18:30)
- [x] ✅ Deployed to production https://sorteos.club (2025-11-10 18:30)

#### Archivos Creados/Modificados Sprint 3.5

```
frontend/
├── src/
│   ├── components/
│   │   ├── layout/
│   │   │   ├── Navbar.tsx                             ✅ NEW
│   │   │   ├── UserMenu.tsx                           ✅ NEW
│   │   │   └── MainLayout.tsx                         ✅ NEW
│   │   └── ui/
│   │       ├── StatsCard.tsx                          ✅ NEW
│   │       ├── EmptyState.tsx                         ✅ NEW
│   │       └── LoadingSpinner.tsx                     ✅ NEW
│   ├── features/
│   │   ├── dashboard/pages/
│   │   │   └── DashboardPage.tsx                      ✅ UPDATED (complete redesign)
│   │   └── raffles/pages/
│   │       ├── MyRafflesPage.tsx                      ✅ NEW
│   │       ├── MyPurchasesPage.tsx                    ✅ NEW
│   │       ├── RafflesListPage.tsx                    ✅ UPDATED (improved filters + search)
│   │       └── RaffleDetailPage.tsx                   ✅ UPDATED (hero design + prominent CTA)
│   ├── App.tsx                                        ✅ UPDATED (MainLayout integration)
│   └── main.tsx                                       ✅ UPDATED (v1.1.0)
├── .dockerignore                                      ✅ NEW (root level)
└── package.json                                       ✅ (unchanged)

Total archivos Sprint 3.5:
- Nuevos: 9 archivos (3 layout + 3 UI components + 2 pages + 1 dockerignore)
- Actualizados: 5 archivos (Dashboard, RafflesList, RaffleDetail, App, main)
```

#### Entregables Completados

**Navegación:**
- ✅ Navbar persistente en todas las páginas protegidas
- ✅ User menu con links rápidos (Dashboard, Mis Sorteos, Mis Compras, Logout)
- ✅ Search bar funcional para buscar sorteos
- ✅ Navegación mobile responsive

**Dashboard:**
- ✅ Bienvenida personalizada con nombre del usuario
- ✅ Quick actions con botones grandes y claros
- ✅ Stats cards con iconos y descripciones
- ✅ Sección de actividad reciente (preparada para datos)
- ✅ Información de cuenta visible

**Experiencia Vendedor:**
- ✅ Página "Mis Sorteos" completa con tabla, filtros y stats
- ✅ Vista clara del progreso de cada sorteo
- ✅ Acciones rápidas (Ver Detalles) en cada sorteo
- ✅ Empty state con CTA para crear primer sorteo

**Experiencia Comprador:**
- ✅ Página "Mis Compras" con historial de participaciones
- ✅ Vista de números comprados por sorteo
- ✅ Status visual de cada compra
- ✅ Resumen de inversión total

**Mejoras Generales:**
- ✅ Componentes reutilizables (StatsCard, EmptyState, LoadingSpinner)
- ✅ Consistencia visual en toda la aplicación
- ✅ Dark mode support completo
- ✅ Responsive design mobile-first
- ✅ Empty states informativos con CTAs
- ✅ Loading states consistentes
- ✅ TypeScript sin errores de compilación

#### Impacto

**Antes del Sprint 3.5:**
- Dashboard vacío sin links útiles
- No había forma de ver "mis sorteos" vs "sorteos disponibles"
- Usuario confundido sobre qué hacer después del login
- Falta de navegación clara

**Después del Sprint 3.5:**
- ✅ Navegación clara y persistente
- ✅ Dashboard útil con acciones rápidas
- ✅ Separación clara: Comprador (Mis Compras) vs Vendedor (Mis Sorteos)
- ✅ Search funcional en navbar
- ✅ User experience profesional y pulida
- ✅ Preparado para Sprint 5-6 (Pagos)

#### Decisiones de Diseño

**Opción Elegida:** Complete UX/UI Overhaul (8-10 horas)
- Layout completo con Navbar persistente
- Todas las páginas mejoradas
- Componentes reutilizables
- Sistema de navegación coherente

**Alternativas Descartadas:**
- Quick fixes (4-5 horas): Demasiado limitado
- Mixed approach: Preferible hacer todo de una vez

#### Próximos Pasos

Con la UX/UI mejorada, ahora podemos continuar con:
1. **Sprint 5-6: Reservas y Pagos** - Implementar flujo de compra
2. Integrar stats reales en Dashboard (cuando tengamos datos)
3. Poblar "Mis Compras" con compras reales (después de Sprint 5-6)
4. Implementar upload de imágenes

---

### Sprint 5-6: Reservas y Pagos 🚧 EN PROGRESO

**Fecha inicio:** 2025-11-11 00:00
**Estado Backend:** 100% completado ✅
**Estado Frontend:** 90% completado ✅
**Última actualización:** 2025-11-11 02:30

#### Tareas Backend
- [x] ✅ Migraciones: reservations, payments, idempotency_keys (2025-11-11 00:05)
  - 000006_create_reservations: tabla con TTL (expires_at), array de number_ids, status enum
  - 000007_create_payments: integración Stripe (payment_intent_id, client_secret, metadata JSONB)
  - 000008_create_idempotency_keys: prevención de duplicados con request fingerprint
- [x] ✅ Sistema de reserva temporal (2025-11-11 00:10)
  - Lock distribuido Redis por número (AcquireMultipleLocks atomic)
  - Crear reserva (status=pending, expires_at=now+5min)
  - Cron job para liberar reservas expiradas (cada 1 minuto)
  - Validación de no duplicados con array overlap operator (&&)
- [x] ✅ Integración con PSP - PayPal (2025-11-11 01:15)
  - Interfaz PaymentProvider abstracta
  - Implementación PayPalProvider con Orders API v2
  - Implementación StripeProvider (opcional/legacy)
  - PayPal configurado como provider por defecto
  - Manejo de webhooks (CHECKOUT.ORDER.APPROVED, PAYMENT.CAPTURE.COMPLETED)
  - Soporte sandbox y producción
  - Idempotencia con Idempotency-Key header
- [x] ✅ Domain entities (2025-11-11 00:08)
  - Reservation: métodos IsExpired, CanBePaid, Confirm, Cancel, Expire
  - Payment: métodos MarkAsSucceeded, MarkAsFailed, Cancel, con metadata JSONB
  - IdempotencyKey: validación de request match con SHA-256
- [x] ✅ Repositorios (2025-11-11 00:12)
  - ReservationRepository: 8 métodos incluye CountActiveReservationsForNumbers
  - PaymentRepository: 6 métodos incluye FindByStripePaymentIntentID
  - IdempotencyKeyRepository: 3 métodos para deduplicación
- [x] ✅ Use Cases (2025-11-11 00:17)
  - CreateReservation: con distributed locks + double-check DB + idempotency
  - CreatePaymentIntent: con Stripe integration + metadata tracking
  - ProcessPaymentWebhook: maneja 3 eventos de Stripe
  - ConfirmReservation, CancelReservation, ExpireReservations
  - GetReservation, GetUserReservations, GetPayment, GetUserPayments
- [x] ✅ HTTP Handlers y Rutas (2025-11-11 00:20)
  - POST /api/v1/reservations - Crear reserva con locks
  - GET /api/v1/reservations/:id - Ver reserva
  - GET /api/v1/reservations/me - Mis reservas
  - POST /api/v1/payments/intent - Crear payment intent (Stripe)
  - GET /api/v1/payments/:id - Ver pago
  - GET /api/v1/payments/me - Mis pagos
  - POST /api/v1/webhooks/stripe - Webhook sin auth (Stripe signed)
- [x] ✅ Background Job (2025-11-11 00:18)
  - ExpireReservationsJob: goroutine con ticker cada 1 minuto
  - Integrado en main.go startup
- [x] ✅ Configuración Payment Provider (2025-11-11 01:15)
  - PaymentConfig struct con provider, clientID, secret, sandbox
  - .env.example actualizado con CONFIG_PAYMENT_PROVIDER=paypal
  - Stripe config mantenida como opcional/legacy
- [x] ✅ Build exitoso con PayPal (2025-11-11 01:15)
  - Dependencias: paypal/v4, stripe-go v76, lib/pq
  - Type conversions corregidas
  - User UUID lookup helper implementado
  - Provider dinámico basado en configuración
  - 0 errores de compilación

#### Tareas Frontend
- [x] ✅ Cart Store con Zustand (2025-11-11 02:00)
  - Estado global del carrito con persistencia localStorage
  - Selección multi-número por raffle
  - Gestión de reservas activas
  - Timer de expiración integrado
- [x] ✅ NumberGrid Multi-selección (2025-11-11 02:05)
  - Toggle de números con click
  - Visual feedback de selección
  - Integrado con cart store
  - Readonly mode para owner/inactive raffles
- [x] ✅ RaffleDetailPage actualizada (2025-11-11 02:10)
  - Botón dinámico "Proceder al Pago"
  - Resumen de selección en tiempo real
  - Botón "Limpiar selección"
  - Navegación a checkout
- [x] ✅ Hooks de API (2025-11-11 02:15)
  - useCreateReservation con idempotency
  - useCreatePaymentIntent con PayPal support
  - useGetReservation con polling si pending
  - useGetPayment, useGetMyPayments
- [x] ✅ ReservationTimer Component (2025-11-11 02:18)
  - Countdown de 5 minutos
  - Visual urgente < 1 minuto
  - Callback onExpire
  - Estados: activo, urgente, expirado
- [x] ✅ Página de Checkout (2025-11-11 02:25)
  - Resumen de pedido con números seleccionados
  - Creación de reserva (POST /reservations)
  - Timer de expiración en tiempo real
  - Redirección a PayPal approval URL
  - Estados: review, reserving, reserved, creating_payment, expired
- [x] ✅ PaymentSuccessPage (2025-11-11 02:27)
  - Mensaje de éxito con confetti
  - Detalles de payment_id y reservation_id
  - Limpieza automática del carrito
  - Links a "Mis Compras" y "Ver Sorteos"
- [x] ✅ PaymentCancelPage (2025-11-11 02:28)
  - Mensaje de cancelación
  - Detección de reserva activa
  - Opción de volver al checkout
  - Link a soporte
- [x] ✅ Router actualizado (2025-11-11 02:30)
  - /checkout (protected)
  - /payment/success (protected)
  - /payment/cancel (protected)

#### Entregables
- [ ] Usuario puede reservar números y pagar con tarjeta ⏳
- [ ] Números no se duplican (prueba con 500 req concurrentes) ⏳
- [ ] Reservas expiradas se liberan automáticamente ⏳ (implementado, pendiente testing)
- [ ] Webhooks procesan pagos correctamente ⏳ (implementado, pendiente testing)

#### Archivos Creados Sprint 5-6 (2025-11-11) - BACKEND RESERVAS Y PAGOS ✅

```
backend/
├── migrations/
│   ├── 000006_create_reservations.up.sql              ✅ NEW
│   ├── 000006_create_reservations.down.sql            ✅ NEW
│   ├── 000007_create_payments.up.sql                  ✅ NEW
│   ├── 000007_create_payments.down.sql                ✅ NEW
│   ├── 000008_create_idempotency_keys.up.sql          ✅ NEW
│   └── 000008_create_idempotency_keys.down.sql        ✅ NEW
├── internal/
│   ├── domain/entities/
│   │   ├── reservation.go                             ✅ NEW
│   │   ├── payment.go                                 ✅ NEW
│   │   └── idempotency_key.go                         ✅ NEW
│   ├── domain/repositories/
│   │   ├── reservation_repository.go                  ✅ NEW
│   │   ├── payment_repository.go                      ✅ NEW
│   │   └── idempotency_key_repository.go              ✅ NEW
│   ├── infrastructure/
│   │   ├── database/
│   │   │   ├── postgres_reservation_repository.go     ✅ NEW
│   │   │   ├── postgres_payment_repository.go         ✅ NEW
│   │   │   └── postgres_idempotency_key_repository.go ✅ NEW
│   │   ├── redis/
│   │   │   └── lock_service.go                        ✅ NEW
│   │   └── payment/
│   │       ├── payment_provider.go                    ✅ NEW (interface)
│   │       ├── paypal_provider.go                     ✅ NEW (2025-11-11 01:15)
│   │       └── stripe_provider.go                     ✅ NEW (legacy)
│   ├── adapters/
│   │   ├── db/
│   │   │   ├── reservation_repository.go              ✅ NEW (wrapper)
│   │   │   ├── payment_repository.go                  ✅ NEW (wrapper)
│   │   │   └── idempotency_key_repository.go          ✅ NEW (wrapper)
│   │   └── redis/
│   │       └── lock_service.go                        ✅ NEW (wrapper)
│   ├── usecases/
│   │   ├── reservation_usecases.go                    ✅ NEW
│   │   └── payment_usecases.go                        ✅ NEW
│   ├── jobs/
│   │   └── expire_reservations_job.go                 ✅ NEW
│   └── adapters/http/
│       └── (handlers integrated in cmd/api/)
├── cmd/api/
│   ├── main.go                                        ✅ UPDATED (+startBackgroundJobs call)
│   ├── payment_routes.go                              ✅ NEW (7 endpoints + webhook)
│   └── jobs.go                                        ✅ NEW (background jobs setup)
├── pkg/config/
│   └── config.go                                      ✅ UPDATED (+PaymentConfig)
├── go.mod                                             ✅ UPDATED (+paypal/v4, +stripe-go, +lib/pq)
├── go.sum                                             ✅ UPDATED
├── .env.example                                       ✅ UPDATED (+PayPal config, +Stripe legacy)
└── Dockerfile                                         ✅ UPDATED (+go mod tidy step)
```

**Total archivos Sprint 5-6 Backend:**
- Migraciones: 6 archivos (3 up + 3 down)
- Domain Entities: 3 archivos (Reservation, Payment, IdempotencyKey)
- Repository Interfaces: 3 archivos
- Repository Implementations: 3 archivos
- Adapter Wrappers: 4 archivos (3 repos + 1 lock service)
- Infrastructure Services: 4 archivos (LockService + PaymentProvider + PayPalProvider + StripeProvider)
- Use Cases: 2 archivos (ReservationUseCases, PaymentUseCases)
- HTTP Routes: 1 archivo (payment_routes.go con 7 endpoints)
- Background Jobs: 2 archivos (expire_reservations_job.go, jobs.go)
- Config: 4 archivos actualizados (main.go, config.go, go.mod, .env.example)
- **Subtotal: 29 archivos creados + 5 actualizados**

**Características Backend Implementadas:**
- ✅ Distributed locks con Redis (atomic multi-lock acquisition)
- ✅ Reservas con TTL de 5 minutos
- ✅ Validación de números disponibles con PostgreSQL array overlap
- ✅ PayPal Orders API v2 integration (provider por defecto)
- ✅ Stripe Payment Intents API (opcional/legacy)
- ✅ Payment Provider abstraction (fácil agregar BAC, SINPE Móvil)
- ✅ Webhook signature verification (PayPal y Stripe)
- ✅ Idempotency keys con SHA-256 fingerprinting
- ✅ Background job para expirar reservas (goroutine + ticker)
- ✅ JSONB metadata en payments y idempotency_keys
- ✅ Conversion de User int64 ID → UUID para nuevas entities
- ✅ Helper function getUserUUID en handlers
- ✅ Rate limiting en reservas (cfg.Business.RateLimitReservePerMinute)
- ✅ Rate limiting en pagos (cfg.Business.RateLimitPaymentPerMinute)
- ✅ Audit logging ready (entities tienen user tracking)
- ✅ Configuración dinámica de payment provider (PayPal/Stripe)
- ✅ Soporte sandbox y producción para PayPal

**Endpoints Backend Implementados:**
- POST /api/v1/reservations - Crear reserva con distributed locks
- GET /api/v1/reservations/:id - Ver reserva (owner only)
- GET /api/v1/reservations/me - Listar mis reservas
- POST /api/v1/payments/intent - Crear payment intent (PayPal/Stripe)
- GET /api/v1/payments/:id - Ver pago (owner only)
- GET /api/v1/payments/me - Listar mis pagos
- POST /api/v1/webhooks/stripe - Webhook (PayPal/Stripe, sin auth, firma verificada)

**Flujo Implementado:**
1. Usuario selecciona números → POST /reservations
2. Backend: adquiere locks en Redis + crea reserva (expires_at = now + 5 min)
3. Usuario procede a pago → POST /payments/intent
4. Backend: crea Order/Payment Intent (PayPal/Stripe) → devuelve approval_url/client_secret
5. Frontend: redirige a PayPal o usa Stripe Elements
6. PayPal/Stripe envía webhook → POST /webhooks/stripe
7. Backend: verifica firma → procesa evento:
   - PAYMENT.CAPTURE.COMPLETED / payment_intent.succeeded: Pago exitoso + confirma reserva
   - PAYMENT.CAPTURE.DENIED / payment_intent.payment_failed: Marca pago como failed
   - payment_intent.canceled: Cancela pago + cancela reserva
8. Background job: cada 1 minuto busca reservas expiradas → marca como expired

#### Archivos Creados Sprint 5-6 (2025-11-11) - FRONTEND CHECKOUT CON PAYPAL ✅

```
frontend/
├── src/
│   ├── store/
│   │   └── cartStore.ts                              ✅ NEW (2025-11-11 02:00)
│   ├── hooks/
│   │   ├── useReservations.ts                        ✅ NEW (2025-11-11 02:15)
│   │   └── usePayments.ts                            ✅ NEW (2025-11-11 02:15)
│   ├── components/
│   │   └── ReservationTimer.tsx                      ✅ NEW (2025-11-11 02:18)
│   ├── features/
│   │   ├── checkout/
│   │   │   └── pages/
│   │   │       ├── CheckoutPage.tsx                  ✅ NEW (2025-11-11 02:25)
│   │   │       ├── PaymentSuccessPage.tsx            ✅ NEW (2025-11-11 02:27)
│   │   │       └── PaymentCancelPage.tsx             ✅ NEW (2025-11-11 02:28)
│   │   └── raffles/
│   │       ├── pages/
│   │       │   └── RaffleDetailPage.tsx              ✅ UPDATED (2025-11-11 02:10)
│   │       └── components/
│   │           └── NumberGrid.tsx                    ✅ UPDATED (2025-11-11 02:05)
│   └── App.tsx                                        ✅ UPDATED (2025-11-11 02:30)
```

**Total archivos Sprint 5-6 Frontend:**
- Cart Store: 1 archivo (cartStore.ts con Zustand + persist)
- API Hooks: 2 archivos (useReservations.ts, usePayments.ts con React Query)
- Components: 1 archivo (ReservationTimer.tsx con countdown)
- Checkout Pages: 3 archivos (CheckoutPage, PaymentSuccessPage, PaymentCancelPage)
- Actualizaciones: 3 archivos (RaffleDetailPage, NumberGrid, App.tsx con routes)
- **Subtotal: 7 archivos creados + 3 actualizados**

**Características Frontend Implementadas:**
- ✅ Cart Store con Zustand (persistencia localStorage)
- ✅ Multi-selección de números con toggle
- ✅ Estado global del carrito por raffle
- ✅ Reserva temporal con timer de 5 minutos
- ✅ Checkout flow multi-step (review → reserving → reserved → payment)
- ✅ Integración PayPal redirect flow
- ✅ Countdown timer con estados (normal, urgente, expirado)
- ✅ Payment success page con confetti
- ✅ Payment cancel page con retry option
- ✅ React Query hooks con auto-refetch para pending reservations
- ✅ Protected routes para checkout y payment pages
- ✅ Limpieza automática del carrito post-pago
- ✅ Detección de reserva expirada en checkout

**Rutas Frontend Implementadas:**
- /raffles/:id - Vista de sorteo con NumberGrid + cart integration
- /checkout - Página de checkout protegida (multi-step flow)
- /payment/success - Página de éxito protegida (con confetti + cart cleanup)
- /payment/cancel - Página de cancelación protegida (con retry)

**Flujo Frontend Implementado:**
1. Usuario navega a /raffles/:id
2. Selecciona números → cart store actualiza selectedNumbers
3. Click "Proceder al Pago" → navega a /checkout
4. CheckoutPage: muestra resumen + botón "Confirmar Reserva"
5. Click confirmar → POST /reservations → setReservation en cart store
6. Timer cuenta regresiva desde 5 minutos
7. Click "Pagar con PayPal" → POST /payments/intent → redirect a approval_url
8. Usuario completa pago en PayPal → redirect a /payment/success?payment_id=xxx
9. PaymentSuccessPage: muestra confetti + limpia cart
10. Usuario puede ver "Mis Compras" o volver a sorteos

**Issues Resueltos (2025-11-11 02:00):**
- ✅ Fixed: go.sum faltaba entradas para lib/pq y stripe-go → Ejecutado go mod tidy
- ✅ Fixed: TypeScript error en CheckoutPage - enabled no existe en useRaffleDetail options
- ✅ Fixed: TypeScript error - Reservation type mismatch (camelCase vs snake_case)
- ✅ Fixed: Missing apiClient module → Creado src/lib/apiClient.ts como re-export
- ✅ Fixed: refetchInterval callback accediendo a data en lugar de query.state.data
- ✅ Build exitoso: Docker image construido sin errores (frontend + backend)

**Issues Resueltos (2025-11-11 04:10) - Testing Phase:**
- ✅ Fixed: 403 error al publicar sorteo → Validación de imágenes temporalmente deshabilitada
  - **Archivo:** `backend/internal/usecase/raffle/publish_raffle.go` (líneas 68-89)
  - **Razón:** Upload de imágenes no implementado aún (Sprint 4 pendiente)
  - **Solución temporal:** Comentadas validaciones de imágenes (imageCount y primaryImage)
  - **TODO:** Re-habilitar validaciones cuando se implemente upload de imágenes
  - **Impacto:** Permite publicar sorteos para testing sin necesidad de imágenes
  - **Nota:** Esto es un **quick fix temporal** para permitir testing E2E del flujo de pagos
  - **Ver:** Sprint 4 en roadmap - "Implementar upload de imágenes" debe completarse antes de producción

**Testing Documentation Created (2025-11-11 02:10):**
- ✅ [TESTING-QUICKSTART.md](./TESTING-QUICKSTART.md) - Guía rápida para empezar (30 min)
- ✅ [testing-strategy.md](./testing-strategy.md) - Estrategia completa de testing (3 niveles)
- ✅ [testing-manual-checklist.md](./testing-manual-checklist.md) - 30 test cases manuales
- ✅ [testing-api-scripts.md](./testing-api-scripts.md) - Scripts cURL para API testing
- ✅ [docker-compose.test.yml](../docker-compose.test.yml) - Entorno de test aislado

**Próximos Pasos:**
1. ✅ **Actualizar roadmap** (esta actualización - 2025-11-11 02:30)
2. ✅ **Correr migraciones** en desarrollo (completado 2025-11-11 00:05)
3. ✅ **Integrar PayPal** como provider por defecto (completado 2025-11-11 01:15)
4. ✅ **Implementar frontend** (NumberGrid multi-select, checkout, PayPal button - completado 2025-11-11 02:30)
5. ✅ **Build Docker image** (completado 2025-11-11 02:00 - frontend + backend sin errores)
6. ✅ **Crear documentación de testing** (completado 2025-11-11 02:10 - 5 archivos)
7. ⏳ **Ejecutar testing manual** (~30 min - usar checklist)
8. ⏳ **Validar con PayPal sandbox credentials** (configurar en .env)
9. ⏳ **Ejecutar testing de API** (~1-2 horas - scripts cURL)
10. ⏳ **Testing de concurrencia** (100 requests simultáneas con script bash)

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
