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
    <section class="repos-island island">
        <header class="island-header-search">
            <div class="search-and-new-button-group"> 
                <div class="search-wrapper">
                    <GlobalSearch />
                </div>
                <a class="btn btn-sm" href="/repo/create">
                    <span class="material-symbols-rounded">add</span> New
                </a>
            </div>
            <div class="title-and-view">
                <h2 class="section-label">Repositories</h2>
                <div class="view-controls">
                    <button class="view-btn {viewMode === 'temporal' ? 'active' : ''}" onclick={() => viewMode = 'temporal'}>
                        <span class="material-symbols-rounded">schedule</span>
                    </button>
                    <button class="view-btn {viewMode === 'list' ? 'active' : ''}" onclick={() => viewMode = 'list'}>
                        <span class="material-symbols-rounded">list</span>
                    </button>
                </div>
            </div>
        </header>

        <div class="island-body">
            {#if loading}
                <div class="skeletons">
                    {#each Array(3) as _}<div class="skeleton-line"></div>{/each}
                </div>
            {:else if allRepos.length === 0}
                <div class="empty-state">
                    <span class="material-symbols-rounded empty-icon">folder_open</span>
                    <h3>Ready for your first project</h3>
                    <p>Create a new repository to get started.</p>
                </div>
            {:else}
                <div class="repo-list-container">
                {#if viewMode === 'temporal'}
                    {@const groups = groupedRepos()}
                    {#if groups.recent && groups.recent.length > 0}
                        <div class="group-header"><span class="group-label">Active this week</span></div>
                        <div class="repo-group-flat">
                            {#each groups.recent as repo}<RepoRow {repo} />{/each}
                        </div>
                    {/if}
                    {#if groups.older && groups.older.length > 0}
                        <div class="group-header" style="margin-top: 1rem;"><span class="group-label">Older</span></div>
                        <div class="repo-group-flat">
                            {#each groups.older as repo}<RepoRow {repo} isOlder={true} />{/each}
                        </div>
                    {/if}
                {:else}
                    <div class="repo-group-flat">
                        {#each groupedRepos().all || [] as repo}<RepoRow {repo} />{/each}
                    </div>
                {/if}
                </div>
            {/if}
        </div>
    </section>
</div>

<style>
  .dashboard-grid {
    display: flex;
    flex-direction: column;
    height: 100%;
  }

  .repos-island {
    height: auto;
    min-height: min-content;
  }

  .island-header-search {
    padding: 2rem 1.5rem 1.5rem;
    display: flex;
    flex-direction: column;
    gap: 1.25rem;
    background: var(--surface-variant);
    border-bottom: 1px solid var(--border-color);
    box-sizing: border-box;
  }

  .search-and-new-button-group {
    display: flex;
    align-items: center;
    gap: 1rem;
  }

  .search-wrapper {
    flex: 1;
    min-width: 0;
  }

  .search-and-new-button-group .btn {
    flex-shrink: 0;
    white-space: nowrap;
  }

  @media (max-width: 1024px) {
    .island-header-search {
      padding: 1.5rem 1rem;
      gap: 1rem;
    }
    .island-body {
      padding: 1rem;
    }
  }

  .title-and-view {
    display: flex;
    align-items: center;
    justify-content: space-between;
  }

  .section-label {
    font-size: 0.75rem;
    font-weight: 900;
    text-transform: uppercase;
    letter-spacing: 0.1em;
    color: var(--text-muted);
    margin: 0;
  }

  .island-body { 
    padding: 1.5rem; 
    flex: 1;
    overflow-y: auto;
  }

  .view-controls {
    display: flex;
    background: var(--md3-surface-3);
    padding: 2px;
    border-radius: 0.5rem;
  }

  .view-btn {
    background: transparent;
    border: none;
    color: var(--text-muted);
    padding: 0.3rem 0.6rem;
    border-radius: 0.4rem;
    cursor: pointer;
    display: flex;
  }
  .view-btn span { font-size: 1.1rem; }
  .view-btn.active { background: var(--md3-surface-1); color: var(--primary-color); }

  .repo-group-flat {
    display: flex;
    flex-direction: column;
    background: var(--surface-variant);
    border-radius: 1.25rem;
    overflow: hidden;
    border: 1px solid var(--border-color);
  }

  .group-label { font-size: 0.7rem; font-weight: 800; text-transform: uppercase; color: var(--primary-color); padding: 0 0.5rem 0.5rem; display: block; }

  .empty-state {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    min-height: 250px;
    text-align: center;
  }

  .empty-state .empty-icon {
    font-size: 4rem;
    opacity: 0.15;
    margin-bottom: 1rem;
  }

  .empty-state h3 {
    margin: 0 0 0.5rem 0;
    color: var(--text-color);
  }

  .empty-state p {
    margin: 0;
    color: var(--text-muted);
    font-size: 0.95rem;
  }

  .skeleton-line { height: 64px; background: var(--md3-surface-3); border-radius: 0.75rem; margin-bottom: 0.5rem; animation: pulse 1.5s infinite; }

  @keyframes pulse { 0%, 100% { opacity: 0.5; } 50% { opacity: 0.8; } }
</style>
