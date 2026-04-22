# Roadmap: Modernización (Go/Gin + Svelte)

## Fase 1: Análisis y Funcionalidades Base (Completado)

- [x] **Integrar detección avanzada de lenguajes:** Reemplazar sistema con `go-enry`.
- [x] **Resaltado de sintaxis en el servidor:** Integrar `chroma` y eliminar JS pesado en cliente.
- [x] **Auditoría de Settings:** Verificar flags `EnableIssues` y `EnablePulls`.

## Fase 2: Desacoplamiento de Lógica (Completado)

- [x] **Crear capa de Servicios:** Extraer lógica de Repos, Issues, PRs y Usuarios.
- [x] **Tipado de Entradas/Salidas:** Definir DTOs independientes del contexto HTTP.

## Fase 3: Migración del Motor HTTP a Gin (Completado)

- [x] **Inicialización y Middleware:** Migrar Auth, Sesiones, Logging y Recovery.
- [x] **Migración a API v2:** Implementar todos los endpoints JSON para el nuevo frontend.
- [x] **Bypass de Instalación:** Pre-configuración automática con SQLite.

## Fase 4: Frontend con SvelteKit (Completado)

- [x] **Migración de Pantallas:** Login, Dashboard, Vista de Repo, Colaboración y Admin.
- [x] **Sistema de Diseño:** Soporte avanzado para temas (Deep Night, Ayu, Monokai, Nord).
- [x] **El Fin del Legacy:** Desconexión total del motor de plantillas de Go (0 templates).

---

## Fase 10: Limpieza de "Grasa" Técnica (En Proceso)

**Objetivo:** Reducir el binario y el ruido de fondo del sistema (CPU/RAM).

### Infraestructura y Base de Datos

- [x] **Drivers Externos:** Eliminar soporte para MySQL, PostgreSQL y MSSQL. Mantener solo **SQLite3**.
- [x] **Optimización SQLite (Modo WAL):**
  - [x] Activar modo **WAL** (Write-Ahead Logging) para alta concurrencia.
  - [x] Configurar `Max Open Conns: 1` para evitar bloqueos de escritura.
  - [x] Configurar `Synchronous: NORMAL` para balance velocidad/seguridad.
- [x] **Unificación de Persistencia (XORM -> GORM):** Migrar el 100% de los modelos restantes. [COMPLETO]
- [x] **Gestión de Memoria (RAM):** [COMPLETO]
  - [x] **GORM .Select():** Optimizar listados para no cargar columnas pesadas en memoria.
  - [x] **Semaforización:** Limitar procesos `git grep` concurrentes para evitar OOM.
- [x] **Caché L2:** Desactivar el caché de nivel 2 del ORM para liberar RAM. [COMPLETO - Resuelto al migrar a GORM]

### Servicios de Fondo y Dependencias

- [x] **Limpieza de Dependencias:** Eliminar rastro de Macaron (`inject`, `binding`) del `go.mod`.
- [x] **Eliminar Mirrors:** Borrar la lógica de sincronización periódica de repos remotos.
- [x] **Eliminar Migraciones Masivas:** Transformado en **Smart Clone** (solo clonado Git, sin APIs externas). [COMPLETO]
- [-] **Eliminar LFS:** Quitar soporte para Large File Storage (`internal/lfsx`). (Pendiente de decisión - No eliminar aún).
- [x] **Update Checker:** Borrar la tarea que busca actualizaciones de Gogs. [COMPLETO]
- [x] **Captcha:** Reemplazado por sistema dual ultra-ligero (**Honeypot** + **Pregunta Configurable**). [COMPLETO]

### Autenticación y Comunicación

- [x] **Modernización de Dependencias (Core):**
  - [x] Eliminar `github.com/unknwon/com` (Sustituido por `internal/strx`). [COMPLETO]
  - [x] Eliminar `github.com/unknwon/cae` (Sustituido por `internal/archiver`). [COMPLETO]
  - [x] Eliminar `github.com/unknwon/i18n` (Sustituido por `internal/i18nx`). [COMPLETO]
- [x] **Simplificación de Autenticación:**
  - [x] Eliminar métodos complejos (PAM, SMTP-Auth, GitHub OAuth).
  - [x] Mantener **Local + LDAP**.
  - [ ] Implementar **OIDC** (OpenID Connect) como estándar moderno. A FUTURO
