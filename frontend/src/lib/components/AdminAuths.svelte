<script lang="ts">
  import { apiFetch } from "$lib/api";
  import { onMount } from "svelte";
  import { i18n } from '$lib/i18n.svelte';
  import { dialog } from '$lib/dialog.svelte';
  import Card from './ui/Card.svelte';
  import Button from './ui/Button.svelte';
  import AdminAuthEditor from './AdminAuthEditor.svelte';
  import TableSkeleton from './ui/TableSkeleton.svelte';
  import PageHeader from './ui/PageHeader.svelte';
  import { ShieldHalf, Plus, Database, ShieldCheck, Settings, Trash2 } from 'lucide-svelte';

  let sources = $state<any[]>([]);
  let loading = $state(true);
  let editing = $state(false);
  let selectedSource = $state<any>(null);

  async function fetchAuths() {
    loading = true;
    editing = false;
    selectedSource = null;
    try {
      const res = await apiFetch("/api/v2/admin/auths");
      if (res.ok) sources = (await res.json()) || [];
    } catch (e) { console.error(e); } finally { loading = false; }
  }

  function editSource(source: any) {
    selectedSource = source;
    editing = true;
  }

  function newSource() {
    selectedSource = null;
    editing = true;
  }

  async function deleteAuth(id: number) {
    if (!await dialog.confirm("¿Seguro que quieres eliminar este origen de autenticación?")) return;
    try {
      const res = await apiFetch(`/api/v2/admin/auths/${id}`, { method: "DELETE" });
      if (res.ok) fetchAuths();
    } catch (e) { console.error(e); }
  }

  onMount(fetchAuths);
</script>

{#if editing}
  <AdminAuthEditor 
    source={selectedSource} 
    onSave={fetchAuths} 
    onCancel={() => editing = false} 
  />
{:else}
  <div class="admin-auths-container">
    <Card class="settings-premium-card" padding="none">
      <PageHeader 
        title={i18n.t('admin.auths') || "Fuentes de Autenticación"}
        subtitle="Gestión centralizada de protocolos LDAP, OAuth2, PAM y servicios externos"
        icon={ShieldHalf}
        compact
      >
        {#snippet actions()}
          <Button variant="primary" icon={Plus} onclick={newSource}>Nueva Fuente</Button>
        {/snippet}
      </PageHeader>

      <div class="admin-auths-body">
        <div class="modern-table-card">
          <table class="modern-table">
            <thead>
              <tr>
                <th>ID</th>
                <th>Nombre del Origen</th>
                <th>Tipo / Protocolo</th>
                <th>Estado</th>
                <th>Acciones</th>
              </tr>
            </thead>
            <tbody>
              {#if loading}
                <TableSkeleton columns={5} rows={3} />
              {:else if sources.length === 0}
                <tr>
                  <td colspan="5" class="empty-state">No hay orígenes de autenticación configurados.</td>
                </tr>
              {:else}
                {#each sources as source}
                  <tr class="hover-row">
                    <td class="id-cell">#{source.id}</td>
                    <td>
                      <div class="source-cell">
                        <div class="source-icon tonal">
                          {#if source.type === 2 || source.type === 5}
                            <Database size={24} />
                          {:else}
                            <ShieldCheck size={24} />
                          {/if}
                        </div>
                        <div class="source-info">
                          <strong>{source.name}</strong>
                          <span>{source.type === 2 ? 'Servidor LDAP' : source.type === 6 ? 'Servicio SMTP' : 'Autenticación Externa'}</span>
                        </div>
                      </div>
                    </td>
                    <td>
                      <span class="type-badge" style="--accent: {source.type === 2 || source.type === 5 ? 'var(--primary-color)' : '#9c27b0'}">
                         {source.type_name || "Protocolo Estándar"}
                      </span>
                    </td>
                    <td>
                      <span class="status-chip" class:active={source.is_actived}>
                        {source.is_actived ? 'Operativo' : 'Inactivo'}
                      </span>
                    </td>
                    <td>
                      <div class="action-buttons">
                        <button class="icon-btn" onclick={() => editSource(source)} title="Configurar">
                          <Settings size={18} />
                        </button>
                        <button class="icon-btn danger" onclick={() => deleteAuth(source.id)} title="Eliminar">
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
      </div>
    </Card>
  </div>
{/if}

<style>
  .admin-auths-body {
    padding: 0 1.5rem 4rem;
    animation: fadeIn 0.4s cubic-bezier(0.4, 0, 0.2, 1);
  }

  @keyframes fadeIn {
    from { opacity: 0; transform: translateY(4px); }
    to { opacity: 1; transform: translateY(0); }
  }
  
  .id-cell { font-family: 'JetBrains Mono', monospace; color: var(--text-muted); font-weight: 700; opacity: 0.5; font-size: 0.85rem; }



  .id-cell { font-family: 'JetBrains Mono', monospace; color: var(--text-muted); font-weight: 700; opacity: 0.5; font-size: 0.85rem; }
  
  .source-cell { display: flex; align-items: center; gap: 1.25rem; }
  .source-icon { width: 44px; height: 44px; display: flex; align-items: center; justify-content: center; border-radius: 0.75rem; color: var(--primary-color); background: rgba(var(--primary-rgb), 0.1); }
  
  .source-info { display: flex; flex-direction: column; }
  .source-info strong { font-size: 1.1rem; font-weight: 800; color: var(--text-color); }
  .source-info span { font-size: 0.8rem; color: var(--text-muted); font-weight: 700; }

  .type-badge { font-size: 0.75rem; font-weight: 900; color: var(--accent); background: rgba(156, 39, 176, 0.1); padding: 0.4rem 1rem; border-radius: 0.5rem; border: 1px solid rgba(156, 39, 176, 0.2); text-transform: uppercase; letter-spacing: 0.05em; }
  
  .status-chip { display: inline-block; padding: 0.3rem 1rem; border-radius: 0.5rem; font-size: 0.7rem; font-weight: 900; background: rgba(255,255,255,0.05); color: var(--text-muted); text-transform: uppercase; border: 1px solid var(--md-sys-color-outline-variant); }
  .status-chip.active { color: #10b981; background: rgba(16, 185, 129, 0.1); border-color: rgba(16, 185, 129, 0.3); }

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

  .empty-state { padding: 4rem; text-align: center; color: var(--text-muted); font-weight: 700; font-size: 1.1rem; opacity: 0.5; }
</style>
