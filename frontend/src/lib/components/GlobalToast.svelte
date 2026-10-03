<script lang="ts">
  import { flip } from "svelte/animate";
  import { fly } from "svelte/transition";
  import { toast } from "$lib/toast.svelte";
</script>

<div class="toast-container">
  {#each toast.list as item (item.id)}
    <div
      class="toast-item {item.type}"
      animate:flip={{ duration: 300 }}
      transition:fly={{ x: 100, duration: 300 }}
      role="alert"
    >
      <span class="material-symbols-rounded icon">
        {#if item.type === 'success'}check_circle
        {:else}info{/if}
      </span>
      <p>{item.message}</p>
      <button class="close-btn" onclick={() => toast.remove(item.id)}>
        <span class="material-symbols-rounded">close</span>
      </button>
    </div>
  {/each}
</div>

<style>
  .toast-container {
    position: fixed;
    bottom: 2rem;
    right: 2rem;
    z-index: 3000;
    display: flex;
    flex-direction: column;
    gap: 0.75rem;
    max-width: 400px;
    width: calc(100% - 4rem);
  }

  .toast-item {
    display: flex;
    align-items: center;
    gap: 1rem;
    padding: 1rem 1.25rem;
    border-radius: 1.5rem;
    background: var(--md-sys-color-surface-container-highest);
    color: var(--md-sys-color-on-surface);
    box-shadow: var(--shadow-lg);
    border: 1px solid var(--md-sys-color-outline-variant);
  }

  .toast-item.success {
    background: var(--md-sys-color-primary-container);
    color: var(--md-sys-color-on-primary-container);
    border-color: var(--md-sys-color-primary);
  }
  .toast-item.success .icon { color: var(--md-sys-color-primary); }

  .toast-item.error {
    background: var(--md-sys-color-error-container);
    color: var(--md-sys-color-on-error-container);
    border-color: var(--md-sys-color-error);
  }
  .toast-item.error .icon { color: var(--md-sys-color-error); }

  .toast-item.warning {
    background: var(--status-warning-wash);
    color: var(--status-warning-strong);
    border-color: var(--status-warning);
  }

  .icon {
    font-size: 1.5rem;
    flex-shrink: 0;
  }

  p {
    margin: 0;
    font-size: 0.9375rem;
    font-weight: 500;
    flex: 1;
  }

  .close-btn {
    background: transparent;
    border: none;
    color: inherit;
    opacity: 0.6;
    cursor: pointer;
    padding: 0.25rem;
    display: flex;
    transition: opacity 0.2s;
  }
  .close-btn:hover {
    opacity: 1;
  }
</style>
