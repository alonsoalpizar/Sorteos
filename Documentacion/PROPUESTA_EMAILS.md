# Propuesta de Sistema de Emails Mejorado para Sorteos

## Estado Actual ✅

Tu sistema **YA TIENE** un sistema de emails funcional con SendGrid que incluye:

1. **Email de Verificación** - Código de 6 dígitos (15 min expiry)
2. **Email de Bienvenida** - Post-verificación exitosa
3. **Email de Reset Password** - Link con token (1 hora expiry)

**Archivo:** `/opt/Sorteos/backend/internal/adapters/notifier/sendgrid.go`

---

## Validación del Código Existente ✅

### ✅ Puntos Fuertes

1. **Arquitectura limpia:**
   - Interface `Notifier` bien definida
   - Separación de responsabilidades
   - Logging completo con zap logger

2. **Seguridad:**
   - Manejo de errores robusto
   - Verificación de status code de SendGrid
   - No expone información sensible en logs

3. **UX:**
   - Templates HTML responsive
   - Fallback a texto plano
   - Diseño profesional con colores corporativos (#3B82F6)

### ⚠️ Áreas de Mejora

1. **URLs Hardcodeadas:**
   ```go
   // Línea 115 - sendgrid.go
   resetURL := fmt.Sprintf("https://sorteos.com/reset-password?token=%s", token)

   // Línea 217 - sendgrid.go
   <a href="https://sorteos.com">Explorar Sorteos</a>
   ```
   **Solución:** Agregar `CONFIG_FRONTEND_URL` al `.env`

2. **API Key No Configurada:**
   ```bash
   # .env línea 46
   CONFIG_SENDGRID_API_KEY=SG.your_sendgrid_api_key_here
   ```
   **Acción requerida:** Obtener API key real de SendGrid

3. **Falta de Templates Específicos para Sorteos:**
   No hay emails para eventos de sorteos (compra, ganador, sorteo próximo, etc.)

---

## Nuevos Emails Propuestos 🚀

### 1. Email de Confirmación de Compra de Números

**Trigger:** Usuario compra números exitosamente
**Cuándo:** Después de payment_intent.succeeded webhook

```go
func (n *SendGridNotifier) SendPurchaseConfirmation(
    email string,
    raffleTitle string,
    numbers []string,
    totalAmount float64,
    drawDate time.Time,
) error
```

**Contenido:**
- ✅ Confirmación de números comprados
- ✅ Detalles del sorteo (título, fecha, premio)
- ✅ Resumen de pago
- ✅ Link para ver sorteo
- ✅ Recordatorio de fecha de sorteo

---

### 2. Email de Recordatorio de Sorteo (24h antes)

**Trigger:** Cron job que ejecuta 24 horas antes del draw_date
**Cuándo:** Diariamente a las 10:00 AM

```go
func (n *SendGridNotifier) SendRaffleReminder(
    email string,
    firstName string,
    raffleTitle string,
    numbers []string,
    drawDate time.Time,
) error
```

**Contenido:**
- ⏰ Recordatorio de sorteo mañana
- 🎟️ Números que el usuario tiene
- 📅 Hora exacta del sorteo
- 🔗 Link para ver sorteo en vivo

**Requerimiento:** Sistema de cron jobs o background worker

---

### 3. Email de Ganador

**Trigger:** Cuando se ejecuta el sorteo y hay un ganador
**Cuándo:** Inmediatamente después del sorteo

```go
func (n *SendGridNotifier) SendWinnerNotification(
    email string,
    firstName string,
    raffleTitle string,
    winnerNumber string,
    prize string,
) error
```

**Contenido:**
- 🎉 ¡Felicidades, ganaste!
- 🏆 Detalles del premio
- 📋 Instrucciones para reclamar premio
- 📞 Contacto de soporte

---

### 4. Email de Sorteo Completado (No Ganador)

**Trigger:** Cuando se ejecuta el sorteo y el usuario NO ganó
**Cuándo:** Inmediatamente después del sorteo

```go
func (n *SendGridNotifier) SendRaffleCompleted(
    email string,
    firstName string,
    raffleTitle string,
    winnerNumber string,
    userNumbers []string,
) error
```

**Contenido:**
- 😔 Gracias por participar
- 🎯 Número ganador revelado
- 🎟️ Tus números
- 🔄 Invitación a participar en otros sorteos

---

### 5. Email de Cancelación de Sorteo

**Trigger:** Cuando un admin cancela un sorteo
**Cuándo:** Inmediatamente al cambiar status a "cancelled"

```go
func (n *SendGridNotifier) SendRaffleCancellation(
    email string,
    firstName string,
    raffleTitle string,
    refundAmount float64,
    reason string,
) error
```

**Contenido:**
- ⚠️ Notificación de cancelación
- 💰 Confirmación de reembolso
- 📝 Razón de cancelación
- 🔗 Link para ver otros sorteos

---

### 6. Email de Reserva Expirada

**Trigger:** Cuando expira una reserva sin pago
**Cuándo:** Job de expiración detecta reserva vencida

```go
func (n *SendGridNotifier) SendReservationExpired(
    email string,
    firstName string,
    raffleTitle string,
    numbers []string,
) error
```

**Contenido:**
- ⏱️ Tu reserva expiró
- 🎟️ Números que tenías reservados
- 🔄 Link para volver a reservar
- ⚡ Mensaje de urgencia (números limitados)

---

### 7. Email Resumen Semanal (Opcional)

**Trigger:** Cron job semanal
**Cuándo:** Lunes a las 9:00 AM

```go
func (n *SendGridNotifier) SendWeeklySummary(
    email string,
    firstName string,
    stats UserWeeklyStats,
) error
```

**Contenido:**
- 📊 Sorteos en los que participó
- 🎯 Sorteos próximos
- 💰 Total gastado esta semana
- ⭐ Sorteos destacados

---

## Sistema de Workers/Ejecutores Recomendado

### Opción 1: **Cron Jobs con Systemd Timers** (Recomendado para tu setup actual)

Ya tienes servicios systemd, podemos agregar timers:

```bash
# /etc/systemd/system/sorteos-email-reminders.timer
[Unit]
Description=Sorteos - Email Reminders Daily Job

[Timer]
OnCalendar=daily
OnCalendar=10:00
Persistent=true

[Install]
WantedBy=timers.target
```

**Ventaja:** No requiere dependencias adicionales, integrado con tu infraestructura.

---

### Opción 2: **Go Cron Job con robfig/cron** (Integrado en tu backend)

Agregar un scheduler interno en Go:

```go
// cmd/api/scheduler.go
package main

import (
    "github.com/robfig/cron/v3"
)

func startScheduler(notifier *notifier.SendGridNotifier, raffleRepo domain.RaffleRepository) {
    c := cron.New()

    // Recordatorios diarios a las 10:00 AM
    c.AddFunc("0 10 * * *", func() {
        sendRaffleReminders(notifier, raffleRepo)
    })

    // Resumen semanal los lunes a las 9:00 AM
    c.AddFunc("0 9 * * 1", func() {
        sendWeeklySummaries(notifier)
    })

    c.Start()
}
```

**Ventaja:** Todo en un solo proceso, fácil de mantener.

---

### Opción 3: **Redis + Bull Queue** (Si escalarás mucho)

Para millones de usuarios eventualmente:

1. Instalar Redis (ya lo tienes)
2. Usar go-workers o similar
3. Encolar trabajos de email

**Ventaja:** Alta escalabilidad, procesamiento asíncrono, reintentos automáticos.

---

## Configuraciones Necesarias

### 1. Variables de Entorno Nuevas

```bash
# Frontend URL (para links en emails)
CONFIG_FRONTEND_URL=https://sorteos.club

# SendGrid API Key REAL
CONFIG_SENDGRID_API_KEY=SG.xxxxxxxxxxxxxxxxxxxxxxxxx

# Habilitar/deshabilitar emails específicos
CONFIG_EMAIL_PURCHASE_CONFIRMATION=true
CONFIG_EMAIL_RAFFLE_REMINDERS=true
CONFIG_EMAIL_WEEKLY_SUMMARY=false

# Límites de emails
CONFIG_EMAIL_RATE_LIMIT_PER_USER_DAILY=20
CONFIG_EMAIL_RATE_LIMIT_GLOBAL_HOURLY=1000
```

---

### 2. Obtener API Key de SendGrid

**Pasos:**
1. Ir a https://app.sendgrid.com/
2. Crear cuenta (gratis hasta 100 emails/día)
3. Settings → API Keys → Create API Key
4. Seleccionar "Full Access"
5. Copiar key y agregar a `.env`

**Plan gratuito:** 100 emails/día
**Plan Essentials:** $19.95/mes - 50,000 emails/mes
**Plan Pro:** $89.95/mes - 100,000 emails/mes

---

### 3. Verificación de Dominio (Opcional pero Recomendado)

Para evitar que emails caigan en spam:

1. Verificar dominio `sorteos.club` en SendGrid
2. Agregar registros DNS (SPF, DKIM, DMARC)
3. Configurar "Sender Authentication"

**Resultado:** Mayor deliverability (99% inbox vs 70% sin verificar)

---

## Plan de Implementación Recomendado

### Fase 1: **Configuración Básica** (30 minutos)
1. ✅ Obtener API key de SendGrid
2. ✅ Agregar `CONFIG_FRONTEND_URL` al `.env`
3. ✅ Modificar URLs hardcodeadas en `sendgrid.go`
4. ✅ Probar emails existentes

### Fase 2: **Emails Transaccionales** (2-3 horas)
1. ✅ Implementar `SendPurchaseConfirmation`
2. ✅ Implementar `SendWinnerNotification`
3. ✅ Implementar `SendRaffleCompleted`
4. ✅ Integrar en flujo de pagos y sorteos

### Fase 3: **Sistema de Recordatorios** (3-4 horas)
1. ✅ Implementar `SendRaffleReminder`
2. ✅ Implementar `SendReservationExpired`
3. ✅ Crear cron job con robfig/cron
4. ✅ Integrar con job de expiración de reservas

### Fase 4: **Emails Promocionales** (Opcional)
1. ✅ Implementar `SendWeeklySummary`
2. ✅ Agregar preferencias de usuario (unsubscribe)
3. ✅ Dashboard de métricas de emails

---

## Métricas y Monitoreo

### KPIs a Trackear:

1. **Open Rate** (tasa de apertura)
   - Target: >20%
   - SendGrid provee estadísticas

2. **Click Rate** (clicks en links)
   - Target: >5%
   - Usar UTM parameters

3. **Bounce Rate** (rebotes)
   - Target: <2%
   - Limpiar lista de emails inválidos

4. **Spam Complaints**
   - Target: <0.1%
   - Agregar link de unsubscribe

### Logging Recomendado:

```go
// Tabla: email_logs
type EmailLog struct {
    ID          int64
    UserID      int64
    EmailType   string  // "purchase_confirmation", "winner", etc
    RecipientEmail string
    SendGridMessageID string
    Status      string  // "sent", "delivered", "opened", "clicked", "bounced"
    SentAt      time.Time
    OpenedAt    *time.Time
    ClickedAt   *time.Time
}
```

---

## Ejemplo de Implementación Completa

### Nuevo Método: Email de Confirmación de Compra

```go
// Agregar a: internal/adapters/notifier/sendgrid.go

// PurchaseDetails contiene detalles de la compra
type PurchaseDetails struct {
    RaffleTitle string
    RaffleID    string
    Numbers     []string
    TotalAmount float64
    DrawDate    time.Time
    Prize       string
}

// SendPurchaseConfirmation envía confirmación de compra
func (n *SendGridNotifier) SendPurchaseConfirmation(
    email, firstName string,
    details *PurchaseDetails,
) error {
    to := mail.NewEmail(firstName, email)
    subject := fmt.Sprintf("¡Compra Confirmada! - %s", details.RaffleTitle)

    // Formatear números
    numbersStr := strings.Join(details.Numbers, ", ")

    // Formatear fecha
    drawDateStr := details.DrawDate.Format("02/01/2006 15:04")

    // Formatear monto
    amountStr := fmt.Sprintf("$%.2f", details.TotalAmount)

    plainTextContent := fmt.Sprintf(`
¡Hola %s!

Tu compra ha sido confirmada exitosamente.

DETALLES DEL SORTEO:
- Sorteo: %s
- Números: %s
- Monto: %s
- Fecha del sorteo: %s

¡Te deseamos mucha suerte!

Ver sorteo: %s/raffles/%s

Saludos,
Equipo de Sorteos Platform
    `, firstName, details.RaffleTitle, numbersStr, amountStr, drawDateStr,
        n.config.FrontendURL, details.RaffleID)

    htmlContent := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Compra Confirmada</title>
</head>
<body style="font-family: Arial, sans-serif; line-height: 1.6; color: #333; margin: 0; padding: 0;">
    <div style="max-width: 600px; margin: 0 auto; padding: 20px; background-color: #f7fafc;">
        <!-- Header -->
        <div style="background-color: #3B82F6; color: white; padding: 30px; text-align: center; border-radius: 8px 8px 0 0;">
            <h1 style="margin: 0; font-size: 28px;">✅ ¡Compra Confirmada!</h1>
        </div>

        <!-- Body -->
        <div style="background-color: white; padding: 30px; border-radius: 0 0 8px 8px; box-shadow: 0 2px 8px rgba(0,0,0,0.1);">
            <p style="font-size: 16px;">¡Hola <strong>%s</strong>!</p>
            <p>Tu compra ha sido <strong>confirmada exitosamente</strong>. Ya estás participando en el sorteo.</p>

            <!-- Raffle Details Card -->
            <div style="background-color: #EFF6FF; border-left: 4px solid #3B82F6; padding: 20px; margin: 25px 0; border-radius: 4px;">
                <h3 style="margin: 0 0 15px 0; color: #1E40AF;">%s</h3>
                <table style="width: 100%%; border-collapse: collapse;">
                    <tr>
                        <td style="padding: 8px 0; color: #64748B; font-size: 14px;">Tus números:</td>
                        <td style="padding: 8px 0; font-weight: bold; font-size: 18px; color: #3B82F6; text-align: right;">%s</td>
                    </tr>
                    <tr>
                        <td style="padding: 8px 0; color: #64748B; font-size: 14px;">Monto pagado:</td>
                        <td style="padding: 8px 0; font-weight: bold; font-size: 16px; text-align: right;">%s</td>
                    </tr>
                    <tr>
                        <td style="padding: 8px 0; color: #64748B; font-size: 14px;">Fecha del sorteo:</td>
                        <td style="padding: 8px 0; font-weight: bold; text-align: right; color: #059669;">📅 %s</td>
                    </tr>
                </table>
            </div>

            <!-- Prize Info -->
            <div style="background-color: #FEF3C7; border-left: 4px solid #F59E0B; padding: 15px; margin: 25px 0; border-radius: 4px;">
                <p style="margin: 0; color: #92400E; font-size: 14px;">🏆 <strong>Premio:</strong> %s</p>
            </div>

            <!-- CTA Button -->
            <div style="text-align: center; margin: 30px 0;">
                <a href="%s/raffles/%s" style="background-color: #3B82F6; color: white; padding: 14px 32px; text-decoration: none; border-radius: 6px; display: inline-block; font-weight: bold; font-size: 16px;">Ver Sorteo en Vivo</a>
            </div>

            <!-- Info Box -->
            <div style="background-color: #F0FDF4; border: 1px solid #BBF7D0; padding: 15px; border-radius: 4px; margin-top: 25px;">
                <p style="margin: 0; color: #065F46; font-size: 13px;">
                    💡 <strong>Tip:</strong> Te enviaremos un recordatorio 24 horas antes del sorteo. ¡Mantente atento!
                </p>
            </div>

            <p style="margin-top: 30px; color: #64748B; font-size: 14px;">
                ¡Te deseamos mucha suerte! 🍀
            </p>
        </div>

        <!-- Footer -->
        <div style="text-align: center; padding: 20px; color: #94A3B8; font-size: 12px;">
            <p style="margin: 5px 0;">Saludos,</p>
            <p style="margin: 5px 0;"><strong>Equipo de Sorteos Platform</strong></p>
            <p style="margin: 15px 0 5px 0;">
                <a href="%s" style="color: #3B82F6; text-decoration: none;">Inicio</a> |
                <a href="%s/profile" style="color: #3B82F6; text-decoration: none;">Mi Cuenta</a> |
                <a href="%s/support" style="color: #3B82F6; text-decoration: none;">Soporte</a>
            </p>
        </div>
    </div>
</body>
</html>
    `, firstName, details.RaffleTitle, numbersStr, amountStr, drawDateStr, details.Prize,
        n.config.FrontendURL, details.RaffleID,
        n.config.FrontendURL, n.config.FrontendURL, n.config.FrontendURL)

    message := mail.NewSingleEmail(n.fromMail, subject, to, plainTextContent, htmlContent)

    response, err := n.client.Send(message)
    if err != nil {
        n.logger.Error("Error sending purchase confirmation email",
            logger.String("email", email),
            logger.Error(err),
        )
        return err
    }

    if response.StatusCode >= 400 {
        n.logger.Error("SendGrid returned error",
            logger.String("email", email),
            logger.Int("status_code", response.StatusCode),
            logger.String("body", response.Body),
        )
        return fmt.Errorf("sendgrid error: %d - %s", response.StatusCode, response.Body)
    }

    n.logger.Info("Purchase confirmation email sent",
        logger.String("email", email),
        logger.String("raffle_id", details.RaffleID),
        logger.Int("status_code", response.StatusCode),
    )

    return nil
}
```

### Integración en Webhook de Stripe

```go
// Agregar en: cmd/api/payment_routes.go o donde manejes el webhook

