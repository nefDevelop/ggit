# Registro de Cambios Estructurales (Ggit vs Gogs)

## ✨ Añadido / Modernizado

- **Backend Web:** Migración completa de Macaron a **Gin** (API v2, Router rápido, Middleware moderno y seguro).
- **Frontend:** Nueva SPA con **SvelteKit** (Proxy Vite, Vanilla CSS, diseño MD3 Expressive).
- **Arquitectura:** Implementación de **Service Layer** para desacoplar la lógica de negocio (BD, Git) de los controladores HTTP.
- **Detección de Lenguajes:** Integración de `go-enry` (el mismo motor que usa GitHub) para estadísticas precisas.
- **Resaltado de Sintaxis:** Renderizado Server-side de código con `chroma`, eliminando la carga de librerías JS pesadas en el cliente.
- **Base de Datos:** Rendimiento optimizado utilizando exclusivamente **SQLite en modo WAL** (alta concurrencia nativa, 0 configuración).
- **Instalación Directa:** Pre-configuración automática (bypass del antiguo instalador de base de datos).
- **UI/UX:** Micro-animaciones, soporte nativo de temas avanzados (Deep Night, Solarized, Ayu, Monokai, Nord).
- **Documentación:** Especificación y soporte integrado para **OpenAPI / Swagger** en la API v2.
- **Monitoreo:** Panel de estado en tiempo real (CPU, Memoria, Goroutines) integrado nativamente en Svelte.

## 🗑️ Eliminado / Deprecado

- **Framework Macaron:** Eliminado junto a sus dependencias acopladas (`macaron/inject`, `macaron/binding`).
- **Plantillas SSR de Go:** Eliminación completa del motor de vistas clásico HTML (`.tmpl`).
- **Gestores de BD Pesados:** Eliminado el soporte y dependencias de MySQL/MariaDB para simplificar el binario.
- **Avatares Externos:** Eliminados **Gravatar** y Libravatar para mejorar la privacidad y depender exclusivamente de almacenamiento local/identicons.
- _(Planificado)_ **Wiki Fantasma:** Eliminación de los repositorios ocultos (`.wiki.git`) a favor de un enfoque _Docs-as-Code_ directamente en la rama principal.
- **Mirrors Automáticos:** Eliminación de clonaciones programadas de terceros para evitar consumo excesivo de CPU y red en segundo plano.
- **Autenticación Exótica:** Eliminación de los métodos de inicio de sesión complejos (PAM, SMTP-Auth, GitHub OAuth) en favor de Local + LDAP. [COMPLETO]
- _(Planificado)_ **OIDC (OpenID Connect):** Implementación de autenticación moderna.
- **Migraciones de APIs Externas:** Limpieza de código de sincronización masiva con GitHub/GitLab.
- **Exceso de Idiomas -> Carga de idiomas lazy desde svelte evitando cargarlos todos en RAM:** Limpieza de archivos de traducciones (`locale`) no mantenidos o de bajo uso para ahorrar uso de memoria RAM en el arranque.


- Captcha, sustituido por sistema mas simple
