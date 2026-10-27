.PHONY: dev backend frontend build build-cli check-theme-tokens check-theme-contrast check-ui build-swagger build-swagger-cdn rebuild clean help install sync test test-frontend test-all lint test-e2e kill-e2e clean-ports demo demo-config seed-demo kill-demo

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
	@echo "  make build-cli         - Compila el cliente de la API v2 (ggit-cli)"
	@echo "  make build-swagger     - Compila con Swagger embebido"
	@echo "  make build-swagger-cdn - Compila con Swagger vía CDN"
	@echo "  make check-theme-tokens   - Verifica el contrato de tokens de los temas"
	@echo "  make check-theme-contrast - Verifica el contraste WCAG de cada tema"
	@echo "  make check-ui             - Ejecuta las dos comprobaciones anteriores"
	@echo "  make test      - Ejecuta los tests de Go"
	@echo "  make test-e2e  - Ejecuta los tests de Playwright (E2E)"
	@echo "  make demo      - Levanta una instancia aislada con datos de ejemplo (puerto 3007)"
	@echo "  make kill-demo - Detiene la instancia de demo"
	@echo "  make install   - Instala las dependencias"

# Verifica que todos los temas definen el mismo contrato de tokens.
check-theme-tokens:
	@python3 scripts/check_theme_tokens.py

# Verifica el contraste WCAG de los pares texto/fondo de cada tema.
check-theme-contrast:
	@python3 scripts/check_theme_contrast.py

check-ui: check-theme-tokens check-theme-contrast

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

# Compila el cliente de la API v2 (no arrastra dependencias del servidor)
build-cli:
	@echo "Compilando ggit-cli..."
	@mkdir -p .bin
	@go build -v -ldflags '-s -w' -o .bin/ggit-cli ./cmd/ggit-cli

build-swagger: TAGS += swagger
build-swagger: build

build-swagger-cdn: TAGS += swagger_cdn
build-swagger-cdn: build

# Tests
# -race exposes data races; several packages mutate the shared database handle
# from tests that call t.Parallel, so this is not optional.
# -timeout fails fast instead of blocking until the ten minute default, which is
# what a killed process used to do.
# -shuffle=on varies test order, which is how the shared handle races surfaced.
test:
	@echo "Ejecutando tests de Go..."
	@GGIT_SKIP_RUN_USER_CHECK=true go test -race -timeout 300s -shuffle=on ./...

test-short:
	@echo "Ejecutando tests de Go (sin los marcados como lentos)..."
	@GGIT_SKIP_RUN_USER_CHECK=true go test -short -timeout 300s ./...

test-frontend:
	@echo "Ejecutando tests unitarios del frontend..."
	@cd frontend && npx vitest run

test-all: test test-frontend

# Entorno E2E. La configuracion se genera en un unico sitio porque rehacerla a
# mano es como se acabo ejecutando la suite contra data/ggit.db, la base real.
# e2e-env-up deja el servidor vivo para poder iterar sobre un solo test sin
# relanzar todo, y test-e2e lo usa y lo baja al terminar.

e2e-config:
	@mkdir -p data
	@cp conf/app.ini conf/app.test.ini
	@# Every data path is redirected under data/test-e2e/. The previous patterns
	@# were "data/data/ggit.db" and "data/repositories", but the values in app.ini
	@# are "data/ggit.db" and "data/repositories", so the database substitution
	@# never matched: the suite ran against data/ggit.db, the real one, while
	@# "rm -rf data/ggit.test.db" deleted a file nobody used.
	@sed -i -E \
	        -e 's#^([[:space:]]*(APP_DATA_PATH|PATH|PROVIDER_CONFIG|OBJECTS_PATH|OBJECTS_TEMP_PATH|TEMP_PATH|AVATAR_UPLOAD_PATH|REPOSITORY_AVATAR_UPLOAD_PATH|ROOT)[[:space:]]*=[[:space:]]*)data/#\1data/test-e2e/#' \
	        -e 's#^([[:space:]]*ROOT[[:space:]]*=[[:space:]]*)data/repositories#\1data/test-e2e/repositories#' \
	        conf/app.test.ini
	@# Fail loudly rather than run against the developer's database again. Both
	@# checks matter: the first catches a rewrite that stopped happening, the
	@# second catches a new data path that nobody remembered to add above.
	@grep -qE '^[[:space:]]*PATH[[:space:]]*=[[:space:]]*data/test-e2e/ggit\.db' conf/app.test.ini \
	    || { echo "ERROR: la base de datos de test no quedo aislada"; exit 1; }
	@! grep -qE '^[[:space:]]*(PATH|ROOT|APP_DATA_PATH|PROVIDER_CONFIG|OBJECTS_PATH|OBJECTS_TEMP_PATH|TEMP_PATH|AVATAR_UPLOAD_PATH|REPOSITORY_AVATAR_UPLOAD_PATH)[[:space:]]*=[[:space:]]*data/(ggit\.db|repositories|data|attachments|avatars|repo-avatars|lfs-objects|tmp)([[:space:]]*)$$' conf/app.test.ini \
	    || { echo "ERROR: queda una ruta de datos real en app.test.ini"; grep -nE '^[[:space:]]*(PATH|ROOT|APP_DATA_PATH|PROVIDER_CONFIG|OBJECTS_PATH|OBJECTS_TEMP_PATH|TEMP_PATH|AVATAR_UPLOAD_PATH|REPOSITORY_AVATAR_UPLOAD_PATH)[[:space:]]*=' conf/app.test.ini; exit 1; }
	@printf '[security]\nINSTALL_LOCK = true\nENABLE_CSRF = false\nCOOKIE_SAMESITE = Lax\nCOOKIE_SECURE = false\n[auth]\nREQUIRE_SIGNIN_VIEW = false\n[server]\nHTTP_PORT = 3006\n' >> conf/app.test.ini

