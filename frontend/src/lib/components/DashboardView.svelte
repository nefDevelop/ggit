<script lang="ts">
  import { apiFetch } from "$lib/api";
  import { authState } from '$lib/auth.svelte';
  import { onMount } from 'svelte';
  import RepoRow from './RepoRow.svelte';
  import GlobalSearch from './GlobalSearch.svelte';
  import ViewHeader from './ui/ViewHeader.svelte';

  let data = $state<{
    repositories: any[],
    collaborative: any[],
    organizations: any[]
  }>({
    repositories: [],
    collaborative: [],
    organizations: []
  });
  let loading = $state(true);
  let viewMode = $state<'list' | 'temporal'>('temporal');
  let searchQuery = $state('');

  const oneWeekAgo = $derived(Date.now() - (7 * 24 * 60 * 60 * 1000));
  const allRepos = $derived([...data.repositories, ...data.collaborative]);
  
  const filteredRepos = $derived.by(() => {
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

  const groupedRepos = $derived.by((): GroupedRepos => {
    const repos = filteredRepos;
    if (viewMode === 'list') {
      return { all: [...repos].sort((a, b) => a.name.localeCompare(b.name)) };
    }
    return {
      recent: repos.filter(r => new Date(r.updated_at).getTime() > oneWeekAgo).sort((a, b) => new Date(b.updated_at).getTime() - new Date(a.updated_at).getTime()),
      older: repos.filter(r => new Date(r.updated_at).getTime() <= oneWeekAgo).sort((a, b) => new Date(b.updated_at).getTime() - new Date(a.updated_at).getTime())
    };
  });

  async function loadDashboard() {
    loading = true;
    try {
      const res = await apiFetch('/api/v2/user/dashboard/data');
      if (res.ok) {
        data = await res.json();
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

    // If activity is disabled, default to list view
    if (authState.user?.disable_news_feed) {
       viewMode = 'list';
    }

    return () => window.removeEventListener('ggit:search', handleGlobalSearch);
  });
</script>

<div class="dashboard-layout">
    <div class="dashboard-main-content">
        <ViewHeader 
            bind:activeTab={viewMode}
            tabs={[
                ...(authState.user?.disable_news_feed ? [] : [{ id: 'temporal', label: 'Activity', icon: 'schedule' }]),
                { id: 'list', label: 'Alphabetical', icon: 'list' }
            ]}
        >
            {#snippet actions()}
                <div class="btn-group">
                  <a class="btn-action-primary" href="/repo/create">
                      <span class="material-symbols-rounded">add</span>
                      <span class="hide-mobile">Nuevo Repositorio</span>
                  </a>
                </div>
            {/snippet}
        </ViewHeader>

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
                    {#each groupedRepos.recent || [] as repo}
                        <RepoRow {repo} />
                    {/each}
                    {#if (groupedRepos.older?.length ?? 0) > 0}
                        <div class="divider-text">Older Projects</div>
                        {#each groupedRepos.older || [] as repo}
                            <RepoRow {repo} isOlder={true} />
                        {/each}
                    {/if}
                {:else}
                    {#each groupedRepos.all || [] as repo}
                        <RepoRow {repo} />
                    {/each}
                {/if}
                </div>
            {/if}
        </div>
    </div>

    {#if data.organizations.length > 0 && !authState.user?.disable_news_feed}
        <aside class="dashboard-sidebar">
            <div class="sidebar-card tonal">
                <header class="sidebar-header">
                    <span class="material-symbols-rounded">corporate_fare</span>
                    <h4>Organizations</h4>
                </header>
                <div class="org-grid">
                    {#each data.organizations as org}
                        <a href="/org/{org.username}" class="org-chip" title={org.full_name || org.username}>
                            <img src={org.avatar_url} alt={org.username} />
                        </a>
                    {/each}
                </div>
            </div>
        </aside>
    {/if}
</div>

<style>
  .dashboard-layout {
    display: flex;
    gap: 2rem;
    height: 100%;
    padding-left: 1.5rem;
    padding-right: 1.5rem;
  }

  .dashboard-main-content { 
    display: flex; 
    flex-direction: column; 
    flex: 1;
    min-width: 0;
  }

  .repository-stream {
    padding: 0 0 4rem;
  }

  .dashboard-sidebar {
    width: 320px;
    flex-shrink: 0;
    display: flex;
    flex-direction: column;
    gap: 1.5rem;
    padding-top: 1.5rem;
  }

  .sidebar-card {
    border-radius: 1.25rem;
    padding: 1.5rem;
    border: 1px solid var(--md-sys-color-outline-variant);
  }

  .sidebar-header {
    display: flex;
    align-items: center;
    gap: 0.75rem;
    margin-bottom: 1.25rem;
    color: var(--text-muted);
  }

  .sidebar-header h4 {
    margin: 0;
    font-size: 0.9rem;
    font-weight: 800;
    text-transform: uppercase;
    letter-spacing: 0.05em;
  }

  .feed-list {
    display: flex;
    flex-direction: column;
    gap: 1.25rem;
  }

  .feed-item {
    display: flex;
    gap: 1rem;
    font-size: 0.85rem;
  }

  .feed-avatar img {
    width: 32px;
    height: 32px;
    border-radius: 8px;
    background: var(--surface-variant);
  }

  .feed-content {
    flex: 1;
    display: flex;
    flex-direction: column;
    gap: 0.25rem;
  }

  .feed-user-line {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }

  .feed-user-link {
    color: var(--text-color);
    text-decoration: none;
  }

  .feed-time {
    font-size: 0.7rem;
    color: var(--text-muted);
    opacity: 0.7;
  }

  .feed-action-line {
    color: var(--text-muted);
    line-height: 1.4;
  }

  :global(.action-summary a) {
    color: var(--primary-color);
    text-decoration: none;
    font-weight: 700;
  }

  :global(.branch-pill) {
    font-family: monospace;
    background: rgba(var(--primary-rgb), 0.1);
    color: var(--primary-color);
    padding: 1px 4px;
    border-radius: 4px;
    font-size: 0.75rem;
  }

  :global(.branch-pill.deleted) {
    background: rgba(239, 68, 68, 0.1);
    color: #ef4444;
    text-decoration: line-through;
  }

  .org-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(40px, 1fr));
    gap: 0.75rem;
  }

  .org-chip img {
    width: 40px;
    height: 40px;
    border-radius: 10px;
    border: 1px solid var(--md-sys-color-outline-variant);
    transition: transform 0.2s;
  }

  .org-chip:hover img {
    transform: scale(1.1) translateY(-2px);
    border-color: var(--primary-color);
  }

  .flat-repo-list { 
    display: flex; 
    flex-direction: column; 
    background: var(--bg-color); 
    border: 1px solid var(--md-sys-color-outline-variant); 
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
    border-top: 1px solid var(--md-sys-color-outline-variant);
    border-bottom: 1px solid var(--md-sys-color-outline-variant);
  }

  .empty-dashboard, .empty-feed { padding: 2rem; text-align: center; color: var(--text-muted); }
  .empty-dashboard span { font-size: 4rem; opacity: 0.3; margin-bottom: 1rem; }

  .skeleton-list { display: flex; flex-direction: column; gap: 1px; background: var(--md-sys-color-outline-variant); border-radius: 1rem; overflow: hidden; }
  .skeleton-row { height: 80px; background: var(--bg-color); animation: pulse 1.5s infinite; }
  .skeleton-feed-item { height: 60px; background: rgba(255,255,255,0.03); border-radius: 10px; animation: pulse 1.5s infinite; }
  
  @keyframes pulse { 0%, 100% { opacity: 1; } 50% { opacity: 0.7; } }

  @media (max-width: 1100px) {
    .dashboard-sidebar { display: none; }
  }

  @media (max-width: 600px) {
    .hide-mobile { display: none; }
  }
</style>
