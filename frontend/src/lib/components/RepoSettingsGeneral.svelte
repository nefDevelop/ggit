<script lang="ts">
  import { apiFetch } from "$lib/api";
  import { onMount, getContext } from 'svelte';
  import { goto } from '$app/navigation';
  
  const repo: any = getContext('repo');
  let settings = $state<any>({});
  let loading = $state(true);
  let saving = $state(false);
  let message = $state('');
  let deleteConfirm = $state('');
  let syncing = $state(false);

  function formatDate(d: string) {
    if (!d || d.startsWith('0001')) return 'Never';
    return new Date(d).toLocaleString();
  }

  async function fetchSettings() {
    loading = true;
    try {
      const res = await apiFetch(`/api/v2/${repo.data.owner.username}/${repo.data.name}/settings`);
      if (res.ok) settings = await res.json();
    } catch (e) {
      console.error(e);
    } finally {
      loading = false;
    }
  }

  async function saveSettings() {
    saving = true;
    message = '';
    try {
      const res = await apiFetch(`/api/v2/${repo.data.owner.username}/${repo.data.name}/settings`, {
        method: 'PATCH',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(settings)
      });
      if (res.ok) {
        message = 'Repository updated successfully!';
        setTimeout(() => message = '', 3000);
        fetchSettings(); // Refresh to get computed fields if any
      }
    } catch (e) {
      message = 'Failed to update repository.';
    } finally {
      saving = false;
    }
  }

  async function syncMirror() {
    syncing = true;
    try {
      const res = await apiFetch(`/api/v2/${repo.data.owner.username}/${repo.data.name}/settings/mirror/sync`, {
        method: 'POST'
      });
      if (res.ok) {
        message = 'Mirror sync triggered!';
        setTimeout(() => message = '', 3000);
        fetchSettings();
      }
    } catch (e) {
      console.error(e);
      message = 'Failed to trigger sync.';
    } finally {
      syncing = false;
    }
  }

  async function deleteRepo() {
    if (deleteConfirm !== settings.name) return;
    try {
      const res = await apiFetch(`/api/v2/${repo.data.owner.username}/${repo.data.name}/settings`, {
        method: 'DELETE',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ name: deleteConfirm })
      });
      if (res.ok) {
        goto('/dashboard');
      }
    } catch (e) {
      console.error(e);
    }
  }

  onMount(fetchSettings);
</script>

