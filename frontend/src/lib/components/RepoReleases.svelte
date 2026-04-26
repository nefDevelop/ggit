<script lang="ts">
  import { apiFetch } from "$lib/api";
  import { onMount } from 'svelte';
  import { dialog } from '$lib/dialog.svelte';
  import Avatar from '$lib/components/Avatar.svelte';
  import ReleaseEditor from './ReleaseEditor.svelte';
  import Button from './ui/Button.svelte';
  import Card from './ui/Card.svelte';
  import Badge from './ui/Badge.svelte';
  
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
    if (!await dialog.confirm("Are you sure you want to delete this release? This will also remove the associated Git tag if it was created through Ggit.")) return;
    
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

  <header class="section-actions-minimal">
     <Button onclick={handleAdd} icon="add_circle">
       Draft Release
     </Button>
  </header>

  {#if loading && releases.length === 0}
    <div class="skeleton-list">
      {#each Array(3) as _}
        <div class="skeleton-item"></div>
      {/each}
    </div>
  {:else if error}
    <div class="error-msg">{error}</div>
  {:else if releases.length === 0}
    <div class="empty-releases">
      <span class="material-symbols-rounded icon">sell</span>
      <h3>No releases yet</h3>
      <p>Start tagging your milestones.</p>
      <Button onclick={handleAdd} icon="rocket_launch">Create Release</Button>
    </div>
  {:else}
    <div class="releases-timeline">
      {#each releases as release}
        <div class="release-item" class:is-draft={release.is_draft}>
          <div class="release-marker"></div>
          
          <Card class="release-card" hover={true} padding="none">
            <header class="release-card-header">
               <div class="header-main-info">
                 <div class="title-row-flat">
                   <h3>{release.title || release.tag_name}</h3>
                   <div class="badge-row">
                     {#if release.is_draft}
                       <Badge variant="neutral">Draft</Badge>
                     {:else if release.is_prerelease}
                       <Badge variant="warning">Pre-release</Badge>
                     {/if}
                     {#if release === releases.find(r => !r.is_draft && !r.is_prerelease)}
                       <Badge variant="success">Latest</Badge>
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
                   <Button variant="ghost" size="sm" icon="edit" onclick={() => handleEdit(release)} />
                   <Button variant="ghost" size="sm" icon="delete" class="danger-icon" onclick={() => handleDelete(release.id)} />
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
                <Card padding="sm" class="assets-box">
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
                </Card>
                
                {#if release.publisher}
                  <div class="publisher-info">
                    <span>Published by</span>
                    <a href={`/${release.publisher.username}`} class="user-chip">
                      <Avatar src={release.publisher.avatar_url} username={release.publisher.username} size={28} />
                      {release.publisher.username}
                    </a>
                  </div>
                {/if}
              </div>
            </div>
          </Card>
        </div>
      {/each}
    </div>

    {#if hasNext || hasPrevious}
      <div class="pagination">
        <Button variant="secondary" onclick={() => fetchReleases("")} disabled={loading || !hasPrevious} icon="first_page">
          Newest
        </Button>
        <Button onclick={() => fetchReleases(nextAfter)} disabled={loading || !hasNext} icon="navigate_next">
          Older releases
        </Button>
      </div>
    {/if}
  {/if}
</div>

<style>
  .repo-releases { display: flex; flex-direction: column; gap: 3rem; min-height: 400px; padding: 1rem; }
  
  .editor-overlay { position: fixed; inset: 0; background: rgba(0,0,0,0.7); backdrop-filter: blur(12px); z-index: 1000; display: flex; align-items: flex-start; justify-content: center; padding: 4rem; overflow-y: auto; }

  .section-actions-minimal { display: flex; justify-content: flex-end; }

  .empty-releases { text-align: center; padding: 6rem 2rem; display: flex; flex-direction: column; align-items: center; gap: 1rem; color: var(--text-muted); }
  .empty-releases .icon { font-size: 4rem; color: var(--md-sys-color-surface-variant); margin-bottom: 1rem; }
  .empty-releases h3 { margin: 0; font-size: 1.5rem; color: var(--md-sys-color-on-surface); }

  .releases-timeline { display: flex; flex-direction: column; position: relative; padding-left: 2.5rem; gap: 1rem; }
  .releases-timeline::before { content: ''; position: absolute; left: 0; top: 1rem; bottom: 1rem; width: 4px; background: var(--md-sys-color-surface-variant); border-radius: 2px; }
  
  .release-item { position: relative; padding-bottom: 2rem; transition: transform 0.2s; }
  .release-item.is-draft { opacity: 0.8; }
  .release-marker { position: absolute; left: -2.85rem; top: 2rem; width: 14px; height: 14px; border-radius: 50%; background: var(--md-sys-color-outline-variant); border: 4px solid var(--md-sys-color-background); z-index: 1; transition: all 0.3s; }
  .release-item:hover .release-marker { background: var(--md-sys-color-primary); transform: scale(1.2); }

  .release-card-header { padding: 2rem; display: flex; justify-content: space-between; align-items: flex-start; background: var(--md-sys-color-background); border-bottom: 1px solid var(--md-sys-color-outline-variant); }
  
  .title-row-flat { display: flex; align-items: center; gap: 1.5rem; flex-wrap: wrap; }
  .title-row-flat h3 { margin: 0; font-size: 1.8rem; font-weight: 800; letter-spacing: -0.02em; color: var(--md-sys-color-on-surface); }
  
  .badge-row { display: flex; gap: 0.5rem; }

  .meta-row { display: flex; gap: 2rem; flex-wrap: wrap; margin-top: 1rem; }
  .meta-row span { display: flex; align-items: center; gap: 0.5rem; font-size: 0.9rem; color: var(--text-muted); font-weight: 600; }
  .meta-row .material-symbols-rounded { font-size: 1.25rem; opacity: 0.7; }
  .meta-row .tag.clickable { color: var(--md-sys-color-primary); cursor: pointer; }
  .meta-row .tag.clickable:hover { text-decoration: underline; }

  .header-actions { display: flex; gap: 0.25rem; }
  :global(.danger-icon:hover) { color: var(--color-danger) !important; }

  .release-body { padding: 2rem; display: grid; grid-template-columns: 1fr 300px; gap: 3rem; }
  
  .release-content { font-size: 1.05rem; line-height: 1.8; color: var(--md-sys-color-on-surface); white-space: pre-wrap; }
  .empty-note { opacity: 0.5; }

  .release-side { display: flex; flex-direction: column; gap: 2rem; }
  
  :global(.assets-box) { background: var(--md-sys-color-background) !important; }
  :global(.assets-box) h4 { margin: 0 0 1rem 0; font-size: 0.85rem; text-transform: uppercase; color: var(--text-muted); letter-spacing: 0.1em; font-weight: 800; }
  .asset-links { display: flex; flex-direction: column; gap: 0.25rem; }
  .asset-item { display: flex; align-items: center; gap: 0.75rem; text-decoration: none; color: var(--md-sys-color-primary); font-weight: 700; font-size: 0.95rem; padding: 0.75rem; border-radius: 0.75rem; transition: all 0.2s; }
  .asset-item:hover { background: var(--md-sys-color-primary-container); transform: translateX(4px); }
  .asset-item.source { color: var(--text-muted); font-weight: 600; }
  .asset-item span { font-size: 1.4rem; opacity: 0.8; }

  .publisher-info { display: flex; flex-direction: column; gap: 0.75rem; padding-left: 0.5rem; }
  .publisher-info span { font-size: 0.8rem; font-weight: 700; text-transform: uppercase; color: var(--text-muted); letter-spacing: 0.05em; }
  .user-chip { display: flex; align-items: center; gap: 0.75rem; font-weight: 700; color: var(--md-sys-color-on-surface); text-decoration: none; }
  .user-chip :global(.avatar-container) { border: 1px solid var(--md-sys-color-outline-variant); }

  .pagination { display: flex; justify-content: center; gap: 1rem; margin-top: 2rem; }
  
  .error-msg { text-align: center; padding: 4rem; color: var(--color-danger); }

  .skeleton-item { height: 250px; border-radius: 1.5rem; background: var(--md-sys-color-surface-variant); animation: pulse 1.5s infinite; margin-bottom: 2rem; }
  @keyframes pulse { 0%, 100% { opacity: 0.4; } 50% { opacity: 0.7; } }

  @media (max-width: 1100px) {
    .release-body { grid-template-columns: 1fr; }
    .editor-overlay { padding: 1rem; }
  }
</style>

