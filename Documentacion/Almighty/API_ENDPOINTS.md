# API Endpoints - Módulo Almighty Admin

**Versión:** 1.0
**Fecha:** 2025-11-18
**Base URL:** `/api/v1/admin`
**Autenticación:** JWT Bearer Token + Role `super_admin`

---

## 1. Autenticación y Permisos

Todos los endpoints requieren:
- Header: `Authorization: Bearer <JWT_TOKEN>`
- User role: `super_admin`
- Rate limit: 10 requests/minute por admin

```bash
# Ejemplo de request
curl -X GET https://sorteos.club/api/v1/admin/users \
  -H "Authorization: Bearer eyJhbGc..." \
  -H "Content-Type: application/json"
```

---

## 2. Gestión de Usuarios

### 2.1 Listar Usuarios

**`GET /api/v1/admin/users`**

Listar todos los usuarios con filtros y paginación.

**Query Parameters:**
```
role          (string, optional) - Filter by role: user, admin, super_admin
status        (string, optional) - Filter by status: active, suspended, banned, deleted
kyc_level     (string, optional) - Filter by KYC: none, email_verified, phone_verified, cedula_verified, full_kyc
search        (string, optional) - Search in name, email, cedula
offset        (int, default: 0)
limit         (int, default: 20, max: 100)
order_by      (string, default: created_at) - created_at, last_login_at, email
order         (string, default: desc) - asc, desc
```

**Response 200:**
```json
{
  "data": [
    {
      "id": 123,
      "uuid": "550e8400-e29b-41d4-a716-446655440000",
      "email": "juan@example.com",
      "first_name": "Juan",
      "last_name": "Pérez",
      "role": "user",
      "status": "active",
      "kyc_level": "email_verified",
      "created_at": "2025-01-15T10:30:00Z",
      "last_login_at": "2025-11-18T08:15:00Z",
      "total_raffles": 5,
      "total_spent": 15000.00
    }
  ],
  "pagination": {
    "total": 1523,
    "offset": 0,
    "limit": 20,
    "has_more": true
  }
}
```

---

### 2.2 Obtener Detalle de Usuario

**`GET /api/v1/admin/users/:id`**

**Response 200:**
```json
{
  "id": 123,
  "uuid": "550e8400-e29b-41d4-a716-446655440000",
  "email": "juan@example.com",
  "phone": "+50612345678",
  "first_name": "Juan",
  "last_name": "Pérez",
  "cedula": "1-2345-6789",
  "role": "user",
  "status": "active",
  "kyc_level": "email_verified",
  "address": {
    "line1": "Av. Central",
    "city": "San José",
    "state": "San José",
    "postal_code": "10101",
    "country": "CR"
  },
  "max_active_raffles": 10,
  "purchase_limit_daily": 100000.00,
  "suspended_by": null,
  "suspended_at": null,
  "suspension_reason": null,
  "kyc_reviewer": null,
  "last_kyc_review": null,
  "created_at": "2025-01-15T10:30:00Z",
  "updated_at": "2025-11-17T14:20:00Z",
  "last_login_at": "2025-11-18T08:15:00Z",
  "last_login_ip": "192.168.1.100",
  "stats": {
    "total_raffles_created": 5,
    "total_tickets_purchased": 23,
    "total_spent": 15000.00,
    "active_reservations": 0
  },
  "recent_activity": [
    {
      "action": "raffle_created",
      "created_at": "2025-11-17T10:00:00Z",
      "description": "Creó rifa: iPhone 15 Pro"
    }
  ]
}
```

---

### 2.3 Actualizar Estado de Usuario

**`PATCH /api/v1/admin/users/:id/status`**

Suspender, activar o banear usuario.

**Request Body:**
```json
{
  "status": "suspended",  // active, suspended, banned
  "reason": "Violación de términos y condiciones - spam"
}
```

