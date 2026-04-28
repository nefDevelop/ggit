# Roadmap & Changelog: Modernización (Go/Gin + Svelte)

## Fase 1: Análisis y Funcionalidades Base ✅
- [x] **Integrar detección avanzada de lenguajes:** Reemplazar sistema con `go-enry`.
- [x] **Resaltado de sintaxis en el servidor:** Integrar `chroma` y eliminar JS pesado en cliente.
- [x] **Auditoría de Settings:** Verificar flags `EnableIssues` y `EnablePulls`.

## Fase 2: Desacoplamiento de Lógica ✅
- [x] **Crear capa de Servicios:** Extraer lógica de Repos, Issues, PRs y Usuarios.
- [x] **Tipado de Entradas/Salidas:** Definir DTOs independientes del contexto HTTP.

## Fase 3: Migración del Motor HTTP a Gin ✅
- [x] **Inicialización y Middleware:** Migrar Auth, Sesiones, Logging y Recovery.
- [x] **Migración a API v2:** Implementar todos los endpoints JSON para el nuevo frontend.
- [x] **Bypass de Instalación:** Pre-configuración automática con SQLite.

## Fase 4: Frontend con SvelteKit ✅
- [x] **Migración de Pantallas:** Login, Dashboard, Vista de Repo, Colaboración y Admin.
- [x] **Sistema de Diseño:** Soporte avanzado para temas (Deep Night, Ayu, Monokai, Nord).
- [x] **El Fin del Legacy:** Desconexión total del motor de plantillas de Go (0 templates).

---

## Fase 10: Limpieza de "Grasa" Técnica ✅
**Objetivo:** Reducir el binario y el ruido de fondo del sistema (CPU/RAM).

### Infraestructura y Base de Datos
- [x] **Drivers Externos:** Eliminar soporte para MySQL, PostgreSQL y MSSQL. Mantener solo **SQLite3**.
- [x] **Optimización SQLite (Modo WAL):**
  - [x] Activar modo **WAL** (Write-Ahead Logging) para alta concurrencia.
  - [x] Configurar `Max Open Conns: 1` para evitar bloqueos de escritura.
  - [x] Configurar `Synchronous: NORMAL` para balance velocidad/seguridad.
- [x] **Unificación de Persistencia (XORM -> GORM):** Migrar el 100% de los modelos restantes.
- [x] **Gestión de Memoria (RAM):**
  - [x] **GORM .Select():** Optimizar listados para no cargar columnas pesadas en memoria.
  - [x] **Semaforización:** Limitar procesos `git grep` concurrentes para evitar OOM.
- [x] **Caché L2:** Desactivar el caché de nivel 2 del ORM para liberar RAM. (Resuelto al migrar a GORM)

### Servicios de Fondo y Dependencias
- [x] **Limpieza de Dependencias:** Eliminar rastro de Macaron (`inject`, `binding`) del `go.mod`.
- [x] **Eliminar Mirrors:** Borrar la lógica de sincronización periódica de repos remotos.
- [x] **Eliminar Migraciones Masivas:** Transformado en **Smart Clone** (solo clonado Git, sin APIs externas).
- [x] **Mantener LFS:** Decidido mantener el soporte para Large File Storage (`internal/lfsx`).
- [x] **Update Checker:** Borrar la tarea que busca actualizaciones de Gogs.
- [x] **Captcha:** Reemplazado por sistema dual ultra-ligero (**Honeypot** + **Pregunta Configurable**).

### Autenticación y Comunicación
- [x] **Modernización de Dependencias (Core):**
  - [x] Eliminar `github.com/unknwon/com` (Sustituido por `internal/strx`).
  - [x] Eliminar `github.com/unknwon/cae` (Sustituido por `internal/archiver`).
  - [x] Eliminar `github.com/unknwon/i18n` (Sustituido por `internal/i18nx`).
- [x] **Simplificación de Autenticación:**
  - [x] Eliminar métodos complejos (PAM, SMTP-Auth, GitHub OAuth).
  - [x] Mantener **Local + LDAP**.
  - [ ] **OIDC (OpenID Connect):** Implementar como estándar moderno. (A FUTURO)
- [x] **Avatares:** Desactivar Gravatar/Libravatar. Usar avatares locales o generados.
- [x] **Webhooks:** Eliminar integraciones específicas (Slack, Discord, DingTalk). Mantener solo **Webhook Genérico (JSON)**.
- [x] **Email Minimalista:** Implementado estilo **Slim HTML** (ligero y con formato esencial).

---

