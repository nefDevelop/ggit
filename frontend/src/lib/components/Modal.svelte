<script lang="ts">
  import { fade, fly } from "svelte/transition";
  
import X from "lucide-svelte/icons/x";


  interface Props {
    show: boolean;
    title: string;
    children?: any;
    onclose: () => void;
    maxWidth?: string;
  }

  let { show, title, children, onclose, maxWidth = "600px" }: Props = $props();

  function handleKeydown(e: KeyboardEvent) {
    if (e.key === "Escape" && show) onclose();
  }
  import { focusTrap } from "$lib/utils/focusTrap";
</script>

<svelte:window onkeydown={handleKeydown} />

{#if show}
  <!-- BACKDROP -->
  <div
    class="modal-backdrop"
    onclick={onclose}
    transition:fade={{ duration: 200 }}
    aria-hidden="true"
  ></div>

  <!-- MODAL CONTENT -->
  <div class="modal-wrapper" transition:fly={{ y: 20, duration: 300 }}>
    <div class="modal-card tonal" style="max-width: {maxWidth}" use:focusTrap role="dialog" aria-modal="true" aria-label={title}>
      <header class="modal-header">
        <h3>{title}</h3>
        <button class="close-btn" onclick={onclose} aria-label="Cerrar">
          <X size={20} />
        </button>
      </header>
      <div class="modal-body">
        {@render children?.()}
      </div>
    </div>
  </div>
{/if}

<style>
  .modal-backdrop {
    position: fixed;
    top: 0;
    left: 0;
    width: 100vw;
    height: 100vh;
    background: rgba(0, 0, 0, 0.8);
    backdrop-filter: blur(8px);
    z-index: 1000;
  }

  .modal-wrapper {
    position: fixed;
    top: 0;
    left: 0;
    width: 100vw;
    height: 100vh;
    display: flex;
    align-items: center;
    justify-content: center;
    padding: 2rem;
    z-index: 1001;
    pointer-events: none;
  }

  .modal-card {
    width: 100%;
    max-width: 600px;
    max-height: 80vh;
    pointer-events: auto;
    display: flex;
    flex-direction: column;
    box-shadow: var(--shadow-modal);
    border: 1px solid rgba(255, 255, 255, 0.1);
  }

  .modal-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 1.5rem 2rem;
    border-bottom: 1px solid var(--md-sys-color-outline-variant);
  }

  .modal-header h3 {
    margin: 0;
    font-size: 1.25rem;
    font-weight: 800;
  }

  .close-btn {
    background: transparent;
    border: none;
    color: var(--text-muted);
    cursor: pointer;
    padding: 0.75rem;
    border-radius: 50%;
    transition: all 0.2s;
  }
  .close-btn:hover {
    background: var(--md-sys-color-surface-variant);
    color: #fff;
  }

  .modal-body {
    padding: 2rem;
    overflow-y: auto;
    font-size: 1rem;
    line-height: 1.6;
    color: var(--text-color);
  }

  @media (max-width: 600px) {
    .modal-wrapper { padding: 1rem; }
    .modal-header { padding: 1rem 1.25rem; }
    .modal-body { padding: 1.25rem; }
  }
</style>
