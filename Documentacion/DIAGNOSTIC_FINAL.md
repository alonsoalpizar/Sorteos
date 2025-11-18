# DIAGNÓSTICO COMPLETO: Schema DB vs Código Admin

**Fecha:** 2025-11-18 18:23
**Backups:**
- Schema: `/tmp/schema_backup.sql`
- Full DB: `/tmp/full_backup.sql`

---

## RESUMEN EJECUTIVO

**Estado:** 4/11 módulos funcionando (36%)
- ✅ Users, Organizers, Reports, Audit
- ❌ Categories, Config, Settlements, Payments, Raffles, Notifications, System

**Problema Principal:** Código asume columnas/tablas que NO existen

---

## ANÁLISIS DETALLADO POR MÓDULO

### 1. ❌ CATEGORIES (5 endpoints)

**Tabla real:** `categories`
```
Columnas: id, name, slug, icon, description, display_order, is_active, created_at, updated_at
```

**Problemas encontrados:**
- ✅ Ya NO busca `deleted_at` (limpiado)
- ✅ Ya usa `icon as icon_url` (corregido)
- ⚠️ Query de raffles puede fallar si tabla raffles tiene problemas

**Acción:** PROBAR - debería funcionar ahora

---

### 2. ❌ CONFIG (3 endpoints)

**Tabla real:** `system_parameters`
```
Columnas: id, key, value, value_type, description, category, is_sensitive, updated_by, created_at, updated_at
```

**Problemas encontrados:**
- ✅ Tabla `system_config` → `system_parameters` (ya corregido)
- ❓ Verificar queries en use cases

**Archivos afectados:**
- `internal/usecase/admin/config/get_system_config.go`
- `internal/usecase/admin/config/list_system_configs.go`
- `internal/usecase/admin/config/update_system_config.go`

**Acción:** REVISAR use cases y PROBAR

---

### 3. ❌ SETTLEMENTS (7 endpoints)

**Tabla real:** `settlements`
```
Columnas: id, uuid, raffle_id, organizer_id, gross_revenue, platform_fee,
platform_fee_percentage, net_payout, status, payment_method, payment_reference,
approved_by, approved_at, paid_at, notes, created_at, updated_at
```

**Nota:** status es ENUM `settlement_status` (pending, approved, rejected, paid)

**Problemas potenciales:**
- Verificar si usa `deleted_at`
- Verificar nombres de columnas en queries

**Acción:** REVISAR use cases

---

### 4. ✅ USERS (6 endpoints) - FUNCIONA

**Tabla real:** `users` (20 columnas)
- Sin problemas detectados

---

### 5. ✅ ORGANIZERS (5 endpoints) - FUNCIONA

**Tabla real:** `organizer_profiles`
- Sin problemas detectados

---

### 6. ❌ PAYMENTS (4 endpoints)

**Tabla real:** `payments`
```
Columnas: id (UUID), reservation_id, user_id, raffle_id, stripe_payment_intent_id,
stripe_client_secret, amount, currency, status, payment_method, error_message,
metadata, created_at, updated_at, paid_at
```

**ADVERTENCIA:** IDs son UUID, NO bigint

**Problemas potenciales:**
- Verificar tipos de datos (UUID vs bigint)
- Verificar si usa `deleted_at`

**Acción:** REVISAR use cases

---

### 7. ❌ RAFFLES (6 endpoints)

**Tabla real:** `raffles`
```
Columnas: id, uuid, user_id, title, description, status (ENUM), price_per_number,
total_numbers, min_number, max_number, draw_date, draw_method, winner_number,
winner_user_id, sold_count, reserved_count, total_revenue, platform_fee_percentage,
platform_fee_amount, net_amount, (continúa...)
```

**Nota:** Tabla tiene 20+ columnas, status es ENUM

**Problemas potenciales:**
- Verificar si usa `deleted_at`
- Verificar campos específicos

**Acción:** REVISAR use cases

---

### 8. ❌ NOTIFICATIONS (5 endpoints)

**PROBLEMA CRÍTICO:** No existe tabla `notification_history` o similar en DB

**Tablas que SÍ existen:**
- audit_logs (para logs de sistema)
- ¿Usar audit_logs para notificaciones?
- ¿Crear tabla nueva?

**Acción:** DECIDIR estrategia

---

### 9. ✅ REPORTS (4 endpoints) - FUNCIONA

**Usa:** Agregaciones de raffles, payments, settlements, users
- Sin problemas detectados

---

### 10. ❌ SYSTEM (6 endpoints)

**Tablas reales:**
- `system_parameters`
- `company_settings`
- `payment_processors`

**Todas existen** ✅

**Problemas potenciales:**
- Verificar queries
- Verificar nombres de campos

**Acción:** REVISAR use cases

---

### 11. ✅ AUDIT (1 endpoint) - FUNCIONA

**Tabla real:** `audit_logs`
- Sin problemas detectados

---

## PLAN DE CORRECCIÓN

### Fase 1: Arreglos Rápidos (30-60 min)

1. **Categories** - Probar si ya funciona
2. **Config** - Revisar/arreglar queries
3. **System** - Revisar/arreglar queries

### Fase 2: Revisión Media (1-2 horas)

4. **Settlements** - Revisar use cases, limpiar deleted_at
5. **Payments** - Verificar tipos UUID, limpiar deleted_at
6. **Raffles** - Verificar campos, limpiar deleted_at

### Fase 3: Decisión Estratégica

7. **Notifications** - DECIDIR: ¿Nueva tabla? ¿Usar audit_logs? ¿Simplificar?

---

## SIGUIENTE PASO RECOMENDADO

**OPCIÓN A (Conservadora):**
Arreglar módulo por módulo, testeando cada uno antes de continuar

**OPCIÓN B (Agresiva):**
Script automatizado para limpiar todos los `deleted_at` de golpe

**OPCIÓN C (Pragmática):**
Enfocarse solo en módulos más críticos (Categories, System, Config) primero

---

**¿Cuál prefieres?**
