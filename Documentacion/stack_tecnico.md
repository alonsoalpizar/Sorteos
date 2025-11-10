# Stack Técnico - Plataforma de Sorteos

**Versión:** 1.0
**Fecha:** 2025-11-10
**Estado:** Definición cerrada y obligatoria

---

## 1. Resumen Ejecutivo

Este documento define el stack tecnológico **cerrado y obligatorio** para la plataforma de sorteos/rifas en línea. Todas las decisiones técnicas están orientadas a garantizar:

- **Alto rendimiento** con manejo de concurrencia masiva
- **Seguridad transaccional** (pagos, reservas, liquidaciones)
- **Escalabilidad** horizontal y vertical
- **Mantenibilidad** con arquitectura limpia y tipado fuerte
- **Observabilidad** completa (logs, métricas, trazas)

---

## 2. Backend

### 2.1 Lenguaje y Runtime

**Go 1.22+**

**¿Por qué Go?**
- Rendimiento nativo comparable a C/C++
- Concurrencia nativa con goroutines y channels
- Compilación estática (binarios sin dependencias)
- Gestión de memoria eficiente (GC optimizado)
- Ideal para APIs de alto tráfico con transacciones críticas
- Ecosistema maduro para fintech y e-commerce

**Reemplazo:** Solo permitido si existe una limitación técnica crítica. Alternativa: Rust (mayor complejidad), Node.js con TypeScript (menor rendimiento).

---

### 2.2 Framework Web

**Gin (gin-gonic/gin)**
Versión mínima: `v1.9.1`

**¿Por qué Gin?**
- Router extremadamente rápido (httprouter bajo el capó)
- Middlewares composables
- Validación integrada con binding
- Soporte para JSON, XML, YAML
- Comunidad activa y amplia documentación

**Alternativas evaluadas:**
- **Chi**: Más minimalista pero menos ecosistema
- **Fiber**: Inspirado en Express, pero menos idiomatic Go
- **Echo**: Similar rendimiento, menor adopción

**Estructura base:**
```go
r := gin.Default()
r.Use(middleware.CORS())
r.Use(middleware.RateLimiter())
r.Use(middleware.Auth())
r.POST("/raffles/:id/reservations", handlers.CreateReservation)
```

---

### 2.3 Dependencias Core

#### Autenticación y Seguridad

**golang-jwt/jwt/v5**
Versión: `v5.2.0+`

- Generación y validación de JWT
- Soporte para RS256, HS256
- Claims personalizados

**Alternativa:** `go-chi/jwtauth` (más acoplado a Chi router)

**Implementación:**
```go
token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
    "user_id": user.ID,
    "role": user.Role,
    "exp": time.Now().Add(15 * time.Minute).Unix(),
})
```

---

#### ORM / Query Builder

**Opción A: GORM (gorm.io/gorm)**
Versión: `v1.25.0+`

- ORM completo con relaciones, hooks, migraciones
- Auto-migraciones (desarrollo)
- Transacciones explícitas

**Opción B: sqlc**
Versión: `v1.25.0+`

- Generación de código Go desde queries SQL
- Type-safe, sin reflection
- Rendimiento superior (consultas compiladas)

**Decisión:** **GORM** para MVP (velocidad de desarrollo), migrar a **sqlc** en módulos críticos (reservas, pagos) si se requiere optimización.

**Ejemplo GORM:**
```go
db.Transaction(func(tx *gorm.DB) error {
    if err := tx.Create(&reservation).Error; err != nil {
        return err
    }
    return tx.Model(&Raffle{}).Where("id = ?", raffleID).
        UpdateColumn("reserved_count", gorm.Expr("reserved_count + ?", len(numbers))).Error
})
```

---

#### Driver PostgreSQL

**jackc/pgx/v5**
Versión: `v5.5.0+`

- Driver PostgreSQL de alto rendimiento
- Soporte para tipos nativos (JSONB, UUID, arrays)
- Connection pooling eficiente
- Compatible con GORM y sqlc

