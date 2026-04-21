<script lang="ts">
  import { onMount } from 'svelte';

  let currentTheme = $state('default');
  let currentScheme = $state('auto');

  const themes = [
    { id: 'default', name: 'Default (Blue)', color: '#0b57d0' },
    { id: 'ayu', name: 'Ayu (Orange)', color: '#ffb454' },
    { id: 'monokai', name: 'Monokai (Green)', color: '#a6e22e' },
    { id: 'nord', name: 'Nord (Ice Blue)', color: '#88c0d0' }
  ];

  const schemes = [
    { id: 'auto', name: 'System Auto', icon: 'brightness_auto' },
    { id: 'light', name: 'Light Mode', icon: 'light_mode' },
    { id: 'dark', name: 'Dark Mode', icon: 'dark_mode' }
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

<div class="appearance-modular-view">
  <header class="section-header-modern">
    <div class="header-info">
      <h3>Appearance</h3>
      <p>Customize the look and feel of your workspace.</p>
    </div>
  </header>
  
  <div class="appearance-content">
    <div class="option-group">
      <h4 class="group-title">Color Mode</h4>
      <div class="scheme-grid">
        {#each schemes as scheme}
          <button 
            type="button" 
            class="scheme-btn" 
            class:active={currentScheme === scheme.id}
            onclick={() => applyScheme(scheme.id)}
          >
            <span class="material-symbols-rounded">{scheme.icon}</span>
            <div class="label-box">
              <span class="name">{scheme.name}</span>
            </div>
          </button>
        {/each}
      </div>
    </div>

    <div class="option-group">
      <h4 class="group-title">Accent Theme</h4>
      <div class="theme-grid">
        {#each themes as theme}
          <button 
            type="button" 
            class="theme-card" 
            class:active={currentTheme === theme.id}
            onclick={() => applyTheme(theme.id)}
          >
            <div class="color-preview" style="background-color: {theme.color};"></div>
            <span class="theme-name">{theme.name}</span>
          </button>
        {/each}
      </div>
    </div>
  </div>
</div>

<style>
  .appearance-modular-view { display: flex; flex-direction: column; gap: 2.5rem; }
  
  .section-header-modern { border-bottom: 1px solid var(--border-color); padding-bottom: 1.5rem; }
  .section-header-modern h3 { margin: 0; font-size: 1.5rem; font-weight: 900; }
  .section-header-modern p { margin: 0.5rem 0 0 0; color: var(--text-muted); font-size: 0.95rem; }

  .appearance-content { display: flex; flex-direction: column; gap: 3rem; }

  .option-group { display: flex; flex-direction: column; gap: 1.25rem; }
  .group-title { font-size: 0.75rem; font-weight: 800; text-transform: uppercase; color: var(--text-muted); letter-spacing: 0.05em; padding-left: 0.5rem; }

  /* Scheme Selector */
  .scheme-grid { display: grid; grid-template-columns: repeat(3, 1fr); gap: 1.5rem; }
  .scheme-btn {
    display: flex; flex-direction: column; align-items: center; justify-content: center; gap: 1rem;
    background: var(--surface-variant);
    border: 1px solid var(--border-color);
    padding: 2rem;
    border-radius: 1.25rem;
    cursor: pointer; transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1);
    color: var(--text-muted);
  }
  .scheme-btn:hover { background: rgba(255,255,255,0.02); transform: translateY(-3px); }
  .scheme-btn.active { border-color: var(--primary-color); background: rgba(var(--primary-rgb), 0.05); color: var(--primary-color); box-shadow: var(--shadow-md); }
  .scheme-btn .material-symbols-rounded { font-size: 1.75rem; }
  .scheme-btn .name { font-weight: 700; font-size: 0.9rem; }

  /* Theme Grid */
  .theme-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(160px, 1fr)); gap: 1.5rem; }
  .theme-card {
    display: flex; align-items: center; gap: 1.25rem;
    background: var(--surface-variant);
    border: 1px solid var(--border-color);
    padding: 1.25rem;
    border-radius: 1.25rem;
    cursor: pointer; transition: all 0.2s; color: var(--text-color);
  }
  .theme-card:hover { transform: translateY(-3px); background: rgba(255,255,255,0.02); }
  .theme-card.active { border-color: var(--primary-color); background: rgba(var(--primary-rgb), 0.05); }
  
  .color-preview { width: 2.5rem; height: 2.5rem; border-radius: 1rem; flex-shrink: 0; box-shadow: var(--shadow-sm); border: 2px solid rgba(255,255,255,0.1); }
  .theme-name { font-weight: 700; font-size: 0.9rem; }

  @media (max-width: 600px) {
    .scheme-grid { grid-template-columns: 1fr; }
    .theme-grid { grid-template-columns: 1fr; }
  }
</style>