func handlePaymentSuccess(payment *domain.Payment, reservation *domain.Reservation, notifier *notifier.SendGridNotifier) {
    // ... código existente de confirmación de pago ...

    // Obtener datos del raffle
    raffle, _ := raffleRepo.FindByUUID(reservation.RaffleID)
    user, _ := userRepo.FindByID(payment.UserID)

    // Enviar email de confirmación
    details := &notifier.PurchaseDetails{
        RaffleTitle: raffle.Title,
        RaffleID:    raffle.UUID,
        Numbers:     reservation.NumberIDs,
        TotalAmount: payment.Amount,
        DrawDate:    raffle.DrawDate,
        Prize:       raffle.PrizeDescription, // Asumiendo que tienes este campo
    }

    if err := notifier.SendPurchaseConfirmation(user.Email, user.FirstName, details); err != nil {
        log.Warn("Failed to send purchase confirmation email", logger.Error(err))
        // No falla la operación si el email no se envía
    }
}
```

---

## Checklist de Implementación

### Configuración Inicial
- [ ] Crear cuenta en SendGrid
- [ ] Obtener API Key
- [ ] Agregar `CONFIG_SENDGRID_API_KEY` real al `.env`
- [ ] Agregar `CONFIG_FRONTEND_URL=https://sorteos.club` al `.env`
- [ ] Verificar dominio en SendGrid (opcional pero recomendado)

