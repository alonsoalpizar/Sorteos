# Parametrización de Reglas de Negocio

**Versión:** 1.0
**Fecha:** 2025-11-10
**Objetivo:** Sistema flexible de configuración de reglas dinámicas

---

## 1. Introducción

Este documento define el catálogo de **parámetros configurables** que controlan el comportamiento del sistema sin necesidad de cambiar código. Incluye:

- Límites operacionales
- Reglas de publicación de sorteos
- Políticas de pago y comisiones
- Thresholds de seguridad

---

## 2. Arquitectura de Parametrización

### 2.1 Tabla de Configuración

```sql
CREATE TABLE system_parameters (
    id BIGSERIAL PRIMARY KEY,
    key VARCHAR(100) UNIQUE NOT NULL,
    value JSONB NOT NULL,
    data_type VARCHAR(20) NOT NULL, -- 'integer', 'decimal', 'boolean', 'string', 'json'
    description TEXT,
    version INT NOT NULL DEFAULT 1,
    updated_by BIGINT REFERENCES users(id),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Historial de cambios
CREATE TABLE parameter_history (
    id BIGSERIAL PRIMARY KEY,
    parameter_key VARCHAR(100) NOT NULL,
    old_value JSONB,
    new_value JSONB,
    changed_by BIGINT REFERENCES users(id),
    changed_at TIMESTAMP NOT NULL DEFAULT NOW()
);
```

---

### 2.2 Acceso a Parámetros (Backend)

```go
type ParameterService struct {
    cache *redis.Client
    repo  ParameterRepository
}

func (s *ParameterService) GetInt(ctx context.Context, key string) (int, error) {
    // 1. Intentar obtener de caché
    cached, err := s.cache.Get(ctx, fmt.Sprintf("param:%s", key)).Result()
    if err == nil {
        return strconv.Atoi(cached)
    }

    // 2. Obtener de DB
    param, err := s.repo.FindByKey(ctx, key)
    if err != nil {
        return 0, err
    }

    // 3. Cachear (5 min)
    s.cache.Set(ctx, fmt.Sprintf("param:%s", key), param.Value, 5*time.Minute)

    return param.AsInt(), nil
}

// Uso
maxActiveRaffles, _ := paramService.GetInt(ctx, "raffle.max_active_per_user")
```

---

## 3. Catálogo de Parámetros

### 3.1 Sorteos (Raffles)

| Key | Tipo | Default | Descripción |
|-----|------|---------|-------------|
| `raffle.max_active_per_user` | int | 10 | Máximo de sorteos activos por usuario |
| `raffle.min_sale_percentage` | decimal | 0.7 | Mínimo 70% vendido para sortear |
| `raffle.allow_reschedule` | boolean | true | Permitir reprogramar fecha de sorteo |
| `raffle.max_reschedules` | int | 2 | Máximo de reprogramaciones |
| `raffle.days_before_draw_to_freeze` | int | 1 | Días antes del sorteo para congelar ediciones |
| `raffle.min_price_per_number` | decimal | 1.00 | Precio mínimo por boleto |
| `raffle.max_price_per_number` | decimal | 1000.00 | Precio máximo por boleto |
| `raffle.min_numbers` | int | 10 | Mínimo de números en un sorteo |
| `raffle.max_numbers` | int | 1000 | Máximo de números en un sorteo |
| `raffle.auto_cancel_if_under_min` | boolean | true | Cancelar auto si < min_sale_percentage al draw_date |

---

### 3.2 Reservas y Compras

| Key | Tipo | Default | Descripción |
|-----|------|---------|-------------|
| `reservation.ttl_minutes` | int | 5 | Tiempo de vida de reserva temporal |
| `reservation.max_numbers_per_user` | int | 10 | Máximo de números por compra |
| `reservation.allow_multiple_reservations` | boolean | true | Usuario puede tener múltiples reservas activas |

---

### 3.3 Pagos y Comisiones

| Key | Tipo | Default | Descripción |
|-----|------|---------|-------------|
| `payment.platform_fee_percentage` | decimal | 0.05 | Comisión de plataforma (5%) |
| `payment.min_platform_fee` | decimal | 0.50 | Comisión mínima en USD |
| `payment.max_platform_fee` | decimal | 100.00 | Comisión máxima en USD |
| `payment.retry_attempts` | int | 3 | Reintentos automáticos en fallo |
| `payment.refund_window_days` | int | 7 | Días para solicitar reembolso |

---

### 3.4 KYC y Verificación

| Key | Tipo | Default | Descripción |
|-----|------|---------|-------------|
| `kyc.min_level_to_create_raffle` | string | email_verified | Nivel mínimo de KYC |
| `kyc.min_level_to_buy` | string | email_verified | Nivel mínimo para comprar |
| `kyc.require_phone_for_winner` | boolean | true | Ganador debe tener teléfono verificado |
| `kyc.auto_verify_email_trusted_domains` | json | ["gmail.com"] | Dominios de confianza |

---

### 3.5 Límites de Seguridad

