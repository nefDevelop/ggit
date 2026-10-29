<script lang="ts">
  import type { Snippet } from 'svelte';

  interface Props {
    variant?: 'neutral' | 'primary' | 'secondary' | 'tertiary' | 'success' | 'danger' | 'warning' | 'ghost';
    size?: 'sm' | 'md';
    class?: string;
    children?: Snippet;
    icon?: string | any;
  }

  let { 
    variant = 'neutral', 
    size = 'md',
    class: className = '',
    children,
    icon
  }: Props = $props();

  const variants: Record<NonNullable<Props['variant']>, string> = {
    neutral: 'badge-neutral',
    primary: 'badge-primary',
    secondary: 'badge-secondary',
    tertiary: 'badge-tertiary',
    success: 'badge-success',
    danger: 'badge-danger',
    warning: 'badge-warning',
    ghost: 'badge-ghost'
  };

  const sizes: Record<NonNullable<Props['size']>, string> = {
    sm: 'badge-sm',
    md: 'badge-md'
  };
</script>

<span class="badge {variants[variant]} {sizes[size]} {className}">
  {#if icon}
    {@const Icon = icon}
    <Icon size={size === 'sm' ? 12 : 14} class="icon" />
  {/if}
  {#if children}
    {@render children()}
  {/if}
</span>

<style>
  .badge {
    display: inline-flex;
    align-items: center;
    gap: 0.35rem;
    padding: 0.25rem 0.75rem;
    border-radius: 2rem;
    font-size: 0.75rem;
    font-weight: 800;
    line-height: 1;
    border: 1px solid transparent;
    user-select: none;
    max-width: 100%;
  }

  .badge-sm { padding: 0.15rem 0.5rem; font-size: 0.7rem; }
  .badge-md { padding: 0.25rem 0.75rem; font-size: 0.75rem; }

  .badge :global(span) {
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }

  .badge-neutral {
    background: var(--md-sys-color-surface-variant);
    color: var(--md-sys-color-on-surface-variant);
  }

  .badge-primary {
    background: var(--md-sys-color-primary-container);
    color: var(--md-sys-color-on-primary-container);
  }

  .badge-secondary {
    background: var(--md-sys-color-secondary-container);
    color: var(--md-sys-color-on-secondary-container);
  }

  .badge-tertiary {
    background: var(--md-sys-color-tertiary-container);
    color: var(--md-sys-color-on-tertiary-container);
  }

  .badge-success {
    background: var(--md-sys-color-secondary-container);
    color: var(--color-success);
    border-color: color-mix(in srgb, var(--color-success), transparent 70%);
  }

  .badge-danger {
    background: var(--md-sys-color-error-container);
    color: var(--md-sys-color-on-error-container);
  }

  .badge-warning {
    background: var(--bg-warning-subtle);
    color: var(--color-warning);
    border-color: color-mix(in srgb, var(--color-warning), transparent 70%);
  }

  .badge-ghost {
    background: transparent;
    border-color: var(--md-sys-color-outline-variant);
    color: var(--md-sys-color-on-surface-variant);
  }

  .icon { font-size: 1rem; }
</style>