**Configuración:**
```go
dsn := "postgres://user:pass@localhost:5432/sorteos?sslmode=disable&pool_max_conns=25"
db, err := gorm.Open(postgres.New(postgres.Config{
    DriverName: "pgx",
    DSN: dsn,
}), &gorm.Config{})
```

---

#### Cliente Redis

**redis/go-redis/v9**
Versión: `v9.5.0+`

**Casos de uso:**
- **Locks distribuidos** para reservas de números
- **Caché** de sorteos activos, configuraciones
- **Rate limiting** por IP/usuario
- **Colas de jobs** diferidos (notificaciones, liquidaciones)
- **Sesiones** y refresh tokens
- **Idempotencia** de pagos (Idempotency-Key → transaction_id)

**Ejemplo lock distribuido:**
```go
lockKey := fmt.Sprintf("lock:raffle:%d:number:%s", raffleID, number)
acquired, err := rdb.SetNX(ctx, lockKey, userID, 30*time.Second).Result()
if !acquired {
    return errors.New("número ya reservado")
}
defer rdb.Del(ctx, lockKey)
```

---

#### Logging

**uber-go/zap**
Versión: `v1.27.0+`

- Logging estructurado de alto rendimiento
- Niveles: Debug, Info, Warn, Error, Fatal
- Campos tipados (evita allocations)
- Integración con OpenTelemetry

**Configuración:**
```go
logger, _ := zap.NewProduction()
defer logger.Sync()
logger.Info("reserva creada",
    zap.String("trace_id", traceID),
    zap.Int64("user_id", userID),
    zap.Int("raffle_id", raffleID),
)
```

---

#### Configuración

**spf13/viper**
Versión: `v1.18.0+`

- Lectura de `.env`, YAML, JSON, TOML
- Variables de entorno con prefijos
- Hot reload (opcional)
- Validación de configuración requerida

**Estructura `.env`:**
```env
CONFIG_ENV=development
CONFIG_DB_HOST=localhost
CONFIG_DB_PORT=5432
CONFIG_REDIS_URL=redis://localhost:6379
CONFIG_JWT_SECRET=super-secret-key
CONFIG_PAYMENT_PROVIDER=stripe
```

**Código:**
```go
viper.SetEnvPrefix("CONFIG")
viper.AutomaticEnv()
viper.SetConfigFile(".env")
viper.ReadInConfig()
```

---

#### Validación

**go-playground/validator/v10**
Versión: `v10.19.0+`

- Validación de structs con tags
- Reglas personalizadas
- Mensajes de error i18n

**Ejemplo:**
```go
type CreateReservationRequest struct {
    RaffleID int      `json:"raffle_id" validate:"required,gt=0"`
    Numbers  []string `json:"numbers" validate:"required,min=1,max=10,dive,len=2"`
}

if err := validator.Validate(req); err != nil {
    return c.JSON(400, gin.H{"errors": err.Error()})
}
```

---

### 2.4 Arquitectura de Módulos Internos

**Estructura hexagonal/clean:**

```
/backend
├── cmd/
│   └── api/
│       └── main.go                    # Entry point
├── internal/
│   ├── domain/                        # Entidades y reglas de negocio
│   │   ├── raffle.go
│   │   ├── user.go
│   │   ├── reservation.go
│   │   └── payment.go
│   ├── usecase/                       # Casos de uso (aplicación)
│   │   ├── create_raffle.go
│   │   ├── reserve_numbers.go
│   │   └── process_payment.go
│   ├── adapters/
│   │   ├── http/                      # Handlers Gin
│   │   │   ├── raffle_handler.go
│   │   │   └── middleware/
│   │   ├── db/                        # Repositorios GORM/sqlc
│   │   │   ├── raffle_repo.go
│   │   │   └── user_repo.go
│   │   ├── payments/                  # Providers (Stripe, local)
│   │   │   ├── provider.go (interfaz)
│   │   │   ├── stripe.go
│   │   │   └── mock.go
│   │   └── notifier/                  # Email, SMS
│   │       ├── email.go
│   │       └── sms.go
├── pkg/                               # Librerías compartidas
│   ├── logger/
│   ├── config/
│   ├── errors/
│   └── validator/
├── migrations/                        # SQL migrations
│   ├── 001_create_users.up.sql
│   └── 001_create_users.down.sql
├── Makefile
├── go.mod
└── go.sum
```

