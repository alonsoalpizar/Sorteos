# Resumen Ejecutivo para Skill de Claude Code

**Proyecto:** Plataforma de Sorteos/Rifas en Línea
**Fecha:** 2025-11-18
**Propósito:** Contexto condensado para diseño de skill

---

## 🎯 RESPUESTAS RÁPIDAS A TUS PREGUNTAS

### 1. Stack Tecnológico

**Backend:**
- Go 1.22+ (Gin framework)
- PostgreSQL 16 (ACID, transaccional)
- Redis 7 (locks distribuidos, cache)
- GORM (ORM)
- JWT (auth), Zap (logging), Viper (config)

**Frontend:**
- React 18 + TypeScript 5.3
- Vite 5 (build tool)
- Tailwind CSS + shadcn/ui
- Zustand (state), React Query (data fetching)
- Axios (HTTP client)

**Infraestructura:**
- Nginx (reverse proxy + SSL)
- systemd (gestión de servicios)
- Instalación nativa (sin Docker desde nov 2025)

**Pagos:**
- Stripe (MVP)
- PayPal (Fase 2)

### 2. Arquitectura Principal

**Estructura de directorios:**
```
/opt/Sorteos/
├── backend/              # Go API (117 archivos)
│   ├── cmd/api/          # Entry point, routes
│   ├── internal/
│   │   ├── domain/       # Entidades (User, Raffle, etc.)
│   │   ├── usecase/      # Lógica de negocio
│   │   └── adapters/     # HTTP, DB, Payments
│   ├── pkg/              # Logger, Config, Errors
│   └── migrations/       # SQL migrations
├── frontend/             # React SPA (67 archivos)
│   ├── src/
│   │   ├── features/     # auth, raffles, dashboard
│   │   ├── components/   # UI components
│   │   └── lib/          # Utilidades
│   └── dist/             # Build (servido por backend)
└── Documentacion/        # 10 docs técnicos
```

**Componentes principales:**
1. **Domain Layer** - Entidades puras sin dependencias
2. **Use Cases** - Lógica de aplicación
3. **Adapters** - HTTP handlers, DB repos, Payments

**Separación backend/frontend:**
- SÍ, similar a DIV
- Backend: API RESTful (puerto 8080)
- Frontend: SPA servido por backend desde `dist/`
- Comunicación: Axios + JWT

### 3. Decisiones Técnicas Importantes

#### Patrones de Diseño:
- **Hexagonal Architecture** (Ports & Adapters)
- **Repository Pattern**
- **Factory Pattern** (payment providers)
- **Strategy Pattern** (lottery sources)

#### Convenciones de Naming:

**Backend Go:**
- Archivos: `snake_case.go`
- Structs: `PascalCase`
- Funciones exportadas: `PascalCase`
- Funciones privadas: `camelCase`

**Frontend TypeScript:**
- Componentes: `PascalCase.tsx`
- Hooks: `useName()`
- Utilidades: `camelCase.ts`

#### Reglas de Validación:

**Críticas:**
- Email: único, formato válido
- Password: 12+ chars, mayúscula, minúscula, número, símbolo
- DrawDate: futuro (mínimo 24h)
- PricePerNumber: ₡100 - ₡10,000

#### Manejo de Errores:

**Backend:**
- Errores tipados (ErrNotFound, ErrUnauthorized, etc.)
- Logging estructurado con Zap
- HTTP status codes estándar

**Frontend:**
- React Query error handling
- Toast notifications
- Interceptor de Axios para 401/429

#### Seguridad:

**Autenticación:**
- JWT: Access token (15 min) + Refresh token (7 días)
- Almacenamiento: Memory (access), HttpOnly Cookie (refresh)
- Rotación de refresh tokens

**Autorización:**
- RBAC: user, admin
- KYC levels: none, email_verified, phone_verified, full_kyc
- Middleware de verificación

**Rate Limiting:**
- Redis Token Bucket
- Límites por endpoint:
  - POST /auth/login: 5 req/min
  - POST /reservations: 10 req/min
  - POST /payments: 5 req/min

**Prevención OWASP:**
- SQL Injection: GORM escapa automáticamente
- XSS: React escapa automáticamente
- CSRF: JWT en headers (no cookies)

### 4. Contexto de Negocio

#### Flujos Principales:

**1. Registrarse y Verificar Email**
```
Usuario → Formulario → Backend crea user → Envía código 6 dígitos
→ Usuario ingresa código → Backend verifica → kyc_level=email_verified
```

**2. Crear Sorteo**
```
Usuario → Formulario (título, precio, números, imágenes)
→ Backend valida → Crea en estado draft → Genera números
→ Usuario publica → Estado cambia a active
```

