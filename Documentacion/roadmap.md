# Roadmap de Desarrollo - Plataforma de Sorteos

**Versión:** 1.0
**Fecha:** 2025-11-10
**Metodología:** Sprints de 2 semanas (Scrum adaptado)

---

## 1. Visión General

Este roadmap define las **3 fases principales** del proyecto, desde el MVP hasta la plataforma completa con aplicaciones móviles nativas. Cada fase incluye hitos medibles, criterios de aceptación y estimaciones realistas.

**Horizonte temporal:**
- **Fase 1 (MVP):** 8-10 semanas
- **Fase 2 (Escalamiento):** 10-12 semanas
- **Fase 3 (Expansión):** 12-16 semanas

---

## 2. Fase 1 - MVP (Producto Mínimo Viable)

**Objetivo:** Lanzar plataforma funcional con un único proveedor de pagos y funcionalidades core.

**Duración estimada:** 8-10 semanas (4-5 sprints)

---

### Sprint 1-2: Infraestructura y Autenticación

#### Tareas Backend
- [ ] Setup proyecto Go con estructura hexagonal
- [ ] Configuración Docker Compose (Postgres, Redis)
- [ ] Migraciones base (users, roles, audit_logs)
- [ ] Sistema de autenticación:
  - Registro con email/teléfono
  - Login con JWT (access + refresh tokens)
  - Middleware de autorización por roles
  - Endpoint de verificación email/SMS (integración Twilio/SendGrid)
- [ ] Logging estructurado con Zap
- [ ] Configuración Viper con .env
- [ ] Rate limiting con Redis

#### Tareas Frontend
- [ ] Setup proyecto Vite + React + TypeScript
- [ ] Configuración Tailwind + shadcn/ui
- [ ] Componentes base (Button, Input, Card, Layout)
- [ ] Páginas: Register, Login, Verify
- [ ] React Query setup con Axios
- [ ] Zustand store para autenticación
- [ ] Protected routes

#### Entregables
- Usuario puede registrarse, verificar cuenta y hacer login
- Tokens JWT funcionales con refresh automático
- Dark mode funcional

---

### Sprint 3-4: Gestión de Sorteos (CRUD Básico)

#### Tareas Backend
- [ ] Migraciones: raffles, raffle_numbers, raffle_images
- [ ] Repositorios GORM para sorteos
- [ ] Casos de uso:
  - CreateRaffle (con validaciones)
  - ListRaffles (paginación, filtros por estado)
  - GetRaffleDetail (con números disponibles)
  - UpdateRaffle (solo owner o admin)
  - SuspendRaffle (admin only)
- [ ] Generación automática de rango de números (00-99 configurable)
- [ ] Upload de imágenes (S3 o local storage)
- [ ] Cache Redis de sorteos activos

#### Tareas Frontend
- [ ] Páginas:
  - Listado de sorteos (grid con filtros)
  - Detalle de sorteo (info, galería, números disponibles)
  - Crear/editar sorteo (formulario multi-step)
- [ ] Componentes:
  - RaffleCard (preview)
  - NumberGrid (visualización 00-99 con estados)
  - ImageUploader
- [ ] Validación con react-hook-form + zod

#### Entregables
- Usuario puede publicar sorteo con detalles completos
- Vista pública de sorteos activos
- Administrador puede suspender sorteos

---

### Sprint 5-6: Reservas y Pagos

#### Tareas Backend
- [ ] Migraciones: reservations, payments, idempotency_keys
- [ ] Sistema de reserva temporal:
  - Lock distribuido Redis por número
  - Crear reserva (status=pending, expires_at=now+5min)
  - Cron job para liberar reservas expiradas
- [ ] Integración con PSP (Stripe como primera opción):
  - Interfaz PaymentProvider
  - Implementación StripeProvider
  - Manejo de webhooks (payment.succeeded, payment.failed)
  - Idempotencia con Idempotency-Key
- [ ] Flujo completo:
  1. POST /raffles/{id}/reservations → crea reserva + lock
  2. POST /payments → intenta cargo con Stripe
  3. Webhook confirma → marca números como sold
  4. Si falla/expira → libera números
- [ ] Tests de concurrencia (vegeta/k6)

#### Tareas Frontend
- [ ] Página de checkout:
  - Selección de números (click en NumberGrid)
  - Carrito temporal (Zustand)
  - Formulario de pago (Stripe Elements)
  - Pantalla de confirmación
- [ ] Componentes:
  - NumberSelector (multi-selección)
  - PaymentForm (iframe Stripe o tarjeta directa)
  - OrderSummary (precio, fees, total)
