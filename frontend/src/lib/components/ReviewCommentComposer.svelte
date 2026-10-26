<script lang="ts">
  import Button from "./ui/Button.svelte";
  import MarkdownRenderer from "./ui/MarkdownRenderer.svelte";
  
  let { onSave, onCancel, initialValue = "" } = $props<{
    onSave: (content: string) => void;
    onCancel: () => void;
    initialValue?: string;
  }>();

  let content = $state("");
  $effect(() => {
    content = initialValue;
  });
  let isPreview = $state(false);
</script>

<div class="inline-comment-composer tonal">
  <nav class="composer-tabs">
    <button class:active={!isPreview} onclick={() => isPreview = false}>Escribir</button>
    <button class:active={isPreview} onclick={() => isPreview = true}>Vista Previa</button>
  </nav>

  <div class="composer-body">
    {#if isPreview}
      <div class="preview-area">
        {#if content.trim()}
          <MarkdownRenderer {content} />
        {:else}
          <p class="empty-preview">Nada que previsualizar.</p>
        {/if}
      </div>
    {:else}
      <textarea 
        bind:value={content} 
        placeholder="Deja un comentario..."
      ></textarea>
    {/if}
  </div>

  <footer class="composer-actions">
    <Button variant="ghost" size="sm" onclick={onCancel}>Cancelar</Button>
    <Button 
      variant="primary" 
      size="sm" 
      disabled={!content.trim()} 
      onclick={() => onSave(content)}
    >
      Añadir comentario
    </Button>
  </footer>
</div>

<style>
  .inline-comment-composer {
    margin: 0.5rem 1rem;
    border-radius: 1rem;
    border: 1px solid var(--md-sys-color-outline-variant);
    background: var(--md-sys-color-surface-container-low);
    display: flex;
    flex-direction: column;
    overflow: hidden;
  }

  .composer-tabs {
    display: flex;
    gap: 1rem;
    padding: 0.5rem 1rem;
    background: var(--sunken-1);
    border-bottom: 1px solid var(--md-sys-color-outline-variant);
  }

  .composer-tabs button {
    background: none;
    border: none;
    color: var(--text-muted);
    font-size: 0.8rem;
    font-weight: 800;
    cursor: pointer;
    padding: 0.25rem 0.5rem;
    border-radius: 4px;
  }

  .composer-tabs button.active {
    color: var(--primary-color);
    background: rgba(var(--primary-rgb), 0.1);
  }

  .composer-body {
    padding: 1rem;
  }

  textarea {
    width: 100%;
    min-height: 100px;
    background: transparent;
    border: none;
    color: var(--text-color);
    font-family: inherit;
    font-size: 0.95rem;
    resize: vertical;
    outline: none;
  }

  .preview-area {
    min-height: 100px;
  }

  .empty-preview {
    color: var(--text-muted);
    font-style: italic;
    font-size: 0.9rem;
  }

  .composer-actions {
    display: flex;
    justify-content: flex-end;
    gap: 0.75rem;
    padding: 0.75rem 1rem;
    background: var(--sunken-1);
    border-top: 1px solid var(--md-sys-color-outline-variant);
  }
</style>
