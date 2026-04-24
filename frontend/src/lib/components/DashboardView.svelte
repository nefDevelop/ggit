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
                      <span>Nuevo Repositorio</span>
                  </a>
                  <a class="btn-action-primary" href="/repo/create?mode=migrate" title="Migrate or Clone Repository">
                      <span class="material-symbols-rounded">cloud_download</span>
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
                    {#each groupedRepos().recent || [] as repo}
                        <RepoRow {repo} />
                    {/each}
                    {#if (groupedRepos().older?.length ?? 0) > 0}
                        <div class="divider-text">Older Projects</div>
                        {#each groupedRepos().older || [] as repo}
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
  }

  .dashboard-main-content { 
    display: flex; 
    flex-direction: column; 
  }

  .repository-stream {
    padding: 0 1.5rem 4rem;
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

  .empty-dashboard { padding: 6rem 2rem; text-align: center; color: var(--text-muted); }
  .empty-dashboard span { font-size: 4rem; opacity: 0.3; margin-bottom: 1rem; }

  .skeleton-list { display: flex; flex-direction: column; gap: 1px; background: var(--md-sys-color-outline-variant); border-radius: 1rem; overflow: hidden; }
  .skeleton-row { height: 80px; background: var(--bg-color); animation: pulse 1.5s infinite; }
  @keyframes pulse { 0%, 100% { opacity: 1; } 50% { opacity: 0.7; } }
</style>
