<script lang="ts">
  import { apiFetch } from "$lib/api";
  import { page } from '$app/state';
  import { goto } from '$app/navigation';
  import { onMount, getContext } from 'svelte';
  import { getLanguageColor } from '$lib/utils/languages';
  import RepoSearchResults from './RepoSearchResults.svelte';
  import MarkdownRenderer from './ui/MarkdownRenderer.svelte';
  import hljs from 'highlight.js';
  import 'highlight.js/styles/github-dark.css';
  
  let { path = '' } = $props();
  let username = $derived(page.params.username);
  let reponame = $derived(page.params.reponame);

  let repoData = $state<any>(null);
  let loading = $state(true);
  let error = $state('');
  
  let searchQuery = $state('');

  let highlightedLines = $derived(() => {
    if (!repoData?.content || !repoData?.is_text_file) return [];
    const language = repoData.name.split('.').pop() || '';
    let html: string;
    try {
      html = hljs.highlight(repoData.content, { language: hljs.getLanguage(language) ? language : 'plaintext' }).value;
    } catch (e) {
      html = hljs.highlightAuto(repoData.content).value;
    }
    return html.split('\n');
  });

  function handleSearch(e: any) {
    searchQuery = e.detail;
  }

  onMount(() => {
    window.addEventListener('ggit:search', handleSearch);
    return () => window.removeEventListener('ggit:search', handleSearch);
  });

  let realPath = $derived(() => {
    let p = path;
    if (p.startsWith('src/')) {
        const parts = p.split('/');
        return parts.slice(2).join('/');
    }
    return p;
  });

  let lastFetched = { user: '', repo: '', path: '' };

  async function fetchRepo() {
    const currentPath = realPath();
    if (lastFetched.user === username && lastFetched.repo === reponame && lastFetched.path === currentPath) return;
    lastFetched = { user: username || '', repo: reponame || '', path: currentPath || '' };
    
    loading = true;
    error = '';
    try {
      const res = await apiFetch(`/api/v2/${username}/${reponame}/view/${currentPath}`);
      if (res.ok) {
        repoData = await res.json();
      } else {
        error = 'Repository content not found.';
      }
    } catch (e) {
      error = 'Failed to load repository.';
    } finally {
      loading = false;
    }
  }

  $effect(() => {
    fetchRepo();
  });

  function formatSize(bytes: number) {
    if (bytes === 0) return '0 B';
    const k = 1024;
    const sizes = ['B', 'KB', 'MB', 'GB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return parseFloat((bytes / Math.pow(k, i)).toFixed(1)) + ' ' + sizes[i];
  }

  const repoContext = getContext<any>('repo');

  // Calculate language percentages - Fallback to global context if local repoData doesn't have them
  const langEntries = $derived(() => {
    // Correctly fallback to repoContext.languages (set in +layout.svelte)
    const langs = repoData?.languages || repoContext.languages;
    if (!langs) return [];
    const total = Object.values(langs as Record<string, number>).reduce((a, b) => a + b, 0);
    if (total === 0) return [];
    return Object.entries(langs as Record<string, number>)
      .map(([name, bytes]) => ({
        name,
        percent: ((bytes / total) * 100).toFixed(1),
        color: getLanguageColor(name)
      }))
      .sort((a, b) => parseFloat(b.percent) - parseFloat(a.percent));
  });

  let showAllLanguages = $state(false);
  let visibleLanguages = $derived(showAllLanguages ? langEntries() : langEntries().slice(0, 10));

  let sidebarCollapsed = $state(false);

  function toggleSidebar() {
    sidebarCollapsed = !sidebarCollapsed;
    localStorage.setItem('repo-sidebar-collapsed', sidebarCollapsed.toString());
  }

  function timeAgo(dateInput: any) {
    if (!dateInput) return '';
    const date = new Date(dateInput);
    const seconds = Math.floor((Date.now() - date.getTime()) / 1000);
    if (seconds < 60) return 'just now';
    const minutes = Math.floor(seconds / 60);
    if (minutes < 60) return `${minutes}m ago`;
    const hours = Math.floor(minutes / 60);
    if (hours < 24) return `${hours}h ago`;
    const days = Math.floor(hours / 24);
    if (days < 30) return `${days}d ago`;
    return date.toLocaleDateString();
  }

  onMount(() => {
    sidebarCollapsed = localStorage.getItem('repo-sidebar-collapsed') === 'true';
  });
</script>

<div class="workspace-area">
  {#if loading}
    <div class="skeleton-grid">
        <div class="skeleton-main"></div>
        <div class="skeleton-side"></div>
    </div>
  {:else if error}
    <div class="error-state">
      <span class="material-symbols-rounded">error</span>
      <p>{error}</p>
    </div>
  {:else if repoData}
    {#if searchQuery}
      <RepoSearchResults query={searchQuery} branch={repoData.branch || 'HEAD'} />
    {:else if repoData.is_bare}
      <div class="empty-state">
        <div class="empty-icon-ring">
          <span class="material-symbols-rounded">terminal</span>
        </div>
        <h2>Push your code</h2>
        <div class="setup-commands">
          <div class="command-block">
            <span class="cmd-label">Quick Setup</span>
            <pre><code>git remote add origin {repoData.repository?.clone_url}
git branch -M main
git push -u origin main</code></pre>
          </div>
        </div>
      </div>
    {:else}
      <div class="technical-grid" class:sidebar-collapsed={sidebarCollapsed}>
        <!-- Main Explorer / Code -->
        <div class="main-column">
          <!-- Last Commit Bar -->
          {#if repoData.last_commit && repoData.type !== 'file'}
            <div class="last-commit-bar">
                <div class="commit-avatar tonal">{repoData.last_commit.author?.name?.[0] || 'U'}</div>
                <div class="commit-info">
                  <div class="commit-summary-row">
                    <a href={`/${username}/${reponame}/commit/${repoData.last_commit.id}`} class="commit-msg-link">
                      <span class="commit-msg">{repoData.last_commit.summary}</span>
                    </a>
                  </div>
                  <div class="commit-meta-row">
                      <span class="author">{repoData.last_commit.author?.name}</span>
                      <span class="separator">committed</span>
                      <span class="time">{timeAgo(repoData.last_commit.author?.when)}</span>
                  </div>
                </div>
                <div class="commit-actions-right">
                  <a href={`/${username}/${reponame}/commit/${repoData.last_commit.id}`} class="sha-link">
                    <code class="sha">{repoData.last_commit.id?.slice(0,7)}</code>
                  </a>
                  <button class="btn-sidebar-toggle" onclick={toggleSidebar} title={sidebarCollapsed ? "Show Sidebar" : "Hide Sidebar"}>
                    <span class="material-symbols-rounded">{sidebarCollapsed ? 'side_navigation' : 'right_panel_close'}</span>
                  </button>
                </div>
            </div>
          {/if}

          <!-- Mobile Only Language Bar -->
          {#if langEntries().length > 0}
            <div class="lang-bar mobile-only">
               {#each langEntries() as lang}
                 <div class="lang-progress" style="width: {lang.percent}%; background: {lang.color}" title="{lang.name}: {lang.percent}%"></div>
               {/each}
            </div>
          {/if}

          <!-- File Content / Table -->
          <div class="content-view">
            {#if repoData.type === 'directory'}
              <table class="file-table">
                <tbody>
                  {#each repoData.entries as entry}
                    <tr class="file-row clickable" onclick={() => goto(`/${username}/${reponame}/src/${repoData.branch}/${realPath() ? realPath() + '/' : ''}${entry.name}`)}>
                      <td class="name-cell">
                        <span class="material-symbols-rounded icon-{entry.type === 'tree' ? 'tree' : 'blob'}">
                          {entry.type === 'tree' ? 'folder' : 'description'}
                        </span>
                        <a href={`/${username}/${reponame}/src/${repoData.branch}/${realPath() ? realPath() + '/' : ''}${entry.name}`} onclick={(e) => e.stopPropagation()}>
                          {entry.name}
                        </a>
                      </td>
                      <td class="msg-cell">
                        <span class="msg-text">{entry.last_commit?.summary || ''}</span>
                      </td>
                      <td class="date-cell">{timeAgo(entry.last_commit?.author?.when) || ''}</td>
                    </tr>
                  {/each}
                </tbody>
              </table>
            {:else if repoData.type === 'file'}
              <div class="file-viewer">
                <div class="file-header">
                   <div class="file-meta">
                     <span class="material-symbols-rounded">draft</span>
                     <span class="filename">{repoData.name}</span>
                     <span class="filesize">{formatSize(repoData.size)}</span>
                   </div>
                    <div class="file-tools side-by-side">
                      {#if repoData.is_text_file}
                        <a href={`/${username}/${reponame}/_edit/${repoData.branch}/${realPath()}`} class="tool-btn">
                          <span class="material-symbols-rounded">edit</span>
                        </a>
                      {/if}
                      <a href={`/api/v2/${username}/${reponame}/raw/${repoData.branch}/${realPath()}`} target="_blank" class="tool-btn">Raw</a>
                      <a href={`/api/v2/${username}/${reponame}/media/${repoData.branch}/${realPath()}`} class="tool-btn" download>
                        <span class="material-symbols-rounded">download</span>
                      </a>
                      <button class="tool-btn btn-sidebar-toggle-inline" onclick={toggleSidebar} title={sidebarCollapsed ? "Show Sidebar" : "Hide Sidebar"}>
                        <span class="material-symbols-rounded">{sidebarCollapsed ? 'side_navigation' : 'right_panel_close'}</span>
                      </button>
                    </div>
                </div>
                <div class="code-box">
                   {#if repoData.is_text_file}
                      <div class="code-grid">
                        <div class="ln-col">
                          {#each highlightedLines() as _, i}
                            <span>{i + 1}</span>
                          {/each}
                        </div>
                        <pre class="code-pre hljs">
                          {#each highlightedLines() as line}
                            <div class="line">{@html line || '&nbsp;'}</div>
                          {/each}
                        </pre>
                      </div>
                   {:else}
                      <div class="binary-placeholder">
                        <span class="material-symbols-rounded">deployed_code</span>
                        <p>Binary file — {formatSize(repoData.size)}</p>
                      </div>
                   {/if}
                </div>

                {#if repoData.readme_content}
                   <div class="file-readme-integrated">
                      <header class="readme-header small">
                         <span class="material-symbols-rounded">menu_book</span>
                         README Preview
                      </header>
                      <div class="readme-content">
                         <MarkdownRenderer 
                            content={repoData.readme_content} 
                            repoContext={{owner: username, name: reponame}}
                         />
                      </div>
                   </div>
                {/if}
              </div>
            {/if}
          </div>

          <!-- README Section (Root view) -->
          {#if repoData.readme}
            <div class="readme-section">
                <header class="readme-header">
                    <span class="material-symbols-rounded">menu_book</span>
                    {repoData.readme.name}
                </header>
                <div class="readme-content">
                    <MarkdownRenderer 
                      content={repoData.readme.content} 
                      repoContext={{owner: username, name: reponame}}
                    />
                </div>
            </div>
          {/if}
        </div>

        <!-- Sidebar -->
        <aside class="side-column" class:hidden-sidebar={sidebarCollapsed}>
          {#if langEntries().length > 0}
            <div class="side-section desktop-only">
                <h4 class="side-title">Languages</h4>
                <div class="lang-bar">
                    {#each langEntries() as lang}
                        <div class="lang-progress" style="width: {lang.percent}%; background: {lang.color}" title="{lang.name}: {lang.percent}%"></div>
                    {/each}
                </div>
                <div class="lang-list">
                    {#each visibleLanguages as lang}
                        <div class="lang-item">
                            <span class="dot" style="background: {lang.color}"></span>
                            <span class="name">{lang.name}</span>
                            <span class="percent">{lang.percent}%</span>
                        </div>
                    {/each}
                </div>
                {#if langEntries().length > 10}
                    <button class="lang-toggle" onclick={() => showAllLanguages = !showAllLanguages}>
                        {showAllLanguages ? 'Show less' : `+ ${langEntries().length - 10} more`}
                    </button>
                {/if}
            </div>
          {/if}

          <div class="side-section">
            <h4 class="side-title">About</h4>
            <p class="about-text">{repoContext.data?.description || 'No description provided.'}</p>
          </div>
        </aside>
      </div>
    {/if}
  {/if}
</div>

<style>
  .workspace-area { height: 100%; }

  .technical-grid {
    display: grid;
    grid-template-columns: 1fr 280px;
    gap: 2.5rem;
    transition: all 0.4s cubic-bezier(0.19, 1, 0.22, 1);
  }

  .technical-grid.sidebar-collapsed {
    grid-template-columns: 1fr 0px;
    gap: 0;
  }

  .main-column { display: flex; flex-direction: column; min-width: 0; }

  /* LAST COMMIT BAR */
  .last-commit-bar {
    display: flex;
    align-items: center;
    gap: 1.25rem;
    padding: 1rem 1.25rem;
    background: var(--surface-variant);
    border-radius: 1rem;
    margin-bottom: 1.5rem;
    border: 1px solid var(--md-sys-color-outline-variant);
    min-width: 0;
  }
  .commit-avatar { width: 32px; height: 32px; border-radius: 8px; display: flex; align-items: center; justify-content: center; font-weight: 800; font-size: 0.9rem; flex-shrink: 0; border: 1px solid var(--md-sys-color-outline-variant); }
  
  .commit-info {
    display: flex;
    flex-direction: column;
    gap: 0.2rem;
    flex: 1;
    min-width: 0;
  }
  
  .commit-summary-row { display: flex; align-items: center; min-width: 0; }
  
  .commit-msg { 
    font-weight: 700; 
    font-size: 0.95rem; 
    color: var(--text-color);
    overflow: hidden; 
    text-overflow: ellipsis; 
    white-space: nowrap;
    display: block;
  }
  
  .commit-meta-row { 
    display: flex; 
    align-items: center;
    gap: 0.4rem; 
    font-size: 0.8rem; 
    color: var(--text-muted);
    opacity: 0.8;
  }
  .author { font-weight: 800; color: var(--text-color); }
  .separator { opacity: 0.5; }

  .commit-msg-link, .sha-link {
    text-decoration: none;
    color: inherit;
    min-width: 0;
  }
  .commit-msg-link:hover .commit-msg {
    color: var(--primary-color);
    text-decoration: underline;
  }
  .sha { font-family: 'JetBrains Mono', monospace; font-size: 0.75rem; font-weight: 800; background: rgba(var(--primary-rgb), 0.1); color: var(--primary-color); padding: 2px 8px; border-radius: 6px; border: 1px solid rgba(var(--primary-rgb), 0.2); transition: all 0.2s; }
  .sha-link:hover .sha { background: var(--primary-color); color: var(--on-primary); }

  .commit-actions-right {
    display: flex;
    align-items: center;
    gap: 1rem;
    flex-shrink: 0;
  }

  .btn-sidebar-toggle {
    background: transparent;
    border: none;
    color: var(--text-muted);
    cursor: pointer;
    display: flex;
    padding: 6px;
    border-radius: 8px;
    transition: all 0.2s;
  }
  .btn-sidebar-toggle:hover {
    background: rgba(var(--primary-rgb), 0.1);
    color: var(--primary-color);
  }
  .btn-sidebar-toggle span {
    font-size: 1.4rem;
  }

  /* FILE TABLE */
  .content-view { 
    background: var(--bg-color); 
    border: 1px solid var(--md-sys-color-outline-variant);
    border-radius: 1rem; 
    overflow: hidden; 
    box-shadow: var(--shadow-sm);
  }
  .file-table { width: 100%; border-collapse: collapse; }
  .file-row { border-bottom: 1px solid var(--md-sys-color-outline-variant); }
  .file-row:last-child { border-bottom: none; }
  .file-row td { padding: 0.85rem 1.25rem; vertical-align: middle; }
  .file-row.clickable { cursor: pointer; }
  .file-row:hover { background: rgba(var(--primary-rgb), 0.04); }

  .name-cell { 
    display: flex; 
    align-items: center; 
    gap: 1rem; 
    min-width: 200px;
  }
  .name-cell a { 
    color: var(--text-color); 
    text-decoration: none; 
    font-weight: 700; 
    font-size: 0.95rem;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }
  .icon-tree { color: #facc15; font-size: 1.4rem; }
  .icon-blob { color: var(--text-muted); font-size: 1.4rem; opacity: 0.7; }
  
  .msg-cell { 
    font-size: 0.85rem; 
    color: var(--text-muted); 
    max-width: 0;
    width: 60%;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }
  .msg-text {
    opacity: 0.7;
  }
  .date-cell { text-align: right; font-size: 0.8rem; color: var(--text-muted); opacity: 0.6; white-space: nowrap; }

  /* VIEWER */
  .file-header { 
    background: var(--surface-variant); 
    padding: 0.85rem 1.25rem; 
    display: flex; 
    justify-content: space-between; 
    align-items: center;
    border-bottom: 1px solid var(--md-sys-color-outline-variant);
  }
  .file-meta { display: flex; align-items: center; gap: 0.75rem; font-weight: 800; font-size: 0.95rem; color: var(--text-color); }
  .filesize { font-weight: 400; opacity: 0.5; font-size: 0.8rem; }
  
  .tool-btn { background: var(--bg-color); border: 1px solid var(--md-sys-color-outline-variant); color: var(--text-muted); padding: 0.4rem 0.85rem; border-radius: 0.6rem; cursor: pointer; font-size: 0.8rem; font-weight: 800; transition: all 0.2s; text-decoration: none; display: flex; align-items: center; gap: 0.5rem; }
  .tool-btn:hover { border-color: var(--primary-color); color: var(--primary-color); background: rgba(var(--primary-rgb), 0.05); }
  .file-tools.side-by-side { display: flex; align-items: center; gap: 0.5rem; }

  .code-box { background: #0d0e11; overflow-x: auto; }
  .code-grid { display: grid; grid-template-columns: 60px 1fr; }
  .ln-col { border-right: 1px solid rgba(255,255,255,0.05); padding: 1.5rem 0; text-align: right; background: rgba(255,255,255,0.01); }
  .ln-col span { display: block; font-family: 'JetBrains Mono', monospace; font-size: 0.8rem; color: var(--text-muted); opacity: 0.3; padding-right: 1.25rem; line-height: 1.6; }
  .code-pre { margin: 0; padding: 1.5rem; font-family: 'JetBrains Mono', monospace; font-size: 0.9rem; line-height: 1.6; color: #e6edf3; }
  .code-pre.hljs { background: transparent; }
  .line { white-space: pre; }

  /* README */
  .readme-section { margin-top: 2rem; background: var(--bg-color); border: 1px solid var(--md-sys-color-outline-variant); border-radius: 1.25rem; overflow: hidden; box-shadow: var(--shadow-sm); }
  .readme-header { padding: 1.25rem 1.5rem; border-bottom: 1px solid var(--md-sys-color-outline-variant); display: flex; align-items: center; gap: 0.75rem; font-weight: 800; font-size: 0.8rem; text-transform: uppercase; letter-spacing: 0.05em; color: var(--text-muted); background: var(--surface-variant); }
  .readme-header.small { padding: 1rem 1.5rem; }
  .readme-content { padding: 2rem; line-height: 1.6; }

  /* SIDEBAR */
  .side-column { display: flex; flex-direction: column; gap: 2.5rem; transition: opacity 0.3s; }
  .side-column.hidden-sidebar { opacity: 0; pointer-events: none; }
  .side-title { font-size: 0.75rem; font-weight: 800; text-transform: uppercase; letter-spacing: 0.1em; color: var(--text-muted); margin-bottom: 1rem; border-bottom: 1px solid var(--md-sys-color-outline-variant); padding-bottom: 0.5rem; }
  .about-text { font-size: 0.95rem; line-height: 1.6; color: var(--text-color); opacity: 0.85; }
  
  .lang-bar { 
    height: 12px; 
    border-radius: 6px; 
    overflow: hidden; 
    display: flex; 
    background: rgba(0,0,0,0.2); 
    margin-bottom: 1.25rem;
    box-shadow: inset 0 1px 3px rgba(0,0,0,0.3);
  }
  .mobile-only { display: none; }
  .lang-progress { height: 100%; transition: width 1s ease; }
  .lang-list { display: flex; flex-direction: column; gap: 0.6rem; }
  .lang-item { display: flex; align-items: center; gap: 0.75rem; font-size: 0.85rem; }
  .lang-item .dot { width: 10px; height: 10px; border-radius: 50%; }
  .lang-item .name { font-weight: 700; }
  .lang-item .percent { color: var(--text-muted); margin-left: auto; font-family: 'JetBrains Mono', monospace; font-size: 0.75rem; }

  .lang-toggle { background: transparent; border: none; color: var(--primary-color); font-size: 0.8rem; font-weight: 800; cursor: pointer; padding: 0.75rem 0; display: flex; align-items: center; gap: 0.5rem; }
  .lang-toggle:hover { text-decoration: underline; }

  @media (max-width: 900px) {
    .technical-grid { grid-template-columns: 1fr; gap: 1.5rem; }
    .mobile-only { display: flex; }
    .desktop-only { display: none; }
  }

  .skeleton-grid { display: grid; grid-template-columns: 1fr 280px; gap: 2.5rem; height: 500px; }
  .skeleton-main, .skeleton-side { background: var(--surface-variant); border-radius: 1rem; animation: pulse 1.5s infinite; }
  @keyframes pulse { 0%, 100% { opacity: 0.5; } 50% { opacity: 0.8; } }

  .empty-state { padding: 5rem 2rem; text-align: center; }
  .empty-icon-ring { width: 80px; height: 80px; border: 2px dashed var(--primary-color); border-radius: 50%; display: flex; align-items: center; justify-content: center; margin: 0 auto 2rem; color: var(--primary-color); }
  .empty-icon-ring span { font-size: 2.5rem; }
  .command-block { background: #000; padding: 1.5rem; border-radius: 1rem; text-align: left; margin-top: 2.5rem; border: 1px solid var(--md-sys-color-outline-variant); }
  .cmd-label { font-size: 0.75rem; font-weight: 800; color: var(--primary-color); display: block; margin-bottom: 0.75rem; text-transform: uppercase; }

  :global(.markdown-body) { color: var(--text-color); font-size: 1rem; line-height: 1.7; }
</style>
