---
title: "Auditoría de consistencia visual"
description: "Estado de la interfaz: espaciado, shells de página, componentes y responsive"
icon: "layout-grid"
---

# Auditoría de consistencia visual

Rama `feature/ui-consistency`. Revisión estática de todo el frontend (62 páginas,
~205 componentes) con evidencia `archivo:línea`. El objetivo no es rediseñar: es que
página, tarjeta y sección respiren igual en toda la aplicación.

## Resumen ejecutivo

Cinco causas explican casi todas las inconsistencias visibles:

1. **No existe una escala de espaciado.** Cero tokens `--space-*`. Hay 13 pasos de `rem`
   crudos y medidas en `px` mezcladas (`styles/layout.css:78-99`). Las utilidades
   `.gap-xs…xl` existen (`styles/utils.css:72-76`) y tienen **0 usos**.
2. **No existe un shell de página.** `#main-content` tiene `padding: 0`
   (`styles/layout.css:60-67`), así que cada página inventa su padding. `.page-container`
   está definido **11 veces** con 4 paddings y 2 `max-width` distintos.
3. **Hay tres sistemas de tarjeta/superficie.** `.island` (`styles/components.css:28`),
   `ui/Card` (`Card.svelte:45-74`) y `.settings-premium-card`
   (`styles/components.css:311`), con sombras y paddings distintos para el mismo rol.
4. **17 breakpoints** para 4 puntos reales (173 reglas `@media`): 600/640/650/768 y
   800/850/900/950/1000 conviven. La fundación solo conoce 1300/1024/600/480.
5. **Padding duplicado** en admin (shell + 12 componentes), organización y perfil, y
   márgenes negativos en commits, así que el contenido "salta" al navegar entre secciones.

## 1. Fundación (tokens, tipografía, superficies)

### 1.1 Espaciado

- No hay tokens de espacio. Valores crudos detectados en la fundación: `0.1, 0.2, 0.25,
  0.5, 0.75, 1, 1.25, 1.5, 1.75, 2, 2.5, 3, 6rem` y `80/90/145/160px`.
- `.gap-*` y `.flex-col/.flex-row` no se usan en ningún componente.
- `Card` tiene su propia escala (`1 / 1.5 / 2.5rem`) que no coincide con nada global.

### 1.2 Tipografía

- 759 declaraciones `font-size` en 142 archivos. Sin tokens de tamaño, peso ni
  `line-height`; `body` no define `line-height` (`styles/base.css:27-35`).
- Tres pesos para el mismo botón primario: 600 (`ui/Button.svelte:117`), 700
  (`components.css:113`), 800 (`components.css:143`).
- Etiquetas de formulario: **7 recetas** (0.75→0.95rem, 600→800, mayúsculas o no).

### 1.3 Superficies y conflictos

- Seis alias para cinco roles MD3 (`--bg-surface`, `--island-bg`, `--surface`,
  `--md3-surface-container`, `--md3-surface-2`…), elegidos a dedo por componente.
- `.tonal` está definido **tres veces** (`tokens.css:166` con `!important`,
  `components.css:44`, `Card.svelte:62`); el `!important` gana y anula los otros dos.
- `--bg-danger-subtle` se usa (`components.css:179,189`) y no está declarado.
- `tables.css:12` conserva `rgba(0,0,0,0.15)` fijo: las tablas se ven oscuras en tema claro.
- `h1` global aplica degradado con `-webkit-text-fill-color: transparent`
  (`base.css:44-53`) y `h2::after` inyecta un divisor (`base.css:65-71`) en cualquier
  tarjeta, aunque el componente no lo pida.
- `html, body { overflow: hidden; 100vw/100vh }` (`base.css:32-33`) obliga a que todo
  scroll ocurra en `#main-content`, y provoca 162 escapes `:global(...)`.

## 2. Shell de página y espaciado

`#main-content` no aporta padding. Convenciones detectadas:

