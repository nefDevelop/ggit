<script lang="ts">
  interface Props {
    checked: boolean;
    label: string;
    description?: string;
    onchange?: (checked: boolean) => void;
    class?: string;
  }

  let { checked = $bindable(), label, description, onchange, class: className }: Props = $props();

  function toggle() {
    checked = !checked;
    if (onchange) onchange(checked);
  }
</script>

<label class="checkbox-container {className}">
  <input 
    type="checkbox" 
    bind:checked={checked} 
    onchange={() => onchange?.(checked)}
  />
  <span class="checkmark" class:checked>
    {#if checked}
      <span class="material-symbols-rounded">check</span>
    {/if}
  </span>
  <div class="checkbox-text">
    <strong>{label}</strong>
    {#if description}
      <p>{description}</p>
    {/if}
  </div>
</label>

<style>
  .checkbox-container {
    display: flex;
    align-items: flex-start;
    gap: 1rem;
    cursor: pointer;
    position: relative;
    padding: 1rem;
    background: var(--md-sys-color-surface-container-low);
    border-radius: 1.25rem;
    transition: all 0.2s;
    border: 1px solid transparent;
    user-select: none;
  }
  .checkbox-container:hover {
    background: var(--md-sys-color-surface-container-high);
    border-color: var(--md-sys-color-outline-variant);
  }
  .checkbox-container:active {
    transform: scale(0.98);
    background: var(--md-sys-color-surface-container-highest);
  }

  input {
    opacity: 0;
    position: absolute;
    cursor: pointer;
  }

  .checkmark {
    width: 1.5rem;
    height: 1.5rem;
    background: var(--md-sys-color-surface);
    border: 2px solid var(--md-sys-color-outline);
    border-radius: 0.5rem;
    flex-shrink: 0;
    display: flex;
    align-items: center;
    justify-content: center;
    transition: all 0.2s;
    margin-top: 0.1rem;
  }

  .checkmark.checked {
    background: var(--md-sys-color-primary);
    border-color: var(--md-sys-color-primary);
  }

  .checkmark span {
    color: var(--md-sys-color-on-primary);
    font-size: 1.1rem;
  }

  .checkbox-text strong {
    display: block;
    font-size: 1rem;
    font-weight: 700;
    color: var(--md-sys-color-on-surface);
  }
  .checkbox-text p {
    font-size: 0.85rem;
    color: var(--md-sys-color-on-surface-variant);
    margin: 0.25rem 0 0;
    line-height: 1.3;
  }
</style>
