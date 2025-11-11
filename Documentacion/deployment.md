# Documentación de Deployment - Sorteos Platform

## Información General

**Dominio Principal**: https://sorteos.club
**Servidor**: 62.171.188.255 (VPS)
**Última Actualización**: 2025-11-10

## Arquitectura de Deployment

```
Internet
    ↓
Nginx (Puerto 80/443) - SSL/TLS Termination
    ↓
Reverse Proxy
    ↓
Docker Container (sorteos-api:8080)
    ├── Go Backend (Gin)
    └── React Frontend (SPA)
    ↓
├── PostgreSQL (sorteos-postgres:5432)
└── Redis (sorteos-redis:6379)
```

## Componentes

### 1. Nginx (Reverse Proxy)

**Archivo de Configuración**: `/etc/nginx/sites-available/sorteos`

**Características**:
- SSL/TLS con Let's Encrypt (certificados automáticos)
- Redirección HTTP → HTTPS
- Redirección www → non-www
- Compresión gzip
- Headers de seguridad (HSTS, X-Frame-Options, etc.)
- Proxy a contenedor Docker en puerto 8080
- Timeout de 60s para operaciones de pago
- Max upload size: 10MB

**Comandos útiles**:
```bash
# Verificar configuración
sudo nginx -t

# Recargar configuración
sudo systemctl reload nginx

# Ver logs
sudo tail -f /var/log/nginx/sorteos_access.log
sudo tail -f /var/log/nginx/sorteos_error.log
```

### 2. Docker Compose Stack

**Archivo**: `/opt/Sorteos/docker-compose.yml`

**Servicios**:

#### sorteos-api (Backend + Frontend)
- **Imagen**: Construida desde `backend/Dockerfile` (multi-stage)
- **Puerto**: 8080 (interno), 80/443 (externo vía Nginx)
- **Variables de entorno**: Ver `.env` en backend
- **Volúmenes**:
  - `./backend/uploads:/app/uploads` (persistente)
- **Health check**: `/health` cada 30s

#### sorteos-postgres (Database)
- **Imagen**: postgres:15-alpine
- **Puerto**: 5432
- **Credenciales**: Ver `.env`
- **Volumen**: `postgres_data` (persistente)
- **Migraciones**: Auto-aplicadas en startup

#### sorteos-redis (Cache)
- **Imagen**: redis:7-alpine
- **Puerto**: 6379
- **Volumen**: `redis_data` (persistente)

**Comandos útiles**:
```bash
cd /opt/Sorteos

# Ver estado de servicios
docker compose ps

# Ver logs
docker compose logs -f api
docker compose logs -f postgres
docker compose logs -f redis

# Rebuild y redeploy
docker compose build api
docker compose up -d api

# Restart todo el stack
docker compose restart

# Stop todo
docker compose down

# Ver recursos
docker stats
```

### 3. SSL/TLS (Let's Encrypt)

**Certificados**: `/etc/letsencrypt/live/sorteos.club/`

**Renovación Automática**: Certbot configurado con cron

**Comandos útiles**:
```bash
# Ver estado de certificados
sudo certbot certificates

# Renovar manualmente (dry-run)
sudo certbot renew --dry-run

# Renovar forzado
sudo certbot renew --force-renewal

# Ver logs de certbot
sudo tail -f /var/log/letsencrypt/letsencrypt.log
```

## URLs y Endpoints

### Frontend (Público)
- **Home/Listado**: https://sorteos.club/raffles
- **Detalle de sorteo**: https://sorteos.club/raffles/:id
- **Login**: https://sorteos.club/login
- **Registro**: https://sorteos.club/register
- **Crear sorteo**: https://sorteos.club/raffles/create (protegido)

### API REST
- **Base URL**: https://sorteos.club/api/v1
- **Health**: https://sorteos.club/health
- **Readiness**: https://sorteos.club/ready

#### Endpoints Disponibles

**Autenticación** (`/api/v1/auth/`)
- `POST /register` - Registro de usuario
- `POST /login` - Iniciar sesión
- `POST /refresh` - Refrescar token
- `POST /verify-email` - Verificar email

