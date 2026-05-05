<script lang="ts">
  import { apiFetch } from '$lib/api';
  import { page } from '$app/state';
  import { fade } from 'svelte/transition';

  let username = $derived(page.params.username);
  let reponame = $derived(page.params.reponame);

  interface Tag {
    name: string;
    sha: string;
    message: string;
    tagger: string;
    created_at: string;
    zip_url: string;
    tar_url: string;
  }

  let tags = $state<Tag[]>([]);
  let loading = $state(true);
  let error = $state('');
  let deleteConfirm = $state('');
  let deleting = $state('');

  // New tag form
  let showNewTag = $state(false);
  let newTagName = $state('');
  let newTagTarget = $state('');
  let creating = $state(false);
  let createError = $state('');

  async function fetchTags() {
    loading = true;
    error = '';
    try {
      const res = await apiFetch(`/api/v2/${username}/${reponame}/tags`);
      if (res.ok) {
        tags = await res.json();
      } else {
        error = 'Failed to load tags.';
      }
    } catch {
      error = 'Network error loading tags.';
    } finally {
      loading = false;
    }
  }

  async function deleteTag(tagName: string) {
    if (deleteConfirm !== tagName) {
      deleteConfirm = tagName;
      return;
    }
    deleting = tagName;
    try {
      const res = await apiFetch(`/api/v2/${username}/${reponame}/tags/${encodeURIComponent(tagName)}`, {
        method: 'DELETE'
      });
      if (res.ok) {
        tags = tags.filter(t => t.name !== tagName);
        deleteConfirm = '';
      } else {
        error = 'Failed to delete tag.';
      }
    } catch {
      error = 'Network error deleting tag.';
    } finally {
      deleting = '';
    }
  }

  async function createTag() {
    if (!newTagName.trim() || !newTagTarget.trim()) return;
    creating = true;
    createError = '';
    try {
      const res = await apiFetch(`/api/v2/${username}/${reponame}/tags`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ name: newTagName.trim(), target: newTagTarget.trim() })
      });
      if (res.ok) {
        const tag = await res.json();
        tags = [tag, ...tags];
        newTagName = '';
        newTagTarget = '';
        showNewTag = false;
      } else {
        const d = await res.json().catch(() => ({}));
        createError = d.message || 'Failed to create tag.';
      }
    } catch {
      createError = 'Network error.';
    } finally {
      creating = false;
    }
  }

  $effect(() => {
    if (username && reponame) fetchTags();
  });
</script>

<svelte:head>
  <title>Tags · {reponame} · {username}</title>
</svelte:head>

