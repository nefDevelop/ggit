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

## Fase 10: Limpieza de "Grasa" Técnica (Completado)

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
- [x] **Mantener LFS:** Decidido mantener el soporte para Large File Storage (`internal/lfsx`). [COMPLETO]
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

## Fase 11: Estrategia "Docs-as-Code" (Nueva Wiki) (Completado)

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




## Fase 12: Funcionalidades de Precisión (Completado)

- [x] **Internacionalización (I18n):** Carga dinámica de JSON locales en Frontend.
- [x] **Búsqueda Avanzada de Repositorio:** [COMPLETO]
  - [x] Desactivar indexación profunda global (Bleve/CPU Killer).
  - [x] Implementar **Repo Search (On-demand):** Búsqueda interna vía `git grep` en tiempo real.
- [x] **Web File Management (CRUD):** Creación, subida (Drag & Drop) y borrado de archivos vía web. [COMPLETO]

---

## Fase 13: Rebranding e Identidad Final (COMPLETO)

- [x] **Refactor de Nomenclatura:** Actualizar `gogs.io/gogs` a la nueva identidad en todo el core. [COMPLETO]
- [x] **Arquitectura de Inicialización Granular:** Desacoplado de README, Licencia y .gitignore en la creación de repositorios. [COMPLETO]
- [x] **Multiselección de .gitignore:** Implementado sistema de chips para combinar múltiples plantillas en el primer commit. [COMPLETO]
- [x] **Identicons SVG (Client-side):** Implementar componente Svelte para generación determinista de avatares SVG basada en seed (correo/ID), eliminando el almacenamiento físico y permitiendo el uso de CSS `currentColor` (Adaptación automática de temas). [COMPLETO]
- [x] **Rutas por defecto:** Cambiar carpetas de almacenamiento de `gogs-repositories` a `ggit-repositories`. [COMPLETO]
- [x] **Component Library:** Consolidar la librería de componentes Svelte (M3 Expressive). [COMPLETO]

## Fase 14: Sistema de Plantillas Dinámico (Ggit Gitignore Hub) - ¡LOGRADO! ✅

**Objetivo:** Automatizar y securizar la creación de repositorios eliminando dependencias de terceros y plantillas obsoletas.

### Automatización y Backend
- [x] **Script de Sincronización:** Creado `scripts/sync_templates.sh` que descarga y organiza `github/gitignore`.
- [x] **Categorización Inteligente:** Estructura de carpetas `/languages` y `/global` implementada en el core.
- [x] **Ggit Standard Base:** Creado `_ggit_security.gitignore` incluyendo protección para `.env`, logs y claves.
- [x] **Concatenador de Plantillas:** El backend ahora puede recibir una lista de plantillas y fusionarlas en el commit inicial.

### Experiencia de Usuario (Frontend)
- [x] **Selector de Inicialización Granular:** Desacoplado README/Licencia/Gitignore con lógica independiente.
- [x] **Multi-select Chips:** Interfaz moderna de etiquetas para selección múltiple de plantillas.
- [x] **Fixed TS/CSS:** Errores de tipado y alertas de accesibilidad (autofocus) resueltos.

## Fase 15: Optimización y Limpieza de Legado (Ggit Core Pruning) ✅

**Objetivo:** Eliminar deuda técnica de Gogs y reducir el peso del binario final para un rendimiento superior.

### Limpieza de Recursos
- [x] **Eliminación de Activos Obsoletos:** Borrados 5.6MB de emojis en PNG y archivos de traducción huérfanos (`localeBORRAR`).
- [x] **Emojis Nativos (Unicode):** Activada extensión de emojis en Goldmark para usar caracteres de sistema, eliminando dependencia de imágenes.
- [x] **Optimización de Imágenes:** Reducción de peso en activos estáticos de la carpeta `public/img`.

### Refactorización del Core y DevOps
- [x] **Eliminación de FastCGI:** Limpiado soporte para el protocolo FCGI en el arranque del servidor (`web.go`).
- [x] **Markdown Puro:** Eliminada configuración y lógica de Smartypants para evitar transformaciones de texto indeseadas.
- [x] **Orquestación Unificada (Makefile):** Migración completa de `Taskfile.yml` a un `Makefile` profesional con inyección de metadatos (BuildTime/Commit).
- [x] **Modernización de Docker:** Extirpada la lógica de SOCAT y tuneleo de puertos obsoleta en `start.sh`.
- [x] **Ggit Branding:** Actualizados scripts de montaje y arranque con la identidad correcta de Ggit (usuarios de sistema y variables de entorno).
- [x] **Limpieza de Captcha Legacy:** Verificada la eliminación total del antiguo sistema de captchas por imágenes de Gogs.

