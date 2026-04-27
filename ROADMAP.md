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

- [x] **Modular Island UI:** Implementada arquitectura de "islas" con radios de `3rem` y _glassmorphism_ avanzado.
- [x] **Iconografía HERO:** Estandarización de iconos de gran formato para mejorar la jerarquía visual en todas las secciones.
- [x] **AdminAuditLog [NUEVO]:** Creación del componente de registro de auditoría integrado en el panel de administración.
- [x] **Monitor & Mantenimiento:** Refactor de vistas críticas con consolas de actividad y métricas en tiempo real.

### Refactorización de Configuración de Usuario (9 Categorías)

- [x] **Simetría Visual:** Aplicado el mismo lenguaje de diseño del Admin a todas las secciones de configuración de usuario.
- [x] **UserSettings Orchestrator:** Rediseño de la barra lateral con indicadores activos y transiciones tonales.
- [x] **Magic Avatar 2.0:** Integración optimizada del selector de avatares deterministas (DiceBear) con previsualización premium.

### Coherencia de Sistema

- [x] **Unificación de Componentes:** Uso sistemático de `Card.svelte` y `Button.svelte` en todo el ecosistema de ajustes.
- [x] \*\*QA de Sincro## Fase 17: Depuración de Funcionalidades Sociales ✅

**Objetivo:** Simplificar la plataforma eliminando funciones de "red social" innecesarias para centrarse en la gestión de código.

- [x] **Eliminación Total de Seguidores:** Borrada la interfaz y toda la lógica de backend para Follow/Unfollow.
- [x] **Limpieza de Base de Datos:** Eliminadas tablas y contadores de seguidores en los modelos de GORM.
- [x] **Optimización de Perfil:** Rediseño de la cabecera de usuario para ahorrar espacio vertical.

---

## Fase 18: Refinado de Gestión y Conectividad Final (Completado) ✅

**Objetivo:** Cerrar brechas entre el frontend y el backend, y asegurar que todas las acciones administrativas sean funcionales.

- [x] **Gestión de Usuarios Bloqueados:** [COMPLETO]
  - [x] Implementar vista en `UserSettings.svelte` para gestionar la lista negra personal.
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

## Fase 20: Estabilidad y Estandarización de UI (Completado) 💎

**Objetivo:** Eliminar la deuda técnica visual y asegurar una experiencia de usuario robusta y coherente.

- [x] **Sistema de Diálogos Global:** Sustitución de `alert/confirm/prompt` nativos por componentes M3.
- [x] **Sistema de Notificaciones (Toasts):** Implementación de avisos no intrusivos globales.
- [x] **Componentes de Formulario Unificados:** Creación de `Checkbox`, `Select` e `Input` reutilizables.
- [x] **Gestión de Errores Centralizada:** Interceptores en la capa de API para feedback automático.
- [x] **Validación de Datos (Zod):** Implementar esquemas de validación en el cliente.
- [x] **Estandarización de Skeletons:** Estados de carga coherentes en todas las tablas y paneles.
- [x] **Unificar**: optimizar los css, y componentes. reutilizar componentes etc. [COMPLETO]
- [x] **Testing de Regresión:** Configurar Playwright para flujos críticos. [CONFIGURADO]

---

## Fase 21: Robustez y Ecosistema (Siguiente Nivel) 🚀

**Objetivo:** Convertir Ggit en una plataforma de desarrollo profesional completa, más allá de un simple servidor de Git.

### 1. Code Review & Branch Protection (Calidad de Código)
- [ ] **Protección de Ramas:** Implementar reglas para evitar pushes directos a `main/master` y requerir PRs.
- [ ] **Review Workflow:** Sistema de aprobaciones (Approve/Request Changes) en Pull Requests.
- [ ] **Comentarios en Línea:** Capacidad de comentar directamente en el diff del código durante una revisión.
- [ ] **Merge Checks:** Bloqueo de fusión si hay conversaciones pendientes o falta de aprobaciones.

### 2. CI/CD & Automatización (Integración Woodpecker/Drone) <- A FUTURO
- [ ] **Infraestructura de Estados (Commit Status):**
  - [ ] **Modelo DB:** Crear tabla `commit_statuses` para persistir estados (pending, success, error, failure) vinculados a un SHA.
  - [ ] **API v2 Endpoints:**
    - `POST /api/v2/repos/:owner/:repo/statuses/:sha`: Para que el CI reporte resultados.
    - `GET /api/v2/repos/:owner/:repo/commits/:ref/statuses`: Para consultar el estado desde el UI.
  - [ ] **UI Integration:** Mostrar iconos de estado ✅/❌/⏳ en la lista de archivos, historial de commits y Pull Requests.
