# 🌐 Acceso al Webmail - SnappyMail - CONFIGURACIÓN FINAL

## ✅ Estado: FUNCIONANDO

**Fecha**: 2025-11-17  
**Servidor**: mail.sorteos.club (62.171.188.255)

---

## 🔗 URLs de Acceso

### Webmail (Usuarios)
**URL**: https://webmail.sorteos.club

### Panel de Administración
**URL**: https://webmail.sorteos.club/?admin

---

## 👤 Credenciales

### Cuentas de Correo

| Email | Contraseña | Propósito |
|-------|------------|-----------|
| noreply@sorteos.club | 9NhNlT4m6FqUbM28FSFuSg== | Correos automáticos |
| info@sorteos.club | +yZ4o7A07toh/4MotrCqTw== | Contacto general |
| soporte@sorteos.club | FQh7jA1Cuth1SP/+oBhopg== | Soporte técnico |
| postmaster@sorteos.club | YKiTy53jeer2LC/UZNripQ== | Administración |

### Panel de Administración

```
Usuario: admin
Contraseña: JmCXgrdua+JA
TOTP: (dejar vacío)
```

---

## ⚙️ Configuración Técnica Aplicada

### IMAP (Lectura)
```
Host: 127.0.0.1 (localhost)
Puerto: 143
Tipo: 0 (Sin TLS - seguro para localhost)
Autenticación: PLAIN/LOGIN
```

### SMTP (Envío)
```
Host: 127.0.0.1 (localhost)
Puerto: 587
Tipo: 1 (STARTTLS)
Autenticación: Obligatoria
```

### Carpetas de Correo
Auto-creadas y suscritas automáticamente:
- ✅ Drafts (Borradores)
- ✅ Sent (Enviados)
- ✅ Trash (Papelera)
- ✅ Junk (Spam)
- ✅ INBOX (Bandeja de entrada)

---

## 🎯 Funcionalidades Confirmadas

- ✅ Login de usuarios funcionando
- ✅ Lectura de correos (IMAP)
- ✅ Envío de correos (SMTP con auth)
- ✅ Carpetas estándar auto-creadas
- ✅ Firma DKIM automática en envíos
- ✅ Certificado SSL válido (Let's Encrypt)
- ✅ Panel de administración accesible

---

## 📱 Configuración para Clientes de Correo

Si prefieres usar Thunderbird, Outlook, o app móvil:

### IMAP (Recibir)
```
Servidor: mail.sorteos.club
Puerto: 993
Seguridad: SSL/TLS
Usuario: tu-email@sorteos.club
Contraseña: [ver tabla arriba]
```

### SMTP (Enviar)
```
Servidor: mail.sorteos.club
Puerto: 587
Seguridad: STARTTLS
Usuario: tu-email@sorteos.club
Contraseña: [ver tabla arriba]
Autenticación: PLAIN
```

---

## 🔧 Archivos de Configuración

### Dominios
- `/var/www/webmail/data/_data_/_default_/domains/sorteos.club.json`
- `/var/www/webmail/data/_data_/_default_/domains/mail.sorteos.club.json`
- `/var/www/webmail/data/_data_/_default_/domains/default.json`

### Configuración General
- `/var/www/webmail/data/_data_/_default_/configs/application.ini`

### Contraseña de Admin
- `/var/www/webmail/data/_data_/_default_/admin_password.txt`

---

## 🆘 Solución de Problemas

### No puedo enviar correos
1. Verifica que Postfix esté activo: `systemctl status postfix`
2. Ver logs: `tail -f /var/log/postfix.log`
3. Verificar autenticación SMTP

### No aparecen las carpetas
1. Desconectar y volver a conectar al webmail
2. Las carpetas se crean automáticamente al primer acceso
3. Verificar: `systemctl status dovecot`

### Admin panel no acepta contraseña
1. Verificar contraseña actual: `cat /var/www/webmail/data/_data_/_default_/admin_password.txt`
2. La contraseña se regenera automáticamente si se borra

---

## ✨ Integración con sorteos.club

El backend ya está configurado para usar este servidor SMTP.

**Archivo**: `/opt/Sorteos/backend/.env`

```env
CONFIG_EMAIL_PROVIDER=smtp
CONFIG_SMTP_HOST=mail.sorteos.club
CONFIG_SMTP_PORT=587
CONFIG_SMTP_USERNAME=noreply@sorteos.club
CONFIG_SMTP_PASSWORD=9NhNlT4m6FqUbM28FSFuSg==
CONFIG_SMTP_USE_TLS=true
CONFIG_SMTP_USE_STARTTLS=true
```

**Estado del servicio**: ✅ Activo

---

**Instalación completada**: 2025-11-17  
**Webmail funcionando**: ✅ SÍ  
**Listo para producción**: ✅ SÍ

