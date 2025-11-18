# Validación del Sistema Real - Sorteos.club

**Fecha:** 2025-11-18 20:35
**Objetivo:** Documentar cómo funciona el sistema REAL para evitar duplicaciones

---

## ✅ CONFIRMADO: Estructura Real del Sistema

### 1. Usuarios y Roles

**Tabla:** `users` (ÚNICA tabla de usuarios)

**Campo `role`:** ENUM con 3 valores
```sql
user_role = {
  'user',        -- Usuario normal (puede comprar tickets)
  'admin',       -- Administrador
  'super_admin'  -- Super administrador
}
```

**IMPORTANTE:**
- ❌ NO hay tabla separada "admins" o "organizers"
- ✅ TODOS están en la tabla `users`
- ✅ El role determina los permisos

**¿Quién es organizador?**
- Cualquier `user` que cree rifas
- NO es un role, es una ACCIÓN
- Si un user crea rifas → se considera organizador

---

### 2. Organizer Profiles (Perfil Adicional)

**Tabla:** `organizer_profiles` (one-to-one con `users`)

**Relación:**
```sql
organizer_profiles.user_id → users.id (FOREIGN KEY, ON DELETE CASCADE)
```

**Constraint único:**
```sql
UNIQUE (user_id)  -- Un user solo puede tener 1 organizer_profile
```

**¿Cuándo se crea?**
- Cuando un user quiere cobrar por sus rifas
- NO todos los users tienen organizer_profile
- Es OPCIONAL y adicional

**Campos clave:**
- `business_name` - Nombre del negocio
- `tax_id` - RUC/Cédula jurídica
- `bank_account_number` - Cuenta bancaria para pagos
- `commission_override` - Comisión personalizada (NULL = usar default)
- `total_payouts` - Total pagado al organizador
- `pending_payout` - Pendiente de pagar
- `verified` - Aprobado por admin para recibir pagos

**Workflow:**
1. User crea cuenta → solo registro en `users`
2. User crea su primera rifa → sigue siendo solo `users`
3. User completa rifa y quiere cobro → admin crea `organizer_profile`
4. Admin verifica datos bancarios → marca `verified = true`
5. Admin procesa liquidación → actualiza `total_payouts`

---

### 3. Raffles (Rifas del Sistema)

**Tabla:** `raffles` (ÚNICA tabla de rifas)

**Relación:**
```sql
raffles.user_id → users.id  -- Quién creó la rifa
raffles.winner_user_id → users.id  -- Quién ganó
```

**¿Qué hace el admin con raffles?**

**Usuario normal puede:**
- Crear su propia rifa
- Ver sus rifas (`/my-raffles`)
- Editar sus rifas draft
- NO puede ver rifas de otros
- NO puede cambiar status manualmente

**Admin puede (endpoints adicionales):**
- Ver TODAS las rifas (no solo las suyas)
- Suspender cualquier rifa (`suspended_by → admin user_id`)
- Cancelar con refund automático
- Hacer sorteo manual (seleccionar ganador)
- Cambiar status forzadamente (con validaciones)
- Ver timeline completo de transacciones

**Campos admin-only en `raffles`:**
```sql
admin_notes TEXT           -- Notas internas del admin
suspended_by BIGINT        -- FK a users (admin que suspendió)
suspended_at TIMESTAMP     -- Cuándo se suspendió
suspension_reason TEXT     -- Por qué se suspendió
```

**NO se crea tabla nueva** - Se usa la misma tabla `raffles` con permisos diferentes.

---

### 4. Categories (Categorías)

**Tabla:** `categories` (ya existe)

**Endpoints públicos (ya existen):**
- `GET /api/v1/categories` - Listar categorías activas

**Endpoints admin (nuevos):**
- `GET /api/v1/admin/categories` - Listar todas (incluidas inactivas)
- `POST /api/v1/admin/categories` - Crear nueva
- `PUT /api/v1/admin/categories/:id` - Editar
- `DELETE /api/v1/admin/categories/:id` - Soft delete (is_active = false)
- `POST /api/v1/admin/categories/reorder` - Cambiar orden

**¿Qué puede hacer el admin?**
- CRUD completo de categorías
- Activar/desactivar
- Reordenar con drag & drop
- Ver count de rifas por categoría

**NO duplicación:** Se usa la MISMA tabla que el sistema público.

---

### 5. Payments (Pagos del Sistema)

**Tabla:** `payments` (ÚNICA tabla de pagos)

**Usuario normal puede:**
- Ver sus propios pagos
- Hacer checkout

**Admin puede:**
- Ver TODOS los pagos
- Procesar refunds (full o partial)
- Ver detalles completos (webhook events, provider data)
- Marcar disputas
- Ver métricas de conversión

**Workflow de refund:**
1. Admin marca payment para refund
2. Backend llama a Stripe/PayPal API
3. Se actualiza `payment.status = 'refunded'`
4. Se liberan los números reservados
5. Se actualiza `raffle.sold_count` y `revenue`

