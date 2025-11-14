# Guía Completa: SMTP Propio vs SendGrid

## Comparación Rápida

| Aspecto | SMTP Propio | SendGrid |
|---------|-------------|----------|
| **Costo** | $0 (si tienes servidor) | $0 - $20/mes (100-50K emails) |
| **Deliverability** | 70-85% (requiere config) | 95-99% (optimizado) |
| **Setup** | Complejo (DNS, servidor) | Fácil (solo API key) |
| **Analytics** | No (debes implementar) | Sí (dashboard completo) |
| **Mantenimiento** | Alto (tú lo gestionas) | Bajo (gestionado) |
| **Escalabilidad** | Limitada por servidor | Ilimitada |
| **Velocidad** | Depende del servidor | Muy rápida (infraestructura global) |
| **Soporte** | Tú mismo | Soporte técnico 24/7 |

---

## Opción 1: SMTP Propio (Tu Servidor MX)

### ✅ Ventajas

1. **Costo Cero** - Si ya tienes un servidor con email
2. **Control Total** - Sobre tu infraestructura
3. **Privacidad** - Los datos no pasan por terceros
4. **Sin Límites Externos** - Solo limitado por tu servidor
5. **Personalización Completa** - Headers, configuraciones, etc.

### ❌ Desventajas

1. **Deliverability Menor** (70-85% vs 95-99%)
   - Mayor riesgo de caer en spam
   - IPs nuevas sin reputación
   - Requiere configurar SPF, DKIM, DMARC correctamente

2. **Configuración Compleja**
   ```bash
   # Debes configurar:
   - Servidor SMTP (Postfix, Exim, etc.)
   - Certificados SSL/TLS
   - Registros DNS (A, MX, PTR, SPF, DKIM, DMARC)
   - Firewall (puertos 25, 465, 587)
   - Antispam/Antivirus
   - Monitoreo de colas
   ```

3. **Mantenimiento Continuo**
   - Actualizar servidor
   - Monitorear colas de email
   - Gestionar bounces
   - Limpiar listas de blacklist
   - Revisar logs de errores

4. **Sin Analytics**
   - No sabes si se abrió el email
   - No sabes si hicieron click
   - Debes implementar tracking manualmente

5. **Problemas de Reputación**
   - Si tu IP está en blacklist, ningún email llega
   - Toma meses construir buena reputación
   - Un solo email de spam puede arruinarlo

### Cuándo Usar SMTP Propio

✅ **Úsalo si:**
- Ya tienes un servidor de correo configurado y funcionando
- Envías menos de 100 emails/día
- No te importan las métricas de apertura/clicks
- Tienes conocimientos técnicos de email servers
- Tu IP/dominio ya tiene buena reputación

❌ **NO lo uses si:**
- Estás empezando desde cero
- Necesitas alta deliverability
- Quieres métricas de emails
- No tienes tiempo para gestionar un servidor SMTP

---

## Opción 2: SendGrid (Servicio Cloud)

### ✅ Ventajas

1. **Setup Ultra Rápido** (5 minutos)
   ```bash
   1. Registrarse en sendgrid.com
   2. Crear API Key
   3. Agregar a .env
   4. Listo!
   ```

2. **Deliverability Excelente** (95-99%)
   - IPs con reputación establecida
   - SPF/DKIM/DMARC pre-configurado
   - Monitoreo automático de bounces
   - Relaciones con Gmail, Outlook, etc.

3. **Analytics Completo**
   - Dashboard visual
   - Open rate (tasa de apertura)
   - Click rate (tasa de clicks)
   - Bounce tracking
   - Spam reports
   - Device/Browser stats
   - Geolocalización

4. **Infraestructura Global**
   - Servidores en múltiples regiones
   - Envío paralelo
   - Failover automático
   - 99.95% uptime SLA

5. **Features Avanzados**
   - A/B Testing de emails
   - Templates visuales (drag & drop)
   - Webhooks para eventos
   - API REST completa
   - Integración con marketing automation

### ❌ Desventajas

1. **Costo** (después del plan gratuito)
   ```
   Free:      100 emails/día (limitado)
   Essentials: $19.95/mes - 50,000 emails
   Pro:        $89.95/mes - 100,000 emails
   Premier:   $449/mes - 1,200,000 emails
   ```

2. **Dependencia de Terceros**
   - Si SendGrid cae, no envías emails
   - Cambios de precios
   - Cambios en ToS

3. **Datos Externos**
   - Los emails pasan por servidores de SendGrid
   - Metadata compartida con SendGrid

