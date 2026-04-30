<script lang="ts">
  import { onMount } from 'svelte';
  import { generateMd3Palette } from '$lib/utils/color';

  let { show = $bindable(false) } = $props();

  const themes = [
    { id: 'default', name: 'Ggit Cobalt', color: '#0b57d0' },
    { id: 'subtle', name: 'Subtle Slate', color: '#607d8b' },
    { id: 'expressive', name: 'MD3 Expressive', color: '#6750a4' },
    { id: 'nord', name: 'Nordic Frost', color: '#5e81ac' },
    { id: 'ayu', name: 'Ayu Mirage', color: '#ff9940' },
    { id: 'monokai', name: 'Retro Monokai', color: '#a6e22e' }
  ];

  let currentTheme = $state('default');
  let currentScheme = $state('light');
  let customColor = $state('#0b57d0');

  const allThemes = $derived([
    ...themes,
    ...(currentTheme === 'custom' || localStorage.getItem('theme') === 'custom' 
      ? [{ id: 'custom', name: 'Personalizado', color: customColor }] 
      : [])
  ]);

  onMount(() => {
    currentTheme = localStorage.getItem('theme') || 'default';
    currentScheme = localStorage.getItem('color-scheme') || (window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light');
    customColor = localStorage.getItem('custom-theme-color') || '#0b57d0';
    applyStyles();
  });

  function setTheme(id: string) {
    currentTheme = id;
    localStorage.setItem('theme', id);
    applyStyles();
  }

  function setScheme(scheme: 'light' | 'dark') {
    currentScheme = scheme;
    localStorage.setItem('color-scheme', scheme);
    applyStyles();
  }

  function applyStyles() {
    document.documentElement.setAttribute('data-theme', currentTheme);
    document.documentElement.setAttribute('data-color-scheme', currentScheme);

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
      };

      for (const [prop, value] of Object.entries(props)) {
        document.documentElement.style.setProperty(prop, value);
      }
      
      const r = parseInt(palette.primary.slice(1, 3), 16);
      const g = parseInt(palette.primary.slice(3, 5), 16);
      const b = parseInt(palette.primary.slice(5, 7), 16);
      document.documentElement.style.setProperty('--primary-rgb', `${r}, ${g}, ${b}`);
    } else {
      // Clear custom props
      const propsToRemove = [
        '--primary-color', '--md-sys-color-primary', '--md-sys-color-on-primary',
        '--md-sys-color-primary-container', '--md-sys-color-on-primary-container',
        '--md-sys-color-secondary', '--md-sys-color-on-secondary',
        '--md-sys-color-secondary-container', '--md-sys-color-on-secondary-container',
        '--md-sys-color-tertiary', '--md-sys-color-on-tertiary',
        '--md-sys-color-tertiary-container', '--md-sys-color-on-tertiary-container',
        '--md-sys-color-error', '--md-sys-color-on-error', '--md-sys-color-error-container', '--md-sys-color-on-error-container',
        '--md-sys-color-background', '--md-sys-color-on-background',
        '--md-sys-color-surface', '--md-sys-color-on-surface',
        '--md-sys-color-surface-variant', '--md-sys-color-on-surface-variant',
        '--md-sys-color-outline', '--md-sys-color-outline-variant',
        '--md-sys-color-surface-container-lowest', '--md-sys-color-surface-container-low',
        '--md-sys-color-surface-container', '--md-sys-color-surface-container-high',
        '--md-sys-color-surface-container-highest', '--primary-rgb'
      ];
      for (const prop of propsToRemove) {
        document.documentElement.style.removeProperty(prop);
      }
    }
  }
</script>

{#if show}
<div class="theme-backdrop" onclick={() => show = false}></div>
<div class="theme-panel island">
  <div class="panel-header">
    <h3>Interfaz y Diseño</h3>
    <button class="close-btn" onclick={() => show = false}>
      <span class="material-symbols-rounded">close</span>
    </button>
  </div>

  <div class="section">
    <span class="section-title">Esquema de Color</span>
    <div class="scheme-toggle">
      <button 
        class:active={currentScheme === 'light'} 
        onclick={() => setScheme('light')}
      >
        <span class="material-symbols-rounded">light_mode</span>
        Claro
      </button>
      <button 
        class:active={currentScheme === 'dark'} 
        onclick={() => setScheme('dark')}
      >
        <span class="material-symbols-rounded">dark_mode</span>
        Oscuro
      </button>
    </div>
  </div>

  <div class="section">
    <span class="section-title">Tonalidad de Acento</span>
    <div class="theme-grid">
      {#each allThemes as theme}
        <button 
          class="theme-option" 
          class:active={currentTheme === theme.id}
          onclick={() => setTheme(theme.id)}
          title={theme.name}
        >
          <div class="color-preview" style="background-color: {theme.color}"></div>
          <span>{theme.name}</span>
        </button>
      {/each}
    </div>
  </div>
</div>
{/if}

<style>
  .theme-backdrop {
    position: fixed;
    inset: 0;
    z-index: 999;
    background: transparent;
  }

  .theme-panel {
    position: fixed;
    bottom: 80px;
    left: 20px;
    width: 280px;
    padding: 1.5rem;
    z-index: 1000;
    box-shadow: var(--shadow-xl);
    border: 1px solid var(--md-sys-color-outline-variant);
  }

  .panel-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 1.5rem;
  }

  .panel-header h3 {
    margin: 0;
    font-size: 1.1rem;
    color: var(--md-sys-color-on-surface);
  }

  .close-btn {
    background: none;
    border: none;
    color: var(--text-muted);
    cursor: pointer;
  }

  .section {
    margin-bottom: 1.5rem;
  }

  .section .section-title {
    display: block;
    font-size: 0.8rem;
    font-weight: 600;
    color: var(--text-muted);
    margin-bottom: 0.75rem;
    text-transform: uppercase;
    letter-spacing: 0.05em;
  }

  .scheme-toggle {
    display: flex;
    gap: 0.5rem;
    background: var(--surface-variant);
    padding: 0.25rem;
    border-radius: 12px;
  }

  .scheme-toggle button {
    flex: 1;
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 0.5rem;
    padding: 0.5rem;
    border-radius: 8px;
    border: none;
    background: none;
    color: var(--md-sys-color-on-surface-variant);
    font-size: 0.9rem;
    cursor: pointer;
    transition: all 0.2s;
  }

  .scheme-toggle button.active {
    background: var(--island-bg);
    color: var(--primary-color);
    box-shadow: var(--shadow-sm);
  }

  .theme-grid {
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
  }

  .theme-option {
    display: flex;
    align-items: center;
    gap: 1rem;
    padding: 0.75rem;
    border-radius: 12px;
    border: 1px solid transparent;
    background: var(--surface-variant);
    color: var(--md-sys-color-on-surface);
    cursor: pointer;
    transition: all 0.2s;
    text-align: left;
  }

  .theme-option:hover {
    background: var(--surface-variant-hover);
  }

  .theme-option.active {
    border-color: var(--primary-color);
    background: var(--primary-surface, var(--surface-variant));
  }

  .color-preview {
    width: 20px;
    height: 20px;
    border-radius: 50%;
    border: 2px solid rgba(255,255,255,0.2);
  }

  .theme-option span {
    font-size: 0.9rem;
    font-weight: 500;
  }
</style>
