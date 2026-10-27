<script lang="ts">
  import { apiFetch } from "$lib/api";
  import { onMount } from "svelte";
  import { page } from "$app/state";
  import GistView from "$lib/components/GistView.svelte";
  import Card from "$lib/components/ui/Card.svelte";
  import Button from "$lib/components/ui/Button.svelte";

  const uuid = $derived(page.params.uuid);
  let gist = $state<any>(null);
  let loading = $state(true);
  let error = $state("");

  async function fetchGist() {
    loading = true;
    try {
      const res = await apiFetch(`/api/v2/gists/${uuid}`);
      if (res.ok) {
        gist = await res.json();
      } else if (res.status === 404) {
        error = "Gist not found";
      } else {
        error = "Failed to load Gist";
      }
    } catch (e) {
      error = "Connection error";
    } finally {
      loading = false;
    }
  }

  onMount(fetchGist);

  $effect(() => {
    if (uuid) fetchGist();
  });
  import CircleAlert from "lucide-svelte/icons/circle-alert";
</script>

<div class="gist-page">
  <div class="container">
    {#if loading}
      <div class="skeleton-view">
        <div class="skeleton-header"></div>
        <div class="skeleton-body"></div>
      </div>
    {:else if error}
      <Card variant="low" padding="lg" class="error-state">
        <CircleAlert size={18} />
        <h3>{error}</h3>
        <p>The gist you are looking for might be private or deleted.</p>
        <Button variant="primary" href="/">Go Home</Button>
      </Card>
    {:else if gist}
      <GistView {gist} />
    {/if}
  </div>
</div>

<style>
  .gist-page { min-height: 100vh; background: var(--bg-color); padding: var(--space-7) var(--space-6); }
  .container { max-width: 900px; margin: 0 auto; }
  
  .skeleton-view { display: flex; flex-direction: column; gap: 2rem; }
  .skeleton-header { height: 100px; background: var(--md-sys-color-surface-container-low); border-radius: 1.5rem; animation: pulse 1.5s infinite; }
  .skeleton-body { height: 400px; background: var(--md-sys-color-surface-container-low); border-radius: 1.5rem; animation: pulse 1.5s infinite; }
  
  @keyframes pulse { 0%, 100% { opacity: 0.5; } 50% { opacity: 0.8; } }

  :global(.error-state) { text-align: center; display: flex; flex-direction: column; align-items: center; gap: 1rem; }
  :global(.error-state span) { font-size: 4rem; color: var(--md-sys-color-error); }

  @media (max-width: 600px) {
    .gist-page { padding: var(--space-6) var(--space-4); }
  }
</style>