- [ ] Manejo de estados:
  - Reserva pendiente (timer 5 min)
  - Pago procesando (spinner)
  - Pago exitoso (confetti + redirect)
  - Pago fallido (reintentar)

#### Entregables
- Usuario puede reservar números y pagar con tarjeta
- Números no se duplican (prueba con 500 req concurrentes)
- Reservas expiradas se liberan automáticamente
- Webhooks procesan pagos correctamente

---

### Sprint 7-8: Selección de Ganador y Backoffice Mínimo

#### Tareas Backend
- [ ] Sistema de selección de ganador:
  - Integración con API Lotería Nacional (o mock)
  - Cron job que consulta resultados en draw_date
  - Marca ganadores en raffle_numbers
  - Notificación por email/SMS al ganador
- [ ] Endpoints backoffice:
  - GET /admin/raffles (listado completo con filtros)
  - PATCH /admin/raffles/{id} (suspender/activar)
  - GET /admin/users (con filtros KYC)
  - POST /admin/settlements (crear liquidación manual)
- [ ] Audit log para todas las acciones de admin

#### Tareas Frontend
- [ ] Panel de usuario (dashboard):
  - Mis sorteos publicados (estados, % vendido)
  - Sorteos en los que participé
  - Sorteos ganados
  - Historial de pagos
- [ ] Panel de admin (backoffice básico):
  - Listado de sorteos con acciones (suspender/activar)
  - Listado de usuarios (verificar/suspender)
  - Vista de liquidaciones pendientes
- [ ] Componentes:
  - DataTable reutilizable (con sorting, paginación)
  - StatusBadge (draft/active/suspended/completed)
  - ActionMenu (suspender, editar, ver detalles)

#### Entregables
- Ganadores se determinan automáticamente según lotería
- Usuario recibe notificación al ganar
- Admin puede gestionar sorteos y usuarios desde backoffice
- Todas las acciones de admin quedan registradas (audit log)

---

### Sprint 9-10: Testing, Optimización y Lanzamiento MVP

#### Tareas
- [ ] Tests de aceptación:
  - Flujo completo end-to-end (Playwright/Cypress)
  - Pruebas de carga (k6): 1000 usuarios concurrentes
  - Pruebas de seguridad (OWASP ZAP)
- [ ] Optimizaciones:
  - Índices de base de datos (EXPLAIN ANALYZE)
  - Lazy loading de imágenes
  - Code splitting en React
  - CDN para assets estáticos
- [ ] Documentación:
  - README con setup instructions
  - API docs (Swagger/OpenAPI)
  - Guía de usuario (screenshots)
- [ ] Deploy a staging:
  - CI/CD pipeline completo
  - Health checks y rollback automático
  - Monitoreo con Prometheus + Grafana
- [ ] Beta testing con 50 usuarios reales
- [ ] Corrección de bugs críticos

#### Entregables
- MVP en producción con dominio custom
- Métricas de rendimiento (p95 < 500ms)
- Documentación completa para usuarios y desarrolladores

---

## 3. Fase 2 - Escalamiento y Funcionalidades Avanzadas

**Objetivo:** Expandir capacidades de la plataforma y preparar para crecimiento.

**Duración estimada:** 10-12 semanas (5-6 sprints)

---

### Sprint 11-12: Múltiples PSPs y Modo "Sin Cobro"

#### Backend
- [ ] Implementar providers adicionales:
  - PayPalProvider
  - LocalCRProvider (procesador de CR por definir)
- [ ] Sistema de routing de pagos:
  - Feature flags por sorteo (Stripe/PayPal/Local)
  - Fallback automático si PSP falla
- [ ] Modo "sin cobro en plataforma":
  - Sorteos gratuitos (owner coordina pago fuera)
  - Solo cobro de suscripción mensual al owner
  - Modelo de suscripción (Stripe Billing)

#### Frontend
- [ ] Selector de método de pago en checkout
- [ ] Modal de suscripción (planes Basic/Pro)
- [ ] Dashboard de owner con estado de suscripción

#### Entregables
- Usuario puede pagar con Stripe, PayPal o método local
- Owners pueden publicar sorteos sin cobro + pagar suscripción

---

### Sprint 13-14: Búsqueda Avanzada y Sistema de Afiliados

#### Backend
- [ ] Full-text search con PostgreSQL (pg_trgm):
  - Búsqueda por título, descripción, categoría
  - Filtros combinados (precio, fecha, % vendido)
  - Ordenamiento por relevancia
