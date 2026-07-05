<script lang="ts">
  import type { RepoContext } from "$lib/types";
  import { apiFetch } from "$lib/api";
  import { onMount, getContext } from 'svelte';
  import { goto } from '$app/navigation';
  import { formatDate } from '$lib/utils/date';
  import { dialog } from '$lib/dialog.svelte';
  import Button from './ui/Button.svelte';
  import Edit from "lucide-svelte/icons/edit";
  import History from "lucide-svelte/icons/history";
  import Save from "lucide-svelte/icons/save";
  import X from "lucide-svelte/icons/x";
  import Trash2 from "lucide-svelte/icons/trash-2";
  import ArrowLeft from "lucide-svelte/icons/arrow-left";
  import MoreVertical from "lucide-svelte/icons/more-vertical";
  import WikiSidebar from "./WikiSidebar.svelte";
  import WikiEditor from "./WikiEditor.svelte";
  import WikiPageViewer from "./WikiPageViewer.svelte";

  let { path = 'Home' } = $props();
  const repo = getContext<RepoContext>('repo');

  let wikiPages = $state<any[]>([]);
  let currentPage = $state<any>(null);
  let history = $state<any[]>([]);
  let loading = $state(true);
  let saving = $state(false);
  let error = $state('');

  let mode = $state<'view' | 'edit' | 'history'>('view');
  let editContent = $state('');
  let editTitle = $state('');
  let commitMessage = $state('');
  let parentPath = $state('');
  let pageOrder = $state<number | null>(null);
  let expandedFolders = $state<Record<string, boolean>>({});

  function parseFrontmatterOrder(content: string): number | null {
    if (!content) return null;
    const match = content.match(/^---\n[\s\S]*?order:\s*(\d+)\s*\n[\s\S]*?\n---/);
    if (match) return parseInt(match[1]);
    return null;
  }

  function cleanContentForPreview(content: string): string {
    if (!content) return '';
    return content.replace(/^---[\s\S]*?\n---\n*/, '');
  }

  function sortTree(items: any[]): any[] {
    return [...items].sort((a, b) => {
      if (a.is_folder && !b.is_folder) return -1;
      if (!a.is_folder && b.is_folder) return 1;
      const aOrder = parseFrontmatterOrder(a.content) ?? 999;
      const bOrder = parseFrontmatterOrder(b.content) ?? 999;
      if (aOrder !== bOrder) return aOrder - bOrder;
      return (a.title || a.name).localeCompare(b.title || b.name);
    });
  }

  async function fetchWiki() {
    loading = true;
    try {
      const res = await apiFetch(`/api/v2/${repo.data?.owner?.username}/${repo.data?.name}/wiki`);
      if (res.ok) {
        const data = await res.json();
        wikiPages = sortTree(data?.pages || []);
        currentPage = data?.current_page || null;
        if (currentPage && currentPage.content) {
          editContent = currentPage.content;
          editTitle = currentPage.title || currentPage.name || '';
          const order = parseFrontmatterOrder(currentPage.content);
          pageOrder = order;
        }
      }
    } catch (e) {
      console.error(e);
    } finally {
      loading = false;
    }
  }

  async function fetchHistory() {
    try {
      const res = await apiFetch(`/api/v2/${repo.data?.owner?.username}/${repo.data?.name}/wiki/history?page=${path}`);
      if (res.ok) history = await res.json();
    } catch (e) {
      console.error(e);
    }
  }

  async function handleSave() {
    if (!editTitle) return;
    saving = true;
    try {
      let content = editContent || '';
      if (pageOrder !== null && pageOrder !== undefined) {
        content = content.replace(/^---[\s\S]*?\n---\n*/, '');
        content = `---\norder: ${pageOrder}\n---\n${content}`;
      }
      const res = await apiFetch(`/api/v2/${repo.data?.owner?.username}/${repo.data?.name}/wiki/${editTitle}`, {
        method: 'PUT',
        body: JSON.stringify({ content, title: editTitle, commit_message: commitMessage, parent_path: parentPath }),
      });
      if (res.ok) {
        mode = 'view';
        await fetchWiki();
      }
    } catch (e) {
      console.error(e);
    } finally {
      saving = false;
    }
  }

  async function handleDelete() {
    if (!await dialog.confirm(`¿Estás seguro de eliminar la página "${path}"?`)) return;
    try {
      const res = await apiFetch(`/api/v2/${repo.data?.owner?.username}/${repo.data?.name}/wiki/${path}`, { method: 'DELETE' });
      if (res.ok) {
        goto(`/${repo.data?.owner?.username}/${repo.data?.name}/wiki/Home`);
        await fetchWiki();
      }
    } catch (e) {
      console.error(e);
    }
  }

  function startEdit() {
    editContent = currentPage?.content || '';
    editTitle = currentPage?.title || currentPage?.name || '';
    const order = parseFrontmatterOrder(currentPage?.content || '');
    pageOrder = order;
    mode = 'edit';
  }

  function startNew() {
    editContent = '';
    editTitle = '';
    commitMessage = '';
    pageOrder = null;
    parentPath = '';
    mode = 'edit';
  }

  function toggleFolder(fullPath: string) {
    expandedFolders = { ...expandedFolders, [fullPath]: !expandedFolders[fullPath] };
  }

  function toggleHistory() {
    if (mode === 'history') {
      mode = 'view';
    } else {
      mode = 'history';
      history = [];
      fetchHistory();
    }
  }

  onMount(() => { if (repo.data) fetchWiki(); });
</script>

