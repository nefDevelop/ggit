<script lang="ts">
  import type { RepoContext } from "$lib/types";
  import { onMount, getContext } from 'svelte';
  import { apiFetch } from '$lib/api';
  import { formatDate } from '$lib/utils/date';
  import MarkdownRenderer from './ui/MarkdownRenderer.svelte';
  import Card from './ui/Card.svelte';
  import Button from './ui/Button.svelte';
  import PageHeader from './ui/PageHeader.svelte';
  

import ArrowLeft from "lucide-svelte/icons/arrow-left";
import Scale from "lucide-svelte/icons/scale";
import FileText from "lucide-svelte/icons/file-text";
import Folder from "lucide-svelte/icons/folder";
import AlertCircle from "lucide-svelte/icons/alert-circle";



  const repo = getContext<RepoContext>('repo');
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
        <div class="empty-icon-ring">
          <AlertCircle size={40} />
        </div>
        <h2>No documentation found</h2>
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
                  {#if isADR(currentFile.name)}
                    <Scale size={32} />
                  {:else}
                    <FileText size={32} />
                  {/if}
                  {currentFile.name}
              </h2>
              <Button variant="secondary" size="sm" icon={ArrowLeft} onclick={goBack}>Back to List</Button>
            </div>
            {#if currentFile.last_commit}
                <div class="last-update">
                    Last updated {formatDate(currentFile.last_commit.author.when)} by {currentFile.last_commit.author.name}
                </div>
            {/if}
        </div>
        
        <Card class="markdown-body-card">
            <MarkdownRenderer content={currentFile.content} />
        </Card>
      </div>
    {:else}
      <div class="flat-list-container">
        {#each entries.filter(e => e.type === 'dir' || e.name.endsWith('.md')) as entry}
          <button class="list-item-row list-item-clickable" onclick={() => loadPath(`${currentPath}/${entry.name}`)}>
            <div class="list-item-icon">
              {#if entry.type === 'dir'}
                <Folder size={20} style="color: var(--primary-color);" />
              {:else if isADR(entry.name)}
                <Scale size={20} style="color: var(--text-muted); opacity: 0.7;" />
              {:else}
                <FileText size={20} style="color: var(--text-muted); opacity: 0.7;" />
              {/if}
            </div>
            <div class="item-main-content">
              <div class="item-title-bar">
                <span class="item-title-link">{entry.name}</span>
              </div>
              {#if entry.last_commit}
                <div class="item-meta-data">{entry.last_commit.summary}</div>
              {/if}
            </div>
            <div class="item-actions-right">
              <span class="material-symbols-rounded" style="color: var(--text-muted); opacity: 0.5; font-size: 1.2rem;">chevron_right</span>
            </div>
          </button>
        {/each}
      </div>
    {/if}
  </div>
</div>

<style>
  .repo-docs-container { display: flex; flex-direction: column; gap: 0; }
  .docs-layout { min-height: 400px; padding: 0; }
  
  .header-main-row { display: flex; justify-content: space-between; align-items: center; width: 100%; }
  .docs-layout { min-height: 400px; }

  .loading-state { display: flex; flex-direction: column; align-items: center; justify-content: center; height: 300px; gap: 1rem; color: var(--text-muted); }
  .spinner { width: 32px; height: 32px; border: 3px solid var(--md-sys-color-outline-variant); border-top-color: var(--primary-color); border-radius: 50%; animation: spin 1s infinite linear; }
  @keyframes spin { to { transform: rotate(360deg); } }

  .hint { margin-top: 1rem; font-size: 0.85rem; opacity: 0.7; }

  .empty-state {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    padding: 6rem 2rem;
    text-align: center;
    gap: 1.25rem;
    color: var(--text-muted);
  }

  .empty-state h2 {
    margin: 0;
    font-size: 1.75rem;
    font-weight: 950;
    color: var(--text-color);
  }

  .empty-state p {
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

  .list-item-clickable {
    width: 100%;
    background: transparent;
    border: none;
    border-bottom: 1px solid var(--md-sys-color-outline-variant);
    cursor: pointer;
    text-align: left;
    font-family: inherit;
  }

  .file-view { display: flex; flex-direction: column; gap: 1.5rem; animation: fadeIn 0.3s ease-out; }
  @keyframes fadeIn { from { opacity: 0; transform: translateY(10px); } to { opacity: 1; transform: translateY(0); } }

  .file-header { display: flex; flex-direction: column; gap: 0.25rem; }
  .file-name { margin: 0; display: flex; align-items: center; gap: 0.75rem; font-size: 1.5rem; font-weight: 900; }
  .last-update { font-size: 0.85rem; color: var(--text-muted); font-weight: 600; }

  :global(.markdown-body-card) { padding: 2.5rem !important; border-radius: 2rem !important; background: var(--bg-color) !important; border: 1px solid var(--md-sys-color-outline-variant) !important; }
</style>