- [ ] **Identidad & Conectividad (OAuth2 Provider):**
  - [ ] Implementar flujo de **OAuth2 / OIDC** para que Woodpecker/Drone puedan autenticar usuarios mediante `ggit`.
  - [ ] Endpoint de descubrimiento `.well-known/openid-configuration`.
- [ ] **Webhooks Seguros:**
  - [ ] Implementar firma de payloads mediante secreto compartido (`HMAC-SHA256`) en la cabecera `X-Ggit-Signature`.
  - [ ] **Webhooks de Organización:** Permitir disparar eventos para todos los repositorios de una organización.

### 3. Observabilidad & APM (Estabilidad Real)

- [ ] **Activación de Prometheus:** Conectar el `MetricsFilter` al router de Gin y exponer `/metrics`.
- [ ] **Instrumentación OpenTelemetry:** Añadir rastreo (tracing) básico para identificar cuellos de botella en operaciones Git pesadas.
- [ ] **Health Checks:** Endpoint `/health` detallado para orquestadores como Kubernetes o Docker Compose.

### 4. Gestión de Proyectos & Colaboración
- [ ] **Tableros Kanban:** Implementar "Ggit Boards" simples para gestionar Issues por columnas (Todo, In Progress, Done).
- [ ] **Búsqueda Global de Código:** Sistema de indexación bajo demanda (tipo `zoekt` o similar) para buscar texto en todos los repositorios del servidor.
- [ ] **Dependencias de Issues:** Capacidad de marcar que un Issue bloquea a otro.

### 5. Seguridad & Cumplimiento (Confianza Total)
- [ ] **Verificación de Firmas GPG:**
  - [ ] Permitir a los usuarios subir sus llaves GPG públicas.
  - [ ] Validar firmas en commits y mostrar el sello "Verified" en el UI.
- [ ] **Escaneo de Secretos:** Detectar accidentalmente claves API o tokens en los pushes (pre-receive hooks).
- [ ] **Soporte OIDC:** Implementar OpenID Connect como estándar para el inicio de sesión único (SSO) empresarial.

### 6. Gestión de Recursos & Resiliencia
- [ ] **Cuotas de Almacenamiento:** Implementar límites de disco configurables por usuario y por repositorio.
- [ ] **LFS File Locking:** Implementar la API de bloqueo de archivos para evitar conflictos de edición en binarios.
- [ ] **Backups Automatizados (Cloud Sync):**
  - [ ] Crear sistema de rotación de backups.
  - [ ] Integrar subida automática a almacenamiento S3/R2 o compatible mediante `rclone`.
- [ ] **Auditoría Extendida:** Registrar cambios en permisos de colaboradores y visibilidad de repositorios en el Audit Log.

---

## Fase 22: Integridad y Autocuración (Pilar de Robustez) 🛡️

**Objetivo:** Garantizar que Ggit sea "incapaz de fallar" o corromper datos, incluso en condiciones adversas (cortes de luz, hardware limitado, alta carga).

### 1. Consistencia Atómica (DB + FS)
- [ ] **Coordinación de Transacciones:** Implementar lógica de "Two-Phase Commit" simplificada para asegurar que un cambio en la base de datos solo se consolide si la operación de Git en disco fue exitosa.
- [ ] **Rollback de Sistema de Archivos:** Capacidad de revertir cambios manuales en el sistema de archivos si la transacción de la base de datos falla.

### 2. Autodiagnóstico y Mantenimiento Proactivo
- [ ] **Ggit Health Sentinel:** Tarea en segundo plano que ejecuta `git fsck` en todos los repositorios de forma rotativa para detectar corrupción de objetos.
- [ ] **Optimización Automática de SQLite:** Programar `PRAGMA optimize`, `VACUUM` y `REINDEX` periódicos para mantener la base de datos ágil y sin fragmentación.
- [ ] **Verificación de Checksums:** Almacenar y verificar checksums de archivos críticos de configuración y hooks.

### 3. Blindaje de Recursos (Anti-OOM)
- [ ] **Process Throttling:** Implementar un limitador de procesos Git concurrentes basado en la RAM disponible (vital para Raspberry Pi).
- [ ] **Timeouts Estrictos:** Asegurar que ningún proceso de Git (como un diff masivo o un grep global) pueda quedar colgado indefinidamente consumiendo CPU.
- [ ] **Modo de Supervivencia (Read-only):** Switch automático a modo de "Solo Lectura" si el espacio en disco baja de un umbral crítico (ej. 50MB) para evitar corrupción de la base de datos SQLite.

---

