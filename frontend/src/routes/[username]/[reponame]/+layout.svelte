<script lang="ts">
  import { page } from '$app/state';
  import { onMount, setContext } from 'svelte';
  
  let { children } = $props();
  
  let username = $derived(page.params.username);
  let reponame = $derived(page.params.reponame);
  let currentPath = $derived(page.url.pathname);

  let repoData = $state<any>(null);
  let branches = $state<any[]>([]);
  let loading = $state(true);
  let error = $state('');

  async function fetchRepoInfo() {
    loading = true;
    error = '';
    try {
      const [repoRes, branchesRes] = await Promise.all([
        fetch(`/api/v2/${username}/${reponame}/view/`),
        fetch(`/api/v2/${username}/${reponame}/branches`)
      ]);
      
      if (repoRes.ok) {
        const data = await repoRes.json();
        repoData = data.repository || data; 
      } else {
        error = 'Repository not found or access denied.';
      }

      if (branchesRes.ok) {
        branches = await branchesRes.json();
      }
    } catch (e) {
      error = 'Failed to load repository information.';
    } finally {
      loading = false;
    }
  }

  onMount(() => {
    fetchRepoInfo();
  });

  setContext('repo', {
    get data() { return repoData; },
    get branches() { return branches; },
    get loading() { return loading; },
    get error() { return error; }
  });

  let activeTab = $derived(
    currentPath.includes('/commits') ? 'commits' :
    currentPath.includes('/issues') ? 'issues' :
    currentPath.includes('/pulls') ? 'pulls' : 'code'
  );

  let showCloneDropdown = $state(false);
  function copyCloneUrl(url: string) {
    navigator.clipboard.writeText(url);
    showCloneDropdown = false;
  }
</script>

