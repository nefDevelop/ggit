<script lang="ts">
  import type { Snippet } from 'svelte';

  interface Props {
    icon?: any;
    title?: string;
    description?: string;
    children?: Snippet;
  }

  let { icon: Icon, title = '', description = '', children }: Props = $props();
</script>

<div class="empty-state">
  {#if Icon}
    <div class="empty-icon"><Icon size={40} /></div>
  {/if}
  {#if title}<h3>{title}</h3>{/if}
  {#if description}<p>{description}</p>{/if}
  {#if children}{@render children()}{/if}
</div>

<style>
  .empty-state {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: var(--space-5, 1.5rem);
    padding: var(--space-8, 4rem) var(--space-6, 2rem);
    text-align: center;
    color: var(--md-sys-color-on-surface-variant);
  }

  .empty-icon {
    color: var(--md-sys-color-on-surface-variant);
    opacity: 0.45;
  }

  h3 {
    margin: 0;
    font-size: 1.25rem;
    font-weight: 800;
    color: var(--md-sys-color-on-surface);
  }

  p {
    margin: 0;
    font-size: 0.95rem;
    line-height: 1.5;
    max-width: 60ch;
  }
</style>
