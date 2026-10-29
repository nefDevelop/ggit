<script lang="ts">
  import type { RepoContext } from "$lib/types";
  import { apiFetch } from "$lib/api";
  import { onMount, getContext } from "svelte";
  import { page } from "$app/state";
  import { goto } from "$app/navigation";

  import Diff from "lucide-svelte/icons/diff";
  import History from "lucide-svelte/icons/history";

  import Button from "$lib/components/ui/Button.svelte";
  import DiffViewer from "$lib/components/DiffViewer.svelte";
  import CommitRow from "$lib/components/CommitRow.svelte";
  import PRComparisonBar from "$lib/components/PRComparisonBar.svelte";
  import PRCreationForm from "$lib/components/PRCreationForm.svelte";

  const repo = getContext<RepoContext>('repo');
  const defaultBranch = $derived(repo?.data?.default_branch || "main");

  const username = $derived(page.params.username);
  const reponame = $derived(page.params.reponame);

  let baseBranch = $state("");
  let headBranch = $state("");
  let headRepoId = $state(0);
  let headRepoOwner = $state("");
  let headRepoName = $state("");

  let metadata = $state<any>(null);
  let comparison = $state<any>(null);
  let loading = $state(true);
  let loadingComp = $state(false);

  let title = $state("");
  let content = $state("");
  let submitting = $state(false);
  let activeTab = $state("diff");

  async function initComparison() {
    loading = true;
    try {
      const res = await apiFetch(`/api/v2/${username}/${reponame}/compare/${defaultBranch}...${defaultBranch}`);
      if (res.ok) {
        metadata = await res.json();
        baseBranch = metadata.base_branch;
        headBranch = metadata.head_branch;
        headRepoId = metadata.head_repo.id;
        headRepoOwner = metadata.head_repo.owner.username;
        headRepoName = metadata.head_repo.name;
        if (metadata.is_same_repo && baseBranch === headBranch) {
          const other = metadata.head_branches.find((b: string) => b !== baseBranch);
          if (other) headBranch = other;
        }
      }
    } catch (e) {
      console.error(e);
    } finally {
      loading = false;
    }
  }

  async function fetchComparison() {
    if (!baseBranch || !headBranch) return;
    loadingComp = true;
    try {
      const headRef = headRepoOwner === username ? headBranch : `${headRepoOwner}:${headBranch}`;
      const res = await apiFetch(`/api/v2/${username}/${reponame}/compare/${baseBranch}...${headRef}`);
      if (res.ok) {
        comparison = await res.json();
        metadata.head_branches = comparison.head_branches;
        if (!title && comparison.num_commits > 0) {
          title = comparison.commits[0].summary;
        }
      } else {
        comparison = null;
      }
    } catch (e) {
      console.error(e);
      comparison = null;
    } finally {
      loadingComp = false;
    }
  }

  async function createPR() {
    if (submitting || !title) return;
    submitting = true;
    try {
      const res = await apiFetch(`/api/v2/${username}/${reponame}/pulls`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ title, content, base_branch: baseBranch, head_branch: headBranch, head_repo_id: headRepoId }),
      });
      if (res.ok) {
        const pr = await res.json();
        goto(`/${username}/${reponame}/pulls/${pr.index}`);
      }
    } catch (e) {
      console.error(e);
    } finally {
      submitting = false;
    }
  }

  $effect(() => {
    if (baseBranch && headBranch && headRepoId) {
      if (baseBranch === headBranch && headRepoId === metadata?.base_repo?.id) {
        comparison = null;
        return;
      }
      fetchComparison();
    }
  });

  onMount(initComparison);
</script>

<div class="pr-creator">
  {#if loading}
    <div class="skeleton-header tonal"></div>
    <div class="workspace-grid skeleton-pulse">
      <div class=""></div>
      <div class=" tonal"></div>
    </div>
  {:else if metadata}
    <PRComparisonBar {username} {reponame} {metadata} bind:baseBranch bind:headBranch bind:headRepoId {comparison} {loadingComp} />

    {#if comparison}
      <div class="workspace-grid">
        <PRCreationForm bind:title bind:content {submitting} onCreate={createPR} />

        <main class="review-panel tonal shadow-md">
          <nav class="panel-tabs">
            <button class:active={activeTab === "diff"} onclick={() => (activeTab = "diff")}>
              <Diff size={20} />
              Changes
              <span class="badge">{comparison.diff?.Files?.length || 0}</span>
            </button>
            <button class:active={activeTab === "commits"} onclick={() => (activeTab = "commits")}>
              <History size={20} />
              Commits
              <span class="badge">{comparison.num_commits}</span>
            </button>
          </nav>

          <div class="panel-content">
            {#if loadingComp}
              <div class="loading-overlay">
                <div class="spinner"></div>
                <p>Updating comparison...</p>
              </div>
            {/if}
            {#if activeTab === "diff"}
              <div class="diff-container">
                <DiffViewer diff={comparison.diff} />
              </div>
            {:else}
              <div class="commits-container">
                {#each comparison.commits as commit}
                  <CommitRow {commit} minimalist={true} />
                {/each}
              </div>
            {/if}
          </div>
        </main>
      </div>
    {/if}
  {/if}
</div>

<style>
  .pr-creator { display: flex; flex-direction: column; gap: 2rem; }
  .workspace-grid {
    display: grid;
    grid-template-columns: 450px 1fr;
    gap: 2rem;
    align-items: start;
    animation: slideUp 0.4s ease-out;
  }
  @media (max-width: 1200px) {
    .workspace-grid { grid-template-columns: 1fr; }
  }
  @keyframes slideUp {
    from { opacity: 0; transform: translateY(20px); }
    to { opacity: 1; transform: translateY(0); }
  }
  .review-panel {
    display: flex;
    flex-direction: column;
    min-height: 800px;
    padding: 0;
    overflow: hidden;
  }
  .panel-tabs {
    display: flex;
    background: var(--sunken-2);
    border-bottom: 1px solid var(--md-sys-color-outline-variant);
  }
  .panel-tabs button {
    flex: 1;
    padding: 1.25rem;
    border: none;
    background: transparent;
    color: var(--text-muted);
    cursor: pointer;
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 0.75rem;
    font-weight: 700;
    transition: all 0.2s;
  }
  .panel-tabs button.active {
    color: var(--primary-color);
    background: rgba(var(--primary-rgb),0.05);
    position: relative;
  }
  .panel-tabs button.active::after {
    content: ""; position: absolute;
    bottom: 0; left: 0; right: 0;
    height: 3px; background: var(--primary-color);
  }
  .badge {
    background: rgba(255,255,255,0.05);
    padding: 0.1rem 0.6rem; border-radius: 2rem;
    font-size: 0.8rem; color: var(--text-color);
  }
  .active .badge { background: var(--primary-color); color: var(--on-primary); }
  .panel-content { flex: 1; position: relative; }
  .loading-overlay {
    position: absolute; inset: 0;
    background: rgba(var(--md3-surface-rgb),0.7);
    display: flex; flex-direction: column;
    align-items: center; justify-content: center;
    gap: 1rem; z-index: 10; backdrop-filter: blur(4px);
  }
  .loading-overlay p { font-weight: 700; color: var(--primary-color); }
  .diff-container, .commits-container { padding: 1rem; }
  .skeleton-header { height: 100px; }
  .skeleton-pulse { opacity: 0.5; pointer-events: none; }
</style>
