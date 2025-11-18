# Base de Datos - Módulo Almighty Admin

**Versión:** 1.0
**Fecha:** 2025-11-18
**DBMS:** PostgreSQL 16

---

## 1. Resumen de Cambios

### 1.1 Tablas Nuevas: 5

- `company_settings` - Configuración de la empresa
- `payment_processors` - Procesadores de pago configurados
- `organizer_profiles` - Perfiles extendidos de organizadores
- `settlements` - Liquidaciones a organizadores
- `system_parameters` - Parámetros configurables del sistema

### 1.2 Tablas Modificadas: 2

- `users` - Agregar campos de administración
- `raffles` - Agregar campos de suspensión

---

## 2. Migración 012: company_settings

```sql
-- Migration: 012_create_company_settings.up.sql
-- Purpose: Almacenar datos maestros de la empresa Sorteos.club

CREATE TABLE company_settings (
    id BIGSERIAL PRIMARY KEY,

    -- Company Info
    company_name VARCHAR(255) NOT NULL DEFAULT 'Sorteos.club',
    tax_id VARCHAR(50), -- RUC o Tax ID

    -- Address
    address_line1 VARCHAR(255),
    address_line2 VARCHAR(255),
    city VARCHAR(100),
    state VARCHAR(100),
    postal_code VARCHAR(20),
    country VARCHAR(2) DEFAULT 'CR', -- ISO 3166-1 alpha-2

    -- Contact
    phone VARCHAR(20),
    email VARCHAR(255),
    website VARCHAR(255) DEFAULT 'https://sorteos.club',
    support_email VARCHAR(255) DEFAULT 'soporte@sorteos.club',

    -- Branding
    logo_url VARCHAR(500),

    -- Audit
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Trigger for updated_at
CREATE TRIGGER update_company_settings_updated_at
    BEFORE UPDATE ON company_settings
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Insert default data
INSERT INTO company_settings (
    company_name,
    country,
    email,
    support_email,
    website
) VALUES (
    'Sorteos.club',
    'CR',
    'info@sorteos.club',
    'soporte@sorteos.club',
    'https://sorteos.club'
);

-- Only one row allowed (singleton pattern)
CREATE UNIQUE INDEX idx_company_settings_singleton ON company_settings ((id IS NOT NULL));

COMMENT ON TABLE company_settings IS 'Configuración global de la empresa (singleton)';
```

**Rollback:**
```sql
-- Migration: 012_create_company_settings.down.sql
DROP TABLE IF EXISTS company_settings CASCADE;
```

---

## 3. Migración 013: payment_processors

```sql
-- Migration: 013_create_payment_processors.up.sql
-- Purpose: Gestionar credenciales de procesadores de pago

CREATE TABLE payment_processors (
    id BIGSERIAL PRIMARY KEY,

    -- Provider Info
    provider VARCHAR(50) NOT NULL, -- 'stripe', 'paypal', 'credix', etc.
    name VARCHAR(255) NOT NULL, -- 'Stripe Production', 'PayPal Sandbox'

    -- Status
    is_active BOOLEAN DEFAULT true,
    is_sandbox BOOLEAN DEFAULT false,

    -- Credentials (ENCRYPTED)
    client_id TEXT, -- PayPal Client ID, Stripe Publishable Key
    secret_key TEXT, -- Encrypted secret key
    webhook_secret TEXT, -- Encrypted webhook verification secret

    -- Configuration
    currency VARCHAR(3) DEFAULT 'CRC', -- ISO 4217
    config JSONB, -- Additional provider-specific config

    -- Audit
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Indexes
CREATE INDEX idx_payment_processors_provider ON payment_processors(provider);
CREATE INDEX idx_payment_processors_active ON payment_processors(is_active);

-- Trigger for updated_at
CREATE TRIGGER update_payment_processors_updated_at
    BEFORE UPDATE ON payment_processors
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Insert existing Stripe configuration (secrets from ENV)
INSERT INTO payment_processors (
    provider,
    name,
    is_active,
    is_sandbox,
    currency,
    config
) VALUES (
    'stripe',
    'Stripe Production',
    true,
    false,
    'CRC',
    '{}'::jsonb
);

-- Insert existing PayPal configuration
INSERT INTO payment_processors (
    provider,
    name,
    is_active,
    is_sandbox,
    currency,
    config
) VALUES (
    'paypal',
    'PayPal Production',
    false, -- Currently not active, stripe is primary
    false,
    'USD',
    '{}'::jsonb
);

COMMENT ON TABLE payment_processors IS 'Configuración de procesadores de pago';
COMMENT ON COLUMN payment_processors.secret_key IS 'Encrypted with AES-256';
COMMENT ON COLUMN payment_processors.webhook_secret IS 'Encrypted with AES-256';
```

