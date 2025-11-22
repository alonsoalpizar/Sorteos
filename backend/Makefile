.PHONY: help run build test test-coverage lint migrate-up migrate-down migrate-create docker-up docker-down clean

# Variables
APP_NAME=sorteos-api
DOCKER_COMPOSE=docker-compose
MIGRATE=migrate
MIGRATIONS_PATH=./migrations
DB_URL=postgres://sorteos_user:sorteos_password@localhost:5432/sorteos_db?sslmode=disable

help: ## Mostrar ayuda
	@echo "Comandos disponibles:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'

run: ## Ejecutar aplicación en desarrollo
	@echo "🚀 Ejecutando $(APP_NAME)..."
	go run ./cmd/api

build: ## Compilar aplicación (directo a producción)
	@echo "🔨 Compilando $(APP_NAME)..."
	go build -o $(APP_NAME) ./cmd/api
	@echo "✅ Binario creado: $(APP_NAME)"

test: ## Ejecutar tests
	@echo "🧪 Ejecutando tests..."
	go test -v -race ./...

test-coverage: ## Ejecutar tests con coverage
	@echo "📊 Ejecutando tests con coverage..."
	go test -v -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "✅ Coverage disponible en coverage.html"

lint: ## Ejecutar linter
	@echo "🔍 Ejecutando golangci-lint..."
	golangci-lint run

tidy: ## Limpiar dependencias
	@echo "🧹 Limpiando dependencias..."
	go mod tidy

deps: ## Descargar dependencias
	@echo "📦 Descargando dependencias..."
	go mod download

migrate-up: ## Aplicar migraciones
	@echo "⬆️  Aplicando migraciones..."
	$(MIGRATE) -path $(MIGRATIONS_PATH) -database "$(DB_URL)" up

migrate-down: ## Revertir última migración
	@echo "⬇️  Revirtiendo última migración..."
	$(MIGRATE) -path $(MIGRATIONS_PATH) -database "$(DB_URL)" down 1

migrate-create: ## Crear nueva migración (usar: make migrate-create name=nombre_migracion)
	@if [ -z "$(name)" ]; then \
		echo "❌ Error: Debes especificar name=nombre_migracion"; \
		exit 1; \
	fi
	@echo "📝 Creando migración: $(name)"
	$(MIGRATE) create -ext sql -dir $(MIGRATIONS_PATH) -seq $(name)

migrate-force: ## Forzar versión de migración (usar: make migrate-force version=1)
	@if [ -z "$(version)" ]; then \
		echo "❌ Error: Debes especificar version=N"; \
		exit 1; \
	fi
	@echo "⚠️  Forzando versión $(version)..."
	$(MIGRATE) -path $(MIGRATIONS_PATH) -database "$(DB_URL)" force $(version)

docker-up: ## Levantar contenedores Docker
	@echo "🐳 Levantando contenedores..."
	$(DOCKER_COMPOSE) up -d
	@echo "✅ Contenedores activos. Esperando 5s..."
	@sleep 5

docker-down: ## Detener contenedores Docker
	@echo "🛑 Deteniendo contenedores..."
	$(DOCKER_COMPOSE) down

docker-logs: ## Ver logs de contenedores
	$(DOCKER_COMPOSE) logs -f

docker-rebuild: ## Reconstruir contenedores
	@echo "🔄 Reconstruyendo contenedores..."
	$(DOCKER_COMPOSE) up -d --build

clean: ## Limpiar archivos generados
	@echo "🧹 Limpiando..."
	rm -f coverage.out coverage.html
	rm -rf tmp/
	@echo "✅ Limpieza completada"

install-tools: ## Instalar herramientas de desarrollo
	@echo "🔧 Instalando herramientas..."
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	@echo "✅ Herramientas instaladas"

dev: docker-up migrate-up run ## Iniciar entorno de desarrollo completo

.DEFAULT_GOAL := help
