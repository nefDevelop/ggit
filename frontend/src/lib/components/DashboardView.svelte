<script lang="ts">
  import { apiFetch } from "$lib/api";
  import { authState } from '$lib/auth.svelte';
  import { onMount } from 'svelte';
  import RepoRow from './RepoRow.svelte';
  import GlobalSearch from './GlobalSearch.svelte';

  let data = $state<{
    repositories: any[],
    collaborative: any[],
    mirrors: any[],
    organizations: any[]
  }>({
    repositories: [],
    collaborative: [],
    mirrors: [],
    organizations: []
  });
  let feeds = $state<any[]>([]);
  let loading = $state(true);
  let viewMode = $state<'list' | 'temporal'>('temporal');
  let searchQuery = $state('');

  const oneWeekAgo = $derived(Date.now() - (7 * 24 * 60 * 60 * 1000));
  const allRepos = $derived([...data.repositories, ...data.collaborative]);
  
  const filteredRepos = $derived(() => {
    if (!searchQuery) return allRepos;
    const q = searchQuery.toLowerCase();
    return allRepos.filter(r => 
      r.name.toLowerCase().includes(q) || 
      (r.description && r.description.toLowerCase().includes(q))
    );
  });

  interface GroupedRepos {
    all?: any[];
    recent?: any[];
    older?: any[];
  }

  const groupedRepos = $derived((): GroupedRepos => {
    const repos = searchQuery ? filteredRepos() : allRepos;
    if (viewMode === 'list') return { all: repos };
    return {
      recent: repos.filter(r => new Date(r.updated_at).getTime() > oneWeekAgo).sort((a, b) => new Date(b.updated_at).getTime() - new Date(a.updated_at).getTime()),
      older: repos.filter(r => new Date(r.updated_at).getTime() <= oneWeekAgo).sort((a, b) => new Date(b.updated_at).getTime() - new Date(a.updated_at).getTime())
    };
  });

  async function loadDashboard() {
    try {
      const dataRes = await apiFetch('/api/v2/user/dashboard/data');
      if (dataRes.ok) {
        data = await dataRes.json();
      }
    } catch (e) {
      console.error("Failed to fetch dashboard data", e);
    } finally {
      loading = false;
    }
  }

  function formatAction(act: any) {
    try {
      switch (act.op_type) {
        case 1: 
          return `<span class="action-summary"><span class="material-symbols-rounded mini">add_box</span> <a href="/${act.repo_user_name}/${act.repo_name}">${act.repo_name}</a></span>`;
        case 5: // ActionCommitRepo
          const pushData = JSON.parse(act.content);
          return `
            <div class="action-summary">
              <span class="action-pill"><span class="material-symbols-rounded mini">history</span> ${pushData.Len}</span>
              <span class="action-pill"><span class="material-symbols-rounded mini">fork_right</span> ${act.ref_name}</span>
              <span class="action-at">@</span>
              <a href="/${act.repo_user_name}/${act.repo_name}" class="repo-link">${act.repo_name}</a>
            </div>
          `;
        case 11: // ActionMergePullRequest
          return `<span class="action-summary"><span class="material-symbols-rounded mini">call_merge</span> merged into <a href="/${act.repo_user_name}/${act.repo_name}">${act.repo_name}</a></span>`;
        default:
          return act.content || 'activity';
      }
    } catch (e) {
      return act.content;
    }
  }

  onMount(() => {
    const handleGlobalSearch = (e: any) => { searchQuery = e.detail; };
    window.addEventListener('ggit:search', handleGlobalSearch);
    loadDashboard();
    return () => window.removeEventListener('ggit:search', handleGlobalSearch);
  });
</script>

