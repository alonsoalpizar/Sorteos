# 📖 LEE ESTO PRIMERO - Plataforma de Sorteos

**Fecha:** 2025-11-18
**Última actualización:** 2025-11-18

---

## 🎯 ¿QUÉ ES ESTO?

Este directorio contiene **toda la documentación técnica** de la Plataforma de Sorteos.

**Total:** ~40 documentos, 500 KB, 8 horas de lectura completa

---

## ⚡ INICIO RÁPIDO

### Si eres nuevo en el proyecto:

**Opción 1 - Contexto Completo (30 min):**
```bash
cat SORTEOS_CONTEXTO_COMPLETO.md
```
→ Lee esto si necesitas entender TODO el proyecto en profundidad

**Opción 2 - Resumen Rápido (5 min):**
```bash
cat RESUMEN_EJECUTIVO_SKILL.md
```
→ Lee esto si necesitas empezar rápido

**Opción 3 - Contexto AI (10 min):**
```bash
cat ../CLAUDE.md
```
→ Lee esto si eres una IA o quieres contexto técnico rápido

---

## 📚 DOCUMENTOS PRINCIPALES

### 🆕 Nuevos (2025-11-18)

| Archivo | Líneas | Tamaño | Tiempo | Propósito |
|---------|--------|--------|--------|-----------|
| **[SORTEOS_CONTEXTO_COMPLETO.md](SORTEOS_CONTEXTO_COMPLETO.md)** | 1,495 | 45 KB | 30 min | **Contexto completo para skill** ⭐ |
| **[RESUMEN_EJECUTIVO_SKILL.md](RESUMEN_EJECUTIVO_SKILL.md)** | 454 | 12 KB | 5 min | **Resumen condensado** ⭐ |
| **[INDICE_DOCUMENTACION.md](INDICE_DOCUMENTACION.md)** | 301 | 13 KB | 10 min | **Índice y navegación** |

### 📖 Fundamentales

| Archivo | Contenido | Audiencia |
|---------|-----------|-----------|
| [stack_tecnico.md](stack_tecnico.md) | Stack completo: Go, React, PostgreSQL, Redis | Todos |
| [arquitecturaIdeaGeneral.md](arquitecturaIdeaGeneral.md) | Arquitectura y concurrencia | Arquitectos, Backend |
| [modulos.md](modulos.md) | 7 módulos con código | Backend |
| [estandar_visual.md](estandar_visual.md) | Design system, componentes UI | Frontend |
| [seguridad.md](seguridad.md) | JWT, RBAC, rate limiting | Backend, DevOps |
| [pagos_integraciones.md](pagos_integraciones.md) | Stripe, webhooks, idempotencia | Backend |
| [roadmap.md](roadmap.md) | Fases de desarrollo | Product, Todos |

---

## 🔍 BÚSQUEDA RÁPIDA

**¿Necesitas información sobre...?**

```
Instalación → ../README.md
Stack → stack_tecnico.md
Arquitectura → SORTEOS_CONTEXTO_COMPLETO.md, modulos.md
Concurrencia → arquitecturaIdeaGeneral.md
Seguridad → seguridad.md
Pagos → pagos_integraciones.md
UI/Colores → estandar_visual.md, .paleta-visual-aprobada.md
Admin → Almighty/README.md
Testing → TESTING-QUICKSTART.md
Deploy → deployment.md
Emails → INDICE_DOCUMENTACION_EMAILS.md
```

---

## 🎓 RUTAS DE APRENDIZAJE

### 1️⃣ Desarrollador Nuevo (2 horas)
```
1. ../README.md (5 min)
2. RESUMEN_EJECUTIVO_SKILL.md (5 min) ⭐
3. stack_tecnico.md - Tu sección (10 min)
4. estandar_visual.md o modulos.md según rol (20 min)
5. Hands-on: Setup + build + test (60 min)
```

### 2️⃣ Arquitecto / Tech Lead (1.5 horas)
```
1. SORTEOS_CONTEXTO_COMPLETO.md (30 min) ⭐
2. arquitecturaIdeaGeneral.md (15 min)
3. modulos.md (30 min)
4. Código: Revisar internal/ (15 min)
```

### 3️⃣ Diseñar Skill de AI (30 min)
```
1. SORTEOS_CONTEXTO_COMPLETO.md (30 min) ⭐
   O alternativamente:
   RESUMEN_EJECUTIVO_SKILL.md (5 min) +
   ../CLAUDE.md (10 min) +
   .paleta-visual-aprobada.md (2 min)
```