| Convención | Páginas | Ejemplo |
| :--- | :--- | :--- |
| `0` | landing, auth, user settings, labels, milestones, releases, pulls/new, forks/stars/watchers | `routes/user/settings/+page.svelte:17` |
| `4rem 2rem` | gists (3), explore/gists | `routes/gists/new/+page.svelte:60` |
| `0 1.5rem 4rem` | gists, explore (4) | `ExploreView.svelte:74` |
| `2rem 1.5rem 4rem` | dashboard, home | `DashboardView.svelte:297` |
| `1.5rem 1.5rem 4rem` | 12 componentes admin | `AdminUsers.svelte:146` |
| `2rem` | org (3), user settings/profile | `routes/org/[orgname]/+page.svelte:20` |
| `3rem 1rem` | releases/new | `routes/.../releases/new/+page.svelte:25` |
| `0 2rem` | perfil de usuario | `routes/[username]/+page.svelte:232` |
| margen negativo `-2rem` | commits | `routes/.../commits/+page.svelte:134` |

Problemas concretos:

- **Doble padding admin**: `AdminView.svelte:86` + los 12 `Admin*.svelte` (`:146`, `:155`,
  `:208`…) aplican `1.5rem 1.5rem 4rem` dos veces (hoy lo salva `.edge-to-edge`).
- **Org y perfil**: página `2rem` + tarjeta `2.5rem` (`OrgSettings.svelte:226`).
- **Repo settings**: layout `1rem` + página `2rem` + `PageHeader` `2.5rem`.
- **Commits**: `margin: -2rem` + `padding: 2rem` para recolocar el contenido.
- **admin/repos/adoption**: fuera del shell con padding, a `0` (`:60`).
- Landing y auth usan `.container`, **clase inexistente** en `src/styles`.

## 3. Componentes

- **`.empty-state` ×25** con cinco paddings (`4rem`, `4rem 2rem`, `5rem`, `5rem 2rem`,
  `6rem 2rem`) y gaps 1/1.25/1.5rem. `SettingsTokens.svelte:238` redefine
  `.empty-state-box` con otros valores.
- **`.skeleton-*` ×~20** definiciones (12 `.skeleton-list` con estructuras distintas).
- **`.section-title`** con cuatro apariencias distintas.
- **`.form-grid`** con gaps 1.5/2/2.5rem.
- **`Select` desalineado con `Input`**: label 0.85/800 mayúsculas vs 0.9/600; fondo
  distinto; `id="select-input"` fijo (IDs duplicados).
- **Selects e inputs a mano**: `RepoSettingsGeneral.svelte:225-234` implementa un
  `<select>` propio; `SettingsProfile.svelte:380` sobreescribe el input.
- **Botones destructivos con 5 variantes** (`danger`, `ghost+.danger-btn`,
  `primary+.danger-btn`, `outline+.danger-outline`, `ghost+.danger-text`).
- **`PageHeader` con 4 usos "stub"** (`RepoPulls:51`, `RepoIssues:58`, `RepoLabels:105`,
  `tags:122`) y ~8 cabeceras hechas a mano.
- **`@keyframes fadeIn/pulse/spin` re-declarados ~30 veces**.
- **`!important` para corregir tarjetas**: `RepoDocs:226`, `SettingsProfile:478,493`,
  `AdminJanitor:152`, `repo/create:555,877`.

## 4. Responsive

- 173 reglas `@media`, 17 umbrales. El móvil real es 600 (79 usos), pero conviven
  640 y 650; el tablet se reparte entre 800/850/900/950/1000.
- Gutter lateral distinto por sección en móvil: repo `0.5rem`
  (`[username]/[reponame]/+layout.svelte:471`), dashboard `0.75rem`
  (`DashboardView.svelte:303`), admin **0** (`admin/+layout.svelte:74`), ajustes
  `1rem`/`1.5rem`.
- **Tablas**: `.modern-table-card { overflow: hidden }` (`tables.css:11`) recorta; el
  wrapper `.table-responsive-wrapper` no lo usa nadie. `AdminNotices` y `AdminAuths` no
  tienen `data-label`, así que el modo apilado muestra etiquetas vacías.
- **Sin `@media`**: `WebhookLogViewer` (grid 350px + 1fr, alto 700px),
  `AdminJanitor` (`minmax(400px,1fr)`), `OrgTeams` (`minmax(350px,…)`),
  `OrgMembers` (`minmax(320px,…)`), `SettingsBlockedUsers`, `PullDetail`, `Modal`.
- **Tap targets < 44px**: botones icono del header del repo (`padding:4px` ≈26px),
  `Button size="sm"` ≈31px, toolbar de `FileContentViewer` 34px, cerrar modal ≈36px.
- Sub-navs con tres estilos: admin y user settings difuminados a `z-index:500`;
  repo settings opaco a `999` (`[reponame]/settings/+page.svelte:206-217`).

