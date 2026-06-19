<script lang="ts">
  import type { Snippet } from 'svelte';

  let { children }: { children?: Snippet } = $props();
  let error = $state<Error | null>(null);

  $effect(() => {
    const cleanup = () => {};
    return cleanup;
  });
</script>

{#if error}
  <div style="padding: 2rem; max-width: 800px; margin: 2rem auto;">
    <h2>Error inesperado</h2>
    <p style="color: var(--md-sys-color-error);">{error.message}</p>
    <pre style="overflow: auto; font-size: 0.8rem; max-height: 400px; background: var(--md-sys-color-surface-container-high); padding: 1rem; border-radius: 8px; margin-top: 1rem;">{error.stack}</pre>
    <button onclick={() => { error = null; window.location.reload(); }} style="margin-top: 1rem;">
      Recargar página
    </button>
  </div>
{:else}
  {#if children}
    {@render children()}
  {/if}
{/if}

<script context="module">
  // Global error handler
  export function handleError(e: Error) {
    console.error('ErrorBoundary caught:', e);
  }
</script>