**Response 200:**
```json
{
  "message": "User status updated successfully",
  "user": {
    "id": 123,
    "status": "suspended",
    "suspended_by": 1,
    "suspended_at": "2025-11-18T10:30:00Z",
    "suspension_reason": "Violación de términos y condiciones - spam"
  }
}
```

---

### 2.4 Actualizar Nivel KYC

**`PATCH /api/v1/admin/users/:id/kyc`**

**Request Body:**
```json
{
  "kyc_level": "full_kyc",
  "notes": "Documentos verificados manualmente"
}
```

**Response 200:**
```json
{
  "message": "KYC level updated successfully",
  "user": {
    "id": 123,
    "kyc_level": "full_kyc",
    "kyc_reviewer": 1,
    "last_kyc_review": "2025-11-18T10:35:00Z"
  }
}
```

---

### 2.5 Forzar Reset de Contraseña

**`POST /api/v1/admin/users/:id/reset-password`**

Envía email con link de reset de contraseña.

**Response 200:**
```json
{
  "message": "Password reset email sent successfully"
}
```

---

### 2.6 Eliminar Usuario (Soft Delete)

**`DELETE /api/v1/admin/users/:id`**

**Response 204 No Content**

---

## 3. Gestión de Organizadores

### 3.1 Listar Organizadores

**`GET /api/v1/admin/organizers`**

**Query Parameters:**
```
verified      (boolean, optional) - true, false
min_revenue   (float, optional) - Minimum total revenue
max_revenue   (float, optional) - Maximum total revenue
date_from     (date, optional) - Filter by creation date
date_to       (date, optional)
offset, limit, order_by, order
```

**Response 200:**
```json
{
  "data": [
    {
      "id": 45,
      "user": {
        "id": 123,
        "email": "juan@example.com",
        "first_name": "Juan",
        "last_name": "Pérez"
      },
      "business_name": "Rifas JP",
      "verified": true,
      "total_raffles": 12,
      "active_raffles": 3,
      "completed_raffles": 8,
      "total_revenue": 1250000.00,
      "total_payouts": 980000.00,
      "pending_payout": 85000.00,
      "commission_override": null,
      "created_at": "2025-01-20T12:00:00Z"
    }
  ],
  "pagination": {
    "total": 87,
    "offset": 0,
    "limit": 20
  }
}
```

---

### 3.2 Detalle de Organizador

**`GET /api/v1/admin/organizers/:id`**

**Response 200:**
```json
{
  "id": 45,
  "user_id": 123,
  "business_name": "Rifas JP",
  "tax_id": "3-101-123456",
  "bank_info": {
    "bank_name": "Banco Nacional",
    "account_type": "checking",
    "account_holder": "Juan Pérez Mora",
    "account_number": "****5678"  // masked
  },
  "payout_schedule": "manual",
  "commission_override": null,
  "total_payouts": 980000.00,
  "pending_payout": 85000.00,
  "verified": true,
  "verified_by": 1,
  "verified_at": "2025-02-01T09:00:00Z",
  "raffles_summary": {
    "total": 12,
    "active": 3,
    "completed": 8,
    "cancelled": 1
  },
  "revenue_breakdown": {
    "total_revenue": 1250000.00,
    "platform_fees": 125000.00,
    "net_revenue": 1125000.00
  },
  "recent_raffles": [/* top 5 */]
}
```

---

### 3.3 Actualizar Perfil de Organizador

**`PUT /api/v1/admin/organizers/:id`**

**Request Body:**
```json
{
  "business_name": "Rifas JP SAR",
  "tax_id": "3-101-123456",
  "bank_name": "Banco Nacional",
  "bank_account_number": "15001001012345678",
  "bank_account_type": "checking",
  "bank_account_holder": "Juan Pérez Mora",
  "payout_schedule": "monthly",
  "verified": true
}
```

**Response 200:**
```json
{
  "message": "Organizer profile updated successfully"
}
```

---

### 3.4 Establecer Comisión Personalizada

**`PATCH /api/v1/admin/organizers/:id/commission`**

