<script lang="ts">
  interface Option {
    value: string | number;
    label: string;
  }

  interface Props {
    value: string | number;
    options: Option[];
    label?: string;
    onchange?: (value: any) => void;
  }

  let { value = $bindable(), options, label, onchange }: Props = $props();

  function handleChange(e: Event) {
    const target = e.target as HTMLSelectElement;
    value = target.value;
    if (onchange) onchange(value);
  }
</script>

<div class="select-container">
  {#if label}
    <label for="select-input">{label}</label>
  {/if}
  <div class="select-wrapper">
    <select id="select-input" {value} onchange={handleChange}>
      {#each options as option}
        <option value={option.value}>{option.label}</option>
      {/each}
    </select>
    <span class="material-symbols-rounded arrow">expand_more</span>
  </div>
</div>

<style>
  .select-container {
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
    width: 100%;
  }

  label {
    font-size: 0.85rem;
    font-weight: 800;
    text-transform: uppercase;
    color: var(--text-muted);
    padding-left: 0.25rem;
  }

  .select-wrapper {
    position: relative;
    width: 100%;
  }

  select {
    width: 100%;
    appearance: none;
    background: var(--md-sys-color-surface-container-low);
    border: 1px solid var(--md-sys-color-outline-variant);
    padding: 0.75rem 3rem 0.75rem 1rem;
    border-radius: 0.75rem;
    color: var(--text-color);
    font-family: inherit;
    font-size: 1rem;
    cursor: pointer;
    transition: all 0.2s;
  }

  select:focus {
    outline: none;
    border-color: var(--primary-color);
    background: var(--md-sys-color-surface-container-high);
  }

  .arrow {
    position: absolute;
    right: 1rem;
    top: 50%;
    transform: translateY(-50%);
    pointer-events: none;
    color: var(--text-muted);
    transition: transform 0.2s;
  }

  select:focus + .arrow {
    transform: translateY(-50%) rotate(180deg);
    color: var(--primary-color);
  }

  option {
    background: var(--md-sys-color-surface-container-highest);
    color: var(--text-color);
  }
</style>