## 5. Decisiones canónicas propuestas

Esto es la especificación a la que debería converger todo. Nada de esto cambia el
lenguaje visual que ya te gusta; solo fija valores.

### 5.1 Escala de espaciado

```css
--space-1: 0.25rem;  /* 4px  */
--space-2: 0.5rem;   /* 8px  */
--space-3: 0.75rem;  /* 12px */
--space-4: 1rem;     /* 16px */
--space-5: 1.5rem;   /* 24px */
--space-6: 2rem;     /* 32px */
--space-7: 3rem;     /* 48px */
--space-8: 4rem;     /* 64px */
```

- Dentro de una tarjeta/formulario: `--space-3/4`.
- Entre secciones de una página: `--space-6`.
- Padding de página: `--space-5` lateral, `--space-7` inferior.
- Escala de tarjeta: `none 0 / sm --space-4 / md --space-5 / lg --space-7`.

### 5.2 Shell de página

Una sola clase global `.page` (definida en `styles/layout.css`) que aporte:

```css
.page {
  display: flex;
  flex-direction: column;
  gap: var(--space-6);
  padding: var(--space-5) var(--space-5) var(--space-7);
  max-width: 1400px;
  margin: 0 auto;
  width: 100%;
}
.page--narrow { max-width: 900px; }
.page--flush  { padding: 0; }
```

Reglas: las páginas no definen padding propio; los componentes que hoy son "página"
(`Admin*`, `DashboardView`, `ExploreView`, `SocialListView`) pasan a `padding: 0` y
el shell pone el aire. Los layouts de repo/admin dejan de duplicar.

### 5.3 Breakpoints

Cuatro valores y nada más:

| Nombre | Valor | Uso |
| :--- | :--- | :--- |
| `sm` | `600px` | móvil |
| `md` | `900px` | tablet |
| `lg` | `1024px` | shell de escritorio |
| `xl` | `1300px` | tres columnas |

CSS no permite variables en `@media`; se documentan como constantes y se sustituyen
los 13 valores restantes. Las reglas 480 se pliegan en 600; 640/650→600; 800/850/950/
1000→900; 1100/1200/1400→1300.

### 5.4 Un solo sistema de tarjeta

- `ui/Card` es la tarjeta canónica. `.island` se mantiene como alias de compatibilidad
  con el mismo padding/sombra, o se migra mecánicamente.
- Prohibido `:global(.card) { padding: … !important }` en componentes.
- Superficies: usar solo `--md-sys-color-surface-container*`; retirar los alias
  duplicados.

### 5.5 Componentes compartidos

- **`EmptyState.svelte`**: `padding: 4rem 2rem; gap: 1.25rem;` icono 3rem, título 1.75rem.
- **`Skeleton`**: variantes `list`, `card`, `form` con altura y gap fijos.
- **`FormField`**: label `0.9rem/600`, gap `0.5rem`, control `0.75rem 1rem`, hint `0.8rem`.
  `Input` y `Select` lo usan; se eliminan los `<select>` a mano.
- **Botón destructivo**: una sola forma, `variant="danger"`; `danger-outline` para filas.
- **`PageHeader`**: siempre con `title`/`icon`; si solo hacen falta acciones, se extrae
  un snippet `PageActions`.

### 5.6 Ritmo

- Gap entre secciones dentro de una página: `2rem`.
- Gap dentro de formularios: `1.5rem`.
- Stats grid: `minmax(280px, 1fr)` + `1.5rem`.
- Un solo `@keyframes fadeIn/pulse/spin` en `styles/`.

## 6. Plan por fases

| Fase | Alcance | Riesgo | Criterio de aceptación |
| :--- | :--- | :--- | :--- |
| 0 | Tokens de espacio, `.page`, breakpoints documentados, `EmptyState`/`Skeleton` compartidos (aditivo) | Bajo | Build y tipos OK; ninguna página cambia aún |
| 1 | Admin (12 páginas): quitar doble padding, migrar a `.page`, unificar vacíos/skeletons | Bajo | Todas las admin con el mismo gutter y bottom |
| 2 | Repo: layout `1rem`→`.page`, issues/pulls/labels/milestones/releases/adr/tags, tablas con scroll y `data-label` | Medio | Sin recortes de tabla a 375px; gutter igual en todas las subpáginas |
| 3 | Ajustes, auth, gists, org, explore, perfil: shell y padding único; formularios a `FormField` | Medio | Dos páginas de ajustes consecutivas se sienten iguales |
| 4 | Responsive: plegar breakpoints, modales, tap targets 44px, sub-navs unificados | Medio | Sin overflow a 320px; iconos ≥44px |
| 5 | Limpieza: `.tonal` único, tokens duplicados, `@keyframes`, `!important`, clases muertas | Bajo | `grep` sin `!important` en tarjetas |

