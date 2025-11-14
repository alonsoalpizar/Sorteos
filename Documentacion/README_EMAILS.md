# Sistema de Emails - Sorteos Platform

## ✅ TL;DR - Resumen Ejecutivo

**Tu sistema YA TIENE emails funcionando**, solo necesitas:

1. **Opción A - SendGrid (5 min):**
   ```bash
   # Obtener API key de sendgrid.com
   # Agregar a .env:
   CONFIG_EMAIL_PROVIDER=sendgrid
   CONFIG_SENDGRID_API_KEY=SG.tu_api_key_aqui
   CONFIG_FRONTEND_URL=https://sorteos.club
   ```

2. **Opción B - Tu SMTP (ya implementado):**
   ```bash
   # Agregar a .env:
   CONFIG_EMAIL_PROVIDER=smtp
   CONFIG_SMTP_HOST=mail.sorteos.club
   CONFIG_SMTP_PORT=587
   CONFIG_SMTP_USERNAME=noreply@sorteos.club
   CONFIG_SMTP_PASSWORD=tu-password
   CONFIG_FRONTEND_URL=https://sorteos.club
   ```

3. **Probar:**
   ```bash
   cd /opt/Sorteos/backend
   ./test_email.sh  # Verifica configuración
   go build -o sorteos-api cmd/api/main.go
   sudo systemctl restart sorteos-api
   ```

---

## 📁 Archivos Creados

| Archivo | Descripción |
|---------|-------------|
| `internal/adapters/notifier/smtp.go` | Implementación SMTP completa |
| `internal/adapters/notifier/notifier.go` | Interface común |
| `pkg/config/config.go` | Configuración actualizada |
| `.env.smtp.example` | Ejemplos de configuración |
| `test_email.sh` | Script de verificación |
| `GUIA_EMAIL_SMTP_VS_SENDGRID.md` | Comparación detallada |
| `PROPUESTA_EMAILS.md` | 7 nuevos tipos de emails |
| `RESUMEN_IMPLEMENTACION_EMAIL.md` | Guía completa |
| `cmd/api/EJEMPLO_ROUTES_MODIFICADO.go` | Ejemplo de código |

---

## 🎯 ¿Qué Puedes Hacer Ahora?

### Emails Actuales (Ya Implementados)
1. ✅ **Verificación de email** - Código de 6 dígitos
2. ✅ **Bienvenida** - Post-verificación
3. ✅ **Reset password** - Link con token

### Nuevos Emails Propuestos
4. 🆕 **Confirmación de compra** - Cuando compran números
5. 🆕 **Notificación de ganador** - ¡Felicidades!
6. 🆕 **Sorteo completado** - Gracias por participar
7. 🆕 **Cancelación de sorteo** - Con reembolso
8. 🆕 **Reserva expirada** - Invitación a reintentar
9. 🆕 **Recordatorio 24h antes** - Con cron job
10. 🆕 **Resumen semanal** - Estadísticas

Ver: `PROPUESTA_EMAILS.md` para código completo.

---

## 📊 Comparación Rápida

| Feature | SendGrid | Tu SMTP |
|---------|----------|---------|
| Setup | 5 min | 1-2 horas |
| Costo | $0-20/mes | $0 |
| Deliverability | 99% | 75-85% |
| Analytics | ✅ | ❌ |
| Mantenimiento | 0 | Alto |

**Recomendación:** Empieza con SendGrid, migra a SMTP si lo necesitas después.

---

## 🚀 Quick Start

### 1. Verificar Configuración

```bash
cd /opt/Sorteos/backend
./test_email.sh sendgrid  # o ./test_email.sh smtp
```

### 2. Configurar .env

```bash
# SendGrid (recomendado para empezar)
CONFIG_EMAIL_PROVIDER=sendgrid
CONFIG_SENDGRID_API_KEY=SG.obtener_de_sendgrid.com
CONFIG_FRONTEND_URL=https://sorteos.club

# O SMTP (si ya tienes servidor)
CONFIG_EMAIL_PROVIDER=smtp
CONFIG_SMTP_HOST=mail.sorteos.club
CONFIG_SMTP_PORT=587
CONFIG_SMTP_USERNAME=noreply@sorteos.club
CONFIG_SMTP_PASSWORD=tu-password-seguro
CONFIG_SMTP_FROM_EMAIL=noreply@sorteos.club
CONFIG_SMTP_FROM_NAME=Plataforma de Sorteos
CONFIG_SMTP_USE_STARTTLS=true
CONFIG_FRONTEND_URL=https://sorteos.club
```

### 3. Modificar routes.go (Solo si usas SMTP)

