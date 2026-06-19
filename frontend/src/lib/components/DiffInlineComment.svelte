<script lang="ts">
  import MarkdownRenderer from "./ui/MarkdownRenderer.svelte";

  let { comment, allowComments = false, onResolve }: {
    comment: any;
    allowComments: boolean;
    onResolve?: (id: number, resolved: boolean) => void;
  } = $props();
</script>

<div class="inline-comment-box tonal" class:resolved={comment.is_resolved}>
  <header class="comment-header">
    <div class="header-left">
      <strong>{comment.reviewer?.username || 'Usuario'}</strong>
      {#if comment.is_resolved}
        <span class="resolved-tag">Resolved</span>
      {/if}
    </div>
    <div class="header-right">
      {#if allowComments}
        <button class="resolve-btn" onclick={() => onResolve?.(comment.id, !comment.is_resolved)}>
          {comment.is_resolved ? 'Unresolve' : 'Resolve'}
        </button>
      {/if}
      <span class="time">{new Date(comment.created_at || Date.now()).toLocaleTimeString()}</span>
    </div>
  </header>
  <div class="comment-body">
    <MarkdownRenderer content={comment.content} />
  </div>
</div>

<style>
  .inline-comment-box {
    margin: 0.5rem 1rem 0.5rem 60px; padding: 1rem;
    border-radius: 0.75rem;
    border: 1px solid var(--md-sys-color-outline-variant);
    background: var(--md-sys-color-surface-container);
  }
  .inline-comment-box.resolved { opacity: 0.7; background: rgba(16,185,129,0.02); }
  .comment-header { display: flex; justify-content: space-between; margin-bottom: 0.5rem; font-size: 0.85rem; }
  .comment-header strong { color: var(--secondary); }
  .comment-header .time { color: var(--text-muted); font-size: 0.75rem; }
  .resolve-btn {
    background: transparent; border: 1px solid var(--md-sys-color-outline);
    border-radius: 4px; color: var(--text-muted);
    font-size: 0.7rem; padding: 2px 6px; cursor: pointer;
    transition: all 0.2s; margin-right: 8px;
  }
  .resolve-btn:hover { background: var(--surface-variant); color: var(--text-color); }
  .resolved-tag {
    background: rgba(16,185,129,0.1); color: #10b981;
    font-size: 0.7rem; font-weight: 700;
    padding: 2px 6px; border-radius: 4px;
    display: flex; align-items: center; gap: 4px; margin-left: 8px;
  }
</style>
