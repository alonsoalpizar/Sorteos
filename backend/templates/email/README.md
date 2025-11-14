# Plantillas de Email - Sorteos Platform

## 📧 Ubicación de las Plantillas

Las plantillas de email se encuentran en:

```
/opt/Sorteos/backend/templates/email/
├── verification.html           # Email de verificación de cuenta
├── welcome.html               # Email de bienvenida
├── password_reset.html        # Email de reset de contraseña
├── purchase_confirmation.html # Email de confirmación de compra
└── README.md                  # Este archivo
```

---

## 🎨 Plantillas Disponibles

### 1. **verification.html** - Verificación de Email
**Variables disponibles:**
```go
{{.FirstName}}        // Nombre del usuario
{{.Code}}             // Código de 6 dígitos
{{.FrontendURL}}      // URL del frontend
{{.VerificationURL}}  // URL directa de verificación (opcional)
```

**Uso:**
```go
data := VerificationEmailData{
    FirstName:   "Juan",
    Code:        "123456",
    FrontendURL: "https://sorteos.club",
}
```

---

### 2. **welcome.html** - Bienvenida
**Variables disponibles:**
```go
{{.FirstName}}    // Nombre del usuario
{{.FrontendURL}}  // URL del frontend
```

**Uso:**
```go
data := WelcomeEmailData{
    FirstName:   "Juan",
    FrontendURL: "https://sorteos.club",
}
```

---

### 3. **password_reset.html** - Reset de Contraseña
**Variables disponibles:**
```go
{{.FirstName}}    // Nombre del usuario (opcional)
{{.ResetURL}}     // URL completa para resetear
{{.FrontendURL}}  // URL del frontend
```

**Uso:**
```go
data := PasswordResetEmailData{
    FirstName:   "Juan",
    ResetURL:    "https://sorteos.club/reset-password?token=abc123",
    FrontendURL: "https://sorteos.club",
}
```

---

### 4. **purchase_confirmation.html** - Confirmación de Compra
**Variables disponibles:**
```go
{{.FirstName}}    // Nombre del usuario
{{.RaffleTitle}}  // Título del sorteo
{{.RaffleID}}     // ID del sorteo (para link)
{{.Numbers}}      // Array de números comprados
{{.TotalAmount}}  // Monto total pagado
{{.DrawDate}}     // Fecha formateada del sorteo
{{.Prize}}        // Descripción del premio
{{.FrontendURL}}  // URL del frontend
```

**Uso:**
```go
data := PurchaseConfirmationData{
    FirstName:   "Juan",
    RaffleTitle: "Gran Sorteo de Navidad",
    RaffleID:    "uuid-123",
    Numbers:     []string{"0001", "0042", "0099"},
    TotalAmount: "$150.00",
    DrawDate:    "25/12/2025 20:00",
    Prize:       "iPhone 15 Pro Max",
    FrontendURL: "https://sorteos.club",
}
```

---

## 🛠️ Cómo Usar las Plantillas

### **Opción 1: Cargar desde Archivos (Recomendado)**

```go
// En tu código Go
import "github.com/sorteos-platform/backend/internal/adapters/notifier"

// Crear loader
loader := notifier.NewTemplateLoader("/opt/Sorteos/backend/templates/email")

// Renderizar plantilla
html, err := loader.RenderTemplate("verification.html", data)
if err != nil {
    log.Error("Error rendering template", err)
}

// Enviar email con el HTML
sendEmail(to, subject, html)
```

---

### **Opción 2: Plantillas Embebidas (Producción)**

Las plantillas se pueden embeber en el binario Go para no depender de archivos externos:

```go
//go:embed templates/*.html
var embeddedTemplates embed.FS

loader := notifier.NewTemplateLoader("") // "" = usar embebidas
```

**Ventaja:** No requiere archivos externos, todo en el binario.

---

## 🎨 Personalizar Plantillas

### **Colores del Tema**

Los colores principales usados:
- **Azul primario:** `#3B82F6` - Botones, headers
- **Verde éxito:** `#10B981` - Confirmaciones, bienvenida
- **Rojo error:** `#EF4444` - Alertas, reset password
- **Amarillo info:** `#F59E0B` - Premios, información
- **Gris texto:** `#333333` - Texto principal
- **Gris secundario:** `#64748B` - Texto secundario

### **Modificar Diseño**

1. Editar el archivo HTML directamente
2. Mantener la estructura de tablas para compatibilidad con clientes de email
3. Usar estilos inline (no CSS externo)
4. Reiniciar backend si usas archivos (no embebidas)

---

## ✉️ Mejores Prácticas

### **1. Usar Tablas para Layout**
```html
<!-- ✅ BIEN: Compatible con todos los clientes -->
<table width="100%">
    <tr>
        <td>Contenido</td>
    </tr>
</table>

<!-- ❌ MAL: No funciona en Outlook -->
<div style="display: flex;">Contenido</div>
```

