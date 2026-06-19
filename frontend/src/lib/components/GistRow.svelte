<script lang="ts">
  

import Code from "lucide-svelte/icons/code";
import ChevronRight from "lucide-svelte/icons/chevron-right";


  import Badge from "./ui/Badge.svelte";

  let { gist } = $props<{ gist: any }>();
</script>

<a href="/gists/{gist.uuid}" class="gist-row">
  <div class="row-icon">
    <Code size={24} />
  </div>
  <div class="row-info">
    <div class="title-line">
      <span class="gist-title">{gist.description || 'Untitled Gist'}</span>
      {#if !gist.is_public}
        <Badge variant="neutral">Secret</Badge>
      {/if}
    </div>
    <div class="meta-line">
      <span>{gist.uuid.split('-')[0]}</span>
      <span class="sep">•</span>
      <span>Updated {new Date(gist.updated_at).toLocaleDateString()}</span>
    </div>
  </div>
  <div class="row-arrow">
    <ChevronRight size={24} />
  </div>
</a>

<style>
  .gist-row {
    display: flex; align-items: center; gap: 1.25rem; padding: 1.25rem;
    text-decoration: none; transition: all 0.2s; border-radius: 1.25rem;
  }
  .gist-row:hover { background: var(--md-sys-color-surface-container); }
  
  .row-icon {
    width: 48px; height: 48px; border-radius: 1rem;
    background: var(--md-sys-color-surface-container-highest);
    display: flex; align-items: center; justify-content: center;
    color: var(--primary-color);
  }
  
  .row-info { flex: 1; min-width: 0; }
  .title-line { display: flex; align-items: center; gap: 0.75rem; margin-bottom: 0.25rem; }
  .gist-title { font-weight: 700; color: var(--md-sys-color-on-surface); font-size: 1.05rem; }
  
  .meta-line { display: flex; align-items: center; gap: 0.5rem; font-size: 0.85rem; color: var(--md-sys-color-on-surface-variant); }
  .sep { opacity: 0.3; }
  
  .row-arrow { color: var(--md-sys-color-on-surface-variant); opacity: 0.3; }
</style>