**Request Body:**
```json
{
  "commission_override": 5.0,  // 5% instead of default 10%
  "reason": "Organizador VIP con alto volumen"
}
```

**Response 200:**
```json
{
  "message": "Commission override set successfully",
  "commission_override": 5.0
}
```

---

### 3.5 Obtener Ingresos de Organizador

**`GET /api/v1/admin/organizers/:id/revenue`**

**Query Parameters:**
```
date_from (date, optional)
date_to   (date, optional)
group_by  (string, optional) - day, month, year
```

**Response 200:**
```json
{
  "total_revenue": 1250000.00,
  "platform_fees": 125000.00,
  "net_revenue": 1125000.00,
  "total_payouts": 980000.00,
  "pending_payout": 85000.00,
  "series": [
    {
      "period": "2025-01",
      "revenue": 450000.00,
      "fees": 45000.00,
      "net": 405000.00
    },
    {
      "period": "2025-02",
      "revenue": 380000.00,
      "fees": 38000.00,
      "net": 342000.00
    }
  ]
}
```

---

## 4. Gestión Avanzada de Rifas

### 4.1 Listar Rifas (Admin View)

**`GET /api/v1/admin/raffles`**

**Query Parameters:**
```
status        (string, optional) - draft, active, suspended, completed, cancelled
organizer_id  (int, optional)
category_id   (int, optional)
search        (string, optional) - Search in title
date_from, date_to, offset, limit
```

**Response 200:**
```json
{
  "data": [
    {
      "id": 567,
      "uuid": "...",
      "title": "iPhone 15 Pro Max",
      "organizer": {
        "id": 123,
        "email": "juan@example.com",
        "name": "Juan Pérez"
      },
      "status": "active",
      "total_numbers": 1000,
      "sold_count": 453,
      "price_per_number": 5000.00,
      "total_revenue": 2265000.00,
      "platform_fee_percentage": 10.0,
      "platform_fee_amount": 226500.00,
      "net_amount": 2038500.00,
      "draw_date": "2025-12-25T20:00:00Z",
      "created_at": "2025-11-01T10:00:00Z",
      "published_at": "2025-11-02T09:00:00Z",
      "suspended_by": null,
      "admin_notes": null
    }
  ],
  "pagination": {...}
}
```

---

### 4.2 Forzar Cambio de Estado

**`PATCH /api/v1/admin/raffles/:id/status`**

**Request Body:**
```json
{
  "status": "suspended",
  "reason": "Imágenes no corresponden al premio anunciado",
  "admin_notes": "Usuario notificado. Requiere actualizar imágenes."
}
```

**Response 200:**
```json
{
  "message": "Raffle status updated successfully",
  "raffle": {
    "id": 567,
    "status": "suspended",
    "suspended_by": 1,
    "suspended_at": "2025-11-18T11:00:00Z",
    "suspension_reason": "Imágenes no corresponden al premio anunciado"
  }
}
```

---

### 4.3 Agregar Notas de Admin

**`POST /api/v1/admin/raffles/:id/notes`**

**Request Body:**
```json
{
  "notes": "Usuario se comprometió a actualizar imágenes en 24h."
}
```

**Response 200:**
```json
{
  "message": "Admin notes added successfully"
}
```

---

### 4.4 Realizar Sorteo Manual

**`POST /api/v1/admin/raffles/:id/draw`**

**Request Body:**
```json
{
  "winner_number": "0543",  // optional, if null will be random
  "notes": "Sorteo realizado manualmente por solicitud del organizador"
}
```

**Response 200:**
```json
{
  "message": "Raffle drawn successfully",
  "winner_number": "0543",
  "winner_user_id": 789,
  "winner_email": "ganador@example.com"
}
```

---

### 4.5 Cancelar Rifa con Reembolsos

**`POST /api/v1/admin/raffles/:id/cancel-refund`**

