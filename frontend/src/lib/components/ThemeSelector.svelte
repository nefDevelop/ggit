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

<div class="settings-card island mt-4">
  <header class="card-header">
    <h3 class="card-title">Appearance</h3>
    <p class="card-desc">Customize the look and feel of your workspace.</p>
  </header>
  
  <div class="theme-options">
    <div class="option-group">
      <h4>Color Mode</h4>
      <div class="button-row">
        {#each schemes as scheme}
          <button 
            type="button" 
            class="selector-btn" 
            class:active={currentScheme === scheme.id}
            onclick={() => applyScheme(scheme.id)}
          >
            <span class="material-symbols-rounded">{scheme.icon}</span>
            {scheme.name}
          </button>
        {/each}
      </div>
    </div>

    <div class="option-group">
      <h4>Accent Theme</h4>
      <div class="theme-grid">
        {#each themes as theme}
          <button 
            type="button" 
            class="theme-card" 
            class:active={currentTheme === theme.id}
            onclick={() => applyTheme(theme.id)}
          >
            <div class="color-preview" style="background-color: {theme.color};"></div>
            <span>{theme.name}</span>
          </button>
        {/each}
      </div>
    </div>
  </div>
</div>

<style>
  .settings-card { display: flex; flex-direction: column; margin-top: 2rem; }
  .card-header { padding: 1.5rem 2rem; border-bottom: 1px solid var(--border-color); }
  .card-title { margin: 0; font-size: 1.25rem; font-weight: 700; }
  .card-desc { margin: 0.5rem 0 0 0; color: var(--text-muted); font-size: 0.9rem; }

  .theme-options { padding: 2rem; display: flex; flex-direction: column; gap: 2rem; }

  .option-group h4 { margin: 0 0 1rem 0; font-size: 0.85rem; font-weight: 800; text-transform: uppercase; color: var(--text-muted); letter-spacing: 0.05em; }

  .button-row { display: flex; gap: 1rem; flex-wrap: wrap; }
  .selector-btn {
    display: flex; align-items: center; gap: 0.5rem;
    background: var(--surface-variant);
    border: 2px solid transparent;
    color: var(--text-color);
    padding: 0.75rem 1.25rem;
    border-radius: 1rem;
    font-weight: 600; font-family: inherit;
    cursor: pointer; transition: all 0.2s ease;
  }
  .selector-btn:hover { background: var(--surface-variant-hover); }
  .selector-btn.active { background: rgba(var(--primary-rgb), 0.1); border-color: var(--primary-color); color: var(--primary-color); }
  .selector-btn .material-symbols-rounded { font-size: 1.25rem; }

  .theme-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(140px, 1fr)); gap: 1rem; }
  .theme-card {
    display: flex; flex-direction: column; align-items: center; gap: 0.75rem;
    background: var(--surface-variant);
    border: 2px solid transparent;
    padding: 1.25rem 1rem;
    border-radius: 1.25rem;
    cursor: pointer; transition: all 0.2s ease; color: var(--text-color); font-weight: 600; font-family: inherit;
  }
  .theme-card:hover { background: var(--surface-variant-hover); transform: translateY(-2px); }
  .theme-card.active { border-color: var(--primary-color); background: rgba(var(--primary-rgb), 0.05); }
  
  .color-preview { width: 3rem; height: 3rem; border-radius: 50%; box-shadow: 0 4px 12px rgba(0,0,0,0.2); border: 2px solid rgba(255,255,255,0.1); }
</style>