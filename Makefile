.PHONY: dev backend frontend build clean help install-deps

# Default target
help:
	@echo "Comandos disponibles:"
	@echo "  make dev       - Ejecuta backend y frontend en paralelo"
	@echo "  make backend   - Compila y ejecuta el servidor Gogs (puerto 3000)"
	@echo "  make frontend  - Ejecuta el servidor de desarrollo Svelte (puerto 5173)"
	@echo "  make build     - Compila el binario de Gogs"
	@echo "  make install   - Instala las dependencias del frontend"

# Limpiar puertos antes de empezar
clean-ports:
	@echo "Limpiando puertos 3000 y 5555..."
	@fuser -k 3000/tcp 5555/tcp || true

# Ejecuta ambos en paralelo usando el operador & de shell
dev: clean-ports
	@echo "Iniciando Gogs Modernized (Desarrollo)..."
	@make -j 2 backend frontend

# Compila y corre el backend
backend: build
	@echo "Arrancando Backend Go en http://127.0.0.1:3000 ..."
	./gogs web -p 3000

# Corre el frontend de Svelte
frontend:
	@echo "Arrancando Frontend Svelte en http://127.0.0.1:5555 ..."
	@cd frontend && npm run dev -- --port 5555

# Compila el binario
build:
	@echo "Compilando Gogs..."
	@go build -o gogs ./cmd/gogs

# Instalación de dependencias
install:
	@echo "Instalando dependencias de Node.js..."
	@cd frontend && npm install

# Limpieza
clean:
	@echo "Limpiando binarios..."
	rm -f gogs
