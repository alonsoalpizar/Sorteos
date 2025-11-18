# Reglas Estrictas para Claude Code - Backend Sorteos

## 🚨 REGLA #1: UN SOLO BINARIO OFICIAL

**NUNCA compilar o copiar binarios en ubicaciones temporales como `/tmp/`**

### Ubicación Oficial del Binario:
```
/opt/Sorteos/backend/sorteos-api
```

### Proceso de Compilación Oficial:

1. **Usar Makefile:**
   ```bash
   cd /opt/Sorteos/backend
   make build
   ```

2. **Ubicación del Binario Compilado:**
   ```
   /opt/Sorteos/backend/bin/sorteos-api
   ```

3. **Para Actualizar Producción:**
   ```bash
   # Detener servicio
   sudo systemctl stop sorteos-api

   # Copiar binario compilado a ubicación oficial
   cp bin/sorteos-api sorteos-api

   # Iniciar servicio
   sudo systemctl start sorteos-api
   ```

### Servicio Systemd:
```
/etc/systemd/system/sorteos-api.service
```

El servicio ejecuta: `/opt/Sorteos/backend/sorteos-api`

### ❌ PROHIBIDO:

- ❌ Compilar en `/tmp/`
- ❌ Crear binarios con nombres diferentes (api-test, api-backup, etc.)
- ❌ Copiar binarios a ubicaciones temporales
- ❌ Mantener múltiples versiones del binario

### ✅ PERMITIDO:

- ✅ Compilar usando `make build` (crea en `bin/sorteos-api`)
- ✅ Copiar desde `bin/sorteos-api` a `sorteos-api` (producción)
- ✅ Crear backup temporal SOLO durante actualización con fecha clara:
  ```bash
  cp sorteos-api sorteos-api.backup-$(date +%Y%m%d-%H%M%S)
  ```
- ✅ Eliminar backups después de verificar que la nueva versión funciona

## 🏗️ Estructura de Compilación

### Makefile Correcto:
```makefile
build:
	go build -o bin/sorteos-api ./cmd/api
```

**Nota:** Compilar TODO el paquete `./cmd/api`, NO solo `cmd/api/main.go`

### Comandos Disponibles:
```bash
make help      # Ver todos los comandos
make build     # Compilar binario
make run       # Ejecutar en desarrollo (go run ./cmd/api)
make test      # Ejecutar tests
make clean     # Limpiar binarios generados
```

## 📁 Estructura de Directorios

```
/opt/Sorteos/backend/
├── bin/                      # Binarios compilados (gitignored)
│   └── sorteos-api          # Binario compilado por make build
├── sorteos-api              # Binario en producción (usado por systemd)
├── cmd/api/                 # Código fuente de la aplicación
│   ├── main.go
│   ├── admin_routes_v2.go
│   ├── routes.go
│   └── ...
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
- [ ] `make build`
- [ ] Verificar compilación exitosa (`ls -lah bin/sorteos-api`)
- [ ] `sudo systemctl stop sorteos-api`
- [ ] `cp bin/sorteos-api sorteos-api`
- [ ] `sudo systemctl start sorteos-api`
- [ ] `sudo systemctl status sorteos-api` (verificar que inicia)
- [ ] `curl http://localhost:8080/health` (verificar respuesta)
- [ ] Eliminar cualquier binario temporal creado

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
