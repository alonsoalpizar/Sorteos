# Operación de Backoffice - Plataforma de Sorteos

**Versión:** 1.0
**Fecha:** 2025-11-10
**Rol:** Administrador (Almighty)

---

## 1. Introducción

Este documento describe las operaciones críticas del backoffice administrativo, incluyendo:

- Gestión de sorteos (suspender, activar, cancelar)
- Gestión de usuarios (verificar KYC, suspender cuentas)
- Liquidaciones y pagos
- Auditoría y reportes
- Configuración del sistema

---

## 2. Panel Principal (Dashboard)

### 2.1 KPIs en Tiempo Real

**Métricas clave:**
- Total de usuarios activos (MAU)
- Sorteos activos / completados hoy
- Ingresos del día / mes
- Tasa de conversión (reserva → pago)
- Pagos fallidos (últimas 24h)
- Disputas activas

**Implementación:**
```sql
-- Vista materializada para performance
CREATE MATERIALIZED VIEW admin_dashboard_kpis AS
SELECT
    (SELECT COUNT(*) FROM users WHERE status = 'active') AS active_users,
    (SELECT COUNT(*) FROM raffles WHERE status = 'active') AS active_raffles,
    (SELECT COALESCE(SUM(amount), 0) FROM payments WHERE status = 'succeeded' AND created_at >= CURRENT_DATE) AS revenue_today,
    (SELECT COUNT(*) FROM payments WHERE status = 'failed' AND created_at >= NOW() - INTERVAL '24 hours') AS failed_payments_24h,
    (SELECT COUNT(*) FROM payments WHERE status = 'disputed' AND created_at >= NOW() - INTERVAL '7 days') AS active_disputes
;

-- Refresh cada 5 minutos (cron job)
REFRESH MATERIALIZED VIEW admin_dashboard_kpis;
```

**Frontend:**
```tsx
function AdminDashboard() {
  const { data: kpis } = useQuery(['admin-kpis'], fetchKPIs, {
    refetchInterval: 60000, // 1 min
  })

  return (
    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
      <StatsCard
        title="Usuarios Activos"
        value={kpis.active_users}
        icon={<Users />}
      />
      <StatsCard
        title="Sorteos Activos"
        value={kpis.active_raffles}
        icon={<TrendingUp />}
      />
      <StatsCard
        title="Ingresos Hoy"
        value={`$${kpis.revenue_today.toFixed(2)}`}
        icon={<DollarSign />}
      />
      <StatsCard
        title="Pagos Fallidos (24h)"
        value={kpis.failed_payments_24h}
        variant={kpis.failed_payments_24h > 10 ? 'danger' : 'default'}
        icon={<AlertTriangle />}
      />
    </div>
  )
}
```

---

## 3. Gestión de Sorteos

### 3.1 Listado de Sorteos

**Filtros:**
- Estado (todos, draft, active, suspended, completed, cancelled)
- Rango de fechas (created_at, draw_date)
- Usuario (por ID o email)
- Monto recaudado (min-max)
- % vendido (min-max)

**Acciones disponibles:**
- Ver detalles
- Suspender / Activar
- Cancelar (con reembolso)
- Editar (admin override)
- Ver transacciones

**Endpoint:**
```go
// GET /admin/raffles
func ListRafflesAdmin(c *gin.Context) {
    filters := RaffleAdminFilters{
        Status:    c.Query("status"),
        UserID:    c.Query("user_id"),
        DrawDate:  c.Query("draw_date"),
        Page:      c.GetInt("page"),
        PageSize:  20,
    }

    raffles, total, err := raffleRepo.ListAdmin(c.Request.Context(), filters)
    if err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }

    c.JSON(200, gin.H{
        "data": raffles,
        "meta": gin.H{
            "total": total,
            "page": filters.Page,
            "page_size": filters.PageSize,
        },
    })
}
```

---

### 3.2 Suspender Sorteo

**Casos de uso:**
- Contenido inapropiado
- Fraude detectado
- Reportes de usuarios
- Violación de términos

**Flujo:**
1. Admin selecciona sorteo
2. Ingresa motivo de suspensión
3. Sistema actualiza `status = suspended`
4. Rechaza nuevas reservas/compras
5. Notifica al owner
6. Registra en audit_logs

