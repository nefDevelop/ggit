<script lang="ts">
  import { page } from '$app/state';
  import { onMount } from 'svelte';

  const reponame = $derived(page.params.reponame);
  const isRepoContext = $derived(!!reponame);

  let query = $state('');
  let searchInput = $state<HTMLInputElement | null>(null);

  function handleKeydown(e: KeyboardEvent) {
    if (e.key === '/' && document.activeElement?.tagName !== 'INPUT') {
      e.preventDefault();
      searchInput?.focus();
    }
  }

  onMount(() => {
    window.addEventListener('keydown', handleKeydown);
    return () => window.removeEventListener('keydown', handleKeydown);
  });

  $effect(() => {
     window.dispatchEvent(new CustomEvent('ggit:search', { detail: query }));
  });
</script>

<div class="search-pill-container">
  <div class="search-pill">
    <span class="material-symbols-rounded icon">
      {isRepoContext ? 'database' : 'search'}
    </span>
    <input 
      bind:this={searchInput}
      bind:value={query}
      placeholder={isRepoContext ? `Search in ${reponame}...` : 'Search repositories...'} 
      class="global-search-input"
    />
    <div class="shortcut-hint">/</div>
  </div>
</div>

<style>
  .search-pill-container {
    width: 100%;
    max-width: 600px;
  }

  .search-pill {
    display: flex;
    align-items: center;
    gap: 0.75rem;
    background: var(--md3-surface-container);
    padding: 0.6rem 1.25rem;
    border-radius: 3rem;
    border: 1px solid var(--border-color);
    transition: all 0.2s ease;
  }

  .search-pill:focus-within {
    background: var(--md3-surface-3);
    border-color: var(--primary-color);
    box-shadow: 0 0 20px rgba(0,0,0,0.3);
    transform: scale(1.02);
  }

  .icon {
    font-size: 1.25rem;
    color: var(--primary-color);
    opacity: 0.8;
  }

  .global-search-input {
    flex: 1;
    background: transparent;
    border: none;
    color: var(--text-color);
    font-size: 0.95rem;
    font-weight: 500;
    outline: none;
  }

  .shortcut-hint {
    font-size: 0.75rem;
    font-weight: 800;
    color: var(--text-muted);
    background: var(--bg-color);
    width: 20px;
    height: 20px;
    display: flex;
    align-items: center;
    justify-content: center;
    border-radius: 4px;
    opacity: 0.5;
  }
</style>