| Key | Tipo | Default | Descripción |
|-----|------|---------|-------------|
| `security.max_login_attempts` | int | 5 | Intentos de login antes de bloqueo |
| `security.lockout_duration_minutes` | int | 30 | Duración de bloqueo |
| `security.session_timeout_minutes` | int | 60 | Timeout de sesión inactiva |
| `security.require_2fa_for_admins` | boolean | true | 2FA obligatorio para admins |

---

### 3.6 Notificaciones

| Key | Tipo | Default | Descripción |
|-----|------|---------|-------------|
| `notification.send_email_on_purchase` | boolean | true | Email al comprar |
| `notification.send_sms_to_winner` | boolean | true | SMS al ganador |
| `notification.reminder_hours_before_draw` | int | 24 | Recordatorio antes del sorteo |

---

## 4. Validaciones Pre-Publicación

### 4.1 Matriz de Validaciones

Al publicar un sorteo, validar:

| Regla | Parámetro | Acción si Falla |
|-------|-----------|-----------------|
| DrawDate es futuro | - | Rechazar |
| DrawDate > now + 7 días | `raffle.min_days_to_draw` | Rechazar |
| Precio >= mínimo | `raffle.min_price_per_number` | Rechazar |
| Total números en rango | `raffle.min_numbers`, `raffle.max_numbers` | Rechazar |
| Usuario < max sorteos activos | `raffle.max_active_per_user` | Rechazar |
| Usuario KYC >= mínimo | `kyc.min_level_to_create_raffle` | Rechazar |
| Tiene al menos 1 imagen | - | Rechazar |

**Implementación:**
```go
func (uc *PublishRaffleUseCase) Validate(ctx context.Context, raffle *Raffle) error {
    params := uc.paramService

    // DrawDate futuro
    if raffle.DrawDate.Before(time.Now()) {
        return errors.New("draw_date debe ser futuro")
    }

    // Precio en rango
    minPrice, _ := params.GetDecimal(ctx, "raffle.min_price_per_number")
    maxPrice, _ := params.GetDecimal(ctx, "raffle.max_price_per_number")
    if raffle.PricePerNumber.LessThan(minPrice) || raffle.PricePerNumber.GreaterThan(maxPrice) {
        return fmt.Errorf("precio debe estar entre %s y %s", minPrice, maxPrice)
    }

    // Máximo sorteos activos
    maxActive, _ := params.GetInt(ctx, "raffle.max_active_per_user")
    activeCount := uc.raffleRepo.CountActive(ctx, raffle.UserID)
    if activeCount >= maxActive {
        return fmt.Errorf("máximo %d sorteos activos alcanzado", maxActive)
    }

    // KYC level
    minKYC, _ := params.GetString(ctx, "kyc.min_level_to_create_raffle")
    user := uc.userRepo.FindByID(ctx, raffle.UserID)
    if user.KYCLevel < minKYC {
        return errors.New("debes verificar tu cuenta para publicar")
    }

    return nil
}
```

---

## 5. Reglas Dinámicas de Auto-Gestión

### 5.1 Auto-Cancelación por Baja Venta

**Regla:** Si al `draw_date` no se alcanzó `min_sale_percentage`, cancelar sorteo automáticamente.

```go
func (uc *CheckRaffleSalesUseCase) Execute(ctx context.Context) error {
    today := time.Now().Truncate(24 * time.Hour)
    raffles := uc.raffleRepo.FindByDrawDate(ctx, today)

    minSalePercentage, _ := uc.paramService.GetDecimal(ctx, "raffle.min_sale_percentage")

    for _, raffle := range raffles {
        salePercentage := decimal.NewFromInt(raffle.SoldCount).
            Div(decimal.NewFromInt(raffle.TotalNumbers))

        if salePercentage.LessThan(minSalePercentage) {
            // Cancelar sorteo
            raffle.Status = RaffleStatusCancelled
            uc.raffleRepo.Update(ctx, raffle)

            // Reembolsar a compradores
            uc.refundAllPurchases(ctx, raffle.ID)

            // Notificar
            uc.notifier.SendEmail(ctx, raffle.Owner.Email, "raffle_cancelled_low_sales", map[string]any{
                "raffle_id": raffle.ID,
                "sold_percentage": salePercentage.StringFixed(2),
            })
        }
    }

    return nil
}
```

**Cron:** Ejecutar diariamente a las 00:00

---

### 5.2 Auto-Reprogramación

**Regla:** Si `allow_reschedule=true` y venta < 50%, owner puede reprogramar (máx 2 veces).

**Validación:**
```go
func (uc *RescheduleRaffleUseCase) Validate(ctx context.Context, raffle *Raffle, newDate time.Time) error {
    allowReschedule, _ := uc.paramService.GetBool(ctx, "raffle.allow_reschedule")
    if !allowReschedule {
        return errors.New("reprogramación no permitida")
    }

    maxReschedules, _ := uc.paramService.GetInt(ctx, "raffle.max_reschedules")
    if raffle.RescheduleCount >= maxReschedules {
        return fmt.Errorf("máximo %d reprogramaciones alcanzadas", maxReschedules)
    }

    // Debe ser al menos 7 días en el futuro
    minDays := 7
    if newDate.Before(time.Now().AddDate(0, 0, minDays)) {
        return fmt.Errorf("nueva fecha debe ser al menos %d días en el futuro", minDays)
    }

    return nil
}
```

