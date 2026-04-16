<script lang="ts">
  import { onMount } from 'svelte';
  
  let notices = $state<any[]>([]);
  let total = $state(0);
  let loading = $state(true);
  let page = $state(1);

  async function fetchNotices() {
    loading = true;
    try {
      const res = await fetch(`/api/v2/admin/notices?page=${page}`);
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
      const res = await fetch(`/api/v2/admin/notices/${id}`, {
        method: 'DELETE'
      });
      if (res.ok) {
        notices = notices.filter(n => n.ID !== id);
        total--;
      }
    } catch (e) {
      console.error(e);
    }
  }

  async function clearAll() {
    if (!confirm('Clear all system notices?')) return;
    try {
      const res = await fetch('/api/v2/admin/notices', { method: 'DELETE' });
      if (res.ok) {
        notices = [];
        total = 0;
      }
    } catch (e) {
      console.error(e);
    }
  }

  onMount(fetchNotices);
</script>

<div class="admin-notices">
  <header class="section-header">
    <div class="header-text">
       <h2>System Notices</h2>
       <p class="subtitle">System events, maintenance logs, and error reports.</p>
    </div>
    <div class="header-actions">
       <span class="count-badge">{total} Notices</span>
       <button class="btn btn-danger-text" onclick={clearAll}>
         <span class="material-symbols-rounded">delete_sweep</span>
         Clear All
       </button>
    </div>
  </header>

  <div class="notices-list">
    {#each notices as notice}
      <div class="notice-item island tonal">
        <div class="notice-icon">
          <span class="material-symbols-rounded">error_outline</span>
        </div>
        <div class="notice-content">
           <div class="notice-meta">
             <span class="type">Notice #{notice.ID}</span>
             <span class="date">{new Date(notice.CreatedUnix * 1000).toLocaleString()}</span>
           </div>
           <p class="desc">{notice.Description}</p>
        </div>
        <button class="btn btn-icon btn-danger-text" onclick={() => deleteNotice(notice.ID)}>
          <span class="material-symbols-rounded">delete</span>
        </button>
      </div>
    {/each}

    {#if !loading && notices.length === 0}
      <div class="empty-state island">
        <span class="material-symbols-rounded">check_circle</span>
        <h3>System healthy</h3>
        <p>No new system notices or alerts.</p>
      </div>
    {/if}
  </div>
</div>

<style>
  .admin-notices { display: flex; flex-direction: column; gap: 2rem; }
  
  .section-header { display: flex; justify-content: space-between; align-items: flex-end; }
  .section-header h2 { margin: 0; font-size: 1.75rem; font-weight: 800; }
  .subtitle { color: var(--text-muted); margin: 0.25rem 0 0 0; }

  .header-actions { display: flex; align-items: center; gap: 1.5rem; }
  .count-badge { font-size: 0.85rem; font-weight: 800; color: var(--text-muted); background: rgba(255,255,255,0.05); padding: 0.25rem 0.75rem; border-radius: 2rem; }

  .notices-list { display: flex; flex-direction: column; gap: 1rem; }
  
  .notice-item { display: flex; align-items: flex-start; gap: 1.25rem; padding: 1.25rem; border-left: 4px solid #ef4444; }
  
  .notice-icon { font-size: 1.5rem; color: #ef4444; margin-top: 0.25rem; }
  
  .notice-content { flex: 1; display: flex; flex-direction: column; gap: 0.5rem; }
  .notice-meta { display: flex; gap: 1rem; font-size: 0.75rem; font-weight: 800; text-transform: uppercase; letter-spacing: 0.05em; color: var(--text-muted); }
  
  .desc { margin: 0; font-size: 0.95rem; line-height: 1.5; color: var(--text-color); font-family: monospace; word-break: break-all; }

  .btn-danger-text { color: #ef4444; display: flex; align-items: center; gap: 0.5rem; font-weight: 700; background: transparent; border: none; cursor: pointer; padding: 0.5rem 1rem; border-radius: 0.75rem; transition: background 0.2s; }
  .btn-danger-text:hover { background: rgba(239, 68, 68, 0.1); }

  .empty-state { text-align: center; padding: 4rem; color: #10b981; border: 2px dashed rgba(16, 185, 129, 0.2); background: rgba(16, 185, 129, 0.02); }
  .empty-state span { font-size: 4rem; margin-bottom: 1rem; opacity: 0.5; }
  .empty-state h3 { margin: 0; }
  .empty-state p { color: var(--text-muted); margin-top: 0.5rem; }
</style>
