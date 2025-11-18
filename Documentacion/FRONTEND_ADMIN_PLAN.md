# Frontend Admin - Plan de Desarrollo con Testing Integrado

## 🎯 Filosofía: Test-First Frontend Development

**Por cada módulo frontend:**
1. ✅ **Test endpoints** - Verificar que funcionen
2. 🔧 **Fix backend** si hay errores
3. 🎨 **Desarrollar UI** - Con confianza de que funciona
4. ✅ **Test integración** - UI + Backend juntos

---

## 📊 Estado Actual del Backend

### Módulos Funcionando (3/11):
- ✅ **Organizers** (200) - 5 endpoints
- ✅ **Reports** (200) - 4 endpoints
- ✅ **Audit** (200) - 1 endpoint

### Módulos con Errores 500 (8/11):
- ❌ **Categories** - Probablemente tabla vacía
- ❌ **Config** - Necesita datos iniciales
- ❌ **Settlements** - Sin liquidaciones
- ❌ **Users** - Posible error en handler
- ❌ **Payments** - Sin pagos en DB
- ❌ **Raffles** - Sin rifas
- ❌ **Notifications** - Error en handler
- ❌ **System** - Error en handler

**Credenciales Admin:**
- Email: `admin@sorteos.com`
- Password: `Admin123456`
- Rol: `super_admin`

---

## 🗺️ Roadmap de Desarrollo Frontend

### Orden Estratégico (de más simple a más complejo):

### **Fase 1: Dashboard & Reports** ⭐ (EMPEZAR AQUÍ)
**Por qué primero:** Endpoints funcionando (200), datos agregados simples

**Módulos:**
1. **Dashboard** (`/admin/reports/dashboard`)
   - Test: `GET /admin/reports/dashboard` ✅
   - UI: Cards con métricas (users, raffles, revenue)
   - Gráficas básicas (Chart.js o Recharts)

2. **Reports** (`/admin/reports/`)
   - Test: `GET /admin/reports/revenue` ✅
   - UI: Filtros por fecha, exportar CSV
   - Tablas de datos

**Estimado:** 4-6 horas

---

### **Fase 2: User Management** 👥
**Por qué segundo:** CRUD básico, 6 endpoints

**Módulo:** Users (`/admin/users`)

**Testing Previo:**
```bash
# Test endpoints
GET /admin/users           # Listar
GET /admin/users/:id       # Ver detalle
PUT /admin/users/:id/status # Cambiar status
PUT /admin/users/:id/kyc    # Actualizar KYC
POST /admin/users/:id/reset-password
DELETE /admin/users/:id
```

**UI:**
- Tabla con búsqueda/filtros
- Modal de detalle de usuario
- Botones de acción (activar/desactivar, KYC, reset password)
- Confirmaciones antes de delete

**Estimado:** 6-8 horas

---

### **Fase 3: Category Management** 📁
**Por qué tercero:** Simple CRUD + reordenamiento

**Módulo:** Categories (`/admin/categories`)

**Testing Previo:**
```bash
# Crear categoría de prueba primero
POST /admin/categories {"name":"Test","description":"Test"}
GET /admin/categories
PUT /admin/categories/:id
POST /admin/categories/reorder
DELETE /admin/categories/:id
```

**UI:**
- Tabla con drag & drop para reordenar
- Modal crear/editar categoría
- Toggle activar/desactivar
- Confirmación antes de eliminar

**Estimado:** 4-6 horas

---

### **Fase 4: Organizer Management** 👔
**Por qué cuarto:** Endpoints funcionando, gestión financiera

**Módulo:** Organizers (`/admin/organizers`)

**Testing Previo:**
```bash
GET /admin/organizers
GET /admin/organizers/:id
PUT /admin/organizers/:id/commission
PUT /admin/organizers/:id/verify
GET /admin/organizers/:id/revenue
```

**UI:**
- Tabla de organizadores
- Detalle con ganancias
- Ajuste de comisión
- Verificación de organizador

**Estimado:** 5-7 horas

---

### **Fase 5: Audit Logs** 📋
**Por qué quinto:** Solo lectura, funcionando

**Módulo:** Audit (`/admin/audit`)

**Testing Previo:**
```bash
GET /admin/audit
GET /admin/audit?action=create&severity=info
```

**UI:**
- Tabla con filtros (admin, action, entity, severity)
- Timeline view opcional
- Búsqueda por texto

**Estimado:** 3-4 horas

---

### **Fase 6: System Configuration** ⚙️
**Por qué sexto:** Configuración crítica

**Módulo:** System (`/admin/system`)

**Testing Previo:**
```bash
# Primero crear datos iniciales en DB
GET /admin/system/parameters
PUT /admin/system/parameters/:key
GET /admin/system/company
PUT /admin/system/company
GET /admin/system/payment-processors
PUT /admin/system/payment-processors/:processor
```

**UI:**
- Formularios de configuración
- Validaciones estrictas
- Confirmación antes de guardar

**Estimado:** 6-8 horas

---

### **Fase 7: Config** 🔧
**Por qué séptimo:** Similar a System

**Módulo:** Config (`/admin/config`)

**Testing + UI similar a System**

**Estimado:** 4-5 horas

---

### **Fase 8: Notifications** 📧
**Por qué octavo:** Envío de emails, más complejo

**Módulo:** Notifications (`/admin/notifications`)

**Testing Previo:**
```bash
POST /admin/notifications/email
POST /admin/notifications/bulk
POST /admin/notifications/templates
POST /admin/notifications/announcements
GET /admin/notifications/history
```

