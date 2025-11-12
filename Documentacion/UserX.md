🗺️ ANÁLISIS DE NAVEGACIÓN - SORTEOS.CLUB
❌ PROBLEMAS ACTUALES
1. Confusión de Roles
Problema: No hay separación clara entre:
- Usuario participante (compra números)
- Usuario creador (publica sorteos)
- Visitante (explora sorteos)

Actual:
Dashboard → Mezcla de todo
"Crear Sorteo" y "Mis Sorteos" al mismo nivel que "Explorar"
2. Navegación Plana
Header actual:
[Logo] | Buscar | Explorar Sorteos | Crear Sorteo | [Avatar]

Problema:
- Todo tiene el mismo peso visual
- "Crear Sorteo" compite con "Explorar"
- No hay contexto de dónde estás
3. Dashboard Confuso
Muestra:
- Acciones Rápidas (Crear + Explorar + Mis Sorteos)
- Resumen de estadísticas
- Actividad Reciente
- Información de la Cuenta

Problema:
- ¿Es un dashboard de creador o de participante?
- Métricas vacías cuando no has creado sorteos
- No hay diferenciación de contexto

✅ ARQUITECTURA PROPUESTA
Concepto: Dual-Mode Navigation
┌─────────────────────────────────────────────┐
│  MODO 1: PARTICIPANTE (default)             │
│  → Explorar sorteos activos                 │
│  → Comprar números                          │
│  → Ver mis participaciones                  │
└─────────────────────────────────────────────┘

┌─────────────────────────────────────────────┐
│  MODO 2: ORGANIZADOR                        │
│  → Crear y gestionar mis sorteos            │
│  → Ver ventas y estadísticas                │
│  → Gestionar ganadores                      │
└─────────────────────────────────────────────┘

🎯 ESTRUCTURA DE NAVEGACIÓN PROPUESTA
A. NAVEGACIÓN PRINCIPAL (Para Participantes)
┌─────────────────────────────────────────────────────┐
│  [Logo]  Explorar   Mis Números   [🔔]   [Avatar]   │
└─────────────────────────────────────────────────────┘

Rutas:
/                    → Landing pública
/explore             → Sorteos activos (PRINCIPAL)
/my-tickets          → Mis participaciones
/raffle/:id          → Detalle de un sorteo
Lógica:

El foco es EXPLORAR sorteos (marketplace)
"Mis Números" muestra tus participaciones
Simple, claro, sin ruido


B. NAVEGACIÓN ORGANIZADOR (Modo Creador)
┌──────────────────────────────────────────────────────────┐
│  [Logo]  📊 Panel  🎲 Sorteos  ➕ Crear   [🔔]  [Avatar] │
└──────────────────────────────────────────────────────────┘

Rutas:
/organizer                     → Dashboard del organizador
/organizer/raffles             → Lista de mis sorteos
/organizer/raffles/new         → Crear sorteo
/organizer/raffles/:id         → Gestionar sorteo
/organizer/raffles/:id/edit    → Editar sorteo
/organizer/sales               → Reporte de ventas
/organizer/winners             → Gestión de ganadores
Lógica:

Es un "área protegida" (sub-app)
Métricas relevantes (ventas, participantes)
Herramientas de gestión


C. CAMBIO DE CONTEXTO (Toggle)
tsx// En el Avatar Dropdown:

┌─────────────────────────────┐
│ CacheClear Test             │
│ cacheclear@example.com      │
├─────────────────────────────┤
│ 👤 Mi Perfil                │
│ 🎟️ Mis Números              │
├─────────────────────────────┤
│ 🎲 Modo Organizador    →    │  ← TOGGLE
├─────────────────────────────┤
│ ⚙️ Configuración            │
│ 🚪 Cerrar Sesión            │
└─────────────────────────────┘

Cuando activas "Modo Organizador":
→ Cambio visual del header
→ Navegación diferente
→ Dashboard diferente
```

---

## 📐 WIREFRAMES DE NAVEGACIÓN

### **1. LANDING (No autenticado)**
```
┌────────────────────────────────────────────────┐
│  [Logo] Sorteos.club    Cómo Funciona  Login  │
└────────────────────────────────────────────────┘

             🎯 Gana con Sorteos Verificables
        Participa en sorteos basados en Lotería Nacional
        
        [Explorar Sorteos →]  [Crear mi Sorteo]
        
        ✓ 100% Transparente   ✓ 24/7   ✓ Seguro