**NO se duplica:** Misma tabla, diferentes permisos de acceso.

---

### 6. Settlements (Liquidaciones a Organizadores)

**Tabla:** `settlements` (nueva - creada para módulo admin)

**Propósito:** Registrar pagos a organizadores por rifas completadas

**Relación:**
```sql
settlements.organizer_id → users.id
settlements.raffle_id → raffles.id
settlements.approved_by → users.id (admin)
```

**Workflow:**
1. Rifa se completa (winner selected)
2. Sistema calcula: gross_revenue - platform_fee = net_payout
3. Admin crea settlement (manual o auto)
4. Settlement status: pending
5. Admin aprueba → status: approved (verifica KYC + banco)
6. Admin procesa pago → status: paid (actualiza organizer_profile.total_payouts)

**Estados:**
```sql
settlement_status = {
  'pending',    -- Creado, esperando revisión
  'approved',   -- Aprobado, listo para pagar
  'paid',       -- Pagado
  'rejected'    -- Rechazado (problema con datos)
}
```

**Esta tabla SÍ es nueva** porque no existía sistema de liquidaciones antes.

---

### 7. Audit Logs (Logs de Auditoría)

**Tabla:** `audit_logs` (ya existe)

**Propósito:** Registrar TODAS las acciones administrativas

**Relación:**
```sql
audit_logs.user_id → users.id     -- Usuario afectado
audit_logs.admin_id → users.id    -- Admin que hizo la acción
```

**Se registra automáticamente:**
- Suspensión de usuarios
- Cambio de KYC level
- Suspensión de rifas
- Aprobación de settlements
- Refunds procesados
- Cambios en system_parameters

**Severity levels:**
```sql
'info'     -- Acciones de lectura
'warning'  -- Acciones de modificación
'error'    -- Errores en operaciones
'critical' -- Operaciones financieras/sensibles
```

---

### 8. System Parameters (Configuración)

**Tabla:** `system_parameters` (nueva)

**Propósito:** Configuración dinámica del sistema

**Ejemplos de parámetros:**
- `platform_fee_percentage` - % de comisión (default: 10.0)
- `max_active_raffles_per_user` - Límite de rifas activas
- `min_raffle_price` - Precio mínimo de ticket
- `max_raffle_duration_days` - Duración máxima de rifa

**Categorías:**
- Business
- Security
- Payment
- Email
- Notifications

**Admin puede:**
- Ver todos los parámetros
- Editar valores (con validación por tipo)
- Ver historial de cambios (audit_logs)

---

### 9. Notifications (Emails)

**Tabla:** `email_notifications` (nueva)

**Propósito:** Historial de emails enviados por admin

**Tipos de notificaciones:**
- Email individual (a un user específico)
- Email bulk (a múltiples users)
- Announcements (a todos los users activos)

**JSONB fields:**
```sql
recipients JSONB  -- Array de {user_id, email, name}
variables JSONB   -- Variables para template: {userName, raffleTitle, etc}
metadata JSONB    -- Datos adicionales
```

**NO reemplaza:** El sistema de emails transaccionales (registro, reset password)
**SÍ agrega:** Emails administrativos/marketing enviados desde panel admin

---

## 🎯 RESUMEN: ¿Qué es nuevo y qué es integración?

### Tablas que YA EXISTEN (integración)
- ✅ `users` - Se administran, no se duplican
- ✅ `raffles` - Se administran con poderes extra
- ✅ `categories` - CRUD admin sobre tabla existente
- ✅ `payments` - Refunds sobre pagos existentes
- ✅ `audit_logs` - Ya existe, se usa para registrar

### Tablas NUEVAS (creadas para admin)
- ➕ `organizer_profiles` - Perfil adicional one-to-one
- ➕ `settlements` - Liquidaciones a organizadores
- ➕ `system_parameters` - Configuración dinámica
- ➕ `email_notifications` - Historial de emails admin
- ➕ `company_settings` - Info de la empresa
- ➕ `payment_processors` - Config de Stripe/PayPal

### Relaciones Clave

```
users (CENTRAL)
├── organizer_profiles (1:1 opcional)
├── raffles (1:N - rifas creadas)
├── raffle_numbers (N:N vía reservations)
├── payments (1:N - compras)
├── settlements (1:N - liquidaciones como organizador)
├── audit_logs (N:N - como user afectado o admin)
└── kyc_documents (1:N)

raffles
├── user_id → users (creador)
├── winner_user_id → users (ganador)
├── suspended_by → users (admin)
├── category_id → categories
└── settlement → settlements (1:1 cuando completa)

settlements
├── organizer_id → users
├── raffle_id → raffles
└── approved_by → users (admin)
```

---

## ⚠️ ERRORES COMUNES A EVITAR

### ❌ Error 1: Crear tabla "admins" separada
**Incorrecto:**
```sql
CREATE TABLE admins (
  id BIGINT,
  email TEXT,
  ...
)
```

**Correcto:**
```sql
-- Usar tabla users existente
SELECT * FROM users WHERE role IN ('admin', 'super_admin');
```

