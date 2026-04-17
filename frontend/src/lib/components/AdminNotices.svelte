<script lang="ts">
  import { apiFetch } from "$lib/api";
  import { onMount } from 'svelte';
  
  let notices = $state<any[]>([]);
  let total = $state(0);
  let loading = $state(true);
  let page = $state(1);

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

  onMount(fetchNotices);
</script>

<div class="admin-notices island tonal">
  <header class="section-header">
    <div class="title">
      <span class="material-symbols-rounded">notifications_active</span>
      <h2>System Notices</h2>
    </div>
    <div class="actions">
       <button class="btn btn-outline" onclick={fetchNotices}>Refresh</button>
       <button class="btn btn-danger" onclick={clearAll}>Clear All</button>
    </div>
  </header>

  <div class="table-container">
    <table class="modern-table">
      <thead>
        <tr>
          <th>ID</th>
          <th>Type</th>
          <th>Description</th>
          <th>Date</th>
          <th>Actions</th>
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
           <tr><td colspan="5" class="empty">No system notices found.</td></tr>
        {:else}
          {#each notices as notice}
            <tr>
              <td>{notice.id}</td>
              <td><span class="type-badge">{notice.type_name}</span></td>
              <td class="desc-cell">{notice.description}</td>
              <td>{new Date(notice.created_at).toLocaleDateString()}</td>
              <td class="actions">
                <button class="icon-btn danger" onclick={() => deleteNotice(notice.id)} title="Delete"><span class="material-symbols-rounded">delete</span></button>
              </td>
            </tr>
          {/each}
        {/if}
      </tbody>
    </table>
  </div>
</div>

<style>
  .admin-notices { padding: 2rem; display: flex; flex-direction: column; gap: 2rem; box-shadow: var(--shadow-xl); }
  
  .section-header { display: flex; justify-content: space-between; align-items: center; }
  .title { display: flex; align-items: center; gap: 0.75rem; }
  .title span { font-size: 1.75rem; color: var(--primary-color); }
  .title h2 { margin: 0; font-size: 1.5rem; font-weight: 800; }

  .actions { display: flex; gap: 1rem; }

  .table-container { overflow-x: auto; }
  .modern-table { width: 100%; border-collapse: collapse; text-align: left; table-layout: fixed; }
  .modern-table th { padding: 1rem; border-bottom: 1px solid var(--border-color); color: var(--text-muted); font-size: 0.75rem; text-transform: uppercase; letter-spacing: 0.1em; font-weight: 800; }
  .modern-table td { padding: 1.25rem 1rem; border-bottom: 1px solid rgba(255,255,255,0.03); font-size: 0.95rem; vertical-align: middle; }

  th:nth-child(1) { width: 50px; }
  th:nth-child(2) { width: 120px; }
  th:nth-child(3) { width: 1fr; }
  th:nth-child(4) { width: 150px; }
  th:nth-child(5) { width: 80px; }

  .type-badge { background: rgba(var(--primary-rgb), 0.1); color: var(--primary-color); padding: 0.2rem 0.6rem; border-radius: 4px; font-weight: 800; font-size: 0.7rem; text-transform: uppercase; }
  
  .desc-cell { overflow: hidden; text-overflow: ellipsis; white-space: nowrap; font-family: var(--font-mono); font-size: 0.85rem; color: var(--text-muted); }

  .actions { display: flex; gap: 0.5rem; }
  .icon-btn { background: transparent; border: none; color: var(--text-muted); cursor: pointer; padding: 0.4rem; border-radius: 0.5rem; transition: all 0.2s; }
  .icon-btn:hover { background: rgba(255,255,255,0.05); color: #fff; }
  .icon-btn.danger:hover { background: rgba(239, 68, 68, 0.1); color: #ef4444; }

  .empty { text-align: center; padding: 3rem; color: var(--text-muted); }
  .skeleton-bar { height: 1.5rem; background: rgba(255,255,255,0.03); border-radius: 0.5rem; animation: pulse 1.5s infinite; }
  @keyframes pulse { 0% { opacity: 0.3; } 50% { opacity: 0.6; } 100% { opacity: 0.3; } }
</style>
