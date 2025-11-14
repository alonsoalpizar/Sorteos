# Resumen: Implementación del Sistema de Emails

## 🎉 ¡Buenas Noticias!

**Tu sistema YA TIENE emails funcionando** con SendGrid. Solo necesitas configurarlo correctamente.

**AHORA TAMBIÉN PUEDES** usar tu propio servidor SMTP/MX si lo prefieres.

---

## Archivos Creados

### 1. **Adaptador SMTP Nuevo**
📁 `/opt/Sorteos/backend/internal/adapters/notifier/smtp.go`
- Implementa envío de emails con SMTP estándar
- Soporta TLS directo (puerto 465) y STARTTLS (puerto 587)
- Compatible con cualquier servidor SMTP

### 2. **Configuración Actualizada**
📁 `/opt/Sorteos/backend/pkg/config/config.go`
- Agregado: `SMTPConfig` struct
- Agregado: `EmailProvider` (seleccionar "sendgrid" o "smtp")
- Agregado: `FrontendURL` para links en emails

### 3. **Ejemplo de Configuración SMTP**
📁 `/opt/Sorteos/backend/.env.smtp.example`
- Ejemplos para Gmail, Office 365, AWS SES, Mailgun, etc.
- Configuración para tu propio servidor MX
- Instrucciones paso a paso

### 4. **Guía Completa de Comparación**
📁 `/opt/Sorteos/GUIA_EMAIL_SMTP_VS_SENDGRID.md`
- Comparación detallada SMTP vs SendGrid
- Ventajas y desventajas de cada opción
- Guía de configuración de servidor SMTP propio
- Configuración de DNS (SPF, DKIM, DMARC)
- Recomendaciones según tu caso de uso

### 5. **Propuesta de Emails Nuevos**
📁 `/opt/Sorteos/PROPUESTA_EMAILS.md`
- 7 nuevos tipos de emails para sorteos
- Sistema de cron jobs para recordatorios
- Ejemplos de código completo

---

## Opción 1: Usar SendGrid (Recomendado para Empezar)

### Paso 1: Obtener API Key

1. Ve a https://app.sendgrid.com/
2. Crear cuenta gratuita (100 emails/día)
3. Settings → API Keys → Create API Key
4. Selecciona "Full Access"
5. Copia la API Key

### Paso 2: Configurar .env

```bash
# En /opt/Sorteos/backend/.env

# Proveedor de email
CONFIG_EMAIL_PROVIDER=sendgrid

# URL del frontend
CONFIG_FRONTEND_URL=https://sorteos.club

# SendGrid
CONFIG_SENDGRID_API_KEY=SG.tu_api_key_real_aqui
CONFIG_SENDGRID_FROM_EMAIL=noreply@sorteos.club
CONFIG_SENDGRID_FROM_NAME=Plataforma de Sorteos
```

### Paso 3: Actualizar sendgrid.go (Opcional)

Actualizar las URLs hardcodeadas para usar `FrontendURL`:

```bash
# En internal/adapters/notifier/sendgrid.go
# Línea 115: Cambiar
resetURL := fmt.Sprintf("https://sorteos.com/reset-password?token=%s", token)

# Por:
resetURL := fmt.Sprintf("%s/reset-password?token=%s", n.config.FrontendURL, token)
```

### Paso 4: Verificar Dominio (Opcional pero Recomendado)

Para mejorar deliverability:
1. SendGrid → Settings → Sender Authentication
2. Authenticate Your Domain
3. Agregar registros DNS proporcionados
4. Esperar verificación

### Paso 5: Probar

```bash
# Recompilar backend
cd /opt/Sorteos/backend
go build -o sorteos-api cmd/api/main.go

# Reiniciar servicio
sudo systemctl restart sorteos-api

# Probar registro de usuario
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "tu-email@example.com",
    "password": "Password123!@#",
    "accepted_terms": true,
    "accepted_privacy": true
  }'

# Deberías recibir un email con el código de verificación
```

### Costos SendGrid

- **Gratis:** 100 emails/día
- **Essentials:** $19.95/mes - 50,000 emails/mes
- **Pro:** $89.95/mes - 100,000 emails/mes

---

## Opción 2: Usar Tu Propio SMTP/MX

### Paso 1: Verificar Servidor SMTP

```bash
# Verificar que tienes servidor SMTP funcionando
telnet mail.sorteos.club 587

# Debería conectar y mostrar:
# 220 mail.sorteos.club ESMTP Postfix
```

