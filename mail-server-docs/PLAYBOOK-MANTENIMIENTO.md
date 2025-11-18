# PLAYBOOK DE MANTENIMIENTO - Servidor de Correo sorteos.club

## 📅 Tareas Diarias Automáticas

### Reportes
- ✅ Reporte diario de logs: `/etc/cron.daily/mail-report`
- ✅ Envío automático a: `postmaster@sorteos.club`
- ✅ Alerta si cola > 10 mensajes

---

## 🔍 Comandos Útiles de Monitoreo

### Verificar Estado de Servicios
```bash
systemctl status postfix dovecot opendkim fail2ban
```

### Ver Cola de Correo
```bash
mailq                    # Ver cola completa
postqueue -p             # Vista detallada
postsuper -d ALL         # Limpiar toda la cola (¡CUIDADO!)
postsuper -d <ID>        # Eliminar mensaje específico
```

### Ver Logs en Tiempo Real
```bash
tail -f /var/log/postfix.log           # Postfix
tail -f /var/log/mail.log              # General
journalctl -u postfix -f               # Postfix (systemd)
journalctl -u dovecot -f               # Dovecot
journalctl -u opendkim -f              # OpenDKIM
```

### Estadísticas de Correo
```bash
pflogsumm /var/log/postfix.log         # Resumen completo
pflogsumm -d today /var/log/postfix.log  # Solo hoy
```

### Verificar Autenticación
```bash
# Ver intentos fallidos
grep "authentication failed" /var/log/postfix.log | tail -20

# Ver IPs bloqueadas por Fail2ban
fail2ban-client status postfix
fail2ban-client status dovecot

# Desbloquear IP
fail2ban-client set postfix unbanip <IP>
```

---

## 🔐 Gestión de Usuarios de Correo

### Crear Nueva Cuenta
```bash
# 1. Crear usuario del sistema
useradd -m -s /bin/bash -c "Descripcion" nombre_usuario

# 2. Establecer contraseña
echo 'nombre_usuario:contraseña_segura' | chpasswd

# 3. Crear buzón Maildir
mkdir -p /home/nombre_usuario/Maildir/{cur,new,tmp}
chown -R nombre_usuario:nombre_usuario /home/nombre_usuario/Maildir
chmod -R 700 /home/nombre_usuario/Maildir
```

### Cambiar Contraseña
```bash
passwd nombre_usuario
```

### Eliminar Cuenta
```bash
userdel -r nombre_usuario  # Elimina usuario y su home
```

### Listar Cuentas Activas
```bash
grep "/home" /etc/passwd | grep -E "(noreply|info|soporte|postmaster)"
```

---

## 🛡️ Seguridad y Fail2ban

### Ver Estadísticas de Bloqueos
```bash
fail2ban-client status               # Jails activas
fail2ban-client status postfix       # Detalles Postfix
fail2ban-client status dovecot       # Detalles Dovecot
fail2ban-client status sshd          # Detalles SSH
```

### Desbloquear IP Específica
```bash
fail2ban-client set postfix unbanip 1.2.3.4
fail2ban-client set dovecot unbanip 1.2.3.4
```

### Añadir IP a Whitelist Permanente
Editar `/etc/fail2ban/jail.local` y añadir en `[DEFAULT]`:
```
ignoreip = 127.0.0.1/8 ::1 TU.IP.AQUI
```
Luego: `systemctl restart fail2ban`

---

## 📊 Verificación DNS y Deliverability

### Verificar Registros DNS
```bash
# SPF
dig TXT sorteos.club +short
dig TXT mail.sorteos.club +short

# DKIM
dig TXT default._domainkey.sorteos.club +short

# DMARC
dig TXT _dmarc.sorteos.club +short

# MX
dig MX sorteos.club +short

# PTR (Reverse DNS)
dig -x 62.171.188.255 +short
```

### Herramientas Online
- SPF: https://mxtoolbox.com/spf.aspx
- DKIM: https://mxtoolbox.com/dkim.aspx
- DMARC: https://mxtoolbox.com/dmarc.aspx
- Mail Tester: https://www.mail-tester.com/
- Blacklist Check: https://mxtoolbox.com/blacklists.aspx

---

## 🔄 Mantenimiento Regular

### Semanal
1. Revisar reportes de `postmaster@sorteos.club`
2. Verificar tamaño de logs: `du -sh /var/log/postfix.log /var/log/mail.log`
3. Verificar espacio en disco: `df -h`
4. Revisar IPs bloqueadas: `fail2ban-client status`

### Mensual
1. Rotar logs manualmente si es necesario: `logrotate -f /etc/logrotate.d/rsyslog`
2. Actualizar sistema: `apt update && apt upgrade`
3. Verificar certificados SSL: `certbot certificates`
4. Revisar reportes DMARC recibidos en `postmaster@sorteos.club`

### Trimestral
1. Backup de claves DKIM: `/opt/Sorteos/mail-server-docs/dkim-backup/`
2. Revisar y ajustar política DMARC (p=none → p=quarantine → p=reject)
3. Auditar cuentas de correo inactivas
4. Revisar rendimiento y optimizar si es necesario

### Anual
1. Considerar rotación de claves DKIM
2. Renovar certificados SSL (automático con certbot, solo verificar)
3. Auditoría de seguridad completa