**Implementación:**
```go
type SuspendRaffleInput struct {
    RaffleID int64
    Reason   string
    AdminID  int64
}

func (uc *SuspendRaffleUseCase) Execute(ctx context.Context, input SuspendRaffleInput) error {
    raffle, err := uc.raffleRepo.FindByID(ctx, input.RaffleID)
    if err != nil {
        return err
    }

    if raffle.Status == RaffleStatusSuspended {
        return errors.New("sorteo ya está suspendido")
    }

    raffle.Status = RaffleStatusSuspended
    raffle.SuspensionReason = input.Reason
    uc.raffleRepo.Update(ctx, raffle)

    // Auditoría
    uc.auditLogger.Log(ctx, AuditLog{
        UserID:     input.AdminID,
        Action:     "suspend_raffle",
        EntityType: "raffle",
        EntityID:   input.RaffleID,
        Metadata: map[string]interface{}{
            "reason": input.Reason,
        },
    })

    // Notificar owner
    uc.notifier.SendEmail(ctx, raffle.Owner.Email, "raffle_suspended", map[string]any{
        "raffle_title": raffle.Title,
        "reason":       input.Reason,
    })

    return nil
}
```

**Frontend:**
```tsx
function SuspendRaffleDialog({ raffleId }: { raffleId: number }) {
  const [reason, setReason] = useState('')
  const suspendMutation = useMutation(suspendRaffle)

  const handleSubmit = () => {
    suspendMutation.mutate({ raffle_id: raffleId, reason })
  }

  return (
    <Dialog>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Suspender Sorteo</DialogTitle>
        </DialogHeader>
        <Textarea
          placeholder="Motivo de la suspensión..."
          value={reason}
          onChange={(e) => setReason(e.target.value)}
          rows={4}
        />
        <DialogFooter>
          <Button variant="destructive" onClick={handleSubmit}>
            Suspender
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}
```

---

### 3.3 Cancelar Sorteo (con Reembolso)

**Casos:**
- Fraude confirmado
- Owner no puede cumplir
- Emergencia legal

**Flujo:**
1. Actualizar `status = cancelled`
2. Obtener todos los pagos del sorteo
3. Procesar reembolsos con PSP
4. Liberar números
5. Notificar a compradores

```go
func (uc *CancelRaffleWithRefundUseCase) Execute(ctx context.Context, raffleID int64, reason string) error {
    payments := uc.paymentRepo.FindByRaffleID(ctx, raffleID)

    for _, payment := range payments {
        if payment.Status == PaymentSucceeded {
            // Reembolsar
            err := uc.paymentProvider.Refund(ctx, payment.ExternalID, payment.Amount)
            if err != nil {
                logger.Error("refund failed", zap.Int64("payment_id", payment.ID), zap.Error(err))
                continue
            }

            payment.Status = PaymentRefunded
            uc.paymentRepo.Update(ctx, payment)
        }
    }

    // Actualizar sorteo
    raffle, _ := uc.raffleRepo.FindByID(ctx, raffleID)
    raffle.Status = RaffleStatusCancelled
    uc.raffleRepo.Update(ctx, raffle)

    // Notificar
    uc.notifier.SendBulkEmail(ctx, payments, "raffle_cancelled_refund", map[string]any{
        "raffle_title": raffle.Title,
        "reason":       reason,
    })

    return nil
}
```

---

## 4. Gestión de Usuarios

### 4.1 Listado de Usuarios

**Filtros:**
- Estado (active, suspended, banned)
- Nivel KYC (none, email_verified, phone_verified, full_kyc)
- Fecha de registro
- Sorteos publicados (> N)
- Compras realizadas (> N)

**Acciones:**
- Ver perfil completo
- Verificar KYC manualmente
- Suspender / Reactivar
- Banear (permanente)
- Ver historial de actividad

---

### 4.2 Verificar KYC Manualmente

**Cuando:** Usuario sube documentos (ID, selfie) para `full_kyc`.

**Flujo:**
1. Admin revisa documentos subidos
2. Verifica autenticidad
3. Aprueba o rechaza
4. Actualiza `user.kyc_level`
5. Notifica al usuario

```go
type VerifyKYCInput struct {
    UserID   int64
    Approved bool
    Reason   string // si rechazado
    AdminID  int64
}

func (uc *VerifyKYCUseCase) Execute(ctx context.Context, input VerifyKYCInput) error {
    user, _ := uc.userRepo.FindByID(ctx, input.UserID)

    if input.Approved {
        user.KYCLevel = KYCFullVerified
        user.KYCVerifiedAt = time.Now()
        user.KYCVerifiedBy = input.AdminID
    } else {
        user.KYCLevel = KYCEmailVerified // downgrade
        user.KYCRejectionReason = input.Reason
    }

    uc.userRepo.Update(ctx, user)

    // Notificar
    template := "kyc_approved"
    if !input.Approved {
        template = "kyc_rejected"
    }

    uc.notifier.SendEmail(ctx, user.Email, template, map[string]any{
        "reason": input.Reason,
    })

    return nil
}
```

---

### 4.3 Suspender Usuario

