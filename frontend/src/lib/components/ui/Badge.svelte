<script lang="ts">
  let { 
    variant = 'neutral', 
    size = 'md',
    class: className = '',
    children,
    icon
  } = $props();

  const variants = {
    neutral: 'badge-neutral',
    primary: 'badge-primary',
    secondary: 'badge-secondary',
    tertiary: 'badge-tertiary',
    success: 'badge-success',
    danger: 'badge-danger',
    warning: 'badge-warning',
    ghost: 'badge-ghost'
  };

  const sizes = {
    sm: 'badge-sm',
    md: 'badge-md'
  };
</script>

<span class="badge {variants[variant]} {sizes[size]} {className}">
  {#if icon}
    <span class="material-symbols-rounded icon">{icon}</span>
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
    white-space: nowrap;
    line-height: 1;
    border: 1px solid transparent;
    user-select: none;
  }

  .badge-sm { padding: 0.15rem 0.5rem; font-size: 0.7rem; }
  .badge-md { padding: 0.25rem 0.75rem; font-size: 0.75rem; }

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

