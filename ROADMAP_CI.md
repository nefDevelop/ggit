# Roadmap: ggit pikoci — CI/CD integrado

**ggit pikoci** es la integración de [PikoCI](https://github.com/PikoCI/pikoci) como motor de CI/CD nativo dentro de ggit. Un solo binario externo (~20 MB), cualquier base de datos, pipelines en HCL.

> Filosofía: ggit gestiona el ciclo de vida de PikoCI como subproceso. El administrador controla todo desde la UI de ggit. Los desarrollores ven los resultados (✅/❌) en commits y PRs.

---

## Fase 1: Fundación — Config + Binary Manager

**Objetivo:** ggit puede descargar, lanzar, monitorizar y detener PikoCI como subproceso gestionado.

### 1.1 Config section `[pikoci]` en `app.ini`

- [ ] Crear `PikoCIOpts` struct en `internal/conf/static.go`
- [ ] Mapear sección en `internal/conf/conf.go`
- [ ] Agregar defaults en `conf/app.ini`

```ini
[pikoci]
ENABLED = false
BINARY_PATH = ~/.ggit/pikoci/pikoci
DATA_DIR = ~/.ggit/pikoci/data
AUTO_DOWNLOAD = true
DOWNLOAD_URL = https://github.com/PikoCI/pikoci/releases/latest/download/pikoci-linux-amd64
PORT = 18080
JWT_SECRET = auto-generado-en-init
```

### 1.2 Binary Manager (`internal/pikoci/manager.go`)

- [ ] `DownloadBinary(url, dest)` — descarga + `chmod +x`, con barra de progreso
- [ ] `StartServer(ctx, opts)` — lanza `exec.Cmd` con flags:
  ```
  pikoci server \
    --db-system sqlite \
    --db-path ~/.ggit/pikoci/data/pikoci.db \
    --http-bind-addr 127.0.0.1:18080 \
    --jwt-secret <generado> \
    --run-worker
  ```
- [ ] `StopServer()` — envía SIGTERM, espera graceful shutdown (5s), SIGKILL si no responde
- [ ] `HealthCheck()` — GET `http://127.0.0.1:18080/healthz`, reintentos
- [ ] Registrar PID en `process.Processes` para visibilidad en admin

### 1.3 GlobalInit hook (`internal/route/init.go`)

- [ ] En `GlobalInit()`, si `pikoci.enabled`:
  - Generar JWT secret si no existe (guardar en `data/` persistente)
  - Llamar `DownloadBinary` si aplica
  - Llamar `StartServer`
  - Goroutine de health check cada 30s (restart automático si muere)
- [ ] Registrar shutdown: `signal.NotifyContext` + `StopServer` en defer

### 1.4 Admin monitor visibility

- [ ] PikoCI aparece como proceso en `GET /admin/monitor` con PID, uptime, estado
- [ ] Botón "Restart PikoCI" en admin panel

### 1.5 CLI command

- [ ] `ggit admin pikoci status` — muestra si está corriendo, PID, puerto, uptime
- [ ] `ggit admin pikoci restart` — restart del subproceso
- [ ] `ggit admin pikoci logs` — tail del stdout/stderr de PikoCI

---

## Fase 2: Admin UI — Enable/Disable CI

**Objetivo:** El administrador enciende CI desde el panel de administración de ggit.

### 2.1 Config API endpoints

- [ ] `GET /admin/pikoci/config` — devuelve la config actual (sin secrets)
- [ ] `PUT /admin/pikoci/config` — actualiza config (enabled, port, etc.)
- [ ] `POST /admin/pikoci/restart` — restart del subproceso

### 2.2 Admin panel UI (SvelteKit)

- [ ] Sección "CI/CD" en el panel de administración
- [ ] Toggle "Enable PikoCI CI/CD"
- [ ] Indicador de estado: RUNNING / STOPPED / ERROR / DOWNLOADING
- [ ] Botón "Restart" con confirmación
- [ ] Config fields: Port, Auto-download, Binary path
- [ ] Log viewer (streaming de logs de PikoCI)

### 2.3 Default pipeline template

- [ ] Template HCL embebido en ggit para el pipeline por defecto
- [ ] Incluye resource type `ggit-status` para reportar commit status
- [ ] Incluye job `test` con un paso `make test` por defecto
- [ ] Se inyectan variables: `$GGIT_URL`, `$GGIT_TOKEN`, `$REPO_FULL_NAME`

### 2.4 Startup pipeline

- [ ] Al iniciar PikoCI, ggit crea un pipeline base que escucha webhooks
- [ ] El pipeline se auto-registra via API de PikoCI (`PUT /api/v1/teams/main/pipelines/ggit-ci`)

---

## Fase 3: Webhook Integration — Push Trigger

**Objetivo:** Cuando alguien hace push a un repo, ggit avisa a PikoCI y se dispara el pipeline.

### 3.1 Resource git + webhook token en PikoCI

- [ ] Al crear el pipeline en PikoCI, se define un resource `git` apuntando al repo local
- [ ] Obtener webhook token via API de PikoCI: `POST /teams/main/pipelines/ggit-ci/resources/repo/webhook_token`
- [ ] Guardar token en DB de ggit asociado al repo

### 3.2 Auto-creación de webhook en ggit

- [ ] Al activar CI, ggit crea automáticamente un webhook en el repo:
  - URL: `http://127.0.0.1:18080/webhooks/<token>`
  - Eventos: Push, Pull Request
  - Content type: JSON
  - Secret: mismo JWT secret (para verificar HMAC si PikoCI lo soporta)
- [ ] Al desactivar CI, se elimina el webhook

### 3.3 Post-receive hook alternativo (fallback)

- [ ] El hook `post-receive` existente en `cmd/ggit/hook.go` ya llama un task trigger URL
- [ ] Como respaldo, también puede notificar a PikoCI directamente:
  ```
  curl -sf -X POST http://127.0.0.1:18080/webhooks/<token>
  ```

### 3.4 Prueba de integración

- [ ] Al hacer push a un repo con CI activo:
  1. ggit recibe el push
  2. Webhook delivery envía POST a PikoCI
  3. PikoCI checkea el resource `git`, detecta nuevo commit
  4. Se dispara el job `test`
  5. Pipeline se ejecuta

---

## Fase 4: Commit Status — Feedback en ggit

**Objetivo:** El estado del build (pending/success/failure) se refleja en commits y PRs de ggit.

### 4.1 Custom notification type `ggit-status`

- [ ] Definir `notification_type` dentro del pipeline de PikoCI:

```hcl
notification_type "ggit-status" {
  params = ["ggit_url", "ggit_token", "repo"]
  notify "exec" {
    path = "/bin/sh"
    args = ["-ec", <<-EOT
      curl -sf -X POST "$param_ggit_url/api/v2/$param_repo/statuses/$notify_sha" \
        -H "Authorization: token $param_ggit_token" \
        -H "Content-Type: application/json" \
        -d '{
          "state": "$notify_state",
          "context": "ci/pikoci",
          "description": "$notify_description",
          "target_url": "$notify_target_url"
        }'
    EOT]
  }
}
```

### 4.2 Pipeline template con status reporting

- [ ] El pipeline template por defecto incluye `notify` steps:
  - Al inicio: `status = "pending"`
  - Al éxito: `state = "success"`
  - Al fallo: `state = "failure"`

### 4.3 Visualización en ggit

- [ ] Los commit statuses ya se muestran en:
  - Vista de commits (✅/❌ al lado del SHA) ✅ existente
  - Vista de PR (checks antes de merge) ✅ existente

---

## Fase 5: Pipeline Customization

**Objetivo:** El administrador puede personalizar el pipeline desde la UI de ggit.

### 5.1 Pipeline editor UI

- [ ] Editor de HCL en el admin panel (textarea con resaltado de sintaxis)
- [ ] Validación de sintaxis (llamar `pikoci validate` o parseo básico)
- [ ] Botón "Test" que corre `pikoci run -p pipeline.hcl -j test` localmente

### 5.2 API endpoints

- [ ] `GET /admin/pikoci/pipeline` — devuelve el pipeline HCL actual
- [ ] `PUT /admin/pikoci/pipeline` — actualiza el pipeline
- [ ] `POST /admin/pikoci/pipeline/validate` — valida sintaxis
- [ ] `POST /admin/pikoci/pipeline/test` — ejecuta dry-run

### 5.3 Variables y secrets

- [ ] `GET /admin/pikoci/vars` — lista variables del pipeline
- [ ] `PUT /admin/pikoci/vars` — actualiza variables (git token, etc.)
- [ ] Los secrets se guardan en la DB de ggit (encriptados)

---

## Fase 6: Per-repo CI (opcional — evaluar antes de implementar)

**Objetivo:** Cada repo puede tener CI activado o desactivado, con su propio pipeline.

### Consideración — ¿Por qué per-repo?

Distintos proyectos necesitan pipelines distintos:
- Un repo Go necesita `make test`, otro Node necesita `npm test`
- Diferentes branches, diferentes reglas
- GitLab CI, Gitea Actions y GitHub Actions son per-repo

### Diseño propuesto

- [ ] Campo `ci_enabled` en modelo `Repository`
- [ ] Toggle "Enable CI" en settings del repo (UI)
- [ ] UI para editar pipeline específico del repo (hereda del template global)
- [ ] El webhook se crea/elimina al activar/desactivar CI en el repo
- [ ] Cada repo tiene su propio resource `git` en PikoCI con su webhook token

### Migración desde Fase 2

- [ ] El pipeline global pasa a ser el "template por defecto" para nuevos repos
- [ ] Al activar CI en un repo, se copia el template y se inyectan las variables del repo
- [ ] Admin puede sobreescribir el pipeline de cada repo desde la UI

---

## Fase 7: Releases automáticas

**Objetivo:** PikoCI puede crear releases en ggit cuando se pushea un tag.

### 7.1 Custom resource type `ggit-release`

```hcl
resource_type "ggit-release" {
  params = ["ggit_url", "ggit_token", "repo"]
  push "exec" {
    path = "/bin/sh"
    args = ["-ec", <<-EOT
      curl -sf -X POST "$param_ggit_url/api/v2/$param_repo/releases" \
        -H "Authorization: token $param_ggit_token" \
        -H "Content-Type: application/json" \
        -d '{
          "tag_name": "$put_tag",
          "target":   "$put_branch",
          "name":     "$put_tag",
          "body":     "$put_body"
        }'
    EOT]
  }
}
```

### 7.2 Pipeline template para tags

- [ ] Resource `git` con `tag = true` para detectar tags nuevos
- [ ] Job `release` que corre build + `put "ggit-release"` en `on_success`
- [ ] Opción de subir artifacts como attachments del release

### 7.3 Changelog automático

- [ ] ggit genera changelog entre tags (usando commits)
- [ ] Se inyecta como `$put_body` en el pipeline

---

## Fase 8: Dashboard de Builds

**Objetivo:** Los usuarios ven el historial de builds en ggit.

### 8.1 Vista de builds en el repo

- [ ] Nueva pestaña "CI" en el repo (visible si CI está activo)
- [ ] Lista de builds: número, commit, branch, status, duración, fecha
- [ ] Link a PikoCI UI para ver logs completos

### 8.2 API de builds

- [ ] `GET /:user/:repo/ci/builds` — lista builds del repo
- [ ] `GET /:user/:repo/ci/builds/:id` — detalle de build
- [ ] `POST /:user/:repo/ci/trigger` — trigger manual de build

### 8.3 Sincronización con PikoCI

- [ ] Webhook de PikoCI hacia ggit cuando un build cambia de estado
- [ ] O polling ligero desde ggit a PikoCI cada N segundos

---

## Resumen de fases

| Fase | Qué | Dependencias | Esfuerzo estimado |
|:-----|:----|:-------------|:-----------------|
| 1 | Binary manager + config + lifecycle | — | Media |
| 2 | Admin UI enable/disable + pipeline template | Fase 1 | Media |
| 3 | Webhook integration + push trigger | Fase 1, 2 | Media |
| 4 | Commit status feedback | Fase 3 | Baja |
| 5 | Pipeline customization (UI editor) | Fase 2, 3, 4 | Alta |
| 6 | Per-repo CI (evaluar) | Fases 1-5 | Media |
| 7 | Releases automáticas | Fase 3, 4, 6 | Media |
| 8 | Dashboard de builds | Fase 3, 4 | Alta |

---

## Notas de arquitectura

### Comunicación ggit ↔ PikoCI

```
ggit (background)              PikoCI (subprocess)
      │                              │
      │  POST /webhooks/<token>       │  (webhook push)
      │─────────────────────────────>│
      │                              │  check git resource
      │                              │  run pipeline job
      │                              │
      │  POST /api/v2/statuses/:sha   │  (commit status)
      │<─────────────────────────────│
      │                              │
      │  POST /api/v2/releases        │  (create release)
      │<─────────────────────────────│
```

### Seguridad

- JWT secret auto-generado, almacenado en disco (`data/pikoci/jwt.secret`)
- PikoCI solo escucha en `127.0.0.1` (no expuesto a la red)
- Token de ggit para API se pasa como variable al pipeline (nunca hardcodeado)
- Webhooks con HMAC-SHA256 verificable

### Gestión de recursos

- PikoCI en idle: ~15-30 MB RAM
- PikoCI ejecutando pipeline: ~50-100 MB RAM
- ggit mata el subproceso al hacer shutdown
- Health check cada 30s con restart automático
- Logs rotados, max 10 MB por archivo

### Naming

El proyecto se llama **ggit pikoci** — dando crédito a ambos proyectos:
- **ggit**: la plataforma Git auto-hospedada
- **PikoCI**: el motor de CI/CD portátil que corre como subproceso