---

## 🚨 Resolución de Problemas Comunes

### Problema: Correos no se envían

<details>
<summary>Solución</summary>

```bash
# 1. Verificar servicios activos
systemctl status postfix dovecot opendkim

# 2. Ver cola de correo
mailq

# 3. Ver logs de error
tail -50 /var/log/postfix.log | grep -i error

# 4. Verificar puertos abiertos
ss -tlnp | grep -E ":(25|465|587)"

# 5. Probar envío manual
echo "Test" | sendmail -v usuario@dominio.com
```
</details>

### Problema: Autenticación SMTP falla

<details>
<summary>Solución</summary>

```bash
# 1. Verificar Dovecot activo
systemctl status dovecot

# 2. Ver logs de autenticación
journalctl -u dovecot -n 50 | grep auth

# 3. Verificar socket de Postfix/Dovecot
ls -lah /var/spool/postfix/private/auth

# 4. Reiniciar servicios
systemctl restart dovecot postfix
```
</details>

### Problema: Correos van a spam

<details>
<summary>Solución</summary>

```bash
# 1. Verificar DKIM está firmando
grep "DKIM-Signature" /var/log/postfix.log

# 2. Verificar registros DNS
dig TXT default._domainkey.sorteos.club +short
dig TXT sorteos.club +short | grep spf
dig TXT _dmarc.sorteos.club +short

# 3. Verificar PTR
dig -x 62.171.188.255 +short

# 4. Probar con mail-tester.com
# Enviar correo a la dirección que te dan

# 5. Verificar blacklists
# Ir a https://mxtoolbox.com/blacklists.aspx
```
</details>

### Problema: Cola de correo saturada

<details>
<summary>Solución</summary>

```bash
# 1. Ver mensajes en cola
postqueue -p | less

# 2. Ver razón de retención
postcat -vq <QUEUE_ID>

# 3. Intentar reenviar todos
postqueue -f

# 4. Eliminar correos específicos
postsuper -d <QUEUE_ID>

# 5. Limpiar cola completa (EXTREMO)
postsuper -d ALL
postsuper -d ALL deferred
```
</details>

---

## 📁 Ubicaciones Importantes

### Archivos de Configuración
```
/etc/postfix/main.cf              - Config principal Postfix
/etc/postfix/master.cf            - Servicios Postfix
/etc/dovecot/dovecot.conf         - Config principal Dovecot
/etc/dovecot/conf.d/              - Configs modulares Dovecot
/etc/opendkim.conf                - Config OpenDKIM
/etc/fail2ban/jail.local          - Config Fail2ban
```

### Logs
```
/var/log/postfix.log              - Logs Postfix
/var/log/mail.log                 - Logs generales de correo
/var/log/fail2ban.log             - Logs Fail2ban
```

### Datos y Buzones
```
/home/*/Maildir/                  - Buzones de correo
/var/spool/postfix/               - Cola y spool de Postfix
/etc/opendkim/keys/               - Claves DKIM
```

### Documentación
```
/opt/Sorteos/mail-server-docs/mail-server-credentials.txt
/opt/Sorteos/mail-server-docs/dns-records-sorteos-club.txt
/opt/Sorteos/mail-server-docs/dkim-backup/
```

---

## 🔄 Backup y Restore

### Hacer Backup Completo
```bash
#!/bin/bash
BACKUP_DIR="/root/mail-backup-$(date +%Y%m%d)"
mkdir -p $BACKUP_DIR

# Configs
cp -r /etc/postfix/ $BACKUP_DIR/
cp -r /etc/dovecot/ $BACKUP_DIR/
cp -r /etc/opendkim/ $BACKUP_DIR/
cp /etc/fail2ban/jail.local $BACKUP_DIR/

# Buzones (puede ser grande)
tar -czf $BACKUP_DIR/mailboxes.tar.gz /home/*/Maildir/

# Crear archivo comprimido
cd /root && tar -czf mail-backup-$(date +%Y%m%d).tar.gz mail-backup-$(date +%Y%m%d)/

echo "Backup completo en: /root/mail-backup-$(date +%Y%m%d).tar.gz"
```

### Restaurar Configuración
```bash
# Descomprimir backup
tar -xzf mail-backup-YYYYMMDD.tar.gz

# Restaurar configs
cp -r mail-backup-YYYYMMDD/postfix/* /etc/postfix/
cp -r mail-backup-YYYYMMDD/dovecot/* /etc/dovecot/
cp -r mail-backup-YYYYMMDD/opendkim/* /etc/opendkim/

# Reiniciar servicios
systemctl restart postfix dovecot opendkim
```

---

## 📞 Contactos y Escalación

**Administrador Principal**:
- Email: postmaster@sorteos.club
- Documentación: `/opt/Sorteos/mail-server-docs/`

**Soporte Crítico**:
1. Verificar logs primero
2. Consultar este playbook
3. Revisar documentación oficial:
   - Postfix: http://www.postfix.org/documentation.html
   - Dovecot: https://doc.dovecot.org/
   - OpenDKIM: http://www.opendkim.org/

---

**Última actualización**: 2025-11-17
**Versión**: 1.0
**Servidor**: mail.sorteos.club (62.171.188.255)
