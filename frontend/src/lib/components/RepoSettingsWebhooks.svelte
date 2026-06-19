<script lang="ts">
  import { apiFetch } from "$lib/api";
  import { onMount, getContext } from 'svelte';
  import WebhookEditor from './WebhookEditor.svelte';
  import WebhookLogViewer from './WebhookLogViewer.svelte';
  import Modal from '$lib/components/Modal.svelte';
  import Card from '$lib/components/ui/Card.svelte';
  import PageHeader from '$lib/components/ui/PageHeader.svelte';
  import Button from '$lib/components/ui/Button.svelte';
  import TableSkeleton from '$lib/components/ui/TableSkeleton.svelte';
  import { dialog } from '$lib/dialog.svelte';
  import { i18n } from '$lib/i18n.svelte';
  

import Webhook from "lucide-svelte/icons/webhook";
import Plus from "lucide-svelte/icons/plus";



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
    if (!await dialog.confirm('¿Estás seguro de que deseas eliminar este webhook?')) return;
    
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

<Card padding="none" class="settings-premium-card">
  <PageHeader 
    title="Webhooks" 
    subtitle="Configura eventos que dispararán peticiones HTTP POST a tus servicios externos."
    icon={Webhook}
  >
    {#snippet actions()}
      <Button variant="primary" icon={Plus} onclick={handleAdd}>Nuevo Webhook</Button>
    {/snippet}
  </PageHeader>

  <div class="settings-content">
    {#if loading}
      <div class="skeleton-list">
        {#each Array(3) as _}
          <div class="skeleton-item tonal"></div>
        {/each}
      </div>
    {:else if error}
      <div class="error-banner">
        <span class="material-symbols-rounded">error</span>
        {error}
      </div>
    {:else if webhooks.length === 0}
      <div class="empty-state island">
        <span class="material-symbols-rounded">webhook</span>
        <p>No hay webhooks configurados para este repositorio.</p>
        <Button variant="outline" size="sm" onclick={handleAdd}>Configurar el primero</Button>
      </div>
    {:else}
      <div class="hooks-stack">
        {#each webhooks as hook}
          <div class="hook-tile tonal">
            <div class="hook-main">
              <div class="hook-url-row">
                <span class="type-badge">{hook.type}</span>
                <code class="url-code">{hook.url}</code>
              </div>
              <div class="hook-meta">
                <div class="status-indicator">
                  {#if hook.last_status === 'success'}
                    <span class="material-symbols-rounded success">check_circle</span>
                    <span class="status-text">Último envío OK</span>
                  {:else if hook.last_status === 'failure'}
                    <span class="material-symbols-rounded failure">error</span>
                    <span class="status-text">Fallo en último envío</span>
                  {:else}
                    <span class="material-symbols-rounded neutral">help</span>
                    <span class="status-text">Sin envíos aún</span>
                  {/if}
                </div>
                <div class="events-list">
                  <span class="material-symbols-rounded">event</span>
                  {hook.events.join(', ')}
                </div>
              </div>
            </div>
            <div class="hook-actions">
              <Button variant="ghost" size="sm" icon="history" onclick={() => handleLogs(hook)} title="Historial" />
              <Button variant="ghost" size="sm" icon="edit" onclick={() => handleEdit(hook)} title="Editar" />
              <Button variant="ghost" size="sm" icon="delete" class="danger-btn" onclick={() => handleDelete(hook.id)} title="Eliminar" />
            </div>
          </div>
        {/each}
      </div>
    {/if}
  </div>
</Card>

<Modal 
  show={showEditor} 
  title={editingHook ? 'Editar Webhook' : 'Nuevo Webhook'} 
  onclose={() => handleEditorClose()}
  maxWidth="800px"
>
  <WebhookEditor 
    hook={editingHook} 
    onclose={handleEditorClose} 
  />
</Modal>

<Modal 
  show={!!showingLogsFor} 
  title="Historial de Webhook" 
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

<style>
  .settings-content { padding: 2.5rem; }
  @media (max-width: 600px) { .settings-content { padding: 1.5rem; } }

  .hooks-stack { display: flex; flex-direction: column; gap: 1rem; }
  
  .hook-tile { padding: 1.25rem 1.75rem; border-radius: 1.5rem; border: 1px solid var(--md-sys-color-outline-variant); display: flex; justify-content: space-between; align-items: center; gap: 2rem; transition: all 0.2s; }
  .hook-tile:hover { transform: translateX(4px); border-color: var(--primary-color); background: rgba(var(--primary-rgb), 0.02); }

  .hook-main { display: flex; flex-direction: column; gap: 0.5rem; flex: 1; min-width: 0; }
  .hook-url-row { display: flex; align-items: center; gap: 1rem; }
  .type-badge { background: var(--primary-color); color: var(--on-primary); font-size: 0.65rem; font-weight: 900; text-transform: uppercase; padding: 0.2rem 0.6rem; border-radius: 6px; letter-spacing: 0.05em; }
  .url-code { font-size: 0.95rem; font-family: 'JetBrains Mono', monospace; color: var(--text-color); overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }

  .hook-meta { display: flex; align-items: center; gap: 2rem; }
  .status-indicator { display: flex; align-items: center; gap: 0.5rem; font-size: 0.85rem; font-weight: 600; }
  .status-indicator .success { color: var(--success); font-size: 1.25rem; }
  .status-indicator .failure { color: var(--error); font-size: 1.25rem; }
  .status-indicator .neutral { color: var(--text-muted); opacity: 0.5; font-size: 1.25rem; }
  .status-text { color: var(--text-muted); }

  .events-list { display: flex; align-items: center; gap: 0.4rem; font-size: 0.8rem; color: var(--text-muted); opacity: 0.8; font-weight: 600; }
  .events-list span { font-size: 1rem; }

  .hook-actions { display: flex; gap: 0.5rem; }
  :global(.danger-btn:hover) { color: var(--error) !important; background: rgba(var(--error-rgb), 0.1) !important; }

  @media (max-width: 800px) {
    .hook-tile { flex-direction: column; align-items: stretch; gap: 1.5rem; }
    .hook-meta { flex-direction: column; align-items: flex-start; gap: 0.75rem; }
    .hook-actions { justify-content: flex-end; border-top: 1px solid var(--md-sys-color-outline-variant); padding-top: 0.75rem; }
  }

  .empty-state { padding: 5rem 2rem; text-align: center; display: flex; flex-direction: column; align-items: center; gap: 1.25rem; }
  .empty-state span { font-size: 4rem; color: var(--primary-color); opacity: 0.3; }
  .empty-state p { font-size: 1.1rem; font-weight: 700; color: var(--text-muted); margin: 0; }

  .skeleton-list { display: flex; flex-direction: column; gap: 1rem; }
  .skeleton-item { height: 80px; border-radius: 1.25rem; animation: pulse 1.5s infinite; }
  @keyframes pulse { 50% { opacity: 0.5; } }

  .error-banner { display: flex; align-items: center; gap: 1rem; padding: 1.5rem; border-radius: 1.5rem; background: rgba(var(--error-rgb), 0.1); color: var(--error); font-weight: 700; }

  @keyframes fadeIn { from { opacity: 0; transform: translateY(10px); } to { opacity: 1; transform: translateY(0); } }
</style>
