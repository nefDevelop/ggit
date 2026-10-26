<script lang="ts">
  import { apiFetch } from "$lib/api";
  import { onMount } from "svelte";
  import GistRow from "$lib/components/GistRow.svelte";
  import Card from "$lib/components/ui/Card.svelte";
  import Button from "$lib/components/ui/Button.svelte";

  let gists = $state<any[]>([]);
  let loading = $state(true);
  let error = $state("");

  async function fetchGists() {
    loading = true;
    try {
      // In a real app, we'd have a global gists endpoint. 
      // For now, let's assume /api/v2/gists returns public ones.
      // Wait, I didn't implement GET /api/v2/gists yet.
      const res = await apiFetch("/api/v2/gists");
      if (res.ok) {
        gists = await res.json();
      } else {
        error = "Failed to load gists";
      }
    } catch (e) {
      error = "Connection error";
    } finally {
      loading = false;
    }
  }

  onMount(fetchGists);
</script>

<div class="explore-gists-page">
  <div class="container">
    <header class="page-header">
      <div class="header-content">
        <h1>Explore Gists</h1>
        <p>Discover useful snippets and notes from the community.</p>
      </div>
      <Button variant="primary" href="/gists/new">
        <span class="material-symbols-rounded">add</span>
        New Gist
      </Button>
    </header>

    {#if loading}
      <div class="skeleton-stack">
        {#each Array(5) as _}
          <div class="skeleton-row"></div>
        {/each}
      </div>
    {:else if gists.length === 0}
      <Card variant="low" padding="lg" class="empty-state">
        <span class="material-symbols-rounded">code_off</span>
        <h3>No public gists found</h3>
        <p>Be the first to share something!</p>
      </Card>
    {:else}
      <div class="gists-stack">
        {#each gists as gist}
          <Card variant="low" padding="none" hover class="row-card">
            <GistRow {gist} />
          </Card>
        {/each}
      </div>
    {/if}
  </div>
</div>

<style>
  .explore-gists-page { min-height: 100vh; background: var(--bg-color); padding: var(--space-7) var(--space-6); }
  .container { max-width: 1000px; margin: 0 auto; }
  
  .page-header { display: flex; align-items: flex-end; justify-content: space-between; margin-bottom: 3rem; }
  .page-header h1 { font-size: 2.5rem; font-weight: 800; color: var(--md-sys-color-on-surface); margin: 0 0 0.5rem 0; letter-spacing: -0.02em; }
  .page-header p { color: var(--md-sys-color-on-surface-variant); font-size: 1.1rem; }

  .gists-stack { display: flex; flex-direction: column; gap: 1rem; }
  
  .skeleton-stack { display: flex; flex-direction: column; gap: 1rem; }
  .skeleton-row { height: 80px; background: var(--md-sys-color-surface-container-low); border-radius: 1.25rem; animation: pulse 1.5s infinite; }
  
  @keyframes pulse { 0%, 100% { opacity: 0.5; } 50% { opacity: 0.8; } }

  :global(.empty-state) { text-align: center; display: flex; flex-direction: column; align-items: center; gap: 1rem; }
  :global(.empty-state span) { font-size: 4rem; color: var(--md-sys-color-on-surface-variant); opacity: 0.3; }

  @media (max-width: 600px) {
    .explore-gists-page { padding: var(--space-6) var(--space-4); }
  }
</style>
