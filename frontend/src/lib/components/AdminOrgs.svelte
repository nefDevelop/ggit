<script lang="ts">
  import { apiFetch } from "$lib/api";
  import { onMount } from 'svelte';
  import Card from './ui/Card.svelte';
  import Button from './ui/Button.svelte';
  import Avatar from './Avatar.svelte';
  import AdminOrgEditor from './AdminOrgEditor.svelte';
  import { i18n } from '$lib/i18n.svelte';
  
  let orgs = $state<any[]>([]);
  let loading = $state(true);
  let page = $state(1);
  let editingOrg = $state<any | null>(null);

  async function fetchOrgs() {
    loading = true;
    try {
      const res = await apiFetch(`/api/v2/admin/orgs?page=${page}`);
      if (res.ok) {
        orgs = await res.json();
      }
    } catch (e) {
      console.error(e);
    } finally {
      loading = false;
    }
  }

  async function deleteOrg(id: number) {
    if (!confirm('¿Estás seguro de que deseas eliminar esta organización? Todos sus contenidos serán eliminados permanentemente.')) return;
    try {
      const res = await apiFetch(`/api/v2/admin/users/${id}`, { method: 'DELETE' });
      if (res.ok) fetchOrgs();
    } catch (e) {
      console.error(e);
    }
  }

  onMount(fetchOrgs);
</script>

<Card class="settings-premium-card" padding="lg">
  <header class="section-header">
    <div class="title-with-icon">
      <span class="material-symbols-rounded hero-icon">corporate_fare</span>
      <div>
        <h2>{i18n.t('admin.organizations')}</h2>
        <p class="subtitle">Gestión y supervisión de las entidades organizacionales</p>
      </div>
    </div>
  </header>

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
          {#each Array(5) as _}
            <tr class="skeleton-row">
              <td colspan="5"><div class="skeleton-line"></div></td>
            </tr>
          {/each}
        {:else}
          {#each orgs as org}
            <tr class="hover-row">
              <td class="id-cell">#{org.id}</td>
              <td>
                <div class="user-cell">
                  <Avatar src={org.avatar_url} username={org.username} size={48} />
                  <div class="names">
                    <strong>{org.username}</strong>
                    <span>Reg. {new Date(org.created_at).toLocaleDateString()}</span>
                  </div>
                </div>
              </td>
              <td>
                <span class="full-name">{org.full_name || '—'}</span>
              </td>
              <td>
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
              <td>
                <div class="action-row">
                  <Button variant="ghost" size="sm" icon="visibility" href="/org/{org.username}" />
                  <Button variant="ghost" size="sm" icon="edit" onclick={() => editingOrg = org} />
                  <Button variant="ghost" size="sm" icon="delete" class="danger-btn" onclick={() => deleteOrg(org.id)} />
                </div>
              </td>
            </tr>
          {/each}
        {/if}
      </tbody>
    </table>
  </div>
</Card>

{#if editingOrg}
  <AdminOrgEditor 
    org={editingOrg} 
    onSave={() => { editingOrg = null; fetchOrgs(); }} 
    onCancel={() => editingOrg = null} 
  />
{/if}

<style>
  :global(.settings-premium-card) { border-radius: 3rem; background: rgba(var(--surface-rgb), 0.5); backdrop-filter: blur(24px); border: 1px solid var(--md-sys-color-outline-variant); }
  
  .section-header { margin-bottom: 4rem; }
  .title-with-icon { display: flex; align-items: center; gap: 2rem; }
  .hero-icon { font-size: 3rem; color: var(--primary-color); background: rgba(var(--primary-rgb), 0.1); padding: 1.25rem; border-radius: 2rem; }
  
  h2 { margin: 0; font-size: 2.25rem; font-weight: 950; letter-spacing: -0.04em; }
  .subtitle { margin: 0.5rem 0 0 0; color: var(--text-muted); font-size: 1.1rem; }

  .modern-table-card { border: 1px solid var(--md-sys-color-outline-variant); border-radius: 2.5rem; overflow: hidden; background: rgba(0,0,0,0.15); box-shadow: inset 0 2px 4px rgba(0,0,0,0.1); }
  
  .modern-table { width: 100%; border-collapse: collapse; text-align: left; }
  .modern-table th { padding: 1.5rem; background: rgba(var(--primary-rgb), 0.03); color: var(--text-muted); font-size: 0.8rem; text-transform: uppercase; letter-spacing: 0.15em; font-weight: 900; border-bottom: 1px solid var(--md-sys-color-outline-variant); }
  .modern-table td { padding: 1.75rem 1.5rem; border-bottom: 1px solid var(--md-sys-color-outline-variant); vertical-align: middle; }
  .modern-table tr:last-child td { border-bottom: none; }

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

  .skeleton-line { height: 2.5rem; background: var(--surface-variant); border-radius: 1rem; animation: pulse 1.5s infinite; }
</style>
