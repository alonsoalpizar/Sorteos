# Arquitectura de Navegación y Capas Visuales - Sorteos.club

**Fecha:** 2025-11-11
**Versión:** 1.0

---

## 1. Separación de Capas Visuales

La plataforma se divide en **3 capas principales** con experiencias de usuario distintas:

### 1.1 FRONTOFFICE - Market/Exploración ("Explorar Sorteos")
**Propósito:** Experiencia pública de marketplace atractiva para comprar números

**Características:**
- Landing page atractiva con sorteos destacados
- Catálogo de sorteos públicos del universo sorteos.club
- Filtros y clasificaciones (categorías, precio, fecha, popularidad)
- Vista de detalle de sorteo optimizada para conversión (compra)
- Header público con opción de Login/Register
- **Estado de login visible**: Si usuario está autenticado, mostrar avatar/nombre (no botones Login/Register)
- Call-to-actions claros: "Comprar números", "Ver más sorteos"

**Rutas:**
```
/                           → Landing page (sorteos destacados)
/explorar                   → Catálogo completo con filtros
/sorteos/:id                → Detalle de sorteo (público)
/categorias/:categoria      → Filtrado por categoría
/sorteos-populares          → Vista de trending
```

**Navegación:**
- **Pública**: Navbar simple con logo, búsqueda, categorías, Login/Register
- **Autenticada**: Logo, búsqueda, categorías, avatar → "Mi cuenta" / "Mis compras" / "Dashboard" / Logout

---

### 1.2 BACKOFFICE - Dashboard Usuario ("Mi Panel")
**Propósito:** Panel de control para usuarios que crean y gestionan sus propios sorteos

**Características:**
- Dashboard con resumen de sorteos propios (estados, % vendido, ingresos)
- CRUD de sorteos propios
- Gestión de números vendidos
- Vista de compradores y ganadores
- Perfil y configuración de cuenta
- Historial de compras donde participé
- Historial de pagos y liquidaciones

**Rutas:**
```
/dashboard                  → Vista general (mis sorteos, resumen ventas)
/dashboard/sorteos          → Listado completo de mis sorteos
/dashboard/sorteos/nuevo    → Crear sorteo
/dashboard/sorteos/:id/edit → Editar sorteo
/dashboard/sorteos/:id      → Detalle con gestión (asignar números, ver compradores)
/dashboard/compras          → Sorteos donde compré números
/dashboard/perfil           → Configuración de cuenta
/dashboard/pagos            → Historial de transacciones
/dashboard/liquidaciones    → Dinero recibido por ventas
```

**Navegación:**
- **Sidebar/Menu**: Dashboard, Mis Sorteos, Mis Compras, Perfil, Liquidaciones
- **Header**: Logo, notificaciones, avatar → configuración/logout
- **Botón destacado**: "Volver al Market" → lleva a `/explorar`

---

### 1.3 ADMIN - Backoffice Administrador ("Panel de Administración")
**Propósito:** Herramientas de moderación y gestión de plataforma

**Características:**
- Listado de todos los sorteos (con filtros de estado, modalidad)
- Suspender/activar sorteos
- Gestión de usuarios (verificar KYC, suspender)
- Liquidaciones pendientes
- Audit log de acciones críticas
- Estadísticas de plataforma

**Rutas:**
```
/admin                      → Dashboard admin (métricas, actividad)
/admin/sorteos              → Listado completo sorteos
/admin/usuarios             → Listado usuarios
/admin/liquidaciones        → Pagos pendientes a usuarios
/admin/auditoria            → Logs de acciones
/admin/config               → Configuración de plataforma
```

**Navegación:**
- **Sidebar/Menu**: Dashboard, Sorteos, Usuarios, Liquidaciones, Auditoría, Config
- **Acceso restringido**: Solo roles `admin` o `super_admin`

---

## 2. Problema Actual: Inconsistencia de Estado de Login

### 2.1 Bug Identificado
**Descripción:**
Usuario está autenticado en `/dashboard` pero al ir a `/` (home) aparecen botones "Login" y "Registro" como si no estuviera autenticado.

**Root Cause:**
- Navbar/Header no valida estado de autenticación consistentemente en todas las rutas
- Probablemente dos componentes diferentes: uno para landing y otro para dashboard

