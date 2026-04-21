<script lang="ts">
  import { apiFetch } from "$lib/api";
  import { onMount, getContext } from 'svelte';
  import WebhookEditor from './WebhookEditor.svelte';
  import WebhookLogViewer from './WebhookLogViewer.svelte';

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
    if (!confirm('Are you sure you want to delete this webhook?')) return;
    
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

<div class="webhooks-manager">
  <div class="header-action-minimal">
    <button class="btn-action-primary" onclick={handleAdd}>
      <span class="material-symbols-rounded">add</span>
      New Webhook
    </button>
  </div>

  {#if loading}
    <div class="skeleton-list">
      <div class="skeleton-item island tonal"></div>
      <div class="skeleton-item island tonal"></div>
    </div>
  {:else if error}
    <div class="error-banner island">
       <span class="material-symbols-rounded">error</span>
       {error}
    </div>
  {:else if webhooks.length === 0}
    <div class="empty-state island tonal">
      <span class="material-symbols-rounded size-xl">webhook</span>
      <h3>No Webhooks Defined</h3>
      <p>You haven't configured any webhooks for this repository yet.</p>
      <button class="btn btn-tonal" onclick={handleAdd}>Create your first webhook</button>
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

  {#if showEditor}
    <div class="modal-overlay" onclick={() => handleEditorClose()}>
      <div class="modal-content" onclick={e => e.stopPropagation()}>
        <WebhookEditor 
          hook={editingHook} 
          onclose={handleEditorClose} 
        />
      </div>
    </div>
  {/if}

  {#if showingLogsFor}
    <div class="modal-overlay" onclick={() => showingLogsFor = null}>
      <div class="modal-content logs-modal" onclick={e => e.stopPropagation()}>
        <WebhookLogViewer 
          hook={showingLogsFor} 
          onclose={() => showingLogsFor = null} 
        />
      </div>
    </div>
  {/if}
</div>

<style>
  .webhooks-manager { display: flex; flex-direction: column; gap: 2rem; animation: fadeIn 0.3s ease-out; }
  
  .section-header { display: flex; justify-content: space-between; align-items: flex-start; }
  .header-content h2 { margin: 0 0 0.5rem 0; font-size: 1.5rem; font-weight: 800; }
  .header-content p { margin: 0; color: var(--text-muted); font-size: 0.95rem; }

  .webhooks-list { display: flex; flex-direction: column; gap: 1rem; }
  .hook-item { display: flex; justify-content: space-between; align-items: center; padding: 1.25rem 1.5rem; transition: transform 0.2s, background 0.2s; }
  .hook-item:hover { transform: translateX(5px); background: rgba(255,255,255,0.03); }

  .hook-main { display: flex; align-items: center; gap: 1rem; margin-bottom: 0.5rem; }
  .hook-type-badge { background: var(--primary-color); color: var(--on-primary); font-size: 0.7rem; font-weight: 800; text-transform: uppercase; padding: 0.2rem 0.5rem; border-radius: 4px; }
  .hook-url { font-size: 0.95rem; opacity: 0.9; }

  .hook-meta { display: flex; align-items: center; gap: 1.5rem; font-size: 0.85rem; color: var(--text-muted); }
  .status-badge { display: flex; align-items: center; gap: 0.4rem; font-weight: 600; }
  .status-badge span { font-size: 1.1rem; }
  .status-badge.success { color: #4ade80; }
  .status-badge.failure { color: #f87171; }
  
  .hook-events { position: relative; padding-left: 1.5rem; }
  .hook-events::before { content: '•'; position: absolute; left: 0.5rem; }

  .hook-actions { display: flex; gap: 0.5rem; }
  .btn-icon { background: transparent; border: none; color: var(--text-muted); padding: 0.5rem; border-radius: 0.5rem; cursor: pointer; transition: all 0.2s; display: flex; }
  .btn-icon:hover { background: rgba(255,255,255,0.05); color: var(--text-color); }
  .btn-icon.delete:hover { background: rgba(248, 113, 113, 0.1); color: #f87171; }

  .empty-state { padding: 4rem; display: flex; flex-direction: column; align-items: center; gap: 1rem; text-align: center; }
  .empty-state .size-xl { font-size: 4rem; opacity: 0.2; }
  .empty-state h3 { margin: 1rem 0 0 0; }
  .empty-state p { color: var(--text-muted); max-width: 400px; margin-bottom: 1rem; }

  .modal-overlay { position: fixed; inset: 0; background: rgba(0,0,0,0.7); backdrop-filter: blur(4px); z-index: 1000; display: flex; align-items: center; justify-content: center; padding: 2rem; }
  .modal-content { background: var(--md3-surface-1); border-radius: 1.5rem; border: 1px solid var(--border-color); width: 100%; max-width: 700px; max-height: 90vh; overflow-y: auto; box-shadow: var(--shadow-2xl); }
  .modal-content.logs-modal { max-width: 1000px; }

  .skeleton-list { display: flex; flex-direction: column; gap: 1rem; }
  .skeleton-item { height: 80px; width: 100%; animation: pulse 1.5s infinite; }
  @keyframes pulse { 0%, 100% { opacity: 0.6; } 50% { opacity: 0.3; } }
  
  @keyframes fadeIn { from { opacity: 0; transform: translateY(10px); } to { opacity: 1; transform: translateY(0); } }
</style>
