# Estructuras de Respuesta Reales del Backend Admin

**Fecha:** 2025-11-18
**Propósito:** Documentar las estructuras REALES que devuelve el backend para mapearlas correctamente en el frontend

---

## Patrón General de Respuesta

**Todos los handlers devuelven:**
```json
{
  "success": true,
  "data": <OutputStruct>
}
```

**En caso de error:**
```json
{
  "success": false,
  "code": "ERROR_CODE",
  "message": "Error description"
}
```

---

## 1. User Management

### GET /api/v1/admin/users (ListUsers)

**Backend Output Struct:**
```go
type ListUsersOutput struct {
    Users      []*domain.User
    Total      int64
    Page       int
    PageSize   int
    TotalPages int
}
```

**Respuesta HTTP:**
```json
{
  "success": true,
  "data": {
    "Users": [
      {
        "id": 1,
        "uuid": "...",
        "email": "user@example.com",
        "first_name": "John",
        "last_name": "Doe",
        "role": "user",
        "status": "active",
        "kyc_level": "email_verified",
        "created_at": "2025-01-01T00:00:00Z",
        ...
      }
    ],
    "Total": 12,
    "Page": 1,
    "PageSize": 20,
    "TotalPages": 1
  }
}
```

**Mapeo Frontend:**
```typescript
const backendData = response.data.data;
return {
  data: backendData.Users || [],
  pagination: {
    page: backendData.Page,
    limit: backendData.PageSize,
    total: backendData.Total,
    total_pages: backendData.TotalPages,
  }
};
```

### GET /api/v1/admin/users/:id (GetUserDetail)

**Backend Output Struct:**
```go
type GetUserDetailOutput struct {
    User          *domain.User
    RaffleStats   *UserRaffleStats
    PaymentStats  *UserPaymentStats
    RecentAudits  []*AuditLog
}
```

**Respuesta HTTP:**
```json
{
  "success": true,
  "data": {
    "User": { ... },
    "RaffleStats": {
      "total_raffles": 5,
      "active_raffles": 2,
      "completed_raffles": 3,
      "total_revenue": 10000.0
    },
    "PaymentStats": {
      "total_payments": 10,
      "total_spent": 5000.0,
      "refund_count": 1
    },
    "RecentAudits": [ ... ]
  }
}
```

**Mapeo Frontend:**
```typescript
const backendData = response.data.data;
return {
  ...backendData.User,
  raffle_stats: backendData.RaffleStats,
  payment_stats: backendData.PaymentStats,
  recent_audit_logs: backendData.RecentAudits,
};
```

---

## 2. Organizer Management

### GET /api/v1/admin/organizers (ListOrganizers)

**Backend Output Struct:**
```go
type ListOrganizersOutput struct {
    Organizers []*OrganizerListItem
    Total      int64
    Page       int
    PageSize   int
    TotalPages int
}

type OrganizerListItem struct {
    UserID        int64
    Name          string
    Email         string
    BusinessName  string
    Verified      bool
    TotalRaffles  int
    TotalRevenue  float64
    PendingPayout float64
}
```

**Mapeo Frontend:**
```typescript
const backendData = response.data.data;
return {
  data: backendData.Organizers || [],
  pagination: {
    page: backendData.Page,
    limit: backendData.PageSize,
    total: backendData.Total,
    total_pages: backendData.TotalPages,
  }
};
```

---

## 3. Raffle Management

### GET /api/v1/admin/raffles (ListRafflesAdmin)

**Backend Output Struct:**
```go
type ListRafflesAdminOutput struct {
    Raffles    []*RaffleAdminListItem
    Total      int64
    Page       int
    PageSize   int
    TotalPages int
}
```

**Mapeo Frontend:** Mismo patrón que users/organizers.

---

## 4. Category Management

### GET /api/v1/admin/categories (ListCategories)

**Backend Output Struct:**
```go
type ListCategoriesOutput struct {
    Categories []*CategoryListItem
    Page       int
    PageSize   int
    TotalCount int64
    TotalPages int
}
```

**Nota:** Usa `TotalCount` en vez de `Total`.

**Mapeo Frontend:**
```typescript
const backendData = response.data.data;
return {
  data: backendData.Categories || [],
  pagination: {
    page: backendData.Page,
    limit: backendData.PageSize,
    total: backendData.TotalCount, // ← Diferente key
    total_pages: backendData.TotalPages,
  }
};
```

### POST /api/v1/admin/categories (CreateCategory)