**Principios:**
- `domain/` **no** tiene dependencias externas (Go puro)
- `usecase/` depende de interfaces de `domain/`
- `adapters/` implementa interfaces y maneja detalles técnicos
- `pkg/` es reutilizable entre proyectos

---

## 3. Frontend

### 3.1 Lenguaje y Runtime

**TypeScript 5.3+**
**Node.js 20 LTS+**

**¿Por qué TypeScript?**
- Type safety en desarrollo
- Detección temprana de errores
- Refactoring seguro
- Ecosistema React 100% compatible
- LSP para mejor DX

---

### 3.2 Build Tool

**Vite 5.0+**

- HMR instantáneo
- Build optimizado (Rollup)
- Soporte nativo para TypeScript, JSX, CSS Modules
- Plugins para PWA, análisis de bundle

**Alternativas descartadas:**
- **Create React App**: Lento, deprecado
- **Webpack**: Configuración compleja
- **Parcel**: Menos ecosistema

---

### 3.3 Framework UI

**React 18.2+**

- Concurrent rendering
- Suspense y Error Boundaries
- Server Components (futuro)
- Ecosistema maduro

---

### 3.4 Librerías Core

#### Routing

**React Router 6.22+**

```tsx
<Routes>
  <Route path="/" element={<HomePage />} />
  <Route path="/raffles/:id" element={<RaffleDetail />} />
  <Route path="/checkout" element={<ProtectedRoute><Checkout /></ProtectedRoute>} />
</Routes>
```

---

#### State Management

**Opción A: Zustand 4.5+** (recomendado para MVP)

- Minimal boilerplate
- TypeScript first
- DevTools

```tsx
const useAuthStore = create<AuthState>((set) => ({
  user: null,
  login: (user) => set({ user }),
  logout: () => set({ user: null }),
}))
```

**Opción B: Redux Toolkit 2.0+**

- Para apps complejas con múltiples slices
- Middleware (thunks, sagas)

---

#### Data Fetching

**TanStack Query (React Query) 5.0+**

- Caché automático
- Refetch automático
- Optimistic updates
- Paginación y scroll infinito

```tsx
const { data, isLoading, error } = useQuery({
  queryKey: ['raffle', id],
  queryFn: () => fetchRaffle(id),
  staleTime: 5 * 60 * 1000, // 5 min
})
```

---

#### HTTP Client

**Axios 1.6+**

- Interceptores para auth (JWT)
- Timeout, retry
- Type-safe con generics de TS

```tsx
const api = axios.create({
  baseURL: import.meta.env.VITE_API_URL,
  timeout: 10000,
})

api.interceptors.request.use((config) => {
  const token = localStorage.getItem('access_token')
  if (token) config.headers.Authorization = `Bearer ${token}`
  return config
})
```

---

#### Forms

**react-hook-form 7.50+ + zod 3.22+**

```tsx
const schema = z.object({
  email: z.string().email(),
  numbers: z.array(z.string().length(2)).min(1).max(10),
})

const { register, handleSubmit, formState: { errors } } = useForm({
  resolver: zodResolver(schema),
})
```

---

#### UI Components

**Tailwind CSS 3.4+ + shadcn/ui**

- Utility-first CSS
- Design tokens customizables
- Componentes accesibles (Radix UI)
- Dark mode con `class` strategy

**Componentes base de shadcn:**
- Button, Input, Select, Card, Table
- Dialog, Toast, Badge, Skeleton
- Form (integrado con react-hook-form)

---

#### Internacionalización

**i18next 23+ + react-i18next**

