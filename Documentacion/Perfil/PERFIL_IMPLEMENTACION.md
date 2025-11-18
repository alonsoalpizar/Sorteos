# Sistema de Perfil de Usuario - Plan de Implementación

## 1. Contexto y Concepto

### ¿Por qué necesitamos un perfil de usuario?

El sistema de Perfil de Usuario es fundamental para completar el ecosistema de Sorteos, permitiendo:

1. **Transparencia y Confianza**: Los organizadores pueden mostrar su identidad verificada, aumentando la confianza de los participantes
2. **Gestión Financiera**: Configurar cuentas IBAN para recibir ganancias de sorteos exitosos
3. **Verificación KYC**: Cumplir con regulaciones al verificar identidad mediante cédula y documentación legal
4. **Motivación del Organizador**: Ver ganancias estimadas en tiempo real incentiva a crear más sorteos
5. **Administración Centralizada**: Un solo lugar para gestionar toda la información personal, financiera y legal

### Flujo del Usuario

```
Usuario Nuevo → Email Verified → Completa Perfil → Sube Cédula → Configura IBAN → Full KYC
                                  (Foto, Datos)    (Ambos lados)   (Para retiros)   (Puede retirar)
```

### Relación con el Sistema de Wallet

- **Balance Available** (₡): Créditos recargados para comprar tickets (NO retirable)
- **Earnings Balance** (₡): Ganancias de sorteos (RETIRABLE vía IBAN)

El perfil conecta la wallet con el mundo real permitiendo liquidaciones a cuentas bancarias del organizador.

---

## 2. Estado Actual

### ✅ Ya Tenemos (80% del trabajo)

- [x] Modelo User con campos de dirección, teléfono, cédula
- [x] Sistema de KYC levels (none → email_verified → phone_verified → cedula_verified → full_kyc)
- [x] Wallet con EarningsBalance separado de BalanceAvailable
- [x] Sistema de carga de imágenes (adaptable para fotos de perfil)
- [x] Tab de "Ganancias" mostrando sorteos activos con desglose
- [x] Infraestructura de autenticación y autorización

### ❌ Nos Falta (20% del trabajo)

- [ ] Campo `profile_photo_url` en tabla users
- [ ] Campo `date_of_birth` en tabla users
- [ ] Campo `iban` (encriptado) en tabla users
- [ ] Tabla `kyc_documents` para almacenar cédula frente/dorso y selfie
- [ ] UseCases y Handlers para gestión de perfil
- [ ] Página de Perfil en frontend con formularios
- [ ] Componente de carga de foto con crop
- [ ] Calculadora de liquidación (ganancia - comisión plataforma - comisión bancaria)

---

## 3. Roadmap de Implementación

### Phase 1: Base de Datos (2-3 horas)
**Objetivo**: Agregar campos necesarios y crear tabla de documentos KYC

- [ ] **Migración 1**: Agregar campos a `users`
  ```sql
  ALTER TABLE users ADD COLUMN profile_photo_url VARCHAR(255);
  ALTER TABLE users ADD COLUMN date_of_birth DATE;
  ALTER TABLE users ADD COLUMN iban VARCHAR(255); -- Encriptado en app layer
  ```

- [ ] **Migración 2**: Crear tabla `kyc_documents`
  ```sql
  CREATE TABLE kyc_documents (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id),
    document_type VARCHAR(50) NOT NULL, -- 'cedula_front', 'cedula_back', 'selfie'
    file_url VARCHAR(255) NOT NULL,
    verification_status VARCHAR(50) DEFAULT 'pending', -- pending, approved, rejected
    uploaded_at TIMESTAMP DEFAULT NOW(),
    verified_at TIMESTAMP,
    rejected_reason TEXT,
    UNIQUE(user_id, document_type)
  );
  ```

- [ ] Ejecutar migraciones con `make migrate-up`

---

### Phase 2: Backend (6-8 horas)
**Objetivo**: Crear lógica de negocio para gestión de perfil