**Rollback:**
```sql
-- Migration: 013_create_payment_processors.down.sql
DROP TABLE IF EXISTS payment_processors CASCADE;
```

---

## 4. Migración 014: organizer_profiles

```sql
-- Migration: 014_create_organizer_profiles.up.sql
-- Purpose: Perfiles extendidos de organizadores con info bancaria y comisiones

CREATE TABLE organizer_profiles (
    id BIGSERIAL PRIMARY KEY,

    -- User Reference
    user_id BIGINT NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,

    -- Business Info
    business_name VARCHAR(255), -- Nombre comercial si es empresa
    tax_id VARCHAR(50), -- RUC o Tax ID de la empresa

    -- Bank Info (ENCRYPTED)
    bank_name VARCHAR(100),
    bank_account_number TEXT, -- Encrypted account number
    bank_account_type VARCHAR(20), -- 'checking', 'savings'
    bank_account_holder VARCHAR(255), -- Nombre del titular

    -- Payout Configuration
    payout_schedule VARCHAR(20) DEFAULT 'manual', -- 'manual', 'weekly', 'monthly'
    commission_override DECIMAL(5,2), -- Custom commission % (NULL = use global default)

    -- Financial Tracking
    total_payouts DECIMAL(12,2) DEFAULT 0.00, -- Total pagado históricamente
    pending_payout DECIMAL(12,2) DEFAULT 0.00, -- Total pendiente de pago

    -- Verification
    verified BOOLEAN DEFAULT false, -- Verificado por admin
    verified_at TIMESTAMP,
    verified_by BIGINT REFERENCES users(id), -- Admin que verificó

    -- Audit
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Indexes
CREATE INDEX idx_organizer_profiles_user_id ON organizer_profiles(user_id);
CREATE INDEX idx_organizer_profiles_verified ON organizer_profiles(verified, created_at DESC);
CREATE INDEX idx_organizer_profiles_total_payouts ON organizer_profiles(total_payouts DESC);

-- Trigger for updated_at
CREATE TRIGGER update_organizer_profiles_updated_at
    BEFORE UPDATE ON organizer_profiles
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

COMMENT ON TABLE organizer_profiles IS 'Perfiles extendidos de organizadores';
COMMENT ON COLUMN organizer_profiles.bank_account_number IS 'Encrypted with AES-256';
COMMENT ON COLUMN organizer_profiles.commission_override IS 'NULL = use global default from system_parameters';
```

**Rollback:**
```sql
-- Migration: 014_create_organizer_profiles.down.sql
DROP TABLE IF EXISTS organizer_profiles CASCADE;
```

---

## 5. Migración 015: settlements