**Razones:**
- Actividad fraudulenta
- Múltiples reportes
- Violación de términos

**Consecuencias:**
- No puede publicar sorteos
- No puede comprar boletos
- Sorteos activos se suspenden
- Sesión cerrada automáticamente

```go
func (uc *SuspendUserUseCase) Execute(ctx context.Context, userID int64, reason string, adminID int64) error {
    user, _ := uc.userRepo.FindByID(ctx, userID)
    user.Status = UserStatusSuspended
    user.SuspensionReason = reason
    uc.userRepo.Update(ctx, user)

    // Suspender sorteos activos del usuario
    raffles := uc.raffleRepo.FindActiveByUser(ctx, userID)
    for _, raffle := range raffles {
        raffle.Status = RaffleStatusSuspended
        uc.raffleRepo.Update(ctx, raffle)
    }

    // Revocar tokens de sesión
    uc.tokenManager.RevokeAllTokens(ctx, userID)

    // Auditoría
    uc.auditLogger.Log(ctx, AuditLog{
        UserID:     adminID,
        Action:     "suspend_user",
        EntityType: "user",
        EntityID:   userID,
        Metadata: map[string]interface{}{
            "reason": reason,
        },
    })

    return nil
}
```

---

## 5. Liquidaciones (Settlements)

### 5.1 Flujo de Liquidación

**Trigger:** Sorteo completado + ganador confirmó recepción de premio.

**Pasos:**
1. Calcular monto bruto (total recaudado)
2. Descontar fees (Stripe + plataforma)
3. Calcular neto
4. Crear registro de liquidación
5. Esperar confirmación de admin
6. Procesar payout (Stripe Connect o manual)

```go
func (uc *CreateSettlementUseCase) Execute(ctx context.Context, raffleID int64) (*Settlement, error) {
    raffle, _ := uc.raffleRepo.FindByID(ctx, raffleID)
    payments := uc.paymentRepo.FindByRaffleID(ctx, raffleID)

    grossAmount := decimal.Zero
    stripeFees := decimal.Zero

    for _, p := range payments {
        if p.Status == PaymentSucceeded {
            grossAmount = grossAmount.Add(p.Amount)
            // Stripe: 2.9% + $0.30
            fee := p.Amount.Mul(decimal.NewFromFloat(0.029)).Add(decimal.NewFromFloat(0.30))
            stripeFees = stripeFees.Add(fee)
        }
    }

    platformFeePercentage, _ := uc.paramService.GetDecimal(ctx, "payment.platform_fee_percentage")
    platformFee := grossAmount.Mul(platformFeePercentage)

    netAmount := grossAmount.Sub(stripeFees).Sub(platformFee)

    settlement := &Settlement{
        RaffleID:    raffleID,
        UserID:      raffle.UserID,
        GrossAmount: grossAmount,
        StripeFees:  stripeFees,
        PlatformFee: platformFee,
        NetAmount:   netAmount,
        Status:      SettlementPending,
    }

    uc.settlementRepo.Create(ctx, settlement)

    return settlement, nil
}
```

---

### 5.2 Aprobar Liquidación (Admin)

**Verificaciones:**
1. Ganador confirmó recepción de premio
2. No hay disputas activas
3. Owner tiene cuenta bancaria vinculada

```go
func (uc *ApproveSettlementUseCase) Execute(ctx context.Context, settlementID int64, adminID int64) error {
    settlement, _ := uc.settlementRepo.FindByID(ctx, settlementID)

    // Verificar que ganador confirmó
    raffle, _ := uc.raffleRepo.FindByID(ctx, settlement.RaffleID)
    if !raffle.WinnerConfirmedReceipt {
        return errors.New("ganador no ha confirmado recepción del premio")
    }

    // Verificar disputas
    disputes := uc.paymentRepo.FindDisputesByRaffleID(ctx, raffle.ID)
    if len(disputes) > 0 {
        return errors.New("hay disputas activas, no se puede liquidar")
    }

    // Procesar payout (Stripe Connect o manual)
    user, _ := uc.userRepo.FindByID(ctx, settlement.UserID)
    if user.StripeConnectAccountID != "" {
        err := uc.paymentProvider.Payout(ctx, user.StripeConnectAccountID, settlement.NetAmount)
        if err != nil {
            return err
        }
    }

    settlement.Status = SettlementApproved
    settlement.ApprovedBy = adminID
    settlement.ApprovedAt = time.Now()
    uc.settlementRepo.Update(ctx, settlement)

    // Notificar owner
    uc.notifier.SendEmail(ctx, user.Email, "settlement_approved", map[string]any{
        "amount": settlement.NetAmount.StringFixed(2),
    })

    return nil
}
```