### Paso 2: Configurar .env

```bash
# En /opt/Sorteos/backend/.env

# Proveedor de email
CONFIG_EMAIL_PROVIDER=smtp

# URL del frontend
CONFIG_FRONTEND_URL=https://sorteos.club

# Tu servidor SMTP
CONFIG_SMTP_HOST=mail.sorteos.club
CONFIG_SMTP_PORT=587
CONFIG_SMTP_USERNAME=noreply@sorteos.club
CONFIG_SMTP_PASSWORD=tu-password-seguro
CONFIG_SMTP_FROM_EMAIL=noreply@sorteos.club
CONFIG_SMTP_FROM_NAME=Plataforma de Sorteos
CONFIG_SMTP_USE_TLS=false
CONFIG_SMTP_USE_STARTTLS=true
CONFIG_SMTP_SKIP_VERIFY=false
```

**Nota:** Usa los valores de tu configuración actual. Consulta `backend/.env.smtp.example` para más ejemplos.

### Paso 3: Modificar routes.go

Actualizar `/opt/Sorteos/backend/cmd/api/routes.go`:

```go
// Buscar la línea donde se crea sendgridNotifier (aprox línea 40)
// Reemplazar:
sendgridNotifier := notifier.NewSendGridNotifier(&cfg.SendGrid, log)

// Por:
var emailNotifier notifier.Notifier

if cfg.EmailProvider == "smtp" {
	emailNotifier = notifier.NewSMTPNotifier(&cfg.SMTP, log)
	log.Info("Using SMTP email provider", logger.String("host", cfg.SMTP.Host))
} else {
	emailNotifier = notifier.NewSendGridNotifier(&cfg.SendGrid, log)
	log.Info("Using SendGrid email provider")
}

// Luego cambiar todas las referencias de sendgridNotifier a emailNotifier
// Por ejemplo línea 47:
registerUseCase := auth.NewRegisterUseCase(userRepo, consentRepo, auditRepo, tokenMgr, emailNotifier, log, cfg.SkipEmailVerification)
```

### Paso 4: Crear Interface Notifier

Crear archivo `/opt/Sorteos/backend/internal/adapters/notifier/notifier.go`:

```go
package notifier

// Notifier es la interface para envío de notificaciones
type Notifier interface {
	SendVerificationEmail(email, code string) error
	SendPasswordResetEmail(email, token string) error
	SendWelcomeEmail(email, firstName string) error
}
```

### Paso 5: Recompilar y Probar

```bash
cd /opt/Sorteos/backend
go build -o sorteos-api cmd/api/main.go

sudo systemctl restart sorteos-api

# Ver logs
sudo journalctl -u sorteos-api -f

# Probar registro
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "Password123!@#",
    "accepted_terms": true,
    "accepted_privacy": true
  }'

# Verificar logs de mail
sudo tail -f /var/log/mail.log
```

### Paso 6: Verificar DNS (Importante para Deliverability)

```bash
# Verificar SPF
dig sorteos.club TXT

# Verificar DKIM
dig default._domainkey.sorteos.club TXT

# Verificar DMARC
dig _dmarc.sorteos.club TXT

# Verificar MX
dig sorteos.club MX

# Verificar reverse DNS
dig -x 62.171.188.255
```

Si alguno falla, consulta `GUIA_EMAIL_SMTP_VS_SENDGRID.md` sección "Configurar DNS".

### Paso 7: Test de Spam Score

Envía un email a: `check-auth@verifier.port25.com`

Recibirás un reporte indicando si pasa SPF, DKIM, DMARC y el spam score.

---

## Comparación Rápida

| Criterio | SendGrid | Tu SMTP |
|----------|----------|---------|
| **Setup** | 5 minutos | 1-2 horas |
| **Costo** | $0-20/mes | $0 (si ya tienes) |
| **Deliverability** | 95-99% | 70-85% |
| **Analytics** | Sí | No |
| **Mantenimiento** | Cero | Alto |
| **Escalabilidad** | Ilimitada | Limitada |

---

## Recomendación

### Para Empezar (MVP):
**Usa SendGrid** plan gratuito
- Setup en 5 minutos
- 100 emails/día suficiente para empezar
- Alta deliverability
- Sin mantenimiento

### Si Ya Tienes Servidor SMTP Configurado:
**Usa tu SMTP**
- Aprovecha infraestructura existente
- Costo cero
- Requiere verificar DNS esté bien configurado

