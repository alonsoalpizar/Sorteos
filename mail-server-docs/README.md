# 📧 Servidor de Correo - sorteos.club

## ✅ Instalación Completada - 2025-11-17

---

## 🎯 Estado del Servidor

**Servidor**: `mail.sorteos.club`
**IP**: `62.171.188.255`
**Sistema**: Ubuntu 24.04 LTS
**Hostname**: `mail.sorteos.club`

### Servicios Instalados y Activos

| Servicio | Puerto | Estado | Función |
|----------|--------|--------|---------|
| **Postfix** | 25 (SMTP) | ✅ Activo | Envío/Recepción correo |
| | 465 (SMTPS) | ✅ Activo | SMTP con TLS implícito |
| | 587 (Submission) | ✅ Activo | SMTP con STARTTLS |
| **Dovecot** | 993 (IMAPS) | ✅ Activo | Acceso IMAP seguro |
| | 995 (POP3S) | ✅ Activo | Acceso POP3 seguro |
| **OpenDKIM** | - | ✅ Activo | Firma digital (DKIM) |
| **Fail2ban** | - | ✅ Activo | Protección anti-fuerza bruta |
| **Nginx** | 443 | ✅ Activo | Webmail (SnappyMail) |

---

## 📁 Cuentas de Correo Creadas

### Cuentas Activas

1. **noreply@sorteos.club**
   - Usuario: `noreply`
   - Propósito: Correos transaccionales (activaciones, notificaciones)
   - Password: Ver `mail-server-credentials.txt`

2. **info@sorteos.club**
   - Usuario: `info`
   - Propósito: Contacto general
   - Password: Ver `mail-server-credentials.txt`

3. **soporte@sorteos.club**
   - Usuario: `soporte`
   - Propósito: Soporte técnico
   - Password: Ver `mail-server-credentials.txt`

4. **postmaster@sorteos.club**
   - Usuario: `postmaster`
   - Propósito: Administración, reportes DMARC
   - Password: Ver `mail-server-credentials.txt`

---

## 🌐 Acceso al Servidor

### SMTP (Para envío desde aplicaciones)
```
Host: mail.sorteos.club
Puerto: 587 (STARTTLS) - RECOMENDADO
Puerto: 465 (SSL/TLS) - Alternativo
Usuario: noreply@sorteos.club (o la cuenta que uses)
Password: [ver credentials.txt]
TLS/SSL: OBLIGATORIO
```

### IMAP (Para leer correos)
```
Host: mail.sorteos.club
Puerto: 993 (IMAPS)
Usuario: usuario@sorteos.club
Password: [ver credentials.txt]
TLS/SSL: OBLIGATORIO
```

### Webmail (Navegador)
```
URL: https://webmail.sorteos.club
Usuario: usuario@sorteos.club
Password: [ver credentials.txt]
```

**NOTA**: Webmail requiere que configures el registro DNS:
```
Tipo: A
Nombre: webmail.sorteos.club
Valor: 62.171.188.255
```

---

## 🔐 Seguridad Implementada

### ✅ Autenticación y Encriptación
- TLS/SSL en todos los puertos (465, 587, 993, 995)
- Certificados Let's Encrypt válidos
- Autenticación SASL obligatoria para envío
- Protocolos seguros: TLS 1.2, TLS 1.3

### ✅ Protección Anti-Spam
- SPF configurado
- DKIM firmando todos los correos salientes
- DMARC con política `p=none` (monitoreo inicial)
- Fail2ban activo contra ataques de fuerza bruta

### ✅ Registros DNS Configurados
- PTR (Reverse DNS): ✅ `62.171.188.255` → `mail.sorteos.club`
- MX: Pendiente de configuración DNS
- SPF: Pendiente de configuración DNS
- DKIM: Clave generada, pendiente DNS
- DMARC: Pendiente de configuración DNS

**Ver archivo** `dns-records-sorteos-club.txt` para registros completos.

---

## 🧪 Pruebas Realizadas

### ✅ Pruebas Internas (100% Exitosas)
- [x] Entrega local de correos (Maildir)
- [x] Firma DKIM activa y funcionando
- [x] Autenticación SMTP puerto 587 con TLS
- [x] Dovecot autenticando correctamente
- [x] Fail2ban bloqueando intentos fallidos

### ⏳ Pruebas Externas (Pendientes DNS)
- [ ] Envío a Gmail
- [ ] Envío a Outlook
- [ ] Verificación SPF/DKIM/DMARC en headers
- [ ] Test en mail-tester.com

---

## 📋 Próximos Pasos CRÍTICOS

### 1. Configurar Registros DNS (URGENTE)
**Archivo**: `dns-records-sorteos-club.txt`