```sql
-- Migration: 015_create_settlements.up.sql
-- Purpose: Registrar liquidaciones y pagos a organizadores

-- Settlement Status ENUM
CREATE TYPE settlement_status AS ENUM (
    'pending',     -- Pendiente de aprobación
    'approved',    -- Aprobado por admin, listo para pagar
    'paid',        -- Pagado al organizador
    'rejected'     -- Rechazado por admin
);

CREATE TABLE settlements (
    id BIGSERIAL PRIMARY KEY,
    uuid UUID UNIQUE DEFAULT uuid_generate_v4(),

    -- References
    raffle_id BIGINT NOT NULL REFERENCES raffles(id) ON DELETE RESTRICT,
    organizer_id BIGINT NOT NULL REFERENCES users(id) ON DELETE RESTRICT,

    -- Amounts
    gross_revenue DECIMAL(12,2) NOT NULL, -- Total vendido
    platform_fee DECIMAL(12,2) NOT NULL, -- Comisión de plataforma
    platform_fee_percentage DECIMAL(5,2) NOT NULL, -- % aplicado
    net_payout DECIMAL(12,2) NOT NULL, -- A pagar al organizador

    -- Status
    status settlement_status DEFAULT 'pending',

    -- Payment Info
    payment_method VARCHAR(50), -- 'bank_transfer', 'paypal', 'sinpe', etc.
    payment_reference VARCHAR(255), -- Número de transferencia, PayPal transaction ID, etc.

    -- Approval
    approved_by BIGINT REFERENCES users(id), -- Admin que aprobó
    approved_at TIMESTAMP,

    -- Payment
    paid_at TIMESTAMP,

    -- Notes
    notes TEXT, -- Notas de admin (ej: razón de rechazo)

    -- Audit
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Indexes
CREATE INDEX idx_settlements_raffle_id ON settlements(raffle_id);
CREATE INDEX idx_settlements_organizer_id ON settlements(organizer_id);
CREATE INDEX idx_settlements_status ON settlements(status, created_at DESC);
CREATE INDEX idx_settlements_approved_by ON settlements(approved_by);

-- Trigger for updated_at
CREATE TRIGGER update_settlements_updated_at
    BEFORE UPDATE ON settlements
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Constraint: raffle can only have one settlement
CREATE UNIQUE INDEX idx_settlements_raffle_unique ON settlements(raffle_id);

COMMENT ON TABLE settlements IS 'Liquidaciones y pagos a organizadores';
COMMENT ON COLUMN settlements.net_payout IS 'gross_revenue - platform_fee';
```

**Rollback:**
```sql
-- Migration: 015_create_settlements.down.sql
DROP TABLE IF EXISTS settlements CASCADE;
DROP TYPE IF EXISTS settlement_status CASCADE;
```

---

## 6. Migración 016: system_parameters

```sql
-- Migration: 016_create_system_parameters.up.sql
-- Purpose: Parámetros de negocio configurables dinámicamente

CREATE TABLE system_parameters (
    id BIGSERIAL PRIMARY KEY,

    -- Parameter Key (unique identifier)
    key VARCHAR(100) UNIQUE NOT NULL,

    -- Value
    value TEXT NOT NULL,
    value_type VARCHAR(20) DEFAULT 'string', -- 'string', 'int', 'float', 'bool', 'json'

    -- Metadata
    description TEXT,
    category VARCHAR(50), -- 'business', 'payment', 'security', 'email', etc.
    is_sensitive BOOLEAN DEFAULT false, -- Si es sensible, no mostrar en logs

    -- Audit
    updated_by BIGINT REFERENCES users(id), -- Último admin que modificó
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Indexes
CREATE INDEX idx_system_parameters_category ON system_parameters(category);
CREATE INDEX idx_system_parameters_updated_by ON system_parameters(updated_by);

-- Trigger for updated_at
CREATE TRIGGER update_system_parameters_updated_at
    BEFORE UPDATE ON system_parameters
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Insert default business parameters
INSERT INTO system_parameters (key, value, value_type, category, description) VALUES
    ('platform_fee_percentage', '10.0', 'float', 'business', 'Comisión de plataforma por defecto (%)'),
    ('max_active_raffles_per_user', '10', 'int', 'business', 'Máximo de rifas activas simultáneas por organizador'),
    ('reservation_ttl_minutes', '5', 'int', 'business', 'Tiempo de expiración de reservas (minutos)'),
    ('min_raffle_numbers', '10', 'int', 'business', 'Mínimo de números en una rifa'),
    ('max_raffle_numbers', '10000', 'int', 'business', 'Máximo de números en una rifa'),
    ('min_price_per_number', '100.0', 'float', 'business', 'Precio mínimo por número (CRC)'),
    ('require_kyc_for_organizers', 'true', 'bool', 'security', 'Requerir email verificado para crear rifas'),
    ('auto_settlement_creation', 'true', 'bool', 'business', 'Crear settlements automáticamente al completar rifa'),
    ('enable_refunds', 'true', 'bool', 'payment', 'Permitir reembolsos de pagos'),
    ('support_email', 'soporte@sorteos.club', 'string', 'email', 'Email de soporte');

COMMENT ON TABLE system_parameters IS 'Parámetros de negocio configurables';
COMMENT ON COLUMN system_parameters.value_type IS 'Tipo de dato: string, int, float, bool, json';
```

