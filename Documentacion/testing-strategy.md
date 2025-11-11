# Estrategia de Testing - Plataforma de Sorteos

**Fecha:** 2025-11-11
**Sprint:** 5-6 (Reservas y Pagos)
**Objetivo:** Validar el flujo completo de compra con PayPal

---

## Niveles de Testing

### Nivel 1: Testing Manual (Prioritario) ⚡
**Duración estimada:** 30-45 minutos
**Herramientas:** Navegador, DevTools, PayPal Sandbox
**Documentación:** [testing-manual-checklist.md](./testing-manual-checklist.md)

**Casos de prueba:**
1. ✅ Registro y login de usuario
2. ✅ Creación de sorteo (draft → published)
3. ✅ Selección de números (1, 5, 10 números)
4. ✅ Creación de reserva (5 min timer)
5. ✅ Pago con PayPal sandbox (exitoso)
6. ✅ Pago cancelado (cancel flow)
7. ✅ Expiración de reserva (timeout)
8. ✅ Números ya vendidos (race condition manual)

---

### Nivel 2: Testing de API (Recomendado) 🔧
**Duración estimada:** 1-2 horas
**Herramientas:** cURL, httpie, o Postman
**Documentación:** [testing-api-scripts.md](./testing-api-scripts.md)

**Casos de prueba:**
1. POST /api/v1/auth/register → 201 Created
2. POST /api/v1/auth/login → 200 OK + tokens
3. POST /api/v1/raffles → 201 Created (draft)
4. PATCH /api/v1/raffles/{id}/publish → 200 OK
5. POST /api/v1/reservations → 201 Created (con distributed locks)
6. POST /api/v1/reservations (duplicado) → 409 Conflict
7. POST /api/v1/payments/intent → 201 Created + PayPal URL
8. GET /api/v1/reservations/me → 200 OK
9. GET /api/v1/payments/me → 200 OK
10. Webhook simulation → 200 OK

**Testing de concurrencia:**
```bash
# 10 requests simultáneas al mismo número
seq 10 | xargs -P10 -I {} curl -X POST http://localhost:8080/api/v1/reservations \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"raffle_id":"uuid","number_ids":["0001"],"session_id":"test-{}"}'

# Solo 1 debería tener éxito (201), el resto 409 Conflict
```

---

### Nivel 3: Testing E2E Automatizado (Avanzado) 🚀
**Duración estimada:** 4-6 horas setup + scripts
**Herramientas:** Playwright (recomendado) o Cypress
**Documentación:** [testing-e2e-playwright.md](./testing-e2e-playwright.md)

**Ventajas:**
- Regression testing automático
- CI/CD integration
- Screenshots y videos de failures
- Cobertura de happy path + edge cases

**Stack recomendado:**
```
frontend/e2e/
├── playwright.config.ts
├── fixtures/
│   ├── users.json
│   └── raffles.json
├── tests/
│   ├── auth.spec.ts
│   ├── raffle-creation.spec.ts
│   ├── number-selection.spec.ts
│   ├── checkout-flow.spec.ts
│   └── payment-flow.spec.ts
└── utils/
    ├── api-helpers.ts
    └── paypal-mock.ts
```

---

## Configuración del Entorno de Testing

### PayPal Sandbox Setup

1. **Crear cuenta sandbox** en https://developer.paypal.com
2. **Crear 2 cuentas de prueba:**
   - Business Account (vendedor)
   - Personal Account (comprador)
3. **Obtener credenciales:**
   - Client ID
   - Secret
4. **Configurar .env.test:**
```bash
CONFIG_PAYMENT_PROVIDER=paypal
CONFIG_PAYMENT_CLIENT_ID=<sandbox_client_id>
CONFIG_PAYMENT_SECRET=<sandbox_secret>
CONFIG_PAYMENT_SANDBOX=true
CONFIG_PAYMENT_SUCCESS_URL=http://localhost:5173/payment/success
CONFIG_PAYMENT_CANCEL_URL=http://localhost:5173/payment/cancel
```

### Docker Compose para Testing

```yaml
# docker-compose.test.yml
services:
  postgres-test:
    image: postgres:16-alpine
    environment:
      POSTGRES_DB: sorteos_test
      POSTGRES_USER: sorteos
      POSTGRES_PASSWORD: sorteos123
    ports:
      - "5433:5432"

  redis-test:
    image: redis:7-alpine
    ports:
      - "6380:6379"

  api-test:
    build: .
    environment:
      DB_HOST: postgres-test
      REDIS_HOST: redis-test
      CONFIG_PAYMENT_SANDBOX: "true"
    depends_on:
      - postgres-test
      - redis-test
    ports:
      - "8081:8080"
```

**Comandos:**
```bash
# Levantar entorno de test
docker compose -f docker-compose.test.yml up -d

# Correr migraciones
docker compose -f docker-compose.test.yml exec api-test ./main migrate up

# Ver logs
docker compose -f docker-compose.test.yml logs -f api-test

# Limpiar entorno
docker compose -f docker-compose.test.yml down -v
```

---

## Métricas de Éxito

### Funcionalidad
- [ ] 100% de casos happy path funcionan
- [ ] 100% de casos de error manejados correctamente
- [ ] 0 race conditions en reservas
- [ ] Timeout de 5 minutos se respeta
- [ ] PayPal redirect funciona correctamente

### Performance
- [ ] Crear reserva: < 500ms (p95)
- [ ] Crear payment intent: < 1s (p95)
- [ ] 100 usuarios concurrentes sin errores
- [ ] 500 requests simultáneas al mismo número: solo 1 éxito

### Seguridad
- [ ] Tokens JWT validan correctamente
- [ ] Webhooks verifican firma
- [ ] No se pueden reservar números de otros usuarios
- [ ] Idempotency keys previenen duplicados

---

## Recomendación de Orden

Para este sprint, te recomiendo:

1. **Hoy/mañana:** Nivel 1 (Testing Manual con checklist) - 30 min
   - Valida que el flujo básico funciona end-to-end
   - Identifica bugs críticos rápidamente

2. **Esta semana:** Nivel 2 (Testing de API) - 2 horas
   - Valida robustez del backend
   - Testing de concurrencia con script bash

3. **Sprint 7-8:** Nivel 3 (E2E automatizado con Playwright) - 6 horas
   - Una vez que el flujo es estable
   - Para regression testing en futuros sprints

---

## Próximos Pasos

1. ✅ Crear checklist de testing manual
2. ✅ Crear scripts de API testing
3. ⏳ Levantar entorno con PayPal sandbox
4. ⏳ Ejecutar testing manual (Nivel 1)
5. ⏳ Ejecutar testing de API (Nivel 2)
6. ⏳ Documentar bugs encontrados
7. ⏳ Iterar y resolver issues

---

## Recursos

- **PayPal Sandbox:** https://developer.paypal.com/dashboard/
- **Playwright Docs:** https://playwright.dev/
- **Docker Compose Testing:** https://docs.docker.com/compose/
- **Redis Locks Testing:** https://redis.io/docs/latest/develop/use/patterns/distributed-locks/
