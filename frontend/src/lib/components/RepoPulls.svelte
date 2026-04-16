<script lang="ts">
  import { onMount, getContext } from 'svelte';
  
  const repo: any = getContext('repo');
  let pulls = $state<any[]>([]);
  let loading = $state(true);
  let error = $state('');
  let stateFilter = $state('open'); // 'open' or 'closed'

  async function fetchPulls() {
    loading = true;
    error = '';
    try {
      const res = await fetch(`/api/v2/${repo.data.owner.username}/${repo.data.name}/pulls?state=${stateFilter}`);
      if (res.ok) {
        pulls = await res.json();
      } else {
        error = 'Failed to load pull requests.';
      }
    } catch (e) {
      error = 'Error connecting to API.';
    } finally {
      loading = false;
    }
  }

  $effect(() => {
    if (repo.data) {
      fetchPulls();
    }
  });

  function getStatusColor(pull: any) {
    if (pull.pull_request?.merged) return '#8b5cf6'; // Purple for merged
    return pull.is_closed ? '#f87171' : '#10b981';
  }
</script>

<div class="pulls-view">
  <div class="pulls-header">
    <div class="filter-bar">
      <button 
        class="filter-btn" 
        class:active={stateFilter === 'open'} 
        onclick={() => stateFilter = 'open'}
      >
        <span class="material-symbols-rounded">merge</span>
        Open
      </button>
      <button 
        class="filter-btn" 
        class:active={stateFilter === 'closed'} 
        onclick={() => stateFilter = 'closed'}
      >
        <span class="material-symbols-rounded">history_edu</span>
        Closed
      </button>
    </div>
    
    <a href="./pulls/new" class="btn btn-primary">
      <span class="material-symbols-rounded">add</span>
      New Pull Request
    </a>
  </div>

  {#if loading}
    <div class="loading-state">
      <div class="spinner"></div>
      <p>Loading pull requests...</p>
    </div>
  {:else if error}
    <div class="error-msg">{error}</div>
  {:else if pulls.length === 0}
    <div class="empty-state">
      <span class="material-symbols-rounded">merge</span>
      <h3>No pull requests found</h3>
      <p>There are no {stateFilter} pull requests in this repository.</p>
    </div>
  {:else}
    <div class="pulls-list island">
      {#each pulls as pull}
        <div class="pull-row">
          <div class="pull-icon" style="color: {getStatusColor(pull)}">
            <span class="material-symbols-rounded">
              {pull.pull_request?.merged ? 'merge' : (pull.is_closed ? 'cancel' : 'merge')}
            </span>
          </div>
          <div class="pull-main">
            <div class="pull-title-row">
              <a href="/{repo.data.owner.username}/{repo.data.name}/pulls/{pull.number}" class="pull-title">
                {pull.title}
              </a>
              {#each pull.labels || [] as label}
                <span class="label-tag" style="background: {label.color}">
                  {label.name}
                </span>
              {/each}
            </div>
            <div class="pull-meta">
              #{pull.number} opened {new Date(pull.created_at).toLocaleDateString()} by <strong>{pull.user?.username}</strong>
            </div>
          </div>
          {#if pull.comments > 0}
            <div class="pull-comments-count">
              <span class="material-symbols-rounded">chat_bubble</span>
              {pull.comments}
            </div>
          {/if}
        </div>
      {/each}
    </div>
  {/if}
</div>

<style>
  .pulls-view {
    display: flex;
    flex-direction: column;
    gap: 1.5rem;
    animation: fadeIn 0.3s ease-out;
  }

  @keyframes fadeIn {
    from { opacity: 0; transform: translateY(10px); }
    to { opacity: 1; transform: translateY(0); }
  }

  .pulls-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }

  .filter-bar {
    display: flex;
    gap: 0.5rem;
    background: rgba(255, 255, 255, 0.05);
    padding: 0.25rem;
    border-radius: 0.75rem;
  }

  .filter-btn {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    padding: 0.5rem 1rem;
    border: none;
    background: transparent;
    color: var(--text-muted);
    font-weight: 600;
    font-size: 0.9rem;
    cursor: pointer;
    border-radius: 0.5rem;
    transition: all 0.2s;
  }

  .filter-btn.active {
    color: var(--text-color);
    background: rgba(255, 255, 255, 0.1);
  }

  .pulls-list {
    display: flex;
    flex-direction: column;
    overflow: hidden;
  }

  .pull-row {
    display: flex;
    padding: 1rem 1.5rem;
    gap: 1rem;
    border-bottom: 1px solid var(--border-color);
  }

  .pull-row:last-child {
    border-bottom: none;
  }

  .pull-main {
    flex: 1;
    display: flex;
    flex-direction: column;
    gap: 0.25rem;
  }

  .pull-title {
    font-size: 1.05rem;
    font-weight: 700;
    color: var(--text-color);
    text-decoration: none;
  }

  .pull-title:hover {
    color: var(--primary-color);
  }

  .label-tag {
    font-size: 0.75rem;
    font-weight: 700;
    padding: 0.1rem 0.6rem;
    border-radius: 2rem;
    color: white;
  }

  .pull-meta {
    font-size: 0.85rem;
    color: var(--text-muted);
  }

  .pull-comments-count {
    display: flex;
    align-items: center;
    gap: 0.4rem;
    font-size: 0.85rem;
    color: var(--text-muted);
  }

  .loading-state, .empty-state {
    display: flex;
    flex-direction: column;
    align-items: center;
    padding: 4rem;
    gap: 1rem;
    color: var(--text-muted);
  }

  .spinner {
    width: 32px;
    height: 32px;
    border: 3px solid rgba(255,255,255,0.1);
    border-top-color: var(--primary-color);
    border-radius: 50%;
    animation: spin 1s infinite linear;
  }

  @keyframes spin {
    to { transform: rotate(360deg); }
  }
</style>
