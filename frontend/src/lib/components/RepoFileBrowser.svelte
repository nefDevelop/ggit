<script lang="ts">
  import { page } from '$app/state';
  import { onMount, getContext } from 'svelte';
  import { getLanguageColor } from '$lib/utils/languages';
  
  let { path = '' } = $props();
  let username = $derived(page.params.username);
  let reponame = $derived(page.params.reponame);

  let repoData = $state<any>(null);
  let loading = $state(true);
  let error = $state('');

  async function fetchRepo() {
    loading = true;
    error = '';
    try {
      const res = await fetch(`/api/v2/${username}/${reponame}/view/${path}`);
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

  const healthStats = $derived({
    commits: 124,
    branches: repoContext?.branches?.length || 1,
    tags: 3,
    license: 'MIT'
  });
</script>

<div class="workspace-area">
  {#if loading}
    <div class="skeleton-grid"></div>
  {:else if error}
    <div class="error-state">
      <span class="material-symbols-rounded">error</span>
      <p>{error}</p>
    </div>
  {:else if repoData}
    {#if repoData.is_bare}
      <div class="empty-state">
        <div class="empty-icon-ring">
          <span class="material-symbols-rounded">terminal</span>
        </div>
        <h2>Push your code</h2>
        <div class="setup-commands">
          <div class="command-block">
            <span class="cmd-label">Quick Setup</span>
            <pre><code>git remote add origin {repoData.repository.clone_url}
git branch -M main
git push -u origin main</code></pre>
          </div>
        </div>
      </div>
    {:else}
      <div class="technical-grid">
        <!-- Main Explorer / Code -->
        <div class="main-column">
          <!-- Branch & Path Row -->
          <div class="browser-controls">
            <div class="branch-selector">
              <span class="material-symbols-rounded">account_tree</span>
              <select 
                class="branch-select"
                value={repoData.branch}
                onchange={(e) => {
                  const newBranch = e.currentTarget.value;
                  window.location.href = `/${username}/${reponame}/${newBranch}`;
                }}
              >
                {#each repoContext.branches as b}
                  <option value={b}>{b}</option>
                {/each}
              </select>
            </div>

            <nav class="breadcrumbs">
              <a href="/{username}/{reponame}" class="crumb-root">
                <span class="material-symbols-rounded">home</span>
              </a>
              {#each path.split('/').filter(p => p) as part, i}
                <span class="sep">/</span>
                <a href="/{username}/{reponame}/{path.split('/').slice(0, i+1).join('/')}" class="crumb-link">
                  {part}
                </a>
              {/each}
            </nav>
          </div>

          <!-- Git Stats (Flat design) -->
          <div class="health-row">
            <a href="/{username}/{reponame}/commits" class="health-item">
              <span class="material-symbols-rounded">history</span>
              <span class="count">{healthStats.commits}</span>
              <span class="label">Commits</span>
            </a>
            <a href="/{username}/{reponame}/branches" class="health-item">
              <span class="material-symbols-rounded">account_tree</span>
              <span class="count">{healthStats.branches}</span>
              <span class="label">Branches</span>
            </a>
            <div class="health-item">
              <span class="material-symbols-rounded">tag</span>
              <span class="count">{healthStats.tags}</span>
              <span class="label">Tags</span>
            </div>
            <div class="health-item license">
              <span class="material-symbols-rounded">balance</span>
              <span class="label">{healthStats.license}</span>
            </div>
          </div>

          <!-- Last Commit (Flat design) -->
          <div class="last-commit-bar">
            <div class="commit-avatar">A</div>
            <div class="commit-info">
               <span class="commit-msg">Refactor technical workspace layout for density</span>
               <span class="commit-meta">
                 <span class="author">nef734</span>
                 <span class="time">2 hours ago</span>
               </span>
            </div>
            <code class="sha">a1b2c3d</code>
          </div>

          <!-- File Content / Table -->
          <div class="content-view">
            {#if repoData.type === 'directory'}
              <table class="file-table">
                <tbody>
                  {#each repoData.entries as entry}
                    <tr class="file-row">
                      <td class="name-cell">
                        <span class="material-symbols-rounded icon-{entry.type === 'tree' ? 'tree' : 'blob'}">
                          {entry.type === 'tree' ? 'folder' : 'description'}
                        </span>
                        <a href="/{username}/{reponame}/{path ? path + '/' : ''}{entry.name}">
                          {entry.name}
                        </a>
                      </td>
                      <td class="msg-cell">
                        <span class="msg-text">Update {entry.name}</span>
                      </td>
                      <td class="date-cell">2d</td>
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
                   <div class="file-tools">
                     <button class="tool-btn">Raw</button>
                     <button class="tool-btn"><span class="material-symbols-rounded">download</span></button>
                   </div>
                </div>
                <div class="code-box">
                   {#if repoData.is_text_file}
                      <div class="code-grid">
                        <div class="ln-col">
                          {#each repoData.highlighted_content as _, i}
                            <span>{i + 1}</span>
                          {/each}
                        </div>
                        <pre class="code-pre">
                          {#each repoData.highlighted_content as line}
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
              </div>
            {/if}
          </div>

          <!-- README Section -->
          <div class="readme-section">
             <header class="readme-header">
                <span class="material-symbols-rounded">menu_book</span>
                README.md
             </header>
             <div class="readme-content">
                <p>Markdown content...</p>
             </div>
          </div>
        </div>

        <!-- Sidebar (Inside technical-grid) -->
        <aside class="side-column">
          <div class="side-section">
            <h4 class="side-title">About</h4>
            <p class="about-text">{repoContext.data?.description || 'No description.'}</p>
          </div>

          <div class="side-section">
            <h4 class="side-title">Languages</h4>
            <div class="lang-bar">
               <div class="lang-progress" style="width: 70%; background: #00ADD8"></div>
               <div class="lang-progress" style="width: 30%; background: #3178C6"></div>
            </div>
            <div class="lang-list">
               <div class="lang-item">
                 <span class="dot" style="background: #00ADD8"></span>
                 <span class="name">Go</span>
                 <span class="percent">70%</span>
               </div>
               <div class="lang-item">
                 <span class="dot" style="background: #3178C6"></span>
                 <span class="name">TypeScript</span>
                 <span class="percent">30%</span>
               </div>
            </div>
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
  }

  .main-column {
    display: flex;
    flex-direction: column;
    min-width: 0;
  }

  /* Browser Controls */
  .browser-controls { display: flex; align-items: center; gap: 1rem; margin-bottom: 1.25rem; }
  .branch-selector { display: flex; align-items: center; gap: 0.5rem; background: var(--md3-surface-container); padding: 0.4rem 0.75rem; border-radius: 0.5rem; border: 1px solid var(--border-color); }
  .branch-select { background: transparent; border: none; color: var(--text-color); font-weight: 700; outline: none; cursor: pointer; }
  
  .breadcrumbs { display: flex; align-items: center; gap: 0.4rem; background: var(--md3-surface-2); padding: 0.4rem 1rem; border-radius: 0.5rem; }
  .crumb-root { color: var(--primary-color); display: flex; }
  .crumb-link { color: var(--text-muted); text-decoration: none; font-weight: 600; font-size: 0.95rem; }
  .sep { opacity: 0.3; }

  /* Integrated Bars */
  .health-row {
    display: flex;
    padding: 0.75rem 0.5rem;
    gap: 2rem;
    border-top: 1px solid var(--border-color);
    border-bottom: 1px solid var(--border-color);
    margin-bottom: 1rem;
  }
  .health-item { display: flex; align-items: center; gap: 0.5rem; color: var(--text-muted); text-decoration: none; font-size: 0.9rem; }
  .health-item .count { color: var(--text-color); font-weight: 700; }
  .health-item.license { margin-left: auto; }

  .last-commit-bar {
    display: flex;
    align-items: center;
    gap: 1rem;
    padding: 0.75rem 1rem;
    background: var(--md3-surface-container);
    border-radius: 0.5rem;
    border: 1px solid var(--border-color);
    margin-bottom: 1rem;
  }
  .commit-avatar { width: 24px; height: 24px; background: var(--primary-color); border-radius: 4px; display: flex; align-items: center; justify-content: center; font-weight: 800; font-size: 0.7rem; color: var(--on-primary); }
  .commit-msg { font-weight: 600; font-size: 0.95rem; flex: 1; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
  .commit-meta { display: flex; gap: 0.5rem; font-size: 0.8rem; color: var(--text-muted); }
  .sha { font-family: monospace; opacity: 0.5; font-size: 0.8rem; }

  /* Content (Table or Viewer) */
  .content-view {
    background: var(--md3-surface-1);
    border: 1px solid var(--border-color);
    border-radius: 0.5rem;
    overflow: hidden;
  }

  .file-table { width: 100%; border-collapse: collapse; }
  .file-row { border-bottom: 1px solid rgba(255,255,255,0.02); }
  .file-row:hover { background: rgba(255,255,255,0.01); }
  .file-row td { padding: 0.6rem 1rem; }
  .name-cell { display: flex; align-items: center; gap: 0.75rem; min-width: 250px; }
  .name-cell a { color: var(--text-color); text-decoration: none; font-weight: 500; }
  .icon-tree { color: #facc15; font-size: 1.2rem; }
  .icon-blob { color: var(--text-muted); font-size: 1.2rem; }
  .msg-cell { font-size: 0.85rem; color: var(--text-muted); opacity: 0.7; }
  .date-cell { text-align: right; font-size: 0.8rem; color: var(--text-muted); opacity: 0.5; }

  /* File Viewer */
  .file-header { background: var(--md3-surface-3); padding: 0.5rem 1rem; display: flex; justify-content: space-between; align-items: center; border-bottom: 1px solid var(--border-color); }
  .file-meta { display: flex; align-items: center; gap: 0.5rem; font-weight: 600; font-size: 0.9rem; }
  .tool-btn { background: transparent; border: 1px solid var(--border-color); color: var(--text-muted); padding: 0.2rem 0.6rem; border-radius: 0.4rem; cursor: pointer; font-size: 0.8rem; }
  .code-box { background: #0d1117; }
  .code-grid { display: grid; grid-template-columns: 45px 1fr; }
  .ln-col { border-right: 1px solid rgba(255,255,255,0.05); padding: 1rem 0; text-align: right; }
  .ln-col span { display: block; font-family: monospace; font-size: 0.75rem; color: var(--text-muted); opacity: 0.3; padding-right: 0.75rem; line-height: 1.5; }
  .code-pre { margin: 0; padding: 1rem; font-family: 'JetBrains Mono', monospace; font-size: 0.85rem; line-height: 1.5; }

  /* Readme */
  .readme-section { margin-top: 1.5rem; border: 1px solid var(--border-color); border-radius: 0.5rem; overflow: hidden; }
  .readme-header { padding: 0.75rem 1.25rem; background: var(--md3-surface-2); border-bottom: 1px solid var(--border-color); display: flex; align-items: center; gap: 0.75rem; font-weight: 700; font-size: 0.9rem; }
  .readme-content { padding: 2rem; }

  /* Sidebar */
  .side-column { display: flex; flex-direction: column; gap: 2.5rem; }
  .side-title { font-size: 0.8rem; font-weight: 800; text-transform: uppercase; letter-spacing: 0.05em; color: var(--text-muted); margin-bottom: 0.75rem; }
  .about-text { font-size: 0.95rem; line-height: 1.5; color: var(--text-color); }
  .lang-bar { height: 8px; border-radius: 4px; overflow: hidden; display: flex; background: var(--md3-surface-3); margin-bottom: 1rem; }
  .lang-item { display: flex; align-items: center; gap: 0.75rem; font-size: 0.85rem; margin-bottom: 0.5rem; }
  .lang-item .dot { width: 8px; height: 8px; border-radius: 50%; }
  .lang-item .name { flex: 1; font-weight: 600; }
</style>