Ver: `cmd/api/EJEMPLO_ROUTES_MODIFICADO.go`

Básicamente cambiar:
```go
sendgridNotifier := notifier.NewSendGridNotifier(&cfg.SendGrid, log)
```

Por:
```go
var emailNotifier notifier.Notifier

if cfg.EmailProvider == "smtp" {
    emailNotifier = notifier.NewSMTPNotifier(&cfg.SMTP, log)
} else {
    emailNotifier = notifier.NewSendGridNotifier(&cfg.SendGrid, log)
}
```

### 4. Recompilar y Probar

```bash
go build -o sorteos-api cmd/api/main.go
sudo systemctl restart sorteos-api

# Probar registro
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "Password123!@#",
    "accepted_terms": true,
    "accepted_privacy": true
  }'
```

---

## 📚 Documentación Completa

1. **Guía de Decisión:** `GUIA_EMAIL_SMTP_VS_SENDGRID.md`
   - Comparación detallada
   - Pros y contras
   - Configuración de servidor SMTP
   - DNS (SPF, DKIM, DMARC)

2. **Implementación:** `RESUMEN_IMPLEMENTACION_EMAIL.md`
   - Paso a paso SendGrid
   - Paso a paso SMTP
   - Troubleshooting
   - Checklist completo

3. **Nuevas Features:** `PROPUESTA_EMAILS.md`
   - 7 nuevos emails
   - Código de ejemplo
   - Sistema de cron jobs
   - Workers y scheduler

4. **Ejemplos de Config:** `.env.smtp.example`
   - Gmail, Office 365, AWS SES
   - Mailgun, Zoho, Mailtrap
   - Tu propio servidor

---

## 🔧 Troubleshooting

### Emails no llegan

```bash
# 1. Verificar configuración
./test_email.sh

# 2. Ver logs
sudo journalctl -u sorteos-api -f

# 3. SendGrid: Verificar Activity en dashboard
# https://app.sendgrid.com/activity

# 4. SMTP: Ver cola de correo
mailq
sudo tail -f /var/log/mail.log
```

### Emails van a spam

**SendGrid:**
- Verificar dominio en SendGrid dashboard
- Configurar Sender Authentication

**SMTP:**
- Verificar SPF, DKIM, DMARC
- Enviar test a: check-auth@verifier.port25.com
- Revisar score en mail-tester.com

---

## 🎓 Recursos

- [SendGrid Docs](https://docs.sendgrid.com/)
- [SMTP RFC 5321](https://tools.ietf.org/html/rfc5321)
- [SPF Setup](https://www.dmarcanalyzer.com/spf/)
- [DKIM Setup](https://www.dmarcanalyzer.com/dkim/)
- [Mail Tester](https://www.mail-tester.com/)

---

## ❓ Preguntas Frecuentes

**¿Cuál uso: SendGrid o SMTP?**
- Empezando: SendGrid
- Ya tienes SMTP: Úsalo
- Producción grande: AWS SES

**¿Cuánto cuesta?**
- SendGrid Free: 100 emails/día
- SendGrid Essentials: $19.95/mes - 50K emails
- SMTP propio: $0 (si ya tienes)
- AWS SES: $0.10 por 1000 emails

**¿Puedo cambiar después?**
Sí, es transparente. Solo cambias `CONFIG_EMAIL_PROVIDER` en `.env`.

**¿Cómo obtengo API key de SendGrid?**
1. https://app.sendgrid.com/
2. Settings → API Keys → Create
3. Full Access
4. Copiar key

**¿Qué puerto SMTP usar?**
- 587 (STARTTLS) - Recomendado
- 465 (TLS directo) - Alternativa
- 25 (sin cifrado) - NO recomendado

---

## 🔐 Seguridad

- ✅ Nunca subas `.env` a Git
- ✅ Rota API keys/passwords regularmente
- ✅ Usa TLS/STARTTLS siempre
- ✅ Implementa rate limiting
- ✅ Valida emails antes de enviar
- ✅ Monitorea bounces y spam reports

---

## 📞 Soporte

¿Necesitas ayuda con?
- ✅ Configuración de DNS
- ✅ Setup de servidor SMTP
- ✅ Implementación de nuevos emails
- ✅ Sistema de cron jobs
- ✅ Troubleshooting

Solo pregunta! Puedo ayudarte paso a paso.

---

## ✨ Próximos Pasos

1. ✅ **Ahora:** Configura SendGrid o SMTP
2. 🔄 **Luego:** Implementa emails de sorteos
3. ⏰ **Después:** Agrega recordatorios automáticos
4. 📊 **Finalmente:** Dashboard de métricas

---

**Happy Emailing! 📧🚀**
