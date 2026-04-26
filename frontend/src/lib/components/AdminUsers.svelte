<script lang="ts">
  import { apiFetch } from "$lib/api";
  import { onMount } from 'svelte';
  import AdminUserEditor from './AdminUserEditor.svelte';
  import Avatar from './Avatar.svelte';
  import Button from './ui/Button.svelte';
  import Card from './ui/Card.svelte';
  import TableSkeleton from './ui/TableSkeleton.svelte';
  import { i18n } from '$lib/i18n.svelte';
  import { dialog } from '$lib/dialog.svelte';
  
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
    if (!await dialog.confirm('¿Estás seguro de que deseas eliminar este usuario? Esta acción es IRREVERSIBLE.')) return;
    try {
      const res = await apiFetch(`/api/v2/admin/users/${id}`, { method: 'DELETE' });
      if (res.ok) fetchUsers();
    } catch (e) {
      console.error(e);
    }
  }

  onMount(fetchUsers);
</script>

<Card class="settings-premium-card" padding="lg">
  <header class="section-header">
    <div class="header-flex">
      <div class="title-with-icon">
        <span class="material-symbols-rounded hero-icon">person</span>
        <div>
          <h2>{i18n.t('admin.users')}</h2>
          <p class="subtitle">Administra las cuentas de usuario, permisos globales y estado del sistema</p>
        </div>
      </div>
      <Button variant="primary" onclick={() => editingUserId = 'new'} icon="person_add">{i18n.t('admin.create_user')}</Button>
    </div>
  </header>

  {#if editingUserId}
    <div class="editor-overlay island tonal">
       <AdminUserEditor 
         userId={editingUserId === 'new' ? null : editingUserId} 
         onSave={fetchUsers} 
         onCancel={() => editingUserId = null} 
       />
    </div>
  {:else}
    <div class="modern-table-card">
      <table class="modern-table">
        <thead>
          <tr>
            <th>ID</th>
            <th>{i18n.t('username')}</th>
            <th>{i18n.t('email')}</th>
            <th>{i18n.t('user.join_on')}</th>
            <th>Admin</th>
            <th>Acciones</th>
          </tr>
        </thead>
        <tbody>
          {#if loading}
            <TableSkeleton columns={6} rows={5} />
          {:else}
            {#each users as user}
              <tr class="hover-row">
                <td class="id-cell" data-label="ID">#{user.id}</td>
                <td data-label={i18n.t('username')}>
                  <div class="user-profile-cell">
                    <Avatar src={user.avatar_url} username={user.username} size={40} class="modular-avatar" />
                    <div class="user-meta">
                      <strong>{user.username}</strong>
                      <span>{user.full_name || 'Sin nombre completo'}</span>
                    </div>
                  </div>
                </td>
                <td class="email-cell" data-label={i18n.t('email')}>{user.email}</td>
                <td class="date-cell" data-label={i18n.t('user.join_on')}>{new Date(user.created_at).toLocaleDateString()}</td>
                <td data-label="Admin">
                  <button class="admin-badge" class:is-admin={user.is_admin} onclick={() => toggleAdmin(user)}>
                    {user.is_admin ? 'Administrador' : 'Usuario'}
                  </button>
                </td>
                <td data-label="Acciones">
                  <div class="action-buttons">
                    <Button variant="ghost" size="sm" icon="edit" onclick={() => editingUserId = user.id} />
                    <Button variant="ghost" size="sm" icon="delete" class="danger-btn" onclick={() => deleteUser(user.id)} />
                  </div>
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
  
  .section-header { margin-bottom: 4rem; }
  .header-flex { display: flex; justify-content: space-between; align-items: center; gap: 2rem; }
  @media (max-width: 600px) {
    .header-flex { flex-direction: column; align-items: stretch; text-align: center; }
    .title-with-icon { flex-direction: column; gap: 1rem; }
  }

  .title-with-icon { display: flex; align-items: center; gap: 2rem; }
  .hero-icon { font-size: 3rem; color: var(--primary-color); background: rgba(var(--primary-rgb), 0.1); padding: 1.25rem; border-radius: 2rem; flex-shrink: 0; }
  
  h2 { margin: 0; font-size: 2.25rem; font-weight: 950; letter-spacing: -0.04em; }
  .subtitle { margin: 0.5rem 0 0 0; color: var(--text-muted); font-size: 1.1rem; }

  @media (max-width: 600px) {
    .section-header { margin-bottom: 2rem; }
    h2 { font-size: 1.5rem; }
    .hero-icon { font-size: 1.5rem; padding: 0.6rem; border-radius: 1rem; }
    .title-with-icon { gap: 0.75rem; }
    .subtitle { font-size: 0.9rem; }
    .editor-overlay { padding: 1.5rem; border-radius: 1.5rem; margin-bottom: 2rem; }
  }

  @media (max-width: 1024px) {
    .user-profile-cell { justify-content: flex-end; }
    .action-buttons { justify-content: flex-end; }
  }

  .id-cell { font-family: 'JetBrains Mono', monospace; color: var(--text-muted); font-weight: 700; opacity: 0.5; }
  
  .user-profile-cell { display: flex; align-items: center; gap: 1.25rem; }
  :global(.modular-avatar) { border-radius: 1rem !important; border: 1px solid var(--md-sys-color-outline-variant); }
  .user-meta { display: flex; flex-direction: column; }
  .user-meta strong { font-size: 1.1rem; font-weight: 800; color: var(--text-color); }
  .user-meta span { font-size: 0.85rem; color: var(--text-muted); font-weight: 600; }

  .email-cell { font-weight: 600; color: var(--text-color); opacity: 0.8; }
  .date-cell { font-size: 0.85rem; font-weight: 700; color: var(--text-muted); }

  .admin-badge { border: none; font-family: inherit; cursor: pointer; padding: 0.35rem 1rem; border-radius: 100px; font-size: 0.75rem; font-weight: 900; text-transform: uppercase; background: rgba(255,255,255,0.05); color: var(--text-muted); border: 1px solid var(--md-sys-color-outline-variant); transition: all 0.3s; }
  .admin-badge.is-admin { background: rgba(var(--primary-rgb), 0.15); color: var(--primary-color); border-color: var(--primary-color); }
  .admin-badge:hover { transform: scale(1.05); filter: brightness(1.2); }

  .action-buttons { display: flex; gap: 0.5rem; }
  :global(.danger-btn:hover) { color: var(--error) !important; background: rgba(var(--error-rgb), 0.1) !important; }

  .editor-overlay { margin-bottom: 3rem; padding: 3rem; border-radius: 3rem; border: 2px solid var(--primary-color); background: rgba(var(--primary-rgb), 0.03); }
</style>
