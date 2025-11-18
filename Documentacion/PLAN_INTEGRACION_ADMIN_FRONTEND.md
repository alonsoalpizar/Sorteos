# 🔧 Plan de Integración - Panel Admin al Frontend Existente

**Fecha:** 2025-11-18 19:30
**Objetivo:** Integrar panel admin al frontend React existente con protección por roles

---

## 📋 Situación Actual

### Frontend Existente ✅
- **Framework:** React 18 + TypeScript + Vite
- **Routing:** React Router v6
- **State:** Zustand (con `isAdmin()` ya implementado)
- **API:** TanStack Query configurado
- **UI:** Tailwind CSS
- **Auth:** Sistema completo con roles (user, admin, super_admin)

### Backend Admin API ✅
- **52 endpoints** admin completamente funcionales
- **11 módulos:** Categories, Config, Settlements, Users, Organizers, Payments, Raffles, Notifications, Reports, System, Audit
- **Testing:** 100% verificado

---

## 🎯 Estrategia de Integración

### Opción Recomendada: **Ruta `/admin` con Layout Separado**

**Por qué esta opción:**
1. ✅ Separación clara entre UI pública y admin
2. ✅ Layout independiente (sidebar diferente, sin navbar público)
3. ✅ Fácil de proteger con middleware
4. ✅ No interfiere con rutas existentes
5. ✅ Escalable para futuros módulos

---

## 📐 Arquitectura Propuesta

### Estructura de Carpetas

```
frontend/src/
├── features/
│   ├── admin/                    ← NUEVO
│   │   ├── components/
│   │   │   ├── AdminLayout.tsx        # Layout específico admin
│   │   │   ├── AdminSidebar.tsx       # Sidebar con 11 módulos
│   │   │   ├── AdminHeader.tsx        # Header con user menu
│   │   │   └── AdminRoute.tsx         # HOC para proteger rutas
│   │   ├── pages/
│   │   │   ├── dashboard/
│   │   │   │   └── DashboardPage.tsx  # Dashboard con métricas
│   │   │   ├── users/
│   │   │   │   ├── UsersListPage.tsx
│   │   │   │   └── UserDetailPage.tsx
│   │   │   ├── categories/
│   │   │   │   └── CategoriesPage.tsx
│   │   │   ├── organizers/
│   │   │   │   └── OrganizersPage.tsx
│   │   │   ├── settlements/
│   │   │   ├── payments/
│   │   │   ├── raffles/
│   │   │   ├── notifications/
│   │   │   ├── reports/
│   │   │   ├── system/
│   │   │   └── audit/
│   │   ├── api/
│   │   │   ├── adminClient.ts         # Axios instance para admin
│   │   │   ├── users.ts               # API calls de users
│   │   │   ├── categories.ts
│   │   │   └── ... (uno por módulo)
│   │   ├── hooks/
│   │   │   ├── useAdminUsers.ts       # React Query hooks
│   │   │   ├── useAdminCategories.ts
│   │   │   └── ...
│   │   └── types/
│   │       └── admin.ts               # TypeScript types
│   ├── auth/                     ← EXISTENTE
│   ├── raffles/                  ← EXISTENTE
│   └── ...
└── App.tsx                       ← MODIFICAR (agregar rutas admin)
```

---

## 🔐 Sistema de Protección

### 1. AdminRoute Component

Crear un componente similar a `ProtectedRoute` pero específico para admin:

```tsx
// src/features/admin/components/AdminRoute.tsx
import { Navigate } from 'react-router-dom';
import { useAuthStore } from '@/store/authStore';

export function AdminRoute({ children }: { children: React.ReactNode }) {
  const isAdmin = useAuthStore((state) => state.isAdmin());
  const isAuthenticated = useAuthStore((state) => state.isAuthenticated);

  if (!isAuthenticated) {
    return <Navigate to="/login" replace />;
  }

  if (!isAdmin) {
    // Redirigir a página de "No autorizado" o dashboard normal
    return <Navigate to="/dashboard" replace />;
  }

  return <>{children}</>;
}
```

### 2. Rutas en App.tsx

```tsx
// Agregar a App.tsx (línea ~264)

{/* Admin routes (super protected) */}
<Route
  path="/admin/*"
  element={
    <AdminRoute>
      <AdminLayout>
        <Routes>
          <Route index element={<Navigate to="/admin/dashboard" replace />} />
          <Route path="dashboard" element={<AdminDashboardPage />} />
          <Route path="users" element={<UsersListPage />} />
          <Route path="users/:id" element={<UserDetailPage />} />
          <Route path="categories" element={<CategoriesPage />} />
          <Route path="organizers" element={<OrganizersPage />} />
          <Route path="settlements" element={<SettlementsPage />} />
          <Route path="payments" element={<PaymentsPage />} />
          <Route path="raffles" element={<AdminRafflesPage />} />
          <Route path="notifications" element={<NotificationsPage />} />
          <Route path="reports" element={<ReportsPage />} />
          <Route path="system" element={<SystemPage />} />
          <Route path="audit" element={<AuditPage />} />
        </Routes>
      </AdminLayout>
    </AdminRoute>
  }
/>
```

