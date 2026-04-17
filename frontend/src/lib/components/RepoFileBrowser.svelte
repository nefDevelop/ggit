<script lang="ts">
  import { apiFetch } from "$lib/api";
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
      const res = await apiFetch(`/api/v2/${username}/${reponame}/view/${path}`);
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

  // Calculate language percentages
  const langEntries = $derived(() => {
    if (!repoData?.languages) return [];
    const total = Object.values(repoData.languages as Record<string, number>).reduce((a, b) => a + b, 0);
    return Object.entries(repoData.languages as Record<string, number>)
      .map(([name, bytes]) => ({
        name,
        percent: ((bytes / total) * 100).toFixed(1),
        color: getLanguageColor(name)
      }))
      .sort((a, b) => parseFloat(b.percent) - parseFloat(a.percent));
  });

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
    {#if repoData.is_bare}
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
      <div class="technical-grid">
        <!-- Main Explorer / Code -->
        <div class="main-column">
          <!-- Last Commit Bar -->
          {#if repoData.last_commit}
            <div class="last-commit-bar">
                <div class="commit-avatar">{repoData.last_commit.author?.name?.[0] || 'U'}</div>
                <div class="commit-info">
                <span class="commit-msg">{repoData.last_commit.summary}</span>
                <span class="commit-meta">
                    <span class="author">{repoData.last_commit.author?.name}</span>
                    <span class="time">{timeAgo(repoData.last_commit.author?.when)}</span>
                </span>
                </div>
                <code class="sha">{repoData.last_commit.id?.slice(0,7)}</code>
            </div>
          {/if}

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
                        <a href="/{username}/{reponame}/src/{repoData.branch}/{path ? path + '/' : ''}{entry.name}">
                          {entry.name}
                        </a>
                      </td>
                      <td class="msg-cell">
                        <span class="msg-text">...</span>
                      </td>
                      <td class="date-cell">--</td>
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
                      {#if repoData.is_text_file}
                        <a href="/{username}/{reponame}/_edit/{repoData.branch}/{path}" class="tool-btn">
                          <span class="material-symbols-rounded">edit</span>
                        </a>
                      {/if}
                      <button class="tool-btn">Raw</button>
                      <button class="tool-btn"><span class="material-symbols-rounded">download</span></button>
                    </div>
                </div>
                <div class="code-box">
                   {#if repoData.is_text_file}
                      <div class="code-grid">
                        <div class="ln-col">
                          {#each repoData.highlighted_content || [] as _, i}
                            <span>{i + 1}</span>
                          {/each}
                        </div>
                        <pre class="code-pre">
                          {#each repoData.highlighted_content || [] as line}
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

                {#if repoData.readme_rendered}
                   <div class="file-readme-integrated">
                      <header class="readme-header small">
                         <span class="material-symbols-rounded">menu_book</span>
                         README Preview
                      </header>
                      <div class="readme-content markdown-body">
                         {@html repoData.readme_rendered}
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
                <div class="readme-content markdown-body">
                    {@html repoData.readme.content}
                </div>
            </div>
          {/if}
        </div>

        <!-- Sidebar -->
        <aside class="side-column">
          {#if langEntries().length > 0}
            <div class="side-section">
                <h4 class="side-title">Languages</h4>
                <div class="lang-bar">
                    {#each langEntries() as lang}
                        <div class="lang-progress" style="width: {lang.percent}%; background: {lang.color}"></div>
                    {/each}
                </div>
                <div class="lang-list">
                    {#each langEntries() as lang}
                        <div class="lang-item">
                            <span class="dot" style="background: {lang.color}"></span>
                            <span class="name">{lang.name}</span>
                            <span class="percent">{lang.percent}%</span>
                        </div>
                    {/each}
                </div>
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
  }

  .main-column { display: flex; flex-direction: column; min-width: 0; }

  .last-commit-bar {
    display: flex;
    align-items: center;
    gap: 1rem;
    padding: 0.75rem 1rem;
    background: rgba(255,255,255,0.03);
    border-radius: 0.5rem;
    margin-bottom: 1.5rem;
  }
  .commit-avatar { width: 24px; height: 24px; background: var(--primary-color); border-radius: 4px; display: flex; align-items: center; justify-content: center; font-weight: 800; font-size: 0.7rem; color: var(--on-primary); }
  .commit-msg { font-weight: 600; font-size: 0.95rem; flex: 1; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
  .commit-meta { display: flex; gap: 0.5rem; font-size: 0.8rem; color: var(--text-muted); }
  .sha { font-family: monospace; opacity: 0.5; font-size: 0.8rem; }

  /* Content */
  .content-view { background: rgba(0,0,0,0.1); border-radius: 0.5rem; overflow: hidden; }
  .file-table { width: 100%; border-collapse: collapse; }
  .file-row { border-bottom: 1px solid rgba(255,255,255,0.02); }
  .file-row td { padding: 0.6rem 1rem; }
  .name-cell { display: flex; align-items: center; gap: 0.75rem; min-width: 250px; }
  .name-cell a { color: var(--text-color); text-decoration: none; font-weight: 500; }
  .icon-tree { color: #facc15; font-size: 1.2rem; }
  .icon-blob { color: var(--text-muted); font-size: 1.2rem; }
  .msg-cell { font-size: 0.85rem; color: var(--text-muted); opacity: 0.6; }
  .date-cell { text-align: right; font-size: 0.8rem; color: var(--text-muted); opacity: 0.4; }

  /* Viewer */
  .file-header { background: var(--md3-surface-3); padding: 0.5rem 1rem; display: flex; justify-content: space-between; align-items: center; }
  .file-meta { display: flex; align-items: center; gap: 0.5rem; font-weight: 600; font-size: 0.9rem; }
  .tool-btn { background: transparent; border: 1px solid var(--border-color); color: var(--text-muted); padding: 0.2rem 0.6rem; border-radius: 0.4rem; cursor: pointer; font-size: 0.8rem; }
  .code-box { background: #0d1117; }
  .code-grid { display: grid; grid-template-columns: 45px 1fr; }
  .ln-col { border-right: 1px solid rgba(255,255,255,0.05); padding: 1rem 0; text-align: right; }
  .ln-col span { display: block; font-family: monospace; font-size: 0.75rem; color: var(--text-muted); opacity: 0.3; padding-right: 0.75rem; line-height: 1.5; }
  .code-pre { margin: 0; padding: 1rem; font-family: 'JetBrains Mono', monospace; font-size: 0.85rem; line-height: 1.5; color: #e6edf3; }

  /* Readme */
  .readme-section { margin-top: 1.5rem; background: rgba(255,255,255,0.01); border-radius: 0.5rem; overflow: hidden; }
  .readme-header { padding: 0.75rem 1.25rem; border-bottom: 1px solid var(--border-color); display: flex; align-items: center; gap: 0.75rem; font-weight: 700; font-size: 0.85rem; text-transform: uppercase; color: var(--text-muted); }
  .readme-header.small { padding: 0.5rem 1rem; background: var(--md3-surface-3); border-top: 1px solid var(--border-color); }
  .readme-content { padding: 2rem; }

  /* Sidebar */
  .side-column { display: flex; flex-direction: column; gap: 2.5rem; }
  .side-title { font-size: 0.75rem; font-weight: 800; text-transform: uppercase; letter-spacing: 0.05em; color: var(--text-muted); margin-bottom: 0.75rem; }
  .about-text { font-size: 0.9rem; line-height: 1.5; color: var(--text-color); opacity: 0.8; }
  .lang-bar { height: 8px; border-radius: 4px; overflow: hidden; display: flex; background: rgba(0,0,0,0.3); margin-bottom: 1rem; }
  .lang-item { display: flex; align-items: center; gap: 0.75rem; font-size: 0.8rem; margin-bottom: 0.4rem; }
  .lang-item .dot { width: 8px; height: 8px; border-radius: 50%; }
  .lang-item .name { flex: 1; font-weight: 600; }
  .lang-item .percent { color: var(--text-muted); }

  /* Skeletons */
  .skeleton-grid { display: grid; grid-template-columns: 1fr 280px; gap: 2.5rem; height: 400px; }
  .skeleton-main { background: var(--md3-surface-2); border-radius: 1rem; animation: pulse 1.5s infinite; }
  .skeleton-side { background: var(--md3-surface-2); border-radius: 1rem; animation: pulse 1.5s infinite; }
  @keyframes pulse { 0%, 100% { opacity: 0.5; } 50% { opacity: 0.8; } }

  .empty-state { padding: 4rem 2rem; text-align: center; }
  .empty-icon-ring { width: 64px; height: 64px; border: 2px dashed var(--primary-color); border-radius: 50%; display: flex; align-items: center; justify-content: center; margin: 0 auto 1.5rem; color: var(--primary-color); }
  .command-block { background: #000; padding: 1.5rem; border-radius: 0.75rem; text-align: left; margin-top: 2rem; }
  .cmd-label { font-size: 0.75rem; font-weight: 800; color: var(--primary-color); display: block; margin-bottom: 0.5rem; }
  code { font-family: 'JetBrains Mono', monospace; }

  /* Markdown Body Overrides */
  :global(.markdown-body) { color: var(--text-color); }
  :global(.markdown-body h1, .markdown-body h2) { border-bottom: 1px solid var(--border-color); padding-bottom: 0.3em; }
  :global(.markdown-body pre) { background-color: rgba(0,0,0,0.3); border-radius: 6px; padding: 16px; }
</style>
