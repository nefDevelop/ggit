# Roadmap: Modernización de Ggit a Stack Moderno (Go/Gin + Svelte)

Este documento describe el plan de acción paso a paso para la refactorización arquitectónica y modernización del repositorio, pasando de una arquitectura tradicional renderizada en el servidor (SSR con Macaron y Templates HTML) a una arquitectura moderna (API REST con Gin y Frontend reactivo con Svelte), además de incorporar nuevas funcionalidades de análisis de código.

## Fase 1: Análisis de Código y Funcionalidades Base (Corto Plazo)
**Objetivo:** Familiarizarse con el motor de Git interno (`internal/gitx`) y añadir valor inmediato sin romper la arquitectura actual.

*   [x] **Integrar detección avanzada de lenguajes:**
    *   Reemplazar/Mejorar el sistema actual con `go-enry` (el estándar de GitHub) para la detección precisa del lenguaje de programación de cada archivo.
    *   Añadir estadísticas de lenguajes al repositorio (la típica barra de porcentajes de lenguajes: 60% Go, 40% JS).
*   [x] **Resaltado de sintaxis en el servidor (Server-side Syntax Highlighting):**
    *   Integrar `chroma` para generar el HTML con los colores de sintaxis directamente desde el backend.
    *   Eliminar dependencias pesadas de JavaScript en el cliente para el resaltado de código, mejorando drásticamente el rendimiento (Time to Interactive).
*   [x] **Auditoría de Settings:**
    *   Verificar y documentar el uso de los flags existentes en base de datos (`EnableIssues`, `EnablePulls`).
    *   Asegurar que la desactivación oculte correctamente las pestañas de UI y bloquee los endpoints correspondientes (seguridad de API).

## Fase 2: Desacoplamiento de la Lógica de Negocio (Medio Plazo)
**Objetivo:** Preparar el terreno para el cambio de framework web. Actualmente, la lógica de negocio (base de datos, Git) está fuertemente acoplada a los controladores de vista de `Macaron` (`internal/route`, `cmd/web.go`).

*   [x] **Crear capa de Servicios (Service Layer):**
    *   Extraer la lógica de manipulación de repositorios, usuarios, issues y pull requests desde los *handlers* hacia paquetes independientes y testeables.
        *   [x] Repositorios
        *   [x] Issues y Pull Requests
        *   [x] Usuarios y Autenticación
*   [x] **Tipado de Entradas/Salidas:**
    *   Definir estructuras claras (DTOs - Data Transfer Objects) para las peticiones y respuestas, independientes del contexto HTTP (`*context.Context`).

## Fase 3: Migración del Motor HTTP a Gin (Medio Plazo)
**Objetivo:** Reemplazar `Macaron` (obsoleto, lento, basado en reflection) por `Gin` (rápido, tipado seguro, estándar de la industria).

*   [x] **Inicialización y Middleware:**
    *   Reescribir `cmd/web.go` para inicializar el motor de `Gin`. [COMPLETO]
    *   Migrar los middlewares esenciales (Autenticación, Sesiones, Logging, Recovery). [COMPLETO]
*   [x] Integrate Gin framework core
*   [x] Port `internal/context` to Gin
*   [x] Implement new CSRF middleware with Gin
*   [x] Resolve `repo` package compilation and naming conflicts
*   [x] Port remaining route packages (`admin`, `user`, etc.)
*   [ ] Finalize API v2 specification
*   [ ] **Migración a API v2 (Endpoints JSON):**
    *   Definir la estructura de paquetes en `internal/route/api/v2`.
    *   **✅ Completados:**
        *   [x] `POST /api/v2/user/login` (Login local)
        *   [x] `POST /api/v2/user/register` (Registro de usuarios)
        *   [x] `GET  /api/v2/user/check` (Validación de sesión actual)
        *   [x] `POST /api/v2/user/logout` (Cierre de sesión)
        *   [x] `POST /api/v2/user/login/two_factor` (Autenticación 2FA)
        *   [x] `POST /api/v2/user/forgot_password` (Solicitar reset)
        *   [x] `POST /api/v2/user/reset_password` (Ejecutar reset)
        *   [x] `POST /api/v2/user/activate` (Activar cuenta/email)
        *   [x] **Ajustes de Usuario:** Ajustar perfil, cambio de contraseña, claves SSH, tokens de acceso (migrar desde `/user/settings/*`).
        *   [x] **Dashboard y Explorar:** Feeds de actividad, listado de repositorios seguidos y globales (migrar desde `/` y `/explore/*`).
        *   [x] **Repositorios:** Ver archivos, listado de commits, branches (migrar desde `/:user/:repo/*`).
        *   [x] **Colaboración:** Creación y listado de Issues, Pull Requests, Wikis y Releases.
        *   [x] **Configuración de Repositorio:** Webhooks, Colaboradores, Ramas protegidas (migrar desde `/:user/:repo/settings/*`).
        *   [x] **Panel de Admin:** Monitorización, Gestión global de usuarios y configuración (migrar desde `/admin/*`).
*   [x] **Transición del SSR:**
    *   Adaptar temporalmente las rutas web para que `Gin` pueda renderizar los viejos `templates/` HTML, permitiendo que la aplicación siga funcionando mientras se construye el nuevo frontend. [COMPLETO]
