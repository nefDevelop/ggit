<script lang="ts">
  import { apiFetch } from "$lib/api";
  import { onMount } from 'svelte';
  import ReleaseEditor from './ReleaseEditor.svelte';
  
  let { username, reponame } = $props();
  
  let releases = $state<any[]>([]);
  let loading = $state(true);
  let error = $state("");
  
  let hasNext = $state(false);
  let nextAfter = $state("");
  let hasPrevious = $state(false);

  let showEditor = $state(false);
  let editingRelease = $state<any>(null);

  async function fetchReleases(after = "") {
    loading = true;
    error = "";
    try {
      const query = after ? `?after=${after}` : "";
      const res = await apiFetch(`/api/v2/${username}/${reponame}/releases${query}`);
      if (res.ok) {
        const data = await res.json();
        releases = data.releases;
        hasNext = data.has_next;
        nextAfter = data.next_after;
        hasPrevious = data.has_previous;
      } else {
        error = "Failed to load releases.";
      }
    } catch (e) {
      console.error(e);
      error = "Network error. Please try again.";
    } finally {
      loading = false;
    }
  }

  function handleAdd() {
    editingRelease = null;
    showEditor = true;
  }

  function handleEdit(rel: any) {
    editingRelease = rel;
    showEditor = true;
  }

  async function handleDelete(id: number) {
    if (!confirm("Are you sure you want to delete this release? This will also remove the associated Git tag if it was created through Ggit.")) return;
    
    try {
      const res = await apiFetch(`/api/v2/${username}/${reponame}/releases/${id}`, {
        method: "DELETE"
      });
      if (res.ok) {
        releases = releases.filter(r => r.id !== id);
      }
    } catch (e) {
      console.error(e);
    }
  }

  function onEditorSaved(newRel: any) {
    showEditor = false;
    fetchReleases();
  }

  onMount(() => fetchReleases());
</script>