**Request Body:**
```json
{
  "reason": "Organizador no puede cumplir con la entrega del premio",
  "refund_all": true
}
```

**Response 200:**
```json
{
  "message": "Raffle cancelled and refunds initiated",
  "total_refunds": 453,
  "total_amount_refunded": 2265000.00,
  "refund_status": "processing"  // async process
}
```

---

### 4.6 Ver Transacciones de Rifa

**`GET /api/v1/admin/raffles/:id/transactions`**

**Response 200:**
```json
{
  "raffle": {...},
  "timeline": [
    {
      "type": "raffle_created",
      "timestamp": "2025-11-01T10:00:00Z",
      "actor": "juan@example.com"
    },
    {
      "type": "raffle_published",
      "timestamp": "2025-11-02T09:00:00Z"
    },
    {
      "type": "reservation_created",
      "timestamp": "2025-11-02T10:15:00Z",
      "user": "comprador1@example.com",
      "numbers": ["0001", "0002"]
    },
    {
      "type": "payment_confirmed",
      "timestamp": "2025-11-02T10:18:00Z",
      "amount": 10000.00,
      "payment_id": "pi_abc123"
    }
  ],
  "summary": {
    "total_reservations": 123,
    "total_payments": 453,
    "total_refunds": 0,
    "conversion_rate": 95.2
  }
}
```

---

## 5. Gestión de Pagos

### 5.1 Listar Pagos

**`GET /api/v1/admin/payments`**

**Query Parameters:**
```
status         (string) - pending, succeeded, failed, refunded
user_id        (int)
raffle_id      (int)
payment_method (string) - stripe, paypal
date_from, date_to, offset, limit
```

**Response 200:**
```json
{
  "data": [
    {
      "id": "550e8400-...",
      "amount": 10000.00,
      "currency": "CRC",
      "status": "succeeded",
      "payment_method": "stripe",
      "stripe_payment_intent_id": "pi_abc123",
      "user": {
        "id": 789,
        "email": "comprador@example.com"
      },
      "raffle": {
        "id": 567,
        "title": "iPhone 15 Pro Max"
      },
      "created_at": "2025-11-18T10:00:00Z",
      "paid_at": "2025-11-18T10:01:23Z"
    }
  ],
  "pagination": {...}
}
```

---

### 5.2 Detalle de Pago

**`GET /api/v1/admin/payments/:id`**

**Response 200:**
```json
{
  "id": "550e8400-...",
  "reservation_id": "...",
  "user_id": 789,
  "raffle_id": 567,
  "amount": 10000.00,
  "currency": "CRC",
  "status": "succeeded",
  "payment_method": "stripe",
  "stripe_payment_intent_id": "pi_abc123",
  "stripe_client_secret": "pi_abc123_secret_xyz",
  "metadata": {
    "stripe_response": {...}
  },
  "created_at": "2025-11-18T10:00:00Z",
  "paid_at": "2025-11-18T10:01:23Z",
  "user": {...},
  "raffle": {...},
  "reservation": {...}
}
```

---

### 5.3 Procesar Reembolso

**`POST /api/v1/admin/payments/:id/refund`**

**Request Body:**
```json
{
  "reason": "Rifa cancelada por incumplimiento del organizador",
  "amount": 10000.00  // optional, default: full refund
}
```

**Response 200:**
```json
{
  "message": "Refund processed successfully",
  "refund_id": "re_xyz789",
  "amount_refunded": 10000.00,
  "status": "succeeded"
}
```

---

## 6. Liquidaciones (Settlements)

### 6.1 Listar Liquidaciones

**`GET /api/v1/admin/settlements`**

**Query Parameters:**
```
status        (string) - pending, approved, paid, rejected
organizer_id  (int)
date_from, date_to, offset, limit
```