- [ ] Sistema de afiliados:
  - Tabla affiliate_links (user_id, code, clicks, conversions)
  - Endpoint para generar link único
  - Tracking de registros por afiliado
  - Cálculo de comisiones

#### Frontend
- [ ] Barra de búsqueda con autocomplete
- [ ] Filtros avanzados (sidebar)
- [ ] Panel de afiliados (generar link, estadísticas)

#### Entregables
- Búsqueda rápida y precisa de sorteos
- Usuarios pueden generar links de afiliado y ganar comisiones

---

### Sprint 15-16: Multilenguaje y Comunicación entre Usuarios

#### Backend
- [ ] i18n en backend (mensajes de error, emails)
- [ ] Sistema de mensajería privada:
  - Tabla messages (sender_id, receiver_id, content, read_at)
  - Notificaciones en tiempo real (WebSockets)

#### Frontend
- [ ] Selector de idioma (Español/Inglés)
- [ ] Inbox de mensajes (estilo chat)
- [ ] Notificaciones en tiempo real (toast)

#### Entregables
- Plataforma disponible en ES/EN
- Usuarios pueden comunicarse vía mensajes privados

---

### Sprint 17-18: Comentarios, Valoraciones e Integración con Redes Sociales

#### Backend
- [ ] Sistema de reviews:
  - Tabla reviews (raffle_id, user_id, rating, comment)
  - Moderación (admin puede ocultar reviews)
- [ ] Open Graph tags dinámicos (meta tags para compartir)

#### Frontend
- [ ] Sección de comentarios en detalle de sorteo
- [ ] Botones de compartir (Facebook, Twitter, WhatsApp)
- [ ] Modal de valoración post-sorteo

#### Entregables
- Usuarios pueden comentar y valorar sorteos
- Compartir en redes sociales genera preview atractivo

---

### Sprint 19-20: Notificaciones en Tiempo Real y Dashboards Avanzados

#### Backend
- [ ] WebSockets para eventos en vivo:
  - Nuevo sorteo publicado
  - Sorteo próximo a cerrarse
  - Ganador anunciado
- [ ] Vistas materializadas para KPIs:
  - Total vendido por sorteo/usuario/período
  - Tasa de conversión reserva → pago
  - Top sorteos por ingresos

#### Frontend
- [ ] Dashboard de owner con gráficos (Chart.js):
  - Ingresos por mes
  - % de vendido por sorteo
  - Tasa de conversión
- [ ] Notificaciones push (PWA)

#### Entregables
- Notificaciones en tiempo real funcionales
- Dashboards con métricas accionables para owners

---

### Sprint 21-22: Optimización y Preparación para Escala

#### Tareas
- [ ] Caching agresivo:
  - CDN para imágenes (CloudFront/Cloudflare)
  - Cache de listados en Redis (invalidación inteligente)
- [ ] Database tuning:
  - Índices compuestos optimizados
  - Particionamiento de tablas grandes (audit_logs)
- [ ] Horizontal scaling:
  - Balanceador de carga (Nginx/HAProxy)
  - Réplicas de lectura en Postgres
- [ ] Pruebas de carga: 10k usuarios concurrentes

#### Entregables
- Plataforma soporta 10k usuarios simultáneos
- Latencia p95 < 300ms en operaciones críticas

---

## 4. Fase 3 - Expansión y Aplicaciones Móviles

**Objetivo:** Alcance global y experiencia móvil nativa.

**Duración estimada:** 12-16 semanas (6-8 sprints)

---

### Sprint 23-26: Aplicación Móvil (React Native)

#### Tareas
- [ ] Setup React Native con TypeScript
- [ ] Compartir lógica con web (custom hooks)
- [ ] Pantallas principales:
  - Login/Register
  - Listado y detalle de sorteos
  - Checkout con Apple Pay / Google Pay
  - Dashboard de usuario
- [ ] Push notifications (FCM)
- [ ] Deep links (abrir sorteo desde notificación)
- [ ] Beta en TestFlight / Google Play Beta

#### Entregables
- Apps nativas iOS + Android en beta pública
- Notificaciones push funcionales

---

### Sprint 27-30: Sorteos Temáticos y Campañas Automatizadas

#### Backend
- [ ] Taxonomía de categorías (Viajes, Tecnología, Moda, etc.)
- [ ] Sistema de tags y recomendaciones
- [ ] Integración con herramienta de marketing automation (HubSpot/Mailchimp):
  - Campañas por email basadas en comportamiento
  - Segmentación de usuarios

