<script lang="ts">
  import type { RepoContext } from "$lib/types";
  import { apiFetch } from "$lib/api";
  import { onMount, getContext } from 'svelte';
  import { formatDate } from "$lib/utils/date";
  import PageActions from "./ui/PageActions.svelte";
  import EmptyState from "./ui/EmptyState.svelte";
  import Button from "./ui/Button.svelte";
  

import Plus from "lucide-svelte/icons/plus";
import GitPullRequest from "lucide-svelte/icons/git-pull-request";


  
  const repo = getContext<RepoContext>('repo');
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
  <PageActions>
      <div class="pill-filter-group-inline">
        <Button 
          variant={stateFilter === 'open' ? 'primary' : 'ghost'} 
          size="sm"
          onclick={() => stateFilter = 'open'}
        >
          Open
        </Button>
        <Button 
          variant={stateFilter === 'closed' ? 'primary' : 'ghost'} 
          size="sm"
          onclick={() => stateFilter = 'closed'}
        >
          Closed
        </Button>
      </div>

      {#if pulls.length > 0 || stateFilter === 'closed'}
        <Button variant="primary" href="./pulls/new" icon={Plus}>
          New Pull Request
        </Button>
      {/if}
    </PageActions>


  {#if loading}
    <div class="skeleton-list"></div>
  {:else if error}
    <div class="error-msg">{error}</div>
  {:else if pulls.length === 0}
    <EmptyState
      icon={GitPullRequest}
      title="No pull requests found"
      description={`There are no ${stateFilter} pull requests in this repository.`}
    >
      {#if stateFilter === 'open'}
        <Button variant="primary" href="./pulls/new" icon={Plus}>
          Create your first pull request
        </Button>
      {/if}
    </EmptyState>
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
              {#if pull.combined_status}
                <span class="ci-status-icon {pull.combined_status}">
                  <span class="material-symbols-rounded mini">
                    {pull.combined_status === 'success' ? 'check_circle' : (pull.combined_status === 'pending' ? 'pending' : 'cancel')}
                  </span>
                </span>
              {/if}
              <a href="/{repo.data.owner.username}/{repo.data.name}/pulls/{pull.number}" class="item-title-link">
                {pull.title}
              </a>
              {#each pull.labels || [] as label}
                <span class="badge-mini" style="background: {label.color}" class:scoped={label.name.includes('::')}>
                  {#if label.name.includes('::')}
                    {@const [key, val] = label.name.split('::')}
                    <span class="key">{key}</span>
                    <span class="val">{val}</span>
                  {:else}
                    {label.name}
                  {/if}
                </span>
              {/each}
            </div>
            <div class="item-meta-data">
              #{pull.number} opened {formatDate(pull.created_at)} by <strong>{pull.user?.username}</strong>
              {#if pull.attention_state === 1 || pull.attention_state === 2}
                <span class="attention-dot" class:author={pull.attention_state === 2} class:reviewer={pull.attention_state === 1} title={pull.attention_state === 1 ? 'Review required' : 'Waiting for author'}></span>
              {/if}
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
    gap: 0;
    animation: fadeIn 0.3s ease-out;
    padding: 0;
  }

  .pill-filter-group-inline {
    display: flex;
    background: var(--surface-variant);
    padding: 4px;
    border-radius: 10px;
    margin-right: 1rem;
  }
  
  .tab-pill {
    background: transparent;
    border: none;
    padding: 0.4rem 0.8rem;
    border-radius: 8px;
    font-size: 0.85rem;
    font-weight: 700;
    color: var(--text-muted);
    cursor: pointer;
    transition: all 0.2s;
  }
  
  .tab-pill.active {
    background: var(--bg-color);
    color: var(--primary-color);
    box-shadow: var(--shadow-sm);
  }

  @keyframes fadeIn {
    from { opacity: 0; transform: translateY(10px); }
    to { opacity: 1; transform: translateY(0); }
  }





  @keyframes spin {
    to { transform: rotate(360deg); }
  }

  .ci-status-icon { display: inline-flex; align-items: center; vertical-align: middle; }
  .ci-status-icon.success { color: var(--status-success); }
  .ci-status-icon.pending { color: var(--color-warning); }
  .ci-status-icon.failure, .ci-status-icon.error { color: var(--error); }

  .attention-dot {
    display: inline-block;
    width: 8px;
    height: 8px;
    border-radius: 50%;
    margin-left: 0.5rem;
    vertical-align: middle;
  }
  .attention-dot.reviewer { background: var(--status-info-strong); box-shadow: 0 0 5px color-mix(in srgb, var(--status-info) 50%, transparent); }
  .attention-dot.author { background: var(--md-sys-color-error); box-shadow: 0 0 5px rgba(186, 26, 26, 0.5); }
</style>
