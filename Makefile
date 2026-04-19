.PHONY: dev backend frontend build clean help install-deps

# Default target
help:
	@echo "Comandos disponibles:"
	@echo "  make dev       - Ejecuta backend y frontend en paralelo"
	@echo "  make backend   - Compila y ejecuta el servidor Gogs (puerto 3005)"
	@echo "  make frontend  - Ejecuta el servidor de desarrollo Svelte (puerto 5173)"
	@echo "  make build     - Compila el binario de Gogs"
	@echo "  make install   - Instala las dependencias del frontend"

# Limpiar puertos antes de empezar
clean-ports:
	@echo "Limpiando puertos 3005 y 5555..."
	@fuser -k 3005/tcp 5555/tcp || true

# Ejecuta ambos en paralelo usando el operador & de shell
dev: clean-ports
	@echo "Iniciando Gogs Modernized (Desarrollo)..."
	@make -j 2 backend frontend

# Compila y corre el backend
backend: build
	@echo "Arrancando Backend Go en http://localhost:3005 ..."
	./gogs web -p 3005

# Corre el frontend de Svelte
frontend:
	@echo "Arrancando Frontend Svelte en http://localhost:5555 ..."
	@cd frontend && npm run dev -- --port 5555

# Compila el binario
build:
	@echo "Compilando Gogs..."
	@go build -o gogs ./cmd/ggit

# Instalación de dependencias
install:
	@echo "Instalando dependencias de Node.js..."
	@cd frontend && npm install

# Limpieza
clean:
	@echo "Limpiando binarios..."
	rm -f gogs
