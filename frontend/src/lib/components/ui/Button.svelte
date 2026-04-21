<script lang="ts">
  let { 
    variant = 'primary', 
    size = 'md', 
    type = 'button', 
    disabled = false, 
    class: className = '',
    onclick,
    children,
    icon
  } = $props();

  const variants = {
    primary: 'btn-primary',
    secondary: 'btn-secondary',
    tertiary: 'btn-tertiary',
    danger: 'btn-danger',
    ghost: 'btn-ghost',
    outline: 'btn-outline'
  };

  const sizes = {
    sm: 'btn-sm',
    md: 'btn-md',
    lg: 'btn-lg'
  };
</script>

<button 
  {type} 
  {disabled} 
  class="btn {variants[variant]} {sizes[size]} {className}" 
  {onclick}
>
  {#if icon}
    <span class="material-symbols-rounded">{icon}</span>
  {/if}
  {#if children}
    {@render children()}
  {/if}
</button>

<style>
  .btn {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    gap: 0.5rem;
    border-radius: 0.75rem;
    font-weight: 600;
    cursor: pointer;
    transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1);
    border: 1px solid transparent;
    font-family: inherit;
    white-space: nowrap;
    text-decoration: none;
    user-select: none;
  }

  .btn:disabled {
    opacity: 0.5;
    cursor: not-allowed;
    filter: grayscale(0.5);
  }

  .btn-sm { padding: 0.4rem 0.8rem; font-size: 0.85rem; border-radius: 0.5rem; }
  .btn-md { padding: 0.6rem 1.25rem; font-size: 0.95rem; }
  .btn-lg { padding: 0.8rem 1.75rem; font-size: 1.1rem; border-radius: 1rem; }

  .btn-primary {
    background: var(--md-sys-color-primary);
    color: var(--md-sys-color-on-primary);
  }
  .btn-primary:hover:not(:disabled) {
    background: color-mix(in srgb, var(--md-sys-color-primary), white 10%);
    box-shadow: var(--shadow-md);
    transform: translateY(-1px);
  }

  .btn-secondary {
    background: var(--md-sys-color-secondary-container);
    color: var(--md-sys-color-on-secondary-container);
  }
  .btn-secondary:hover:not(:disabled) {
    background: color-mix(in srgb, var(--md-sys-color-secondary-container), var(--md-sys-color-on-secondary-container) 8%);
  }

  .btn-tertiary {
    background: var(--md-sys-color-tertiary-container);
    color: var(--md-sys-color-on-tertiary-container);
  }
  .btn-tertiary:hover:not(:disabled) {
    background: color-mix(in srgb, var(--md-sys-color-tertiary-container), var(--md-sys-color-on-tertiary-container) 8%);
  }

  .btn-danger {
    background: var(--md-sys-color-error);
    color: var(--md-sys-color-on-error);
  }
  .btn-danger:hover:not(:disabled) {
    background: color-mix(in srgb, var(--md-sys-color-error), white 10%);
    box-shadow: var(--shadow-md);
    transform: translateY(-1px);
  }

  .btn-ghost {
    background: transparent;
    color: var(--md-sys-color-on-surface);
  }
  .btn-ghost:hover:not(:disabled) {
    background: var(--md-sys-color-surface-variant);
  }

  .btn-outline {
    background: transparent;
    color: var(--md-sys-color-primary);
    border-color: var(--md-sys-color-outline);
  }
  .btn-outline:hover:not(:disabled) {
    background: var(--md-sys-color-primary-container);
    border-color: var(--md-sys-color-primary);
  }

  .material-symbols-rounded { 
    font-size: 1.25rem; 
    font-variation-settings: 'FILL' 0, 'wght' 400, 'GRAD' 0, 'opsz' 24;
  }
</style>

