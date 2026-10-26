<script lang="ts">
  import { apiFetch } from "$lib/api";
  import { page } from '$app/state';
  import { onMount } from 'svelte';
  

import Search from "lucide-svelte/icons/search";
import SearchX from "lucide-svelte/icons/search-x";
import FileText from "lucide-svelte/icons/file-text";



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
        <Search size={18} />
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
        <div class="empty-icon-ring">
          <SearchX size={40} />
        </div>
        <h2>No matches found</h2>
        <p>No matches found for this keyword in the current branch.</p>
    </div>
  {:else}
    <div class="results-list">
      {#each results as match}
        <a href={`/${username}/${reponame}/src/${branch}/${match.filename}#L${match.line}`} class="match-item">
          <div class="match-file">
            <FileText size={16} />
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
    background: var(--md-sys-color-secondary-container);
    color: var(--md-sys-color-on-secondary-container);
  }

  .header-main { display: flex; align-items: center; gap: 0.75rem; }
  .title { font-weight: 700; font-size: 1rem; }
  .count-badge { background: var(--secondary); color: var(--on-secondary); padding: 0.1rem 0.5rem; border-radius: 1rem; font-size: 0.75rem; font-weight: 800; }

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
  .filename { font-weight: 700; color: var(--primary-color); font-size: 0.95rem; }
  .line-num { font-size: 0.8rem; color: var(--text-muted); opacity: 0.6; margin-left: 0.5rem; }

  .match-content {
    margin: 0;
    font-family: 'JetBrains Mono', monospace;
    font-size: 0.85rem;
    color: var(--text-color);
    background: var(--sunken-2);
    padding: 0.75rem 1rem;
    border-radius: 0.5rem;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
    border-left: 2px solid var(--tertiary);
  }

  .loading-state, .no-results {
    padding: 6rem 2rem;
    text-align: center;
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 1.25rem;
    color: var(--text-muted);
  }

  .no-results h2 {
    margin: 0;
    font-size: 1.75rem;
    font-weight: 950;
    color: var(--text-color);
  }

  .no-results p {
    max-width: 400px;
    margin: 0;
    font-size: 1.1rem;
    line-height: 1.5;
  }

  .empty-icon-ring {
    width: 80px;
    height: 80px;
    border: 2px dashed var(--md-sys-color-outline-variant);
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;
    margin-bottom: 1rem;
    color: var(--primary-color);
    background: rgba(var(--primary-rgb), 0.05);
  }

  .spinner {
    width: 32px;
    height: 32px;
    border: 3px solid var(--md-sys-color-tertiary-container);
    border-top-color: var(--tertiary);
    border-radius: 50%;
    animation: spin 0.8s linear infinite;
  }
  @keyframes spin { to { transform: rotate(360deg); } }

</style>
