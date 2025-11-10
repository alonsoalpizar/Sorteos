# Sorteos Platform - Frontend

Aplicación React + TypeScript construida con Vite para la plataforma de sorteos.

## Tecnologías

- **React 18** - Biblioteca UI
- **TypeScript** - Tipado estático
- **Vite 5** - Build tool y dev server
- **Tailwind CSS** - Estilos utility-first
- **TanStack Query** - Manejo de estado del servidor
- **Zustand** - Estado global de cliente
- **React Hook Form + Zod** - Formularios y validación
- **Axios** - Cliente HTTP con interceptors
- **React Router v6** - Enrutamiento

## Estructura del Proyecto

```
src/
├── components/ui/       # Componentes reutilizables (Button, Input, Card, etc.)
├── features/           # Módulos por funcionalidad
│   ├── auth/          # Autenticación (login, register, verify)
│   └── dashboard/     # Dashboard principal
├── lib/               # Utilidades y configuraciones
│   ├── api.ts        # Cliente Axios con interceptors
│   ├── utils.ts      # Funciones auxiliares (cn, formatters)
│   └── queryClient.ts # Configuración React Query
├── store/            # Stores de Zustand
│   └── authStore.ts  # Estado de autenticación
├── hooks/            # Custom hooks
│   └── useAuth.ts    # Hooks de autenticación
├── api/              # Clientes API
│   └── auth.ts       # Endpoints de autenticación
├── types/            # Definiciones TypeScript
│   └── auth.ts       # Tipos de autenticación
├── App.tsx           # Componente raíz con routing
├── main.tsx          # Entry point
└── index.css         # Estilos globales y Tailwind
```

## Paleta de Colores Aprobada

⚠️ **IMPORTANTE**: Este proyecto usa ÚNICAMENTE la siguiente paleta:

### Colores Permitidos ✅
- **Primary**: Blue #3B82F6 (confianza, profesionalismo)
- **Secondary**: Slate #64748B (elegancia corporativa)
- **Success**: Green #10B981 (confirmaciones)
- **Warning**: Amber #F59E0B (alertas)
- **Destructive**: Red #EF4444 (errores)

### Colores PROHIBIDOS ❌
- Purple (#8B5CF6)
- Pink (#EC4899)
- Magenta
- Fuchsia
- Violet
- Rainbow gradients

## Instalación

```bash
# Instalar dependencias
npm install

# o con yarn
yarn install
```

## Scripts Disponibles

```bash
# Desarrollo (puerto 5173)
npm run dev

# Build de producción
npm run build

# Preview del build
npm run preview

# Type checking
npm run type-check

# Lint
npm run lint
```

## Variables de Entorno

Crear archivo `.env` en la raíz del frontend:

```env
# URL del backend API (opcional, usa proxy por defecto)
VITE_API_URL=http://localhost:8080/api
```

## Proxy de Desarrollo

El proyecto está configurado para hacer proxy de `/api` hacia el backend en `localhost:8080`.

Esto se configura en `vite.config.ts`:

```typescript
server: {
  port: 5173,
  proxy: {
    '/api': {
      target: 'http://localhost:8080',
      changeOrigin: true,
    },
  },
}
```

## Páginas Implementadas

### Autenticación

1. **Register** (`/register`)
   - Formulario completo con validación Zod
   - Campos: nombre, apellido, email, teléfono, contraseña
   - Checkboxes GDPR (términos, privacidad, marketing)
   - Validación de contraseña (12+ caracteres, mayúsculas, minúsculas, números, símbolos)

2. **Login** (`/login`)
   - Email y contraseña
   - Manejo de errores
   - Redirección automática si ya está autenticado

3. **Verify Email** (`/verify-email`)
   - Código de 6 dígitos
   - Auto-login después de verificación
   - Mensaje de expiración (15 minutos)

4. **Dashboard** (`/dashboard`)
   - Área protegida (requiere autenticación)
   - Muestra información del usuario
   - Nivel KYC, rol, estado

## Características del Cliente API

### Gestión de Tokens

- **Access Token**: Corta duración (15 min), enviado en header Authorization
- **Refresh Token**: Larga duración (7 días), almacenado en localStorage
- **Refresh Automático**: Cuando access token expira, se renueva automáticamente
- **Logout Automático**: Si refresh falla, redirige a login

### Interceptors de Axios

```typescript
// Request interceptor
// Agrega Bearer token a todas las peticiones

// Response interceptor
// Maneja errores 401
// Refresca tokens automáticamente
// Reintenta peticiones fallidas
```

## Componentes UI

### Button
- Variantes: default, destructive, outline, secondary, ghost, link
- Tamaños: default, sm, lg, icon
- Loading state con spinner

### Input
- Variante con mensaje de error integrado
- Soporte para todos los tipos HTML

### Label
- Indicador de requerido (*)
- Accesibilidad mejorada

### Card
- Card, CardHeader, CardTitle, CardDescription, CardContent, CardFooter
- Composable components

### Alert
- Variantes: default, success, warning, destructive, info
- Con título y descripción

### Badge
- Para mostrar estados (KYC level, user status, role)
- Variantes con colores aprobados

## Store de Autenticación (Zustand)

```typescript
const { user, isAuthenticated } = useAuthStore();

// Actions
setUser(user);
setAuth(user, accessToken, refreshToken);
logout();

// Helpers
hasMinimumKYC('email_verified');
isAdmin();
isEmailVerified();
```

## Hooks de Autenticación

```typescript
// Mutations
const { mutate } = useRegister();
const { mutate } = useLogin();
const { mutate } = useVerifyEmail();
const { mutate } = useLogout();

// Queries
const { data: user } = useCurrentUser();

// State
const isAuthenticated = useIsAuthenticated();
const user = useUser();
const isAdmin = useIsAdmin();
```

## Protected Routes

```tsx
<Route
  path="/dashboard"
  element={
    <ProtectedRoute requireEmailVerification>
      <DashboardPage />
    </ProtectedRoute>
  }
/>
```

## Validaciones

### Email
- Formato válido
- Único en sistema (validado en backend)

### Password
- Mínimo 12 caracteres
- Al menos 1 mayúscula
- Al menos 1 minúscula
- Al menos 1 número
- Al menos 1 símbolo

### Cédula (opcional)
- 7-10 dígitos
- Solo números

### Teléfono (opcional)
- Formato E.164 (+573001234567)

## Dark Mode

El proyecto tiene soporte completo para dark mode usando Tailwind CSS:

- Variables CSS definidas en `index.css`
- Clase `.dark` en root para activar
- Todos los componentes adaptados

## Próximos Pasos

Según el roadmap (Sprint 3-4):

1. **Gestión de Sorteos**
   - Crear sorteo (formulario)
   - Listar sorteos (grid/tabla)
   - Ver detalle de sorteo
   - Editar/eliminar sorteo (admin)

2. **Sistema de Participación**
   - Reservar números
   - Confirmar pago
   - Ver mis participaciones

3. **Perfil de Usuario**
   - Editar perfil
   - Cambiar contraseña
   - Completar KYC (subir cédula)

## Deployment

```bash
# Build optimizado
npm run build

# Los archivos estáticos estarán en dist/
# Pueden servirse con cualquier servidor estático (nginx, vercel, netlify, etc.)
```

## Soporte

Para más información, revisar la documentación en `/opt/Sorteos/Documentacion/`.
