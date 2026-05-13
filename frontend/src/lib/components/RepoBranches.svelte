<script lang="ts">
  import { apiFetch } from "$lib/api";
  import { onMount } from 'svelte';
  import { dialog } from '$lib/dialog.svelte';
  import { Star, GitFork, History, Trash2, AlertCircle } from 'lucide-svelte';
  
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
    if (!await dialog.confirm(`Are you sure you want to delete branch "${name}"?`)) return;
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
  <div class="branches-content">
  {#if loading}
    <div class="skeleton-table"></div>
  {:else if data}
    <!-- Primary Branch -->
    <div class="branch-list-container">
      {#if data.default}
        <div class="branch-row-item primary-row">
          <div class="branch-main-info">
            <Star size={18} class="branch-icon-acc" />
            <div class="branch-name-stack">
              <a href={`/${username}/${reponame}/src/branch/${data.default.name}`} class="branch-name-text">{data.default.name}</a>
              <div class="badges-row">
                <span class="badge-v2 primary">Primary</span>
                {#if data.default.is_protected}
                  <span class="badge-v2 protected">Protected</span>
                {/if}
              </div>
            </div>
          </div>
          <div class="branch-meta-info">
             <span class="meta-label">Updated {data.default.commit?.committer?.when ? new Date(data.default.commit.committer.when).toLocaleDateString() : 'No commits'}</span>
          </div>
        </div>
      {:else}
        <div class="empty-state-simple">No branches found.</div>
      {/if}

      <!-- Active Branches -->
      {#if data.active.length > 0}
        <div class="list-divider">Active Development</div>
        {#each data.active as branch}
          <div class="branch-row-item">
            <div class="branch-main-info">
              <GitFork size={18} class="branch-icon" />
              <div class="branch-name-stack">
                <a href={`/${username}/${reponame}/src/branch/${branch.name}`} class="branch-name-text">{branch.name}</a>
                {#if branch.is_protected}
                  <span class="badge-v2 protected">Protected</span>
                {/if}
                <span class="commit-shorthand">{branch.commit?.summary || ''}</span>
              </div>
            </div>
            <div class="branch-meta-info">
               <span class="meta-label">{branch.commit?.committer?.when ? new Date(branch.commit.committer.when).toLocaleDateString() : 'Recently'}</span>
               <button class="btn-icon-minimal delete" disabled={branch.is_protected} onclick={() => deleteBranch(branch.name)}>
                 <Trash2 size={16} />
               </button>
            </div>
          </div>
        {/each}
      {/if}

      <!-- Stale Branches -->
      {#if data.stale.length > 0}
        <div class="list-divider">Legacy (Stale)</div>
        {#each data.stale as branch}
          <div class="branch-row-item stale-row">
            <div class="branch-main-info">
              <History size={18} class="branch-icon" />
              <div class="branch-name-stack">
                <a href={`/${username}/${reponame}/src/branch/${branch.name}`} class="branch-name-text">{branch.name}</a>
              </div>
            </div>
            <div class="branch-meta-info">
               <span class="meta-label">Old: {branch.commit?.committer?.when ? new Date(branch.commit.committer.when).toLocaleDateString() : 'N/A'}</span>
               <button class="btn-icon-minimal delete" onclick={() => deleteBranch(branch.name)}>
                 <Trash2 size={16} />
               </button>
            </div>
          </div>
        {/each}
      {/if}
    </div>
  {:else}
    <div class="empty-state">
      <AlertCircle size={48} />
      <h3>No branches found</h3>
      <p>This repository has no branches yet.</p>
    </div>
  {/if}
  </div>
</div>

<style>
  .repo-branches { display: flex; flex-direction: column; gap: 0; }
  .branches-content { padding: 0; }
  
  .branch-list-container {
    border: 1px solid var(--md-sys-color-outline-variant);
    background: var(--bg-color);
    border-radius: 1rem;
    overflow: hidden;
  }

  .branch-row-item {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 1rem 1.75rem;
    border-bottom: 1px solid var(--md-sys-color-outline-variant);
    transition: background 0.2s;
  }
  .branch-row-item:last-child { border-bottom: none; }
  .branch-row-item:hover { background: rgba(var(--primary-rgb), 0.04); }

  .primary-row { background: rgba(var(--primary-rgb), 0.02); }
  .stale-row { opacity: 0.6; }

  .list-divider {
    background: var(--surface-variant);
    padding: 0.5rem 1.75rem;
    font-size: 0.7rem;
    font-weight: 900;
    text-transform: uppercase;
    color: var(--text-muted);
    letter-spacing: 0.1em;
    border-bottom: 1px solid var(--md-sys-color-outline-variant);
  }

  .branch-main-info { display: flex; align-items: center; gap: 1.25rem; }
  :global(.branch-icon) { color: var(--text-muted); opacity: 0.5; }
  :global(.branch-icon-acc) { color: var(--primary-color); }
  
  .branch-name-stack { display: flex; align-items: center; gap: 0.75rem; }
  .branch-name-text { font-weight: 800; font-size: 1.05rem; color: var(--text-color); text-decoration: none; transition: color 0.2s; }
  .branch-name-text:hover { color: var(--primary-color); }
  .commit-shorthand { font-size: 0.8rem; color: var(--text-muted); opacity: 0.7; font-weight: 500; max-width: 250px; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }

  .badges-row { display: flex; gap: 0.5rem; }
  .badge-v2 { font-size: 0.6rem; font-weight: 900; text-transform: uppercase; padding: 2px 8px; border-radius: 20px; letter-spacing: 0.05em; }
  .badge-v2.primary { background: var(--primary-color); color: var(--on-primary); }
  .badge-v2.protected { border: 1px solid var(--md-sys-color-outline-variant); color: var(--text-muted); }

  .branch-meta-info { display: flex; align-items: center; gap: 1.5rem; }
  .meta-label { font-size: 0.8rem; font-weight: 700; color: var(--text-muted); }

  .skeleton-table { height: 400px; background: var(--surface-variant); border-radius: 1rem; animation: pulse 1.5s infinite; }
  @keyframes pulse { 0%, 100% { opacity: 0.3; } 50% { opacity: 0.6; } }

  .empty-state { display: flex; flex-direction: column; align-items: center; justify-content: center; gap: 1rem; padding: 4rem 2rem; text-align: center; color: var(--text-muted); opacity: 0.6; }
  .empty-state :global(svg) { font-size: 3rem; }

  @media (max-width: 600px) {
    .branch-row-item {
      padding: 0.75rem 1rem;
      flex-direction: column;
      align-items: flex-start;
      gap: 0.5rem;
    }
    .branch-meta-info {
      width: 100%;
      justify-content: space-between;
    }
    .branch-main-info {
      gap: 0.75rem;
    }
    .list-divider {
      padding: 0.5rem 1rem;
    }
  }
</style>