### Cuándo Usar SendGrid

✅ **Úsalo si:**
- Estás empezando un proyecto nuevo
- Necesitas alta deliverability
- Quieres analytics y métricas
- Prefieres enfocarte en tu producto
- Envías más de 100 emails/día
- Necesitas escalabilidad

❌ **NO lo uses si:**
- Tienes restricciones de privacidad extremas
- No quieres pagar después de 100 emails/día
- Ya tienes infraestructura SMTP funcionando perfectamente

---

## Alternativas Populares

### AWS SES (Amazon Simple Email Service)

**Pros:**
- **Muy barato:** $0.10 por 1000 emails
- Infraestructura de Amazon
- Alta deliverability

**Contras:**
- Configuración más técnica
- Analytics básico
- Requiere cuenta AWS

**Recomendado para:** Proyectos con mucho volumen (50K+ emails/mes)

---

### Mailgun

**Pros:**
- Similar a SendGrid
- Buen API
- Analytics completo

**Contras:**
- Más caro que SendGrid
- UI menos amigable

**Precio:** $35/mes - 50,000 emails

---

### Postmark

**Pros:**
- Excelente para emails transaccionales
- Deliverability del 98%
- Soporte premium

**Contras:**
- Más caro
- No tiene marketing automation

**Precio:** $15/mes - 10,000 emails

---

### Mailtrap (Solo Testing)

**Pros:**
- Gratis
- Perfecto para desarrollo
- Captura emails sin enviarlos

**Contras:**
- NO envía emails reales
- Solo para testing

---

## Configuración de Tu Propio SMTP

Si decides usar tu propio servidor, aquí está la guía completa:

### 1. Servidor SMTP

Instalar Postfix (Ubuntu/Debian):

```bash
sudo apt update
sudo apt install postfix mailutils

# Configurar como "Internet Site"
# Hostname: mail.sorteos.club
```

### 2. Certificado SSL/TLS

```bash
# Usar Let's Encrypt
sudo apt install certbot
sudo certbot certonly --standalone -d mail.sorteos.club

# Configurar en Postfix
sudo nano /etc/postfix/main.cf
```

Agregar:
```
smtpd_tls_cert_file=/etc/letsencrypt/live/mail.sorteos.club/fullchain.pem
smtpd_tls_key_file=/etc/letsencrypt/live/mail.sorteos.club/privkey.pem
smtpd_use_tls=yes
smtpd_tls_security_level=may
```

### 3. Autenticación SASL

```bash
sudo apt install libsasl2-modules sasl2-bin

# Crear usuarios
sudo saslpasswd2 -c noreply@sorteos.club
```

### 4. Configurar DNS

**Registro A:**
```
mail.sorteos.club. IN A 62.171.188.255
```

**Registro MX:**
```
sorteos.club. IN MX 10 mail.sorteos.club.
```

**Registro PTR (Reverse DNS):**
```
255.188.171.62.in-addr.arpa. IN PTR mail.sorteos.club.
```
*Nota: Esto se configura con tu proveedor de hosting/VPS*

**Registro SPF:**
```
sorteos.club. IN TXT "v=spf1 mx a ip4:62.171.188.255 ~all"
```

**Registro DKIM:**
```bash
# Generar claves DKIM
sudo apt install opendkim opendkim-tools
sudo opendkim-genkey -d sorteos.club -s default

# Ver clave pública
sudo cat /etc/opendkim/keys/sorteos.club/default.txt
```

Agregar a DNS:
```
default._domainkey.sorteos.club. IN TXT "v=DKIM1; k=rsa; p=MIGfMA0GCSqGSIb3DQEBAQUAA4..."
```

**Registro DMARC:**
```
_dmarc.sorteos.club. IN TXT "v=DMARC1; p=quarantine; rua=mailto:dmarc@sorteos.club; pct=100"
```

### 5. Configurar Firewall

```bash
sudo ufw allow 25/tcp   # SMTP
sudo ufw allow 465/tcp  # SMTPS (TLS)
sudo ufw allow 587/tcp  # Submission (STARTTLS)
```

### 6. Probar Configuración

```bash
# Test básico
echo "Test email" | mail -s "Test" destinatario@example.com

# Test SMTP auth
telnet mail.sorteos.club 587
```

### 7. Monitoreo

```bash
# Ver cola de emails
mailq

# Ver logs
sudo tail -f /var/log/mail.log

# Estadísticas
sudo pflogsumm /var/log/mail.log
```

### 8. Verificar Spam Score