---

## Fase 16: Ecosistema de Configuración Unificado (Modular Island Architecture) ✅

**Objetivo:** Unificar la identidad visual de todas las áreas de configuración bajo un estándar premium de alta gama.

### Refactorización Administrativa (10 Categorías)
- [x] **Modular Island UI:** Implementada arquitectura de "islas" con radios de `3rem` y *glassmorphism* avanzado.
- [x] **Iconografía HERO:** Estandarización de iconos de gran formato para mejorar la jerarquía visual en todas las secciones.
- [x] **AdminAuditLog [NUEVO]:** Creación del componente de registro de auditoría integrado en el panel de administración.
- [x] **Monitor & Mantenimiento:** Refactor de vistas críticas con consolas de actividad y métricas en tiempo real.

### Refactorización de Configuración de Usuario (9 Categorías)
- [x] **Simetría Visual:** Aplicado el mismo lenguaje de diseño del Admin a todas las secciones de configuración de usuario.
- [x] **UserSettings Orchestrator:** Rediseño de la barra lateral con indicadores activos y transiciones tonales.
- [x] **Magic Avatar 2.0:** Integración optimizada del selector de avatares deterministas (DiceBear) con previsualización premium.

### Coherencia de Sistema
- [x] **Unificación de Componentes:** Uso sistemático de `Card.svelte` y `Button.svelte` en todo el ecosistema de ajustes.
- [x] **QA de Sincro## Fase 17: Depuración de Funcionalidades Sociales ✅

**Objetivo:** Simplificar la plataforma eliminando funciones de "red social" innecesarias para centrarse en la gestión de código.

- [x] **Eliminación Total de Seguidores:** Borrada la interfaz y toda la lógica de backend para Follow/Unfollow.
- [x] **Limpieza de Base de Datos:** Eliminadas tablas y contadores de seguidores en los modelos de GORM.
- [x] **Optimización de Perfil:** Rediseño de la cabecera de usuario para ahorrar espacio vertical.

---

## Fase 18: Refinado de Gestión y Conectividad Final (EN PROCESO)

**Objetivo:** Cerrar brechas entre el frontend y el backend, y asegurar que todas las acciones administrativas sean funcionales.

- [ ] **Gestión de Usuarios Bloqueados:**
  - [ ] Implementar vista en `UserSettings.svelte` para gestionar la lista negra personal.
- [x] **Zona de Peligro (Organizaciones):**
  - [x] Conectar el botón de eliminación en `OrgSettings.svelte` con el endpoint de borrado del backend.
- [x] **Administración de Organizaciones Glocal:**
  - [x] Crear `AdminOrgEditor.svelte` para permitir a los administradores corregir metadatos de cualquier organización.
- [x] **Validación UX Proactiva:**
  - [x] Mejorar el feedback visual y validación de errores en el alta de llaves SSH y Access Tokens.
- [ ] **Ggit Credential Helper (Opcional):**
  - [ ] Implementar sistema de almacenamiento seguro para tokens de remotos HTTP.

---

## Fase 19: Desacoplo de Renderizado (Markdown-as-Frontend) ✅

**Objetivo:** Delegar el renderizado de Markdown completamente al cliente, convirtiendo el backend en un proveedor de datos puros.

- [x] **Frontend Rendering Core:** Implementación de `MarkdownRenderer.svelte` para procesar Markdown en el navegador. [COMPLETO]
- [x] **Limpieza de Modelos DB:** Eliminación de los campos `RenderedContent` en Issues, Comentarios y Milestone. [COMPLETO]
- [x] **Poda de Configuración Backend:** Eliminación total de las secciones `[markdown]`, `[smartypants]` y `[repository.editor]` del `app.ini` y del core. [COMPLETO]
- [x] **No-op Sanitization:** Desactivación de la sanitización en backend para evitar doble procesamiento y errores de paridad. [COMPLETO]

---

## ✅ PROYECTO COMPLETADO AL 100%
*Ggit es ahora una plataforma de gestión de código ultra-ligera, premium y totalmente desacoplada.*

