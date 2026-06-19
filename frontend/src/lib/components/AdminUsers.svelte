<script lang="ts">
  

import User from "lucide-svelte/icons/user";
import UserPlus from "lucide-svelte/icons/user-plus";
import Edit from "lucide-svelte/icons/edit";
import Trash2 from "lucide-svelte/icons/trash-2";


  import { apiFetch } from "$lib/api";
  import { onMount } from 'svelte';
  import AdminUserEditor from './AdminUserEditor.svelte';
  import Avatar from './Avatar.svelte';
  import Button from './ui/Button.svelte';
  import Card from './ui/Card.svelte';
  import TableSkeleton from './ui/TableSkeleton.svelte';
  import { i18n } from '$lib/i18n.svelte';
  import { dialog } from '$lib/dialog.svelte';
  import PageHeader from './ui/PageHeader.svelte';
  
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
        users = (await res.json()) || [];
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

<div class="admin-users-container">
  <PageHeader 
    title={i18n.t('admin.users') || "User Management"}
    subtitle="Administra las cuentas de usuario, permisos globales y estado del sistema"
    icon={User}
  >
    {#snippet actions()}
      <Button variant="primary" onclick={() => editingUserId = 'new'} icon={UserPlus}>{i18n.t('admin.create_user')}</Button>
    {/snippet}
  </PageHeader>

  <div class="admin-users-body">
  {#if editingUserId}
    <div class="editor-overlay">
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
                    <button class="icon-btn" onclick={() => editingUserId = user.id} title="Editar">
                      <Edit size={18} />
                    </button>
                    <button class="icon-btn danger" onclick={() => deleteUser(user.id)} title="Eliminar">
                      <Trash2 size={18} />
                    </button>
                  </div>
                </td>
              </tr>
            {/each}
          {/if}
        </tbody>
      </table>
    </div>
  {/if}
  </div>
</div>

<style>
  .admin-users-body {
    padding: 0 1.5rem 4rem;
    animation: fadeIn 0.4s cubic-bezier(0.4, 0, 0.2, 1);
  }

  @keyframes fadeIn {
    from { opacity: 0; transform: translateY(4px); }
    to { opacity: 1; transform: translateY(0); }
  }
  
  @media (max-width: 600px) {
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

  .admin-badge { border: none; font-family: inherit; cursor: pointer; padding: 0.35rem 1rem; border-radius: 0.75rem; font-size: 0.75rem; font-weight: 900; text-transform: uppercase; background: rgba(255,255,255,0.05); color: var(--text-muted); border: 1px solid var(--md-sys-color-outline-variant); transition: all 0.3s; }
  .admin-badge.is-admin { background: rgba(var(--primary-rgb), 0.15); color: var(--primary-color); border-color: var(--primary-color); }
  .admin-badge:hover { transform: scale(1.05); filter: brightness(1.2); }

  .action-buttons { display: flex; gap: 0.5rem; }
  .icon-btn {
    background: transparent;
    border: none;
    color: var(--text-muted);
    cursor: pointer;
    padding: 0.5rem;
    display: flex;
    align-items: center;
    justify-content: center;
    border-radius: 0.5rem;
    transition: all 0.2s;
  }
  .icon-btn:hover {
    background: var(--md-sys-color-surface-variant);
    color: var(--text-color);
  }
  .icon-btn.danger:hover {
    background: rgba(var(--error-rgb), 0.1);
    color: var(--error);
  }

  .editor-overlay { margin-bottom: 3rem; padding: 3rem; border-radius: 2rem; border: 1px solid var(--md-sys-color-outline-variant); background: var(--md-sys-color-surface-container-high); }
</style>