#### Frontend
- [ ] Landing pages por categoría
- [ ] Recomendaciones personalizadas
- [ ] Builder de campañas (admin)

#### Entregables
- Sorteos organizados por temas
- Campañas automatizadas de email marketing

---

### Sprint 31-34: Analytics Avanzado y A/B Testing

#### Backend
- [ ] Integración con Google Analytics 4
- [ ] Events tracking personalizado
- [ ] Sistema de feature flags (LaunchDarkly/Unleash)

#### Frontend
- [ ] Dashboards de analytics para owners
- [ ] A/B testing en páginas clave (checkout, landing)

#### Entregables
- Análisis detallado de comportamiento de usuarios
- Optimización basada en datos (A/B tests)

---

### Sprint 35-38: Programa de Fidelización y Gamificación

#### Backend
- [ ] Sistema de puntos y niveles:
  - Puntos por compra, referido, compartir
  - Niveles (Bronce, Plata, Oro)
  - Recompensas (descuentos, boletos gratis)
- [ ] Tabla de logros (achievements)

#### Frontend
- [ ] Perfil con badges y nivel actual
- [ ] Marketplace de recompensas
- [ ] Animaciones de logros desbloqueados

#### Entregables
- Sistema de fidelización activo
- Incremento en retención de usuarios (meta: +20%)

---

## 5. Hitos Críticos

| Hito | Fecha Estimada | Criterio de Éxito |
|------|----------------|-------------------|
| MVP Lanzado | Semana 10 | 100 sorteos publicados, 500 usuarios registrados |
| 1 PSP Adicional | Semana 14 | 30% de pagos con PSP alternativo |
| App Móvil Beta | Semana 26 | 1000 descargas en beta |
| 10k Usuarios Activos | Semana 32 | 10k MAU con < 300ms p95 latency |

---

## 6. Riesgos y Mitigaciones

| Riesgo | Probabilidad | Impacto | Mitigación |
|--------|--------------|---------|------------|
| Integración PSP falla | Media | Alto | Mock provider para tests, fallback automático |
| Doble venta de números | Baja | Crítico | Tests de concurrencia en CI, locks distribuidos |
| Escalado de DB | Media | Alto | Réplicas de lectura, caché agresivo |
| Retraso en app móvil | Alta | Medio | Priorizar web, liberar móvil en Fase 3.5 si necesario |

---

## 7. Recursos Necesarios

### Equipo Mínimo (Fase 1)
- 1 Backend Developer (Go)
- 1 Frontend Developer (React)
- 1 Full-Stack Developer (Go + React)
- 1 DevOps (part-time)
- 1 QA (part-time)

### Equipo Fase 2-3
- +1 Backend Developer
- +1 Mobile Developer (React Native)
- +1 UX/UI Designer
- DevOps full-time

---

## 8. Presupuesto Estimado (Infraestructura)

**Fase 1 (MVP):**
- AWS/DigitalOcean: $100-200/mes
- Stripe fees: 2.9% + $0.30 por transacción
- SendGrid: $15/mes (40k emails)
- Twilio: ~$0.01/SMS

**Fase 2:**
- Infra: $300-500/mes (réplicas, CDN)
- Multiple PSPs: fees variables

**Fase 3:**
- Infra: $800-1200/mes (app móvil, analytics)

---

## 9. Métricas de Éxito por Fase

**Fase 1 (MVP):**
- 500 usuarios registrados
- 100 sorteos publicados
- 70% tasa de conversión reserva → pago
- 0 incidentes de doble venta

**Fase 2:**
- 5000 usuarios activos mensuales (MAU)
- 3 PSPs integrados
- NPS > 40

**Fase 3:**
- 20k MAU
- Apps móviles con 4.5+ estrellas
- 80% retención mensual

---

## 10. Dependencias Externas

- **API Lotería Nacional de Costa Rica:** Confirmación de disponibilidad y documentación
- **PSP Local (CR):** Identificar y firmar contrato antes de Sprint 11
- **Revisión legal:** Términos, privacidad, compliance con regulaciones de sorteos

---

## 11. Próximos Pasos Inmediatos

1. **Definir stack de desarrollo:** ✅ (ver [stack_tecnico.md](./stack_tecnico.md))
2. **Crear estructura de carpetas:** En progreso
3. **Setup repositorio Git:** Pendiente
4. **Diseño de base de datos:** Siguiente
5. **Sprint 0 (setup):** Semana 1

---

**Actualizado:** 2025-11-10
**Próxima revisión:** Después de Sprint 2