```tsx
// es.json
{
  "raffle.buy_ticket": "Comprar boleto",
  "raffle.sold_out": "Agotado"
}

// Uso
const { t } = useTranslation()
<Button>{t('raffle.buy_ticket')}</Button>
```

---

#### Testing

**Vitest 1.2+** (unit tests)
**MSW 2.0+** (mock API)
**Testing Library** (component tests)

```tsx
// mock handler
rest.post('/api/reservations', (req, res, ctx) => {
  return res(ctx.json({ id: 123, status: 'pending' }))
})
```

---

### 3.5 Estructura de Carpetas

```
/frontend
├── public/
│   └── assets/
├── src/
│   ├── app/
│   │   ├── App.tsx
│   │   ├── router.tsx
│   │   └── providers.tsx          # React Query, Auth, i18n
│   ├── features/
│   │   ├── auth/
│   │   │   ├── components/
│   │   │   ├── hooks/
│   │   │   ├── api.ts
│   │   │   └── store.ts
│   │   ├── raffles/
│   │   └── checkout/
│   ├── components/                # Shared UI
│   │   ├── ui/                    # shadcn components
│   │   ├── Layout.tsx
│   │   └── ProtectedRoute.tsx
│   ├── lib/
│   │   ├── api.ts                 # Axios instance
│   │   ├── utils.ts
│   │   └── constants.ts
│   ├── styles/
│   │   └── globals.css
│   ├── types/
│   │   └── api.ts
│   └── main.tsx
├── vite.config.ts
├── tailwind.config.js
├── tsconfig.json
└── package.json
```

---

## 4. Base de Datos

### 4.1 Motor

**PostgreSQL 15.5+**

**¿Por qué PostgreSQL?**
- ACID compliant (crítico para transacciones de pago)
- Índices avanzados (B-tree, GIN, GiST)
- JSONB para datos semi-estructurados
- Transacciones con niveles de aislamiento configurables
- Vistas materializadas para KPIs
- Extensions: `uuid-ossp`, `pg_trgm` (búsqueda fuzzy)

**Configuración recomendada:**
```sql
-- postgresql.conf
max_connections = 100
shared_buffers = 4GB
effective_cache_size = 12GB
work_mem = 64MB
maintenance_work_mem = 1GB
```

---

### 4.2 Migraciones

**golang-migrate/migrate**
Versión: `v4.17+`

```bash
migrate -path ./migrations -database "postgres://user:pass@localhost/sorteos" up
```

**Ejemplo migration:**
```sql
-- 001_create_raffles.up.sql
CREATE TABLE raffles (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id),
    title VARCHAR(255) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'draft',
    draw_date TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    CONSTRAINT check_status CHECK (status IN ('draft', 'active', 'suspended', 'completed'))
);

CREATE INDEX idx_raffles_status_draw_date ON raffles(status, draw_date);
```

---

### 4.3 Esquema Crítico (resumen)

**Tablas principales:**
- `users` (id, email, phone, role, kyc_level, created_at)
- `raffles` (id, user_id, title, description, draw_date, status, lottery_source)
- `raffle_numbers` (raffle_id, number, user_id, status, reserved_at, sold_at)
- `reservations` (id, raffle_id, user_id, numbers[], status, expires_at, idempotency_key)
- `payments` (id, reservation_id, provider, amount, status, metadata JSONB)
- `settlements` (id, raffle_id, user_id, gross_amount, fees, net_amount, status)
- `audit_logs` (id, user_id, action, entity_type, entity_id, ip, metadata JSONB)

**Índices críticos:**
```sql
CREATE UNIQUE INDEX idx_reservations_idempotency ON reservations(idempotency_key) WHERE status != 'cancelled';
CREATE INDEX idx_raffle_numbers_available ON raffle_numbers(raffle_id, status) WHERE status = 'available';
```

---

## 5. Concurrencia, Caché y Colas

### 5.1 Redis

**Versión:** 7.2+

**Modos:**
- **Standalone** (desarrollo/staging)
- **Sentinel** (producción con HA)
- **Cluster** (escala horizontal si > 100k usuarios activos)

