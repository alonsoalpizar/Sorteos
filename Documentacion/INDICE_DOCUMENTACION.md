# Índice de Documentación - Plataforma de Sorteos

**Proyecto:** Sistema de Sorteos/Rifas en Línea
**Propietario:** Ing. Alonso Alpízar
**Fecha:** 2025-11-18

---

## 📚 DOCUMENTOS PRINCIPALES

### 🚀 Para Empezar

| Documento | Descripción | Cuándo Leer | Tamaño |
|-----------|-------------|-------------|--------|
| **[README.md](../README.md)** | Información general del proyecto, setup e instalación | Primero - Onboarding | 5 min |
| **[CLAUDE.md](../CLAUDE.md)** | Contexto rápido para AI, stack actual, comandos útiles | AI/Dev - Referencia rápida | 10 min |
| **[SORTEOS_CONTEXTO_COMPLETO.md](SORTEOS_CONTEXTO_COMPLETO.md)** | **NUEVO** - Respuestas detalladas a todas las preguntas clave | Diseño de skill, onboarding completo | 30 min |
| **[RESUMEN_EJECUTIVO_SKILL.md](RESUMEN_EJECUTIVO_SKILL.md)** | **NUEVO** - Versión condensada para referencia rápida | Quick reference antes de codear | 5 min |

---

## 🏗️ Arquitectura y Stack

| Documento | Contenido Clave | Audiencia |
|-----------|-----------------|-----------|
| **[arquitecturaIdeaGeneral.md](arquitecturaIdeaGeneral.md)** | Visión general, problema de concurrencia, solución de locks distribuidos | Arquitectos, Backend devs |
| **[stack_tecnico.md](stack_tecnico.md)** | Stack completo (Go, React, PostgreSQL, Redis), dependencias, versiones | Todos los devs |
| **[modulos.md](modulos.md)** | 7 módulos del sistema con código, casos de uso, interfaces | Backend devs, arquitectos |

**Tiempo total:** 60 minutos

---

## 🎨 Diseño y Frontend

| Documento | Contenido Clave | Audiencia |
|-----------|-----------------|-----------|
| **[estandar_visual.md](estandar_visual.md)** | Design system completo, componentes shadcn/ui, colores | Frontend devs, diseñadores |
| **[.paleta-visual-aprobada.md](.paleta-visual-aprobada.md)** | Referencia rápida de colores permitidos/prohibidos | Todos los devs |
| **[FloatingCheckout.md](FloatingCheckout.md)** | Especificación del botón flotante de checkout | Frontend devs |
| **[arquitectura-navegacion.md](arquitectura-navegacion.md)** | Estructura de navegación y rutas | Frontend devs |

**Tiempo total:** 40 minutos

---

## 🔐 Seguridad y Pagos

| Documento | Contenido Clave | Audiencia |
|-----------|-----------------|-----------|
| **[seguridad.md](seguridad.md)** | JWT, RBAC, rate limiting, prevención OWASP Top 10 | Backend devs, DevOps |
| **[pagos_integraciones.md](pagos_integraciones.md)** | Stripe, PayPal, webhooks, idempotencia | Backend devs |
| **[parametrizacion_reglas.md](parametrizacion_reglas.md)** | 80+ parámetros configurables del sistema | Backend devs, product |

**Tiempo total:** 50 minutos

---

## 📋 Operaciones y Negocio

| Documento | Contenido Clave | Audiencia |
|-----------|-----------------|-----------|
| **[operacion_backoffice.md](operacion_backoffice.md)** | Dashboard admin, liquidaciones, moderación | Backend devs, admins |
| **[terminos_y_condiciones_impacto.md](terminos_y_condiciones_impacto.md)** | GDPR, PCI DSS, cumplimiento legal | Legal, backend devs |

**Tiempo total:** 30 minutos

---

## 📅 Planificación

| Documento | Contenido Clave | Audiencia |
|-----------|-----------------|-----------|
| **[roadmap.md](roadmap.md)** | Fases de desarrollo, MVP, Fase 2, Fase 3 | Product managers, todos |
| **[deployment.md](deployment.md)** | Guía de despliegue, migración Docker → Local | DevOps |
| **[DEPLOYMENT_VALIDATION.md](DEPLOYMENT_VALIDATION.md)** | Validación de migración a instalación local | DevOps |

