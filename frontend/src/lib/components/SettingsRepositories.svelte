<script lang="ts">
  import { onMount } from 'svelte';
  import RepoRow from './RepoRow.svelte';
  import { apiFetch } from '$lib/api';
  
  let repos = $state<any[]>([]);
  let loading = $state(true);

  async function fetchRepos() {
    loading = true;
    try {
      const res = await apiFetch('/api/v2/user/settings/repos');
      if (res.ok) {
        repos = await res.json();
      }
    } catch (e) {
      console.error(e);
    } finally {
      loading = false;
    }
  }

  async function leaveRepo(id: number, name: string) {
    if (!confirm(`Are you sure you want to leave the collaboration in ${name}?`)) return;
    try {
      const res = await apiFetch(`/api/v2/user/settings/repos/${id}`, {
        method: 'DELETE'
      });
      if (res.ok) {
        await fetchRepos();
      }
    } catch (e) {
      console.error(e);
    }
  }

  onMount(fetchRepos);
</script>

<div class="settings-card island">
  <header class="card-header">
    <h3 class="card-title">Repositories</h3>
    <p class="card-desc">List of repositories you own or collaborate in.</p>
  </header>

  <div class="card-body">
    {#if loading}
      <div class="skeletons">
        {#each Array(3) as _}<div class="skeleton-line"></div>{/each}
      </div>
    {:else if repos.length === 0}
      <div class="empty-state">
        <span class="material-symbols-rounded">folder_open</span>
        <p>You don't have any repositories yet.</p>
      </div>
    {:else}
      <div class="repo-list">
        {#each repos as repo}
          <div class="repo-entry">
            <div class="repo-main">
               <RepoRow {repo} />
            </div>
            {#if repo.owner.username !== repo.context_user} <!-- If collaborative -->
              <button class="btn btn-ghost error btn-sm" onclick={() => leaveRepo(repo.id, repo.full_name)}>
                <span class="material-symbols-rounded">logout</span> Leave
              </button>
            {/if}
          </div>
        {/each}
      </div>
    {/if}
  </div>
</div>

<style>
  .settings-card { display: flex; flex-direction: column; }
  .card-header { padding: 1.5rem 2rem; border-bottom: 1px solid var(--md-sys-color-outline-variant); }
  .card-title { margin: 0; font-size: 1.25rem; font-weight: 700; }
  .card-desc { margin: 0.5rem 0 0; font-size: 0.9rem; color: var(--text-muted); }

  .card-body { padding: 1.5rem; }

  .repo-list { display: flex; flex-direction: column; gap: 0.75rem; }
  .repo-entry { display: flex; align-items: center; gap: 1rem; }
  .repo-main { flex: 1; }

  .skeleton-line { height: 80px; background: var(--md3-surface-3); border-radius: 1rem; margin-bottom: 0.5rem; animation: pulse 1.5s infinite; }
  @keyframes pulse { 0%, 100% { opacity: 0.5; } 50% { opacity: 0.8; } }

  .empty-state { text-align: center; padding: 4rem; color: var(--text-muted); }
  .empty-state span { font-size: 3rem; opacity: 0.2; }

  .btn-ghost.error { color: #ef4444; }
  .btn-ghost.error:hover { background: rgba(239, 68, 68, 0.1); }
</style>