### **2. Estilos Inline**
```html
<!-- ✅ BIEN -->
<p style="color: #333; font-size: 16px;">Texto</p>

<!-- ❌ MAL -->
<style>p { color: #333; }</style>
<p>Texto</p>
```

### **3. Texto Alternativo**
Siempre incluir versión de texto plano además del HTML.

### **4. Responsive**
```html
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<table width="600" style="max-width: 100%;">
```

---

## 🧪 Probar Plantillas

### **Método 1: Herramientas Online**
- [Litmus](https://litmus.com/) - Testing en múltiples clientes
- [Email on Acid](https://www.emailonacid.com/) - Preview en tiempo real
- [Mailtrap](https://mailtrap.io/) - Sandbox para desarrollo

### **Método 2: Enviar Email de Prueba**
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"test@gmail.com","password":"Test123@","accepted_terms":true,"accepted_privacy":true}'
```

### **Método 3: Renderizar sin Enviar**
```go
// En tests
html, _ := loader.RenderTemplate("verification.html", testData)
fmt.Println(html) // Ver HTML generado
```

---

## 📁 Estructura de Directorios

```
backend/
├── templates/
│   └── email/                    # Plantillas editables
│       ├── verification.html
│       ├── welcome.html
│       ├── password_reset.html
│       ├── purchase_confirmation.html
│       └── README.md             # Este archivo
│
└── internal/adapters/notifier/
    ├── templates/                # Copias para embeber
    │   ├── verification.html
    │   ├── welcome.html
    │   └── password_reset.html
    │
    ├── template_loader.go        # Loader de plantillas
    ├── smtp.go                   # Envío por SMTP
    └── sendgrid.go              # Envío por SendGrid
```

---

## 🚀 Agregar Nueva Plantilla

### **Paso 1: Crear HTML**
```bash
cd /opt/Sorteos/backend/templates/email
nano nueva_plantilla.html
```

### **Paso 2: Definir Estructura de Datos**
```go
// En template_loader.go
type NuevaPlantillaData struct {
    Campo1 string
    Campo2 int
    // ...
}
```

### **Paso 3: Crear Método en Notifier**
```go
// En smtp.go o sendgrid.go
func (n *SMTPNotifier) SendNuevaPlantilla(email string, data *NuevaPlantillaData) error {
    html, err := n.templateLoader.RenderTemplate("nueva_plantilla.html", data)
    if err != nil {
        return err
    }

    return n.sendEmail(email, "Asunto", plainText, html)
}
```

### **Paso 4: Copiar para Embeber (Opcional)**
```bash
cp nueva_plantilla.html ../internal/adapters/notifier/templates/
```

---

## 🔄 Actualizar Plantillas en Producción

### **Si usas archivos:**
1. Editar archivo HTML
2. Los cambios se aplican inmediatamente (próximo email)

### **Si usas embebidas:**
1. Editar archivo HTML
2. Copiar a `internal/adapters/notifier/templates/`
3. Recompilar backend: `go build`
4. Reiniciar servicio: `sudo systemctl restart sorteos-api`

---

## 📊 Métricas de Email

Para trackear opens/clicks, agregar parámetros UTM:

```html
<a href="{{.FrontendURL}}/raffles?utm_source=email&utm_medium=purchase_confirmation&utm_campaign=transactional">
    Ver Sorteo
</a>
```

---

## 🎯 Roadmap de Plantillas

- [ ] reminder_24h.html - Recordatorio 24h antes del sorteo
- [ ] winner_notification.html - Notificación de ganador
- [ ] raffle_completed.html - Sorteo completado (no ganaste)
- [ ] reservation_expired.html - Reserva expirada
- [ ] weekly_summary.html - Resumen semanal
- [ ] raffle_cancelled.html - Cancelación de sorteo

---

## 💡 Tips

1. **Mantén simple el diseño** - Los clientes de email tienen soporte limitado de CSS
2. **Prueba en múltiples clientes** - Gmail, Outlook, Apple Mail, etc.
3. **Incluye siempre texto plano** - Algunos usuarios prefieren texto
4. **Usa colores accesibles** - Contraste suficiente para lectura
5. **Optimiza peso** - Evita imágenes pesadas en línea
6. **Agrega unsubscribe link** - Para emails promocionales

---

## 📧 Contacto

¿Necesitas ayuda con las plantillas?
- Revisa [PROPUESTA_EMAILS.md](../../../PROPUESTA_EMAILS.md)
- Consulta [GUIA_EMAIL_SMTP_VS_SENDGRID.md](../../../GUIA_EMAIL_SMTP_VS_SENDGRID.md)
