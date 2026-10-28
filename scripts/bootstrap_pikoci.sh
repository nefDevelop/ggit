#!/usr/bin/env bash
#
# Prepara PikoCI y lo enlaza con ggit: genera data/pikoci/pikoci.env, levanta
# el compose combinado y escribe API_TOKEN y BOT_TOKEN en data/conf/app.ini.
#
#   podman-compose -f docker-compose.yml -f docker-compose.pikoci.yml up -d
#
# Es idempotente: si app.ini ya tiene ambos tokens no crea ninguno nuevo.
#
# Variables de entorno admitidas:
#   CONTAINER_RUNTIME  docker|podman (se autodetecta)
#   PIKOCI_IMAGE       imagen a usar para generar el hash y arrancar
#   PIKOCI_CONTAINER   nombre del contenedor de PikoCI (pikoci)
#   GGIT_CONTAINER     nombre del contenedor de ggit (ggit)
#   PIKOCI_TEAM        equipo de PikoCI (main)
#   PIKOCI_HTTP_PORT   puerto publicado de la UI/API (18080)
#   PIKOCI_PUBLIC_URL  URL pública usada en los enlaces de estado
#   PIKOCI_INTERNAL_URL URL que ggit usa para hablar con PikoCI
#   GGIT_HTTP_PORT     puerto publicado de ggit (3005)
#   GGIT_REPO          OWNER/NAME al que añadir el bot como colaborador
#   APP_INI            ruta de app.ini en el host (data/conf/app.ini)
#   ENV_FILE           ruta de pikoci.env en el host (data/pikoci/pikoci.env)
#   SKIP_UP            true para no levantar el compose

set -euo pipefail

ROOT_DIR=$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)
cd "$ROOT_DIR"