---

## 🎨 AdminLayout Design

### Características:

1. **Sidebar izquierdo fijo:**
   - Logo + título "Panel Admin"
   - Menú con 11 módulos
   - Iconos + labels
   - Highlight en ruta activa
   - Colapsible en mobile

2. **Header superior:**
   - Título de página actual
   - User menu (nombre + rol + logout)
   - Breadcrumbs opcionales

3. **Área de contenido:**
   - Padding consistente
   - Max-width para legibilidad
   - Scroll independiente

### Colores sugeridos:
```
Sidebar:    bg-slate-900 (dark mode style)
Hover:      bg-slate-800
Active:     bg-blue-600
Text:       text-slate-300
Content:    bg-slate-50
```

---

## 📦 Módulos UI Necesarios

### Componentes shadcn/ui a instalar:

```bash
npx shadcn-ui@latest add button
npx shadcn-ui@latest add card
npx shadcn-ui@latest add table
npx shadcn-ui@latest add dialog
npx shadcn-ui@latest add dropdown-menu
npx shadcn-ui@latest add input
npx shadcn-ui@latest add label
npx shadcn-ui@latest add select
npx shadcn-ui@latest add badge
npx shadcn-ui@latest add alert
npx shadcn-ui@latest add tabs
npx shadcn-ui@latest add separator
```

### Adicionales:
```bash
npm install recharts
npm install date-fns
npm install @tanstack/react-table
```

---

## 🚀 Plan de Implementación (Paso a Paso)

### Fase 1: Setup Base (2-3 horas)

1. **Crear estructura de carpetas**
   ```bash
   mkdir -p src/features/admin/{components,pages,api,hooks,types}
   mkdir -p src/features/admin/pages/{dashboard,users,categories,organizers,settlements,payments,raffles,notifications,reports,system,audit}
   ```

2. **AdminRoute.tsx** - Componente de protección

3. **AdminLayout.tsx** - Layout base con sidebar

4. **AdminSidebar.tsx** - Menú lateral

5. **Agregar rutas a App.tsx**

6. **Probar acceso:** `/admin` debe redirigir a `/admin/dashboard`

---

### Fase 2: Dashboard (4-6 horas)

**Módulo más importante: Visión general del sistema**

**Endpoints a usar:**
- `GET /api/v1/admin/reports/dashboard`

**Componentes:**
1. **MetricCard.tsx** - Card con métrica (Total Users, Total Raffles, Total Revenue, etc.)
2. **RevenueChart.tsx** - Gráfica de ingresos (Recharts)
3. **RecentActivityTable.tsx** - Últimas actividades

**Layout:**
```
+------------------+
| Métricas (4 cards en grid)     |
+------------------+
| Gráfica Revenue  |
+------------------+
| Recent Activity  |
+------------------+
```

---

### Fase 3: Users Management (6-8 horas)

**CRUD completo de usuarios**

**Endpoints:**
- `GET /api/v1/admin/users`
- `GET /api/v1/admin/users/:id`
- `PUT /api/v1/admin/users/:id/status`
- `PUT /api/v1/admin/users/:id/kyc`
- `POST /api/v1/admin/users/:id/reset-password`
- `DELETE /api/v1/admin/users/:id`

**Componentes:**
1. **UsersTable.tsx** - Tabla con paginación + filtros
2. **UserDetailModal.tsx** - Modal con info completa
3. **UserActionsMenu.tsx** - Dropdown con acciones
4. **StatusBadge.tsx** - Badge según status/kyc

**Features:**
- Búsqueda por nombre/email
- Filtros: status, role, kyc_level
- Paginación
- Acciones: Ver detalle, Cambiar status, Actualizar KYC, Reset password, Eliminar

---

### Fase 4: Categories (4-6 horas)

**Gestión de categorías con drag & drop**

**Endpoints:**
- `GET /api/v1/admin/categories`
- `POST /api/v1/admin/categories`
- `PUT /api/v1/admin/categories/:id`
- `POST /api/v1/admin/categories/reorder`
- `DELETE /api/v1/admin/categories/:id`

**Componentes:**
1. **CategoriesTable.tsx** - Con drag & drop (react-beautiful-dnd)
2. **CategoryFormModal.tsx** - Crear/editar
3. **DeleteConfirmDialog.tsx** - Confirmación

**Features:**
- Reordenar con drag & drop
- Toggle active/inactive
- Editar inline o modal
- Ver count de raffles por categoría

---

### Fase 5-10: Resto de Módulos (30-50 horas)

Seguir patrón similar:
1. Crear página
2. Implementar API calls
3. Crear hooks React Query
4. UI con tabla + acciones
5. Testing

---

## 🎨 UI/UX Guidelines

### Consistencia Visual