- [x] **Avatares:** Desactivar Gravatar/Libravatar. Usar avatares locales o generados.
- [x] **Webhooks:** Eliminar integraciones específicas (Slack, Discord, DingTalk). Mantener solo **Webhook Genérico (JSON)**. [COMPLETO]
- [x] **Email Minimalista:** Implementado estilo **Slim HTML** (ligero y con formato esencial). [COMPLETO]

---

## Fase 11: Estrategia "Docs-as-Code" (Nueva Wiki)

**Objetivo:** Eliminar la wiki tradicional a favor de contenido dentro del repositorio.

- [x] **Eliminar Wiki Legacy:** Borrar el sistema de repositorios independientes `.wiki.git`. [COMPLETO]
- [x] **Implementar "Wiki-Folder":** [COMPLETO]
  - En cada push, se miran los frontmatter de la documentacion.
  - [x] **Detector de Carpeta:** Activar pestaña Wiki si existe `/docs` o `.docs`. (carpeta seleccionable en ajustes del repositocio, por defecto `/docs`)
  - [x] **Sidebar Dinámico:** Generar menú lateral leyendo YAML/TOML Frontmatter de los `.md`.
  - [x] **Caché por Commit:** Solo recalcular el sidebar si el `CommitSHA` cambia.
  - [x] **Wikilinks:** Soporte para `[[Enlace]]` mediante Goldmark en el backend. (Global: Se ha migrado todo el motor de Markdown a Goldmark).

<!-- ---

title: Instalacion
order: 2 

---
se puede poner en todos el order, con 1 solo ya se posiciona todos los de la carpeta. si no esta puesto por orden alfabetico.  -->




## Fase 12: Funcionalidades de Precisión

- [x] **Internacionalización (I18n):** Carga dinámica de JSON locales en Frontend.
- [ ] **Búsqueda Avanzada de Repositorio:**
  - [x] Desactivar indexación profunda global (Bleve/CPU Killer).
  - [x] Implementar **Repo Search (On-demand):** Búsqueda interna vía `git grep` en tiempo real.
- [ ] **Web File Management (CRUD):** Creación, subida (Drag & Drop) y borrado de archivos vía web.

---

## Fase 13: Rebranding e Identidad Final

- [ ] **Refactor de Nomenclatura:** Actualizar `gogs.io/gogs` a la nueva identidad en todo el core.
- [x] **Rutas por defecto:** Cambiar carpetas de almacenamiento de `gogs-repositories` a `ggit-repositories`. [COMPLETO]
- [ ] **Component Library:** Consolidar la librería de componentes Svelte (M3 Expressive).

---






4. Caché del Sidebar por Commit

- Análisis: Construir la jerarquía del menú desde los .md es costoso. Solo debería hacerse una vez por cada Push.
- La mejor alternativa: github.com/hashicorp/golang-lru/v2
  - Por qué: Es una caché en memoria LRU (Least Recently Used) genérica (usando Generics de Go 1.18+) ultrarrápida y segura para
    concurrencia.
  - La lógica sería: clave = RepoID_CommitSHA. Si un usuario entra a la Wiki, buscas esa clave en la LRU. Si el SHA del último
    commit de la rama principal no ha cambiado, devuelves el Sidebar precalculado en $O(1)$. Si cambia (alguien hace push), la clave
    antigua desaparecerá naturalmente de la caché.

Resumen de viabilidad:
GORM .Select() (Prioridad 🔥):

Viabilidad: Total. Es la mejora con mayor impacto inmediato y menor riesgo.
Acción: Refactorizar los endpoints de listas para no cargar columnas pesadas (como descripciones largas o configuraciones) si solo vamos a mostrar el nombre y el ID.
Gestión de Procesos (Semaforización) (Prioridad ⚡):

Viabilidad: Muy alta.
Acción: Implementar un límite (semáforo) para procesos pesados como git grep. Esto evita que una ráfaga de usuarios consuma toda la RAM del sistema operativo simultáneamente.
go-git/v5 para Metadatos (Prioridad 🛠):

Viabilidad: Media/Alta.
Acción: Usarlo específicamente para el escaneo de la carpeta /docs (Wiki) y el árbol de archivos rápido. Para operaciones pesadas (grep, diff), seguiremos usando procesos externos para no saturar el Heap de Go.
Eliminación de Mirrors:

Estado: Ya está prácticamente completado. He verificado que la lógica de sincronización ya no está en cron