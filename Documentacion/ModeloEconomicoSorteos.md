# Modelo Económico y Contable - Plataforma de Sorteos

**Documento de Arquitectura Financiera**  
Fecha: 12 de noviembre, 2025  
Versión: 1.0

---

## Tabla de Contenidos

1. [Visión General del Modelo](#1-visión-general-del-modelo)
2. [Actores y Responsabilidades](#2-actores-y-responsabilidades)
3. [Estructura de Costos y Comisiones](#3-estructura-de-costos-y-comisiones)
4. [Flujo de Fondos Completo](#4-flujo-de-fondos-completo)
5. [Matemáticas del Sistema de Prepago](#5-matemáticas-del-sistema-de-prepago)
6. [Arquitectura Contable](#6-arquitectura-contable)
7. [Sistema de Ledger (Libro Mayor)](#7-sistema-de-ledger-libro-mayor)
8. [Reconciliación y Solvencia](#8-reconciliación-y-solvencia)
9. [Casos de Uso Detallados](#9-casos-de-uso-detallados)
10. [Implementación Técnica](#10-implementación-técnica)

---

## 1. Visión General del Modelo

### Filosofía del Modelo

La plataforma opera como un **marketplace puro tipo Uber**: conecta organizadores de rifas con participantes, procesando pagos y reteniendo comisiones por el servicio.

**Principios fundamentales:**

1. **Transparencia total**: Cada actor sabe exactamente qué paga y qué recibe
2. **Costos distribuidos justamente**: 
   - El participante paga los costos transaccionales fijos
   - El organizador paga los costos variables (% procesador + % plataforma)
3. **Sistema de prepago opcional**: Para optimizar experiencia de usuarios frecuentes
4. **Liquidez garantizada**: El dinero real siempre cubre las obligaciones

### Modelo de Negocio

```
┌─────────────────────────────────────────────────────────────┐
│                    PLATAFORMA DE SORTEOS                     │
│                                                              │
│  PARTICIPANTE ←→ PLATAFORMA ←→ ORGANIZADOR                  │
│                                                              │
│  • Paga número     • Procesa pagos    • Recibe liquidación  │
│  + fee fijo        • Retiene comisión  - comisión plataforma│
│                    • Transfiere fondos - comisión procesador │
└─────────────────────────────────────────────────────────────┘
```

---

## 2. Actores y Responsabilidades

### 2.1 Participante (Comprador de Números)

**Responsabilidad económica:**
- Paga el precio nominal del número (ej: ₡1,000)
- Asume el costo fijo transaccional ($0.40 ≈ ₡200)

**Opciones de pago:**
1. **Pago directo**: Paga con tarjeta por cada número
2. **Wallet prepago**: Recarga saldo previamente y consume sin fees adicionales

### 2.2 Organizador (Dueño de la Rifa)

**Responsabilidad económica:**
- Asume el costo del procesador de tarjetas (5% del monto)
- Asume la comisión de la plataforma (6% del monto)

**Recibe:**
- 89% del precio nominal de cada número vendido
- Ejemplo: Si el número vale ₡1,000 → recibe ₡890

### 2.3 Plataforma

**Responsabilidades:**
- Procesar pagos con tarjeta
- Gestionar wallets de usuarios
- Administrar fondos en custodia
- Liquidar a organizadores
- Garantizar solvencia del sistema

**Ingresos:**
- Comisión sobre ventas (6% del precio nominal)
- Margen sobre fees transaccionales (en modelo de prepago)

---

## 3. Estructura de Costos y Comisiones

### 3.1 Parámetros del Sistema

```go
type PlatformFees struct {
    // Costo del procesador de tarjetas
    ProcessorPercentRate decimal.Decimal // 5% del monto procesado
    ProcessorFixedFeeUSD decimal.Decimal // $0.40 por transacción
    
    // Comisión de la plataforma
    PlatformCommissionRate decimal.Decimal // 6% del precio nominal
    
    // Tipo de cambio (para convertir USD → CRC)
    ExchangeRate decimal.Decimal // ej: 500 CRC por USD
}

// Valores actuales
var CurrentFees = PlatformFees{
    ProcessorPercentRate:   decimal.NewFromFloat(0.05),  // 5%
    ProcessorFixedFeeUSD:   decimal.NewFromFloat(0.40),  // $0.40
    PlatformCommissionRate: decimal.NewFromFloat(0.06),  // 6%
    ExchangeRate:           decimal.NewFromInt(500),     // ₡500/$
}
```

### 3.2 Desglose de un Pago Directo

**Ejemplo: Número de ₡1,000**

```
PARTICIPANTE PAGA:
├─ Precio del número: ₡1,000
├─ Fee transaccional: ₡200 ($0.40)
└─ TOTAL COBRADO: ₡1,200

PROCESADOR DE TARJETAS COBRA:
├─ 5% de ₡1,200 = ₡60
└─ Fee ya incluido: ₡200
   TOTAL FEES: ₡260

PLATAFORMA RECIBE NETO:
├─ Cobrado: ₡1,200
├─ Menos procesador: -₡260
└─ NETO: ₡940

DISTRIBUCIÓN DEL NETO:
├─ Para organizador (89%): ₡890
├─ Para plataforma (11%): ₡50
└─ TOTAL: ₡940 ✓
```

### 3.3 Cálculo Detallado

```
Sea P = Precio nominal del número = ₡1,000
Sea F = Fee fijo transaccional = ₡200
Sea r_proc = Tasa procesador = 5%
Sea r_plat = Comisión plataforma = 6%

PAGO DEL PARTICIPANTE:
Total = P + F = ₡1,200

COSTO PROCESADOR:
Costo = (P + F) × r_proc + F
      = ₡1,200 × 0.05 + ₡200
      = ₡60 + ₡200
      = ₡260

NETO PARA DISTRIBUIR:
Neto = Total - Costo_Procesador
     = ₡1,200 - ₡260
     = ₡940

ORGANIZADOR RECIBE:
Organizador = P × (1 - r_proc - r_plat)
            = ₡1,000 × (1 - 0.05 - 0.06)
            = ₡1,000 × 0.89
            = ₡890

PLATAFORMA RETIENE:
Plataforma = Neto - Organizador
           = ₡940 - ₡890
           = ₡50
```

---

## 4. Flujo de Fondos Completo

### 4.1 Escenario Completo: Rifa de 100 Números

**Configuración:**
- Rifa de 100 números a ₡1,000 cada uno
- Valor nominal total: ₡100,000
- Todos los participantes usan wallet prepago

### 4.2 Fase 1: Recargas de Wallet

**100 usuarios recargan para tener ₡10,000 de crédito cada uno**

```
Usuario individual:
├─ Desea crédito: ₡10,000
├─ Debe pagar: ₡10,737 (calculado con fórmula)
└─ Procesador cobra: ₡737 (5% + $0.40)

Sistema completo (100 usuarios):
├─ Total procesado en tarjetas: ₡1,073,700
├─ Procesador se lleva: ₡73,700
├─ Ingresa a caja plataforma: ₡1,000,000
└─ Créditos otorgados: ₡1,000,000 (100 × ₡10,000)
```

**Estado del sistema después de recargas:**

```
╔════════════════════════════════════════╗
║ CAJA DE LA PLATAFORMA (dinero real)   ║
║ Balance: ₡1,000,000                    ║
╚════════════════════════════════════════╝
         ↓ Este dinero cubre ↓
╔════════════════════════════════════════╗
║ WALLETS DE USUARIOS (créditos)        ║
║ 100 usuarios × ₡10,000 = ₡1,000,000   ║
╚════════════════════════════════════════╝
```

### 4.3 Fase 2: Compra de Números con Wallet

**100 usuarios compran 1 número de ₡1,000 cada uno**

```
Por cada compra:
├─ Se debita ₡1,000 del wallet del usuario
├─ Se registra venta para el organizador
└─ NO se mueve dinero real (solo contabilidad)

Después de 100 compras:
├─ Wallets usuarios: ₡1,000,000 → ₡900,000
├─ Ventas del organizador: ₡100,000 (nominal)
└─ Caja plataforma: ₡1,000,000 (sin cambios aún)
```

**Estado del sistema después de compras:**

```
╔════════════════════════════════════════╗
║ CAJA DE LA PLATAFORMA                  ║
║ Balance: ₡1,000,000 (sin cambios)      ║
╚════════════════════════════════════════╝
         ↓ Debe cubrir ↓
╔════════════════════════════════════════╗
║ WALLETS DE USUARIOS                    ║
║ Saldo: ₡900,000                        ║
╚════════════════════════════════════════╝
         +
╔════════════════════════════════════════╗
║ DEUDA CON ORGANIZADOR                  ║
║ Por pagar: ₡89,000 (89% × ₡100,000)   ║
╚════════════════════════════════════════╝
         +
╔════════════════════════════════════════╗
║ GANANCIA PLATAFORMA                    ║
║ Comisión: ₡11,000 (11% × ₡100,000)    ║
╚════════════════════════════════════════╝

Verificación:
₡900,000 + ₡89,000 + ₡11,000 = ₡1,000,000 ✓
```

### 4.4 Fase 3: Liquidación al Organizador

**Se completa la rifa, se transfiere dinero real al organizador**

```
Liquidación:
├─ Números vendidos: 100 × ₡1,000 = ₡100,000
├─ Organizador recibe: ₡89,000 (89%)
└─ Se transfiere de caja plataforma → cuenta organizador

Estado final del sistema:
├─ Caja plataforma: ₡1,000,000 - ₡89,000 = ₡911,000
├─ Wallets usuarios: ₡900,000
└─ Ganancia plataforma acumulada: ₡11,000

Verificación de solvencia:
₡911,000 (caja) = ₡900,000 (wallets) + ₡11,000 (ganancia) ✓
```

---

## 5. Matemáticas del Sistema de Prepago

### 5.1 El Problema

El usuario quiere recibir exactamente ₡D de crédito, pero el procesador cobra fees sobre el monto procesado.

**No podemos cobrar ₡D directamente** porque:
```
Si cobramos ₡10,000:
├─ Procesador cobra 5%: ₡500
├─ Procesador cobra fijo: ₡200
├─ Total fees: ₡700
└─ Nos quedan: ₡9,300 (¡no alcanza para dar ₡10,000!)
```

### 5.2 La Fórmula Correcta

**Necesitamos calcular cuánto cobrar (C) para que después de fees queden los ₡D deseados.**

```
Variables:
  C = Monto a cobrar (lo que buscamos)
  D = Crédito deseado (₡10,000)
  r = Tasa del procesador (0.05 = 5%)
  f = Fee fijo en CRC (₡200)

Ecuación:
  Lo que queda después de fees = Crédito deseado
  C - (C × r) - f = D
  C × (1 - r) - f = D
  C × (1 - r) = D + f
  
  C = (D + f) / (1 - r)
```

### 5.3 Ejemplos de Cálculo

#### Caso 1: Usuario quiere ₡10,000 de crédito

```
D = ₡10,000
r = 0.05
f = ₡200

C = (₡10,000 + ₡200) / (1 - 0.05)
  = ₡10,200 / 0.95
  = ₡10,736.84
  ≈ ₡10,737 (redondeado)

Verificación:
  Cobrado: ₡10,737
  Fee 5%: ₡10,737 × 0.05 = ₡536.85
  Fee fijo: ₡200
  Total fees: ₡736.85
  Neto: ₡10,737 - ₡737 = ₡10,000 ✓
```

#### Caso 2: Usuario quiere ₡50,000 de crédito

```
C = (₡50,000 + ₡200) / 0.95
  = ₡50,200 / 0.95
  = ₡52,842.10
  ≈ ₡52,842

Verificación:
  Cobrado: ₡52,842
  Fee 5%: ₡2,642
  Fee fijo: ₡200
  Neto: ₡50,000 ✓
```

### 5.4 Tabla de Referencia

| Crédito Deseado | Monto a Cobrar | Fees Totales | % Fee Efectivo |
|-----------------|----------------|--------------|----------------|
| ₡5,000          | ₡5,368         | ₡368         | 7.36%          |
| ₡10,000         | ₡10,737        | ₡737         | 7.37%          |
| ₡20,000         | ₡21,263        | ₡1,263       | 6.32%          |
| ₡50,000         | ₡52,842        | ₡2,842       | 5.68%          |
| ₡100,000        | ₡105,474       | ₡5,474       | 5.47%          |

**Observación:** El % efectivo de fee disminuye con montos mayores (economía de escala del fee fijo).

---

## 6. Arquitectura Contable

### 6.1 Tipos de Cuentas

El sistema maneja **tres tipos fundamentales de cuentas**:

```go
// 1. Cuenta de Caja de la Plataforma (REAL)
// Representa el dinero físico que la plataforma tiene
type PlatformCashAccount struct {
    Balance decimal.Decimal
    // Este es dinero REAL en tu cuenta bancaria
}

// 2. Wallets de Usuarios (VIRTUAL)
// Representan créditos que los usuarios pueden gastar
type UserWallet struct {
    UserID  uuid.UUID
    Balance decimal.Decimal
    // Este es un "vale" o "crédito virtual"
}

// 3. Cuentas por Pagar a Organizadores (PASIVO)
// Representan dinero que DEBES a organizadores
type OrganizerBalance struct {
    OrganizerID     uuid.UUID
    PendingPayout   decimal.Decimal
    // Este es una DEUDA de la plataforma
}
```

### 6.2 Ecuación Fundamental de Solvencia

**En todo momento debe cumplirse:**

```
Caja de la Plataforma ≥ Wallets de Usuarios + Deudas con Organizadores
```

**O de forma más precisa:**

```
Activos (Caja) = Pasivos (Wallets + Organizadores) + Capital (Ganancia)
```

### 6.3 Diagrama de Cuentas

```
┌──────────────────────────────────────────────────────────┐
│                        ACTIVOS                            │
├──────────────────────────────────────────────────────────┤
│ platform_cash                            ₡1,000,000      │
│   ├─ Dinero en cuenta bancaria                           │
│   └─ DEBE SER ≥ suma de pasivos                          │
└──────────────────────────────────────────────────────────┘

┌──────────────────────────────────────────────────────────┐
│                        PASIVOS                            │
├──────────────────────────────────────────────────────────┤
│ user_wallets_liability                   ₡900,000        │
│   ├─ Suma de todos los wallets                           │
│   └─ Dinero que DEBES a usuarios                         │
│                                                           │
│ organizer_balances_liability             ₡89,000         │
│   ├─ Suma de cuentas por pagar                           │
│   └─ Dinero que DEBES a organizadores                    │
└──────────────────────────────────────────────────────────┘

┌──────────────────────────────────────────────────────────┐
│                     PATRIMONIO                            │
├──────────────────────────────────────────────────────────┤
│ platform_revenue                         ₡11,000         │
│   ├─ Ganancia acumulada                                  │
│   └─ Capital propio de la plataforma                     │
└──────────────────────────────────────────────────────────┘

VERIFICACIÓN:
Activos = Pasivos + Patrimonio
₡1,000,000 = ₡900,000 + ₡89,000 + ₡11,000 ✓
```

---

## 7. Sistema de Ledger (Libro Mayor)

### 7.1 Estructura de la Tabla

El ledger es un **registro inmutable** de todas las transacciones financieras del sistema.

```sql
CREATE TABLE ledger_entries (
    -- Identificación
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    entry_type VARCHAR(50) NOT NULL,
    
    -- Monto de la transacción
    amount DECIMAL(12,2) NOT NULL,
    
    -- Contabilidad de doble entrada
    debit_account_type VARCHAR(50) NOT NULL,
    debit_account_id UUID,
    credit_account_type VARCHAR(50) NOT NULL,
    credit_account_id UUID,
    
    -- Referencia a la operación origen
    reference_type VARCHAR(50),
    reference_id UUID,
    
    -- Metadatos
    description TEXT,
    metadata JSONB,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    
    -- Índices para consultas eficientes
    INDEX idx_created (created_at DESC),
    INDEX idx_debit_account (debit_account_type, debit_account_id),
    INDEX idx_credit_account (credit_account_type, credit_account_id),
    INDEX idx_reference (reference_type, reference_id)
);

-- Tipos de cuentas válidos
CREATE TYPE account_type AS ENUM (
    'platform_cash',      -- Caja de la plataforma
    'platform_liability', -- Pasivos de la plataforma
    'platform_revenue',   -- Ingresos de la plataforma
    'user_wallet',        -- Wallet individual de usuario
    'organizer_balance',  -- Cuenta por pagar a organizador
    'external'            -- Cuenta externa (procesador, banco)
);

-- Tipos de entradas válidos
CREATE TYPE entry_type AS ENUM (
    'user_recharge_cash',    -- Dinero entra de recarga
    'user_recharge_credit',  -- Se acredita wallet
    'number_purchase_debit', -- Se debita wallet
    'number_sale_credit',    -- Se acredita organizador
    'platform_commission',   -- Comisión de plataforma
    'organizer_payout'       -- Pago a organizador
);
```

### 7.2 Principios del Ledger

1. **Inmutabilidad**: Una vez creada, una entrada NUNCA se modifica ni elimina
2. **Doble entrada**: Cada transacción tiene un débito y un crédito
3. **Trazabilidad**: Cada entrada referencia la operación que la originó
4. **Auditabilidad**: El balance de cualquier cuenta se puede reconstruir desde el inicio

### 7.3 Asientos Contables por Operación

#### A. Usuario Recarga ₡10,737 → Recibe ₡10,000 de Crédito

**Asiento 1: Dinero entra a caja**
```sql
INSERT INTO ledger_entries (
    entry_type,
    amount,
    debit_account_type,
    debit_account_id,
    credit_account_type,
    credit_account_id,
    reference_type,
    reference_id,
    description
) VALUES (
    'user_recharge_cash',
    10737.00,
    'external',                    -- DÉBITO: viene del procesador
    NULL,
    'platform_cash',               -- CRÉDITO: entra a nuestra caja
    '00000000-0000-0000-0000-000000000001',
    'recharge',
    'recharge-uuid-12345',
    'Recarga via tarjeta de usuario-uuid-789'
);
```

**Asiento 2: Se acredita wallet del usuario**
```sql
INSERT INTO ledger_entries (
    entry_type,
    amount,
    debit_account_type,
    debit_account_id,
    credit_account_type,
    credit_account_id,
    reference_type,
    reference_id,
    description
) VALUES (
    'user_recharge_credit',
    10000.00,
    'platform_liability',          -- DÉBITO: aumenta nuestro pasivo
    '00000000-0000-0000-0000-000000000001',
    'user_wallet',                 -- CRÉDITO: aumenta el wallet
    'user-uuid-789',
    'recharge',
    'recharge-uuid-12345',
    'Crédito otorgado por recarga'
);
```

**Efecto neto:**
- `platform_cash`: +₡10,737
- `user_wallet[user-789]`: +₡10,000
- Diferencia: ₡737 (margen que cubre los fees del procesador)

#### B. Usuario Compra Número de ₡1,000 con Wallet

**Asiento 1: Se debita wallet**
```sql
INSERT INTO ledger_entries (
    entry_type,
    amount,
    debit_account_type,
    debit_account_id,
    credit_account_type,
    credit_account_id,
    reference_type,
    reference_id,
    description
) VALUES (
    'number_purchase_debit',
    1000.00,
    'user_wallet',                 -- DÉBITO: baja el wallet
    'user-uuid-789',
    'platform_liability',          -- CRÉDITO: baja el pasivo
    '00000000-0000-0000-0000-000000000001',
    'raffle_number',
    'purchase-uuid-456',
    'Compra número 42 en rifa-uuid-123'
);
```

**Asiento 2: Se acredita cuenta del organizador**
```sql
INSERT INTO ledger_entries (
    entry_type,
    amount,
    debit_account_type,
    debit_account_id,
    credit_account_type,
    credit_account_id,
    reference_type,
    reference_id,
    description
) VALUES (
    'number_sale_credit',
    890.00,                        -- 89% del precio nominal
    'platform_cash',               -- DÉBITO: comprometemos caja
    '00000000-0000-0000-0000-000000000001',
    'organizer_balance',           -- CRÉDITO: deuda con organizador
    'organizer-uuid-555',
    'raffle_number',
    'purchase-uuid-456',
    'Venta número 42 (89% para organizador)'
);
```

**Asiento 3: Comisión de plataforma**
```sql
INSERT INTO ledger_entries (
    entry_type,
    amount,
    debit_account_type,
    debit_account_id,
    credit_account_type,
    credit_account_id,
    reference_type,
    reference_id,
    description
) VALUES (
    'platform_commission',
    110.00,                        -- 11% del precio nominal
    'platform_cash',               -- DÉBITO: no sale, queda en caja
    '00000000-0000-0000-0000-000000000001',
    'platform_revenue',            -- CRÉDITO: ingreso
    '00000000-0000-0000-0000-000000000001',
    'raffle_number',
    'purchase-uuid-456',
    'Comisión plataforma (11%)'
);
```

**Efecto neto:**
- `user_wallet[user-789]`: -₡1,000
- `organizer_balance[org-555]`: +₡890
- `platform_revenue`: +₡110
- `platform_cash`: sin cambio (se reasigna internamente)

#### C. Liquidación al Organizador

**Asiento único: Pago**
```sql
INSERT INTO ledger_entries (
    entry_type,
    amount,
    debit_account_type,
    debit_account_id,
    credit_account_type,
    credit_account_id,
    reference_type,
    reference_id,
    description
) VALUES (
    'organizer_payout',
    89000.00,                      -- 100 números × ₡890
    'organizer_balance',           -- DÉBITO: se salda la deuda
    'organizer-uuid-555',
    'platform_cash',               -- CRÉDITO: sale de caja
    '00000000-0000-0000-0000-000000000001',
    'payout',
    'payout-uuid-999',
    'Liquidación rifa-uuid-123 completada'
);
```

**Efecto neto:**
- `platform_cash`: -₡89,000
- `organizer_balance[org-555]`: -₡89,000 (queda en 0)

---

## 8. Reconciliación y Solvencia

### 8.1 Servicio de Reconciliación

```go
type ReconciliationService struct {
    db *gorm.DB
}

func (rs *ReconciliationService) ReconcileAll() (*ReconciliationReport, error) {
    report := &ReconciliationReport{
        Timestamp: time.Now(),
    }
    
    // 1. Calcular caja de plataforma
    rs.db.Raw(`
        SELECT SUM(
            CASE 
                WHEN credit_account_type = 'platform_cash' THEN amount
                WHEN debit_account_type = 'platform_cash' THEN -amount
                ELSE 0
            END
        ) as balance
        FROM ledger_entries
    `).Scan(&report.PlatformCash)
    
    // 2. Calcular total de wallets
    rs.db.Raw(`
        SELECT SUM(
            CASE 
                WHEN credit_account_type = 'user_wallet' THEN amount
                WHEN debit_account_type = 'user_wallet' THEN -amount
                ELSE 0
            END
        ) as balance
        FROM ledger_entries
    `).Scan(&report.TotalUserWallets)
    
    // 3. Calcular deuda con organizadores
    rs.db.Raw(`
        SELECT SUM(
            CASE 
                WHEN credit_account_type = 'organizer_balance' THEN amount
                WHEN debit_account_type = 'organizer_balance' THEN -amount
                ELSE 0
            END
        ) as balance
        FROM ledger_entries
    `).Scan(&report.TotalOrganizerDebt)
    
    // 4. Calcular ingresos de plataforma
    rs.db.Raw(`
        SELECT SUM(amount) as revenue
        FROM ledger_entries
        WHERE credit_account_type = 'platform_revenue'
    `).Scan(&report.PlatformRevenue)
    
    // 5. Verificar ecuación fundamental
    expectedCash := report.TotalUserWallets.
        Add(report.TotalOrganizerDebt).
        Add(report.PlatformRevenue)
    
    report.ExpectedCash = expectedCash
    report.Discrepancy = report.PlatformCash.Sub(expectedCash)
    report.Balanced = report.Discrepancy.Abs().LessThan(decimal.NewFromFloat(0.01))
    
    return report, nil
}

type ReconciliationReport struct {
    Timestamp          time.Time
    PlatformCash       decimal.Decimal // Lo que TENEMOS
    TotalUserWallets   decimal.Decimal // Lo que DEBEMOS a usuarios
    TotalOrganizerDebt decimal.Decimal // Lo que DEBEMOS a organizadores
    PlatformRevenue    decimal.Decimal // Lo que HEMOS GANADO
    ExpectedCash       decimal.Decimal // Lo que DEBERÍA haber en caja
    Discrepancy        decimal.Decimal // Diferencia (debe ser ≈0)
    Balanced           bool            // ¿Todo cuadra?
}
```

### 8.2 Checks de Solvencia

```go
func (rs *ReconciliationService) CheckSolvency() (*SolvencyCheck, error) {
    var cash, liabilities decimal.Decimal
    
    // Obtener caja
    rs.db.Raw(`
        SELECT SUM(
            CASE 
                WHEN credit_account_type = 'platform_cash' THEN amount
                WHEN debit_account_type = 'platform_cash' THEN -amount
                ELSE 0
            END
        ) FROM ledger_entries
    `).Scan(&cash)
    
    // Obtener pasivos totales
    rs.db.Raw(`
        SELECT SUM(balance) FROM (
            SELECT SUM(
                CASE 
                    WHEN credit_account_type = 'user_wallet' THEN amount
                    WHEN debit_account_type = 'user_wallet' THEN -amount
                    ELSE 0
                END
            ) as balance
            FROM ledger_entries
            
            UNION ALL
            
            SELECT SUM(
                CASE 
                    WHEN credit_account_type = 'organizer_balance' THEN amount
                    WHEN debit_account_type = 'organizer_balance' THEN -amount
                    ELSE 0
                END
            ) as balance
            FROM ledger_entries
        ) totals
    `).Scan(&liabilities)
    
    check := &SolvencyCheck{
        Cash:        cash,
        Liabilities: liabilities,
        Solvent:     cash.GreaterThanOrEqual(liabilities),
        Ratio:       cash.Div(liabilities), // Debe ser ≥ 1.0
    }
    
    if !check.Solvent {
        // ¡ALERTA CRÍTICA! No hay suficiente dinero
        rs.sendCriticalAlert(check)
    }
    
    return check, nil
}

type SolvencyCheck struct {
    Cash        decimal.Decimal
    Liabilities decimal.Decimal
    Solvent     bool
    Ratio       decimal.Decimal // Cash / Liabilities
}
```

### 8.3 Alertas y Monitoreo

```go
// Monitoreo continuo de solvencia
func (rs *ReconciliationService) MonitorSolvency(ctx context.Context) {
    ticker := time.NewTicker(5 * time.Minute)
    defer ticker.Stop()
    
    for {
        select {
        case <-ctx.Done():
            return
        case <-ticker.C:
            check, err := rs.CheckSolvency()
            if err != nil {
                log.Error("Error checking solvency", "error", err)
                continue
            }
            
            // Umbrales de alerta
            if !check.Solvent {
                rs.sendCriticalAlert(check) // CRÍTICO: insolvencia
            } else if check.Ratio.LessThan(decimal.NewFromFloat(1.1)) {
                rs.sendWarningAlert(check) // WARNING: margen bajo (<10%)
            }
            
            // Registrar métricas
            rs.recordMetrics(check)
        }
    }
}
```

---

## 9. Casos de Uso Detallados

### 9.1 Caso: Usuario Nuevo - Primera Compra

**Escenario:**
- Usuario Juan ve una rifa de un iPhone
- Número cuesta ₡1,000
- Juan no tiene wallet, paga directo

**Flujo:**

1. **Juan selecciona número**
   ```
   Frontend muestra:
   ┌─────────────────────────────────┐
   │ Número: 42                      │
   │ Precio: ₡1,000                  │
   │ Cargo por servicio: ₡200        │
   │ ─────────────────────────────   │
   │ TOTAL: ₡1,200                   │
   │                                 │
   │ [Pagar con tarjeta]             │
   └─────────────────────────────────┘
   ```

2. **Procesamiento del pago**
   ```
   Juan paga ₡1,200 con tarjeta
   ├─ Procesador cobra 5%: ₡60
   ├─ Procesador cobra fijo: ₡200
   └─ Neto para plataforma: ₡940
   ```

3. **Asientos contables** (ver sección 7.3)

4. **Resultado**
   ```
   ✓ Juan tiene el número 42 reservado
   ✓ Organizador tiene ₡890 acumulados (por cobrar)
   ✓ Plataforma tiene ₡50 de ganancia
   ```

### 9.2 Caso: Usuario Frecuente - Usa Wallet

**Escenario:**
- María participa en muchas rifas
- Decide recargar ₡50,000 para ahorrar en fees

**Flujo:**

1. **María recarga wallet**
   ```
   María quiere: ₡50,000 de crédito
   Sistema calcula: debe pagar ₡52,842
   
   Frontend muestra:
   ┌─────────────────────────────────┐
   │ Vas a recargar: ₡50,000         │
   │                                 │
   │ Monto a pagar: ₡52,842          │
   │ (incluye procesamiento ₡2,842)  │
   │                                 │
   │ Tu saldo quedará en: ₡50,000    │
   │                                 │
   │ Con este saldo puedes comprar   │
   │ hasta 50 números sin pagar      │
   │ cargos adicionales              │
   │                                 │
   │ [Confirmar recarga]             │
   └─────────────────────────────────┘
   ```

2. **Procesamiento**
   ```
   Paga ₡52,842
   ├─ Procesador cobra ₡2,842
   └─ Plataforma recibe ₡50,000
       └─ Acredita wallet de María: ₡50,000
   ```

3. **María compra 5 números**
   ```
   Por cada número de ₡1,000:
   ├─ Se debita ₡1,000 de su wallet
   ├─ NO paga fee adicional
   └─ Saldo: ₡50,000 → ₡45,000
   
   Ahorro vs pago directo:
   ├─ Directo: 5 × ₡1,200 = ₡6,000
   ├─ Con wallet: 5 × ₡1,000 = ₡5,000
   └─ Ahorro: ₡1,000 (4 × ₡200 de fees no pagados)
   ```

### 9.3 Caso: Organizador - Liquidación

**Escenario:**
- Pedro organiza rifa de una moto
- 200 números a ₡5,000 cada uno
- Se venden todos

**Flujo:**

1. **Ventas acumuladas**
   ```
   200 números × ₡5,000 = ₡1,000,000 (nominal)
   
   Por cada venta:
   ├─ Pedro acumula: ₡4,450 (89%)
   └─ Plataforma retiene: ₡550 (11%)
   ```

2. **Balance de Pedro**
   ```
   Después de 200 ventas:
   ├─ Total por cobrar: ₡890,000
   ├─ Estado: PENDIENTE
   └─ Esperando cierre de rifa
   ```

3. **Cierre y liquidación**
   ```
   Se define ganador
   ├─ Rifa cambia a estado: COMPLETED
   ├─ Se activa liquidación automática
   └─ O Pedro solicita pago manual
   
   Transferencia:
   ├─ De: Caja de plataforma
   ├─ A: Cuenta bancaria de Pedro
   ├─ Monto: ₡890,000
   └─ Tiempo: 24-48 horas hábiles
   ```

4. **Comprobante**
   ```
   ┌─────────────────────────────────────┐
   │ Liquidación Rifa #12345             │
   │                                     │
   │ Números vendidos: 200               │
   │ Precio unitario: ₡5,000             │
   │ Total vendido: ₡1,000,000           │
   │                                     │
   │ Desglose:                           │
   │  • Costo procesador (5%): ₡50,000   │
   │  • Comisión plataforma (6%): ₡60,000│
   │  ───────────────────────────────    │
   │  Total a recibir: ₡890,000          │
   │                                     │
   │ Estado: PROCESADO                   │
   │ Fecha: 2025-11-12 10:30             │
   │                                     │
   │ [Descargar comprobante]             │
   └─────────────────────────────────────┘
   ```

---

## 10. Implementación Técnica

### 10.1 Estructura de Servicios

```go
// Servicio principal de wallet
type WalletService struct {
    db                *gorm.DB
    ledger            *LedgerService
    paymentProcessor  *PaymentProcessorClient
    calculator        *RechargeCalculator
}

// Servicio de ledger contable
type LedgerService struct {
    db *gorm.DB
}

// Servicio de liquidaciones
type PayoutService struct {
    db          *gorm.DB
    ledger      *LedgerService
    bankClient  *BankTransferClient
}
```

### 10.2 Operación: Recarga de Wallet

```go
func (ws *WalletService) ProcessRecharge(
    ctx context.Context,
    userID uuid.UUID,
    desiredCredit decimal.Decimal,
    paymentMethodID string,
) (*RechargeResult, error) {
    
    // 1. Calcular cuánto cobrar
    breakdown := ws.calculator.CalculateChargeForCredit(desiredCredit)
    
    // 2. Procesar pago con tarjeta
    payment, err := ws.paymentProcessor.ChargeCard(
        ctx,
        paymentMethodID,
        breakdown.ChargeAmount,
        map[string]interface{}{
            "type":    "wallet_recharge",
            "user_id": userID.String(),
            "credit":  desiredCredit.String(),
        },
    )
    if err != nil {
        return nil, fmt.Errorf("payment failed: %w", err)
    }
    
    // 3. Registrar en ledger y acreditar wallet
    return ws.db.Transaction(func(tx *gorm.DB) (*RechargeResult, error) {
        
        // 3a. Asiento: dinero entra a caja
        if err := ws.ledger.RecordEntry(ctx, tx, &LedgerEntry{
            EntryType:          "user_recharge_cash",
            Amount:             breakdown.ChargeAmount,
            DebitAccountType:   "external",
            CreditAccountType:  "platform_cash",
            CreditAccountID:    PlatformAccountID,
            ReferenceType:      "recharge",
            ReferenceID:        payment.ID,
            Description:        fmt.Sprintf("Recarga de usuario %s", userID),
        }); err != nil {
            return nil, err
        }
        
        // 3b. Asiento: se acredita wallet
        if err := ws.ledger.RecordEntry(ctx, tx, &LedgerEntry{
            EntryType:          "user_recharge_credit",
            Amount:             desiredCredit,
            DebitAccountType:   "platform_liability",
            DebitAccountID:     PlatformAccountID,
            CreditAccountType:  "user_wallet",
            CreditAccountID:    userID,
            ReferenceType:      "recharge",
            ReferenceID:        payment.ID,
            Description:        "Crédito otorgado por recarga",
        }); err != nil {
            return nil, err
        }
        
        // 3c. Actualizar balance del wallet
        if err := tx.Exec(`
            INSERT INTO user_wallets (user_id, balance)
            VALUES (?, ?)
            ON CONFLICT (user_id) DO UPDATE
            SET balance = user_wallets.balance + EXCLUDED.balance,
                updated_at = NOW()
        `, userID, desiredCredit).Error; err != nil {
            return nil, err
        }
        
        return &RechargeResult{
            PaymentID:     payment.ID,
            AmountCharged: breakdown.ChargeAmount,
            CreditGranted: desiredCredit,
            NewBalance:    ws.getWalletBalance(tx, userID),
        }, nil
    })
}
```

### 10.3 Operación: Compra de Número con Wallet

```go
func (ws *WalletService) PurchaseNumberWithWallet(
    ctx context.Context,
    userID uuid.UUID,
    raffleID uuid.UUID,
    numberID uuid.UUID,
    price decimal.Decimal,
) error {
    
    return ws.db.Transaction(func(tx *gorm.DB) error {
        
        // 1. Verificar y bloquear wallet
        var wallet UserWallet
        if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
            Where("user_id = ?", userID).
            First(&wallet).Error; err != nil {
            return fmt.Errorf("wallet not found: %w", err)
        }
        
        if wallet.Balance.LessThan(price) {
            return ErrInsufficientBalance
        }
        
        // 2. Calcular distribución
        organizerAmount := price.Mul(decimal.NewFromFloat(0.89)) // 89%
        platformAmount := price.Mul(decimal.NewFromFloat(0.11))  // 11%
        
        // 3. Asiento: debitar wallet
        purchaseID := uuid.New()
        if err := ws.ledger.RecordEntry(ctx, tx, &LedgerEntry{
            EntryType:          "number_purchase_debit",
            Amount:             price,
            DebitAccountType:   "user_wallet",
            DebitAccountID:     userID,
            CreditAccountType:  "platform_liability",
            CreditAccountID:    PlatformAccountID,
            ReferenceType:      "raffle_number",
            ReferenceID:        purchaseID,
            Description:        fmt.Sprintf("Compra número %s", numberID),
        }); err != nil {
            return err
        }
        
        // 4. Asiento: acreditar organizador
        organizer := ws.getOrganizerID(tx, raffleID)
        if err := ws.ledger.RecordEntry(ctx, tx, &LedgerEntry{
            EntryType:          "number_sale_credit",
            Amount:             organizerAmount,
            DebitAccountType:   "platform_cash",
            DebitAccountID:     PlatformAccountID,
            CreditAccountType:  "organizer_balance",
            CreditAccountID:    organizer,
            ReferenceType:      "raffle_number",
            ReferenceID:        purchaseID,
            Description:        "Venta de número (89% organizador)",
        }); err != nil {
            return err
        }
        
        // 5. Asiento: comisión plataforma
        if err := ws.ledger.RecordEntry(ctx, tx, &LedgerEntry{
            EntryType:          "platform_commission",
            Amount:             platformAmount,
            DebitAccountType:   "platform_cash",
            DebitAccountID:     PlatformAccountID,
            CreditAccountType:  "platform_revenue",
            CreditAccountID:    PlatformAccountID,
            ReferenceType:      "raffle_number",
            ReferenceID:        purchaseID,
            Description:        "Comisión plataforma (11%)",
        }); err != nil {
            return err
        }
        
        // 6. Actualizar wallet
        newBalance := wallet.Balance.Sub(price)
        if err := tx.Model(&wallet).Update("balance", newBalance).Error; err != nil {
            return err
        }
        
        // 7. Reservar número en la rifa
        return ws.reserveNumber(tx, raffleID, numberID, userID, purchaseID)
    })
}
```

### 10.4 Operación: Liquidación a Organizador

```go
func (ps *PayoutService) ProcessPayout(
    ctx context.Context,
    raffleID uuid.UUID,
) (*PayoutResult, error) {
    
    return ps.db.Transaction(func(tx *gorm.DB) (*PayoutResult, error) {
        
        // 1. Obtener balance del organizador para esta rifa
        organizerID := ps.getOrganizerID(tx, raffleID)
        
        var pendingAmount decimal.Decimal
        ps.db.Raw(`
            SELECT SUM(amount)
            FROM ledger_entries
            WHERE credit_account_type = 'organizer_balance'
              AND credit_account_id = ?
              AND reference_type = 'raffle_number'
              AND reference_id IN (
                  SELECT id FROM raffle_numbers WHERE raffle_id = ?
              )
        `, organizerID, raffleID).Scan(&pendingAmount)
        
        // 2. Verificar que haya balance por pagar
        if pendingAmount.LessThanOrEqual(decimal.Zero) {
            return nil, ErrNoBalanceToPay
        }
        
        // 3. Verificar solvencia
        check, _ := ps.ledger.CheckSolvency()
        if !check.Solvent {
            return nil, ErrInsufficientFunds
        }
        
        // 4. Iniciar transferencia bancaria
        transfer, err := ps.bankClient.InitiateTransfer(ctx, &BankTransfer{
            RecipientID: organizerID,
            Amount:      pendingAmount,
            Reference:   fmt.Sprintf("Liquidación rifa %s", raffleID),
        })
        if err != nil {
            return nil, fmt.Errorf("bank transfer failed: %w", err)
        }
        
        // 5. Registrar pago en ledger
        payoutID := uuid.New()
        if err := ps.ledger.RecordEntry(ctx, tx, &LedgerEntry{
            EntryType:          "organizer_payout",
            Amount:             pendingAmount,
            DebitAccountType:   "organizer_balance",
            DebitAccountID:     organizerID,
            CreditAccountType:  "platform_cash",
            CreditAccountID:    PlatformAccountID,
            ReferenceType:      "payout",
            ReferenceID:        payoutID,
            Description:        fmt.Sprintf("Liquidación rifa %s", raffleID),
            Metadata: map[string]interface{}{
                "transfer_id": transfer.ID,
                "raffle_id":   raffleID.String(),
            },
        }); err != nil {
            return nil, err
        }
        
        // 6. Marcar rifa como liquidada
        if err := tx.Exec(`
            UPDATE raffles
            SET payout_status = 'COMPLETED',
                payout_date = NOW(),
                payout_id = ?
            WHERE id = ?
        `, payoutID, raffleID).Error; err != nil {
            return nil, err
        }
        
        return &PayoutResult{
            PayoutID:    payoutID,
            TransferID:  transfer.ID,
            Amount:      pendingAmount,
            OrganizerID: organizerID,
            Status:      "COMPLETED",
        }, nil
    })
}
```

### 10.5 Job de Reconciliación Nocturna

```go
func (rs *ReconciliationService) DailyReconciliationJob(ctx context.Context) {
    // Ejecutar todos los días a las 2:00 AM
    schedule := cron.New()
    schedule.AddFunc("0 2 * * *", func() {
        
        log.Info("Starting daily reconciliation")
        
        // 1. Generar reporte de reconciliación
        report, err := rs.ReconcileAll()
        if err != nil {
            log.Error("Reconciliation failed", "error", err)
            rs.sendAlertToAdmin(err)
            return
        }
        
        // 2. Verificar que todo cuadre
        if !report.Balanced {
            log.Error("Accounts not balanced!", 
                "discrepancy", report.Discrepancy.String())
            rs.sendCriticalAlert(report)
        }
        
        // 3. Verificar solvencia
        check, err := rs.CheckSolvency()
        if err != nil {
            log.Error("Solvency check failed", "error", err)
            return
        }
        
        if !check.Solvent {
            log.Error("System is INSOLVENT!", 
                "cash", check.Cash.String(),
                "liabilities", check.Liabilities.String())
            rs.sendCriticalAlert(check)
        }
        
        // 4. Guardar snapshot diario
        snapshot := &DailySnapshot{
            Date:               time.Now(),
            PlatformCash:       report.PlatformCash,
            TotalUserWallets:   report.TotalUserWallets,
            TotalOrganizerDebt: report.TotalOrganizerDebt,
            PlatformRevenue:    report.PlatformRevenue,
            Balanced:           report.Balanced,
            Solvent:            check.Solvent,
        }
        
        if err := rs.db.Create(snapshot).Error; err != nil {
            log.Error("Failed to save snapshot", "error", err)
        }
        
        log.Info("Daily reconciliation completed", 
            "balanced", report.Balanced,
            "solvent", check.Solvent)
    })
    
    schedule.Start()
}
```

---

## Resumen Ejecutivo

### Modelo Económico

1. **Participante** paga precio nominal + fee fijo de transacción
2. **Organizador** recibe 89% del precio nominal (después de costos)
3. **Plataforma** retiene 11% (6% comisión + margen sobre fees)

### Sistema de Prepago (Wallet)

- Usuario recarga monto X
- Fórmula: `Cobrar = (CréditoDeseado + FeeFijo) / (1 - TasaProcesador)`
- Ahorra fees en compras subsecuentes
- Mejora UX para usuarios frecuentes

### Arquitectura Contable

- **Ledger de doble entrada**: Registro inmutable de todas las transacciones
- **Tres tipos de cuentas**: Cash (activo), Wallets (pasivo), Organizadores (pasivo)
- **Ecuación fundamental**: `Caja ≥ Wallets + Organizadores`
- **Reconciliación diaria**: Verificar que todo cuadre

### Garantías del Sistema

1. **Solvencia**: Siempre hay dinero suficiente para pagar obligaciones
2. **Trazabilidad**: Cada peso se puede rastrear desde origen hasta destino
3. **Auditabilidad**: El ledger permite auditorías completas
4. **Transparencia**: Todos los actores conocen los costos por adelantado

---

**Fin del documento**

Este documento debe mantenerse actualizado conforme evolucione el sistema.