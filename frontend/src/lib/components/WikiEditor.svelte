<script lang="ts">
  let { editTitle = $bindable(), editContent = $bindable(), commitMessage = $bindable(), parentPath = $bindable(), pageOrder = $bindable(), wikiPages = [] }: {
    editTitle: string;
    editContent: string;
    commitMessage: string;
    parentPath: string;
    pageOrder: number | null;
    wikiPages: any[];
  } = $props();

  function extractFoldersFromSlantedPaths(items: any[]): string[] {
    const foldersSet = new Set<string>();
    for (const item of items) {
      if (item.is_folder) {
        foldersSet.add(item.name);
      } else if (item.name && item.name.includes('/')) {
        const parts = item.name.split('/');
        parts.pop();
        let cumulative = '';
        for (const part of parts) {
          cumulative = cumulative ? `${cumulative}/${part}` : part;
          foldersSet.add(cumulative);
        }
      }
    }
    return Array.from(foldersSet).sort();
  }

  const existingFolders = $derived(extractFoldersFromSlantedPaths(wikiPages));
</script>

<div class="edit-area">
  <div class="wiki-metadata-row">
    <div class="meta-field">
      <label for="parent-folder">Ubicación / Carpeta Padre</label>
      <select id="parent-folder" bind:value={parentPath}>
        <option value="">(Raíz / Root)</option>
        {#each existingFolders as folder}
          <option value={folder}>{folder}</option>
        {/each}
      </select>
    </div>
    <div class="meta-field">
      <label for="page-order">Orden de la página</label>
      <input type="number" id="page-order" bind:value={pageOrder} placeholder="Ej: 1, 2..." />
    </div>
  </div>

  <textarea class="wiki-editor" bind:value={editContent} placeholder="Escribe el contenido en Markdown..."></textarea>
  <div class="commit-msg-box">
    <label for="commit">Mensaje del commit (opcional)</label>
    <input type="text" id="commit" bind:value={commitMessage} placeholder="¿Qué has cambiado?" />
  </div>
</div>

<style>
  .edit-area { display: flex; flex-direction: column; gap: 1rem; }
  .wiki-metadata-row { display: flex; gap: 1rem; flex-wrap: wrap; }
  .meta-field { flex: 1; min-width: 200px; display: flex; flex-direction: column; gap: 0.25rem; }
  .meta-field label { font-size: 0.75rem; font-weight: 700; color: var(--text-muted); }
  .meta-field select, .meta-field input {
    background: rgba(0,0,0,0.2); border: 1px solid var(--md-sys-color-outline-variant);
    color: inherit; padding: 0.5rem; border-radius: 0.5rem; font-family: inherit;
  }
  .wiki-editor {
    width: 100%; min-height: 400px; resize: vertical;
    background: rgba(0,0,0,0.2); border: 1px solid var(--md-sys-color-outline-variant);
    color: inherit; padding: 1rem; border-radius: 0.75rem;
    font-family: "JetBrains Mono", monospace; font-size: 0.9rem;
    box-sizing: border-box; outline: none;
  }
  .commit-msg-box { display: flex; flex-direction: column; gap: 0.25rem; }
  .commit-msg-box label { font-size: 0.75rem; font-weight: 700; color: var(--text-muted); }
  .commit-msg-box input {
    background: rgba(0,0,0,0.2); border: 1px solid var(--md-sys-color-outline-variant);
    color: inherit; padding: 0.5rem 1rem; border-radius: 0.5rem; font-family: inherit;
  }
</style>
