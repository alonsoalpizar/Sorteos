# Términos, Condiciones e Impacto Legal

**Versión:** 1.0
**Fecha:** 2025-11-10
**Advertencia:** Este documento es técnico. Requiere revisión legal profesional.

---

## 1. Introducción

Este documento describe el **impacto técnico** de los términos y condiciones en el sistema, incluyendo:

- Consentimientos requeridos
- Cumplimiento legal (GDPR, PCI DSS, etc.)
- Políticas de privacidad y cookies
- Términos de uso y limitaciones de responsabilidad

**Nota:** Este NO es el documento legal final. Debe ser redactado por un abogado especializado en:
- Legislación de sorteos/rifas en Costa Rica
- E-commerce internacional
- Protección de datos (GDPR, CCPA)
- Pagos electrónicos (PCI DSS)

---

## 2. Consentimientos Requeridos (GDPR Compliance)

### 2.1 Registro de Usuario

**Checkboxes obligatorios:**

```tsx
<form onSubmit={handleRegister}>
  <Input name="email" />
  <Input name="password" />

  {/* Obligatorio - Bloquea submit */}
  <Checkbox required>
    Acepto los <Link to="/terms">Términos y Condiciones</Link>
  </Checkbox>

  {/* Obligatorio - Bloquea submit */}
  <Checkbox required>
    He leído la <Link to="/privacy">Política de Privacidad</Link>
  </Checkbox>

  {/* Opcional - Marketing */}
  <Checkbox>
    Deseo recibir emails promocionales (puedes cancelar en cualquier momento)
  </Checkbox>

  <Button type="submit">Registrarse</Button>
</form>
```

**Almacenamiento en DB:**
```sql
CREATE TABLE user_consents (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id),
    consent_type VARCHAR(50) NOT NULL, -- 'terms', 'privacy', 'marketing'
    consent_version VARCHAR(20) NOT NULL, -- '1.0', '1.1'
    granted BOOLEAN NOT NULL,
    ip_address INET,
    user_agent TEXT,
    granted_at TIMESTAMP NOT NULL DEFAULT NOW()
);
```

**Código backend:**
```go
type RegisterUserInput struct {
    Email              string
    Password           string
    AcceptedTerms      bool
    AcceptedPrivacy    bool
    MarketingOptIn     bool
}

func (uc *RegisterUserUseCase) Execute(ctx context.Context, input RegisterUserInput) error {
    if !input.AcceptedTerms || !input.AcceptedPrivacy {
        return errors.New("debes aceptar términos y política de privacidad")
    }

    user := createUser(...)

    // Registrar consentimientos
    uc.consentRepo.Create(ctx, &UserConsent{
        UserID:        user.ID,
        ConsentType:   "terms",
        ConsentVersion: "1.0",
        Granted:       true,
        IPAddress:     c.ClientIP(),
        UserAgent:     c.Request.UserAgent(),
    })

    uc.consentRepo.Create(ctx, &UserConsent{
        UserID:        user.ID,
        ConsentType:   "privacy",
        ConsentVersion: "1.0",
        Granted:       true,
    })

    if input.MarketingOptIn {
        uc.consentRepo.Create(ctx, &UserConsent{
            UserID:      user.ID,
            ConsentType: "marketing",
            Granted:     true,
        })
    }

    return nil
}
```

---

### 2.2 Actualización de Términos

**Cuando términos cambian:**
1. Incrementar versión (`v1.0` → `v1.1`)
2. Al próximo login, mostrar banner:
   - "Hemos actualizado nuestros términos"
   - Usuario debe aceptar nuevamente
   - Si rechaza → logout automático

```tsx
function TermsUpdateBanner() {
  const { user } = useAuth()
  const [show, setShow] = useState(false)

  useEffect(() => {
    // Verificar si usuario aceptó última versión
    const latestVersion = '1.1'
    const userAcceptedVersion = user.consents.find(c => c.type === 'terms')?.version

    if (userAcceptedVersion !== latestVersion) {
      setShow(true)
    }
  }, [user])

  const handleAccept = async () => {
    await api.post('/users/me/consents', {
      consent_type: 'terms',
      consent_version: '1.1',
      granted: true,
    })
    setShow(false)
  }

  if (!show) return null

  return (
    <Alert variant="warning">
      <AlertTitle>Nuevos Términos y Condiciones</AlertTitle>
      <AlertDescription>
        Hemos actualizado nuestros términos. Por favor revísalos y acepta para continuar.
      </AlertDescription>
      <div className="mt-4 flex gap-2">
        <Button variant="outline" asChild>
          <Link to="/terms">Leer Términos</Link>
        </Button>
        <Button onClick={handleAccept}>Aceptar</Button>
      </div>
    </Alert>
  )
}
```

