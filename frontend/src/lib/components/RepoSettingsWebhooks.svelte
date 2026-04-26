<script lang="ts">
  import { apiFetch } from "$lib/api";
  import { onMount, getContext } from 'svelte';
  import WebhookEditor from './WebhookEditor.svelte';
  import WebhookLogViewer from './WebhookLogViewer.svelte';
  import Modal from '$lib/components/Modal.svelte';
  import TableSkeleton from '$lib/components/ui/TableSkeleton.svelte';
  import { dialog } from '$lib/dialog.svelte';

  const repo: any = getContext('repo');
  
  let webhooks = $state<any[]>([]);
  let loading = $state(true);
  let error = $state('');

  let editingHook = $state<any>(null);
  let showingLogsFor = $state<any>(null);
  let showEditor = $state(false);

  async function fetchWebhooks() {
    loading = true;
    try {
      const res = await apiFetch(`/api/v2/${repo.data.owner.username}/${repo.data.name}/settings/hooks`);
      if (res.ok) {
        webhooks = await res.json();
      } else {
        error = 'Failed to load webhooks.';
      }
    } catch (e) {
      error = 'An error occurred while fetching webhooks.';
    } finally {
      loading = false;
    }
  }

  function handleAdd() {
    editingHook = null;
    showEditor = true;
  }

  function handleEdit(hook: any) {
    editingHook = hook;
    showEditor = true;
  }

  function handleLogs(hook: any) {
    showingLogsFor = hook;
  }

  async function handleDelete(id: number) {
    if (!await dialog.confirm('Are you sure you want to delete this webhook?')) return;
    
    try {
      const res = await apiFetch(`/api/v2/${repo.data.owner.username}/${repo.data.name}/settings/hooks/${id}`, {
        method: 'DELETE'
      });
      if (res.ok) {
        webhooks = webhooks.filter(w => w.id !== id);
      }
    } catch (e) {
      console.error(e);
    }
  }

  function handleEditorClose(refresh = false) {
    showEditor = false;
    editingHook = null;
    if (refresh) fetchWebhooks();
  }

  onMount(fetchWebhooks);
</script>

<div class="settings-view-flat">
  <div class="header-action-minimal">
    <button class="btn-action-primary" onclick={handleAdd}>
      <span class="material-symbols-rounded">add</span>
      New Webhook
    </button>
  </div>

  {#if loading}
    <div class="hooks-list-flat">
      <table style="width: 100%;">
        <tbody>
          <TableSkeleton columns={1} rows={3} />
        </tbody>
      </table>
    </div>
  {:else if error}
    <div class="error-msg">{error}</div>
  {:else if webhooks.length === 0}
    <div class="empty-state-flat">
      <span class="material-symbols-rounded">webhook</span>
      <p>No webhooks configured for this repository.</p>
    </div>
  {:else}
    <div class="hooks-list-flat">
      {#each webhooks as hook}
        <div class="hook-row-flat">
          <div class="hook-core">
            <div class="hook-url-row">
              <span class="type-badge">{hook.type}</span>
              <code class="url-code">{hook.url}</code>
            </div>
            <div class="hook-status-meta">
              <div class="status-indicator">
                {#if hook.last_status === 'success'}
                  <span class="material-symbols-rounded success">check_circle</span>
                {:else if hook.last_status === 'failure'}
                  <span class="material-symbols-rounded failure">error</span>
                {:else}
                  <span class="material-symbols-rounded neutral">help</span>
                {/if}
              </div>
              <span class="events-summary">{hook.events.join(', ')}</span>
            </div>
          </div>
          <div class="hook-row-actions">
            <button class="btn-icon-minimal" onclick={() => handleLogs(hook)} title="Logs">
              <span class="material-symbols-rounded">history</span>
            </button>
            <button class="btn-icon-minimal" onclick={() => handleEdit(hook)} title="Edit">
              <span class="material-symbols-rounded">edit</span>
            </button>
            <button class="btn-icon-minimal delete" onclick={() => handleDelete(hook.id)} title="Delete">
              <span class="material-symbols-rounded">delete</span>
            </button>
          </div>
        </div>
      {/each}
    </div>
  {/if}

  <Modal 
    show={showEditor} 
    title={editingHook ? 'Edit Webhook' : 'Add New Webhook'} 
    onclose={() => handleEditorClose()}
  >
    <WebhookEditor 
      hook={editingHook} 
      onclose={handleEditorClose} 
    />
  </Modal>

  <Modal 
    show={!!showingLogsFor} 
    title="Webhook Logs" 
    onclose={() => showingLogsFor = null}
    maxWidth="1000px"
  >
    {#if showingLogsFor}
      <WebhookLogViewer 
        hook={showingLogsFor} 
        onclose={() => showingLogsFor = null} 
      />
    {/if}
  </Modal>
</div>

<style>
  .settings-view-flat { display: flex; flex-direction: column; gap: 2rem; padding: 1rem; animation: fadeIn 0.3s ease-out; }
  
  .header-action-minimal { display: flex; justify-content: flex-end; }

  .hooks-list-flat { display: flex; flex-direction: column; background: var(--bg-color); border: 1px solid var(--md-sys-color-outline-variant); border-radius: 1rem; overflow: hidden; }
  
  .hook-row-flat { display: flex; justify-content: space-between; align-items: center; padding: 1.25rem 1.5rem; border-bottom: 1px solid var(--md-sys-color-outline-variant); transition: background 0.2s; }
  .hook-row-flat:last-child { border-bottom: none; }
  .hook-row-flat:hover { background: var(--surface-variant); }

  .hook-core { display: flex; flex-direction: column; gap: 0.5rem; flex: 1; min-width: 0; }
  .hook-url-row { display: flex; align-items: center; gap: 1rem; }
  .type-badge { background: var(--primary-color); color: var(--on-primary); font-size: 0.7rem; font-weight: 800; text-transform: uppercase; padding: 0.2rem 0.6rem; border-radius: 4px; letter-spacing: 0.05em; }
  .url-code { font-size: 0.9rem; opacity: 0.9; font-family: 'JetBrains Mono', monospace; }

  .hook-status-meta { display: flex; align-items: center; gap: 1.5rem; font-size: 0.85rem; color: var(--text-muted); }
  .status-indicator { display: flex; align-items: center; }
  .status-indicator .success { color: var(--success); }
  .status-indicator .failure { color: var(--error); }
  .status-indicator .neutral { color: var(--text-muted); opacity: 0.5; }
  
  .events-summary { opacity: 0.7; }

  .hook-row-actions { display: flex; gap: 0.5rem; }

  .empty-state-flat { text-align: center; padding: 4rem; color: var(--text-muted); opacity: 0.5; display: flex; flex-direction: column; align-items: center; gap: 1rem; }
  .empty-state-flat span { font-size: 3rem; }

  @keyframes fadeIn { from { opacity: 0; transform: translateY(10px); } to { opacity: 1; transform: translateY(0); } }
</style>
