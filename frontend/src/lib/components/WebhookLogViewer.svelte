<script lang="ts">
  import type { RepoContext } from "$lib/types";
  import { apiFetch } from "$lib/api";
  import { onMount, getContext } from 'svelte';
  

import RefreshCw from "lucide-svelte/icons/refresh-cw";
import Info from "lucide-svelte/icons/info";
import CheckCircle2 from "lucide-svelte/icons/check-circle-2";
import AlertCircle from "lucide-svelte/icons/alert-circle";
import RotateCcw from "lucide-svelte/icons/rotate-ccw";
import MousePointerClick from "lucide-svelte/icons/mouse-pointer-click";


  import Button from './ui/Button.svelte';

  let { hook, onclose } = $props();
  const repo = getContext<RepoContext>('repo');

  let logs = $state<any[]>([]);
  let loading = $state(true);
  let error = $state('');
  
  let selectedLog = $state<any>(null);

  async function fetchLogs() {
    loading = true;
    try {
      const res = await apiFetch(`/api/v2/${repo.data.owner.username}/${repo.data.name}/settings/hooks/${hook.id}/history`);
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

  async function redeliver(uuid: string) {
    try {
      const res = await apiFetch(`/api/v2/${repo.data.owner.username}/${repo.data.name}/settings/hooks/${hook.id}/redeliver/${uuid}`, {
        method: 'POST'
      });
      if (res.ok) {
        fetchLogs();
      }
    } catch (e) {
      console.error(e);
    }
  }

  import { formatDateTime } from '$lib/utils/date';

  function formatDate(iso: string) {
    return formatDateTime(iso);
  }

  onMount(fetchLogs);
</script>

<div class="webhook-logs">
  <div class="modal-body-split">
     <aside class="logs-sidebar">
        <header class="sidebar-header">
           <span class="url-hint">{hook.url}</span>
           <button class="btn-icon-minimal" onclick={fetchLogs} title="Refresh">
             <RefreshCw size={18} />
           </button>
        </header>
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
            <Info size={32} />
            <p>No delivery attempts yet.</p>
          </div>
        {:else}
          <div class="log-list">
            {#each logs as log}
              <button 
                class="log-item" 
                class:active={selectedLog?.id === log.id}
                class:success={log.is_succeed}
                class:failure={!log.is_succeed}
                onclick={() => selectedLog = log}
              >
                {#if log.is_succeed}
                  <CheckCircle2 size={20} class="status-icon" />
                {:else}
                  <AlertCircle size={20} class="status-icon" />
                {/if}
                <div class="log-info">
                   <div class="log-top">
                     <span class="uuid">{log.uuid.substring(0, 8)}</span>
                     {#if log.response_status}
                       <span class="status-code">{log.response_status}</span>
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
              <header class="details-header">
                <div class="header-main">
                  <h3>Detalles del Envío</h3>
                  <span class="full-uuid">UUID: {selectedLog.uuid}</span>
                </div>
                <Button variant="secondary" size="sm" icon={RotateCcw} onclick={() => redeliver(selectedLog.uuid)}>
                  Reenviar
                </Button>
              </header>

              <section class="details-section">
                <h4>Request</h4>
                <div class="code-block">
                  <div class="headers">
                    {#each Object.entries(selectedLog.request_headers) as [key, val]}
                      <div class="header-row"><strong>{key}:</strong> {val}</div>
                    {/each}
                  </div>
                  <div class="body-preview">
                     <h5>Payload</h5>
                     <pre>{selectedLog.request_body || '(Empty body)'}</pre>
                  </div>
                </div>
              </section>

              <section class="details-section">
                <h4>Response</h4>
                {#if selectedLog.response_status}
                  <div class="code-block">
                    <div class="status-line" class:success={selectedLog.is_succeed} class:failure={!selectedLog.is_succeed}>
                        HTTP/1.1 {selectedLog.response_status}
                    </div>
                    <div class="headers">
                      {#each Object.entries(selectedLog.response_headers) as [key, val]}
                        <div class="header-row"><strong>{key}:</strong> {val}</div>
                      {/each}
                    </div>
                    <div class="body-preview">
                       <h5>Body</h5>
                       <pre>{selectedLog.response_body || '(Empty body)'}</pre>
                    </div>
                  </div>
                {:else}
                  <div class="empty-response">No response received (Timeout or Network Error).</div>
                {/if}
              </section>
           </div>
        {:else}
           <div class="no-selection">
              <MousePointerClick size={64} />
              <p>Selecciona un envío para ver los detalles</p>
           </div>
        {/if}
     </main>
  </div>
</div>

<style>
  .webhook-logs { display: flex; flex-direction: column; height: 100%; }
  
  .sidebar-header { display: flex; justify-content: space-between; align-items: center; padding: 1rem 1.5rem; border-bottom: 1px solid var(--md-sys-color-outline-variant); background: rgba(0,0,0,0.05); }
  .url-hint { font-size: 0.75rem; opacity: 0.5; font-family: monospace; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; max-width: 200px; }

  .modal-body-split { display: grid; grid-template-columns: 350px 1fr; height: 700px; max-height: 80vh; overflow: hidden; }

  .logs-sidebar { border-right: 1px solid var(--md-sys-color-outline-variant); background: rgba(0,0,0,0.1); overflow-y: auto; }
  .log-list { display: flex; flex-direction: column; }
  
  .log-item { display: flex; align-items: center; gap: 1rem; padding: 1rem 1.5rem; background: transparent; border: none; border-bottom: 1px solid var(--md-sys-color-outline-variant); color: var(--text-color); cursor: pointer; text-align: left; transition: all 0.2s; width: 100%; }
  .log-item:hover { background: rgba(255,255,255,0.03); }
  .log-item.active { background: rgba(var(--primary-rgb), 0.1); color: var(--primary-color); border-right: 3px solid var(--primary-color); }
  
  :global(.status-icon) { flex-shrink: 0; }
  .log-item.success :global(.status-icon) { color: #4ade80; }
  .log-item.failure :global(.status-icon) { color: #f87171; }
  
  .log-info { flex: 1; display: flex; flex-direction: column; gap: 0.25rem; }
  .log-top { display: flex; justify-content: space-between; align-items: center; }
  .uuid { font-family: monospace; font-size: 0.8rem; opacity: 0.7; }
  .status-code { font-size: 0.75rem; font-weight: 800; background: rgba(0,0,0,0.2); padding: 0.1rem 0.4rem; border-radius: 4px; }
  
  .log-bottom { display: flex; justify-content: space-between; align-items: baseline; }
  .event-type { font-size: 0.85rem; font-weight: 700; text-transform: uppercase; letter-spacing: 0.05em; }
  .date { font-size: 0.7rem; color: var(--text-muted); }

  .log-details { overflow-y: auto; padding: 2rem; position: relative; }
  .no-selection { height: 100%; display: flex; flex-direction: column; align-items: center; justify-content: center; gap: 1rem; color: var(--text-muted); opacity: 0.5; }
  .no-selection :global(svg) { font-size: 4rem; }

  .details-content { display: flex; flex-direction: column; gap: 2rem; }
  
  .details-header { display: flex; justify-content: space-between; align-items: center; background: rgba(var(--primary-rgb), 0.05); padding: 1.5rem; border-radius: 1rem; border: 1px solid rgba(var(--primary-rgb), 0.1); }
  .header-main h3 { margin: 0; font-size: 1.25rem; font-weight: 800; }
  .full-uuid { font-size: 0.75rem; opacity: 0.5; font-family: monospace; }

  .details-section h4 { margin: 0 0 1rem 0; font-size: 0.9rem; text-transform: uppercase; letter-spacing: 0.1em; color: var(--text-muted); border-bottom: 1px solid var(--md-sys-color-outline-variant); padding-bottom: 0.5rem; font-weight: 800; }
  
  .code-block { background: #000; border-radius: 1rem; padding: 1.5rem; font-family: 'JetBrains Mono', 'Fira Code', monospace; font-size: 0.85rem; border: 1px solid var(--md-sys-color-outline-variant); }
  .header-row { margin-bottom: 0.2rem; white-space: pre-wrap; word-break: break-all; }
  .header-row strong { color: var(--primary-color); }
  
  .status-line { font-weight: 800; font-size: 1rem; margin-bottom: 1rem; padding-bottom: 0.5rem; border-bottom: 1px solid rgba(255,255,255,0.1); }
  .status-line.success { color: #4ade80; }
  .status-line.failure { color: #f87171; }
  
  .body-preview { margin-top: 1.5rem; padding-top: 1rem; border-top: 1px solid rgba(255,255,255,0.1); }
  .body-preview h5 { margin: 0 0 0.75rem 0; color: var(--text-muted); font-size: 0.75rem; text-transform: uppercase; }
  .body-preview pre { margin: 0; white-space: pre-wrap; word-break: break-all; opacity: 0.9; line-height: 1.5; }

  .empty-logs { padding: 3rem; text-align: center; color: var(--text-muted); }
  .empty-logs :global(svg) { font-size: 3rem; opacity: 0.3; margin-bottom: 1rem; }

  .skeleton-log-item { height: 60px; background: rgba(255,255,255,0.05); margin: 0.5rem 1rem; border-radius: 0.5rem; animation: pulse 1.5s infinite; }
  @keyframes pulse { 0%, 100% { opacity: 0.6; } 50% { opacity: 0.3; } }

  .btn-icon-minimal { background: transparent; border: none; color: var(--text-muted); cursor: pointer; display: flex; padding: 0.25rem; border-radius: 50%; transition: all 0.2s; }
  .btn-icon-minimal:hover { background: rgba(255,255,255,0.05); color: var(--text-color); }
</style>
