<script lang="ts">
  import { toast } from "$lib/toast.svelte";
  import { page } from '$app/state';

  let { commit, minimalist = false } = $props();

  let copied = $state(false);
  let username = $derived(page.params.username);
  let reponame = $derived(page.params.reponame);

  function timeAgo(dateStr: string) {
    if (!dateStr) return "unknown";
    const date = new Date(dateStr);
    const now = new Date();
    const diff = now.getTime() - date.getTime();
    const seconds = Math.floor(diff / 1000);
    const minutes = Math.floor(seconds / 60);
    const hours = Math.floor(minutes / 60);
    const days = Math.floor(hours / 24);

    if (days > 0) return `${days}d ago`;
    if (hours > 0) return `${hours}h ago`;
    if (minutes > 0) return `${minutes}m ago`;
    return "just now";
  }

  function copyHash(e: Event) {
    e.preventDefault();
    e.stopPropagation();
    if (typeof navigator !== 'undefined' && navigator.clipboard) {
      navigator.clipboard.writeText(commit.id);
      copied = true;
      toast.success("Commit hash copiado al portapapeles", 2000);
      setTimeout(() => copied = false, 2000);
    }
  }
</script>

<a href={`/${username}/${reponame}/commit/${commit.id}`} class="commit-row tonal hover-lift" class:minimalist>
  <div class="commit-avatar">
    {commit.author?.name?.[0] || "U"}
  </div>
  <div class="commit-content">
    <div class="commit-summary">{commit.summary}</div>
    <div class="commit-meta">
      <span class="author">{commit.author?.name}</span>
      <span class="time">{timeAgo(commit.author?.when)}</span>
    </div>
  </div>
  <div class="commit-actions">
    <button class="sha-btn" title="Copiar" onclick={copyHash}>
      <code class="sha">{commit.id?.slice(0, 7)}</code>
    </button>
    {#if !minimalist}
      <button class="btn btn-icon" onclick={copyHash} title="Copy Hash">
        <span class="material-symbols-rounded">{copied ? 'check' : 'content_copy'}</span>
      </button>
    {/if}
  </div>
</a>

<style>
  .commit-row {
    display: flex;
    align-items: center;
    gap: 1rem;
    padding: 0.75rem 1.25rem;
    border: 1px solid var(--md-sys-color-outline-variant);
    margin-bottom: 0.5rem;
    background: transparent;
    text-decoration: none;
    color: inherit;
    cursor: pointer;
  }

  .minimalist {
    padding: 0.5rem 1rem;
    gap: 0.75rem;
    border: none;
    border-bottom: 1px solid var(--md-sys-color-outline-variant);
    border-radius: 0;
    background: transparent;
    box-shadow: none;
  }
  .minimalist:hover {
    background: rgba(255, 255, 255, 0.02);
  }

  .commit-avatar {
    width: 32px;
    height: 32px;
    background: var(--primary-color);
    color: var(--on-primary);
    border-radius: 8px;
    display: flex;
    align-items: center;
    justify-content: center;
    font-weight: 800;
    font-size: 0.9rem;
  }

  .minimalist .commit-avatar {
    width: 24px;
    height: 24px;
    border-radius: 6px;
    font-size: 0.75rem;
  }

  .commit-content {
    flex: 1;
    display: flex;
    flex-direction: column;
    min-width: 0;
  }

  .commit-summary {
    font-weight: 600;
    font-size: 0.95rem;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }

  .minimalist .commit-summary {
    font-size: 0.9rem;
  }

  .commit-meta {
    display: flex;
    gap: 0.75rem;
    font-size: 0.8rem;
    color: var(--text-muted);
  }

  .sha {
    font-family: monospace;
    background: rgba(255, 255, 255, 0.05);
    padding: 0.2rem 0.5rem;
    border-radius: 4px;
    color: var(--primary-color);
    font-size: 0.85rem;
  }

  .sha-btn {
    background: transparent;
    border: none;
    padding: 0;
    cursor: pointer;
    border-radius: 4px;
    transition: all 0.2s;
  }
  .sha-btn:hover .sha {
    background: rgba(var(--primary-rgb), 0.15);
  }

  .commit-actions {
    display: flex;
    align-items: center;
    gap: 0.5rem;
  }
  @media (max-width: 600px) {
    .commit-row {
      flex-wrap: wrap;
      gap: 0.5rem;
    }
    .commit-content {
      order: 1;
      width: calc(100% - 40px);
    }
    .commit-actions {
      order: 2;
      width: 100%;
      justify-content: flex-end;
      padding-top: 0.25rem;
      border-top: 1px solid rgba(255, 255, 255, 0.03);
    }
  }
</style>
