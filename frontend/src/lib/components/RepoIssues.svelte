<script lang="ts">
  import { apiFetch } from "$lib/api";
  import { onMount, getContext } from 'svelte';
  
  const repo: any = getContext('repo');
  let issues = $state<any[]>([]);
  let loading = $state(true);
  let error = $state('');
  let stateFilter = $state('open'); // 'open' or 'closed'

  async function fetchIssues() {
    loading = true;
    error = '';
    try {
      const res = await apiFetch(`/api/v2/${repo.data.owner.username}/${repo.data.name}/issues?state=${stateFilter}`);
      if (res.ok) {
        issues = await res.json();
      } else {
        error = 'Failed to load issues.';
      }
    } catch (e) {
      error = 'Error connecting to API.';
    } finally {
      loading = false;
    }
  }

  $effect(() => {
    if (repo.data) {
      fetchIssues();
    }
  });

  function getStatusColor(issue: any) {
    return issue.is_closed ? '#f87171' : '#10b981';
  }
</script>

<div class="issues-view">
  <div class="issues-header">
    <div class="filter-bar">
      <button 
        class="filter-btn" 
        class:active={stateFilter === 'open'} 
        onclick={() => stateFilter = 'open'}
      >
        <span class="material-symbols-rounded">error</span>
        {repo.data?.open_issues_count || 0} Open
      </button>
      <button 
        class="filter-btn" 
        class:active={stateFilter === 'closed'} 
        onclick={() => stateFilter = 'closed'}
      >
        <span class="material-symbols-rounded">check_circle</span>
        {repo.data?.closed_issues_count || 0} Closed
      </button>
    </div>
    
    <a href="./issues/new" class="btn btn-primary">
      <span class="material-symbols-rounded">add</span>
      New Issue
    </a>
  </div>

  {#if loading}
    <div class="loading-state">
      <div class="spinner"></div>
      <p>Loading issues...</p>
    </div>
  {:else if error}
    <div class="error-msg">{error}</div>
  {:else if issues.length === 0}
    <div class="empty-state">
      <span class="material-symbols-rounded">error_outline</span>
      <h3>No issues found</h3>
      <p>There are no {stateFilter} issues in this repository.</p>
    </div>
  {:else}
    <div class="issues-list island">
      {#each issues as issue}
        <div class="issue-row">
          <div class="issue-icon" style="color: {getStatusColor(issue)}">
            <span class="material-symbols-rounded">
              {issue.is_closed ? 'check_circle' : 'error'}
            </span>
          </div>
          <div class="issue-main">
            <div class="issue-title-row">
              <a href="/{repo.data.owner.username}/{repo.data.name}/issues/{issue.number}" class="issue-title">
                {issue.title}
              </a>
              {#each issue.labels || [] as label}
                <span class="label-tag" style="background: {label.color}; --label-bg: {label.color}">
                  {label.name}
                </span>
              {/each}
            </div>
            <div class="issue-meta">
              #{issue.number} opened {new Date(issue.created_at).toLocaleDateString()} by <strong>{issue.user?.username}</strong>
            </div>
          </div>
          {#if issue.comments > 0}
            <div class="issue-comments-count">
              <span class="material-symbols-rounded">chat_bubble</span>
              {issue.comments}
            </div>
          {/if}
        </div>
      {/each}
    </div>
  {/if}
</div>

<style>
  .issues-view {
    display: flex;
    flex-direction: column;
    gap: 1.5rem;
    animation: fadeIn 0.3s ease-out;
  }

  @keyframes fadeIn {
    from { opacity: 0; transform: translateY(10px); }
    to { opacity: 1; transform: translateY(0); }
  }

  .issues-header {
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

  .filter-btn:hover {
    color: var(--text-color);
    background: rgba(255, 255, 255, 0.05);
  }

  .filter-btn.active {
    color: var(--text-color);
    background: rgba(255, 255, 255, 0.1);
  }

  .issues-list {
    display: flex;
    flex-direction: column;
    overflow: hidden;
  }

  .issue-row {
    display: flex;
    padding: 1rem 1.5rem;
    gap: 1rem;
    border-bottom: 1px solid var(--border-color);
    transition: background 0.2s;
  }

  .issue-row:last-child {
    border-bottom: none;
  }

  .issue-row:hover {
    background: rgba(255, 255, 255, 0.02);
  }

  .issue-icon {
    padding-top: 0.2rem;
  }

  .issue-main {
    flex: 1;
    display: flex;
    flex-direction: column;
    gap: 0.25rem;
  }

  .issue-title-row {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    flex-wrap: wrap;
  }

  .issue-title {
    font-size: 1.05rem;
    font-weight: 700;
    color: var(--text-color);
    text-decoration: none;
  }

  .issue-title:hover {
    color: var(--primary-color);
  }

  .label-tag {
    font-size: 0.75rem;
    font-weight: 700;
    padding: 0.1rem 0.6rem;
    border-radius: 2rem;
    color: white;
    text-shadow: 0 1px 2px rgba(0,0,0,0.3);
    border: 1px solid rgba(255,255,255,0.1);
  }

  .issue-meta {
    font-size: 0.85rem;
    color: var(--text-muted);
  }

  .issue-comments-count {
    display: flex;
    align-items: center;
    gap: 0.4rem;
    font-size: 0.85rem;
    font-weight: 600;
    color: var(--text-muted);
  }

  .issue-comments-count span {
    font-size: 1.1rem;
  }

  .loading-state, .empty-state {
    display: flex;
    flex-direction: column;
    align-items: center;
    padding: 4rem;
    gap: 1rem;
    color: var(--text-muted);
    text-align: center;
  }

  .empty-state span {
    font-size: 3rem;
    opacity: 0.3;
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

  .error-msg {
    color: #f87171;
    padding: 2rem;
    text-align: center;
    background: rgba(248, 113, 113, 0.1);
    border-radius: 1rem;
  }
</style>
