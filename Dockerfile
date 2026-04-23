# Etapa 1: Compilar el Frontend (Svelte)
FROM node:20-alpine AS frontendbuilder
WORKDIR /app/ggit/frontend
COPY frontend/package*.json ./
RUN npm install
COPY frontend/ ./
# Aseguramos que se usa el adaptador estático para generar los archivos de la SPA
RUN sed -i "s|import adapter from '@sveltejs/adapter-auto'|import adapter from '@sveltejs/adapter-static'|g" svelte.config.js
RUN npm run build

# Etapa 2: Compilar el Backend (Go)
FROM golang:1.26-alpine3.23 AS binarybuilder
RUN apk --no-cache add build-base git linux-pam-dev bash
WORKDIR /app/ggit
COPY . .
RUN ./docker/build/install-task.sh
RUN chmod +x scripts/sync_templates.sh && ./scripts/sync_templates.sh
# El tag bindata es vital para que las plantillas internas funcionen si falla el frontend
RUN TAGS="cert pam bindata" task build

# Etapa 3: Imagen Final (Alpine)
FROM alpine:3.23
RUN apk --no-cache add \
  bash \
  ca-certificates \
  curl \
  git \
  linux-pam \
  openssh \
  s6 \
  shadow \
  socat \
  tzdata \
  rsync \
  "zlib>1.3.2"

ENV GGIT_CUSTOM=/data/ggit
WORKDIR /app/ggit

# Copiar configuraciones y scripts de inicio
COPY docker ./docker
# Copiar el binario compilado en la etapa 2
COPY --from=binarybuilder /app/ggit/.bin/ggit .
# Copiar el frontend compilado en la etapa 1
COPY --from=frontendbuilder /app/ggit/frontend/build ./frontend/build

RUN ./docker/build/finalize.sh

# Configurar volúmenes y puertos
VOLUME ["/data", "/backup"]
EXPOSE 22 3005

# Healthcheck apuntando al nuevo puerto
HEALTHCHECK CMD (curl --noproxy localhost -o /dev/null -sS http://localhost:3005/healthcheck) || exit 1

ENTRYPOINT ["/app/ggit/docker/start.sh"]
CMD ["/usr/bin/s6-svscan", "/app/ggit/docker/s6/"]
