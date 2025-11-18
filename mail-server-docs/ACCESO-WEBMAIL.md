# 🌐 Acceso al Webmail - SnappyMail

## ✅ Configuración Completada

El webmail está **completamente configurado** y listo para usar.

---

## 🔗 URL de Acceso

**Webmail**: https://webmail.sorteos.club

---

## 👤 Credenciales de Acceso

### Cuentas Disponibles

| Email | Contraseña | Propósito |
|-------|-----------|-----------|
| noreply@sorteos.club | 9NhNlT4m6FqUbM28FSFuSg== | Emails automáticos (no leer) |
| info@sorteos.club | +yZ4o7A07toh/4MotrCqTw== | Contacto general |
| soporte@sorteos.club | FQh7jA1Cuth1SP/+oBhopg== | Soporte técnico |
| postmaster@sorteos.club | YKiTy53jeer2LC/UZNripQ== | Administración |

---

## 🔧 Configuración Aplicada

SnappyMail está configurado para usar:

### IMAP (Lectura de correos)
```
Host: localhost
Puerto: 993
Tipo: SSL/TLS (type=2)
Autenticación: PLAIN/LOGIN
```

### SMTP (Envío de correos)
```
Host: localhost
Puerto: 587
Tipo: STARTTLS (type=1)
Autenticación: Obligatoria
```

---

## 📝 Cómo Iniciar Sesión

1. **Ir a**: https://webmail.sorteos.club

2. **Ingresar credenciales**:
   - Usuario: `info@sorteos.club` (o cualquier cuenta)
   - Contraseña: (ver tabla arriba)

3. **¡Listo!** Ya puedes ver y enviar correos desde el navegador

---

## ⚙️ Configuración Técnica

### Archivos de Configuración

- **Dominio sorteos.club**: `/var/www/webmail/data/_data_/_default_/domains/sorteos.club.json`
- **Dominio mail.sorteos.club**: `/var/www/webmail/data/_data_/_default_/domains/mail.sorteos.club.json`

### Cambios Realizados

✅ Puerto IMAP: 143 → 993 (SSL/TLS)
✅ Puerto SMTP: 25 → 587 (STARTTLS)
✅ Tipo conexión IMAP: 0 → 2 (SSL/TLS)
✅ Tipo conexión SMTP: 0 → 1 (STARTTLS)
✅ SMTP Auth: false → true
✅ SSL allow_self_signed: true

---

## 🔐 Certificado SSL

El webmail usa certificado Let's Encrypt válido que cubre:
- sorteos.club
- www.sorteos.club
- mail.sorteos.club
- webmail.sorteos.club

**Válido hasta**: 2026-02-15

---

## 🎯 Funcionalidades

- ✅ Lectura de correos (IMAP)
- ✅ Envío de correos (SMTP)
- ✅ Gestión de carpetas
- ✅ Búsqueda de mensajes
- ✅ Adjuntos
- ✅ Firma DKIM automática en envíos
- ✅ Acceso seguro HTTPS

---

## 📱 Acceso desde Clientes de Correo

Si prefieres usar Thunderbird, Outlook o móvil en lugar del webmail:

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
```

---

## 🆘 Solución de Problemas

### Error: "Can't connect to host"
- Verificar que Dovecot esté activo: `systemctl status dovecot`
- Verificar puerto 993 abierto: `ss -tlnp | grep :993`

### Error: "Authentication failed"
- Verificar que usas el email completo (ej: info@sorteos.club)
- Verificar contraseña correcta
- Ver logs: `journalctl -u dovecot -f`

### No puedo enviar correos
- Verificar que Postfix esté activo: `systemctl status postfix`
- Ver logs: `tail -f /var/log/postfix.log`

---

**Fecha de configuración**: 2025-11-17
**Estado**: ✅ Operativo
