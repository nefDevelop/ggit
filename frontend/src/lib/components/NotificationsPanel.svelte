<script lang="ts">
  import { fade, fly } from "svelte/transition";
  import { toast } from "$lib/toast.svelte";
  import { authState } from "$lib/auth.svelte";
  import { apiFetch } from "$lib/api";
  import { onMount } from "svelte";
  import Button from "./ui/Button.svelte";
  import { X, Eraser, BellOff, CheckCircle2, AlertCircle, AlertTriangle, Info, RefreshCw } from "lucide-svelte";

  interface Props {
    show: boolean;
    onclose: () => void;
  }

  let { show = $bindable(false), onclose }: Props = $props();

  let activeTab = $state<'local' | 'system'>('local');
  let systemNotices = $state<any[]>([]);
  let loadingSystemNotices = $state(false);

  async function fetchSystemNotices() {
    if (!authState.user?.is_admin) return;
    loadingSystemNotices = true;
    try {
      const res = await apiFetch("/api/v2/admin/notices?page=1");
      if (res.ok) {
        const data = await res.json();
        systemNotices = data?.notices || [];
      }
    } catch (e) {
      console.error("Failed to fetch system notices", e);
    } finally {
      loadingSystemNotices = false;
    }
  }

  $effect(() => {
    if (show) {
      toast.clearUnread();
      if (authState.user?.is_admin && activeTab === 'system') {
        fetchSystemNotices();
      }
    }
  });

  function closePanel() {
    onclose();
  }
</script>