**Tiempo total:** 45 minutos

---

## 🧪 Testing y QA

| Documento | Contenido Clave | Audiencia |
|-----------|-----------------|-----------|
| **[TESTING-QUICKSTART.md](TESTING-QUICKSTART.md)** | Guía rápida de testing | QA, devs |
| **[testing-strategy.md](testing-strategy.md)** | Estrategia de testing completa | QA lead |
| **[testing-manual-checklist.md](testing-manual-checklist.md)** | Checklist de pruebas manuales | QA |
| **[testing-api-scripts.md](testing-api-scripts.md)** | Scripts curl para probar API | Backend devs |

**Tiempo total:** 30 minutos

---

## 📧 Emails

| Documento | Contenido Clave | Audiencia |
|-----------|-----------------|-----------|
| **[INDICE_DOCUMENTACION_EMAILS.md](INDICE_DOCUMENTACION_EMAILS.md)** | Índice de docs de emails | Backend devs |
| **[README_EMAILS.md](README_EMAILS.md)** | Sistema de emails general | Backend devs |
| **[GUIA_EMAIL_SMTP_VS_SENDGRID.md](GUIA_EMAIL_SMTP_VS_SENDGRID.md)** | Comparación SMTP vs SendGrid | DevOps |
| **[GUIA_PLANTILLAS_EMAIL.md](GUIA_PLANTILLAS_EMAIL.md)** | Plantillas HTML de emails | Frontend devs |
| **[PROPUESTA_EMAILS.md](PROPUESTA_EMAILS.md)** | Propuesta original de emails | Histórico |
| **[RESUMEN_IMPLEMENTACION_EMAIL.md](RESUMEN_IMPLEMENTACION_EMAIL.md)** | Resumen de implementación | Backend devs |

**Tiempo total:** 40 minutos

---

## 🔄 Reservaciones

| Documento | Contenido Clave | Audiencia |
|-----------|-----------------|-----------|
| **[README_RESERVACIONES.md](README_RESERVACIONES.md)** | Sistema de reservaciones general | Backend devs |
| **[PLAN_REFACTORIZACION_RESERVAS.md](PLAN_REFACTORIZACION_RESERVAS.md)** | Plan de refactorización | Backend lead |
| **[CAMBIOS_REFACTORIZACION_RESERVAS.md](CAMBIOS_REFACTORIZACION_RESERVAS.md)** | Cambios implementados | Backend devs |
| **[fase_reservas_concurrencia.md](fase_reservas_concurrencia.md)** | Manejo de concurrencia | Backend devs |
| **[fix_reservacion_limpieza.md](fix_reservacion_limpieza.md)** | Fix de limpieza automática | Backend devs |
| **[fix_mensaje_confuso_reserva_propia.md](fix_mensaje_confuso_reserva_propia.md)** | Fix UX de reservas | Frontend devs |

**Tiempo total:** 50 minutos

---

## 📝 Otros

| Documento | Contenido Clave | Audiencia |
|-----------|-----------------|-----------|
| **[UserX.md](UserX.md)** | Propuesta de UX mejorada | Diseñadores, product |
| **[NEXT_SESSION.md](../NEXT_SESSION.md)** | TODOs para próxima sesión | Todos |
| **[USUARIOS_PRUEBA.md](../USUARIOS_PRUEBA.md)** | Usuarios de prueba para testing | QA |
| **[SPRINT_5-6_PROGRESS.md](../backend/SPRINT_5-6_PROGRESS.md)** | Progreso de sprints recientes | Product managers |

**Tiempo total:** 20 minutos

---

## 🗂️ Documentación por Módulo

### Almighty (Admin Panel)

| Documento | Descripción |
|-----------|-------------|
| **[Almighty/README.md](Almighty/README.md)** | Descripción general del panel admin |
| **[Almighty/ARQUITECTURA_ALMIGHTY.md](Almighty/ARQUITECTURA_ALMIGHTY.md)** | Arquitectura técnica |
| **[Almighty/API_ENDPOINTS.md](Almighty/API_ENDPOINTS.md)** | Endpoints de administración |
| **[Almighty/BASE_DE_DATOS.md](Almighty/BASE_DE_DATOS.md)** | Esquema de DB para admin |
| **[Almighty/CHECKLIST_IMPLEMENTACION.md](Almighty/CHECKLIST_IMPLEMENTACION.md)** | Checklist de tareas |
| **[Almighty/ROADMAP_ALMIGHTY.md](Almighty/ROADMAP_ALMIGHTY.md)** | Roadmap del módulo |

