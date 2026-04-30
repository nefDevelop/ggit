<script lang="ts">
  /**
   * ChipNavigation.svelte
   * Componente unificado para navegación por "píldoras" (chips) con scroll horizontal.
   * Soporta tanto navegación por rutas (enlaces) como cambio de pestañas interno (botones).
   */
  interface Tab {
    id: string;
    label: string;
    icon: string;
    path?: string;
  }

  let { 
    items, 
    activeId, 
    title = '',
    titleIcon = '',
    onSelect,
    class: className = ""
  }: {
    items: Tab[];
    activeId: string;
    title?: string;
    titleIcon?: string;
    onSelect?: (id: string) => void;
    class?: string;
  } = $props();
</script>

<div class="chip-nav-container {className}">
  {#if title}
    <div class="chip-nav-header">
      {#if titleIcon}
        <span class="material-symbols-rounded">{titleIcon}</span>
      {/if}
      <h3 class="chip-nav-title">{title}</h3>
    </div>
  {/if}

  <nav class="chip-nav-rail">
    {#each items as item}
      {#if item.path}
        <a 
          href={item.path} 
          class="chip-item" 
          class:active={activeId === item.id || activeId === item.path}
        >
          <span class="material-symbols-rounded">{item.icon}</span>
          <span class="label">{item.label}</span>
        </a>
      {:else}
        <button 
          type="button"
          class="chip-item" 
          class:active={activeId === item.id}
          onclick={() => onSelect?.(item.id)}
        >
          <span class="material-symbols-rounded">{item.icon}</span>
          <span class="label">{item.label}</span>
        </button>
      {/if}
    {/each}
  </nav>
</div>

<style>
  .chip-nav-container {
    display: flex;
    flex-direction: column;
    gap: 1.25rem;
    width: 100%;
    box-sizing: border-box;
  }

  .chip-nav-header {
    display: flex;
    align-items: center;
    gap: 1rem;
    padding-left: 0.5rem;
  }

  .chip-nav-header span {
    font-size: 1.75rem;
    color: var(--primary-color);
  }

  .chip-nav-title {
    font-size: 1.1rem;
    font-weight: 900;
    margin: 0;
    letter-spacing: -0.02em;
    color: var(--text-color);
  }

  .chip-nav-rail {
    display: flex;
    gap: 0.5rem;
    overflow-x: auto;
    padding: 0.25rem 0 1rem;
    margin: 0 -1rem;
    padding-left: 1rem;
    scrollbar-width: none;
    -webkit-overflow-scrolling: touch;
  }

  .chip-nav-rail::-webkit-scrollbar {
    display: none;
  }

  .chip-item {
    display: flex;
    align-items: center;
    gap: 0.6rem;
    padding: 0.75rem 1.25rem;
    border-radius: 100px;
    background: rgba(255, 255, 255, 0.04);
    color: var(--text-muted);
    text-decoration: none;
    font-size: 0.9rem;
    font-weight: 700;
    white-space: nowrap;
    border: 1px solid var(--md-sys-color-outline-variant);
    transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
    cursor: pointer;
    flex-shrink: 0;
    font-family: inherit;
  }

  .chip-item:hover {
    background: rgba(255, 255, 255, 0.08);
    color: var(--text-color);
  }

  .chip-item.active {
    background: var(--primary-color);
    color: white;
    border-color: var(--primary-color);
    box-shadow: 0 4px 12px rgba(var(--primary-rgb), 0.2);
  }

  .chip-item span.material-symbols-rounded {
    font-size: 1.25rem;
  }

  @media (max-width: 1024px) {
    .chip-nav-header {
      display: none;
    }
    .chip-nav-container {
      gap: 0;
    }
    .chip-nav-rail {
      padding-bottom: 0.5rem;
    }
  }
</style>
