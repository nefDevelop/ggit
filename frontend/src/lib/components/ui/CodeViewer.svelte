<script lang="ts">
  import { onMount } from 'svelte';
  import hljs from 'highlight.js';
  import 'highlight.js/styles/github-dark.css';

  let { 
    code = '', 
    language = 'text',
    showLineNumbers = true
  } = $props<{
    code: string;
    language?: string;
    showLineNumbers?: boolean;
  }>();

  let lineNumbers = $derived((code.match(/\n/g) || []).length + 1);
  let highlightedCode = $state('');

  function updateHighlight() {
    try {
      const result = hljs.highlight(code, { language: language || 'text' });
      highlightedCode = result.value;
    } catch (e) {
      highlightedCode = code;
    }
  }

  $effect(() => {
    updateHighlight();
  });

  onMount(updateHighlight);
</script>

<div class="code-viewer">
  {#if showLineNumbers}
    <div class="gutter">
      {#each Array(lineNumbers) as _, i}
        <div class="line-number">{i + 1}</div>
      {/each}
    </div>
  {/if}
  <div class="code-container">
    <pre><code class="hljs">{@html highlightedCode}</code></pre>
  </div>
</div>

<style>
  .code-viewer {
    display: flex;
    background: #0d1117;
    border-radius: 0.75rem;
    overflow: hidden;
    font-family: 'JetBrains Mono', 'Fira Code', monospace;
    font-size: 0.9rem;
    line-height: 1.5rem;
  }

  .gutter {
    padding: 1.5rem 0.5rem;
    background: #090c10;
    border-right: 1px solid #30363d;
    display: flex;
    flex-direction: column;
    align-items: flex-end;
    user-select: none;
    min-width: 3.5rem;
    flex-shrink: 0;
  }

  .line-number {
    color: #484f58;
    height: 1.5rem;
    font-size: 0.85rem;
  }

  .code-container {
    flex: 1;
    min-width: 0;
    overflow-x: auto;
    padding: 1.5rem;
  }

  pre {
    margin: 0;
    padding: 0;
    background: transparent;
  }

  code {
    font-family: inherit;
    background: transparent;
    padding: 0;
  }

  /* Scrollbar styling */
  .code-container::-webkit-scrollbar { width: 10px; height: 10px; }
  .code-container::-webkit-scrollbar-track { background: transparent; }
  .code-container::-webkit-scrollbar-thumb { background: #30363d; border-radius: 10px; border: 2px solid #0d1117; }
  .code-container::-webkit-scrollbar-thumb:hover { background: #484f58; }
</style>
