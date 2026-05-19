<script lang="ts">
  import { onMount, tick } from 'svelte';
  import hljs from 'highlight.js';
  import 'highlight.js/styles/github-dark.css';

  let { 
    value = $bindable(), 
    placeholder = '', 
    language = 'text',
    minHeight = '300px'
  } = $props<{
    value: string;
    placeholder?: string;
    language?: string;
    minHeight?: string;
  }>();

  let lineNumbers = $derived((value.match(/\n/g) || []).length + 1);
  let textarea: HTMLTextAreaElement;
  let gutter: HTMLDivElement;
  let backdrop: HTMLDivElement;

  function syncScroll() {
    if (textarea) {
      if (gutter) gutter.scrollTop = textarea.scrollTop;
      if (backdrop) {
        backdrop.scrollTop = textarea.scrollTop;
        backdrop.scrollLeft = textarea.scrollLeft;
      }
    }
  }

  function handleKeydown(e: KeyboardEvent) {
    if (e.key === 'Tab') {
      e.preventDefault();
      const start = textarea.selectionStart;
      const end = textarea.selectionEnd;
      value = value.substring(0, start) + '  ' + value.substring(end);
      setTimeout(() => {
        textarea.selectionStart = textarea.selectionEnd = start + 2;
        updateHighlight();
      }, 0);
    }
  }

  let highlightedCode = $state('');

  function updateHighlight() {
    const code = value + (value.endsWith('\n') ? ' ' : '');
    try {
      const result = hljs.highlight(code, { language: language || 'text' });
      highlightedCode = result.value;
    } catch (e) {
      // Fallback to auto-detect or plain text if language is unknown/fails
      try {
        const result = hljs.highlightAuto(code);
        highlightedCode = result.value;
      } catch (e2) {
        highlightedCode = code.replace(/&/g, "&amp;").replace(/</g, "&lt;").replace(/>/g, "&gt;");
      }
    }
  }

  $effect(() => {
    updateHighlight();
  });

  onMount(() => {
    updateHighlight();
  });
</script>

<div class="code-editor" style="--min-height: {minHeight}">
  <div class="editor-gutter" bind:this={gutter}>
    {#each Array(lineNumbers) as _, i}
      <div class="line-number">{i + 1}</div>
    {/each}
  </div>
  <div class="editor-container">
    <div class="backdrop" bind:this={backdrop}>
      <pre><code class="hljs">{@html highlightedCode}</code></pre>
    </div>
    <textarea
      bind:this={textarea}
      bind:value
      {placeholder}
      spellcheck="false"
      onscroll={syncScroll}
      onkeydown={handleKeydown}
      oninput={updateHighlight}
    ></textarea>
  </div>
</div>

<style>
  .code-editor {
    display: flex;
    background: #0d1117;
    border: 1px solid var(--md-sys-color-outline-variant);
    border-radius: 1rem;
    overflow: hidden;
    min-height: var(--min-height);
    font-family: 'JetBrains Mono', 'Fira Code', monospace;
    position: relative;
  }

  .editor-gutter {
    padding: 1.5rem 0.5rem;
    background: #090c10;
    border-right: 1px solid #30363d;
    display: flex;
    flex-direction: column;
    align-items: flex-end;
    user-select: none;
    overflow: hidden;
    min-width: 3.5rem;
    flex-shrink: 0;
  }

  .line-number {
    font-size: 0.85rem;
    line-height: 1.5rem;
    color: #484f58;
    height: 1.5rem;
    font-family: inherit;
  }

  .editor-container {
    position: relative;
    flex: 1;
    min-width: 0;
  }

  .backdrop, textarea {
    width: 100%;
    height: 100%;
    margin: 0;
    padding: 1.5rem;
    box-sizing: border-box;
    font-family: inherit;
    font-size: 0.95rem;
    line-height: 1.5rem;
    tab-size: 2;
    white-space: pre;
    word-break: normal;
    overflow-wrap: normal;
  }

  .backdrop {
    position: absolute;
    top: 0;
    left: 0;
    pointer-events: none;
    color: #e6edf3; /* Base color for unhighlighted text */
    overflow: hidden;
    z-index: 1;
  }

  .backdrop pre, .backdrop code {
    margin: 0;
    padding: 0;
    font-family: inherit;
    font-size: inherit;
    line-height: inherit;
    background: transparent;
  }

  textarea {
    position: relative;
    z-index: 2;
    background: transparent;
    color: transparent;
    caret-color: #e6edf3;
    border: none;
    resize: none;
    outline: none;
    display: block;
    overflow: auto;
  }

  /* To make text visible but allow highlight to show through */
  textarea {
    color: rgba(230, 237, 243, 0.1); /* Very faint so we can see the highlight below */
    /* Actually better: set textarea color to transparent and rely on backdrop for everything except caret */
    color: transparent;
  }

  textarea::selection {
    background: rgba(var(--primary-rgb), 0.3);
  }

  /* Scrollbar styling */
  textarea::-webkit-scrollbar { width: 10px; height: 10px; }
  textarea::-webkit-scrollbar-track { background: transparent; }
  textarea::-webkit-scrollbar-thumb { background: #30363d; border-radius: 10px; border: 2px solid #0d1117; }
  textarea::-webkit-scrollbar-thumb:hover { background: #484f58; }
</style>