### Correcciones al Código Existente
- [ ] Modificar `sendgrid.go:115` para usar `CONFIG_FRONTEND_URL`
- [ ] Modificar `sendgrid.go:217` para usar `CONFIG_FRONTEND_URL`
- [ ] Agregar campo `FrontendURL` a `SendGridConfig` struct
- [ ] Probar emails existentes (verificación, bienvenida, reset)

### Nuevos Emails (Prioridad Alta)
- [ ] Implementar `SendPurchaseConfirmation`
- [ ] Implementar `SendWinnerNotification`
- [ ] Implementar `SendRaffleCompleted`
- [ ] Integrar en webhook de Stripe/PayPal
- [ ] Integrar en flujo de ejecución de sorteo

### Nuevos Emails (Prioridad Media)
- [ ] Implementar `SendRaffleCancellation`
- [ ] Implementar `SendReservationExpired`
- [ ] Integrar en job de expiración de reservas

### Sistema de Workers
- [ ] Evaluar opción de scheduler (systemd timer vs cron vs robfig/cron)
- [ ] Implementar job de recordatorios de sorteo (24h antes)
- [ ] Implementar job de resumen semanal (opcional)

### Mejoras Avanzadas (Opcional)
- [ ] Crear tabla `email_logs` para tracking
- [ ] Implementar preferencias de usuario (unsubscribe)
- [ ] Agregar rate limiting de emails por usuario
- [ ] Dashboard de métricas de emails
- [ ] A/B testing de templates
- [ ] Templates con SendGrid Dynamic Templates (más visual)