## Fase 23: Paridad Funcional y Reconexión (Gogs-to-Svelte Sync) 🔗

**Objetivo:** Asegurar que ninguna funcionalidad útil del Gogs original se haya quedado "huérfana" o incompleta en la migración a Svelte.

### 1. Auditoría de Paridad (UI Gap Analysis)
- [x] **Mapeo de Rutas Legacy:** Revisar las antiguas plantillas `.tmpl` de Gogs para identificar opciones de configuración faltantes. [COMPLETO]
  - [x] Sincronización de opciones de Repository (Public Wiki/Issues, Manual Merge, Repo Avatar).
  - [x] Sincronización de opciones de User (AvatarEmail, DisableNewsFeed).
  - [x] Sincronización de opciones de Organization (Private Org, Org Avatar).
  - [x] Sincronización de opciones de Admin (AllowGitHook, AllowImportLocal en User Editor).
- [ ] **Funcionalidades de Repositorio "Huérfanas":**
  - [ ] **Gestión de Colaboradores:** Interfaz para añadir usuarios por nombre y asignar roles.
  - [ ] **Deploy Keys:** Interfaz para claves SSH exclusivas del repositorio.
  - [ ] **Server-side Git Hooks:** Editor de scripts para hooks de Git en el servidor.
  - [ ] **Gestión de Ramas:** Listado de ramas, borrado y selección de rama por defecto.
  - [x] **Branch Protection:** Configuración de reglas de protección para ramas críticas. [COMPLETO]
  - [ ] **Mirror Settings:** Configuración de sincronización para repositorios espejo.
  - [ ] **LFS Management:** Visualización y control de objetos Large File Storage.
- [ ] **Gestión de Organizaciones Avanzada:**
  - [ ] **Team Editor:** Interfaz detallada para añadir miembros y repositorios a un equipo.
  - [ ] **Gestión de Miembros:** Cambiar roles (Propietario vs Miembro) y expulsión masiva.
- [ ] **Motor de Migración e Importación:**
  - [ ] **Wizard de Migración:** UI para importar desde GitHub/GitLab/Gitea usando tokens y clones espejo.
  - [ ] **Importación de Metadatos:** Capacidad de importar Issues y Labels durante la migración inicial.
- [ ] **UI de Metadatos de Issues:**
  - [ ] **Label Manager:** Editor de etiquetas dinámico (colores M3).
  - [ ] **Milestone Planner:** Gestión de hitos y progreso porcentual.
- [ ] **Integraciones de Webhooks:**
  - [ ] **Payload Presets:** UI específica para webhooks de Slack, Discord y MS Teams.
  - [ ] **Webhook Tester:** Botón para disparar eventos de prueba (Test Delivery).
  - [ ] **Delivery Log:** Historial de envíos con re-intento manual (Redeliver).
- [x] **Completado de Componentes "Cortos":** Identificar componentes Svelte que están visualmente presentes pero cuya acciones aún no llaman a la API. [COMPLETO]

### 2. Sincronización API-Frontend
- [x] **Wiki Externa (Legacy Link):** Reinstaurar la capacidad de redirigir la pestaña Wiki a una URL externa. [COMPLETO]
- [ ] **Wiki Onboarding:** Mostrar una guía visual y plantillas sugeridas cuando la carpeta de la wiki (`/docs`) está vacía.
- [ ] **Soporte _sidebar.md:** Implementar renderizado automático del archivo `_sidebar.md` como menú lateral global.
- [x] **Consolidación de API v2:** Asegurar que todos los campos de los modelos de base de datos sean editables vía API. [COMPLETO]
- [x] **Reconexión de Flags de Repositorio:** Verificar que opciones de visibilidad y módulos funcionen bidireccionalmente. [COMPLETO]

### 3. Validación de Flujos Completos
- [ ] **Vistas de Código Avanzadas:**
  - [ ] **Blame View:** Implementar la vista de autoría línea por línea.
  - [ ] **Compare & Pull Request Flow:** UI para comparar ramas y previsualizar diffs antes de crear PRs.
  - [ ] **Historias de Archivo:** Navegación optimizada por el historial de un solo archivo.
- [ ] **Funcionalidades Sociales & Feed:**
  - [ ] **RSS/Atom Feeds:** Re-habilitar los endpoints de feed para actividad global y de repo.
  - [ ] **Buscador Global (Explore):** Filtros avanzados por lenguaje y popularidad.
