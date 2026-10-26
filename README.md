# ggit

Forge Git autoalojada. **Es un fork de [Gogs](https://github.com/gogs/gogs)**, no un
proyecto independiente: la mayor parte del backend es código de Gogs (MIT). Lo que
hay aquí es el trabajo de migración y las functionalities añadidas. Ver
[NOTICE](NOTICE) para el desglose y [LICENSE](LICENSE) para la licencia.

Sin relación institucional con el proyecto Gogs.

## Qué es de Gogs y qué es de aquí

| | Gogs upstream | ggit |
| :--- | :--- | :--- |
| Framework HTTP | Macaron | **Gin v1.12** |
| Persistencia | XORM | **GORM**, 14 migraciones nuevas |
| Bases de datos soportadas | MySQL, PostgreSQL, SQLite, MSSQL | **Solo SQLite**, con WAL y un único escritor |
| Frontend | Plantillas Go renderizadas en servidor | **`frontend/` SvelteKit 2 + Svelte 5**, 205 componentes, desacoplado |
| API | v1 | **v2** |
| Acceso a datos | Directo desde los handlers | **`internal/services/`** |

El frontend no tiene equivalente en Gogs: aquí es una aplicación aparte que habla
con la API por HTTP, con Playwright para E2E, Vitest para unitarios y MSW para
mocking de red.

### Migraciones de persistencia

XORM → GORM no fue un reemplazo de biblioteca. Cada Entidad de Gogs usa tipos de
columna propios que GORM no conoce, así que las columnas se reconstruyeron
explícitamente en 14 migraciones (v20 a v34) con `AutoMigrate`, verificadas contra
la evolution del esquema de Gogs.

## Funcionalidades añadidas

- **Code Review.** Aprobaciones por turnos y comentarios en línea sobre el diff,
  con protección de rama. `internal/database/pull_review.go` y
  `pull_review_comment.go`, con sus tests.
- **Docs-as-Code.** El wiki se sirve desde un directorio `/docs` del propio
  repositorio, en lugar de un `.wiki.git` paralelo que había que sincronizar a
  mano.
- **Secret Guardian.** Detección de claves de API y `.env` en los pushes, con
  pagina de configuración en el repo. Son unas 350 líneas entre servicio, test y
  rutas.
- **Kanban e Inbox unificado.** Vistas de tablero y bandeja multi-repositorio
  sobre la API de incidencias existente, no un subsistema aparte.

## Puesta en marcha

```bash
# Backend
go run main.go

# Frontend (en otra terminal)
cd frontend && npm install && npm run dev
```

Con Docker:

```bash
docker compose up -d
```

## Cliente CLI

`cmd/ggit-cli` es un cliente de la API v2, separado del binario del servidor:

```bash
make build-cli
export GGIT_HOST=https://ggit.734038.xyz
# Crea un token con tu contraseña y lo deja guardado en la configuración
.bin/ggit-cli auth login --username pi --password -
.bin/ggit-cli auth status
.bin/ggit-cli repo list
.bin/ggit-cli repo create mi-repo --private --auto-init
.bin/ggit-cli repo create otro --org acme
.bin/ggit-cli issue list pi/ggit --state open
.bin/ggit-cli pr list pi/ggit
.bin/ggit-cli pr create pi/ggit --title "Título" --head feature/x
.bin/ggit-cli issue create pi/ggit --title "Título" --body "Texto"
.bin/ggit-cli issue comment pi/ggit 12 --body-file nota.md
.bin/ggit-cli issue edit pi/ggit 12 --title "Nuevo título" --assignee 3
.bin/ggit-cli issue close pi/ggit 12
.bin/ggit-cli pr review pi/ggit 4 approve --body "LGTM"
.bin/ggit-cli pr merge pi/ggit 4 --yes
.bin/ggit-cli repo delete mi-repo --yes
.bin/ggit-cli auth logout
```

La configuración también puede vivir en `~/.config/ggit-cli/config.json`
(`{"host": "...", "token": "..."}`); las variables de entorno tienen prioridad.
Todos los comandos aceptan `--json`.

## CI con PikoCI

ggit integra [PikoCI](https://github.com/PikoCI/pikoci) como motor de CI. El servidor de
PikoCI corre aparte (recomendado: contenedor) y ggit publica el pipeline de cada repositorio
desde su `.ggit/ci.hcl`.

1. Configura la sección `[pikoci]` de `app.ini`:

```ini
[pikoci]
ENABLED       = true
URL           = http://127.0.0.1:18080
PUBLIC_URL    = https://ci.example.com
TEAM          = main
API_TOKEN     = <token de la API de PikoCI>
BOT_TOKEN     = <token del usuario bot de ggit>
PIPELINE_FILE = .ggit/ci.hcl
```

2. Crea el usuario bot y su token (se muestra una sola vez):

```bash
ggit admin pikoci bot --repo pi/ggit
```

3. Copia el token en `BOT_TOKEN` y crea un token de API en la interfaz de PikoCI para
   `API_TOKEN`.

4. Añade `.ggit/ci.hcl` al repositorio. La plantilla por defecto, con el reporte de estado
   ya incluido, está en `GET /api/v2/admin/pikoci/pipeline-template`.

5. Activa CI en **Ajustes → CI/CD** del repositorio. Desde ahí puedes republicar el
   pipeline o desactivarlo. Los pushes a la rama por defecto disparan el build y el estado
   aparece en el commit y en el PR.

Referencia completa: [docs/advancing/pikoci-ci.mdx](docs/advancing/pikoci-ci.mdx).

## Documentación

`docs/` contiene la documentación del producto. El sitio que se despliega desde
ahí usa Mintlify (`docs/docs.json`).

## Licencia

MIT. Ver [LICENSE](LICENSE), que conserva el aviso de copyright original de Gogs
tal como exige esa licencia.
