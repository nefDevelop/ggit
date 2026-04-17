<script lang="ts">
  import { onMount } from 'svelte';
  
  let { q = '' } = $props();
  
  let repos = $state<any[]>([]);
  let loading = $state(true);

  async function fetchRepos() {
    loading = true;
    try {
      const res = await fetch(`/api/v2/explore/repos?q=${encodeURIComponent(q)}`);
      if (res.ok) {
        const data = await res.json();
        repos = data.repos || [];
      }
    } catch (e) {
      console.error(e);
    } finally {
      loading = false;
    }
  }

  $effect(() => {
    fetchRepos();
  });

  onMount(fetchRepos);
</script>

<div class="explore-repos">
  {#if loading && repos.length === 0}
    <div class="skeleton-grid">
      {#each Array(6) as _}
        <div class="skeleton-card island"></div>
      {/each}
    </div>
  {:else if repos.length > 0}
    <div class="repos-grid">
      {#each repos as repo}
        <a href="/{repo.owner?.username || 'unknown'}/{repo.name}" class="repo-card island tonal hover-lift">
          <div class="repo-header">
            <span class="material-symbols-rounded repo-icon">
              {repo.private ? 'lock' : repo.fork ? 'fork_right' : 'book'}
            </span>
            <div class="repo-names">
              <span class="owner">{repo.owner?.username || 'unknown'} /</span>
              <strong class="name">{repo.name}</strong>
            </div>
          </div>
          <p class="repo-desc">{repo.description || 'No description provided.'}</p>
          <div class="repo-stats">
            <div class="stat">
              <span class="material-symbols-rounded">star</span>
              {repo.stars_count}
            </div>
            <div class="stat">
              <span class="material-symbols-rounded">fork_right</span>
              {repo.forks_count}
            </div>
            {#if repo.language}
              <div class="stat">
                <span class="lang-dot"></span>
                {repo.language}
              </div>
            {/if}
            <span class="updated">Updated {new Date(repo.updated_at).toLocaleDateString()}</span>
          </div>
        </a>
      {/each}
    </div>
  {:else}
    <div class="empty-state island tonal">
      <span class="material-symbols-rounded">search_off</span>
      <h3>No repositories found</h3>
      <p>Try searching for something else.</p>
    </div>
  {/if}
</div>

<style>
  .repos-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(350px, 1fr)); gap: 1.5rem; }
  .repo-card { display: flex; flex-direction: column; gap: 1rem; padding: 2rem; text-decoration: none; color: inherit; height: 100%; border: 1px solid var(--border-color); transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1); }
  .repo-card:hover { border-color: var(--primary-color); background: rgba(var(--primary-rgb), 0.05); }

  .repo-header { display: flex; align-items: center; gap: 1rem; }
  .repo-icon { font-size: 1.75rem; color: var(--primary-color); opacity: 0.8; }
  .repo-names { font-size: 1.25rem; font-weight: 500; }
  .owner { opacity: 0.6; font-weight: 400; }
  .name { font-weight: 800; color: #fff; }

  .repo-desc { font-size: 0.95rem; color: var(--text-muted); margin: 0; line-height: 1.6; flex: 1; display: -webkit-box; -webkit-line-clamp: 3; line-clamp: 3; -webkit-box-orient: vertical; overflow: hidden; }

  .repo-stats { display: flex; align-items: center; gap: 1.5rem; font-size: 0.85rem; color: var(--text-muted); margin-top: 1.5rem; padding-top: 1rem; border-top: 1px solid var(--border-color); }
  .stat { display: flex; align-items: center; gap: 0.4rem; font-weight: 700; }
  .stat span { font-size: 1.1rem; color: var(--primary-color); }
  .lang-dot { width: 10px; height: 10px; border-radius: 50%; background: var(--primary-color); box-shadow: 0 0 8px var(--primary-color); }
  .updated { margin-left: auto; opacity: 0.4; font-size: 0.75rem; }

  .empty-state { padding: 5rem; text-align: center; display: flex; flex-direction: column; align-items: center; gap: 1rem; }
  .empty-state span { font-size: 4rem; color: var(--primary-color); opacity: 0.3; }
  .empty-state h3 { margin: 0; }
  .empty-state p { color: var(--text-muted); }

  .skeleton-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(350px, 1fr)); gap: 1.5rem; }
  .skeleton-card { height: 180px; animation: pulse 1.5s infinite; }
  @keyframes pulse { 0% { opacity: 0.3; } 50% { opacity: 0.5; } 100% { opacity: 0.3; } }
</style>
