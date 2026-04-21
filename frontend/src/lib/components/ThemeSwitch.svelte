<script lang="ts">
  import { onMount } from 'svelte';

  let { show = $bindable(false) } = $props();

  const themes = [
    { id: 'default', name: 'Ggit Cobalt', color: '#0b57d0' },
    { id: 'expressive', name: 'MD3 Expressive', color: '#6750a4' },
    { id: 'nord', name: 'Nordic Frost', color: '#5e81ac' },
    { id: 'ayu', name: 'Ayu Mirage', color: '#ff9940' },
    { id: 'monokai', name: 'Retro Monokai', color: '#a6e22e' }
  ];

  let currentTheme = $state('default');
  let currentScheme = $state('light');

  onMount(() => {
    currentTheme = localStorage.getItem('theme') || 'default';
    currentScheme = localStorage.getItem('color-scheme') || (window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light');
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
  }
</script>

{#if show}
<div class="theme-panel island">
  <div class="panel-header">
    <h3>Interfaz y Diseño</h3>
    <button class="close-btn" on:click={() => show = false}>
      <span class="material-symbols-rounded">close</span>
    </button>
  </div>

  <div class="section">
    <label>Esquema de Color</label>
    <div class="scheme-toggle">
      <button 
        class:active={currentScheme === 'light'} 
        on:click={() => setScheme('light')}
      >
        <span class="material-symbols-rounded">light_mode</span>
        Claro
      </button>
      <button 
        class:active={currentScheme === 'dark'} 
        on:click={() => setScheme('dark')}
      >
        <span class="material-symbols-rounded">dark_mode</span>
        Oscuro
      </button>
    </div>
  </div>

  <div class="section">
    <label>Tonalidad de Acento</label>
    <div class="theme-grid">
      {#each themes as theme}
        <button 
          class="theme-option" 
          class:active={currentTheme === theme.id}
          on:click={() => setTheme(theme.id)}
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

  .section label {
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
