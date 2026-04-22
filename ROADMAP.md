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
- [x] - github.com/unknwon/com: Es una librería de utilidades antigua de la era Gogs. Gran parte de lo que hace ya lo hace la librería estándar de Go o strx.

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

---

He analizado las dependencias actuales del proyecto leyendo tu go.mod y explorando el código. Aquí tienes un análisis técnico profundo
sobre las librerías actuales y las mejores alternativas para implementar la funcionalidad "Wiki-Folder", buscando el máximo rendimiento
y mantenibilidad:

1. Motor de Markdown y Extensiones (Wikilinks)

- Lo que usas actualmente: github.com/russross/blackfriday (v1.6.0).
- Análisis: Esta versión es muy antigua (Blackfriday v1 está esencialmente deprecado). Es rápido, pero su diseño es monolítico.
  Extenderlo para soportar [[Wikilinks]] o interceptar nodos te obligaría a hacer "hacks" procesando el string de texto a base de
  Expresiones Regulares antes o después del renderizado, lo cual es lento y muy propenso a errores.
- La mejor alternativa: github.com/yuin/goldmark
  - Por qué: Es el estándar de oro actual en el ecosistema Go (lo usa Hugo y Gitea). Es 100% compatible con CommonMark y, lo más
    importante, está basado en AST (Abstract Syntax Tree).
  - Para implementar [[Wikilinks]], solo necesitas escribir un pequeño Parser para el AST. Esto hace el parsing en una sola pasada,
    de forma nativa, limpia y extremadamente eficiente.

2. Extracción del Frontmatter (YAML/TOML para el Sidebar)

- Lo que tendrías que hacer ahora: Usar regexp para cortar todo lo que está entre --- y --- al inicio del .md, y luego pasarle ese
  string a gopkg.in/yaml.v3 (que he visto que tienes como dependencia indirecta en el go.mod).
- La mejor alternativa: github.com/yuin/goldmark-meta (una extensión oficial de Goldmark).
  - Por qué: Extrae el Frontmatter automáticamente durante la fase de parseo del AST, sin penalizar el rendimiento. Cuando llamas a
    la función de renderizado, directamente te devuelve un map[string]any con los datos (ej: Title, Order) y el contenido limpio
    listo para HTML. Ahorras ciclos de CPU y simplificas el código enormemente.

3. Lectura del Repositorio Git (Detector de /docs y archivos)

- Lo que usas actualmente: github.com/gogs/git-module (v1.8.7).
- Análisis: Esta librería funciona haciendo llamadas al binario del sistema (git rev-list, git cat-file, etc.) usando os/exec. Para
  leer un archivo suelto está bien, pero para montar el "Sidebar Dinámico", tendrías que leer el contenido de 10, 20 o 50 archivos .md
  de golpe. Hacer spawn de 50 procesos git en el servidor por cada visita hundiría el rendimiento.
- La mejor alternativa: github.com/go-git/go-git/v5.
  - Por qué: Es una implementación de Git puramente en Go. Lee los archivos directamente del disco (o memoria) parseando los objetos
    de Git (blobs, trees) sin levantar procesos externos.
  - Con go-git, recorrer la carpeta docs/ y extraer el contenido de 50 archivos .md toma literalmente microsegundos. Sería
    infinitamente más escalable para Ggit.

4. Caché del Sidebar por Commit

- Análisis: Construir la jerarquía del menú desde los .md es costoso. Solo debería hacerse una vez por cada Push.
- La mejor alternativa: github.com/hashicorp/golang-lru/v2
  - Por qué: Es una caché en memoria LRU (Least Recently Used) genérica (usando Generics de Go 1.18+) ultrarrápida y segura para
    concurrencia.
  - La lógica sería: clave = RepoID_CommitSHA. Si un usuario entra a la Wiki, buscas esa clave en la LRU. Si el SHA del último
    commit de la rama principal no ha cambiado, devuelves el Sidebar precalculado en $O(1)$. Si cambia (alguien hace push), la clave
    antigua desaparecerá naturalmente de la caché.

Resumen de la Estrategia Óptima

Si yo tuviera que diseñar esta arquitectura para Ggit hoy, la propondría así:

1.  Sustituir Blackfriday por Goldmark en todo el proyecto de Ggit. Esto modernizaría enormemente todo el renderizado de repositorios,
    no solo la Wiki.
2.  Usar go-git (junto al git-module que ya tienes para compatibilidad) solo para la función de mapear el directorio de la Wiki a gran
    velocidad en memoria.
3.  Almacenar la estructura del Sidebar generada en un LRU Cache atado al Commit SHA.

