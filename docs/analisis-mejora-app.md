# Análisis de mejora — ggit

> Basado en revisión de código completa del repositorio `/home/nef734/Documentos/Proyectos/ggit`.
> Fecha: 2026-06-19

---

## 🔴 Crítico

| # | Problema | Archivo | Impacto |
|---|----------|---------|---------|
| 1 | SECRET_KEY hardcodeado `RhgOvRxKpJ1A8ub` | `conf/app.ini:153` | Cualquiera con acceso al repo puede descifrar cookies, sesiones, 2FA |
| 2 | Versiones irreales: Go `1.25.9`, TS `6.0.2`, Vite `8.0.7` | `go.mod:3`, `package.json:31-32` | `go mod tidy` y `npm install` fallan para otros entornos |
| 3 | XSS: `Safe()` y `Str2HTML()` sin sanitizar | `template/template.go:131-136` | Permiten HTML arbitrario en respuestas (ya corregido en Sprint 1) |
| 4 | Sin tipos compartidos frontend/backend | `api.ts` usa `any`; backend 637 líneas de tipos sin conectar | Cada cambio de API hay que reflejarlo a mano |
| 5 | CI/CD roto para el fork | `.github/workflows/*.yml` apunta a `gogs/gogs` | Releases y Docker nunca se construyen para `nef734/ggit` |

## 🟠 Alto impacto

### Backend

| # | Problema | Archivo | Detalle |
|---|----------|---------|---------|
| 6 | 65+ referencias Gogs residuales | Múltiples archivos | Comentarios, imports, configs, CI apuntan a `gogs.io/gogs` |
| 7 | Error wrapping inconsistente | `database/*.go` (~100 ocurrencias) | `%v` en vez de `%w` → `errors.Is/As` no funciona |
| 8 | `context.TODO()` en producción (~35 sitios) | `database/*.go` | Operaciones no cancelables, sin deadline |
| 9 | Swagger obsoleto | `docs/docs.go` vs `v2.go` | Documenta ~15 de ~150 endpoints reales, rutas con prefijo incorrecto |
| 10 | Sin límite de request body | `route/api/v2/*.go` | DoS por JSON gigante |
| 11 | Respuestas API inconsistentes | `context/api.go:28,37,49` | Unas usan `{message}`, otras `{ok: true}`, otras objetos planos |

### Frontend

| # | Problema | Archivo | Detalle |
|---|----------|---------|---------|
| 12 | Sin AbortController en llamadas API | Todos los componentes con `apiFetch` | Memory leaks, race conditions al navegar rápido |
| 13 | `getContext('repo')` tipado como `any` | 25+ componentes | Cero type-safety en el dato compartido más crítico |
| 14 | 8 componentes >400 líneas | `RepoFileBrowser` (921), `WikiView` (865), `IssueDetail` (683), `Navigation` (674), `PRCreator` (615), `DashboardView` (525), `DiffViewer` (512), `repo/create/+page.svelte` (1069) | Violan SRP, difíciles de testear |
| 15 | `timeAgo` duplicado en 4 componentes | `IssueDetail`, `RepoRow`, `RepoFileBrowser`, `commits/+page` | Debería ser utilidad compartida |
| 16 | `$derived(() => ...)` incorrecto | `RepoRow`, `RepoFileBrowser` | Envuelve función en vez de valor; debería ser `$derived.by()` |
| 17 | Hardcoded strings sin i18n | `login/+page.svelte`, `repo/create/+page.svelte`, `+error.svelte` | Cientos de líneas de UI sin traducción |
| 18 | `console.log` de debug en producción | `+layout.svelte:34-55`, `DashboardView.svelte:39-51` | Fuga de info, ruido en consola |
| 19 | Sin ESLint/Prettier/Husky | — | Sin formato automático, sin hooks pre-commit |

## 🟡 Medio

