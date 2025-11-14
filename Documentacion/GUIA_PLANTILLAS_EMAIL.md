# Guía de Plantillas de Email - Sorteos Platform

## 📧 ¿Dónde Crear las Plantillas?

Te he preparado **3 opciones** para gestionar tus plantillas de email. Elige la que mejor se adapte a tu equipo:

---

## ✨ Opción 1: Archivos HTML (RECOMENDADA) ⭐

### **Ubicación**
```
/opt/Sorteos/backend/templates/email/
├── verification.html           # ✅ Creada
├── welcome.html               # ✅ Creada
├── password_reset.html        # ✅ Creada
├── purchase_confirmation.html # ✅ Creada
└── README.md                  # ✅ Documentación completa
```

### **Ventajas**
- ✅ **Edición directa** - Solo editas el HTML
- ✅ **Sin recompilar** - Los cambios se aplican inmediatamente
- ✅ **Versionable** - Git rastrea cambios
- ✅ **Fácil para diseñadores** - No necesitan saber Go
- ✅ **Preview rápido** - Abre en navegador para ver

### **Cómo Usar**

**1. Editar plantilla:**
```bash
nano /opt/Sorteos/backend/templates/email/verification.html
```

**2. Guardar y listo:**
El próximo email usará la nueva versión automáticamente.

### **Variables Disponibles**

En tus plantillas HTML puedes usar:

```html
<!-- Email de Verificación -->
{{.FirstName}}        <!-- Nombre del usuario -->
{{.Code}}             <!-- Código de 6 dígitos -->
{{.FrontendURL}}      <!-- https://sorteos.club -->

<!-- Email de Bienvenida -->
{{.FirstName}}
{{.FrontendURL}}

<!-- Reset de Contraseña -->
{{.ResetURL}}         <!-- Link completo con token -->
{{.FrontendURL}}

<!-- Confirmación de Compra -->
{{.FirstName}}
{{.RaffleTitle}}      <!-- Nombre del sorteo -->
{{.RaffleID}}         <!-- UUID del sorteo -->
{{.Numbers}}          <!-- ["0001", "0042"] -->
{{.TotalAmount}}      <!-- "$150.00" -->
{{.DrawDate}}         <!-- "25/12/2025 20:00" -->
{{.Prize}}            <!-- "iPhone 15 Pro Max" -->
{{.FrontendURL}}
```

### **Ejemplo de Plantilla**

```html
<!DOCTYPE html>
<html>
<body style="font-family: Arial, sans-serif;">
    <div style="max-width: 600px; margin: 0 auto;">
        <h2 style="color: #3B82F6;">¡Hola {{.FirstName}}!</h2>
        <p>Tu código de verificación es: <strong>{{.Code}}</strong></p>
        <a href="{{.FrontendURL}}" style="color: #3B82F6;">
            Ir a Sorteos Platform
        </a>
    </div>
</body>
</html>
```

---

## 🔧 Opción 2: Plantillas Embebidas (Producción)

Para **ambientes de producción** donde prefieres un solo binario sin archivos externos.

### **Ubicación**
```
/opt/Sorteos/backend/internal/adapters/notifier/templates/
├── verification.html      # ✅ Copiada
├── welcome.html          # ✅ Copiada
└── password_reset.html   # ✅ Copiada
```

### **Cómo Funciona**

Las plantillas se **embeben en el binario Go** durante la compilación:

```go
//go:embed templates/*.html
var embeddedTemplates embed.FS
```

### **Ventajas**
- ✅ **Un solo archivo** - Todo en el binario
- ✅ **Portable** - No depende de archivos externos
- ✅ **Más rápido** - No lee del disco
- ✅ **Seguro** - No se pueden modificar en runtime

### **Desventajas**
- ❌ **Requiere recompilar** - Cada cambio necesita rebuild
- ❌ **Menos flexible** - No puedes cambiar en caliente

### **Cuándo Usar**
- Deploy a producción
- Ambientes containerizados (Docker)
- Cuando la portabilidad es crítica

