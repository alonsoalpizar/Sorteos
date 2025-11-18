# Módulo Almighty Admin - Documentación

**Sistema:** Sorteos.club
**Versión:** 1.0
**Fecha:** 2025-11-18

---

## Qué es el Módulo Almighty

Panel de administración para super-admins que permite:

- ✅ Gestionar usuarios (suspender, cambiar KYC, resetear contraseñas)
- ✅ Gestionar organizadores (perfiles, comisiones personalizadas)
- ✅ Control de rifas (suspender, cancelar con refund, sorteos manuales)
- ✅ Liquidaciones a organizadores (aprobar, pagar)
- ✅ Dashboard con métricas globales
- ✅ Reportes financieros exportables
- ✅ Configuración del sistema (parámetros, categorías, procesadores de pago)

---

## Documentos Principales

### 1. [ROADMAP_ALMIGHTY.md](ROADMAP_ALMIGHTY.md) 📅
**Para qué sirve:** Guía completa de implementación dividida en 8 fases semanales.

**Contiene:**
- Descripción de cada fase (Fundación, Usuarios, Organizadores, etc.)
- Tareas específicas por fase
- Criterios de aceptación
- Estimaciones de tiempo

**Cuándo usarlo:** Para entender el plan completo y orden de implementación.

---

### 2. [CHECKLIST_IMPLEMENTACION.md](CHECKLIST_IMPLEMENTACION.md) ✅
**Para qué sirve:** Lista práctica de tareas para ir marcando día a día.

**Contiene:**
- 217 tareas organizadas por semana
- Checkboxes para marcar progreso
- Espacio para anotar fechas de completado
- Sección de notas y bloqueadores

**Cuándo usarlo:** Trabajo diario - ir marcando tareas completadas.

---

### 3. [BASE_DE_DATOS.md](BASE_DE_DATOS.md) 🗄️
**Para qué sirve:** Referencia técnica de base de datos.

**Contiene:**
- 7 migraciones SQL completas (012-018)
- Diagrama ER de relaciones
- Queries comunes
- Scripts de backfill para datos existentes

**Cuándo usarlo:** Al crear migraciones y consultar esquemas de tablas.

---

### 4. [API_ENDPOINTS.md](API_ENDPOINTS.md) 🔌
**Para qué sirve:** Especificación completa de la API REST.

**Contiene:**
- 52 endpoints documentados
- Request/Response examples
- Query parameters
- Códigos de error

**Cuándo usarlo:** Al implementar handlers y al consumir la API desde frontend.

---

### 5. [ARQUITECTURA_ALMIGHTY.md](ARQUITECTURA_ALMIGHTY.md) 🏗️
**Para qué sirve:** Entender cómo se integra el módulo al sistema existente.

**Contiene:**
- Diagrama de capas (Hexagonal Architecture)
- Flujo de datos (ejemplo: suspender usuario)
- Decisiones arquitectónicas
- Integración con sistema existente

**Cuándo usarlo:** Al diseñar nuevos componentes y entender el flujo general.

---

## Inicio Rápido

### Para el Implementador

1. **Leer primero:** [ROADMAP_ALMIGHTY.md](ROADMAP_ALMIGHTY.md) - Entender las 8 fases
2. **Usar diariamente:** [CHECKLIST_IMPLEMENTACION.md](CHECKLIST_IMPLEMENTACION.md) - Ir marcando tareas
3. **Consultar cuando necesites:**
   - [BASE_DE_DATOS.md](BASE_DE_DATOS.md) - Al crear migraciones
   - [API_ENDPOINTS.md](API_ENDPOINTS.md) - Al implementar endpoints
   - [ARQUITECTURA_ALMIGHTY.md](ARQUITECTURA_ALMIGHTY.md) - Al diseñar componentes

---

## Resumen Técnico

### Stack Tecnológico

**Backend:**
- Go 1.22+ (Gin framework)
- Arquitectura Hexagonal
- PostgreSQL 16
- Redis 7

**Frontend:**
- React 18 + TypeScript
- Vite
- shadcn/ui + Tailwind CSS
- React Query

### Números Clave

- **7 migraciones** de base de datos (012-018)
- **5 tablas nuevas** + 2 tablas modificadas
- **52 endpoints** API REST
- **47 casos de uso** de negocio
- **12 páginas** principales en frontend
- **7-8 semanas** de desarrollo estimado

---

## Flujo de Implementación Sugerido

```
Semana 1: Base de datos (migraciones + repositorios)
         ↓
Semana 2: Gestión de usuarios (backend)
         ↓
Semana 3: Gestión de organizadores (backend)
         ↓
Semana 4: Gestión de rifas y pagos (backend)
         ↓
Semana 5: Liquidaciones (backend)
         ↓
Semana 6: Reportes financieros (backend)
         ↓
Semana 7: Frontend completo + Configuración
         ↓
Semana 8: Testing + Despliegue
```

---

## Primeros Pasos

### 1. Ejecutar Migraciones (Semana 1)

```bash
cd /opt/Sorteos/backend
# Crear archivos de migración según BASE_DE_DATOS.md
migrate create -ext sql -dir migrations -seq create_company_settings
migrate create -ext sql -dir migrations -seq create_payment_processors
migrate create -ext sql -dir migrations -seq create_organizer_profiles
migrate create -ext sql -dir migrations -seq create_settlements
migrate create -ext sql -dir migrations -seq create_system_parameters
migrate create -ext sql -dir migrations -seq add_raffle_admin_fields
migrate create -ext sql -dir migrations -seq add_user_admin_fields

# Ejecutar migraciones
make migrate-up
```

### 2. Crear Estructura de Carpetas

```bash
# Backend - Casos de uso
mkdir -p /opt/Sorteos/backend/internal/usecase/admin/{user,organizer,raffle,payment,settlement,reports,system,category}

# Backend - Handlers
mkdir -p /opt/Sorteos/backend/internal/adapters/http/handler/admin

# Frontend - Admin module
mkdir -p /opt/Sorteos/frontend/src/features/admin/{layout,pages,components}
mkdir -p /opt/Sorteos/frontend/src/hooks
```

### 3. Empezar con User Management (Ejemplo)

```bash
# Crear primer caso de uso
touch /opt/Sorteos/backend/internal/usecase/admin/user/list_users.go

# Crear handler
touch /opt/Sorteos/backend/internal/adapters/http/handler/admin/user_handler.go

# Configurar rutas
# Editar /opt/Sorteos/backend/cmd/api/routes.go
```

---

## Contacto y Soporte

**Proyecto:** Sorteos.club
**Owner:** Tu nombre
**Fecha de inicio:** ____
**Fecha estimada de fin:** ____ (8 semanas después)

---

## Changelog

| Fecha | Versión | Cambios |
|-------|---------|---------|
| 2025-11-18 | 1.0 | Documentación inicial creada |

---

**¡Éxito con la implementación! 🚀**
