<script lang="ts">
  import { apiFetch } from "$lib/api";
  import { i18n } from '$lib/i18n.svelte';
  import { onMount } from 'svelte';

  let auditLogs = $state<any[]>([]);
  let loading = $state(true);
  let error = $state('');
  let page = $state(1);
  let hasMore = $state(true);

  async function fetchAuditLogs(pageNum: number = 1) {
    loading = true;
    error = '';
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
      } else {
        error = 'Failed to load audit logs.';
      }
    } catch (e) {
      error = 'An error occurred while fetching audit logs.';
    } finally {
      loading = false;
    }
  }

  onMount(() => {
    fetchAuditLogs();
  });

  function formatDate(unix: number) {
    return new Date(unix * 1000).toLocaleString();
  }

  function getOpColor(opType: number) {
    // Basic color coding for operations
    if (opType === 1 || opType === 19) return '#10b981'; // Green for creation
    if (opType === 17 || opType === 18) return '#ef4444'; // Red for deletion
    if (opType === 5) return '#3b82f6'; // Blue for Push
    return 'var(--text-muted)';
  }
</script>

<div class="admin-view-header">
  <div class="header-content">
    <h1>{i18n.t('admin.audit_log') || 'Audit Log'}</h1>
    <p class="subtitle">Tracking system activity and administrative changes.</p>
  </div>
  <button class="btn btn-tonal" onclick={() => fetchAuditLogs(1)} disabled={loading}>
    <span class="material-symbols-rounded" class:spin={loading}>refresh</span>
    Refresh
  </button>
</div>

{#if error}
  <div class="error-card admin-card">
    <span class="material-symbols-rounded">warning</span>
    <p>{error}</p>
  </div>
{/if}

<div class="audit-table-container admin-card shadow-sm">
  <table class="audit-table">
    <thead>
      <tr>
        <th>Time</th>
        <th>Action</th>
        <th>Actor</th>
        <th>Repository</th>
        <th>Details</th>
      </tr>
    </thead>
    <tbody>
      {#each auditLogs as log}
        <tr>
          <td class="time-cell">{formatDate(log.created_at)}</td>
          <td>
            <span class="op-badge" style="--op-color: {getOpColor(log.op_type)}">
              {log.op_name}
            </span>
          </td>
          <td class="actor-cell">
            <span class="material-symbols-rounded">person</span>
            {log.actor_name}
          </td>
          <td class="repo-cell">
            <a href="/{log.repo_path}">{log.repo_path}</a>
          </td>
          <td class="content-cell">
            {#if log.content}
                <code class="mini-code">{log.content.slice(0, 50)}{log.content.length > 50 ? '...' : ''}</code>
            {:else}
                <span class="muted">--</span>
            {/if}
          </td>
        </tr>
      {/each}
      {#if auditLogs.length === 0 && !loading}
        <tr>
          <td colspan="5" class="empty-row">No audit logs found.</td>
        </tr>
      {/if}
    </tbody>
  </table>
  
  {#if hasMore}
    <div class="table-actions">
        <button class="btn btn-ghost load-more" onclick={() => fetchAuditLogs(page + 1)} disabled={loading}>
            {loading ? 'Loading...' : 'Load More Results'}
        </button>
    </div>
  {/if}
</div>

<style>
  .admin-view-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 2rem;
  }
  .header-content h1 { margin: 0; font-size: 2rem; font-weight: 800; letter-spacing: -0.02em; }
  .subtitle { color: var(--text-muted); margin: 0.25rem 0 0; }

  .audit-table-container {
    background: var(--md3-surface-container);
    border-radius: 1.25rem;
    overflow: hidden;
    border: 1px solid var(--border-color);
  }

  .audit-table {
    width: 100%;
    border-collapse: collapse;
    font-size: 0.9rem;
  }

  .audit-table th {
    text-align: left;
    padding: 1rem 1.5rem;
    background: rgba(255, 255, 255, 0.02);
    color: var(--text-muted);
    font-size: 0.75rem;
    text-transform: uppercase;
    font-weight: 800;
    letter-spacing: 0.05em;
    border-bottom: 1px solid var(--border-color);
  }

  .audit-table td {
    padding: 1rem 1.5rem;
    border-bottom: 1px solid rgba(255, 255, 255, 0.03);
    vertical-align: middle;
  }

  .time-cell { font-family: monospace; color: var(--text-muted); font-size: 0.8rem; }
  
  .op-badge {
    padding: 0.25rem 0.6rem;
    border-radius: 2rem;
    font-size: 0.75rem;
    font-weight: 700;
    background: rgba(var(--op-color, 255, 255, 255), 0.1);
    color: var(--op-color);
    border: 1px solid rgba(var(--op-color, 255, 255, 255), 0.2);
    white-space: nowrap;
  }

  .actor-cell { display: flex; align-items: center; gap: 0.5rem; font-weight: 600; }
  .actor-cell .material-symbols-rounded { font-size: 1.2rem; opacity: 0.5; }

  .repo-cell a { color: var(--primary-color); text-decoration: none; font-weight: 600; }
  .repo-cell a:hover { text-decoration: underline; }

  .content-cell .mini-code {
    background: rgba(0,0,0,0.3);
    padding: 0.2rem 0.4rem;
    border-radius: 0.4rem;
    font-size: 0.8rem;
    opacity: 0.8;
  }

  .empty-row { padding: 4rem; text-align: center; color: var(--text-muted); font-style: italic; }

  .table-actions {
    padding: 1rem;
    display: flex;
    justify-content: center;
    background: rgba(255, 255, 255, 0.01);
  }

  .load-more { width: 100%; max-width: 400px; }

  .spin { animation: spin 1s linear infinite; }
  @keyframes spin { from { transform: rotate(0deg); } to { transform: rotate(360deg); } }

  .btn-tonal {
    background: var(--md3-secondary-container);
    color: var(--md3-on-secondary-container);
  }

  .admin-card { padding: 2rem; background: var(--md3-surface-container); border-radius: 1.5rem; }
  .error-card { display: flex; align-items: center; gap: 1rem; color: #f87171; border: 1px solid rgba(248, 113, 113, 0.2); margin-bottom: 2rem; }
</style>
