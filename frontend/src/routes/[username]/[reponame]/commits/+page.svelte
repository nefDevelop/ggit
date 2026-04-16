<script lang="ts">
  import { page } from '$app/state';
  import { onMount } from 'svelte';

  let username = $derived(page.params.username);
  let reponame = $derived(page.params.reponame);

  let commitsData = $state<any>(null);
  let loading = $state(true);
  let error = $state('');

  async function fetchCommits() {
    loading = true;
    error = '';
    try {
      const res = await fetch(`/api/v2/${username}/${reponame}/commits`);
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

  function formatDate(timestamp: number) {
    return new Date(timestamp * 1000).toLocaleDateString(undefined, {
      year: 'numeric',
      month: 'short',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit'
    });
  }

  function truncateSHA(sha: string) {
    return sha.substring(0, 7);
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
    <div class="island error-island">
      <span class="material-symbols-rounded">error</span>
      <p>{error}</p>
      <button class="btn btn-tonal" onclick={fetchCommits}>Retry</button>
    </div>
  {:else if commitsData}
    <div class="commits-list island">
      {#each commitsData.commits as c}
        <div class="commit-item">
          <div class="commit-avatar">
            {#if c.user?.avatar_url}
              <img src={c.user.avatar_url} alt={c.Author.Name} />
            {:else}
              <div class="avatar-placeholder">
                {c.Author.Name.charAt(0).toUpperCase()}
              </div>
            {/if}
          </div>
          
          <div class="commit-info">
            <div class="commit-header">
              <span class="commit-message">{c.Message.split('\n')[0]}</span>
              {#if c.Message.split('\n').slice(1).join('\n').trim()}
                <span class="material-symbols-rounded more-icon" title="View extended message">more_horiz</span>
              {/if}
            </div>
            
            <div class="commit-meta">
              <span class="author-name">{c.Author.Name}</span>
              <span class="text-muted">committed on</span>
              <span class="commit-date">{formatDate(c.Author.When)}</span>
            </div>
          </div>

          <div class="commit-actions">
            <code class="sha-badge">{truncateSHA(c.ID)}</code>
            <button class="icon-btn">
              <span class="material-symbols-rounded">code</span>
            </button>
          </div>
        </div>
      {/each}
    </div>

    {#if commitsData.commits.length === 0}
      <div class="empty-state">
        <span class="material-symbols-rounded">history</span>
        <p>No commits found in this repository.</p>
      </div>
    {/if}
  {/if}
</div>

<style>
  .commits-view {
    display: flex;
    flex-direction: column;
    gap: 1.5rem;
  }

  .commits-list {
    padding: 0;
    overflow: hidden;
  }

  .commit-item {
    display: flex;
    align-items: center;
    padding: 1rem 1.5rem;
    gap: 1.25rem;
    border-bottom: 1px solid var(--border-color);
    transition: background 0.2s;
  }

  .commit-item:last-child {
    border-bottom: none;
  }

  .commit-item:hover {
    background-color: rgba(255, 255, 255, 0.03);
  }

  .commit-avatar img, .avatar-placeholder {
    width: 40px;
    height: 40px;
    border-radius: 50%;
    object-fit: cover;
  }

  .avatar-placeholder {
    background: var(--md3-surface-2);
    display: flex;
    align-items: center;
    justify-content: center;
    font-weight: 700;
    color: var(--primary-color);
  }

  .commit-info {
    flex: 1;
    display: flex;
    flex-direction: column;
    gap: 0.25rem;
    min-width: 0;
  }

  .commit-header {
    display: flex;
    align-items: center;
    gap: 0.5rem;
  }

  .commit-message {
    font-weight: 600;
    font-size: 1.05rem;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }

  .more-icon {
    font-size: 1.2rem;
    color: var(--text-muted);
    cursor: pointer;
  }

  .commit-meta {
    font-size: 0.9rem;
    display: flex;
    gap: 0.4rem;
    align-items: center;
  }

  .author-name {
    font-weight: 700;
    color: var(--text-color);
  }

  .commit-actions {
    display: flex;
    align-items: center;
    gap: 1rem;
  }

  .sha-badge {
    font-family: 'JetBrains Mono', monospace;
    background: var(--md3-surface-2);
    padding: 0.3rem 0.6rem;
    border-radius: 0.5rem;
    font-size: 0.85rem;
    color: var(--primary-color);
    border: 1px solid var(--border-color);
  }

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
  }

  .empty-state .material-symbols-rounded {
    font-size: 3rem;
    margin-bottom: 1rem;
    opacity: 0.3;
  }

  .error-island {
    padding: 3rem;
    text-align: center;
  }
</style>
