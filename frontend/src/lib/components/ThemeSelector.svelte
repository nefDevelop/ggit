<script lang="ts">
  import { onMount } from 'svelte';
  import { i18n } from '$lib/i18n.svelte';
  import Card from './ui/Card.svelte';

  let currentTheme = $state('default');
  let currentScheme = $state('auto');

  const themes = [
    { id: 'default', name: 'Ggit Dynamic', color: '#0b57d0' },
    { id: 'subtle', name: 'Subtle Slate', color: '#607d8b' },
    { id: 'ayu', name: 'Ayu Sunrise', color: '#ffb454' },
    { id: 'monokai', name: 'Monokai Pro', color: '#a6e22e' },
    { id: 'nord', name: 'Nordic Frost', color: '#88c0d0' }
  ];

  const schemes = [
    { id: 'auto', name: 'Sistema Dual', icon: 'brightness_auto' },
    { id: 'light', name: 'Modo Claro', icon: 'light_mode' },
    { id: 'dark', name: 'Modo Oscuro', icon: 'dark_mode' }
  ];

  onMount(() => {
    currentTheme = localStorage.getItem('theme') || 'default';
    currentScheme = localStorage.getItem('color-scheme') || 'auto';
  });

  function applyTheme(themeId: string) {
    currentTheme = themeId;
    localStorage.setItem('theme', themeId);
    document.documentElement.setAttribute('data-theme', themeId);
  }

  function applyScheme(schemeId: string) {
    currentScheme = schemeId;
    localStorage.setItem('color-scheme', schemeId);

    if (schemeId === 'dark' || (schemeId === 'auto' && window.matchMedia('(prefers-color-scheme: dark)').matches)) {
      document.documentElement.setAttribute('data-color-scheme', 'dark');
    } else {
      document.documentElement.setAttribute('data-color-scheme', 'light');
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
            <div class="color-disk shadow-sm" style="background-color: {theme.color};"></div>
            <div class="theme-info">
              <span class="name">{theme.name}</span>
              <span class="hex">{theme.color}</span>
            </div>
            {#if currentTheme === theme.id}
               <span class="material-symbols-rounded check-icon">check_circle</span>
            {/if}
          </button>
        {/each}
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
    padding: 1.5rem 1.75rem; border-radius: 2rem; border: 1px solid var(--md-sys-color-outline-variant);
    background: rgba(0,0,0,0.1); transition: all 0.3s;
  }
  .theme-modular-tile:hover { transform: scale(1.02); background: rgba(255,255,255,0.02); border-color: var(--primary-color); }
  .theme-modular-tile.active { border-color: var(--primary-color); background: rgba(var(--primary-rgb), 0.05); }

  .color-disk { width: 3.5rem; height: 3.5rem; border-radius: 50%; border: 4px solid rgba(255,255,255,0.1); flex-shrink: 0; }
  .theme-info { display: flex; flex-direction: column; flex: 1; text-align: left; }
  .theme-info .name { font-weight: 900; font-size: 1.1rem; color: var(--text-color); }
  .theme-info .hex { font-family: 'JetBrains Mono', monospace; font-size: 0.75rem; color: var(--text-muted); opacity: 0.6; }

  .check-icon { color: var(--primary-color); font-size: 1.75rem; }

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