<script lang="ts">
  import { apiFetch } from "$lib/api";
  import { onMount } from 'svelte';
  import CommitPanel from './CommitPanel.svelte';
  
  let { username, reponame, treePath, branch, isNewFile = false, onComplete } = $props();
  
  let content = $state('');
  let loading = $state(true);
  let showCommit = $state(false);
  let isSaving = $state(false);

  async function fetchFile() {
    if (isNewFile) {
       loading = false;
       return;
    }
    loading = true;
    try {
      // Reusing repository raw view logic or similar
      const res = await apiFetch(`/${username}/${reponame}/raw/${branch}/${treePath}`);
      if (res.ok) content = await res.text();
    } catch (e) {
      console.error(e);
    } finally {
      loading = false;
    }
  }

  async function handleCommit(commitData: any) {
    isSaving = true;
    try {
      const res = await apiFetch(`/api/v2/${username}/${reponame}/editor/edit`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          tree_path: treePath,
          content: content,
          commit_summary: commitData.summary,
          commit_message: commitData.message,
          is_new_file: isNewFile,
          branch: commitData.commitChoice === 'branch' ? commitData.newBranchName : branch
        })
      });

      if (res.ok) {
        onComplete();
      }
    } catch (e) {
      console.error(e);
    } finally {
      isSaving = false;
    }
  }

  onMount(fetchFile);
</script>

<div class="web-editor-container">
  <header class="editor-toolbar island tonal">
    <div class="toolbar-left">
      <span class="material-symbols-rounded">description</span>
      <div class="file-path">
        <span class="repo-name">{reponame}</span> / 
        <span class="path">{treePath}</span>
      </div>
    </div>
    <div class="toolbar-right">
       <button class="btn btn-primary" onclick={() => showCommit = true}>
         Commit changes...
       </button>
    </div>
  </header>

  <main class="editor-main island">
    {#if loading}
      <div class="loader">Loading file...</div>
    {:else}
      <textarea 
        bind:value={content} 
        spellcheck="false"
        placeholder="Start typing..."
      ></textarea>
    {/if}
  </main>

  {#if showCommit}
    <div class="overlay">
       <CommitPanel 
          {isSaving}
          onCommit={handleCommit} 
          onCancel={() => showCommit = false} 
       />
    </div>
  {/if}
</div>

<style>
  .web-editor-container { display: flex; flex-direction: column; height: 80vh; gap: 1rem; position: relative; }
  
  .editor-toolbar { display: flex; justify-content: space-between; align-items: center; padding: 1rem 2rem; }
  .toolbar-left { display: flex; align-items: center; gap: 1rem; }
  .file-path { font-family: monospace; font-size: 1rem; }
  .repo-name { color: var(--text-muted); }
  .path { font-weight: 800; color: var(--primary-color); }

  .editor-main { flex: 1; padding: 0; overflow: hidden; display: flex; border: 1px solid var(--md-sys-color-outline-variant); }
  textarea { flex: 1; background: #0d1117; color: #e6edf3; border: none; padding: 2rem; font-family: 'JetBrains Mono', 'Fira Code', monospace; font-size: 1rem; line-height: 1.6; resize: none; overflow-y: auto; outline: none; }
  
  .overlay { position: absolute; inset: 0; background: rgba(0,0,0,0.8); backdrop-filter: blur(8px); display: flex; align-items: center; justify-content: center; z-index: 100; animation: fadeIn 0.2s ease-out; }
  @keyframes fadeIn { from { opacity: 0; } to { opacity: 1; } }

  .loader { padding: 4rem; text-align: center; width: 100%; color: var(--text-muted); }
</style>
