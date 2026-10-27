<script lang="ts">
  interface Option {
    value: string | number;
    label: string;
  }

  interface Props {
    value: string | number;
    options: Option[];
    label?: string;
    hint?: string;
    error?: string;
    class?: string;
    onchange?: (value: any) => void;
  }

  let {
    value = $bindable(),
    options,
    label,
    hint,
    error,
    class: className = '',
    onchange
  }: Props = $props();

  const id = $props.id();

  function handleChange(e: Event) {
    const target = e.target as HTMLSelectElement;
    value = target.value;
    if (onchange) onchange(value);
  }
  import ChevronDown from "lucide-svelte/icons/chevron-down";
</script>

<div class="select-container {className}">
  {#if label}
    <label for={id}>{label}</label>
  {/if}
  <div class="select-wrapper" class:has-error={!!error}>
    <select id={id} {value} onchange={handleChange} aria-invalid={error ? 'true' : undefined}>
      {#each options as option}
        <option value={option.value}>{option.label}</option>
      {/each}
    </select>
    <ChevronDown class="arrow" size={18} />
  </div>
  {#if error}
    <span class="error-text">{error}</span>
  {:else if hint}
    <span class="hint-text">{hint}</span>
  {/if}
</div>

<style>
  .select-container {
    display: flex;
    flex-direction: column;
    gap: var(--space-2, 0.5rem);
    width: 100%;
  }

  label {
    font-size: 0.9rem;
    font-weight: 600;
    color: var(--md-sys-color-on-surface-variant);
    margin-left: var(--space-1, 0.25rem);
  }

  .select-wrapper {
    position: relative;
    width: 100%;
  }

  select {
    width: 100%;
    appearance: none;
    background: var(--md-sys-color-surface-container-high);
    border: 1px solid var(--md-sys-color-outline-variant);
    padding: 0.75rem 2.75rem 0.75rem 1rem;
    border-radius: 0.75rem;
    color: var(--md-sys-color-on-surface);
    font-family: inherit;
    font-size: 1rem;
    cursor: pointer;
    transition: all 0.2s;
    outline: none;
    box-sizing: border-box;
  }

  select:focus {
    border-color: var(--md-sys-color-primary);
    background: var(--md-sys-color-surface-container-highest);
    box-shadow: 0 0 0 3px color-mix(in srgb, var(--md-sys-color-primary), transparent 85%);
  }

  select:disabled {
    opacity: 0.5;
    cursor: not-allowed;
    background: var(--md-sys-color-surface-container-low);
  }

  .has-error select {
    border-color: var(--md-sys-color-error);
  }

  .has-error select:focus {
    box-shadow: 0 0 0 3px color-mix(in srgb, var(--md-sys-color-error), transparent 85%);
  }

  .arrow {
    position: absolute;
    right: 1rem;
    top: 50%;
    transform: translateY(-50%);
    pointer-events: none;
    color: var(--md-sys-color-on-surface-variant);
    transition: transform 0.2s;
  }

  select:focus + :global(.arrow) {
    transform: translateY(-50%) rotate(180deg);
    color: var(--md-sys-color-primary);
  }

  option {
    background: var(--md-sys-color-surface-container-highest);
    color: var(--md-sys-color-on-surface);
  }

  .hint-text,
  .error-text {
    font-size: 0.8rem;
    margin-left: var(--space-1, 0.25rem);
  }

  .hint-text {
    color: var(--md-sys-color-on-surface-variant);
  }

  .error-text {
    color: var(--md-sys-color-error);
  }
</style>
