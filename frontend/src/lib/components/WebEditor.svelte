<script lang="ts">
  import { apiFetch } from "$lib/api";
  import { onMount } from "svelte";
  import CommitPanel from "./CommitPanel.svelte";
  import CodeEditor from "./ui/CodeEditor.svelte";

  let {
    username,
    reponame,
    treePath,
    branch,
    isNewFile = false,
    onComplete,
  } = $props();

  let content = $state("");
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
      const res = await apiFetch(
        `/api/v2/${username}/${reponame}/raw/${branch}/${treePath}`,
      );
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
      const res = await apiFetch(
        `/api/v2/${username}/${reponame}/editor/edit`,
        {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({
            tree_path: treePath,
            content: content,
            commit_summary: commitData.summary,
            commit_message: commitData.message,
            is_new_file: isNewFile,
            branch:
              commitData.commitChoice === "branch"
                ? commitData.newBranchName
                : branch,
          }),
        },
      );

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
  <div class="editor-glass glass">
    <header class="editor-toolbar">
      <div class="toolbar-left">
        <div class="bread-container">
          <span class="material-symbols-rounded mini-icon">edit_note</span>
          <div class="file-path">
            <span class="repo-name">{reponame}</span>
            <span class="sep">/</span>
            <span class="path">{treePath}</span>
          </div>
        </div>
      </div>
      <div class="toolbar-right">
        <div class="btn-group">
          <button class="btn-action-primary" onclick={() => (showCommit = true)} disabled={loading}>
            <span class="material-symbols-rounded">cloud_upload</span>
            <span>Commit Changes</span>
          </button>
        </div>
      </div>
    </header>

    <main class="editor-main">
      {#if loading}
        <div class="loader-container">
          <div class="spinner"></div>
          <p>Preparing editor...</p>
        </div>
      {:else}
        <div class="editor-wrapper">
          <CodeEditor 
            bind:value={content} 
            placeholder="Start typing..." 
          />
        </div>
      {/if}
    </main>
  </div>

  {#if showCommit}
    <div class="overlay">
      <CommitPanel
        {isSaving}
        onCommit={handleCommit}
        onCancel={() => (showCommit = false)}
      />
    </div>
  {/if}
</div>

<style>
  .web-editor-container {
    display: flex;
    flex-direction: column;
    height: calc(100vh - 12rem);
    position: relative;
    animation: fadeIn 0.4s ease-out;
  }

  .editor-glass {
    display: flex;
    flex-direction: column;
    height: 100%;
    border-radius: 1.5rem;
    overflow: hidden;
    background: var(--surface);
    border: 1px solid var(--md-sys-color-outline-variant);
  }

  .editor-toolbar {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 1rem 1.5rem;
    background: rgba(var(--primary-rgb), 0.05);
    border-bottom: 1px solid var(--md-sys-color-outline-variant);
  }

  .bread-container {
    display: flex;
    align-items: center;
    gap: 0.75rem;
    padding: 0.5rem 1rem;
    background: var(--md3-surface-2);
    border-radius: 1rem;
    border: 1px solid var(--md-sys-color-outline-variant);
  }

  .mini-icon { font-size: 1.25rem; color: var(--primary-color); }

  .file-path {
    font-family: inherit;
    font-size: 0.95rem;
    font-weight: 500;
    display: flex;
    align-items: center;
    gap: 0.5rem;
  }
  .repo-name { color: var(--text-muted); }
  .sep { opacity: 0.3; }
  .path { font-weight: 800; color: var(--text-color); }

  .editor-main {
    flex: 1;
    display: flex;
    overflow: hidden;
  }

  .editor-wrapper {
    flex: 1;
    display: flex;
  }


  .loader-container {
    flex: 1;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    gap: 1rem;
    color: var(--text-muted);
  }

  .spinner {
    width: 40px;
    height: 40px;
    border: 3px solid rgba(var(--primary-rgb), 0.1);
    border-top-color: var(--primary-color);
    border-radius: 50%;
    animation: spin 1s linear infinite;
  }

  @keyframes spin { to { transform: rotate(360deg); } }
  @keyframes fadeIn {
    from { opacity: 0; transform: translateY(10px); }
    to { opacity: 1; transform: translateY(0); }
  }

  .overlay {
    position: fixed;
    inset: 0;
    background: rgba(0, 0, 0, 0.4);
    backdrop-filter: blur(12px);
    display: flex;
    align-items: center;
    justify-content: center;
    z-index: 1000;
    animation: fadeInHalf 0.3s ease-out;
  }

  @keyframes fadeInHalf { from { opacity: 0; } to { opacity: 1; } }
</style>
