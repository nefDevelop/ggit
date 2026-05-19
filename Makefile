.PHONY: dev backend frontend build build-swagger build-swagger-cdn rebuild clean help install sync test lint test-e2e kill-e2e clean-ports

# Variables de compilación
PKG_PATH := github.com/nef734/ggit/internal/conf
BUILD_TIME := $(shell date -u '+%Y-%m-%d %I:%M:%S %Z')
BUILD_COMMIT := $(shell git rev-parse HEAD 2>/dev/null || echo "unknown")
TAGS := cert pam bindata go_json
LDFLAGS := -s -w -X "$(PKG_PATH).BuildTime=$(BUILD_TIME)" -X "$(PKG_PATH).BuildCommit=$(BUILD_COMMIT)"

# Default target
help:
	@echo "Comandos disponibles:"
	@echo "  make dev       - Ejecuta backend y frontend en paralelo"
	@echo "  make backend   - Ejecuta el servidor Ggit (puerto 3005)"
	@echo "  make frontend  - Ejecuta el servidor de desarrollo Svelte (puerto 5555)"
	@echo "  make build             - Compila el binario de Ggit"
	@echo "  make build-swagger     - Compila con Swagger embebido"
	@echo "  make build-swagger-cdn - Compila con Swagger vía CDN"
	@echo "  make test      - Ejecuta los tests de Go"
	@echo "  make test-e2e  - Ejecuta los tests de Playwright (E2E)"
	@echo "  make install   - Instala las dependencias"

# Sincroniza plantillas de gitignore
sync:
	@echo "Sincronizando plantillas gitignore..."
	@chmod +x scripts/sync_templates.sh && ./scripts/sync_templates.sh

# Limpiar puertos antes de empezar
clean-ports:
	@echo "Limpiando puertos de desarrollo y tests..."
	@fuser -k 3005/tcp 5555/tcp 3006/tcp 5556/tcp || true

# Ejecuta ambos en paralelo
dev: clean-ports
	@echo "Iniciando Ggit Modernized (Desarrollo)..."
	@make -j 2 backend frontend

# Compila y corre el backend
backend: build
	@echo "Arrancando Backend Go en http://localhost:3005 ..."
	@GGIT_WORK_DIR=$$(pwd) GGIT_SKIP_RUN_USER_CHECK=true ./.bin/ggit web -c conf/app.ini -p 3005

# Corre el frontend de Svelte
frontend:
	@echo "Arrancando Frontend Svelte en http://localhost:5555 ..."
	@cd frontend && npm run dev -- --port 5555

# Compila el binario
build:
	@echo "Compilando Ggit..."
	@mkdir -p .bin
	@go build -v -tags "$(shell echo $(TAGS) | tr ' ' ',')" -ldflags '$(LDFLAGS)' -o .bin/ggit ./cmd/ggit

build-swagger: TAGS += swagger
build-swagger: build

build-swagger-cdn: TAGS += swagger_cdn
build-swagger-cdn: build

# Tests
test:
	@echo "Ejecutando tests..."
	@GGIT_SKIP_RUN_USER_CHECK=true go test -v ./...

# Tests de Playwright (E2E)
test-e2e: clean-ports build
	@echo "Preparando entorno aislado para tests E2E..."
	@mkdir -p data
	@# Crear config temporal para tests
	@cp conf/app.ini conf/app.test.ini
	@sed -i -e 's|/data/data|data/data|g' \
	        -e 's|/data/repositories|ggit-repositories-test|g' \
	        -e 's|data/data/ggit.db|data/ggit.test.db|g' conf/app.test.ini
	@echo "[security]" >> conf/app.test.ini
	@echo "INSTALL_LOCK = true" >> conf/app.test.ini
	@echo "ENABLE_CSRF = false" >> conf/app.test.ini
	@echo "COOKIE_SAMESITE = Lax" >> conf/app.test.ini
	@echo "COOKIE_SECURE = false" >> conf/app.test.ini
	@echo "[auth]" >> conf/app.test.ini
	@echo "REQUIRE_SIGNIN_VIEW = false" >> conf/app.test.ini
	@echo "[server]" >> conf/app.test.ini
	@echo "HTTP_PORT = 3006" >> conf/app.test.ini
	@rm -rf data/gogs.test.db
	@rm -rf ggit-repositories-test
	@mkdir -p ggit-repositories-test
	@rm -rf frontend/tests/.auth/*
	@rm -rf .bin/data # Limpiar rastro de ejecuciones fallidas
	@echo "Asegurando usuario admin para tests..."
	@GGIT_WORK_DIR=$$(pwd) GGIT_SKIP_RUN_USER_CHECK=true ./.bin/ggit admin create-user --name ggit-admin --password admin1234 --email admin@example.com --admin -c conf/app.test.ini
	@echo "Verificando dependencias de Playwright..."
	@cd frontend && npx playwright install
	@echo "Construyendo frontend para producción (test)..."
	@cd frontend && npm run build
	@echo "Levantando entorno para pruebas E2E (Servido por Backend)..."
	@# Levantar backend (Puerto 3006) usando la config de test. Servirá la SPA desde frontend/build.
	@GGIT_WORK_DIR=$$(pwd) GGIT_SKIP_RUN_USER_CHECK=true ./.bin/ggit web -c conf/app.test.ini -p 3006 > .backend.log 2>&1 & echo $$! > .backend.pid
	@echo "Esperando a que el backend responda (http://localhost:3006)..."
	@timeout 60s bash -c 'until curl -s http://localhost:3006 > /dev/null; do sleep 1; done' || (echo "Error: El backend no arrancó" && make kill-e2e && exit 1)
	@echo "Ejecutando Playwright..."
	@cd frontend && npx playwright test --reporter=list || (echo "--- LOGS DE ERROR BACKEND ---" && cat ../.backend.log && cd .. && make kill-e2e && exit 1)
	@rm -f conf/app.test.ini
	@cd .. && make kill-e2e

kill-e2e:
	@echo "Cerrando servidores de prueba..."
	@if [ -f .backend.pid ]; then kill $$(cat .backend.pid) || true; fi
	@if [ -f .frontend.pid ]; then kill $$(cat .frontend.pid) || true; fi
	@rm -f .backend.pid .frontend.pid
	@fuser -k 3006/tcp 5556/tcp || true

# Limpieza
clean:
	@echo "Limpiando binarios..."
	rm -rf .bin
