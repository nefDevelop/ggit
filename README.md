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
| Persistencia | XORM | **GORM**, 13 migraciones nuevas |
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
explícitamente en 13 migraciones (v20 a v33) con `AutoMigrate`, verificadas contra
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

## Documentación

`docs/` contiene la documentación del producto. El sitio que se despliega desde
ahí usa VitePress.

## Licencia

MIT. Ver [LICENSE](LICENSE), que conserva el aviso de copyright original de Gogs
tal como exige esa licencia.