#### 2.1 Domain Layer

- [ ] Actualizar `internal/domain/user.go`:
  - Agregar campos `ProfilePhotoURL`, `DateOfBirth`, `IBAN`
  - Método `CanWithdraw() bool` (verifica full_kyc + IBAN configurado)

- [ ] Crear `internal/domain/kyc_document.go`:
  - Struct `KYCDocument` con campos completos
  - Enum `DocumentType` (CedulaFront, CedulaBack, Selfie)
  - Enum `VerificationStatus` (Pending, Approved, Rejected)

#### 2.2 Repository Layer

- [ ] Extender `internal/adapters/db/user_repository.go`:
  - `UpdateProfile(ctx, userID, profileData)`
  - `UpdateProfilePhoto(ctx, userID, photoURL)`
  - `UpdateIBAN(ctx, userID, encryptedIBAN)`
  - `GetUserProfile(ctx, userID)` - retorna user completo

- [ ] Crear `internal/adapters/db/kyc_document_repository.go`:
  - `UploadDocument(ctx, userID, docType, fileURL)`
  - `GetDocuments(ctx, userID)`
  - `UpdateVerificationStatus(ctx, docID, status, reason)`

#### 2.3 UseCase Layer

- [ ] Crear `internal/usecase/profile/update_profile.go`:
  - Validar campos (fecha de nacimiento razonable, formato de dirección)
  - Actualizar perfil del usuario

- [ ] Crear `internal/usecase/profile/upload_photo.go`:
  - Reutilizar infraestructura de `image/` para subir foto
  - Actualizar campo `profile_photo_url`

- [ ] Crear `internal/usecase/profile/configure_iban.go`:
  - Validar formato IBAN costarricense (CR + 22 dígitos)
  - Encriptar IBAN antes de guardar
  - Requerir `kyc_level >= cedula_verified`

- [ ] Crear `internal/usecase/profile/upload_kyc_document.go`:
  - Subir imagen del documento
  - Guardar en tabla `kyc_documents`
  - Si se suben todos (frente+dorso+selfie), actualizar `kyc_level → full_kyc`

#### 2.4 Handler Layer

- [ ] Crear `internal/adapters/http/handler/profile_handler.go`:
  - `GET /api/v1/profile` - Obtener perfil completo
  - `PUT /api/v1/profile` - Actualizar datos personales
  - `POST /api/v1/profile/photo` - Subir foto de perfil
  - `POST /api/v1/profile/iban` - Configurar IBAN
  - `POST /api/v1/profile/kyc/:document_type` - Subir documento KYC
  - `GET /api/v1/profile/kyc` - Ver documentos KYC cargados

- [ ] Registrar rutas en router con middleware de autenticación

---

### Phase 3: Frontend (8-10 horas)
**Objetivo**: Crear interfaz de usuario para gestión de perfil

#### 3.1 Estructura de Archivos

```
frontend/src/features/profile/
├── components/
│   ├── ProfilePage.tsx           # Página principal con tabs
│   ├── PersonalInfoForm.tsx      # Nombre, dirección, teléfono, fecha nacimiento
│   ├── ProfilePhotoUpload.tsx    # Upload con crop de imagen
│   ├── KYCDocuments.tsx          # Carga de cédula (frente/dorso) y selfie
│   ├── BankAccountForm.tsx       # Configuración de IBAN
│   ├── WalletSummary.tsx         # Resumen de wallet con link
│   └── SettlementCalculator.tsx  # Calculadora de liquidación
├── hooks/
│   ├── useProfile.ts             # React Query para GET /profile
│   ├── useUpdateProfile.ts       # Mutation para PUT /profile
│   ├── useUploadPhoto.ts         # Mutation para POST /profile/photo
│   ├── useConfigureIBAN.ts       # Mutation para POST /profile/iban
│   └── useUploadKYCDocument.ts   # Mutation para POST /profile/kyc/:type
└── types/
    └── profile.ts                # TypeScript interfaces
```

#### 3.2 API Client