**Sorteos** (`/api/v1/raffles/`)
- `GET /` - Listar sorteos (público)
- `GET /:id` - Detalle de sorteo (público)
- `POST /` - Crear sorteo (requiere auth + email verificado)
- `PUT /:id` - Actualizar sorteo (requiere ownership)
- `POST /:id/publish` - Publicar sorteo (requiere ownership)
- `DELETE /:id` - Eliminar sorteo (requiere ownership)
- `POST /:id/suspend` - Suspender sorteo (requiere admin)

**Usuario**
- `GET /api/v1/profile` - Perfil del usuario autenticado
- `GET /api/v1/admin/users` - Listar usuarios (solo admin)

## Usuarios de Prueba

Ver archivo completo: [`USUARIOS_PRUEBA.md`](../USUARIOS_PRUEBA.md)

### Usuario Normal
- **Email**: test@sorteos.com
- **Password**: TestPassword123!
- **Permisos**: Crear sorteos, comprar números

### Usuario Admin
- **Email**: admin@sorteos.com
- **Password**: Admin123456!
- **Permisos**: Todos los permisos + administración

## Monitoreo y Logs

### Logs de Aplicación
```bash
# Backend API
docker compose logs -f api

# Ver últimas 100 líneas
docker compose logs --tail=100 api

# Seguir errores específicos
docker compose logs -f api | grep ERROR
```

### Logs de Nginx
```bash
# Access log
sudo tail -f /var/log/nginx/sorteos_access.log

# Error log
sudo tail -f /var/log/nginx/sorteos_error.log

# Analizar requests por status code
sudo awk '{print $9}' /var/log/nginx/sorteos_access.log | sort | uniq -c | sort -rn
```

### Logs de PostgreSQL
```bash
docker compose logs -f postgres

# Conexiones activas
docker compose exec postgres psql -U sorteos_user -d sorteos_db -c "SELECT * FROM pg_stat_activity;"
```

### Logs de Redis
```bash
docker compose logs -f redis

# Conectar a Redis CLI
docker compose exec redis redis-cli
```

## Backup y Restore

### Base de Datos PostgreSQL

**Backup Manual**:
```bash
# Crear backup
docker compose exec postgres pg_dump -U sorteos_user sorteos_db > backup_$(date +%Y%m%d_%H%M%S).sql

# Con compresión
docker compose exec postgres pg_dump -U sorteos_user sorteos_db | gzip > backup_$(date +%Y%m%d_%H%M%S).sql.gz
```

**Restore**:
```bash
# Desde archivo SQL
docker compose exec -T postgres psql -U sorteos_user sorteos_db < backup.sql

# Desde archivo comprimido
gunzip -c backup.sql.gz | docker compose exec -T postgres psql -U sorteos_user sorteos_db
```

**Backup Automatizado** (cron):
```bash
# Editar crontab
sudo crontab -e

# Agregar línea para backup diario a las 2am
0 2 * * * cd /opt/Sorteos && docker compose exec postgres pg_dump -U sorteos_user sorteos_db | gzip > /opt/backups/sorteos_$(date +\%Y\%m\%d).sql.gz
```

### Archivos Subidos (Uploads)
```bash
# Backup de uploads
tar -czf uploads_backup_$(date +%Y%m%d).tar.gz /opt/Sorteos/backend/uploads/

# Restore
tar -xzf uploads_backup.tar.gz -C /opt/Sorteos/backend/
```

## Troubleshooting

### Problema: Container no inicia
```bash
# Ver logs detallados
docker compose logs api

# Ver estado del container
docker compose ps

# Verificar variables de entorno
docker compose exec api env | grep CONFIG

# Restart forzado
docker compose down && docker compose up -d
```

### Problema: Error de conexión a DB
```bash
# Verificar que PostgreSQL esté corriendo
docker compose ps postgres

# Verificar conectividad desde el container
docker compose exec api ping postgres

# Verificar credenciales
docker compose exec postgres psql -U sorteos_user -d sorteos_db -c "SELECT 1;"
```

### Problema: SSL no funciona
```bash
# Verificar certificados
sudo certbot certificates

# Ver logs de Nginx
sudo tail -f /var/log/nginx/sorteos_error.log

# Test manual
curl -vI https://sorteos.club
```