## Fase 11: Estrategia "Docs-as-Code" (Nueva Wiki) ✅
**Objetivo:** Eliminar la wiki tradicional a favor de contenido dentro del repositorio.
- [x] **Eliminar Wiki Legacy:** Borrar el sistema de repositorios independientes `.wiki.git`.
- [x] **Implementar "Wiki-Folder":** 
  - [x] **Detector de Carpeta:** Activar pestaña Wiki si existe `/docs` o `.docs` (carpeta seleccionable en ajustes).
  - [x] **Sidebar Dinámico:** Generar menú lateral leyendo YAML/TOML Frontmatter de los `.md`.
  - [x] **Caché por Commit:** Solo recalcular el sidebar si el `CommitSHA` cambia.
  - [x] **Wikilinks:** Soporte para `[[Enlace]]` mediante Goldmark en el backend.

---

## Fase 12: Funcionalidades de Precisión ✅
- [x] **Internacionalización (I18n):** Carga dinámica de JSON locales en Frontend.
- [x] **Búsqueda Avanzada de Repositorio:**
  - [x] Desactivar indexación profunda global (Bleve/CPU Killer).
  - [x] Implementar **Repo Search (On-demand):** Búsqueda interna vía `git grep` en tiempo real.
- [x] **Web File Management (CRUD):** Creación, subida (Drag & Drop) y borrado de archivos vía web.

---

## Fase 13: Rebranding e Identidad Final ✅
- [x] **Refactor de Nomenclatura:** Actualizar `gogs.io/gogs` a la nueva identidad `ggit` en todo el core.
- [x] **Arquitectura de Inicialización Granular:** Desacoplado de README, Licencia y .gitignore en la creación de repositorios.
- [x] **Multiselección de .gitignore:** Implementado sistema de chips para combinar múltiples plantillas.
- [x] **Identicons SVG (Client-side):** Componente Svelte para generación determinista de avatares SVG basada en seed.
- [x] **Rutas por defecto:** Cambiar carpetas de almacenamiento a `ggit-repositories`.
- [x] **Component Library:** Consolidar la librería de componentes Svelte (M3 Expressive).

---

## Fase 14: Sistema de Plantillas Dinámico ✅
**Objetivo:** Automatizar y securizar la creación de repositorios eliminando dependencias de terceros.
- [x] **Script de Sincronización:** Creado `scripts/sync_templates.sh` para organizar `github/gitignore`.
- [x] **Categorización Inteligente:** Estructura de carpetas `/languages` y `/global`.
- [x] **Ggit Standard Base:** Creado `_ggit_security.gitignore` con protección para `.env`, logs y claves.
- [x] **Concatenador de Plantillas:** El backend fusiona múltiples plantillas en el commit inicial.
- [x] **Selector Multi-select:** Interfaz moderna de etiquetas para selección múltiple.

---

## Fase 15: Optimización y Limpieza de Legado ✅
**Objetivo:** Eliminar deuda técnica y reducir el peso del binario.
- [x] **Eliminación de Activos Obsoletos:** Borrados 5.6MB de emojis en PNG y archivos huérfanos.
- [x] **Emojis Nativos (Unicode):** Uso de caracteres de sistema vía Goldmark.
- [x] **Orquestación Unificada (Makefile):** Migración completa a `Makefile` con inyección de metadatos.
- [x] **Modernización de Docker:** Limpieza de lógica de SOCAT y tuneleo obsoleto.
- [x] **Markdown Puro:** Eliminada lógica de Smartypants para evitar transformaciones indeseadas.

---

## Fase 16: Ecosistema de Configuración Unificado ✅
**Objetivo:** Unificar la identidad visual bajo un estándar premium (Modular Island Architecture).
- [x] **Modular Island UI:** Arquitectura de "islas" con radios de `3rem` y glassmorphism.
- [x] **AdminAuditLog:** Registro de auditoría integrado en el panel de administración.
- [x] **Monitor & Mantenimiento:** Consolas de actividad y métricas en tiempo real.
- [x] **Magic Avatar 2.0:** Selector de avatares deterministas (DiceBear) con previsualización.
- [x] **QA de Sincronización:** Verificación de coherencia entre estados de UI y backend.

---

## Fase 17: Depuración de Funcionalidades Sociales ✅
**Objetivo:** Simplificar la plataforma eliminando funciones de "red social" innecesarias.
- [x] **Eliminación Total de Seguidores:** Borrada la interfaz y lógica de backend para Follow/Unfollow.
- [x] **Limpieza de Base de Datos:** Eliminadas tablas y contadores de seguidores en GORM.
- [x] **Optimización de Perfil:** Rediseño de la cabecera de usuario para ahorro de espacio.

