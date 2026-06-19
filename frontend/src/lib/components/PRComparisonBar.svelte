<script lang="ts">
  import BranchSelector from "$lib/components/BranchSelector.svelte";
  import Folder from "lucide-svelte/icons/folder";
  import ArrowLeft from "lucide-svelte/icons/arrow-left";
  import Info from "lucide-svelte/icons/info";
  import CheckCircle2 from "lucide-svelte/icons/check-circle-2";
  import AlertCircle from "lucide-svelte/icons/alert-circle";
  import GitFork from "lucide-svelte/icons/git-fork";

  let { username = "", reponame = "", metadata = null, baseBranch = $bindable(), headBranch = $bindable(), headRepoId = $bindable(), comparison = null, loadingComp = false }: {
    username: string;
    reponame: string;
    metadata: any;
    baseBranch: string;
    headBranch: string;
    headRepoId: number;
    comparison: any;
    loadingComp: boolean;
  } = $props();

  function handleHeadRepoChange(e: Event) {
    const id = parseInt((e.target as HTMLSelectElement).value);
    const repo = metadata?.all_head_repos?.find((r: any) => r.id === id);
    if (repo) {
      headRepoId = repo.id;
      headBranch = repo.default_branch || "master";
    }
  }
</script>

<header class="comparison-bar tonal shadow-lg">
  <div class="selector-group">
    <div class="repo-info">
      <span class="label">Base Repository</span>
      <div class="repo-name">
        <Folder size={18} />
        {username}/{reponame}
      </div>
    </div>
    <BranchSelector
      label="Base Branch"
      branches={(metadata?.base_branches || []).map((b: string) => ({ name: b }))}
      bind:selected={baseBranch}
    />
  </div>

  <div class="comparison-arrow">
    <ArrowLeft size={24} />
  </div>

  <div class="selector-group">
    <div class="repo-selector">
      <span class="label">Head Repository</span>
      <select onchange={handleHeadRepoChange} value={headRepoId}>
        {#each metadata?.all_head_repos || [] as repo}
          <option value={repo.id}>{repo.owner.username}/{repo.name}</option>
        {/each}
      </select>
    </div>
    <BranchSelector
      label="Compare Branch"
      branches={(metadata?.head_branches || []).map((b: string) => ({ name: b }))}
      bind:selected={headBranch}
      icon={GitFork}
    />
  </div>

  <div class="comparison-status">
    {#if baseBranch === headBranch && metadata?.head_repo?.id === metadata?.base_repo?.id}
      <div class="status-chip neutral">
        <Info size={18} />
        Choose different branches
      </div>
    {:else if comparison}
      <div class="status-chip success">
        <CheckCircle2 size={18} />
        Able to merge
      </div>
    {:else if !loadingComp}
      <div class="status-chip error">
        <AlertCircle size={18} />
        Conflict detected
      </div>
    {/if}
  </div>
</header>

<style>
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
  .repo-info, .repo-selector {
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
    background: rgba(255,255,255,0.02);
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
  .repo-selector select {
    background: rgba(255,255,255,0.02);
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
  .status-chip.success { background: rgba(34,197,94,0.1); color: #22c55e; }
  .status-chip.neutral { background: rgba(255,255,255,0.05); color: var(--text-muted); }
  .status-chip.error { background: rgba(239,68,68,0.1); color: #ef4444; }
  @media (max-width: 800px) {
    .comparison-bar { padding: 1rem; gap: 1rem; }
    .selector-group { min-width: 100%; }
    .comparison-arrow { width: 100%; justify-content: center; margin-top: 0; transform: rotate(-90deg); }
    .comparison-status { width: 100%; justify-content: center; }
  }
</style>
