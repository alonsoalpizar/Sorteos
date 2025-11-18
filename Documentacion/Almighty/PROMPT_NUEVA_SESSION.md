# Prompt para Nueva Sesión - Integración Frontend Admin

## CONTEXTO CRÍTICO (Leer PRIMERO)

Este es un proyecto de **INTEGRACIÓN**, NO de creación desde cero.

### Sistema Existente (Lo que YA existe y funciona)

**Frontend React:** `/opt/Sorteos/frontend/`
- ✅ React 18 + TypeScript + Vite funcionando
- ✅ Sistema de autenticación completo con roles (user, admin, super_admin)
- ✅ Zustand store con `isAdmin()` implementado
- ✅ React Router v6 configurado en [App.tsx](file:///opt/Sorteos/frontend/src/App.tsx)
- ✅ TanStack Query para API calls
- ✅ Tailwind CSS + shadcn/ui

**Backend API:** `https://mail.sorteos.club/api/v1/`
- ✅ 52 endpoints admin funcionales (`/api/v1/admin/*`)
- ✅ 30 endpoints públicos funcionando
- ✅ Sistema completo de usuarios, rifas, pagos, categorías

**Base de Datos:** PostgreSQL con 19 tablas
- ✅ `users` - TODOS los usuarios (role: user/admin/super_admin)
- ✅ `raffles` - Todas las rifas del sistema
- ✅ `categories` - Categorías de rifas (ya existe CRUD público)
- ✅ `payments` - Pagos del sistema
- ✅ `organizer_profiles` - Perfil adicional para users con role organizer
- ✅ `settlements` - Liquidaciones a organizadores
- ✅ `audit_logs` - Logs de auditoría
- ✅ Y 12 tablas más...

---

## REGLAS FUNDAMENTALES

### 1. NO CREAR DUPLICADOS
❌ **NUNCA** crear nuevas tablas si ya existen
❌ **NUNCA** crear nuevos endpoints si ya existen
❌ **NUNCA** crear un sistema separado

✅ **SIEMPRE** usar las tablas existentes
✅ **SIEMPRE** usar los endpoints existentes
✅ **SIEMPRE** integrar al frontend React existente

### 2. ENTENDER EL SISTEMA REAL

**Pregunta del usuario clave:**
> "Gestión de Organizadores... está pensado con los organizadores actuales del sistema, que al final todos los usuarios son?"

**Respuesta:**
- En la tabla `users`, hay un campo `role` que puede ser: `user`, `admin`, `super_admin`
- Los organizadores son `users` normales que crean rifas
- La tabla `organizer_profiles` es ADICIONAL (one-to-one con `users`)
- NO son entidades separadas, es el MISMO user con perfil adicional

**Otro ejemplo - Raffles:**
- Ya existe tabla `raffles` con TODAS las rifas del sistema
- Ya existe CRUD público de raffles
- El admin simplemente tiene PODERES ADICIONALES:
  - Suspender cualquier rifa
  - Cancelar con refund
  - Cambiar status forzadamente
  - Ver todas las rifas (incluidas de otros)
  - Hacer sorteo manual

**Categories:**
- Ya existe tabla `categories`
- Ya existe endpoint público `GET /api/v1/categories`
- El admin tiene endpoints adicionales:
  - `POST /api/v1/admin/categories` (crear)
  - `PUT /api/v1/admin/categories/:id` (editar)
  - `DELETE /api/v1/admin/categories/:id` (eliminar)

### 3. ARQUITECTURA DE INTEGRACIÓN

```
Sistema Sorteos.club (EXISTENTE)
├── Frontend React (/opt/Sorteos/frontend/)
│   ├── /login                    ← Ya existe
│   ├── /register                 ← Ya existe
│   ├── /raffles                  ← Ya existe (público)
│   ├── /raffles/:id              ← Ya existe (detalle)
│   ├── /my-raffles               ← Ya existe (usuario logueado)
│   ├── /dashboard                ← Ya existe (usuario logueado)
│   └── /admin/*                  ← POR AGREGAR (solo super_admin)
│       ├── /admin/dashboard      ← NUEVO: Ver métricas globales
│       ├── /admin/users          ← NUEVO: Administrar tabla `users`
│       ├── /admin/raffles        ← NUEVO: Administrar tabla `raffles` (misma del sistema)
│       ├── /admin/categories     ← NUEVO: Administrar tabla `categories` (misma)
│       ├── /admin/organizers     ← NUEVO: Ver users + organizer_profiles
│       └── ... (resto de módulos)
│
├── Backend API (mail.sorteos.club/api/v1/)
│   ├── /auth/*                   ← Ya existe
│   ├── /raffles/*                ← Ya existe (público)
│   ├── /categories               ← Ya existe (público)
│   ├── /profile/*                ← Ya existe
│   └── /admin/*                  ← YA FUNCIONAL (52 endpoints)
│       ├── GET /admin/users                    ← Usa tabla `users`
│       ├── GET /admin/raffles                  ← Usa tabla `raffles`
│       ├── PUT /admin/raffles/:id/status       ← Modifica tabla `raffles`
│       ├── GET /admin/organizers               ← JOIN users + organizer_profiles
│       └── ... (48 endpoints más)
│
└── Base de Datos PostgreSQL
    ├── users (ÚNICA tabla de usuarios)
    ├── raffles (ÚNICA tabla de rifas)
    ├── categories (ÚNICA tabla de categorías)
    └── ... (16 tablas más)
```

---

## OBJETIVO DE LA TAREA

**Integrar panel admin al frontend React existente**

**NO hacer:**
- ❌ Crear un frontend separado
- ❌ Crear nuevas tablas
- ❌ Crear nuevos endpoints (ya están)
- ❌ Sistema independiente

**SÍ hacer:**
- ✅ Agregar rutas `/admin/*` a [App.tsx](file:///opt/Sorteos/frontend/src/App.tsx) existente
- ✅ Crear componente `AdminRoute` para proteger por rol
- ✅ Crear `AdminLayout` (sidebar + header específico para admin)
- ✅ Crear páginas que CONSUMAN los endpoints admin existentes
- ✅ Reutilizar componentes UI existentes cuando sea posible
- ✅ Integrar con el auth store existente

---

## ARCHIVOS CLAVE A LEER

### Antes de empezar, leer:

1. **[ROADMAP_ALMIGHTY.md](file:///opt/Sorteos/Documentacion/Almighty/ROADMAP_ALMIGHTY.md)**
   - Plan oficial de 7-8 semanas
   - Estado actual: Backend 100% completo

2. **[API_ENDPOINTS.md](file:///opt/Sorteos/Documentacion/Almighty/API_ENDPOINTS.md)**
   - 52 endpoints admin documentados
   - Request/response examples

3. **[BASE_DE_DATOS.md](file:///opt/Sorteos/Documentacion/Almighty/BASE_DE_DATOS.md)**
   - 19 tablas del sistema
   - Relaciones entre tablas

4. **[App.tsx](file:///opt/Sorteos/frontend/src/App.tsx)**
   - Routing actual del frontend
   - Estructura existente

5. **[authStore.ts](file:///opt/Sorteos/frontend/src/store/authStore.ts)**
   - Sistema de auth existente
   - Método `isAdmin()` ya implementado

### Durante desarrollo, consultar:

- `/opt/Sorteos/Documentacion/Almighty/` - Toda la documentación oficial
- `/opt/Sorteos/backend/internal/usecase/admin/` - Use cases implementados
- `/opt/Sorteos/frontend/src/` - Frontend existente

---

## VALIDACIÓN CRÍTICA (Hacer ANTES de crear componentes)

### Paso 1: Validar tablas existentes

```bash
# Conectar a DB
psql -U sorteos_user -d sorteos_db

# Ver todas las tablas
\dt

# Ver estructura de tabla específica (ejemplo)
\d users
\d raffles
\d organizer_profiles
\d categories
```

**Preguntas a responder:**
- ¿Qué campos tiene la tabla `users`?
- ¿Cómo se relaciona `users` con `organizer_profiles`?
- ¿Qué campos tiene la tabla `raffles`?
- ¿La tabla `categories` ya existe?

### Paso 2: Probar endpoints admin

```bash
# Ejecutar script de testing
bash /opt/Sorteos/Documentacion/Almighty/test_admin_endpoints.sh

# O probar manualmente (ejemplo)
TOKEN="..." # Token de super_admin
curl -H "Authorization: Bearer $TOKEN" \
  https://mail.sorteos.club/api/v1/admin/users
```

**Confirmar:**
- ✅ Endpoint devuelve data de tabla `users` real
- ✅ Los campos coinciden con el schema de DB
- ✅ No hay duplicación de datos

### Paso 3: Verificar frontend existente

```bash
cd /opt/Sorteos/frontend
npm run dev
```

**Navegar a:**
- `/raffles` - Ver que funciona
- `/my-raffles` - Ver que usa la misma tabla
- `/dashboard` - Ver el dashboard actual de usuarios

**Entender:**
- ¿Cómo se consumen las APIs actuales?
- ¿Qué componentes UI existen y se pueden reutilizar?
- ¿Cómo funciona el routing actual?

---

## PLAN DE TRABAJO (Fase 1 - Semana 1)

### Día 1: Setup Base (4 horas)

1. **Crear estructura de carpetas**
   ```bash
   mkdir -p /opt/Sorteos/frontend/src/features/admin/{components,pages,hooks,types}
   mkdir -p /opt/Sorteos/frontend/src/features/admin/pages/{dashboard,users,categories}
   ```

2. **Crear AdminRoute.tsx** (protección por rol)
   - Usar `useAuthStore` existente
   - Verificar `isAdmin()` que YA existe

3. **Crear AdminLayout.tsx** (sidebar + header)
   - Sidebar con 11 módulos
   - Reutilizar componentes UI existentes

4. **Agregar rutas a App.tsx**
   - Agregar `/admin/*` después de línea 263
   - Usar patrón similar a rutas existentes

5. **Crear DashboardPage simple**
   - Solo mostrar "Admin Dashboard" de momento
   - Verificar que carga correctamente

### Día 2-3: Users Management (8 horas)

1. **Crear UsersListPage.tsx**
   - Tabla con data de `GET /api/v1/admin/users`
   - Mostrar: id, name, email, role, status, kyc_level
   - Filtros: role, status, search

2. **Crear UserDetailPage.tsx**
   - Data de `GET /api/v1/admin/users/:id`
   - Mostrar perfil completo
   - Acciones: suspender, cambiar KYC, etc.

3. **Crear hooks**
   - `useAdminUsers()` - TanStack Query hook
   - `useUpdateUserStatus()` - Mutation hook
   - Reutilizar patrón de otros hooks existentes

### Día 4-5: Dashboard con métricas (8 horas)

1. **Completar DashboardPage**
   - Consumir `GET /api/v1/admin/reports/dashboard`
   - KPI cards: total users, raffles, revenue
   - Gráfico de ingresos (Recharts)

2. **Verificar integración**
   - Las métricas deben coincidir con la DB real
   - Los números deben ser consistentes con el sistema

---

## SEGUIMIENTO Y DOCUMENTACIÓN

### Actualizar en cada hito:

1. **[ROADMAP_ALMIGHTY.md](file:///opt/Sorteos/Documentacion/Almighty/ROADMAP_ALMIGHTY.md)**
   - Actualizar % de "Páginas Frontend"
   - Marcar tareas completadas

2. **Crear STATUS updates**
   - `STATUS_FRONTEND_SEMANA_1.md` cuando complete Semana 1
   - Incluir: páginas creadas, endpoints usados, issues encontrados

3. **Git commits descriptivos**
   - Commit por feature completado
   - Ejemplo: "feat(admin): Add users list page with filters"

---

## PREGUNTAS A RESOLVER ANTES DE EMPEZAR

### 1. Sobre Organizadores
**P:** ¿Cómo funcionan los organizadores en el sistema real?
**A:** (Leer tabla `users` + `organizer_profiles` en DB)

**P:** ¿Un user puede ser organizador Y admin?
**A:** (Ver campo `role` en tabla users)

### 2. Sobre Raffles
**P:** ¿Qué puede hacer un admin que un usuario normal no puede?
**A:** (Leer use cases en `/opt/Sorteos/backend/internal/usecase/admin/raffle/`)

**P:** ¿Se crea una rifa nueva o se administra la existente?
**A:** Se administra la existente (tabla `raffles` única)

### 3. Sobre Categories
**P:** ¿Ya existe gestión de categorías?
**A:** (Verificar endpoints públicos vs admin)

---

## CHECKLIST INICIAL (Primera hora de trabajo)

- [ ] Leer ROADMAP_ALMIGHTY.md completo
- [ ] Leer API_ENDPOINTS.md (al menos sección de Users)
- [ ] Conectar a DB y ver estructura de tabla `users`
- [ ] Probar endpoint `GET /api/v1/admin/users` con curl
- [ ] Ver cómo funciona el routing en App.tsx actual
- [ ] Ver cómo funciona el auth store existente
- [ ] Entender que NO se crean tablas nuevas
- [ ] Entender que todo es integración, no separación

---

## COMANDO DE INICIO

```bash
# 1. Ver estado del backend
cd /opt/Sorteos/backend
./sorteos-api  # Debe estar corriendo

# 2. Ver DB
psql -U sorteos_user -d sorteos_db
\dt  # Ver todas las tablas
\d users  # Ver estructura de users
\q

# 3. Probar endpoint admin
bash /opt/Sorteos/Documentacion/Almighty/test_admin_endpoints.sh

# 4. Iniciar frontend dev
cd /opt/Sorteos/frontend
npm run dev

# 5. Leer documentación
cat /opt/Sorteos/Documentacion/Almighty/ROADMAP_ALMIGHTY.md
cat /opt/Sorteos/Documentacion/Almighty/API_ENDPOINTS.md
```

---

## MENSAJE FINAL

Este NO es un proyecto nuevo. Es **integración** de funcionalidades admin al sistema Sorteos.club existente.

**La clave del éxito:** Entender primero qué existe, luego integrar sin duplicar.

**Tiempo estimado:** 7-8 semanas (según ROADMAP_ALMIGHTY.md)
**Fase actual:** Semana 1 - Setup base + Users Management

---

**Creado:** 2025-11-18 20:30
**Para:** Nueva sesión de desarrollo frontend admin
**Contexto:** Backend 100% completo, frontend existente funcional
