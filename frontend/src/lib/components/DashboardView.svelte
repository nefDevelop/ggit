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

  const oneWeekAgo = $derived(Date.now() - (7 * 24 * 60 * 60 * 1000));

  const groupedRepos = $derived(() => {
    const all = [...data.repositories, ...data.collaborative];
    if (viewMode === 'list') return { all };

    return {
      recent: all.filter(r => new Date(r.updated_at).getTime() > oneWeekAgo).sort((a, b) => new Date(b.updated_at).getTime() - new Date(a.updated_at).getTime()),
      older: all.filter(r => new Date(r.updated_at).getTime() <= oneWeekAgo).sort((a, b) => new Date(b.updated_at).getTime() - new Date(a.updated_at).getTime())
    };
  });

  onMount(async () => {
    try {
      const [dataRes, feedsRes] = await Promise.all([
        fetch('/api/v2/user/dashboard/data'),
        fetch('/api/v2/user/dashboard/feeds')
      ]);

      if (dataRes.ok) {
        data = await dataRes.json();
      }
      if (feedsRes.ok) {
        const feedsData = await feedsRes.json();
        feeds = feedsData.feeds;
      }
    } catch (e) {
      console.error("Failed to fetch dashboard data", e);
    } finally {
      loading = false;
    }
  });

  function getCsrfToken() {
    if (typeof document === 'undefined') return '';
    const match = document.cookie.match(new RegExp('(^| )XSRF-TOKEN=([^;]+)'));
    return match ? match[2] : '';
  }
</script>

<div class="dashboard-content">
  <main class="island content">
    <div class="dashboard-grid">
      <section class="repos-island">
        <header class="content-header">
          <div class="header-main">
            <h2 class="expressive-heading">Repositories</h2>
            <div class="view-controls">
              <button 
                class="view-btn {viewMode === 'temporal' ? 'active' : ''}" 
                onclick={() => viewMode = 'temporal'}
                title="Temporal View"
              >
                <span class="material-symbols-rounded">schedule</span>
                <span class="btn-text">Temporal</span>
              </button>
              <button 
                class="view-btn {viewMode === 'list' ? 'active' : ''}" 
                onclick={() => viewMode = 'list'}
                title="Full List"
              >
                <span class="material-symbols-rounded">list</span>
                <span class="btn-text">List</span>
              </button>
            </div>
          </div>
          <a class="btn btn-sm btn-primary" href="/repo/create">
             <span class="material-symbols-rounded">add</span> New
          </a>
        </header>

        {#if loading}
          <div class="skeletons">
            {#each Array(3) as _}
              <div class="skeleton-line"></div>
            {/each}
          </div>
        {:else if data.repositories.length === 0 && data.collaborative.length === 0}
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
          <div class="repo-list">
            {#if viewMode === 'temporal'}
              {@const groups = groupedRepos()}
              {#if groups.recent.length > 0}
                <div class="group-label">Active this week</div>
                {#each groups.recent as repo}
                  <RepoRow {repo} />
                {/each}
              {/if}

              {#if groups.older.length > 0}
                <div class="group-label">Older</div>
                {#each groups.older as repo}
                  <RepoRow {repo} />
                {/each}
              {/if}
            {:else}
              {#each data.repositories as repo}
                <RepoRow {repo} />
              {/each}
              {#if data.collaborative.length > 0}
                <div class="collaboration-divider">Collaborative</div>
                {#each data.collaborative as repo}
                  <RepoRow {repo} />
                {/each}
              {/if}
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
    max-width: 1400px;
    margin: 0 auto;
    width: 100%;
    padding: 1.5rem;
  }

  .expressive-heading {
    font-size: 1.5rem;
    font-weight: 700;
    letter-spacing: -0.02em;
    margin: 0;
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
    gap: 0.75rem;
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
    transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1);
  }

  .view-btn span { font-size: 1.2rem; }
  .btn-text { font-size: 0.85rem; font-weight: 600; }

  .view-btn.active {
    background: var(--md3-surface-1);
    color: var(--primary-color);
    box-shadow: 0 2px 8px rgba(0,0,0,0.12);
  }

  .dashboard-grid {
    display: grid;
    grid-template-columns: 1fr 340px;
    gap: 2rem;
  }

  .repos-island, .activity-island {
    background: var(--md3-surface-2);
    padding: 1.5rem;
    border-radius: 1.5rem;
    border: 1px solid var(--border-color);
  }

  .repo-list {
    display: flex;
    flex-direction: column;
    background: var(--md3-surface-1);
    border-radius: 1rem;
    overflow: hidden;
    border: 1px solid var(--border-color);
  }

  .group-label {
    padding: 1.25rem 1rem 0.5rem 1rem;
    font-size: 0.75rem;
    font-weight: 800;
    text-transform: uppercase;
    letter-spacing: 0.1em;
    color: var(--primary-color);
    background: var(--md3-surface-1);
  }

  .collaboration-divider {
    padding: 0.75rem 1rem;
    background: var(--md3-surface-3);
    font-size: 0.75rem;
    font-weight: 700;
    text-transform: uppercase;
    color: var(--text-muted);
    border-top: 1px solid var(--border-color);
  }

  .feed-list {
    display: flex;
    flex-direction: column;
    gap: 1.5rem;
  }

  .feed-item {
    display: flex;
    gap: 1rem;
    font-size: 0.9rem;
  }

  .feed-avatar { width: 36px; height: 36px; border-radius: 50%; }
  .feed-user { font-weight: 600; color: var(--primary-color); }
  .feed-time { display: block; font-size: 0.75rem; color: var(--text-muted); margin-top: 0.25rem; }

  .skeleton-line { height: 70px; background: var(--md3-surface-3); border-radius: 0.75rem; margin-bottom: 0.75rem; animation: pulse 1.5s infinite; }
  .skeleton-feed { height: 50px; background: var(--md3-surface-3); border-radius: 0.75rem; margin-bottom: 1rem; animation: pulse 1.5s infinite; }

  @keyframes pulse { 0%, 100% { opacity: 0.5; } 50% { opacity: 0.8; } }

  @media (max-width: 1100px) { .dashboard-grid { grid-template-columns: 1fr; } }
</style>
