<script lang="ts">
  import GlobalSearch from '../GlobalSearch.svelte';
  import type { Snippet } from 'svelte';

  interface Tab {
    id: string;
    label: string;
    icon: any;
    hasBadge?: boolean;
  }

  interface Props {
    tabs?: Tab[];
    activeTab?: string;
    showSearch?: boolean;
    noPills?: boolean;
    actions?: Snippet;
    onSelect?: (id: string) => void;
  }

  let { 
    tabs = [], 
    activeTab = $bindable(''), 
    showSearch = true,
    noPills = false,
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
              class:no-pills={noPills}
              class:active={activeTab === tab.id}
              onclick={() => { activeTab = tab.id; onSelect?.(tab.id); }}
            >
              <div class="tab-content-wrapper">
                {#if typeof tab.icon === 'string'}
                  <span class="material-symbols-rounded mini">{tab.icon}</span>
                {:else if tab.icon}
                  {@const Icon = tab.icon}
                  <Icon size={16} />
                {/if}
                <span class="tab-label-text">{tab.label}</span>
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
    min-width: 0;
    overflow-x: auto;
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
  .tab-pill.active { background: var(--primary-color); color: var(--md-sys-color-on-primary); border-color: var(--primary-color); }
  
  .tab-pill.no-pills {
    border: none;
    border-radius: 0.75rem;
    padding: 0.6rem 1.25rem;
    background: transparent;
    color: var(--text-muted);
    opacity: 0.7;
  }
  .tab-pill.no-pills:hover {
    background: var(--md-sys-color-surface-container-low);
    opacity: 1;
    color: var(--text-color);
  }
  .tab-pill.no-pills.active {
    background: var(--md-sys-color-secondary-container);
    color: var(--md-sys-color-on-secondary-container);
    opacity: 1;
    box-shadow: none;
  }
  
  .tab-content-wrapper {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    position: relative;
  }

  @media (max-width: 1024px) {
    .view-header-abstract { padding: 0 1rem 1.5rem; }
    .search-and-action-row {
      display: grid;
      grid-template-columns: 1fr;
      gap: 1rem;
    }
    .search-flex {
      width: 100%;
      max-width: 100%;
    }
    .view-mode-selector-row {
      width: 100%;
      overflow-x: auto;
    }
    .view-mode-selector {
      display: flex;
      flex-wrap: nowrap;
      gap: 0.5rem;
      padding-bottom: 0.25rem;
    }
    .tab-pill {
      flex: 0 0 auto;
    }
    .tab-pill:not(.active) .tab-label-text {
      display: inline;
    }
    .tab-pill:not(.active), .tab-pill.active {
      padding: 0.5rem 1rem;
      border-radius: 2rem;
    }
  }
</style>
