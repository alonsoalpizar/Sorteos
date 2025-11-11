# Testing Quick Start Guide 🚀

**Sprint 5-6:** Reservas y Pagos con PayPal
**Tiempo estimado:** 30 minutos (manual) + 1 hora (API scripts)

---

## Opción 1: Testing Manual (Recomendado para empezar) ⚡

### Step 1: Levantar el entorno
```bash
cd /opt/Sorteos

# Levantar servicios
docker compose up -d

# Verificar que todo esté corriendo
docker compose ps

# Ver logs del API
docker compose logs -f api

# Verificar migraciones aplicadas
docker compose logs api | grep "migration"
```

### Step 2: Configurar PayPal Sandbox (5 min)

1. **Ir a:** https://developer.paypal.com/dashboard/
2. **Login** con tu cuenta PayPal
3. **Apps & Credentials** → Sandbox
4. **Create App:**
   - App Name: `Sorteos Testing`
   - App Type: Merchant
5. **Copiar credenciales:**
   - Client ID
   - Secret
6. **Actualizar `.env`:**
```bash
CONFIG_PAYMENT_PROVIDER=paypal
CONFIG_PAYMENT_CLIENT_ID=<tu_client_id_aquí>
CONFIG_PAYMENT_SECRET=<tu_secret_aquí>
CONFIG_PAYMENT_SANDBOX=true
```

7. **Reiniciar API:**
```bash
docker compose restart api
```

8. **Crear cuentas de prueba:**
   - Sandbox → Accounts → Create account
   - Crear 1 Business account (vendedor)
   - Crear 1 Personal account (comprador)

### Step 3: Ejecutar Checklist Manual (30 min)

Abrir y seguir: **[testing-manual-checklist.md](./testing-manual-checklist.md)**

**URLs:**
- Frontend: http://localhost:5173
- API: http://localhost:8080
- PayPal Sandbox: https://www.sandbox.paypal.com

**Flujo:**
1. Registrar usuario
2. Crear sorteo (draft → publish)
3. Seleccionar 3-5 números
4. Crear reserva (timer 5 min)
5. Pagar con PayPal sandbox
6. Verificar success page

---

## Opción 2: Testing de API con Scripts (1-2 horas) 🔧

### Step 1: Instalar herramientas
```bash
# Ubuntu/Debian
sudo apt install curl jq apache2-bench

# macOS
brew install curl jq

# Verificar instalación
curl --version
jq --version
```

### Step 2: Setup variables y helpers
```bash
# Copiar helpers
cat > ~/sorteos-test-helpers.sh << 'EOF'
export API_URL="http://localhost:8080/api/v1"
export TOKEN=""

alias pj='python3 -m json.tool'

post_auth() {
  curl -X POST "$API_URL$1" \
    -H "Authorization: Bearer $TOKEN" \
    -H "Content-Type: application/json" \
    -d "$2" -s | pj
}

get_auth() {
  curl -X GET "$API_URL$1" \
    -H "Authorization: Bearer $TOKEN" -s | pj
}
EOF

source ~/sorteos-test-helpers.sh
```

### Step 3: Ejecutar scripts de testing

Seguir: **[testing-api-scripts.md](./testing-api-scripts.md)**

**Test suites incluidos:**
1. ✅ Autenticación (register, login, me)
2. ✅ CRUD Sorteos (create, publish, list)
3. ✅ Reservas (create, conflict, idempotency)
4. ✅ Pagos (create intent, webhook simulation)
5. ✅ Concurrencia (10 requests simultáneas)
6. ✅ Validaciones y errores
7. ✅ Performance benchmarks

---

## Opción 3: Entorno de Test Aislado (Avanzado) 🐳

### Usar docker-compose.test.yml

**Ventaja:** Entorno separado del desarrollo, con datos independientes.

```bash
# Levantar entorno de test
docker compose -f docker-compose.test.yml up -d

# Verificar servicios
docker compose -f docker-compose.test.yml ps

# Ver logs
docker compose -f docker-compose.test.yml logs -f api-test

# Ejecutar migrations
docker compose -f docker-compose.test.yml exec api-test /app/main migrate up

# API Test: http://localhost:8081
# Frontend Test: http://localhost:5174
# Postgres Test: localhost:5433
# Redis Test: localhost:6380
```

