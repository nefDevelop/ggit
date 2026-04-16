<script lang="ts">
  import { onMount } from 'svelte';
  import { page } from '$app/state';
  import { goto } from '$app/navigation';

  let { username, reponame, initialData = null } = $props();

  let branches = $state<any[]>([]);
  let tag_name = $state(initialData?.TagName || '');
  let tag_target = $state(initialData?.Target || '');
  let title = $state(initialData?.Title || '');
  let content = $state(initialData?.Note || '');
  let is_prerelease = $state(initialData?.IsPrerelease || false);
  let is_draft = $state(initialData?.IsDraft || false);
  
  let loadingBranches = $state(true);
  let submitting = $state(false);

  async function fetchBranches() {
    try {
      const res = await fetch(`/api/v2/${username}/${reponame}/branches`);
      if (res.ok) {
        const data = await res.json();
        branches = [data.default, ...data.active, ...data.stale];
        if (!tag_target) tag_target = data.default.name;
      }
    } finally {
      loadingBranches = false;
    }
  }

  async function saveRelease() {
    if (!tag_name || !tag_target) return;
    submitting = true;
    
    const url = initialData 
      ? `/api/v2/${username}/${reponame}/releases/${initialData.ID}`
      : `/api/v2/${username}/${reponame}/releases`;
    
    const method = initialData ? 'PATCH' : 'POST';

    try {
      const res = await fetch(url, {
        method,
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          tag_name,
          target: tag_target,
          title,
          content,
          prerelease: is_prerelease,
          draft: is_draft ? "1" : "" // API expectations vary, but usually this is how it works
        })
      });

      if (res.ok) {
        goto(`/${username}/${reponame}/releases`);
      }
    } catch (e) {
      console.error(e);
    } finally {
      submitting = false;
    }
  }

  onMount(fetchBranches);
</script>

<div class="release-editor island tonal">
  <header class="editor-header">
     <h2>{initialData ? 'Edit Release' : 'Draft a new release'}</h2>
     <p>Tag your software versions for your users to download.</p>
  </header>

  <form onsubmit={(e) => { e.preventDefault(); saveRelease(); }}>
    <div class="form-row">
      <div class="field">
        <label for="tag">Tag version</label>
        <input id="tag" type="text" bind:value={tag_name} placeholder="e.g. v1.0.0" required disabled={!!initialData} />
      </div>
      <div class="field">
        <label for="target">Target branch</label>
        <select id="target" bind:value={tag_target} disabled={!!initialData}>
          {#each branches as b}
            <option value={b.name}>{b.name}</option>
          {/each}
        </select>
      </div>
    </div>

    <div class="field">
      <label for="title">Release title</label>
      <input id="title" type="text" bind:value={title} placeholder="Release title" />
    </div>

    <div class="field">
      <label for="content">Describe this release</label>
      <textarea id="content" bind:value={content} placeholder="Markdown supported..."></textarea>
    </div>

    <div class="options-row island">
       <label class="checkbox-label">
         <input type="checkbox" bind:checked={is_prerelease} />
         <div class="check-text">
            <strong>This is a pre-release</strong>
            <span>We’ll point out that this release is not yet ready for production.</span>
         </div>
       </label>
    </div>

    <div class="footer-actions">
       <button type="button" class="btn btn-outline" onclick={() => goto(`/${username}/${reponame}/releases`)}>Cancel</button>
       <div class="submit-group">
          <label class="draft-toggle">
            <input type="checkbox" bind:checked={is_draft} />
            Keep as draft
          </label>
          <button type="submit" class="btn btn-primary lg" disabled={submitting}>
            {submitting ? 'Saving...' : initialData ? 'Update Release' : 'Publish Release'}
          </button>
       </div>
    </div>
  </form>
</div>

<style>
  .release-editor { display: flex; flex-direction: column; gap: 2.5rem; padding: 3rem; max-width: 900px; margin: 0 auto; box-shadow: var(--shadow-xl); }
  
  .editor-header h2 { margin: 0; font-size: 1.75rem; font-weight: 800; }
  .editor-header p { color: var(--text-muted); margin: 0.5rem 0 0 0; }

  form { display: flex; flex-direction: column; gap: 1.5rem; }
  
  .form-row { display: grid; grid-template-columns: 1fr 1fr; gap: 1.5rem; }
  
  .field { display: flex; flex-direction: column; gap: 0.5rem; }
  .field label { font-size: 0.85rem; font-weight: 700; color: var(--text-muted); text-transform: uppercase; letter-spacing: 0.05em; }
  .field input, .field select, .field textarea { background: rgba(0,0,0,0.2); border: 1px solid var(--border-color); color: #fff; padding: 0.75rem 1rem; border-radius: 0.75rem; font-size: 1rem; font-family: inherit; }
  .field textarea { height: 300px; resize: vertical; }

  .options-row { background: rgba(255,152,0,0.05); border: 1px solid rgba(255,152,0,0.2); padding: 1.5rem; border-radius: 1rem; }
  
  .checkbox-label { display: flex; align-items: flex-start; gap: 1rem; cursor: pointer; }
  .checkbox-label input { width: 1.25rem; height: 1.25rem; margin-top: 0.25rem; }
  .check-text { display: flex; flex-direction: column; gap: 0.1rem; }
  .check-text strong { font-size: 0.95rem; }
  .check-text span { font-size: 0.85rem; color: var(--text-muted); }

  .footer-actions { border-top: 1px solid var(--border-color); padding-top: 2rem; display: flex; justify-content: space-between; align-items: center; }
  .submit-group { display: flex; align-items: center; gap: 2rem; }
  
  .draft-toggle { display: flex; align-items: center; gap: 0.5rem; font-size: 0.9rem; font-weight: 600; cursor: pointer; color: var(--text-muted); }
  
  @media (max-width: 600px) {
    .form-row { grid-template-columns: 1fr; }
    .footer-actions { flex-direction: column; gap: 1.5rem; align-items: stretch; }
    .submit-group { flex-direction: column; align-items: center; }
  }
</style>
