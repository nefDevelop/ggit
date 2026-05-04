# Blueprint Técnico Refinado: ggit (Mapeo por Zonas Web UI)

Este documento organiza las funcionalidades según dónde aparecerán en la interfaz web de ggit y si se integran en secciones ya existentes o crean zonas nuevas.

---

## 🏗️ 1. Flujo de Trabajo en Pull Requests (Vistas de PR)
*Mejoras integradas directamente en el ciclo de revisión de código.*

| Funcionalidad | Zona Web UI | Integración | Descripción |
| :--- | :--- | :--- | :--- |
| **Estados de Atención** | Cabecera y Lista de PRs | **Existente:** Vista de PR | Semáforo de turno (Tu turno / Esperando). |
| **Merge Checks** | Cuadro de Merge (Footer) | **Existente:** Vista de PR | Indica qué falta para poder fusionar. |
| **Commit Status API** | Lista de Commits y PRs | **Existente:** Vistas de historial | Iconos de CI (✅/❌) junto a los commits. |
| **Scoped Labels** | Panel de Etiquetas | **Existente:** Issues/PRs | Etiquetas tipo `key::value` automáticas. |

---

## ✂️ 2. Cirugía de Git en la Web (Herramientas de Código)
*Funciones para manipular el repositorio desde el navegador.*

| Funcionalidad | Zona Web UI | Integración | Descripción |
| :--- | :--- | :--- | :--- |
| **Web Cherry-Pick** | Vista de Commit individual | **Existente:** Botón en Commit | Copiar un commit a otra rama con un clic. |
| **Conflict Resolver** | Vista de PR (con conflicto) | **Existente:** Vista de PR | Editor visual para arreglar choques de código. |
| **Web-based Rebase** | Vista de PR / Ramas | **Existente:** Opciones de Merge | Reordenar commits antes de fusionar. |
| **Web Stash Manager** | Visor de Archivos / Sidebar | **Nueva Zona:** "Herramientas" | Guardar cambios web temporales sin commit. |
| **Reflog Explorer** | Configuración -> Mantenimiento | **Existente:** Maintenance | Recuperar commits "perdidos" tras errores. |

---

## ⚡ 3. Productividad & Navegación (Dashboard y Visores)
*Optimización de la experiencia de uso diaria.*

| Funcionalidad | Zona Web UI | Integración | Descripción |
| :--- | :--- | :--- | :--- |
| **Unified Inbox** | Dashboard Principal | **Nueva Pestaña:** Inbox | Todas tus tareas de todos los repos en un sitio. |
| **Inline Blame** | Visor de Archivos (Código) | **Existente:** File View | Info de autor tenue al final de cada línea. |
| **Hierarchical Proj.** | Dashboard / Explorar | **Existente:** Repo Lists | Listado de repos organizado por carpetas. |

---

## 📂 4. Organización & Planificación (Gestión de Proyectos)
*Estructura y seguimiento del progreso.*

| Funcionalidad | Zona Web UI | Integración | Descripción |
| :--- | :--- | :--- | :--- |
| **Kanban Virtual** | Lista de Issues | **Existente:** Toggle Vista | Cambiar de modo Lista a modo Tablero. |
| **Hill Charts** | Vista de Hitos (Milestones) | **Existente:** Milestones | Gráfica de incertidumbre vs. ejecución. |
| **Gists Integrados** | Menú Global / Perfil | **Nueva Zona:** Gists | Snippets de código rápidos y personales. |
| **Technotes / ADR** | Nueva pestaña en Repo | **Nueva Pestaña:** Docs | Diario técnico y decisiones de arquitectura. |

---

## 🛡️ 5. Core, Seguridad & Janitor (Administración)
*Mantenimiento y robustez del sistema.*

| Funcionalidad | Zona Web UI | Integración | Descripción |
| :--- | :--- | :--- | :--- |
| **Stale Detection** | **Janitor** (Repo & Admin) | **Existente:** Janitor | Detecta ramas e issues muertas. |
| **Blindaje Anti-OOM** | Panel Admin -> Monitor | **Existente:** Monitor | Control de RAM para procesos Git pesados. |
| **Atomic Consist.** | Invisible (Backend) | -- | Garantía de integridad Git + DB. |
| **Pattern Store** | Invisible (Refactor) | -- | Mejora de la arquitectura de código. |

---

## 🔍 Resumen de Integración
1.  **Aprovechamiento:** El 70% de las funciones se inyectan en vistas que ya existen (PRs, Commits, Janitor), reduciendo el desorden visual.
2.  **Nuevas Zonas:** Solo se crean pestañas nuevas para conceptos potentes como los **Gists**, el **Inbox Global** o las **Technotes**.
3.  **Janitor:** Se convierte en el "cerebro de mantenimiento", avisando de qué limpiar (Stale) además de hacerlo (GC).
