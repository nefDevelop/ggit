# Sprint 1 — Limpieza de código muerto (template engine)

## Objetivo

Eliminar funciones inseguras y código no utilizado del paquete `internal/template/`. El motor de templates HTML del lado servidor está desconectado desde la migración a SvelteKit SPA.

## Cambios realizados

### 1. `Safe()` — Eliminada

- **Archivo:** `internal/template/template.go:131-133`
- **Motivo:** Función sin ningún caller en todo el código base. Convertía cualquier string a `template.HTML` sin sanitizar, permitiendo XSS si alguien la usara en el futuro.
- **Acción:** Eliminada la función y su registro en FuncMap.

### 2. `Sanitize` — De no-op a `html.EscapeString`

- **Archivo:** `internal/template/template.go:68`
- **Antes:** `func(s string) string { return s }` (devolvía el input sin cambios, siendo un sanitizador nulo)
- **Después:** `html.EscapeString` (escapa `& < > " '`)
- **Impacto:** Template engine desconectado, este cambio no afecta la app en producción.

### 3. `Str2HTML` — Ahora escapa HTML

- **Archivo:** `internal/template/template.go:135-137`
- **Antes:** `template.HTML(raw)` (pasaba crudo)
- **Después:** `template.HTML(html.EscapeString(raw))`
- **Nota:** Solo se usa en templates desconectados. La función homónima en `internal/email/email.go` NO se modificó porque el engine de emails sigue activo.

### 4. Import `"html"` agregado

Necesario para `html.EscapeString`.

## Verificación

- `go build ./internal/template/` → OK
- `go vet ./internal/template/` → OK