### 2.2 Solución Propuesta
**Implementación:**
1. **Un solo componente Navbar** que valide `useAuthStore()` en todas las rutas
2. **Conditional rendering basado en `user`:**
   ```tsx
   {user ? (
     <UserMenu user={user} /> // Avatar + dropdown
   ) : (
     <>
       <Button variant="outline" onClick={() => navigate('/login')}>Login</Button>
       <Button onClick={() => navigate('/register')}>Registro</Button>
     </>
   )}
   ```
3. **Persistencia de sesión:** JWT refresh en localStorage ya implementado

**Ubicación de archivos:**
- `/opt/Sorteos/frontend/src/components/layout/Navbar.tsx` → revisar lógica
- `/opt/Sorteos/frontend/src/components/layout/MainLayout.tsx` → verificar uso consistente

---

## 3. Modalidades de Sorteo (2 Tipos)

### 3.1 Modalidad: **GESTIONADO** (Commission-based)
**Descripción:**
Sorteos con procesamiento de pago automático (PayPal/Stripe). Plataforma cobra comisión.

**Características:**
- Pago inmediato con tarjeta/PayPal
- Comisión de plataforma: 5% (configurable)
- Liquidación automática después del sorteo
- Webhook confirma pago en tiempo real
- Integración completa con payment providers

**Flujo:**
1. Usuario selecciona números → Reserva 5 min
2. Pago con PayPal/Stripe → Confirmación inmediata
3. Número confirmado (status = sold)
4. Plataforma retiene comisión → Liquida a creador después del draw_date

**Base de datos:**
```sql
raffle.management_type = 'managed' -- ENUM: managed, self_managed
raffle.platform_fee_percentage = 0.05 -- (5%)
```

---

### 3.2 Modalidad: **AUTOGESTIONADO** (Self-managed)
**Descripción:**
Sorteos donde el creador gestiona pagos manualmente (efectivo, SINPE Móvil, transferencia). Sin comisión, pago por suscripción mensual.

**Características:**
- **Sin comisión de plataforma** → 0% fee
- Pago por suscripción mensual (ej: $10/mes, sorteos ilimitados)
- Creador confirma pagos manualmente
- Reserva con timeout más largo (24 horas)
- Creador puede asignar números directamente desde backoffice
- Números quedan "pendientes de confirmación" hasta que creador los active

**Flujo:**
1. Usuario selecciona números → Reserva 24 horas (status = reserved_pending_manual)
2. Usuario paga fuera de plataforma (SINPE, efectivo, transferencia)
3. Creador recibe comprobante → Confirma desde backoffice
4. Número confirmado (status = sold)
5. Si timeout de 24h → número liberado automáticamente

**Casos de uso:**
- Sorteos privados entre amigos/familia
- Grupos de WhatsApp/Telegram
- Comunidades pequeñas
- Negocios locales

**Funcionalidad Backoffice (Creador):**
- Vista de números reservados pendientes
- Botón "Confirmar pago" con campo de notas (ej: "SINPE recibido desde 8888-8888")
- Asignar número manualmente a usuario (buscar por email/nombre)
- Liberar reservas manualmente si no pagan

**Base de datos:**
```sql
raffle.management_type = 'self_managed'
raffle.platform_fee_percentage = 0.00 -- Sin comisión
raffle.reservation_ttl_hours = 24 -- Timeout más largo

-- Nuevo status para números
raffle_number.status = 'reserved_pending_manual' -- Esperando confirmación manual
```

**Suscripción:**
```sql
-- Nueva tabla
users.subscription_plan = 'self_managed' -- NULL, 'self_managed'
users.subscription_expires_at = '2025-12-11' -- Fecha de vencimiento
users.subscription_status = 'active' -- active, expired, canceled

-- Validación en create_raffle:
IF raffle.management_type = 'self_managed' THEN
  REQUIRE user.subscription_plan = 'self_managed'
    AND user.subscription_status = 'active'
    AND user.subscription_expires_at > NOW()
END IF
```

---

## 4. Clasificaciones y Taxonomía (Fase 2)

**Propósito:** Organizar sorteos por categorías para mejorar descubrimiento

**Categorías sugeridas:**
- **Electrónica:** iPhones, laptops, consolas, tablets
- **Vehículos:** Carros, motos, bicicletas
- **Viajes:** Paquetes turísticos, boletos aéreos
- **Dinero en efectivo:** Jackpots, premios monetarios
- **Experiencias:** Conciertos, cenas, eventos deportivos
- **Hogar:** Electrodomésticos, muebles
- **Otros**

