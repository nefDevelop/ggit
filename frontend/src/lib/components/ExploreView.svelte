<script lang="ts">
  import ExploreRepos from './ExploreRepos.svelte';
  import ExploreUsers from './ExploreUsers.svelte';
  import ExploreOrganizations from './ExploreOrganizations.svelte';

  let activeTab = $state('repos');
  let searchQuery = $state('');

  const tabs = [
    { id: 'repos', label: 'Repositories', icon: 'account_tree' },
    { id: 'users', label: 'Users', icon: 'person_search' },
    { id: 'organizations', label: 'Organizations', icon: 'corporate_fare' }
  ];
</script>

<div class="explore-container">
  <header class="explore-header">
    <div class="header-content">
       <h1>Explore Ggit</h1>
       <p class="subtitle">Discover interesting projects and connect with other developers.</p>
    </div>
    <div class="search-island island tonal">
      <span class="material-symbols-rounded search-icon">search</span>
      <input 
        type="text" 
        bind:value={searchQuery} 
        placeholder="Search repositories, users or organizations..." 
        oninput={() => {}} 
      />
    </div>
  </header>

  <nav class="explore-nav">
    {#each tabs as tab}
      <button 
        class="explore-tab" 
        class:active={activeTab === tab.id}
        onclick={() => activeTab = tab.id}
      >
        <span class="material-symbols-rounded">{tab.icon}</span>
        {tab.label}
      </button>
    {/each}
  </nav>

  <main class="explore-content">
    {#if activeTab === 'repos'}
      <ExploreRepos q={searchQuery} />
    {:else if activeTab === 'users'}
      <ExploreUsers q={searchQuery} />
    {:else if activeTab === 'organizations'}
      <ExploreOrganizations q={searchQuery} />
    {/if}
  </main>
</div>

<style>
  .explore-container { display: flex; flex-direction: column; gap: 2rem; }
  
  .explore-header { display: flex; flex-direction: column; align-items: center; text-align: center; gap: 2rem; padding: 2rem 0; }
  .header-content h1 { font-size: 2.5rem; font-weight: 800; margin: 0; background: linear-gradient(135deg, var(--primary-color), #8b5cf6); -webkit-background-clip: text; -webkit-text-fill-color: transparent; }
  .subtitle { color: var(--text-muted); font-size: 1.1rem; max-width: 500px; margin: 0.5rem auto 0; }

  .search-island { display: flex; align-items: center; gap: 1rem; padding: 0.75rem 1.5rem; width: 100%; max-width: 600px; border-radius: 2rem; border: 1px solid var(--md-sys-color-outline-variant); }
  .search-island input { flex: 1; background: transparent; border: none; font-size: 1.1rem; color: var(--text-color); font-family: inherit; }
  .search-island input:focus { outline: none; }
  .search-icon { color: var(--text-muted); font-size: 1.5rem; }

  .explore-nav { display: flex; justify-content: center; gap: 1rem; border-bottom: 1px solid var(--md-sys-color-outline-variant); }
  .explore-tab { background: transparent; border: none; padding: 1rem 1.5rem; color: var(--text-muted); font-weight: 700; font-size: 1.05rem; cursor: pointer; display: flex; align-items: center; gap: 0.5rem; position: relative; transition: all 0.2s; }
  .explore-tab:hover { color: var(--text-color); }
  .explore-tab.active { color: var(--primary-color); }
  .explore-tab.active::after { content: ''; position: absolute; bottom: 0; left: 0; right: 0; height: 3px; background: var(--primary-color); border-top-left-radius: 3px; border-top-right-radius: 3px; }

  .explore-content { animation: slideUp 0.4s ease-out; }
  @keyframes slideUp { from { opacity: 0; transform: translateY(20px); } to { opacity: 1; transform: translateY(0); } }
</style>
