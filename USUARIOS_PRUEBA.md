# Usuarios de Prueba - Plataforma de Sorteos

**Fecha de creación:** 2025-11-10
**Servidor:** http://62.171.188.255:8080

---

## 1. Usuario Normal (Verificado)

### Credenciales
- **Email:** `test@sorteos.com`
- **Password:** `TestPassword123!`
- **ID:** 1
- **UUID:** `2bbe7e86-4b71-4500-9f41-46ff135ced95`
- **Rol:** `user`
- **KYC Level:** `email_verified` ✅
- **Estado:** `active`

### Permisos
- ✅ Crear sorteos (máximo 10 activos)
- ✅ Listar sorteos públicos
- ✅ Ver detalle de sorteos
- ✅ Publicar sus propios sorteos
- ✅ Actualizar sus propios sorteos
- ✅ Eliminar sus propios sorteos (soft delete)
- ❌ Suspender sorteos (solo admin)

### Tokens de Acceso (Válidos por 15 minutos)
```
Access Token:
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJlbWFpbCI6InRlc3RAc29ydGVvcy5jb20iLCJyb2xlIjoidXNlciIsImt5Y19sZXZlbCI6ImVtYWlsX3ZlcmlmaWVkIiwiaXNzIjoic29ydGVvcy1wbGF0Zm9ybSIsInN1YiI6IjEiLCJleHAiOjE3NjI3NTk4NjksIm5iZiI6MTc2Mjc1ODk2OSwiaWF0IjoxNzYyNzU4OTY5fQ.EexyZ7d7cmPKUOhQ2EvH18kZt5NHgKQvxqdlbhxxXhY

Refresh Token:
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJlbWFpbCI6InRlc3RAc29ydGVvcy5jb20iLCJyb2xlIjoiIiwia3ljX2xldmVsIjoiIiwiaXNzIjoic29ydGVvcy1wbGF0Zm9ybSIsInN1YiI6IjEiLCJleHAiOjE3NjMzNjM3NjksIm5iZiI6MTc2Mjc1ODk2OSwiaWF0IjoxNzYyNzU4OTY5fQ.1sRRZEbhzeF8MjujcAsQ8Na_M9ZuWE1fbFLTidJJido
```

---

## 2. Usuario Administrador (Verificado)

### Credenciales
- **Email:** `admin@sorteos.com`
- **Password:** `Admin123456!`
- **ID:** 2
- **UUID:** `957f7e6c-7b69-462b-b32c-387939e99a0a`
- **Rol:** `admin` 👑
- **KYC Level:** `email_verified` ✅
- **Estado:** `active`

### Permisos
- ✅ Todos los permisos del usuario normal
- ✅ Suspender cualquier sorteo
- ✅ Actualizar cualquier sorteo
- ✅ Eliminar cualquier sorteo
- ✅ Ver/editar usuarios
- ✅ Acceso completo al backoffice

### Tokens de Acceso (Válidos por 15 minutos)
```
Access Token:
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoyLCJlbWFpbCI6ImFkbWluQHNvcnRlb3MuY29tIiwicm9sZSI6ImFkbWluIiwia3ljX2xldmVsIjoiZW1haWxfdmVyaWZpZWQiLCJpc3MiOiJzb3J0ZW9zLXBsYXRmb3JtIiwic3ViIjoiMiIsImV4cCI6MTc2Mjc1OTg3OSwibmJmIjoxNzYyNzU4OTc5LCJpYXQiOjE3NjI3NTg5Nzl9.rXmGBbQxHSomGeX-7uy8d1Wsc5kkgg__iQ7PjkhdrUc

Refresh Token:
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoyLCJlbWFpbCI6ImFkbWluQHNvcnRlb3MuY29tIiwicm9sZSI6IiIsImt5Y19sZXZlbCI6IiIsImlzcyI6InNvcnRlb3MtcGxhdGZvcm0iLCJzdWIiOiIyIiwiZXhwIjoxNzYzMzYzNzc5LCJuYmYiOjE3NjI3NTg5NzksImlhdCI6MTc2Mjc1ODk3OX0.UBVAVk2ujkVafUfZ8oNnAFKMdj8aA4uU14crby_XdRc
```

---

## 3. Ejemplos de Uso de la API

### 3.1 Login (Obtener nuevos tokens)

```bash
# Usuario normal
curl -X POST http://62.171.188.255:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@sorteos.com",
    "password": "TestPassword123!"
  }'

# Admin
curl -X POST http://62.171.188.255:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@sorteos.com",
    "password": "Admin123456!"
  }'
```

### 3.2 Listar Sorteos (Público - No requiere autenticación)

```bash
curl http://62.171.188.255:8080/api/v1/raffles
```

### 3.3 Crear Sorteo (Requiere autenticación + KYC)

```bash
# Primero hacer login para obtener el access_token
TOKEN="<access_token_aquí>"

curl -X POST http://62.171.188.255:8080/api/v1/raffles \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "title": "Sorteo iPhone 15 Pro",
    "description": "Sorteo de un iPhone 15 Pro Max 256GB nuevo en caja sellada. El sorteo se realizará mediante la Lotería Nacional de Costa Rica.",
    "price_per_number": 5000,
    "total_numbers": 100,
    "draw_date": "2025-12-25T18:00:00Z",
    "draw_method": "loteria_nacional_cr"
  }'
```