**UI:**
- Formulario de email individual
- Selector de usuarios para bulk
- Editor de templates
- Historial de envíos

**Estimado:** 8-10 horas

---

### **Fase 9: Raffle Management** 🎫
**Por qué noveno:** Gestión compleja de rifas

**Módulo:** Raffles (`/admin/raffles`)

**Testing Previo:**
```bash
GET /admin/raffles
GET /admin/raffles/:id/transactions
PUT /admin/raffles/:id/status
POST /admin/raffles/:id/draw
POST /admin/raffles/:id/notes
POST /admin/raffles/:id/cancel
```

**UI:**
- Tabla de rifas con estados
- Vista de transacciones
- Acciones admin (cancelar, forzar sorteo)
- Sistema de notas

**Estimado:** 8-10 horas

---

### **Fase 10: Payment Management** 💳
**Por qué décimo:** Manejo de dinero, crítico

**Módulo:** Payments (`/admin/payments`)

**Testing Previo:**
```bash
GET /admin/payments
GET /admin/payments/:id
POST /admin/payments/:id/refund
POST /admin/payments/:id/dispute
```

**UI:**
- Tabla de pagos
- Detalle de transacción
- Proceso de reembolso
- Gestión de disputas

**Estimado:** 6-8 horas

---

### **Fase 11: Settlements** 💰
**Por qué último:** Liquidaciones financieras, más complejo

**Módulo:** Settlements (`/admin/settlements`)

**Testing Previo:**
```bash
GET /admin/settlements
GET /admin/settlements/:id
POST /admin/settlements
PUT /admin/settlements/:id/approve
PUT /admin/settlements/:id/reject
PUT /admin/settlements/:id/payout
POST /admin/settlements/auto-create
```

**UI:**
- Workflow de aprobación
- Detalles de liquidación
- Historial de pagos
- Auto-creación masiva

**Estimado:** 8-10 horas

---

## 📦 Stack Tecnológico Frontend

### Core:
- **React 18** + **TypeScript**
- **Vite** (build tool)
- **React Router v6**

### UI Components:
- **shadcn/ui** (componentes base)
- **Tailwind CSS** (estilos)
- **Lucide React** (iconos)

### Data Management:
- **TanStack Query** (React Query v5) - API calls
- **Zustand** (estado global ligero)

### Forms:
- **React Hook Form** + **Zod** (validación)

### Tables:
- **TanStack Table** (tablas avanzadas)

### Charts:
- **Recharts** (gráficas)

### Utils:
- **date-fns** (fechas)
- **axios** (HTTP client)

---

## 🏗️ Estructura de Archivos

```
frontend/src/
├── features/
│   └── admin/
│       ├── components/
│       │   ├── Layout/
│       │   │   ├── AdminLayout.tsx
│       │   │   ├── Sidebar.tsx
│       │   │   └── Header.tsx
│       │   ├── Dashboard/
│       │   │   ├── DashboardPage.tsx
│       │   │   ├── MetricCard.tsx
│       │   │   └── RevenueChart.tsx
│       │   ├── Users/
│       │   │   ├── UsersPage.tsx
│       │   │   ├── UserTable.tsx
│       │   │   ├── UserDetailModal.tsx
│       │   │   └── UserActionsMenu.tsx
│       │   ├── Categories/
│       │   ├── Organizers/
│       │   ├── Audit/
│       │   ├── System/
│       │   ├── Config/
│       │   ├── Notifications/
│       │   ├── Raffles/
│       │   ├── Payments/
│       │   └── Settlements/
│       ├── hooks/
│       │   ├── useDashboard.ts
│       │   ├── useUsers.ts
│       │   └── ...
│       ├── api/
│       │   ├── adminClient.ts
│       │   ├── users.ts
│       │   └── ...
│       └── types/
│           └── admin.ts
├── components/
│   └── ui/  (shadcn components)
├── lib/
│   ├── axios.ts
│   └── queryClient.ts
└── App.tsx
```

---

## ✅ Checklist por Módulo

Para cada módulo, seguir:

- [ ] **Backend Testing**
  - [ ] Test endpoints con curl/Postman
  - [ ] Fix errores 500 si existen
  - [ ] Crear datos de prueba en DB
  - [ ] Documentar comportamiento esperado

- [ ] **API Client**
  - [ ] Crear funciones en `api/[module].ts`
  - [ ] Definir tipos TypeScript
  - [ ] Crear hooks React Query

- [ ] **UI Components**
  - [ ] Layout básico
  - [ ] Tabla/Lista principal
  - [ ] Formularios
  - [ ] Modales
  - [ ] Acciones

- [ ] **Integration Testing**
  - [ ] Test flujo completo en browser
  - [ ] Verificar errores manejados
  - [ ] Test responsivo
  - [ ] Deploy

---

## 🎯 Objetivo Final

**11 módulos admin funcionando al 100%**

- Dashboard con métricas en tiempo real
- Gestión completa de usuarios
- CRUD de categorías
- Gestión de organizadores
- Logs de auditoría
- Configuración del sistema
- Envío de notificaciones
- Gestión de rifas
- Manejo de pagos
- Aprobación de liquidaciones

**Estimado Total:** 65-85 horas de desarrollo

---

**Última actualización:** 2025-11-18
**Estado:** Testing inicial completado (3/11 módulos OK)
**Siguiente paso:** Comenzar Fase 1 - Dashboard & Reports