Registros mínimos obligatorios:
```
✅ PTR: Ya configurado
⏳ A (mail.sorteos.club): 62.171.188.255
⏳ SPF (sorteos.club): v=spf1 +a +mx +ip4:62.171.188.255 ~all
⏳ SPF (mail.sorteos.club): v=spf1 a mx ip4:62.171.188.255 ~all
⏳ DKIM (default._domainkey.sorteos.club): [ver archivo DNS]
⏳ DMARC (_dmarc.sorteos.club): v=DMARC1; p=none; rua=mailto:postmaster@sorteos.club
```

### 2. Hacer Pruebas Externas
Cuando DNS propague (1-2 horas):
```bash
# Desde el servidor, enviar a tu Gmail personal
swaks --to tu@gmail.com \
  --from noreply@sorteos.club \
  --server localhost \
  --port 587 \
  --auth PLAIN \
  --auth-user noreply@sorteos.club \
  --auth-password "[password]" \
  --tls \
  --header "Subject: Test Servidor Nuevo" \
  --body "Verificar headers SPF/DKIM/DMARC"
```

Luego:
1. Revisar el correo en Gmail
2. "Mostrar original" → Verificar headers
3. Debe mostrar: SPF=pass, DKIM=pass, DMARC=pass

### 3. Cambiar MX (Solo cuando todo funcione)
```
Tipo: MX
Nombre: sorteos.club
Prioridad: 10
Valor: mail.sorteos.club
```

### 4. Actualizar Aplicación sorteos.club
Editar: `/opt/Sorteos/backend/.env`
```env
CONFIG_EMAIL_PROVIDER=smtp
CONFIG_SMTP_HOST=mail.sorteos.club
CONFIG_SMTP_PORT=587
CONFIG_SMTP_USERNAME=noreply@sorteos.club
CONFIG_SMTP_PASSWORD=[password de noreply]
CONFIG_SMTP_USE_TLS=true
CONFIG_SMTP_FROM_EMAIL=noreply@sorteos.club
```

Reiniciar:
```bash
systemctl restart sorteos-api
```

---

## 📊 Monitoreo y Logs

### Comandos Útiles
```bash
# Ver estado de servicios
systemctl status postfix dovecot opendkim

# Ver cola de correo
mailq

# Logs en tiempo real
tail -f /var/log/postfix.log

# Estadísticas diarias
pflogsumm -d today /var/log/postfix.log

# IPs bloqueadas por Fail2ban
fail2ban-client status postfix
```

### Reportes Automáticos
- **Frecuencia**: Diario
- **Destino**: postmaster@sorteos.club
- **Contenido**: Resumen de logs, alertas de cola

---

## 📚 Documentación Disponible

| Archivo | Descripción |
|---------|-------------|
| `mail-server-credentials.txt` | Todas las contraseñas de cuentas |
| `dns-records-sorteos-club.txt` | Registros DNS completos para copiar |
| `dkim-backup/` | Backup de claves DKIM |
| `PLAYBOOK-MANTENIMIENTO.md` | Guía completa de administración |
| `README.md` | Este archivo |

---

## 🆘 Soporte y Ayuda

### Problemas Comunes
Consulta: `PLAYBOOK-MANTENIMIENTO.md`

### Logs Importantes
```
/var/log/postfix.log    - Logs de Postfix
/var/log/mail.log       - Logs generales
/var/log/fail2ban.log   - Intentos bloqueados
```

### Verificar DNS
```bash
dig TXT sorteos.club +short
dig TXT default._domainkey.sorteos.club +short
dig MX sorteos.club +short
dig -x 62.171.188.255 +short
```

---

## 🎉 Resumen de lo Logrado

✅ Servidor de correo completo instalado y funcionando
✅ 4 cuentas de correo creadas y operativas
✅ DKIM firmando correos automáticamente
✅ TLS/SSL en todos los servicios
✅ Fail2ban protegiendo contra ataques
✅ Webmail (SnappyMail) instalado
✅ Monitoreo automático configurado
✅ Documentación completa generada
✅ Backups de claves DKIM realizados

---

## ⚠️ IMPORTANTE

**Antes de usar en producción**:
1. Configurar TODOS los registros DNS
2. Esperar propagación (1-24 horas)
3. Hacer pruebas de envío a Gmail/Outlook
4. Verificar score en mail-tester.com (debe ser >8/10)
5. Solo entonces cambiar el MX
6. Actualizar .env de la aplicación

**Progresión DMARC Recomendada**:
- Semana 1-2: `p=none` (solo monitorear)
- Semana 3-4: `p=quarantine` (enviar a spam los que fallan)
- Semana 5+: `p=reject` (rechazar los que fallan)

---

**Instalación realizada**: 2025-11-17
**Servidor**: mail.sorteos.club
**Estado**: ✅ Listo para configuración DNS
