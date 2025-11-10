# Seguridad - Plataforma de Sorteos

**Versión:** 1.0
**Fecha:** 2025-11-10
**Nivel de Criticidad:** ALTA (manejo de pagos y datos personales)

---

## 1. Principios de Seguridad

1. **Defense in Depth**: Múltiples capas de seguridad
2. **Least Privilege**: Permisos mínimos necesarios
3. **Fail Secure**: En caso de error, denegar acceso
4. **Security by Design**: Seguridad desde el diseño, no agregada después
5. **Zero Trust**: Nunca confiar, siempre verificar

---

## 2. Autenticación y Autorización

### 2.1 JWT (JSON Web Tokens)

**Access Token:**
- Algoritmo: `HS256` (HMAC-SHA256) para MVP, migrar a `RS256` (RSA) en producción
- Expiración: **15 minutos**
- Payload:
```json
{
  "user_id": 12345,
  "role": "user",
  "kyc_level": "email_verified",
  "iat": 1699632000,
  "exp": 1699632900
}
```

**Refresh Token:**
- Expiración: **7 días**
- Almacenado en Redis con TTL
- Rotación obligatoria al usar (invalida token anterior)
- Revocable por `jti` (JWT ID único)

**Implementación Go:**
```go
func GenerateAccessToken(userID int64, role string) (string, error) {
    secret := []byte(config.JWTSecret)
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "user_id": userID,
        "role": role,
        "exp": time.Now().Add(15 * time.Minute).Unix(),
        "iat": time.Now().Unix(),
    })
    return token.SignedString(secret)
}
```

**Almacenamiento en Frontend:**
- Access token: **Memory** (variable en React state)
- Refresh token: **HttpOnly Cookie** (previene XSS)
- Nunca en `localStorage` (vulnerable a XSS)

---

### 2.2 RBAC (Control de Acceso Basado en Roles)

**Roles:**
- `user`: Usuario estándar
- `admin`: Administrador (Almighty)

**Middleware de Autorización:**
```go
func RequireRole(allowedRoles ...string) gin.HandlerFunc {
    return func(c *gin.Context) {
        claims := c.MustGet("claims").(*TokenClaims)

        for _, role := range allowedRoles {
            if claims.Role == role {
                c.Next()
                return
            }
        }

        c.JSON(403, gin.H{"error": "forbidden"})
        c.Abort()
    }
}

// Uso
r.PATCH("/admin/raffles/:id", RequireRole("admin"), handlers.SuspendRaffle)
```

**Validación de Ownership:**
```go
func RequireOwnership(resourceType string) gin.HandlerFunc {
    return func(c *gin.Context) {
        claims := c.MustGet("claims").(*TokenClaims)
        resourceID := c.Param("id")

        // Verificar que claims.UserID == resource.OwnerID
        // ...
    }
}
```

---

### 2.3 Verificación de Identidad (KYC)

**Niveles:**
1. `none`: Registrado pero no verificado
2. `email_verified`: Email confirmado
3. `phone_verified`: Teléfono confirmado
4. `full_kyc`: ID + selfie + verificación manual (futuro)

**Restricciones por nivel:**
- Crear sorteo: `>= email_verified`
- Comprar boletos: `>= email_verified`
- Retirar fondos: `>= full_kyc` (futuro)

**Tokens de verificación:**
```go
// Token de verificación (1 hora de vida)
token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
    "user_id": userID,
    "type": "email_verification",
    "exp": time.Now().Add(1 * time.Hour).Unix(),
})
```

---

## 3. Rate Limiting

### 3.1 Límites por Endpoint

| Endpoint | Límite | Ventana | Razón |
|----------|--------|---------|-------|
| POST /auth/login | 5 req/min | Por IP | Prevenir brute force |
| POST /auth/register | 3 req/hora | Por IP | Prevenir spam de cuentas |
| POST /raffles/:id/reservations | 10 req/min | Por user_id | Prevenir abuse |
| POST /payments | 5 req/min | Por user_id | Prevenir intentos fraudulentos |
| GET /raffles | 60 req/min | Por IP | Tráfico normal |

### 3.2 Implementación con Redis

