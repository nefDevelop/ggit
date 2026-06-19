<script lang="ts">
  import MarkdownRenderer from "./ui/MarkdownRenderer.svelte";
  import Button from "./ui/Button.svelte";
  import BookText from "lucide-svelte/icons/book-text";
  import Plus from "lucide-svelte/icons/plus";

  let { page = null, repoOwner = "", repoName = "", onCreatePage }: {
    page: any;
    repoOwner: string;
    repoName: string;
    onCreatePage?: () => void;
  } = $props();
</script>

{#if page}
  <MarkdownRenderer content={page.content} repoContext={{owner: repoOwner, name: repoName}} />
{:else}
  <div class="empty-wiki-state">
    <div class="empty-icon-ring"><BookText size={40} /></div>
    <h2>Wiki Vacía</h2>
    <p>Este repositorio aún no tiene documentación. Las wikis son geniales para guías de usuario y documentación técnica.</p>
    <div style="margin-top: 1rem;">
      <Button variant="primary" onclick={onCreatePage} icon={Plus}>Crear Primera Página</Button>
    </div>
  </div>
{/if}

<style>
  .empty-wiki-state {
    display: flex; flex-direction: column; align-items: center;
    justify-content: center; padding: 4rem 2rem; gap: 1rem;
    text-align: center;
  }
  .empty-icon-ring {
    width: 80px; height: 80px; border-radius: 50%;
    background: var(--md-sys-color-primary-container);
    color: var(--md-sys-color-on-primary-container);
    display: flex; align-items: center; justify-content: center;
  }
  .empty-wiki-state h2 { margin: 0; font-size: 1.5rem; font-weight: 800; }
  .empty-wiki-state p { margin: 0; color: var(--text-muted); max-width: 400px; }
</style>
