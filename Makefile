# Che, este es nuestro Makefile re copado para el classifier
# Con todos los chiches para development y production

# Mark all non-file targets as PHONY
.PHONY: all build run test clean lint format swagger help dev-tools deps docker-* db-* debug* pprof-* coverage bench

# Variables, customizalas si queres che
BINARY_NAME=classifier
BUILD_DIR=./bin
COVERAGE_DIR=./coverage
SWAGGER_DIR=./docs/swagger
DOCKER_COMPOSE=docker compose
GO=go

# Build flags
BUILD_FLAGS=-trimpath -ldflags="-s -w"
DEV_FLAGS=-race -gcflags="all=-N -l"

# Default target when just running 'make'
.DEFAULT_GOAL := help

help: ## Che, mostra todos los comandos disponibles
	@echo 'Dale, estos son todos los comandos que podes usar:'
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

# Main commands
all: clean deps lint test build ## La posta: hace todo el build pipeline

build: ## Build del binario nomás
	@echo "📦 Buildeando la app..."
	mkdir -p $(BUILD_DIR)
	$(GO) build $(BUILD_FLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/web

build-dev: ## Build para development con race detection y debugging
	@echo "🔧 Buildeando para development..."
	mkdir -p $(BUILD_DIR)
	$(GO) build $(DEV_FLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/web

run: ## Correr en modo development
	@echo "🚀 Arrancando la app..."
	$(GO) run $(DEV_FLAGS) ./cmd/web

run-prod: build ## Correr en modo production
	@echo "🚀 Arrancando en producción..."
	$(BUILD_DIR)/$(BINARY_NAME)

# Test commands
test: ## Correr todos los tests
	@echo "🧪 Corriendo tests..."
	$(GO) test -race -count=1 ./...

test-verbose: ## Tests con output verboso
	@echo "🔍 Corriendo tests con todos los detalles..."
	$(GO) test -v -race -count=1 ./...

coverage: ## Generar reporte de coverage
	@echo "📊 Generando coverage report..."
	mkdir -p $(COVERAGE_DIR)
	$(GO) test -coverprofile=$(COVERAGE_DIR)/coverage.out ./...
	$(GO) tool cover -html=$(COVERAGE_DIR)/coverage.out -o $(COVERAGE_DIR)/coverage.html
	@echo "📊 Coverage report generado en $(COVERAGE_DIR)/coverage.html"

bench: ## Correr benchmarks
	@echo "🏃 Corriendo benchmarks..."
	$(GO) test -bench=. -benchmem ./...

# Code quality
lint: ## Correr el linter
	@echo "🔍 Chequeando el código..."
	golangci-lint run --fix

format: ## Formatear el código
	@echo "✨ Formateando el código..."
	find . -name '*.go' -not -path "./vendor/*" -exec gofmt -s -w {} \;
	find . -name '*.go' -not -path "./vendor/*" -exec goimports -w {} \;

# Documentation
swagger: ## Generar docs de Swagger
	@echo "📚 Generando documentación API..."
	mkdir -p $(SWAGGER_DIR)
	swagger generate spec -o $(SWAGGER_DIR)/swagger.json --scan-models
	@echo "📚 Swagger docs generados en $(SWAGGER_DIR)/swagger.json"

# Database
db-init: ## Inicializar la base de datos
	@echo "🗃️ Inicializando la base de datos..."
	mysql -u appuser -p classifiersdb < init.sql

db-reset: ## Reset total de la base de datos
	@echo "🗑️ Reseteando la base de datos..."
	mysql -u appuser -p -e "DROP DATABASE IF EXISTS classifiersdb; CREATE DATABASE classifiersdb;"
	make db-init

# Docker commands
docker-build: ## Buildear imagen de Docker
	@echo "🐳 Buildeando imagen Docker..."
	$(DOCKER_COMPOSE) build

docker-up: ## Levantar todos los servicios
	@echo "🐳 Levantando servicios..."
	$(DOCKER_COMPOSE) up -d

docker-down: ## Bajar todos los servicios
	@echo "🐳 Bajando servicios..."
	$(DOCKER_COMPOSE) down

docker-logs: ## Ver logs de los containers
	@echo "📝 Mostrando logs..."
	$(DOCKER_COMPOSE) logs -f

# Debug tools
debug: build-dev ## Arrancar con el debugger
	@echo "🔧 Iniciando debugger..."
	dlv exec $(BUILD_DIR)/$(BINARY_NAME)

debug-test: ## Debug de tests
	@echo "🔧 Debuggeando tests..."
	dlv test ./...

# Performance analysis
pprof-cpu: ## CPU profiling
	@echo "📈 Generando CPU profile..."
	$(GO) test -cpuprofile=$(COVERAGE_DIR)/cpu.prof -bench=. ./...
	$(GO) tool pprof $(COVERAGE_DIR)/cpu.prof

pprof-mem: ## Memory profiling
	@echo "📊 Generando memory profile..."
	$(GO) test -memprofile=$(COVERAGE_DIR)/mem.prof -bench=. ./...
	$(GO) tool pprof $(COVERAGE_DIR)/mem.prof

pprof-trace: ## Execution tracing
	@echo "🔍 Generando execution trace..."
	$(GO) test -trace=$(COVERAGE_DIR)/trace.out -bench=. ./...
	$(GO) tool trace $(COVERAGE_DIR)/trace.out

# Development setup
dev-tools: ## Instalar herramientas de desarrollo
	@echo "🔧 Instalando herramientas..."
	$(GO) install golang.org/x/tools/cmd/goimports@latest
	$(GO) install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	$(GO) install github.com/go-delve/delve/cmd/dlv@latest
	$(GO) install github.com/swaggo/swag/cmd/swag@latest
	$(GO) install github.com/go-swagger/go-swagger/cmd/swagger@latest

# Dependencies
deps: ## Instalar dependencias del proyecto
	@echo "📦 Instalando dependencias..."
	$(GO) mod download
	$(GO) mod tidy

# Clean up
clean: ## Limpiar binarios y archivos temporales
	@echo "🧹 Limpiando todo..."
	rm -rf $(BUILD_DIR)
	rm -rf $(COVERAGE_DIR)
	rm -rf $(SWAGGER_DIR)
	$(GO) clean
	$(DOCKER_COMPOSE) down -v

# Misc
check-updates: ## Checkear updates de dependencias
	@echo "🔍 Buscando actualizaciones..."
	$(GO) list -u -m all

generate: ## Correr go generate
	@echo "🔨 Generando código..."
	$(GO) generate ./...

version: ## Mostrar versión de Go y dependencias
	@echo "ℹ️ Versiones:"
	@$(GO) version
	@echo "Módulos:"
	@$(GO) list -m all