---

## Fase 18: Refinado de Gestión y Conectividad ✅
- [x] **Gestión de Usuarios Bloqueados:** Vista en `UserSettings.svelte` para lista negra personal.
- [x] **Zona de Peligro (Organizaciones):** Conexión del borrado con el backend.
- [x] **Administración Glocal:** `AdminOrgEditor.svelte` para edición administrativa de organizaciones.
- [x] **Validación UX Proactiva:** Mejor feedback en alta de llaves SSH y Access Tokens.

---

## Fase 19: Desacoplo de Renderizado ✅
**Objetivo:** Delegar el renderizado de Markdown completamente al cliente.
- [x] **Frontend Rendering Core:** Implementación de `MarkdownRenderer.svelte`.
- [x] **Limpieza de Modelos DB:** Eliminación de campos `RenderedContent`.
- [x] **Poda de Configuración:** Eliminación de secciones `[markdown]` y `[smartypants]` del core.

---

## Fase 20: Estabilidad y Estandarización de UI ✅
- [x] **Sistema de Diálogos Global:** Sustitución de `alert/confirm/prompt` por componentes M3.
- [x] **Sistema de Notificaciones (Toasts):** Avisos no intrusivos globales.
- [x] **Componentes de Formulario Unificados:** `Checkbox`, `Select`, e `Input` reutilizables.
- [x] **Validación de Datos (Zod):** Esquemas de validación en el cliente.
- [x] **Testing de Regresión:** Configuración de Playwright para flujos críticos.
- [x] **Optimización CSS:** Reutilización de componentes y limpieza de estilos.

---

## Fase 21: Robustez y Ecosistema (Siguiente Nivel) 🚀
**Objetivo:** Convertir Ggit en una plataforma de desarrollo profesional completa.

### 1. Code Review & Branch Protection ✅
- [x] **Protección de Ramas:** Reglas para evitar pushes directos y requerir PRs.
- [x] **Review Workflow:** Sistema de aprobaciones (Approve/Request Changes) en Pull Requests.
- [x] **Comentarios en Línea:** Capacidad de comentar directamente en el diff del código.
- [ ] **Merge Checks:** Bloqueo de fusión si hay conversaciones pendientes o falta de aprobaciones.

### 2. CI/CD & Automatización
- [ ] **Commit Status API:** Modelo DB y API para reportar resultados de CI (iconos ✅/❌).
- [ ] **OAuth2/OIDC Provider:** Flujo para que el CI pueda autenticar usuarios vía `ggit`.
- [x] **Webhooks Seguros:** Firma de payloads mediante `HMAC-SHA256`.

### 3. Observabilidad & APM ✅
- [x] **Prometheus & Health Checks:** Exposición de `/metrics` y `/health` detallados.
- [x] **Instrumentación OpenTelemetry:** Rastreo básico para identificar cuellos de botella.

### 4. Gestión de Proyectos
- [ ] **Tableros Kanban:** Gestión de Issues por columnas (Todo, In Progress, Done).
- [ ] **Búsqueda Global de Código:** Indexación bajo demanda para buscar texto en todo el servidor.

---

## Fase 22: Integridad y Autocuración 🛡️
**Objetivo:** Garantizar la consistencia entre Base de Datos y Sistema de Archivos.
- [ ] **Consistencia Atómica:** Lógica de "Two-Phase Commit" para operaciones Git + DB.
- [ ] **Ggit Health Sentinel:** Tarea en segundo plano para ejecutar `git fsck` rotativo.
- [ ] **Self-Healing Hooks:** Reparación automática de hooks corruptos o ausentes.
- [x] **Dangling Objects Cleanup:** Identificar y limpiar rastro de archivos de repositorios borrados.
- [ ] **Blindaje Anti-OOM:** Limitador de procesos Git concurrentes basado en RAM disponible.
- [x] **Modo Supervivencia:** Switch a "Solo Lectura" si el espacio en disco es crítico.
- [x] **Dashboard Leído:** Animación de parpadeo en contadores de dashboard para items nuevos.
- [ ] **Ggit Credential Helper:** Almacenamiento seguro para tokens de remotos HTTP.

---

## Fase 23: Paridad Funcional (Gogs-to-Svelte Sync) 🔗
**Objetivo:** Recuperar funcionalidades útiles del legado en la nueva interfaz.
- [x] **Auditoría de Paridad UI:** Sincronización de opciones de Repo, User y Org.
- [ ] **Gestión de Colaboradores:** Interfaz para añadir usuarios y asignar roles.
- [ ] **Deploy Keys & Git Hooks:** Editores web para claves y scripts de servidor.
- [ ] **Label & Milestone Manager:** Editores dinámicos para metadatos de Issues.
- [ ] **Webhook Tester & Logs:** Herramientas de diagnóstico para integraciones.

