<script lang="ts">
  import { apiFetch } from "$lib/api";
  import { onMount } from "svelte";
  import { i18n } from '$lib/i18n.svelte';
  import Card from './ui/Card.svelte';
  import Button from './ui/Button.svelte';

  let sources = $state<any[]>([]);
  let loading = $state(true);

  async function fetchAuths() {
    loading = true;
    try {
      const res = await apiFetch("/api/v2/admin/auths");
      if (res.ok) sources = await res.json();
    } catch (e) { console.error(e); } finally { loading = false; }
  }

  async function deleteAuth(id: number) {
    if (!confirm("¿Seguro que quieres eliminar este origen de autenticación?")) return;
    try {
      const res = await apiFetch(`/api/v2/admin/auths/${id}`, { method: "DELETE" });
      if (res.ok) fetchAuths();
    } catch (e) { console.error(e); }
  }

  onMount(fetchAuths);
</script>

<Card class="settings-premium-card" padding="lg">
  <header class="section-header">
    <div class="header-flex">
      <div class="title-with-icon">
        <span class="material-symbols-rounded hero-icon">lock_person</span>
        <div>
          <h2>Fuentes de Autenticación</h2>
          <p class="subtitle">Gestión centralizada de protocolos LDAP, OAuth2, PAM y servicios externos</p>
        </div>
      </div>
      <Button variant="primary" icon="add_moderator">Nueva Fuente</Button>
    </div>
  </header>

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
          {#each Array(3) as _}
            <tr class="skeleton-row">
              <td colspan="5"><div class="skeleton-line"></div></td>
            </tr>
          {/each}
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
                    <span class="material-symbols-rounded">{source.type === 2 || source.type === 5 ? "hub" : "verified_user"}</span>
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
                  <Button variant="ghost" size="sm" icon="settings" />
                  <Button variant="ghost" size="sm" icon="delete" class="danger-btn" onclick={() => deleteAuth(source.id)} />
                </div>
              </td>
            </tr>
          {/each}
        {/if}
      </tbody>
    </table>
  </div>
</Card>

<style>
  
  .section-header { margin-bottom: 4rem; }
  .header-flex { display: flex; justify-content: space-between; align-items: center; gap: 2rem; }
  .title-with-icon { display: flex; align-items: center; gap: 2rem; }
  .hero-icon { font-size: 3rem; color: var(--primary-color); background: rgba(var(--primary-rgb), 0.1); padding: 1.25rem; border-radius: 2rem; }
  h2 { margin: 0; font-size: 2.25rem; font-weight: 950; letter-spacing: -0.04em; }
  .subtitle { margin: 0.5rem 0 0 0; color: var(--text-muted); font-size: 1.1rem; }



  .id-cell { font-family: 'JetBrains Mono', monospace; color: var(--text-muted); font-weight: 700; opacity: 0.5; font-size: 0.85rem; }
  
  .source-cell { display: flex; align-items: center; gap: 1.25rem; }
  .source-icon { width: 44px; height: 44px; display: flex; align-items: center; justify-content: center; border-radius: 1rem; color: var(--primary-color); background: rgba(var(--primary-rgb), 0.1); }
  .source-icon span { font-size: 1.75rem; }
  
  .source-info { display: flex; flex-direction: column; }
  .source-info strong { font-size: 1.1rem; font-weight: 800; color: var(--text-color); }
  .source-info span { font-size: 0.8rem; color: var(--text-muted); font-weight: 700; }

  .type-badge { font-size: 0.75rem; font-weight: 900; color: var(--accent); background: rgba(156, 39, 176, 0.1); padding: 0.4rem 1rem; border-radius: 100px; border: 1px solid rgba(156, 39, 176, 0.2); text-transform: uppercase; letter-spacing: 0.05em; }
  
  .status-chip { display: inline-block; padding: 0.3rem 1rem; border-radius: 100px; font-size: 0.7rem; font-weight: 900; background: rgba(255,255,255,0.05); color: var(--text-muted); text-transform: uppercase; border: 1px solid var(--md-sys-color-outline-variant); }
  .status-chip.active { color: #10b981; background: rgba(16, 185, 129, 0.1); border-color: rgba(16, 185, 129, 0.3); }

  .action-buttons { display: flex; gap: 0.5rem; }
  :global(.danger-btn:hover) { color: var(--error) !important; background: rgba(var(--error-rgb), 0.1) !important; }

  .empty-state { padding: 4rem; text-align: center; color: var(--text-muted); font-weight: 700; font-size: 1.1rem; opacity: 0.5; }
  .skeleton-line { height: 2rem; background: var(--surface-variant); border-radius: 0.75rem; animation: pulse 1.5s infinite; }
</style>