- [ ] **Ciclo de Vida de Issues/PRs:** Asegurar que todas las acciones (cerrar, reabrir, etiquetar, asignar) estén operativas al 100% en la nueva interfaz.
- [x] **Gestión de Sesiones Activas:** Implementar vista para listar y revocar sesiones de usuario (navegadores/dispositivos). [COMPLETO]
- [x] **Danger Zone de Usuario:** Implementar el flujo de borrado permanente de cuenta. [COMPLETO]
- [x] **Configuraciones de Seguridad Avanzadas:** Reconectar opciones de borrado de cuenta, visibilidad y transferencia. [COMPLETO]

---

## Fase 24: Estandarización de la Capa de Datos (EN PROCESO) 🗄️

**Objetivo:** Consolidar la estructura de la base de datos y asegurar que el esquema sea robusto y extensible.

- [x] **Gestión de Sesiones (DB Store):** Implementar almacenamiento de sesiones en base de datos para persistencia entre reinicios. [COMPLETO]
- [x] **Migraciones Atómicas:** Implementar sistema de migraciones v26 para cambios en el esquema de sesiones. [COMPLETO]
- [ ] **Auditoría de Índices:** Revisar y optimizar los índices de las tablas críticas (UserRepository, Actions, Issues).
- [ ] **Consolidación de Repositorios (Store Pattern):** Unificar todos los accesos a la base de datos bajo el patrón `Store`.

---

## Fase 25: Paridad de Infraestructura y Auth (Final Sync) 🔐

**Objetivo:** Completar las herramientas de administración de bajo nivel y autenticación corporativa.

- [ ] **Gestión de Fuentes de Autenticación:** Interfaz en el panel de administración para añadir/editar LDAP, OAuth2 (GitHub, Google), y PAM.
- [ ] **Release Assets & Management:**
  - [ ] **Upload Center:** Interfaz para subir binarios y archivos a cada Release.
  - [ ] **Release States:** Soporte para estados "Borrador" y "Pre-release".
- [ ] **Monitorización Administrativa:**
  - [ ] **Cron Scheduler:** Vista para monitorizar y ejecutar manualmente tareas programadas.
  - [ ] **Process Tracker:** Listado de procesos Git activos en el servidor para diagnóstico.
  - [ ] **System Notices:** Centro de avisos globales del sistema para comunicación con usuarios.
- [ ] **Mantenimiento de Repositorios:** Implementar los botones de "Regenerar hooks" y "Resincronizar llaves" en la administración.
- [ ] **Branding & Personalización:**
  - [ ] **Instance Customizer:** UI para subir logotipos, favicon y editar textos del footer.
  - [ ] **Global Injection:** Editor para CSS y JavaScript personalizado del lado del cliente.
- [ ] **Red y Seguridad de Bajo Nivel:**
  - [ ] **SSH Embedded Server:** Configuración de interfaz, puerto y llaves del servidor SSH.
  - [ ] **Reverse Proxy Config:** Gestión de Trusted Proxies y detección de IPs reales.
- [ ] **Ajustes Globales Dinámicos:** Permitir que los ajustes de `app.ini` editables vía UI se apliquen en caliente o notifiquen reinicio.
- [ ] **Exportación de Datos:** Botón para generar un backup portable de un repositorio o una organización completa.

---

## Fase 21: Robustez y Ecosistema (Siguiente Nivel) 🚀

**Objetivo:** Convertir Ggit en una plataforma de desarrollo profesional completa.

### 1. Code Review & Branch Protection
- [x] **Protección de Ramas:** Implementar reglas para evitar pushes directos a `main/master` y requerir PRs. [COMPLETO]
- [x] **Review Workflow:** Sistema de aprobaciones (Approve/Request Changes) en Pull Requests. [COMPLETO]
- [x] **Comentarios en Línea:** Capacidad de comentar directamente en el diff del código. [COMPLETO]

### 2. CI/CD & Automatización (Integración Woodpecker/Drone)
- [ ] **Commit Status API:** Modelo DB y API para reportar resultados de CI (iconos ✅/❌ en el UI).
- [ ] **OAuth2 Provider:** Implementar flujo OIDC para que el CI pueda autenticar usuarios vía `ggit`.
- [ ] **Webhooks Seguros:** Firma de payloads mediante HMAC-SHA256.

---

## Fase 22: Integridad y Autocuración 🛡️

**Objetivo:** Garantizar que los datos del servidor y el disco estén siempre sincronizados.

- [ ] **Background FS Check:** Tarea programada que valida que cada registro en la DB tenga su correspondiente carpeta `.git` en disco.
- [ ] **Self-Healing Hooks:** Reparación automática de hooks corruptos o ausentes al detectar un acceso al repositorio.
- [ ] **Dangling Objects Cleanup:** Identificar y limpiar archivos de base de datos vinculados a repositorios borrados por errores de sistema.
