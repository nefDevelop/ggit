<script lang="ts">
  interface Props {
    width?: string;
    height?: string;
    radius?: string;
    class?: string;
    circle?: boolean;
    variant?: 'list' | 'form' | 'card';
    rows?: number;
  }

  let {
    width = '100%',
    height = '1rem',
    radius = '0.5rem',
    class: className = '',
    circle = false,
    variant,
    rows = 3
  }: Props = $props();

  const style = $derived(`
    width: ${width}; 
    height: ${height}; 
    border-radius: ${circle ? '50%' : radius};
  `);
</script>

{#if variant === 'card'}
  <div class="skeleton skeleton-card {className}"></div>
{:else if variant}
  <div class="skeleton-stack {className}">
    {#each Array(rows) as _}
      <div class="skeleton {variant === 'form' ? 'skeleton-line' : 'skeleton-row'}"></div>
    {/each}
  </div>
{:else}
  <div class="skeleton {className}" {style}></div>
{/if}

<style>
  .skeleton {
    background: var(--md-sys-color-surface-variant);
    animation: pulse 1.5s infinite ease-in-out;
    display: inline-block;
    vertical-align: middle;
  }

  .skeleton-card {
    display: block;
    width: 100%;
    height: 260px;
    border-radius: var(--radius, 1.5rem);
  }

  .skeleton-stack {
    display: flex;
    flex-direction: column;
    gap: var(--space-3, 0.75rem);
    width: 100%;
  }

  .skeleton-row {
    display: block;
    height: 3.5rem;
    border-radius: var(--space-3, 0.75rem);
  }

  .skeleton-line {
    display: block;
    height: 1rem;
    border-radius: var(--space-2, 0.5rem);
  }

  .skeleton-line:last-child {
    width: 60%;
  }

  @keyframes pulse {
    0% {
      opacity: 0.3;
    }
    50% {
      opacity: 0.7;
    }
    100% {
      opacity: 0.3;
    }
  }
</style>