- [ ] Crear `src/api/profile.ts`:
  - `getProfile()`
  - `updateProfile(data)`
  - `uploadProfilePhoto(file)`
  - `configureIBAN(iban)`
  - `uploadKYCDocument(documentType, file)`
  - `getKYCDocuments()`

#### 3.3 Componentes Principales

- [ ] **ProfilePage.tsx**:
  - Layout con tabs: "Información Personal", "Documentación KYC", "Cuenta Bancaria", "Resumen Wallet"
  - Mostrar indicador de nivel KYC actual
  - Botones de acción según estado (ej: si no tiene IBAN, mostrar alerta)

- [ ] **PersonalInfoForm.tsx**:
  - Campos: Nombre, Apellidos, Fecha de Nacimiento, Teléfono
  - Dirección: Provincia, Cantón, Distrito, Detalles
  - Validación de formulario con React Hook Form
  - Submit actualiza perfil

- [ ] **ProfilePhotoUpload.tsx**:
  - Drag & drop o click para subir
  - Preview con crop circular
  - Librería: `react-image-crop` o `react-easy-crop`
  - Comprimir imagen antes de subir

- [ ] **KYCDocuments.tsx**:
  - Tres cards: Cédula Frente, Cédula Dorso, Selfie
  - Cada card muestra estado: No cargado | Pendiente | Aprobado | Rechazado
  - Upload individual con preview
  - Mostrar razón de rechazo si aplica

- [ ] **BankAccountForm.tsx**:
  - Input para IBAN con formato CR12-3456-7890-1234-5678-90
  - Validación de formato IBAN costarricense
  - Requerir `kyc_level >= cedula_verified`
  - Mostrar alerta si no cumple requisitos

- [ ] **WalletSummary.tsx**:
  - Cards mostrando Balance Available y Earnings Balance
  - Link a `/wallet` para ver detalles completos
  - Badge mostrando si puede retirar (full_kyc + IBAN configurado)

- [ ] **SettlementCalculator.tsx**:
  - Input: Monto a retirar (máximo: earnings_balance)
  - Cálculos:
    ```
    Monto a retirar:        ₡100,000.00
    - Comisión Plataforma:  ₡      0.00  (0% en retiros)
    - Comisión Bancaria:    ₡  1,500.00  (1.5% mínimo ₡500)
    = Recibirás:            ₡ 98,500.00
    ```
  - Botón "Solicitar Retiro" (si cumple requisitos)

#### 3.4 Routing

- [ ] Agregar ruta en `src/App.tsx`:
  ```tsx
  <Route path="/profile" element={<ProtectedRoute><ProfilePage /></ProtectedRoute>} />
  ```

- [ ] Agregar link en navbar/sidebar:
  ```tsx
  <NavLink to="/profile">Mi Perfil</NavLink>
  ```

---

### Phase 4: Integración y Testing (4-6 horas)
**Objetivo**: Verificar funcionamiento end-to-end

- [ ] **Testing Backend**:
  - Probar endpoints con Postman/curl
  - Verificar validaciones (IBAN inválido, fecha nacimiento futura, etc.)
  - Probar actualización de `kyc_level` al subir documentos completos
  - Verificar encriptación de IBAN en DB
*******Único ajuste sugerido:
 Testing, cuando pruebes la integración, podrás usar los endpoints admin que implementé para:
Aprobar/rechazar documentos KYC desde el panel admin
Ver el perfil completo del usuario desde perspectiva admin
Actualizar KYC level manualmente si es necesario

- [ ] **Testing Frontend**:
  - Flujo completo: Usuario nuevo → Completar perfil → Subir documentos → Configurar IBAN
  - Validar que campos required funcionen
  - Probar crop de foto de perfil
  - Verificar que calculadora muestre montos correctos

- [ ] **Integración Wallet**:
  - Desde tab "Ganancias" en Wallet, mostrar link a configurar IBAN si no lo tiene
  - En ProfilePage, mostrar earnings_balance actual
  - Verificar que no pueda retirar si no tiene full_kyc + IBAN

