<script lang="ts">
  import { onMount, getContext } from 'svelte';
  import { apiFetch } from '$lib/api';
  import MarkdownRenderer from './ui/MarkdownRenderer.svelte';
  import Card from './ui/Card.svelte';
  import Button from './ui/Button.svelte';
  import PageHeader from './ui/PageHeader.svelte';

  const repo: any = getContext('repo');
  let adrFolder = $derived(repo.data.adr_folder_name || 'ADR');
  let currentPath = $state(repo.data.adr_folder_name || 'ADR');
  let entries = $state<any[]>([]);
  let currentFile = $state<any>(null);
  let loading = $state(true);
  let error = $state('');

  async function loadPath(path: string) {
    loading = true;
    error = '';
    try {
      const res = await apiFetch(`/api/v2/${repo.data.owner.username}/${repo.data.name}/view/${path}`);
      if (res.ok) {
        const data = await res.json();
        if (data.type === 'directory') {
          entries = data.entries;
          currentFile = null;
          currentPath = path;
        } else {
          currentFile = data;
          currentPath = path;
        }
      } else {
        // If 'ADR' folder doesn't exist, try 'docs' as fallback for tech notes
        if (path === adrFolder) {
           return loadPath('docs');
        }
        error = `ADR folder "${path}" not found. Create a folder to track architectural decisions.`;
      }
    } catch (e) {
      error = 'Failed to load architectural records.';
    } finally {
      loading = false;
    }
  }

  function isADR(filename: string) {
      return true; // Everything in this component is treated with the ADR style
  }

  onMount(() => {
    // Check if the folder exists in root before trying to load it
    const hasADR = repo.entries?.some((e: any) => e.name === adrFolder && e.type === 'dir');
    const hasDocs = repo.entries?.some((e: any) => e.name === 'docs' && e.type === 'dir');

    if (hasADR) {
      loadPath(adrFolder);
    } else if (hasDocs) {
      loadPath('docs');
    } else {
      loading = false;
      error = `No documentation folder found (looked for "${adrFolder}" or "docs").`;
    }
  });

  function goBack() {
    const parts = currentPath.split('/');
    parts.pop();
    loadPath(parts.join('/') || adrFolder);
  }
</script>

<div class="repo-docs-container">

  <div class="docs-layout">
    {#if loading}
      <div class="loading-state">
        <div class="spinner"></div>
        <p>Loading documentation...</p>
      </div>
    {:else if error}
      <div class="empty-state">
        <span class="material-symbols-rounded">error_outline</span>
        <h3>No documentation found</h3>
        <p>{error}</p>
        <div class="hint">
            Tip: Create a <code>docs/</code> or <code>ADR/</code> folder in your root branch to see it here.
        </div>
      </div>
    {:else if currentFile}
      <div class="file-view">
        <div class="file-header">
            <div class="header-main-row">
              <h2 class="file-name">
                  <span class="material-symbols-rounded">{isADR(currentFile.name) ? 'gavel' : 'description'}</span>
                  {currentFile.name}
              </h2>
              <Button variant="secondary" size="sm" icon="arrow_back" onclick={goBack}>Back to List</Button>
            </div>
            {#if currentFile.last_commit}
                <div class="last-update">
                    Last updated {new Date(currentFile.last_commit.author.when).toLocaleDateString()} by {currentFile.last_commit.author.name}
                </div>
            {/if}
        </div>
        
        <Card class="markdown-body-card">
            <MarkdownRenderer content={currentFile.content} />
        </Card>
      </div>
    {:else}
      <div class="docs-grid">
        {#each entries.filter(e => e.type === 'dir' || e.name.endsWith('.md')) as entry}
          <button class="doc-tile tonal" onclick={() => loadPath(`${currentPath}/${entry.name}`)}>
            <div class="tile-icon">
              <span class="material-symbols-rounded">
                {entry.type === 'dir' ? 'folder' : (isADR(entry.name) ? 'gavel' : 'description')}
              </span>
            </div>
            <div class="tile-info">
              <span class="tile-name">{entry.name}</span>
              {#if entry.last_commit}
                <span class="tile-meta">{entry.last_commit.summary}</span>
              {/if}
            </div>
          </button>
        {/each}
      </div>
    {/if}
  </div>
</div>

<style>
  .repo-docs-container { display: flex; flex-direction: column; gap: 0; margin: -2rem; }
  .docs-layout { min-height: 400px; padding: 2rem; }
  
  .header-main-row { display: flex; justify-content: space-between; align-items: center; width: 100%; }
  .docs-layout { min-height: 400px; }

  .loading-state { display: flex; flex-direction: column; align-items: center; justify-content: center; height: 300px; gap: 1rem; color: var(--text-muted); }
  .spinner { width: 32px; height: 32px; border: 3px solid var(--md-sys-color-outline-variant); border-top-color: var(--primary-color); border-radius: 50%; animation: spin 1s infinite linear; }
  @keyframes spin { to { transform: rotate(360deg); } }

  .hint { margin-top: 1rem; font-size: 0.85rem; opacity: 0.7; }

  .docs-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(280px, 1fr)); gap: 1rem; }
  .doc-tile { border: none; padding: 1.25rem; border-radius: 1.25rem; display: flex; align-items: center; gap: 1.25rem; text-align: left; transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1); cursor: pointer; }
  .doc-tile:hover { transform: translateY(-2px); background: var(--md-sys-color-surface-container-high); box-shadow: var(--shadow-md); }
  
  .tile-icon { background: rgba(var(--primary-rgb), 0.1); color: var(--primary-color); padding: 0.75rem; border-radius: 1rem; display: flex; }
  .tile-info { display: flex; flex-direction: column; gap: 0.15rem; min-width: 0; }
  .tile-name { font-weight: 800; color: var(--text-color); font-size: 0.95rem; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
  .tile-meta { font-size: 0.75rem; color: var(--text-muted); overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }

  .file-view { display: flex; flex-direction: column; gap: 1.5rem; animation: fadeIn 0.3s ease-out; }
  @keyframes fadeIn { from { opacity: 0; transform: translateY(10px); } to { opacity: 1; transform: translateY(0); } }

  .file-header { display: flex; flex-direction: column; gap: 0.25rem; }
  .file-name { margin: 0; display: flex; align-items: center; gap: 0.75rem; font-size: 1.5rem; font-weight: 900; }
  .file-name span { font-size: 2rem; color: var(--primary-color); }
  .last-update { font-size: 0.85rem; color: var(--text-muted); font-weight: 600; }

  :global(.markdown-body-card) { padding: 2.5rem !important; border-radius: 2rem !important; background: var(--bg-color) !important; border: 1px solid var(--md-sys-color-outline-variant) !important; }
</style>