**Response 200:**
```json
{
  "data": [
    {
      "id": 12,
      "uuid": "...",
      "raffle": {
        "id": 567,
        "title": "iPhone 15 Pro Max"
      },
      "organizer": {
        "id": 123,
        "email": "juan@example.com",
        "name": "Juan Pérez"
      },
      "gross_revenue": 5000000.00,
      "platform_fee": 500000.00,
      "platform_fee_percentage": 10.0,
      "net_payout": 4500000.00,
      "status": "pending",
      "created_at": "2025-11-15T20:05:00Z"
    }
  ],
  "summary": {
    "pending_count": 5,
    "pending_amount": 8500000.00,
    "approved_count": 2,
    "approved_amount": 3200000.00
  }
}
```

---

### 6.2 Crear Liquidación

**`POST /api/v1/admin/settlements`**

**Request Body:**
```json
{
  "raffle_id": 567,
  "notes": "Liquidación manual solicitada por organizador"
}
```

**Response 201:**
```json
{
  "message": "Settlement created successfully",
  "settlement": {
    "id": 13,
    "uuid": "...",
    "gross_revenue": 5000000.00,
    "net_payout": 4500000.00,
    "status": "pending"
  }
}
```

---

### 6.3 Aprobar Liquidación

**`PATCH /api/v1/admin/settlements/:id/approve`**

**Response 200:**
```json
{
  "message": "Settlement approved successfully",
  "settlement": {
    "id": 13,
    "status": "approved",
    "approved_by": 1,
    "approved_at": "2025-11-18T12:00:00Z"
  }
}
```

---

### 6.4 Rechazar Liquidación

**`PATCH /api/v1/admin/settlements/:id/reject`**

**Request Body:**
```json
{
  "reason": "Datos bancarios incorrectos. Por favor actualizar."
}
```

**Response 200:**
```json
{
  "message": "Settlement rejected",
  "settlement": {
    "id": 13,
    "status": "rejected",
    "notes": "Datos bancarios incorrectos. Por favor actualizar."
  }
}
```

---

### 6.5 Marcar como Pagado

**`PATCH /api/v1/admin/settlements/:id/mark-paid`**

**Request Body:**
```json
{
  "payment_method": "bank_transfer",
  "payment_reference": "SINPE-2025111812345678"
}
```

**Response 200:**
```json
{
  "message": "Settlement marked as paid",
  "settlement": {
    "id": 13,
    "status": "paid",
    "payment_method": "bank_transfer",
    "payment_reference": "SINPE-2025111812345678",
    "paid_at": "2025-11-18T12:30:00Z"
  }
}
```

---

## 7. Categorías

### 7.1 Listar Categorías (Admin)

**`GET /api/v1/admin/categories`**

**Response 200:**
```json
{
  "data": [
    {
      "id": 1,
      "name": "Electrónica",
      "slug": "electronica",
      "icon": "📱",
      "description": "Smartphones, laptops, etc.",
      "display_order": 1,
      "is_active": true,
      "raffles_count": 45,
      "created_at": "2025-01-10T00:00:00Z"
    }
  ]
}
```

---

### 7.2 Crear Categoría

**`POST /api/v1/admin/categories`**

**Request Body:**
```json
{
  "name": "Viajes",
  "icon": "✈️",
  "description": "Paquetes turísticos y experiencias"
}
```

**Response 201:**
```json
{
  "message": "Category created successfully",
  "category": {
    "id": 8,
    "slug": "viajes",
    ...
  }
}
```

---

### 7.3 Actualizar Categoría

**`PUT /api/v1/admin/categories/:id`**

**Request Body:**
```json
{
  "name": "Viajes y Turismo",
  "icon": "🌍",
  "description": "Paquetes turísticos, vuelos y experiencias",
  "is_active": true
}
```

**Response 200:**
```json
{
  "message": "Category updated successfully"
}
```

---

### 7.4 Eliminar Categoría

**`DELETE /api/v1/admin/categories/:id`**

**Response 204 No Content**

---

### 7.5 Reordenar Categorías

**`POST /api/v1/admin/categories/reorder`**

