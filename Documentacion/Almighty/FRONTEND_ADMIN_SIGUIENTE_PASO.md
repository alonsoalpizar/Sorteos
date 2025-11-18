# Frontend Admin - Próximo Paso Concreto

**Fecha:** 2025-11-18 20:00
**Sesión:** Continuación después de completar backend admin 100%

---

## Estado Actual CONFIRMADO ✅

### Backend
- ✅ **52 endpoints admin funcionales** (11/11 módulos working)
- ✅ Todas las rutas en `/api/v1/admin/*`
- ✅ Middleware de auth con `RequireRole("admin", "super_admin")`
- ✅ Probado y funcionando en producción

### Frontend Existente
- ✅ React 18 + TypeScript + Vite
- ✅ React Router v6 configurado en [App.tsx](file:///opt/Sorteos/frontend/src/App.tsx)
- ✅ Zustand auth store con `isAdmin()` implementado
- ✅ TanStack Query configurado
- ✅ MainLayout existente
- ✅ Tailwind CSS + shadcn/ui

### Lo que NO existe aún
- ❌ Rutas `/admin/*` en App.tsx
- ❌ AdminRoute HOC para proteger rutas admin
- ❌ AdminLayout (sidebar + header para panel admin)
- ❌ Páginas del panel admin (Dashboard, Users, etc.)
- ❌ API hooks para endpoints admin

---

## Próximo Paso INMEDIATO

**Objetivo:** Agregar infraestructura base para panel admin integrado

### Paso 1: Crear AdminRoute.tsx (15 min)

**Ubicación:** `/opt/Sorteos/frontend/src/features/admin/components/AdminRoute.tsx`

```tsx
import { Navigate } from 'react-router-dom';
import { useAuthStore } from '@/store/authStore';

export function AdminRoute({ children }: { children: React.ReactNode }) {
  const isAdmin = useAuthStore((state) => state.isAdmin());
  const isAuthenticated = useAuthStore((state) => state.isAuthenticated);

  if (!isAuthenticated) {
    return <Navigate to="/login" replace />;
  }

  if (!isAdmin) {
    return <Navigate to="/dashboard" replace />;
  }

  return <>{children}</>;
}
```

### Paso 2: Crear AdminLayout.tsx (30 min)

**Ubicación:** `/opt/Sorteos/frontend/src/features/admin/components/AdminLayout.tsx`

Layout básico con:
- Sidebar izquierdo con menú de 11 módulos
- Header con user menu
- Área de contenido principal

### Paso 3: Crear DashboardPage.tsx simple (15 min)

**Ubicación:** `/opt/Sorteos/frontend/src/features/admin/pages/dashboard/DashboardPage.tsx`

Página inicial simple que solo muestre "Admin Dashboard" (de momento).

### Paso 4: Agregar rutas a App.tsx (10 min)

Agregar después de la línea 263 (antes del 404 redirect):

```tsx
{/* Admin routes */}
<Route
  path="/admin/*"
  element={
    <AdminRoute>
      <AdminLayout>
        <Routes>
          <Route index element={<Navigate to="/admin/dashboard" replace />} />
          <Route path="dashboard" element={<AdminDashboardPage />} />
        </Routes>
      </AdminLayout>
    </AdminRoute>
  }
/>
```

### Paso 5: Verificar en navegador (5 min)

1. Iniciar dev server: `npm run dev`
2. Login como admin
3. Navegar a `/admin`
4. Debe mostrar dashboard básico

**Tiempo total:** ~75 minutos (1.5 horas)

---

## Después de la Base

Una vez confirmado que `/admin` funciona, continuar con:

1. **Semana 1-2:** Dashboard + Users Management
2. **Semana 3:** Categories + Organizers
3. **Semana 4:** Raffles + Payments
4. **Semana 5:** Settlements + Reports
5. **Semana 6-7:** Notifications + System Config + Audit

**Tiempo total estimado:** 7-8 semanas (según ROADMAP_ALMIGHTY.md)

---

## Comandos Útiles

```bash
# Frontend dev
cd /opt/Sorteos/frontend
npm run dev

# Backend (ya corriendo)
# API en: https://mail.sorteos.club/api/v1

# Testing endpoints admin
bash /opt/Sorteos/Documentacion/Almighty/test_admin_endpoints.sh
```

---

## Referencias Clave

- Backend admin endpoints: [API_ENDPOINTS.md](file:///opt/Sorteos/Documentacion/Almighty/API_ENDPOINTS.md)
- Roadmap oficial: [ROADMAP_ALMIGHTY.md](file:///opt/Sorteos/Documentacion/Almighty/ROADMAP_ALMIGHTY.md)
- App.tsx actual: [App.tsx:1](file:///opt/Sorteos/frontend/src/App.tsx#L1)
- Auth store: [authStore.ts](file:///opt/Sorteos/frontend/src/store/authStore.ts)

---

**Próxima sesión debe empezar:** Creando AdminRoute.tsx
