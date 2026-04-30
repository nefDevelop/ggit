<script lang="ts">
  import { onMount } from 'svelte';
  import { i18n } from '$lib/i18n.svelte';
  import Card from './ui/Card.svelte';
  import { generateMd3Palette } from '$lib/utils/color';

  let currentTheme = $state('default');
  let currentScheme = $state('auto');
  let customColor = $state('#0b57d0');
  let derivedPalette = $derived(generateMd3Palette(customColor, currentScheme === 'dark' || (currentScheme === 'auto' && typeof window !== 'undefined' && window.matchMedia('(prefers-color-scheme: dark)').matches)));

  const themes = [
    { 
      id: 'default', 
      name: 'Ggit Dynamic', 
      palette: ['#0b57d0', '#00639b', '#565e71', '#f0f2f5'] 
    },
    { 
      id: 'ayu', 
      name: 'Ayu Sunrise', 
      palette: ['#ff9940', '#5c6166', '#73d031', '#fafafa'] 
    },
    { 
      id: 'monokai', 
      name: 'Monokai Pro', 
      palette: ['#a6e22e', '#ccccc6', '#fd971f', '#272822'] 
    },
    { 
      id: 'nord', 
      name: 'Nordic Frost', 
      palette: ['#5e81ac', '#81a1c1', '#88c0d0', '#eceff4'] 
    }
  ];

  const schemes = [
    { id: 'auto', name: 'Sistema Dual', icon: 'brightness_auto' },
    { id: 'light', name: 'Modo Claro', icon: 'light_mode' },
    { id: 'dark', name: 'Modo Oscuro', icon: 'dark_mode' }
  ];

  onMount(() => {
    currentTheme = localStorage.getItem('theme') || 'default';
    currentScheme = localStorage.getItem('color-scheme') || 'auto';
    customColor = localStorage.getItem('custom-theme-color') || '#0b57d0';
    
    if (currentTheme === 'custom') {
      applyCustomColor(customColor);
    }
  });

  function applyTheme(themeId: string) {
    currentTheme = themeId;
    localStorage.setItem('theme', themeId);
    document.documentElement.setAttribute('data-theme', themeId);
    
    if (themeId !== 'custom') {
      // Clear custom color overrides if switching to preset
      document.documentElement.style.removeProperty('--primary-color');
      document.documentElement.style.removeProperty('--md-sys-color-primary');
      const propsToRemove = [
        '--md-sys-color-on-primary', '--md-sys-color-primary-container', '--md-sys-color-on-primary-container',
        '--md-sys-color-secondary', '--md-sys-color-on-secondary', '--md-sys-color-secondary-container', '--md-sys-color-on-secondary-container',
        '--md-sys-color-tertiary', '--md-sys-color-on-tertiary', '--md-sys-color-tertiary-container', '--md-sys-color-on-tertiary-container',
        '--md-sys-color-error', '--md-sys-color-on-error', '--md-sys-color-error-container', '--md-sys-color-on-error-container',
        '--md-sys-color-background', '--md-sys-color-on-background', '--md-sys-color-surface', '--md-sys-color-on-surface',
        '--md-sys-color-surface-variant', '--md-sys-color-on-surface-variant', '--md-sys-color-outline', '--md-sys-color-outline-variant',
        '--md-sys-color-surface-container-lowest', '--md-sys-color-surface-container-low', '--md-sys-color-surface-container',
        '--md-sys-color-surface-container-high', '--md-sys-color-surface-container-highest', '--primary-rgb'
      ];
      for (const prop of propsToRemove) {
        document.documentElement.style.removeProperty(prop);
      }
    } else {
      applyCustomColor(customColor);
    }
  }

  function applyCustomColor(hex: string) {
    customColor = hex;
    localStorage.setItem('custom-theme-color', hex);
    
    if (currentTheme === 'custom') {
      const isDark = currentScheme === 'dark' || (currentScheme === 'auto' && window.matchMedia('(prefers-color-scheme: dark)').matches);
      const palette = generateMd3Palette(hex, isDark);
      
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
    }
  }

  function applyScheme(schemeId: string) {
    currentScheme = schemeId;
    localStorage.setItem('color-scheme', schemeId);

    if (schemeId === 'dark' || (schemeId === 'auto' && window.matchMedia('(prefers-color-scheme: dark)').matches)) {
      document.documentElement.setAttribute('data-color-scheme', 'dark');
    } else {
      document.documentElement.setAttribute('data-color-scheme', 'light');
    }

    if (currentTheme === 'custom') {
      applyCustomColor(customColor);
    }
  }
