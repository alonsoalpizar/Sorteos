# 📧 Sistema de Emails - Sorteos.club

## Resumen del Sistema

Tu aplicación ya cuenta con un **sistema completo de notificaciones por email** integrado con Go (backend).

---

## 🏗️ Arquitectura del Sistema

### 1. **Infraestructura de Correo**

✅ **Servidor SMTP Propio**: `mail.sorteos.club` (62.171.188.255)
- Postfix 3.8.6 (MTA)
- Dovecot 2.3.21 (MDA/IMAP)
- OpenDKIM con firma de emails
- SPF, DKIM, DMARC configurados
- SSL/TLS con Let's Encrypt

✅ **Webmail**: https://webmail.sorteos.club
- SnappyMail 2.38.2
- Tema DarkShine
- Idioma: Español
- Admin: `admin` / `Admin2025!`

---

## 📂 Estructura de Código

### **Backend (Go)**

```
/opt/Sorteos/backend/
├── internal/adapters/notifier/
│   ├── notifier.go              # Interface del notifier
│   ├── smtp.go                  # Implementación SMTP ✅
│   ├── sendgrid.go              # Implementación SendGrid (deprecado)
│   ├── template_loader.go       # Cargador de plantillas con embed
│   └── templates/
│       ├── verification.html    # ✅ Email de verificación
│       ├── welcome.html         # ✅ Email de bienvenida
│       ├── password_reset.html  # ✅ Reset de contraseña
│       └── purchase_confirmation.html # ✅ Confirmación de compra
│
└── pkg/config/
    └── smtp.go                  # Configuración SMTP
```

---

## 📧 Plantillas Disponibles

### 1. **verification.html** - Verificación de Cuenta
**Variables:**
- `{{.FirstName}}` - Nombre del usuario
- `{{.Code}}` - Código de verificación de 6 dígitos
- `{{.FrontendURL}}` - URL del frontend
- `{{.VerificationURL}}` - Link directo de verificación (opcional)

**Uso:**
```go
data := VerificationEmailData{
    FirstName: "Juan",
    Code: "123456",
    FrontendURL: "https://sorteos.club",
}
```

