<script lang="ts">
  import { apiFetch } from "$lib/api";
  import { onMount } from 'svelte';
  import AdminAuthEditor from './AdminAuthEditor.svelte';
  
  let sources = $state<any[]>([]);
  let loading = $state(true);
  let showEditor = $state(false);
  let editingAuthId = $state<number | null>(null);

  async function fetchAuths() {
    loading = true;
    try {
      const res = await apiFetch(`/api/v2/admin/auths`);
      if (res.ok) {
        sources = await res.json();
      }
    } catch (e) {
      console.error(e);
    } finally {
      loading = false;
    }
  }

  async function deleteAuth(id: number) {
    if (!confirm('Are you sure you want to delete this login source?')) return;
    try {
      const res = await apiFetch(`/api/v2/admin/auths/${id}`, { method: 'DELETE' });
      if (res.ok) fetchAuths();
    } catch (e) {
      console.error(e);
    }
  }

  function handleSave() {
    showEditor = false;
    editingAuthId = null;
    fetchAuths();
  }

  function handleCancel() {
    showEditor = false;
    editingAuthId = null;
  }

  function startEdit(id: number) {
    editingAuthId = id;
    showEditor = true;
  }

  function startCreate() {
    editingAuthId = null;
    showEditor = true;
  }

  onMount(fetchAuths);
</script>

{#if showEditor}
  <AdminAuthEditor authId={editingAuthId} onSave={handleSave} onCancel={handleCancel} />
{/if}

<div class="admin-auth island island-interactive tonal">
  <header class="section-header">
    <div class="title">
      <span class="material-symbols-rounded">vpn_key</span>
      <h2>Authentication Sources</h2>
    </div>
    <button class="btn btn-primary">Add New Source</button>
  </header>

  <div class="table-container">
    <table class="modern-table">
      <thead>
        <tr>
          <th>ID</th>
          <th>Name</th>
          <th>Type</th>
          <th>Enabled</th>
          <th>Created</th>
          <th>Actions</th>
        </tr>
      </thead>
      <tbody>
        {#if loading}
          {#each Array(3) as _}
            <tr class="skeleton-row">
              <td colspan="6"><div class="skeleton-bar"></div></td>
            </tr>
          {/each}
        {:else if sources.length === 0}
           <tr><td colspan="6" class="empty">No authentication sources configured.</td></tr>
        {:else}
          {#each sources as auth}
            <tr>
              <td>{auth.id}</td>
              <td><b>{auth.name}</b></td>
              <td><span class="type-badge">{auth.type_name}</span></td>
              <td>
                <span class="status-dot" class:active={auth.is_activated}></span>
                {auth.is_activated ? 'Enabled' : 'Disabled'}
              </td>
              <td>{new Date(auth.created_at).toLocaleDateString()}</td>
              <td class="actions">
                <button class="icon-btn" onclick={() => startEdit(auth.id)} title="Edit"><span class="material-symbols-rounded">settings</span></button>
                <button class="icon-btn danger" onclick={() => deleteAuth(auth.id)} title="Delete"><span class="material-symbols-rounded">delete</span></button>
              </td>
            </tr>
          {/each}
        {/if}
      </tbody>
    </table>
  </div>
</div>

<style>
  .admin-auth { padding: 2rem; display: flex; flex-direction: column; gap: 2rem; box-shadow: var(--shadow-xl); }
  
  .section-header { display: flex; justify-content: space-between; align-items: center; }
  .title { display: flex; align-items: center; gap: 0.75rem; }
  .title span { font-size: 1.75rem; color: var(--primary-color); }
  .title h2 { margin: 0; font-size: 1.5rem; font-weight: 800; }

  .table-container { overflow-x: auto; }
  .modern-table { width: 100%; border-collapse: collapse; text-align: left; }
  .modern-table th { padding: 1rem; border-bottom: 1px solid var(--md-sys-color-outline-variant); color: var(--text-muted); font-size: 0.75rem; text-transform: uppercase; letter-spacing: 0.1em; font-weight: 800; }
  .modern-table td { padding: 1.25rem 1rem; border-bottom: 1px solid rgba(255,255,255,0.03); font-size: 0.95rem; }

  .type-badge { background: rgba(var(--primary-rgb), 0.1); color: var(--primary-color); padding: 0.2rem 0.6rem; border-radius: 4px; font-weight: 800; font-size: 0.7rem; }
  
  .status-dot { display: inline-block; width: 0.6rem; height: 0.6rem; border-radius: 50%; background: #ef4444; margin-right: 0.5rem; }
  .status-dot.active { background: #22c55e; box-shadow: 0 0 8px #22c55e; }

  .actions { display: flex; gap: 0.5rem; }
  .icon-btn { background: transparent; border: none; color: var(--text-muted); cursor: pointer; padding: 0.4rem; border-radius: 0.5rem; transition: all 0.2s; }
  .icon-btn:hover { background: rgba(255,255,255,0.05); color: #fff; }
  .icon-btn.danger:hover { background: rgba(239, 68, 68, 0.1); color: #ef4444; }

  .empty { text-align: center; padding: 3rem; color: var(--text-muted); }
  .hidden { display: none; }
  .skeleton-bar { height: 1.5rem; background: rgba(255,255,255,0.03); border-radius: 0.5rem; animation: pulse 1.5s infinite; }
  @keyframes pulse { 0% { opacity: 0.3; } 50% { opacity: 0.6; } 100% { opacity: 0.3; } }
</style>