<div class="dashboard-grid full-width">
    <!-- Repositories List -->
    <div class="dashboard-main-content">
        <header class="content-controls">
            <div class="search-row"> 
                <div class="search-box">
                    <GlobalSearch />
                </div>
                <a class="btn-action-primary" href="/repo/create">
                    <span class="material-symbols-rounded">add</span>
                    <span>New Project</span>
                </a>
            </div>
            <div class="view-mode-selector">
                <button class="tab-pill" class:active={viewMode === 'temporal'} onclick={() => viewMode = 'temporal'}>
                    <span class="material-symbols-rounded mini">schedule</span>
                    Activity
                </button>
                <button class="tab-pill" class:active={viewMode === 'list'} onclick={() => viewMode = 'list'}>
                    <span class="material-symbols-rounded mini">list</span>
                    Alphabetical
                </button>
            </div>
        </header>

        <div class="repository-stream">
            {#if loading}
                <div class="skeleton-list">
                    {#each Array(5) as _}<div class="skeleton-row"></div>{/each}
                </div>
            {:else if allRepos.length === 0}
                <div class="empty-dashboard">
                    <span class="material-symbols-rounded">rocket_launch</span>
                    <h3>Start your first project</h3>
                    <p>No repositories found. Create one to begin your journey.</p>
                </div>
            {:else}
                <div class="flat-repo-list">
                {#if viewMode === 'temporal'}
                    {#each groupedRepos().recent || [] as repo}
                        <RepoRow {repo} />
                    {/each}
                    {#if groupedRepos().older?.length > 0}
                        <div class="divider-text">Older Projects</div>
                        {#each groupedRepos().older as repo}
                            <RepoRow {repo} isOlder={true} />
                        {/each}
                    {/if}
                {:else}
                    {#each groupedRepos().all || [] as repo}
                        <RepoRow {repo} />
                    {/each}
                {/if}
                </div>
            {/if}
        </div>
    </div>
</div>

<style>
  .dashboard-grid { 
    display: flex; 
    flex-direction: column; 
    height: 100%; 
    padding: 1rem;
  }

  .dashboard-main-content { 
    display: flex; 
    flex-direction: column; 
    gap: 2rem;
  }

  .content-controls { 
    display: flex; 
    flex-direction: column; 
    gap: 1.5rem; 
    padding-bottom: 2rem;
    border-bottom: 1px solid var(--border-color);
  }

  .search-row { display: flex; gap: 1rem; align-items: center; }
  .search-box { flex: 1; }

  .error-toast {
    background: var(--error-container);
    color: var(--on-error-container);
    padding: 1rem 1.5rem;
    border-radius: 1rem;
    margin-bottom: 2rem;
    display: flex;
    align-items: center;
    gap: 1rem;
    font-weight: 700;
    box-shadow: var(--shadow-sm);
    animation: shake 0.4s ease-out;
  }

  .view-mode-selector { display: flex; gap: 0.5rem; }
  .tab-pill { 
    background: transparent; 
    border: 1px solid var(--border-color); 
    padding: 0.5rem 1rem; 
    border-radius: 2rem; 
    color: var(--text-muted); 
    font-weight: 700; 
    font-size: 0.85rem; 
    cursor: pointer; 
    display: flex; 
    align-items: center; 
    gap: 0.5rem; 
    transition: all 0.2s;
  }
  .selector-btn.active { background: var(--primary-color); color: var(--on-primary); border-color: var(--primary-color); }
  .selector-btn .mini { font-size: 1.1rem; }

  .flat-repo-list { 
    display: flex; 
    flex-direction: column; 
    background: var(--bg-color); 
    border: 1px solid var(--border-color); 
    border-radius: 1rem; 
    overflow: hidden; 
  }

  .divider-text { 
    padding: 1.5rem 1.25rem 0.75rem; 
    font-size: 0.7rem; 
    font-weight: 900; 
    text-transform: uppercase; 
    color: var(--text-muted); 
    letter-spacing: 0.1em; 
    background: var(--surface-variant);
    border-top: 1px solid var(--border-color);
    border-bottom: 1px solid var(--border-color);
  }

  .empty-dashboard { padding: 6rem 2rem; text-align: center; color: var(--text-muted); }
  .empty-dashboard span { font-size: 4rem; opacity: 0.3; margin-bottom: 1rem; }

  .skeleton-list { display: flex; flex-direction: column; gap: 1px; background: var(--border-color); border-radius: 1rem; overflow: hidden; }
  .skeleton-row { height: 80px; background: var(--bg-color); animation: pulse 1.5s infinite; }
  @keyframes pulse { 0%, 100% { opacity: 1; } 50% { opacity: 0.7; } }
</style>