**Request Body:**
```json
{
  "order": [1, 3, 2, 8, 4]  // Array of category IDs in new order
}
```

**Response 200:**
```json
{
  "message": "Categories reordered successfully"
}
```

---

## 8. Reportes

### 8.1 Dashboard Global

**`GET /api/v1/admin/reports/dashboard`**

**Response 200:**
```json
{
  "users": {
    "total": 15234,
    "active": 14500,
    "suspended": 234,
    "banned": 50,
    "new_this_week": 123
  },
  "organizers": {
    "total": 456,
    "verified": 234,
    "pending_verification": 222
  },
  "raffles": {
    "total": 1245,
    "active": 345,
    "completed": 789,
    "suspended": 12
  },
  "revenue": {
    "today": 125000.00,
    "this_week": 980000.00,
    "this_month": 4500000.00,
    "this_year": 45000000.00,
    "all_time": 128000000.00
  },
  "platform_fees": {
    "this_month": 450000.00,
    "this_year": 4500000.00
  },
  "settlements": {
    "pending_count": 8,
    "pending_amount": 3200000.00
  },
  "charts": {
    "revenue_last_30_days": [
      {"date": "2025-10-19", "revenue": 45000.00},
      {"date": "2025-10-20", "revenue": 67000.00}
    ],
    "raffles_by_category": [
      {"category": "Electrónica", "count": 234},
      {"category": "Vehículos", "count": 89}
    ]
  }
}
```

---

### 8.2 Reporte de Ingresos

**`GET /api/v1/admin/reports/revenue`**

**Query Parameters:**
```
date_from  (date, required)
date_to    (date, required)
group_by   (string) - day, month, year
category_id (int, optional)
organizer_id (int, optional)
```

**Response 200:**
```json
{
  "total_revenue": 4500000.00,
  "platform_fees": 450000.00,
  "net_revenue": 4050000.00,
  "series": [
    {
      "period": "2025-11-01",
      "revenue": 150000.00,
      "fees": 15000.00
    }
  ]
}
```

---

### 8.3 Reporte de Liquidaciones

**`GET /api/v1/admin/reports/liquidations`**

**Response 200:**
```json
{
  "data": [
    {
      "raffle_id": 567,
      "raffle_title": "iPhone 15 Pro Max",
      "organizer": "Juan Pérez",
      "gross_revenue": 5000000.00,
      "platform_fee": 500000.00,
      "net_payout": 4500000.00,
      "settlement_status": "paid",
      "completed_at": "2025-11-15T20:00:00Z"
    }
  ]
}
```

---

### 8.4 Exportar Reporte

**`POST /api/v1/admin/reports/export`**

**Request Body:**
```json
{
  "report_type": "revenue",  // revenue, liquidations, payouts
  "format": "csv",  // csv, xlsx, pdf
  "date_from": "2025-01-01",
  "date_to": "2025-11-30",
  "filters": {
    "category_id": 1
  }
}
```

**Response 200:**
```json
{
  "message": "Report generated successfully",
  "download_url": "/api/v1/admin/reports/download/abc123xyz",
  "expires_at": "2025-11-18T13:00:00Z"
}
```

---

## 9. Configuración del Sistema

### 9.1 Listar Parámetros

**`GET /api/v1/admin/system/parameters`**

**Query Parameters:**
```
category (string, optional) - business, payment, security, email
```

**Response 200:**
```json
{
  "data": [
    {
      "id": 1,
      "key": "platform_fee_percentage",
      "value": "10.0",
      "value_type": "float",
      "category": "business",
      "description": "Comisión de plataforma por defecto (%)",
      "updated_by": 1,
      "updated_at": "2025-11-01T10:00:00Z"
    }
  ],
  "grouped_by_category": {
    "business": [...],
    "payment": [...],
    "security": [...]
  }
}
```

---

### 9.2 Actualizar Parámetro

**`PUT /api/v1/admin/system/parameters/:key`**