</script>

<Card class="settings-premium-card" padding="lg">
  <header class="section-header">
    <div class="title-with-icon">
      <span class="material-symbols-rounded hero-icon">palette</span>
      <div>
        <h2>{i18n.t('settings.appearance') || 'Apariencia'}</h2>
        <p class="subtitle">Personaliza la interfaz y el esquema de colores de tu espacio</p>
      </div>
    </div>
  </header>

  <div class="modular-stack">
    <!-- COLOR SCHEME -->
    <div class="config-module">
      <div class="module-tag">
        <span class="material-symbols-rounded">contrast</span>
        Modo de Color
      </div>
      <div class="scheme-grid">
        {#each schemes as scheme}
          <button 
            type="button" 
            class="scheme-modular-card island" 
            class:active={currentScheme === scheme.id}
            onclick={() => applyScheme(scheme.id)}
          >
            <div class="scheme-icon-box">
              <span class="material-symbols-rounded">{scheme.icon}</span>
            </div>
            <span class="scheme-name">{scheme.name}</span>
            {#if currentScheme === scheme.id}
              <div class="active-dot"></div>
            {/if}
          </button>
        {/each}
      </div>
    </div>

    <!-- ACCENT THEME -->
    <div class="config-module">
      <div class="module-tag">
        <span class="material-symbols-rounded">colorize</span>
        Tema de Acento
      </div>
      <div class="theme-modular-grid">
        {#each themes as theme}
          <button 
            type="button" 
            class="theme-modular-tile tonal" 
            class:active={currentTheme === theme.id}
            onclick={() => applyTheme(theme.id)}
          >
            <div class="palette-preview">
               <div class="seed-color shadow-sm" style="background-color: {theme.palette[0]};"></div>
               <div class="derived-dots">
                  <div class="dot" style="background-color: {theme.palette[1]};"></div>
                  <div class="dot" style="background-color: {theme.palette[2]};"></div>
                  <div class="dot" style="background-color: {theme.palette[3]};"></div>
               </div>
            </div>
            <div class="theme-info">
              <span class="name">{theme.name}</span>
              <span class="hex">{theme.palette[0]}</span>
            </div>
            {#if currentTheme === theme.id}
               <span class="material-symbols-rounded check-icon">check_circle</span>
            {/if}
          </button>
        {/each}

        <!-- CUSTOM THEME BOX -->
        <button 
          type="button" 
          class="theme-modular-tile tonal custom-tile" 
          class:active={currentTheme === 'custom'}
          onclick={() => applyTheme('custom')}
        >
          <div class="palette-preview">
             <div class="seed-color shadow-sm custom-preview" style="background-color: {customColor};">
                <input 
                  type="color" 
                  bind:value={customColor} 
                  oninput={(e) => applyCustomColor(e.currentTarget.value)} 
                  title="Elegir color"
                />
             </div>
             <div class="derived-dots">
                <div class="dot" style="background-color: {derivedPalette.secondary};" title="Secondary"></div>
                <div class="dot" style="background-color: {derivedPalette.tertiary};" title="Tertiary"></div>
                <div class="dot" style="background-color: {derivedPalette.primaryContainer};" title="Container"></div>
             </div>
          </div>
          <div class="theme-info">
            <span class="name">Personalizado</span>
            <span class="hex">{customColor}</span>
          </div>
          {#if currentTheme === 'custom'}
             <span class="material-symbols-rounded check-icon">check_circle</span>
          {/if}
        </button>
      </div>
    </div>
  </div>
</Card>

<style>

  .section-header { margin-bottom: 4rem; }
  .title-with-icon { display: flex; align-items: center; gap: 2rem; }
  .hero-icon { font-size: 3rem; color: var(--primary-color); background: rgba(var(--primary-rgb), 0.1); padding: 1.25rem; border-radius: 2rem; }
  h2 { margin: 0; font-size: 2.25rem; font-weight: 950; letter-spacing: -0.04em; }
  .subtitle { margin: 0.5rem 0 0 0; color: var(--text-muted); font-size: 1.1rem; }

  .modular-stack { display: flex; flex-direction: column; gap: 4rem; }
  .config-module { display: flex; flex-direction: column; gap: 1.5rem; }

  .module-tag { font-size: 0.8rem; font-weight: 950; text-transform: uppercase; color: var(--primary-color); letter-spacing: 0.15em; display: flex; align-items: center; gap: 0.75rem; border-bottom: 2px solid var(--md-sys-color-outline-variant); padding-bottom: 0.75rem; width: fit-content; opacity: 0.8; }
  .module-tag span { font-size: 1.25rem; }

  /* Scheme Selector */
  .scheme-grid { display: grid; grid-template-columns: repeat(3, 1fr); gap: 1.25rem; }
  .scheme-modular-card {
    position: relative; display: flex; flex-direction: column; align-items: center; gap: 1.25rem;
    padding: 3rem 2rem; border-radius: 2.5rem; border: 1px solid var(--md-sys-color-outline-variant);
    background: rgba(255,255,255,0.01); transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);
  }
  .scheme-modular-card:hover { transform: translateY(-8px); background: rgba(255,255,255,0.03); border-color: var(--primary-color); }
  .scheme-modular-card.active { border-color: var(--primary-color); background: rgba(var(--primary-rgb), 0.1); box-shadow: var(--shadow-xl); }

  .scheme-icon-box { width: 64px; height: 64px; border-radius: 1.5rem; background: rgba(var(--primary-rgb), 0.05); display: flex; align-items: center; justify-content: center; color: var(--primary-color); }
  .scheme-icon-box span { font-size: 2.5rem; }
  .scheme-name { font-weight: 900; font-size: 1.1rem; color: var(--text-color); }

  .active-dot { position: absolute; top: 1.5rem; right: 1.5rem; width: 12px; height: 12px; background: var(--primary-color); border-radius: 50%; box-shadow: 0 0 15px var(--primary-color); }

  /* Theme Grid */
  .theme-modular-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(280px, 1fr)); gap: 1rem; }
  .theme-modular-tile {
    display: flex; align-items: center; gap: 1.5rem;
    padding: 1.25rem 1.5rem; border-radius: 2rem; border: 1px solid var(--md-sys-color-outline-variant);
    background: rgba(0,0,0,0.1); transition: all 0.3s;
    cursor: pointer;
    text-align: left;
  }
  .theme-modular-tile:hover { transform: scale(1.02); background: rgba(255,255,255,0.02); border-color: var(--primary-color); }
  .theme-modular-tile.active { border-color: var(--primary-color); background: rgba(var(--primary-rgb), 0.05); }

  .palette-preview { display: flex; align-items: center; gap: 0.75rem; position: relative; }
  .seed-color { width: 3.5rem; height: 3.5rem; border-radius: 50%; border: 3px solid rgba(255,255,255,0.1); flex-shrink: 0; position: relative; overflow: hidden; }
  
  .derived-dots { display: flex; flex-direction: column; gap: 4px; }
  .dot { width: 10px; height: 10px; border-radius: 50%; border: 1px solid rgba(255,255,255,0.05); }
  .dot.auto { background: var(--text-muted); opacity: 0.3; }

  .custom-preview input { position: absolute; inset: 0; opacity: 0; cursor: pointer; width: 100%; height: 100%; }

  .theme-info { display: flex; flex-direction: column; flex: 1; min-width: 0; }
  .theme-info .name { font-weight: 900; font-size: 1rem; color: var(--text-color); white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }
  .theme-info .hex { font-family: 'JetBrains Mono', monospace; font-size: 0.7rem; color: var(--text-muted); opacity: 0.6; }

  .check-icon { color: var(--primary-color); font-size: 1.5rem; }

  @media (max-width: 800px) {
    .section-header { margin-bottom: 2rem; }
    h2 { font-size: 1.5rem; }
    .hero-icon { font-size: 1.5rem; padding: 0.6rem; border-radius: 1rem; }
    .title-with-icon { gap: 0.75rem; }
    .subtitle { font-size: 0.9rem; }

    .modular-stack { gap: 2rem; }
    .scheme-grid { grid-template-columns: 1fr; }
    .theme-modular-grid { grid-template-columns: 1fr; }
    .scheme-modular-card { padding: 1.5rem; border-radius: 1.5rem; }
    .theme-modular-tile { padding: 1rem; border-radius: 1.5rem; }
  }
</style>