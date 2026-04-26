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
  let feeds = $state<any[]>([]);
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
      const [dataRes, feedRes] = await Promise.all([
        apiFetch('/api/v2/user/dashboard/data'),
        apiFetch('/api/v2/user/dashboard/feeds')
      ]);
      
      if (dataRes.ok) {
        data = await dataRes.json();
      }
      if (feedRes.ok) {
        const feedData = await feedRes.json();
        feeds = feedData.feeds || [];
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
        case 6: // ActionCreateIssue
          return `<span class="action-summary"><span class="material-symbols-rounded mini">error</span> opened issue <a href="/${act.repo_user_name}/${act.repo_name}/issues/${act.content.split('|')[0]}">${act.repo_name}#${act.content.split('|')[0]}</a></span>`;
        case 7: // ActionCreatePullRequest
          return `<span class="action-summary"><span class="material-symbols-rounded mini">merge</span> opened pull request <a href="/${act.repo_user_name}/${act.repo_name}/pulls/${act.content.split('|')[0]}">${act.repo_name}#${act.content.split('|')[0]}</a></span>`;
        case 10: // ActionCommentIssue
          return `<span class="action-summary"><span class="material-symbols-rounded mini">chat_bubble</span> commented on <a href="/${act.repo_user_name}/${act.repo_name}/issues/${act.content.split('|')[0]}">${act.repo_name}#${act.content.split('|')[0]}</a></span>`;
        case 16: // ActionCreateBranch
          return `<span class="action-summary"><span class="material-symbols-rounded mini">add_circle</span> created branch <span class="branch-pill">${act.ref_name}</span> at <a href="/${act.repo_user_name}/${act.repo_name}">${act.repo_name}</a></span>`;
        case 17: // ActionDeleteBranch
          return `<span class="action-summary"><span class="material-symbols-rounded mini">delete</span> deleted branch <span class="branch-pill deleted">${act.ref_name}</span> from <a href="/${act.repo_user_name}/${act.repo_name}">${act.repo_name}</a></span>`;
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

<div class="dashboard-layout">
    <div class="dashboard-main-content">
        <ViewHeader 
            bind:activeTab={viewMode}
            tabs={[
                { id: 'temporal', label: 'Activity', icon: 'schedule' },
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

    <aside class="dashboard-sidebar">
        <div class="sidebar-card news-feed-card tonal">
            <header class="sidebar-header">
                <span class="material-symbols-rounded">rss_feed</span>
                <h4>News Feed</h4>
            </header>
            
            <div class="feed-list">
                {#if loading}
                    {#each Array(3) as _}<div class="skeleton-feed-item"></div>{/each}
                {:else if feeds.length === 0}
                    <div class="empty-feed">
                        <p>No recent activity in your network.</p>
                    </div>
                {:else}
                    {#each feeds as act}
                        <div class="feed-item">
                            <div class="feed-avatar">
                                <img src={act.act_avatar} alt={act.act_user_name} />
                            </div>
                            <div class="feed-content">
                                <div class="feed-user-line">
                                    <a href="/${act.act_user_name}" class="feed-user-link"><strong>${act.act_user_name}</strong></a>
                                    <span class="feed-time">{new Date(act.created_at).toLocaleDateString()}</span>
                                </div>
                                <div class="feed-action-line">
                                    {@html formatAction(act)}
                                </div>
                            </div>
                        </div>
                    {/each}
                {/if}
            </div>
        </div>

        {#if data.organizations.length > 0}
            <div class="sidebar-card tonal">
                <header class="sidebar-header">
                    <span class="material-symbols-rounded">corporate_fare</span>
                    <h4>Organizations</h4>
                </header>
                <div class="org-grid">
                    {#each data.organizations as org}
                        <a href="/org/${org.username}" class="org-chip" title={org.full_name || org.username}>
                            <img src={org.avatar_url} alt={org.username} />
                        </a>
                    {/each}
                </div>
            </div>
        {/if}
    </aside>
</div>

<style>
  .dashboard-layout { 
    display: grid; 
    grid-template-columns: 1fr 320px;
    gap: 2rem;
    height: 100%; 
    padding-right: 1.5rem;
  }

  .dashboard-main-content { 
    display: flex; 
    flex-direction: column; 
    min-width: 0;
  }

  .repository-stream {
    padding: 0 0 4rem;
  }

  .dashboard-sidebar {
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
    .dashboard-layout { grid-template-columns: 1fr; }
    .dashboard-sidebar { display: none; }
  }

  @media (max-width: 600px) {
    .hide-mobile { display: none; }
  }
</style>