### **Proceso de Actualización**

```bash
# 1. Editar plantilla
nano templates/email/verification.html

# 2. Copiar a directorio de embebido
cp templates/email/*.html internal/adapters/notifier/templates/

# 3. Recompilar
go build -o bin/sorteos-api ./cmd/api/

# 4. Reiniciar
sudo systemctl restart sorteos-api
```

---

## ☁️ Opción 3: SendGrid Dynamic Templates (No Recomendado)

**Nota:** Requiere cuenta de SendGrid de pago ($19.95/mes mínimo).

### **Ventajas**
- ✅ **Editor visual** - Drag & drop, sin código
- ✅ **A/B Testing** - Prueba variantes
- ✅ **Analytics** - Open rate, click rate
- ✅ **Sin mantenimiento** - SendGrid lo gestiona

### **Desventajas**
- ❌ **Costo mensual** - $19.95/mes mínimo
- ❌ **Vendor lock-in** - Dependes de SendGrid
- ❌ **Menos control** - No tienes el HTML

### **No lo recomiendo porque:**
Ya tienes SMTP propio sin costo. SendGrid solo agrega gastos innecesarios.

---

## 🎨 Diseño de las Plantillas Creadas

### **Características Profesionales**

1. **Responsive** - Se adapta a móvil y desktop
2. **Compatible** - Funciona en Gmail, Outlook, Apple Mail
3. **Branded** - Colores de tu marca (#3B82F6)
4. **Accesible** - Alto contraste, texto legible
5. **HTML + Texto Plano** - Fallback para clientes antiguos

### **Colores del Sistema**

```css
Azul Principal:    #3B82F6  (Botones, headers)
Verde Éxito:       #10B981  (Confirmaciones)
Rojo Alerta:       #EF4444  (Urgente, reset)
Amarillo Info:     #F59E0B  (Premios, tips)
Gris Texto:        #333333  (Principal)
Gris Secundario:   #64748B  (Secundario)
```

### **Estructura Común**

Todas las plantillas tienen:
- **Header colorido** con ícono y título
- **Body** con contenido principal
- **Cajas destacadas** para información importante
- **Botones CTA** (Call To Action)
- **Footer** con links y copyright

---

## 📋 Plantillas Disponibles

### 1. **verification.html** - Verificación de Cuenta
- **Cuándo:** Usuario se registra
- **Contiene:** Código de 6 dígitos grande
- **Color:** Azul (#3B82F6)
- **CTA:** Opcional link de verificación

### 2. **welcome.html** - Bienvenida
- **Cuándo:** Email verificado exitosamente
- **Contiene:** Checklist de cuenta activada
- **Color:** Verde (#10B981)
- **CTA:** Explorar Sorteos

### 3. **password_reset.html** - Reset de Contraseña
- **Cuándo:** Usuario olvidó contraseña
- **Contiene:** Link con token de reset
- **Color:** Rojo (#EF4444)
- **CTA:** Botón de restablecer

### 4. **purchase_confirmation.html** - Confirmación de Compra
- **Cuándo:** Pago exitoso de números
- **Contiene:** Detalles de compra, números, fecha sorteo
- **Color:** Azul (#3B82F6)
- **CTA:** Ver Sorteo en Vivo

---

## 🛠️ Cómo Modificar una Plantilla

### **Ejemplo: Cambiar Color del Header**

```bash
# 1. Abrir plantilla
nano /opt/Sorteos/backend/templates/email/verification.html

# 2. Buscar (Ctrl+W):
background-color: #3B82F6

# 3. Cambiar por tu color:
background-color: #FF6B6B

# 4. Guardar (Ctrl+O, Enter, Ctrl+X)

# 5. ¡Listo! Próximo email usará el nuevo color
```

### **Ejemplo: Agregar Tu Logo**

```html
<!-- En el header, antes del h1 -->
<tr>
    <td style="text-align: center; padding: 20px 0;">
        <img src="https://sorteos.club/logo.png"
             alt="Sorteos Platform"
             width="150"
             style="max-width: 150px;">
    </td>
</tr>
```

**Nota:** El logo debe estar hosteado online (no puede ser archivo local).

---

## 🧪 Probar Plantillas

### **Método 1: Enviar Email Real**
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"test@gmail.com","password":"Test123@","accepted_terms":true,"accepted_privacy":true}'
```

### **Método 2: Ver HTML en Navegador**
```bash
# Copiar plantilla al directorio público del frontend
cp templates/email/verification.html ../frontend/public/test-email.html

# Abrir en navegador
# http://localhost:5173/test-email.html
```

### **Método 3: Herramientas Online**
- [Litmus](https://litmus.com/) - Testing profesional
- [Mailtrap](https://mailtrap.io/) - Sandbox de desarrollo (gratis)
- [Email on Acid](https://www.emailonacid.com/) - Preview

---

## 📊 Mejores Prácticas

### **✅ DO (Hacer)**
1. Usar **tablas** para layout (compatibilidad con Outlook)
2. **Estilos inline** (no CSS externo)
3. **Width máximo 600px** (estándar de emails)
4. **Incluir texto plano** además de HTML
5. **Alt text en imágenes**
6. **Probar en múltiples clientes** (Gmail, Outlook, Apple Mail)

### **❌ DON'T (No Hacer)**
1. ~~CSS en `<style>` tags~~ (Outlook lo ignora)
2. ~~JavaScript~~ (Bloqueado por seguridad)
3. ~~Video embebido~~ (No funciona)
4. ~~Flexbox o Grid~~ (Soporte limitado)
5. ~~Imágenes de fondo~~ (No en Outlook)
6. ~~Fuentes web complejas~~ (Stick to Arial, Verdana, Georgia)

---

## 🚀 Agregar Nueva Plantilla

### **Paso 1: Crear HTML**
```bash
cd /opt/Sorteos/backend/templates/email
nano raffle_reminder.html
```

### **Paso 2: Usar Plantilla Base**
Copia `verification.html` como base y modifica:
- Header color y título
- Variables `{{.Nombre}}`
- Contenido del body
- Botón CTA

### **Paso 3: Definir Datos en Go**
```go
// En internal/adapters/notifier/template_loader.go

type RaffleReminderData struct {
    FirstName   string
    RaffleTitle string
    DrawDate    string
    Numbers     []string
    FrontendURL string
}
```

### **Paso 4: Crear Método de Envío**
```go
// En internal/adapters/notifier/smtp.go

func (n *SMTPNotifier) SendRaffleReminder(
    email string,
    data *RaffleReminderData,
) error {
    // Cargar plantilla
    html, err := n.templateLoader.RenderTemplate(
        "raffle_reminder.html",
        data,
    )
    if err != nil {
        return err
    }

    // Texto plano
    plainText := fmt.Sprintf(`
Hola %s,

Te recordamos que el sorteo "%s" será mañana a las %s.

Tus números: %s

Saludos,
Sorteos Platform
    `, data.FirstName, data.RaffleTitle, data.DrawDate,
       strings.Join(data.Numbers, ", "))

    // Enviar
    return n.sendEmail(
        email,
        "Recordatorio: Sorteo mañana - " + data.RaffleTitle,
        plainText,
        html,
    )
}
```

---

## 🔄 Flujo de Trabajo Recomendado

### **Para Desarrollo**
```
1. Editar templates/email/*.html
2. Probar localmente (enviar email de test)
3. Ajustar diseño según resultados
4. Repetir hasta perfecto
```

### **Para Producción**
```
1. Finalizar diseño en development
2. Copiar a internal/adapters/notifier/templates/
3. Recompilar backend
4. Deploy a producción
5. Monitor logs de emails enviados
```

---

## 📈 Roadmap de Plantillas Futuras

### **Próximas a Implementar**
- [ ] **raffle_reminder.html** - Recordatorio 24h antes
- [ ] **winner_notification.html** - ¡Ganaste!
- [ ] **raffle_completed.html** - Sorteo finalizado
- [ ] **reservation_expired.html** - Reserva expirada
- [ ] **raffle_cancelled.html** - Cancelación con reembolso
- [ ] **weekly_summary.html** - Resumen semanal
- [ ] **account_suspended.html** - Suspensión de cuenta

Ver código completo en: [PROPUESTA_EMAILS.md](PROPUESTA_EMAILS.md)

---

## 📁 Estructura de Archivos Final

```
/opt/Sorteos/backend/
├── templates/
│   └── email/                          # ← EDITAR AQUÍ
│       ├── verification.html           # ✅ Email verificación
│       ├── welcome.html                # ✅ Email bienvenida
│       ├── password_reset.html         # ✅ Reset password
│       ├── purchase_confirmation.html  # ✅ Confirmación compra
│       └── README.md                   # ✅ Documentación
│
├── internal/adapters/notifier/
│   ├── templates/                      # ← Para embeber
│   │   ├── verification.html
│   │   ├── welcome.html
│   │   └── password_reset.html
│   │
│   ├── template_loader.go              # ✅ Loader de plantillas
│   ├── smtp.go                         # ✅ Envío SMTP
│   ├── sendgrid.go                     # Envío SendGrid
│   └── notifier.go                     # Interface común
│
└── .env
    CONFIG_SENDGRID_TEMPLATES_DIR=/opt/Sorteos/backend/templates/email
```

---

## 💡 Tips Finales

1. **Backup antes de editar** - `cp verification.html verification.html.bak`
2. **Usa editor con syntax highlighting** - VSCode, Sublime, nano con colores
3. **Prueba en móvil** - Envía a tu Gmail y abre en celular
4. **Mantén consistencia** - Mismo header/footer en todos
5. **Documenta cambios** - Commit en Git con mensaje descriptivo
6. **Mide resultados** - Agrega UTM params para tracking

---

## 🆘 Troubleshooting

### **El HTML no se renderiza correctamente**
- Verifica que las variables `{{.Variable}}` coincidan con los nombres en Go
- Usa estilos inline, no CSS externo
- Prueba en [Litmus](https://litmus.com/) para ver en qué cliente falla

### **Los cambios no se aplican**
- Si usas archivos: Verifica la ruta en `.env`
- Si usas embebidas: Recompila el backend
- Revisa logs: `sudo journalctl -u sorteos-api -f`

### **Email va a spam**
- Verifica SPF/DKIM/DMARC en DNS
- No uses palabras spam ("GRATIS", "GANADOR", excesivos !!!)
- Mantén balance texto/imágenes (70% texto, 30% imágenes)

---

## 📚 Recursos Adicionales

- **Documentación interna:** [templates/email/README.md](backend/templates/email/README.md)
- **Propuesta de emails:** [PROPUESTA_EMAILS.md](PROPUESTA_EMAILS.md)
- **Guía SMTP vs SendGrid:** [GUIA_EMAIL_SMTP_VS_SENDGRID.md](GUIA_EMAIL_SMTP_VS_SENDGRID.md)
- **Tutorial HTML Email:** https://www.campaignmonitor.com/dev-resources/guides/coding-html-emails/
- **Email Client Support:** https://www.caniemail.com/

---

## ✨ Resumen Ejecutivo

**Tu sistema de plantillas:**
- ✅ **4 plantillas profesionales** creadas y listas
- ✅ **Fácil de editar** - Solo HTML, sin recompilar
- ✅ **Responsive** - Funciona en móvil y desktop
- ✅ **Compatible** - Gmail, Outlook, Apple Mail
- ✅ **Documentado** - README completo incluido
- ✅ **Extensible** - Fácil agregar nuevas plantillas

**Ubicación principal:**
```
/opt/Sorteos/backend/templates/email/
```

**Para editar:**
```bash
nano /opt/Sorteos/backend/templates/email/verification.html
```

**¡Listo para usar!** 🚀
