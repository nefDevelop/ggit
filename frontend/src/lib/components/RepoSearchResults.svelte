<script lang="ts">
  import { apiFetch } from "$lib/api";
  import { page } from '$app/state';
  import { onMount, getContext } from 'svelte';

  interface Props {
    query?: string;
    branch?: string;
  }

  let { query = '', branch = 'HEAD' }: Props = $props();
  let username = $derived(page.params.username);
  let reponame = $derived(page.params.reponame);

  let results = $state<any[]>([]);
  let loading = $state(false);
  let error = $state('');

  async function performSearch() {
    if (!query || query.length < 2) {
      results = [];
      return;
    }
    loading = true;
    error = '';
    try {
      const res = await apiFetch(`/api/v2/${username}/${reponame}/search?q=${encodeURIComponent(query)}&ref=${branch}`);
      if (res.ok) {
        results = await res.json();
      } else {
        error = 'Search failed.';
      }
    } catch (e) {
      error = 'Failed to connect to search service.';
    } finally {
      loading = false;
    }
  }

  $effect(() => {
    const timer = setTimeout(performSearch, 300);
    return () => clearTimeout(timer);
  });
</script>

<div class="search-results-overlay">
  <header class="results-header">
    <div class="header-main">
        <span class="material-symbols-rounded">search</span>
        <span class="title">Results for "{query}"</span>
        <span class="count-badge">{results.length} matches</span>
    </div>
  </header>

  {#if loading && results.length === 0}
    <div class="loading-state">
      <div class="spinner"></div>
      <p>Grep-ing repository info...</p>
    </div>
  {:else if error}
    <div class="error-msg">{error}</div>
  {:else if results.length === 0 && query.length >= 2 && !loading}
    <div class="no-results">
        <span class="material-symbols-rounded">search_off</span>
        <p>No matches found for this keyword.</p>
    </div>
  {:else}
    <div class="results-list">
      {#each results as match}
        <a href={`/${username`}/${reponame}/src/{branch}/{match.filename}#L{match.line}`} class="match-item">
          <div class="match-file">
            <span class="material-symbols-rounded">description</span>
            <span class="filename">{match.filename}</span>
            <span class="line-num">Line {match.line}</span>
          </div>
          <pre class="match-content"><code>{match.content.trim()}</code></pre>
        </a>
      {/each}
    </div>
  {/if}
</div>

<style>
  .search-results-overlay {
    background: var(--bg-color);
    min-height: 200px;
  }

  .results-header {
    padding: 1rem 1.5rem;
    border-bottom: 1px solid var(--md-sys-color-outline-variant);
    display: flex;
    justify-content: space-between;
    align-items: center;
    background: var(--surface-variant);
  }

  .header-main { display: flex; align-items: center; gap: 0.75rem; }
  .title { font-weight: 700; font-size: 1rem; }
  .count-badge { background: var(--primary-color); color: var(--on-primary); padding: 0.1rem 0.5rem; border-radius: 1rem; font-size: 0.75rem; font-weight: 800; }

  .results-list { display: flex; flex-direction: column; }

  .match-item {
    padding: 1.25rem 1.5rem;
    border-bottom: 1px solid var(--md-sys-color-outline-variant);
    text-decoration: none;
    color: inherit;
    transition: background 0.2s;
  }
  .match-item:hover { background: var(--surface-variant); }

  .match-file { display: flex; align-items: center; gap: 0.6rem; margin-bottom: 0.75rem; }
  .match-file .material-symbols-rounded { font-size: 1.2rem; color: var(--text-muted); }
  .filename { font-weight: 700; color: var(--primary-color); font-size: 0.95rem; }
  .line-num { font-size: 0.8rem; color: var(--text-muted); opacity: 0.6; margin-left: 0.5rem; }

  .match-content {
    margin: 0;
    font-family: 'JetBrains Mono', monospace;
    font-size: 0.85rem;
    color: var(--text-color);
    background: rgba(0,0,0,0.2);
    padding: 0.75rem 1rem;
    border-radius: 0.5rem;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
    border-left: 2px solid var(--primary-color);
  }

  .loading-state, .no-results {
    padding: 4rem 2rem;
    text-align: center;
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 1rem;
    color: var(--text-muted);
  }

  .spinner {
    width: 32px;
    height: 32px;
    border: 3px solid var(--primary-container);
    border-top-color: var(--primary-color);
    border-radius: 50%;
    animation: spin 0.8s linear infinite;
  }
  @keyframes spin { to { transform: rotate(360deg); } }

  .no-results .material-symbols-rounded { font-size: 3rem; opacity: 0.3; }
</style>
