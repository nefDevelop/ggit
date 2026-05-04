# Roadmap & Changelog: Modernización (Go/Gin + Svelte)

## Fase 1: Análisis y Funcionalidades Base ✅
- [x] **Integrar detección avanzada de lenguajes:** Reemplazar sistema con `go-enry`. **Tested**.
- [x] **Resaltado de sintaxis en el servidor:** Integrar `chroma` y eliminar JS pesado en cliente. **Tested**.
- [x] **Auditoría de Settings:** Verificar flags `EnableIssues` y `EnablePulls`. **Tested**.
- [x] **Crear capa de Servicios:** Extraer lógica de Repos, Issues, PRs y Usuarios.
- [x] **Tipado de Entradas/Salidas:** Definir DTOs independientes del contexto HTTP.
- [x] **Inicialización y Middleware:** Migrar Auth, Sesiones, Logging y Recovery.
- [x] **Migración a API v2:** Implementar todos los endpoints JSON para el nuevo frontend.
- [x] **Bypass de Instalación:** Pre-configuración automática con SQLite. **Tested**.

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
- [x] **Mantener LFS:** Decidido mantener el soporte para Large File Storage (`internal/lfsx`). **Tested**.
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
- [x] **Web File Management (CRUD):** Creación, subida (Drag & Drop) y borrado de archivos vía web. **Tested**.

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
- [x] **AdminAuditLog:** Registro de auditoría integrado en el panel de administración. **Tested**.
- [x] **Monitor & Mantenimiento:** Consolas de actividad y métricas en tiempo real.
- [x] **Magic Avatar 2.0:** Selector de avatares deterministas (DiceBear) con previsualización. **Tested**.
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
- [x] **Protección de Ramas:** Reglas para evitar pushes directos y requerir PRs. **Tested**.
- [x] **Review Workflow:** Sistema de aprobaciones (Approve/Request Changes) en Pull Requests. **Tested**.
- [x] **Comentarios en Línea:** Capacidad de comentar directamente en el diff del código. **Tested**.
- [ ] **Merge Checks:** Bloqueo de fusión si hay conversaciones pendientes o falta de aprobaciones.

### 2. CI/CD & Automatización
- [ ] **Commit Status API:** Modelo DB y API para reportar resultados de CI (iconos ✅/❌).
- [ ] **OAuth2/OIDC Provider:** Flujo para que el CI pueda autenticar usuarios vía `ggit`.
- [x] **Webhooks Seguros:** Firma de payloads mediante `HMAC-SHA256`. **Tested**.
- [x] **Prometheus & Health Checks:** Exposición de `/metrics` y `/health` detallados. **Tested**.

### 3. Observabilidad & APM ✅
- [x] **Prometheus & Health Checks:** Exposición de `/metrics` y `/health` detallados. **Tested**.
- [x] **Instrumentación OpenTelemetry:** Rastreo básico para identificar cuellos de botella.

### 4. Gestión de Proyectos
- [x] **Tableros Kanban:** Gestión de Issues por columnas (Todo, In Progress, Done). **Integrated 🗂️**

---

## Fase 22: Integridad y Autocuración 🛡️
**Objetivo:** Garantizar la consistencia entre Base de Datos y Sistema de Archivos.
- [ ] **Consistencia Atómica:** Lógica de "Two-Phase Commit" para operaciones Git + DB.
- [x] **Ggit Health Sentinel:** Orquestación de `git fsck` automático basado en actividad. **Tested**.
- [x] **Self-Healing Hooks:** Reparación automática de hooks corruptos o ausentes. **Integrated 🛡️**
- [x] **Dangling Objects Cleanup:** Identificar y limpiar rastro de archivos de repositorios borrados. **Tested**.
- [x] **Blindaje Anti-OOM:** Limitador de procesos Git concurrentes basado en RAM disponible. **Integrated 🛡️**
- [x] **Modo Supervivencia:** Switch a "Solo Lectura" si el espacio en disco es crítico. **Tested**.
- [x] **Dashboard Leído:** Animación de parpadeo en contadores de dashboard para items nuevos.
- [ ] **Ggit Credential Helper:** Almacenamiento seguro para tokens de remotos HTTP.
- [x] **Secret Guardian 🛡️:** Detección proactiva de secretos (API keys, .env) en tiempo de push. **Tested**.

---

## Fase 23: Paridad Funcional (Gogs-to-Svelte Sync) 🔗
**Objetivo:** Recuperar funcionalidades útiles del legado en la nueva interfaz.
- [x] **Auditoría de Paridad UI:** Sincronización de opciones de Repo, User y Org.
- [x] **Gestión de Colaboradores:** Interfaz para añadir usuarios y asignar roles (Read, Write, Admin).
- [x] **Deploy Keys:** Editor de llaves de despliegue para CI/CD.
- [x] **Git Hooks Manager:** Editor de scripts del lado del servidor para automatización avanzada. **Tested**.
- [x] **Label & Milestone Manager:** Editores dinámicos para metadatos de Issues.
- [x] **Webhook Tester & Logs:** Herramientas de diagnóstico para integraciones. **Tested**.