### ❌ Error 2: Crear tabla "admin_raffles" separada
**Incorrecto:** Tabla separada para rifas administradas

**Correcto:** Usar misma tabla `raffles` con JOIN si es admin:
```sql
-- Admin ve TODAS las rifas
SELECT * FROM raffles;

-- User ve solo las suyas
SELECT * FROM raffles WHERE user_id = $1;
```

### ❌ Error 3: Duplicar endpoints de categories
**Incorrecto:**
- `/api/v1/categories` (público)
- `/api/v1/admin/admin-categories` (admin)

**Correcto:**
- `/api/v1/categories` (GET público - solo activas)
- `/api/v1/admin/categories` (CRUD admin - todas)

### ❌ Error 4: No verificar relaciones existentes
**Incorrecto:** Asumir que organizer_id es un id diferente a user_id

**Correcto:**
```typescript
// organizer_id ES un user_id
interface Settlement {
  organizer_id: number;  // FK a users.id
}
```

---

## ✅ VALIDACIÓN ANTES DE CREAR COMPONENTES

### Checklist de Validación

**Para cada módulo admin, verificar:**

1. **¿La tabla ya existe?**
   - [ ] Conectar a DB: `psql -U postgres sorteos_db`
   - [ ] Ver estructura: `\d nombre_tabla`
   - [ ] Ver relaciones: `\d+ nombre_tabla`

2. **¿Los endpoints ya existen?**
   - [ ] Revisar [API_ENDPOINTS.md](file:///opt/Sorteos/Documentacion/Almighty/API_ENDPOINTS.md)
   - [ ] Probar con curl o script de testing
   - [ ] Verificar request/response format

3. **¿Hay duplicación con sistema público?**
   - [ ] Ver rutas existentes en [App.tsx](file:///opt/Sorteos/frontend/src/App.tsx)
   - [ ] Verificar si ya existe página similar
   - [ ] Identificar qué es NUEVO vs qué es ADMIN-ONLY

4. **¿Las relaciones están claras?**
   - [ ] Identificar FKs en schema
   - [ ] Entender cascadas (ON DELETE CASCADE, etc)
   - [ ] Ver constraints únicos

---

## 📊 MATRIZ DE INTEGRACIÓN

| Módulo Admin | Tabla Principal | ¿Ya existe? | Tipo de Integración |
|--------------|-----------------|-------------|---------------------|
| Users | `users` | ✅ SÍ | Admin powers sobre existente |
| Organizers | `users` + `organizer_profiles` | ✅ users / ➕ profiles | JOIN de 2 tablas |
| Raffles | `raffles` | ✅ SÍ | Admin powers sobre existente |
| Categories | `categories` | ✅ SÍ | CRUD admin sobre existente |
| Payments | `payments` | ✅ SÍ | Refunds sobre existente |
| Settlements | `settlements` | ➕ NUEVA | Funcionalidad nueva |
| Notifications | `email_notifications` | ➕ NUEVA | Historial de emails admin |
| System Config | `system_parameters` | ➕ NUEVA | Configuración dinámica |
| Audit | `audit_logs` | ✅ SÍ | Lectura de logs existentes |
| Reports | Múltiples (aggregations) | ✅ SÍ | Queries sobre existentes |
| Dashboard | Múltiples (KPIs) | ✅ SÍ | Métricas sobre existentes |

---

## 🔍 EJEMPLO CONCRETO: Módulo Users

### Paso 1: Ver tabla real
```sql
\d users
-- Confirmar campos: id, email, role, status, kyc_level, etc.
```

### Paso 2: Ver endpoint admin
```bash
curl -H "Authorization: Bearer $TOKEN" \
  https://mail.sorteos.club/api/v1/admin/users
```

### Paso 3: Entender response
```json
{
  "data": [
    {
      "id": 1,
      "email": "user@example.com",
      "first_name": "John",
      "last_name": "Doe",
      "role": "user",
      "status": "active",
      "kyc_level": "basic",
      "created_at": "2025-01-01T00:00:00Z"
    }
  ],
  "pagination": {...}
}
```

### Paso 4: Crear componente UsersListPage.tsx
```typescript
// Frontend consume endpoint existente
const { data } = useQuery({
  queryKey: ['admin', 'users'],
  queryFn: () => fetch('/api/v1/admin/users', {
    headers: { Authorization: `Bearer ${token}` }
  })
});

// Muestra data de tabla real
<Table>
  {data.data.map(user => (
    <TableRow key={user.id}>
      <TableCell>{user.email}</TableCell>
      <TableCell>{user.role}</TableCell>
      <TableCell>{user.status}</TableCell>
    </TableRow>
  ))}
</Table>
```

**NO se crea:** Nueva tabla, nuevo endpoint, nuevo sistema
**SÍ se crea:** UI para administrar data existente

---

**Documento creado:** 2025-11-18 20:35
**Propósito:** Evitar duplicaciones y entender el sistema real
**Uso:** Leer ANTES de crear cualquier componente frontend admin