PIKOCI_IMAGE=${PIKOCI_IMAGE:-ghcr.io/pikoci/pikoci:v0.8.0}
PIKOCI_CONTAINER=${PIKOCI_CONTAINER:-pikoci}
GGIT_CONTAINER=${GGIT_CONTAINER:-ggit}
PIKOCI_TEAM=${PIKOCI_TEAM:-main}
PIKOCI_HTTP_PORT=${PIKOCI_HTTP_PORT:-18080}
PIKOCI_PUBLIC_URL=${PIKOCI_PUBLIC_URL:-http://localhost:${PIKOCI_HTTP_PORT}}
PIKOCI_INTERNAL_URL=${PIKOCI_INTERNAL_URL:-http://pikoci:8080}
GGIT_HTTP_PORT=${GGIT_HTTP_PORT:-3005}
GGIT_REPO=${GGIT_REPO:-}
APP_INI=${APP_INI:-data/conf/app.ini}
ENV_FILE=${ENV_FILE:-data/pikoci/pikoci.env}
SKIP_UP=${SKIP_UP:-false}

info() { printf '%s\n' "$*"; }
warn() { printf 'aviso: %s\n' "$*" >&2; }
die() { printf 'error: %s\n' "$*" >&2; exit 1; }

detect_runtime() {
	if [ -n "${CONTAINER_RUNTIME:-}" ]; then
		RUNTIME=$CONTAINER_RUNTIME
	elif command -v podman >/dev/null 2>&1; then
		RUNTIME=podman
	elif command -v docker >/dev/null 2>&1; then
		RUNTIME=docker
	else
		die "no encuentro podman ni docker"
	fi
	command -v "$RUNTIME" >/dev/null 2>&1 || die "el runtime $RUNTIME no está instalado"
}

detect_compose() {
	if command -v docker >/dev/null 2>&1 && docker compose version >/dev/null 2>&1; then
		COMPOSE=(docker compose)
	elif command -v podman-compose >/dev/null 2>&1; then
		COMPOSE=(podman-compose)
	elif command -v docker-compose >/dev/null 2>&1; then
		COMPOSE=(docker-compose)
	else
		die "no encuentro docker compose ni podman-compose"
	fi
}

read_env() {
	sed -nE "s/^$1=(.*)\$/\1/p" "$ENV_FILE" | head -n1
}

generate_env_file() {
	local admin_password hash
	mkdir -p "$(dirname "$ENV_FILE")"

	if [ -s "$ENV_FILE" ]; then
		JWT_SECRET=$(read_env JWT_SECRET)
		ADMIN_PASSWORD=$(read_env PIKOCI_ADMIN_PASSWORD)
		PIKOCI_USERS=$(read_env USERS)
		if [ -z "$JWT_SECRET" ] || [ -z "$ADMIN_PASSWORD" ] || [ -z "$PIKOCI_USERS" ]; then
			die "$ENV_FILE existe pero le faltan JWT_SECRET, PIKOCI_ADMIN_PASSWORD o USERS"
		fi
		return
	fi

	JWT_SECRET=$(openssl rand -hex 32)
	admin_password=$(openssl rand -base64 32 | tr -dc 'A-Za-z0-9' | head -c 24)
	hash=$("$RUNTIME" run --rm "$PIKOCI_IMAGE" user-password -u admin -p "$admin_password" | sed -n 's/^admin://p')
	[ -n "$hash" ] || die "no se pudo generar el hash de la contraseña de admin"

	ADMIN_PASSWORD=$admin_password
	PIKOCI_USERS="admin:$hash"

	{
		printf 'JWT_SECRET=%s\n' "$JWT_SECRET"
		printf 'DB_SYSTEM=sqlite\n'
		printf 'DB_NAME=/data/pikoci.db\n'
		printf 'RUN_WORKER=true\n'
		printf 'PIKOCI_ADMIN_PASSWORD=%s\n' "$ADMIN_PASSWORD"
		printf 'USERS=%s\n' "$PIKOCI_USERS"
	} >"$ENV_FILE"
	chmod 600 "$ENV_FILE"
	info "generado $ENV_FILE (guárdalo: contiene la contraseña de admin de PikoCI)"
}

wait_http() {
	local url=$1 name=$2
	for _ in $(seq 1 60); do
		if curl -fsS -m 2 "$url" >/dev/null 2>&1; then
			return 0
		fi
		sleep 2
	done
	die "$name no responde en $url"
}

wait_for_file() {
	local path=$1 name=$2
	for _ in $(seq 1 30); do
		[ -f "$path" ] && return 0
		sleep 2
	done
	die "no aparece $path; arranca $name una vez para que genere la configuración"
}

create_api_token() {
	local ids id out token

	# El CLI devuelve 1 tras un login correcto; el create de abajo lo verifica.
	"$RUNTIME" exec "$PIKOCI_CONTAINER" pikoci client -u localhost:8080 login \
		--username admin --password "$ADMIN_PASSWORD" >/dev/null 2>&1 || true

	ids=$("$RUNTIME" exec "$PIKOCI_CONTAINER" sh -c \
		'pikoci client -u localhost:8080 api-tokens list | jq -r ".[] | select(.name == \"ggit\") | .id"')
	for id in $ids; do
		"$RUNTIME" exec "$PIKOCI_CONTAINER" pikoci client -u localhost:8080 \
			api-tokens delete --id "$id" >/dev/null
	done

	out=$("$RUNTIME" exec "$PIKOCI_CONTAINER" pikoci client -u localhost:8080 \
		api-tokens create --name ggit --team-canonical "$PIKOCI_TEAM" --role maintain)
	token=$(printf '%s\n' "$out" | sed -nE 's/^[[:space:]]*"token":[[:space:]]*"([^"]+)".*/\1/p' | head -n1)
	[ -n "$token" ] || die "no se pudo crear el API token de PikoCI: $out"
	printf '%s' "$token"
}

create_bot_token() {
	local args=() out token
	[ -n "$GGIT_REPO" ] && args=(--repo "$GGIT_REPO")

	out=$("$RUNTIME" exec -u git "$GGIT_CONTAINER" /app/ggit/ggit admin pikoci bot \
		${args[@]+"${args[@]}"} 2>&1) || true
	token=$(printf '%s\n' "$out" | sed -nE 's/^Token "[^"]*": ([^[:space:]]+).*/\1/p' | tail -n1)
	if [ -z "$token" ]; then
		printf '%s\n' "$out" >&2
		die "no se pudo crear el token del bot de ggit"
	fi
	if [ -n "$GGIT_REPO" ] && ! printf '%s' "$out" | grep -q 'es colaborador'; then
		warn "el token se creó, pero no se pudo añadir el bot a $GGIT_REPO"
	fi
	printf '%s' "$token"
}

write_app_ini() {
	local api_token=$1 bot_token=$2 tmp
	mkdir -p "$(dirname "$APP_INI")"
	# Se escribe en un temporal y se renombra: el directorio puede ser escribible
	# aunque app.ini pertenezca a otro UID (contenedor).
	tmp=$(mktemp "$(dirname "$APP_INI")/.app.ini.XXXXXX")
	{
		awk 'BEGIN { skip = 0 } /^\[pikoci\]/ { skip = 1; next } /^\[/ { skip = 0 } skip == 0 { print }' \
			"$APP_INI"
		printf '\n[pikoci]\n'
		printf 'ENABLED = true\n'
		printf 'MODE = sidecar\n'
		printf 'URL = %s\n' "$PIKOCI_INTERNAL_URL"
		printf 'PUBLIC_URL = %s\n' "$PIKOCI_PUBLIC_URL"
		printf 'TEAM = %s\n' "$PIKOCI_TEAM"
		printf 'API_TOKEN = %s\n' "$api_token"
		printf 'BOT_TOKEN = %s\n' "$bot_token"
	} >"$tmp"
	chmod 644 "$tmp"
	mv "$tmp" "$APP_INI"
}

detect_runtime
detect_compose
generate_env_file

if [ "$SKIP_UP" != "true" ]; then
	info "levantando ggit y PikoCI..."
	"${COMPOSE[@]}" -f docker-compose.yml -f docker-compose.pikoci.yml up -d
fi

wait_http "http://127.0.0.1:${PIKOCI_HTTP_PORT}/health" "PikoCI"
wait_http "http://127.0.0.1:${GGIT_HTTP_PORT}/health" "ggit"
wait_for_file "$APP_INI" "ggit"

existing_api=$(sed -nE 's/^[[:space:]]*API_TOKEN[[:space:]]*=[[:space:]]*([^[:space:]]+).*/\1/p' "$APP_INI" | head -n1)
existing_bot=$(sed -nE 's/^[[:space:]]*BOT_TOKEN[[:space:]]*=[[:space:]]*([^[:space:]]+).*/\1/p' "$APP_INI" | head -n1)

api_token=$existing_api
bot_token=$existing_bot

if [ -n "$api_token" ] && [ -n "$bot_token" ]; then
	info "app.ini ya tiene los dos tokens; no se genera nada nuevo"
	exit 0
fi

if [ -z "$api_token" ]; then
	info "creando API token en PikoCI..."
	api_token=$(create_api_token)
fi

if [ -z "$bot_token" ]; then
	info "creando token del bot de ggit..."
	bot_token=$(create_bot_token)
fi

write_app_ini "$api_token" "$bot_token"
info "configuración [pikoci] escrita en $APP_INI"

info "reiniciando ggit para aplicar la configuración..."
"$RUNTIME" restart "$GGIT_CONTAINER" >/dev/null
wait_http "http://127.0.0.1:${GGIT_HTTP_PORT}/health" "ggit"

info "listo"
info "  PikoCI:     $PIKOCI_PUBLIC_URL"
info "  ggit -> CI: $PIKOCI_INTERNAL_URL (equipo $PIKOCI_TEAM)"
info "  admin PikoCI: admin / $(read_env PIKOCI_ADMIN_PASSWORD)  ($ENV_FILE)"