**Token Bucket Algorithm:**
```go
func RateLimitMiddleware(maxRequests int, window time.Duration) gin.HandlerFunc {
    return func(c *gin.Context) {
        key := fmt.Sprintf("ratelimit:%s:%s", c.Request.URL.Path, c.ClientIP())

        count, err := rdb.Incr(ctx, key).Result()
        if err != nil {
            c.AbortWithStatus(500)
            return
        }

        if count == 1 {
            rdb.Expire(ctx, key, window)
        }

        if count > int64(maxRequests) {
            c.Header("X-RateLimit-Limit", strconv.Itoa(maxRequests))
            c.Header("X-RateLimit-Remaining", "0")
            c.JSON(429, gin.H{"error": "too many requests"})
            c.Abort()
            return
        }

        c.Header("X-RateLimit-Limit", strconv.Itoa(maxRequests))
        c.Header("X-RateLimit-Remaining", strconv.Itoa(maxRequests-int(count)))
        c.Next()
    }
}
```

---

## 4. Validación y Sanitización

### 4.1 Validación de Entrada

**Backend (Go):**
```go
type CreateRaffleRequest struct {
    Title       string          `json:"title" validate:"required,min=5,max=200"`
    Description string          `json:"description" validate:"required,min=20,max=2000"`
    DrawDate    time.Time       `json:"draw_date" validate:"required,future"`
    Price       decimal.Decimal `json:"price" validate:"required,gt=0,lte=10000"`
}

func ValidateRequest(c *gin.Context, req interface{}) error {
    if err := c.ShouldBindJSON(req); err != nil {
        return err
    }

    validate := validator.New()
    if err := validate.Struct(req); err != nil {
        return err
    }

    return nil
}
```

**Frontend (React + Zod):**
```tsx
const createRaffleSchema = z.object({
  title: z.string().min(5).max(200),
  description: z.string().min(20).max(2000),
  drawDate: z.date().refine((date) => date > new Date(), {
    message: "La fecha debe ser futura",
  }),
  price: z.number().positive().max(10000),
})

const { register, handleSubmit, formState: { errors } } = useForm({
  resolver: zodResolver(createRaffleSchema),
})
```

### 4.2 Sanitización

**HTML/XSS:**
- **Nunca** renderizar HTML no sanitizado
- Usar `textContent` en lugar de `innerHTML`
- React escapa por defecto, pero evitar `dangerouslySetInnerHTML`

**SQL Injection:**
- **Siempre** usar prepared statements (GORM lo hace automáticamente)
- Nunca concatenar strings en queries

**Command Injection:**
- Nunca ejecutar comandos shell con input del usuario
- Si es necesario, usar whitelist estricta

---

## 5. Seguridad de Datos

### 5.1 Encriptación en Tránsito