<div class="repo-releases">
  {#if showEditor}
    <div class="editor-overlay">
      <ReleaseEditor 
        {username} 
        {reponame} 
        release={editingRelease} 
        onSaved={onEditorSaved} 
        onCancel={() => showEditor = false} 
      />
    </div>
  {/if}

  <header class="section-actions">
     <div class="header-info">
       <h3>Versions & Releases</h3>
       <p>History of tags and published versions.</p>
     </div>
     <button class="btn btn-primary lg" onclick={handleAdd}>
       <span class="material-symbols-rounded">add_circle</span>
       Draft new release
     </button>
  </header>

  {#if loading && releases.length === 0}
    <div class="skeleton-timeline">
      {#each Array(3) as _}
        <div class="skeleton-item island tonal"></div>
      {/each}
    </div>
  {:else if error}
    <div class="error-state island tonal">
      <span class="material-symbols-rounded">error</span>
      <p>{error}</p>
      <button class="btn btn-tonal" onclick={() => fetchReleases()}>Retry</button>
    </div>
  {:else if releases.length === 0}
    <div class="empty-state island tonal">
      <div class="icon-blob">
        <span class="material-symbols-rounded">sell</span>
      </div>
      <h3>No releases yet</h3>
      <p>Start tagging your milestones to see them here.</p>
      <button class="btn btn-primary" onclick={handleAdd}>Create your first release</button>
    </div>
  {:else}
    <div class="releases-timeline">
      {#each releases as release}
        <div class="release-item" class:is-draft={release.is_draft}>
          <div class="release-marker"></div>
          
          <div class="release-card island tonal">
            <header class="release-header">
               <div class="header-main">
                 <div class="title-row">
                   <h2>{release.title || release.tag_name}</h2>
                   <div class="badges">
                     {#if release.is_draft}
                       <span class="badge draft">Draft</span>
                     {:else if release.is_prerelease}
                       <span class="badge pre">Pre-release</span>
                     {/if}
                     {#if release === releases.find(r => !r.is_draft && !r.is_prerelease)}
                       <span class="badge latest">Latest</span>
                     {/if}
                   </div>
                 </div>

                 <div class="meta-row">
                   <span class="tag clickable" title="View tag on Git">
                     <span class="material-symbols-rounded">sell</span>
                     {release.tag_name}
                   </span>
                   {#if release.target}
                     <span class="target">
                       <span class="material-symbols-rounded">merge_type</span>
                       {release.target}
                     </span>
                   {/if}
                   <span class="date">
                     <span class="material-symbols-rounded">calendar_today</span>
                     {new Date(release.created_at * 1000).toLocaleDateString(undefined, { dateStyle: 'long' })}
                   </span>
                 </div>
               </div>

               {#if release.id}
                 <div class="header-actions">
                   <button class="btn-icon" onclick={() => handleEdit(release)} title="Edit Release">
                     <span class="material-symbols-rounded">edit</span>
                   </button>
                   <button class="btn-icon delete" onclick={() => handleDelete(release.id)} title="Delete Release">
                     <span class="material-symbols-rounded">delete</span>
                   </button>
                 </div>
               {/if}
            </header>

            <div class="release-body">
              <div class="release-content markdown-body">
                {#if release.note}
                  {release.note}
                {:else}
                  <i class="empty-note">No description available for this version.</i>
                {/if}
              </div>

              <div class="release-side">
                <div class="assets-box island tonal surface-2">
                   <h4>Assets</h4>
                   <div class="asset-links">
                     {#if release.attachments && release.attachments.length > 0}
                       {#each release.attachments as asset}
                         <a href={asset.download_url} class="asset-item" download>
                           <span class="material-symbols-rounded">get_app</span>
                           <span class="name">{asset.name}</span>
                         </a>
                       {/each}
                     {/if}
                     <a href={`/${username}/${reponame}/archive/${release.tag_name}.zip`} class="asset-item source">
                       <span class="material-symbols-rounded">zip_box</span>
                       <span class="name">Source code (zip)</span>
                     </a>
                     <a href={`/${username}/${reponame}/archive/${release.tag_name}.tar.gz`} class="asset-item source">
                       <span class="material-symbols-rounded">archive</span>
                       <span class="name">Source code (tar.gz)</span>
                     </a>
                   </div>
                </div>
                
                {#if release.publisher}
                  <div class="publisher-info">
                    <span>Published by</span>
                    <a href={`/${release.publisher.username}`} class="user-chip">
                      <img src={release.publisher.avatar_url} alt={release.publisher.username} />
                      {release.publisher.username}
                    </a>
                  </div>
                {/if}
              </div>
            </div>
          </div>
        </div>
      {/each}
    </div>

    {#if hasNext || hasPrevious}
      <div class="pagination">
        <button class="btn btn-tonal" onclick={() => fetchReleases("")} disabled={loading || !hasPrevious}>
          <span class="material-symbols-rounded">first_page</span>
          Newest
        </button>
        <button class="btn btn-primary" onclick={() => fetchReleases(nextAfter)} disabled={loading || !hasNext}>
          Older releases
          <span class="material-symbols-rounded">navigate_next</span>
        </button>
      </div>
    {/if}
  {/if}
</div>

<style>
  .repo-releases { display: flex; flex-direction: column; gap: 3rem; min-height: 400px; padding: 1rem; }
  
  .editor-overlay { position: fixed; inset: 0; background: rgba(0,0,0,0.8); backdrop-filter: blur(8px); z-index: 1000; display: flex; align-items: flex-start; justify-content: center; padding: 4rem; overflow-y: auto; }

  .section-actions { display: flex; justify-content: space-between; align-items: flex-start; }
  .header-info h3 { margin: 0; font-size: 1.75rem; font-weight: 800; }
  .header-info p { margin: 0.5rem 0 0 0; color: var(--text-muted); font-size: 1rem; }
  .btn.lg { height: 56px; padding: 0 1.5rem; gap: 0.75rem; border-radius: 1.25rem; font-weight: 800; }

  .empty-state { text-align: center; padding: 6rem 2rem; display: flex; flex-direction: column; align-items: center; gap: 1.5rem; color: var(--text-muted); border-radius: 2rem; }
  .icon-blob { width: 120px; height: 120px; background: rgba(var(--primary-rgb), 0.1); border-radius: 40%; display: flex; align-items: center; justify-content: center; margin-bottom: 1rem; }
  .empty-state span { font-size: 4rem; color: var(--primary-color); }
  .empty-state h3 { margin: 0; font-size: 1.5rem; color: var(--text-color); }

  .releases-timeline { display: flex; flex-direction: column; position: relative; padding-left: 2.5rem; gap: 1rem; }
  .releases-timeline::before { content: ''; position: absolute; left: 0; top: 1rem; bottom: 1rem; width: 4px; background: var(--md3-surface-5); border-radius: 2px; }
  
  .release-item { position: relative; padding-bottom: 2rem; transition: transform 0.2s; }
  .release-item.is-draft { opacity: 0.8; }
  .release-marker { position: absolute; left: -2.85rem; top: 2rem; width: 14px; height: 14px; border-radius: 50%; background: var(--md3-surface-5); border: 4px solid var(--bg-color); z-index: 1; transition: all 0.3s; }
  .release-item:hover .release-marker { background: var(--primary-color); transform: scale(1.2); }

  .release-card { padding: 0; display: flex; flex-direction: column; overflow: hidden; border: 1px solid transparent; transition: all 0.3s; }
  .release-card:hover { border-color: var(--primary-color); background: var(--md3-surface-2); box-shadow: 0 8px 32px rgba(0,0,0,0.2); }
  
  .release-header { padding: 2rem; display: flex; justify-content: space-between; align-items: flex-start; background: rgba(255,255,255,0.02); }
  .header-main { display: flex; flex-direction: column; gap: 1rem; }
  
  .title-row { display: flex; align-items: center; gap: 1.5rem; flex-wrap: wrap; }
  .title-row h2 { margin: 0; font-size: 1.8rem; font-weight: 800; letter-spacing: -0.02em; }
  
  .badges { display: flex; gap: 0.5rem; }
  .badge { font-size: 0.75rem; font-weight: 800; text-transform: uppercase; padding: 0.3rem 0.8rem; border-radius: 0.5rem; }
  .badge.latest { background: #10b981; color: #fff; box-shadow: 0 0 16px rgba(16, 185, 129, 0.3); }
  .badge.pre { background: #f59e0b; color: #fff; }
  .badge.draft { background: var(--md3-surface-5); color: var(--text-muted); }

  .meta-row { display: flex; gap: 2rem; flex-wrap: wrap; }
  .meta-row span { display: flex; align-items: center; gap: 0.5rem; font-size: 0.9rem; color: var(--text-muted); font-weight: 600; }
  .meta-row .material-symbols-rounded { font-size: 1.25rem; opacity: 0.7; }
  .meta-row .tag.clickable { color: var(--primary-color); cursor: pointer; }
  .meta-row .tag.clickable:hover { text-decoration: underline; }

  .header-actions { display: flex; gap: 0.5rem; }
  .btn-icon.delete:hover { background: rgba(239, 68, 68, 0.1); color: #ef4444; }

  .release-body { padding: 2rem; display: grid; grid-template-columns: 1fr 300px; gap: 3rem; border-top: 1px solid var(--border-color); }
  
  .release-content { font-size: 1.05rem; line-height: 1.8; color: var(--text-color); white-space: pre-wrap; }
  .empty-note { opacity: 0.4; }

  .release-side { display: flex; flex-direction: column; gap: 2rem; }
  
  .assets-box { padding: 1.5rem; border-radius: 1.25rem; }
  .assets-box h4 { margin: 0 0 1rem 0; font-size: 0.85rem; text-transform: uppercase; color: var(--text-muted); letter-spacing: 0.1em; font-weight: 800; }
  .asset-links { display: flex; flex-direction: column; gap: 0.25rem; }
  .asset-item { display: flex; align-items: center; gap: 0.75rem; text-decoration: none; color: var(--primary-color); font-weight: 700; font-size: 0.95rem; padding: 0.75rem; border-radius: 0.75rem; transition: all 0.2s; }
  .asset-item:hover { background: rgba(var(--primary-rgb), 0.1); transform: translateX(4px); }
  .asset-item.source { color: var(--text-muted); font-weight: 600; }
  .asset-item span { font-size: 1.4rem; opacity: 0.8; }

  .publisher-info { display: flex; flex-direction: column; gap: 0.75rem; padding-left: 0.5rem; }
  .publisher-info span { font-size: 0.8rem; font-weight: 700; text-transform: uppercase; color: var(--text-muted); letter-spacing: 0.05em; }
  .user-chip { display: flex; align-items: center; gap: 0.75rem; font-weight: 700; color: var(--text-color); text-decoration: none; }
  .user-chip img { width: 28px; height: 28px; border-radius: 50%; background: var(--md3-surface-5); }

  .pagination { display: flex; justify-content: center; gap: 1rem; margin-top: 2rem; }
  
  .error-state { text-align: center; padding: 4rem; display: flex; flex-direction: column; align-items: center; gap: 1rem; color: #ef4444; }
  .error-state span { font-size: 3rem; }

  .skeleton-timeline { display: flex; flex-direction: column; gap: 1.5rem; padding-left: 2.5rem; }
  .skeleton-item { height: 250px; border-radius: 1.5rem; animation: pulse 1.5s infinite; }
  @keyframes pulse { 0%, 100% { opacity: 0.4; } 50% { opacity: 0.7; } }

  @media (max-width: 1100px) {
    .release-body { grid-template-columns: 1fr; }
    .editor-overlay { padding: 1rem; }
  }
</style>
