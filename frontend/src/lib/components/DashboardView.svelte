<script lang="ts">
  import { authState } from '$lib/auth.svelte';
  import { onMount } from 'svelte';
  import RepoRow from './RepoRow.svelte';

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
  let searchInput: HTMLInputElement;

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

  const groupedRepos = $derived(() => {
    const repos = searchQuery ? filteredRepos() : allRepos;
    if (viewMode === 'list') return { all: repos };

    return {
      recent: repos.filter(r => new Date(r.updated_at).getTime() > oneWeekAgo).sort((a, b) => new Date(b.updated_at).getTime() - new Date(a.updated_at).getTime()),
      older: repos.filter(r => new Date(r.updated_at).getTime() <= oneWeekAgo).sort((a, b) => new Date(b.updated_at).getTime() - new Date(a.updated_at).getTime())
    };
  });

  // Shortcut for search
  function handleKeydown(e: KeyboardEvent) {
    if (e.key === '/' && document.activeElement?.tagName !== 'INPUT') {
      e.preventDefault();
      searchInput?.focus();
    }
  }

  onMount(async () => {
    window.addEventListener('keydown', handleKeydown);
    try {
      const [dataRes, feedsRes] = await Promise.all([
        fetch('/api/v2/user/dashboard/data'),
        fetch('/api/v2/user/dashboard/feeds')
      ]);

      if (dataRes.ok) data = await dataRes.json();
      if (feedsRes.ok) {
        const feedsData = await feedsRes.json();
        feeds = feedsData.feeds;
      }
    } catch (e) {
      console.error("Failed to fetch dashboard data", e);
    } finally {
      loading = false;
    }
    return () => window.removeEventListener('keydown', handleKeydown);
  });
</script>

