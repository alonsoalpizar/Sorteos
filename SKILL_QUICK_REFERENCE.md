# Skill sorteos-context - Referencia Rápida

**Ubicación:** `/opt/.claude/skills/sorteos-context/`
**Estado:** ✅ Instalado
**Auto-activa en:** `/opt/Sorteos/`

---

## 🚀 Comandos Rápidos

```bash
# Ver archivo principal del skill (7 reglas críticas)
cat /opt/.claude/skills/sorteos-context/SKILL.md

# Ver arquitectura hexagonal
cat /opt/.claude/skills/sorteos-context/references/architecture.md

# Ver reglas de negocio
cat /opt/.claude/skills/sorteos-context/references/business-rules.md

# Ver estado actual
cat /opt/.claude/skills/sorteos-context/references/current-status.md

# Validar proyecto
/opt/.claude/skills/sorteos-context/scripts/validate-structure.sh
```

---

## 🚨 TOP 7 REGLAS (Memorizar)

1. **❌ COLORES:** NUNCA morado/rosa → SOLO azul/gris
2. **🏛️ ARQUITECTURA:** domain NO importa GORM/Gin
3. **🔒 LOCKS:** Redis SETNX OBLIGATORIO en reservas
4. **🔑 IDEMPOTENCIA:** Header Idempotency-Key en pagos
5. **🖥️ NATIVO:** NO Docker → systemd
6. **📝 NAMING:** snake_case Go, PascalCase React
7. **✅ VALIDACIÓN:** Backend + Frontend (dual)

---

## 📁 Estructura del Skill

```
/opt/.claude/skills/sorteos-context/
├── SKILL.md              # ← SIEMPRE CARGAR
├── README.md
├── references/           # ← BAJO DEMANDA
│   ├── architecture.md
│   ├── business-rules.md
│   └── current-status.md
└── scripts/
    └── validate-structure.sh
```

---

## 🎯 Cuándo Cargar Cada Referencia

| Referencia | Cuándo Cargar |
|-----------|---------------|
| **SKILL.md** | **SIEMPRE** (auto-activa) |
| **architecture.md** | Trabajas en capas, separación de concerns |
| **business-rules.md** | Implementas lógica de negocio (reservas, pagos) |
| **current-status.md** | Necesitas saber qué está implementado |

---

## ⚡ Checklist Pre-Código

Antes de escribir código, verificar:

- [ ] ¿Usas colores? → NO morado/rosa
- [ ] ¿Importas en domain/? → NO GORM/Gin
- [ ] ¿Implementas reservas? → Locks Redis
- [ ] ¿Implementas pagos? → Idempotency-Key
- [ ] ¿Sugieres Docker? → NO, usar systemd
- [ ] ¿Naming correcto? → snake_case/PascalCase
- [ ] ¿Validación? → Backend + Frontend

---

## 📚 Más Documentación

- `/opt/.claude/skills/INSTALACION_SKILL.md` - Guía completa
- `/opt/Sorteos/Documentacion/` - Docs del proyecto completo
- `/opt/Sorteos/CLAUDE.md` - Contexto técnico rápido

---

**Última actualización:** 2025-11-18