<div class="workspace-shell">
  {#if loading}
    <div class="workspace-container">
      <div class="skeleton-island island"></div>
    </div>
  {:else if error}
    <div class="workspace-container">
      <div class="island error-island">
        <span class="material-symbols-rounded">error</span>
        <p>{error}</p>
        <a href="/dashboard" class="btn btn-tonal">Back to Dashboard</a>
      </div>
    </div>
  {:else if repoData}
    <div class="workspace-container">
      <!-- THE MASTER ISLAND -->
      <div class="master-island island">
        
        <!-- Technical Header (Integrated) -->
        <header class="tech-header">
          <div class="breadcrumb-nav">
            <a href="/dashboard" class="bc-item home-icon">
              <span class="material-symbols-rounded">grid_view</span>
            </a>
            <span class="bc-sep">/</span>
            <a href="/{username}" class="bc-item">{username}</a>
            <span class="bc-sep">/</span>
            <h1 class="bc-item active">{reponame}</h1>
            {#if repoData.private}
              <span class="lock-icon material-symbols-rounded">lock</span>
            {/if}
          </div>

          <div class="header-actions">
            <div class="clone-zone">
              <button class="btn btn-sm btn-primary clone-btn" onclick={() => showCloneDropdown = !showCloneDropdown}>
                <span class="material-symbols-rounded">download</span>
                Clone
              </button>
              {#if showCloneDropdown}
                <div class="clone-dropdown island shadow-lg">
                  <div class="clone-option">
                    <span class="label">HTTPS</span>
                    <code class="url">{repoData.clone_url}</code>
                    <button class="copy-icon" onclick={() => copyCloneUrl(repoData.clone_url)}>
                      <span class="material-symbols-rounded">content_copy</span>
                    </button>
                  </div>
                </div>
              {/if}
            </div>
            <a href="/{username}/{reponame}/settings" class="btn-icon" title="Settings">
              <span class="material-symbols-rounded">settings</span>
            </a>
          </div>
        </header>

        <!-- Technical Tabs (Integrated) -->
        <nav class="tech-tabs">
          <a href="/{username}/{reponame}" class="tab" class:active={activeTab === 'code'}>
            <span class="material-symbols-rounded">code</span>
            Code
          </a>
          <a href="/{username}/{reponame}/commits" class="tab" class:active={activeTab === 'commits'}>
            <span class="material-symbols-rounded">history</span>
            Commits
          </a>
          <a href="/{username}/{reponame}/issues" class="tab" class:active={activeTab === 'issues'}>
            <span class="material-symbols-rounded">error</span>
            Issues
            {#if repoData.open_issues_count > 0}
              <span class="badge">{repoData.open_issues_count}</span>
            {/if}
          </a>
          <a href="/{username}/{reponame}/pulls" class="tab" class:active={activeTab === 'pulls'}>
            <span class="material-symbols-rounded">merge</span>
            PRs
          </a>
        </nav>

        <!-- Content Area (Children render here) -->
        <div class="master-content">
          {@render children()}
        </div>
      </div>
    </div>
  {/if}
</div>

<style>
  .workspace-shell {
    height: 100%;
    width: 100%;
    padding: 0;
  }

  .workspace-container {
    height: 100%;
    display: flex;
    flex-direction: column;
  }

  .master-island {
    flex: 1;
    display: flex;
    flex-direction: column;
    background-color: var(--md3-surface-container);
    border-radius: var(--radius);
    overflow: hidden;
    /* We don't want internal scrolls if the page can handle it, 
       but for sticky sidebars we need a good layout */
  }

  .tech-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 1.25rem 1.5rem;
    background: rgba(255,255,255,0.02);
    border-bottom: 1px solid var(--border-color);
  }

  .breadcrumb-nav { display: flex; align-items: center; gap: 0.5rem; font-size: 1.1rem; }
  .bc-item { color: var(--text-muted); text-decoration: none; font-weight: 500; transition: color 0.2s; }
  .bc-item.active { color: var(--text-color); font-weight: 800; font-size: 1.2rem; }
  .bc-sep { color: var(--text-muted); opacity: 0.3; }
  .home-icon span { font-size: 1.2rem; }
  .lock-icon { font-size: 1rem; color: var(--primary-color); opacity: 0.8; }

  .header-actions { display: flex; align-items: center; gap: 1rem; }
  .clone-zone { position: relative; }
  .clone-dropdown {
    position: absolute;
    top: calc(100% + 0.5rem);
    right: 0;
    width: 350px;
    z-index: 100;
    padding: 1rem;
    background: var(--md3-surface-2);
  }

  .clone-option {
    display: flex;
    align-items: center;
    gap: 0.75rem;
    background: var(--md3-surface-container);
    padding: 0.5rem;
    border-radius: 0.5rem;
    border: 1px solid var(--border-color);
  }
  .clone-option code { flex: 1; font-size: 0.8rem; overflow: hidden; text-overflow: ellipsis; opacity: 0.8; }

  .btn-icon { color: var(--text-muted); display: flex; padding: 0.5rem; border-radius: 50%; }
  .btn-icon:hover { background: var(--md3-surface-3); color: var(--primary-color); }

  /* Integrated Tabs */
  .tech-tabs {
    display: flex;
    padding: 0 1rem;
    background: rgba(255,255,255,0.01);
    border-bottom: 1px solid var(--border-color);
  }

  .tab {
    display: flex;
    align-items: center;
    gap: 0.6rem;
    padding: 1rem 1.25rem;
    color: var(--text-muted);
    text-decoration: none;
    font-size: 0.9rem;
    font-weight: 600;
    position: relative;
  }

  .tab:hover { color: var(--text-color); }
  .tab.active { color: var(--primary-color); }
  .tab.active::after {
    content: '';
    position: absolute;
    bottom: -1px;
    left: 1.25rem;
    right: 1.25rem;
    height: 3px;
    background: var(--primary-color);
    border-radius: 4px 4px 0 0;
  }

  .badge { background: var(--primary-color); color: var(--on-primary); font-size: 0.7rem; padding: 0.1rem 0.4rem; border-radius: 1rem; margin-left: 0.25rem; }

  .master-content {
    flex: 1;
    padding: 1.5rem;
    overflow-y: auto;
  }

  .skeleton-island { height: 100%; width: 100%; animation: pulse 1.5s infinite; }
  @keyframes pulse { 0%, 100% { opacity: 0.6; } 50% { opacity: 0.3; } }

  .error-island { height: 100%; display: flex; flex-direction: column; align-items: center; justify-content: center; gap: 1rem; }
</style>