---

## Fase 24: Estandarización de la Capa de Datos 🗄️
- [x] **Gestión de Sesiones (DB Store):** Almacenamiento persistente de sesiones en base de datos. **Tested**.
- [x] **Migraciones Atómicas:** Sistema de migraciones v26 implementado.
- [x] **Auditoría de Índices:** Optimización de tablas críticas (UserRepository, Actions). **Tested**.
- [ ] **Pattern Store:** Unificar accesos a DB bajo un patrón de repositorio común.

---

## Fase 25: Infraestructura y Auth Final 🔐
- [x] **Gestión de Fuentes de Auth:** UI para LDAP, OAuth2 (GitHub, Google) y PAM.
- [x] **Release Management:** Upload center para binarios y estados de release (Draft/Pre). **Tested**.
- [x] **Monitorización Admin:** Cron scheduler y process tracker en el panel admin.
- [x] **Respaldo de Metadatos:** Exportación ligera de DB y adjuntos (sin repositorios).
- [x] **Adopción de Repositorios (Hardened 🛡️):** Función para vincular carpetas Git existentes con validación estricta de Path Traversal. **Tested**.
- [x] **Seguridad de Búsqueda:** Mitigación de Parameter Injection en `git grep` mediante separadores de argumentos. **Tested**.

---

## Fase 26: The Lean Developer Hub (Filosofía de Bajo Impacto) 🚀
**Objetivo:** Reducir la fricción del desarrollador integrando herramientas profesionales en la UI existente.

### 1. Flujo de Trabajo en Pull Requests (Vistas de PR)
- [x] **Estados de Atención:** Semáforo de turno (Autor/Revisor) integrado en la cabecera del PR.
- [x] **Merge Checks:** Validación visual de hilos resueltos y aprobaciones en el cuadro de merge.
- [x] **Commit Status API:** Integración de iconos de CI (✅/❌) en listas de commits y PRs.
- [x] **Scoped Labels:** Soporte para etiquetas excluyentes `key::value` en el panel lateral.

### 2. Cirugía de Git en la Web (Herramientas de Código)
- [ ] **Web Cherry-Pick:** Botón directo en la vista de commit para copiar cambios a otra rama.
- [ ] **Conflict Resolver:** Editor visual integrado en la vista de PR para resolver choques de código.
- [ ] **Web-based Rebase:** Interfaz para reordenar o hacer squash de commits desde la web.
- [ ] **Web Stash Manager:** Nueva zona de "Herramientas" para guardar cambios temporales sin commit.
- [ ] **Reflog Explorer:** Sección en Mantenimiento para recuperar commits mediante el historial de HEAD.

### 3. Productividad & Navegación
- [ ] **Unified Inbox:** Nueva pestaña en el Dashboard para gestionar tareas de todos los repositorios.
- [ ] **Inline Blame:** Visualización de autor y fecha al final de cada línea en el visor de archivos.
- [ ] **Hierarchical Projects:** Soporte para organizar repositorios en carpetas en el Dashboard.

### 4. Organización & Planificación
- [x] **Kanban Virtual:** Selector de vista en la lista de Issues para pasar de modo Lista a Tablero.
- [x] **Hill Charts:** Integración de gráficas de incertidumbre en la vista de Hitos (Milestones). **Integrated 🏔️**
- [ ] **Gists Integrados:** Nueva zona global para gestionar snippets de código personales.
- [x] **Technotes / ADR:** Nueva pestaña "Docs" en el repositorio para diarios técnicos y decisiones. **Integrated 📝**

### 5. Core & Janitor (Mantenimiento)
- [x] **Stale Detection:** Nueva regla en el **Janitor** existente para detectar ramas e issues muertas. **Healed 🛡️**
- [x] **Blindaje Anti-OOM:** Monitor de recursos dinámico y limitador de concurrencia.
- [ ] **Atomic Consist. & Pattern Store:** Refactorización interna para integridad y mantenibilidad.

---

## Estimación de Dificultad (Tareas Pendientes)

