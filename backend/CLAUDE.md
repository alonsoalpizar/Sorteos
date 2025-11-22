# Reglas Estrictas para Claude Code - Backend Sorteos

## 🚀 COMPILACION Y DEPLOY RAPIDO

### Backend (Go):
```bash
cd /opt/Sorteos/backend
sudo systemctl stop sorteos-api && \
go build -o sorteos-api ./cmd/api && \
sudo systemctl start sorteos-api
```

### Frontend (Vite/React):
```bash
cd /opt/Sorteos/frontend
npm run build
```
**Nota:** El backend Go sirve el frontend directamente desde `/opt/Sorteos/frontend/dist/` via symlink. No es necesario copiar archivos a ningún otro lugar.

### Todo junto (Backend + Frontend):
```bash
# Frontend (primero para que esté listo cuando el backend reinicie)
cd /opt/Sorteos/frontend && npm run build

# Backend
cd /opt/Sorteos/backend && sudo systemctl stop sorteos-api && \
go build -o sorteos-api ./cmd/api && sudo systemctl start sorteos-api
```

---

## 🚨 REGLA #1: UN SOLO BINARIO OFICIAL

**NUNCA compilar o copiar binarios en ubicaciones temporales como `/tmp/`**

### Ubicacion Oficial del Binario:
```
/opt/Sorteos/backend/sorteos-api
```

### Servicio Systemd:
```
/etc/systemd/system/sorteos-api.service
ExecStart=/opt/Sorteos/backend/sorteos-api
```

### Proceso de Compilacion Oficial:

```bash
cd /opt/Sorteos/backend
sudo systemctl stop sorteos-api
go build -o sorteos-api ./cmd/api
sudo systemctl start sorteos-api
```

**Nota:** Se compila directamente en `sorteos-api` (ubicación de producción). No se usa carpeta `bin/` intermedia.

### Verificar Deploy:
```bash
sudo systemctl status sorteos-api
curl http://localhost:8080/health
```

### ❌ PROHIBIDO:

- ❌ Compilar en `/tmp/`
- ❌ Crear binarios con nombres diferentes (api-test, api-backup, etc.)
- ❌ Copiar binarios a ubicaciones temporales
- ❌ Mantener múltiples versiones del binario

### ✅ PERMITIDO:

- ✅ Compilar directamente: `go build -o sorteos-api ./cmd/api`
- ✅ Usar `make build` si se prefiere (actualizar Makefile para compilar directo)
- ✅ Crear backup temporal SOLO si es necesario:
  ```bash
  cp sorteos-api sorteos-api.backup-$(date +%Y%m%d-%H%M%S)
  ```
- ✅ Eliminar backups después de verificar que la nueva versión funciona

## 🏗️ Estructura de Compilación

### Makefile:
```makefile
build:
	go build -o sorteos-api ./cmd/api
```

**Nota:** Compilar TODO el paquete `./cmd/api`, NO solo `cmd/api/main.go`

### Comandos Disponibles:
```bash
make help      # Ver todos los comandos
make build     # Compilar binario directo a producción
make run       # Ejecutar en desarrollo (go run ./cmd/api)
make test      # Ejecutar tests
```

## 📁 Estructura de Directorios

```
/opt/Sorteos/backend/
├── sorteos-api              # Binario en producción (usado por systemd)
├── cmd/api/                 # Código fuente de la aplicación
│   ├── main.go
│   ├── admin_routes_v2.go
│   ├── routes.go
│   └── ...
├── frontend/                # Symlink a ../frontend
├── Makefile                 # Build script oficial
└── CLAUDE.md               # Este archivo
```

## 🔍 Verificación del Servicio

```bash
# Ver status
sudo systemctl status sorteos-api

# Ver logs
sudo journalctl -u sorteos-api -f

# Verificar binario en uso
ps aux | grep sorteos-api

# Ver qué binario está corriendo
sudo lsof -p $(pgrep sorteos-api) | grep sorteos-api
```

## 📝 Checklist de Actualización

Cuando se actualice el backend:

- [ ] `cd /opt/Sorteos/backend`
- [ ] `git pull` (si aplica)
- [ ] `sudo systemctl stop sorteos-api`
- [ ] `go build -o sorteos-api ./cmd/api`
- [ ] `sudo systemctl start sorteos-api`
- [ ] `sudo systemctl status sorteos-api` (verificar que inicia)
- [ ] `curl http://localhost:8080/health` (verificar respuesta)

## ⚠️ Resolución de Problemas

Si el servicio no inicia:

```bash
# Ver error específico
sudo journalctl -u sorteos-api -n 50 --no-pager

# Ejecutar binario directamente para ver error completo
./sorteos-api

# Verificar permisos
ls -lah sorteos-api
# Debe ser: -rwxr-xr-x root root
```

## 🎯 Estado Actual

**Endpoints Admin:** 52/52 (100%) ✅

**Distribución:**
- Categories: 5 endpoints
- Config: 3 endpoints
- Settlements: 7 endpoints
- Users: 6 endpoints
- Organizers: 5 endpoints
- Payments: 4 endpoints
- Raffles: 6 endpoints
- Notifications: 5 endpoints
- Reports: 4 endpoints
- System: 6 endpoints
- Audit: 1 endpoint

**Última actualización:** 2025-11-18
**Binario:** 27MB (compilado con Go 1.22+)
