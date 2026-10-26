# Roadmap: ggit pikoci — CI/CD integrado

**ggit pikoci** es la integración de [PikoCI](https://github.com/PikoCI/pikoci) como motor
de CI/CD dentro de ggit. PikoCI es un binario único (o contenedor) con pipelines en HCL,
base de datos propia y UI web.

> Filosofía revisada: el CI **no** se incrusta en el proceso de ggit. ggit configura,
> dispara y muestra resultados; PikoCI ejecuta. El aislamiento manda sobre la comodidad.

Verificado contra **PikoCI v0.8.0** (commit `6ec42e1`, 2026-10-25). El clon de consulta
vive en `tmp/pikoci/` (ignorado por git); sus docs están en `tmp/pikoci/docs/`.

---

## Decisiones de arquitectura (leer antes de tocar código)

### 1. Aislamiento: sidecar por defecto, subproceso como modo reducido

Los jobs de PikoCI son **código arbitrario ejecutado como subproceso del host** con los
runners `exec`/`shell`; `docker` es opcional y no hay sandbox. Si ggit lo lanza como hijo
con su mismo usuario, un pipeline puede leer `conf/app.ini` (con `SECRET_KEY`), la SQLite y
todos los repos bare.

- **Producción (recomendado): contenedor/sidecar.** Volumen propio, red interna, sin acceso
  al sistema de ficheros de ggit. ggit solo consulta estado y envía triggers.
- **Instalaciones single-binary: subproceso** con UID/GID distinto
  (`SysProcAttr.Credential`), `XDG_DATA_HOME` propio y límites de recursos. Requiere
  **apagado ordenado en ggit**, que hoy no existe (`cmd/ggit/web.go` usa `ListenAndServe`
  sin `Shutdown`).

Ningún modo se implementa sin resolver antes el aislamiento.

### 2. Trigger: llamada directa desde el push, no webhooks de ggit

La entrega de webhooks de ggit bloquea loopback por SSRF (`netx.IsBlockedLocalHostname`,
`internal/database/webhook.go`). Además, PikoCI ignora el payload y no verifica HMAC: el
webhook solo re-lanza el `check` del resource.

Por tanto: el push llama **en proceso** a un servicio de ggit que hace
`POST /webhooks/<token>` contra PikoCI. Sin webhook de ggit, sin allowlist de loopback.

### 3. Estado: commit statuses nativos de ggit

PikoCI notifica con `notification_type` + `notify "exec"` y ggit ya tiene
`POST /api/v2/{owner}/{repo}/statuses/{sha}` (`internal/route/api/v2/repo.go`), con
visualización en commits y PRs. No hace falta backend nuevo para esto.

### 4. Pipelines: per-repo desde el propio repo

En vez de un HCL global editable en la UI, el pipeline vive en `.ggit/ci.hcl` del repo y
ggit lo publica en PikoCI al activar CI. Encaja con la filosofía docs-as-code del wiki.
Ojo: la API espera el HCL **en base64** dentro de `{"config": ..., "vars": {...}}`.

### 5. Despliegue: mismo PC o PC de CI

Son dos decisiones separadas: dónde vive el **servidor** de PikoCI y dónde corren los
**workers** (donde se ejecuta el código).

| Escenario | Servidor | Workers | Cuándo |
| :--- | :--- | :--- | :--- |
| Mismo PC, contenedor | sidecar junto a ggit | embebido (`--run-worker=true`) | Simple; builds ligeros. Docker publica el puerto solo en loopback. |
| Mismo PC, binario | subproceso con UID propio | embebido | Solo si no hay contenedores; exige apagado ordenado. |
| PC de CI | contenedor o binario en el otro PC | embebido o workers extra | Builds pesados o separación física de la forja. |

Recomendado para empezar: **mismo PC, PikoCI en contenedor**. El pipeline no puede leer
`conf/app.ini`, la SQLite ni los repos bare porque no comparte volúmenes, y el mapeo
`127.0.0.1:18080:18080` resuelve que PikoCI no tenga flag de bind.

En el PC de CI, ggit y PikoCI se ven por red interna: PikoCI necesita alcanzar ggit por
HTTPS para clonar, y ggit necesita alcanzar la URL de PikoCI (`[pikoci] URL`). Firewall
para que el puerto solo acepte tráfico de ggit.

Workers remotos (`pikoci worker --pikoci-url ... --worker-token ...`): permiten tener la
máquina de builds separada del servidor, con `--run-worker=false` en el servidor y routing
por `--tags`. Para una máquina air-gapped, las toolchains e imágenes Docker deben estar
preinstaladas o en un mirror local; el worker solo necesita red hacia ggit y PikoCI.

En ambos modos la integración es idéntica: ggit solo conoce una `URL` y habla HTTP.


---

## Correcciones respecto a la versión anterior del roadmap

| Antes | Realidad en v0.8.0 |
| :--- | :--- |
| `--http-bind-addr 127.0.0.1:18080` | No existe. Solo `--port`, y escucha en `:port` (todas las interfaces). Exige firewall o red de contenedor. |
| `--db-path ~/.ggit/pikoci/data/pikoci.db` | No existe. La ruta SQLite es fija en el data dir XDG; se controla con `XDG_DATA_HOME` del proceso. `--db-name` se ignora en sqlite. |
| `GET /healthz` | Es `GET /health` (JSON con estado y ping a la base de datos). También `/version` y `/metrics` (Prometheus). |
| `PUT /api/v1/teams/main/pipelines/...` | Sin prefijo `/api/v1`: `PUT /teams/main/pipelines/<nombre>`, body `{"config": "<HCL base64>", "vars": {...}}`. |
| Webhooks con HMAC-SHA256 verificable | No hay HMAC ni validación de payload. `POST /webhooks/<token>` ignora el body; el token es el único secreto. |
| `$notify_sha`, `$notify_state`, `$notify_target_url` automáticos | No se rellenan solos: existen como `$notify_<attr>` solo si el paso `notify` declara ese atributo. Para el SHA usa `$version_ref` (en `on_trigger` no hay checkout). |
| `pikoci validate` | El comando es `pikoci pipeline validate <file.hcl>`. |
| ggit crea un webhook hacia `127.0.0.1` | Bloqueado por el guard SSRF de ggit. Se sustituye por llamada directa (decisión 2). |
| Fase de dashboard de builds en ggit | Duplica la UI de PikoCI. Se descarta: enlace o proxy. |

También conviene saber: los tokens de API de PikoCI son JWT o `pko_`, y **crear un token
requiere un JWT** (`POST /login`), así que ggit necesita guardar una credencial de servicio
creada una vez. Y el usuario `admin`/`admin123` se siembra por defecto: hay que
sobreescribirlo con `--users`.

---

## Fase 0 — Spike de viabilidad (antes de escribir código)

**Objetivo:** validar las tres incógnitas que pueden tumbar el diseño.

- [ ] Lanzar PikoCI a mano con `--db-system sqlite`, `--jwt-secret`, `XDG_DATA_HOME` propio
      y `--users ci:<hash-bcrypt>`.
- [ ] **Repo privado de ggit**: `pull` inyecta `https://oauth2:<token>@…` y ggit acepta el
      token como contraseña (`internal/route/repo/http.go`), pero el fallback de `check`
      (`git ls-remote` en `pikoci/builtin/resource_types/git.hcl`) **no inyecta
      credenciales**. Verificar si el `check` funciona; si no, evaluar URL con credenciales
      vía secret store o añadir un provider `ggit` al resource type.
- [ ] Reportar estado a `POST /api/v2/{owner}/{repo}/statuses/{sha}` desde un `notify "exec"`
      con `state`, `context`, `description`, `target_url` y `sha = "$version_ref"`.
- [ ] Pipeline mínimo HCL con `get → task → notify` y `pikoci pipeline validate`.

**Criterio de salida:** push a un repo privado → build ejecutado → estado visible en el
commit de ggit.

---

## Fase 1 — Configuración y ciclo de vida

**Objetivo:** ggit conoce, arranca (si aplica), vigila y detiene PikoCI.

### 1.1 Sección `[pikoci]` en `app.ini`

- [ ] `PikoCIOpts` en `internal/conf/static.go` (patrón de `[cron]`) y mapeo en
      `internal/conf/conf.go`.
- [ ] Defaults en `conf/app.ini`.
- [ ] Secretos generados y persistidos como `SECRET_KEY` (`ensureSecretKey`), nunca en
      `conf/app.ini`.

```ini
[pikoci]
ENABLED     = false
MODE        = sidecar   ; sidecar | subprocess
URL         = http://127.0.0.1:18080
BINARY_PATH =           ; solo MODE=subprocess
DATA_DIR    =           ; se exporta como XDG_DATA_HOME al hijo
CONCURRENCY = 2
PUBLIC_URL  =           ; enlace a la UI de PikoCI para el admin
```

### 1.2 Servicio `internal/services/pikoci`

- [ ] `var CI *Service` + `Init(db, conf)` llamado desde `internal/services/services.go`,
      siguiendo el patrón de `internal/services/repo`.
- [ ] `Status(ctx)` (GET `/health`), `Restart(ctx)`, `LogTail(n)`.
- [ ] Modo subprocess: supervisor con backoff + `process.Add("pikoci", cmd)` para que
      aparezca en el monitor admin (`GetAdminMonitor`, `KillAdminProcess`).
- [ ] Modo sidecar: sin gestión de proceso, solo health y credenciales.

### 1.3 Apagado ordenado (prerrequisito del modo subprocess)

- [ ] `signal.NotifyContext` en `cmd/ggit/web.go` y `server.Shutdown(ctx)`.
- [ ] `CI.Stop(ctx)`: `SIGQUIT` (drain, `--drain-timeout`) y `SIGKILL` si no responde.

### 1.4 Diagnóstico

- [ ] `ggit admin pikoci status` (modo, URL, salud, versión) y `restart` (subprocess).
- [ ] Logs de PikoCI con rotación (10 MB/archivo).

**Criterio de salida:** el admin ve el estado real de PikoCI desde ggit y el arranque/apagado
no deja procesos huérfanos.

---

## Fase 2 — Trigger en el push

**Objetivo:** un push dispara el `check` del resource sin webhooks de ggit.

- [ ] `CI.NotifyPush(ctx, repo, ref)` → `POST /webhooks/<token>` de PikoCI.
- [ ] Obtener y guardar el token con
      `POST /teams/{team}/pipelines/{pipeline}/resources/{resource}/webhook_token`
      (rol `maintain`).
- [ ] Llamada desde `cmd/ggit/hook.go` (push por SSH) y desde el handler de push HTTP.
- [ ] Deduplicar por `(repo, sha)` para no lanzar builds repetidos.
- [ ] Si PikoCI no responde, el push **no** falla: se registra y se sigue.

**Criterio de salida:** push → build en PikoCI en segundos, sin webhook de ggit configurado.

---

## Fase 3 — Estado en commits y PRs

**Objetivo:** ✅/❌ en el commit y checks en el PR.

- [ ] Usuario bot `pikoci` + token de acceso, colaborador de los repos con CI.
      (Los tokens de ggit son de usuario completo, sin scopes: el bot solo debe tener
      acceso a esos repos.)
- [ ] `notification_type "ggit-status"` con `notify "exec"` y `curl`, declarando:
      `state`, `sha = "$version_ref"`, `context = "ci/pikoci"`, `description`,
      `target_url` (enlace al build en PikoCI).
- [ ] Estados usados: `pending` al empezar, `success`/`failure`/`error` al terminar.
- [ ] Verificar que el PR muestra `combined_status` (`internal/database/commit_status.go`).

**Criterio de salida:** PR bloqueado por estado `failure` y mergeable tras `success`, con
la protección de rama existente.

---

## Fase 4 — Pipeline per-repo desde `.ggit/ci.hcl`

**Objetivo:** cada repo define su CI en su propio código.

- [ ] `EnablePikoCI bool` en `Repository` (`internal/database/repo.go`), espejo en
      `internal/domain/repository.go` y en los conversores.
- [ ] Migración explícita: `AutoMigrate` **no** añade columnas a tablas existentes; hace
      falta un `NewMigration` al final de `internal/database/migrations/migrations.go`.
- [ ] Al activar CI en un repo: leer `.ggit/ci.hcl` del commit, sustituir variables
      (`$GGIT_URL`, `$REPO_FULL_NAME`) y `PUT /teams/main/pipelines/<owner-repo>` con el HCL
      en base64.
- [ ] Re-publicar al cambiar el fichero; borrar el pipeline al desactivar CI.
- [ ] Validación previa con `pikoci pipeline validate` antes de publicar.
- [ ] Endpoint `GET/PUT /:user/:repo/ci/pipeline` + toggle en ajustes del repo.

**Criterio de salida:** dos repos con stacks distintos (Go y Node) ejecutan pipelines
distintos definidos en su propio `.ggit/ci.hcl`.

---

## Fase 5 — UI admin mínima

- [ ] Sección "CI/CD" en `frontend/src/routes/admin/` registrada en `admin/+layout.svelte`.
- [ ] Estado (RUNNING / STOPPED / UNREACHABLE), modo, versión, salud de la base de datos.
- [ ] Botón "Restart" (solo modo subprocess) con confirmación.
- [ ] Enlace a la UI de PikoCI (`PUBLIC_URL`).
- [ ] Aviso visible cuando `MODE=subprocess` (riesgo de aislamiento).

---

## Backlog (no ahora)

- Releases automáticas al pushear un tag: **hecho** como plantilla opcional
  (`GET /api/v2/admin/pikoci/release-template`), validada con PikoCI v0.8.0.
- Aprobaciones manuales (`approve`, ya soportado por PikoCI).
- Dashboard de builds en ggit: **descartado**; enlazar a PikoCI.
- Scopes por repo en `AccessToken` (hoy son de usuario completo).
- Proxy de la UI de PikoCI bajo autenticación de ggit (hoy implica doble login).

---

## Seguridad

- **Aislamiento primero**: contenedor o UID separado; el pipeline no debe poder leer
  `conf/app.ini`, la base de datos ni los repos bare de ggit.
- PikoCI escucha en todas las interfaces: firewall o red de contenedor, nunca expuesto.
- Cambiar `admin`/`admin123` con `--users` (hash bcrypt) y rotar `--jwt-secret` y
  `--secret-key`.
- Los webhooks de PikoCI no firman: token largo, secreto, rotable.
- Token de ggit para CI: usuario bot, mínimo acceso, rotación periódica.
- Versionar PikoCI: fijar release y verificar `SHA256SUMS`; no `latest` a ciegas.

---

## Puntos de integración en ggit

| Necesidad | Dónde |
| :--- | :--- |
| Config y secretos | `internal/conf/static.go`, `internal/conf/conf.go`, `ensureSecretKey` |
| Arranque/parada | `internal/route/init.go` (`GlobalInit`), `cmd/ggit/web.go` |
| Registro de proceso | `internal/process/manager.go` (`Add`, `Kill`) |
| Servicio nuevo | `internal/services/pikoci` + `internal/services/services.go` |
| Trigger en push | `cmd/ggit/hook.go`, handler de push HTTP |
| Commit status | `POST /api/v2/:user/:repo/statuses/:sha` (`internal/route/api/v2/repo.go`) |
| Webhooks (solo si hicieran falta) | `internal/database/webhook.go` (guard SSRF, HMAC saliente) |
| Modelo y migraciones | `internal/database/repo.go`, `internal/database/migrations/` |
| Admin API | `internal/route/api/v2/admin.go`, `internal/route/api/v2/v2.go` |
| Frontend admin | `frontend/src/routes/admin/*`, `frontend/src/lib/api.ts` |

---

## Referencias de PikoCI v0.8.0

Con el clon en `tmp/pikoci/`:

- Flags del servidor: `tmp/pikoci/cmd/server.go`.
- Rutas HTTP: `tmp/pikoci/pikoci/transport/http/handler.go`.
- Resource `git` (auth, `tag`, `pr`, `check`):
  `tmp/pikoci/pikoci/builtin/resource_types/git.hcl`.
- Notificaciones y variables: `tmp/pikoci/docs/Notifications.md`.
- Pipelines y primitivas: `tmp/pikoci/docs/Pipeline.md`.
- Base de datos y migraciones: `tmp/pikoci/docs/Database.md`.
- Despliegue y releases: `tmp/pikoci/docs/Deployment.md`, `tmp/pikoci/deploy/pipeline.hcl`.

---

## Resumen de fases

| Fase | Qué | Depende de | Riesgo |
| :--- | :--- | :--- | :--- |
| 0 | Spike: repo privado, status, notify | — | Alto (decide el diseño) |
| 1 | Config, servicio, ciclo de vida | 0 | Medio |
| 2 | Trigger directo en push | 1 | Bajo |
| 3 | Commit status y checks en PR | 0, 2 | Bajo |
| 4 | Pipeline per-repo desde `.ggit/ci.hcl` | 1, 2 | Medio |
| 5 | UI admin mínima | 1 | Bajo |

---

## Estado de implementación

Rama `feature/pikoci-ci`.

| Fase | Estado |
| :--- | :--- |
| 0 | **Pendiente**: requiere ggit y PikoCI vivos |
| 1.1 Configuración | Hecho |
| 1.2 Servicio sidecar | Hecho (health y estado) |
| 1.2 Subprocess | Pendiente (opcional, decidido sidecar) |
| 1.3 Apagado ordenado | Hecho |
| 1.4 CLI `admin pikoci status` | Hecho |
| 2 Trigger en push | Hecho (SSH y HTTP) |
| 3 Estado en commits | Hecho: bot, token y plantilla validada con `pikoci pipeline validate` |
| 4 Aprovisionamiento | Hecho: sync, webhook token, auto-republicación al cambiar el fichero y toggle por repo |
| 5 UI | Tarjeta de estado en admin y pestaña CI/CD en ajustes del repo |
| 7 Releases | Plantilla validada y servida por API; se activa como pipeline aparte |

Cambios de diseño respecto al plan original:

- No se añade `EnablePikoCI` a `Repository`: la fila en `pikoci_repo` **es** el estado
  de CI activado, así que no hay migración ni segunda fuente de verdad.
- Se descarta el dedupe por `(repo, sha)`: el webhook solo relanza el `check` y PikoCI no
  construye si la versión no cambió.
- `API_TOKEN` y `BOT_TOKEN` viven en `app.ini`, como el resto de credenciales de ggit.
- La plantilla solo se republica cuando el push toca `.ggit/ci.hcl` en la rama por
  defecto, para no resetear el estado del resource en cada push.

Pendiente real:

- Fase 0: verificar el `check` de repos privados y el reporte de estados con instancias
  reales.
- Fase 5 restante: editar configuración, restart y logs desde el panel.
- Supervisor del modo subprocess, si alguna vez se quiere.

