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

<div class="gists-home">
  <div class="container">
    <PageHeader 
      title="Gists" 
      subtitle="Snippet management and public code sharing"
      icon={ScrollText}
    >
      {#snippet actions()}
        <a href="/gists/new" class="btn-primary">
          <Plus size={18} />
          <span>New Gist</span>
        </a>
      {/snippet}
    </PageHeader>

    <ViewHeader 
      bind:activeTab 
      tabs={[
        { id: "mine", label: "My Gists", icon: User },
        { id: "public", label: "Explore", icon: Search }
      ]}
    />

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
  </div>
</div>

<style>
  .gists-home { padding: 2rem; min-height: 100vh; background: var(--bg-color); }
  .container { max-width: 900px; margin: 0 auto; }

  .gist-list { display: flex; flex-direction: column; gap: 1rem; margin-top: 2rem; }

  .btn-primary {
    display: flex; align-items: center; gap: 0.5rem;
    background: var(--primary-color); color: white;
    padding: 0.75rem 1.25rem; border-radius: 1rem;
    text-decoration: none; font-weight: 700;
    transition: all 0.2s;
  }
  .btn-primary:hover { transform: translateY(-2px); box-shadow: var(--shadow-md); opacity: 0.9; }

  .skeleton-list { display: flex; flex-direction: column; gap: 1rem; }
  .skeleton-row { height: 80px; background: var(--md-sys-color-surface-container-low); border-radius: 1.25rem; animation: pulse 1.5s infinite; }
  @keyframes pulse { 0%, 100% { opacity: 0.5; } 50% { opacity: 0.8; } }

  .empty-state { text-align: center; padding: 4rem 2rem; color: var(--text-muted); display: flex; flex-direction: column; align-items: center; gap: 1rem; }
  .empty-state h3 { margin: 0; color: var(--text-color); }
</style>