**Todos los módulos deben tener:**
1. **Header de página:**
   - Título grande
   - Breadcrumbs (opcional)
   - Botón de acción primaria (si aplica)

2. **Área de filtros:**
   - Inputs de búsqueda
   - Selects de filtrado
   - Date pickers
   - Botón "Limpiar filtros"

3. **Tabla/Grid principal:**
   - Paginación estándar
   - Loading states (skeleton)
   - Empty states (cuando no hay data)
   - Error states

4. **Acciones:**
   - Dropdown menu con opciones
   - Modales para editar/crear
   - Confirmaciones para delete

---

## 🔑 Acceso al Panel Admin

### Para Usuarios Normales:
- **NO ven** link a `/admin` en el navbar
- Si intentan acceder manualmente → Redirigidos a `/dashboard`

### Para Admins (role: admin o super_admin):
- **SÍ ven** link "Admin Panel" en dropdown del user menu
- Acceso completo a `/admin/*`

### Modificar MainLayout (navbar):

```tsx
// En src/components/layout/MainLayout.tsx
// Agregar al user dropdown menu:

{isAdmin() && (
  <DropdownMenuItem asChild>
    <Link to="/admin">
      <Shield className="mr-2 h-4 w-4" />
      Panel Admin
    </Link>
  </DropdownMenuItem>
)}
```

---

## 📊 Ventajas de Esta Arquitectura

### ✅ Pros:

1. **Separación clara:** Admin UI separado del resto
2. **Seguro:** Protección en frontend + backend
3. **Escalable:** Fácil agregar nuevos módulos
4. **Mantenible:** Código organizado por feature
5. **Performance:** Lazy loading de rutas admin
6. **UX consistente:** Layout y componentes reutilizables

### ⚠️ Consideraciones:

1. **Duplicación mínima:** Reutilizar componentes UI base (buttons, cards)
2. **API consistency:** Un solo axios instance para admin
3. **Error handling:** Toasts consistentes para errores
4. **Loading states:** Skeleton loaders en tablas

---

## 🎯 Orden de Desarrollo Recomendado

### Semana 1: Setup + Dashboard + Users
```
Día 1-2: Setup base + AdminLayout + Routing (6-8h)
Día 3-4: Dashboard con métricas (4-6h)
Día 5-7: Users management completo (6-8h)
```

### Semana 2: Categories + Organizers + Audit
```
Día 8-9:  Categories (4-6h)
Día 10-11: Organizers (5-7h)
Día 12-13: Audit Logs (3-4h)
```

### Semana 3: Config + Notifications
```
Día 14-16: System + Config (6-8h)
Día 17-19: Notifications (8-10h)
```

### Semana 4: Raffles + Payments
```
Día 20-22: Raffles management (8-10h)
Día 23-25: Payments (6-8h)
```

### Semana 5: Settlements + Polish
```
Día 26-28: Settlements (8-10h)
Día 29-30: Testing final + bug fixes (4-6h)
```

---

## 🧪 Testing Strategy

### Checklist por Módulo:

- [ ] Carga inicial (loading states)
- [ ] Tabla con datos
- [ ] Paginación funciona
- [ ] Filtros funcionan
- [ ] Búsqueda funciona
- [ ] Crear nuevo (si aplica)
- [ ] Editar existente
- [ ] Eliminar con confirmación
- [ ] Error handling (red toast)
- [ ] Success handling (green toast)
- [ ] Responsive en mobile
- [ ] Permisos (solo admin ve)

---

## 🚀 Próximos Pasos Inmediatos

### 1. Setup Inicial (HOY)
```bash
# Crear estructura
cd /opt/Sorteos/frontend/src/features
mkdir -p admin/{components,pages,api,hooks,types}

# Instalar dependencias faltantes (si las hay)
cd /opt/Sorteos/frontend
npm install recharts date-fns @tanstack/react-table
```

### 2. Crear Componentes Base
1. `AdminRoute.tsx` - Protección
2. `AdminLayout.tsx` - Layout principal
3. `AdminSidebar.tsx` - Menú lateral

### 3. Agregar Rutas a App.tsx

### 4. Crear Dashboard Simple
- Fetch data de `/admin/reports/dashboard`
- Mostrar 4 cards con métricas

### 5. Probar en Navegador
- Login como admin
- Navegar a `/admin`
- Ver dashboard con datos reales

---

## 📝 Notas Importantes

### Backend ya está listo ✅
- 52 endpoints funcionales
- CORS configurado
- JWT auth implementado
- Rate limiting activo

### Frontend necesita:
- [ ] AdminLayout con sidebar
- [ ] Rutas protegidas por rol
- [ ] API clients para admin endpoints
- [ ] React Query hooks
- [ ] UI components por módulo

### Tiempo estimado total:
**5 semanas** (65-85 horas) - Desarrollo completo de 11 módulos

---

**Documento creado:** 2025-11-18 19:30
**Próxima acción:** Crear estructura de carpetas y AdminLayout base
**Responsable:** Frontend team