**Request Body:**
```json
{
  "value": "12.0",
  "reason": "Ajuste de comisión por incremento de costos operativos"
}
```

**Response 200:**
```json
{
  "message": "Parameter updated successfully",
  "parameter": {
    "key": "platform_fee_percentage",
    "value": "12.0",
    "updated_by": 1,
    "updated_at": "2025-11-18T13:00:00Z"
  }
}
```

---

### 9.3 Obtener Configuración de Empresa

**`GET /api/v1/admin/system/company`**

**Response 200:**
```json
{
  "company_name": "Sorteos.club",
  "tax_id": "3-101-123456",
  "address": {
    "line1": "...",
    "city": "San José",
    "country": "CR"
  },
  "phone": "+50612345678",
  "email": "info@sorteos.club",
  "support_email": "soporte@sorteos.club",
  "website": "https://sorteos.club",
  "logo_url": "https://sorteos.club/logo.png"
}
```

---

### 9.4 Actualizar Configuración de Empresa

**`PUT /api/v1/admin/system/company`**

**Request Body:**
```json
{
  "company_name": "Sorteos.club S.A.",
  "tax_id": "3-101-123456",
  "phone": "+50612345678",
  ...
}
```

**Response 200:**
```json
{
  "message": "Company settings updated successfully"
}
```

---

### 9.5 Listar Procesadores de Pago

**`GET /api/v1/admin/system/payment-processors`**

**Response 200:**
```json
{
  "data": [
    {
      "id": 1,
      "provider": "stripe",
      "name": "Stripe Production",
      "is_active": true,
      "is_sandbox": false,
      "currency": "CRC",
      "client_id": "pk_live_***", // masked
      "created_at": "2025-01-01T00:00:00Z"
    }
  ]
}
```

---

### 9.6 Actualizar Procesador de Pago

**`PUT /api/v1/admin/system/payment-processors/:id`**

**Request Body:**
```json
{
  "name": "Stripe Production v2",
  "is_active": true,
  "secret_key": "sk_live_...",  // will be encrypted
  "webhook_secret": "whsec_..."  // will be encrypted
}
```

**Response 200:**
```json
{
  "message": "Payment processor updated successfully"
}
```

---

## 10. Auditoría

### 10.1 Listar Logs de Auditoría

**`GET /api/v1/admin/audit`**

**Query Parameters:**
```
action     (string, optional) - user_suspended, settlement_approved, etc.
severity   (string, optional) - info, warning, error, critical
user_id    (int, optional)
admin_id   (int, optional)
date_from, date_to, offset, limit
```

**Response 200:**
```json
{
  "data": [
    {
      "id": 12345,
      "action": "user_suspended",
      "severity": "warning",
      "description": "User suspended for ToS violation",
      "user_id": 789,
      "admin_id": 1,
      "entity_type": "user",
      "entity_id": 789,
      "ip_address": "192.168.1.1",
      "metadata": {
        "status_before": "active",
        "status_after": "suspended",
        "reason": "Spam"
      },
      "created_at": "2025-11-18T10:30:00Z"
    }
  ],
  "pagination": {...}
}
```

---

## 11. Códigos de Error

| Code | Message | Descripción |
|------|---------|-------------|
| 401 | Unauthorized | JWT inválido o expirado |
| 403 | Forbidden | Usuario no tiene rol super_admin |
| 404 | Not Found | Recurso no encontrado |
| 400 | Bad Request | Validación de inputs falló |
| 409 | Conflict | Operación no permitida (ej: admin se suspende a sí mismo) |
| 429 | Too Many Requests | Rate limit excedido |
| 500 | Internal Server Error | Error del servidor |

**Formato de Error:**
```json
{
  "error": {
    "code": "FORBIDDEN",
    "message": "You do not have permission to perform this action",
    "details": {
      "required_role": "super_admin",
      "your_role": "admin"
    }
  }
}
```

---

**Total de Endpoints:** 52
**Última actualización:** 2025-11-18
