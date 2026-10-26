<script lang="ts">
  import { apiFetch } from "$lib/api";
  import { page } from '$app/state';
  import { onMount } from 'svelte';
  import Avatar from "$lib/components/Avatar.svelte";
  import CommitStatusIcon from "$lib/components/CommitStatusIcon.svelte";

  let username = $derived(page.params.username);
  let reponame = $derived(page.params.reponame);

  let commitsData = $state<any>(null);
  let loading = $state(true);
  let error = $state('');

  async function fetchCommits() {
    loading = true;
    error = '';
    try {
      const res = await apiFetch(`/api/v2/${username}/${reponame}/commits/`);
      if (res.ok) {
        commitsData = await res.json();
      } else {
        error = 'Failed to load commits history.';
      }
    } catch (e) {
      error = 'An error occurred while fetching commits.';
    } finally {
      loading = false;
    }
  }

  onMount(() => {
    fetchCommits();
  });

  function formatDate(dateValue: any) {
    if (!dateValue) return 'Unknown';
    const d = typeof dateValue === 'number' ? new Date(dateValue * 1000) : new Date(dateValue);
    return d.toLocaleDateString(undefined, {
      year: 'numeric',
      month: 'short',
      day: 'numeric'
    });
  }

  function truncateSHA(sha: any) {
    if (!sha) return '-------';
    return String(sha).substring(0, 7);
  }
</script>

<div class="commits-view">
  {#if loading}
    <div class="skeleton-list">
      {#each Array(5) as _}
        <div class="skeleton-item"></div>
      {/each}
    </div>
  {:else if error}
    <div class="tonal error-tonal">
      <span class="material-symbols-rounded">error</span>
      <p>{error}</p>
      <button class="btn btn-tonal" onclick={fetchCommits}>Retry</button>
    </div>
  {:else if commitsData}
    <div class="commits-list">
      {#each commitsData.commits as c}
        {@const author = c.Author || c.author}
        {@const message = c.Message || c.message}
        {@const id = c.ID || c.id || c.sha}
        
        <div class="commit-row">
          <!-- Left: User Identity -->
          <div class="commit-avatar-cell">
            <Avatar 
              src={c.user?.avatar_url} 
              username={author?.Name || author?.name || "User"} 
              seed={author?.Name || author?.name}
              size={32} 
              class="user-avatar"
            />
          </div>
          
          <!-- Middle: Message & Meta -->
          <div class="commit-info-cell">
            <div class="commit-header">
              <CommitStatusIcon status={c.status} size={18} class="status-icon-margin" />
              <a href={`/${username}/${reponame}/commit/${id}`} class="commit-message-link">
                <span class="commit-message">{String(message || '').split('\n')[0]}</span>
              </a>
              {#if String(message || '').split('\n').length > 1}
                <span class="material-symbols-rounded more-icon">more_horiz</span>
              {/if}
            </div>
            
            <div class="commit-meta">
              <span class="author-name">{author?.Name || author?.name || 'Anonymous'}</span>
              <span class="text-sep">committed on</span>
              <span class="commit-date">{formatDate(author?.When || author?.when)}</span>
            </div>
          </div>

          <!-- Right: Technical context -->
          <div class="commit-actions-cell">
            <a href={`/${username}/${reponame}/commit/${id}`} class="id-tag-link">
              <div class="id-tag">
                <span class="material-symbols-rounded mini">fingerprint</span>
                <code class="sha-text">{truncateSHA(id)}</code>
              </div>
            </a>
            <a href={`/api/v2/${username}/${reponame}/archive/${id}.zip`} class="btn-icon" download title="Download ZIP at this commit">
              <span class="material-symbols-rounded">folder_zip</span>
            </a>
          </div>
        </div>
      {/each}
    </div>

    {#if commitsData.commits.length === 0}
      <div class="empty-state">
        <span class="material-symbols-rounded">error_outline</span>
        <h3>No commits found</h3>
        <p>There are no commits in this repository yet.</p>
      </div>
    {/if}
  {/if}
</div>

<style>
  .commits-view {
    display: flex;
    flex-direction: column;
    gap: var(--space-5);
  }

  .commits-list {
    display: flex;
    flex-direction: column;
    border: 1px solid var(--border-color);
    background: var(--bg-color);
    border-radius: 1rem;
    overflow: hidden;
  }

  .commit-row {
    display: flex;
    align-items: center;
    padding: 0.85rem 1.5rem;
    gap: 1.25rem;
    border-bottom: 1px solid var(--border-color);
    transition: background 0.2s;
  }

  .commit-row:last-child { border-bottom: none; }
  .commit-row:nth-child(even) { background: var(--sunken-1); }
  .commit-row:hover { background: rgba(var(--primary-rgb), 0.04); }

  :global(.user-avatar) {
    width: 32px;
    height: 32px;
    border-radius: 50%;
    object-fit: cover;
    display: block;
    border: 1px solid var(--border-color);
  }

  .commit-info-cell {
    flex: 1;
    display: flex;
    flex-direction: column;
    min-width: 0;
  }

  :global(.status-icon-margin) {
    margin-right: 0.75rem;
    align-self: center;
  }

  .commit-message { font-weight: 700; font-size: 0.95rem; color: var(--text-color); white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }
  .more-icon { font-size: 1rem; color: var(--text-muted); opacity: 0.5; }

  .commit-meta { font-size: 0.8rem; display: flex; gap: 0.5rem; align-items: center; color: var(--text-muted); margin-top: 2px; }
  .author-name { font-weight: 800; color: var(--text-color); }
  .text-sep { opacity: 0.5; font-weight: 500; }

  .commit-actions-cell { display: flex; align-items: center; gap: 1rem; }

  .id-tag {
    display: flex;
    align-items: center;
    gap: 5px;
    background: var(--surface-variant);
    padding: 2px 10px;
    border-radius: 20px;
    border: 1px solid var(--border-color);
    color: var(--primary-color);
    transition: all 0.2s;
  }
  .id-tag:hover { background: var(--primary-color); color: white; }
  .id-tag .mini { font-size: 0.9rem; opacity: 0.7; }
  .sha-text { font-family: 'JetBrains Mono', monospace; font-size: 0.75rem; font-weight: 800; }

  .commit-message-link, .id-tag-link {
    text-decoration: none;
    color: inherit;
  }
  .commit-message-link:hover .commit-message {
    color: var(--primary-color);
    text-decoration: underline;
  }

  .btn-icon {
    background: transparent;
    border: none;
    color: var(--text-muted);
    cursor: pointer;
    padding: 6px;
    border-radius: 8px;
    display: flex;
    transition: all 0.2s;
  }
  .btn-icon:hover { background: var(--surface-variant); color: var(--primary-color); }

  /* Skeletons */
  .skeleton-list {
    display: flex;
    flex-direction: column;
    gap: 1rem;
  }

  .skeleton-item {
    height: 80px;
    background: var(--md3-surface-2);
    border-radius: var(--radius);
    animation: pulse 1.5s infinite;
  }

  @keyframes pulse {
    0% { opacity: 0.6; }
    50% { opacity: 0.3; }
    100% { opacity: 0.6; }
  }

  .empty-state {
    text-align: center;
    padding: 4rem;
    color: var(--text-muted);
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 1rem;
  }

  .empty-state .material-symbols-rounded {
    font-size: 3rem;
    opacity: 0.3;
  }

  .error-tonal {
    padding: 3rem;
    text-align: center;
  }
</style>