**Rollback:**
```sql
-- Migration: 016_create_system_parameters.down.sql
DROP TABLE IF EXISTS system_parameters CASCADE;
```

---

## 7. Migración 017: Modificar tabla raffles

```sql
-- Migration: 017_add_raffle_admin_fields.up.sql
-- Purpose: Agregar campos de administración a tabla raffles

ALTER TABLE raffles
    ADD COLUMN suspension_reason TEXT,
    ADD COLUMN suspended_by BIGINT REFERENCES users(id),
    ADD COLUMN suspended_at TIMESTAMP,
    ADD COLUMN admin_notes TEXT;

-- Index for suspended_by
CREATE INDEX idx_raffles_suspended_by ON raffles(suspended_by) WHERE suspended_by IS NOT NULL;

-- Comments
COMMENT ON COLUMN raffles.suspension_reason IS 'Razón de suspensión (visible para organizador)';
COMMENT ON COLUMN raffles.suspended_by IS 'Admin que suspendió la rifa';
COMMENT ON COLUMN raffles.admin_notes IS 'Notas privadas de admin (no visibles para organizador)';
```

**Rollback:**
```sql
-- Migration: 017_add_raffle_admin_fields.down.sql
ALTER TABLE raffles
    DROP COLUMN IF EXISTS suspension_reason,
    DROP COLUMN IF EXISTS suspended_by,
    DROP COLUMN IF EXISTS suspended_at,
    DROP COLUMN IF EXISTS admin_notes;
```

---

## 8. Migración 018: Modificar tabla users

```sql
-- Migration: 018_add_user_admin_fields.up.sql
-- Purpose: Agregar campos de administración a tabla users

ALTER TABLE users
    ADD COLUMN suspension_reason TEXT,
    ADD COLUMN suspended_by BIGINT REFERENCES users(id),
    ADD COLUMN suspended_at TIMESTAMP,
    ADD COLUMN last_kyc_review TIMESTAMP,
    ADD COLUMN kyc_reviewer BIGINT REFERENCES users(id);

-- Indexes
CREATE INDEX idx_users_suspended_by ON users(suspended_by) WHERE suspended_by IS NOT NULL;
CREATE INDEX idx_users_kyc_reviewer ON users(kyc_reviewer) WHERE kyc_reviewer IS NOT NULL;

-- Comments
COMMENT ON COLUMN users.suspension_reason IS 'Razón de suspensión/ban';
COMMENT ON COLUMN users.suspended_by IS 'Admin que suspendió al usuario';
COMMENT ON COLUMN users.last_kyc_review IS 'Última revisión de KYC por admin';
COMMENT ON COLUMN users.kyc_reviewer IS 'Admin que revisó/cambió KYC';
```

**Rollback:**
```sql
-- Migration: 018_add_user_admin_fields.down.sql
ALTER TABLE users
    DROP COLUMN IF EXISTS suspension_reason,
    DROP COLUMN IF EXISTS suspended_by,
    DROP COLUMN IF EXISTS suspended_at,
    DROP COLUMN IF EXISTS last_kyc_review,
    DROP COLUMN IF EXISTS kyc_reviewer;
```

---

## 9. Diagrama de Relaciones (ER Diagram)