Quick wins (una sola sesión): `.page` en admin, `EmptyState`, tablas con
`overflow-x`, `Select` alineado con `Input`, botón destructivo único.

## 7. Cómo medir que se cumple

- `grep -r "padding: 1.5rem 1.5rem 4rem" lib/components` → 0.
- `grep -rn ":global(.card)" src` → 0.
- Páginas con `padding:` propio (excluyendo `.page--flush`) → 0.
- `@media` con valores fuera de {600, 900, 1024, 1300} → 0 (salvo casos justificados).
- Toda `.modern-table` dentro de `.table-responsive-wrapper`.
- Iconos de acción con área ≥44×44px.

## Estado de implementación (rama `feature/ui-consistency`)

- **Fase 0 hecha**: escala `--space-1..8` en `tokens.css`, clase `.page` (con
  `--narrow/--wide/--flush` y versión móvil) en `layout.css`, `EmptyState.svelte` y
  `Skeleton.svelte` (variantes `list/form/card`, compatible con el uso previo).
- **Piloto admin hecho**: el layout de admin es el único dueño del padding
  (`1.5rem 1.5rem 4rem`, móvil `1rem 1rem 5rem`) y los 12 `Admin*` pasan a `padding: 0`.
  `AdminJanitor` usa ya `EmptyState` y `Skeleton`.
- **Fase 2 (parcial)**: gutters del repo unificados con admin (`--space-5`/`--space-7`,
  móvil `--space-4`/`--space-6`), fuera el `margin: -2rem` de commits, padding inferior
  duplicado de ajustes eliminado, `WikiView` sin doble padding, tablas con scroll
  horizontal y fondo del tema, y `data-label` restaurado en `AdminNotices` y `AdminAuths`.
- **Fase 2b**: `PageActions.svelte` sustituye las cuatro cabeceras "stub" (pulls, issues,
  labels, tags) y `EmptyState` reemplaza los estados vacíos de pulls, issues y labels,
  borrando su CSS local. El CTA de issues deja de ser un `<a>` manual.
- **Fase 2c**: `EmptyState` reemplaza también los estados vacíos de documentación,
  milestones, releases, forks/stars/watchers y el de error del navegador de archivos;
  se elimina el CSS local de `.empty-state` y `.empty-icon-ring` donde ya no se usa.
  `RepoFileBrowser` conserva su bloque especial "Push your code".
- **Fase 4 (parcial)**: grids de Janitor, equipos y miembros con una columna en móvil;
  modales con padding móvil y cierre de 44px; panel de notificaciones y popover de ramas
  acotados al viewport; el visor de logs de webhooks pasa a una columna por debajo de 900px.
- **Fase 3 (parcial)**: dashboard, explore, gists (listado, nuevo, vista, edición y
  explore/gists) y organización usan la escala `--space`; las páginas de organización
  adoptan `.page.page--medium` y eliminan su `.page-container` local. Reglas móviles
  añadidas a explore, gists y org.
- **Fase 3 (cierre)**: ajustes de usuario con `--space` y hueco inferior para el sub-nav
  móvil; las cuatro páginas de auth sin regla móvil (`sign_up`, `activate`,
  `reset_password`, `forgot_password`) reducen el padding de la tarjeta por debajo de 600px.
- **Detectado**: `AdminView.svelte` no lo importa ninguna ruta (código muerto); conserva
  el padding duplicado y se puede eliminar en una fase de limpieza.

Pendiente: fases 2-5 del plan.

## Anexo: de dónde salen los números

- 62 `+page.svelte`; 173 reglas `@media`; 17 umbrales distintos.
- 759 `font-size`; 162 usos de `:global(...)`; ~25 `.empty-state`.
- 11 definiciones locales de `.page-container`; 3 sistemas de tarjeta.
- `.gap-*`, `.flex-col`, `.text-truncate`, `.w-full`: 0 usos.
