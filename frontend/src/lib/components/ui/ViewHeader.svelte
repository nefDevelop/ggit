<script lang="ts">
  import GlobalSearch from '../GlobalSearch.svelte';

  interface Tab {
    id: string;
    label: string;
    icon: string;
  }

  let { 
    tabs = [], 
    activeTab = $bindable(''), 
    showSearch = true,
    actions
  } = $props();

</script>

<header class="content-controls view-header-abstract">
  <div class="search-and-action-row">
    {#if showSearch}
      <div class="search-flex">
        <GlobalSearch />
      </div>
    {/if}
    
    {#if actions}
      <div class="actions-flex">
        {@render actions()}
      </div>
    {/if}
  </div>

  {#if tabs.length > 0}
    <div class="view-mode-selector-row">
      <nav class="view-mode-selector">
        {#each tabs as tab}
          <button 
            class="tab-pill" 
            class:active={activeTab === tab.id}
            onclick={() => activeTab = tab.id}
          >
            <span class="material-symbols-rounded mini">{tab.icon}</span>
            {tab.label}
          </button>
        {/each}
      </nav>
    </div>
  {/if}
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

  .view-mode-selector { display: flex; gap: 0.5rem; }
  
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
  
  @media (max-width: 800px) {
    .search-and-action-row { flex-wrap: wrap; }
    .search-flex { min-width: 100%; order: 1; }
    .actions-flex { order: 0; margin-left: auto; }
  }
</style>
