<script lang="ts">
  import { apiFetch } from "$lib/api";
  import { onMount } from 'svelte';
  import { dialog } from '$lib/dialog.svelte';
  import { formatDate } from '$lib/utils/date';
  import Card from './ui/Card.svelte';
  import Button from './ui/Button.svelte';
  import Avatar from './Avatar.svelte';
  import AdminOrgEditor from './AdminOrgEditor.svelte';
  import { i18n } from '$lib/i18n.svelte';
  import TableSkeleton from './ui/TableSkeleton.svelte';
  import PageHeader from './ui/PageHeader.svelte';
  

import Building2 from "lucide-svelte/icons/building-2";
import Eye from "lucide-svelte/icons/eye";
import Edit2 from "lucide-svelte/icons/edit-2";
import Trash2 from "lucide-svelte/icons/trash-2";


  
  let orgs = $state<any[]>([]);
  let loading = $state(true);
  let page = $state(1);
  let editingOrg = $state<any | null>(null);

  async function fetchOrgs() {
    loading = true;
    try {
      const res = await apiFetch(`/api/v2/admin/orgs?page=${page}`);
      if (res.ok) {
        const json = await res.json();
        orgs = Array.isArray(json) ? json : [];
      }
    } catch (e) {
      console.error(e);
    } finally {
      loading = false;
    }
  }

  async function deleteOrg(id: number) {
    if (!await dialog.confirm('¿Estás seguro de que deseas eliminar esta organización? Todos sus contenidos serán eliminados permanentemente.')) return;
    try {
      const res = await apiFetch(`/api/v2/admin/users/${id}`, { method: 'DELETE' });
      if (res.ok) fetchOrgs();
    } catch (e) {
      console.error(e);
    }
  }

  onMount(fetchOrgs);
</script>

<div class="admin-orgs-container">
  <PageHeader 
    title={i18n.t('admin.organizations') || "Organizations"}
    subtitle="Gestión y supervisión de las entidades organizacionales"
    icon={Building2}
  />

  <div class="admin-orgs-body">

  <div class="modern-table-card">
    <table class="modern-table">
      <thead>
        <tr>
          <th>ID</th>
          <th>{i18n.t('organization')}</th>
          <th>{i18n.t('settings.full_name')}</th>
          <th>Métricas</th>
          <th>Acciones</th>
        </tr>
      </thead>
      <tbody>
        {#if loading}
          <TableSkeleton columns={5} rows={5} />
        {:else}
          {#each orgs as org}
            <tr class="hover-row">
              <td class="id-cell" data-label="ID">#{org.id}</td>
              <td data-label={i18n.t('organization')}>
                <div class="user-cell">
                  <Avatar src={org.avatar_url} username={org.username} size={48} />
                  <div class="names">
                    <strong>{org.username}</strong>
                    <span>Reg. {formatDate(org.created_at)}</span>
                  </div>
                </div>
              </td>
              <td data-label={i18n.t('settings.full_name')}>
                <span class="full-name">{org.full_name || '—'}</span>
              </td>
              <td data-label="Métricas">
                <div class="metrics-grid">
                  <div class="metric">
                    <span class="val">{org.num_teams}</span>
                    <span class="lab">Equipos</span>
                  </div>
                  <div class="metric">
                    <span class="val">{org.num_members}</span>
                    <span class="lab">Miembros</span>
                  </div>
                </div>
              </td>
              <td data-label="Acciones">
                <div class="action-row">
                  <Button variant="ghost" size="sm" icon={Eye} href="/org/{org.username}" />
                  <Button variant="ghost" size="sm" icon={Edit2} onclick={() => editingOrg = org} />
                  <Button variant="ghost" size="sm" icon={Trash2} class="danger-btn" onclick={() => deleteOrg(org.id)} />
                </div>
              </td>
            </tr>
          {/each}
        {/if}
      </tbody>
    </table>
  </div>
  </div>
</div>

{#if editingOrg}
  <AdminOrgEditor 
    org={editingOrg} 
    onSave={() => { editingOrg = null; fetchOrgs(); }} 
    onCancel={() => editingOrg = null} 
  />
{/if}

<style>
  .admin-orgs-body {
    padding: 0 1.5rem 4rem;
    animation: fadeIn 0.4s cubic-bezier(0.4, 0, 0.2, 1);
  }

  @keyframes fadeIn {
    from { opacity: 0; transform: translateY(4px); }
    to { opacity: 1; transform: translateY(0); }
  }
  
  @media (max-width: 600px) {
    /* Mobile adjustments if needed */
  }

  @media (max-width: 1024px) {
    .user-cell { justify-content: flex-end; }
    .metrics-grid { justify-content: flex-end; }
    .action-row { justify-content: flex-end; }
  }


  .id-cell { font-family: monospace; color: var(--text-muted); font-weight: 700; opacity: 0.5; }

  .user-cell { display: flex; align-items: center; gap: 1.25rem; }
  .names { display: flex; flex-direction: column; }
  .names strong { color: var(--text-color); line-height: 1.2; font-weight: 800; font-size: 1.1rem; }
  .names span { font-size: 0.85rem; color: var(--text-muted); }

  .full-name { font-weight: 600; color: var(--text-color); }

  .metrics-grid { display: flex; gap: 2.5rem; }
  .metric { display: flex; flex-direction: column; align-items: center; min-width: 60px; }
  .metric .val { font-size: 1.25rem; font-weight: 900; color: var(--primary-color); }
  .metric .lab { font-size: 0.7rem; text-transform: uppercase; font-weight: 800; color: var(--text-muted); letter-spacing: 0.05em; }

  .action-row { display: flex; gap: 0.75rem; }
  :global(.danger-btn:hover) { color: var(--error) !important; background: rgba(var(--error-rgb), 0.15) !important; }
</style>
