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
- [x] **Caché L2:** Desactivar el caché de nivel 2 del ORM para liberar RAM. [COMPLETO - Resuelto al migrar a GORM]

### Servicios de Fondo y Dependencias

- [x] **Limpieza de Dependencias:** Eliminar rastro de Macaron (`inject`, `binding`) del `go.mod`.
- [x] **Eliminar Mirrors:** Borrar la lógica de sincronización periódica de repos remotos.
- [x] **Eliminar Migraciones Masivas:** Transformado en **Smart Clone** (solo clonado Git, sin APIs externas). [COMPLETO]
- [-] **Eliminar LFS:** Quitar soporte para Large File Storage (`internal/lfsx`). (Pendiente de decisión - No eliminar aún).
- [x] **Update Checker:** Borrar la tarea que busca actualizaciones de Gogs. [COMPLETO]
- [x] **Captcha:** Reemplazado por sistema dual ultra-ligero (**Honeypot** + **Pregunta Configurable**). [COMPLETO]

### Autenticación y Comunicación

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
- [ ] **Implementar "Wiki-Folder":**
  - [ ] **Detector de Carpeta:** Activar pestaña Wiki si existe `/docs` o `.docs`.
  - [ ] **Sidebar Dinámico:** Generar menú lateral leyendo YAML/TOML Frontmatter de los `.md`.
  - [ ] **Caché por Commit:** Solo recalcular el sidebar si el `CommitSHA` cambia.
  - [ ] **Wikilinks:** Soporte para `[[Enlace]]` mediante Goldmark en el backend.
  - En cada push, se miran los frontmatter de la documentacion. con

---
title: Instalacion
parent: Introduccion
y algun metodo mas? como se crearia ucon frontm

---

## Fase 12: Funcionalidades de Precisión

- [x] **Internacionalización (I18n):** Carga dinámica de JSON locales en Frontend.
- [ ] **Búsqueda Avanzada de Repositorio:**
  - [x] Desactivar indexación profunda global (Bleve/CPU Killer).
  - [ ] Implementar **Repo Search (On-demand):** Búsqueda interna vía `git grep` en tiempo real.
- [ ] **Web File Management (CRUD):** Creación, subida (Drag & Drop) y borrado de archivos vía web.

---

## Fase 13: Rebranding e Identidad Final

- [ ] **Refactor de Nomenclatura:** Actualizar `gogs.io/gogs` a la nueva identidad en todo el core.
- [x] **Rutas por defecto:** Cambiar carpetas de almacenamiento de `gogs-repositories` a `ggit-repositories`. [COMPLETO]
- [ ] **Component Library:** Consolidar la librería de componentes Svelte (M3 Expressive).
