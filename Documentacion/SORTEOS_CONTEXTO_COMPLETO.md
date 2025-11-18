# Contexto Técnico Completo - Plataforma de Sorteos

**Proyecto:** Sistema de Sorteos/Rifas en Línea
**Propietario:** Ing. Alonso Alpízar
**Fecha:** 2025-11-18
**Versión:** 2.0
**Estado:** MVP en desarrollo - Sistema de autenticación y gestión de sorteos implementados

---

## 📋 ÍNDICE

1. [Stack Tecnológico](#1-stack-tecnológico)
2. [Arquitectura Principal](#2-arquitectura-principal)
3. [Decisiones Técnicas Importantes](#3-decisiones-técnicas-importantes)
4. [Contexto de Negocio](#4-contexto-de-negocio)
5. [Estado Actual del Desarrollo](#5-estado-actual-del-desarrollo)

---

## 1. STACK TECNOLÓGICO

### 1.1 Backend

#### Lenguaje y Runtime
- **Go 1.22+** - Lenguaje principal
  - **¿Por qué Go?**
    - Rendimiento nativo comparable a C/C++
    - Concurrencia nativa con goroutines y channels (crítico para reservas simultáneas)
    - Compilación estática (binario sin dependencias)
    - Gestión de memoria eficiente con GC optimizado
    - Ideal para APIs de alto tráfico con transacciones críticas
    - Ecosistema maduro para fintech y e-commerce

#### Framework Web
- **Gin (gin-gonic/gin) v1.9.1+**
  - Router extremadamente rápido (httprouter bajo el capó)
  - Middlewares composables
  - Validación integrada con binding
  - Soporte para JSON, XML, YAML
  - Comunidad activa y amplia documentación

#### Dependencias Core

**Autenticación y Seguridad:**
- `golang-jwt/jwt/v5` (v5.2.0+) - Generación y validación de JWT
- `bcrypt` - Hashing de contraseñas (cost 12)

**ORM / Query Builder:**
- **GORM (gorm.io/gorm) v1.25.0+** - Para MVP (velocidad de desarrollo)
- **sqlc** (futuro) - Para módulos críticos si se requiere optimización

**Driver PostgreSQL:**
- `jackc/pgx/v5` (v5.5.0+) - Driver de alto rendimiento
- Soporte para tipos nativos (JSONB, UUID, arrays)
- Connection pooling eficiente

**Cliente Redis:**
- `redis/go-redis/v9` (v9.5.0+)
- Casos de uso:
  - Locks distribuidos para reservas
  - Caché de sorteos activos
  - Rate limiting por IP/usuario
  - Sesiones y refresh tokens
  - Idempotencia de pagos

**Logging:**
- `uber-go/zap` (v1.27.0+) - Logging estructurado de alto rendimiento
- Niveles: Debug, Info, Warn, Error, Fatal
- Campos tipados (evita allocations)

**Configuración:**
- `spf13/viper` (v1.18.0+)
- Lectura de `.env`, YAML, JSON, TOML
- Variables de entorno con prefijos

**Validación:**
- `go-playground/validator/v10` (v10.19.0+)
- Validación de structs con tags
- Reglas personalizadas

**Migraciones:**
- `golang-migrate/migrate` (v4.17+)
- Migraciones SQL versionadas

#### Archivos Go Implementados
- **Total:** 117 archivos .go
- **Estructura:**
  - `cmd/api/` - Entry point, routes, jobs
  - `internal/domain/` - Entidades y reglas de negocio
  - `internal/usecase/` - Casos de uso (auth, raffle, admin, image)
  - `internal/adapters/` - HTTP handlers, DB repositories, payments, notifiers
  - `pkg/` - Utilidades compartidas (logger, config, errors, crypto)

### 1.2 Frontend

#### Lenguaje y Runtime
- **TypeScript 5.3+** - Type safety en desarrollo
- **Node.js 20 LTS+** - Runtime para desarrollo y build

#### Build Tool
- **Vite 5.0+**
  - HMR instantáneo
  - Build optimizado (Rollup)
  - Soporte nativo para TypeScript, JSX, CSS Modules
  - Build de producción: ~10 segundos

#### Framework UI
- **React 18.2+**
  - Concurrent rendering
  - Suspense y Error Boundaries
  - Hooks modernos

#### Librerías Core

**Routing:**
- `react-router-dom` (v6.22+)

**State Management:**
- `zustand` (v4.5+) - Minimal boilerplate, TypeScript first
- Uso: Estado de autenticación, carrito de compra

**Data Fetching:**
- `@tanstack/react-query` (v5.0+)
  - Caché automático
  - Refetch automático
  - Optimistic updates
  - Paginación

**HTTP Client:**
- `axios` (v1.6+)
  - Interceptores para auth (JWT)
  - Timeout, retry
  - Type-safe con generics de TS

**Forms:**
- `react-hook-form` (v7.50+) + `zod` (v3.22+)
  - Validación declarativa
  - TypeScript first

**UI Components:**
- **Tailwind CSS 3.4+** - Utility-first CSS
- **shadcn/ui** - Componentes accesibles basados en Radix UI
  - Button, Input, Select, Card, Table
  - Dialog, Toast, Badge, Skeleton
  - Form (integrado con react-hook-form)

**Internacionalización (Fase 2):**
- `i18next` (v23+) + `react-i18next` (v14+)

#### Archivos TypeScript Implementados
- **Total:** 67 archivos .ts/.tsx
- **Estructura:**
  - `app/` - Router y providers
  - `features/` - Módulos (auth, raffles, dashboard)
  - `components/ui/` - Componentes shadcn/ui
  - `components/` - Componentes de negocio (NumberGrid, ImageUploader)
  - `lib/` - Utilidades, API client
  - `store/` - Zustand stores
  - `api/` - Clientes API tipados

### 1.3 Base de Datos

#### Motor Principal
- **PostgreSQL 16** (instalación nativa local)
  - **¿Por qué PostgreSQL?**
    - ACID compliant (crítico para transacciones de pago)
    - Índices avanzados (B-tree, GIN, GiST)
    - JSONB para datos semi-estructurados
    - Transacciones con niveles de aislamiento configurables
    - Vistas materializadas para KPIs
    - Extensions: `uuid-ossp`, `pg_trgm` (búsqueda fuzzy)
  - **Puerto:** 5432
  - **Base de datos:** sorteos_db
  - **Usuario:** sorteos_user

#### Cache y Concurrencia
- **Redis 7.2+** (instalación nativa local)
  - **Modos:**
    - Standalone (desarrollo/staging)
    - Sentinel (producción con HA - futuro)
  - **Puerto:** 6379
  - **Casos de uso críticos:**
    - Locks distribuidos (SETNX) para reservas de números
    - Caché de sorteos activos (TTL: 5-10 min)
    - Rate limiting (Token Bucket)
    - Idempotencia de pagos (24h TTL)

### 1.4 Infraestructura

#### Servidor Web
- **Nginx** - Reverse proxy + SSL
  - Proxy: `https://sorteos.club` → `localhost:8080`
  - SSL/TLS con Let's Encrypt
  - Servir archivos estáticos (delegado al backend Go)

#### Gestión de Servicios
- **systemd** - Todos los servicios gestionados nativamente
  - `postgresql.service` - Base de datos
  - `redis-server.service` - Cache y locks
  - `sorteos-api.service` - Backend Go
  - `nginx.service` - Reverse proxy

#### Pagos
- **Stripe** (MVP - Fase 1)
  - Payment Intents API
  - Webhooks con verificación de firma
  - Tokens para tarjetas (PCI DSS delegado)
- **PayPal** (Fase 2)
- **Procesador local Costa Rica** (Fase 2)

#### Notificaciones
- **SMTP Propio** (sorteos.club)
  - Dovecot + Postfix
  - DKIM, SPF, DMARC configurados
  - Plantillas de email transaccionales
- **SendGrid** (futuro - emails masivos)
- **Twilio** (futuro - SMS)

### 1.5 Migración Reciente: Docker → Local

**Estado anterior (Docker):**
- 6 paquetes Docker + dependencias (464 MB overhead)
- Rebuild frontend: 3+ minutos
- Debugging complejo (logs en contenedores)

**Estado actual (Nativo):**
- PostgreSQL 16 instalado nativamente
- Redis 7 instalado nativamente
- Backend Go como servicio systemd
- Frontend servido por backend desde `dist/`
- Rebuild frontend: **10 segundos**
- Logs centralizados en journalctl
- Stack nativo, rápido y mantenible

---

## 2. ARQUITECTURA PRINCIPAL

### 2.1 Estructura de Directorios

```
/opt/Sorteos/
├── backend/                          # API en Go
│   ├── cmd/
│   │   └── api/                      # Entry point
│   │       ├── main.go               # Inicialización del servidor
│   │       ├── routes.go             # Definición de rutas
│   │       ├── payment_routes.go     # Rutas de pagos
│   │       └── jobs.go               # Cron jobs (limpieza reservas)
│   ├── internal/
│   │   ├── domain/                   # Entidades y reglas de negocio
│   │   │   ├── user.go
│   │   │   ├── raffle.go
│   │   │   ├── reservation.go
│   │   │   └── payment.go
│   │   ├── usecase/                  # Casos de uso (lógica de aplicación)
│   │   │   ├── auth/                 # Registro, login, verificación
│   │   │   ├── raffle/               # Crear, publicar, listar sorteos
│   │   │   ├── category/             # Gestión de categorías
│   │   │   ├── image/                # Subida de imágenes
│   │   │   └── admin/                # Operaciones administrativas
│   │   └── adapters/                 # Adaptadores externos
│   │       ├── http/                 # Handlers Gin
│   │       ├── db/                   # Repositorios GORM
│   │       ├── redis/                # Cliente Redis
│   │       ├── payments/             # Providers (Stripe, PayPal)
│   │       └── notifier/             # Email, SMS
│   ├── pkg/                          # Librerías compartidas
│   │   ├── config/                   # Configuración (Viper)
│   │   ├── logger/                   # Logger (Zap)
│   │   ├── errors/                   # Errores personalizados
│   │   └── crypto/                   # Password hashing, códigos
│   ├── migrations/                   # Migraciones SQL
│   │   ├── 001_create_users_table.up.sql
│   │   ├── 002_create_raffles_table.up.sql
│   │   └── ... (migraciones versionadas)
│   ├── uploads/                      # Archivos subidos (imágenes)
│   ├── .env                          # Variables de entorno
│   ├── Makefile                      # Comandos útiles
│   ├── go.mod                        # Dependencias Go
│   └── sorteos-api                   # Binario compilado
├── frontend/                         # SPA en React + TypeScript
│   ├── src/
│   │   ├── app/                      # Router y providers
│   │   │   ├── App.tsx
│   │   │   └── router.tsx
│   │   ├── features/                 # Módulos por funcionalidad
│   │   │   ├── auth/                 # Login, registro, verificación
│   │   │   │   ├── components/
│   │   │   │   ├── pages/
│   │   │   │   └── api/
│   │   │   ├── raffles/              # Listado, detalle, creación
│   │   │   │   ├── components/
│   │   │   │   └── pages/
│   │   │   └── dashboard/            # Dashboard usuario/admin
│   │   ├── components/               # Componentes compartidos
│   │   │   ├── ui/                   # shadcn/ui (Button, Card, etc.)
│   │   │   ├── layout/               # Navbar, Footer
│   │   │   ├── NumberGrid.tsx        # Grid de números de sorteo
│   │   │   ├── ImageUploader.tsx     # Subida de imágenes
│   │   │   └── ReservationTimer.tsx  # Timer de reserva
│   │   ├── lib/                      # Utilidades
│   │   │   ├── api.ts                # Cliente Axios
│   │   │   ├── queryClient.ts        # React Query config
│   │   │   └── utils.ts              # Helpers
│   │   ├── store/                    # Zustand stores
│   │   │   ├── authStore.ts          # Estado de autenticación
│   │   │   └── cartStore.ts          # Carrito de números
│   │   ├── types/                    # Definiciones TypeScript
│   │   └── main.tsx                  # Entry point
│   ├── dist/                         # Build de producción (servido por backend)
│   ├── public/                       # Assets estáticos
│   ├── package.json
│   ├── vite.config.ts
│   └── tsconfig.json
└── Documentacion/                    # Documentación técnica (10 docs)
    ├── arquitecturaIdeaGeneral.md
    ├── stack_tecnico.md
    ├── roadmap.md
    ├── modulos.md
    ├── estandar_visual.md
    ├── seguridad.md
    ├── pagos_integraciones.md
    ├── parametrizacion_reglas.md
    ├── operacion_backoffice.md
    └── CLAUDE.md                     # Contexto rápido para AI
```

### 2.2 Componentes Principales del Sistema

#### Backend (Arquitectura Hexagonal)

**Capa 1: Domain (Núcleo de negocio)**
- Entidades puras sin dependencias externas
- Reglas de negocio invariantes
- Interfaces que definen contratos
- Ejemplo: `User`, `Raffle`, `Reservation`, `Payment`

**Capa 2: Use Cases (Lógica de aplicación)**
- Orquestación de reglas de negocio
- Implementación de casos de uso
- Ejemplo: `RegisterUser`, `CreateRaffle`, `ReserveNumbers`, `ProcessPayment`

**Capa 3: Adapters (Implementaciones técnicas)**
- **Driving (entradas):** HTTP handlers (Gin)
- **Driven (salidas):**
  - Repositories (GORM/PostgreSQL)
  - Payment Providers (Stripe, PayPal)
  - Notifiers (Email SMTP, SMS)
  - Cache (Redis)

**Beneficios de esta arquitectura:**
- Testabilidad: Cada capa es testeable aisladamente
- Independencia: Cambiar base de datos no afecta el dominio
- Extensibilidad: Agregar nuevos PSPs sin tocar use cases
- Mantenibilidad: Separación clara de responsabilidades

#### Frontend (Feature-based)

**Estructura por features:**
- Cada módulo (`auth`, `raffles`, `dashboard`) contiene:
  - Components: Componentes específicos del módulo
  - Pages: Páginas completas
  - API: Clientes API tipados
  - Hooks: Custom hooks del módulo
  - Types: Tipos TypeScript específicos

**Componentes UI compartidos:**
- `components/ui/`: shadcn/ui components (Button, Card, Input, etc.)
- `components/layout/`: Layout components (Navbar, Footer)
- `components/`: Business components (NumberGrid, ImageUploader)

**Estado:**
- **Local:** useState, useReducer
- **Cliente (Global):** Zustand stores
- **Servidor:** React Query (TanStack Query)

### 2.3 Separación Backend/Frontend

**Similitud con DIV:** SÍ, hay separación total

**Backend:**
- API RESTful en Go (puerto 8080)
- Endpoints: `/api/v1/*`
- Autenticación: JWT en header Authorization
- Servir frontend desde `/frontend/dist/`

**Frontend:**
- SPA en React + TypeScript
- Build con Vite → archivos estáticos en `dist/`
- Comunicación con backend vía Axios
- Proxy de desarrollo (Vite) para `/api` → `localhost:8080`

**Flujo de deployment:**
1. Build frontend: `npm run build` → `dist/`
2. Backend Go sirve archivos estáticos desde `dist/`
3. Nginx proxy reverso: `https://sorteos.club` → `localhost:8080`
4. Backend maneja tanto API como serving del frontend

---

## 3. DECISIONES TÉCNICAS IMPORTANTES

### 3.1 Patrones de Diseño Específicos

#### 1. Hexagonal Architecture (Ports & Adapters)

**Aplicación:**
```go
// Domain: Interfaz (Port)
type PaymentProvider interface {
    Authorize(ctx context.Context, input AuthorizeInput) (*AuthorizeOutput, error)
    Capture(ctx context.Context, paymentID string) error
    Refund(ctx context.Context, paymentID string, amount decimal.Decimal) error
}

// Adapters: Implementaciones
type StripeProvider struct { ... }
type PayPalProvider struct { ... }
type LocalCRProvider struct { ... }

// Use Case depende de la interfaz, no de la implementación
type ProcessPaymentUseCase struct {
    provider PaymentProvider  // Inyección de dependencia
}
```

**Beneficio:** Cambiar de Stripe a PayPal no requiere modificar use cases

#### 2. Repository Pattern

**Aplicación:**
```go
// Domain: Interfaz
type UserRepository interface {
    Create(ctx context.Context, user *User) error
    FindByEmail(ctx context.Context, email string) (*User, error)
    Update(ctx context.Context, user *User) error
}

// Adapter: Implementación GORM
type PostgresUserRepository struct {
    db *gorm.DB
}

// Tests: Mock implementation
type MockUserRepository struct {
    users map[int64]*User
}
```

#### 3. Factory Pattern

**Aplicación:**
```go
func NewPaymentProvider(providerType string, config Config) PaymentProvider {
    switch providerType {
    case "stripe":
        return NewStripeProvider(config.Stripe)
    case "paypal":
        return NewPayPalProvider(config.PayPal)
    default:
        return NewMockProvider()
    }
}
```

#### 4. Strategy Pattern

**Aplicación:** Diferentes fuentes de lotería
```go
type LotterySource interface {
    GetResult(date string) (*LotteryResult, error)
}

type LoteriaNacionalCR struct { ... }
type ManualDraw struct { ... }
```

### 3.2 Convenciones de Naming

#### Backend (Go)

**Archivos:**
- Snake case: `user_repository.go`, `create_raffle.go`
- Test files: `*_test.go`

**Paquetes:**
- Lowercase, singular: `package user`, `package raffle`

**Structs y Types:**
- PascalCase: `type User struct`, `type RaffleStatus string`
- Exportados: Primera letra mayúscula
- Privados: Primera letra minúscula

**Funciones y métodos:**
- PascalCase exportados: `func CreateUser(...)`
- camelCase privados: `func validateEmail(...)`

**Constantes:**
- PascalCase: `const MaxReservationTime = 5 * time.Minute`

**Ejemplo completo:**
```go
// internal/usecase/raffle/create_raffle.go
package raffle

type CreateRaffleUseCase struct {
    raffleRepo   domain.RaffleRepository
    imageStorage domain.ImageStorage
    logger       *zap.Logger
}

func NewCreateRaffleUseCase(deps Dependencies) *CreateRaffleUseCase {
    return &CreateRaffleUseCase{
        raffleRepo:   deps.RaffleRepo,
        imageStorage: deps.ImageStorage,
        logger:       deps.Logger,
    }
}

func (uc *CreateRaffleUseCase) Execute(ctx context.Context, input CreateRaffleInput) (*domain.Raffle, error) {
    // Lógica del caso de uso
}
```

#### Frontend (TypeScript/React)

**Archivos:**
- PascalCase para componentes: `LoginPage.tsx`, `NumberGrid.tsx`
- camelCase para utilidades: `utils.ts`, `apiClient.ts`

**Componentes:**
- PascalCase: `function LoginPage() { ... }`

**Hooks:**
- Prefijo `use`: `useAuth()`, `useRaffles()`

**Types/Interfaces:**
- PascalCase: `interface User { ... }`, `type RaffleStatus = '...'`

**Constantes:**
- SCREAMING_SNAKE_CASE: `const API_BASE_URL = '...'`

**Ejemplo completo:**
```typescript
// features/auth/pages/LoginPage.tsx
import { useAuth } from '@/hooks/useAuth'
import { Button } from '@/components/ui/Button'

interface LoginFormData {
  email: string
  password: string
}

export function LoginPage() {
  const { login, isLoading } = useAuth()

  const handleSubmit = async (data: LoginFormData) => {
    await login(data)
  }

  return <form>...</form>
}
```

### 3.3 Reglas de Validación

#### Backend (Go)

**Validación con tags:**
```go
type CreateRaffleRequest struct {
    Title       string          `json:"title" validate:"required,min=5,max=200"`
    Description string          `json:"description" validate:"required,min=20,max=2000"`
    DrawDate    time.Time       `json:"draw_date" validate:"required,future"`
    Price       decimal.Decimal `json:"price" validate:"required,gt=0,lte=10000"`
}
```

**Validaciones personalizadas:**
```go
// Validar que DrawDate sea futuro
func validateFutureDate(fl validator.FieldLevel) bool {
    date := fl.Field().Interface().(time.Time)
    return date.After(time.Now())
}

validate.RegisterValidation("future", validateFutureDate)
```

#### Frontend (TypeScript)

**Validación con Zod:**
```typescript
import { z } from 'zod'

const registerSchema = z.object({
  email: z.string().email('Email inválido'),
  password: z.string()
    .min(12, 'Mínimo 12 caracteres')
    .regex(/[A-Z]/, 'Debe contener mayúscula')
    .regex(/[a-z]/, 'Debe contener minúscula')
    .regex(/[0-9]/, 'Debe contener número')
    .regex(/[^A-Za-z0-9]/, 'Debe contener símbolo'),
  phone: z.string().regex(/^\+\d{10,15}$/, 'Formato E.164'),
})

type RegisterFormData = z.infer<typeof registerSchema>
```

**Validaciones críticas:**
- Email: Formato válido + único en sistema
- Password: Mínimo 12 chars, mayúscula, minúscula, número, símbolo
- Teléfono: Formato E.164 (+573001234567)
- Cédula: 7-10 dígitos solo números (Costa Rica)

### 3.4 Manejo de Errores

#### Backend (Go)

**Tipos de errores:**
```go
// pkg/errors/errors.go
var (
    ErrNotFound           = errors.New("resource not found")
    ErrUnauthorized       = errors.New("unauthorized")
    ErrForbidden          = errors.New("forbidden")
    ErrBadRequest         = errors.New("bad request")
    ErrInternalServer     = errors.New("internal server error")
    ErrConflict           = errors.New("conflict")
    ErrNumberAlreadyReserved = errors.New("number already reserved")
)

// Errores con contexto
type AppError struct {
    Err     error
    Code    int
    Message string
    Details map[string]interface{}
}

func (e *AppError) Error() string {
    return e.Message
}
```

**Manejo en handlers:**
```go
func (h *RaffleHandler) CreateRaffle(c *gin.Context) {
    raffle, err := h.useCase.Execute(c.Request.Context(), input)
    if err != nil {
        switch {
        case errors.Is(err, ErrNotFound):
            c.JSON(404, gin.H{"error": err.Error()})
        case errors.Is(err, ErrBadRequest):
            c.JSON(400, gin.H{"error": err.Error()})
        case errors.Is(err, ErrUnauthorized):
            c.JSON(401, gin.H{"error": err.Error()})
        default:
            logger.Error("unexpected error", zap.Error(err))
            c.JSON(500, gin.H{"error": "internal server error"})
        }
        return
    }

    c.JSON(201, raffle)
}
```

**Logging estructurado:**
```go
logger.Error("failed to create raffle",
    zap.Error(err),
    zap.Int64("user_id", userID),
    zap.String("title", input.Title),
    zap.String("trace_id", traceID),
)
```

#### Frontend (TypeScript)

**Manejo con React Query:**
```typescript
const { mutate: createRaffle, error, isError } = useMutation({
  mutationFn: (data: CreateRaffleData) => api.createRaffle(data),
  onError: (error: AxiosError<ApiError>) => {
    if (error.response?.status === 400) {
      toast.error(error.response.data.message)
    } else if (error.response?.status === 401) {
      toast.error('Debes iniciar sesión')
      navigate('/login')
    } else {
      toast.error('Error inesperado. Intenta de nuevo.')
    }
  },
  onSuccess: (raffle) => {
    toast.success('Sorteo creado exitosamente')
    navigate(`/raffles/${raffle.id}`)
  }
})
```

**Interceptor de Axios:**
```typescript
api.interceptors.response.use(
  response => response,
  error => {
    if (error.response?.status === 401) {
      // Token expirado, intentar refresh
      return refreshTokenAndRetry(error.config)
    }

    if (error.response?.status === 429) {
      toast.error('Demasiadas solicitudes. Espera un momento.')
    }

    return Promise.reject(error)
  }
)
```

### 3.5 Seguridad

#### Autenticación JWT

**Access Token (15 minutos):**
```go
func GenerateAccessToken(userID int64, role string, kycLevel string) (string, error) {
    claims := jwt.MapClaims{
        "user_id":   userID,
        "role":      role,
        "kyc_level": kycLevel,
        "exp":       time.Now().Add(15 * time.Minute).Unix(),
        "iat":       time.Now().Unix(),
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(config.JWTSecret))
}
```

**Refresh Token (7 días):**
- Almacenado en Redis con TTL
- Rotación obligatoria al usar (invalida anterior)
- Revocable por `jti` (JWT ID único)

**Middleware de autenticación:**
```go
func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.JSON(401, gin.H{"error": "missing authorization header"})
            c.Abort()
            return
        }

        tokenString := strings.TrimPrefix(authHeader, "Bearer ")
        claims, err := ValidateToken(tokenString)
        if err != nil {
            c.JSON(401, gin.H{"error": "invalid token"})
            c.Abort()
            return
        }

        c.Set("user_id", claims.UserID)
        c.Set("role", claims.Role)
        c.Next()
    }
}
```

#### Rate Limiting

**Implementación con Redis:**
```go
func RateLimitMiddleware(maxRequests int, window time.Duration) gin.HandlerFunc {
    return func(c *gin.Context) {
        userID := c.GetInt64("user_id")
        key := fmt.Sprintf("ratelimit:%s:%d", c.Request.URL.Path, userID)

        count, _ := rdb.Incr(ctx, key).Result()
        if count == 1 {
            rdb.Expire(ctx, key, window)
        }

        if count > int64(maxRequests) {
            c.JSON(429, gin.H{"error": "too many requests"})
            c.Abort()
            return
        }

        c.Next()
    }
}
```

**Límites por endpoint:**
- `POST /auth/login`: 5 req/min por IP (prevenir brute force)
- `POST /auth/register`: 3 req/hora por IP (prevenir spam)
- `POST /raffles/:id/reservations`: 10 req/min por user_id
- `POST /payments`: 5 req/min por user_id
- `GET /raffles`: 60 req/min por IP

#### Prevención OWASP Top 10

**SQL Injection:**
- GORM escapa automáticamente parámetros
- Nunca concatenar strings en queries
```go
// ✅ Correcto
db.Where("email = ?", email).First(&user)

// ❌ Incorrecto
db.Raw(fmt.Sprintf("SELECT * FROM users WHERE email = '%s'", email))
```

**XSS:**
- Frontend: React escapa automáticamente
- Backend: Validación de HTML en inputs
```go
func SanitizeHTML(input string) string {
    p := bluemonday.StrictPolicy()
    return p.Sanitize(input)
}
```

**CSRF:**
- SPA sin cookies de sesión (JWT en header)
- Estado en memoria o localStorage
- Si se usan cookies: SameSite=Strict

**Mass Assignment:**
- Binding selectivo en Go
```go
type UpdateUserRequest struct {
    FirstName string `json:"first_name"`
    LastName  string `json:"last_name"`
    // NO incluir role, kyc_level (campos sensibles)
}
```

---

## 4. CONTEXTO DE NEGOCIO

### 4.1 Flujos Principales

#### 1. Registrarse y Verificar Email

**Actores:** Usuario nuevo

**Flujo:**
1. Usuario completa formulario de registro
   - Nombre, apellido, email, teléfono, contraseña
   - Acepta términos y condiciones (GDPR)
2. Backend valida datos y crea usuario (status=active, kyc_level=none)
3. Backend genera código de verificación de 6 dígitos
4. Backend envía email con código (expira en 15 minutos)
5. Usuario ingresa código en frontend
6. Backend valida código y actualiza kyc_level=email_verified
7. Backend genera access_token + refresh_token
8. Frontend guarda tokens y redirige a dashboard

**Reglas de negocio:**
- Email debe ser único en el sistema
- Contraseña: mínimo 12 caracteres, mayúscula, minúscula, número, símbolo
- Código de verificación: 6 dígitos, válido por 15 minutos
- Sin verificación de email, no se puede crear sorteos ni comprar boletos

**Estados:**
```
Usuario nuevo → Registrado (kyc_level=none) → Email verificado (kyc_level=email_verified)
```

#### 2. Crear Sorteo

**Actores:** Usuario con kyc_level >= email_verified

**Flujo:**
1. Usuario completa formulario de creación:
   - Título, descripción
   - Categoría (electrónica, vehículos, etc.)
   - Fecha de sorteo
   - Fuente de lotería (Lotería Nacional CR, manual)
   - Rango de números (ej: 00-99)
   - Precio por número
   - Imágenes (mínimo 1, máximo 5)
2. Backend valida parámetros (ver `parametrizacion_reglas.md`)
3. Backend crea sorteo en estado `draft`
4. Backend genera números disponibles (ej: 00, 01, ..., 99)
5. Backend sube imágenes a storage (filesystem local o S3 futuro)
6. Backend retorna sorteo_id
7. Usuario puede:
   - Publicar sorteo (cambia a `active`)
   - Editar draft
   - Eliminar draft

**Reglas de negocio:**
- DrawDate debe ser futuro (mínimo 24 horas)
- PricePerNumber: mínimo ₡100, máximo ₡10,000
- Máximo 10 sorteos activos por usuario (parámetro configurable)
- Imágenes: formatos JPG/PNG, tamaño máximo 2 MB cada una

**Estados:**
```
draft → active → (suspended) → completed/cancelled
```

#### 3. Comprar Boleto (Reservar y Pagar)

**Actores:** Usuario con kyc_level >= email_verified

**Flujo (crítico - alta concurrencia):**

**Fase 1: Reserva (5 minutos)**
1. Usuario ve sorteo activo y grid de números
2. Usuario selecciona números disponibles (máximo 10)
3. Frontend genera UUID como idempotency_key
4. Frontend POST `/raffles/:id/reservations`
5. Backend ejecuta lógica de concurrencia:
   ```
   a. Validar raffle.status == active
   b. Verificar idempotencia (si existe, retornar reserva anterior)
   c. Adquirir locks distribuidos en Redis (SETNX, TTL=30s):
      - lock:raffle:123:num:01
      - lock:raffle:123:num:15
   d. Si algún lock falla → liberar todos → error 409 "número ya reservado"
   e. Si todos los locks OK → crear reserva en DB (transacción):
      - INSERT INTO reservations (status=pending, expires_at=now+5min)
      - UPDATE raffle_numbers SET status=reserved, user_id=X
   f. Liberar locks
   g. Guardar reserva en Redis (TTL=5min)
   ```
6. Frontend recibe reservation_id y muestra timer de 5 minutos
7. Frontend redirige a checkout

**Fase 2: Pago (Stripe)**
8. Frontend muestra formulario de pago (Stripe Elements)
9. Usuario ingresa datos de tarjeta
10. Frontend tokeniza tarjeta con Stripe.js (no envía a backend)
11. Frontend POST `/payments` con:
    - reservation_id
    - payment_method_id (token de Stripe)
    - idempotency_key (mismo UUID)
12. Backend:
    ```
    a. Verificar idempotencia en Redis (24h TTL)
    b. Si existe payment_id → retornar pago anterior (200 OK)
    c. Crear PaymentIntent en Stripe:
       - amount = reservation.numbers.length * raffle.price_per_number
       - metadata: { reservation_id, user_id }
    d. Si pago requiere acción (3D Secure):
       - Retornar action_url
       - Frontend redirige a Stripe
    e. Si pago exitoso inmediatamente:
       - Webhook de Stripe confirma (async)
       - O verificar status en backend
    ```
13. Webhook de Stripe llega a `/webhooks/stripe`:
    ```
    a. Verificar firma del webhook
    b. Extraer payment_intent.id y metadata
    c. Buscar reservation_id en metadata
    d. Transacción:
       - UPDATE payments SET status=succeeded
       - UPDATE reservations SET status=confirmed
       - UPDATE raffle_numbers SET status=sold, sold_at=now
    e. Enviar email de confirmación al usuario
    ```
14. Frontend polling cada 2s para verificar pago confirmado
15. Al confirmar → mostrar comprobante con números comprados

**Fase 3: Limpieza automática (Cron job cada 1 minuto)**
```
a. Buscar reservas con status=pending y expires_at < now
b. Para cada reserva expirada:
   - UPDATE reservations SET status=expired
   - UPDATE raffle_numbers SET status=available, user_id=NULL
```

**Reglas de negocio críticas:**
- Locks distribuidos obligatorios para prevenir doble venta
- Reserva expira exactamente a los 5 minutos
- Idempotencia en reservas y pagos (mismo UUID → mismo resultado)
- Números solo cambian a `sold` cuando pago está confirmado
- Si pago falla → liberar números automáticamente

**Prevención de problemas:**
```
Problema: 2 usuarios clickean el mismo número simultáneamente
Solución: Lock distribuido en Redis (SETNX) - solo uno adquiere el lock

Problema: Usuario paga dos veces por error (doble click)
Solución: Idempotency-Key en Redis (24h TTL) - retorna pago anterior

Problema: Webhook de Stripe llega tarde (después de 5 min)
Solución: Verificar si reserva ya expiró - si sí, hacer refund automático

Problema: Backend crashea mientras tiene locks
Solución: Locks con TTL de 30s - se liberan automáticamente
```

#### 4. Procesar Sorteo y Seleccionar Ganador

**Actores:** Cron job (ejecuta diariamente a las 00:00 UTC)

**Flujo:**
1. Buscar sorteos con `draw_date <= today` y `status=active`
2. Para cada sorteo:
   ```
   a. Consultar API de Lotería Nacional de Costa Rica
   b. Obtener número ganador del día
   c. Extraer últimos 2 dígitos (o según configuración)
   d. Buscar número ganador en raffle_numbers
   e. Si número fue vendido:
      - raffle.winner_id = raffle_numbers.user_id
      - raffle.winning_number = "42"
      - raffle.status = completed
      - Enviar email/SMS al ganador
      - Enviar email al owner del sorteo
      - Crear settlement (calcular neto después de comisión)
   f. Si número NO fue vendido:
      - raffle.winner_id = NULL
      - raffle.status = completed
      - Enviar email al owner (no hubo ganador)
   ```

**Reglas de negocio:**
- Fuente oficial: Lotería Nacional de Costa Rica
- Si API falla → reintentar 3 veces (cada hora)
- Si falla definitivamente → sorteo pasa a `manual_draw` (admin interviene)
- Comisión de la plataforma: 5-10% (configurable por sorteo)
- Settlement automático: transferencia a cuenta del owner (Fase 2)

#### 5. Gestión de Backoffice (Admin - Almighty)

**Actores:** Usuario con role=admin

**Funcionalidades:**
1. **Gestión de Sorteos:**
   - Ver todos los sorteos (activos, suspendidos, completados)
   - Suspender sorteo (con razón → envía email al owner)
   - Forzar cambio de estado
   - Sorteo manual de ganador
   - Cancelar con reembolso

2. **Gestión de Usuarios:**
   - Ver lista de usuarios
   - Verificar KYC manualmente
   - Suspender/banear usuario
   - Ver historial de compras

3. **Transacciones:**
   - Ver todas las transacciones
   - Ver pagos fallidos
   - Procesar reembolsos

4. **Liquidaciones:**
   - Ver pendientes
   - Crear settlement manual
   - Marcar como pagado

5. **Auditoría:**
   - Ver logs de todas las acciones admin
   - Filtrar por fecha, acción, usuario

**Reglas de negocio:**
- Todas las acciones admin se registran en `audit_logs`
- Suspender sorteo → notificar al owner vía email
- Cancelar sorteo con ventas → reembolso automático a compradores
- Settlement requiere aprobación manual (Fase 1)

### 4.2 Reglas de Negocio Críticas

#### Concurrencia y Reservas
1. **Máximo 10 números por reserva** (previene acaparamiento)
2. **Reserva expira en 5 minutos exactos** (libera números para otros)
3. **Lock distribuido obligatorio** (previene doble venta al 100%)
4. **Idempotencia en reservas** (mismo UUID → misma reserva)
5. **Limpieza automática cada 1 minuto** (libera reservas expiradas)

#### Pagos
1. **Idempotencia obligatoria** (header `Idempotency-Key`)
2. **TTL de idempotencia: 24 horas** (mismo pago no se crea dos veces)
3. **Webhooks con verificación de firma** (seguridad Stripe)
4. **Números solo `sold` cuando pago confirmado** (no con pending)
5. **Refund automático si webhook llega post-expiración**

#### KYC y Trust Levels
1. **none:** Solo puede ver sorteos
2. **email_verified:** Puede crear sorteos y comprar boletos
3. **phone_verified:** (Futuro) Puede comprar hasta ₡50,000
4. **full_kyc:** (Futuro) Puede retirar fondos y crear sorteos premium

#### Sorteos
1. **Máximo 10 sorteos activos por usuario** (evita spam)
2. **DrawDate mínimo: 24 horas en el futuro** (tiempo para ventas)
3. **Precio por número: ₡100 - ₡10,000** (rango razonable)
4. **Mínimo 1 imagen, máximo 5** (presentación adecuada)
5. **Solo owner puede editar/publicar** (seguridad)
6. **Admin puede suspender cualquier sorteo** (moderación)

#### Comisiones y Settlements
1. **Comisión de plataforma: 5-10%** (configurable por sorteo)
2. **Mínimo 60% de números vendidos para realizar sorteo** (parámetro)
3. **Si no se alcanza mínimo → cancelar y reembolsar** (automático)
4. **Settlement automático en Fase 2** (Stripe Connect)
5. **Retiro mínimo: ₡10,000** (evita micro-transacciones)

### 4.3 Integraciones Externas

#### 1. Stripe (Pagos)
- **Producto:** Payment Intents API
- **Webhooks:**
  - `payment_intent.succeeded`
  - `payment_intent.payment_failed`
  - `charge.dispute.created` (chargeback)
- **Seguridad:**
  - Verificación de firma (`Stripe-Signature` header)
  - Idempotencia con `Idempotency-Key`
- **PCI DSS:** Delegado a Stripe (no almacenamos tarjetas)

#### 2. Lotería Nacional de Costa Rica (Fuente de sorteo)
- **API:** (En investigación - puede requerir scraping)
- **Alternativa:** Entrada manual por admin (Fase 1)
- **Backup:** Si API falla → manual draw

#### 3. SendGrid (Emails - Fase 2)
- **Actualmente:** SMTP propio (sorteos.club)
- **Plantillas:**
  - Verificación de email
  - Confirmación de compra
  - Notificación de ganador
  - Confirmación de sorteo (owner)
  - Alertas admin

#### 4. Twilio (SMS - Fase 2)
- Verificación de teléfono
- Notificación de ganador
- 2FA (futuro)

---

## 5. ESTADO ACTUAL DEL DESARROLLO

### 5.1 Fase Actual

**Sprint:** MVP - Autenticación y Gestión de Sorteos
**Duración estimada:** 8-10 semanas
**Progreso:** ~60% completado

### 5.2 Funcionalidades Completadas ✅

#### Backend (Go)

**Autenticación y Usuarios:**
- [x] Registro de usuarios con validación
- [x] Login con JWT (access + refresh tokens)
- [x] Verificación de email con código de 6 dígitos
- [x] Refresh token con rotación
- [x] Logout (invalidar tokens)
- [x] Middleware de autenticación
- [x] Middleware de autorización (roles)
- [x] RBAC (user, admin)
- [x] KYC levels (none, email_verified)
- [x] Hashing de contraseñas con bcrypt

**Gestión de Sorteos:**
- [x] Crear sorteo (draft)
- [x] Listar sorteos (paginado, filtros)
- [x] Ver detalle de sorteo
- [x] Actualizar sorteo (owner only)
- [x] Publicar sorteo (draft → active)
- [x] Generación automática de números

**Categorías:**
- [x] Listar categorías predefinidas

**Imágenes:**
- [x] Subir imágenes (filesystem local)
- [x] Eliminar imágenes
- [x] Establecer imagen principal
- [x] Validación de formatos (JPG, PNG)

**Admin (Almighty):**
- [x] Listar todos los sorteos
- [x] Cancelar sorteo con reembolso
- [x] Forzar cambio de estado
- [x] Sorteo manual de ganador
- [x] Ver transacciones de sorteo
- [x] Logs de auditoría

**Infraestructura:**
- [x] Configuración con Viper (.env)
- [x] Logging estructurado con Zap
- [x] Manejo de errores customizados
- [x] CORS configurado
- [x] Rate limiting (básico)
- [x] Migraciones SQL (10 archivos)
- [x] Health checks (/health, /ready)
- [x] Servicio systemd (sorteos-api)

**Sistema de Emails:**
- [x] SMTP propio configurado (sorteos.club)
- [x] Plantillas HTML para emails
- [x] Verificación de email
- [x] Confirmación de registro
- [x] DKIM, SPF, DMARC configurados

#### Frontend (React + TypeScript)

**Autenticación:**
- [x] Página de registro con validación completa
- [x] Página de login
- [x] Página de verificación de email (código 6 dígitos)
- [x] Manejo de tokens (access + refresh)
- [x] Refresh automático de tokens
- [x] Logout
- [x] Protected routes
- [x] Redirección automática si no autenticado

**Gestión de Sorteos:**
- [x] Listar sorteos (grid view)
- [x] Ver detalle de sorteo
- [x] Crear sorteo (formulario multi-step)
- [x] Subir imágenes
- [x] Grid de números (visualización)
- [x] Filtros por categoría

**Dashboard:**
- [x] Dashboard básico de usuario
- [x] Mostrar información de perfil
- [x] Ver mis sorteos creados
- [x] Ver mis participaciones

**UI Components:**
- [x] Button (variants: default, destructive, outline, secondary, ghost)
- [x] Input con validación
- [x] Card (Header, Content, Footer)
- [x] Alert (success, warning, error, info)
- [x] Badge (estados)
- [x] LoadingSpinner
- [x] EmptyState
- [x] PasswordStrength indicator
- [x] Navbar con UserMenu
- [x] Layout principal

**Estado:**
- [x] Zustand store para auth
- [x] Zustand store para carrito (preparado)
- [x] React Query para data fetching
- [x] Interceptor de Axios para auth

**Utilidades:**
- [x] API client configurado
- [x] Helpers de formato (fecha, moneda)
- [x] Validaciones con Zod

### 5.3 Funcionalidades en Progreso 🚧

#### Backend

**Sistema de Reservas:**
- [ ] Endpoint POST /raffles/:id/reservations
- [ ] Locks distribuidos con Redis (SETNX)
- [ ] Manejo de concurrencia (1000+ requests simultáneos)
- [ ] Liberación automática de reservas expiradas (cron job)
- [ ] Idempotencia de reservas

**Sistema de Pagos:**
- [ ] Integración completa de Stripe
- [ ] Endpoint POST /payments
- [ ] Webhook handler con verificación de firma
- [ ] Idempotencia de pagos
- [ ] Manejo de 3D Secure
- [ ] Refunds automáticos

**Sorteo de Ganadores:**
- [ ] Cron job diario
- [ ] Integración con API de Lotería Nacional CR
- [ ] Sorteo manual (admin)
- [ ] Notificaciones a ganadores

**Settlements:**
- [ ] Cálculo de comisiones
- [ ] Creación de settlements
- [ ] Transferencias (Fase 2 - Stripe Connect)

#### Frontend

**Checkout Flow:**
- [ ] Página de selección de números
- [ ] Timer de reserva (5 min)
- [ ] Integración de Stripe Elements
- [ ] Página de pago
- [ ] Confirmación de compra
- [ ] Comprobante digital

**Dashboard Avanzado:**
- [ ] Ver mis números comprados
- [ ] Historial de compras
- [ ] Ver sorteos ganados
- [ ] Estadísticas personales

**Admin Panel:**
- [ ] Dashboard de administración
- [ ] Gestión de usuarios
- [ ] Gestión de sorteos
- [ ] Ver transacciones
- [ ] Logs de auditoría

### 5.4 Funcionalidades Pendientes (Backlog)

#### Fase 2 (Semanas 11-22)

**Múltiples PSPs:**
- [ ] Integración de PayPal
- [ ] Procesador local Costa Rica
- [ ] Selector de método de pago

**Modo sin cobro:**
- [ ] Sorteos gratuitos (sponsor)
- [ ] Sistema de suscripción premium

**Búsqueda avanzada:**
- [ ] Filtros por precio, fecha, categoría
- [ ] Búsqueda por texto
- [ ] Ordenamiento múltiple

**Sistema de afiliados:**
- [ ] Códigos de referido
- [ ] Comisiones a afiliados
- [ ] Dashboard de afiliado

**Multilenguaje:**
- [ ] Español (completo)
- [ ] Inglés (traducciones)
- [ ] i18next configurado

**Comunicación entre usuarios:**
- [ ] Chat vendedor-comprador
- [ ] Preguntas en sorteos
- [ ] Notificaciones en tiempo real

#### Fase 3 (Semanas 23-38)

**Aplicación móvil:**
- [ ] React Native (iOS + Android)
- [ ] Push notifications
- [ ] Compartir sorteos

**Dashboards en tiempo real:**
- [ ] WebSockets
- [ ] Actualización live de números
- [ ] Contador de ventas en vivo

**Marketing automatizado:**
- [ ] Emails de recordatorio
- [ ] Campañas segmentadas
- [ ] A/B testing

**Programa de fidelización:**
- [ ] Sistema de puntos
- [ ] Niveles de usuario
- [ ] Recompensas

### 5.5 Problemas Conocidos y Áreas de Mejora

#### Bugs Conocidos
1. **Timer de reserva no sincroniza con backend** (frontend)
   - Prioridad: Alta
   - Fix estimado: 2 horas

2. **Imágenes no se eliminan del filesystem al borrar sorteo** (backend)
   - Prioridad: Media
   - Fix estimado: 1 hora

3. **Refresh token rotation puede fallar en condiciones de concurrencia** (backend)
   - Prioridad: Alta
   - Fix estimado: 4 horas

#### Deuda Técnica
1. **Tests unitarios limitados** (~20% coverage)
   - Objetivo: 80% coverage
   - Esfuerzo: 2 semanas

2. **Documentación de API (Swagger)** pendiente
   - Herramienta: swag
   - Esfuerzo: 1 semana

3. **Logs de auditoría no implementados en todos los endpoints**
   - Esfuerzo: 3 días

4. **Rate limiting básico (sin diferenciación por endpoint)**
   - Objetivo: Implementar límites específicos
   - Esfuerzo: 2 días

#### Mejoras de Performance
1. **Caché de listados de sorteos en Redis** (implementar)
   - Impacto: Reducción 70% en queries a DB
   - Esfuerzo: 1 día

2. **Lazy loading de imágenes en frontend**
   - Impacto: Mejor UX en listados
   - Esfuerzo: 1 día

3. **Optimización de queries con índices compuestos**
   - Impacto: Queries 3x más rápidas
   - Esfuerzo: 2 días

4. **CDN para imágenes** (Fase 2)
   - Impacto: Carga 5x más rápida
   - Esfuerzo: 1 semana (migración a S3 + CloudFront)

#### Seguridad
1. **Migrar JWT de HS256 a RS256** (producción)
   - Razón: Mejor seguridad con claves asimétricas
   - Esfuerzo: 1 día

2. **Implementar 2FA** (Fase 2)
   - Método: TOTP (Google Authenticator)
   - Esfuerzo: 1 semana

3. **Scan de vulnerabilidades automatizado** (CI/CD)
   - Herramientas: Trivy, Snyk
   - Esfuerzo: 2 días

### 5.6 Métricas de Desarrollo

**Código:**
- Backend Go: 117 archivos, ~15,000 líneas
- Frontend TS/TSX: 67 archivos, ~8,000 líneas
- Documentación: 10 archivos MD, 181 KB

**Commits:**
- Total: ~350 commits
- Frecuencia: 15-20 commits/semana
- Branches: main, development, feature/*

**Stack Health:**
- PostgreSQL 16: ✅ Activo
- Redis 7: ✅ Activo
- Backend API: ✅ Activo (uptime 99.5%)
- Nginx: ✅ Activo
- SSL: ✅ Válido (Let's Encrypt)

**Performance:**
- Tiempo de build frontend: 10 segundos
- Tiempo de compilación backend: 5 segundos
- Tiempo de startup backend: 2 segundos
- Response time promedio API: 120ms

---

## 📊 RESUMEN EJECUTIVO

### Stack en una línea
**Go + Gin + PostgreSQL + Redis + React + TypeScript + Vite + Tailwind + shadcn/ui**

### Arquitectura en una línea
**Hexagonal (backend) + Feature-based (frontend) + Instalación nativa (sin Docker)**

### Flujo crítico en una línea
**Reserva con locks distribuidos (Redis) → Pago con Stripe → Webhook confirma → Números sold**

### Estado actual en una línea
**MVP 60% completo - Auth y Sorteos ✅ - Pagos y Reservas 🚧**

### Próximo hito
**Implementar sistema completo de reservas con concurrencia + integración de Stripe**
**Estimado: 3-4 semanas**

---

## 📞 CONTACTO Y REFERENCIAS

**Propietario:** Ing. Alonso Alpízar
**Despliegue:** https://sorteos.club
**Documentación completa:** `/opt/Sorteos/Documentacion/`

**Archivos clave de referencia:**
1. `CLAUDE.md` - Contexto rápido para AI
2. `arquitecturaIdeaGeneral.md` - Visión general y concurrencia
3. `stack_tecnico.md` - Tecnologías detalladas
4. `modulos.md` - 7 módulos con código
5. `roadmap.md` - Plan de desarrollo completo

---

**Última actualización:** 2025-11-18
**Versión:** 2.0
**Generado para:** Diseño de skill de Claude Code
