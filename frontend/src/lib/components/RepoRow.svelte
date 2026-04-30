<script lang="ts">
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

  let avatarError = $state(false);

  import { toast } from '$lib/toast.svelte';

  let copied = $state(false);
  function copyToClipboard(text: string) {
    navigator.clipboard.writeText(text);
    copied = true;
    toast.success("Enlace de clonación copiado");
    setTimeout(() => copied = false, 2000);
  }
</script>

<div class="repo-row {isOlder ? 'older-row' : ''} {isNew() ? 'new-update' : ''}" onclick={() => goto(`/${repo.full_name}`)} role="link" tabindex="0" onkeydown={(e) => e.key === 'Enter' && goto(`/${repo.full_name}`)}>
  {#if isNew()}
    <span class="new-dot-indicator"></span>
  {/if}
  
  <div class="repo-avatar-container">
    {#if repo.avatar_url && !avatarError}
      <img src={repo.avatar_url} alt={repo.name} class="repo-avatar" onerror={() => avatarError = true} />
    {/if}
    <div class="repo-icon-placeholder" style="background-color: {getLanguageColor(repo.language || 'Plain Text')}20; color: {getLanguageColor(repo.language || 'Plain Text')}">
      <span class="material-symbols-rounded">inventory_2</span>
    </div>
  </div>

  <div class="main-info">
    <div class="title-line">
      <a href="/{repo.full_name}" class="repo-name" onclick={(e) => e.stopPropagation()}>
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
        <span class="material-symbols-rounded status-icon" title="Private">lock</span>
      {/if}
      {#if repo.fork}
        <span class="material-symbols-rounded status-icon" title="Fork">account_tree</span>
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
      <span class="material-symbols-rounded {isRecentCommit() ? 'blink-active' : ''}">history</span>
      {repo.num_commits || 0}
    </div>
    <div class="stat-item {repo.num_branches > 0 ? 'active' : 'dim'}" title="Branches">
      <span class="material-symbols-rounded {isRecentBranch() ? 'blink-active' : ''}">account_tree</span>
      {repo.num_branches || 0}
    </div>
    <div class="stat-item {repo.open_issues_count > 0 ? 'active' : 'dim'}" title="Issues">
      <span class="material-symbols-rounded {isRecentIssue() ? 'blink-active' : ''}">error</span>
      {repo.open_issues_count || 0}
    </div>
    <div class="stat-item {repo.open_pulls_count > 0 ? 'active' : 'dim'}" title="Pull Requests">
      <span class="material-symbols-rounded {isRecentPull() ? 'blink-active' : ''}">merge</span>
      {repo.open_pulls_count || 0}
    </div>
  </div>

  <div class="actions-info">
    <button class="btn-icon-minimal" onclick={(e) => { e.stopPropagation(); copyToClipboard(repo.clone_url); }} title="Copy Clone URL">
      <span class="material-symbols-rounded" class:success-icon={copied}>{copied ? 'check' : 'content_copy'}</span>
    </button>
    <span class="time">{timeAgo(repo.updated_at)}</span>
  </div>
</div>

<style>
  .repo-row {
    display: grid;
    grid-template-columns: 56px 1fr 140px 220px 120px;
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
  }

  .repo-row:last-child { border-bottom: none; }

  .repo-row:hover {
    background: rgba(var(--primary-rgb), 0.05);
    z-index: 10;
    box-shadow: inset 0 0 0 1px rgba(var(--primary-rgb), 0.1);
  }

  .repo-avatar-container {
    width: 48px;
    height: 48px;
    position: relative;
    border-radius: 12px;
    overflow: hidden;
    background: var(--md-sys-color-surface-variant);
    display: flex;
    align-items: center;
    justify-content: center;
    border: 1px solid var(--md-sys-color-outline-variant);
  }

  .repo-avatar {
    width: 100%;
    height: 100%;
    object-fit: cover;
    position: absolute;
    top: 0;
    left: 0;
    z-index: 2;
  }

  .repo-icon-placeholder {
    width: 100%;
    height: 100%;
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 1.5rem;
    z-index: 1;
  }

  .new-update {
    border-left: 3px solid var(--tertiary);
  }

  .new-dot-indicator {
    width: 8px;
    height: 8px;
    background: var(--tertiary);
    border-radius: 50%;
    position: absolute;
    left: 8px;
    top: 50%;
    transform: translateY(-50%);
    box-shadow: 0 0 12px var(--tertiary);
    animation: pulse-dot 2s infinite;
    z-index: 5;
  }

  @keyframes pulse-dot {
    0% { transform: translateY(-50%) scale(1); opacity: 1; }
    50% { transform: translateY(-50%) scale(1.6); opacity: 0.4; }
    100% { transform: translateY(-50%) scale(1); opacity: 1; }
  }

  /* BLINK ANIMATIONS FOR COUNTERS */
  @keyframes activity-blink {
    0%, 100% { opacity: 1; transform: scale(1); filter: drop-shadow(0 0 0px var(--tertiary)); }
    50% { opacity: 0.7; transform: scale(1.1); filter: drop-shadow(0 0 5px var(--tertiary)); color: var(--tertiary); }
  }

  .blink-active {
    animation: activity-blink 1.5s infinite ease-in-out;
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

  .desc-line { display: flex; align-items: center; gap: 0.5rem; min-width: 0; }
  .status-icon { font-size: 1rem; color: var(--text-muted); opacity: 0.6; flex-shrink: 0; }
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
  .repo-desc { font-size: 0.85rem; color: var(--text-muted); margin: 0; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }

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
  
  .stat-item .material-symbols-rounded {
    font-size: 1.1rem;
  }

  .stat-item.dim { opacity: 0.2; }
  .stat-item.active { color: var(--primary-color); font-weight: 700; }

  .actions-info { display: flex; align-items: center; justify-content: flex-end; gap: 1rem; }
  .success-icon { color: var(--md-sys-color-primary) !important; }
  .time { font-size: 0.8rem; color: var(--text-muted); font-variant-numeric: tabular-nums; flex-shrink: 0; }

  @media (max-width: 1200px) {
    .repo-row { 
      grid-template-columns: 56px 1fr 140px 180px 100px; 
    }
    .stats-info { gap: 0.75rem; }
  }

  @media (max-width: 1000px) {
    .repo-row { 
      grid-template-columns: 56px 1fr 140px 100px; 
    }
    .stats-info { display: none; }
  }

  @media (max-width: 800px) {
    .repo-row { grid-template-columns: 56px 1fr 80px; }
    .tech-info { display: none; }
  }
</style>