```

**Acciones claras:**
- Explorar → Marketplace público
- Crear → Registro + Onboarding de organizador

---

### **2. EXPLORAR (Participante - Autenticado)**
```
┌──────────────────────────────────────────────────────┐
│  [🎲] Sorteos    Explorar  Mis Números  [🔔]  [CT▼]  │
└──────────────────────────────────────────────────────┘

[🔍 Buscar sorteos...]    [Filtros ▼]

┌─────────────────────┐  ┌─────────────────────┐
│ iPhone 15 Pro       │  │ PlayStation 5       │
│ ₡50,000             │  │ ₡25,000             │
│ 🎟️ 234/500          │  │ 🎟️ 89/200           │
│ ⏰ 3 días           │  │ ⏰ 5 horas          │
│ [Participar →]      │  │ [Participar →]      │
└─────────────────────┘  └─────────────────────┘

Sidebar:
├─ Categorías
│  ├─ 📱 Electrónica
│  ├─ 🏍️ Vehículos
│  └─ 💰 Efectivo
├─ Estado
│  ├─ 🟢 Activos
│  └─ ⏳ Próximos
└─ Precio
   ├─ < ₡10,000
   └─ ₡10,000 - ₡50,000
```

**Características:**
- Browse de sorteos activos
- Filtros claros
- CTAs directos para participar

---

### **3. MIS NÚMEROS (Participante)**
```
┌──────────────────────────────────────────────────────┐
│  [🎲] Sorteos    Explorar  Mis Números  [🔔]  [CT▼]  │
└──────────────────────────────────────────────────────┘

Mis Participaciones

Tabs: [Activos]  [Finalizados]  [Ganados 🎉]

┌─────────────────────────────────────────────────┐
│ iPhone 15 Pro                                   │
│ Números: #0234, #0567, #0891                    │
│ Total: ₡1,500    |    Sorteo: 15 Dic 8:00 PM   │
│ Estado: ⏳ En espera                            │
└─────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────┐
│ PlayStation 5                                   │
│ Números: #0042                                  │
│ Total: ₡500    |    Sorteo: 18 Dic 9:00 PM     │
│ Estado: 🟢 Activo                               │
└─────────────────────────────────────────────────┘
```

**Características:**
- Historial de participaciones
- Status claro de cada sorteo
- Separación por estado

---

### **4. ORGANIZADOR - DASHBOARD**
```
┌────────────────────────────────────────────────────────┐
│ [🎲] Panel  Sorteos  ➕ Crear   [🔔]  [Volver a Participar] │
└────────────────────────────────────────────────────────┘

Panel de Organizador

┌─────────────────┐ ┌─────────────────┐ ┌─────────────────┐
│ Sorteos Activos │ │ Ventas del Mes  │ │ Participantes   │
│      3          │ │   ₡45,000       │ │     234         │
└─────────────────┘ └─────────────────┘ └─────────────────┘

Sorteos Recientes
┌──────────────────────────────────────────────────────┐
│ iPhone 15 Pro        🟢 Activo     234/500  [Gestionar] │
│ PlayStation 5        ⏳ Pendiente   89/200  [Gestionar] │
│ MacBook Air          🔴 Finalizado 500/500  [Ver]       │
└──────────────────────────────────────────────────────┘

[➕ Crear Nuevo Sorteo]
```

**Características:**
- Métricas de organizador
- Acceso rápido a gestión
- Separación de contexto clara

---

### **5. ORGANIZAR - MIS SORTEOS**
```
┌────────────────────────────────────────────────────────┐
│ [🎲] Panel  Sorteos  ➕ Crear   [🔔]  [Volver]          │
└────────────────────────────────────────────────────────┘

Mis Sorteos

[➕ Crear Sorteo]   Filtros: [Todos ▼]  [Buscar...]

┌──────────────────────────────────────────────────────┐
│ iPhone 15 Pro                            🟢 Activo   │
│ Premio: ₡50,000  |  Vendidos: 234/500  |  ₡11,700    │
│ Sorteo: 15 Dic 8:00 PM  |  Basado en: Lotería Nacional │
│                                                        │
│ [📊 Ver Reporte]  [✏️ Editar]  [🎲 Sortear]          │
└──────────────────────────────────────────────────────┘

┌──────────────────────────────────────────────────────┐
│ PlayStation 5                         ⏳ Programado  │
│ Premio: ₡25,000  |  Vendidos: 0/200   |  ₡0          │
│ Inicia: 20 Dic 6:00 PM                               │
│                                                        │
│ [✏️ Editar]  [🗑️ Eliminar]                          │
└──────────────────────────────────────────────────────┘
```

**Características:**
- Lista de sorteos con actions
- Status claros
- Accesos rápidos a gestión

---

### **6. CREAR SORTEO (Flujo completo)**
```
┌────────────────────────────────────────────────────────┐
│ [🎲] Panel  Sorteos  ➕ Crear   [🔔]  [Volver]          │
└────────────────────────────────────────────────────────┘

