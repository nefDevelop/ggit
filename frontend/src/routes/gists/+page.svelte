<script lang="ts">
  import { apiFetch } from "$lib/api";
  import { onMount } from "svelte";
  import { authState } from "$lib/auth.svelte";
  import GistRow from "$lib/components/GistRow.svelte";
  import PageHeader from "$lib/components/ui/PageHeader.svelte";
  import ViewHeader from "$lib/components/ui/ViewHeader.svelte";
  import { ScrollText, Plus, User, Search } from "lucide-svelte";

  let myGists = $state<any[]>([]);
  let publicGists = $state<any[]>([]);
  let loading = $state(true);
  let activeTab = $state("mine");

  async function fetchData() {
    loading = true;
    try {
      const [pubRes, myRes] = await Promise.all([
        apiFetch("/api/v2/gists"),
        authState.isAuthenticated 
          ? apiFetch(`/api/v2/users/${authState.user?.username}/gists`)
          : Promise.resolve({ ok: false } as any)
      ]);

      if (pubRes.ok) publicGists = await pubRes.json();
      if (myRes.ok) myGists = await myRes.json();
    } catch (e) {
      console.error("Failed to fetch gists", e);
    } finally {
      loading = false;
    }
  }

  onMount(fetchData);

  const displayGists = $derived(activeTab === "mine" ? myGists : publicGists);
</script>

<div class="dashboard-grid full-width">
  <div class="dashboard-main-content">
    <PageHeader 
      title="Gists" 
      subtitle="Snippet management and public code sharing"
      icon={ScrollText}
    >
      {#snippet actions()}
        <a href="/gists/new" class="btn-action-primary">
          <Plus size={18} />
          <span class="hide-mobile">New Gist</span>
        </a>
      {/snippet}
    </PageHeader>
  
    <ViewHeader 
      bind:activeTab 
      noPills={true}
      tabs={[
        { id: "mine", label: "My Gists", icon: User },
        { id: "public", label: "Explore", icon: Search }
      ]}
    />
  
    <main class="gists-content">
      <div class="gist-list">
      {#if loading}
        <div class="skeleton-list">
          {#each Array(5) as _}
            <div class="skeleton-row"></div>
          {/each}
        </div>
      {:else if displayGists.length === 0}
        <div class="empty-state">
          <ScrollText size={48} />
          <h3>No gists found</h3>
          <p>{activeTab === 'mine' ? 'You haven\'t created any gists yet.' : 'No public gists available.'}</p>
        </div>
      {:else}
        {#each displayGists as gist}
          <GistRow {gist} />
        {/each}
      {/if}
    </div>
  </main>
</div>
</div>

<style>
  .dashboard-grid { display: flex; flex-direction: column; height: 100%; }
  .dashboard-main-content { display: flex; flex-direction: column; }
  .gists-content { padding: 0 1.5rem 4rem; animation: fadeIn 0.4s cubic-bezier(0.4, 0, 0.2, 1); }
  @keyframes fadeIn { from { opacity: 0; transform: translateY(4px); } to { opacity: 1; transform: translateY(0); } }

  .gist-list { display: flex; flex-direction: column; gap: 1rem; }

  .skeleton-list { display: flex; flex-direction: column; gap: 1rem; }
  .skeleton-row { height: 80px; background: var(--md-sys-color-surface-container-low); border-radius: 1.25rem; animation: pulse 1.5s infinite; }
  @keyframes pulse { 0%, 100% { opacity: 0.5; } 50% { opacity: 0.8; } }

  .empty-state { text-align: center; padding: 4rem 2rem; color: var(--text-muted); display: flex; flex-direction: column; align-items: center; gap: 1rem; }
  .empty-state h3 { margin: 0; color: var(--text-color); }

  @media (max-width: 600px) {
    .hide-mobile { display: none; }
  }
</style>
