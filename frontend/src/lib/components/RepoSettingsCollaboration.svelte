<script lang="ts">
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
      const res = await fetch(`/api/v2/${repo.data.owner.username}/${repo.data.name}/settings/collaborators`);
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
      const res = await fetch(`/api/v2/${repo.data.owner.username}/${repo.data.name}/settings/collaborators/${newCollaborator}`, {
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
      const res = await fetch(`/api/v2/${repo.data.owner.username}/${repo.data.name}/settings/collaborators/${username}`, {
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

<div class="settings-card island">
  <header class="card-header">
    <h3 class="card-title">Manage Collaborators</h3>
  </header>
  
  <div class="card-body">
    <p class="desc">Collaborators have read/write access to this repository.</p>

    <div class="collaborator-list">
      {#each collaborators as col}
        <div class="col-row island tonal">
          <img src={col.avatar_url} alt={col.username} class="avatar" />
          <div class="col-info">
            <strong>{col.full_name || col.username}</strong>
            <span class="col-meta">@{col.username}</span>
          </div>
          <div class="perm-badge">Write</div>
          <button class="btn btn-icon btn-danger-text" onclick={() => removeCollaborator(col.username)}>
            <span class="material-symbols-rounded">person_remove</span>
          </button>
        </div>
      {/each}
      {#if !loading && collaborators.length === 0}
        <div class="empty-state">No collaborators added yet.</div>
      {/if}
    </div>

    <hr class="divider" />

    <form class="add-col-form" onsubmit={(e) => { e.preventDefault(); addCollaborator(); }}>
      <h4>Add Collaborator</h4>
      <div class="form-group">
        <label for="col_name">Username</label>
        <div class="inline-form">
          <input type="text" id="col_name" bind:value={newCollaborator} placeholder="Enter username..." required />
          <button type="submit" class="btn btn-primary" disabled={submitting || !newCollaborator}>
            {submitting ? 'Adding...' : 'Add User'}
          </button>
        </div>
        {#if message}
           <span class="error-msg">{message}</span>
        {/if}
      </div>
    </form>
  </div>
</div>

<style>
  .settings-card { display: flex; flex-direction: column; }
  .card-header { padding: 1.5rem 2rem; border-bottom: 1px solid var(--border-color); }
  .card-title { margin: 0; font-size: 1.25rem; font-weight: 700; }
  
  .card-body { padding: 2rem; }
  .desc { color: var(--text-muted); margin-bottom: 2rem; }

  .collaborator-list { display: flex; flex-direction: column; gap: 1rem; margin-bottom: 2rem; }
  .col-row { display: flex; align-items: center; gap: 1.25rem; padding: 0.75rem 1.25rem; }
  .avatar { width: 40px; height: 40px; border-radius: 50%; background: #222; }
  
  .col-info { flex: 1; display: flex; flex-direction: column; gap: 0.1rem; }
  .col-info strong { font-size: 0.95rem; }
  .col-meta { font-size: 0.75rem; color: var(--text-muted); }

  .perm-badge { background: rgba(var(--primary-rgb), 0.1); color: var(--primary-color); padding: 0.2rem 0.6rem; border-radius: 4px; font-size: 0.7rem; font-weight: 800; text-transform: uppercase; }

  .divider { border: 0; border-top: 1px solid var(--border-color); margin: 2rem 0; }

  .add-col-form { display: flex; flex-direction: column; gap: 1rem; }
  .add-col-form h4 { margin: 0; font-size: 1.1rem; margin-bottom: 0.5rem; }
  
  .inline-form { display: flex; gap: 1rem; }
  .inline-form input { flex: 1; }
  
  input { background: rgba(0,0,0,0.2); border: 1px solid var(--border-color); padding: 0.75rem 1rem; border-radius: 0.75rem; color: var(--text-color); font-family: inherit; }
  input:focus { outline: none; border-color: var(--primary-color); background: rgba(0,0,0,0.4); }

  .error-msg { font-size: 0.8rem; color: #ef4444; margin-top: 0.5rem; }
  .empty-state { text-align: center; padding: 2rem; color: var(--text-muted); border: 2px dashed var(--border-color); border-radius: 1rem; }
  
  .btn-danger-text { color: #ef4444; }
  .btn-danger-text:hover { background: rgba(239, 68, 68, 0.1); }
</style>
