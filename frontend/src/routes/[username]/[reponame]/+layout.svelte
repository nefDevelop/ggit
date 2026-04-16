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
        repoData = data.repository;
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

  // Provide repo data and branches to children
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
</script>

<div class="repo-shell">
  {#if loading}
    <div class="repo-container">
      <div class="skeleton-header"></div>
      <div class="skeleton-content"></div>
    </div>
  {:else if error}
    <div class="repo-container">
      <div class="island error-island">
        <span class="material-symbols-rounded">error</span>
        <p>{error}</p>
        <a href="/dashboard" class="btn btn-tonal">Back to Dashboard</a>
      </div>
    </div>
  {:else if repoData}
    <div class="repo-container">
      <!-- Repo Header Island -->
      <header class="repo-header-island island">
        <div class="repo-title-row">
          <div class="repo-name-stack">
            <div class="repo-path">
              <a href="/{username}" class="owner-link">{username}</a> /
            </div>
            <h1 class="repo-name">{reponame}</h1>
            {#if repoData.is_private}
              <span class="badge private-badge">
                <span class="material-symbols-rounded">lock</span> Private
              </span>
            {/if}
          </div>
          
          <div class="repo-actions">
            <div class="stat-group">
              <div class="stat-item">
                <span class="material-symbols-rounded">star</span>
                <span>{repoData.num_stars || 0}</span>
              </div>
              <div class="stat-item">
                <span class="material-symbols-rounded">fork_right</span>
                <span>{repoData.num_forks || 0}</span>
              </div>
            </div>
            <button class="btn">
              <span class="material-symbols-rounded">download</span>
              Clone
            </button>
          </div>
        </div>
        
        {#if repoData.description}
          <p class="repo-description">{repoData.description}</p>
        {/if}

        <!-- Tabs Navigation -->
        <nav class="repo-tabs">
          <a href="/{username}/{reponame}" class="tab-item" class:active={activeTab === 'code'}>
            <span class="material-symbols-rounded">code</span>
            Code
          </a>
          <a href="/{username}/{reponame}/commits" class="tab-item" class:active={activeTab === 'commits'}>
            <span class="material-symbols-rounded">history</span>
            Commits
          </a>
          <a href="/{username}/{reponame}/issues" class="tab-item" class:active={activeTab === 'issues'}>
            <span class="material-symbols-rounded">error_outline</span>
            Issues
          </a>
          <a href="/{username}/{reponame}/pulls" class="tab-item" class:active={activeTab === 'pulls'}>
            <span class="material-symbols-rounded">merge</span>
            Pull Requests
          </a>
        </nav>
      </header>

      <main class="repo-content">
        {@render children()}
      </main>
    </div>
  {/if}
</div>

<style>
  .repo-shell {
    min-height: 100vh;
    padding: 2rem 0;
  }

  .repo-container {
    max-width: 1200px;
    margin: 0 auto;
    display: flex;
    flex-direction: column;
    gap: 1.5rem;
    padding: 0 1rem;
  }

  .repo-header-island {
    padding: 2rem 2rem 0;
    overflow: hidden;
  }

  .repo-title-row {
    display: flex;
    justify-content: space-between;
    align-items: flex-start;
    margin-bottom: 1.5rem;
  }

  .repo-name-stack {
    display: flex;
    flex-direction: column;
    gap: 0.25rem;
  }

  .repo-path {
    font-size: 1.1rem;
    color: var(--text-muted);
  }

  .owner-link {
    color: var(--primary-color);
    text-decoration: none;
    font-weight: 500;
  }

  .repo-name {
    font-size: 2.25rem;
    font-weight: 800;
    margin: 0;
    letter-spacing: -0.04em;
  }

  .private-badge {
    background-color: rgba(255, 255, 255, 0.05);
    padding: 0.25rem 0.75rem;
    border-radius: 2rem;
    font-size: 0.8rem;
    display: inline-flex;
    align-items: center;
    gap: 0.4rem;
    margin-top: 0.5rem;
    width: fit-content;
  }

  .repo-description {
    margin: 0.5rem 0 2rem;
    font-size: 1.1rem;
    color: var(--text-color);
    max-width: 800px;
    opacity: 0.8;
  }

  .repo-actions {
    display: flex;
    align-items: center;
    gap: 1.5rem;
  }

  .stat-group {
    display: flex;
    gap: 1rem;
  }

  .stat-item {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    background-color: rgba(255, 255, 255, 0.03);
    padding: 0.5rem 0.8rem;
    border-radius: 0.75rem;
    font-weight: 600;
  }

  /* Tabs Navigation */
  .repo-tabs {
    display: flex;
    gap: 2rem;
    margin-top: 1rem;
  }

  .tab-item {
    padding: 1rem 0;
    color: var(--text-muted);
    text-decoration: none;
    font-weight: 600;
    display: flex;
    align-items: center;
    gap: 0.6rem;
    position: relative;
    transition: color 0.2s;
  }

  .tab-item:hover {
    color: var(--text-color);
  }

  .tab-item.active {
    color: var(--primary-color);
  }

  .tab-item.active::after {
    content: '';
    position: absolute;
    bottom: 0;
    left: 0;
    right: 0;
    height: 3px;
    background-color: var(--primary-color);
    border-radius: 3px 3px 0 0;
  }

  .repo-content {
    display: flex;
    flex-direction: column;
    gap: 1.5rem;
  }

  /* Skeletons */
  .skeleton-header {
    height: 200px;
    background: var(--md3-surface-2);
    border-radius: var(--radius);
    animation: pulse 1.5s infinite;
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
