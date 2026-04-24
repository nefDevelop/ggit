<script lang="ts">
  import { apiFetch } from "$lib/api";
  import { i18n } from '$lib/i18n.svelte';
  import { onMount } from 'svelte';
  import Card from './ui/Card.svelte';
  import Button from './ui/Button.svelte';

  let auditLogs = $state<any[]>([]);
  let loading = $state(true);
  let page = $state(1);
  let hasMore = $state(true);

  async function fetchAuditLogs(pageNum: number = 1) {
    loading = true;
    try {
      const res = await apiFetch(`/api/v2/admin/audit-log?page=${pageNum}`);
      if (res.ok) {
        const data = await res.json();
        if (pageNum === 1) {
          auditLogs = data;
        } else {
          auditLogs = [...auditLogs, ...data];
        }
        hasMore = data.length > 0;
        page = pageNum;
      }
    } catch (e) {
      console.error(e);
    } finally {
      loading = false;
    }
  }

  onMount(() => { fetchAuditLogs(); });

  function formatDate(unix: number) {
    return new Date(unix * 1000).toLocaleString();
  }

  function getOpColor(opType: number) {
    if (opType === 1 || opType === 19) return '#10b981'; 
    if (opType === 17 || opType === 18) return '#ef4444';
    if (opType === 5) return '#3b82f6';
    if (opType === 20) return '#f59e0b';
    return '#888';
  }
</script>

<Card class="settings-premium-card" padding="lg">
  <header class="section-header">
    <div class="header-flex">
      <div class="title-with-icon">
        <span class="material-symbols-rounded hero-icon">history_edu</span>
        <div>
          <h2>Registro de Auditoría</h2>
          <p class="subtitle">Historial de eventos críticos, cambios de seguridad y acciones administrativas</p>
        </div>
      </div>
      <Button variant="ghost" icon="refresh" onclick={() => fetchAuditLogs(1)} disabled={loading}>Actualizar</Button>
    </div>
  </header>

  <div class="modern-table-card">
    <table class="modern-table">
      <thead>
        <tr>
          <th>Tiempo del Evento</th>
          <th>Acción Realizada</th>
          <th>Usuario Responsable</th>
          <th>Entidad / Repositorio</th>
          <th>Detalles Técnicos</th>
        </tr>
      </thead>
      <tbody>
        {#if loading && auditLogs.length === 0}
          {#each Array(5) as _}
            <tr class="skeleton-row">
              <td colspan="5"><div class="skeleton-line"></div></td>
            </tr>
          {/each}
        {:else if auditLogs.length === 0}
          <tr><td colspan="5" class="empty-state">No se han encontrado registros de auditoría históricos.</td></tr>
        {:else}
          {#each auditLogs as log}
            <tr class="hover-row">
              <td class="time-cell" data-label="Tiempo del Evento">{formatDate(log.created_at)}</td>
              <td data-label="Acción Realizada">
                <span class="op-capsule" style="--op-color: {getOpColor(log.op_type)}">
                  {log.op_name}
                </span>
              </td>
              <td data-label="Usuario Responsable">
                <div class="actor-brand">
                  <span class="material-symbols-rounded">person</span>
                  <strong>{log.actor_name}</strong>
                </div>
              </td>
              <td class="repo-cell" data-label="Entidad / Repositorio">
                {#if log.repo_path}
                   <a href="/{log.repo_path}" class="repo-link">{log.repo_path}</a>
                {:else}
                   <span class="muted-value">—</span>
                {/if}
              </td>
              <td data-label="Detalles Técnicos">
                {#if log.content}
                   <code class="mini-snippet">{log.content.slice(0, 60)}{log.content.length > 60 ? '...' : ''}</code>
                {:else}
                   <span class="null-value">No disponible</span>
                {/if}
              </td>
            </tr>
          {/each}
        {/if}
      </tbody>
    </table>
  </div>

  {#if hasMore}
    <div class="footer-actions">
       <Button variant="primary" icon="expand_more" class="load-more-btn" onclick={() => fetchAuditLogs(page + 1)} disabled={loading}>
          {loading ? 'Consultando...' : 'Ver registros anteriores'}
       </Button>
    </div>
  {/if}
</Card>

<style>
  :global(.settings-premium-card) { border-radius: 3rem; background: rgba(var(--surface-rgb), 0.5); backdrop-filter: blur(24px); border: 1px solid var(--md-sys-color-outline-variant); }
  
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
  }

  .modern-table-card { border: 1px solid var(--md-sys-color-outline-variant); border-radius: 2.5rem; overflow: hidden; background: rgba(0,0,0,0.15); box-shadow: inset 0 2px 4px rgba(0,0,0,0.1); }
  .modern-table { width: 100%; border-collapse: collapse; text-align: left; }
  .modern-table th { padding: 1.5rem; background: rgba(var(--primary-rgb), 0.03); color: var(--text-muted); font-size: 0.8rem; text-transform: uppercase; letter-spacing: 0.15em; font-weight: 900; border-bottom: 1px solid var(--md-sys-color-outline-variant); }
  .modern-table td { padding: 1.75rem 1.5rem; border-bottom: 1px solid var(--md-sys-color-outline-variant); vertical-align: middle; }
  .modern-table tr:last-child td { border-bottom: none; }

  @media (max-width: 1024px) {
    .modern-table-card { border-radius: 1.5rem; }
    .time-cell, .op-capsule, .actor-brand, .repo-link, .mini-snippet { text-align: right; }
    .actor-brand { justify-content: flex-end; }
  }

  .time-cell { font-family: 'JetBrains Mono', monospace; font-size: 0.85rem; color: var(--text-muted); font-weight: 700; opacity: 0.8; }
  
  .op-capsule { padding: 0.4rem 1rem; border-radius: 100px; font-size: 0.75rem; font-weight: 900; background: rgba(var(--op-color, 255, 255, 255), 0.1); color: var(--op-color, #fff); border: 1px solid rgba(var(--op-color, 255, 255, 255), 0.1); text-transform: uppercase; letter-spacing: 0.05em; }

  .actor-brand { display: flex; align-items: center; gap: 0.75rem; }
  .actor-brand span { color: var(--primary-color); opacity: 0.6; font-size: 1.25rem; }
  .actor-brand strong { color: var(--text-color); font-weight: 800; font-size: 1rem; }

  .repo-link { color: var(--primary-color); text-decoration: none; font-weight: 800; font-family: 'JetBrains Mono', monospace; font-size: 0.85rem; }
  .repo-link:hover { text-decoration: underline; }

  .mini-snippet { background: rgba(0,0,0,0.25); padding: 0.4rem 0.75rem; border-radius: 0.85rem; font-size: 0.85rem; color: #999; font-family: 'JetBrains Mono', monospace; border: 1px solid rgba(255,255,255,0.05); }
  .muted-value { color: #555; font-weight: 800; }
  .null-value { color: #444; font-style: italic; font-size: 0.85rem; }

  .footer-actions { margin-top: 3rem; display: flex; justify-content: center; }
  :global(.load-more-btn) { width: auto; min-width: 300px; border-radius: 1.5rem !important; }

  .empty-state { padding: 5rem; text-align: center; color: var(--text-muted); font-weight: 700; opacity: 0.5; }
  .skeleton-line { height: 2rem; background: var(--surface-variant); border-radius: 0.75rem; animation: pulse 1.5s infinite; }
</style>
