<script lang="ts">
  import Card from "./ui/Card.svelte";
  import Button from "./ui/Button.svelte";
  import Input from "./ui/Input.svelte";
  import Badge from "./ui/Badge.svelte";

  let { onSave, loading = false } = $props<{ 
    onSave: (data: any) => void,
    loading?: boolean 
  }>();

  let description = $state("");
  let isPublic = $state(true);
  let files = $state([{ filename: "file1.txt", content: "", language: "text" }]);

  function addFile() {
    files.push({ filename: `file${files.length + 1}.txt`, content: "", language: "text" });
  }

  function removeFile(index: number) {
    if (files.length > 1) {
      files.splice(index, 1);
    }
  }

  function handleSubmit() {
    onSave({
      description,
      is_public: isPublic,
      files
    });
  }
</script>

<div class="gist-editor">
  <div class="editor-header">
    <Input 
      placeholder="Gist description..." 
      bind:value={description} 
      class="desc-input"
    />
  </div>

  <div class="files-stack">
    {#each files as file, i}
      <Card variant="low" padding="none" class="file-card">
        <div class="file-header">
          <Input 
            placeholder="Filename including extension..." 
            bind:value={file.filename} 
            class="filename-input"
          />
          <div class="spacer"></div>
          {#if files.length > 1}
            <button class="delete-btn" onclick={() => removeFile(i)}>
              <span class="material-symbols-rounded">delete</span>
            </button>
          {/if}
        </div>
        <div class="file-body">
          <textarea 
            placeholder="Enter file content..." 
            bind:value={file.content}
            spellcheck="false"
          ></textarea>
        </div>
      </Card>
    {/each}
  </div>

  <div class="editor-footer">
    <Button variant="secondary" onclick={addFile} icon="add">
      Add File
    </Button>
    <div class="spacer"></div>
    <div class="footer-meta">
      <div class="visibility-toggle">
        <label>
          <input type="checkbox" bind:checked={isPublic} />
          <span>Public Gist</span>
        </label>
      </div>
      <Badge variant="secondary">{files.length} {files.length === 1 ? 'file' : 'files'}</Badge>
    </div>
    <Button variant="primary" onclick={handleSubmit} {loading}>
      {loading ? 'Creating...' : 'Create Gist'}
    </Button>
  </div>
</div>

<style>
  .gist-editor { display: flex; flex-direction: column; gap: 1.5rem; }
  .editor-header { margin-bottom: 0.5rem; }
  
  :global(.desc-input) { font-size: 1.2rem !important; font-weight: 600; }
  
  .files-stack { display: flex; flex-direction: column; gap: 1.5rem; }
  
  .file-header { 
    display: flex; align-items: center; gap: 1rem; padding: 0.5rem 1rem;
    background: var(--md-sys-color-surface-container); border-bottom: 1px solid var(--md-sys-color-outline-variant);
    border-radius: 1.25rem 1.25rem 0 0;
  }
  
  :global(.filename-input) { background: transparent !important; border: none !important; font-family: 'JetBrains Mono', monospace; font-weight: 600; }
  
  .spacer { flex: 1; }
  
  .delete-btn { 
    background: transparent; border: none; color: var(--md-sys-color-error); 
    cursor: pointer; padding: 0.5rem; border-radius: 50%; transition: all 0.2s;
  }
  .delete-btn:hover { background: var(--md-sys-color-error-container); }

  .file-body { padding: 0; }
  textarea {
    width: 100%; min-height: 300px; padding: 1.5rem;
    background: var(--md-sys-color-surface-container-low); border: none;
    color: var(--md-sys-color-on-surface); font-family: 'JetBrains Mono', monospace;
    font-size: 0.95rem; line-height: 1.6; outline: none;
    resize: vertical; border-radius: 0 0 1.25rem 1.25rem;
  }

  .editor-footer { display: flex; align-items: center; gap: 1.5rem; padding-top: 1rem; }
  
  .footer-meta { display: flex; align-items: center; gap: 1.5rem; }
  .visibility-toggle { display: flex; align-items: center; gap: 0.5rem; color: var(--md-sys-color-on-surface-variant); font-weight: 600; }
  .visibility-toggle label { display: flex; align-items: center; gap: 0.5rem; cursor: pointer; }
</style>
