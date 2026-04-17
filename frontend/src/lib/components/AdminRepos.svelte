<script lang="ts">
  import { onMount } from 'svelte';
  
  let repos = $state<any[]>([]);
  let loading = $state(true);
  let page = $state(1);

  async function fetchRepos() {
    loading = true;
    try {
      const res = await fetch(`/api/v2/admin/repos?page=${page}`);
      if (res.ok) {
        repos = await res.json();
      }
    } catch (e) {
      console.error(e);
    } finally {
      loading = false;
    }
  }

  async function deleteRepo(id: number, name: string) {
    if (!confirm(`Are you sure you want to delete repository "${name}"? This action cannot be undone.`)) {
      return;
    }

    try {
      const res = await fetch(`/api/v2/admin/repos/${id}`, {
        method: 'DELETE'
      });
      if (res.ok) {
        repos = repos.filter(r => r.id !== id);
      } else {
        alert('Failed to delete repository.');
      }
    } catch (e) {
      console.error(e);
      alert('Error connecting to server.');
    }
  }

  onMount(fetchRepos);
</script>

<div class="admin-repos island tonal">
  <header class="section-header">
    <div class="title">
      <span class="material-symbols-rounded">database</span>
      <h2>Repository Oversight</h2>
    </div>
  </header>

  <div class="table-container">
    <table class="modern-table">
      <thead>
        <tr>
          <th>ID</th>
          <th>Repository</th>
          <th>Owner</th>
          <th>Stars</th>
          <th>Forks</th>
          <th>Private</th>
          <th>Actions</th>
        </tr>
      </thead>
      <tbody>
        {#if loading}
          {#each Array(5) as _}
            <tr class="skeleton-row">
              <td colspan="7"><div class="skeleton-bar"></div></td>
            </tr>
          {/each}
        {:else}
          {#each repos as repo}
            <tr>
              <td>{repo.id}</td>
              <td>
                <div class="repo-cell">
                   <strong class="repo-name">{repo.name}</strong>
                   <span class="repo-desc text-muted">{repo.description || 'No description'}</span>
                </div>
              </td>
              <td>{repo.owner?.username}</td>
              <td>{repo.num_stars}</td>
              <td>{repo.num_forks}</td>
              <td>
                <span class="status-badge" class:private={repo.is_private}>
                  {repo.is_private ? 'Private' : 'Public'}
                </span>
              </td>
              <td class="actions">
                <a href="/{repo.owner?.username}/{repo.name}" class="icon-btn" title="View"><span class="material-symbols-rounded">visibility</span></a>
                <button class="icon-btn danger" title="Delete" onclick={() => deleteRepo(repo.id, repo.name)}>
                  <span class="material-symbols-rounded">delete</span>
                </button>
              </td>
            </tr>
          {/each}
        {/if}
      </tbody>
    </table>
  </div>
</div>

<style>
  .admin-repos { padding: 2rem; display: flex; flex-direction: column; gap: 2rem; box-shadow: var(--shadow-xl); }
  
  .section-header { display: flex; justify-content: space-between; align-items: center; }
  .title { display: flex; align-items: center; gap: 0.75rem; }
  .title span { font-size: 1.75rem; color: var(--primary-color); }
  .title h2 { margin: 0; font-size: 1.5rem; font-weight: 800; }

  .table-container { overflow-x: auto; }
  .modern-table { width: 100%; border-collapse: collapse; text-align: left; }
  .modern-table th { padding: 1rem; border-bottom: 1px solid var(--border-color); color: var(--text-muted); font-size: 0.75rem; text-transform: uppercase; letter-spacing: 0.1em; font-weight: 800; }
  .modern-table td { padding: 1.25rem 1rem; border-bottom: 1px solid rgba(255,255,255,0.03); font-size: 0.95rem; }

  .repo-cell { display: flex; flex-direction: column; gap: 0.1rem; }
  .repo-name { font-size: 1rem; color: #fff; font-weight: 700; }
  .repo-desc { font-size: 0.8rem; opacity: 0.6; display: -webkit-box; -webkit-line-clamp: 1; line-clamp: 1; -webkit-box-orient: vertical; overflow: hidden; }

  .status-badge { font-size: 0.7rem; font-weight: 800; text-transform: uppercase; padding: 0.2rem 0.6rem; border-radius: 4px; background: rgba(34, 197, 94, 0.1); color: #22c55e; }
  .status-badge.private { background: rgba(239, 68, 68, 0.1); color: #ef4444; }

  .actions { display: flex; gap: 0.5rem; }
  .icon-btn { background: transparent; border: none; color: var(--text-muted); cursor: pointer; padding: 0.4rem; border-radius: 0.5rem; transition: all 0.2s; text-decoration: none; display: flex; align-items: center; }
  .icon-btn:hover { background: rgba(255,255,255,0.05); color: #fff; }
  .icon-btn.danger:hover { background: rgba(239, 68, 68, 0.1); color: #ef4444; }

  .skeleton-bar { height: 1.5rem; background: rgba(255,255,255,0.03); border-radius: 0.5rem; animation: pulse 1.5s infinite; }
  @keyframes pulse { 0% { opacity: 0.3; } 50% { opacity: 0.6; } 100% { opacity: 0.3; } }
</style>
