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

<div class="dashboard-grid full-width">
  <div class="dashboard-main-content">
    <PageHeader 
        title="Explorar"
        subtitle="Descubre nuevos repositorios, usuarios y comunidades"
        icon="explore"
    />
    <ViewHeader {tabs} bind:activeTab />

    <main class="explore-content">
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
  .dashboard-grid { 
    display: flex; 
    flex-direction: column; 
    height: 100%; 
  }

  .dashboard-main-content { 
    display: flex; 
    flex-direction: column; 
  }

  .explore-content { padding: 0 1.5rem 4rem; animation: fadeIn 0.4s cubic-bezier(0.4, 0, 0.2, 1); }
  @keyframes fadeIn { from { opacity: 0; transform: translateY(4px); } to { opacity: 1; transform: translateY(0); } }
</style>
