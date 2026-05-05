<script lang="ts">
  import { apiFetch } from "$lib/api";
  import { onMount } from "svelte";
  import Card from "$lib/components/ui/Card.svelte";
  import Button from "$lib/components/ui/Button.svelte";
  import Badge from "$lib/components/ui/Badge.svelte";
  import { marked } from "marked";
  import hljs from "highlight.js";
  import "highlight.js/styles/github-dark.css";

  let { gist } = $props<{ gist: any }>();
  let files = $derived(gist.files || []);

  function renderMarkdown(content: string) {
    return marked(content);
  }

  onMount(() => {
    hljs.highlightAll();
  });
</script>

<div class="gist-view">
  <header class="gist-header">
    <div class="header-main">
      <div class="title-row">
        <h1>{gist.description || "Untitled Gist"}</h1>
        {#if !gist.is_public}
          <Badge variant="neutral">Secret</Badge>
        {/if}
      </div>
      <div class="meta-row">
        <span class="user">@{gist.user_id}</span>
        <span class="sep">•</span>
        <span class="date">{new Date(gist.created_at).toLocaleDateString()}</span>
        <span class="sep">•</span>
        <Badge variant="secondary" size="sm">{files.length} {files.length === 1 ? 'file' : 'files'}</Badge>
      </div>
    </div>
  </header>

  <div class="gist-files">
    {#each files as file}
      <Card variant="low" padding="none" class="file-card">
        <div class="file-header">
          <span class="material-symbols-rounded">description</span>
          <span class="filename">{file.filename}</span>
          <div class="spacer"></div>
          <Badge variant="ghost">{file.language || 'Plain Text'}</Badge>
        </div>
        <div class="file-content">
          {#if file.filename.endsWith('.md')}
            <div class="markdown-body">
              {@html renderMarkdown(file.content)}
            </div>
          {:else}
            <pre><code class="language-{file.language || 'text'}">{file.content}</code></pre>
          {/if}
        </div>
      </Card>
    {/each}
  </div>
</div>

<style>
  .gist-view { display: flex; flex-direction: column; gap: 2rem; }
  .gist-header { padding-bottom: 1.5rem; border-bottom: 1px solid var(--md-sys-color-outline-variant); }
  .title-row { display: flex; align-items: center; gap: 1rem; margin-bottom: 0.5rem; }
  .title-row h1 { margin: 0; font-size: 1.75rem; font-weight: 800; color: var(--md-sys-color-on-surface); }
  
  .meta-row { display: flex; align-items: center; gap: 0.75rem; color: var(--md-sys-color-on-surface-variant); font-size: 0.9rem; }
  .sep { opacity: 0.3; }

  .gist-files { display: flex; flex-direction: column; gap: 2rem; }
  
  .file-header { 
    display: flex; align-items: center; gap: 0.75rem; padding: 0.75rem 1.25rem; 
    background: var(--md-sys-color-surface-container); border-bottom: 1px solid var(--md-sys-color-outline-variant);
    border-radius: 1.25rem 1.25rem 0 0;
  }
  .file-header span { color: var(--primary-color); }
  .filename { font-weight: 600; color: var(--md-sys-color-on-surface); }
  .spacer { flex: 1; }

  .file-content { padding: 1.5rem; background: var(--md-sys-color-surface-container-low); border-radius: 0 0 1.25rem 1.25rem; }
  
  pre { margin: 0; font-family: 'JetBrains Mono', monospace; font-size: 0.9rem; overflow-x: auto; }
  code { border-radius: 0.5rem; }

  .markdown-body { color: var(--md-sys-color-on-surface); line-height: 1.6; }
</style>