---

## 3. GDPR - Derechos del Usuario

### 3.1 Derecho al Acceso (Export Data)

**Endpoint:** `GET /users/me/data-export`

**Retorna:**
```json
{
  "user": {
    "id": 12345,
    "email": "user@example.com",
    "phone": "+50612345678",
    "created_at": "2025-01-01T00:00:00Z"
  },
  "raffles": [
    { "id": 1, "title": "iPhone 15", "status": "completed" }
  ],
  "purchases": [
    { "raffle_id": 2, "numbers": ["01", "15"], "amount": 10.00, "date": "2025-01-15" }
  ],
  "payments": [
    { "id": 123, "amount": 10.00, "status": "succeeded", "date": "2025-01-15" }
  ],
  "consents": [
    { "type": "terms", "version": "1.0", "granted_at": "2025-01-01" }
  ]
}
```

**Implementación:**
```go
func ExportUserData(c *gin.Context) {
    userID := c.MustGet("user_id").(int64)

    data := map[string]interface{}{
        "user":      userRepo.FindByID(ctx, userID),
        "raffles":   raffleRepo.FindByUserID(ctx, userID),
        "purchases": reservationRepo.FindByUserID(ctx, userID),
        "payments":  paymentRepo.FindByUserID(ctx, userID),
        "consents":  consentRepo.FindByUserID(ctx, userID),
    }

    c.JSON(200, data)
}
```

---

### 3.2 Derecho al Olvido (Delete Account)

**Endpoint:** `DELETE /users/me`

**Proceso:**
1. Marcar user.status = `deleted`
2. **Anonimizar** (no eliminar físicamente):
   - email → `deleted_12345@anonymous.com`
   - phone → `NULL`
   - password_hash → `NULL`
   - Conservar transacciones (requerimiento fiscal: 7 años)
3. Eliminar sesiones activas
4. Cancelar sorteos activos
5. Notificar vía email (confirmación)

```go
func (uc *DeleteUserUseCase) Execute(ctx context.Context, userID int64) error {
    user, _ := uc.userRepo.FindByID(ctx, userID)

    // Anonimizar
    user.Email = fmt.Sprintf("deleted_%d@anonymous.com", userID)
    user.Phone = ""
    user.PasswordHash = ""
    user.Status = UserStatusDeleted
    uc.userRepo.Update(ctx, user)

    // Cancelar sorteos activos
    raffles := uc.raffleRepo.FindActiveByUser(ctx, userID)
    for _, raffle := range raffles {
        uc.cancelRaffleUseCase.Execute(ctx, raffle.ID, "usuario eliminó su cuenta")
    }

    // Revocar sesiones
    uc.tokenManager.RevokeAllTokens(ctx, userID)

    // Notificar (último email)
    uc.notifier.SendEmail(ctx, user.Email, "account_deleted", nil)

    // Auditoría
    logger.Info("user_deleted", zap.Int64("user_id", userID))

    return nil
}
```

**Retención de datos (fiscal):**
- Transacciones de pago: 7 años
- Audit logs: 2 años
- Datos anonimizados: Indefinido

---

## 4. PCI DSS Compliance (Pagos)

### 4.1 Delegación a PSP

**Nunca almacenar:**
- Número completo de tarjeta (PAN)
- CVV/CVC
- PIN

**Permitido:**
- Últimos 4 dígitos (para UI)
- Token de Stripe (`pm_xxx`)
- Marca de tarjeta (Visa, Mastercard)

**Implementación:**
```go
type PaymentMethod struct {
    ID         int64
    UserID     int64
    Provider   string // "stripe"
    ExternalID string // "pm_1xxx" (token de Stripe)
    Last4      string // "4242"
    Brand      string // "visa"
    IsDefault  bool
}

// NUNCA almacenar esto:
type ForbiddenData struct {
    CardNumber string // ❌
    CVV        string // ❌
}
```

**Stripe Elements (Frontend):**
```tsx
import { CardElement } from '@stripe/react-stripe-js'

// CardElement NO envía datos de tarjeta al backend
// Solo tokeniza y envía pm_xxx al backend
<CardElement />
```

