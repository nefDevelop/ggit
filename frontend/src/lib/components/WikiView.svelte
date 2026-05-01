<script lang="ts">
  import { apiFetch } from "$lib/api";
  import { onMount, getContext } from 'svelte';
  import { goto } from '$app/navigation';
  import { dialog } from '$lib/dialog.svelte';
  import Button from './ui/Button.svelte';
  import { 
    FileText as Description, 
    Edit, 
    History as HistoryIcon, 
    Plus as Add, 
    Save, 
    X as Close, 
    Trash2 as Delete,
    ArrowLeft as ArrowBack,
    MoreVertical as MoreVert,
    Folder,
    ChevronRight
  } from 'lucide-svelte';
  import { slide } from 'svelte/transition';
  import MarkdownRenderer from './ui/MarkdownRenderer.svelte';
  
  let { path = 'Home' } = $props();
  const repo: any = getContext('repo');
  
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

  function parseFrontmatterOrder(content: string): number | null {
    const match = content.match(/^---\n[\s\S]*?order:\s*(\d+)\s*\n[\s\S]*?\n---/);
    if (match) return parseInt(match[1]);
    return null;
  }

  function extractFoldersFromSlantedPaths(items: any[]): string[] {
    const foldersSet = new Set<string>();
    for (const item of items) {
      if (item.is_folder) {
        foldersSet.add(item.name);
      } else if (item.name && item.name.includes('/')) {
        const parts = item.name.split('/');
        parts.pop();
        let cumulative = '';
        for (const part of parts) {
          cumulative = cumulative ? `${cumulative}/${part}` : part;
          foldersSet.add(cumulative);
        }
      }
    }
    return Array.from(foldersSet).sort();
  }

  async function fetchWiki() {
    loading = true;
    error = '';
    try {
      const listRes = await apiFetch(`/api/v2/${repo.data.owner.username}/${repo.data.name}/wiki`);
      if (listRes.ok) wikiPages = (await listRes.json()) || [];
      
      const pageRes = await apiFetch(`/api/v2/${repo.data.owner.username}/${repo.data.name}/wiki/${path}`);
      if (pageRes.ok) {
        currentPage = await pageRes.json();
        editContent = currentPage.content;
        editTitle = currentPage.name;
      } else if (wikiPages.length === 0) {
        currentPage = null;
        editContent = '';
        editTitle = '';
      } else {
        error = 'Wiki page not found.';
      }
    } catch (e) {
      error = 'Failed to load wiki.';
    } finally {
      loading = false;
    }
  }

  async function fetchHistory() {
    try {
      const res = await apiFetch(`/api/v2/${repo.data.owner.username}/${repo.data.name}/wiki/history/${path}`);
      if (res.ok) history = await res.json();
    } catch (e) {
      console.error('Failed to load history');
    }
  }

  async function handleSave() {
    saving = true;
    try {
      const finalTitle = parentPath ? `${parentPath}/${editTitle}` : editTitle;
      let finalContent = editContent;
      if (pageOrder !== null && pageOrder !== undefined && !isNaN(pageOrder)) {
        // Strip out existing frontmatter if any
        const cleaned = finalContent.replace(/^---\n[\s\S]*?\n---\n/, '');
        finalContent = `---\norder: ${pageOrder}\n---\n${cleaned}`;
      }

      const endpoint = mode === 'edit' && currentPage?.id ? 
        `/api/v2/${repo.data.owner.username}/${repo.data.name}/wiki/${path}` :
        `/api/v2/${repo.data.owner.username}/${repo.data.name}/wiki`;

      const res = await apiFetch(endpoint, {
        method: mode === 'edit' && currentPage?.id ? 'PUT' : 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          title: finalTitle,
          content: finalContent,
          message: commitMessage || `Update ${finalTitle}`
        })
      });

      if (res.ok) {
        const data = await res.json();
        goto(`/${repo.data.owner.username}/${repo.data.name}/wiki/${data.url || finalTitle}`);
        mode = 'view';
      } else {
        const data = await res.json();
        error = data.message || 'Error saving page';
      }
    } catch (e) {
      error = 'Connection error';
    } finally {
      saving = false;
    }
  }

  async function handleDelete() {
    if (!await dialog.confirm('Are you sure you want to delete this page?')) return;
    try {
      const res = await apiFetch(`/api/v2/${repo.data.owner.username}/${repo.data.name}/wiki/${path}`, {
        method: 'DELETE'
      });
      if (res.ok) {
        goto(`/${repo.data.owner.username}/${repo.data.name}/wiki`);
      }
    } catch (e) {
      error = 'Failed to delete page';
    }
  }

  function startEdit() {
    mode = 'edit';
    const full = currentPage?.name || '';
    if (full.includes('/')) {
      const parts = full.split('/');
      editTitle = parts.pop() || '';
      parentPath = parts.join('/');
    } else {
      editTitle = full;
      parentPath = '';
    }
    editContent = currentPage?.content || '';
    pageOrder = parseFrontmatterOrder(editContent);
  }

  function startNew() {
    mode = 'edit';
    currentPage = null;
    editTitle = '';
    editContent = '';
    parentPath = '';
    pageOrder = null;
  }

  let expandedFolders = $state<Record<string, boolean>>({});

  function toggleFolder(folderPath: string) {
    expandedFolders[folderPath] = !expandedFolders[folderPath];
  }

  $effect(() => {
    if (repo.data && path) {
      fetchWiki();
    }
  });

  $effect(() => {
    if (mode === 'history' && path) {
        fetchHistory();
    }
  });

  const formatDate = (ts: number) => new Date(ts * 1000).toLocaleDateString('es-ES', {
    day: 'numeric', month: 'short', year: 'numeric'
  });