**Implementación:**
```sql
-- Nueva tabla
CREATE TABLE raffle_categories (
  id BIGSERIAL PRIMARY KEY,
  name VARCHAR(100) NOT NULL UNIQUE,
  slug VARCHAR(100) NOT NULL UNIQUE,
  icon VARCHAR(50), -- emoji o icon name
  display_order INT DEFAULT 0
);

-- Relación
ALTER TABLE raffles ADD COLUMN category_id BIGINT REFERENCES raffle_categories(id);
```

**Frontend:**
- Filtro por categoría en `/explorar`
- Landing page con secciones por categoría
- Iconos visuales para cada categoría

---

## 5. Plan de Implementación

### 5.1 Fase Inmediata (Sprint 6.5 - Fix navegación)
- [ ] Fix: Navbar consistente en todas las rutas (mostrar estado de login)
- [ ] Reorganizar rutas: separar claramente `/explorar` (frontoffice) de `/dashboard` (backoffice)
- [ ] Mejorar landing page `/` con CTAs claros
- [ ] Documentar arquitectura de navegación

### 5.2 Fase 2 (Sprint 7-8 - Modalidad Autogestionada)
- [ ] Backend: Agregar campo `management_type` ENUM a `raffles`
- [ ] Backend: Crear tabla `user_subscriptions`
- [ ] Backend: Validación de suscripción en `create_raffle`
- [ ] Backend: Endpoint para confirmar pagos manuales: `PATCH /dashboard/sorteos/:id/numeros/:numero/confirmar`
- [ ] Backend: Endpoint para asignar número manualmente: `POST /dashboard/sorteos/:id/numeros/:numero/asignar`
- [ ] Backend: Timeout de 24h para reservas autogestionadas
- [ ] Frontend: Toggle en create raffle: "Tipo de sorteo" (Gestionado / Autogestionado)
- [ ] Frontend: Backoffice muestra tabla de reservas pendientes
- [ ] Frontend: Botón "Confirmar pago" con modal de notas
- [ ] Frontend: Buscador de usuarios para asignación manual

### 5.3 Fase 3 (Sprint 9-10 - Clasificaciones)
- [ ] Backend: Tabla `raffle_categories`
- [ ] Backend: Seed con categorías iniciales
- [ ] Backend: Filtrado por categoría en listado
- [ ] Frontend: Selector de categoría en create/edit raffle
- [ ] Frontend: Filtros por categoría en `/explorar`
- [ ] Frontend: Landing page con secciones por categoría

---

## 6. Estructura de Archivos Propuesta

### 6.1 Frontend - Separación de Concerns
```
frontend/src/
├── features/
│   ├── landing/              # Frontoffice - Landing/Market
│   │   ├── pages/
│   │   │   ├── LandingPage.tsx
│   │   │   └── ExplorePage.tsx
│   │   └── components/
│   │       ├── HeroSection.tsx
│   │       ├── FeaturedRaffles.tsx
│   │       └── CategoryFilter.tsx
│   │
│   ├── raffles/              # Frontoffice - Detalle público
│   │   └── pages/
│   │       └── RaffleDetailPage.tsx  # Vista pública (compra)
│   │
│   ├── dashboard/            # Backoffice - Usuario
│   │   ├── pages/
│   │   │   ├── DashboardPage.tsx      # Overview
│   │   │   ├── MyRafflesPage.tsx      # Listado sorteos propios
│   │   │   ├── RaffleManagePage.tsx   # Gestión detallada (asignar números)
│   │   │   ├── MyPurchasesPage.tsx
│   │   │   └── ProfilePage.tsx
│   │   └── components/
│   │       ├── RaffleStats.tsx
│   │       ├── PendingReservations.tsx # Para autogestionados
│   │       └── ManualAssignment.tsx
│   │
│   └── admin/                # Backoffice - Admin
│       ├── pages/
│       │   ├── AdminDashboard.tsx
│       │   ├── AdminRaffles.tsx
│       │   └── AdminUsers.tsx
│       └── components/
│           └── DataTable.tsx
│
├── components/
│   └── layout/
│       ├── Navbar.tsx        # ÚNICO navbar para toda la app
│       ├── PublicLayout.tsx  # Layout para landing/explorar
│       ├── DashboardLayout.tsx # Layout con sidebar para dashboard
│       └── AdminLayout.tsx   # Layout para admin
│
└── App.tsx                   # Routing con layouts
```