---

## 6. Auditoría y Logs

### 6.1 Eventos Auditables

**Críticos (requieren aprobación):**
- Suspender / banear usuario
- Cancelar sorteo con reembolso
- Aprobar liquidación
- Cambiar parámetros del sistema

**Informativos:**
- Login de admin
- Visualización de datos sensibles (PII)
- Cambios en sorteos
- Verificación de KYC

**Tabla:**
```sql
SELECT
    al.id,
    u.email AS admin_email,
    al.action,
    al.entity_type,
    al.entity_id,
    al.ip_address,
    al.metadata,
    al.created_at
FROM audit_logs al
JOIN users u ON al.user_id = u.id
WHERE al.created_at >= NOW() - INTERVAL '7 days'
ORDER BY al.created_at DESC;
```

---

### 6.2 Visualización en Backoffice

```tsx
function AuditLogsPage() {
  const { data: logs } = useQuery(['audit-logs'], fetchAuditLogs)

  return (
    <Table>
      <TableHeader>
        <TableRow>
          <TableHead>Fecha</TableHead>
          <TableHead>Admin</TableHead>
          <TableHead>Acción</TableHead>
          <TableHead>Entidad</TableHead>
          <TableHead>IP</TableHead>
        </TableRow>
      </TableHeader>
      <TableBody>
        {logs.map((log) => (
          <TableRow key={log.id}>
            <TableCell>{formatDate(log.created_at)}</TableCell>
            <TableCell>{log.admin_email}</TableCell>
            <TableCell>
              <Badge variant={getActionVariant(log.action)}>
                {log.action}
              </Badge>
            </TableCell>
            <TableCell>
              {log.entity_type} #{log.entity_id}
            </TableCell>
            <TableCell className="font-mono text-xs">
              {log.ip_address}
            </TableCell>
          </TableRow>
        ))}
      </TableBody>
    </Table>
  )
}
```

---

## 7. Reportes

### 7.1 Reporte de Ingresos

**Filtros:**
- Rango de fechas
- Grupo (por día, semana, mes)
- PSP (Stripe, PayPal, etc.)

**Métricas:**
- Ingresos brutos
- Fees totales
- Ingresos netos
- Número de transacciones

```sql
SELECT
    DATE_TRUNC('day', created_at) AS date,
    provider,
    COUNT(*) AS transaction_count,
    SUM(amount) AS gross_amount,
    SUM(amount * 0.029 + 0.30) AS fees,
    SUM(amount * 0.95) AS net_amount
FROM payments
WHERE status = 'succeeded'
    AND created_at BETWEEN '2025-01-01' AND '2025-01-31'
GROUP BY date, provider
ORDER BY date DESC;
```

---

### 7.2 Reporte de Conversión

**Métricas:**
- Reservas creadas
- Reservas confirmadas (pagadas)
- Reservas expiradas
- Tasa de conversión

```sql
SELECT
    COUNT(*) FILTER (WHERE status = 'pending') AS reservations_created,
    COUNT(*) FILTER (WHERE status = 'confirmed') AS reservations_confirmed,
    COUNT(*) FILTER (WHERE status = 'expired') AS reservations_expired,
    ROUND(COUNT(*) FILTER (WHERE status = 'confirmed')::NUMERIC / COUNT(*) * 100, 2) AS conversion_rate
FROM reservations
WHERE created_at >= NOW() - INTERVAL '30 days';
```

---

## 8. Configuración del Sistema

### 8.1 Gestión de Parámetros

Ver: [parametrizacion_reglas.md](./parametrizacion_reglas.md)

**Acciones:**
- Editar valores
- Ver historial de cambios
- Restaurar valor anterior

---

### 8.2 Gestión de Notificaciones

**Templates de email:**
- Editar contenido (Markdown)
- Previsualizar
- Enviar email de prueba

---

## 9. Alertas y Monitoreo

### 9.1 Alertas Críticas

**Configurar en Prometheus:**
- Pagos fallidos > 20% en 1 hora
- Reservas con doble venta (debe ser 0)
- Tiempo de respuesta > 2s (p95)
- Disputas > 5 en 24 horas

**Notificación:**
- Email a equipo técnico
- Slack webhook
- Dashboard con indicador rojo

---

## 10. Checklist de Operación Diaria

- [ ] Revisar KPIs del dashboard
- [ ] Verificar pagos fallidos (> 10 → investigar)
- [ ] Aprobar KYC pendientes
- [ ] Revisar reportes de usuarios
- [ ] Aprobar liquidaciones pendientes
- [ ] Revisar logs de auditoría (acciones sospechosas)
- [ ] Verificar alertas de Prometheus

---

**Actualizado:** 2025-11-10