---

## Fase 24: Estandarización de la Capa de Datos 🗄️
- [x] **Gestión de Sesiones (DB Store):** Almacenamiento persistente de sesiones en base de datos.
- [x] **Migraciones Atómicas:** Sistema de migraciones v26 implementado.
- [ ] **Auditoría de Índices:** Optimización de tablas críticas (UserRepository, Actions).
- [ ] **Pattern Store:** Unificar accesos a DB bajo un patrón de repositorio común.

---

## Fase 25: Infraestructura y Auth Final 🔐
- [ ] **Gestión de Fuentes de Auth:** UI para LDAP, OAuth2 (GitHub, Google) y PAM.
- [ ] **Release Management:** Upload center para binarios y estados de release (Draft/Pre).
- [ ] **Monitorización Admin:** Cron scheduler y process tracker en el panel admin.
- [ ] **Branding Global:** Instance customizer para logos, favicons y custom CSS/JS.
- [ ] **Exportación Portable:** Backup de repositorios/organizaciones en un click.
---

## Estimación de Dificultad (Tareas Pendientes)

| Tarea | Fase | Dificultad | Explicación |
| :--- | :---: | :---: | :--- |
| **Webhooks Seguros** | 21 | Completado | Implementar firmas HMAC-SHA256 en las cabeceras de salida. |
| **Dangling Objects Cleanup** | 22 | Completado | Script para eliminar registros de DB que ya no tienen repo en disco. |
| **Modo Supervivencia** | 22 | Completado | Lógica de bloqueo de escritura basada en umbral de espacio en disco. |
| **Dashboard Leído** | 22 | Completado | Implementación de animaciones CSS/Svelte para estados de lectura. |
| **Branding Global** | 25 | Baja | Interfaz simple para modificar logos y estilos personalizados. |
| **OIDC (OpenID Connect)** | 10 | Media | Integración de cliente estándar para autenticación externa. |
| **Merge Checks** | 21 | Media | Validación de hilos de conversación resueltos antes de permitir merge. |
| **Commit Status API** | 21 | Media | Creación de endpoints y modelos para reportes de CI externos. |
| **Ggit Health Sentinel** | 22 | Media | Orquestación de `git fsck` periódico sin degradar el rendimiento. |
| **Self-Healing Hooks** | 22 | Media | Detección y recreación automática de scripts de hooks faltantes. |
| **Blindaje Anti-OOM** | 22 | Media | Limitador de concurrencia dinámico basado en la RAM disponible. |
| **Ggit Credential Helper** | 22 | Media | Sistema de cifrado y almacenamiento para tokens de remotos HTTP. |
| **Gestión de Colaboradores** | 23 | Media | Nueva UI de búsqueda y gestión de roles granulares. |
| **Deploy Keys & Git Hooks** | 23 | Media | Editores web con validación de sintaxis y seguridad en servidor. |
| **Label & Milestone Manager** | 23 | Media | Refactorización de gestión de metadatos a componentes Material 3. |
| **Webhook Tester & Logs** | 23 | Media | Sistema de reintentos y visualización de payloads de entrega. |
| **Auditoría de Índices** | 24 | Media | Optimización de consultas críticas mediante índices SQLite precisos. |
| **Fuentes de Auth UI** | 25 | Media | Interfaz administrativa para configurar proveedores LDAP/OAuth2. |
| **Release Management** | 25 | Media | Lógica de subida y persistencia de activos binarios en lanzamientos. |
| **Monitorización Admin** | 25 | Media | Consola de control para tareas cron y procesos Git activos. |
| **Exportación Portable** | 25 | Media | Empaquetado atómico de base de datos y archivos Git para backups. |
| **OAuth2/OIDC Provider** | 21 | Alta | Implementar Ggit como proveedor de identidad completo (Server-side). |
| **Tableros Kanban** | 21 | Alta | Sistema interactivo de gestión de estados con drag-and-drop. |
| **Búsqueda Global** | 21 | Alta | Indexación profunda de código con alto consumo de recursos. |
| **Consistencia Atómica** | 22 | Alta | Coordinación de transacciones (2PC) entre SQLite y el sistema de archivos. |
| **Pattern Store** | 24 | Alta | Refactorización estructural de toda la capa de persistencia del sistema. |
