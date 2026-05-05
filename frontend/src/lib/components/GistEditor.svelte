<script lang="ts">
  import Card from "./ui/Card.svelte";
  import Button from "./ui/Button.svelte";
  import Input from "./ui/Input.svelte";
  import Badge from "./ui/Badge.svelte";
  import CodeEditor from "./ui/CodeEditor.svelte";

  let { gist = $bindable(null), onsave, ondelete, saving = false, isEdit = false } = $props<{ 
    gist?: any,
    onsave: () => void,
    ondelete?: () => void,
    saving?: boolean,
    isEdit?: boolean
  }>();

  let description = $state("");
  let isPublic = $state(true);
  let files = $state([{ filename: "file1.txt", content: "", language: "text" }]);

  $effect(() => {
    if (gist) {
      description = gist.description || "";
      isPublic = gist.is_public ?? true;
      if (gist.files && gist.files.length > 0) {
        files = JSON.parse(JSON.stringify(gist.files)); // Deep copy to avoid direct mutation
      }
    }
  });

  function addFile() {
    files.push({ filename: `file${files.length + 1}.txt`, content: "", language: "text" });
  }

  function removeFile(index: number) {
    if (files.length > 1) {
      files.splice(index, 1);
    }
  }

  function detectLanguage(filename: string) {
    const ext = filename.split('.').pop()?.toLowerCase();
    const map: Record<string, string> = {
      'js': 'javascript', 'ts': 'typescript', 'py': 'python', 'go': 'go',
      'md': 'markdown', 'json': 'json', 'html': 'html', 'css': 'css',
      'rs': 'rust', 'c': 'c', 'cpp': 'cpp', 'java': 'java', 'rb': 'ruby',
      'sh': 'bash', 'yaml': 'yaml', 'yml': 'yaml', 'xml': 'xml', 'sql': 'sql'
    };
    return map[ext || ''] || 'text';
  }

  $effect(() => {
    files.forEach(f => {
      f.language = detectLanguage(f.filename);
    });
  });

  function handleSubmit() {
    // Update the gist prop or a local object to send back
    if (gist) {
      gist.description = description;
      gist.is_public = isPublic;
      gist.files = files;
    }
    onsave();
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
          <CodeEditor 
            bind:value={file.content} 
            placeholder="Enter file content..." 
            language={file.language}
          />
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
    {#if isEdit && ondelete}
      <Button variant="ghost" onclick={ondelete} class="delete-gist-btn">
        <span class="material-symbols-rounded">delete</span>
        Eliminar Gist
      </Button>
    {/if}
    <Button variant="primary" onclick={handleSubmit} loading={saving}>
      {#if isEdit}
        {saving ? 'Updating...' : 'Update Gist'}
      {:else}
        {saving ? 'Creating...' : 'Create Gist'}
      {/if}
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

  .editor-footer { display: flex; align-items: center; gap: 1.5rem; padding-top: 1rem; }
  
  .footer-meta { display: flex; align-items: center; gap: 1.5rem; }
  .visibility-toggle { display: flex; align-items: center; gap: 0.5rem; color: var(--md-sys-color-on-surface-variant); font-weight: 600; }
  .visibility-toggle label { display: flex; align-items: center; gap: 0.5rem; cursor: pointer; }

  :global(.delete-gist-btn) { color: var(--md-sys-color-error) !important; }
  :global(.delete-gist-btn:hover) { background: var(--md-sys-color-error-container) !important; }
</style>
