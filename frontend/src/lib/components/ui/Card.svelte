<script lang="ts">
  import type { Snippet } from 'svelte';

  interface Props {
    class?: string;
    children?: Snippet;
    padding?: 'none' | 'sm' | 'md' | 'lg';
    hover?: boolean;
    variant?: 'lowest' | 'low' | 'container' | 'high' | 'highest';
    tonal?: boolean;
  }

  let { 
    class: className = '',
    children,
    padding = 'md',
    hover = false,
    variant = 'container',
    tonal = false
  }: Props = $props();

  const paddings: Record<NonNullable<Props['padding']>, string> = {
    none: 'p-0',
    sm: 'p-sm',
    md: 'p-md',
    lg: 'p-lg'
  };

  const variants: Record<NonNullable<Props['variant']>, string> = {
    lowest: 'bg-lowest',
    low: 'bg-low',
    container: 'bg-container',
    high: 'bg-high',
    highest: 'bg-highest'
  };
</script>

<div class="card {paddings[padding]} {variants[variant]} {hover ? 'hoverable' : ''} {tonal ? 'tonal' : ''} {className}">
  {#if children}
    {@render children()}
  {/if}
</div>

<style>
  .card {
    border: 1px solid var(--md-sys-color-outline-variant);
    border-radius: 1.5rem;
    box-shadow: var(--shadow-sm);
    overflow: hidden;
    transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
    color: var(--md-sys-color-on-surface);
    max-width: 100%;
    width: 100%;
    box-sizing: border-box;
  }

  .bg-lowest { background: var(--md-sys-color-surface-container-lowest); }
  .bg-low { background: var(--md-sys-color-surface-container-low); }
  .bg-container { background: var(--md-sys-color-surface-container); }
  .bg-high { background: var(--md-sys-color-surface-container-high); }
  .bg-highest { background: var(--md-sys-color-surface-container-highest); }
  .tonal {
    background: var(--md-sys-color-secondary-container);
    color: var(--md-sys-color-on-secondary-container);
  }

  .hoverable:hover {
    transform: translateY(-4px);
    box-shadow: var(--shadow-lg);
    border-color: var(--md-sys-color-primary);
    background: var(--md-sys-color-surface-container-high);
  }

  .p-0 { padding: 0; }
  .p-sm { padding: 1rem; }
  .p-md { padding: 1.5rem; }
  .p-lg { padding: 2.5rem; }

  @media (max-width: 768px) {
    .p-md { padding: 1rem; }
    .p-lg { padding: 1.25rem; }
  }
</style>

