<script lang="ts">
  import { apiFetch } from "$lib/api";
  import { onMount } from 'svelte';
  import AdminUserEditor from './AdminUserEditor.svelte';
  import Button from './ui/Button.svelte';
  import Card from './ui/Card.svelte';
  import Badge from './ui/Badge.svelte';
  
  let users = $state<any[]>([]);
  let loading = $state(true);
  let page = $state(1);
  let editingUserId = $state<number | null | 'new'>(null);

  async function fetchUsers() {
    loading = true;
    editingUserId = null;
    try {
      const res = await apiFetch(`/api/v2/admin/users?page=${page}`);
      if (res.ok) {
        users = await res.json();
      }
    } catch (e) {
      console.error(e);
    } finally {
      loading = false;
    }
  }

  async function toggleAdmin(user: any) {
    try {
      const res = await apiFetch(`/api/v2/admin/users/${user.id}`, {
        method: 'PATCH',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ admin: !user.is_admin })
      });
      if (res.ok) fetchUsers();
    } catch (e) {
      console.error(e);
    }
  }

  async function deleteUser(id: number) {
    if (!confirm('Are you sure you want to delete this user? This action is IRREVERSIBLE.')) return;
    try {
      const res = await apiFetch(`/api/v2/admin/users/${id}`, { method: 'DELETE' });
      if (res.ok) fetchUsers();
    } catch (e) {
      console.error(e);
    }
  }

  onMount(fetchUsers);
</script>

<Card class="admin-users">
  <header class="section-header">
    <div class="title">
      <span class="material-symbols-rounded">group</span>
      <h2>User Management</h2>
    </div>
    <Button onclick={() => editingUserId = 'new'} icon="person_add">New User</Button>
  </header>

  {#if editingUserId}
    <div class="editor-overlay">
       <AdminUserEditor 
         userId={editingUserId === 'new' ? null : editingUserId} 
         onSave={fetchUsers} 
         onCancel={() => editingUserId = null} 
       />
    </div>
  {/if}

  <div class="table-container">
    <table class="modern-table">
      <thead>
        <tr>
          <th>ID</th>
          <th>User</th>
          <th>Email</th>
          <th>Joined</th>
          <th>Admin</th>
          <th>Actions</th>
        </tr>
      </thead>
      <tbody>
        {#if loading}
          {#each Array(5) as _}
            <tr class="skeleton-row">
              <td colspan="6"><div class="skeleton-bar"></div></td>
            </tr>
          {/each}
        {:else}
          {#each users as user}
            <tr>
              <td>{user.id}</td>
              <td>
                <div class="user-cell">
                  <img src={user.avatar_url} alt="" class="small-avatar" />
                  <div class="names">
                    <strong>{user.username}</strong>
                    <span>{user.full_name || ''}</span>
                  </div>
                </div>
              </td>
              <td>{user.email}</td>
              <td>{new Date(user.created_at).toLocaleDateString()}</td>
              <td>
                <Button 
                   size="sm"
                   variant={user.is_admin ? 'primary' : 'secondary'}
                   onclick={() => toggleAdmin(user)}
                >
                  {user.is_admin ? 'Admin' : 'User'}
                </Button>
              </td>
              <td class="actions">
                <Button variant="ghost" size="sm" icon="edit" onclick={() => editingUserId = user.id} />
                <Button variant="ghost" size="sm" icon="delete" class="danger-icon" onclick={() => deleteUser(user.id)} />
              </td>
            </tr>
          {/each}
        {/if}
      </tbody>
    </table>
  </div>
</Card>

<style>
  .admin-users { padding: 2rem; display: flex; flex-direction: column; gap: 2rem; }
  
  .section-header { display: flex; justify-content: space-between; align-items: center; }
  .title { display: flex; align-items: center; gap: 0.75rem; }
  .title span { font-size: 1.75rem; color: var(--md-sys-color-primary); }
  .title h2 { margin: 0; font-size: 1.5rem; font-weight: 800; }

  .table-container { overflow-x: auto; }
  .modern-table { width: 100%; border-collapse: collapse; text-align: left; }
  .modern-table th { padding: 1rem; border-bottom: 1px solid var(--md-sys-color-outline-variant); color: var(--text-muted); font-size: 0.75rem; text-transform: uppercase; letter-spacing: 0.1em; font-weight: 800; }
  .modern-table td { padding: 1.25rem 1rem; border-bottom: 1px solid var(--md-sys-color-surface-variant); font-size: 0.95rem; }

  .user-cell { display: flex; align-items: center; gap: 1rem; }
  .small-avatar { width: 32px; height: 32px; border-radius: 50%; background: var(--md-sys-color-background); border: 1px solid var(--md-sys-color-outline-variant); }
  .names { display: flex; flex-direction: column; }
  .names strong { color: var(--md-sys-color-on-surface); line-height: 1.2; }
  .names span { font-size: 0.8rem; color: var(--text-muted); }

  .actions { display: flex; gap: 0.25rem; }
  :global(.danger-icon:hover) { color: var(--color-danger) !important; }

  .skeleton-bar { height: 1.5rem; background: var(--md-sys-color-surface-variant); border-radius: 0.5rem; animation: pulse 1.5s infinite; }
  @keyframes pulse { 0% { opacity: 0.3; } 50% { opacity: 0.6; } 100% { opacity: 0.3; } }

  .editor-overlay { margin-bottom: 2rem; border-bottom: 1px solid var(--md-sys-color-outline-variant); padding-bottom: 2rem; animation: slideDown 0.3s ease-out; }
  @keyframes slideDown { from { opacity: 0; transform: translateY(-20px); } to { opacity: 1; transform: translateY(0); } }
</style>

