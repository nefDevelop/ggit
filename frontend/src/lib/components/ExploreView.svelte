<script lang="ts">
  import { onMount } from 'svelte';
  import ExploreRepos from './ExploreRepos.svelte';
  import ExploreUsers from './ExploreUsers.svelte';
  import ExploreOrganizations from './ExploreOrganizations.svelte';
  import ExploreTopics from './ExploreTopics.svelte';
  import GlobalSearch from './GlobalSearch.svelte';
  import ViewHeader from './ui/ViewHeader.svelte';
  import PageHeader from './ui/PageHeader.svelte';

  let { defaultTab = 'repos' } = $props();
  let activeTab = $state(defaultTab);
  $effect(() => { activeTab = defaultTab; });
  let searchQuery = $state('');

  const tabs = [
    { id: 'repos', label: 'Repositories', icon: 'account_tree' },
    { id: 'users', label: 'Users', icon: 'person_search' },
    { id: 'organizations', label: 'Organizations', icon: 'corporate_fare' },
    { id: 'topics', label: 'Topics', icon: 'tag' }
  ];

  onMount(() => {
    const handleGlobalSearch = (e: any) => { searchQuery = e.detail; };
    window.addEventListener('ggit:search', handleGlobalSearch);
    return () => window.removeEventListener('ggit:search', handleGlobalSearch);
  });
</script>

<div class="explore-page">
  <div class="explore-content-inner">
    <PageHeader 
        title="Explorar"
        subtitle="Descubre nuevos repositorios, usuarios y comunidades"
        icon="explore"
    />
    <ViewHeader {tabs} bind:activeTab />

    <main class="explore-main-stream">
      {#if activeTab === 'repos'}
        <ExploreRepos q={searchQuery} />
      {:else if activeTab === 'users'}
        <ExploreUsers q={searchQuery} />
      {:else if activeTab === 'organizations'}
        <ExploreOrganizations q={searchQuery} />
      {:else if activeTab === 'topics'}
        <ExploreTopics q={searchQuery} />
      {/if}
    </main>
  </div>
</div>

<style>
  .explore-page { display: flex; flex-direction: column; gap: 0; margin: -2rem; min-height: 100%; }
  .explore-content-inner { padding: 0; display: flex; flex-direction: column; }
  .explore-main-stream { padding: 2rem; animation: fadeIn 0.4s cubic-bezier(0.4, 0, 0.2, 1); width: 100%; box-sizing: border-box; }
  
  @keyframes fadeIn { from { opacity: 0; transform: translateY(4px); } to { opacity: 1; transform: translateY(0); } }
</style>
