# Roadmap: Modernización (Go/Gin + Svelte)

## Fase 1: Análisis de Código y Funcionalidades Base (Corto Plazo)

**Objetivo:** Familiarizarse con el motor de Git interno (`internal/gitx`) y añadir valor inmediato sin romper la arquitectura actual.

- [x] **Integrar detección avanzada de lenguajes:**
  - Reemplazar/Mejorar el sistema actual con `go-enry` (el estándar de GitHub) para la detección precisa del lenguaje de programación de cada archivo.
  - Añadir estadísticas de lenguajes al repositorio (la típica barra de porcentajes de lenguajes: 60% Go, 40% JS).
- [x] **Resaltado de sintaxis en el servidor (Server-side Syntax Highlighting):**
  - Integrar `chroma` para generar el HTML con los colores de sintaxis directamente desde el backend.
  - Eliminar dependencias pesadas de JavaScript en el cliente para el resaltado de código, mejorando drásticamente el rendimiento (Time to Interactive).
- [x] **Auditoría de Settings:**
  - Verificar y documentar el uso de los flags existentes en base de datos (`EnableIssues`, `EnablePulls`).
  - Asegurar que la desactivación oculte correctamente las pestañas de UI y bloquee los endpoints correspondientes (seguridad de API).

## Fase 2: Desacoplamiento de la Lógica de Negocio (Medio Plazo)

**Objetivo:** Preparar el terreno para el cambio de framework web. Actualmente, la lógica de negocio (base de datos, Git) está fuertemente acoplada a los controladores de vista de `Macaron` (`internal/route`, `cmd/web.go`).

- [x] **Crear capa de Servicios (Service Layer):**
  - Extraer la lógica de manipulación de repositorios, usuarios, issues y pull requests desde los _handlers_ hacia paquetes independientes y testeables.
    - [x] Repositorios
    - [x] Issues y Pull Requests
    - [x] Usuarios y Autenticación
- [x] **Tipado de Entradas/Salidas:**
  - Definir estructuras claras (DTOs - Data Transfer Objects) para las peticiones y respuestas, independientes del contexto HTTP (`*context.Context`).

## Fase 3: Migración del Motor HTTP a Gin (Medio Plazo)

**Objetivo:** Reemplazar `Macaron` (obsoleto, lento, basado en reflection) por `Gin` (rápido, tipado seguro, estándar de la industria).

- [x] **Inicialización y Middleware:**
  - Reescribir `cmd/web.go` para inicializar el motor de `Gin`. [COMPLETO]
  - Migrar los middlewares esenciales (Autenticación, Sesiones, Logging, Recovery). [COMPLETO]
- [x] Integrate Gin framework core
- [x] Port `internal/context` to Gin
- [x] Implement new CSRF middleware with Gin
- [x] Resolve `repo` package compilation and naming conflicts
- [x] Port remaining route packages (`admin`, `user`, etc.)
- [x] Finalize API v2 specification
- [x] **Migración a API v2 (Endpoints JSON):**
  - Definir la estructura de paquetes en `internal/route/api/v2`.
  - **✅ Completados:**
    - [x] `POST /api/v2/user/login` (Login local)
    - [x] `POST /api/v2/user/register` (Registro de usuarios)
    - [x] `GET  /api/v2/user/check` (Validación de sesión actual)
    - [x] `POST /api/v2/user/logout` (Cierre de sesión)
    - [x] `POST /api/v2/user/login/two_factor` (Autenticación 2FA)
    - [x] `POST /api/v2/user/forgot_password` (Solicitar reset)
    - [x] `POST /api/v2/user/reset_password` (Ejecutar reset)
    - [x] `POST /api/v2/user/activate` (Activar cuenta/email)
    - [x] **Ajustes de Usuario:** Ajustar perfil, cambio de contraseña, claves SSH, tokens de acceso (migrar desde `/user/settings/*`).
    - [x] **Dashboard y Explorar:** Feeds de actividad, listado de repositorios seguidos y globales (migrar desde `/` y `/explore/*`).
    - [x] **Repositorios:** Ver archivos, listado de commits, branches (migrar desde `/:user/:repo/*`).
    - [x] **Colaboración:** Creación y listado de Issues, Pull Requests, Wikis y Releases.
    - [x] **Configuración de Repositorio:** Webhooks, Colaboradores, Ramas protegidas (migrar desde `/:user/:repo/settings/*`).
    - [x] **Panel de Admin:** Monitorización, Gestión global de usuarios y configuración (migrar desde `/admin/*`).
- [x] **Transición del SSR:**
  - Adaptar temporalmente las rutas web para que `Gin` pueda renderizar los viejos `templates/` HTML, permitiendo que la aplicación siga funcionando mientras se construye el nuevo frontend. [COMPLETO]
