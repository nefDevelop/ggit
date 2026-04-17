<script lang="ts">
  import { apiFetch } from "$lib/api";
  import { onMount, getContext } from 'svelte';

  let { hook, onclose } = $props();
  const repo: any = getContext('repo');

  let logs = $state<any[]>([]);
  let loading = $state(true);
  let error = $state('');
  
  let selectedLog = $state<any>(null);

  async function fetchLogs() {
    loading = true;
    try {
      const res = await apiFetch(`/api/v2/${repo.data.owner.username}/${repo.data.name}/settings/hooks/${hook.id}/logs`);
      if (res.ok) {
        logs = await res.json();
      } else {
        error = 'Failed to load delivery logs.';
      }
    } catch (e) {
      error = 'An error occurred while fetching logs.';
    } finally {
      loading = false;
    }
  }

  function formatDate(ts: number) {
    const d = new Date(ts / 1000); // Microseconds to milliseconds
    return d.toLocaleString();
  }

  onMount(fetchLogs);
</script>

<div class="webhook-logs">
  <header class="modal-header">
    <div class="title-zone">
      <span class="material-symbols-rounded">history</span>
      <div class="title-meta">
        <h3>Delivery History</h3>
        <code>{hook.url}</code>
      </div>
    </div>
    <div class="header-actions">
        <button class="btn btn-icon btn-tonal" onclick={fetchLogs} title="Refresh">
          <span class="material-symbols-rounded">refresh</span>
        </button>
        <button class="btn-close" onclick={() => onclose()}>
          <span class="material-symbols-rounded">close</span>
        </button>
    </div>
  </header>

  <div class="modal-body-split">
     <aside class="logs-sidebar">
        {#if loading}
          <div class="skeleton-logs">
            {#each Array(5) as _}
              <div class="skeleton-log-item"></div>
            {/each}
          </div>
        {:else if error}
          <div class="error-msg">{error}</div>
        {:else if logs.length === 0}
          <div class="empty-logs">
            <span class="material-symbols-rounded">info</span>
            <p>No delivery attempts yet.</p>
          </div>
        {:else}
          <div class="log-list">
            {#each logs as log}
              <button 
                class="log-item" 
                class:active={selectedLog?.id === log.id}
                class:success={log.succeed}
                class:failure={!log.succeed}
                onclick={() => selectedLog = log}
              >
                <span class="material-symbols-rounded status-icon">
                  {log.succeed ? 'check_circle' : 'error'}
                </span>
                <div class="log-info">
                   <div class="log-top">
                     <span class="uuid">{log.uuid.substring(0, 8)}</span>
                     {#if log.status}
                       <span class="status-code">{log.status}</span>
                     {/if}
                   </div>
                   <div class="log-bottom">
                      <span class="event-type">{log.event_type}</span>
                      <span class="date">{formatDate(log.delivered_at)}</span>
                   </div>
                </div>
              </button>
            {/each}
          </div>
        {/if}
     </aside>

     <main class="log-details">
        {#if selectedLog}
           <div class="details-content">
              <section class="details-section">
                <h4>Request</h4>
                <div class="code-block">
                  <div class="headers">
                    {#each Object.entries(selectedLog.request.headers) as [key, val]}
                      <div class="header-row"><strong>{key}:</strong> {val}</div>
                    {/each}
                  </div>
                </div>
              </section>

              <section class="details-section">
                <h4>Response</h4>
                {#if selectedLog.response}
                  <div class="code-block">
                    <div class="status-line" class:success={selectedLog.succeed} class:failure={!selectedLog.succeed}>
                        HTTP/1.1 {selectedLog.response.status}
                    </div>
                    <div class="headers">
                      {#each Object.entries(selectedLog.response.headers) as [key, val]}
                        <div class="header-row"><strong>{key}:</strong> {val}</div>
                      {/each}
                    </div>
                    <div class="body-preview">
                       <h5>Body</h5>
                       <pre>{selectedLog.response.body || '(Empty body)'}</pre>
                    </div>
                  </div>
                {:else}
                  <div class="empty-response">No response received (Timeout or Network Error).</div>
                {/if}
              </section>
           </div>
        {:else}
           <div class="no-selection">
              <span class="material-symbols-rounded">ads_click</span>
              <p>Select a delivery to view details</p>
           </div>
        {/if}
     </main>
  </div>
</div>

<style>
  .webhook-logs { display: flex; flex-direction: column; height: 100%; }
  
  .modal-header { padding: 1rem 2rem; display: flex; justify-content: space-between; align-items: center; border-bottom: 1px solid var(--border-color); }
  .title-zone { display: flex; align-items: center; gap: 1rem; }
  .title-zone span { font-size: 2rem; color: var(--primary-color); }
  .title-meta h3 { margin: 0; font-size: 1.15rem; font-weight: 800; }
  .title-meta code { font-size: 0.8rem; opacity: 0.6; }
  
  .header-actions { display: flex; align-items: center; gap: 1rem; }

  .modal-body-split { display: grid; grid-template-columns: 350px 1fr; height: 700px; max-height: 80vh; overflow: hidden; }

  .logs-sidebar { border-right: 1px solid var(--border-color); background: rgba(0,0,0,0.1); overflow-y: auto; }
  .log-list { display: flex; flex-direction: column; }
  
  .log-item { display: flex; align-items: center; gap: 1rem; padding: 1rem 1.5rem; background: transparent; border: none; border-bottom: 1px solid var(--border-color); color: var(--text-color); cursor: pointer; text-align: left; transition: all 0.2s; width: 100%; }
  .log-item:hover { background: rgba(255,255,255,0.03); }
  .log-item.active { background: rgba(var(--primary-rgb), 0.1); color: var(--primary-color); border-right: 3px solid var(--primary-color); }
  
  .status-icon { font-size: 1.5rem; }
  .log-item.success .status-icon { color: #4ade80; }
  .log-item.failure .status-icon { color: #f87171; }
  
  .log-info { flex: 1; display: flex; flex-direction: column; gap: 0.25rem; }
  .log-top { display: flex; justify-content: space-between; align-items: center; }
  .uuid { font-family: monospace; font-size: 0.8rem; opacity: 0.7; }
  .status-code { font-size: 0.75rem; font-weight: 800; background: rgba(0,0,0,0.2); padding: 0.1rem 0.4rem; border-radius: 4px; }
  
  .log-bottom { display: flex; justify-content: space-between; align-items: baseline; }
  .event-type { font-size: 0.85rem; font-weight: 700; text-transform: uppercase; letter-spacing: 0.05em; }
  .date { font-size: 0.7rem; color: var(--text-muted); }

  .log-details { overflow-y: auto; padding: 2rem; position: relative; }
  .no-selection { height: 100%; display: flex; flex-direction: column; align-items: center; justify-content: center; gap: 1rem; color: var(--text-muted); opacity: 0.5; }
  .no-selection span { font-size: 4rem; }

  .details-content { display: flex; flex-direction: column; gap: 2.5rem; }
  .details-section h4 { margin: 0 0 1rem 0; font-size: 1rem; text-transform: uppercase; letter-spacing: 0.1em; color: var(--text-muted); border-bottom: 1px solid var(--border-color); padding-bottom: 0.5rem; }
  
  .code-block { background: #000; border-radius: 1rem; padding: 1.5rem; font-family: 'JetBrains Mono', 'Fira Code', monospace; font-size: 0.85rem; border: 1px solid var(--border-color); }
  .header-row { margin-bottom: 0.2rem; white-space: pre-wrap; word-break: break-all; }
  .header-row strong { color: var(--primary-color); }
  
  .status-line { font-weight: 800; font-size: 1rem; margin-bottom: 1rem; padding-bottom: 0.5rem; border-bottom: 1px solid rgba(255,255,255,0.1); }
  .status-line.success { color: #4ade80; }
  .status-line.failure { color: #f87171; }
  
  .body-preview { margin-top: 1.5rem; padding-top: 1rem; border-top: 1px solid rgba(255,255,255,0.1); }
  .body-preview h5 { margin: 0 0 0.75rem 0; color: var(--text-muted); font-size: 0.75rem; text-transform: uppercase; }
  .body-preview pre { margin: 0; white-space: pre-wrap; word-break: break-all; opacity: 0.9; line-height: 1.5; }

  .empty-logs { padding: 3rem; text-align: center; color: var(--text-muted); }
  .empty-logs span { font-size: 3rem; opacity: 0.3; margin-bottom: 1rem; }
</style>
