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
      const [dataRes, feedsRes] = await Promise.all([
        apiFetch('/api/v2/user/dashboard/data'),
        apiFetch('/api/v2/user/dashboard/feeds')
      ]);
      if (dataRes.ok) {
        data = await dataRes.json();
        console.log("Dashboard Data Loaded:", data);
      }
      if (feedsRes.ok) {
        const feedsData = await feedsRes.json();
        feeds = feedsData.feeds;
        console.log("Dashboard Feeds Loaded:", feeds.length);
      }
    } catch (e) {
      console.error("Failed to fetch dashboard data", e);
    } finally {
      loading = false;
    }
  }

  onMount(() => {
    const handleGlobalSearch = (e: any) => { searchQuery = e.detail; };
    window.addEventListener('ggit:search', handleGlobalSearch);
    loadDashboard();
    return () => window.removeEventListener('ggit:search', handleGlobalSearch);
  });
</script>

<div class="dashboard-grid" style:grid-template-columns={authState.user?.disable_news_feed ? '1fr' : '1.5fr 1fr'}>
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

    {#if !authState.user?.disable_news_feed}
    <div class="sidebar-column">
        <!-- Organizations Section -->
        {#if (data.organizations || []).length > 0}
            <section class="orgs-island island tonal">
                <header class="simple-header">
                    <h2 class="section-label">Organizations</h2>
                </header>
                <div class="island-body">
                    <div class="org-list">
                        {#each data.organizations as org}
                            <a href="/org/{org.username}" class="org-item hover-lift">
                                <img src={org.avatar_url} alt="" class="org-avatar" />
                                <span class="org-name">{org.username}</span>
                            </a>
                        {/each}
                    </div>
                </div>
            </section>
        {/if}

        <section class="activity-island island">
            <header class="simple-header">
                <h2 class="section-label">Recent Activity</h2>
            </header>
            <div class="island-body">
                {#if loading}
                    <div class="skeletons">
                        {#each Array(5) as _}<div class="skeleton-feed"></div>{/each}
                    </div>
                {:else if feeds.length === 0}
                    <div class="empty-state mini"><p class="text-muted">No activity.</p></div>
                {:else}
                    <div class="feed-list">
                    {#each feeds as act}
                        <div class="feed-item">
                        <img src={act.act_avatar} alt="" class="feed-avatar" />
                        <div class="feed-content">
                            <span class="feed-user">{act.act_user_name}</span>
                            <span class="feed-action">{@html act.content}</span>
                            <span class="feed-time">{new Date(act.created_unix * 1000).toLocaleDateString()}</span>
                        </div>
                        </div>
                    {/each}
                    </div>
                {/if}
            </div>
        </section>
    </div>
    {/if}
</div>

<style>
  .dashboard-grid {
    display: grid;
    gap: 1.5rem;
    align-items: stretch;
    height: 100%;
  }

  .repos-island, .activity-island {
    display: flex;
    flex-direction: column;
    padding: 0; /* Removed !important, adjust specificity if needed */
    overflow: hidden;
  }

  .repos-island {
    height: 100%;
  }

  .activity-island {
    max-height: 40vh;
    overflow-y: auto;
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
    min-width: 0; /* Previene que el buscador rompa el layout de flexbox */
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

  .header-tools {
    display: flex;
    justify-content: space-between;
    align-items: center;
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
    overflow-y: auto; /* Evita que el contenido se salga por debajo de la isla */
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

  .simple-header {
    padding: 1.5rem;
  }

  .repos-island, .activity-island {
    box-shadow: none; /* Removed !important, adjust specificity if needed */
    border: 1px solid var(--border-color);
    background: transparent; /* Removed !important, adjust specificity if needed */
  }

  .repo-group-flat {
    display: flex;
    flex-direction: column;
    background: var(--surface-variant);
    border-radius: 1.25rem;
    overflow: hidden;
    border: 1px solid var(--border-color);
  }

  .group-label { font-size: 0.7rem; font-weight: 800; text-transform: uppercase; color: var(--primary-color); padding: 0 0.5rem 0.5rem; display: block; }
  .feed-list { display: flex; flex-direction: column; gap: 1.25rem; }
  .feed-item { display: flex; gap: 1rem; font-size: 0.85rem; }
  .feed-avatar { width: 32px; height: 32px; border-radius: 50%; opacity: 0.8; }
  .feed-user { font-weight: 700; color: var(--text-color); }
  .feed-time { color: var(--text-muted); font-size: 0.75rem; }

  .empty-state {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    min-height: 250px;
    text-align: center;
  }

  .empty-state.mini {
    min-height: 100px;
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
  .skeleton-feed { height: 40px; background: var(--md3-surface-3); border-radius: 0.5rem; margin-bottom: 1rem; animation: pulse 1.5s infinite; }

  .sidebar-column { display: flex; flex-direction: column; gap: 1.5rem; }
  .orgs-island { margin-bottom: 0; }
  .org-list { display: flex; flex-direction: column; gap: 0.5rem; }
  .org-item { display: flex; align-items: center; gap: 0.75rem; text-decoration: none; color: var(--text-color); padding: 0.5rem; border-radius: 0.75rem; transition: background 0.2s; }
  .org-item:hover { background: rgba(255,255,255,0.05); }
  .org-avatar { width: 24px; height: 24px; border-radius: 4px; }
  .org-name { font-weight: 600; font-size: 0.9rem; }

  @keyframes pulse { 0%, 100% { opacity: 0.5; } 50% { opacity: 0.8; } }
  @media (max-width: 1024px) { 
    .dashboard-grid { grid-template-columns: 1fr; } /* Removed !important, adjust specificity if needed */
    /* .hero-header is not used in the template, removed */
  }
</style>
