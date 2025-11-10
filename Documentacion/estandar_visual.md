# Estándar Visual - Plataforma de Sorteos

**Versión:** 1.0
**Fecha:** 2025-11-10
**Design System:** Basado en Tailwind CSS + shadcn/ui

---

## 1. Introducción

Este documento define el estándar visual obligatorio para toda la plataforma (web y móvil). El objetivo es garantizar:

- **Consistencia** en toda la experiencia de usuario
- **Accesibilidad** WCAG 2.1 nivel AA
- **Escalabilidad** mediante design tokens y componentes reutilizables
- **Mantenibilidad** con una única fuente de verdad
- **Profesionalismo** con una paleta sobria que transmite credibilidad y confianza

**Tecnología base:**
- **Tailwind CSS 3.4+** para utilidades
- **shadcn/ui** para componentes base (construidos sobre Radix UI)
- **CSS Variables** para theming

**⚠️ RESTRICCIÓN OBLIGATORIA DE COLORES:**
- **PROHIBIDO:** Morado, púrpura, violeta, magenta, tonalidades arcoíris
- **PERMITIDO:** Azules, grises, verdes (éxito), naranjas (advertencia), rojos (error)
- **Objetivo:** Diseño sobrio, aesthetic, profesional y de alta calidad que transmite credibilidad

---

## 2. Design Tokens

### 2.1 Colores

**Filosofía de Color:**
- Paleta sobria y profesional
- Enfoque en azules (confianza, seguridad financiera) y grises (elegancia)
- Verde solo para estados de éxito
- Naranja/ámbar para advertencias controladas
- Rojo solo para errores críticos

#### Paleta Principal (Light Mode)

```css
:root {
  /* Primary - Azul Corporativo (CTAs, links, confianza)
     Transmite: Seguridad, confiabilidad, profesionalismo */
  --color-primary-50: #eff6ff;
  --color-primary-100: #dbeafe;
  --color-primary-200: #bfdbfe;
  --color-primary-300: #93c5fd;
  --color-primary-400: #60a5fa;
  --color-primary-500: #3b82f6; /* Base - Azul confiable */
  --color-primary-600: #2563eb;
  --color-primary-700: #1d4ed8;
  --color-primary-800: #1e40af;
  --color-primary-900: #1e3a8a;

  /* Secondary - Azul Oscuro/Slate (alternativa profesional)
     Transmite: Seriedad, corporativo, elegancia */
  --color-secondary-50: #f8fafc;
  --color-secondary-100: #f1f5f9;
  --color-secondary-200: #e2e8f0;
  --color-secondary-300: #cbd5e1;
  --color-secondary-400: #94a3b8;
  --color-secondary-500: #64748b; /* Base - Slate profesional */
  --color-secondary-600: #475569;
  --color-secondary-700: #334155;
  --color-secondary-800: #1e293b;
  --color-secondary-900: #0f172a;

  /* Neutral - Grises (interfaces, textos, bordes)
     Transmite: Neutralidad, limpieza, minimalismo */
  --color-neutral-50: #fafafa;
  --color-neutral-100: #f5f5f5;
  --color-neutral-200: #e5e5e5;
  --color-neutral-300: #d4d4d4;
  --color-neutral-400: #a3a3a3;
  --color-neutral-500: #737373;
  --color-neutral-600: #525252;
  --color-neutral-700: #404040;
  --color-neutral-800: #262626;
  --color-neutral-900: #171717;

  /* Semantic Colors - USO RESTRINGIDO */
  --color-success: #10b981;  /* Verde esmeralda - Solo para confirmaciones */
  --color-warning: #f59e0b;  /* Ámbar - Solo para advertencias */
  --color-error: #ef4444;    /* Rojo - Solo para errores críticos */
  --color-info: #3b82f6;     /* Azul primary - Para información */

  /* Backgrounds - Paleta limpia y profesional */
  --bg-primary: #ffffff;
  --bg-secondary: #f9fafb;
  --bg-tertiary: #f3f4f6;
  --bg-elevated: #ffffff;  /* Cards elevadas con shadow */

  /* Text - Jerarquía clara */
  --text-primary: #111827;    /* Títulos y texto principal */
  --text-secondary: #6b7280;  /* Texto secundario */
  --text-tertiary: #9ca3af;   /* Texto terciario, placeholders */
  --text-disabled: #d1d5db;   /* Texto deshabilitado */
}
```

