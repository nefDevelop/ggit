<script lang="ts">
  import { apiFetch } from "$lib/api";
  import { onMount, getContext } from "svelte";
  import PageHeader from "./ui/PageHeader.svelte";
  import IssueBoard from "./IssueBoard.svelte";
  import Button from "./ui/Button.svelte";
  import { List, LayoutDashboard, Plus, AlertCircle, CheckCircle2, MessageSquare } from "lucide-svelte";

  const repo: any = getContext("repo");
  let issues = $state<any[]>([]);
  let loading = $state(true);
  let error = $state("");
  let stateFilter = $state("open"); // 'open' or 'closed'
  let view = $state<'list' | 'board'>('list');

  async function fetchIssues() {
    loading = true;
    error = "";
    try {
      const res = await apiFetch(
        `/api/v2/${repo.data.owner.username}/${repo.data.name}/issues?state=${stateFilter}`,
      );
      if (res.ok) {
        issues = await res.json();
      } else {
        error = "Failed to load issues.";
      }
    } catch (e) {
      error = "Error connecting to API.";
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
    return issue.is_closed ? "var(--error)" : "var(--success)";
  }
</script>

<div class="issues-view">
  <PageHeader title="" icon="" compact>
    {#snippet actions()}
      <div class="pill-filter-group-inline">
        <Button
          variant={stateFilter === "open" ? "primary" : "ghost"}
          size="sm"
          onclick={() => (stateFilter = "open")}
        >
          {repo.data?.open_issues_count || 0} Open
        </Button>
        <Button
          variant={stateFilter === "closed" ? "primary" : "ghost"}
          size="sm"
          onclick={() => (stateFilter = "closed")}
        >
          {repo.data?.closed_issues_count || 0} Closed
        </Button>
      </div>

      <div class="view-toggle">
        <Button 
          variant={view === 'list' ? 'tonal' : 'ghost'} 
          size="sm"
          icon={List}
          onclick={() => view = 'list'}
        />
        <Button 
          variant={view === 'board' ? 'tonal' : 'ghost'} 
          size="sm"
          icon={LayoutDashboard}
          onclick={() => view = 'board'}
        />
      </div>
      <a
        href="/{repo.data.owner.username}/{repo.data.name}/issues/new"
        class="btn-action-primary"
      >
        <Plus size={18} />
        New Issue
      </a>
    {/snippet}
  </PageHeader>


  {#if loading}
    <div class="loading-state">
      <div class="spinner"></div>
      <p>Loading issues...</p>
    </div>
  {:else if error}
    <div class="error-msg">{error}</div>
  {:else if issues.length === 0}
    <div class="empty-state">
      <AlertCircle size={48} />
      <h3>No issues found</h3>
      <p>There are no {stateFilter} issues in this repository.</p>
    </div>
  {:else}
    {#if view === 'list'}
      <div class="issues-list-modern">
        {#each issues as issue}
          <div class="issue-row">
            <div class="issue-icon" style="color: {getStatusColor(issue)}">
              {#if issue.is_closed}
                <CheckCircle2 size={18} />
              {:else}
                <AlertCircle size={18} />
              {/if}
            </div>
            <div class="issue-main">
              <div class="issue-title-row">
                <a
                  href="/{repo.data.owner.username}/{repo.data
                    .name}/issues/{issue.number}"
                  class="issue-title"
                >
                  {issue.title}
                </a>
                {#each issue.labels || [] as label}
                  <span
                    class="label-tag"
                    style="background: {label.color}"
                    class:scoped={label.name.includes("::")}
                  >
                    {#if label.name.includes("::")}
                      {@const [key, val] = label.name.split("::")}
                      <span class="key">{key}</span>
                      <span class="val">{val}</span>
                    {:else}
                      {label.name}
                    {/if}
                  </span>
                {/each}
              </div>
              <div class="issue-meta">
                #{issue.number} opened {new Date(
                  issue.created_at,
                ).toLocaleDateString()} by <strong>{issue.user?.username}</strong>
              </div>
            </div>
            {#if issue.comments > 0}
              <div class="issue-comments-count">
                <MessageSquare size={16} />
                {issue.comments}
              </div>
            {/if}
          </div>
        {/each}
      </div>
    {:else}
      <IssueBoard {issues} />
    {/if}
  {/if}
</div>

<style>
  .issues-view {
    display: flex;
    flex-direction: column;
    gap: 0;
    animation: fadeIn 0.3s ease-out;
    margin: -2rem; /* Negate parent padding for full-width header if needed, but let's just use gap 0 */
  }

  /* ... existing styles ... */
  .pill-filter-group-inline {
    display: flex;
    background: var(--surface-variant);
    padding: 4px;
    border-radius: 10px;
    margin-right: 1rem;
  }

  @keyframes fadeIn {
    from {
      opacity: 0;
      transform: translateY(10px);
    }
    to {
      opacity: 1;
      transform: translateY(0);
    }
  }

  .issue-row {
    display: flex;
    padding: 1.25rem 1.5rem;
    gap: 1.25rem;
    border-bottom: 1px solid var(--md-sys-color-outline-variant);
    transition: background 0.2s;
  }
  .issue-row:hover {
    background: var(--surface-variant);
  }

  .issue-icon {
    padding-top: 0.2rem;
  }

  .issue-main {
    flex: 1;
    display: flex;
    flex-direction: column;
    gap: 0.25rem;
    min-width: 0;
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
    overflow: hidden;
    text-overflow: ellipsis;
    display: block;
    max-width: 100%;
  }
  .issue-title:hover {
    color: var(--primary-color);
  }

  .label-tag {
    font-size: 0.65rem;
    font-weight: 800;
    padding: 0.15rem 0.5rem;
    border-radius: 4px;
    color: white;
    text-transform: uppercase;
    letter-spacing: 0.05em;
    display: inline-flex;
    align-items: center;
  }

  .label-tag.scoped {
    padding: 0;
    overflow: hidden;
    gap: 0;
  }
  .label-tag.scoped .key {
    background: rgba(0, 0, 0, 0.25);
    padding: 0.15rem 0.5rem;
    opacity: 0.9;
  }
  .label-tag.scoped .val {
    padding: 0.15rem 0.5rem;
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

  .loading-state,
  .empty-state {
    display: flex;
    flex-direction: column;
    align-items: center;
    padding: 4rem;
    gap: 1rem;
    color: var(--text-muted);
    text-align: center;
  }

  .spinner {
    width: 32px;
    height: 32px;
    border: 3px solid var(--md-sys-color-outline-variant);
    border-top-color: var(--primary-color);
    border-radius: 50%;
    animation: spin 1s infinite linear;
  }

  @keyframes spin {
    to {
      transform: rotate(360deg);
    }
  }

  .error-msg {
    color: var(--error);
    padding: 2rem;
    text-align: center;
    background: var(--error-container);
    border-radius: 1rem;
  }

  .view-toggle {
    display: flex;
    background: var(--surface-variant);
    padding: 4px;
    border-radius: 10px;
    margin-right: 0.5rem;
  }
</style>
