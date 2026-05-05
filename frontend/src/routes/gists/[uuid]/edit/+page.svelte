<script lang="ts">
  import { apiFetch } from "$lib/api";
  import { onMount } from "svelte";
  import { page } from "$app/state";
  import { goto } from "$app/navigation";
  import GistEditor from "$lib/components/GistEditor.svelte";
  import Card from "$lib/components/ui/Card.svelte";
  import Button from "$lib/components/ui/Button.svelte";
  import { toast } from "$lib/toast.svelte";

  const uuid = $derived(page.params.uuid);
  let gist = $state<any>(null);
  let loading = $state(true);
  let saving = $state(false);

  async function fetchGist() {
    loading = true;
    try {
      const res = await apiFetch(`/api/v2/gists/${uuid}`);
      if (res.ok) {
        gist = await res.json();
      } else {
        toast.error("Failed to load gist");
        goto("/gists");
      }
    } catch (e) {
      toast.error("Connection error");
    } finally {
      loading = false;
    }
  }

  async function handleSave() {
    saving = true;
    try {
      const res = await apiFetch(`/api/v2/gists/${uuid}`, {
        method: "PUT",
        body: JSON.stringify({
          description: gist.description,
          is_public: gist.is_public,
          files: gist.files
        })
      });

      if (res.ok) {
        toast.success("Gist updated successfully!");
        goto(`/gists/${uuid}`);
      } else {
        const err = await res.text();
        toast.error(`Failed to update: ${err}`);
      }
    } catch (e) {
      toast.error("Network error");
    } finally {
      saving = false;
    }
  }

  async function handleDelete() {
    if (!confirm("Are you sure you want to delete this gist?")) return;
    try {
      const res = await apiFetch(`/api/v2/gists/${uuid}`, { method: "DELETE" });
      if (res.ok) {
        toast.success("Gist deleted");
        goto("/gists");
      } else {
        toast.error("Failed to delete gist");
      }
    } catch (e) {
      toast.error("Network error");
    }
  }

  onMount(fetchGist);
</script>

<div class="gist-edit-page">
  <div class="container">
    {#if loading}
      <div class="skeleton-view">
        <div class="skeleton-header"></div>
        <div class="skeleton-body"></div>
      </div>
    {:else if gist}
      <GistEditor bind:gist onsave={handleSave} ondelete={handleDelete} {saving} isEdit={true} />
    {/if}
  </div>
</div>

<style>
  .gist-edit-page { min-height: 100vh; background: var(--bg-color); padding: 4rem 2rem; }
  .container { max-width: 900px; margin: 0 auto; }
  
  .skeleton-view { display: flex; flex-direction: column; gap: 2rem; }
  .skeleton-header { height: 100px; background: var(--md-sys-color-surface-container-low); border-radius: 1.5rem; animation: pulse 1.5s infinite; }
  .skeleton-body { height: 400px; background: var(--md-sys-color-surface-container-low); border-radius: 1.5rem; animation: pulse 1.5s infinite; }
  @keyframes pulse { 0%, 100% { opacity: 0.5; } 50% { opacity: 0.8; } }
</style>
