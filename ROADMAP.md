# Roadmap: Modernización de Gogs a Stack Moderno (Go/Gin + Svelte)

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
*   [ ] Port remaining route packages (`admin`, `user`, etc.)
*   [ ] Finalize API v2 specification
*   [ ] **Diseño de la API v2 (API First):**
    *   Definir la estructura de paquetes en `internal/route/api/v2`.
    *   Comenzar a implementar los nuevos endpoints bajo el prefijo `/api/v2/`.
    *   Estos endpoints consumirán la nueva capa de Servicios (Fase 2) y devolverán exclusivamente JSON.
*   [x] **Transición del SSR:**
    *   Adaptar temporalmente las rutas web para que `Gin` pueda renderizar los viejos `templates/` HTML, permitiendo que la aplicación siga funcionando mientras se construye el nuevo frontend. [COMPLETO]

## Fase 4: Nuevo Frontend con SvelteKit (Largo Plazo)
**Objetivo:** Reemplazar el antiguo SSR basado en plantillas Go por una Single Page Application (SPA) o aplicación SvelteKit ultra-rápida y moderna.

*   [ ] **Inicializar proyecto SvelteKit:**
    *   [x] Initialize SvelteKit project in `frontend/`
    *   [ ] Configure Vite proxy for Gin backend
    *   [ ] Implement base design system with Vanilla CSS
*   [ ] **Configuración de Proxy y API:**
    *   [x] Implement `internal/ginx/csrf.go` (Gin-compatible CSRF)
    *   [x] Refactor `internal/context/context.go` to wrap `*gin.Context`
    *   [x] Stabilize `internal/route/repo` package (Resolved `stdctx` vs `gcontext` conflicts)
    *   [ ] Update `cmd/web.go` to use Gin router
*   [ ] **Implementar Autenticación y Estado Global:**
    *   Conectar el frontend con los endpoints de login/registro de Gin (manejo de JWT o cookies seguras).
*   [ ] **Migración Progresiva de Pantallas:**
    *   Dashboard / Explorar repositorios.
    *   Vista de Repositorio (Árbol de archivos, Historial de Commits).
    *   Vista de Código (Integrando el HTML generado por `chroma` en la Fase 1).
    *   Issues y Pull Requests.
    *   Panel de Configuración de Usuario y Repositorio.
*   [ ] **Eliminación del Legacy:**
    *   Una vez que el frontend SvelteKit cubra el 100% de las funcionalidades, eliminar el directorio `templates/` y la lógica de renderizado HTML del backend Gin.
    *   Construir el frontend en estático y servirlo desde Gin (o usar un servidor web/CDN delante).