**Casos de uso detallados:**

#### A. Locks Distribuidos (Reservas)
```go
// Patrón: SETNX con TTL
lockKey := fmt.Sprintf("lock:raffle:%d:num:%s", raffleID, number)
acquired := rdb.SetNX(ctx, lockKey, userID, 30*time.Second)
```

#### B. Caché de Sorteos Activos
```go
cacheKey := fmt.Sprintf("raffle:%d", id)
rdb.Set(ctx, cacheKey, raffleJSON, 10*time.Minute)
```

#### C. Rate Limiting (Token Bucket)
```go
key := fmt.Sprintf("ratelimit:user:%d:reserve", userID)
count := rdb.Incr(ctx, key)
if count == 1 {
    rdb.Expire(ctx, key, 1*time.Minute)
}
if count > 10 { // máx 10 reservas/min
    return errors.New("rate limit exceeded")
}
```

#### D. Idempotencia de Pagos
```go
key := fmt.Sprintf("payment:idempotency:%s", idempotencyKey)
exists := rdb.Exists(ctx, key)
if exists {
    // Retornar pago existente
}
rdb.Set(ctx, key, paymentID, 24*time.Hour)
```

---

## 6. Mensajería y Notificaciones

### 6.1 Arquitectura

**Patrón:** Adaptador con drivers intercambiables

**Interfaz Go:**
```go
type Notifier interface {
    SendEmail(ctx context.Context, to string, template string, data map[string]any) error
    SendSMS(ctx context.Context, to string, message string) error
}
```

**Drivers:**
- **Email**: SendGrid, AWS SES, SMTP
- **SMS**: Twilio, AWS SNS

**Configuración:**
```env
CONFIG_NOTIFIER_EMAIL_DRIVER=sendgrid
CONFIG_NOTIFIER_SMS_DRIVER=twilio
CONFIG_SENDGRID_API_KEY=xxx
CONFIG_TWILIO_ACCOUNT_SID=xxx
```

---

## 7. Infraestructura

### 7.1 Containerización

**Docker 24+ + Docker Compose**

**Servicios:**
```yaml
# docker-compose.yml
services:
  api:
    build: ./backend
    ports: ["8080:8080"]
    env_file: .env
    depends_on: [postgres, redis]

  postgres:
    image: postgres:15-alpine
    volumes: ["pgdata:/var/lib/postgresql/data"]
    environment:
      POSTGRES_DB: sorteos

  redis:
    image: redis:7-alpine
    command: redis-server --appendonly yes
```

---

### 7.2 CI/CD

**Etapas obligatorias:**

1. **Lint** (golangci-lint, ESLint)
2. **Tests** (go test -race, vitest)
3. **Build** (binario Go, bundle Vite)
4. **Migrations dry-run** (verificar sintaxis)
5. **Security scan** (Trivy, npm audit)
6. **Deploy** (staging → smoke tests → production)

**GitHub Actions ejemplo:**
```yaml
- name: Run tests
  run: make test
- name: Build
  run: make build
- name: Migrate
  run: migrate -path ./migrations -database $DB_URL up
```

---

### 7.3 Configuración (12-Factor)

**Principios:**
- **Nunca** secretos en código
- Variables con prefijo `CONFIG_*`
- `.env.example` en repo, `.env` en `.gitignore`
- Secrets en vault (AWS Secrets Manager, HashiCorp Vault)

---

## 8. Observabilidad

### 8.1 Logging

**Estructura obligatoria:**
```json
{
  "level": "info",
  "timestamp": "2025-11-10T10:30:00Z",
  "trace_id": "abc123",
  "user_id": 456,
  "message": "reserva creada",
  "raffle_id": 789,
  "numbers": ["01", "15"]
}
```

---

### 8.2 Métricas

**Prometheus + Grafana**

**Métricas clave:**
- `http_requests_total{method, path, status}`
- `reservation_duration_seconds`
- `payment_success_rate`
- `active_reservations_gauge`