**3. Comprar Boleto (CRÍTICO - Alta Concurrencia)**
```
FASE 1 - Reserva:
Usuario selecciona números → POST /reservations
→ Backend:
  1. Lock distribuido en Redis (SETNX)
  2. Verificar disponibilidad en DB
  3. Crear reserva (expires_at = now + 5min)
  4. Liberar lock
→ Frontend muestra timer 5 min

FASE 2 - Pago:
Usuario ingresa tarjeta → Stripe.js tokeniza
→ POST /payments con payment_method_id
→ Backend crea PaymentIntent en Stripe
→ Webhook confirma pago → Números pasan a sold

FASE 3 - Limpieza:
Cron job cada 1 min libera reservas expiradas
```

**4. Sorteo de Ganador**
```
Cron job diario → Consulta Lotería Nacional CR
→ Extrae número ganador → Busca en raffle_numbers
→ Si vendido: marca ganador + notifica
→ Si no vendido: winner_id=NULL
→ Crea settlement
```

**5. Backoffice Admin (Almighty)**
```
Admin puede:
- Suspender sorteos (con razón)
- Cancelar con reembolso
- Sorteo manual
- Ver transacciones
- Gestionar usuarios
- Ver logs de auditoría
```

#### Reglas de Negocio Críticas:

**Concurrencia:**
- Máximo 10 números por reserva
- Reserva expira en 5 minutos exactos
- Lock distribuido obligatorio (SETNX en Redis)
- Idempotencia en reservas (UUID)

**Pagos:**
- Idempotencia obligatoria (header Idempotency-Key)
- TTL de 24h para idempotencia
- Webhooks con verificación de firma
- Refund automático si webhook llega post-expiración

**KYC:**
- none: Solo ver sorteos
- email_verified: Crear sorteos y comprar
- phone_verified: (Futuro) Límites mayores
- full_kyc: (Futuro) Retirar fondos

**Sorteos:**
- Máximo 10 sorteos activos por usuario
- DrawDate mínimo: 24h en futuro
- Comisión: 5-10% (configurable)
- Mínimo 60% vendido para realizar sorteo

#### Integraciones Externas:

1. **Stripe** - Pagos con Payment Intents + Webhooks
2. **Lotería Nacional CR** - Fuente de sorteo oficial
3. **SMTP propio** - Emails transaccionales (sorteos.club)
4. **SendGrid** - (Fase 2) Emails masivos
5. **Twilio** - (Fase 2) SMS

### 5. Estado Actual del Desarrollo

**Fase:** MVP (60% completado)
**Duración:** 8-10 semanas
**Progreso:** Semana 6

#### ✅ Completado:

**Backend:**
- [x] Auth completo (registro, login, verificación, JWT)
- [x] CRUD de sorteos
- [x] Gestión de imágenes
- [x] Admin panel (suspender, cancelar, sorteo manual)
- [x] Sistema de emails SMTP
- [x] 10 migraciones SQL
- [x] Servicio systemd

**Frontend:**
- [x] Registro y login
- [x] Verificación de email
- [x] Listar y ver sorteos
- [x] Crear sorteo con imágenes
- [x] Dashboard usuario
- [x] Protected routes
- [x] 20+ componentes UI (shadcn/ui)

#### 🚧 En Progreso:

**Backend:**
- [ ] Sistema de reservas con locks distribuidos
- [ ] Integración completa de Stripe
- [ ] Webhooks con verificación de firma
- [ ] Cron job para limpieza de reservas
- [ ] Sorteo automático de ganadores

**Frontend:**
- [ ] Checkout flow completo
- [ ] Timer de reserva
- [ ] Stripe Elements integrado
- [ ] Comprobante digital
- [ ] Dashboard avanzado

#### ❌ Pendiente (Backlog):

**Fase 2:**
- [ ] PayPal integration
- [ ] Búsqueda avanzada
- [ ] Afiliados
- [ ] Multilenguaje (i18next)
- [ ] Chat usuario-vendedor

**Fase 3:**
- [ ] App móvil (React Native)
- [ ] WebSockets (tiempo real)
- [ ] Marketing automatizado
- [ ] Programa de fidelización

#### 🐛 Problemas Conocidos:

1. **Timer de reserva no sincroniza** (Alta prioridad)
2. **Imágenes no se borran al eliminar sorteo** (Media prioridad)
3. **Refresh token rotation bajo concurrencia** (Alta prioridad)

#### 📊 Deuda Técnica:

1. Tests unitarios (~20% coverage, objetivo 80%)
2. Documentación Swagger pendiente
3. Logs de auditoría incompletos
4. Rate limiting básico (mejorar granularidad)

#### ⚡ Mejoras de Performance:

1. Caché de listados en Redis (reducción 70% queries)
2. Lazy loading imágenes
3. Índices compuestos en DB
4. CDN para imágenes (Fase 2)

---

## 🔑 CONCEPTOS CLAVE PARA EL SKILL

### Problema Central del Sistema

**Doble venta de números de sorteo en alta concurrencia**

**Solución (3 capas):**
1. **Lock distribuido en Redis** (SETNX, TTL 30s)
2. **Verificación en PostgreSQL** (transacción ACID)
3. **Reserva temporal** (5 min para pagar)

**Código ejemplo:**
```go
// 1. Adquirir lock
lockKey := fmt.Sprintf("lock:raffle:%d:num:%s", raffleID, number)
acquired := rdb.SetNX(ctx, lockKey, userID, 30*time.Second)
if !acquired {
    return errors.New("número ya reservado")
}
defer rdb.Del(ctx, lockKey)

// 2. Verificar en DB (transacción)
db.Transaction(func(tx *gorm.DB) error {
    // Verificar disponibilidad
    // Crear reserva
    // Actualizar números a reserved
})

// 3. TTL automático en Redis
rdb.Set(ctx, fmt.Sprintf("reservation:%d", resID), res, 5*time.Minute)
```

### Flujo Crítico Simplificado

```
Usuario clickea número
    ↓
Lock Redis (30s)
    ↓
Verificar DB (transacción)
    ↓
Crear reserva (5 min)
    ↓
Liberar lock
    ↓
Timer 5 min (frontend)
    ↓
Usuario paga (Stripe)
    ↓
Webhook confirma
    ↓
Números → sold
```

### Entidades Principales

```go
User {
  id, email, password_hash, role, kyc_level, status
}

Raffle {
  id, user_id, title, status, draw_date, price_per_number, total_numbers
}

RaffleNumber {
  raffle_id, number, user_id, status (available/reserved/sold)
}

Reservation {
  id, raffle_id, user_id, numbers[], status, expires_at, idempotency_key
}

Payment {
  id, reservation_id, provider, amount, status, external_id, idempotency_key
}
```

### Restricciones Visuales (CRÍTICO)

**PROHIBIDO:**
- Morado, púrpura, violeta
- Rosa, pink, magenta
- Fucsia, gradientes neón

**PERMITIDO:**
- Azul #3B82F6 (primary)
- Slate #64748B (secondary)
- Verde #10B981 (success)
- Ámbar #F59E0B (warning)
- Rojo #EF4444 (error)

**Referencias:** Stripe, Linear, Vercel, Coinbase

---

## 📋 CHECKLIST PARA SKILL

**El skill debe conocer:**
- [x] Stack tecnológico completo
- [x] Estructura de directorios
- [x] Arquitectura hexagonal
- [x] Patrones de diseño usados
- [x] Convenciones de naming
- [x] Reglas de validación
- [x] Manejo de errores
- [x] Seguridad (JWT, rate limiting)
- [x] Flujos de negocio críticos
- [x] Problema de concurrencia y solución
- [x] Integraciones externas
- [x] Estado actual (completado vs pendiente)
- [x] Restricciones visuales
- [x] Comandos útiles (systemd, build, deploy)

**El skill debe poder:**
- [ ] Generar código Go siguiendo arquitectura hexagonal
- [ ] Generar componentes React siguiendo convenciones
- [ ] Sugerir fixes para problemas de concurrencia
- [ ] Proponer mejoras de performance
- [ ] Validar código contra reglas de negocio
- [ ] Generar tests unitarios
- [ ] Documentar endpoints (Swagger)
- [ ] Sugerir índices de DB según queries

---

## 🚀 PRÓXIMOS PASOS CON EL SKILL

### Prioridad 1: Sistema de Reservas
- Implementar locks distribuidos
- Tests de concurrencia (1000 usuarios simultáneos)
- Cron job de limpieza

### Prioridad 2: Integración de Pagos
- Stripe Payment Intents completo
- Webhooks con verificación
- Idempotencia en todos los flows

### Prioridad 3: Tests y Documentación
- Coverage 20% → 80%
- Swagger completo
- Tests de carga (k6)

---

**Documento completo:** [SORTEOS_CONTEXTO_COMPLETO.md](SORTEOS_CONTEXTO_COMPLETO.md)
**Referencias:** `/opt/Sorteos/Documentacion/` (10 docs)

**Última actualización:** 2025-11-18
**Versión:** 1.0