**⚠️ COLORES PROHIBIDOS (NO USAR JAMÁS):**
```css
/* ❌ PROHIBIDO - No usar en ningún contexto */
--color-purple: PROHIBIDO;   /* Morado/Púrpura */
--color-violet: PROHIBIDO;   /* Violeta */
--color-magenta: PROHIBIDO;  /* Magenta */
--color-pink: PROHIBIDO;     /* Rosa/Pink */
--color-fuchsia: PROHIBIDO;  /* Fucsia */
/* Cualquier gradiente arcoíris está prohibido */
```

#### Dark Mode

**Filosofía Dark Mode:**
- Fondo oscuro profesional (slate/azul muy oscuro)
- Contraste WCAG AAA en textos
- Azules más claros para mantener legibilidad
- Sin colores saturados que cansen la vista

```css
.dark {
  /* Primary ajustado para dark mode */
  --color-primary-400: #60a5fa;
  --color-primary-500: #3b82f6;
  --color-primary-600: #2563eb;

  /* Backgrounds oscuros profesionales */
  --bg-primary: #0f172a;      /* Slate 900 - Fondo principal */
  --bg-secondary: #1e293b;    /* Slate 800 - Fondo secundario */
  --bg-tertiary: #334155;     /* Slate 700 - Fondo terciario */
  --bg-elevated: #1e293b;     /* Cards con ligero elevation */

  /* Textos con contraste AAA */
  --text-primary: #f1f5f9;    /* Slate 100 - Texto principal */
  --text-secondary: #cbd5e1;  /* Slate 300 - Texto secundario */
  --text-tertiary: #94a3b8;   /* Slate 400 - Texto terciario */
  --text-disabled: #64748b;   /* Slate 500 - Deshabilitado */

  /* Neutral invertidos */
  --color-neutral-700: #e5e5e5;
  --color-neutral-800: #f5f5f5;
}
```

**⚠️ Dark Mode también prohibe:**
- Gradientes neón o fluorescentes
- Colores saturados tipo cyberpunk
- Tonalidades moradas/púrpuras
- Mantener sobriedad y elegancia

**Configuración Tailwind:**
```js
// tailwind.config.js
module.exports = {
  darkMode: 'class',
  theme: {
    extend: {
      colors: {
        primary: {
          50: 'var(--color-primary-50)',
          100: 'var(--color-primary-100)',
          200: 'var(--color-primary-200)',
          300: 'var(--color-primary-300)',
          400: 'var(--color-primary-400)',
          500: 'var(--color-primary-500)',
          600: 'var(--color-primary-600)',
          700: 'var(--color-primary-700)',
          800: 'var(--color-primary-800)',
          900: 'var(--color-primary-900)',
        },
        secondary: {
          50: 'var(--color-secondary-50)',
          100: 'var(--color-secondary-100)',
          // ... resto de tonos
        },
        neutral: {
          50: 'var(--color-neutral-50)',
          // ... resto de tonos
        },
        // ⚠️ IMPORTANTE: NO agregar purple, violet, pink, magenta, fuchsia
      },
    },
  },
}
```

---

### 2.1.1 Referencias Visuales (Aesthetic Profesional)

**Inspiración de diseño (sitios de referencia):**
- **Stripe.com** - Paleta azul/slate profesional
- **Linear.app** - Minimalismo con grises y azules
- **Vercel.com** - Dark mode elegante sin colores saturados
- **Coinbase.com** - Confianza financiera con azules corporativos
- **Notion.com** - UI limpia y profesional

**Características clave:**
- Espacios en blanco generosos
- Tipografía clara y legible (Inter, sans-serif)
- Sombras sutiles (no exageradas)
- Bordes suaves (radius 8-16px)
- Animaciones suaves (150-300ms)

**❌ Anti-referencias (NO seguir):**
- Sitios con gradientes arcoíris
- Diseños tipo "gaming" con neón
- Interfaces con colores saturados
- Paletas vibrantes tipo Material Design 1.0

---

### 2.1.2 Checklist de Validación de Colores

Antes de implementar cualquier componente, verificar:

- [ ] **No usa morado/púrpura** en ninguna tonalidad
- [ ] **No usa violeta** (ni #8B5CF6 ni similares)
- [ ] **No usa magenta/fucsia** (ni #EC4899 ni similares)
- [ ] **No usa rosa/pink** (ni #F472B6 ni similares)
- [ ] **No usa gradientes arcoíris** (rainbow gradients)
- [ ] **Colores principales son azul + gris** (primary + neutral)
- [ ] **Verde solo para success** (confirmaciones, éxito)
- [ ] **Naranja solo para warnings** (advertencias)
- [ ] **Rojo solo para errors** (errores críticos)
- [ ] **Contraste WCAG AA mínimo** (4.5:1 para texto)
- [ ] **Paleta transmite profesionalismo** (no colores "juguetones")

**Herramienta de validación:**
```tsx
// utils/validateColor.ts
const PROHIBITED_COLORS = [
  '#8B5CF6', // violet-500
  '#A855F7', // purple-500
  '#EC4899', // pink-500
  '#F472B6', // pink-400
  '#D946EF', // fuchsia-500
  // Agregar cualquier tono morado/rosa/magenta
]

export function isColorProhibited(hexColor: string): boolean {
  // Verificar si el color está en la lista prohibida
  // o si está en el rango HSL de morados (270-320 grados)
  return PROHIBITED_COLORS.includes(hexColor.toUpperCase())
}
```

---

### 2.2 Tipografía

#### Fuentes

**Principal:** Inter (sans-serif)
**Monospace:** Fira Code (para códigos, números de sorteo)

```css
@import url('https://fonts.googleapis.com/css2?family=Inter:wght@400;500;600;700&display=swap');

:root {
  --font-sans: 'Inter', system-ui, -apple-system, sans-serif;
  --font-mono: 'Fira Code', monospace;
}
```

#### Escalas

```css
:root {
  /* Font Sizes */
  --text-xs: 0.75rem;      /* 12px */
  --text-sm: 0.875rem;     /* 14px */
  --text-base: 1rem;       /* 16px */
  --text-lg: 1.125rem;     /* 18px */
  --text-xl: 1.25rem;      /* 20px */
  --text-2xl: 1.5rem;      /* 24px */
  --text-3xl: 1.875rem;    /* 30px */
  --text-4xl: 2.25rem;     /* 36px */
  --text-5xl: 3rem;        /* 48px */

  /* Line Heights */
  --leading-tight: 1.25;
  --leading-normal: 1.5;
  --leading-relaxed: 1.75;

  /* Font Weights */
  --font-normal: 400;
  --font-medium: 500;
  --font-semibold: 600;
  --font-bold: 700;
}
```

**Uso en Tailwind:**
```jsx
<h1 className="text-4xl font-bold text-primary-900">Título</h1>
<p className="text-base text-neutral-700 leading-relaxed">Párrafo</p>
```

---

### 2.3 Espaciado

**Escala modular (base 4px):**

```js
// tailwind.config.js
spacing: {
  '0': '0',
  '1': '0.25rem',  // 4px
  '2': '0.5rem',   // 8px
  '3': '0.75rem',  // 12px
  '4': '1rem',     // 16px
  '6': '1.5rem',   // 24px
  '8': '2rem',     // 32px
  '12': '3rem',    // 48px
  '16': '4rem',    // 64px
}
```

**Uso:**
- Margin/padding entre secciones: `py-16` (64px)
- Padding de cards: `p-6` (24px)
- Gap en grids: `gap-4` (16px)
- Spacing entre elementos inline: `space-x-2` (8px)

---

### 2.4 Bordes y Sombras

#### Radios de Borde

```css
:root {
  --radius-sm: 0.25rem;   /* 4px */
  --radius-md: 0.5rem;    /* 8px */
  --radius-lg: 0.75rem;   /* 12px */
  --radius-xl: 1rem;      /* 16px */
  --radius-full: 9999px;  /* Círculo */
}
```

**Uso:**
- Buttons: `rounded-lg` (12px)
- Cards: `rounded-xl` (16px)
- Inputs: `rounded-md` (8px)
- Badges: `rounded-full`

#### Sombras

```css
:root {
  --shadow-sm: 0 1px 2px 0 rgb(0 0 0 / 0.05);
  --shadow-md: 0 4px 6px -1px rgb(0 0 0 / 0.1);
  --shadow-lg: 0 10px 15px -3px rgb(0 0 0 / 0.1);
  --shadow-xl: 0 20px 25px -5px rgb(0 0 0 / 0.1);
}
```

**Uso:**
- Cards: `shadow-md`
- Modals: `shadow-xl`
- Dropdowns: `shadow-lg`

---

## 3. Grid System

**Layout principal:** 12 columnas con max-width containers

```jsx
<div className="container mx-auto px-4 max-w-7xl">
  <div className="grid grid-cols-12 gap-6">
    <div className="col-span-12 md:col-span-8">Main</div>
    <div className="col-span-12 md:col-span-4">Sidebar</div>
  </div>
</div>
```

**Breakpoints:**
```js
screens: {
  'sm': '640px',
  'md': '768px',
  'lg': '1024px',
  'xl': '1280px',
  '2xl': '1536px',
}
```

---

## 4. Componentes Base (shadcn/ui)

### 4.1 Button

**Variantes:**
- `default`: Primario (azul)
- `secondary`: Secundario (gris)
- `outline`: Borde con fondo transparente
- `ghost`: Sin fondo, solo hover
- `destructive`: Rojo para acciones peligrosas

**Tamaños:**
- `sm`: 32px altura
- `md`: 40px altura (default)
- `lg`: 48px altura

**Ejemplo:**
```tsx
import { Button } from '@/components/ui/button'

<Button variant="default" size="lg">
  Comprar Boleto
</Button>

<Button variant="outline" size="md">
  Ver Detalles
</Button>
```

**Estados:**
- Hover: `brightness(90%)`
- Active: `brightness(85%)`
- Disabled: `opacity: 0.5, cursor: not-allowed`
- Loading: Spinner + `cursor: wait`

---

### 4.2 Input

**Variantes:**
- Default
- Error (borde rojo)
- Success (borde verde)
- Disabled

**Ejemplo:**
```tsx
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'

<div>
  <Label htmlFor="email">Correo electrónico</Label>
  <Input
    id="email"
    type="email"
    placeholder="tu@email.com"
    className="mt-1"
  />
</div>
```

**Con error:**
```tsx
<Input
  type="text"
  className="border-error focus:ring-error"
  aria-invalid="true"
/>
<p className="text-sm text-error mt-1">Este campo es requerido</p>
```

---

### 4.3 Select

```tsx
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'

<Select>
  <SelectTrigger className="w-full">
    <SelectValue placeholder="Selecciona categoría" />
  </SelectTrigger>
  <SelectContent>
    <SelectItem value="tech">Tecnología</SelectItem>
    <SelectItem value="travel">Viajes</SelectItem>
  </SelectContent>
</Select>
```

---

### 4.4 Card

```tsx
import {
  Card,
  CardHeader,
  CardTitle,
  CardDescription,
  CardContent,
  CardFooter,
} from '@/components/ui/card'

<Card>
  <CardHeader>
    <CardTitle>iPhone 15 Pro</CardTitle>
    <CardDescription>Sorteo el 25 de diciembre</CardDescription>
  </CardHeader>
  <CardContent>
    <img src="..." alt="iPhone" className="rounded-lg" />
    <p className="mt-4 text-neutral-700">Descripción...</p>
  </CardContent>
  <CardFooter className="flex justify-between">
    <span className="text-2xl font-bold">$5/boleto</span>
    <Button>Participar</Button>
  </CardFooter>
</Card>
```

---

### 4.5 Table

```tsx
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table'

<Table>
  <TableHeader>
    <TableRow>
      <TableHead>Sorteo</TableHead>
      <TableHead>Estado</TableHead>
      <TableHead>Vendidos</TableHead>
    </TableRow>
  </TableHeader>
  <TableBody>
    <TableRow>
      <TableCell className="font-medium">iPhone 15</TableCell>
      <TableCell><Badge variant="success">Activo</Badge></TableCell>
      <TableCell>45/100</TableCell>
    </TableRow>
  </TableBody>
</Table>
```

---

### 4.6 Dialog (Modal)

```tsx
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from '@/components/ui/dialog'

<Dialog>
  <DialogTrigger asChild>
    <Button>Abrir Modal</Button>
  </DialogTrigger>
  <DialogContent className="sm:max-w-md">
    <DialogHeader>
      <DialogTitle>Confirmar Compra</DialogTitle>
      <DialogDescription>
        Estás a punto de comprar 3 boletos por $15
      </DialogDescription>
    </DialogHeader>
    <div className="mt-4">
      <Button className="w-full">Confirmar</Button>
    </div>
  </DialogContent>
</Dialog>
```

---

### 4.7 Toast (Notificaciones)

```tsx
import { useToast } from '@/components/ui/use-toast'

function Component() {
  const { toast } = useToast()

  return (
    <Button
      onClick={() => {
        toast({
          title: "Pago exitoso",
          description: "Tus boletos han sido confirmados",
          variant: "success",
        })
      }}
    >
      Mostrar Toast
    </Button>
  )
}
```

---

### 4.8 Badge

```tsx
import { Badge } from '@/components/ui/badge'

<Badge variant="default">Activo</Badge>
<Badge variant="secondary">Borrador</Badge>
<Badge variant="destructive">Suspendido</Badge>
<Badge variant="outline">Completado</Badge>
```

---

### 4.9 Skeleton (Loading)

```tsx
import { Skeleton } from '@/components/ui/skeleton'

<Card>
  <Skeleton className="h-48 w-full" /> {/* Imagen */}
  <div className="p-4 space-y-2">
    <Skeleton className="h-6 w-3/4" /> {/* Título */}
    <Skeleton className="h-4 w-1/2" /> {/* Descripción */}
  </div>
</Card>
```

---

### 4.10 EmptyState (Sin Datos)

**Componente custom:**
```tsx
interface EmptyStateProps {
  title: string
  description?: string
  icon?: React.ReactNode
  action?: React.ReactNode
}

export function EmptyState({ title, description, icon, action }: EmptyStateProps) {
  return (
    <div className="flex flex-col items-center justify-center py-12 text-center">
      {icon && <div className="mb-4 text-neutral-400">{icon}</div>}
      <h3 className="text-lg font-semibold text-neutral-900">{title}</h3>
      {description && (
        <p className="mt-2 text-sm text-neutral-600 max-w-sm">{description}</p>
      )}
      {action && <div className="mt-6">{action}</div>}
    </div>
  )
}

// Uso
<EmptyState
  icon={<SearchIcon className="w-12 h-12" />}
  title="No hay sorteos activos"
  description="Sé el primero en crear un sorteo"
  action={<Button>Crear Sorteo</Button>}
/>
```

---

## 5. Componentes Específicos del Dominio

### 5.1 RaffleCard (Preview de Sorteo)

```tsx
interface RaffleCardProps {
  id: number
  title: string
  imageUrl: string
  price: number
  soldPercentage: number
  drawDate: string
  status: 'active' | 'completed' | 'suspended'
}

export function RaffleCard({ id, title, imageUrl, price, soldPercentage, drawDate, status }: RaffleCardProps) {
  return (
    <Card className="overflow-hidden hover:shadow-lg transition-shadow">
      <div className="relative">
        <img src={imageUrl} alt={title} className="w-full h-48 object-cover" />
        <Badge
          variant={status === 'active' ? 'default' : 'secondary'}
          className="absolute top-2 right-2"
        >
          {status}
        </Badge>
      </div>
      <CardContent className="p-4">
        <h3 className="font-semibold text-lg text-neutral-900 truncate">{title}</h3>
        <div className="mt-2 flex items-center justify-between">
          <span className="text-2xl font-bold text-primary-600">${price}</span>
          <span className="text-sm text-neutral-600">por boleto</span>
        </div>
        <div className="mt-3">
          <div className="flex justify-between text-xs text-neutral-600 mb-1">
            <span>Vendidos</span>
            <span>{soldPercentage}%</span>
          </div>
          <div className="h-2 bg-neutral-200 rounded-full overflow-hidden">
            <div
              className="h-full bg-primary-500 transition-all"
              style={{ width: `${soldPercentage}%` }}
            />
          </div>
        </div>
        <div className="mt-3 text-xs text-neutral-600">
          Sorteo: {new Date(drawDate).toLocaleDateString('es-CR')}
        </div>
      </CardContent>
      <CardFooter className="p-4 pt-0">
        <Button variant="default" className="w-full" asChild>
          <Link to={`/raffles/${id}`}>Ver Detalles</Link>
        </Button>
      </CardFooter>
    </Card>
  )
}
```

---

### 5.2 NumberGrid (Selección de Números)

```tsx
interface NumberGridProps {
  raffleId: number
  totalNumbers: number // 100 (00-99)
  availableNumbers: string[]
  selectedNumbers: string[]
  onSelect: (number: string) => void
}

export function NumberGrid({ totalNumbers, availableNumbers, selectedNumbers, onSelect }: NumberGridProps) {
  const numbers = Array.from({ length: totalNumbers }, (_, i) =>
    i.toString().padStart(2, '0')
  )

  return (
    <div className="grid grid-cols-10 gap-2">
      {numbers.map((num) => {
        const isAvailable = availableNumbers.includes(num)
        const isSelected = selectedNumbers.includes(num)

        return (
          <button
            key={num}
            onClick={() => isAvailable && onSelect(num)}
            disabled={!isAvailable}
            className={cn(
              'aspect-square rounded-lg font-mono font-semibold text-sm transition-all',
              'focus:outline-none focus:ring-2 focus:ring-primary-500',
              isSelected && 'bg-primary-500 text-white shadow-md scale-105',
              !isSelected && isAvailable && 'bg-white border-2 border-neutral-300 hover:border-primary-400',
              !isAvailable && 'bg-neutral-100 text-neutral-400 cursor-not-allowed'
            )}
          >
            {num}
          </button>
        )
      })}
    </div>
  )
}
```

---

### 5.3 OrderSummary (Resumen de Orden)

```tsx
interface OrderSummaryProps {
  selectedNumbers: string[]
  pricePerNumber: number
  platformFee: number
  total: number
}

export function OrderSummary({ selectedNumbers, pricePerNumber, platformFee, total }: OrderSummaryProps) {
  return (
    <Card>
      <CardHeader>
        <CardTitle>Resumen de Compra</CardTitle>
      </CardHeader>
      <CardContent className="space-y-3">
        <div className="flex justify-between text-sm">
          <span className="text-neutral-600">Números seleccionados</span>
          <span className="font-medium">{selectedNumbers.length}</span>
        </div>
        <div className="flex justify-between text-sm">
          <span className="text-neutral-600">Precio por boleto</span>
          <span className="font-medium">${pricePerNumber.toFixed(2)}</span>
        </div>
        <div className="flex justify-between text-sm">
          <span className="text-neutral-600">Comisión de plataforma</span>
          <span className="font-medium">${platformFee.toFixed(2)}</span>
        </div>
        <Separator />
        <div className="flex justify-between text-lg font-bold">
          <span>Total</span>
          <span className="text-primary-600">${total.toFixed(2)}</span>
        </div>
        <div className="pt-3">
          <Button className="w-full" size="lg">
            Proceder al Pago
          </Button>
        </div>
      </CardContent>
    </Card>
  )
}
```

---

## 6. Estados de Componentes

### 6.1 Hover

- Buttons: `brightness(90%)` + `shadow-md`
- Cards: `shadow-lg`
- Links: `underline` + `text-primary-700`

### 6.2 Focus

- Inputs: `ring-2 ring-primary-500`
- Buttons: `ring-2 ring-primary-500 ring-offset-2`

### 6.3 Active

- Buttons: `brightness(85%)` + `scale(0.98)`

### 6.4 Disabled

- Opacity: `0.5`
- Cursor: `not-allowed`
- No hover effects

### 6.5 Loading

- Skeleton loaders
- Spinner en buttons: `<Loader2 className="animate-spin" />`

---

## 7. Animaciones

**Transiciones suaves:**
```css
.transition-all {
  transition: all 150ms cubic-bezier(0.4, 0, 0.2, 1);
}
```

**Animaciones específicas:**
```tsx
// Fade in
<div className="animate-in fade-in duration-300">...</div>

// Slide from bottom
<div className="animate-in slide-in-from-bottom-4 duration-500">...</div>

// Scale on hover
<Card className="hover:scale-105 transition-transform">...</Card>
```

---

## 8. Accesibilidad (WCAG AA)

### 8.1 Contraste de Colores

**Mínimos requeridos:**
- Texto normal (16px): Ratio 4.5:1
- Texto grande (18px+): Ratio 3:1
- Componentes interactivos: Ratio 3:1

**Verificación:**
- Usar herramientas como [WebAIM Contrast Checker](https://webaim.org/resources/contrastchecker/)

### 8.2 Keyboard Navigation

- Todos los elementos interactivos deben ser accesibles con `Tab`
- Focus visible: `focus:ring-2`
- Escape cierra modals
- Enter/Space activan buttons

### 8.3 ARIA Labels

```tsx
<Button aria-label="Cerrar modal">
  <XIcon className="w-4 h-4" />
</Button>

<Input
  type="text"
  aria-invalid={!!error}
  aria-describedby="error-message"
/>
{error && <p id="error-message" className="text-error">{error}</p>}
```

### 8.4 Screen Readers

- Usar semantic HTML (`<nav>`, `<main>`, `<aside>`)
- Textos alternativos en imágenes
- `aria-live` para notificaciones dinámicas

---

## 9. Dark Mode

**Toggle:**
```tsx
import { Moon, Sun } from 'lucide-react'
import { useTheme } from 'next-themes'

function ThemeToggle() {
  const { theme, setTheme } = useTheme()

  return (
    <Button
      variant="ghost"
      size="sm"
      onClick={() => setTheme(theme === 'dark' ? 'light' : 'dark')}
    >
      <Sun className="rotate-0 scale-100 transition-all dark:-rotate-90 dark:scale-0" />
      <Moon className="absolute rotate-90 scale-0 transition-all dark:rotate-0 dark:scale-100" />
    </Button>
  )
}
```

**Persistencia:**
- Guardar preferencia en `localStorage`
- Respetar preferencia del sistema: `prefers-color-scheme`

---

## 10. Responsive Design

**Mobile First:**
- Diseñar primero para móvil (320px+)
- Usar breakpoints progresivos

**Ejemplo:**
```tsx
<div className="
  px-4 py-8           // Mobile
  md:px-8 md:py-12    // Tablet
  lg:px-16 lg:py-16   // Desktop
">
  <h1 className="text-2xl md:text-4xl lg:text-5xl">
    Título Responsivo
  </h1>
</div>
```

---

## 11. Iconos

**Librería:** Lucide React (fork de Feather Icons)

```tsx
import { Search, ShoppingCart, User, Bell, TrendingUp } from 'lucide-react'

<Button>
  <ShoppingCart className="w-4 h-4 mr-2" />
  Comprar
</Button>
```

**Tamaños estándar:**
- `w-4 h-4` (16px) - Inline en buttons
- `w-5 h-5` (20px) - En cards
- `w-6 h-6` (24px) - Headers
- `w-12 h-12` (48px) - Empty states

---

## 12. Imágenes

**Optimización:**
- WebP con fallback a JPEG
- Lazy loading: `loading="lazy"`
- Responsive images: `srcset`

**Aspect Ratios:**
- Preview de sorteos: `16:9`
- Avatares: `1:1`
- Banners: `21:9`

```tsx
<img
  src="/raffle.webp"
  alt="iPhone 15 Pro"
  className="w-full h-48 object-cover rounded-lg"
  loading="lazy"
/>
```

---

## 13. Checklist de Implementación

- [ ] Configurar Tailwind con design tokens
- [ ] Instalar shadcn/ui components
- [ ] Crear componentes custom (RaffleCard, NumberGrid, OrderSummary)
- [ ] Implementar dark mode con next-themes
- [ ] Agregar Lucide React para iconos
- [ ] Configurar font loading (Inter)
- [ ] Verificar contraste de colores (WCAG AA)
- [ ] Testear keyboard navigation
- [ ] Validar responsive en dispositivos reales
- [ ] Documentar componentes en Storybook (opcional)

---

**Ver también:**
- [Tailwind CSS Docs](https://tailwindcss.com/docs)
- [shadcn/ui Components](https://ui.shadcn.com)
- [WCAG Guidelines](https://www.w3.org/WAI/WCAG21/quickref/)