### Para Producción a Gran Escala (>50K emails/mes):
**AWS SES**
- $0.10 por 1000 emails
- Muy económico
- Escalable

---

## Próximos Pasos Opcionales

Una vez que tengas emails funcionando, puedes:

1. **Implementar nuevos emails de sorteos**
   - Email de confirmación de compra
   - Email de ganador
   - Recordatorios 24h antes

   Ver: `PROPUESTA_EMAILS.md`

2. **Agregar sistema de cron jobs**
   - Recordatorios automáticos
   - Resúmenes semanales

3. **Implementar tracking de emails**
   - Tabla `email_logs` en BD
   - Dashboard de métricas

4. **A/B Testing de templates**
   - Probar diferentes diseños
   - Optimizar open rate

---

## Troubleshooting

### Emails no llegan con SendGrid

```bash
# 1. Verificar API Key
echo $CONFIG_SENDGRID_API_KEY

# 2. Ver logs del backend
sudo journalctl -u sorteos-api -f

# 3. Verificar en SendGrid dashboard
# Activity → Email Activity
```

### Emails no llegan con SMTP

```bash
# 1. Verificar conectividad
telnet mail.sorteos.club 587

# 2. Ver cola de emails
mailq

# 3. Ver logs
sudo tail -f /var/log/mail.log

# 4. Test manual
echo "Test" | mail -s "Test" test@example.com

# 5. Verificar DNS
dig sorteos.club MX
dig sorteos.club TXT  # SPF

# 6. Verificar spam score
# Enviar a check-auth@verifier.port25.com
```

### Emails caen en spam

**Con SendGrid:**
1. Verificar dominio en SendGrid
2. Configurar Sender Authentication
3. Evitar palabras spam en subject/body

**Con SMTP:**
1. Configurar SPF, DKIM, DMARC
2. Verificar reverse DNS (PTR)
3. Construir reputación gradualmente
4. No enviar muchos emails de golpe

---

## Checklist de Implementación

### SendGrid
- [ ] Crear cuenta en SendGrid
- [ ] Obtener API Key
- [ ] Agregar a `.env`
- [ ] Configurar `CONFIG_FRONTEND_URL`
- [ ] Recompilar backend
- [ ] Probar registro de usuario
- [ ] Verificar dominio (opcional)

### SMTP Propio
- [ ] Verificar servidor SMTP funciona
- [ ] Configurar credenciales en `.env`
- [ ] Modificar `routes.go`
- [ ] Crear `notifier.go` interface
- [ ] Recompilar backend
- [ ] Verificar DNS (SPF, DKIM, DMARC)
- [ ] Test de spam score
- [ ] Probar registro de usuario

---

## Archivos a Modificar (Solo si usas SMTP)

1. `/opt/Sorteos/backend/.env`
2. `/opt/Sorteos/backend/cmd/api/routes.go`
3. `/opt/Sorteos/backend/internal/adapters/notifier/notifier.go` (crear)

---

## Comandos Útiles

```bash
# Recompilar backend
cd /opt/Sorteos/backend
go build -o sorteos-api cmd/api/main.go

# Reiniciar servicio
sudo systemctl restart sorteos-api

# Ver logs
sudo journalctl -u sorteos-api -f

# Verificar que está corriendo
sudo systemctl status sorteos-api

# Ver configuración
cat /opt/Sorteos/backend/.env | grep EMAIL
cat /opt/Sorteos/backend/.env | grep SMTP
cat /opt/Sorteos/backend/.env | grep SENDGRID

# Test de conectividad SMTP
telnet mail.sorteos.club 587

# Ver cola de emails (si usas SMTP)
mailq

# Ver logs de mail (si usas SMTP)
sudo tail -f /var/log/mail.log
```

---

## Soporte

Si necesitas ayuda con:
- Configuración de DNS
- Setup de servidor SMTP
- Implementación de nuevos emails
- Sistema de cron jobs
- Troubleshooting

Solo pregúntame! 🚀

---

## Documentos de Referencia

- `PROPUESTA_EMAILS.md` - Nuevos emails para sorteos
- `GUIA_EMAIL_SMTP_VS_SENDGRID.md` - Comparación detallada
- `.env.smtp.example` - Ejemplos de configuración SMTP
- `internal/adapters/notifier/smtp.go` - Implementación SMTP
- `internal/adapters/notifier/sendgrid.go` - Implementación SendGrid
