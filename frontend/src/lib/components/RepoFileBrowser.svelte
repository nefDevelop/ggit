<script lang="ts">
  import { 
    AlertCircle, Terminal, History, PanelLeftOpen, PanelRightClose, 
    Folder, FileText, File, Edit, Download, Binary, BookText,
    Rocket, ChevronRight
  } from "lucide-svelte";
  import { apiFetch } from "$lib/api";
  import { page } from '$app/state';
  import { goto } from '$app/navigation';
  import { onMount, getContext } from 'svelte';
  import { getLanguageColor } from '$lib/utils/languages';
  import RepoSearchResults from './RepoSearchResults.svelte';
  import MarkdownRenderer from './ui/MarkdownRenderer.svelte';
  import hljs from 'highlight.js';
  
  let { path = '' } = $props();
  let username = $derived(page.params.username);
  let reponame = $derived(page.params.reponame);

  let repoData = $state<any>(null);
  let loading = $state(true);
  let error = $state('');
  
  let searchQuery = $state('');

  let codeRows = $derived(() => {
    if (!repoData?.content || !repoData?.is_text_file) return [];
    
    let lines: string[];
    if (repoData.highlighted_content) {
      lines = repoData.highlighted_content.split('\n');
    } else {
      const language = repoData.name.split('.').pop() || '';
      let html: string;
      
      // Skip highlighting for very large files to prevent UI thread lock
      if (repoData.content.length > 100000) {
        html = repoData.content.replace(/[&<>"']/g, (m: string) => ({'&':'&amp;','<':'&lt;','>':'&gt;','"':'&quot;',"'":"&#39;"}[m] || m));
      } else {
        try {
          html = hljs.highlight(repoData.content, { language: hljs.getLanguage(language) ? language : 'plaintext' }).value;
        } catch (e) {
          html = hljs.highlightAuto(repoData.content).value;
        }
      }
      lines = html.split('\n');
    }

    return lines.map((content, i) => ({
      num: i + 1,
      content: content || '&nbsp;',
      blame: showBlame ? blameData[i] : null
    }));
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

  let blameData = $state<any[]>([]);
  let showBlame = $state(false);
  let blameLoading = $state(false);

  async function toggleBlame() {
    if (showBlame) {
      showBlame = false;
      return;
    }

    if (blameData.length > 0) {
      showBlame = true;
      return;
    }
    
    blameLoading = true;
    try {
      const res = await apiFetch(`/api/v2/${username}/${reponame}/blame/${realPath()}?ref=${repoData.branch}`);
      if (res.ok) {
        blameData = await res.json();
        showBlame = true;
      }
    } catch (e) {
      console.error('Failed to fetch blame', e);
    } finally {
      blameLoading = false;
    }
  }
</script>

<div class="workspace-area">
  {#if loading}
    <div class="technical-grid skeleton-mode">
      <div class="main-column">
        <div class="skeleton-pulse last-commit-skeleton"></div>
        <div class="content-view skeleton-table-view">
          {#each Array(8) as _}
            <div class="skeleton-row-item">
              <div class="skeleton-pulse icon-placeholder"></div>
              <div class="skeleton-pulse text-placeholder"></div>
              <div class="skeleton-pulse msg-placeholder"></div>
            </div>
          {/each}
        </div>
      </div>
      <aside class="side-column">
        <div class="side-section">
          <div class="skeleton-pulse title-placeholder"></div>
          <div class="skeleton-pulse text-placeholder" style="height: 4rem;"></div>
        </div>
        <div class="side-section">
          <div class="skeleton-pulse title-placeholder"></div>
          <div class="skeleton-pulse lang-bar-placeholder"></div>
        </div>
      </aside>
    </div>
  {:else if error}
    <div class="empty-state">
      <AlertCircle size={48} />
      <h3>No content found</h3>
      <p>{error || 'This repository is empty or the path does not exist.'}</p>
    </div>
  {:else if repoData}
    {#if searchQuery}
      <RepoSearchResults query={searchQuery} branch={repoData.branch || 'HEAD'} />
    {:else if repoData.is_bare}
      <div class="empty-state-push">
        <div class="empty-icon-ring">
          <Terminal size={40} />
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
                    {#if sidebarCollapsed}
                      <PanelLeftOpen size={20} />
                    {:else}
                      <PanelRightClose size={20} />
                    {/if}
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
                        {#if entry.type === 'tree'}
                          <Folder size={18} class="icon-tree" />
                        {:else}
                          <FileText size={18} class="icon-blob" />
                        {/if}
                        <a href={`/${username}/${reponame}/src/${repoData.branch}/${realPath() ? realPath() + '/' : ''}${entry.name}`} data-sveltekit-preload-data="hover" onclick={(e) => e.stopPropagation()}>
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
                     <File size={18} />
                     <span class="filename">{repoData.name}</span>
                     <span class="filesize">{formatSize(repoData.size)}</span>
                   </div>
                    <div class="file-tools side-by-side">
                      {#if repoData.is_text_file}
                        <button class="tool-btn" class:active={showBlame} onclick={toggleBlame} disabled={blameLoading}>
                           {#if blameLoading}
                             <div class="mini-spinner"></div>
                           {:else}
                             <History size={16} />
                             Blame
                           {/if}
                        </button>
                        <a href={`/${username}/${reponame}/_edit/${repoData.branch}/${realPath()}`} class="tool-btn">
                          <Edit size={16} />
                        </a>
                      {/if}
                      <a href={`/api/v2/${username}/${reponame}/raw/${repoData.branch}/${realPath()}`} target="_blank" class="tool-btn">Raw</a>
                      <a href={`/api/v2/${username}/${reponame}/media/${repoData.branch}/${realPath()}`} class="tool-btn" download title="Download">
                        <Download size={16} />
                      </a>
                      <button class="tool-btn btn-sidebar-toggle-inline" onclick={toggleSidebar} title={sidebarCollapsed ? "Show Sidebar" : "Hide Sidebar"}>
                        {#if sidebarCollapsed}
                          <PanelLeftOpen size={18} />
                        {:else}
                          <PanelRightClose size={18} />
                        {/if}
                      </button>
                    </div>
                </div>
                <div class="code-box">
                    {#if repoData.is_text_file}
                       <div class="code-container" class:with-blame={showBlame}>
                          {#each codeRows() as row}
                            <div class="code-line-row">
                               {#if showBlame && row.blame}
                                  <div class="blame-cell" title="{row.blame.author} - {new Date(row.blame.time * 1000).toLocaleString()}">
                                     <span class="blame-author">{row.blame.author.split(' ')[0]}</span>
                                     <span class="blame-time">{timeAgo(row.blame.time * 1000)}</span>
                                  </div>
                               {/if}
                               <div class="ln-cell">{row.num}</div>
                               <div class="code-cell">
                                 <pre class="hljs"><code>{@html row.content}</code></pre>
                               </div>
                            </div>
                          {/each}
                       </div>
                    {:else}
                      <div class="binary-placeholder">
                        <Binary size={48} />
                        <p>Binary file — {formatSize(repoData.size)}</p>
                      </div>
                   {/if}
                </div>

                {#if repoData.readme_content}
                   <div class="file-readme-integrated">
                      <header class="readme-header small">
                         <BookText size={16} />
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
                    <BookText size={18} />
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
        
        <aside class="side-column" class:hidden-sidebar={sidebarCollapsed}>

          <div class="side-section">
            <h4 class="side-title">About</h4>
            <p class="about-text">{repoContext.data?.description || 'No description provided.'}</p>
          </div>

          <div class="side-section">
            <h4 class="side-title">Releases</h4>
            {#if repoContext.stats.releases > 0}
               <a href={`/${username}/${reponame}/releases`} class="side-item-link">
                  <div class="side-item-main">
                     <Rocket size={18} class="icon-primary" />
                     <span class="side-item-text"><b>{repoContext.stats.releases}</b> releases</span>
                  </div>
                  <ChevronRight size={18} />
               </a>
            {:else}
               <p class="about-text">No releases published.</p>
            {/if}
          </div>

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
    padding: 2rem;
    animation: fadeIn 0.4s ease-out;
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
    background: var(--md-sys-color-secondary-container);
    border-radius: 1rem;
    margin-bottom: 1.5rem;
    border: 1px solid var(--md-sys-color-outline-variant);
    min-width: 0;
  }
  .commit-avatar {
    width: 32px;
    height: 32px;
    background: var(--md-sys-color-secondary);
    color: var(--md-sys-color-on-secondary);
    border-radius: 8px;
    display: flex;
    align-items: center;
    justify-content: center;
    font-weight: 800;
    font-size: 0.9rem;
    flex-shrink: 0;
    border: 1px solid var(--md-sys-color-outline-variant);
  }
  
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
  .sha {
    font-family: monospace;
    background: var(--md-sys-color-tertiary-container);
    padding: 0.2rem 0.5rem;
    border-radius: 4px;
    color: var(--md-sys-color-on-tertiary-container);
    font-size: 0.85rem;
  }
  .sha-link:hover .sha { background: var(--tertiary); color: var(--on-tertiary); }

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

  /* FILE TABLE */
  .content-view { 
    background: var(--md-sys-color-secondary-container);
    border: 3px solid var(--md-sys-color-outline-variant);
    border-radius: 1rem; 
    overflow: hidden; 
    box-shadow: var(--shadow-lg);
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
  :global(.icon-tree) { color: #facc15; }
  :global(.icon-blob) { color: var(--text-muted); opacity: 0.7; }
  
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
  
  .tool-btn { background: var(--bg-color); border: 1px solid var(--md-sys-color-outline-variant); color: var(--text-muted); padding: 0 0.85rem; border-radius: 0.6rem; cursor: pointer; font-size: 0.8rem; font-weight: 800; transition: all 0.2s; text-decoration: none; display: flex; align-items: center; justify-content: center; gap: 0.5rem; height: 34px; box-sizing: border-box; }
  .tool-btn:hover { border-color: var(--primary-color); color: var(--primary-color); background: rgba(var(--primary-rgb), 0.05); }
  .file-tools.side-by-side { display: flex; align-items: center; gap: 0; }
  .file-tools.side-by-side .tool-btn { border-radius: 0; border-left-width: 0; }
  .file-tools.side-by-side .tool-btn:first-child { border-top-left-radius: 0.6rem; border-bottom-left-radius: 0.6rem; border-left-width: 1px; }
  .file-tools.side-by-side .tool-btn:last-child { border-top-right-radius: 0.6rem; border-bottom-right-radius: 0.6rem; }  .code-box { background: var(--md-sys-color-surface-container-lowest); overflow-x: auto; border: 1px solid var(--md-sys-color-outline-variant); margin: 1rem; border-radius: 0.75rem; }
  .code-container { display: flex; flex-direction: column; }
  .code-line-row { display: flex; align-items: stretch; min-height: 1.6rem; }
  .code-line-row:hover { background: rgba(var(--primary-rgb), 0.05); }
  
  .ln-cell { width: 60px; min-width: 60px; text-align: right; padding-right: 1.25rem; font-family: 'JetBrains Mono', monospace; font-size: 0.85rem; line-height: 1.6; color: var(--text-muted); opacity: 0.4; user-select: none; border-right: 1px solid var(--md-sys-color-outline-variant); background: var(--md-sys-color-surface-container-low); display: flex; align-items: flex-start; justify-content: flex-end; padding-top: 2px; }
  
  .code-cell { flex: 1; padding: 0 1rem; display: flex; align-items: flex-start; overflow-x: auto; }
  .code-cell pre { margin: 0; padding: 0; font-family: 'JetBrains Mono', monospace; font-size: 0.85rem; line-height: 1.6; background: transparent; }
  .code-cell pre code { background: transparent; padding: 0; }

  .blame-cell { width: 140px; min-width: 140px; border-right: 1px solid var(--md-sys-color-outline-variant); background: var(--md-sys-color-surface-container-low); padding: 0 0.75rem; font-size: 0.75rem; display: flex; align-items: center; gap: 0.5rem; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; opacity: 0.6; }
  .blame-author { font-weight: 800; color: var(--primary-color); }
  .blame-time { color: var(--text-muted); font-size: 0.7rem; }

  .tool-btn.active { background: var(--primary-color) !important; color: white !important; }
  .mini-spinner { width: 14px; height: 14px; border: 2px solid rgba(255,255,255,0.3); border-top-color: white; border-radius: 50%; animation: spin 0.8s linear infinite; }
  @keyframes spin { to { transform: rotate(360deg); } }

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
  
  .side-item-link {
    display: flex;
    align-items: center;
    justify-content: space-between;
    text-decoration: none;
    color: var(--text-color);
    padding: 0.75rem 1rem;
    background: rgba(var(--primary-rgb), 0.05);
    border-radius: 12px;
    border: 1px solid var(--md-sys-color-outline-variant);
    transition: all 0.2s;
  }
  .side-item-link:hover {
    background: rgba(var(--primary-rgb), 0.1);
    border-color: var(--primary-color);
    transform: translateY(-1px);
  }
  .side-item-main {
    display: flex;
    align-items: center;
    gap: 0.75rem;
  }
  :global(.icon-primary) {
    color: var(--primary-color);
  }
  .side-item-text {
    font-size: 0.9rem;
  }
  
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

  .skeleton-mode { pointer-events: none; }
  .skeleton-pulse { background: rgba(var(--primary-rgb), 0.05); animation: pulse 1.5s infinite ease-in-out; border-radius: 0.5rem; }
  
  .last-commit-skeleton { height: 60px; width: 100%; margin-bottom: 1.5rem; border-radius: 1rem; }
  .skeleton-table-view { background: transparent; border-color: var(--md-sys-color-outline-variant); }
  
  .skeleton-row-item { display: flex; align-items: center; gap: 1rem; padding: 1rem 1.25rem; border-bottom: 1px solid var(--md-sys-color-outline-variant); }
  .icon-placeholder { width: 20px; height: 20px; flex-shrink: 0; }
  .text-placeholder { height: 1rem; width: 150px; flex-shrink: 0; }
  .msg-placeholder { height: 0.8rem; width: 100%; opacity: 0.5; }
  
  .title-placeholder { height: 1rem; width: 100px; margin-bottom: 1rem; }
  .lang-bar-placeholder { height: 12px; width: 100%; border-radius: 6px; }

  @keyframes pulse { 0% { opacity: 0.3; } 50% { opacity: 0.7; } 100% { opacity: 0.3; } }

  .empty-state { padding: 5rem 2rem; text-align: center; display: flex; flex-direction: column; align-items: center; gap: 1rem; color: var(--text-muted); }
  .empty-state-push { padding: 5rem 2rem; text-align: center; }
  .empty-icon-ring { width: 80px; height: 80px; border: 2px dashed var(--tertiary); border-radius: 50%; display: flex; align-items: center; justify-content: center; margin: 0 auto 2rem; color: var(--tertiary); }
  .command-block { background: #000; padding: 1.5rem; border-radius: 1rem; text-align: left; margin-top: 2.5rem; border: 1px solid var(--md-sys-color-outline-variant); }
  .cmd-label { font-size: 0.75rem; font-weight: 800; color: var(--primary-color); display: block; margin-bottom: 0.75rem; text-transform: uppercase; }

  :global(.markdown-body) { color: var(--text-color); font-size: 1rem; line-height: 1.7; }
</style>
