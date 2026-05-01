<script lang="ts">
  import { apiFetch } from "$lib/api";
  import { onMount } from "svelte";
  import { i18n } from '$lib/i18n.svelte';
  import { dialog } from '$lib/dialog.svelte';
  import Card from './ui/Card.svelte';
  import Button from './ui/Button.svelte';
  import PageHeader from './ui/PageHeader.svelte';
  import Modal from "./Modal.svelte";

  let notices = $state<any[]>([]);
  let total = $state(0);
  let loading = $state(true);
  let page = $state(1);
  const pageSize = 50;

  let selectedNotice = $state<any>(null);
  let showModal = $state(false);

  async function fetchNotices() {
    loading = true;
    try {
      const res = await apiFetch(`/api/v2/admin/notices?page=${page}`);
      if (res.ok) {
        const data = await res.json();
        notices = data.notices || [];
        total = data.total || 0;
      }
    } catch (e) { console.error(e); } finally { loading = false; }
  }

  async function deleteNotice(id: number) {
    if (!await dialog.confirm("¿Eliminar este aviso?")) return;
    try {
      const res = await apiFetch(`/api/v2/admin/notices/${id}`, { method: "DELETE" });
      if (res.ok) fetchNotices();
    } catch (e) { console.error(e); }
  }

  async function clearAll() {
    if (!await dialog.confirm("¿Seguro que quieres borrar TODOS los avisos?")) return;
    try {
      const res = await apiFetch(`/api/v2/admin/notices`, { method: "DELETE" });
      if (res.ok) fetchNotices();
    } catch (e) { console.error(e); }
  }

  function openDetail(notice: any) {
    selectedNotice = notice;
    showModal = true;
  }

  const totalPages = $derived(Math.ceil(total / pageSize));

  onMount(fetchNotices);

  $effect(() => {
    if (page) fetchNotices();
  });
</script>