**Backend Output Struct:**
```go
type CreateCategoryOutput struct {
    CategoryID  int64
    Name        string
    Description string
    IconURL     string
    IsActive    bool
    CreatedAt   string
    Message     string
}
```

**Respuesta HTTP:**
```json
{
  "success": true,
  "data": {
    "CategoryID": 123,
    "Name": "Electrónica",
    "Description": "...",
    "IconURL": "...",
    "IsActive": true,
    "CreatedAt": "2025-11-18T...",
    "Message": "Categoría creada exitosamente"
  }
}
```

**Mapeo Frontend:**
```typescript
// Retornar directamente response.data.data
return response.data.data;
```

---

## 5. Payment Management

### GET /api/v1/admin/payments (ListPaymentsAdmin)

**Backend Output Struct:**
```go
type ListPaymentsAdminOutput struct {
    Payments   []*PaymentAdminListItem
    Total      int64
    Page       int
    PageSize   int
    TotalPages int
}
```

**Mapeo Frontend:** Mismo patrón estándar.

---

## 6. Settlement Management

### GET /api/v1/admin/settlements (ListSettlements)

**Backend Output Struct:**
```go
type ListSettlementsOutput struct {
    Settlements []*SettlementListItem
    Total       int64
    Page        int
    PageSize    int
    TotalPages  int
}
```

**Mapeo Frontend:** Mismo patrón estándar.

---

## 7. Audit Logs

### GET /api/v1/admin/audit/logs (ListAuditLogs)

**Backend Output Struct:**
```go
type ListAuditLogsOutput struct {
    Logs          []*AuditLog
    Total         int64
    Page          int
    PageSize      int
    TotalPages    int
    // Estadísticas adicionales
    InfoCount     int64
    WarningCount  int64
    ErrorCount    int64
    CriticalCount int64
}
```

**Mapeo Frontend:**
```typescript
const backendData = response.data.data;
return {
  data: backendData.Logs || [],
  pagination: {
    page: backendData.Page,
    limit: backendData.PageSize,
    total: backendData.Total,
    total_pages: backendData.TotalPages,
  },
  stats: {
    info: backendData.InfoCount,
    warning: backendData.WarningCount,
    error: backendData.ErrorCount,
    critical: backendData.CriticalCount,
  }
};
```

---

## 8. System Config

### GET /api/v1/admin/config (ListSystemConfigs)

**Backend Output Struct:**
```go
type ListSystemConfigsOutput struct {
    Configs    []*ConfigListItem
    TotalCount int // ← No usa paginación
}
```

**Nota:** NO tiene paginación.

**Mapeo Frontend:**
```typescript
const backendData = response.data.data;
return {
  data: backendData.Configs || [],
  total: backendData.TotalCount,
};
```

---

## Resumen de Patrones de Mapeo

### Patrón Estándar (Lista con Paginación)
```typescript
const backendData = response.data.data;
return {
  data: backendData.<ArrayField> || [], // Users, Organizers, Raffles, Payments, Settlements, Logs
  pagination: {
    page: backendData.Page,
    limit: backendData.PageSize,
    total: backendData.Total || backendData.TotalCount, // Verificar key
    total_pages: backendData.TotalPages,
  }
};
```

### Patrón Detalle (Sin Paginación)
```typescript
const backendData = response.data.data;
return backendData; // O mapear campos específicos
```

### Patrón Creación/Actualización
```typescript
return response.data.data; // Devuelve el objeto creado/actualizado
```

---

## Checklist de Mapeo por Módulo

| Módulo | Array Key | Total Key | Paginación | Status |
|--------|-----------|-----------|------------|--------|
| Users | `Users` | `Total` | ✅ | ✅ Mapeado |
| Organizers | `Organizers` | `Total` | ✅ | ⏳ Pendiente |
| Raffles | `Raffles` | `Total` | ✅ | ⏳ Pendiente |
| Categories | `Categories` | `TotalCount` | ✅ | ⏳ Pendiente |
| Payments | `Payments` | `Total` | ✅ | ⏳ Pendiente |
| Settlements | `Settlements` | `Total` | ✅ | ⏳ Pendiente |
| Audit Logs | `Logs` | `Total` | ✅ | ⏳ Pendiente |
| System Config | `Configs` | `TotalCount` | ❌ | ⏳ Pendiente |
| Notifications | Varía | - | Varía | ⏳ Pendiente |
| Reports | Varía | - | ❌ | ⏳ Pendiente |

---

**Actualizado:** 2025-11-18 22:00
**Próximo paso:** Implementar mapeo en adminApi.ts para todos los módulos antes de crear las páginas
