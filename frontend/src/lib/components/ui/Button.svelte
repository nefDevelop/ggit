<script lang="ts">
  import type { Snippet } from 'svelte';

  const variants = {
    primary: 'btn-primary',
    secondary: 'btn-secondary',
    tertiary: 'btn-tertiary',
    danger: 'btn-danger',
    ghost: 'btn-ghost',
    outline: 'btn-outline',
    tonal: 'btn-tonal'
  };

  const sizes = {
    sm: 'btn-sm',
    md: 'btn-md',
    lg: 'btn-lg'
  };

  interface Props {
    variant?: keyof typeof variants;
    size?: keyof typeof sizes;
    type?: 'button' | 'submit' | 'reset';
    loading?: boolean;
    disabled?: boolean;
    class?: string;
    onclick?: (e: MouseEvent) => void;
    children?: Snippet;
    icon?: any;
    href?: string;
    title?: string;
    "aria-label"?: string;
  }

  let { 
    variant = 'primary', 
    size = 'md', 
    type = 'button', 
    loading = false,
    disabled = false, 
    class: className = '',
    onclick,
    children,
    icon,
    href,
    title,
    "aria-label": ariaLabel
  }: Props = $props();
</script>

{#if href}
  <a 
    {href}
    {title}
    aria-label={ariaLabel}
    class="btn {variants[variant]} {sizes[size]} {className}" 
    onclick={(e) => !disabled && onclick?.(e as any)}
  >
    {#if typeof icon === 'string'}
      <span class="material-symbols-rounded">{icon}</span>
    {:else if icon}
      {@const Icon = icon}
      <Icon size={18} />
    {/if}
    {#if children}
      {@render children()}
    {/if}
  </a>
{:else}
  <button 
    {type} 
    disabled={disabled || loading} 
    {title}
    aria-label={ariaLabel}
    class="btn {variants[variant]} {sizes[size]} {className}" 
    onclick={(e) => !loading && onclick?.(e)}
  >
    {#if loading}
      <span class="spinner"></span>
    {:else}
      {#if typeof icon === 'string'}
        <span class="material-symbols-rounded">{icon}</span>
      {:else if icon}
        {@const Icon = icon}
        <Icon size={18} />
      {/if}
      {#if children}
        {@render children()}
      {/if}
    {/if}
  </button>
{/if}

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
  
  .btn-tonal {
    background: var(--md-sys-color-secondary-container);
    color: var(--md-sys-color-on-secondary-container);
  }
  .btn-tonal:hover:not(:disabled) {
    background: color-mix(in srgb, var(--md-sys-color-secondary-container), var(--md-sys-color-on-secondary-container) 8%);
  }

  .material-symbols-rounded { 
    font-size: 1.25rem; 
    font-variation-settings: 'FILL' 0, 'wght' 400, 'GRAD' 0, 'opsz' 24;
  }

  .spinner {
    width: 1.25rem;
    height: 1.25rem;
    border: 2px solid currentColor;
    border-top-color: transparent;
    border-radius: 50%;
    animation: spin 0.8s linear infinite;
  }

  @keyframes spin {
    from { transform: rotate(0deg); }
    to { transform: rotate(360deg); }
  }
</style>