- [ ] **UX Polish**:
  - Mensajes de éxito/error claros
  - Loading spinners durante uploads
  - Tooltips explicando cada campo
  - Responsive design (mobile-friendly)

- [ ] **Deploy**:
  - `cd /opt/Sorteos/backend && make build && sudo systemctl restart sorteos-backend`
  - `cd /opt/Sorteos/frontend && npm run build && sudo cp -r dist/* /var/www/sorteos.club/`
  - Verificar en producción

---

## 4. Decisiones Técnicas

### Almacenamiento de IBAN
- **Decisión**: Encriptar IBAN en application layer antes de guardar en DB
- **Razón**: Datos financieros sensibles, cumplir con mejores prácticas de seguridad
- **Implementación**: Usar librería Go `crypto/aes` con key en variable de entorno

### Tabla kyc_documents vs Campos en users
- **Decisión**: Tabla separada `kyc_documents`
- **Razón**: Permite almacenar múltiples versiones si se rechazan documentos, historial de verificación
- **Alternativa descartada**: Campos `cedula_front_url`, `cedula_back_url` en users (menos flexible)

### Foto de Perfil
- **Decisión**: Reutilizar infraestructura de `image/` existente para raffles
- **Razón**: No reinventar la rueda, ya maneja upload, storage, validaciones
- **Path**: `/uploads/profile-photos/{user_id}/{timestamp}.jpg`

### Validación de IBAN
- **Formato CR**: `CR` + 22 dígitos (total 24 caracteres)
- **Validación Frontend**: Regex pattern
- **Validación Backend**: Algoritmo MOD-97 para checksum IBAN
- **Ejemplo válido**: `CR12345678901234567890`

### Comisiones de Retiro
- **Plataforma**: 0% (ya se cobró 10% al finalizar sorteo)
- **Bancaria**: 1.5% con mínimo ₡500, máximo ₡5,000
- **Monto mínimo retiro**: ₡5,000

### Niveles KYC Requeridos
| Acción                     | KYC Level Mínimo   |
|----------------------------|--------------------|
| Comprar tickets            | email_verified     |
| Crear sorteos              | email_verified     |
| Ver ganancias estimadas    | email_verified     |
| Configurar IBAN            | cedula_verified    |
| Retirar ganancias          | full_kyc           |

---

## 5. Seguimiento de Progreso

### Checklist General

#### Backend
- [ ] Migraciones ejecutadas
- [ ] Domain models actualizados
- [ ] Repositories implementados
- [ ] UseCases implementados
- [ ] Handlers implementados
- [ ] Rutas registradas
- [ ] Tests unitarios pasando

#### Frontend
- [ ] Estructura de archivos creada
- [ ] API client implementado
- [ ] Hooks de React Query listos
- [ ] Componentes implementados
- [ ] Routing configurado
- [ ] Build exitoso sin errores

#### Integración
- [ ] Endpoints funcionando en Postman
- [ ] Frontend conectando correctamente
- [ ] Flujo completo funcional
- [ ] Deploy a producción exitoso
- [ ] Testing en producción OK

---

## 6. Estimación de Tiempo

| Fase               | Horas Estimadas | Horas Reales |
|--------------------|-----------------|--------------|
| Phase 1: Database  | 2-3 hrs         | _____        |
| Phase 2: Backend   | 6-8 hrs         | _____        |
| Phase 3: Frontend  | 8-10 hrs        | _____        |
| Phase 4: Testing   | 4-6 hrs         | _____        |
| **TOTAL**          | **20-27 hrs**   | **_____**    |

---

## 7. Próximos Pasos

1. Revisar y aprobar este plan
2. Comenzar con Phase 1 (Database migrations)
3. Probar migraciones en desarrollo
4. Continuar con Phase 2 (Backend)
5. Iterar según feedback

---

**Última actualización**: 2025-11-18
**Estado**: Plan inicial - Pendiente de aprobación