| # | Problema | Archivo | Detalle |
|---|----------|---------|---------|
| 20 | Race condition: `NextIssueIndex` sin mutex | `database/repo.go:1010` | Dos issues pueden recibir el mismo índice |
| 21 | N+1 query: `DeleteRepository` | `database/repo.go:2021-2033` | Carga issues uno por uno para borrar comentarios |
| 22 | `IsErrUploadNotExist` chequea tipo incorrecto | `database/error.go:277` | `ErrAttachmentNotExist` en vez de `ErrUploadNotExist` |
| 23 | Panic en producción | `database/database.go:119`, `context/repo.go:88` | `panic("unreachable")` en paths alcanzables |
| 24 | Gorutinas sin límite en `getRepoDirectory` | `route/api/v2/repo.go:348-358` | Una gorutina por entrada de directorio |
| 25 | SecretKey usado sin domain separation | `tool/tool.go:111` | SHA1(data + SecretKey) para códigos de verificación |
| 26 | Cookies CSRF no HttpOnly ni Secure | `ginx/csrf.go:53` | Token accesible desde JS, enviado por HTTP plano |
| 27 | Reverse proxy auto-registra usuarios | `context/auth.go:289-316` | Atacante con SSRF puede crear cuentas admin |
| 28 | Database password en texto plano al instalar | `route/api/v2/install.go:135` | Escrita a `app.ini` sin cifrar |
| 29 | Toast con `duration` muerto | `toast.svelte.ts:7,17` | Se acepta duración pero nunca se auto-descartan |
| 30 | Solo 12 `aria-label` en 33k LOC | Todos los componentes | Accesibilidad muy baja |
| 31 | `role="button"` en `<div>` sin soporte teclado | 8 componentes | Espacio no activa el botón |
| 32 | CSS: dos sistemas de variables mezclados | Todos los componentes | `--text-color` (Gogs legacy) y `--md-sys-color-*` (MD3) conviven |
| 33 | CSS: selectores muertos | `RepoRow`, varios | Clases que no existen en el DOM |

## 🔵 Bajo

| # | Problema | Archivo | Detalle |
|---|----------|---------|---------|
| 34 | Artefactos del migration | — | Resuelto: se eliminaron `fix_imports.py`, `migrate_errors.py`, `scratch_test_ignores.go`, `SELECT`, `temporary_features.md`, el `package-lock.json` raíz y el symlink `log -> /data/log`, que al clonar creaba un enlace a una ruta inexistente en otra máquina |
| 35 | Debug prints en producción | `ginx/template.go:17,53,58` | `fmt.Printf("[DEBUG]")` en template setup |
| 36 | Trace logging excesivo | `route/api/v2/repo.go` | 15+ `log.Trace()` en handler de alta frecuencia |
| 37 | README desactualizado | `README.md` | Referencias a `ggit.io`, enlaces rotos |
| 38 | Sin `.env.example` | — | No hay mecanismo estándar para secrets |
| 39 | Sin `.go-version` / `.nvmrc` | — | No hay pinning de versiones de herramientas |
| 40 | Dependabot solo para Go, no npm | `.github/dependabot.yml` | Frontend no recibe actualizaciones automáticas |

---

## Sprints sugeridos

| Sprint | Área | Estado |
|--------|------|--------|
| 1 | Template cleanup (Safe/Sanitize/Str2HTML) | ✅ |
| 2 | SECRET_KEY + versiones irreales | ⏳ |
| 3 | CI/CD: workflows apuntar a nef734/ggit | ⏳ |
| 4 | ESLint + Prettier + Husky frontend | ⏳ |
| 5 | AbortController + tipar getContext('repo') | ⏳ |
| 6 | Dividir componentes grandes | ⏳ |
| 7 | Backend: %v→%w, context, body limits | ⏳ |
| 8 | i18n completo + console.log + $derived.by | ⏳ |
| 9 | Referencias Gogs (65+ archivos) | ✅ |
| 10 | Variables CSS + aria-labels | ⏳ |

### Extras completados

| Sprint | Área |
|--------|------|
| A1-A9 | Git operations layer (merge locking, gorutinas, context, TOCTOU, errores, rate limit) |
| B1-B7 | Database N+1, context.TODO, IsErrUploadNotExist, NextIssueIndex, full table scans |
| D1 | highlight.js → `$lib/hljs` |
| D2 | lucide-svelte imports individuales (71 archivos) |
| D3 | Google Fonts duplicado eliminado |
| — | Parpadeo blanco corregido (stopPropagation → preventDefault) |
| — | hooks.client.ts: reload removido en error de chunks |

---

## Análisis A: Git operations layer

Cubre cómo la app interactúa con repositorios git: clonado, push, receive-pack, hooks, LFS, merge conflicts, concurrencia.

| # | Problema | Archivo | Detalle |
|---|----------|---------|---------|
| A1 | Sin rate limiting en operaciones git | `internal/route/repo/repo.go` | Push/clone sin límite de concurrencia, DoS via git operations |
| A2 | Locking inconsistente en merges | `database/pull.go:455` | FIXME documentado: "make sure two merges does not happen at same time" |
| A3 | LFS sin autenticación en batch | `internal/route/lfs/batch.go` | Verificar si LFS batch requiere auth o permite acceso anónimo |
| A4 | Git hooks sin validación de contenido | `internal/route/repo/repo.go` | Hooks ejecutados del lado servidor sin sanitizar |
| A5 | Sin límite de tamaño en objetos LFS | `internal/route/lfs/` | No hay enforce de `lfs.maxobjectsize` |
| A6 | Gorutinas sin control en `getRepoDirectory` | `route/api/v2/repo.go:348-358` | Una gorutina por entrada → DoS en repos grandes |
| A7 | Sin verificación de integridad en receive-pack | `internal/route/repo/repo.go` | No se validan objetos antes de aceptar push |

