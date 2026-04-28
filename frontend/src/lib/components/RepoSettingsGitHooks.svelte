<script lang="ts">
  import { apiFetch } from "$lib/api";
  import { onMount, getContext } from "svelte";
  import Modal from "$lib/components/Modal.svelte";
  import Button from "$lib/components/ui/Button.svelte";
  import { toast } from "$lib/toast.svelte";

  const repo: any = getContext("repo");

  let hooks = $state<any[]>([]);
  let loading = $state(true);
  let error = $state("");

  let editingHook = $state<any>(null);
  let hookContent = $state("");
  let isSubmitting = $state(false);

  async function fetchHooks() {
    loading = true;
    try {
      const res = await apiFetch(
        `/api/v2/${repo.data.owner.username}/${repo.data.name}/settings/hooks/git`,
      );
      if (res.ok) {
        hooks = await res.json();
      } else {
        error = "Failed to load git hooks.";
      }
    } catch (e) {
      error = "An error occurred while fetching git hooks.";
    } finally {
      loading = false;
    }
  }

  function handleEdit(hook: any) {
    editingHook = hook;
    hookContent = hook.content;
  }

  async function handleSave() {
    isSubmitting = true;
    try {
      const res = await apiFetch(
        `/api/v2/${repo.data.owner.username}/${repo.data.name}/settings/hooks/git/${editingHook.name}`,
        {
          method: "PATCH",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({ content: hookContent }),
        },
      );

      if (res.ok) {
        toast.success(`Hook ${editingHook.name} updated successfully.`);
        editingHook = null;
        fetchHooks();
      } else {
        const data = await res.json();
        toast.error(data.message || "Failed to update hook.");
      }
    } catch (e) {
      console.error(e);
      toast.error("An error occurred while saving the hook.");
    } finally {
      isSubmitting = false;
    }
  }

  onMount(fetchHooks);
</script>

<div class="git-hooks-manager">
  <header class="section-header">
    <div class="header-content">
      <h2>Git Hooks</h2>
      <p>
        Server-side Git hooks allow you to execute custom scripts on the server when certain actions are performed (e.g., pushing commits).
      </p>
    </div>
  </header>

  {#if loading}
    <div class="skeleton-list">
      {#each Array(3) as _}
        <div class="skeleton-item tonal"></div>
      {/each}
    </div>
  {:else if error}
    <div class="error-banner">
      <span class="material-symbols-rounded">error</span>
      {error}
    </div>
  {:else}
    <div class="hooks-list">
      {#each hooks as hook}
        <div class="hook-card tonal">
          <div class="hook-info">
            <div class="hook-name-row">
              <span class="material-symbols-rounded">terminal</span>
              <span class="hook-name">{hook.name}</span>
              {#if hook.is_active}
                <span class="badge badge-active">Active</span>
              {:else}
                <span class="badge badge-inactive">Not Configured</span>
              {/if}
            </div>
            <p class="hook-description">
              {#if hook.name === 'pre-receive'}
                Executed before any refs are updated. Use it to enforce project policies.
              {:else if hook.name === 'update'}
                Executed once for each ref being updated. Similar to pre-receive but per-ref.
              {:else if hook.name === 'post-receive'}
                Executed after all refs have been updated. Use it for notifications or deployments.
              {/if}
            </p>
          </div>
          <div class="hook-actions">
            <Button variant="secondary" size="sm" icon="edit" onclick={() => handleEdit(hook)}>
              Edit Script
            </Button>
          </div>
        </div>
      {/each}
    </div>
  {/if}

  <Modal 
    show={!!editingHook} 
    title={`Edit Git Hook: ${editingHook?.name}`} 
    onclose={() => editingHook = null}
    maxWidth="900px"
  >
    {#if editingHook}
      <div class="hook-editor">
        <div class="editor-help">
          <span class="material-symbols-rounded">info</span>
          <p>This script will be executed as a shell script on the server. Make sure it has a proper shebang (e.g., <code>#!/bin/bash</code>).</p>
        </div>
        
        <textarea 
          bind:value={hookContent} 
          class="hook-textarea"
          placeholder="#!/bin/bash\n\n# Your custom hook logic here..."
          spellcheck="false"
        ></textarea>

        <footer class="modal-footer">
          <Button variant="ghost" onclick={() => editingHook = null}>Cancel</Button>
          <Button variant="primary" loading={isSubmitting} onclick={handleSave}>
            Save Hook
          </Button>
        </footer>
      </div>
    {/if}
  </Modal>
</div>

<style>
  .git-hooks-manager { display: flex; flex-direction: column; gap: 2rem; animation: fadeIn 0.3s ease-out; }
  
  .section-header { display: flex; justify-content: space-between; align-items: flex-start; }
  .header-content h2 { margin: 0 0 0.5rem 0; font-size: 1.5rem; font-weight: 800; }
  .header-content p { margin: 0; color: var(--text-muted); font-size: 0.95rem; max-width: 800px; }

  .hooks-list { display: flex; flex-direction: column; gap: 1rem; }
  .hook-card { display: flex; justify-content: space-between; align-items: center; padding: 1.5rem 2rem; }
  
  .hook-info { display: flex; flex-direction: column; gap: 0.5rem; }
  .hook-name-row { display: flex; align-items: center; gap: 0.75rem; }
  .hook-name-row span:first-child { color: var(--primary-color); opacity: 0.7; }
  .hook-name { font-weight: 800; font-size: 1.1rem; font-family: 'JetBrains Mono', monospace; }
  
  .badge { font-size: 0.65rem; font-weight: 800; text-transform: uppercase; padding: 0.15rem 0.6rem; border-radius: 4px; letter-spacing: 0.05em; }
  .badge-active { background: var(--success); color: var(--on-success); }
  .badge-inactive { background: rgba(255,255,255,0.05); color: var(--text-muted); }

  .hook-description { margin: 0; font-size: 0.9rem; color: var(--text-muted); opacity: 0.8; }

  .hook-editor { display: flex; flex-direction: column; gap: 1.5rem; }
  .editor-help { display: flex; align-items: center; gap: 1rem; padding: 1rem; background: rgba(var(--primary-rgb), 0.05); border-radius: 0.75rem; border: 1px solid rgba(var(--primary-rgb), 0.1); color: var(--text-muted); font-size: 0.85rem; }
  .editor-help code { background: rgba(0,0,0,0.2); padding: 0.1rem 0.3rem; border-radius: 4px; }

  .hook-textarea { width: 100%; height: 400px; background: #000; color: #fff; border: 1px solid var(--md-sys-color-outline-variant); border-radius: 1rem; padding: 1.5rem; font-family: 'JetBrains Mono', 'Fira Code', monospace; font-size: 0.9rem; line-height: 1.5; resize: vertical; outline: none; }
  .hook-textarea:focus { border-color: var(--primary-color); box-shadow: 0 0 0 2px rgba(var(--primary-rgb), 0.2); }

  .modal-footer { display: flex; justify-content: flex-end; gap: 1rem; padding-top: 1rem; border-top: 1px solid var(--md-sys-color-outline-variant); }

  .skeleton-list { display: flex; flex-direction: column; gap: 1rem; }
  .skeleton-item { height: 100px; width: 100%; border-radius: 1rem; animation: pulse 1.5s infinite; }
  @keyframes pulse { 0%, 100% { opacity: 0.6; } 50% { opacity: 0.3; } }

  @keyframes fadeIn { from { opacity: 0; transform: translateY(10px); } to { opacity: 1; transform: translateY(0); } }
</style>
