<script lang="ts">
  import { onMount } from 'svelte';
  import { generateMd3Palette } from '$lib/utils/color';

  let { show = $bindable(false) } = $props();

  const themes = [
    { id: 'default', name: 'Cobalt', color: '#0b57d0' },
    { id: 'expressive', name: 'Violet', color: '#6750a4' },
    { id: 'nord', name: 'Arctic', color: '#5e81ac' },
    { id: 'ayu', name: 'Sunset', color: '#ff9940' },
    { id: 'monokai', name: 'Retro', color: '#a6e22e' }
  ];

  let currentTheme = $state('default');
  let currentScheme = $state('light');
  let customColor = $state('#0b57d0');

  // Derived palette for the custom preview dots
  let derivedPalette = $derived(generateMd3Palette(customColor, currentScheme === 'dark' || (currentScheme === 'auto' && typeof window !== 'undefined' && window.matchMedia('(prefers-color-scheme: dark)').matches)));

  onMount(() => {
    currentTheme = localStorage.getItem('theme') || 'default';
    currentScheme = localStorage.getItem('color-scheme') || (window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light');
    customColor = localStorage.getItem('custom-theme-color') || '#0b57d0';
    applyStyles();
  });

  function setTheme(id: string) {
    currentTheme = id;
    localStorage.setItem('theme', id);
    document.documentElement.setAttribute('data-theme', id);
    applyStyles();
  }

  function setScheme(scheme: 'light' | 'dark') {
    currentScheme = scheme;
    localStorage.setItem('color-scheme', scheme);
    document.documentElement.setAttribute('data-color-scheme', scheme);
    applyStyles();
  }

  function applyStyles() {
    if (currentTheme === 'custom') {
      const isDark = currentScheme === 'dark' || (currentScheme === 'auto' && window.matchMedia('(prefers-color-scheme: dark)').matches);
      const palette = generateMd3Palette(customColor, isDark);
      
      const props: Record<string, string> = {
        '--primary-color': palette.primary,
        '--md-sys-color-primary': palette.primary,
        '--md-sys-color-on-primary': palette.onPrimary,
        '--md-sys-color-primary-container': palette.primaryContainer,
        '--md-sys-color-on-primary-container': palette.onPrimaryContainer,
        '--md-sys-color-secondary': palette.secondary,
        '--md-sys-color-on-secondary': palette.onSecondary,
        '--md-sys-color-secondary-container': palette.secondaryContainer,
        '--md-sys-color-on-secondary-container': palette.onSecondaryContainer,
        '--md-sys-color-tertiary': palette.tertiary,
        '--md-sys-color-on-tertiary': palette.onTertiary,
        '--md-sys-color-tertiary-container': palette.tertiaryContainer,
        '--md-sys-color-on-tertiary-container': palette.onTertiaryContainer,
        '--md-sys-color-error': palette.error,
        '--md-sys-color-on-error': palette.onError,
        '--md-sys-color-error-container': palette.errorContainer,
        '--md-sys-color-on-error-container': palette.onErrorContainer,
        '--md-sys-color-background': palette.background,
        '--md-sys-color-on-background': palette.onBackground,
        '--md-sys-color-surface': palette.surface,
        '--md-sys-color-on-surface': palette.onSurface,
        '--md-sys-color-surface-variant': palette.surfaceVariant,
        '--md-sys-color-on-surface-variant': palette.onSurfaceVariant,
        '--md-sys-color-outline': palette.outline,
        '--md-sys-color-outline-variant': palette.outlineVariant,
        '--md-sys-color-surface-container-lowest': palette.surfaceContainerLowest,
        '--md-sys-color-surface-container-low': palette.surfaceContainerLow,
        '--md-sys-color-surface-container': palette.surfaceContainer,
        '--md-sys-color-surface-container-high': palette.surfaceContainerHigh,
        '--md-sys-color-surface-container-highest': palette.surfaceContainerHighest,
        '--color-success': palette.tertiary,
        '--color-warning': '#f59e0b',
        '--success': palette.tertiary,
        '--success-container': palette.tertiaryContainer,
      };

      for (const [prop, value] of Object.entries(props)) {
        document.documentElement.style.setProperty(prop, value);
      }
      
      const r = parseInt(palette.primary.slice(1, 3), 16);
      const g = parseInt(palette.primary.slice(3, 5), 16);
      const b = parseInt(palette.primary.slice(5, 7), 16);
      document.documentElement.style.setProperty('--primary-rgb', `${r}, ${g}, ${b}`);

      const rS = parseInt(palette.surface.slice(1, 3), 16);
      const gS = parseInt(palette.surface.slice(3, 5), 16);
      const bS = parseInt(palette.surface.slice(5, 7), 16);
      document.documentElement.style.setProperty('--surface-rgb', `${rS}, ${gS}, ${bS}`);

      const hoverMix = isDark ? 'white 5%' : 'black 5%';
      document.documentElement.style.setProperty('--surface-variant-hover', `color-mix(in srgb, ${palette.surfaceVariant}, ${hoverMix})`);
      document.documentElement.style.setProperty('--bg-success-subtle', `color-mix(in srgb, ${palette.tertiary}, transparent 85%)`);
      document.documentElement.style.setProperty('--bg-warning-subtle', `color-mix(in srgb, #f59e0b, transparent 85%)`);
    } else {
      const propsToRemove = [
        '--primary-color', '--primary-rgb', '--surface-rgb',
        '--md-sys-color-primary', '--md-sys-color-on-primary', '--md-sys-color-primary-container', '--md-sys-color-on-primary-container',
        '--md-sys-color-secondary', '--md-sys-color-on-secondary', '--md-sys-color-secondary-container', '--md-sys-color-on-secondary-container',
        '--md-sys-color-tertiary', '--md-sys-color-on-tertiary', '--md-sys-color-tertiary-container', '--md-sys-color-on-tertiary-container',
        '--md-sys-color-error', '--md-sys-color-on-error', '--md-sys-color-error-container', '--md-sys-color-on-error-container',
        '--md-sys-color-background', '--md-sys-color-on-background', '--md-sys-color-surface', '--md-sys-color-on-surface',
        '--md-sys-color-surface-variant', '--md-sys-color-on-surface-variant', '--md-sys-color-outline', '--md-sys-color-outline-variant',
        '--md-sys-color-surface-container-lowest', '--md-sys-color-surface-container-low', '--md-sys-color-surface-container',
        '--md-sys-color-surface-container-high', '--md-sys-color-surface-container-highest',
        '--color-success', '--color-warning', '--success', '--success-container',
        '--surface-variant-hover', '--bg-success-subtle', '--bg-warning-subtle'
      ];
      for (const prop of propsToRemove) {
        document.documentElement.style.removeProperty(prop);
      }
    }
  }

  function handleCustomInput(e: any) {
    customColor = e.target.value;
    localStorage.setItem('custom-theme-color', customColor);
    if (currentTheme === 'custom') applyStyles();
  }
</script>

{#if show}
<div class="theme-backdrop" onclick={() => show = false} onkeydown={(e) => e.key === 'Escape' && (show = false)} role="presentation"></div>
<div class="theme-panel island">
  <div class="panel-header">
    <div class="title-box">
      <span class="material-symbols-rounded">palette</span>
      <h3>Apariencia</h3>
    </div>
    <button class="close-btn" onclick={() => show = false}>
      <span class="material-symbols-rounded">close</span>
    </button>
  </div>

  <div class="config-stack">
    <!-- MODO DE COLOR -->
    <div class="section">
      <span class="section-tag">Esquema</span>
      <div class="scheme-toggle">
        <button 
          class:active={currentScheme === 'light'} 
          onclick={() => setScheme('light')}
          title="Modo Claro"
        >
          <span class="material-symbols-rounded">light_mode</span>
        </button>
        <button 
          class:active={currentScheme === 'dark'} 
          onclick={() => setScheme('dark')}
          title="Modo Oscuro"
        >
          <span class="material-symbols-rounded">dark_mode</span>
        </button>
      </div>
    </div>

    <!-- TEMAS -->
    <div class="section">
      <span class="section-tag">Tema de Acento</span>
      <div class="theme-grid">
        {#each themes as theme}
          <button 
            class="theme-tile tonal" 
            class:active={currentTheme === theme.id}
            onclick={() => setTheme(theme.id)}
            title={theme.name}
          >
            <div class="color-dot shadow-sm" style="background-color: {theme.color}"></div>
            <span class="theme-label">{theme.name}</span>
          </button>
        {/each}

        <!-- CUSTOM OPTION -->
        <button 
          class="theme-tile tonal custom-tile" 
          class:active={currentTheme === 'custom'}
          onclick={() => setTheme('custom')}
        >
          <div class="custom-color-box shadow-sm" style="background-color: {customColor}">
            <input type="color" value={customColor} oninput={handleCustomInput} title="Color personalizado" />
            <div class="derived-indicators">
              <div class="dot" style="background-color: {derivedPalette.secondary}"></div>
              <div class="dot" style="background-color: {derivedPalette.tertiary}"></div>
            </div>
          </div>
          <span class="theme-label">Custom</span>
        </button>
      </div>
    </div>
  </div>

  <div class="panel-footer">
    <p>Los cambios se aplican al instante</p>
  </div>
</div>
{/if}

<style>
  .theme-backdrop {
    position: fixed;
    inset: 0;
    z-index: 999;
    background: rgba(0,0,0,0.2);
    backdrop-filter: blur(4px);
  }

  .theme-panel {
    position: fixed;
    bottom: 90px;
    left: 24px;
    width: 300px;
    padding: 1.5rem;
    z-index: 1000;
    box-shadow: var(--shadow-xl);
    border: 1px solid var(--md-sys-color-outline-variant);
    background: var(--md-sys-color-surface-container);
    border-radius: 2rem;
    animation: slideUp 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  }

  @keyframes slideUp {
    from { transform: translateY(20px); opacity: 0; }
    to { transform: translateY(0); opacity: 1; }
  }

  .panel-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 2rem;
  }

  .title-box {
    display: flex;
    align-items: center;
    gap: 0.75rem;
    color: var(--primary-color);
  }

  .title-box span { font-size: 1.5rem; }
  .title-box h3 {
    margin: 0;
    font-family: var(--font-heading);
    font-size: 1.25rem;
    font-weight: 900;
    color: var(--text-base);
  }

  .close-btn {
    width: 32px;
    height: 32px;
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;
    background: var(--surface-variant);
    border: none;
    color: var(--text-muted);
    cursor: pointer;
    transition: all 0.2s;
  }

  .close-btn:hover { background: var(--surface-variant-hover); color: var(--text-base); }

  .config-stack {
    display: flex;
    flex-direction: column;
    gap: 1.75rem;
  }

  .section-tag {
    display: block;
    font-size: 0.7rem;
    font-weight: 900;
    text-transform: uppercase;
    letter-spacing: 0.1em;
    color: var(--text-muted);
    margin-bottom: 0.75rem;
    opacity: 0.8;
  }

  .scheme-toggle {
    display: flex;
    gap: 0.5rem;
    background: var(--surface-variant);
    padding: 0.4rem;
    border-radius: 1rem;
  }

  .scheme-toggle button {
    flex: 1;
    height: 40px;
    border-radius: 0.75rem;
    border: none;
    background: none;
    color: var(--text-muted);
    cursor: pointer;
    transition: all 0.2s;
    display: flex;
    align-items: center;
    justify-content: center;
  }

  .scheme-toggle button.active {
    background: var(--md-sys-color-surface-container-lowest);
    color: var(--primary-color);
    box-shadow: var(--shadow-sm);
  }

  .theme-grid {
    display: grid;
    grid-template-columns: repeat(2, 1fr);
    gap: 0.75rem;
  }

  .theme-tile {
    display: flex;
    align-items: center;
    gap: 0.75rem;
    padding: 0.75rem;
    border-radius: 1.25rem;
    border: 1px solid var(--md-sys-color-outline-variant);
    background: var(--md-sys-color-surface-container-low);
    cursor: pointer;
    transition: all 0.2s;
    text-align: left;
  }

  .theme-tile:hover {
    transform: translateY(-2px);
    background: var(--surface-variant-hover);
    border-color: var(--primary-color);
  }

  .theme-tile.active {
    border-color: var(--primary-color);
    background: rgba(var(--primary-rgb), 0.1);
  }

  .color-dot {
    width: 24px;
    height: 24px;
    border-radius: 50%;
    border: 2px solid rgba(255,255,255,0.1);
  }

  .theme-label {
    font-size: 0.85rem;
    font-weight: 700;
    color: var(--text-base);
  }

  /* Custom Option */
  .custom-color-box {
    width: 24px;
    height: 24px;
    border-radius: 50%;
    position: relative;
    overflow: hidden;
    border: 2px solid rgba(255,255,255,0.1);
  }

  .custom-color-box input {
    position: absolute;
    inset: 0;
    opacity: 0;
    cursor: pointer;
  }

  .derived-indicators {
    position: absolute;
    right: -2px;
    top: 50%;
    transform: translateY(-50%);
    display: flex;
    flex-direction: column;
    gap: 2px;
  }

  .derived-indicators .dot {
    width: 6px;
    height: 6px;
    border-radius: 50%;
  }

  .panel-footer {
    margin-top: 2rem;
    text-align: center;
  }

  .panel-footer p {
    font-size: 0.75rem;
    color: var(--text-muted);
    opacity: 0.6;
    margin: 0;
  }
</style>
