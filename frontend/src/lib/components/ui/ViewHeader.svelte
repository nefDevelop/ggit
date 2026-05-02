<script lang="ts">
  import GlobalSearch from '../GlobalSearch.svelte';
  import type { Snippet } from 'svelte';

  interface Tab {
    id: string;
    label: string;
    icon: string;
    hasBadge?: boolean;
  }

  interface Props {
    tabs?: Tab[];
    activeTab?: string;
    showSearch?: boolean;
    actions?: Snippet;
    onSelect?: (id: string) => void;
  }

  let { 
    tabs = [], 
    activeTab = $bindable(''), 
    showSearch = true,
    actions,
    onSelect
  }: Props = $props();

</script>

<header class="content-controls view-header-abstract">
  <div class="search-and-action-row">
    {#if showSearch}
      <div class="search-flex">
        <GlobalSearch />
      </div>
    {/if}

    {#if tabs.length > 0}
      <div class="view-mode-selector-row">
        <nav class="view-mode-selector">
          {#each tabs as tab}
            <button 
              class="tab-pill" 
              class:active={activeTab === tab.id}
              onclick={() => { activeTab = tab.id; onSelect?.(tab.id); }}
            >
              <div class="tab-content-wrapper">
                <span class="material-symbols-rounded mini">{tab.icon}</span>
                <span class="tab-label-text">{tab.label}</span>
                {#if tab.hasBadge}
                  <span class="tab-badge-dot"></span>
                {/if}
              </div>
            </button>
          {/each}
        </nav>
      </div>
    {/if}
    
    {#if actions}
      <div class="actions-flex">
        {@render actions()}
      </div>
    {/if}
  </div>
</header>

<style>
  .view-header-abstract {
    display: flex;
    flex-direction: column;
    gap: 1.25rem;
    padding: 0 1.5rem 2rem;
  }

  .search-and-action-row {
    display: flex;
    align-items: center;
    gap: 1.5rem;
    padding: 1.5rem 0 0.5rem;
    width: 100%;
  }

  .search-flex {
    flex: 1;
    max-width: 600px;
  }

  .actions-flex {
    display: flex;
    gap: 1rem;
    align-items: center;
  }

  .view-mode-selector-row {
    display: flex;
    align-items: center;
  }

  .view-mode-selector { 
    display: flex; 
    gap: 0.5rem; 
    overflow-x: auto; 
    scrollbar-width: none; 
    -webkit-overflow-scrolling: touch;
    width: 100%;
  }
  .view-mode-selector::-webkit-scrollbar { display: none; }

  .tab-pill { 
    background: transparent; 
    border: 1px solid var(--md-sys-color-outline-variant); 
    padding: 0.5rem 1rem; 
    border-radius: 2rem; 
    color: var(--text-muted); 
    font-weight: 700; 
    font-size: 0.85rem; 
    cursor: pointer; 
    display: flex; 
    align-items: center; 
    gap: 0.5rem; 
    transition: all 0.2s;
    white-space: nowrap;
    font-family: inherit;
  }
  .tab-pill:hover { background: var(--md-sys-color-surface-container-low); color: var(--text-color); }
  .tab-pill.active { background: var(--primary-color); color: white; border-color: var(--primary-color); }
  
  .tab-content-wrapper {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    position: relative;
  }

  .tab-badge-dot {
    width: 6px;
    height: 6px;
    background: #ff5252;
    border-radius: 50%;
    position: absolute;
    top: -2px;
    right: -8px;
    box-shadow: 0 0 8px #ff5252;
    animation: pulse-badge 2s infinite;
  }

  @keyframes pulse-badge {
    0% { transform: scale(1); opacity: 1; }
    50% { transform: scale(1.5); opacity: 0.5; }
    100% { transform: scale(1); opacity: 1; }
  }

  @media (max-width: 800px) {
    .view-header-abstract { padding: 0 1rem 1.5rem; }
    .search-and-action-row { flex-wrap: nowrap; gap: 0.75rem; }
    .search-flex { min-width: 0; flex: 1; order: 0; }
    .actions-flex { order: 1; width: auto; justify-content: flex-end; }
    .tab-pill:not(.active) .tab-label-text { display: none; }
    .tab-pill:not(.active) { padding: 0.5rem 0.65rem; border-radius: 100px; }
    .tab-pill.active { padding: 0.5rem 1rem; border-radius: 100px; }
  }
</style>
