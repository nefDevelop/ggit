<script lang="ts">
  import { apiFetch } from "$lib/api";
  import { onMount } from 'svelte';
  import AdminUserEditor from './AdminUserEditor.svelte';
  import Avatar from './Avatar.svelte';
  import Button from './ui/Button.svelte';
  import Card from './ui/Card.svelte';
  import { i18n } from '$lib/i18n.svelte';
  
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
      <h2>{i18n.t('admin.users')}</h2>
    </div>
    <Button onclick={() => editingUserId = 'new'} icon="person_add">{i18n.t('admin.create_user')}</Button>
  </header>

  {#if editingUserId}
    <div class="editor-overlay">
       <AdminUserEditor 
         userId={editingUserId === 'new' ? null : editingUserId} 
         onSave={fetchUsers} 
         onCancel={() => editingUserId = null} 
       />
    </div>
  {:else}
    <div class="table-container">
      <table class="modern-table">
        <thead>
          <tr>
            <th>ID</th>
            <th>{i18n.t('username')}</th>
            <th>{i18n.t('email')}</th>
            <th>{i18n.t('user.join_on')}</th>
            <th>Admin</th>
            <th>{i18n.t('admin.save_user')}</th>
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
                    <Avatar src={user.avatar_url} username={user.username} size={32} class="small-avatar" />
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
  {/if}
</Card>

<style>
  .admin-users { padding: 2.5rem; display: flex; flex-direction: column; gap: 2.5rem; border-radius: 2.5rem; }
  
  .section-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 0.5rem; }
  .title { display: flex; align-items: center; gap: 1rem; }
  .title span { font-size: 2rem; color: var(--primary-color); background: rgba(var(--primary-rgb), 0.1); padding: 0.75rem; border-radius: 1rem; }
  .title h2 { margin: 0; font-size: 1.75rem; font-weight: 900; letter-spacing: -0.02em; }

  .table-container { overflow-x: auto; border-radius: 1.5rem; background: rgba(var(--surface-rgb), 0.3); border: 1px solid var(--md-sys-color-outline-variant); }
  .modern-table { width: 100%; border-collapse: collapse; text-align: left; }
  .modern-table th { padding: 1.25rem 1.5rem; border-bottom: 1px solid var(--md-sys-color-outline-variant); color: var(--text-muted); font-size: 0.75rem; text-transform: uppercase; letter-spacing: 0.1em; font-weight: 800; background: rgba(0,0,0,0.05); }
  .modern-table td { padding: 1.25rem 1.5rem; border-bottom: 1px solid var(--md-sys-color-outline-variant); font-size: 0.95rem; }
  .modern-table tr:last-child td { border-bottom: none; }

  .user-cell { display: flex; align-items: center; gap: 1rem; }
  .names { display: flex; flex-direction: column; }
  .names strong { color: var(--text-color); line-height: 1.2; font-weight: 700; }
  .names span { font-size: 0.8rem; color: var(--text-muted); }

  .actions { display: flex; gap: 0.5rem; }
  :global(.danger-icon:hover) { color: var(--error) !important; }

  .skeleton-bar { height: 1.5rem; background: var(--surface-variant); border-radius: 0.5rem; animation: pulse 1.5s infinite; }
  @keyframes pulse { 0% { opacity: 0.3; } 50% { opacity: 0.6; } 100% { opacity: 0.3; } }

  .editor-overlay { animation: fadeIn 0.3s ease-out; }
  @keyframes fadeIn { from { opacity: 0; transform: scale(0.98); } to { opacity: 1; transform: scale(1); } }
</style>
