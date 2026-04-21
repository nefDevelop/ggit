<script lang="ts">
  let { branches = [], selected = $bindable(), label = "", icon = "account_tree" } = $props();

  let isOpen = $state(false);
  let filter = $state("");

  let filteredBranches = $derived(
    branches.filter((b: any) => b.name.toLowerCase().includes(filter.toLowerCase()))
  );

  function select(name: string) {
    selected = name;
    isOpen = false;
    filter = "";
  }

  // Close on outside click is tricky without a directive, but we'll use a simple approach
</script>

<div class="branch-selector">
  {#if label}
    <span class="selector-label">{label}</span>
  {/if}
  <div class="selector-control">
    <button class="selector-trigger island island-interactive tonal" onclick={() => isOpen = !isOpen}>
      <span class="material-symbols-rounded">{icon}</span>
      <span class="current-name">{selected || 'Select branch...'}</span>
      <span class="material-symbols-rounded expand">expand_more</span>
    </button>

    {#if isOpen}
      <div class="selector-dropdown island island-interactive shadow-xl">
        <div class="search-box">
          <span class="material-symbols-rounded">search</span>
          <input 
            type="text" 
            placeholder="Filter branches..." 
            bind:value={filter}
            contenteditable="true"
            onkeydown={(e) => { if(e.key === 'Escape') isOpen = false; }}
          />
        </div>
        <div class="branch-list">
          {#each filteredBranches as b}
            <button 
                class="branch-item" 
                class:active={b.name === selected}
                onclick={() => select(b.name)}
            >
              <span class="material-symbols-rounded">check</span>
              <span class="name">{b.name}</span>
            </button>
          {/each}
          {#if filteredBranches.length === 0}
            <div class="empty">No branches found</div>
          {/if}
        </div>
      </div>
      <div class="backdrop" onclick={() => isOpen = false} onkeydown={(e) => { if(e.key === 'Escape') isOpen = false; }} role="button" tabindex="0"></div>
    {/if}
  </div>
</div>

<style>
  .branch-selector { display: flex; flex-direction: column; gap: 0.25rem; position: relative; }
  .selector-label { font-size: 0.7rem; font-weight: 800; text-transform: uppercase; color: var(--text-muted); opacity: 0.7; padding-left: 0.5rem; }
  
  .selector-control { position: relative; z-index: 10; }
  
  .selector-trigger { 
    display: flex; 
    align-items: center; 
    gap: 0.75rem; 
    padding: 0.6rem 1rem; 
    border: 1px solid var(--md-sys-color-outline-variant); 
    border-radius: 0.75rem;
    cursor: pointer;
    background: rgba(255,255,255,0.02);
    min-width: 200px;
    width: 100%;
    transition: all 0.2s;
  }
  .selector-trigger:hover { background: rgba(255,255,255,0.05); border-color: var(--primary-color); }
  
  .current-name { flex: 1; text-align: left; font-weight: 600; font-size: 0.95rem; }
  .expand { font-size: 1.25rem; opacity: 0.5; }
  
  .selector-dropdown { 
    position: absolute; 
    top: calc(100% + 0.5rem); 
    left: 0; 
    width: 300px; 
    max-height: 400px; 
    z-index: 1000; 
    display: flex; 
    flex-direction: column; 
    overflow: hidden; 
    background: var(--md3-surface-2);
    border: 1px solid var(--md-sys-color-outline-variant);
    box-shadow: 0 12px 40px rgba(0,0,0,0.4);
    animation: slideDown 0.2s cubic-bezier(0, 0, 0.2, 1);
  }
  
  @keyframes slideDown { from { opacity: 0; transform: translateY(-10px); } to { opacity: 1; transform: translateY(0); } }
  
  .search-box { display: flex; align-items: center; gap: 0.75rem; padding: 0.75rem 1rem; border-bottom: 1px solid var(--md-sys-color-outline-variant); background: rgba(0,0,0,0.2); }
  .search-box span { font-size: 1.25rem; color: var(--text-muted); }
  .search-box input { background: transparent; border: none; color: #fff; width: 100%; outline: none; font-size: 0.95rem; }
  
  .branch-list { overflow-y: auto; display: flex; flex-direction: column; padding: 0.5rem; }
  .branch-item { 
    display: flex; 
    align-items: center; 
    gap: 0.75rem; 
    padding: 0.75rem; 
    border: none; 
    background: transparent; 
    color: var(--text-muted); 
    cursor: pointer; 
    border-radius: 0.5rem;
    text-align: left;
    transition: all 0.15s;
  }
  .branch-item:hover { background: rgba(255,255,255,0.05); color: #fff; }
  .branch-item .material-symbols-rounded { font-size: 1.1rem; opacity: 0; }
  .branch-item.active { background: rgba(var(--primary-rgb), 0.1); color: var(--primary-color); font-weight: 700; }
  .branch-item.active .material-symbols-rounded { opacity: 1; }
  
  .empty { padding: 2rem; text-align: center; color: var(--text-muted); font-size: 0.9rem; font-style: italic; }
  
  .backdrop { position: fixed; inset: 0; z-index: -1; background: transparent; }
</style>