| Tarea | Fase | Dificultad | Consumo RAM | Explicación |
| :--- | :---: | :---: | :---: | :--- |
| **Webhooks Seguros** | 21 | Completado | Despreciable | Implementar firmas HMAC-SHA256 en las cabeceras de salida. |
| **Dangling Objects Cleanup** | 22 | Completado | Bajo | Script para eliminar registros de DB que ya no tienen repo en disco. |
| **Modo Supervivencia** | 22 | Completado | Despreciable | Lógica de bloqueo de escritura basada en umbral de espacio en disco. |
| **Dashboard Leído** | 22 | Completado | Despreciable | Implementación de animaciones CSS/Svelte para estados de lectura. |
| **OIDC (OpenID Connect)** | 10 | Media | Bajo | Integración de cliente estándar para autenticación externa. |
| **Merge Checks** | 26 | Completado | Bajo | Integración visual de validaciones antes de permitir merge. |
| **Commit Status API** | 26 | Completado | Bajo | Creación de endpoints y modelos para reportes de CI externos. |
| **Ggit Health Sentinel** | 22 | Completado | Alto | Orquestación de `git fsck` periódico sin degradar el rendimiento. |
| **Self-Healing Hooks** | 22 | Completado | Bajo | Detección y recreación automática de scripts de hooks faltantes. |
| **Blindaje Anti-OOM** | 26 | Completado | Medio | Limitador de concurrencia dinámico basado en la RAM disponible. |
| **Ggit Credential Helper** | 22 | Media | Bajo | Sistema de cifrado y almacenamiento para tokens de remotos HTTP. |
| **Gestión de Colaboradores** | 23 | Completado | Bajo | Nueva UI de búsqueda y gestión de roles granulares. |
| **Deploy Keys & Git Hooks** | 23 | Completado | Bajo | Editores web con validación de sintaxis y seguridad en servidor. |
| **Label & Milestone Manager** | 23 | Completado | Bajo | Refactorización de gestión de metadatos a componentes Material 3. |
| **Webhook Tester & Logs** | 23 | Completado | Medio | Sistema de reintentos y visualización de payloads de entrega. |
| **Auditoría de Índices** | 24 | Completado | Medio | Optimización de consultas SQLite precisas. **Tested**. |
| **Fuentes de Auth UI** | 25 | Completado | Bajo | Interfaz administrativa para configurar proveedores LDAP/OAuth2. |
| **Release Management** | 25 | Completado | Medio | Lógica de subida y persistencia de activos binarios en lanzamientos. |
| **Monitorización Admin** | 25 | Completado | Bajo | Consola de control para tareas cron y procesos Git activos. |
| **Respaldo de Metadatos** | 25 | Completado | Bajo | Exportación ligera de DB y adjuntos (sin repositorios). |
| **Adopción de Repositorios** | 22 | Completado | Bajo | Función para vincular carpetas Git existentes a la DB. |
| **Tableros Kanban** | 26 | Completado | Medio | Selector de vista interactivo en la lista de Issues. |
| **Consistencia Atómica** | 26 | Alta | Medio | Coordinación de transacciones entre SQLite y el sistema de archivos. |
| **Pattern Store** | 26 | Alta | Medio | Refactorización estructural de la capa de persistencia. |
| **Estados de Atención** | 26 | Completado | **0 MB** | Motor de turnos para PRs integrado en la cabecera. |
| **Unified Inbox** | 26 | Media | Bajo | Nueva pestaña global para gestión de tareas multi-repo. |
| **Stale Detection** | 26 | Baja | **0 MB** | Nueva regla integrada en el motor de Janitor existente. |
| **Gists Integrados** | 26 | Media | **0 MB** | Nueva zona global para gestión de snippets de código. |

- [x] **Analisis de basura y de seguridad**: 
---

## 🚩 Pendientes (Priorizados para 1.x)


### Estabilidad y Robustez (Inmediato)
- [x] **Mitigación SSRF en Webhooks:** Filtro de IPs locales (127.0.0.1, 10.0.x, etc) en el Tester HTTP. **Tested**.
- [x] **Commit Status API:** Integración de estados de CI externos. **Tested**.
- [x] **Auditoría de Índices:** Optimización final de consultas SQL. **Tested**.
- [x] **Blindaje Anti-OOM:** Limitador de concurrencia Git dinámico. **Healed 🛡️**
- [x] **Self-Healing Hooks:** Sistema de autorreparación de scripts en disco. **Healed 🛡️**

### Conectividad y Ecosistema (Medio Plazo)
- [ ] **OIDC (OpenID Connect):** Soporte estándar para login externo.
- [ ] **Ggit Credential Helper:** Almacenamiento seguro de credenciales para remotos HTTP. 
- [ ] **Pattern Store:** Refactorización de la capa de acceso a datos (Clean Architecture). SI

### Funcionalidades Avanzadas (Largo Plazo / 2.0)
- [x] **Estados de Atención:** Implementar motor de turnos en PRs.
- [x] **Kanban Virtual:** Vista de tablero basada en etiquetas.
- [ ] **Consistencia Atómica (2PC):** Coordinación estricta entre DB y sistema de archivos.


---