**Características:**
- Header azul (#3B82F6)
- Código en fuente monoespaciada grande
- Advertencia de expiración (15 min)
- Responsive

---

### 2. **welcome.html** - Bienvenida Post-Verificación
**Variables:**
- `{{.FirstName}}` - Nombre del usuario
- `{{.FrontendURL}}` - URL del frontend

**Uso:**
```go
data := WelcomeEmailData{
    FirstName: "Juan",
    FrontendURL: "https://sorteos.club",
}
```

**Características:**
- Header verde (#10B981)
- Lista de features disponibles
- CTA "Explorar Sorteos"
- Tips de uso

---

### 3. **password_reset.html** - Reset de Contraseña
**Variables:**
- `{{.FirstName}}` - Nombre del usuario (opcional)
- `{{.ResetURL}}` - Link de reset con token
- `{{.FrontendURL}}` - URL del frontend

**Uso:**
```go
data := PasswordResetEmailData{
    FirstName: "Juan",
    ResetURL: "https://sorteos.club/reset-password?token=xyz",
    FrontendURL: "https://sorteos.club",
}
```

**Características:**
- Header rojo (#EF4444)
- Warning de seguridad
- CTA "Restablecer Contraseña"
- Link alternativo (fallback)
- Advertencia de expiración (1 hora)

---

### 4. **purchase_confirmation.html** - Confirmación de Compra
**Variables:**
- `{{.FirstName}}` - Nombre del usuario
- `{{.RaffleTitle}}` - Nombre del sorteo
- `{{.RaffleID}}` - ID del sorteo
- `{{.Numbers}}` - Slice de números comprados
- `{{.TotalAmount}}` - Monto total pagado
- `{{.DrawDate}}` - Fecha del sorteo
- `{{.Prize}}` - Descripción del premio
- `{{.FrontendURL}}` - URL del frontend

**Uso:**
```go
data := PurchaseConfirmationData{
    FirstName: "Juan",
    RaffleTitle: "MacBook Pro M3",
    RaffleID: "1234",
    Numbers: []string{"00042", "00043"},
    TotalAmount: "$50.00",
    DrawDate: "25 de Diciembre, 2025",
    Prize: "MacBook Pro M3 14\"",
    FrontendURL: "https://sorteos.club",
}
```

---

## ⚙️ Configuración Actual (.env)

```env
CONFIG_EMAIL_PROVIDER=smtp
CONFIG_SMTP_HOST=mail.sorteos.club
CONFIG_SMTP_PORT=587
CONFIG_SMTP_USERNAME=noreply@sorteos.club
CONFIG_SMTP_PASSWORD=9NhNlT4m6FqUbM28FSFuSg==
CONFIG_SMTP_FROM_EMAIL=noreply@sorteos.club
CONFIG_SMTP_FROM_NAME=Plataforma de Sorteos
CONFIG_SMTP_USE_TLS=true
CONFIG_SMTP_USE_STARTTLS=true
CONFIG_SMTP_SKIP_VERIFY=false
CONFIG_FRONTEND_URL=https://sorteos.club
```

---

## 🚀 Cómo Usar el Sistema

### Ejemplo de envío desde el backend:

```go
// En tu handler o use case
notifier := // obtener instancia del notifier

// Enviar email de verificación
err := notifier.SendVerificationEmail(
    "usuario@example.com",
    "123456",
)

// Enviar email de bienvenida
err := notifier.SendWelcomeEmail(
    "usuario@example.com",
    "Juan",
)

// Enviar reset de contraseña
err := notifier.SendPasswordResetEmail(
    "usuario@example.com",
    "token_here",
)
```

---

## 📬 Usuarios de Email Configurados

| Email | Contraseña | Propósito |
|-------|------------|-----------|
| noreply@sorteos.club | 9NhNlT4m6FqUbM28FSFuSg== | Emails automáticos del sistema |
| info@sorteos.club | +yZ4o7A07toh/4MotrCqTw== | Consultas generales |
| soporte@sorteos.club | FQh7jA1Cuth1SP/+oBhopg== | Soporte técnico |
| postmaster@sorteos.club | YKiTy53jeer2LC/UZNripQ== | Administrador de correo |

**Credenciales completas en:** `/opt/Sorteos/mail-server-docs/mail-server-credentials.txt`

---

## 🛠️ Gestión de Usuarios

**Script interactivo:**
```bash
sudo /opt/Sorteos/scripts/manage-email-users.sh
```

**Funciones:**
1. Crear nuevos usuarios de email
2. Listar usuarios existentes
3. Cambiar contraseñas
4. Eliminar usuarios

---

## 🎨 Sistema de Plantillas

### Template Loader (Go Embed)

El sistema usa **Go embed** para incrustar las plantillas en el binario:

```go
//go:embed templates/*.html
var embeddedTemplates embed.FS
```

**Ventajas:**
- ✅ No necesitas copiar templates al servidor
- ✅ Todo está en el binario compilado
- ✅ Caching automático
- ✅ Fallback a filesystem si existe directorio

---

## 📊 Deliverability & Reputación

**Estado Actual:**
- ✅ SPF: PASS
- ✅ DKIM: PASS
- ✅ DMARC: PASS
- ⚠️ Reputación: Nueva (emails pueden ir a spam)

**Timeline de Mejora:**
- Semana 1-2: Emails van a spam (NORMAL)
- Semana 2-4: Mejora gradual
- Semana 4-8: Mayoría llega a inbox
- Mes 2-3: Reputación consolidada

**Acelerador de Reputación:**
- Warm-up: Empezar con 10-20 emails/día
- Incrementar gradualmente
- Mantener engagement alto (respuestas)
- Evitar quejas de spam

---

## 🧪 Testing

### Test desde Webmail:
1. Login: https://webmail.sorteos.club
2. Usuario: `noreply@sorteos.club`
3. Password: `9NhNlT4m6FqUbM28FSFuSg==`
4. Compose → Enviar email de prueba

### Test desde Backend:
```bash
cd /opt/Sorteos/backend
go run cmd/api/test_email.go
```

### Test con mail-tester.com:
```bash
bash /tmp/test-mail-tester.sh
```

---

## 📝 Próximas Plantillas Sugeridas

Plantillas que podrías crear según las necesidades de sorteos:

1. **raffle_created.html** - Confirmación de sorteo creado
2. **raffle_cancelled.html** - Cancelación de sorteo
3. **winner_notification.html** - Notificación de ganador
4. **draw_reminder.html** - Recordatorio de sorteo próximo
5. **payment_received.html** - Confirmación de pago
6. **referral_reward.html** - Premio por referido
7. **account_suspension.html** - Suspensión de cuenta
8. **monthly_summary.html** - Resumen mensual de actividad

---

## 📚 Recursos Adicionales

- [SnappyMail Docs](https://snappymail.eu/docs/)
- [Postfix Configuration](http://www.postfix.org/documentation.html)
- [Go Email Templates Best Practices](https://golang.org/pkg/html/template/)
- [Mail-Tester](https://www.mail-tester.com/)

---

## 🔐 Seguridad

- ✅ Contraseñas almacenadas de forma segura
- ✅ SMTP con autenticación
- ✅ TLS/STARTTLS habilitado
- ✅ DKIM firma todos los emails
- ✅ SPF protege contra spoofing
- ✅ DMARC monitorea deliverability

---

**Última actualización:** 18 de Noviembre, 2025
**Documentado por:** Claude Code
