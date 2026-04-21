<script lang="ts">
  import { apiFetch } from "$lib/api";
  import { onMount } from 'svelte';
  import { fade, fly } from 'svelte/transition';
  import Modal from './Modal.svelte';
  
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
    } catch (e) {
      console.error(e);
    } finally {
      loading = false;
    }
  }

  async function deleteNotice(id: number) {
    if (!confirm('Delete this notice?')) return;
    try {
      const res = await apiFetch(`/api/v2/admin/notices/${id}`, { method: 'DELETE' });
      if (res.ok) fetchNotices();
    } catch (e) {
      console.error(e);
    }
  }

  async function clearAll() {
    if (!confirm('Are you sure you want to clear ALL notices?')) return;
    try {
      const res = await apiFetch(`/api/v2/admin/notices`, { method: 'DELETE' });
      if (res.ok) fetchNotices();
    } catch (e) {
      console.error(e);
    }
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

<div class="notices-view" in:fly={{ y: 20, duration: 500 }}>
  <header class="admin-header">
    <div class="title-group">
      <h1>System Notices</h1>
      <p>System alerts, errors, and background task logs.</p>
    </div>
    <div class="header-actions">
       <button class="btn btn-outline" onclick={fetchNotices}>
         <span class="material-symbols-rounded">refresh</span>
         Refresh
       </button>
       <button class="btn btn-danger" onclick={clearAll}>
         <span class="material-symbols-rounded">delete_sweep</span>
         Clear All
       </button>
    </div>
  </header>

  <div class="admin-notices island tonal">
    <div class="table-container">
      <table class="modern-table">
        <thead>
          <tr>
            <th class="col-id">ID</th>
            <th class="col-type">Type</th>
            <th class="col-desc">Description</th>
            <th class="col-date">Date</th>
            <th class="col-actions">Actions</th>
          </tr>
        </thead>
        <tbody>
          {#if loading}
            {#each Array(5) as _}
              <tr class="skeleton-row">
                <td colspan="5"><div class="skeleton-bar"></div></td>
              </tr>
            {/each}
          {:else if notices.length === 0}
             <tr><td colspan="5" class="empty-state">No system notices found.</td></tr>
          {:else}
            {#each notices as notice}
              <tr class="notice-row" onclick={() => openDetail(notice)}>
                <td><span class="id-tag">{notice.id}</span></td>
                <td><span class="type-badge tier-{notice.type}">{notice.type_name}</span></td>
                <td class="desc-cell">
                  <div class="desc-content">{notice.description}</div>
                </td>
                <td class="date-cell">{new Date(notice.created_at).toLocaleString()}</td>
                <td class="actions-cell" onclick={(e) => e.stopPropagation()}>
                  <button class="icon-btn danger" onclick={() => deleteNotice(notice.id)} title="Delete">
                    <span class="material-symbols-rounded">delete</span>
                  </button>
                </td>
              </tr>
            {/each}
          {/if}
        </tbody>
      </table>
    </div>

    <!-- PAGINATION -->
    {#if total > pageSize}
      <footer class="pagination-footer">
        <span class="page-info">Showing {notices.length} of {total} notices</span>
        <div class="pagination-controls">
          <button 
            class="btn-nav" 
            disabled={page === 1} 
            onclick={() => page--}
          >
            <span class="material-symbols-rounded">chevron_left</span>
          </button>
          <span class="current-page">Page {page} of {totalPages}</span>
          <button 
            class="btn-nav" 
            disabled={page >= totalPages} 
            onclick={() => page++}
          >
            <span class="material-symbols-rounded">chevron_right</span>
          </button>
        </div>
      </footer>
    {/if}
  </div>
</div>

<!-- DETAIL MODAL -->
<Modal show={showModal} title="Notice Details" onclose={() => showModal = false}>
  {#if selectedNotice}
    <div class="notice-detail">
      <div class="detail-header">
        <span class="type-badge tier-{selectedNotice.type} big">{selectedNotice.type_name}</span>
        <span class="detail-date">{new Date(selectedNotice.created_at).toLocaleString()}</span>
      </div>
      <div class="detail-body">
        <h4>Content</h4>
        <pre class="detail-text">{selectedNotice.description}</pre>
      </div>
    </div>
  {/if}
</Modal>

<style>
  .notices-view { display: flex; flex-direction: column; gap: 2.5rem; }

  .admin-header { display: flex; justify-content: space-between; align-items: flex-end; }
  .admin-header h1 { margin: 0; font-size: 2.5rem; font-weight: 900; }
  .admin-header p { margin: 0.5rem 0 0 0; color: var(--text-muted); font-size: 1.1rem; }
  .header-actions { display: flex; gap: 1rem; }

  .admin-notices { padding: 0; overflow: hidden; border: 1px solid rgba(255,255,255,0.05); }
  
  .table-container { width: 100%; overflow-x: auto; }
  .modern-table { width: 100%; border-collapse: collapse; text-align: left; }
  .modern-table th { background: rgba(255,255,255,0.01); padding: 1.25rem 1.5rem; font-size: 0.75rem; font-weight: 800; text-transform: uppercase; letter-spacing: 0.1em; color: var(--text-muted); border-bottom: 1px solid var(--md-sys-color-outline-variant); }
  .modern-table td { padding: 1.25rem 1.5rem; border-bottom: 1px solid rgba(255,255,255,0.02); vertical-align: middle; }

  .notice-row { cursor: pointer; transition: background 0.2s; }
  .notice-row:hover { background: rgba(255,255,255,0.03); }

  .id-tag { color: var(--text-muted); font-family: var(--font-mono); font-size: 0.85rem; }
  
  .type-badge { padding: 0.25rem 0.75rem; border-radius: 6px; font-size: 0.7rem; font-weight: 800; text-transform: uppercase; }
  .type-badge.tier-1 { background: rgba(59, 130, 246, 0.1); color: #3b82f6; } /* info */
  .type-badge.tier-2 { background: rgba(245, 158, 11, 0.1); color: #f59e0b; } /* warning */
  .type-badge.tier-3 { background: rgba(239, 68, 68, 0.1); color: #ef4444; } /* error */
  .type-badge.big { font-size: 0.9rem; padding: 0.5rem 1rem; }

  .desc-cell { max-width: 400px; }
  .desc-content { overflow: hidden; text-overflow: ellipsis; white-space: nowrap; font-size: 0.9rem; color: #ddd; font-family: var(--font-mono); }
  
  .date-cell { white-space: nowrap; font-size: 0.85rem; color: var(--text-muted); }

  .icon-btn { background: transparent; border: none; color: var(--text-muted); padding: 0.5rem; border-radius: 50%; cursor: pointer; transition: all 0.2s; }
  .icon-btn:hover { background: rgba(255,255,255,0.05); color: #fff; }
  .icon-btn.danger:hover { color: #ef4444; background: rgba(239, 68, 68, 0.1); }

  .pagination-footer { display: flex; justify-content: space-between; align-items: center; padding: 1.25rem 1.5rem; background: rgba(0,0,0,0.2); }
  .page-info { font-size: 0.85rem; color: var(--text-muted); }
  .pagination-controls { display: flex; align-items: center; gap: 1rem; }
  .btn-nav { background: var(--surface-variant); border: 1px solid var(--md-sys-color-outline-variant); color: white; padding: 0.4rem; border-radius: 0.5rem; cursor: pointer; transition: all 0.2s; }
  .btn-nav:hover:not(:disabled) { transform: scale(1.1); background: var(--primary-color); }
  .btn-nav:disabled { opacity: 0.3; cursor: not-allowed; }
  .current-page { font-size: 0.9rem; font-weight: 600; }

  /* DETAIL MODAL STYLES */
  .notice-detail { display: flex; flex-direction: column; gap: 2rem; }
  .detail-header { display: flex; justify-content: space-between; align-items: center; }
  .detail-date { font-size: 0.9rem; color: var(--text-muted); }
  .detail-body h4 { margin: 0 0 1rem 0; font-size: 0.8rem; font-weight: 800; text-transform: uppercase; color: var(--text-muted); }
  .detail-text { background: rgba(0,0,0,0.3); padding: 1.5rem; border-radius: 12px; font-family: 'JetBrains Mono', monospace; font-size: 0.9rem; white-space: pre-wrap; word-break: break-all; border: 1px solid var(--md-sys-color-outline-variant); margin: 0; color: #6ee7b7; }

  .empty-state { padding: 4rem; text-align: center; color: var(--text-muted); font-style: italic; }
</style>