<div class="wiki-container">
  <div class="wiki-grid" style={wikiPages.length === 0 ? 'display: block;' : ''}>
    {#if wikiPages.length > 0}
      <WikiSidebar {wikiPages} currentPath={path} {expandedFolders} repoOwner={repo.data?.owner?.username || ''} repoName={repo.data?.name || ''}
        onNewPage={startNew} onToggleFolder={toggleFolder} />
    {/if}

    <main class="wiki-main">
      {#if loading}
        <div class="wiki-skeleton">
          <div class="skeleton-header"></div>
          <div class="skeleton-body"></div>
        </div>
      {:else}
        {#if wikiPages.length > 0 || mode !== 'view'}
          <header class="wiki-content-header">
            <div class="header-left">
              {#if mode === 'edit'}
                <input type="text" class="title-input" bind:value={editTitle} placeholder="Título de la página..." />
              {:else}
                <h1>{currentPage?.name || 'Wiki'}</h1>
              {/if}
            </div>

            <div class="header-right">
              {#if mode === 'view' || mode === 'history'}
                <div class="tab-group">
                  <button class="tab-btn" class:active={mode === 'view'} onclick={() => mode = 'view'}>Página</button>
                  <button class="tab-btn" class:active={mode === 'history'} onclick={toggleHistory}>Historial</button>
                </div>
              {/if}

              <div class="action-group">
                {#if mode === 'view'}
                  <Button variant="tonal" onclick={startEdit} icon={Edit}>Editar</Button>
                {:else if mode === 'edit'}
                  <Button variant="ghost" onclick={() => mode = 'view'} disabled={saving} icon={X}>Cancelar</Button>
                  <Button variant="primary" onclick={handleSave} disabled={saving || !editTitle} loading={saving} icon={Save}>Guardar</Button>
                {/if}
                {#if mode === 'view' && path !== 'Home' && wikiPages.length > 0}
                  <button class="btn-danger-icon" onclick={handleDelete} title="Eliminar Página" aria-label="Eliminar Página">
                    <Trash2 size={18} />
                  </button>
                {/if}
              </div>
            </div>
          </header>
        {/if}

        <div class="wiki-body" style={wikiPages.length === 0 && mode === 'view' ? 'padding: 0; display: flex; align-items: center; justify-content: center; min-height: 60vh;' : ''}>
          {#if mode === 'view'}
            <WikiPageViewer page={currentPage || null}
              repoOwner={repo.data?.owner?.username || ''} repoName={repo.data?.name || ''}
              onCreatePage={startNew} />
          {:else if mode === 'edit'}
            <WikiEditor bind:editTitle bind:editContent bind:commitMessage bind:parentPath bind:pageOrder {wikiPages} />
          {:else if mode === 'history'}
            <div class="history-list">
              {#each history as entry}
                <div class="history-item">
                  <div class="history-node"></div>
                  <div class="history-content">
                    <span class="history-msg">{entry.message || entry.title}</span>
                    <span class="history-meta">{entry.author} · {formatDate(entry.date * 1000 || entry.date)}</span>
                  </div>
                </div>
              {/each}
            </div>
          {/if}
        </div>
      {/if}
    </main>
  </div>
</div>

<style>
  .wiki-container { padding: 1.5rem; }
  .wiki-grid { display: flex; gap: 2rem; align-items: flex-start; }
  .wiki-main { flex: 1; min-width: 0; }
  .wiki-content-header {
    display: flex; justify-content: space-between; align-items: center;
    flex-wrap: wrap; gap: 1rem; margin-bottom: 1.5rem;
  }
  .header-left { flex: 1; min-width: 200px; }
  .header-left h1 { margin: 0; font-size: 1.75rem; font-weight: 800; }
  .title-input {
    width: 100%; font-size: 1.5rem; font-weight: 800;
    background: rgba(0,0,0,0.2); border: 1px solid var(--md-sys-color-outline-variant);
    color: inherit; padding: 0.5rem 1rem; border-radius: 0.5rem;
    font-family: inherit; box-sizing: border-box;
  }
  .header-right { display: flex; align-items: center; gap: 0.75rem; }
  .tab-group { display: flex; background: rgba(0,0,0,0.1); border-radius: 0.5rem; padding: 2px; }
  .tab-btn {
    padding: 0.4rem 1rem; border: none; background: transparent;
    color: var(--text-muted); cursor: pointer; border-radius: 0.4rem;
    font-weight: 700; font-size: 0.85rem; transition: all 0.2s;
  }
  .tab-btn.active { background: var(--md-sys-color-surface-container-high); color: var(--text-color); }
  .action-group { display: flex; gap: 0.5rem; align-items: center; }
  .btn-danger-icon {
    background: transparent; border: none; color: var(--md-sys-color-error);
    cursor: pointer; padding: 0.5rem; border-radius: 0.5rem;
    transition: all 0.2s; opacity: 0.7;
  }
  .btn-danger-icon:hover { opacity: 1; background: rgba(var(--error-rgb),0.1); }
  .wiki-body { min-height: 400px; }
  .wiki-skeleton { display: flex; flex-direction: column; gap: 1rem; }
  .skeleton-header { height: 40px; background: var(--md-sys-color-surface-container-high); border-radius: 0.75rem; }
  .skeleton-body { height: 300px; background: var(--md-sys-color-surface-container-high); border-radius: 0.75rem; }
  .history-list { display: flex; flex-direction: column; gap: 0.5rem; }
  .history-item { display: flex; gap: 1rem; padding: 0.75rem 0; border-left: 2px solid var(--md-sys-color-outline-variant); padding-left: 1rem; }
  .history-node { width: 10px; height: 10px; background: var(--primary-color); border-radius: 50%; margin-top: 0.3rem; }
  .history-content { display: flex; flex-direction: column; gap: 0.25rem; }
  .history-msg { font-weight: 600; }
  .history-meta { font-size: 0.8rem; color: var(--text-muted); }
</style>