<div class="admin-notices-container">
  <PageHeader 
    title={i18n.t('admin.notices') || "Avisos del Sistema"}
    subtitle="Registros de errores, alertas y tareas en segundo plano"
    icon="notifications"
  >
    {#snippet actions()}
      <div class="header-actions">
        <Button variant="ghost" icon="refresh" onclick={fetchNotices}>Actualizar</Button>
        <Button variant="danger" icon="delete_sweep" onclick={clearAll}>Limpiar Todo</Button>
      </div>
    {/snippet}
  </PageHeader>

  <div class="admin-notices-body">
  {#if loading}
    <div class="skeleton-monitor"></div>
  {:else}
    <div class="modern-table-card">
      <table class="modern-table">
        <thead>
          <tr>
            <th class="col-id">ID</th>
            <th class="col-type">Tipo</th>
            <th class="col-desc">Descripción</th>
            <th class="col-date">Fecha</th>
            <th class="col-actions">Acciones</th>
          </tr>
        </thead>
        <tbody>
          {#if notices.length === 0}
            <tr><td colspan="5" class="empty-state">No se han encontrado avisos del sistema.</td></tr>
          {:else}
            {#each notices as notice}
              <tr class="hover-row clickable" onclick={() => openDetail(notice)}>
                <td class="id-cell">#{notice.id}</td>
                <td>
                  <span class="type-badge tier-{notice.type}">{notice.type_name}</span>
                </td>
                <td>
                  <div class="desc-preview">{notice.description}</div>
                </td>
                <td>
                  <span class="date-tag">{new Date(notice.created_at).toLocaleString()}</span>
                </td>
                <td>
                  <Button variant="ghost" size="sm" icon="delete" class="danger-btn" onclick={(e) => { e.stopPropagation(); deleteNotice(notice.id); }} />
                </td>
              </tr>
            {/each}
          {/if}
        </tbody>
      </table>
    </div>

    {#if total > pageSize}
      <footer class="pagination-footer">
        <span class="count-info">Mostrando {notices.length} de {total} avisos</span>
        <div class="nav-controls">
          <Button variant="ghost" size="sm" icon="chevron_left" disabled={page === 1} onclick={() => page--} />
          <span class="page-current">Página {page} de {totalPages}</span>
          <Button variant="ghost" size="sm" icon="chevron_right" disabled={page >= totalPages} onclick={() => page++} />
        </div>
      </footer>
    {/if}
    </div>
  </div>

<Modal show={showModal} title="Detalle del Aviso" onclose={() => showModal = false}>
  {#if selectedNotice}
    <div class="premium-detail">
      <div class="detail-top">
        <span class="type-badge tier-{selectedNotice.type} big">{selectedNotice.type_name}</span>
        <span class="detail-time">Generado el {new Date(selectedNotice.created_at).toLocaleString()}</span>
      </div>
      <div class="detail-content">
        <span class="detail-label">Contenido del Log</span>
        <div class="code-box">
          {selectedNotice.description}
        </div>
      </div>
    </div>
  {/if}
</Modal>

<style>
  .admin-notices-body {
    padding: 0 1.5rem 4rem;
    animation: fadeIn 0.4s cubic-bezier(0.4, 0, 0.2, 1);
  }

  @keyframes fadeIn {
    from { opacity: 0; transform: translateY(4px); }
    to { opacity: 1; transform: translateY(0); }
  }
  .section-header { margin-bottom: 4rem; }
  .header-flex { display: flex; justify-content: space-between; align-items: center; gap: 2rem; }
  .title-with-icon { display: flex; align-items: center; gap: 2rem; }
  .hero-icon { font-size: 3rem; color: var(--primary-color); background: rgba(var(--primary-rgb), 0.1); padding: 1.25rem; border-radius: 2rem; }
  h2 { margin: 0; font-size: 2.25rem; font-weight: 950; letter-spacing: -0.04em; }
  .subtitle { margin: 0.5rem 0 0 0; color: var(--text-muted); font-size: 1.1rem; }
  .header-actions { display: flex; gap: 0.75rem; }


  .clickable { cursor: pointer; transition: background 0.3s cubic-bezier(0.4, 0, 0.2, 1); }
  .clickable:hover { background: rgba(255,255,255,0.03); }

  .id-cell { font-family: 'JetBrains Mono', monospace; font-weight: 700; color: var(--text-muted); opacity: 0.5; }

  .type-badge { padding: 0.4rem 0.8rem; border-radius: 0.75rem; font-size: 0.7rem; font-weight: 900; text-transform: uppercase; background: #333; }
  .type-badge.tier-1 { background: rgba(59, 130, 246, 0.12); color: #3b82f6; }
  .type-badge.tier-2 { background: rgba(245, 158, 11, 0.12); color: #f59e0b; }
  .type-badge.tier-3 { background: rgba(239, 68, 68, 0.12); color: #ef4444; }
  .type-badge.big { font-size: 0.9rem; padding: 0.5rem 1.25rem; border-radius: 1rem; }

  .desc-preview { font-family: 'JetBrains Mono', monospace; font-size: 0.85rem; color: #ddd; max-width: 500px; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; opacity: 0.8; }
  .date-tag { font-size: 0.85rem; color: var(--text-muted); font-weight: 600; }

  .pagination-footer { margin-top: 2rem; display: flex; justify-content: space-between; align-items: center; padding: 1.5rem; background: rgba(0,0,0,0.1); border-radius: 2rem; border: 1px solid var(--md-sys-color-outline-variant); }
  .count-info { font-size: 0.85rem; color: var(--text-muted); font-weight: 600; }
  .nav-controls { display: flex; align-items: center; gap: 1.5rem; }
  .page-current { font-weight: 800; font-size: 0.9rem; }

  .premium-detail { display: flex; flex-direction: column; gap: 2rem; padding: 0.5rem; }
  .detail-top { display: flex; justify-content: space-between; align-items: center; }
  .detail-time { color: var(--text-muted); font-size: 0.9rem; font-weight: 600; }
  .detail-content { display: flex; flex-direction: column; gap: 0.75rem; }
  .detail-content .detail-label { font-size: 0.75rem; text-transform: uppercase; font-weight: 900; color: var(--primary-color); letter-spacing: 0.1em; }
  .code-box { background: #08080c; border: 1px solid #1a1a24; border-radius: 1.25rem; padding: 1.5rem; font-family: 'JetBrains Mono', monospace; line-height: 1.6; font-size: 0.95rem; color: #6ee7b7; white-space: pre-wrap; word-break: break-all; }

  .skeleton-monitor { height: 600px; background: var(--surface-variant); border-radius: 2.5rem; animation: pulse 1.5s infinite; }
</style>