*   [x] **Bypass de Instalación:**
    *   Implementar pre-configuración automática (SQLite + Install Lock) para saltar el instalador legado y arrancar directamente en entorno moderno. [COMPLETO]

## Fase 4: Nuevo Frontend con SvelteKit (Largo Plazo)
**Objetivo:** Reemplazar el antiguo SSR basado en plantillas Go por una Single Page Application (SPA) o aplicación SvelteKit ultra-rápida y moderna.

*   [ ] **Inicializar proyecto SvelteKit:**
    *   [x] Initialize SvelteKit project in `frontend/`
    *   [x] Configure Vite proxy for Gin backend
    *   [x] Implement base design system with Vanilla CSS
*   [ ] **Configuración de Proxy y API:**
    *   [x] Implement `internal/ginx/csrf.go` (Gin-compatible CSRF)
    *   [x] Refactor `internal/context/context.go` to wrap `*gin.Context`
    *   [x] Stabilize `internal/route/repo` package (Resolved `stdctx` vs `gcontext` conflicts)
    *   [x] Update `cmd/web.go` to use Gin router
*   [x] **Implementar Autenticación y Estado Global:**
    *   Conectar el frontend con los endpoints de login/registro de Gin (manejo de JWT o cookies seguras).
*   [ ] **Migración Progresiva de Pantallas:**
    *   [x] **Login/Registro:** Flujos de autenticación fluidos. [COMPLETO]
    *   [x] **Dashboard:** Vista de islas con repositorios y feeds de actividad. [COMPLETO]
    *   [x] **Technical Audit Log:** Refactorizar el Feed de actividad hacia eventos técnicos (Pushes/Merges) con opción de desactivarlo en Settings. [COMPLETO]
    *   [x] **Vista de Repositorio:** Navegador de archivos, visor de código (Chroma) e historial de commits. [COMPLETO]
    *   [x] **Colaboración:** Interfaz de Issues, Pull Requests y Wiki. [COMPLETO]
    *   [x] **Ajustes:** Panel de Configuración de Usuario y Repositorio. [COMPLETO]
    *   [x] **Explorar:** Descubrimiento global de proyectos. [COMPLETO]
    *   [x] **Administración:** Panel de Control Global para administradores. [COMPLETO]
    *   [x] **Organizaciones:** Gestión de equipos, miembros y permisos a nivel organizacional. [COMPLETO]
    *   [x] **Avanzado de Repo:** Releases, etiquetas, visor de ramas y editor web. [COMPLETO]

## Fase 9: Sistema de Diseño y Pulido (Paralelo)
**Objetivo:** Garantizar que Ggit se sienta como una aplicación premium de 2026.
    *   [ ] Micro-animaciones en botones y transiciones de página.
    *   [ ] Refinado de sombras
    *   [x] Soporte avanzado para temas (Deep Night / Solarized, Ayu, Monokai, Nord).

## Fase 10: Limpieza de Legacy y Consistencia UX
**Objetivo:** Eliminar el 100% de los rastro de Gogs antiguo y unificar la experiencia.

*   [x] **Páginas de Error:** Modernizar 404 y 500 con Svelte. [COMPLETO]
*   [x] **Migración:** Interfaz de importación de repositorios (Phase 10.1). [COMPLETO]
*   [x] **Wiki Avanzado:** Edición e historial modernizados. [COMPLETO]
*   [x] **PR Diff Viewer:** Visor de cambios interactivo y de alto rendimiento. [COMPLETO]
*   [x] **Fase 10.3: Capa Social y Organizacional:**
    *   [x] **Perfiles Pro:** Pestañas de Stars, Followers, Following y actividad extendida.
    *   [x] **Watchers y Forks:** Listados de interactuación social en repositorios.
    *   [x] **Org Intelligence:** Gestión refinada de equipos y miembros.
*   [ ] **Fase 10.4: Herramientas de Precisión:**
    *   [ ] **PR Creator:** Selector de ramas y preview de comparación (master...branch).
    *   [ ] **Release Suite:** Gestión de etiquetas, assets y notas de versión.
    *   [ ] **Webhook Logs:** Historial de entregas y depuración de hooks.

## Fase 11: El Fin del Legacy (0 Templates)
**Objetivo:** Desconexión total del motor de plantillas de Go.

*   [ ] **Audit de Plantillas:** Eliminar archivos `.tmpl` redundantes (Meta: 0 legacy templates).
*   [ ] **Modernización de Emails:** Restilo de los 7 templates de mail a MD3 o sistema moderno.
* [x] **Admin Monitor:** Live status del servidor (Go routines, memoria, CPU) en Svelte. [COMPLETO]
*   [x] **Finalizar API v2:** Documentación OpenAPI/Swagger completa.

*   [ ] **Component Library (M3 Expressive):**
    *   Crear una librería interna de Svelte para Botones, Inputs, Modals y Toasts con el nuevo sistema de colores y elevaciones.
*   [ ] **Micro-interacciones:**
    *   Añadir transiciones suaves entre rutas, esqueletos de carga (skeleton loaders) y feedback háptico visual.

