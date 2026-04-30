<script lang="ts">
  import { fade, fly } from "svelte/transition";
  import { dialog } from "$lib/dialog.svelte";

  let inputValue = $state("");

  $effect(() => {
    if (dialog.current?.type === 'prompt') {
      inputValue = dialog.current.defaultValue || "";
    }
  });

  function handleKeydown(e: KeyboardEvent) {
    if (!dialog.current) return;
    if (e.key === "Escape") dialog.close(null);
    if (e.key === "Enter") {
       if (dialog.current.type === 'prompt') dialog.close(inputValue);
       else dialog.close(true);
    }
  }
</script>

<svelte:window onkeydown={handleKeydown} />

{#if dialog.current}
  <!-- BACKDROP -->
  <div
    class="dialog-backdrop"
    onclick={() => dialog.close(null)}
    transition:fade={{ duration: 200 }}
    aria-hidden="true"
  ></div>

  <!-- DIALOG CONTENT -->
  <div class="dialog-wrapper" transition:fly={{ y: 20, duration: 300 }}>
    <div class="dialog-card tonal">
      <header class="dialog-header">
        <h3>{dialog.current.title}</h3>
      </header>
      
      <div class="dialog-body">
        <p>{dialog.current.message}</p>
        
        {#if dialog.current.type === 'prompt'}
          <input 
            type="text" 
            bind:value={inputValue} 
            class="dialog-input"
            placeholder="Escribe aquí..."
          />
        {/if}
      </div>

      <footer class="dialog-footer">
        {#if dialog.current.type !== 'alert'}
          <button class="btn btn-text" onclick={() => dialog.close(null)}>
            {dialog.current.cancelText || 'Cancelar'}
          </button>
        {/if}
        
        <button class="btn btn-primary" onclick={() => dialog.close(dialog.current?.type === 'prompt' ? inputValue : true)}>
          {dialog.current.confirmText || 'Aceptar'}
        </button>
      </footer>
    </div>
  </div>
{/if}

<style>
  .dialog-backdrop {
    position: fixed;
    top: 0;
    left: 0;
    width: 100vw;
    height: 100vh;
    background: rgba(0, 0, 0, 0.4);
    backdrop-filter: blur(4px);
    z-index: 2000;
  }

  .dialog-wrapper {
    position: fixed;
    top: 0;
    left: 0;
    width: 100vw;
    height: 100vh;
    display: flex;
    align-items: center;
    justify-content: center;
    padding: 2rem;
    z-index: 2001;
    pointer-events: none;
  }

  .dialog-card {
    width: 100%;
    max-width: 450px;
    background: var(--md-sys-color-surface-container-high);
    border-radius: 1.75rem;
    pointer-events: auto;
    display: flex;
    flex-direction: column;
    box-shadow: var(--shadow-xl);
    border: 1px solid var(--md-sys-color-outline-variant);
    overflow: hidden;
  }

  .dialog-header {
    padding: 1.5rem 1.5rem 0.5rem;
  }

  .dialog-header h3 {
    margin: 0;
    font-size: 1.5rem;
    font-weight: 400;
    color: var(--md-sys-color-on-surface);
  }

  .dialog-body {
    padding: 1rem 1.5rem;
    font-size: 1rem;
    color: var(--md-sys-color-on-surface-variant);
  }

  .dialog-input {
    width: 100%;
    margin-top: 1rem;
    padding: 0.75rem 1rem;
    border-radius: 0.5rem;
    border: 1px solid var(--md-sys-color-outline);
    background: var(--md-sys-color-surface);
    color: var(--md-sys-color-on-surface);
    font-family: inherit;
    font-size: 1rem;
  }
  .dialog-input:focus {
    outline: 2px solid var(--md-sys-color-primary);
    border-color: transparent;
  }

  .dialog-footer {
    padding: 1.5rem;
    display: flex;
    justify-content: flex-end;
    gap: 0.5rem;
  }

  .btn {
    padding: 0.625rem 1.5rem;
    border-radius: 100px;
    font-weight: 500;
    font-size: 0.875rem;
    cursor: pointer;
    border: none;
    transition: all 0.2s;
  }

  .btn-text {
    background: transparent;
    color: var(--md-sys-color-primary);
  }
  .btn-text:hover {
    background: rgba(var(--primary-rgb), 0.08);
  }

  .btn-primary {
    background: var(--md-sys-color-primary);
    color: var(--md-sys-color-on-primary);
  }
  .btn-primary:hover {
    box-shadow: var(--shadow-md);
    opacity: 0.9;
  }
</style>
