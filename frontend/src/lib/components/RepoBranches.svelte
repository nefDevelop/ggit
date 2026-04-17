<script lang="ts">
  import { apiFetch } from "$lib/api";
  import { onMount } from 'svelte';
  
  let { username, reponame } = $props();
  
  let data = $state<any>(null);
  let loading = $state(true);

  async function fetchBranches() {
    loading = true;
    try {
      const res = await apiFetch(`/api/v2/${username}/${reponame}/branches`);
      if (res.ok) data = await res.json();
    } catch (e) {
      console.error(e);
    } finally {
      loading = false;
    }
  }

  async function deleteBranch(name: string) {
    if (!confirm(`Are you sure you want to delete branch "${name}"?`)) return;
    try {
      const res = await apiFetch(`/api/v2/${username}/${reponame}/branches/${name}`, { method: 'DELETE' });
      if (res.ok) fetchBranches();
    } catch (e) {
      console.error(e);
    }
  }

  onMount(fetchBranches);
</script>

<div class="repo-branches">
  <header class="section-header">
    <h2>Branches</h2>
    <p class="subtitle">Branches are where you develop features and fix bugs.</p>
  </header>

  {#if loading}
    <div class="skeleton-table"></div>
  {:else if data}
    <!-- Default Branch -->
    <section class="branch-section island tonal">
      <h3 class="section-title">Default Branch</h3>
      <div class="branch-row">
        <div class="branch-info">
          <span class="material-symbols-rounded icon">merge_type</span>
          <div class="names">
            <strong>{data.default.name}</strong>
            <span class="desc">Default branch for the repository.</span>
          </div>
          {#if data.default.is_protected}
            <span class="badge protected">Protected</span>
          {/if}
        </div>
        <div class="branch-meta">
           <span class="date">Updated {data.default.commit?.committer?.when ? new Date(data.default.commit.committer.when).toLocaleDateString() : 'Unknown'}</span>
           <span class="btn btn-outline disabled">Delete</span>
        </div>
      </div>
    </section>

    <!-- Active Branches -->
    {#if data.active.length > 0}
      <section class="branch-section island">
        <h3 class="section-title">Active Branches</h3>
        {#each data.active as branch}
          <div class="branch-row">
            <div class="branch-info">
              <span class="material-symbols-rounded icon">merge_type</span>
              <div class="names">
                <strong>{branch.name}</strong>
                <span class="commit-msg">{branch.commit.summary}</span>
              </div>
              {#if branch.is_protected}
                <span class="badge protected">Protected</span>
              {/if}
            </div>
            <div class="branch-meta">
               <span class="date">{branch.commit?.committer?.when ? new Date(branch.commit.committer.when).toLocaleDateString() : 'Recently'}</span>
               <button class="btn btn-danger-text" disabled={branch.is_protected} onclick={() => deleteBranch(branch.name)}>
                 <span class="material-symbols-rounded">delete</span>
               </button>
            </div>
          </div>
        {/each}
      </section>
    {/if}

    <!-- Stale Branches -->
    {#if data.stale.length > 0}
      <section class="branch-section island dimmed">
        <h3 class="section-title">Stale Branches</h3>
        {#each data.stale as branch}
          <div class="branch-row">
            <div class="branch-info">
              <span class="material-symbols-rounded icon">merge_type</span>
              <div class="names">
                <strong>{branch.name}</strong>
                <span class="commit-msg">{branch.commit.summary}</span>
              </div>
            </div>
            <div class="branch-meta">
               <span class="date">Old: {branch.commit?.committer?.when ? new Date(branch.commit.committer.when).toLocaleDateString() : 'N/A'}</span>
               <button class="btn btn-danger-text" onclick={() => deleteBranch(branch.name)}>
                 <span class="material-symbols-rounded">delete</span>
               </button>
            </div>
          </div>
        {/each}
      </section>
    {/if}
  {/if}
</div>

<style>
  .repo-branches { display: flex; flex-direction: column; gap: 2.5rem; }
  
  .section-header h2 { margin: 0; font-size: 1.75rem; font-weight: 800; }
  .subtitle { color: var(--text-muted); margin: 0.25rem 0 0 0; }

  .branch-section { padding: 0; overflow: hidden; display: flex; flex-direction: column; }
  .section-title { padding: 1.25rem 2rem; margin: 0; font-size: 1rem; background: rgba(255,255,255,0.02); border-bottom: 1px solid var(--border-color); }

  .branch-row { display: flex; align-items: center; justify-content: space-between; padding: 1.25rem 2rem; border-bottom: 1px solid var(--border-color); }
  .branch-row:last-child { border-bottom: none; }

  .branch-info { display: flex; align-items: center; gap: 1.25rem; }
  .icon { color: var(--primary-color); opacity: 0.6; }
  
  .names { display: flex; flex-direction: column; }
  .names strong { font-size: 1.1rem; }
  .desc, .commit-msg { font-size: 0.85rem; color: var(--text-muted); }

  .badge { font-size: 0.7rem; font-weight: 800; text-transform: uppercase; background: rgba(var(--primary-rgb), 0.1); color: var(--primary-color); padding: 0.2rem 0.6rem; border-radius: 4px; margin-left: 0.5rem; }

  .branch-meta { display: flex; align-items: center; gap: 2rem; }
  .branch-meta .date { font-size: 0.85rem; color: var(--text-muted); font-weight: 600; min-width: 100px; text-align: right; }

  .btn-danger-text { color: #ef4444; border: none; background: transparent; cursor: pointer; padding: 0.5rem; border-radius: 0.5rem; transition: background 0.2s; }
  .btn-danger-text:hover:not(:disabled) { background: rgba(239, 68, 68, 0.1); }
  .btn-danger-text:disabled { opacity: 0.3; cursor: not-allowed; }

  .dimmed { opacity: 0.6; grayscale: 0.5; }
  .skeleton-table { height: 400px; animation: pulse 1.5s infinite; }
  @keyframes pulse { 0%, 100% { opacity: 0.3; } 50% { opacity: 0.6; } }
</style>