---

### 4.2 Comunicación Segura

- **TLS 1.3** obligatorio
- **HSTS** activado
- **Certificate Pinning** en app móvil (futuro)

---

## 5. Cookies y Tracking

### 5.1 Banner de Cookies

**GDPR requiere:**
- Informar sobre uso de cookies
- Permitir rechazo (excepto cookies esenciales)
- Categorías: Esenciales, Analíticas, Marketing

```tsx
function CookieBanner() {
  const [preferences, setPreferences] = useState({
    essential: true, // siempre activadas
    analytics: false,
    marketing: false,
  })

  const handleAcceptAll = () => {
    setPreferences({ essential: true, analytics: true, marketing: true })
    saveCookiePreferences(preferences)
  }

  const handleSavePreferences = () => {
    saveCookiePreferences(preferences)
  }

  return (
    <div className="fixed bottom-0 left-0 right-0 bg-neutral-900 text-white p-6">
      <h3 className="font-semibold">Este sitio usa cookies</h3>
      <p className="text-sm mt-2">
        Usamos cookies esenciales para el funcionamiento del sitio y cookies
        opcionales para mejorar tu experiencia.
      </p>

      <div className="mt-4 space-y-2">
        <Checkbox checked disabled>Esenciales (requeridas)</Checkbox>
        <Checkbox
          checked={preferences.analytics}
          onCheckedChange={(checked) => setPreferences({ ...preferences, analytics: checked })}
        >
          Analíticas (Google Analytics)
        </Checkbox>
        <Checkbox
          checked={preferences.marketing}
          onCheckedChange={(checked) => setPreferences({ ...preferences, marketing: checked })}
        >
          Marketing (Meta Pixel, Google Ads)
        </Checkbox>
      </div>

      <div className="mt-4 flex gap-2">
        <Button onClick={handleAcceptAll}>Aceptar Todas</Button>
        <Button variant="outline" onClick={handleSavePreferences}>
          Guardar Preferencias
        </Button>
      </div>
    </div>
  )
}
```

---

### 5.2 Almacenamiento de Preferencias

```tsx
function saveCookiePreferences(preferences: CookiePreferences) {
  // Guardar en localStorage (no afecta a terceros)
  localStorage.setItem('cookie_preferences', JSON.stringify(preferences))

  // Cargar scripts según preferencias
  if (preferences.analytics) {
    loadGoogleAnalytics()
  }

  if (preferences.marketing) {
    loadMetaPixel()
    loadGoogleAds()
  }
}
```

---

## 6. Políticas del Sistema

### 6.1 Términos de Uso (Resumen Técnico)

**Secciones requeridas:**
1. **Elegibilidad**: Mayores de 18 años
2. **Uso permitido**: Publicar sorteos legales, comprar boletos
3. **Uso prohibido**:
   - Fraude, lavado de dinero
   - Sorteos ilegales (armas, drogas)
   - Manipulación de resultados
4. **Comisiones**: 5% del total recaudado
5. **Responsabilidad**: Plataforma no garantiza premios (responsabilidad del owner)
6. **Resolución de disputas**: Arbitraje en Costa Rica
7. **Cancelación**: Admin puede suspender cuentas sin previo aviso

**Validación técnica:**
```go
func (uc *CreateRaffleUseCase) Validate(raffle *Raffle) error {
    // Verificar edad (18+)
    user := uc.userRepo.FindByID(ctx, raffle.UserID)
    if user.Age < 18 {
        return errors.New("debes ser mayor de 18 años")
    }

    // Verificar contenido prohibido (palabras clave)
    prohibitedWords := []string{"arma", "droga", "apuesta"}
    for _, word := range prohibitedWords {
        if strings.Contains(strings.ToLower(raffle.Title), word) {
            return errors.New("contenido prohibido detectado")
        }
    }

    return nil
}
```

---

### 6.2 Política de Privacidad (Resumen Técnico)

**Datos recopilados:**
- Email, teléfono, nombre
- Datos de pago (delegados a Stripe)
- Dirección IP, user agent
- Actividad en la plataforma (sorteos, compras)

**Uso de datos:**
- Procesar transacciones
- Verificar identidad (KYC)
- Enviar notificaciones
- Analítica (Google Analytics)
- Marketing (con consentimiento)