- [x] **Bypass de Instalación:**
  - Implementar pre-configuración automática (SQLite + Install Lock) para saltar el instalador legado y arrancar directamente en entorno moderno. [COMPLETO]

## Fase 4: Nuevo Frontend con SvelteKit (Largo Plazo)

**Objetivo:** Reemplazar el antiguo SSR basado en plantillas Go por una Single Page Application (SPA) o aplicación SvelteKit ultra-rápida y moderna.

- [x] **Inicializar proyecto SvelteKit:**
  - [x] Initialize SvelteKit project in `frontend/`
  - [x] Configure Vite proxy for Gin backend
  - [x] Implement base design system with Vanilla CSS
- [x] **Configuración de Proxy y API:**
  - [x] Implement `internal/ginx/csrf.go` (Gin-compatible CSRF)
  - [x] Refactor `internal/context/context.go` to wrap `*gin.Context`
  - [x] Stabilize `internal/route/repo` package (Resolved `stdctx` vs `gcontext` conflicts)
  - [x] Update `cmd/web.go` to use Gin router
- [x] **Implementar Autenticación y Estado Global:**
  - Conectar el frontend con los endpoints de login/registro de Gin (manejo de JWT o cookies seguras).
- [x] **Migración Progresiva de Pantallas:**
  - [x] **Login/Registro:** Flujos de autenticación fluidos. [COMPLETO]
  - [x] **Dashboard:** Vista de islas con repositorios y feeds de actividad. [COMPLETO]
  - [x] **Technical Audit Log:** Refactorizar el Feed de actividad hacia eventos técnicos (Pushes/Merges) con opción de desactivarlo en Settings. [COMPLETO]
  - [x] **Vista de Repositorio:** Navegador de archivos, visor de código (Chroma) e historial de commits. [COMPLETO]
  - [x] **Colaboración:** Interfaz de Issues, Pull Requests y Wiki. [COMPLETO]
  - [x] **Ajustes:** Panel de Configuración de Usuario y Repositorio. [COMPLETO]
  - [x] **Discovery & Exploration Suite** (Explore Repos/Users/Orgs/Topics). [COMPLETO]
  - [x] **Users and Organizations Suite** (Profiles, Dashboards, Settings, Teams, Members). [COMPLETO]
  - [x] **Administración:** Panel de Control Global para administradores. [COMPLETO]
  - [x] **Avanzado de Repo:** Releases, etiquetas, visor de ramas y editor web. [COMPLETO]

## Fase 9: Sistema de Diseño y Pulido (Paralelo)

**Objetivo:** Garantizar que Ggit se sienta como una aplicación premium de 2026.
_ [x] Micro-animaciones en botones y transiciones de página.
_ [x] Refinado de sombras \* [x] Soporte avanzado para temas (Deep Night / Solarized, Ayu, Monokai, Nord).

## Fase 10: Limpieza de Legacy y Consistencia UX

**Objetivo:** Eliminar el 100% de los rastro de Ggit antiguo y unificar la experiencia.

- [x] **Páginas de Error:** Modernizar 404 y 500 con Svelte. [COMPLETO]
- [x] **Migración:** Interfaz de importación de repositorios (Phase 10.1). [COMPLETO]
- [x] **Wiki Avanzado:** Edición e historial modernizados. [COMPLETO]
- [x] **PR Diff Viewer:** Visor de cambios interactivo y de alto rendimiento. [COMPLETO]
- [x] **Fase 10.3: Capa Social y Organizacional:**
  - [x] **Perfiles Pro:** Pestañas de Stars, Followers, Following y actividad extendida.
  - [x] **Watchers y Forks:** Listados de interactuación social en repositorios.
  - [x] **Org Intelligence:** Gestión refinada de equipos y miembros.
- [x] **Fase 10.4: Herramientas de Precisión:**
  - [x] **PR Creator**: Selector de ramas y preview de comparación (master...branch).
  - [x] **Release Suite:** Gestión de etiquetas, assets y notas de versión.
  - [x] **Webhook Logs:** Historial de entregas y depuración de hooks.

- [x] **Fase 11: El Fin del Legacy (0 Templates)**
      **Objetivo:** Desconexión total del motor de plantillas de Go.

- [x] **Audit de Plantillas:** Eliminar archivos `.tmpl` redundantes (Meta: 0 legacy templates [COMPLETO] en Admin). [COMPLETO]
- [x] **Modernización de Emails:** Restilo de los 7 templates de mail a MD3. [COMPLETO]
- [x] **Admin Monitor:** Live status del servidor (Go routines, memoria, CPU) en Svelte. [COMPLETO]
- [x] **Finalizar API v2:** Documentación OpenAPI/Swagger completa. [COMPLETO]

- [ ] **Component Library (M3 Expressive):**
  - Crear una librería interna de Svelte para Botones, Inputs, Modals y Toasts con el nuevo sistema de colores y elevaciones.
