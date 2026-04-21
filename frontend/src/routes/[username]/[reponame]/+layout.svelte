<script lang="ts">
  import { apiFetch } from "$lib/api";
  import { page } from '$app/state';
  import { onMount, setContext } from 'svelte';
  import GlobalSearch from '$lib/components/GlobalSearch.svelte';
  
  let { children } = $props();
  
  let username = $derived(page.params.username);
  let reponame = $derived(page.params.reponame);
  let currentPath = $derived(page.url.pathname);

  let repoData = $state<any>(null);
  let gitStats = $state<any>({ commits: 0, branches: 0, tags: 0 });
  let loading = $state(true);
  let error = $state('');

  async function fetchRepoInfo() {
    loading = true;
    error = '';
    try {
      // We use the view endpoint for the root to get the repository object AND stats
      const repoRes = await apiFetch(`/api/v2/${username}/${reponame}/view/`);
      
      if (repoRes.ok) {
        const res = await repoRes.json();
        repoData = res.repository;
        if (res.stats) {
            gitStats = res.stats;
        }
      } else {
        error = 'Repository not found or access denied.';
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
    get stats() { return gitStats; },
    get loading() { return loading; }
  });

  let activeTab = $derived(
    currentPath.includes('/commits') ? 'commits' :
    currentPath.includes('/issues') ? 'issues' :
    currentPath.includes('/pulls') ? 'pulls' :
    currentPath.includes('/wiki') ? 'wiki' :
    currentPath.includes('/branches') ? 'branches' : 'code'
  );

  let showCloneDropdown = $state(false);
  function copyCloneUrl(url: string) {
    navigator.clipboard.writeText(url);
    showCloneDropdown = false;
  }
</script>

<div class="workspace-shell">
  {#if loading}
    <div class="skeleton-master island"></div>
  {:else if error}
    <div class="island error-island">
      <span class="material-symbols-rounded">error</span>
      <p>{error}</p>
      <a href="/dashboard" class="btn btn-tonal">Back to Dashboard</a>
    </div>
  {:else if repoData}
    <div class="master-island island">
      <!-- HEADER STRATOSPHERE -->
      <div class="repo-header-region">
        <!-- ROW 1: Integrated Search -->
        <div class="island-search-bar">
          <GlobalSearch />
        </div>

        <!-- ROW 2: Identity & Meta -->
        <header class="repo-identity-bar">
          <div class="breadcrumb-nav">
            <a href="/dashboard" class="bc-home"><span class="material-symbols-rounded">grid_view</span></a>
            <span class="bc-sep">/</span>
            <a href="/{username}" class="bc-item">{username}</a>
            <span class="bc-sep">/</span>
            <h1 class="bc-item active">{reponame}</h1>
            {#if repoData.private}<span class="lock-icon material-symbols-rounded">lock</span>{/if}
          </div>

          <div class="repo-meta-actions">
            {#if repoData.license}
              <div class="license-tag">
                <span class="material-symbols-rounded">balance</span>
                {repoData.license}
              </div>
            {/if}
            
            <div class="clone-zone">
              <button class="btn btn-sm btn-primary" onclick={() => showCloneDropdown = !showCloneDropdown}>
                <span class="material-symbols-rounded">download</span>
                Clone
              </button>
              {#if showCloneDropdown}
                <div class="clone-dropdown island shadow-lg">
                  <div class="clone-option">
                    <span class="label">HTTPS</span>
                    <code class="url">{repoData.clone_url}</code>
                    <button class="copy-btn" onclick={() => copyCloneUrl(repoData.clone_url)}>
                      <span class="material-symbols-rounded">content_copy</span>
                    </button>
                    </div>
                    <div class="archive-section">                  <span class="dropdown-header">Download Source</span>
                    <div class="archive-grid">
                      <a href="/api/v2/{username}/{reponame}/archive/{repoData.default_branch}.zip" class="archive-btn tonal">
                        <span class="material-symbols-rounded">zip_box</span>
                        <span>ZIP</span>
                      </a>
                      <a href="/api/v2/{username}/{reponame}/archive/{repoData.default_branch}.tar.gz" class="archive-btn tonal">
                        <span class="material-symbols-rounded">package</span>
                        <span>TAR.GZ</span>
                      </a>
                    </div>
                  </div>
                </div>
              {/if}
            </div>
            <a href="/{username}/{reponame}/settings" class="settings-link" title="Settings">
              <span class="material-symbols-rounded">settings</span>
            </a>
          </div>
        </header>

        <!-- ROW 3: Rich Technical Tabs (Stats integrated) -->
        <nav class="rich-tabs">
          <a href="/{username}/{reponame}" class="tab" class:active={activeTab === 'code'}>
            <span class="material-symbols-rounded">code</span>
            {#if activeTab === 'code'}<span class="tab-label">Code</span>{/if}
          </a>
          <a href="/{username}/{reponame}/commits" class="tab" class:active={activeTab === 'commits'}>
            <span class="material-symbols-rounded">history</span>
            <span class="tab-count">{gitStats.commits}</span>
            {#if activeTab === 'commits'}<span class="tab-label">Commits</span>{/if}
          </a>
          <a href="/{username}/{reponame}/branches" class="tab" class:active={activeTab === 'branches'}>
            <span class="material-symbols-rounded">account_tree</span>
            <span class="tab-count">{gitStats.branches}</span>
            {#if activeTab === 'branches'}<span class="tab-label">Branches</span>{/if}
          </a>
          <a href="/{username}/{reponame}/issues" class="tab" class:active={activeTab === 'issues'}>
            <span class="material-symbols-rounded">error</span>
            <span class="tab-count">{repoData.open_issues_count || 0}</span>
            {#if activeTab === 'issues'}<span class="tab-label">Issues</span>{/if}
          </a>
          <a href="/{username}/{reponame}/pulls" class="tab" class:active={activeTab === 'pulls'}>
            <span class="material-symbols-rounded">merge</span>
            <span class="tab-count">{repoData.open_pulls_count || 0}</span>
            {#if activeTab === 'pulls'}<span class="tab-label">PRs</span>{/if}
          </a>
          <a href="/{username}/{reponame}/wiki" class="tab" class:active={activeTab === 'wiki'}>
            <span class="material-symbols-rounded">book</span>
            {#if activeTab === 'wiki'}<span class="tab-label">Wiki</span>{/if}
          </a>
        </nav>
      </div>

      <!-- MAIN CONTENT -->
      <div class="master-body">
        {@render children()}
      </div>
    </div>
  {/if}
</div>

<style>
  .workspace-shell { height: 100%; width: 100%; display: flex; flex-direction: column; }
  .master-island { flex: 1; display: flex; flex-direction: column; overflow: hidden; }

  /* Row 1: Search */
  .island-search-bar {
    padding: 1.5rem 1.5rem 0.75rem;
    background: rgba(255,255,255,0.01);
  }

  /* Row 2: Identity */
  .repo-identity-bar {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 0.75rem 1.5rem;
  }
  .breadcrumb-nav { display: flex; align-items: center; gap: 0.5rem; }
  .bc-item { color: var(--text-muted); text-decoration: none; font-weight: 500; font-size: 1.1rem; }
  .bc-item.active { color: var(--text-color); font-weight: 800; font-size: 1.25rem; }
  .bc-home { color: var(--text-muted); display: flex; opacity: 0.6; }
  .bc-sep { opacity: 0.2; }
  .lock-icon { font-size: 1rem; color: var(--primary-color); margin-left: 0.5rem; }

  .repo-meta-actions { display: flex; align-items: center; gap: 1.5rem; }
  .license-tag { display: flex; align-items: center; gap: 0.4rem; font-size: 0.8rem; font-weight: 700; color: var(--text-muted); }
  .license-tag span { font-size: 1.1rem; }
  .settings-link { color: var(--text-muted); display: flex; padding: 4px; transition: color 0.2s; }
  .settings-link:hover { color: var(--primary-color); }

  .clone-zone { position: relative; }
  .clone-dropdown { position: absolute; top: calc(100% + 0.5rem); right: 0; width: 320px; z-index: 100; padding: 1rem; background: var(--md3-surface-2); border-radius: 1rem; border: 1px solid var(--border-color); }
  .clone-option { display: flex; align-items: center; gap: 0.75rem; background: rgba(0,0,0,0.2); padding: 0.5rem; border-radius: 0.5rem; }
  .clone-option code { flex: 1; font-size: 0.8rem; overflow: hidden; text-overflow: ellipsis; opacity: 0.7; }
  .copy-btn { background: transparent; border: none; color: var(--text-muted); cursor: pointer; display: flex; }

  .archive-section { margin-top: 1.25rem; padding-top: 1rem; border-top: 1px solid rgba(255,255,255,0.05); }
  .dropdown-header { display: block; font-size: 0.7rem; font-weight: 800; text-transform: uppercase; color: var(--text-muted); margin-bottom: 0.75rem; letter-spacing: 0.05em; }
  .archive-grid { display: grid; grid-template-columns: 1fr 1fr; gap: 0.75rem; }
  .archive-btn { display: flex; align-items: center; justify-content: center; gap: 0.5rem; padding: 0.6rem; border-radius: 0.75rem; text-decoration: none; color: var(--text-color); font-size: 0.85rem; font-weight: 600; transition: all 0.2s; }
  .archive-btn:hover { background: rgba(var(--primary-rgb), 0.1); color: var(--primary-color); transform: translateY(-2px); }
  .archive-btn span { font-size: 1.25rem; }

  /* Row 3: Rich Tabs */
  .repo-header-region {
    background-color: var(--surface-variant);
    border-bottom: 1px solid var(--border-color);
  }

  .rich-tabs {
    display: flex;
    gap: 0.25rem;
    padding: 0 1rem;
  }
  .tab {
    display: flex;
    align-items: center;
    gap: 0.6rem;
    padding: 0.85rem 1.25rem;
    color: var(--text-muted);
    text-decoration: none;
    font-size: 0.85rem;
    font-weight: 700;
    transition: all 0.2s;
    border-radius: 0.75rem 0.75rem 0 0;
    margin-bottom: -1px; /* Align with bottom border */
  }
  .tab:hover { color: var(--text-color); background: rgba(0,0,0,0.03); }
  .tab.active { color: var(--primary-color); background: var(--bg-color); position: relative; border: 1px solid var(--border-color); border-bottom-color: var(--bg-color); }
  
  .tab-count { font-variant-numeric: tabular-nums; background: rgba(0,0,0,0.05); padding: 0.1rem 0.4rem; border-radius: 4px; font-size: 0.75rem; color: var(--text-muted); }
  .tab.active .tab-count { background: var(--primary-container); color: var(--on-primary-container); }
  .tab-label { text-transform: uppercase; letter-spacing: 0.05em; font-size: 0.75rem; }

  /* Responsive Tabs: Only show label for the active tab when space is tight */
  @media (max-width: 800px) {
    .tab:not(.active) .tab-label { display: none; }
    .tab:not(.active) .tab-count { display: none; }
    .tab { padding: 0.85rem 1rem; }
  }

  .master-body { flex: 1; overflow-y: auto; padding: 2rem; background: var(--bg-color); }

  .skeleton-master { flex: 1; animation: pulse 1.5s infinite; }
  @keyframes pulse { 0%, 100% { opacity: 0.6; } 50% { opacity: 0.3; } }

  .error-island { height: 100%; display: flex; flex-direction: column; align-items: center; justify-content: center; gap: 1rem; color: #f87171; }
</style>