**TLS 1.3:**
- Certificado SSL/TLS válido (Let's Encrypt)
- HSTS (HTTP Strict Transport Security): `max-age=31536000; includeSubDomains`
- Redirección automática HTTP → HTTPS

**Nginx Config:**
```nginx
server {
    listen 443 ssl http2;
    ssl_certificate /etc/letsencrypt/live/sorteos.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/sorteos.com/privkey.pem;
    ssl_protocols TLSv1.3;
    add_header Strict-Transport-Security "max-age=31536000; includeSubDomains" always;
}
```

---

### 5.2 Encriptación en Reposo

**Base de Datos:**
- PostgreSQL con encriptación a nivel de disco (LUKS o AWS RDS encryption)
- Datos sensibles (tarjetas) **nunca** se almacenan (usar tokens de Stripe)
- Contraseñas: **bcrypt** con costo 12

**Contraseñas:**
```go
import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) (string, error) {
    hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
    return string(hash), err
}

func VerifyPassword(hash, password string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    return err == nil
}
```

**Datos PII (Personally Identifiable Information):**
- Email, teléfono, dirección: cifrado a nivel de aplicación (opcional)
- Backup de DB cifrado con clave rotable

---

### 5.3 Secrets Management

**Variables de Entorno:**
```bash
# .env (NUNCA commitear)
CONFIG_JWT_SECRET=xxxx-yyyy-zzzz-random-64-chars
CONFIG_DB_PASSWORD=xxxx
CONFIG_STRIPE_SECRET_KEY=sk_live_xxxx
CONFIG_SENDGRID_API_KEY=SG.xxxx
```

**Rotación de Secretos:**
- JWT Secret: Rotación manual cada 90 días
- DB Password: Rotación automática con AWS Secrets Manager
- API Keys de terceros: Rotación según política del proveedor

---

## 6. Protección contra Vulnerabilidades OWASP Top 10

### 6.1 A01: Broken Access Control

**Mitigación:**
- Middleware de autorización en **todos** los endpoints protegidos
- Validar ownership en operaciones CRUD
- Logs de auditoría para acciones sensibles

### 6.2 A02: Cryptographic Failures

**Mitigación:**
- TLS 1.3 obligatorio
- bcrypt para passwords
- No almacenar datos de tarjetas (usar Stripe tokens)

### 6.3 A03: Injection (SQL, NoSQL, Command)

**Mitigación:**
- GORM con prepared statements
- Validación estricta de entrada
- Whitelist en lugar de blacklist

### 6.4 A04: Insecure Design

**Mitigación:**
- Threat modeling en fase de diseño
- Peer review de features sensibles (pagos, reservas)
- Tests de seguridad en CI

### 6.5 A05: Security Misconfiguration

**Mitigación:**
- Desactivar endpoints de debug en producción
- Configuración de CORS estricta:
```go
r.Use(cors.New(cors.Config{
    AllowOrigins: []string{"https://sorteos.com"},
    AllowMethods: []string{"GET", "POST", "PATCH", "DELETE"},
    AllowHeaders: []string{"Authorization", "Content-Type"},
    AllowCredentials: true,
}))
```
- Headers de seguridad:
```go
c.Header("X-Content-Type-Options", "nosniff")
c.Header("X-Frame-Options", "DENY")
c.Header("X-XSS-Protection", "1; mode=block")
c.Header("Content-Security-Policy", "default-src 'self'")
```

### 6.6 A06: Vulnerable and Outdated Components

**Mitigación:**
- Dependabot / Renovate para actualizaciones automáticas
- `go mod verify` y `npm audit` en CI
- Trivy para escaneo de vulnerabilidades en Docker images

### 6.7 A07: Identification and Authentication Failures

**Mitigación:**
- MFA (Multi-Factor Authentication) para admins (futuro)
- Bloqueo de cuenta tras 5 intentos fallidos
- Logout en todos los dispositivos al cambiar contraseña

### 6.8 A08: Software and Data Integrity Failures

**Mitigación:**
- Verificar firmas de webhooks (Stripe signature)
- Subresource Integrity (SRI) en CDN:
```html
<script src="https://cdn.example.com/lib.js"
  integrity="sha384-oqVuAfXRKap7fdgcCY5uykM6+R9GqQ8K/ux..."
  crossorigin="anonymous"></script>
```

### 6.9 A09: Security Logging and Monitoring Failures

**Mitigación:**
- Logs estructurados con trace_id
- Alertas en Prometheus para eventos críticos:
  - Login fallido > 10 veces en 5 min
  - Payment failure > 20% en 1 hora
  - Reserva con doble venta (debe ser 0)

### 6.10 A10: Server-Side Request Forgery (SSRF)

**Mitigación:**
- Validar URLs en webhooks
- Whitelist de dominios permitidos para callbacks
- No seguir redirects automáticamente

---

## 7. Seguridad en Pagos

**Ver:** [pagos_integraciones.md](./pagos_integraciones.md) para detalles completos.

**Resumen:**
- **Nunca** almacenar números de tarjeta completos
- Usar tokens de Stripe/PayPal
- PCI DSS compliance delegado a PSP
- Webhooks firmados (HMAC-SHA256)
- Idempotencia con `Idempotency-Key`

---

## 8. Auditoría y Logs

### 8.1 Logs de Auditoría

**Eventos auditables:**
- Login/logout
- Cambio de contraseña
- Creación/suspensión de sorteo
- Pago procesado
- Acción de admin (suspender usuario, etc.)

**Tabla audit_logs:**
```sql
CREATE TABLE audit_logs (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT REFERENCES users(id),
    action VARCHAR(100) NOT NULL,
    entity_type VARCHAR(50),
    entity_id BIGINT,
    ip_address INET,
    user_agent TEXT,
    metadata JSONB,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_audit_logs_user_action ON audit_logs(user_id, action, created_at);
```

**Logging Go:**
```go
func LogAudit(ctx context.Context, userID int64, action string, metadata map[string]interface{}) {
    logger.Info("audit_log",
        zap.Int64("user_id", userID),
        zap.String("action", action),
        zap.String("ip", c.ClientIP()),
        zap.Any("metadata", metadata),
    )

    db.Create(&AuditLog{
        UserID: userID,
        Action: action,
        Metadata: metadata,
        IPAddress: c.ClientIP(),
        // ...
    })
}
```

---

### 8.2 Logging Estructurado

**Formato:**
```json
{
  "level": "info",
  "timestamp": "2025-11-10T10:30:00Z",
  "trace_id": "abc123",
  "user_id": 456,
  "action": "reserve_numbers",
  "raffle_id": 789,
  "numbers": ["01", "15"],
  "duration_ms": 124
}
```

**Correlación con trace_id:**
- Generar UUID en cada request
- Propagar en headers (`X-Trace-ID`)
- Incluir en todos los logs de esa request

---

## 9. GDPR y Privacidad

### 9.1 Cumplimiento GDPR

**Derechos del usuario:**
1. **Derecho al acceso**: GET /users/me/data (exportar datos)
2. **Derecho a la rectificación**: PATCH /users/me
3. **Derecho al olvido**: DELETE /users/me (anonimizar, no eliminar físicamente)
4. **Portabilidad**: JSON export de todos los datos

**Consentimiento:**
- Checkbox obligatorio en registro: "Acepto términos y condiciones"
- Opción de newsletter separada (opt-in)

### 9.2 Retención de Datos

- Datos de usuario activo: Indefinido
- Datos de usuario eliminado: Anonimizar tras 30 días
- Logs de auditoría: 2 años (requerimiento fiscal)
- Transacciones de pago: 7 años (requerimiento legal)

---

## 10. Penetration Testing

### 10.1 Tests Automáticos (CI)

**OWASP ZAP:**
```yaml
# .github/workflows/security.yml
- name: ZAP Scan
  uses: zaproxy/action-baseline@v0.7.0
  with:
    target: 'https://staging.sorteos.com'
```

**Trivy (Scan de vulnerabilidades):**
```yaml
- name: Run Trivy
  run: trivy image --severity HIGH,CRITICAL backend:latest
```

### 10.2 Tests Manuales (Trimestrales)

- [ ] Brute force en login
- [ ] SQL Injection en formularios
- [ ] XSS en campos de texto
- [ ] CSRF en acciones críticas
- [ ] Doble venta con requests concurrentes
- [ ] Escalación de privilegios (user → admin)

---

## 11. Incident Response Plan

**Niveles de incidentes:**
1. **Crítico**: Doble venta, brecha de datos, acceso no autorizado
2. **Alto**: Payment failure > 50%, DoS
3. **Medio**: Endpoint caído, bug no crítico
4. **Bajo**: Typo en UI, bug cosmético

**Protocolo:**
1. Detectar (alertas Prometheus)
2. Contener (rollback, feature flag off)
3. Investigar (logs, traces)
4. Remediar (fix, deploy)
5. Comunicar (usuarios afectados, email/toast)
6. Post-mortem (documento público, acciones correctivas)

---

## 12. Checklist de Seguridad (Pre-Production)

### Autenticación
- [ ] JWT con expiración corta (15 min)
- [ ] Refresh tokens rotativos
- [ ] Logout revoca refresh token
- [ ] Rate limiting en login

### Autorización
- [ ] Middleware RBAC en endpoints sensibles
- [ ] Validación de ownership
- [ ] Auditoría de acciones de admin

### Datos
- [ ] TLS 1.3 en producción
- [ ] Passwords con bcrypt (costo 12)
- [ ] Secrets en env (no hardcoded)
- [ ] Backup de DB cifrado

### Validación
- [ ] Validación en backend (nunca confiar en frontend)
- [ ] Sanitización de HTML
- [ ] Prepared statements (GORM)

### Headers
- [ ] HSTS
- [ ] X-Frame-Options: DENY
- [ ] X-Content-Type-Options: nosniff
- [ ] CSP (Content-Security-Policy)

### Dependencias
- [ ] npm audit sin vulnerabilidades HIGH
- [ ] go mod verify OK
- [ ] Dependabot activo

### Monitoring
- [ ] Alertas de login fallido
- [ ] Alertas de payment failure
- [ ] Logs con trace_id
- [ ] Dashboard de seguridad (Grafana)

---

**Actualizado:** 2025-11-10
**Próxima revisión:** Tras primer pentest
