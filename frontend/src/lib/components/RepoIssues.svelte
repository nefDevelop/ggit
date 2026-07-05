<script lang="ts">
  import type { RepoContext } from "$lib/types";
  import { apiFetch } from "$lib/api";
  import { onMount, getContext } from "svelte";
  import { formatDate } from "$lib/utils/date";
  import PageHeader from "./ui/PageHeader.svelte";
  import IssueBoard from "./IssueBoard.svelte";
  import Button from "./ui/Button.svelte";
  

import List from "lucide-svelte/icons/list";
import LayoutDashboard from "lucide-svelte/icons/layout-dashboard";
import Plus from "lucide-svelte/icons/plus";
import AlertCircle from "lucide-svelte/icons/alert-circle";
import CheckCircle2 from "lucide-svelte/icons/check-circle-2";
import MessageSquare from "lucide-svelte/icons/message-square";



  const repo = getContext<RepoContext>('repo');
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
  <PageHeader title="" icon="" compact noBorder>
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
      {#if issues.length > 0 || stateFilter === 'closed'}
        <a
          href="/{repo.data.owner.username}/{repo.data.name}/issues/new"
          class="btn-action-primary"
        >
          <Plus size={18} />
          New Issue
        </a>
      {/if}
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
      <div class="empty-icon-ring">
        <AlertCircle size={40} />
      </div>
      <h2>No issues found</h2>
      <p>There are no {stateFilter} issues in this repository.</p>
      {#if stateFilter === 'open'}
        <a href="/{repo.data.owner.username}/{repo.data.name}/issues/new" class="btn-action-primary" style="margin-top: 1rem; text-decoration: none;">
          <Plus size={18} />
          Create your first issue
        </a>
      {/if}
    </div>
  {:else}
    {#if view === 'list'}
      <div class="flat-list-container">
        {#each issues as issue}
          <div class="list-item-row">
            <div class="list-item-icon" style="color: {getStatusColor(issue)}">
              {#if issue.is_closed}
                <CheckCircle2 size={18} />
              {:else}
                <AlertCircle size={18} />
              {/if}
            </div>
            <div class="item-main-content">
              <div class="item-title-bar">
                <a
                  href="/{repo.data.owner.username}/{repo.data.name}/issues/{issue.number}"
                  class="item-title-link"
                >
                  {issue.title}
                </a>
                {#each issue.labels || [] as label}
                  <span
                    class="badge-mini"
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
              <div class="item-meta-data">
                #{issue.number} opened {formatDate(issue.created_at)} by <strong>{issue.user?.username}</strong>
              </div>
            </div>
            {#if issue.comments > 0}
              <div class="item-actions-right">
                <div class="issue-comments-count">
                  <MessageSquare size={16} />
                  {issue.comments}
                </div>
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
    padding: 0;
  }

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
    padding: 6rem 2rem;
    gap: 1.25rem;
    color: var(--text-muted);
    text-align: center;
  }

  .empty-state h2 {
    margin: 0;
    font-size: 1.75rem;
    font-weight: 950;
    color: var(--text-color);
  }

  .empty-state p {
    max-width: 400px;
    margin: 0;
    font-size: 1.1rem;
    line-height: 1.5;
  }

  .empty-icon-ring {
    width: 80px;
    height: 80px;
    border: 2px dashed var(--md-sys-color-outline-variant);
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;
    margin-bottom: 1rem;
    color: var(--primary-color);
    background: rgba(var(--primary-rgb), 0.05);
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
