# Auditoría: Schema Real vs Código Implementado

**Fecha:** 2025-11-18
**Problema:** Endpoints admin fallan porque el código asume esquemas incorrectos

---

## Tablas en DB (19 tablas reales):

```
audit_logs
categories
company_settings
idempotency_keys
kyc_documents
organizer_profiles
payment_processors
payments
raffle_images
raffle_numbers
raffles
reservations
schema_migrations
settlements
system_parameters
user_consents
users
wallet_transactions
wallets
```

---

## Análisis por Módulo:

### 1. CATEGORIES ❌ (Problemas encontrados)

**Tabla real:** `categories`
```sql
id, name, slug, icon, description, display_order, is_active, created_at, updated_at
```

**Código asume:**
- ❌ Columna `deleted_at` (NO EXISTE)
- ❌ Columna `icon_url` (real es `icon`)
- ❌ En query de raffles: `WHERE deleted_at IS NULL` (tabla raffles tampoco tiene)

**Archivos afectados:**
- `internal/usecase/admin/category/list_categories.go`

---

### 2. CONFIG/SYSTEM_PARAMETERS ❌

**Tabla real:** `system_parameters`

**Código asume:**
- ❌ Tabla llamada `system_config` (NO EXISTE)
- ✅ Cambié a `system_parameters` pero aún faltan validaciones

**Archivos afectados:**
- `internal/usecase/admin/config/*.go` (3 archivos)

---

### 3. SETTLEMENTS ❓ (Requiere verificación)

**Tabla real:** `settlements`

**Verificar esquema:**
