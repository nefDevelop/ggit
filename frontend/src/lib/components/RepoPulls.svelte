<script lang="ts">
  import { apiFetch } from "$lib/api";
  import { onMount, getContext } from 'svelte';
  import PageHeader from "./ui/PageHeader.svelte";
  import Button from "./ui/Button.svelte";
  import { Plus } from "lucide-svelte";
  
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
  <PageHeader title="" icon="" compact>
    {#snippet actions()}
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

      <Button variant="primary" href="./pulls/new" icon={Plus}>
        New Pull Request
      </Button>    {/snippet}
  </PageHeader>


  {#if loading}
    <div class="skeleton-list"></div>
  {:else if error}
    <div class="error-msg">{error}</div>
  {:else if pulls.length === 0}
    <div class="empty-state">
      <span class="material-symbols-rounded">error_outline</span>
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
              #{pull.number} opened {new Date(pull.created_at).toLocaleDateString()} by <strong>{pull.user?.username}</strong>
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
    margin: -2rem;
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

  .ci-status-icon { display: inline-flex; align-items: center; vertical-align: middle; }
  .ci-status-icon.success { color: #10b981; }
  .ci-status-icon.pending { color: #f59e0b; }
  .ci-status-icon.failure, .ci-status-icon.error { color: #ef4444; }

  .attention-dot {
    display: inline-block;
    width: 8px;
    height: 8px;
    border-radius: 50%;
    margin-left: 0.5rem;
    vertical-align: middle;
  }
  .attention-dot.reviewer { background: #0061a4; box-shadow: 0 0 5px rgba(0, 97, 164, 0.5); }
  .attention-dot.author { background: #ba1a1a; box-shadow: 0 0 5px rgba(186, 26, 26, 0.5); }

  .badge-mini {
    font-size: 0.65rem;
    font-weight: 850;
    padding: 0.1rem 0.5rem;
    border-radius: 4px;
    color: white;
    text-transform: uppercase;
    display: inline-flex;
    align-items: center;
  }
  .badge-mini.scoped {
    padding: 0;
    overflow: hidden;
  }
  .badge-mini.scoped .key {
    background: rgba(0, 0, 0, 0.2);
    padding: 0.1rem 0.5rem;
  }
  .badge-mini.scoped .val {
    padding: 0.1rem 0.5rem;
  }
</style>
