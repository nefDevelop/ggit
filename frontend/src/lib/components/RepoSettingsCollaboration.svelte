<script lang="ts">
  import { apiFetch } from "$lib/api";
  import { onMount, getContext } from 'svelte';
  
  const repo: any = getContext('repo');
  let collaborators = $state<any[]>([]);
  let loading = $state(true);
  let newCollaborator = $state('');
  let submitting = $state(false);
  let message = $state('');

  async function fetchCollaborators() {
    loading = true;
    try {
      const res = await apiFetch(`/api/v2/${repo.data.owner.username}/${repo.data.name}/settings/collaborators`);
      if (res.ok) collaborators = await res.json();
    } catch (e) {
      console.error(e);
    } finally {
      loading = false;
    }
  }

  async function addCollaborator() {
    if (!newCollaborator) return;
    submitting = true;
    message = '';
    try {
      const res = await apiFetch(`/api/v2/${repo.data.owner.username}/${repo.data.name}/settings/collaborators/${newCollaborator}`, {
        method: 'POST'
      });
      if (res.ok) {
        newCollaborator = '';
        fetchCollaborators();
      } else {
        const body = await res.json();
        message = body.message || 'User not found.';
      }
    } catch (e) {
      message = 'Failed to add collaborator.';
    } finally {
      submitting = false;
    }
  }

  async function removeCollaborator(username: string) {
    if (!confirm(`Remove ${username} as collaborator?`)) return;
    try {
      const res = await apiFetch(`/api/v2/${repo.data.owner.username}/${repo.data.name}/settings/collaborators/${username}`, {
        method: 'DELETE'
      });
      if (res.ok) {
        fetchCollaborators();
      }
    } catch (e) {
      console.error(e);
    }
  }

  onMount(fetchCollaborators);
</script>

<div class="settings-view-flat">
  <div class="collaborator-list-modern">
    {#each collaborators as col}
      <div class="collab-row-flat">
        <img src={col.avatar_url} alt={col.username} class="mini-avatar" />
        <div class="collab-info">
          <strong>{col.full_name || col.username}</strong>
          <span class="user-handle">@{col.username}</span>
        </div>
        <div class="permission-pill">Write</div>
        <button class="btn-danger-icon-minimal" onclick={() => removeCollaborator(col.username)}>
          <span class="material-symbols-rounded">person_remove</span>
        </button>
      </div>
    {/each}
    {#if !loading && collaborators.length === 0}
      <div class="empty-collabs">No collaborators added yet.</div>
    {/if}
  </div>

  <div class="add-section-flat">
    <form class="add-collab-form-modern" onsubmit={(e) => { e.preventDefault(); addCollaborator(); }}>
      <div class="input-action-group">
        <input type="text" bind:value={newCollaborator} placeholder="Add someone by username..." required />
        <button type="submit" class="btn-action-primary" disabled={submitting || !newCollaborator}>
          {submitting ? 'Adding...' : 'Add'}
        </button>
      </div>
      {#if message}
         <span class="status-error">{message}</span>
      {/if}
    </form>
  </div>
</div>

<style>
  .settings-view-flat { display: flex; flex-direction: column; gap: 2rem; padding: 1rem; }

  .collaborator-list-modern { display: flex; flex-direction: column; background: var(--bg-color); border: 1px solid var(--md-sys-color-outline-variant); border-radius: 1rem; overflow: hidden; }
  
  .collab-row-flat { display: flex; align-items: center; gap: 1.25rem; padding: 1rem 1.5rem; border-bottom: 1px solid var(--md-sys-color-outline-variant); transition: background 0.2s; }
  .collab-row-flat:last-child { border-bottom: none; }
  .collab-row-flat:hover { background: var(--surface-variant); }
  
  .mini-avatar { width: 36px; height: 36px; border-radius: 50%; background: var(--surface-2); }
  
  .collab-info { flex: 1; display: flex; flex-direction: column; gap: 0.1rem; }
  .collab-info strong { font-size: 0.95rem; font-weight: 700; color: var(--text-color); }
  .user-handle { font-size: 0.75rem; color: var(--text-muted); opacity: 0.7; }

  .permission-pill { background: rgba(var(--primary-rgb), 0.1); color: var(--primary-color); padding: 0.2rem 0.75rem; border-radius: 2rem; font-size: 0.7rem; font-weight: 800; text-transform: uppercase; letter-spacing: 0.05em; }

  .btn-danger-icon-minimal { background: transparent; border: none; color: var(--error-color, #ef4444); cursor: pointer; padding: 0.5rem; border-radius: 50%; transition: all 0.2s; opacity: 0.4; }
  .btn-danger-icon-minimal:hover { opacity: 1; background: rgba(239, 68, 68, 0.1); }

  .add-section-flat { margin-top: 1rem; padding: 1.5rem; background: var(--surface-variant); border-radius: 1rem; border: 1px solid var(--md-sys-color-outline-variant); }
  
  .input-action-group { display: flex; gap: 1rem; }
  .input-action-group input { flex: 1; background: var(--bg-color); border: 1px solid var(--md-sys-color-outline-variant); padding: 0.75rem 1.25rem; border-radius: 0.75rem; color: var(--text-color); font-family: inherit; font-size: 0.95rem; }
  .input-action-group input:focus { outline: none; border-color: var(--primary-color); }

  .btn-action-primary { background: var(--primary-color); color: var(--on-primary); border: none; padding: 0 1.5rem; border-radius: 0.75rem; font-weight: 700; cursor: pointer; transition: all 0.2s; }
  .btn-action-primary:hover:not(:disabled) { background: var(--primary-hover); transform: translateY(-1px); }
  .btn-action-primary:disabled { opacity: 0.5; cursor: not-allowed; }

  .status-error { font-size: 0.8rem; color: #ef4444; margin-top: 0.75rem; display: block; font-weight: 600; }
  .empty-collabs { text-align: center; padding: 3rem; color: var(--text-muted); font-size: 0.9rem; opacity: 0.6; }
</style>