**Tiempo total:** 60 minutos

---

## 🎯 RUTAS DE LECTURA RECOMENDADAS

### 1️⃣ Onboarding Completo (4 horas)

**Para:** Desarrolladores nuevos en el proyecto

**Orden sugerido:**
1. [README.md](../README.md) - 5 min
2. [SORTEOS_CONTEXTO_COMPLETO.md](SORTEOS_CONTEXTO_COMPLETO.md) - 30 min ⭐
3. [stack_tecnico.md](stack_tecnico.md) - 20 min
4. [arquitecturaIdeaGeneral.md](arquitecturaIdeaGeneral.md) - 15 min
5. [modulos.md](modulos.md) - 30 min
6. [estandar_visual.md](estandar_visual.md) - 15 min
7. [seguridad.md](seguridad.md) - 20 min
8. [pagos_integraciones.md](pagos_integraciones.md) - 20 min
9. [roadmap.md](roadmap.md) - 15 min
10. Hands-on: Setup local + build + test

### 2️⃣ Quick Start Backend (1.5 horas)

**Para:** Backend developer que necesita empezar rápido

**Orden sugerido:**
1. [CLAUDE.md](../CLAUDE.md) - 10 min
2. [RESUMEN_EJECUTIVO_SKILL.md](RESUMEN_EJECUTIVO_SKILL.md) - 5 min ⭐
3. [stack_tecnico.md](stack_tecnico.md) - Backend section - 10 min
4. [modulos.md](modulos.md) - Módulos relevantes - 20 min
5. [seguridad.md](seguridad.md) - 15 min
6. [README_RESERVACIONES.md](README_RESERVACIONES.md) - 15 min
7. Código: Explorar `backend/internal/usecase/`

### 3️⃣ Quick Start Frontend (1 hora)

**Para:** Frontend developer que necesita empezar rápido

**Orden sugerido:**
1. [RESUMEN_EJECUTIVO_SKILL.md](RESUMEN_EJECUTIVO_SKILL.md) - 5 min ⭐
2. [stack_tecnico.md](stack_tecnico.md) - Frontend section - 10 min
3. [estandar_visual.md](estandar_visual.md) - 15 min
4. [.paleta-visual-aprobada.md](.paleta-visual-aprobada.md) - 2 min
5. [FloatingCheckout.md](FloatingCheckout.md) - 5 min
6. Código: Explorar `frontend/src/components/ui/`

### 4️⃣ Arquitectura Profunda (2 horas)

**Para:** Arquitectos, tech leads

**Orden sugerido:**
1. [SORTEOS_CONTEXTO_COMPLETO.md](SORTEOS_CONTEXTO_COMPLETO.md) - Sección Arquitectura - 20 min ⭐
2. [arquitecturaIdeaGeneral.md](arquitecturaIdeaGeneral.md) - 15 min
3. [modulos.md](modulos.md) - 30 min
4. [fase_reservas_concurrencia.md](fase_reservas_concurrencia.md) - 15 min
5. [pagos_integraciones.md](pagos_integraciones.md) - 20 min
6. [Almighty/ARQUITECTURA_ALMIGHTY.md](Almighty/ARQUITECTURA_ALMIGHTY.md) - 10 min
7. Código: Revisar arquitectura hexagonal en `internal/`

### 5️⃣ Diseño de Skill para Claude (30 minutos)

**Para:** Configurar AI assistant

**Orden sugerido:**
1. **[SORTEOS_CONTEXTO_COMPLETO.md](SORTEOS_CONTEXTO_COMPLETO.md) - 30 min** ⭐
   - O alternativamente:
2. **[RESUMEN_EJECUTIVO_SKILL.md](RESUMEN_EJECUTIVO_SKILL.md) - 5 min** ⭐
3. [CLAUDE.md](../CLAUDE.md) - 10 min
4. [.paleta-visual-aprobada.md](.paleta-visual-aprobada.md) - 2 min