### Problema: API lenta
```bash
# Ver recursos del container
docker stats sorteos-api

# Ver conexiones activas a DB
docker compose exec postgres psql -U sorteos_user -d sorteos_db -c "SELECT COUNT(*) FROM pg_stat_activity;"

# Verificar Redis
docker compose exec redis redis-cli ping
docker compose exec redis redis-cli INFO stats
```

## Proceso de Deploy de Nuevas Versiones

### 1. Deploy de Backend
```bash
cd /opt/Sorteos

# Pull cambios (si aplica)
git pull origin main

# Rebuild imagen
docker compose build api

# Deploy con zero-downtime
docker compose up -d api

# Verificar logs
docker compose logs -f api

# Rollback si es necesario
docker compose down api
docker compose up -d api
```

### 2. Deploy de Frontend
El frontend está integrado en el contenedor de backend, por lo que sigue el mismo proceso que arriba.

### 3. Migraciones de Base de Datos
```bash
# Las migraciones se aplican automáticamente en startup
# Para aplicar manualmente:

cd /opt/Sorteos/backend

# Crear nueva migración
migrate create -ext sql -dir migrations -seq nombre_migracion

# Aplicar migraciones
docker compose exec api sh -c "migrate -path ./migrations -database 'postgres://sorteos_user:sorteos_password@postgres:5432/sorteos_db?sslmode=disable' up"

# Rollback última migración
docker compose exec api sh -c "migrate -path ./migrations -database 'postgres://sorteos_user:sorteos_password@postgres:5432/sorteos_db?sslmode=disable' down 1"
```

## Seguridad

### Headers de Seguridad Configurados
- ✅ HSTS (Strict-Transport-Security)
- ✅ X-Frame-Options: SAMEORIGIN
- ✅ X-Content-Type-Options: nosniff
- ✅ X-XSS-Protection: 1; mode=block
- ✅ Referrer-Policy: strict-origin-when-cross-origin

### Best Practices Implementadas
- ✅ SSL/TLS con certificados válidos
- ✅ Redirección forzada HTTP → HTTPS
- ✅ Container runs as non-root user
- ✅ Secrets en variables de entorno (no en código)
- ✅ Rate limiting implementado (Redis)
- ✅ CORS configurado correctamente
- ✅ JWT con refresh tokens
- ✅ SQL injection prevention (GORM prepared statements)
- ✅ Soft delete en lugar de eliminación física

### Pendientes de Implementar
- ⏳ Firewall (UFW) con whitelist de puertos
- ⏳ Fail2ban para protección contra brute force
- ⏳ Backup automatizado a almacenamiento externo
- ⏳ Monitoring con Prometheus + Grafana
- ⏳ Alertas automáticas (email/slack)
- ⏳ WAF (Web Application Firewall)

## Performance

### Optimizaciones Implementadas
- ✅ Gzip compression (nivel 6)
- ✅ HTTP/2
- ✅ Keep-alive connections
- ✅ Browser caching para assets estáticos (1 año)
- ✅ Redis caching
- ✅ Database connection pooling
- ✅ Frontend code splitting y tree shaking

### Métricas Actuales
- Frontend build size: 393.92 kB (gzipped: 120.77 kB)
- CSS build size: 24.00 kB (gzipped: 4.99 kB)
- Backend binary size: ~15 MB
- Container memory usage: ~50-100 MB
- Response time promedio: <100ms (local), <200ms (API)

## Contacto y Soporte

**Repositorio**: (Agregar URL del repositorio)
**Documentación adicional**: `/opt/Sorteos/Documentacion/`
**Issues**: (Agregar URL de issues)

## Changelog

### 2025-11-10 - Sprint 3-4 Deploy
- ✅ Implementado sistema completo de gestión de sorteos (backend)
- ✅ Implementado frontend React con páginas de listado, detalle y creación
- ✅ Integrado frontend en contenedor Docker con backend
- ✅ Configurado Nginx como reverse proxy
- ✅ Dominio sorteos.club configurado con SSL
- ✅ Usuarios de prueba creados y verificados

### Anterior
- ✅ Sprint 1-2: Infraestructura y autenticación
- ✅ Configuración inicial de PostgreSQL, Redis, Docker Compose
- ✅ Sistema de autenticación con JWT
- ✅ Verificación de email
- ✅ Rate limiting