</script>

<div class="wiki-container">
  <div class="wiki-grid" style={wikiPages.length === 0 ? 'display: block;' : ''}>
    <!-- Sidebar -->
    {#if wikiPages.length > 0}
    <aside class="wiki-sidebar">
      <div class="sidebar-header-minimal">
        <Button class="w-full" onclick={startNew} icon="add">
          New Page
        </Button>
      </div>
      
      <nav class="wiki-nav">
        {@render renderTree(wikiPages)}
      </nav>

{#snippet renderTree(items: any[], depth = 0, parentPath = '')}
    {#each items as pageItem}
      {@const fullPath = parentPath ? `${parentPath}/${pageItem.name}` : pageItem.name}
      {#if pageItem.is_folder}
        <div class="folder-section">
          <button 
            class="nav-item folder-label" 
            style="padding-left: {depth * 0.75 + 0.8}rem"
            onclick={() => toggleFolder(fullPath)}
          >
            <ChevronRight size={14} class="chevron {expandedFolders[fullPath] ? 'open' : ''}" />
            <Folder size={14} class="opacity-60" />
            <span>{pageItem.title}</span>
          </button>
          {#if expandedFolders[fullPath]}
            <div class="folder-children" transition:slide>
                {@render renderTree(pageItem.children || [], depth + 1, fullPath)}
            </div>
          {/if}
        </div>
      {:else}
        <a 
          href="/{repo.data.owner.username}/{repo.data.name}/wiki/{pageItem.name}" 
          class="nav-item"
          class:active={path === pageItem.name}
          style="padding-left: {depth * 0.75 + 0.8}rem"
        >
          <Description size={14} class="opacity-60" />
          <span>{pageItem.title || pageItem.name}</span>
        </a>
      {/if}
    {/each}
{/snippet}
    </aside>
    {/if}

    <!-- Content Area -->
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
                <button class="tab-btn" class:active={mode === 'view'} onclick={() => mode = 'view'}>
                    Página
                </button>
                <button class="tab-btn" class:active={mode === 'history'} onclick={() => mode = 'history'}>
                    Historial
                </button>
            </div>
            {/if}
            
            <div class="action-group">
                {#if mode === 'view'}
                    <button class="btn-icon-text" onclick={startEdit}>
                        <Edit size={16} />
                        Editar
                    </button>
                {:else if mode === 'edit'}
                    <button class="btn-icon-text variant" onclick={() => { if (wikiPages.length === 0) { mode = 'view'; } else { mode = 'view'; } }}>
                        <Close size={16} />
                        Cancelar
                    </button>
                    <button class="btn-primary" onclick={handleSave} disabled={saving || !editTitle}>
                        <Save size={16} />
                        {saving ? 'Guardando...' : 'Guardar'}
                    </button>
                {/if}
                
                {#if mode === 'view' && path !== 'Home' && wikiPages.length > 0}
                    <button class="btn-danger-icon" onclick={handleDelete} title="Eliminar Página">
                        <Delete size={18} />
                    </button>
                {/if}
            </div>
          </div>
        </header>
        {/if}

        <div class="wiki-body" style={wikiPages.length === 0 && mode === 'view' ? 'padding: 0; display: flex; align-items: center; justify-content: center; min-height: 60vh;' : ''}>
          {#if mode === 'view'}
              {#if currentPage}
                <MarkdownRenderer
                  content={currentPage?.content}
                  repoContext={{owner: repo?.data?.owner?.username || repo?.owner?.username || repo?.owner?.name || '', name: repo?.data?.name || repo?.name || ''}}
                />
              {:else}
                <div class="empty-wiki-state">
                  <span class="material-symbols-rounded" style="font-size: 4rem; opacity: 0.5;">menu_book</span>
                  <h2 style="margin: 0.5rem 0;">Wiki Vacía</h2>
                  <p style="margin-bottom: 1.5rem;">Este repositorio aún no tiene documentación.</p>
                  <Button variant="primary" onclick={startNew} icon="add">Crear Primera Página</Button>
                </div>
              {/if}
          {:else if mode === 'edit'}
            {@const existingFolders = extractFoldersFromSlantedPaths(wikiPages)}
            <div class="edit-area">
                <div class="wiki-metadata-row">
                  <div class="meta-field">
                    <label for="parent-folder">Ubicación / Carpeta Padre</label>
                    <select id="parent-folder" bind:value={parentPath}>
                      <option value="">(Raíz / Root)</option>
                      {#each existingFolders as folder}
                        <option value={folder}>{folder}</option>
                      {/each}
                    </select>
                  </div>
                  <div class="meta-field">
                    <label for="page-order">Orden de la página</label>
                    <input type="number" id="page-order" bind:value={pageOrder} placeholder="Ej: 1, 2..." />
                  </div>
                </div>

                <textarea 
                    class="wiki-editor" 
                    bind:value={editContent} 
                    placeholder="Escribe el contenido en Markdown..."
                ></textarea>
                <div class="commit-msg-box">
                    <label for="commit">Mensaje del commit (opcional)</label>
                    <input type="text" id="commit" bind:value={commitMessage} placeholder="¿Qué has cambiado?" />
                </div>
            </div>
          {:else if mode === 'history'}
            <div class="history-list">
                {#each history as entry}
                    <div class="history-item">
                        <div class="history-node"></div>
                        <div class="history-info">
                            <div class="history-msg">{entry.message}</div>
                            <div class="history-meta">
                                <b>{entry.author}</b> • {formatDate(entry.created)} • <code class="hash">{entry.id.substring(0, 8)}</code>
                            </div>
                        </div>
                    </div>
                {/each}
                {#if history.length === 0}
                    <div class="empty-history">Cargando historial...</div>
                {/if}
            </div>
          {/if}
        </div>
      {/if}
    </main>
  </div>
</div>

<style>
  .wiki-container {
    padding: 1rem;
    max-width: 1400px;
    margin: 0 auto;
    animation: fadeIn 0.4s ease-out;
  }

  .wiki-grid {
    display: grid;
    grid-template-columns: 280px 1fr;
    gap: 1.5rem;
    align-items: start;
  }
  @media (max-width: 900px) {
    .wiki-grid {
      grid-template-columns: 1fr;
    }
  }

  /* Sidebar */
  .wiki-sidebar {
    padding: 1.5rem;
    position: sticky;
    top: 1rem;
  }

  .sidebar-header-minimal {
    margin-bottom: 1.5rem;
  }

  .wiki-nav {
    display: flex;
    flex-direction: column;
    gap: 0.25rem;
  }

  .nav-item {
    display: flex;
    align-items: center;
    gap: 0.75rem;
    padding: 0.6rem 0.8rem;
    border-radius: 12px;
    text-decoration: none;
    color: var(--md-sys-color-on-surface-variant);
    font-size: 0.95rem;
    transition: all 0.2s;
  }

  .nav-item:hover {
    background: var(--md-sys-color-surface-variant);
    color: var(--md-sys-color-on-surface);
  }

  .nav-item.active {
    background: var(--md-sys-color-primary-container);
    color: var(--md-sys-color-primary);
    font-weight: 600;
  }

  .folder-label {
    cursor: pointer;
    opacity: 0.8;
    background: transparent !important;
    font-size: 0.8rem;
    font-weight: 800;
    text-transform: uppercase;
    letter-spacing: 0.05em;
    margin-top: 0.5rem;
    width: 100%;
    border: none;
    text-align: left;
  }

  :global(.chevron) {
    transition: transform 0.2s;
  }

  :global(.chevron.open) {
    transform: rotate(90deg);
  }

  .folder-children {
    display: flex;
    flex-direction: column;
    gap: 0.1rem;
    border-left: 1px solid var(--md-sys-color-outline-variant);
    margin-left: 0.2rem;
  }

  :global(.opacity-60) {
    opacity: 0.6;
  }

  /* Main Content */
  .wiki-main {
    min-height: 80vh;
    display: flex;
    flex-direction: column;
  }

  .wiki-content-header {
    padding: 1.5rem 2rem;
    border-bottom: 1px solid var(--md-sys-color-outline-variant);
    display: flex;
    justify-content: space-between;
    align-items: center;
    flex-wrap: wrap;
    gap: 1rem;
  }

  .header-left h1 {
    font-size: 1.75rem;
    margin: 0;
    font-weight: 700;
  }

  .title-input {
    font-size: 1.75rem;
    font-weight: 700;
    background: transparent;
    border: none;
    border-bottom: 2px solid var(--md-sys-color-primary);
    color: var(--md-sys-color-on-surface);
    width: 100%;
    outline: none;
  }

  .header-right {
    display: flex;
    align-items: center;
    gap: 2rem;
    max-width: 100%;
  }

  .tab-group {
    background: var(--md-sys-color-surface-variant);
    padding: 0.25rem;
    border-radius: 100px;
    display: flex;
    overflow-x: auto;
    scrollbar-width: none;
  }
  .tab-group::-webkit-scrollbar { display: none; }

  .tab-btn {
    padding: 0.4rem 1.25rem;
    border-radius: 100px;
    border: none;
    background: transparent;
    color: var(--md-sys-color-on-surface-variant);
    font-size: 0.9rem;
    font-weight: 600;
    cursor: pointer;
    transition: all 0.2s;
  }

  .tab-btn.active {
    background: var(--md-sys-color-surface-container);
    color: var(--md-sys-color-primary);
    box-shadow: var(--shadow-sm);
  }

  .action-group {
    display: flex;
    gap: 0.75rem;
    align-items: center;
  }

  .btn-icon-text {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    background: var(--md-sys-color-surface-variant);
    color: var(--md-sys-color-on-surface);
    border: none;
    padding: 0.5rem 1rem;
    border-radius: 12px;
    font-weight: 600;
    cursor: pointer;
    transition: all 0.2s;
  }

  .btn-icon-text.variant {
    background: transparent;
    color: var(--md-sys-color-on-surface-variant);
  }

  .btn-primary {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    background: var(--md-sys-color-primary);
    color: white;
    border: none;
    padding: 0.5rem 1.25rem;
    border-radius: 12px;
    font-weight: 600;
    cursor: pointer;
    transition: all 0.2s;
  }

  .btn-primary:hover:not(:disabled) {
    background: var(--primary-hover);
    transform: translateY(-1px);
  }

  .btn-danger-icon {
    color: var(--md-sys-color-error);
    background: transparent;
    border: none;
    cursor: pointer;
    padding: 0.5rem;
    border-radius: 10px;
    transition: background 0.2s;
  }

  .btn-danger-icon:hover {
    background: color-mix(in srgb, var(--md-sys-color-error) 10%, transparent);
  }

  .wiki-body {
    padding: 2rem;
    flex: 1;
  }

  .empty-wiki-state {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    height: 100%;
    min-height: 400px;
    text-align: center;
    color: var(--text-muted);
  }

  /* Edit Area */
  .edit-area {
    display: flex;
    flex-direction: column;
    gap: 1.5rem;
    height: 100%;
  }

  .wiki-metadata-row {
    display: flex;
    gap: 1rem;
    flex-wrap: wrap;
  }

  .meta-field {
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
    flex: 1;
    min-width: 200px;
  }

  .meta-field label {
    font-size: 0.85rem;
    font-weight: 600;
    color: var(--md-sys-color-on-surface-variant);
  }

  .meta-field select,
  .meta-field input {
    background: var(--md-sys-color-surface-variant);
    border: 1px solid var(--md-sys-color-outline-variant);
    border-radius: 12px;
    padding: 0.75rem 1rem;
    color: var(--md-sys-color-on-surface);
    font-size: 0.95rem;
    outline: none;
    transition: all 0.2s;
  }

  .meta-field select:focus,
  .meta-field input:focus {
    border-color: var(--md-sys-color-primary);
    background: rgba(var(--primary-rgb), 0.05);
  }

  .wiki-editor {
    flex: 1;
    min-height: 400px;
    background: var(--md-sys-color-surface-variant);
    border: 1px solid var(--md-sys-color-outline-variant);
    border-radius: 16px;
    padding: 1.5rem;
    font-family: 'JetBrains Mono', monospace;
    font-size: 1rem;
    color: var(--md-sys-color-on-surface);
    resize: vertical;
    outline: none;
  }

  .commit-msg-box {
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
  }

  .commit-msg-box label {
    font-size: 0.85rem;
    font-weight: 600;
    color: var(--md-sys-color-on-surface-variant);
  }

  .commit-msg-box input {
    background: var(--md-sys-color-surface-variant);
    border: 1px solid var(--md-sys-color-outline-variant);
    border-radius: 10px;
    padding: 0.75rem 1rem;
    color: var(--md-sys-color-on-surface);
  }

  /* History List */
  .history-list {
    display: flex;
    flex-direction: column;
    padding-left: 1rem;
    border-left: 2px solid var(--md-sys-color-outline-variant);
  }

  .history-item {
    position: relative;
    padding-bottom: 2rem;
    padding-left: 1.5rem;
  }

  .history-node {
    position: absolute;
    left: -27px;
    top: 5px;
    width: 12px;
    height: 12px;
    background: var(--md-sys-color-primary);
    border: 3px solid var(--md-sys-color-surface-container);
    border-radius: 50%;
  }

  .history-msg {
    font-weight: 600;
    color: var(--md-sys-color-on-surface);
    margin-bottom: 0.25rem;
  }

  .history-meta {
    font-size: 0.85rem;
    color: var(--md-sys-color-on-surface-variant);
  }

  .hash {
    background: var(--md-sys-color-surface-variant);
    padding: 0.1rem 0.4rem;
    border-radius: 4px;
    font-family: monospace;
  }

  /* Skeletons */
  .skeleton-header {
    height: 40px;
    width: 300px;
    background: var(--md-sys-color-surface-variant);
    border-radius: 8px;
    margin-bottom: 2rem;
    animation: pulse 1.5s infinite;
  }

  .skeleton-body {
    height: 300px;
    background: var(--md-sys-color-surface-variant);
    border-radius: 12px;
    animation: pulse 1.5s infinite;
  }

  @keyframes pulse {
    0%, 100% { opacity: 0.5; }
    50% { opacity: 0.8; }
  }

  /* Code Block Styling - Integrated with Ggit Theme */
  :global(.markdown-body pre) {
    background-color: var(--md-sys-color-surface-container-low, #f0f2f5);
    padding: 1.5rem;
    border-radius: 16px;
    overflow-x: auto;
    margin: 1.5rem 0;
    border: 1px solid var(--md-sys-color-outline-variant, rgba(0,0,0,0.05));
    font-family: 'JetBrains Mono', 'Fira Code', ui-monospace, monospace;
    font-size: 0.85rem;
    line-height: 1.6;
    box-shadow: inset 0 1px 4px rgba(0,0,0,0.02);
  }

  :global(.markdown-body code.hljs) {
    padding: 0;
    background: transparent;
    color: var(--md-sys-color-on-surface-variant);
  }

  /* Customizing Highlight.js colors to match Ggit's vibe */
  :global(.hljs-keyword, .hljs-selector-tag) { color: var(--md-sys-color-primary, #6750a4); font-weight: 600; }
  :global(.hljs-string, .hljs-doctag) { color: #1e8e3e; } /* Nature Green */
  :global(.hljs-title, .hljs-section, .hljs-selector-id) { color: var(--md-sys-color-tertiary, #7d5260); font-weight: bold; }
  :global(.hljs-comment) { color: var(--md-sys-color-outline, #79747e); font-style: italic; }
  :global(.hljs-number, .hljs-literal) { color: #b05c00; }
  :global(.hljs-type, .hljs-class) { color: var(--md-sys-color-secondary, #625b71); }
  :global(.hljs-attr, .hljs-variable) { color: #8e24aa; }

  @keyframes fadeIn {
    from { opacity: 0; transform: translateY(10px); }
    to { opacity: 1; transform: translateY(0); }
  }

  /* MDX Callouts */
  :global(.markdown-body Warning), 
  :global(.markdown-body Note), 
  :global(.markdown-body Tip), 
  :global(.markdown-body Important), 
  :global(.markdown-body Caution),
  :global(.markdown-body info),
  :global(.markdown-body notice) {
    display: block;
    padding: 1rem 1.25rem;
    margin: 1.5rem 0;
    border-radius: 12px;
    border-left: 6px solid;
    font-size: 0.95rem;
    line-height: 1.6;
    box-shadow: 0 4px 12px rgba(0,0,0,0.05);
  }

  :global(.markdown-body Warning) {
    background: #fff8e1;
    border-color: #ffc107;
    color: #856404;
  }
  :global(.markdown-body Warning::before) { content: "⚠️ Warning"; font-weight: 800; display: block; margin-bottom: 0.5rem; text-transform: uppercase; font-size: 0.75rem; }

  :global(.markdown-body Note), :global(.markdown-body info) {
    background: #e3f2fd;
    border-color: #2196f3;
    color: #0c5460;
  }
  :global(.markdown-body Note::before) { content: "📝 Note"; font-weight: 800; display: block; margin-bottom: 0.5rem; text-transform: uppercase; font-size: 0.75rem; }

  :global(.markdown-body Tip), :global(.markdown-body notice) {
    background: #e8f5e9;
    border-color: #4caf50;
    color: #155724;
  }
  :global(.markdown-body Tip::before) { content: "💡 Tip"; font-weight: 800; display: block; margin-bottom: 0.5rem; text-transform: uppercase; font-size: 0.75rem; }

  :global(.markdown-body Important) {
    background: #f3e5f5;
    border-color: #9c27b0;
    color: #4a148c;
  }
  :global(.markdown-body Important::before) { content: "🚀 Important"; font-weight: 800; display: block; margin-bottom: 0.5rem; text-transform: uppercase; font-size: 0.75rem; }

  :global(.markdown-body Caution) {
    background: #ffebee;
    border-color: #f44336;
    color: #721c24;
  }
  :global(.markdown-body Caution::before) { content: "🛑 Caution"; font-weight: 800; display: block; margin-bottom: 0.5rem; text-transform: uppercase; font-size: 0.75rem; }

  @media (max-width: 900px) {
    .wiki-grid {
      grid-template-columns: 1fr;
    }
    
    .wiki-sidebar {
      position: relative;
      top: 0;
    }

    .header-right {
        width: 100%;
        justify-content: space-between;
    }
  }
</style>
