<script lang="ts">
  import { onMount, getContext } from 'svelte';
  import { 
    Description, 
    Edit, 
    History as HistoryIcon, 
    Add, 
    Save, 
    Close, 
    Delete,
    ArrowBack,
    MoreVert
  } from 'lucide-svelte';
  
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

  async function fetchWiki() {
    loading = true;
    error = '';
    try {
      const listRes = await fetch(`/api/v2/${repo.data.owner.username}/${repo.data.name}/wiki`);
      if (listRes.ok) wikiPages = await listRes.json();
      
      const pageRes = await fetch(`/api/v2/${repo.data.owner.username}/${repo.data.name}/wiki/${path}`);
      if (pageRes.ok) {
        currentPage = await pageRes.json();
        editContent = currentPage.content;
        editTitle = currentPage.name;
      } else if (path === 'Home' || path === '') {
        currentPage = { name: 'Home', content: '# Welcome to the Wiki\n\nThis is your repository\'s wiki. Start documenting!' };
        editContent = currentPage.content;
        editTitle = 'Home';
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
      const res = await fetch(`/api/v2/${repo.data.owner.username}/${repo.data.name}/wiki/${path}/history`);
      if (res.ok) history = await res.json();
    } catch (e) {
      console.error('Failed to load history');
    }
  }

  async function handleSave() {
    saving = true;
    try {
      const method = mode === 'edit' && currentPage?.id ? 'PUT' : 'POST';
      const endpoint = mode === 'edit' ? 
        `/api/v2/${repo.data.owner.username}/${repo.data.name}/wiki/${path}` :
        `/api/v2/${repo.data.owner.username}/${repo.data.name}/wiki`;

      const res = await fetch(endpoint, {
        method: mode === 'edit' ? 'PUT' : 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          title: editTitle,
          content: editContent,
          message: commitMessage || `Update ${editTitle}`
        })
      });

      if (res.ok) {
        const data = await res.json();
        window.location.href = `/${repo.data.owner.username}/${repo.data.name}/wiki/${data.url || editTitle}`;
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
    if (!confirm('Are you sure you want to delete this page?')) return;
    try {
      const res = await fetch(`/api/v2/${repo.data.owner.username}/${repo.data.name}/wiki/${path}`, {
        method: 'DELETE'
      });
      if (res.ok) {
        window.location.href = `/${repo.data.owner.username}/${repo.data.name}/wiki`;
      }
    } catch (e) {
      error = 'Failed to delete page';
    }
  }

  function startEdit() {
    mode = 'edit';
    editTitle = currentPage?.name || '';
    editContent = currentPage?.content || '';
  }

  function startNew() {
    mode = 'edit';
    currentPage = null;
    editTitle = '';
    editContent = '';
    path = 'New';
  }

  $effect(() => {
    if (repo.data) {
        fetchWiki();
    }
  });

  $effect(() => {
    if (mode === 'history' && path) {
        fetchHistory();
    }
  });

  const formatDate = (ts: number) => new Date(ts * 1000).toLocaleDateString('es-ES', {
    member: 'short', day: 'numeric', month: 'short', year: 'numeric'
  });
</script>

<div class="wiki-container">
  <div class="wiki-grid">
    <!-- Sidebar -->
    <aside class="wiki-sidebar island">
      <div class="sidebar-header">
        <h2>Wiki</h2>
        <button class="btn-new" onclick={startNew} title="New Page">
          <Add size={20} />
        </button>
      </div>
      
      <nav class="wiki-nav">
        {#each wikiPages as pageItem}
          <a 
            href="/{repo.data.owner.username}/{repo.data.name}/wiki/{pageItem.name}" 
            class="nav-item"
            class:active={path === pageItem.name}
          >
            <Description size={16} />
            <span>{pageItem.name}</span>
          </a>
        {/each}
        {#if wikiPages.length === 0}
          <div class="empty-state">No hay páginas aún</div>
        {/if}
      </nav>
    </aside>

    <!-- Content Area -->
    <main class="wiki-main island">
      {#if loading}
        <div class="wiki-skeleton">
            <div class="skeleton-header"></div>
            <div class="skeleton-body"></div>
        </div>
      {:else}
        <header class="wiki-content-header">
          <div class="header-left">
            {#if mode === 'edit'}
                <input type="text" class="title-input" bind:value={editTitle} placeholder="Título de la página..." />
            {:else}
                <h1>{currentPage?.name || 'Wiki'}</h1>
            {/if}
          </div>

          <div class="header-right">
            <div class="tab-group">
                <button class="tab-btn" class:active={mode === 'view'} onclick={() => mode = 'view'}>
                    Página
                </button>
                <button class="tab-btn" class:active={mode === 'history'} onclick={() => mode = 'history'}>
                    Historial
                </button>
            </div>
            
            <div class="action-group">
                {#if mode === 'view'}
                    <button class="btn-icon-text" onclick={startEdit}>
                        <Edit size={16} />
                        Editar
                    </button>
                {:else if mode === 'edit'}
                    <button class="btn-icon-text variant" onclick={() => mode = 'view'}>
                        <Close size={16} />
                        Cancelar
                    </button>
                    <button class="btn-primary" onclick={handleSave} disabled={saving || !editTitle}>
                        <Save size={16} />
                        {saving ? 'Guardando...' : 'Guardar'}
                    </button>
                {/if}
                
                {#if mode === 'view' && path !== 'Home'}
                    <button class="btn-danger-icon" onclick={handleDelete} title="Eliminar Página">
                        <Delete size={18} />
                    </button>
                {/if}
            </div>
          </div>
        </header>

        <div class="wiki-body">
          {#if mode === 'view'}
            <div class="markdown-body">
              {@html currentPage?.content_html || '<i>Cargando contenido...</i>'}
            </div>
          {:else if mode === 'edit'}
            <div class="edit-area">
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

  /* Sidebar */
  .wiki-sidebar {
    padding: 1.5rem;
    position: sticky;
    top: 1rem;
  }

  .sidebar-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 1.5rem;
  }

  .sidebar-header h2 {
    font-size: 1.25rem;
    margin: 0;
    font-weight: 600;
  }

  .btn-new {
    width: 36px;
    height: 36px;
    border-radius: 10px;
    background: var(--primary-container);
    color: var(--on-primary-container);
    border: none;
    cursor: pointer;
    display: flex;
    align-items: center;
    justify-content: center;
    transition: all 0.2s;
  }

  .btn-new:hover {
    transform: scale(1.05);
    background: var(--primary);
    color: white;
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
    color: var(--text-secondary);
    font-size: 0.95rem;
    transition: all 0.2s;
  }

  .nav-item:hover {
    background: var(--surface-variant);
    color: var(--text-primary);
  }

  .nav-item.active {
    background: var(--primary-container);
    color: var(--primary);
    font-weight: 600;
  }

  /* Main Content */
  .wiki-main {
    min-height: 80vh;
    display: flex;
    flex-direction: column;
  }

  .wiki-content-header {
    padding: 1.5rem 2rem;
    border-bottom: 1px solid var(--border-color);
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
    border-bottom: 2px solid var(--primary);
    color: var(--text-primary);
    width: 100%;
    outline: none;
  }

  .header-right {
    display: flex;
    align-items: center;
    gap: 2rem;
  }

  .tab-group {
    background: var(--surface-variant);
    padding: 0.25rem;
    border-radius: 100px;
    display: flex;
  }

  .tab-btn {
    padding: 0.4rem 1.25rem;
    border-radius: 100px;
    border: none;
    background: transparent;
    color: var(--text-secondary);
    font-size: 0.9rem;
    font-weight: 600;
    cursor: pointer;
    transition: all 0.2s;
  }

  .tab-btn.active {
    background: var(--island-bg);
    color: var(--primary);
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
    background: var(--surface-variant);
    color: var(--text-primary);
    border: none;
    padding: 0.5rem 1rem;
    border-radius: 12px;
    font-weight: 600;
    cursor: pointer;
    transition: all 0.2s;
  }

  .btn-icon-text.variant {
    background: transparent;
    color: var(--text-secondary);
  }

  .btn-primary {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    background: var(--primary);
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
    color: var(--error);
    background: transparent;
    border: none;
    cursor: pointer;
    padding: 0.5rem;
    border-radius: 10px;
    transition: background 0.2s;
  }

  .btn-danger-icon:hover {
    background: rgba(var(--error-rgb), 0.1);
  }

  .wiki-body {
    padding: 2rem;
    flex: 1;
  }

  /* Edit Area */
  .edit-area {
    display: flex;
    flex-direction: column;
    gap: 1.5rem;
    height: 100%;
  }

  .wiki-editor {
    flex: 1;
    min-height: 400px;
    background: var(--surface-variant);
    border: 1px solid var(--border-color);
    border-radius: 16px;
    padding: 1.5rem;
    font-family: 'JetBrains Mono', monospace;
    font-size: 1rem;
    color: var(--text-primary);
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
    color: var(--text-secondary);
  }

  .commit-msg-box input {
    background: var(--surface-variant);
    border: 1px solid var(--border-color);
    border-radius: 10px;
    padding: 0.75rem 1rem;
    color: var(--text-primary);
  }

  /* History List */
  .history-list {
    display: flex;
    flex-direction: column;
    padding-left: 1rem;
    border-left: 2px solid var(--border-color);
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
    background: var(--primary);
    border: 3px solid var(--island-bg);
    border-radius: 50%;
  }

  .history-msg {
    font-weight: 600;
    color: var(--text-primary);
    margin-bottom: 0.25rem;
  }

  .history-meta {
    font-size: 0.85rem;
    color: var(--text-secondary);
  }

  .hash {
    background: var(--surface-variant);
    padding: 0.1rem 0.4rem;
    border-radius: 4px;
    font-family: monospace;
  }

  /* Skeletons */
  .skeleton-header {
    height: 40px;
    width: 300px;
    background: var(--surface-variant);
    border-radius: 8px;
    margin-bottom: 2rem;
    animation: pulse 1.5s infinite;
  }

  .skeleton-body {
    height: 300px;
    background: var(--surface-variant);
    border-radius: 12px;
    animation: pulse 1.5s infinite;
  }

  @keyframes pulse {
    0%, 100% { opacity: 0.5; }
    50% { opacity: 0.8; }
  }

  @keyframes fadeIn {
    from { opacity: 0; transform: translateY(10px); }
    to { opacity: 1; transform: translateY(0); }
  }

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
