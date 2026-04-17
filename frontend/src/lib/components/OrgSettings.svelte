<script lang="ts">
  import { apiFetch } from "$lib/api";
  import { onMount } from 'svelte';
  
  let { orgName } = $props();
  
  let org = $state<any>({});
  let loading = $state(true);
  let saving = $state(false);
  let message = $state('');

  async function fetchOrg() {
    loading = true;
    try {
      const res = await apiFetch(`/api/v2/orgs/${orgName}`);
      if (res.ok) org = await res.json();
    } catch (e) {
      console.error(e);
    } finally {
      loading = false;
    }
  }

  async function saveOrg() {
    saving = true;
    message = '';
    try {
      const res = await apiFetch(`/api/v2/orgs/${orgName}`, {
        method: 'PATCH',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(org)
      });
      if (res.ok) {
        message = 'Organization updated successfully!';
        setTimeout(() => message = '', 3000);
      }
    } catch (e) {
      message = 'Failed to update organization.';
    } finally {
      saving = false;
    }
  }

  onMount(fetchOrg);
</script>

<div class="org-settings island">
  <header class="card-header">
    <h3 class="card-title">Organization Settings</h3>
  </header>
  
  {#if loading}
    <div class="skeleton-form"></div>
  {:else}
    <form class="settings-form" onsubmit={(e) => { e.preventDefault(); saveOrg(); }}>
      <div class="form-group">
        <label for="full_name">Organization Name</label>
        <input type="text" id="full_name" bind:value={org.full_name} placeholder="Public name" />
      </div>

      <div class="form-group">
        <label for="desc">Description</label>
        <textarea id="desc" bind:value={org.description} placeholder="A short bio of the organization"></textarea>
      </div>

      <div class="form-group">
        <label for="website">Website</label>
        <input type="url" id="website" bind:value={org.website} placeholder="https://org.example.com" />
      </div>

      <div class="form-group">
        <label for="location">Location</label>
        <input type="text" id="location" bind:value={org.location} placeholder="City, Country" />
      </div>

      <div class="form-actions">
        {#if message}
           <span class="status-msg success">{message}</span>
        {/if}
        <button type="submit" class="btn btn-primary" disabled={saving}>
          {saving ? 'Saving...' : 'Update Settings'}
        </button>
      </div>
    </form>
  {/if}
</div>

<div class="org-settings island danger">
  <header class="card-header">
    <h3 class="card-title red">Danger Zone</h3>
  </header>
  <div class="card-body danger-body">
    <div class="danger-item">
      <div class="danger-text">
        <strong>Delete this organization</strong>
        <p>This will permanently delete the organization and all its repositories and teams. This cannot be undone.</p>
      </div>
      <button class="btn btn-danger">Delete Organization</button>
    </div>
  </div>
</div>

<style>
  .org-settings { display: flex; flex-direction: column; overflow: hidden; margin-bottom: 2rem; }
  .card-header { padding: 1.5rem 2rem; border-bottom: 1px solid var(--border-color); }
  .card-title { margin: 0; font-size: 1.25rem; font-weight: 700; }
  .card-title.red { color: #ef4444; }

  .settings-form { padding: 2rem; display: flex; flex-direction: column; gap: 1.5rem; }

  .form-group { display: flex; flex-direction: column; gap: 0.5rem; }
  .form-group label { font-size: 0.85rem; font-weight: 800; text-transform: uppercase; color: var(--text-muted); padding-left: 0.25rem; }
  
  input, textarea { background: rgba(0,0,0,0.2); border: 1px solid var(--border-color); padding: 0.75rem 1rem; border-radius: 0.75rem; color: var(--text-color); font-family: inherit; font-size: 1rem; transition: all 0.2s; }
  textarea { min-height: 80px; }
  input:focus, textarea:focus { outline: none; border-color: var(--primary-color); background: rgba(0,0,0,0.4); }

  .form-actions { display: flex; justify-content: flex-end; align-items: center; gap: 1.5rem; margin-top: 2rem; padding-top: 1.5rem; border-top: 1px solid var(--border-color); }
  .status-msg { font-size: 0.9rem; color: #10b981; }

  .danger { border: 1px solid rgba(239, 68, 68, 0.2); }
  .danger-body { padding: 2rem; }
  .danger-item { display: flex; justify-content: space-between; align-items: center; }
  .danger-text strong { color: #ef4444; font-size: 1.1rem; display: block; margin-bottom: 0.25rem; }
  .danger-text p { color: var(--text-muted); margin: 0; font-size: 0.9rem; }

  .skeleton-form { height: 300px; animation: pulse 1.5s infinite; }
  @keyframes pulse { 0%, 100% { opacity: 0.3; } 50% { opacity: 0.6; } }
</style>
