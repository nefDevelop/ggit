<script lang="ts">
  import { apiFetch } from "$lib/api";
  import { onMount, getContext } from "svelte";
  import { page } from "$app/state";
  import { goto } from "$app/navigation";
  import DiffViewer from "$lib/components/DiffViewer.svelte";
  import CommitRow from "$lib/components/CommitRow.svelte";
  import BranchSelector from "$lib/components/BranchSelector.svelte";

  const repo: any = getContext("repo");
  const defaultBranch = $derived(repo?.data?.default_branch || "main");

  const username = $derived(page.params.username);
  const reponame = $derived(page.params.reponame);

  // State
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
  let activeTab = $state("diff"); // 'diff' or 'commits'

  // Fetch initial comparison data
  async function initComparison() {
    loading = true;
    try {
      // We start by comparing the default branch with itself or a suggested branch
      // The backend will provide available branches and repositories (forks)
      const res = await apiFetch(
        `/api/v2/${username}/${reponame}/compare/${defaultBranch}...${defaultBranch}`,
      );
      if (res.ok) {
        metadata = await res.json();
        baseBranch = metadata.base_branch;
        headBranch = metadata.head_branch;
        headRepoId = metadata.head_repo.id;
        headRepoOwner = metadata.head_repo.owner.username;
        headRepoName = metadata.head_repo.name;

        // If there's only one repo and we are comparing same branch, try to find another branch
        if (metadata.is_same_repo && baseBranch === headBranch) {
          const other = metadata.head_branches.find(
            (b: string) => b !== baseBranch,
          );
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
      // Format: base...headUser:headBranch
      const headRef =
        headRepoOwner === username
          ? headBranch
          : `${headRepoOwner}:${headBranch}`;
      const res = await apiFetch(
        `/api/v2/${username}/${reponame}/compare/${baseBranch}...${headRef}`,
      );
      if (res.ok) {
        comparison = await res.json();
        // Update local metadata if branches/repos changed in the response
        metadata.head_branches = comparison.head_branches;

        // Auto-fill title if empty
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
    if (!title) return;
    submitting = true;
    try {
      const res = await apiFetch(`/api/v2/${username}/${reponame}/pulls`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          title,
          content,
          base_branch: baseBranch,
          head_branch: headBranch,
          head_repo_id: headRepoId,
        }),
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

  function handleHeadRepoChange(e: Event) {
    const id = parseInt((e.target as HTMLSelectElement).value);
    const repo = metadata.all_head_repos.find((r: any) => r.id === id);
    if (repo) {
      headRepoId = repo.id;
      headRepoOwner = repo.owner.username;
      headRepoName = repo.name;
      // When repo changes, we should probably reset head branch to default
      headBranch = repo.default_branch || "master";
    }
  }

  $effect(() => {
    if (baseBranch && headBranch && headRepoId) {
      if (baseBranch === headBranch && headRepoId === metadata.base_repo.id) {
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
    <header class="comparison-bar tonal shadow-lg">
      <div class="selector-group">
        <div class="repo-info">
          <span class="label">Base Repository</span>
          <div class="repo-name">
            <span class="material-symbols-rounded">folder</span>
            {username}/{reponame}
          </div>
        </div>
        <BranchSelector
          label="Base Branch"
          branches={metadata.base_branches.map((b: string) => ({ name: b }))}
          bind:selected={baseBranch}
        />
      </div>

      <div class="comparison-arrow">
        <span class="material-symbols-rounded">arrow_back</span>
      </div>

      <div class="selector-group">
        <div class="repo-selector">
          <span class="label">Head Repository</span>
          <select class="" onchange={handleHeadRepoChange} value={headRepoId}>
            {#each metadata.all_head_repos as repo}
              <option value={repo.id}>{repo.owner.username}/{repo.name}</option>
            {/each}
          </select>
        </div>
        <BranchSelector
          label="Compare Branch"
          branches={metadata.head_branches.map((b: string) => ({ name: b }))}
          bind:selected={headBranch}
          icon="call_split"
        />
      </div>

      <div class="comparison-status">
        {#if baseBranch === headBranch && metadata.head_repo.id === metadata.base_repo.id}
          <div class="status-chip neutral">
            <span class="material-symbols-rounded">info</span>
            Choose different branches
          </div>
        {:else if comparison}
          <div class="status-chip success">
            <span class="material-symbols-rounded">check_circle</span>
            Able to merge
          </div>
        {:else if !loadingComp}
          <div class="status-chip error">
            <span class="material-symbols-rounded">error</span>
            Conflict detected
          </div>
        {/if}
      </div>
    </header>

    {#if comparison}
      <div class="workspace-grid">
        <section class="creation-details shadow-xl">
          <div class="form-header">
            <h3>New Pull Request</h3>
            <p class="text-muted">
              Fill in the details to submit your changes for review.
            </p>
          </div>

          <div class="form-fields">
            <div class="field">
              <label for="pr-title">Title</label>
              <input
                id="pr-title"
                type="text"
                bind:value={title}
                placeholder="A concise summary of changes"
                required
              />
            </div>
            <div class="field">
              <label for="pr-body">Description</label>
              <textarea
                id="pr-body"
                bind:value={content}
                placeholder="Explain what this PR does and why it's needed..."
              ></textarea>
            </div>
          </div>

          <div class="form-actions">
            <button
              class="btn btn-primary lg"
              onclick={createPR}
              disabled={!title || submitting}
            >
              {#if submitting}
                <span class="spinner-sm"></span>
                Creating...
              {:else}
                <span class="material-symbols-rounded">merge</span>
                Create Pull Request
              {/if}
            </button>
          </div>
        </section>

        <main class="review-panel tonal shadow-md">
          <nav class="panel-tabs">
            <button
              class:active={activeTab === "diff"}
              onclick={() => (activeTab = "diff")}
            >
              <span class="material-symbols-rounded">difference</span>
              Changes
              <span class="badge">{comparison.diff?.Files?.length || 0}</span>
            </button>
            <button
              class:active={activeTab === "commits"}
              onclick={() => (activeTab = "commits")}
            >
              <span class="material-symbols-rounded">history</span>
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
  .pr-creator {
    display: flex;
    flex-direction: column;
    gap: 2rem;
  }

  .comparison-bar {
    padding: 1.5rem 2.5rem;
    display: flex;
    align-items: center;
    gap: 2rem;
    background: var(--md3-surface-1);
    border-radius: 1.5rem;
    flex-wrap: wrap;
  }
  .selector-group {
    display: flex;
    align-items: flex-end;
    gap: 1rem;
    flex: 1;
    min-width: 300px;
  }
  .repo-info,
  .repo-selector {
    display: flex;
    flex-direction: column;
    gap: 0.25rem;
    flex: 1;
    min-width: 0;
  }
  .label {
    font-size: 0.7rem;
    font-weight: 800;
    text-transform: uppercase;
    color: var(--text-muted);
    opacity: 0.7;
    padding-left: 0.5rem;
  }

  .repo-name {
    padding: 0.6rem 1rem;
    background: rgba(255, 255, 255, 0.02);
    border: 1px solid var(--md-sys-color-outline-variant);
    border-radius: 0.75rem;
    font-weight: 700;
    display: flex;
    align-items: center;
    gap: 0.75rem;
    color: var(--text-muted);
    overflow: hidden;
    white-space: nowrap;
    text-overflow: ellipsis;
  }
  .repo-name span:last-child {
    overflow: hidden;
    text-overflow: ellipsis;
  }
  .repo-selector select {
    background: rgba(255, 255, 255, 0.02);
    border: 1px solid var(--md-sys-color-outline-variant);
    color: #fff;
    padding: 0.6rem 1rem;
    border-radius: 0.75rem;
    font-weight: 600;
    width: 100%;
    outline: none;
  }

  .comparison-arrow {
    display: flex;
    align-items: center;
    margin-top: 1rem;
    color: var(--primary-color);
    opacity: 0.5;
  }

  .comparison-status {
    width: 220px;
    display: flex;
    justify-content: flex-end;
  }
  .status-chip {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    padding: 0.5rem 1rem;
    border-radius: 2rem;
    font-size: 0.85rem;
    font-weight: 700;
  }
  .status-chip.success {
    background: rgba(34, 197, 94, 0.1);
    color: #22c55e;
  }
  .status-chip.neutral {
    background: rgba(255, 255, 255, 0.05);
    color: var(--text-muted);
  }
  .status-chip.error {
    background: rgba(239, 68, 68, 0.1);
    color: #ef4444;
  }

  .workspace-grid {
    display: grid;
    grid-template-columns: 450px 1fr;
    gap: 2rem;
    align-items: start;
    animation: slideUp 0.4s ease-out;
  }

  @media (max-width: 1200px) {
    .workspace-grid {
      grid-template-columns: 1fr;
    }
    .creation-details {
      padding: 1.5rem;
    }
  }

  @media (max-width: 800px) {
    .comparison-bar {
      padding: 1rem;
      gap: 1rem;
    }
    .selector-group {
      min-width: 100%;
    }
    .comparison-arrow {
      width: 100%;
      justify-content: center;
      margin-top: 0;
      transform: rotate(-90deg);
    }
    .comparison-status {
      width: 100%;
      justify-content: center;
    }
  }
  @keyframes slideUp {
    from {
      opacity: 0;
      transform: translateY(20px);
    }
    to {
      opacity: 1;
      transform: translateY(0);
    }
  }

  .creation-details {
    padding: 2.5rem;
    display: flex;
    flex-direction: column;
    gap: 2rem;
  }
  .form-header h3 {
    margin: 0;
    font-size: 1.5rem;
    font-weight: 800;
  }
  .form-header p {
    margin: 0.25rem 0 0;
    font-size: 0.95rem;
  }

  .form-fields {
    display: flex;
    flex-direction: column;
    gap: 1.5rem;
  }
  .field {
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
  }
  .field label {
    font-size: 0.8rem;
    font-weight: 700;
    color: var(--text-muted);
  }
  .field input {
    background: rgba(0, 0, 0, 0.2);
    border: 1px solid var(--md-sys-color-outline-variant);
    color: #fff;
    padding: 1rem;
    border-radius: 1rem;
    font-size: 1.1rem;
    width: 100%;
    outline: none;
  }
  .field input:focus {
    border-color: var(--primary-color);
    background: rgba(var(--primary-rgb), 0.05);
  }
  .field textarea {
    background: rgba(0, 0, 0, 0.2);
    border: 1px solid var(--md-sys-color-outline-variant);
    color: #fff;
    padding: 1rem;
    border-radius: 1rem;
    font-size: 1rem;
    height: 350px;
    resize: none;
    width: 100%;
    outline: none;
    font-family: inherit;
  }
  .field textarea:focus {
    border-color: var(--primary-color);
    background: rgba(var(--primary-rgb), 0.05);
  }

  .form-actions {
    margin-top: 1rem;
  }
  .btn.lg {
    height: 64px;
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 0.75rem;
    font-size: 1.1rem;
    font-weight: 800;
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
    background: rgba(0, 0, 0, 0.2);
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
    background: rgba(var(--primary-rgb), 0.05);
    position: relative;
  }
  .panel-tabs button.active::after {
    content: "";
    position: absolute;
    bottom: 0;
    left: 0;
    right: 0;
    height: 3px;
    background: var(--primary-color);
  }
  .badge {
    background: rgba(255, 255, 255, 0.05);
    padding: 0.1rem 0.6rem;
    border-radius: 2rem;
    font-size: 0.8rem;
    color: var(--text-color);
  }
  .active .badge {
    background: var(--primary-color);
    color: var(--on-primary);
  }

  .panel-content {
    flex: 1;
    position: relative;
  }
  .loading-overlay {
    position: absolute;
    inset: 0;
    background: rgba(var(--md3-surface-rgb), 0.7);
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    gap: 1rem;
    z-index: 10;
    backdrop-filter: blur(4px);
  }
  .loading-overlay p {
    font-weight: 700;
    color: var(--primary-color);
  }

  .diff-container,
  .commits-container {
    padding: 1rem;
  }

  .skeleton-header {
    height: 100px;
  }
  .skeleton-pulse {
    opacity: 0.5;
    pointer-events: none;
  }
  .spinner-sm {
    width: 18px;
    height: 18px;
    border: 2px solid rgba(255, 255, 255, 0.2);
    border-top-color: #fff;
    border-radius: 50%;
    animation: spin 0.8s linear infinite;
  }
  @keyframes spin {
    to {
      transform: rotate(360deg);
    }
  }
</style>
