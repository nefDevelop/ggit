<script lang="ts">
  import { apiFetch } from '$lib/api';
  import { page } from '$app/state';
  import { fade } from 'svelte/transition';
  import PageHeader from "$lib/components/ui/PageHeader.svelte";
  import Button from "$lib/components/ui/Button.svelte";
  import { Plus, Tag, FolderDown, Archive, Trash2, Rocket, AlertCircle } from "lucide-svelte";

  let username = $derived(page.params.username);
  let reponame = $derived(page.params.reponame);

  interface RepoTag {
    name: string;
    sha: string;
    message: string;
    tagger: string;
    created_at: string;
    zip_url: string;
    tar_url: string;
  }

  let tags = $state<RepoTag[]>([]);
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
  <PageHeader title="" icon="" compact>
    {#snippet actions()}
      <Button variant="primary" onclick={() => { showNewTag = !showNewTag; createError = ''; }} icon={Plus}>
        New tag
      </Button>
    {/snippet}
  </PageHeader>

  <div class="tags-content">
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
          <Button variant="ghost" size="sm" onclick={() => showNewTag = false}>Cancel</Button>
          <Button variant="primary" size="sm" onclick={createTag} disabled={creating || !newTagName || !newTagTarget}>
            {#if creating}
              <span class="material-symbols-rounded spin">sync</span>
            {:else}
              <Tag size={16} />
            {/if}
            Create tag
          </Button>
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
      <div class="empty-state">
        <AlertCircle size={48} />
        <h3>No tags found</h3>
        <p>There are no tags in this repository. Tags mark specific points in the history as important — use them for releases.</p>
        <Button variant="primary" size="sm" onclick={() => showNewTag = true} icon={Plus}>
          Create first tag
        </Button>
      </div>
    {:else}
      <div class="flat-list-container">
        {#each tags as tag (tag.name)}
          <div class="list-item-row" transition:fade={{ duration: 120 }}>
            <div class="list-item-icon">
              <Tag size={18} style="color: var(--text-muted); opacity: 0.5;" />
            </div>
            
            <div class="item-main-content">
              <div class="item-title-bar">
                <a href={`/${username}/${reponame}/releases`} class="item-title-link">
                  {tag.name}
                </a>
              </div>
              <div class="item-meta-data">
                {#if tag.sha}
                  <code class="tag-sha">{tag.sha.slice(0, 8)}</code>
                {/if}
              </div>
            </div>

            <div class="item-actions-right">
              <a href={tag.zip_url} class="btn-icon-minimal" title="Download ZIP">
                <FolderDown size={16} />
              </a>
              <a href={tag.tar_url} class="btn-icon-minimal" title="Download TAR.GZ">
                <Archive size={16} />
              </a>
              <a href={`/${username}/${reponame}/releases`} class="btn-icon-minimal" title="Release">
                <Rocket size={16} />
              </a>
              <button
                class="btn-icon-minimal delete"
                class:confirm={deleteConfirm === tag.name}
                onclick={() => deleteTag(tag.name)}
                disabled={deleting === tag.name}
                title={deleteConfirm === tag.name ? 'Click again to confirm delete' : 'Delete tag'}
              >
                {#if deleting === tag.name}
                  <span class="material-symbols-rounded spin">sync</span>
                {:else if deleteConfirm === tag.name}
                  <span style="color: #fb923c; font-weight: bold; font-size: 0.8rem;">Sure?</span>
                {:else}
                  <Trash2 size={16} />
                {/if}
              </button>
            </div>
          </div>
        {/each}
      </div>
    {/if}
  </div>
</div>

<style>
  .tags-page { display: flex; flex-direction: column; gap: 0; }
  .tags-content { padding: 0; display: flex; flex-direction: column; gap: 1.25rem; }

  /* New tag form */
  .new-tag-form { padding: 1.5rem; border-radius: 1rem; border: 1px solid var(--primary-color); display: flex; flex-direction: column; gap: 1rem; margin-bottom: 1.25rem; }
  .new-tag-form h3 { margin: 0; font-size: 1rem; font-weight: 700; }
  .form-row { display: grid; grid-template-columns: 1fr 1fr; gap: 1rem; }
  .form-group { display: flex; flex-direction: column; gap: 0.4rem; }
  .form-group label { font-size: 0.8rem; font-weight: 700; color: var(--text-muted); text-transform: uppercase; letter-spacing: 0.05em; }
  .form-control { background: rgba(0,0,0,0.15); border: 1px solid var(--border-color); border-radius: 10px; padding: 0.6rem 0.9rem; color: var(--text-color); font-size: 0.9rem; outline: none; transition: border-color 0.2s; }
  .form-control:focus { border-color: var(--primary-color); }
  .form-error { color: #f87171; font-size: 0.85rem; margin: 0; }
  .form-actions { display: flex; gap: 0.75rem; justify-content: flex-end; }

  .tag-sha { font-size: 0.75rem; font-family: 'JetBrains Mono', monospace; background: rgba(0,0,0,0.15); padding: 0.15rem 0.5rem; border-radius: 6px; color: var(--text-muted); }

  .delete.confirm { color: #fb923c; background: rgba(251, 146, 60, 0.1); }

  .empty-state { display: flex; flex-direction: column; align-items: center; justify-content: center; gap: 1rem; padding: 4rem 2rem; border-radius: 1rem; text-align: center; color: var(--text-muted); opacity: 0.6; }
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
  }
</style>
