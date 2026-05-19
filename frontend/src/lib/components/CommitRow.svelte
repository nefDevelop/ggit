<script lang="ts">
  import { toast } from "$lib/toast.svelte";
  import { page } from '$app/state';
  import { Check, Copy } from "lucide-svelte";
  import Avatar from './Avatar.svelte';

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
  // commit.user is the matched Ggit account (has real AvatarURL).
  // commit.author.avatar_url is set for last_commit bars.
  const authorAvatarUrl = $derived(
    commit.user?.avatar_url || commit.author?.avatar_url || ''
  );
</script>

<a href={`/${username}/${reponame}/commit/${commit.id}`} class="commit-row tonal hover-lift" class:minimalist>
  <Avatar
    src={authorAvatarUrl}
    username={commit.author?.name || 'U'}
    seed={commit.author?.email || commit.author?.name || 'unknown'}
    size={32}
    type="user"
  />
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
        {#if copied}
          <Check size={16} />
        {:else}
          <Copy size={16} />
        {/if}
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
    background: var(--md-sys-color-tertiary-container);
    padding: 0.2rem 0.5rem;
    border-radius: 4px;
    color: var(--md-sys-color-on-tertiary-container);
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