**Compartimos datos con:**
- PSP (Stripe, PayPal)
- Email provider (SendGrid)
- SMS provider (Twilio)
- **NO** vendemos datos a terceros

**Retención:**
- Cuenta activa: Indefinido
- Cuenta eliminada: Anonimizado tras 30 días
- Transacciones: 7 años (fiscal)

---

## 7. Limitación de Responsabilidad

**Cláusulas críticas (técnicas):**

**1. Disponibilidad del servicio:**
> "La plataforma se ofrece 'AS IS' sin garantías. No garantizamos disponibilidad 100%. Mantenimiento programado con aviso de 24h."

**Implementación:**
```tsx
// Banner de mantenimiento
function MaintenanceBanner({ scheduledAt, duration }) {
  return (
    <Alert variant="warning">
      Mantenimiento programado: {scheduledAt} (duración: {duration})
    </Alert>
  )
}
```

**2. Entrega de premios:**
> "La plataforma facilita el sorteo. La entrega del premio es responsabilidad exclusiva del organizador. No garantizamos calidad, autenticidad ni entrega."

**Implementación:**
```tsx
// Disclaimer en página de sorteo
<Alert variant="info">
  <InfoIcon />
  Este sorteo es organizado por {owner.name}. La plataforma no garantiza
  la entrega del premio. Contacta al organizador para detalles.
</Alert>
```

**3. Resultados de sorteos:**
> "Los resultados se basan en fuentes oficiales (Lotería Nacional de CR). No somos responsables por errores en la fuente."

---

## 8. Cumplimiento por Jurisdicción

### 8.1 Costa Rica

**Requisitos:**
- Licencia de operación (Ministerio de Hacienda)
- Impuesto sobre sorteos/rifas (TBD)
- Restricciones de premios (no alcohol, tabaco)

**Acción técnica:**
- Filtrar categorías prohibidas
- Retención de impuestos en liquidaciones

---

### 8.2 Internacional (si se expande)

**GDPR (Europa):**
- ✅ Consentimientos explícitos
- ✅ Derecho al olvido
- ✅ Exportar datos
- ✅ DPO (Data Protection Officer) si > 250 empleados

**CCPA (California):**
- Permitir "No vender mis datos"
- Revelar qué datos se recopilan

---

## 9. Auditorías y Cumplimiento

### 9.1 Auditoría Anual

**Verificar:**
- [ ] Todos los consentimientos registrados en DB
- [ ] Términos actualizados (última versión)
- [ ] Política de privacidad publicada
- [ ] Cookies con opt-out funcional
- [ ] Datos PII cifrados o anonimizados
- [ ] Backup de DB cifrado
- [ ] Logs de auditoría completos

---

### 9.2 Reporte de Transparencia (Anual)

**Publicar:**
- Número de usuarios activos
- Número de sorteos realizados
- Número de solicitudes de eliminación de datos (GDPR)
- Número de disputas resueltas
- Tiempo de respuesta promedio

---

## 10. Checklist Legal (Pre-Launch)

- [ ] Contratar abogado especializado en e-commerce
- [ ] Redactar términos y condiciones finales
- [ ] Redactar política de privacidad
- [ ] Redactar política de cookies
- [ ] Obtener licencia de operación (si requerida en CR)
- [ ] Registrar marca comercial
- [ ] Contratar seguro de responsabilidad (opcional)
- [ ] Configurar DPO o email de contacto para privacidad
- [ ] Publicar términos en /terms (accesible sin login)
- [ ] Publicar política de privacidad en /privacy

---

## 11. Contacto Legal

**Email de privacidad:** privacy@sorteos.com
**Email de soporte:** support@sorteos.com

**Formulario de GDPR:**
```tsx
<form action="/privacy/request">
  <Select name="request_type">
    <SelectItem value="export">Exportar mis datos</SelectItem>
    <SelectItem value="delete">Eliminar mi cuenta</SelectItem>
    <SelectItem value="correct">Corregir datos incorrectos</SelectItem>
    <SelectItem value="complaint">Presentar queja</SelectItem>
  </Select>
  <Textarea name="details" placeholder="Detalles..." />
  <Button type="submit">Enviar Solicitud</Button>
</form>
```

**Tiempo de respuesta:** Máximo 30 días (GDPR)

---

**Actualizado:** 2025-11-10
**Próxima revisión legal:** Antes del lanzamiento
