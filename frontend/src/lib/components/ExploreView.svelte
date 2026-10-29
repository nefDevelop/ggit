<script lang="ts">
  import { onMount } from 'svelte';
  import ExploreRepos from './ExploreRepos.svelte';
  import ExploreUsers from './ExploreUsers.svelte';
  import ExploreOrganizations from './ExploreOrganizations.svelte';
  import ExploreTopics from './ExploreTopics.svelte';
  import GlobalSearch from './GlobalSearch.svelte';
  import ViewHeader from './ui/ViewHeader.svelte';
  import PageHeader from './ui/PageHeader.svelte';
  

import Compass from "lucide-svelte/icons/compass";
import GitBranch from "lucide-svelte/icons/git-branch";
import Users from "lucide-svelte/icons/users";
import Building2 from "lucide-svelte/icons/building-2";
import Hash from "lucide-svelte/icons/hash";



  let { defaultTab = 'repos' } = $props();
  let activeTab = $state('repos');
  $effect(() => { if (defaultTab) activeTab = defaultTab; });
  let searchQuery = $state('');

  const tabs = [
    { id: 'repos', label: 'Repositories', icon: GitBranch },
    { id: 'users', label: 'Users', icon: Users },
    { id: 'organizations', label: 'Organizations', icon: Building2 },
    { id: 'topics', label: 'Topics', icon: Hash }
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
        icon={Compass}
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

  .explore-content { padding: 0 var(--space-5) var(--space-7); animation: fadeIn 0.4s cubic-bezier(0.4, 0, 0.2, 1); }
  @keyframes fadeIn { from { opacity: 0; transform: translateY(4px); } to { opacity: 1; transform: translateY(0); } }

  @media (max-width: 600px) {
    .explore-content { padding: 0 var(--space-4) var(--space-6); }
  }
</style>