---

## Preguntas Frecuentes

### ¿Necesito instalar algo más?
No. Ya tienes `github.com/sendgrid/sendgrid-go` en tu `go.mod`.

### ¿Cuánto cuesta SendGrid?
- **Gratis:** 100 emails/día para siempre
- **Essentials:** $19.95/mes - 50,000 emails
- **Pro:** $89.95/mes - 100,000 emails

### ¿Puedo usar otro proveedor en lugar de SendGrid?
Sí. Puedes implementar:
- **AWS SES** (más barato, $0.10 por 1000 emails)
- **Mailgun** (similar a SendGrid)
- **Postmark** (excelente para transaccionales)
- **SMTP directo** (si tienes servidor de correo)

Solo necesitas implementar la interface `Notifier`.

### ¿Los emails van a caer en spam?
Si verificas el dominio en SendGrid y configuras SPF/DKIM, tendrás ~99% de deliverability.

### ¿Puedo probar sin SendGrid?
Sí, usa `CONFIG_SKIP_EMAIL_VERIFICATION=true` en desarrollo. O usa MailTrap.io para testing.

---

## Recursos Adicionales

- [SendGrid API Docs](https://docs.sendgrid.com/api-reference/mail-send/mail-send)
- [Go SendGrid Library](https://github.com/sendgrid/sendgrid-go)
- [Email Best Practices](https://sendgrid.com/blog/email-best-practices/)
- [SPF/DKIM Setup](https://docs.sendgrid.com/ui/account-and-settings/how-to-set-up-domain-authentication)

---

## Contacto y Soporte

Si necesitas ayuda implementando cualquiera de estas mejoras, puedo:
1. Escribir el código completo de los nuevos emails
2. Configurar el sistema de cron jobs
3. Crear las migraciones de base de datos necesarias
4. Integrar con tu flujo de pagos y sorteos

¿Por dónde quieres empezar?