**Testing de API contra entorno de test:**
```bash
export API_URL="http://localhost:8081/api/v1"
# ... ejecutar scripts normalmente
```

**Limpiar entorno:**
```bash
docker compose -f docker-compose.test.yml down -v
```

---

## Checklist Pre-Testing ✅

Antes de empezar, verificar:

- [ ] Docker y Docker Compose instalados
- [ ] Puerto 8080 (API) libre
- [ ] Puerto 5173 (Frontend) libre
- [ ] Cuenta PayPal Developer creada
- [ ] Credenciales sandbox configuradas en `.env`
- [ ] `docker compose up -d` ejecutado sin errores
- [ ] Migraciones aplicadas (ver logs)
- [ ] Frontend accesible en http://localhost:5173

---

## Troubleshooting 🔍

### API no levanta
```bash
# Ver logs completos
docker compose logs api

# Verificar migraciones
docker compose exec api /app/main migrate status

# Verificar conexión a DB
docker compose exec postgres psql -U sorteos -d sorteos_db -c '\dt'

# Verificar conexión a Redis
docker compose exec redis redis-cli ping
```

### Frontend no carga
```bash
# Ver logs
docker compose logs frontend

# Rebuild sin cache
docker compose build frontend --no-cache
docker compose up -d frontend

# Verificar variables de entorno
docker compose exec frontend env | grep VITE
```

### PayPal redirect no funciona
```bash
# Verificar configuración en .env
cat .env | grep CONFIG_PAYMENT

# Ver logs del API durante la creación del payment intent
docker compose logs -f api

# Verificar que SANDBOX=true
# Verificar que Client ID y Secret son de Sandbox (no Live)
```

### Números no se reservan (race condition)
```bash
# Verificar Redis está corriendo
docker compose exec redis redis-cli ping

# Ver locks activos
docker compose exec redis redis-cli KEYS "lock:*"

# Ver TTL de locks
docker compose exec redis redis-cli TTL "lock:raffle:UUID:number:0001"
```

---

## Resultados Esperados 🎯

### Testing Manual (30 min)
- ✅ **30/30 test cases** pasando
- ⏱️ **Performance:** Todas las páginas < 2s
- 🐛 **Bugs críticos:** 0
- 📊 **Cobertura:** Happy path completo

### Testing de API (1-2 horas)
- ✅ **30/30 endpoints** funcionando
- ⏱️ **Performance:** < 500ms (p95)
- 🔒 **Concurrencia:** 0 duplicados en 100 requests
- 🔑 **Idempotency:** Funciona correctamente

---

## Próximos Pasos

1. **Ejecutar Opción 1** (Testing Manual) → 30 min ⚡
2. **Documentar bugs** encontrados en GitHub Issues
3. **Resolver bugs críticos** (si existen)
4. **Re-ejecutar tests** después de fixes
5. **Ejecutar Opción 2** (Testing API) → 1-2 horas 🔧
6. **Actualizar roadmap** con resultados

---

## Documentación Completa

- **Estrategia General:** [testing-strategy.md](./testing-strategy.md)
- **Checklist Manual:** [testing-manual-checklist.md](./testing-manual-checklist.md)
- **Scripts de API:** [testing-api-scripts.md](./testing-api-scripts.md)
- **Roadmap:** [roadmap.md](./roadmap.md)

---

## Ayuda

**Dudas o problemas?**
1. Revisar logs: `docker compose logs -f api`
2. Verificar DB: `docker compose exec postgres psql -U sorteos -d sorteos_db`
3. Ver Redis: `docker compose exec redis redis-cli MONITOR`
4. Leer documentación: [testing-strategy.md](./testing-strategy.md)

**Listo para empezar?** ✨

```bash
# 🚀 Let's go!
cd /opt/Sorteos
docker compose up -d
open http://localhost:5173
```
