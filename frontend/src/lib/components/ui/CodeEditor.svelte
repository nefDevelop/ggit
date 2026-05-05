<script lang="ts">
  import { onMount } from 'svelte';

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

  function syncScroll() {
    if (gutter && textarea) {
      gutter.scrollTop = textarea.scrollTop;
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
      }, 0);
    }
  }
</script>

<div class="code-editor" style="--min-height: {minHeight}">
  <div class="editor-gutter" bind:this={gutter}>
    {#each Array(lineNumbers) as _, i}
      <div class="line-number">{i + 1}</div>
    {/each}
  </div>
  <textarea
    bind:this={textarea}
    bind:value
    {placeholder}
    spellcheck="false"
    onscroll={syncScroll}
    onkeydown={handleKeydown}
  ></textarea>
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
    min-width: 3rem;
  }

  .line-number {
    font-size: 0.85rem;
    line-height: 1.6;
    color: #484f58;
    height: 1.6rem; /* Must match textarea line-height */
  }

  textarea {
    flex: 1;
    background: transparent;
    color: #e6edf3;
    border: none;
    padding: 1.5rem;
    font-family: inherit;
    font-size: 0.95rem;
    line-height: 1.6;
    resize: vertical;
    outline: none;
    white-space: pre;
    overflow-wrap: normal;
    overflow-x: auto;
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
