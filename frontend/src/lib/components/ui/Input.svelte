<script lang="ts">
  let { 
    value = $bindable(), 
    placeholder = '', 
    type = 'text', 
    label = '', 
    error = '', 
    disabled = false,
    class: className = '',
    multiline = false,
    rows = 3,
    icon = '',
    min = undefined,
    max = undefined,
    step = undefined
  } = $props();
</script>

<div class="field {className}">
  {#if label}
    <label class="label">{label}</label>
  {/if}
  
  <div class="input-container {error ? 'has-error' : ''} {icon ? 'with-icon' : ''}">
    {#if icon}
      <span class="material-symbols-rounded input-icon">{icon}</span>
    {/if}
    
    {#if multiline}
      <textarea 
        bind:value 
        {placeholder} 
        {disabled} 
        {rows}
      ></textarea>
    {:else}
      <input 
        {type} 
        bind:value 
        {placeholder} 
        {disabled} 
        {min}
        {max}
        {step}
      />
    {/if}
  </div>

  {#if error}
    <span class="error-text">{error}</span>
  {/if}
</div>

<style>
  .field {
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
    width: 100%;
  }

  .label {
    font-size: 0.9rem;
    font-weight: 600;
    color: var(--md-sys-color-on-surface-variant);
    margin-left: 0.25rem;
  }

  .input-container {
    position: relative;
    width: 100%;
    display: flex;
    align-items: center;
  }

  .input-icon {
    position: absolute;
    left: 1rem;
    color: var(--md-sys-color-on-surface-variant);
    font-size: 1.25rem;
    pointer-events: none;
    z-index: 1;
  }

  .with-icon input, .with-icon textarea {
    padding-left: 2.75rem;
  }

  input, textarea {
    width: 100%;
    background: var(--md-sys-color-surface-container-high);
    border: 1px solid var(--md-sys-color-outline-variant);
    color: var(--md-sys-color-on-surface);
    padding: 0.75rem 1rem;
    border-radius: 0.75rem;
    font-size: 1rem;
    font-family: inherit;
    transition: all 0.2s;
    outline: none;
    box-sizing: border-box;
  }

  input:focus, textarea:focus {
    border-color: var(--md-sys-color-primary);
    background: var(--md-sys-color-surface-container-highest);
    box-shadow: 0 0 0 3px rgba(var(--primary-rgb), 0.15);
    box-shadow: 0 0 0 3px color-mix(in srgb, var(--md-sys-color-primary), transparent 85%);
  }

  input:disabled, textarea:disabled {
    opacity: 0.5;
    cursor: not-allowed;
    background: var(--md-sys-color-surface-container-low);
  }

  .has-error input, .has-error textarea {
    border-color: var(--md-sys-color-error);
  }

  .has-error input:focus, .has-error textarea:focus {
    box-shadow: 0 0 0 3px color-mix(in srgb, var(--md-sys-color-error), transparent 85%);
  }

  .error-text {
    font-size: 0.8rem;
    color: var(--md-sys-color-error);
    font-weight: 500;
    margin-left: 0.25rem;
  }

  textarea {
    resize: vertical;
    min-height: 80px;
  }

  ::placeholder {
    color: var(--md-sys-color-on-surface-variant);
    opacity: 0.5;
  }
</style>

