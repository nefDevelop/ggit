<script lang="ts">
  import { apiFetch } from "$lib/api";
  import { onMount, getContext } from 'svelte';
  import { page } from '$app/state';
  import DiffViewer from '$lib/components/DiffViewer.svelte';
  import Avatar from '$lib/components/Avatar.svelte';
  import CommitStatusIcon from '$lib/components/CommitStatusIcon.svelte';
  
import ExternalLink from "lucide-svelte/icons/external-link";
import CheckCircle2 from "lucide-svelte/icons/check-circle-2";
import XCircle from "lucide-svelte/icons/x-circle";
import AlertCircle from "lucide-svelte/icons/alert-circle";


  const repo: any = getContext('repo');
  let sha = $derived(page.params.sha);

  let commitData = $state<any>(null);
  let diffData = $state<any>(null);
  let loading = $state(true);
  let error = $state('');

  async function fetchCommit() {
    loading = true;
    error = '';
    try {
      const res = await apiFetch(`/api/v2/${repo.data.owner.username}/${repo.data.name}/commit/${sha}`);
      if (res.ok) {
        const data = await res.json();
        commitData = data.commit;
        diffData = data.diff;
      } else {
        error = 'Commit not found.';
      }
    } catch (e) {
      error = 'Failed to load commit details.';
    } finally {
      loading = false;
    }
  }

  $effect(() => {
    if (repo.data && sha) {
      fetchCommit();
    }
  });

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
</script>

<div class="commit-view-container">
  {#if loading}
    <div class="skeleton-header"></div>
    <div class="skeleton-body"></div>
  {:else if error}
    <div class="error-msg">{error}</div>
  {:else if commitData}
    <div class="commit-header tonal">
      <div class="commit-title-row">
        <h2>{commitData.summary}</h2>
        <code class="sha">{commitData.id?.slice(0, 7)}</code>
      </div>
      {#if commitData.message && commitData.message !== commitData.summary}
        <pre class="commit-message">{commitData.message}</pre>
      {/if}
      <div class="commit-meta">
        <div class="author-info">
          <Avatar
            src={commitData.author?.avatar_url || ''}
            username={commitData.author?.name || 'U'}
            seed={commitData.author?.email || commitData.author?.name || 'unknown'}
            size={32}
            type="user"
          />
          <span><strong>{commitData.author?.name}</strong> committed {timeAgo(commitData.author?.when)}</span>
        </div>
        
        {#if commitData.status}
          <div class="status-summary">
            <CommitStatusIcon status={commitData.status} size={20} />
            <span class="status-label">{commitData.status}</span>
          </div>
        {/if}
      </div>

      {#if commitData.statuses && commitData.statuses.length > 0}
        <div class="statuses-list island tonal">
          <header class="statuses-header">Commit Statuses</header>
          {#each commitData.statuses as s}
            <div class="status-detail-row">
              <CommitStatusIcon status={s.State} size={18} />
              <div class="status-text-info">
                <span class="status-context"><b>{s.Context}</b> — {s.Description || 'No description'}</span>
              </div>
              {#if s.TargetURL}
                <a href={s.TargetURL} target="_blank" class="status-link">
                  Details <ExternalLink size={14} />
                </a>
              {/if}
            </div>
          {/each}
        </div>
      {/if}
    </div>

    {#if diffData}
      <DiffViewer diff={diffData} />
    {/if}
  {/if}
</div>

<style>
  .commit-view-container {
    display: flex;
    flex-direction: column;
    gap: 1.5rem;
    max-width: 1200px;
    margin: 0 auto;
    width: 100%;
    animation: fadeIn 0.4s ease-out;
  }

  .commit-header {
    padding: 1.5rem;
    border-radius: 1.5rem;
    border: 1px solid var(--md-sys-color-outline-variant);
    display: flex;
    flex-direction: column;
    gap: 1rem;
  }

  .commit-title-row {
    display: flex;
    justify-content: space-between;
    align-items: flex-start;
    gap: 1rem;
  }

  .commit-title-row h2 {
    margin: 0;
    font-size: 1.5rem;
    font-weight: 800;
    color: var(--md-sys-color-on-surface);
    word-break: break-word;
  }

  .sha {
    font-family: monospace;
    background: rgba(var(--primary-rgb), 0.1);
    color: var(--primary-color);
    padding: 0.4rem 0.8rem;
    border-radius: 8px;
    font-size: 0.9rem;
    font-weight: bold;
    flex-shrink: 0;
  }

  .commit-message {
    margin: 0;
    white-space: pre-wrap;
    font-family: inherit;
    font-size: 0.95rem;
    color: var(--md-sys-color-on-surface-variant);
    background: rgba(0,0,0,0.02);
    padding: 1rem;
    border-radius: 12px;
    border-left: 4px solid var(--md-sys-color-outline-variant);
  }

  .commit-meta {
    display: flex;
    align-items: center;
    gap: 1rem;
    padding-top: 1rem;
    border-top: 1px solid var(--md-sys-color-outline-variant);
  }

  .author-info {
    display: flex;
    align-items: center;
    gap: 0.75rem;
    font-size: 0.9rem;
    color: var(--md-sys-color-on-surface-variant);
  }

  .status-summary {
    margin-left: auto;
    display: flex;
    align-items: center;
    gap: 0.5rem;
    background: var(--surface-variant);
    padding: 0.4rem 1rem;
    border-radius: 20px;
    border: 1px solid var(--md-sys-color-outline-variant);
  }
  .status-label { font-weight: 800; font-size: 0.85rem; text-transform: capitalize; }

  .statuses-list {
    display: flex;
    flex-direction: column;
    margin-top: 1rem;
    border: 1px solid var(--md-sys-color-outline-variant);
    border-radius: 12px;
    overflow: hidden;
  }
  .statuses-header {
    background: var(--surface-variant);
    padding: 0.5rem 1rem;
    font-size: 0.75rem;
    font-weight: 800;
    text-transform: uppercase;
    letter-spacing: 0.05em;
    color: var(--text-muted);
    border-bottom: 1px solid var(--md-sys-color-outline-variant);
  }
  .status-detail-row {
    display: flex;
    align-items: center;
    gap: 1rem;
    padding: 0.75rem 1rem;
    border-bottom: 1px solid var(--md-sys-color-outline-variant);
  }
  .status-detail-row:last-child { border-bottom: none; }
  .status-text-info { flex: 1; min-width: 0; }
  .status-context { font-size: 0.9rem; }
  .status-link {
    display: flex;
    align-items: center;
    gap: 0.4rem;
    font-size: 0.8rem;
    font-weight: 700;
    color: var(--primary-color);
    text-decoration: none;
  }
  .status-link:hover { text-decoration: underline; }



  .error-msg {
    padding: 2rem;
    text-align: center;
    color: var(--md-sys-color-error);
    background: rgba(var(--error-rgb), 0.1);
    border-radius: 1rem;
  }

  .skeleton-header {
    height: 120px;
    background: var(--md-sys-color-surface-variant);
    border-radius: 1.5rem;
    animation: pulse 1.5s infinite;
  }

  .skeleton-body {
    height: 400px;
    background: var(--md-sys-color-surface-variant);
    border-radius: 1.5rem;
    animation: pulse 1.5s infinite;
  }

  @keyframes pulse {
    0%, 100% { opacity: 0.5; }
    50% { opacity: 0.8; }
  }

  @keyframes fadeIn {
    from { opacity: 0; transform: translateY(10px); }
    to { opacity: 1; transform: translateY(0); }
  }
</style>