---

## 📊 LO QUE DEBES SABER (MÍNIMO)

### Stack en una línea:
**Go + Gin + PostgreSQL + Redis + React + TypeScript + Vite + Tailwind + shadcn/ui**

### Problema central:
**Doble venta de números → Locks distribuidos (Redis SETNX)**

### Arquitectura:
**Hexagonal (backend) + Feature-based (frontend) + Instalación nativa (sin Docker)**

### Estado actual:
**MVP 60% - Auth ✅ Sorteos ✅ Reservas 🚧 Pagos 🚧**

### Restricción visual:
**NO morado/rosa, SÍ azul/gris/verde/ámbar/rojo**

---

## 🗂️ ORGANIZACIÓN DE DOCUMENTOS

```
Documentacion/
├── LEEME_PRIMERO.md              ← ESTÁS AQUÍ
├── SORTEOS_CONTEXTO_COMPLETO.md  ← PRINCIPAL ⭐
├── RESUMEN_EJECUTIVO_SKILL.md    ← QUICK REF ⭐
├── INDICE_DOCUMENTACION.md       ← ÍNDICE COMPLETO
│
├── stack_tecnico.md              ← Tecnologías
├── arquitecturaIdeaGeneral.md    ← Arquitectura
├── modulos.md                    ← Módulos del sistema
├── estandar_visual.md            ← Design system
├── seguridad.md                  ← Seguridad
├── pagos_integraciones.md        ← Pagos
├── roadmap.md                    ← Roadmap
│
├── Almighty/                     ← Admin panel (6 docs)
├── [Testing docs...]             ← Tests (4 docs)
├── [Email docs...]               ← Emails (6 docs)
├── [Reservas docs...]            ← Reservaciones (6 docs)
└── [Otros...]                    ← 10+ docs más
```

---

## 🚀 COMANDOS ÚTILES

### Ver documentos principales:
```bash
# Contexto completo
cat SORTEOS_CONTEXTO_COMPLETO.md | less

# Resumen rápido
cat RESUMEN_EJECUTIVO_SKILL.md | less

# Índice
cat INDICE_DOCUMENTACION.md | less
```

### Buscar en toda la documentación:
```bash
grep -r "keyword" .
```

### Ver estadísticas:
```bash
find . -name "*.md" | wc -l  # Total de archivos
du -sh .                      # Tamaño total
```

---

## ❓ PREGUNTAS FRECUENTES

### ¿Por dónde empiezo?
→ Lee **RESUMEN_EJECUTIVO_SKILL.md** (5 min)

### ¿Necesito leer todo?
→ No. Usa el **INDICE_DOCUMENTACION.md** para navegar según tu rol

### ¿Dónde está el código?
→ Backend: `/opt/Sorteos/backend/`
→ Frontend: `/opt/Sorteos/frontend/`

### ¿Cómo contribuyo?
→ Lee primero SORTEOS_CONTEXTO_COMPLETO.md sección "Decisiones Técnicas"

### ¿Hay ejemplos de código?
→ Sí, en **modulos.md** hay 7 módulos con código completo

---

## 📞 AYUDA

**Propietario:** Ing. Alonso Alpízar
**Proyecto:** https://sorteos.club
**Ubicación:** `/opt/Sorteos/Documentacion/`

**Si estás perdido:**
1. Lee este archivo (LEEME_PRIMERO.md)
2. Lee RESUMEN_EJECUTIVO_SKILL.md
3. Usa INDICE_DOCUMENTACION.md para navegar

---

## ✅ CHECKLIST INICIAL

Antes de empezar a codear, asegúrate de:

- [ ] Leer RESUMEN_EJECUTIVO_SKILL.md (5 min)
- [ ] Conocer el stack (Go + React + PostgreSQL + Redis)
- [ ] Entender el problema de concurrencia (locks distribuidos)
- [ ] Conocer la arquitectura (hexagonal)
- [ ] Revisar colores permitidos (.paleta-visual-aprobada.md)
- [ ] Setup local completado (../README.md)
- [ ] Build y tests funcionando

---

**Última actualización:** 2025-11-18
**Versión:** 1.0
**Próxima actualización:** Cuando haya cambios significativos en arquitectura o stack

---

**🎯 TIP:** Si solo tienes 5 minutos, lee **RESUMEN_EJECUTIVO_SKILL.md**
