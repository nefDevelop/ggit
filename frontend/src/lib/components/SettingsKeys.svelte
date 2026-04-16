<script lang="ts">
  import { onMount } from 'svelte';
  
  let keys = $state<any[]>([]);
  let loading = $state(true);
  let newKeyTitle = $state('');
  let newKeyContent = $state('');
  let submitting = $state(false);

  async function fetchKeys() {
    loading = true;
    try {
      const res = await fetch('/api/v2/user/settings/ssh_keys');
      if (res.ok) keys = await res.json();
    } catch (e) {
      console.error(e);
    } finally {
      loading = false;
    }
  }

  async function addKey() {
    if (!newKeyTitle || !newKeyContent) return;
    submitting = true;
    try {
      const res = await fetch('/api/v2/user/settings/ssh_keys', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ title: newKeyTitle, content: newKeyContent })
      });
      if (res.ok) {
        const key = await res.json();
        keys = [...keys, key];
        newKeyTitle = '';
        newKeyContent = '';
      }
    } catch (e) {
      console.error(e);
    } finally {
      submitting = false;
    }
  }

  async function deleteKey(id: number) {
    if (!confirm('Are you sure you want to delete this key?')) return;
    try {
      const res = await fetch(`/api/v2/user/settings/ssh_keys/${id}`, {
        method: 'DELETE'
      });
      if (res.ok) {
        keys = keys.filter(k => k.id !== id);
      }
    } catch (e) {
      console.error(e);
    }
  }

  onMount(fetchKeys);
</script>

<div class="settings-card island">
  <header class="card-header">
    <h3 class="card-title">SSH Keys</h3>
  </header>
  
  <div class="card-body">
    <p class="desc">SSH keys allow you to establish a secure connection between your computer and Ggit.</p>

    <div class="keys-list">
      {#each keys as key}
        <div class="key-row island tonal">
          <div class="key-icon">
            <span class="material-symbols-rounded">vpn_key</span>
          </div>
          <div class="key-info">
            <strong>{key.title}</strong>
            <code>{key.fingerprint || 'Fingerprint pending'}</code>
            <span class="key-meta">Added on {new Date(key.created_at).toLocaleDateString()}</span>
          </div>
          <button class="btn btn-icon btn-danger-text" onclick={() => deleteKey(key.id)}>
            <span class="material-symbols-rounded">delete</span>
          </button>
        </div>
      {/each}
      {#if !loading && keys.length === 0}
        <div class="empty-state">No SSH keys added yet.</div>
      {/if}
    </div>

    <hr class="divider" />

    <form class="add-key-form" onsubmit={(e) => { e.preventDefault(); addKey(); }}>
      <h4>Add SSH Key</h4>
      <div class="form-group">
        <label for="key_title">Title</label>
        <input type="text" id="key_title" bind:value={newKeyTitle} placeholder="Work Laptop" required />
      </div>
      <div class="form-group">
        <label for="key_content">Key Content</label>
        <textarea id="key_content" bind:value={newKeyContent} placeholder="Begins with 'ssh-rsa', 'ssh-ed25519', etc." required></textarea>
      </div>
      <div class="form-actions">
        <button type="submit" class="btn btn-primary" disabled={submitting || !newKeyTitle || !newKeyContent}>
          {submitting ? 'Adding...' : 'Add Key'}
        </button>
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

  .keys-list { display: flex; flex-direction: column; gap: 1rem; margin-bottom: 2rem; }
  .key-row { display: flex; align-items: center; gap: 1.5rem; padding: 1rem 1.5rem; }
  .key-icon { font-size: 1.5rem; color: var(--primary-color); opacity: 0.8; }
  .key-info { flex: 1; display: flex; flex-direction: column; gap: 0.25rem; min-width: 0; }
  .key-info strong { font-size: 1rem; }
  .key-info code { font-size: 0.75rem; background: rgba(0,0,0,0.3); padding: 0.1rem 0.4rem; border-radius: 4px; opacity: 0.6; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }
  .key-meta { font-size: 0.75rem; color: var(--text-muted); opacity: 0.5; }

  .divider { border: 0; border-top: 1px solid var(--border-color); margin: 2rem 0; }

  .add-key-form { display: flex; flex-direction: column; gap: 1.5rem; }
  .add-key-form h4 { margin: 0; font-size: 1.1rem; }
  .form-group { display: flex; flex-direction: column; gap: 0.5rem; }
  .form-group label { font-size: 0.8rem; font-weight: 800; text-transform: uppercase; color: var(--text-muted); }
  
  input, textarea { background: rgba(0,0,0,0.2); border: 1px solid var(--border-color); padding: 0.75rem 1rem; border-radius: 0.75rem; color: var(--text-color); font-family: inherit; }
  textarea { min-height: 120px; font-family: monospace; font-size: 0.9rem; }
  input:focus, textarea:focus { outline: none; border-color: var(--primary-color); background: rgba(0,0,0,0.4); }

  .form-actions { display: flex; justify-content: flex-end; }
  .empty-state { text-align: center; padding: 2rem; color: var(--text-muted); border: 2px dashed var(--border-color); border-radius: 1rem; }
  
  .btn-danger-text { color: #ef4444; }
  .btn-danger-text:hover { background: rgba(239, 68, 68, 0.1); }
</style>