**Endpoint:**
```go
import "github.com/prometheus/client_golang/prometheus/promhttp"
r.GET("/metrics", gin.WrapH(promhttp.Handler()))
```

---

### 8.3 Trazas

**OpenTelemetry**

- Propagación de `trace_id` en headers
- Spans por operación (DB, Redis, HTTP)
- Exportar a Jaeger/Tempo

---

## 9. Seguridad

### 9.1 Dependencias

- **Actualizaciones automáticas**: Dependabot, Renovate
- **Auditorías**: `go mod verify`, `npm audit`
- **CVE scanning**: Trivy, Snyk

### 9.2 Secrets

- **Rotación**: JWT refresh tokens cada 7 días
- **Encriptación**: TLS 1.3 en tránsito, AES-256 en reposo

---

## 10. Reemplazo de Tecnologías

| Componente | Actual | Alternativa | Condición |
|------------|--------|-------------|-----------|
| Go | 1.22+ | Rust, Node.js | Solo si limitación crítica |
| Gin | v1.9+ | Chi, Fiber | Solo con aprobación arquitecto |
| PostgreSQL | 15+ | MySQL 8 | No recomendado (menos features) |
| Redis | 7+ | Memcached | No (locks distribuidos necesarios) |
| React | 18+ | Vue 3, Svelte | Solo para módulos aislados |
| GORM | v1.25+ | sqlc | Permitido en módulos críticos |

---

## 11. Versiones y Compatibilidad

**Matriz de compatibilidad:**
```
Go 1.22.0       → PostgreSQL 15.5 (pgx v5)
Gin 1.9.1       → Go 1.20+
React 18.2.0    → Node.js 20 LTS
TypeScript 5.3  → React 18+
Redis 7.2       → go-redis v9
```

---

## 12. Makefile (Backend)

```makefile
.PHONY: run test build migrate-up migrate-down

run:
	go run cmd/api/main.go

test:
	go test -v -race -cover ./...

build:
	go build -o bin/api cmd/api/main.go

migrate-up:
	migrate -path ./migrations -database $(DB_URL) up

migrate-down:
	migrate -path ./migrations -database $(DB_URL) down 1

lint:
	golangci-lint run

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down
```

---

## 13. Dependencias Completas

### Backend (go.mod)
```go
require (
    github.com/gin-gonic/gin v1.9.1
    github.com/golang-jwt/jwt/v5 v5.2.0
    gorm.io/gorm v1.25.5
    gorm.io/driver/postgres v1.5.4
    github.com/jackc/pgx/v5 v5.5.0
    github.com/redis/go-redis/v9 v9.5.0
    go.uber.org/zap v1.27.0
    github.com/spf13/viper v1.18.0
    github.com/go-playground/validator/v10 v10.19.0
    github.com/prometheus/client_golang v1.19.0
)
```

### Frontend (package.json)
```json
{
  "dependencies": {
    "react": "^18.2.0",
    "react-dom": "^18.2.0",
    "react-router-dom": "^6.22.0",
    "@tanstack/react-query": "^5.0.0",
    "zustand": "^4.5.0",
    "axios": "^1.6.0",
    "react-hook-form": "^7.50.0",
    "zod": "^3.22.0",
    "i18next": "^23.7.0",
    "react-i18next": "^14.0.0"
  },
  "devDependencies": {
    "@types/react": "^18.2.0",
    "typescript": "^5.3.0",
    "vite": "^5.0.0",
    "vitest": "^1.2.0",
    "tailwindcss": "^3.4.0",
    "eslint": "^8.56.0",
    "msw": "^2.0.0"
  }
}
```

---

## 14. Consideraciones Finales

- Este stack está **cerrado** salvo justificación técnica crítica
- Toda nueva dependencia debe pasar revisión de seguridad
- Priorizar librerías con mantenimiento activo (commits < 6 meses)
- Documentar **por qué** se elige una librería en este archivo

---

**Próximos pasos:**
Ver [roadmap.md](./roadmap.md) para fases de implementación.