Envía un email a: check-auth@verifier.port25.com

Recibirás un reporte con:
- SPF: PASS ✅
- DKIM: PASS ✅
- DMARC: PASS ✅
- Spam Score

---

## Implementación en Tu Proyecto

### Configurar .env para SMTP

```bash
# Cambiar proveedor a SMTP
CONFIG_EMAIL_PROVIDER=smtp

# Frontend URL
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

### Modificar routes.go

El código ya está listo, solo necesitas cambiar el notifier:

```go
// En cmd/api/routes.go

var emailNotifier Notifier

if cfg.EmailProvider == "smtp" {
    emailNotifier = notifier.NewSMTPNotifier(&cfg.SMTP, log)
} else {
    emailNotifier = notifier.NewSendGridNotifier(&cfg.SendGrid, log)
}

// Usar emailNotifier en vez de sendgridNotifier
registerUseCase := auth.NewRegisterUseCase(
    userRepo, consentRepo, auditRepo, tokenMgr,
    emailNotifier, // <-- Aquí
    log, cfg.SkipEmailVerification,
)
```

### Probar

```bash
# Compilar
cd /opt/Sorteos/backend
go build -o sorteos-api cmd/api/main.go

# Ejecutar
./sorteos-api

# Registrar un usuario
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "Password123!@#",
    "accepted_terms": true,
    "accepted_privacy": true
  }'

# Revisar logs
tail -f /var/log/mail.log  # (si usas tu SMTP)
```

---

## Recomendación Final

### Para Desarrollo/MVP:
**SendGrid Plan Gratuito** (100 emails/día)
- Setup en 5 minutos
- Sin complejidad
- Perfecto para empezar

### Para Producción (< 1000 emails/día):
**SendGrid Essentials** ($19.95/mes - 50K emails)
- Excelente deliverability
- Analytics completo
- Sin mantenimiento

### Para Producción (> 50K emails/mes):
**AWS SES** ($0.10 por 1000 emails)
- Muy económico
- Escalable
- Requiere algo más de configuración

### Solo si Ya Tienes Infraestructura:
**Tu Propio SMTP**
- Si ya tienes servidor de correo funcionando
- Si tienes experiencia con email servers
- Si el costo es crítico

---

## Migración de SendGrid a SMTP

Si decides cambiar después, es fácil:

1. Configurar tu servidor SMTP
2. Cambiar en `.env`:
   ```bash
   CONFIG_EMAIL_PROVIDER=smtp  # Era "sendgrid"
   ```
3. Reiniciar backend
4. Listo!

El código sigue funcionando igual porque ambos implementan la misma interface `Notifier`.

---

## FAQs

### ¿Puedo usar Gmail como SMTP?
Sí, pero con limitaciones:
- Máximo 500 emails/día
- Solo para desarrollo/testing
- Debes habilitar "App Passwords"

### ¿Necesito un servidor dedicado para SMTP?
No necesariamente. Puedes usar:
- VPS compartido (si permite puerto 25)
- Servidor dedicado
- Servicio managed como AWS SES

### ¿Qué pasa si mi IP está en blacklist?
Debes solicitar delisting en:
- Spamhaus.org
- Barracuda Central
- SpamCop
- MXToolbox Blacklist Check

Puede tomar días o semanas.

### ¿Cuál es más seguro?
Ambos son seguros si se configuran correctamente:
- SMTP: TLS/STARTTLS + contraseñas fuertes
- SendGrid: API Key + HTTPS

### ¿Puedo combinar ambos?
Sí! Puedes:
- SMTP para emails internos
- SendGrid para emails a clientes

Solo necesitas crear dos instancias de notifier.

---

## Recursos Adicionales

- [SendGrid Docs](https://docs.sendgrid.com/)
- [Postfix Tutorial](https://www.postfix.org/documentation.html)
- [Mail Tester](https://www.mail-tester.com/) - Test spam score
- [MXToolbox](https://mxtoolbox.com/) - DNS/Email diagnostics
- [DMARC Analyzer](https://dmarc.org/)

---

## Conclusión

**Para tu proyecto Sorteos Platform:**

Si estás empezando → **SendGrid**
- Fácil, rápido, confiable
- Plan gratuito suficiente para MVP
- Puedes cambiar después

Si ya tienes servidor SMTP → **SMTP Propio**
- Aprovecha lo que tienes
- Costo cero
- Requiere conocimiento técnico

**Mi recomendación personal:**
Empieza con SendGrid, y si después creces mucho (>50K emails/mes), migra a AWS SES para reducir costos.