### 6.2 Backend - Nuevos Endpoints
```
backend/internal/adapters/http/handler/
├── raffle/
│   ├── confirm_manual_payment_handler.go   # PATCH /dashboard/sorteos/:id/numeros/:numero/confirmar
│   └── assign_number_handler.go            # POST /dashboard/sorteos/:id/numeros/:numero/asignar
│
└── subscription/
    ├── create_subscription_handler.go      # POST /subscriptions (Stripe subscription)
    └── list_subscriptions_handler.go       # GET /subscriptions/me
```

---

## 7. User Stories - Modalidad Autogestionada

### 7.1 Como Creador de Sorteo Autogestionado
```gherkin
Scenario: Crear sorteo autogestionado
  Given soy usuario con suscripción activa
  When creo un sorteo y selecciono "Autogestionado"
  Then el sorteo se crea con management_type = 'self_managed'
  And platform_fee_percentage = 0.00
  And reservation_ttl_hours = 24

Scenario: Confirmar pago manual
  Given tengo un sorteo autogestionado
  And un usuario reservó el número 42 (status = reserved_pending_manual)
  And recibí SINPE Móvil desde 8888-8888
  When voy a Dashboard → Mi Sorteo → Reservas Pendientes
  And hago clic en "Confirmar pago" para número 42
  And escribo notas: "SINPE recibido 8888-8888"
  Then número 42 cambia a status = 'sold'
  And usuario recibe email de confirmación

Scenario: Asignar número manualmente
  Given mi hermano me llamó para reservar el número 15
  When voy a Dashboard → Mi Sorteo → Asignar Número
  And busco usuario por email "hermano@example.com"
  And selecciono número 15
  Then número 15 se asigna a mi hermano (status = sold)
  And él recibe email con confirmación
```

### 7.2 Como Comprador en Sorteo Autogestionado
```gherkin
Scenario: Reservar número en sorteo autogestionado
  Given estoy viendo un sorteo autogestionado
  When selecciono número 42 y hago clic en "Reservar"
  Then número queda reservado por 24 horas
  And veo mensaje: "Tienes 24 horas para confirmar pago con el organizador"
  And veo datos de contacto del organizador
  And veo instrucciones de pago (SINPE, transferencia, etc.)

Scenario: Timeout de reserva
  Given reservé número 42 hace 24 horas
  And no confirmé pago
  When corre el background job de limpieza
  Then mi reserva se cancela automáticamente
  And número 42 vuelve a disponible
```

---

## 8. Preguntas Abiertas / Decisiones Pendientes

### 8.1 Precio de Suscripción Autogestionada
- ¿Cuánto cobrar por mes? (ej: $10/mes USD)
- ¿Limitar cantidad de sorteos activos? (ej: 5 sorteos máx al mismo tiempo)
- ¿Trial gratuito de 7 días?

### 8.2 Verificación de Identidad (KYC)
- ¿Requerir KYC para crear sorteos?
- ¿KYC solo para sorteos gestionados (con pago automático)?
- ¿Límite de monto sin KYC? (ej: hasta $500 sin verificación)

### 8.3 Categorías Iniciales
- ¿Qué categorías priorizar en MVP?
- ¿Permitir sorteos sin categoría al inicio?
- ¿Tags adicionales? (ej: #trending #nuevos #terminan-pronto)

---

## 9. Métricas de Éxito

### 9.1 Navegación
- **Bounce rate** en landing page < 60%
- **Conversion rate** explorar → detalle → compra > 5%
- **Tiempo en dashboard** > 3 minutos (engagement)

### 9.2 Modalidades
- **% sorteos autogestionados** vs gestionados
- **Tasa de confirmación manual** de pagos > 80%
- **Churn rate** suscripciones < 15% mensual

### 9.3 Clasificaciones
- **% usuarios que usan filtro de categoría** > 40%
- **Clicks en categorías** en landing page

---

## 10. Referencias

- **Archivo relacionado:** `/opt/Sorteos/Documentacion/roadmap.md` (Sprint 7-8)
- **Componente Navbar:** `/opt/Sorteos/frontend/src/components/layout/Navbar.tsx`
- **Zustand Auth Store:** `/opt/Sorteos/frontend/src/store/authStore.ts`
- **Backend Routes:** `/opt/Sorteos/backend/cmd/api/raffle_routes.go`

---

**Última actualización:** 2025-11-11
**Autor:** Claude Code
**Estado:** Draft para revisión