Crear Nuevo Sorteo

Progress: ●━━━━ 1/4  [Información Básica]

┌──────────────────────────────────────────────────────┐
│                                                        │
│  📸 Foto del Premio                                   │
│     [Subir imagen]                                    │
│                                                        │
│  🎁 Nombre del Premio *                               │
│     [iPhone 15 Pro 256GB]                             │
│                                                        │
│  📝 Descripción                                       │
│     [Describe el premio en detalle...]                │
│                                                        │
│  🏷️ Categoría *                                       │
│     [Electrónica ▼]                                   │
│                                                        │
│              [← Cancelar]  [Siguiente →]              │
└──────────────────────────────────────────────────────┘

Pasos:
1. Información Básica
2. Configuración de Números
3. Fecha y Lotería
4. Revisión y Publicar
```

**Características:**
- Multi-step claro
- Progress indicator
- Validación por paso
- Guardado automático

---

## 🎯 DIAGRAMA DE FLUJO
```
                    ┌──────────────┐
                    │   Landing    │
                    └──────┬───────┘
                           │
              ┌────────────┴────────────┐
              │                         │
       ┌──────▼──────┐          ┌──────▼──────┐
       │   Explorar  │          │  Registrarse │
       │  (público)  │          └──────┬───────┘
       └──────┬──────┘                 │
              │                  ┌─────▼─────┐
              │                  │   Login   │
              │                  └─────┬─────┘
              │                        │
       ┌──────▼────────────────────────▼──────┐
       │                                       │
┌──────▼──────────┐              ┌────────────▼─────────┐
│   PARTICIPANTE  │              │     ORGANIZADOR      │
├─────────────────┤              ├──────────────────────┤
│ • Explorar      │              │ • Panel              │
│ • Mis Números   │◄────toggle───┤ • Mis Sorteos        │
│ • Comprar       │              │ • Crear Sorteo       │
└─────────────────┘              │ • Gestionar          │
                                 └──────────────────────┘

📋 RESUMEN DE CAMBIOS RECOMENDADOS
Cambios Inmediatos:
markdown1. Separar Dashboard:
   - /dashboard → Para organizadores con sorteos
   - / o /explore → Landing para participantes

2. Reorganizar Header:
   - Participante: Explorar | Mis Números
   - Organizador: Panel | Sorteos | Crear

3. Agregar Toggle de Contexto:
   - En dropdown de avatar
   - Cambia navegación completa

4. Multi-step en Crear Sorteo:
   - Paso 1: Info básica
   - Paso 2: Configuración
   - Paso 3: Fecha/Lotería
   - Paso 4: Publicar

5. "Mis Sorteos" separado de "Mis Números":
   - Mis Sorteos → Organizador
   - Mis Números → Participante

🎨 COMPONENTES NECESARIOS
tsx// 1. ContextToggle
<ContextToggle 
  current="participant" 
  onChange={(mode) => navigate(mode === 'organizer' ? '/organizer' : '/explore')}
/>

// 2. NavBar condicional
<NavBar mode={userContext} />

// 3. MultiStepForm
<MultiStepForm 
  steps={[BasicInfo, NumberConfig, Schedule, Review]}
  onComplete={handlePublish}
/>

// 4. RaffleCard (2 versiones)
<RaffleCard.Browse />      // Para explorar
<RaffleCard.Manage />      // Para gestionar

// 5. EmptyState contextual
<EmptyState.Participant />
<EmptyState.Organizer />
```

---

## 💡 PROPUESTA DE IMPLEMENTACIÓN

### Fase 1: Separación de Contextos (3-4 días)
```
✅ Crear rutas /organizer/*
✅ Dual navigation (ParticipantNav + OrganizerNav)
✅ Context toggle en avatar dropdown
✅ Redirect lógico basado en rol
```

### Fase 2: Reorganizar Dashboard (2-3 días)
```
✅ Dashboard de organizador con métricas relevantes
✅ "Explorar" como landing principal
✅ "Mis Números" para participaciones
✅ Empty states contextuales
```

### Fase 3: Crear Sorteo Multi-Step (4-5 días)
```
✅ Wizard de 4 pasos
✅ Validación por paso
✅ Preview antes de publicar
✅ Guardado automático (draft)