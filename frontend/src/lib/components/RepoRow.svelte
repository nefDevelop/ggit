<script lang="ts">
  import { 
    Package, Lock, GitFork, History, GitBranch, AlertCircle, 
    GitPullRequest, Check, Copy 
  } from "lucide-svelte";
  import { getLanguageColor } from '$lib/utils/languages';
  import { goto } from '$app/navigation';

  let { repo, isOlder = false } = $props();

  function timeAgo(dateStr: string) {
    const date = new Date(dateStr);
    const seconds = Math.floor((Date.now() - date.getTime()) / 1000);
    if (seconds < 60) return 'just now';
    const minutes = Math.floor(seconds / 60);
    if (minutes < 60) return `${minutes}m ago`;
    const hours = Math.floor(minutes / 60);
    if (hours < 24) return `${hours}h ago`;
    const days = Math.floor(hours / 24);
    if (days < 30) return `${days}d ago`;
    return date.toLocaleDateString();
  }

  // Calculate language percentages for the mini bar
  const langEntries = $derived(() => {
    if (!repo?.languages) return [];
    const total = Object.values(repo.languages as Record<string, number>).reduce((a, b) => a + b, 0);
    return Object.entries(repo.languages as Record<string, number>)
      .map(([name, bytes]) => ({
        name,
        percent: ((bytes / total) * 100).toFixed(1),
        color: getLanguageColor(name)
      }))
      .sort((a, b) => parseFloat(b.percent) - parseFloat(a.percent));
  });
  
  const isNew = $derived(() => {
    if (!repo?.updated_at) return false;
    const date = new Date(repo.updated_at);
    return (Date.now() - date.getTime()) < (24 * 60 * 60 * 1000);
  });

  const recentThresh = 24 * 60 * 60 * 1000; // 24 hours
  const isRecentCommit = $derived(() => repo.last_commit_at && (Date.now() - new Date(repo.last_commit_at).getTime()) < recentThresh);
  const isRecentIssue = $derived(() => repo.last_issue_at && (Date.now() - new Date(repo.last_issue_at).getTime()) < recentThresh);
  const isRecentPull = $derived(() => repo.last_pull_at && (Date.now() - new Date(repo.last_pull_at).getTime()) < recentThresh);
  const isRecentBranch = $derived(() => repo.last_branch_at && (Date.now() - new Date(repo.last_branch_at).getTime()) < recentThresh);

  import { toast } from '$lib/toast.svelte';
  import Button from './ui/Button.svelte';
  import Avatar from './Avatar.svelte';

  let copied = $state(false);
  function copyToClipboard(text: string) {
    navigator.clipboard.writeText(text);
    copied = true;
    toast.success("Enlace de clonación copiado");
    setTimeout(() => copied = false, 2000);
  }
</script>

