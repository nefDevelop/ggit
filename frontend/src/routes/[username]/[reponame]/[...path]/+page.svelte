<script lang="ts">
  import { page } from '$app/state';
  import { onMount, getContext } from 'svelte';
  import { authState } from '$lib/auth.svelte';

  let username = $derived(page.params.username);
  let reponame = $derived(page.params.reponame);
  let path = $derived(page.params.path || '');

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
        error = 'Repository not found or access denied.';
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
</script>

<div class="repo-browser-view">
  {#if loading}
    <div class="skeleton-content"></div>
  {:else if error}
    <div class="island error-island">
      <span class="material-symbols-rounded">error</span>
      <p>{error}</p>
      <a href="/{username}/{reponame}" class="btn btn-tonal">Go home</a>
    </div>
  {:else if repoData}
    <!-- Navigation / Breadcrumbs -->
    <div class="breadcrumbs-row">
      <nav class="breadcrumbs">
        <a href="/{username}/{reponame}" class="crumb-item">
          <span class="material-symbols-rounded">home</span>
        </a>
        {#each path.split('/').filter(p => p) as part, i}
          <span class="separator">/</span>
          <a href="/{username}/{reponame}/{path.split('/').slice(0, i+1).join('/')}" class="crumb-item">
            {part}
          </a>
        {/each}
      </nav>
      
      <div class="branch-selector-container">
        <select 
          class="branch-selector" 
          value={repoData.branch}
          onchange={(e) => {
            const newBranch = e.currentTarget.value;
            window.location.href = `/${username}/${reponame}/${newBranch}`;
          }}
        >
          {#each getContext<any>('repo').branches as b}
            <option value={b}>{b}</option>
          {/each}
        </select>
        <span class="material-symbols-rounded branch-icon">account_tree</span>
      </div>
    </div>

    <!-- Explorer Island -->
    {#if repoData.type === 'directory'}
      <div class="explorer-island island">
        <table class="file-table">
          <thead>
            <tr>
              <th>Name</th>
              <th>Message</th>
              <th>Size</th>
            </tr>
          </thead>
          <tbody>
            {#each repoData.entries as entry}
              <tr class="file-row">
                <td class="name-cell">
                  <span class="material-symbols-rounded icon-{entry.type === 'tree' ? 'folder' : 'file'}">
                    {entry.type === 'tree' ? 'folder_open' : 'article'}
                  </span>
                  <a href="/{username}/{reponame}/{path ? path + '/' : ''}{entry.name}">
                    {entry.name}
                  </a>
                </td>
                <td class="msg-cell text-muted">Initial commit</td>
                <td class="size-cell text-muted">
                  {entry.type === 'tree' ? '-' : formatSize(entry.size)}
                </td>
              </tr>
            {/each}
          </tbody>
        </table>
      </div>
    {:else if repoData.type === 'file'}
      <div class="file-viewer-island island">
        <div class="file-header">
          <div class="file-info">
            <span class="material-symbols-rounded">description</span>
            <span class="file-name">{repoData.name}</span>
            <span class="file-size text-muted">({formatSize(repoData.size)})</span>
          </div>
          <div class="file-actions">
            <button class="icon-btn" title="Copy raw content">
              <span class="material-symbols-rounded">content_copy</span>
            </button>
            <a href="/api/v1/repos/{username}/{reponame}/raw/{path}" class="icon-btn" title="Download raw" download>
              <span class="material-symbols-rounded">download</span>
            </a>
          </div>
        </div>
        
        <div class="code-container">
          {#if repoData.is_text_file}
            <div class="code-viewer">
              <div class="line-numbers">
                {#each repoData.highlighted_content as _, i}
                  <div class="line-number">{i + 1}</div>
                {/each}
              </div>
              <pre class="chroma">
                {#each repoData.highlighted_content as line}
                  <div class="code-line">{@html line || '&nbsp;'}</div>
                {/each}
              </pre>
            </div>
          {:else}
            <div class="binary-view">
              <span class="material-symbols-rounded">view_in_ar</span>
              <p>Binary file content cannot be displayed.</p>
              <a href="/api/v1/repos/{username}/{reponame}/raw/{path}" class="btn btn-tonal">Download File</a>
            </div>
          {/if}
        </div>
      </div>
    {/if}
  {/if}
</div>

<style>
  .repo-browser-view {
    display: flex;
    flex-direction: column;
    gap: 1.5rem;
  }

  /* Breadcrumbs */
  .breadcrumbs-row {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 0 0.5rem;
  }

  .breadcrumbs {
    display: flex;
    align-items: center;
    gap: 0.5rem;
  }

  .crumb-item {
    color: var(--primary-color);
    text-decoration: none;
    font-weight: 500;
    display: flex;
    align-items: center;
  }

  .crumb-item:hover {
    text-decoration: underline;
  }

  .separator {
    color: var(--text-muted);
    opacity: 0.5;
  }

  .branch-selector-container {
    position: relative;
    display: flex;
    align-items: center;
  }

  .branch-selector {
    appearance: none;
    background: var(--md3-surface-2);
    padding: 0.5rem 1rem 0.5rem 2.5rem;
    border-radius: 2rem;
    font-size: 0.9rem;
    font-weight: 600;
    color: var(--text-color);
    border: 1px solid var(--border-color);
    cursor: pointer;
    transition: all 0.2s;
    min-width: 140px;
  }

  .branch-selector:hover {
    background: var(--md3-surface-3);
    border-color: var(--primary-color);
  }

  .branch-selector:focus {
    outline: none;
    border-color: var(--primary-color);
    box-shadow: 0 0 0 2px rgba(var(--primary-rgb), 0.2);
  }

  .branch-icon {
    position: absolute;
    left: 1rem;
    pointer-events: none;
    font-size: 1.25rem;
    color: var(--primary-color);
  }

  /* Explorer Table */
  .explorer-island {
    padding: 0.5rem;
  }

  .file-table {
    width: 100%;
    border-collapse: collapse;
  }

  .file-table th {
    text-align: left;
    padding: 1rem;
    font-size: 0.8rem;
    text-transform: uppercase;
    color: var(--text-muted);
    border-bottom: 1px solid var(--border-color);
  }

  .file-row {
    transition: background 0.2s;
  }

  .file-row:hover {
    background-color: rgba(255, 255, 255, 0.03);
  }

  .file-row td {
    padding: 0.85rem 1rem;
  }

  .name-cell {
    display: flex;
    align-items: center;
    gap: 1rem;
  }

  .name-cell a {
    color: var(--text-color);
    text-decoration: none;
    font-weight: 500;
  }

  .name-cell a:hover {
    color: var(--primary-color);
  }

  .icon-folder { color: #facc15; }
  .icon-file { color: var(--primary-color); }

  /* File Viewer */
  .file-viewer-island {
    padding: 0;
    overflow: hidden;
  }

  .file-header {
    background-color: rgba(255, 255, 255, 0.02);
    padding: 1rem 1.5rem;
    border-bottom: 1px solid var(--border-color);
    display: flex;
    justify-content: space-between;
    align-items: center;
  }

  .file-info {
    display: flex;
    align-items: center;
    gap: 0.75rem;
  }

  .file-name {
    font-weight: 700;
  }

  .code-viewer {
    display: flex;
    background-color: #0d1117; /* High-contrast code background */
    min-height: 100%;
  }

  .line-numbers {
    padding: 1.5rem 1rem;
    text-align: right;
    background-color: rgba(255, 255, 255, 0.02);
    border-right: 1px solid var(--border-color);
    user-select: none;
    min-width: 3.5rem;
  }

  .line-number {
    font-family: 'JetBrains Mono', monospace;
    font-size: 0.85rem;
    color: var(--text-muted);
    opacity: 0.5;
    line-height: 1.5;
    height: 1.5rem;
  }

  .code-container {
    overflow-x: auto;
    padding: 0;
  }

  .chroma {
    margin: 0;
    padding: 1.5rem;
    font-family: 'JetBrains Mono', 'Fira Code', monospace;
    font-size: 0.9rem;
    line-height: 1.5;
    flex: 1;
    overflow-x: auto;
  }

  .code-line {
    height: 1.5rem;
    white-space: pre;
  }
  
  .binary-view {
    padding: 4rem;
    text-align: center;
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 1rem;
  }

  .binary-view .material-symbols-rounded {
    font-size: 4rem;
    opacity: 0.3;
  }

  .skeleton-content {
    height: 400px;
    background: var(--md3-surface-2);
    border-radius: var(--radius);
    animation: pulse 1.5s infinite;
  }

  @keyframes pulse {
    0% { opacity: 0.6; }
    50% { opacity: 0.3; }
    100% { opacity: 0.6; }
  }

  .error-island {
    text-align: center;
    color: #f87171;
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 1rem;
    padding: 4rem;
  }
</style>