```
┌─────────────────┐
│     users       │
├─────────────────┤
│ id (PK)         │◄───┐
│ email           │    │
│ role            │    │
│ status          │    │
│ suspended_by    │────┘ (self-reference)
│ kyc_level       │
│ kyc_reviewer    │────┐
└─────────────────┘    │
         │             │
         │ 1:1         │
         ▼             │
┌─────────────────────────┐
│ organizer_profiles      │
├─────────────────────────┤
│ id (PK)                 │
│ user_id (FK, UNIQUE)    │
│ business_name           │
│ bank_account_number     │ (encrypted)
│ commission_override     │
│ total_payouts           │
│ pending_payout          │
│ verified                │
│ verified_by (FK)        │───┘
└─────────────────────────┘
         │
         │ 1:N
         ▼
┌─────────────────┐
│   settlements   │
├─────────────────┤
│ id (PK)         │
│ raffle_id (FK)  │───┐
│ organizer_id(FK)│   │
│ gross_revenue   │   │
│ platform_fee    │   │
│ net_payout      │   │
│ status          │   │ 1:1
│ approved_by(FK) │   │
│ paid_at         │   │
└─────────────────┘   │
                      ▼
┌─────────────────┐
│    raffles      │
├─────────────────┤
│ id (PK)         │
│ user_id (FK)    │ (organizador)
│ status          │
│ total_revenue   │
│ platform_fee_%  │
│ suspended_by(FK)│
│ admin_notes     │
└─────────────────┘
         │
         │ N:1
         ▼
┌─────────────────┐
│   categories    │
├─────────────────┤
│ id (PK)         │
│ name            │
│ slug            │
│ is_active       │
│ display_order   │
└─────────────────┘

┌───────────────────┐
│ company_settings  │ (singleton)
├───────────────────┤
│ id (PK)           │
│ company_name      │
│ tax_id            │
│ address           │
│ logo_url          │
└───────────────────┘

┌───────────────────────┐
│ payment_processors    │
├───────────────────────┤
│ id (PK)               │
│ provider              │
│ is_active             │
│ secret_key            │ (encrypted)
│ webhook_secret        │ (encrypted)
└───────────────────────┘

┌───────────────────┐
│ system_parameters │
├───────────────────┤
│ id (PK)           │
│ key (UNIQUE)      │
│ value             │
│ value_type        │
│ category          │
│ updated_by (FK)   │───┐
└───────────────────┘   │
                        │
                        │ N:1
                        ▼
                  ┌─────────┐
                  │  users  │
                  └─────────┘
```

---

## 10. Queries Comunes

### 10.1 Listar organizadores con métricas

```sql
SELECT
    u.id,
    u.email,
    u.first_name,
    u.last_name,
    op.business_name,
    op.verified,
    op.total_payouts,
    op.pending_payout,
    COUNT(DISTINCT r.id) as total_raffles,
    COUNT(DISTINCT CASE WHEN r.status = 'active' THEN r.id END) as active_raffles,
    SUM(CASE WHEN r.status = 'completed' THEN r.total_revenue ELSE 0 END) as total_revenue
FROM users u
LEFT JOIN organizer_profiles op ON u.id = op.user_id
LEFT JOIN raffles r ON u.id = r.user_id AND r.deleted_at IS NULL
WHERE u.role = 'user' -- Usuarios que han creado rifas
GROUP BY u.id, op.id
HAVING COUNT(DISTINCT r.id) > 0
ORDER BY total_revenue DESC
LIMIT 50;
```

### 10.2 Dashboard KPIs

```sql
-- Total users by status
SELECT
    status,
    COUNT(*) as count
FROM users
WHERE deleted_at IS NULL
GROUP BY status;

-- Revenue this month
SELECT
    SUM(total_revenue) as revenue_this_month,
    SUM(platform_fee_amount) as platform_fees_collected
FROM raffles
WHERE status = 'completed'
    AND completed_at >= DATE_TRUNC('month', NOW())
    AND completed_at < DATE_TRUNC('month', NOW() + INTERVAL '1 month');

-- Pending settlements
SELECT
    COUNT(*) as pending_count,
    SUM(net_payout) as pending_amount
FROM settlements
WHERE status = 'pending';
```

### 10.3 Settlements pendientes por organizador

```sql
SELECT
    s.id,
    s.uuid,
    r.title as raffle_title,
    u.email as organizer_email,
    u.first_name || ' ' || u.last_name as organizer_name,
    s.gross_revenue,
    s.platform_fee,
    s.net_payout,
    s.created_at
FROM settlements s
JOIN raffles r ON s.raffle_id = r.id
JOIN users u ON s.organizer_id = u.id
WHERE s.status = 'pending'
ORDER BY s.created_at ASC;
```

### 10.4 Audit logs de acciones críticas

```sql
SELECT
    al.id,
    al.action,
    al.severity,
    al.description,
    u_admin.email as admin_email,
    u_target.email as target_user_email,
    al.created_at
FROM audit_logs al
LEFT JOIN users u_admin ON al.admin_id = u_admin.id
LEFT JOIN users u_target ON al.user_id = u_target.id
WHERE al.severity IN ('warning', 'critical')
    AND al.created_at >= NOW() - INTERVAL '7 days'
ORDER BY al.created_at DESC
LIMIT 100;
```

