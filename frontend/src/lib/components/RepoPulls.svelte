<script lang="ts">
  import { apiFetch } from "$lib/api";
  import { onMount, getContext } from 'svelte';
  import PageHeader from "./ui/PageHeader.svelte";
  
  const repo: any = getContext('repo');
  let pulls = $state<any[]>([]);
  let loading = $state(true);
  let error = $state('');
  let stateFilter = $state('open'); // 'open' or 'closed'

  async function fetchPulls() {
    loading = true;
    error = '';
    try {
      const res = await apiFetch(`/api/v2/${repo.data.owner.username}/${repo.data.name}/pulls?state=${stateFilter}`);
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
    if (pull.pull_request?.merged) return 'var(--tertiary)'; 
    return pull.is_closed ? 'var(--md-sys-color-error)' : 'var(--primary-color)';
  }
</script>

<div class="pulls-view">
  <PageHeader title="Pull Requests" icon="merge" compact>
    {#snippet actions()}
      <a href="./pulls/new" class="btn-action-primary">
        <span class="material-symbols-rounded">add</span>
        New Pull Request
      </a>
    {/snippet}
  </PageHeader>

  <div class="pill-filter-group">
    <button 
      class="tab-pill" 
      class:active={stateFilter === 'open'} 
      onclick={() => stateFilter = 'open'}
    >
      <span class="material-symbols-rounded mini">merge</span>
      Open
    </button>
    <button 
      class="tab-pill" 
      class:active={stateFilter === 'closed'} 
      onclick={() => stateFilter = 'closed'}
    >
      <span class="material-symbols-rounded mini">history_edu</span>
      Closed
    </button>
  </div>

  {#if loading}
    <div class="skeleton-list"></div>
  {:else if error}
    <div class="error-msg">{error}</div>
  {:else if pulls.length === 0}
    <div class="empty-state">
      <span class="material-symbols-rounded">merge</span>
      <h3>No pull requests found</h3>
      <p>There are no {stateFilter} pull requests in this repository.</p>
    </div>
  {:else}
    <div class="flat-list-container">
      {#each pulls as pull}
        <div class="list-item-row">
          <div class="status-indicator" style="color: {getStatusColor(pull)}">
            <span class="material-symbols-rounded">
              {pull.pull_request?.merged ? 'merge' : (pull.is_closed ? 'cancel' : 'merge')}
            </span>
          </div>
          <div class="item-main-content">
            <div class="item-title-bar">
              <a href="/{repo.data.owner.username}/{repo.data.name}/pulls/{pull.number}" class="item-title-link">
                {pull.title}
              </a>
              {#each pull.labels || [] as label}
                <span class="badge-mini" style="background: {label.color}">
                  {label.name}
                </span>
              {/each}
            </div>
            <div class="item-meta-data">
              #{pull.number} opened {new Date(pull.created_at).toLocaleDateString()} by <strong>{pull.user?.username}</strong>
            </div>
          </div>
          {#if pull.comments > 0}
            <div class="item-comment-badge">
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

  .empty-state {
    display: flex;
    flex-direction: column;
    align-items: center;
    padding: 4rem;
    gap: 1rem;
    color: var(--text-muted);
  }

  @keyframes spin {
    to { transform: rotate(360deg); }
  }
</style>
