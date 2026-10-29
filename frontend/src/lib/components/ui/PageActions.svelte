<script lang="ts">
  import type { Snippet } from 'svelte';

  interface Props {
    children?: Snippet;
    class?: string;
  }

  let { children, class: className = '' }: Props = $props();
</script>

<div class="page-actions {className}">
  {@render children?.()}
</div>

<style>
  .page-actions {
    display: flex;
    justify-content: flex-end;
    align-items: center;
    gap: var(--space-3, 0.75rem);
    flex-wrap: wrap;
  }
</style>