<div class="settings-card">
  
  {#if loading}
    <div class="skeleton-form"></div>
  {:else}
    <form class="settings-form" onsubmit={(e) => { e.preventDefault(); saveSettings(); }}>
      <div class="form-group">
        <label for="repo_name">Repository Name</label>
        <input type="text" id="repo_name" bind:value={settings.name} required />
      </div>

      <div class="form-group">
        <label for="repo_desc">Description</label>
        <textarea id="repo_desc" bind:value={settings.description} placeholder="Short description..."></textarea>
      </div>

      <div class="form-group">
        <label for="repo_site">Website</label>
        <input type="url" id="repo_site" bind:value={settings.website} placeholder="https://example.com" />
      </div>

      <div class="options-grid">
         <label class="checkbox-container">
            <input type="checkbox" bind:checked={settings.private} />
            <span class="checkmark"></span>
            <div class="checkbox-text">
              <strong>Make Private</strong>
              <p>Only visible to collaborators.</p>
            </div>
         </label>
         <label class="checkbox-container">
            <input type="checkbox" bind:checked={settings.enable_issues} />
            <span class="checkmark"></span>
            <div class="checkbox-text">
              <strong>Issues</strong>
              <p>Lightweight issue tracking.</p>
            </div>
         </label>
         <label class="checkbox-container">
            <input type="checkbox" bind:checked={settings.enable_pulls} />
            <span class="checkmark"></span>
            <div class="checkbox-text">
              <strong>Pull Requests</strong>
              <p>Code reviews and merges.</p>
            </div>
         </label>
         <label class="checkbox-container">
            <input type="checkbox" bind:checked={settings.enable_wiki} />
            <span class="checkmark"></span>
            <div class="checkbox-text">
              <strong>Wiki</strong>
              <p>Project documentation.</p>
            </div>
         </label>
      </div>

      <div class="form-actions">
        {#if message}
           <span class="status-msg success">{message}</span>
        {/if}
        <button type="submit" class="btn btn-primary" disabled={saving}>
          {saving ? 'Saving...' : 'Save Settings'}
        </button>
      </div>
    </form>
  {/if}
</div>

{#if settings.is_mirror}
<div class="settings-card island island-interactive mirror-card">
  <header class="card-header">
    <h3 class="card-title">Mirror Settings</h3>
  </header>
  <div class="card-body mirror-body">
    <div class="form-group">
      <label for="mirror_address">Remote Address</label>
      <input type="text" id="mirror_address" bind:value={settings.mirror_address} placeholder="https://github.com/user/repo.git" />
    </div>

    <div class="form-group">
      <label for="mirror_interval">Sync Interval (Hours)</label>
      <input type="number" id="mirror_interval" bind:value={settings.mirror_interval} min="1" />
    </div>

    <div class="mirror-meta">
      <div class="meta-item">
        <span class="meta-label">Last Sync:</span>
        <span class="meta-value">{formatDate(settings.last_sync)}</span>
      </div>
      <div class="meta-item">
        <span class="meta-label">Next Sync:</span>
        <span class="meta-value">{formatDate(settings.next_sync)}</span>
      </div>
    </div>

    <div class="form-actions space-between">
       <button class="btn btn-secondary" onclick={syncMirror} disabled={syncing}>
         {#if syncing}
           Syncing...
         {:else}
           Sync Now
         {/if}
       </button>
       <button class="btn btn-primary" onclick={saveSettings} disabled={saving}>
         {saving ? 'Saving...' : 'Save Mirror Settings'}
       </button>
    </div>
  </div>
</div>
{/if}

<div class="settings-card island island-interactive danger">
  <header class="card-header">
    <h3 class="card-title red">Danger Zone</h3>
  </header>
  <div class="card-body danger-body">
    <div class="danger-item">
      <div class="danger-text">
        <strong>Delete this repository</strong>
        <p>Once you delete a repository, there is no going back. Please be certain.</p>
      </div>
      <div class="delete-action">
         <input type="text" bind:value={deleteConfirm} placeholder="Type repository name to confirm" />
         <button class="btn btn-danger" disabled={deleteConfirm !== settings.name} onclick={deleteRepo}>
           Delete Repository
         </button>
      </div>
    </div>
  </div>
</div>

<style>
  .settings-card { display: flex; flex-direction: column; overflow: hidden; }
  .card-header { padding: 1.5rem 2rem; border-bottom: 1px solid var(--md-sys-color-outline-variant); }
  .card-title { margin: 0; font-size: 1.25rem; font-weight: 700; }
  .card-title.red { color: #ef4444; }

  .settings-form { padding: 2rem; display: flex; flex-direction: column; gap: 1.5rem; }

  .form-group { display: flex; flex-direction: column; gap: 0.5rem; }
  .form-group label { font-size: 0.85rem; font-weight: 800; text-transform: uppercase; color: var(--text-muted); padding-left: 0.25rem; }
  
  input, textarea { background: rgba(0,0,0,0.2); border: 1px solid var(--md-sys-color-outline-variant); padding: 0.75rem 1rem; border-radius: 0.75rem; color: var(--text-color); font-family: inherit; font-size: 1rem; transition: all 0.2s; }
  textarea { min-height: 80px; }
  input:focus, textarea:focus { outline: none; border-color: var(--primary-color); background: rgba(0,0,0,0.4); }

  .options-grid { display: grid; grid-template-columns: 1fr 1fr; gap: 1rem; margin-top: 1rem; }

  .checkbox-container { display: flex; align-items: flex-start; gap: 1rem; cursor: pointer; position: relative; padding: 1rem; background: rgba(255,255,255,0.02); border-radius: 1rem; transition: background 0.2s; }
  .checkbox-container:hover { background: rgba(255,255,255,0.04); }
  .checkbox-container input { visibility: hidden; width: 0; height: 0; position: absolute; }
  
  .checkmark { width: 1.25rem; height: 1.25rem; background: rgba(255,255,255,0.05); border: 2px solid var(--md-sys-color-outline-variant); border-radius: 0.4rem; flex-shrink: 0; position: relative; transition: all 0.2s; margin-top: 0.2rem; }
  .checkbox-container input:checked + .checkmark { background: var(--primary-color); border-color: var(--primary-color); }
  .checkmark::after { content: 'check'; font-family: 'Material Symbols Rounded'; position: absolute; top: 50%; left: 50%; transform: translate(-50%, -50%); color: white; opacity: 0; font-size: 0.9rem; }
  .checkbox-container input:checked + .checkmark::after { opacity: 1; }

  .checkbox-text strong { display: block; font-size: 0.95rem; margin-bottom: 0.1rem; }
  .checkbox-text p { font-size: 0.8rem; color: var(--text-muted); margin: 0; line-height: 1.2; }

  .form-actions { display: flex; justify-content: flex-end; align-items: center; gap: 1.5rem; margin-top: 2rem; padding-top: 1.5rem; border-top: 1px solid var(--md-sys-color-outline-variant); }
  .status-msg { font-size: 0.9rem; color: #10b981; }

  .danger { margin-top: 2rem; border: 1px solid rgba(239, 68, 68, 0.2); }
  .danger-body { padding: 2rem; }
  .danger-item { display: flex; flex-direction: column; gap: 1.5rem; }
  .danger-text strong { color: #ef4444; font-size: 1.1rem; }
  .danger-text p { color: var(--text-muted); margin: 0.25rem 0 0 0; }
  
  .delete-action { display: flex; gap: 1rem; }
  .delete-action input { flex: 1; border-color: rgba(239, 68, 68, 0.2); }
  .delete-action input:focus { border-color: #ef4444; }

  .skeleton-form { height: 300px; animation: pulse 1.5s infinite; }
  @keyframes pulse { 0%, 100% { opacity: 0.3; } 50% { opacity: 0.6; } }

  .mirror-card { margin-top: 2rem; }
  .mirror-body { padding: 2rem; display: flex; flex-direction: column; gap: 1.5rem; }
  .mirror-meta { display: grid; grid-template-columns: 1fr 1fr; gap: 1rem; padding: 1rem; background: rgba(0,0,0,0.2); border-radius: 0.75rem; border: 1px solid var(--md-sys-color-outline-variant); }
  .meta-item { display: flex; flex-direction: column; gap: 0.25rem; }
  .meta-label { font-size: 0.75rem; font-weight: 800; text-transform: uppercase; color: var(--text-muted); }
  .meta-value { font-size: 0.95rem; font-weight: 600; }
  .form-actions.space-between { justify-content: space-between; }
</style>