**Consecuencias:**
- Notificar a compradores
- Opción de solicitar reembolso (7 días)
- Incrementar `raffle.reschedule_count`

---

## 6. Cambio de Parámetros en Producción

### 6.1 Workflow de Cambio

1. **Admin** modifica parámetro en backoffice
2. Sistema guarda en `parameter_history`
3. Invalida caché Redis
4. Logs de auditoría registran cambio
5. Notificación a equipo de desarrollo (opcional)

**Implementación:**
```go
func (s *ParameterService) Update(ctx context.Context, key string, newValue interface{}, adminID int64) error {
    // Obtener valor anterior
    param, _ := s.repo.FindByKey(ctx, key)
    oldValue := param.Value

    // Actualizar
    param.Value = newValue
    param.Version++
    param.UpdatedBy = adminID
    s.repo.Update(ctx, param)

    // Guardar historial
    s.repo.CreateHistory(ctx, &ParameterHistory{
        ParameterKey: key,
        OldValue:     oldValue,
        NewValue:     newValue,
        ChangedBy:    adminID,
    })

    // Invalidar caché
    s.cache.Del(ctx, fmt.Sprintf("param:%s", key))

    // Log de auditoría
    logger.Warn("parameter_changed",
        zap.String("key", key),
        zap.Any("old", oldValue),
        zap.Any("new", newValue),
        zap.Int64("admin_id", adminID),
    )

    return nil
}
```

---

### 6.2 UI de Backoffice

```tsx
function ParametersPage() {
  const { data: parameters } = useQuery(['parameters'], fetchParameters)

  return (
    <Table>
      <TableHeader>
        <TableRow>
          <TableHead>Parámetro</TableHead>
          <TableHead>Valor Actual</TableHead>
          <TableHead>Descripción</TableHead>
          <TableHead>Acciones</TableHead>
        </TableRow>
      </TableHeader>
      <TableBody>
        {parameters.map((param) => (
          <TableRow key={param.key}>
            <TableCell className="font-mono">{param.key}</TableCell>
            <TableCell>
              <Badge>{param.value}</Badge>
            </TableCell>
            <TableCell className="text-sm text-neutral-600">
              {param.description}
            </TableCell>
            <TableCell>
              <Button
                variant="outline"
                size="sm"
                onClick={() => openEditDialog(param)}
              >
                Editar
              </Button>
            </TableCell>
          </TableRow>
        ))}
      </TableBody>
    </Table>
  )
}
```

---

## 7. Validación de Cambios

**Antes de guardar un cambio:**
- [ ] Validar tipo de dato (int, decimal, boolean, etc.)
- [ ] Validar rango (ej: `min_sale_percentage` debe estar entre 0 y 1)
- [ ] Simular impacto (ej: cambiar `max_active_per_user` a 5 → ¿cuántos usuarios se verían afectados?)
- [ ] Confirmar con admin antes de aplicar

---

## 8. Ejemplos de Uso

### 8.1 Escenario: Aumentar Comisión

**Parámetro:** `payment.platform_fee_percentage`
**Valor actual:** 0.05 (5%)
**Nuevo valor:** 0.07 (7%)

**Pasos:**
1. Admin cambia valor en backoffice
2. Sistema invalida caché
3. Nuevas compras usan 7%
4. Compras ya procesadas mantienen 5% (histórico)

---

### 8.2 Escenario: Deshabilitar Reprogramación Temporal

**Parámetro:** `raffle.allow_reschedule`
**Valor actual:** `true`
**Nuevo valor:** `false`

**Efecto:**
- Endpoint `PATCH /raffles/:id/reschedule` retorna 403
- UI oculta botón "Reprogramar"

---

## 9. Tests de Parametrización

```go
func TestParameterService(t *testing.T) {
    t.Run("obtener parámetro int", func(t *testing.T) {
        service := NewParameterService(...)
        value, err := service.GetInt(ctx, "raffle.max_active_per_user")
        assert.NoError(t, err)
        assert.Equal(t, 10, value)
    })

    t.Run("cachear parámetro", func(t *testing.T) {
        // Primera llamada → DB
        service.GetInt(ctx, "raffle.max_active_per_user")

        // Segunda llamada → Redis cache
        cached, _ := rdb.Get(ctx, "param:raffle.max_active_per_user").Result()
        assert.Equal(t, "10", cached)
    })

    t.Run("invalidar caché al actualizar", func(t *testing.T) {
        service.Update(ctx, "raffle.max_active_per_user", 15, adminID)

        // Verificar que caché fue invalidado
        _, err := rdb.Get(ctx, "param:raffle.max_active_per_user").Result()
        assert.Error(t, err) // cache miss
    })
}
```

---

## 10. Próximos Pasos

1. Crear migración de `system_parameters`
2. Seedear parámetros default
3. Implementar `ParameterService` en backend
4. UI de backoffice para gestionar parámetros
5. Documentar impacto de cada parámetro

---

**Actualizado:** 2025-11-10
