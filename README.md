# Plataforma de Sorteos

Sistema de sorteos/rifas en línea con gestión de pagos, reservas de números y backoffice administrativo.

## 📋 Documentación

Toda la documentación técnica se encuentra en [/Documentacion](/Documentacion/):

- **[arquitecturaIdeaGeneral.md](Documentacion/arquitecturaIdeaGeneral.md)** - Visión general del sistema
- **[stack_tecnico.md](Documentacion/stack_tecnico.md)** - Stack tecnológico completo (Go, React, PostgreSQL, Redis)
- **[roadmap.md](Documentacion/roadmap.md)** - Plan de desarrollo por fases
- **[modulos.md](Documentacion/modulos.md)** - Módulos del sistema y casos de uso
- **[estandar_visual.md](Documentacion/estandar_visual.md)** - Design system y componentes UI
  - ⚠️ **RESTRICCIÓN:** NO usar morado, púrpura, violeta, rosa, magenta (paleta profesional azul/gris)
  - Ver [paleta-visual-aprobada.md](Documentacion/.paleta-visual-aprobada.md) para referencia rápida
- **[seguridad.md](Documentacion/seguridad.md)** - Políticas de seguridad (JWT, RBAC, rate limiting)
- **[pagos_integraciones.md](Documentacion/pagos_integraciones.md)** - Sistema de pagos (Stripe, webhooks, idempotencia)
- **[parametrizacion_reglas.md](Documentacion/parametrizacion_reglas.md)** - Parámetros configurables
- **[operacion_backoffice.md](Documentacion/operacion_backoffice.md)** - Operación administrativa
- **[terminos_y_condiciones_impacto.md](Documentacion/terminos_y_condiciones_impacto.md)** - Cumplimiento legal (GDPR, PCI DSS)

## 🏗️ Estructura del Proyecto

```
/opt/Sorteos/
├── backend/              # API en Go
│   ├── cmd/
│   │   └── api/          # Entry point
│   ├── internal/
│   │   ├── domain/       # Entidades y reglas de negocio
│   │   ├── usecase/      # Casos de uso
│   │   └── adapters/     # HTTP, DB, Payments, Notifier
│   ├── pkg/              # Librerías compartidas
│   └── migrations/       # Migraciones SQL
├── frontend/             # SPA en React + TypeScript
│   ├── src/
│   │   ├── app/          # Router y providers
│   │   ├── features/     # Módulos (auth, raffles, checkout)
│   │   ├── components/   # Componentes UI (shadcn/ui)
│   │   └── lib/          # Utilidades
│   └── public/
└── Documentacion/        # Docs técnicas
```

## 🚀 Stack Tecnológico

### Backend
- **Go 1.22+** con Gin
- **PostgreSQL 15+** (base de datos principal)
- **Redis 7+** (cache, locks distribuidos, rate limiting)
- **GORM** (ORM) o **sqlc** (type-safe queries)
- **Zap** (logging), **Viper** (config), **JWT** (auth)

### Frontend
- **React 18+** con **TypeScript**
- **Vite** (build tool)
- **TanStack Query** (data fetching)
- **Zustand** (state management)
- **Tailwind CSS + shadcn/ui** (UI components)
- **React Hook Form + Zod** (validación)

### Infraestructura
- **Docker + Docker Compose**
- **Nginx** (reverse proxy)
- **Let's Encrypt** (SSL/TLS)
- **Prometheus + Grafana** (monitoreo)

### Pagos
- **Stripe** (primary PSP)
- **PayPal** (Fase 2)
- Procesador local CR (Fase 2)

## 📦 Instalación y Setup

### Prerrequisitos

- Go 1.22+
- Node.js 20 LTS+
- PostgreSQL 15+
- Redis 7+
- Docker + Docker Compose

### Backend

```bash
cd backend

# Instalar dependencias
go mod download

# Copiar variables de entorno
cp .env.example .env
# Editar .env con tus credenciales

# Ejecutar migraciones
make migrate-up

# Ejecutar servidor de desarrollo
make run
# API disponible en http://localhost:8080
```

### Frontend

```bash
cd frontend

# Instalar dependencias
npm install

# Copiar variables de entorno
cp .env.example .env
# Editar .env con tus credenciales

# Ejecutar servidor de desarrollo
npm run dev
# App disponible en http://localhost:5173
```

### Docker Compose (Full Stack)

```bash
# Levantar todos los servicios
docker-compose up -d

# Ver logs
docker-compose logs -f

# Detener
docker-compose down
```

## 🧪 Tests

### Backend
```bash
cd backend
make test           # Tests unitarios
make test-coverage  # Con coverage
```

### Frontend
```bash
cd frontend
npm run test        # Vitest
npm run test:ui     # UI de tests
```

## 📊 Monitoreo

- **API Metrics:** http://localhost:8080/metrics (Prometheus)
- **Grafana:** http://localhost:3000 (admin/admin)
- **Logs:** `docker-compose logs -f api`

## 🔐 Seguridad

- Autenticación: JWT (access + refresh tokens)
- Autorización: RBAC (user, admin)
- Rate limiting: Redis (5-60 req/min según endpoint)
- Encriptación: TLS 1.3, bcrypt (passwords), tokens de Stripe (tarjetas)
- Compliance: GDPR, PCI DSS (delegado a Stripe)

Ver [seguridad.md](Documentacion/seguridad.md) para detalles completos.

## 📈 Roadmap

### MVP (8-10 semanas) ✅
- Auth & perfil de usuarios
- Publicación de sorteos
- Reserva y compra de números
- Procesamiento de pagos (Stripe)
- Selección de ganador
- Backoffice básico

### Fase 2 (10-12 semanas)
- Múltiples PSPs (PayPal, local CR)
- Modo "sin cobro" + suscripción
- Búsqueda avanzada y afiliados
- Multilenguaje (ES/EN)
- Comunicación entre usuarios

### Fase 3 (12-16 semanas)
- Aplicación móvil (React Native)
- Dashboards en tiempo real
- Campañas de marketing automatizadas
- Programa de fidelización

Ver [roadmap.md](Documentacion/roadmap.md) para detalles completos.

## 🤝 Contribución

1. Fork el proyecto
2. Crea una rama feature (`git checkout -b feature/AmazingFeature`)
3. Commit tus cambios (`git commit -m 'Add: AmazingFeature'`)
4. Push a la rama (`git push origin feature/AmazingFeature`)
5. Abre un Pull Request

## 📄 Licencia

Todos los derechos de autor e idea de negocio pertenecen a **Ing. Alonso Alpízar**.

## 📧 Contacto

- Email: [contacto@sorteos.com](mailto:contacto@sorteos.com)
- Soporte: [support@sorteos.com](mailto:support@sorteos.com)
- Privacidad: [privacy@sorteos.com](mailto:privacy@sorteos.com)

---

**Versión:** 1.0
**Última actualización:** 2025-11-10