**Contexto mínimo crítico:**
- Stack: Go + React + PostgreSQL + Redis
- Arquitectura: Hexagonal (backend), Feature-based (frontend)
- Problema central: Concurrencia en reservas → Locks distribuidos (Redis)
- Restricción visual: NO morado/rosa, SÍ azul/gris
- Estado: MVP 60%, reservas + pagos en desarrollo

---

## 📊 ESTADÍSTICAS DE DOCUMENTACIÓN

**Total de archivos:** ~40 documentos
**Tamaño total:** ~500 KB
**Tiempo de lectura completo:** ~8 horas
**Última actualización:** 2025-11-18

**Por categoría:**
- Arquitectura y Stack: 8 docs (90 min)
- Frontend y Diseño: 5 docs (40 min)
- Seguridad y Pagos: 3 docs (50 min)
- Operaciones: 2 docs (30 min)
- Testing: 4 docs (30 min)
- Emails: 6 docs (40 min)
- Reservaciones: 6 docs (50 min)
- Almighty: 6 docs (60 min)
- Planificación: 3 docs (45 min)
- Otros: 4 docs (20 min)

---

## 🔍 BÚSQUEDA RÁPIDA

**¿Buscas información sobre...?**

| Tema | Ver Documento |
|------|---------------|
| Instalación y setup | [README.md](../README.md) |
| Stack completo | [stack_tecnico.md](stack_tecnico.md) |
| Arquitectura hexagonal | [modulos.md](modulos.md), [SORTEOS_CONTEXTO_COMPLETO.md](SORTEOS_CONTEXTO_COMPLETO.md) |
| Concurrencia y locks | [arquitecturaIdeaGeneral.md](arquitecturaIdeaGeneral.md), [fase_reservas_concurrencia.md](fase_reservas_concurrencia.md) |
| JWT y seguridad | [seguridad.md](seguridad.md) |
| Integración Stripe | [pagos_integraciones.md](pagos_integraciones.md) |
| Componentes UI | [estandar_visual.md](estandar_visual.md) |
| Colores permitidos | [.paleta-visual-aprobada.md](.paleta-visual-aprobada.md) |
| Panel de admin | [Almighty/README.md](Almighty/README.md) |
| Testing | [TESTING-QUICKSTART.md](TESTING-QUICKSTART.md) |
| Despliegue | [deployment.md](deployment.md) |
| Roadmap | [roadmap.md](roadmap.md) |
| Emails | [INDICE_DOCUMENTACION_EMAILS.md](INDICE_DOCUMENTACION_EMAILS.md) |
| Contexto AI/Skill | **[SORTEOS_CONTEXTO_COMPLETO.md](SORTEOS_CONTEXTO_COMPLETO.md)** ⭐ |

---

## ⭐ DOCUMENTOS DESTACADOS (MUST READ)

**Top 5 para empezar:**
1. **[SORTEOS_CONTEXTO_COMPLETO.md](SORTEOS_CONTEXTO_COMPLETO.md)** - Responde todas las preguntas clave
2. **[RESUMEN_EJECUTIVO_SKILL.md](RESUMEN_EJECUTIVO_SKILL.md)** - Versión condensada
3. **[CLAUDE.md](../CLAUDE.md)** - Contexto técnico rápido
4. **[stack_tecnico.md](stack_tecnico.md)** - Stack tecnológico detallado
5. **[modulos.md](modulos.md)** - Módulos del sistema con código

**Top 3 para AI/Skill:**
1. **[SORTEOS_CONTEXTO_COMPLETO.md](SORTEOS_CONTEXTO_COMPLETO.md)** ⭐⭐⭐
2. **[RESUMEN_EJECUTIVO_SKILL.md](RESUMEN_EJECUTIVO_SKILL.md)** ⭐⭐
3. **[CLAUDE.md](../CLAUDE.md)** ⭐

---

## 📞 CONTACTO

**Propietario:** Ing. Alonso Alpízar
**Proyecto:** https://sorteos.club
**Documentación:** `/opt/Sorteos/Documentacion/`

---

**Última actualización:** 2025-11-18
**Versión:** 1.0
**Mantenido por:** Equipo de desarrollo Sorteos