{#if show}
  <!-- svelte-ignore a11y_click_events_have_key_events -->
  <!-- svelte-ignore a11y_no_noninteractive_element_interactions -->
  <div
    class="panel-backdrop"
    onclick={closePanel}
    transition:fade={{ duration: 250 }}
    role="button"
    tabindex="0"
  >
    <!-- svelte-ignore a11y_no_static_element_interactions -->
    <div
      class="panel-content"
      onclick={(e) => e.stopPropagation()}
      transition:fly={{ x: -320, duration: 300, opacity: 1 }}
      role="none"
    >
      <header class="panel-header">
        <div class="header-top">
          <h3>Centro de Avisos</h3>
          <Button variant="ghost" size="sm" icon={X} onclick={closePanel} aria-label="Cerrar panel" />
        </div>

        {#if authState.user?.is_admin}
          <div class="tab-header">
            <button
              class="tab-btn"
              class:active={activeTab === 'local'}
              onclick={() => activeTab = 'local'}
            >
              Recientes ({toast.history.length})
            </button>
            <button
              class="tab-btn"
              class:active={activeTab === 'system'}
              onclick={() => {
                activeTab = 'system';
                fetchSystemNotices();
              }}
            >
              Sistema ({systemNotices.length})
            </button>
          </div>
        {/if}
      </header>

      <div class="panel-body">
        {#if activeTab === 'local'}
          <div class="list-section">
            <div class="section-actions">
              <span class="section-title">Historial de sesión</span>
              {#if toast.history.length > 0}
                <Button
                  variant="ghost" 
                  size="sm"
                  icon={Eraser}
                  onclick={() => toast.clearHistory()}
                >
                  Limpiar historial
                </Button>
              {/if}
            </div>

            {#if toast.history.length === 0}
              <div class="empty-state">
                <BellOff size={40} class="empty-icon" />
                <p>No hay notificaciones recientes.</p>
              </div>
            {:else}
              <div class="notices-list">
                {#each toast.history as item (item.id)}
                  <div class="notice-item {item.type}">
                    <div class="type-icon">
                      {#if item.type === 'success'}<CheckCircle2 size={18} />
                      {:else if item.type === 'error'}<AlertCircle size={18} />
                      {:else if item.type === 'warning'}<AlertTriangle size={18} />
                      {:else}<Info size={18} />{/if}
                    </div>
                    <div class="notice-text">
                      <p class="message">{item.message}</p>
                      <span class="time">{new Date(item.timestamp).toLocaleTimeString()}</span>
                    </div>
                    <Button variant="ghost" size="sm" icon={X} onclick={() => toast.remove(item.id)} />
                  </div>
                {/each}
              </div>
            {/if}
          </div>
        {:else if activeTab === 'system'}
          <div class="list-section">
            <div class="section-actions">
              <span class="section-title">Avisos globales del sistema</span>
              <Button variant="ghost" size="sm" icon={RefreshCw} onclick={fetchSystemNotices}>
                Actualizar
              </Button>
            </div>

            {#if loadingSystemNotices}
              <div class="skeleton-list">
                <div class="skeleton-item"></div>
                <div class="skeleton-item"></div>
                <div class="skeleton-item"></div>
              </div>
            {:else if systemNotices.length === 0}
              <div class="empty-state">
                <CheckCircle2 size={40} class="empty-icon" />
                <p>No se encontraron avisos de sistema.</p>
              </div>
            {:else}
              <div class="notices-list">
                {#each systemNotices as item}
                  <div class="notice-item system tier-{item.type}">
                    <div class="type-icon">
                      {#if item.type === 3}<AlertCircle size={18} />
                      {:else if item.type === 2}<AlertTriangle size={18} />
                      {:else}<Info size={18} />{/if}
                    </div>
                    <div class="notice-text">
                      <span class="system-badge">{item.type_name}</span>
                      <p class="message">{item.description}</p>
                      <span class="time">{new Date(item.created_at).toLocaleString()}</span>
                    </div>
                  </div>
                {/each}
              </div>
            {/if}
          </div>
        {/if}
      </div>
    </div>
  </div>
{/if}

<style>
  .panel-backdrop {
    position: absolute;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    background: rgba(0, 0, 0, 0.4);
    backdrop-filter: blur(4px);
    z-index: 2100;
    display: flex;
  }

  .panel-content {
    background: var(--md-sys-color-surface-container);
    color: var(--md-sys-color-on-surface);
    width: 320px;
    height: 100%;
    box-shadow: 10px 0 30px rgba(0, 0, 0, 0.25);
    border-right: 1px solid var(--md-sys-color-outline-variant);
    display: flex;
    flex-direction: column;
    animation: none;
  }

  .panel-header {
    padding: 1.5rem 1.25rem 0.5rem;
    display: flex;
    flex-direction: column;
    gap: 1rem;
    border-bottom: 1px solid var(--md-sys-color-outline-variant);
  }

  .header-top {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }

  h3 {
    margin: 0;
    font-size: 1.15rem;
    font-weight: 800;
    letter-spacing: -0.02em;
    color: var(--md-sys-color-on-surface);
  }

  .close-btn {
    color: var(--md-sys-color-on-surface-variant);
  }
  .close-btn:hover {
    color: var(--md-sys-color-error);
    background: rgba(239, 68, 68, 0.1);
  }

  .icon-btn {
    background: transparent;
    border: none;
    cursor: pointer;
    padding: 0.5rem;
    display: flex;
    border-radius: 50%;
    transition: background 0.2s, color 0.2s;
  }

  .tab-header {
    display: flex;
    background: var(--md-sys-color-surface-container-high);
    padding: 0.25rem;
    border-radius: 1rem;
    gap: 0.25rem;
  }

  .tab-btn {
    flex: 1;
    background: transparent;
    border: none;
    color: var(--md-sys-color-on-surface-variant);
    font-size: 0.75rem;
    font-weight: 700;
    padding: 0.5rem;
    border-radius: 0.75rem;
    cursor: pointer;
    transition: all 0.2s ease;
  }

  .tab-btn.active {
    background: var(--md-sys-color-surface-container-highest);
    color: var(--md-sys-color-primary);
  }

  .panel-body {
    flex: 1;
    overflow-y: auto;
    padding: 1.25rem;
  }

  .list-section {
    display: flex;
    flex-direction: column;
    gap: 1rem;
  }

  .section-actions {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }

  .section-title {
    font-size: 0.75rem;
    text-transform: uppercase;
    font-weight: 800;
    letter-spacing: 0.05em;
    color: var(--md-sys-color-on-surface-variant);
    opacity: 0.7;
  }

  .action-text-btn {
    background: transparent;
    border: none;
    color: var(--md-sys-color-primary);
    cursor: pointer;
    font-size: 0.75rem;
    font-weight: 700;
    display: flex;
    align-items: center;
    gap: 0.25rem;
    padding: 0.25rem 0.5rem;
    border-radius: 0.5rem;
    transition: background 0.2s;
  }

  .action-text-btn:hover {
    background: rgba(var(--primary-rgb), 0.08);
  }

  .mini-icon {
    font-size: 1rem;
  }

  .empty-state {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    padding: 3rem 1rem;
    text-align: center;
    color: var(--md-sys-color-on-surface-variant);
    opacity: 0.6;
    gap: 0.5rem;
  }

  .empty-icon {
    font-size: 2.5rem;
  }

  .empty-state p {
    margin: 0;
    font-size: 0.85rem;
    font-weight: 600;
  }

  .notices-list {
    display: flex;
    flex-direction: column;
    gap: 0.75rem;
  }

  .notice-item {
    display: flex;
    align-items: flex-start;
    gap: 0.75rem;
    padding: 0.85rem;
    border-radius: 1rem;
    background: var(--md-sys-color-surface-container-highest);
    color: var(--md-sys-color-on-surface);
    border: 1px solid var(--md-sys-color-outline-variant);
    position: relative;
    animation: slideDownItem 0.3s ease-out;
  }

  @keyframes slideDownItem {
    from { opacity: 0; transform: translateY(-8px); }
    to { opacity: 1; transform: translateY(0); }
  }

  .notice-item.success {
    border-color: var(--md-sys-color-primary);
    background: var(--md-sys-color-primary-container);
    color: var(--md-sys-color-on-primary-container);
  }
  .notice-item.success .type-icon {
    color: var(--md-sys-color-primary);
  }

  .notice-item.error {
    border-color: var(--md-sys-color-error);
    background: var(--md-sys-color-error-container);
    color: var(--md-sys-color-on-error-container);
  }
  .notice-item.error .type-icon {
    color: var(--md-sys-color-error);
  }

  .notice-item.warning {
    border-color: #ffb74d;
    background: rgba(255, 183, 77, 0.1);
  }
  .notice-item.warning .type-icon {
    color: #ffa726;
  }

  .type-icon {
    font-size: 1.25rem;
    flex-shrink: 0;
    margin-top: 0.1rem;
  }

  .notice-text {
    flex: 1;
    display: flex;
    flex-direction: column;
    gap: 0.25rem;
    min-width: 0;
  }

  .message {
    margin: 0;
    font-size: 0.825rem;
    line-height: 1.35;
    font-weight: 500;
    word-break: break-word;
  }

  .time {
    font-size: 0.7rem;
    opacity: 0.6;
    font-weight: 600;
  }

  .remove-btn {
    background: transparent;
    border: none;
    color: inherit;
    opacity: 0.4;
    cursor: pointer;
    padding: 0.15rem;
    display: flex;
    border-radius: 50%;
    transition: opacity 0.2s;
    align-self: flex-start;
  }

  .remove-btn:hover {
    opacity: 1;
  }

  .remove-btn span {
    font-size: 1.1rem;
  }

  /* SYSTEM NOTICES SPECIFICS */
  .notice-item.system {
    border-left: 4px solid var(--md-sys-color-primary);
  }

  .notice-item.system.tier-3 {
    border-left-color: var(--md-sys-color-error);
  }

  .notice-item.system.tier-2 {
    border-left-color: #ffa726;
  }

  .system-badge {
    font-size: 0.65rem;
    font-weight: 800;
    text-transform: uppercase;
    color: var(--md-sys-color-primary);
    opacity: 0.8;
  }

  /* Skeleton loading styles */
  .skeleton-list {
    display: flex;
    flex-direction: column;
    gap: 0.75rem;
  }

  .skeleton-item {
    height: 70px;
    background: var(--md-sys-color-surface-container-highest);
    border-radius: 1rem;
    opacity: 0.6;
    animation: skeleton-pulse 1.5s infinite ease-in-out;
  }

  @keyframes skeleton-pulse {
    0% { opacity: 0.3; }
    50% { opacity: 0.7; }
    100% { opacity: 0.3; }
  }
</style>
