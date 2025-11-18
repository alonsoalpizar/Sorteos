# Diagnóstico Completo: Schema DB Real vs Código Implementado

**Fecha:** 2025-11-18 18:22
**Objetivo:** Identificar TODAS las discrepancias entre schema de DB y código

---

## Metodología:

1. ✅ Backup completo de DB en `/tmp/schema_backup.sql` y `/tmp/full_backup.sql`
2. ⏳ Extraer columnas reales de cada tabla usada por admin endpoints
3. ⏳ Comparar contra queries en use cases
4. ⏳ Generar plan de corrección

---

## Tablas Involucradas en Endpoints Admin (11 módulos):

### Módulo Categories (5 endpoints)
**Usa tablas:** categories, raffles

### Módulo Config (3 endpoints)
**Usa tablas:** system_parameters

### Módulo Settlements (7 endpoints)
**Usa tablas:** settlements

### Módulo Users (6 endpoints)
**Usa tablas:** users

### Módulo Organizers (5 endpoints)
**Usa tablas:** organizer_profiles, users

### Módulo Payments (4 endpoints)
**Usa tablas:** payments

### Módulo Raffles (6 endpoints)
**Usa tablas:** raffles

### Módulo Notifications (5 endpoints)
**Usa tablas:** ¿notification_history? (NO EXISTE)

### Módulo Reports (4 endpoints)
**Usa tablas:** raffles, payments, settlements, users

### Módulo System (6 endpoints)
**Usa tablas:** system_parameters, company_settings, payment_processors

### Módulo Audit (1 endpoint)
**Usa tablas:** audit_logs

---

## Análisis en Progreso...