## Análisis B: Base de datos — esquema, migraciones, índices

Cubre modelos GORM, migraciones, índices, integridad referencial, datos huérfanos.

| # | Problema | Archivo | Detalle |
|---|----------|---------|---------|
| B1 | Sin foreign keys explícitas en GORM | `database/*.go` | GORM no crea FKs por defecto, integridad referencial solo a nivel app |
| B2 | Migraciones sin versionado claro | `database/migrations/migrations.go` | Migraciones con comments de Gogs v0.7.33, v0.9.141, v0.12.0 — difíciles de seguir |
| B3 | Carga de issues uno por uno | `database/repo.go:2021-2033` | N+1 query en DeleteRepository (comentarios y attachments por issue) |
| B4 | Orphan cleanup no implementado | `database/repo.go` | No hay garbage collection de attachments huérfanos, comentarios de issues eliminados |
| B5 | `context.TODO()` en ~35 queries | `database/*.go` | Queries sin contexto de cancelación ni timeout |
| B6 | Race condition en issue indexing | `database/repo.go:1010` | FIXME: sin mutex, dos issues pueden recibir mismo índice |
| B7 | Sin índices compuestos en tablas grandes | `database/models/` | Falta revisar índices en `issue`, `pull_request`, `comment`, `action` |
| B8 | Carga completa de objetos en vez de IDs | `database/org.go:398` | FIXME: "only need to get IDs here, not all fields of repository" |

## Análisis C: Auth & permisos — flujo completo

Cubre login, sesiones, OIDC, 2FA, tokens API, SSH keys, modelo de permisos, middleware.

| # | Problema | Archivo | Detalle |
|---|----------|---------|---------|
| C1 | Reverse proxy auto-registra usuarios | `context/auth.go:289-316` | Atacante con SSRF puede crear cuentas admin |
| C2 | CSRF cookie no HttpOnly ni Secure | `ginx/csrf.go:53` | Token accesible desde JS, enviado por HTTP plano |
| C3 | SecretKey usado sin domain separation | `tool/tool.go:111` | SHA1(data + SecretKey) para códigos, mismo key para todo |
| C4 | 2FA sin rate limiting | `route/api/v2/user_settings.go:698` | Fuerza bruta contra 2FA sin protección |
| C5 | Sesión `gogs_session` nombre legacy | `ginx/ginx.go:52` | Ya renombrado a `ggit_session` en Sprint de Gogs cleanup |
| C6 | Token API sin expiración forzada | `database/access_token.go` | Verificar si tokens tienen TTL o son eternos |
| C7 | Permisos chequeados por duplicado | `context/repo.go` | Middleware `RequireRepoAdmin` + checks manuales en handlers |
| C8 | OIDC sin validación de `nonce` | `route/api/v2/auth.go` | Posible replay attack si falta nonce |

## Análisis D: Frontend — rendimiento y bundle

Cubre bundle size, code splitting, lazy loading, dependencias pesadas, caché SPA, carga inicial.

| # | Problema | Archivo | Detalle | Estado |
|---|----------|---------|---------|--------|
| D1 | highlight.js + marked embebidos sin lazy loading | `package.json:37,39` | ~200KB+ de syntax highlighter y markdown parser en bundle principal | ✅ |
| D2 | lucide-svelte sin tree-shaking completo | `package.json:25` | Iconos importados como individuales (71 archivos) | ✅ |
| D3 | Google Fonts duplicado | `base.css:1` | `@import` eliminado, carga solo desde `app.html` | ✅ |
| D4 | Sin lazy loading de rutas | `routes/` | SvelteKit genera 62 route nodes por ruta | ✅ |
| D5 | Sin compresión Brotli en build | `vite.config.ts` | `vite-plugin-compression` agregado, genera `.br` | ✅ |
| D6 | Carga inline de fuentes/iconos | `app.html` | Ya óptimo (preconnect + display=swap) | ✅ |
| D7 | Sin métricas de performance | — | No hay Web Vitals, ni lazy loading de imágenes, ni responsive images | ⏳ |
| D8 | Chunk splitting manual sin verificar | `vite.config.ts:16-22` | marked y dompurify separados, warning bajado a 500KB | ✅ |