- [ ] **Micro-interacciones:**
  - Añadir transiciones suaves entre rutas, esqueletos de carga (skeleton loaders) y feedback háptico visual.

## Fase 12: Brechas de Migración y Refinado Final

**Objetivo:** Alcanzar la paridad total con Ggit eliminando el código legacy y completando las herramientas de precisión.

- [x] **Sistema de Internacionalización (I18n):**
  - Implementar motor de traducciones en Svelte (soporte para JSON locales).
  - Migrar textos hardcoded a claves de traducción.
- [x] **Gestión Completa de Mirrors:**
  - Añadir controles de intervalo de sincronización en Settings.
  - Botón de sincronización forzada ("Sync Now") en la UI.
- [x] **Herramientas de Admin Avanzadas:**
  - Panel de mantenimiento (Limpieza de archivos huérfanos, REINDEX).
  - Gestión de avisos del sistema (System Notices).
- [x] **Depuración de API v1:**
  - [x] Auditar puntos finales de `v1` para asegurar paridad en `v2`. [COMPLETO]
  - [x] Desactivar/Eliminar `internal/route/api/v1` una vez confirmada la migración. [COMPLETO]
- [x] **Deploy Keys de Repositorio:**
  - Añadir gestión de claves SSH específicas para repositorios (Solo lectura/Lectura-Escritura).
- [ ] **Web File Management (CRUD):**
  - Implementar creación de nuevos archivos directamente desde la interfaz.
  - Soporte para subida de archivos (Drag & Drop).
  - Opción de eliminar archivos y carpetas.

## Fase 13: Inteligencia y Visualización (Próximo)

**Objetivo:** Transformar Ggit de una herramienta de gestión a una plataforma de insights profundos.

- [ ] **Repository Insights:**
  - Generación de gráficos de actividad (Commit Frequency) en SVG/Svelte.
  - Listado dinámico de contribuidores con estadísticas de impacto.
- [ ] **Advanced Discovery:**
  - Integración de filtros de tópicos y lenguajes en la búsqueda global (`GlobalSearch`).
  - Sugerencias inteligentes de búsqueda basadas en el contexto.

## Fase 14: Desconexión Total y Rebranding

**Objetivo:** Eliminar cualquier rastro técnico de Ggit y consolidar la identidad de Ggit.

- [ ] **Unificación de Persistencia (XORM -> GORM):**
  - Migrar el 100% de los modelos restantes a GORM para eliminar la dependencia de XORM.
- [ ] **Limpieza de Dependencias:**
  - Eliminar rastros de Macaron (`inject`, `binding`) del `go.mod` mediante refactorización de servicios.
- [ ] **Refactor de Nomenclatura:**
  - Actualizar `gogs.io/gogs` a `ggit.io/ggit` (o similar) en todo el core.
  - Cambiar rutas por defecto (`gogs-repositories` -> `ggit-repositories`) en la configuración.

## Fase 15: Optimización de Rendimiento y Simplificación (Core)

**Objetivo:** Convertir el fork en una plataforma ultra-ligera eliminando "grasa" técnica y procesos innecesarios, orientándose a equipos de alto rendimiento.

- [x] **Optimización de Base de Datos:**
  - [x] Activar modo **WAL** en SQLite para alta concurrencia. [COMPLETO]
  - [x] Eliminar soporte para **MySQL/MariaDB** para simplificar el binario y las dependencias. [COMPLETO]
- [ ] **Estrategia "Docs-as-Code" (Nueva Wiki):**
  - [ ] Eliminar el sistema de repositorios "fantasma" (`.wiki.git`) que complica respaldos y duplica el I/O en disco.
  - [ ] Implementar la "Carpeta Mágica": La Wiki leerá automáticamente de `/docs` o `.docs` del repo principal.
  - [ ] Añadir soporte para **Wikilinks** `[[Enlace]]` en el renderizado de Markdown.
- [ ] **Limpieza de Procesos en Segundo Plano:**
  - [ ] Eliminar **Migraciones Masivas** (`internal/repox`): Quitar lógica pesada de llamadas a APIs de GitHub/GitLab y tareas cron de clonación externa.
  - [x] Eliminar **Mirrors**
  - [ ] Eliminar **LFS** (`internal/lfsx`): Simplificar la gestión de archivos y tablas de base de datos.
  - [x] Reducir **Idiomas**: Dejar solo los principales (Ej. Inglés, Español, Chino) en `conf/locale` para aligerar la memoria RAM. -> termino adaptandose a svcelte los idiomas y evitando la carga en RAM desde el inicio.
  - [ ] Desactivar **Búsqueda Global de Contenido**: Evitar indexación profunda de repositorios que castigue la CPU.
- [ ] **Simplificación de Autenticación:**
  - [ ] Mantener solo **Local + LDAP/OIDC**. Eliminar métodos obsoletos y complejos (PAM, SMTP-Auth).
