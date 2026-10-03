<script lang="ts">
  import type { RepoContext } from "$lib/types";
  import { apiFetch } from "$lib/api";
  import { onMount, getContext } from 'svelte';
  import { goto } from '$app/navigation';
  import { formatDate } from '$lib/utils/date';
  import { dialog } from '$lib/dialog.svelte';
  import Button from './ui/Button.svelte';
  import History from "lucide-svelte/icons/history";
  import ArrowLeft from "lucide-svelte/icons/arrow-left";
  import MoreVertical from "lucide-svelte/icons/more-vertical";
  import FileText from "lucide-svelte/icons/file-text";
  import X from "lucide-svelte/icons/x";
  import WikiSidebar from "./WikiSidebar.svelte";
  import WikiPageViewer from "./WikiPageViewer.svelte";

  let { path = 'Home' } = $props();
  const repo = getContext<RepoContext>('repo');

  let wikiPages = $state<any[]>([]);
  let currentPage = $state<any>(null);
  let history = $state<any[]>([]);
  let loading = $state(true);
  let saving = $state(false);
  let error = $state('');

  // The wiki is read only: the backend answers 501 to every write because the
  // pages live in the repository as docs, not in the database.
  let mode = $state<'view' | 'history'>('view');
  let showDocsHint = $state(false);
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
        // ListWikiPages responde un array plano, no un objeto con "pages". Leyendo
        // data.pages la wiki se quedaba siempre vacia.
        wikiPages = sortTree(Array.isArray(data) ? data : data?.pages || []);
      }
    } catch (e) {
      console.error(e);
    } finally {
      loading = false;
    }
  }

  // The wiki is served from the repository as docs, so a page is a file in the
  // tree. GetWikiPage answers {name, path, content}; without this the viewer was
  // always handed a null page and no content ever appeared.
  async function loadPage(pagePath: string) {
    if (!repo.data?.owner?.username || !repo.data?.name) return;

    currentPage = null;
    const res = await apiFetch(
      `/api/v2/${repo.data.owner.username}/${repo.data.name}/wiki/${encodeURIComponent(pagePath)}`,
    );

    if (!res.ok) {
      // A page that is not in the tree is not an error worth a toast: the file
      // simply does not exist yet, and the viewer has an empty state for it.
      currentPage = null;
      return;
    }

    const data = await res.json();
    currentPage = {
      name: data.name || pagePath,
      title: data.name || pagePath,
      path: data.path || pagePath,
      content: data.content || '',
    };
  }

  async function fetchHistory() {
    try {
      // GetWikiPage distingue el historial por el prefijo de ruta "history/".
      // Con ?page= el parametro llega como "history" a secas y el backend buscaba
      // una pagina llamada "history", con un 404.
      const res = await apiFetch(`/api/v2/${repo.data?.owner?.username}/${repo.data?.name}/wiki/history/${encodeURIComponent(path)}`);
      if (res.ok) history = await res.json();
    } catch (e) {
      console.error(e);
    }
  }

  // Both entry points that used to open the editor now explain why there is none.
  function explainDocsAsCode() {
    showDocsHint = true;
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

  // Navigating between pages swaps the path prop without remounting, so the
  // content has to follow it.
  $effect(() => {
    if (!repo.data || !path) return;
    loadPage(path);
  });
</script>

<div class="wiki-container">
  <div class="wiki-grid" style={wikiPages.length === 0 ? 'display: block;' : ''}>
    {#if wikiPages.length > 0}
      <WikiSidebar {wikiPages} currentPath={path} {expandedFolders} repoOwner={repo.data?.owner?.username || ''} repoName={repo.data?.name || ''}
        onNewPage={explainDocsAsCode} onToggleFolder={toggleFolder} />
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
              <h1>{currentPage?.title || currentPage?.name || 'Wiki'}</h1>
            </div>

            <div class="header-right">
              {#if mode === 'view' || mode === 'history'}
                <div class="tab-group">
                  <button class="tab-btn" class:active={mode === 'view'} onclick={() => mode = 'view'}>Página</button>
                  <button class="tab-btn" class:active={mode === 'history'} onclick={toggleHistory}>Historial</button>
                </div>
              {/if}

              {#if mode === 'view' && currentPage}
                <div class="action-group">
                  <Button variant="ghost" onclick={explainDocsAsCode} icon={FileText}>
                    Editar en el repositorio
                  </Button>
                </div>
              {/if}
            </div>
          </header>
        {/if}

        {#if showDocsHint}
          <div class="docs-hint">
            <FileText size={18} />
            <div>
              <strong>Esta wiki se edita en el repositorio</strong>
              <p>
                Las páginas son ficheros Markdown dentro de la carpeta
                <code>{repo.data?.wiki_folder_name || 'docs'}</code>. El servidor
                rechaza cualquier escritura aqui a proposito, asi que los cambios
                llegan haciendo push y se ven en cuanto se fusionan.
              </p>
            </div>
            <button class="docs-hint-close" onclick={() => showDocsHint = false} aria-label="Cerrar aviso">
              <X size={16} />
            </button>
          </div>
        {/if}

        <div class="wiki-body" style={wikiPages.length === 0 && mode === 'view' ? 'padding: 0; display: flex; align-items: center; justify-content: center; min-height: 60vh;' : ''}>
          {#if mode === 'view'}
            <WikiPageViewer page={currentPage || null}
              repoOwner={repo.data?.owner?.username || ''} repoName={repo.data?.name || ''}
              onCreatePage={explainDocsAsCode} />
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
  .docs-hint {
    display: flex;
    gap: 0.85rem;
    align-items: flex-start;
    margin-bottom: 1rem;
    padding: 0.9rem 1.1rem;
    border: 1px solid var(--status-info-muted);
    border-radius: var(--radius, 1rem);
    background: var(--status-info-muted);
    color: var(--md-sys-color-on-surface);
  }
  .docs-hint strong { display: block; margin-bottom: 0.2rem; }
  .docs-hint p { margin: 0; font-size: 0.9rem; color: var(--md-sys-color-on-surface-variant); }
  .docs-hint code {
    font-family: var(--font-mono);
    background: var(--output-surface-raised);
    padding: 0.1rem 0.3rem;
    border-radius: 0.3rem;
  }
  .docs-hint-close {
    margin-left: auto;
    display: grid;
    place-items: center;
    padding: 0.25rem;
    border: none;
    border-radius: 0.4rem;
    background: transparent;
    color: inherit;
    cursor: pointer;
  }
  .docs-hint-close:hover { background: var(--status-info-muted); }

  .wiki-container { padding: 1.5rem; }
  .wiki-grid { display: flex; gap: 2rem; align-items: flex-start; }
  .wiki-main { flex: 1; min-width: 0; }
  .wiki-content-header {
    display: flex; justify-content: space-between; align-items: center;
    flex-wrap: wrap; gap: 1rem; margin-bottom: 1.5rem;
  }
  .header-left { flex: 1; min-width: 200px; }
  .header-left h1 { margin: 0; font-size: 1.75rem; font-weight: 800; }
  .header-right { display: flex; align-items: center; gap: 0.75rem; }
  .tab-group { display: flex; background: var(--md-sys-color-surface-container-lowest); border-radius: 0.5rem; padding: 2px; }
  .tab-btn {
    padding: 0.4rem 1rem; border: none; background: transparent;
    color: var(--text-muted); cursor: pointer; border-radius: 0.4rem;
    font-weight: 700; font-size: 0.85rem; transition: all 0.2s;
  }
  .tab-btn.active { background: var(--md-sys-color-surface-container-high); color: var(--text-color); }
  .action-group { display: flex; gap: 0.5rem; align-items: center; }
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
