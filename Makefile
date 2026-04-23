.PHONY: dev backend frontend build rebuild clean help install sync test lint

# Variables de compilación
PKG_PATH := gogs.io/gogs/internal/conf
BUILD_TIME := $(shell date -u '+%Y-%m-%d %I:%M:%S %Z')
BUILD_COMMIT := $(shell git rev-parse HEAD 2>/dev/null || echo "unknown")
TAGS := cert pam bindata
LDFLAGS := -s -w -X "$(PKG_PATH).BuildTime=$(BUILD_TIME)" -X "$(PKG_PATH).BuildCommit=$(BUILD_COMMIT)"

# Default target
help:
	@echo "Comandos disponibles:"
	@echo "  make dev       - Ejecuta backend y frontend en paralelo"
	@echo "  make backend   - Ejecuta el servidor Ggit (puerto 3005)"
	@echo "  make frontend  - Ejecuta el servidor de desarrollo Svelte (puerto 5555)"
	@echo "  make build     - Compila el binario de Ggit"
	@echo "  make sync      - Sincroniza plantillas de gitignore"
	@echo "  make test      - Ejecuta los tests de Go"
	@echo "  make install   - Instala las dependencias del frontend"

# Sincroniza plantillas de gitignore
sync:
	@echo "Sincronizando plantillas gitignore..."
	@chmod +x scripts/sync_templates.sh && ./scripts/sync_templates.sh

# Limpiar puertos antes de empezar
clean-ports:
	@echo "Limpiando puertos 3005 y 5555..."
	@fuser -k 3005/tcp 5555/tcp || true

# Ejecuta ambos en paralelo
dev: clean-ports
	@echo "Iniciando Ggit Modernized (Desarrollo)..."
	@make -j 2 backend frontend

# Compila y corre el backend
backend: build
	@echo "Arrancando Backend Go en http://localhost:3005 ..."
	./ggit web -p 3005

# Corre el frontend de Svelte
frontend:
	@echo "Arrancando Frontend Svelte en http://localhost:5555 ..."
	@cd frontend && npm run dev -- --port 5555

# Compila el binario
build:
	@echo "Compilando Ggit (Commit: $(BUILD_COMMIT))..."
	@go build -v -tags "$(TAGS)" -ldflags '$(LDFLAGS)' -o ggit ./cmd/ggit

# Fuerza la recompilación completa
rebuild:
	@echo "Forzando recompilación completa (Commit: $(BUILD_COMMIT))..."
	@go build -a -v -tags "$(TAGS)" -ldflags '$(LDFLAGS)' -o ggit ./cmd/ggit

# Instalación de dependencias
install:
	@echo "Instalando dependencias de Node.js..."
	@cd frontend && npm install

# Tests
test:
	@echo "Ejecutando tests..."
	@go test -v ./...

# Limpieza
clean:
	@echo "Limpiando binarios..."
	rm -f ggit