<div class="dashboard-content">
  <main class="island content">
    <div class="dashboard-grid">
      <section class="repos-island">
        <header class="content-header">
          <div class="header-main">
            <div class="search-row">
              <div class="search-container pill">
                <span class="material-symbols-rounded search-icon">search</span>
                <input 
                  bind:this={searchInput}
                  bind:value={searchQuery}
                  placeholder="Search projects... (Press /)" 
                  class="search-input"
                />
                {#if searchQuery}
                  <button class="clear-search" onclick={() => searchQuery = ''}>
                    <span class="material-symbols-rounded">close</span>
                  </button>
                {/if}
              </div>
              <a class="btn btn-sm btn-primary" href="/repo/create">
                 <span class="material-symbols-rounded">add</span> New
              </a>
            </div>

            <div class="title-row">
              <h2 class="expressive-heading">Repositories</h2>
              <div class="view-controls">
                <button 
                  class="view-btn {viewMode === 'temporal' ? 'active' : ''}" 
                  onclick={() => viewMode = 'temporal'}
                >
                  <span class="material-symbols-rounded">schedule</span>
                  <span class="btn-text">Temporal</span>
                </button>
                <button 
                  class="view-btn {viewMode === 'list' ? 'active' : ''}" 
                  onclick={() => viewMode = 'list'}
                >
                  <span class="material-symbols-rounded">list</span>
                  <span class="btn-text">List</span>
                </button>
              </div>
            </div>
          </div>
        </header>

        {#if loading}
          <div class="skeletons">
            {#each Array(3) as _}
              <div class="skeleton-line"></div>
            {/each}
          </div>
        {:else if allRepos.length === 0}
          <div class="empty-state">
            <div class="empty-icon">
              <span class="material-symbols-rounded" style="font-size: 3rem;">folder_open</span>
            </div>
            <h3>No repositories yet</h3>
            <p class="text-muted">Start by creating your first project.</p>
            <br/>
            <a class="btn" href="/repo/create">Create Repository</a>
          </div>
        {:else}
          <div class="repo-list-container">
            {#if viewMode === 'temporal'}
              {@const groups = groupedRepos()}
              
              {#if groups.recent.length > 0}
                <div class="group-header">
                  <span class="group-label">Active this week</span>
                  <span class="group-count">{groups.recent.length}</span>
                </div>
                <div class="repo-group recent">
                  {#each groups.recent as repo}
                    <RepoRow {repo} />
                  {/each}
                </div>
              {/if}

              {#if groups.older.length > 0}
                <div class="group-header">
                  <span class="group-label">Older</span>
                  <span class="group-count">{groups.older.length}</span>
                </div>
                <div class="repo-group older">
                  {#each groups.older as repo}
                    <RepoRow {repo} isOlder={true} />
                  {/each}
                </div>
              {/if}
            {:else}
              <div class="repo-group">
                {#each groupedRepos().all as repo}
                  <RepoRow {repo} />
                {/each}
              </div>
            {/if}
          </div>
        {/if}
      </section>

      <section class="activity-island">
        <header class="content-header">
          <h2 class="expressive-heading">Activity</h2>
        </header>

        {#if loading}
          <div class="skeletons">
            {#each Array(5) as _}
              <div class="skeleton-feed"></div>
            {/each}
          </div>
        {:else if feeds.length === 0}
          <div class="empty-state mini">
            <p class="text-muted">No recent activity.</p>
          </div>
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
      </section>
    </div>
  </main>
</div>

<style>
  .dashboard-content {
    width: 100%;
    box-sizing: border-box;
  }

  .expressive-heading {
    font-size: 1.25rem;
    font-weight: 800;
    letter-spacing: -0.02em;
    color: var(--text-color);
  }

  .content {
    background: transparent;
    border: none;
    box-shadow: none;
    padding: 0;
  }

  .content-header {
    display: flex;
    justify-content: space-between;
    align-items: flex-start;
    margin-bottom: 2rem;
  }

  .header-main {
    display: flex;
    flex-direction: column;
    gap: 1.5rem;
    flex: 1;
    width: 100%;
  }

  .search-row {
    display: flex;
    align-items: center;
    justify-content: space-between;
    width: 100%;
    margin-bottom: 0.5rem;
    gap: 1.5rem;
  }

  .title-row {
    display: flex;
    align-items: center;
    justify-content: space-between;
    width: 100%;
  }

  /* Search Terminal Style */
  .search-container {
    position: relative;
    display: flex;
    align-items: center;
    background: var(--md3-surface-3);
    border-radius: 0.75rem;
    padding: 0 0.75rem;
    border: 1px solid transparent;
    transition: all 0.2s ease;
    flex: 1;
    max-width: 500px;
  }

  .search-container.pill {
    border-radius: 2rem;
    background: var(--md3-surface-container);
    border: 1px solid var(--border-color);
  }

  .search-container:focus-within {
    background: var(--md3-surface-1);
    border-color: var(--primary-color);
    box-shadow: 0 0 0 3px rgba(var(--primary-rgb), 0.1);
  }

  .search-icon {
    font-size: 1.2rem;
    color: var(--text-muted);
  }

  .search-input {
    background: transparent;
    border: none;
    padding: 0.75rem;
    color: var(--text-color);
    font-size: 0.9rem;
    width: 100%;
    font-family: inherit;
  }

  .search-input:focus { outline: none; }

  .clear-search {
    background: transparent;
    border: none;
    color: var(--text-muted);
    cursor: pointer;
    display: flex;
    padding: 0.25rem;
  }

  .view-controls {
    display: flex;
    background: var(--md3-surface-3);
    padding: 3px;
    border-radius: 0.75rem;
    width: fit-content;
  }

  .view-btn {
    background: transparent;
    border: none;
    color: var(--text-muted);
    padding: 0.4rem 0.8rem;
    border-radius: 0.6rem;
    cursor: pointer;
    display: flex;
    align-items: center;
    gap: 0.5rem;
    transition: all 0.2s ease;
  }

  .view-btn.active {
    background: var(--md3-surface-1);
    color: var(--primary-color);
    box-shadow: 0 2px 8px rgba(0,0,0,0.1);
  }

  .btn-text { font-size: 0.8rem; font-weight: 700; text-transform: uppercase; }

  /* Grid & Layout */
  .dashboard-grid {
    display: grid;
    grid-template-columns: 1fr 340px;
    gap: 1.5rem; /* Consistent gap between islands */
  }

  .repos-island, .activity-island {
    background: var(--md3-surface-2);
    padding: 1.5rem;
    border-radius: 1.5rem;
    border: 1px solid var(--border-color);
    height: fit-content;
  }

  /* Rhythm & Partitioning */
  .repo-list-container {
    display: flex;
    flex-direction: column;
    gap: 1.5rem;
  }

  .repo-group {
    display: flex;
    flex-direction: column;
    background: var(--md3-surface-1);
    border-radius: 1rem;
    overflow: hidden;
    border: 1px solid var(--border-color);
  }

  .repo-group.recent {
    gap: 2px;
    background: transparent;
    border: none;
  }

  .group-header {
    display: flex;
    align-items: center;
    gap: 0.75rem;
    padding: 0 0.5rem;
  }

  .group-label {
    font-size: 0.75rem;
    font-weight: 800;
    text-transform: uppercase;
    letter-spacing: 0.1em;
    color: var(--primary-color);
  }

  .group-count {
    font-size: 0.7rem;
    background: var(--md3-surface-3);
    padding: 0.1rem 0.5rem;
    border-radius: 1rem;
    color: var(--text-muted);
  }

  .feed-list { display: flex; flex-direction: column; gap: 1.5rem; }
  .feed-item { display: flex; gap: 1rem; font-size: 0.9rem; }
  .feed-avatar { width: 36px; height: 36px; border-radius: 50%; }
  .feed-user { font-weight: 600; color: var(--primary-color); }
  .feed-time { display: block; font-size: 0.75rem; color: var(--text-muted); margin-top: 0.25rem; }

  .skeleton-line { height: 70px; background: var(--md3-surface-3); border-radius: 0.75rem; margin-bottom: 0.75rem; animation: pulse 1.5s infinite; }
  .skeleton-feed { height: 50px; background: var(--md3-surface-3); border-radius: 0.75rem; margin-bottom: 1rem; animation: pulse 1.5s infinite; }

  @keyframes pulse { 0%, 100% { opacity: 0.5; } 50% { opacity: 0.8; } }

  @media (max-width: 1100px) { .dashboard-grid { grid-template-columns: 1fr; } }
</style>