<div class="tags-page">
  <div class="header-actions-right">
    <button class="btn btn-primary btn-sm" onclick={() => { showNewTag = !showNewTag; createError = ''; }}>
      <span class="material-symbols-rounded">add_circle</span>
      New tag
    </button>
  </div>

  {#if showNewTag}
    <div class="new-tag-form tonal" transition:fade={{ duration: 150 }}>
      <h3>Create a new tag</h3>
      <div class="form-row">
        <div class="form-group">
          <label for="tag-name">Tag name</label>
          <input id="tag-name" class="form-control" placeholder="v1.0.0" bind:value={newTagName} />
        </div>
        <div class="form-group">
          <label for="tag-target">Target branch or commit SHA</label>
          <input id="tag-target" class="form-control" placeholder="main or abc1234" bind:value={newTagTarget} />
        </div>
      </div>
      {#if createError}
        <p class="form-error">{createError}</p>
      {/if}
      <div class="form-actions">
        <button class="btn btn-ghost btn-sm" onclick={() => showNewTag = false}>Cancel</button>
        <button class="btn btn-primary btn-sm" onclick={createTag} disabled={creating || !newTagName || !newTagTarget}>
          {#if creating}
            <span class="material-symbols-rounded spin">sync</span>
          {:else}
            <span class="material-symbols-rounded">sell</span>
          {/if}
          Create tag
        </button>
      </div>
    </div>
  {/if}

  {#if error}
    <div class="alert-error">{error}</div>
  {/if}

  {#if loading}
    <div class="skeleton-list">
      {#each [1,2,3] as _}
        <div class="skeleton-row tonal"></div>
      {/each}
    </div>
  {:else if tags.length === 0}
    <div class="empty-state tonal">
      <span class="material-symbols-rounded empty-icon">sell</span>
      <h3>No tags yet</h3>
      <p>Tags mark specific points in the history as important — use them for releases.</p>
      <button class="btn btn-primary btn-sm" onclick={() => showNewTag = true}>
        <span class="material-symbols-rounded">add</span>
        Create first tag
      </button>
    </div>
  {:else}
    <div class="tag-list">
      {#each tags as tag (tag.name)}
        <div class="tag-row tonal" transition:fade={{ duration: 120 }}>
          <div class="tag-info">
            <a href={`/${username}/${reponame}/releases`} class="tag-name">
              {tag.name}
            </a>
            {#if tag.sha}
              <code class="tag-sha">{tag.sha.slice(0, 8)}</code>
            {/if}
          </div>

          <div class="tag-actions">
            <a href={tag.zip_url} class="btn btn-ghost btn-xs" title="Download ZIP">
              <span class="material-symbols-rounded">folder_zip</span>
              zip
            </a>
            <a href={tag.tar_url} class="btn btn-ghost btn-xs" title="Download TAR.GZ">
              <span class="material-symbols-rounded">inventory_2</span>
              tar.gz
            </a>
            <a href={`/${username}/${reponame}/releases`} class="btn btn-ghost btn-xs">
              <span class="material-symbols-rounded">rocket_launch</span>
              Release
            </a>
            <button
              class="btn btn-ghost btn-xs danger"
              class:confirm={deleteConfirm === tag.name}
              onclick={() => deleteTag(tag.name)}
              disabled={deleting === tag.name}
              title={deleteConfirm === tag.name ? 'Click again to confirm delete' : 'Delete tag'}
            >
              {#if deleting === tag.name}
                <span class="material-symbols-rounded spin">sync</span>
              {:else if deleteConfirm === tag.name}
                <span class="material-symbols-rounded">warning</span>
                Confirm?
              {:else}
                <span class="material-symbols-rounded">delete</span>
              {/if}
            </button>
          </div>
        </div>
      {/each}
    </div>
  {/if}
</div>

<style>
  .tags-page { max-width: 900px; margin: 0 auto; display: flex; flex-direction: column; gap: 1.25rem; }

  .header-actions-right { display: flex; justify-content: flex-end; margin-bottom: 0.5rem; }

  .page-header { display: flex; align-items: center; justify-content: space-between; }
  .header-left { display: flex; align-items: center; gap: 0.75rem; }
  .header-icon { font-size: 1.5rem; color: var(--primary-color); }
  h2 { margin: 0; font-size: 1.25rem; font-weight: 800; display: flex; align-items: center; gap: 0.5rem; }
  .count-chip { font-size: 0.8rem; background: var(--primary-container); color: var(--on-primary-container); padding: 0.1rem 0.5rem; border-radius: 999px; font-weight: 700; }

  /* New tag form */
  .new-tag-form { padding: 1.5rem; border-radius: 1rem; border: 1px solid var(--primary-color); display: flex; flex-direction: column; gap: 1rem; }
  .new-tag-form h3 { margin: 0; font-size: 1rem; font-weight: 700; }
  .form-row { display: grid; grid-template-columns: 1fr 1fr; gap: 1rem; }
  .form-group { display: flex; flex-direction: column; gap: 0.4rem; }
  .form-group label { font-size: 0.8rem; font-weight: 700; color: var(--text-muted); text-transform: uppercase; letter-spacing: 0.05em; }
  .form-control { background: rgba(0,0,0,0.15); border: 1px solid var(--border-color); border-radius: 10px; padding: 0.6rem 0.9rem; color: var(--text-color); font-size: 0.9rem; outline: none; transition: border-color 0.2s; }
  .form-control:focus { border-color: var(--primary-color); }
  .form-error { color: #f87171; font-size: 0.85rem; margin: 0; }
  .form-actions { display: flex; gap: 0.75rem; justify-content: flex-end; }

  /* Tag list */
  .tag-list { display: flex; flex-direction: column; gap: 0.5rem; }
  .tag-row { display: flex; align-items: center; justify-content: space-between; padding: 1rem 1.25rem; border-radius: 0.875rem; border: 1px solid var(--border-color); transition: border-color 0.2s; gap: 1rem; }
  .tag-row:hover { border-color: var(--primary-color); }
  .tag-info { display: flex; align-items: center; gap: 0.75rem; min-width: 0; }
  .tag-name { display: flex; align-items: center; gap: 0.4rem; font-weight: 700; font-size: 0.95rem; text-decoration: none; color: var(--primary-color); }
  .tag-name span { font-size: 1.1rem; }
  .tag-sha { font-size: 0.75rem; font-family: 'JetBrains Mono', monospace; background: rgba(0,0,0,0.15); padding: 0.15rem 0.5rem; border-radius: 6px; color: var(--text-muted); }
  .tag-actions { display: flex; align-items: center; gap: 0.5rem; flex-shrink: 0; }

  .btn-xs { padding: 0.3rem 0.7rem; font-size: 0.78rem; border-radius: 8px; display: flex; align-items: center; gap: 0.3rem; }
  .btn-xs span { font-size: 1rem; }
  .danger { color: #f87171; }
  .danger:hover { background: rgba(248, 113, 113, 0.1); }
  .danger.confirm { color: #fb923c; border-color: #fb923c; }

  /* Empty */
  .empty-state { display: flex; flex-direction: column; align-items: center; justify-content: center; gap: 1rem; padding: 4rem 2rem; border-radius: 1rem; border: 1px dashed var(--border-color); text-align: center; }
  .empty-icon { font-size: 3rem; color: var(--primary-color); opacity: 0.4; }
  .empty-state h3 { margin: 0; font-size: 1.1rem; }
  .empty-state p { margin: 0; color: var(--text-muted); font-size: 0.9rem; max-width: 400px; }

  .skeleton-list { display: flex; flex-direction: column; gap: 0.5rem; }
  .skeleton-row { height: 60px; border-radius: 0.875rem; animation: pulse 1.5s infinite; }
  @keyframes pulse { 0%, 100% { opacity: 0.5; } 50% { opacity: 0.25; } }

  .alert-error { background: rgba(248,113,113,0.1); border: 1px solid #f87171; border-radius: 0.75rem; padding: 0.75rem 1rem; color: #f87171; font-size: 0.9rem; }

  .spin { animation: spin 1s linear infinite; }
  @keyframes spin { to { transform: rotate(360deg); } }

  @media (max-width: 600px) {
    .form-row { grid-template-columns: 1fr; }
    .tag-row { flex-direction: column; align-items: flex-start; }
    .tag-actions { width: 100%; flex-wrap: wrap; }
  }
</style>
