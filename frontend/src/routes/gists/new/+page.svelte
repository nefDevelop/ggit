<script lang="ts">
  import { apiFetch } from "$lib/api";
  import { goto } from "$app/navigation";
  import GistEditor from "$lib/components/GistEditor.svelte";

  let loading = $state(false);
  let error = $state("");

  let gist = $state({
    description: "",
    is_public: true,
    files: [{ filename: "file1.txt", content: "", language: "text" }]
  });

  async function handleSave() {
    loading = true;
    error = "";
    try {
      const res = await apiFetch("/api/v2/gists", {
        method: "POST",
        body: JSON.stringify(gist)
      });
      
      if (res.ok) {
        const result = await res.json();
        goto(`/gists/${result.uuid}`);
      } else {
        const errData = await res.json();
        error = errData.message || "Failed to create gist";
      }
    } catch (e) {
      error = "Connection error";
    } finally {
      loading = false;
    }
  }
</script>

<div class="new-gist-page">
  <div class="container">
    <header class="page-header">
      <div class="header-content">
        <h1>Create New Gist</h1>
        <p>Instantly share code snippets and notes.</p>
      </div>
    </header>

    {#if error}
      <div class="error-banner">
        <span class="material-symbols-rounded">error</span>
        {error}
      </div>
    {/if}

    <GistEditor bind:gist onsave={handleSave} saving={loading} />
  </div>
</div>

<style>
  .new-gist-page { min-height: 100vh; background: var(--bg-color); padding: 4rem 2rem; }
  .container { max-width: 900px; margin: 0 auto; }
  
  .page-header { margin-bottom: 3rem; }
  .page-header h1 { font-size: 2.5rem; font-weight: 800; color: var(--md-sys-color-on-surface); margin: 0 0 0.5rem 0; letter-spacing: -0.02em; }
  .page-header p { color: var(--md-sys-color-on-surface-variant); font-size: 1.1rem; }

  .error-banner { 
    padding: 1rem 1.5rem; background: var(--md-sys-color-error-container); 
    color: var(--md-sys-color-on-error-container); border-radius: 1rem;
    display: flex; align-items: center; gap: 0.75rem; margin-bottom: 2rem;
    font-weight: 600;
  }
</style>
