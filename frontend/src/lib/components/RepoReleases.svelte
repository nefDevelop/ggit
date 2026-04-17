<script lang="ts">
  import { apiFetch } from "$lib/api";
  import { onMount } from 'svelte';
  
  let { username, reponame } = $props();
  
  let releases = $state<any[]>([]);
  let loading = $state(true);

  async function fetchReleases() {
    loading = true;
    try {
      const res = await apiFetch(`/api/v2/${username}/${reponame}/releases`);
      if (res.ok) releases = await res.json();
    } catch (e) {
      console.error(e);
    } finally {
      loading = false;
    }
  }

  onMount(fetchReleases);
</script>

<div class="repo-releases">
  <header class="section-actions">
     <button class="btn btn-primary">
       <span class="material-symbols-rounded">publish</span>
       Draft new release
     </button>
  </header>

  {#if loading}
    <div class="skeleton-list"></div>
  {:else if releases.length === 0}
    <div class="empty-state island">
      <span class="material-symbols-rounded">sell</span>
      <h3>No releases yet</h3>
      <p>Tags and releases will appear here once you create them.</p>
    </div>
  {:else}
    <div class="releases-timeline">
      {#each releases as release}
        <div class="release-item">
          <div class="release-marker"></div>
          <div class="release-card island tonal">
            <header class="release-header">
               <div class="title-row">
                 <h2>{release.Title || release.TagName}</h2>
                 {#if release.IsDraft}
                   <span class="badge draft">Draft</span>
                 {:else if release.IsPrerelease}
                   <span class="badge pre">Pre-release</span>
                 {:else}
                   <span class="badge latest">Latest</span>
                 {/if}
               </div>
               <div class="meta-row">
                 <span class="tag">
                   <span class="material-symbols-rounded">sell</span>
                   {release.TagName}
                 </span>
                 <span class="target">
                   <span class="material-symbols-rounded">merge_type</span>
                   {release.Target}
                 </span>
                 <span class="date">
                   <span class="material-symbols-rounded">calendar_today</span>
                   {new Date(release.created_at || Date.now()).toLocaleDateString()}
                 </span>
               </div>
            </header>

            <div class="release-content">
              {release.Note || 'No description provided.'}
            </div>

            <div class="release-assets">
               <h4>Assets</h4>
               <div class="asset-links">
                 <a href="#" class="asset-item">
                   <span class="material-symbols-rounded">zip_box</span>
                   Source code (zip)
                 </a>
                 <a href="#" class="asset-item">
                   <span class="material-symbols-rounded">archive</span>
                   Source code (tar.gz)
                 </a>
               </div>
            </div>
          </div>
        </div>
      {/each}
    </div>
  {/if}
</div>

<style>
  .repo-releases { display: flex; flex-direction: column; gap: 2rem; }
  .section-actions { display: flex; justify-content: flex-end; }
  
  .empty-state { text-align: center; padding: 4rem; color: var(--text-muted); }
  .empty-state span { font-size: 4rem; opacity: 0.2; }

  .releases-timeline { display: flex; flex-direction: column; position: relative; padding-left: 2rem; }
  .releases-timeline::before { content: ''; position: absolute; left: 0; top: 0; bottom: 0; width: 2px; background: var(--border-color); }
  
  .release-item { position: relative; padding-bottom: 3rem; }
  .release-marker { position: absolute; left: -2.35rem; top: 1.5rem; width: 12px; height: 12px; border-radius: 50%; background: var(--primary-color); border: 3px solid var(--bg-color); z-index: 1; }

  .release-card { padding: 2rem; display: flex; flex-direction: column; gap: 1.5rem; }
  
  .title-row { display: flex; align-items: center; gap: 1rem; }
  .title-row h2 { margin: 0; font-size: 1.5rem; font-weight: 800; }
  
  .badge { font-size: 0.75rem; font-weight: 800; text-transform: uppercase; padding: 0.2rem 0.6rem; border-radius: 4px; }
  .badge.latest { background: rgba(16, 185, 129, 0.1); color: #10b981; }
  .badge.pre { background: rgba(245, 158, 11, 0.1); color: #f59e0b; }
  .badge.draft { background: rgba(255, 255, 255, 0.05); color: var(--text-muted); }

  .meta-row { display: flex; gap: 1.5rem; }
  .meta-row span { display: flex; align-items: center; gap: 0.4rem; font-size: 0.85rem; color: var(--text-muted); font-weight: 600; }
  .meta-row span span { font-size: 1.1rem; opacity: 0.7; }

  .release-content { font-size: 1rem; line-height: 1.6; color: var(--text-color); padding: 1rem; background: rgba(0,0,0,0.1); border-radius: 0.75rem; border: 1px solid var(--border-color); white-space: pre-wrap; }

  .release-assets h4 { margin: 0 0 1rem 0; font-size: 0.9rem; text-transform: uppercase; color: var(--text-muted); letter-spacing: 0.05em; }
  .asset-links { display: flex; flex-direction: column; gap: 0.5rem; }
  .asset-item { display: flex; align-items: center; gap: 0.75rem; text-decoration: none; color: var(--primary-color); font-weight: 600; font-size: 0.95rem; padding: 0.5rem; border-radius: 0.5rem; transition: background 0.2s; }
  .asset-item:hover { background: rgba(var(--primary-rgb), 0.05); }
  .asset-item span { font-size: 1.25rem; }

  .skeleton-list { height: 400px; animation: pulse 1.5s infinite; }
  @keyframes pulse { 0%, 100% { opacity: 0.3; } 50% { opacity: 0.6; } }
</style>