e2e-env-up: clean-ports build e2e-config
	@echo "Preparando entorno aislado para tests E2E..."
	@rm -rf data/test-e2e
	@mkdir -p data/test-e2e/repositories
	@rm -rf frontend/tests/.auth/*
	@rm -rf frontend/test-results
	@echo "Asegurando usuario admin para tests..."
	@GGIT_WORK_DIR=$$(pwd) GGIT_SKIP_RUN_USER_CHECK=true ./.bin/ggit admin create-user --name ggit-admin --password admin1234 --email admin@example.com --admin -c conf/app.test.ini
	@echo "Asegurando usuario regular para tests..."
	@GGIT_WORK_DIR=$$(pwd) GGIT_SKIP_RUN_USER_CHECK=true ./.bin/ggit admin create-user --name regular-user --password password123 --email regular@example.com -c conf/app.test.ini
	@cd frontend && npm run build
	@echo "Levantando backend de pruebas en el puerto 3006..."
	@GGIT_WORK_DIR=$$(pwd) GGIT_SKIP_RUN_USER_CHECK=true ./.bin/ggit web -c conf/app.test.ini -p 3006 > .backend.log 2>&1 & echo $$! > .backend.pid
	@timeout 60s bash -c 'until curl -s http://localhost:3006 > /dev/null; do sleep 1; done' \
	    || { echo "Error: El backend no arrancó"; cat .backend.log; make kill-e2e; exit 1; }
	@cd frontend && npx playwright test --project=setup --reporter=list
	@echo "Entorno listo. Para bajarlo: make kill-e2e"

# Tests de Playwright (E2E)
test-e2e: e2e-env-up
	@echo "Ejecutando Playwright..."
	@cd frontend && npx playwright test --reporter=list || (echo "--- PLAYWRIGHT FAILED ---" && cd .. && make kill-e2e && exit 1)
	@cd .. && make kill-e2e
	@rm -f conf/app.test.ini

kill-e2e:
	@echo "Cerrando servidores de prueba..."
	@if [ -f .backend.pid ]; then kill $$(cat .backend.pid) || true; fi
	@if [ -f .frontend.pid ]; then kill $$(cat .frontend.pid) || true; fi
	@rm -f .backend.pid .frontend.pid
	@fuser -k 3006/tcp 5556/tcp || true

# Entorno de demo. Levanta una instancia aislada con datos de ejemplo para tomar
# capturas, reutilizando el patron de e2e-config: la configuracion se genera en
# un unico sitio y apunta a data/demo/, nunca a data/ggit.db. El servidor queda
# en segundo plano (como e2e-env-up) para poder capturar la UI; se baja con
# "make kill-demo".

demo-config:
	@mkdir -p data
	@cp conf/app.ini conf/app.demo.ini
	@# Todas las rutas de datos se redirigen bajo data/demo/ para no tocar la
	@# base real. Los valores de app.ini son "data/ggit.db" y
	@# "data/repositories", por eso el prefijo que se sustituye es "data/".
	@sed -i -E \
	        -e 's#^([[:space:]]*(APP_DATA_PATH|PATH|PROVIDER_CONFIG|OBJECTS_PATH|OBJECTS_TEMP_PATH|TEMP_PATH|AVATAR_UPLOAD_PATH|REPOSITORY_AVATAR_UPLOAD_PATH|ROOT)[[:space:]]*=[[:space:]]*)data/#\1data/demo/#' \
	        -e 's#^([[:space:]]*ROOT[[:space:]]*=[[:space:]]*)data/repositories#\1data/demo/repositories#' \
	        conf/app.demo.ini
	@# Falla ruidosamente en vez de arrancar contra la base del desarrollador.
	@grep -qE '^[[:space:]]*PATH[[:space:]]*=[[:space:]]*data/demo/ggit\.db' conf/app.demo.ini \
	    || { echo "ERROR: la base de datos del demo no quedo aislada"; exit 1; }
	@! grep -qE '^[[:space:]]*(PATH|ROOT|APP_DATA_PATH|PROVIDER_CONFIG|OBJECTS_PATH|OBJECTS_TEMP_PATH|TEMP_PATH|AVATAR_UPLOAD_PATH|REPOSITORY_AVATAR_UPLOAD_PATH)[[:space:]]*=[[:space:]]*data/(ggit\.db|repositories|data|attachments|avatars|repo-avatars|lfs-objects|tmp)([[:space:]]*)$$' conf/app.demo.ini \
	    || { echo "ERROR: queda una ruta de datos real en app.demo.ini"; grep -nE '^[[:space:]]*(PATH|ROOT|APP_DATA_PATH|PROVIDER_CONFIG|OBJECTS_PATH|OBJECTS_TEMP_PATH|TEMP_PATH|AVATAR_UPLOAD_PATH|REPOSITORY_AVATAR_UPLOAD_PATH)[[:space:]]*=' conf/app.demo.ini; exit 1; }
	@# CSRF se desactiva porque es un entorno local y aislado para capturas, no
	@# una instalacion expuesta. La API con Basic Auth lo exime igualmente.
	@# [log] ROOT_PATH se fija bajo data/demo/ porque, al no existir "custom/",
	@# CustomDir() cae en "data/" y por defecto los logs irian a data/log/.
	@printf '[security]\nINSTALL_LOCK = true\nENABLE_CSRF = false\nCOOKIE_SAMESITE = Lax\nCOOKIE_SECURE = false\n[auth]\nREQUIRE_SIGNIN_VIEW = false\n[server]\nHTTP_PORT = 3007\nEXTERNAL_URL = http://localhost:3007/\nLOCAL_ROOT_URL = http://localhost:3007/\n[log]\nROOT_PATH = data/demo/log\n' >> conf/app.demo.ini

# Siembra el dataset de ejemplo en la instancia de demo ya arrancada
# (scripts/seed_demo.sh). Usa Basic Auth sobre la API v2, exenta de CSRF.
# Rellena usuarios, organizacion, varios repositorios, archivos, etiquetas,
# hitos, issues con comentarios, rama + PR + review, CI statuses, tag,
# release, proteccion de rama y un gist.
seed-demo:
	@echo "Sembrando datos de ejemplo..."
	@chmod +x scripts/seed_demo.sh && ./scripts/seed_demo.sh http://localhost:3007

demo: build demo-config
	@echo "Limpiando instancia de demo previa..."
	@if [ -f .demo.pid ]; then kill $$(cat .demo.pid) || true; fi
	@fuser -k 3007/tcp || true
	@rm -rf data/demo
	@mkdir -p data/demo/repositories
	@# El hook pre-receive escribe su log en <ROOT_PATH>/hooks/ antes de que
	@# SetEngine() cree ese subdirectorio, asi que en una instalacion fresca el
	@# primer push falla si no existe. Se crea aqui para poder sembrar commits.
	@mkdir -p data/demo/log/hooks
	@echo "Compilando frontend (necesario para servir la SPA en produccion)..."
	@cd frontend && npm run build
	@echo "Creando usuario administrador de demo..."
	@GGIT_WORK_DIR=$$(pwd) GGIT_SKIP_RUN_USER_CHECK=true ./.bin/ggit admin create-user \
	    --name demo --password demo1234 --email demo@example.com --admin -c conf/app.demo.ini
	@echo "Levantando demo aislado en el puerto 3007..."
	@GGIT_WORK_DIR=$$(pwd) GGIT_SKIP_RUN_USER_CHECK=true ./.bin/ggit web -c conf/app.demo.ini -p 3007 > .demo.log 2>&1 & echo $$! > .demo.pid
	@timeout 60s bash -c 'until curl -s http://localhost:3007/health > /dev/null; do sleep 1; done' \
	    || { echo "ERROR: la demo no arranco"; cat .demo.log; $(MAKE) kill-demo; exit 1; }
	@$(MAKE) seed-demo || { echo "ERROR: no se pudo sembrar el repositorio"; $(MAKE) kill-demo; exit 1; }
	@echo ""
	@echo "Demo listo:     http://localhost:3007"
	@echo "Usuario:        demo"
	@echo "Contrasena:     demo1234"
	@echo "Log:            .demo.log"
	@echo "Para detenerla: make kill-demo"

kill-demo:
	@echo "Cerrando instancia de demo..."
	@if [ -f .demo.pid ]; then kill $$(cat .demo.pid) || true; fi
	@rm -f .demo.pid
	@fuser -k 3007/tcp || true

# Limpieza
clean:
	@echo "Limpiando binarios..."
	rm -rf .bin