<div class="repo-row {isOlder ? 'older-row' : ''} {isNew() ? 'new-update' : ''}" onclick={() => goto(`/${repo.full_name}`)} role="link" tabindex="0" onkeydown={(e) => e.key === 'Enter' && goto(`/${repo.full_name}`)}>
  
  <div class="repo-avatar-wrap">
    <Avatar 
      src={repo.avatar_url} 
      username={repo.name} 
      seed={repo.id?.toString()} 
      size={48} 
      type="repo" 
    />
  </div>

  <div class="main-info">
    <div class="title-line">
      <a href="/{repo.full_name}" class="repo-name" data-sveltekit-preload-data="hover" onclick={(e) => e.stopPropagation()}>
        {repo.name}
      </a>
      {#if repo.empty}
        <span class="badge-empty">Empty</span>
      {/if}
    </div>

    <!-- REAL MINI LANGUAGE BAR ABOVE DESCRIPTION -->
    {#if langEntries().length > 0}
      <div class="mini-lang-bar">
        {#each langEntries() as lang}
          <div class="lang-segment" style="width: {lang.percent}%; background: {lang.color}" title="{lang.name} {lang.percent}%"></div>
        {/each}
      </div>
    {/if}

    <div class="desc-line">
      {#if repo.private}
        <Lock size={14} class="status-icon" />
      {/if}
      {#if repo.fork}
        <GitFork size={14} class="status-icon" />
      {/if}
      <p class="repo-desc">
        {repo.description || 'No description provided.'}
      </p>
    </div>
  </div>

  <div class="tech-info">
    {#if repo.language}
      <div class="lang">
        <span class="lang-dot" style="background-color: {getLanguageColor(repo.language)}"></span>
        {repo.language}
      </div>
    {:else if repo.empty}
      <div class="setup-hint">Ready for code</div>
    {/if}
  </div>

  <div class="stats-info">
    <div class="stat-item {repo.num_commits > 0 ? 'active' : 'dim'}" title="Commits">
      <History size={16} class={isRecentCommit() ? 'blink-active' : ''} />
      {repo.num_commits || 0}
    </div>
    <div class="stat-item {repo.num_branches > 0 ? 'active' : 'dim'}" title="Branches">
      <GitBranch size={16} class={isRecentBranch() ? 'blink-active' : ''} />
      {repo.num_branches || 0}
    </div>
    <div class="stat-item {repo.open_issues_count > 0 ? 'active' : 'dim'}" title="Issues">
      <AlertCircle size={16} class={isRecentIssue() ? 'blink-active' : ''} />
      {repo.open_issues_count || 0}
    </div>
    <div class="stat-item {repo.open_pulls_count > 0 ? 'active' : 'dim'}" title="Pull Requests">
      <GitPullRequest size={16} class={isRecentPull() ? 'blink-active' : ''} />
      {repo.open_pulls_count || 0}
    </div>
  </div>

  <div class="actions-info">
    <Button 
      variant="ghost" 
      size="sm" 
      icon={copied ? Check : Copy} 
      class={copied ? 'success-text' : ''}
      onclick={(e) => { e.stopPropagation(); copyToClipboard(repo.clone_url); }} 
      aria-label="Copy Clone URL" 
    />
    <span class="time">{timeAgo(repo.updated_at)}</span>
  </div>
</div>

<style>
  .repo-row {
    display: grid;
    grid-template-columns: 56px 1fr 80px 130px 90px;
    align-items: center;
    padding: 1.15rem 1.5rem;
    background: transparent;
    transition: all 0.25s cubic-bezier(0.34, 1.56, 0.64, 1);
    gap: 1.25rem;
    position: relative;
    border: none;
    cursor: pointer;
    min-width: 0;
    border-bottom: 1px solid var(--md-sys-color-outline-variant);
    max-width: 850px;
    margin: 0 auto;
    width: 100%;
    box-sizing: border-box;
  }

  .repo-row:last-child { border-bottom: none; }

  .repo-row:hover {
    background: rgba(var(--primary-rgb), 0.05);
    z-index: 10;
    box-shadow: inset 0 0 0 1px rgba(var(--primary-rgb), 0.1);
  }

  .repo-avatar-wrap {
    width: 48px;
    height: 48px;
    display: flex;
    align-items: center;
    justify-content: center;
    flex-shrink: 0;
  }

  .new-update {
    background: linear-gradient(90deg, rgba(var(--primary-rgb), 0.08), transparent) !important;
  }

  /* SMOOTHER BLINK ANIMATIONS FOR COUNTERS */
  @keyframes activity-blink {
    0%, 100% { opacity: 1; transform: scale(1); filter: drop-shadow(0 0 0px var(--tertiary)); }
    50% { opacity: 0.8; transform: scale(1.05); filter: drop-shadow(0 0 4px var(--tertiary)); color: var(--tertiary); }
  }

  :global(.blink-active) {
    animation: activity-blink 3.5s infinite ease-in-out;
  }

  .older-row { padding: 0.75rem 1.25rem; opacity: 0.8; }
  .main-info { display: flex; flex-direction: column; min-width: 0; gap: 0.35rem; }

  /* MINI LANG BAR */
  .mini-lang-bar {
    display: flex;
    height: 4px;
    width: 120px;
    border-radius: 2px;
    overflow: hidden;
    background: var(--md-sys-color-surface-variant);
    margin: 0.1rem 0;
  }
  .lang-segment { height: 100%; }

  .desc-line { display: flex; align-items: flex-start; gap: 0.5rem; min-width: 0; }
  :global(.status-icon) { color: var(--text-muted); opacity: 0.6; flex-shrink: 0; margin-top: 2px; }
  .repo-name { 
    font-weight: 700; 
    color: var(--primary-color); 
    text-decoration: none; 
    font-size: 1.1rem; 
    display: inline-flex; 
    align-items: center; 
    gap: 0.5rem; 
    max-width: 100%;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }
  .repo-desc {
    font-size: 0.85rem;
    color: var(--text-muted);
    margin: 0;
    display: -webkit-box;
    -webkit-line-clamp: 2;
    line-clamp: 2;
    -webkit-box-orient: vertical;
    overflow: hidden;
    text-overflow: ellipsis;
  }

  .badge-empty { font-size: 0.65rem; background: var(--md3-surface-2); padding: 0.1rem 0.5rem; border-radius: 2rem; font-weight: 800; letter-spacing: 0.05em; }
  .lang { display: flex; align-items: center; gap: 0.5rem; font-size: 0.85rem; font-weight: 600; }
  .lang-dot { width: 10px; height: 10px; border-radius: 50%; }

  .stats-info { 
    display: flex; 
    align-items: center;
    gap: 1.25rem; 
    flex-wrap: nowrap;
  }
  
  .stat-item { 
    display: flex; 
    align-items: center; 
    gap: 0.35rem; 
    font-size: 0.85rem; 
    white-space: nowrap;
  }
  

  .stat-item.dim { opacity: 0.2; }
  .stat-item.active { color: var(--primary-color); font-weight: 700; }

  .actions-info { display: flex; align-items: center; justify-content: flex-end; gap: 1rem; }
  :global(.success-icon) { color: var(--md-sys-color-primary) !important; }
  .time { font-size: 0.8rem; color: var(--text-muted); font-variant-numeric: tabular-nums; flex-shrink: 0; }

  @media (max-width: 1200px) {
    .repo-row { 
      grid-template-columns: 56px 1fr 80px 130px 90px; 
    }
    .stats-info { gap: 0.75rem; }
  }

  @media (max-width: 1000px) {
    .repo-row { 
      grid-template-columns: 56px 1fr 80px 90px; 
    }
    .stats-info { display: none; }
  }

  @media (max-width: 800px) {
    .repo-row { grid-template-columns: 56px 1fr 80px; }
    .tech-info { display: none; }
  }

  @media (max-width: 480px) {
    .repo-row { padding: 0.75rem 0.5rem; gap: 0.75rem; grid-template-columns: 48px 1fr 40px; }
    .repo-name { font-size: 1rem; }
    .repo-desc { font-size: 0.8rem; }
    .actions-info .time { display: none; }
  }
</style>