### 3.4 Ver Detalle de Sorteo

```bash
# Por ID
curl http://62.171.188.255:8080/api/v1/raffles/1

# Con números incluidos
curl "http://62.171.188.255:8080/api/v1/raffles/1?include_numbers=true"

# Con imágenes incluidas
curl "http://62.171.188.255:8080/api/v1/raffles/1?include_images=true"
```

### 3.5 Publicar Sorteo (Requiere autenticación)

```bash
TOKEN="<access_token_aquí>"

curl -X POST http://62.171.188.255:8080/api/v1/raffles/1/publish \
  -H "Authorization: Bearer $TOKEN"
```

### 3.6 Actualizar Sorteo (Solo owner o admin)

```bash
TOKEN="<access_token_aquí>"

curl -X PUT http://62.171.188.255:8080/api/v1/raffles/1 \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "title": "Sorteo iPhone 15 Pro Max - ACTUALIZADO",
    "description": "Nueva descripción actualizada",
    "draw_date": "2025-12-31T18:00:00Z"
  }'
```

### 3.7 Suspender Sorteo (Solo admin)

```bash
ADMIN_TOKEN="<admin_access_token_aquí>"

curl -X POST http://62.171.188.255:8080/api/v1/raffles/1/suspend \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -d '{
    "reason": "Sorteo suspendido por incumplimiento de términos de servicio"
  }'
```

### 3.8 Eliminar Sorteo (Solo si no tiene ventas)

```bash
TOKEN="<access_token_aquí>"

curl -X DELETE http://62.171.188.255:8080/api/v1/raffles/1 \
  -H "Authorization: Bearer $TOKEN"
```

### 3.9 Refresh Token (Renovar access token)

```bash
REFRESH_TOKEN="<refresh_token_aquí>"

curl -X POST http://62.171.188.255:8080/api/v1/auth/refresh \
  -H "Content-Type: application/json" \
  -d '{
    "refresh_token": "'"$REFRESH_TOKEN"'"
  }'
```

---

## 4. Endpoints Disponibles

### Autenticación
- `POST /api/v1/auth/register` - Registrar usuario
- `POST /api/v1/auth/login` - Iniciar sesión
- `POST /api/v1/auth/refresh` - Renovar token
- `POST /api/v1/auth/verify-email` - Verificar email

### Sorteos (Públicos)
- `GET /api/v1/raffles` - Listar sorteos (paginado, con filtros)
- `GET /api/v1/raffles/:id` - Detalle de sorteo

### Sorteos (Autenticados + KYC email_verified)
- `POST /api/v1/raffles` - Crear sorteo (rate limit: 10/hora)
- `PUT /api/v1/raffles/:id` - Actualizar sorteo (solo owner/admin)
- `POST /api/v1/raffles/:id/publish` - Publicar sorteo
- `DELETE /api/v1/raffles/:id` - Eliminar sorteo (solo owner/admin)

### Sorteos (Solo Admin)
- `POST /api/v1/raffles/:id/suspend` - Suspender sorteo

### Otros
- `GET /health` - Health check
- `GET /ready` - Readiness check

---

## 5. Códigos de Error Comunes

- `400` - Validación fallida (campos inválidos)
- `401` - No autorizado (token inválido o expirado)
- `403` - Prohibido (sin permisos suficientes)
- `404` - No encontrado
- `409` - Conflicto (email ya existe, etc.)
- `429` - Rate limit excedido
- `500` - Error interno del servidor

---

## 6. Notas Importantes

1. **Tokens expiran cada 15 minutos** - Usa el refresh token para renovarlos
2. **Rate Limiting:**
   - Registro: 5 intentos/minuto por IP
   - Login: 5 intentos/minuto por IP
   - Crear sorteo: 10 sorteos/hora por usuario
3. **KYC Levels:**
   - `none` - No puede crear sorteos
   - `email_verified` - Puede crear sorteos ✅
   - `phone_verified` - Límites más altos (futuro)
   - `id_verified` - Sin límites (futuro)
4. **Estados de Sorteo:**
   - `draft` - Borrador (editable)
   - `active` - Publicado y activo
   - `suspended` - Suspendido por admin
   - `completed` - Finalizado con ganador
   - `cancelled` - Cancelado

---

## 7. Troubleshooting

### Token expirado
Si recibes error 401, haz login de nuevo o usa el refresh token:
```bash
curl -X POST http://62.171.188.255:8080/api/v1/auth/refresh \
  -H "Content-Type: application/json" \
  -d '{"refresh_token": "tu_refresh_token"}'
```

### Verificar estado de usuario
```sql
docker exec sorteos-postgres psql -U sorteos_user -d sorteos_db \
  -c "SELECT id, email, role, kyc_level, status, email_verified FROM users;"
```

### Ver logs del backend
```bash
docker logs -f sorteos-api
```

---

**Actualizado:** 2025-11-10 07:16
**Versión API:** v1
**Base URL:** http://62.171.188.255:8080