---

## 11. Optimizaciones

### 11.1 Materialized Views para Reportes

```sql
-- Vista materializada para métricas de organizadores
CREATE MATERIALIZED VIEW mv_organizer_metrics AS
SELECT
    u.id as user_id,
    u.email,
    op.business_name,
    op.verified,
    COUNT(DISTINCT r.id) as total_raffles,
    COUNT(DISTINCT CASE WHEN r.status = 'completed' THEN r.id END) as completed_raffles,
    SUM(CASE WHEN r.status = 'completed' THEN r.total_revenue ELSE 0 END) as total_revenue,
    SUM(CASE WHEN r.status = 'completed' THEN r.platform_fee_amount ELSE 0 END) as total_fees_paid,
    op.total_payouts,
    op.pending_payout
FROM users u
LEFT JOIN organizer_profiles op ON u.id = op.user_id
LEFT JOIN raffles r ON u.id = r.user_id AND r.deleted_at IS NULL
WHERE u.role = 'user'
GROUP BY u.id, op.id;

-- Índice en la vista
CREATE INDEX idx_mv_organizer_metrics_revenue ON mv_organizer_metrics(total_revenue DESC);

-- Refresh automático cada hora (con cron job o pg_cron)
-- SELECT cron.schedule('refresh-organizer-metrics', '0 * * * *', 'REFRESH MATERIALIZED VIEW mv_organizer_metrics');
```

### 11.2 Particionamiento de audit_logs (futuro)

```sql
-- Si la tabla audit_logs crece mucho, considerar particionamiento por mes
-- CREATE TABLE audit_logs_2025_11 PARTITION OF audit_logs
--     FOR VALUES FROM ('2025-11-01') TO ('2025-12-01');
```

---

## 12. Backfill Scripts (Migración de datos existentes)

### 12.1 Crear organizer_profiles para usuarios con rifas

```sql
-- Script para crear organizer_profiles para todos los usuarios que han creado rifas
INSERT INTO organizer_profiles (user_id, verified, created_at, updated_at)
SELECT DISTINCT
    r.user_id,
    false as verified, -- Por defecto no verificados, admin debe revisar
    NOW(),
    NOW()
FROM raffles r
LEFT JOIN organizer_profiles op ON r.user_id = op.user_id
WHERE op.id IS NULL -- Solo crear si no existe
    AND r.deleted_at IS NULL;
```

### 12.2 Crear settlements para rifas completed sin settlement

```sql
-- Script para crear settlements para rifas completadas que no tienen settlement
INSERT INTO settlements (
    raffle_id,
    organizer_id,
    gross_revenue,
    platform_fee,
    platform_fee_percentage,
    net_payout,
    status,
    created_at,
    updated_at
)
SELECT
    r.id,
    r.user_id,
    r.total_revenue,
    r.platform_fee_amount,
    r.platform_fee_percentage,
    r.net_amount,
    'pending' as status, -- Requiere aprobación de admin
    r.completed_at,
    NOW()
FROM raffles r
LEFT JOIN settlements s ON r.id = s.raffle_id
WHERE r.status = 'completed'
    AND s.id IS NULL -- No tiene settlement
    AND r.total_revenue > 0;
```

---

## 13. Validaciones y Constraints

### 13.1 Check Constraints

```sql
-- Validar que net_payout = gross_revenue - platform_fee
ALTER TABLE settlements
    ADD CONSTRAINT chk_settlements_net_payout
    CHECK (net_payout = gross_revenue - platform_fee);

-- Validar que platform_fee_percentage esté entre 0 y 50%
ALTER TABLE settlements
    ADD CONSTRAINT chk_settlements_fee_percentage
    CHECK (platform_fee_percentage >= 0 AND platform_fee_percentage <= 50);

-- Validar que commission_override esté entre 0 y 50%
ALTER TABLE organizer_profiles
    ADD CONSTRAINT chk_organizer_commission_range
    CHECK (commission_override IS NULL OR (commission_override >= 0 AND commission_override <= 50));
```

---

**Última actualización:** 2025-11-18
**Total de migraciones:** 7 (012-018)
**Estado:** Pendiente de ejecución
