<script lang="ts">
  import { getLanguageColor } from '$lib/utils/languages';

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

  let copied = $state(false);
  function copyToClipboard(text: string) {
    navigator.clipboard.writeText(text);
    copied = true;
    setTimeout(() => copied = false, 2000);
  }
</script>

<div class="repo-row {isOlder ? 'older-row' : ''}">
  <div class="main-info">
    <div class="title-line">
      <a href="/{repo.full_name}" class="repo-name">
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
    <div class="stat-item dim" title="Commits">
      <span class="material-symbols-rounded">history</span>
      {repo.num_commits || 0}
    </div>
    <div class="stat-item dim" title="Branches">
      <span class="material-symbols-rounded">account_tree</span>
      {repo.num_branches || 0}
    </div>
    <div class="stat-item {repo.open_issues_count > 0 ? 'active' : 'dim'}" title="Issues">
      <span class="material-symbols-rounded">error</span>
      {repo.open_issues_count || 0}
    </div>
    <div class="stat-item {repo.open_pulls_count > 0 ? 'active' : 'dim'}" title="Pull Requests">
      <span class="material-symbols-rounded">merge</span>
      {repo.open_pulls_count || 0}
    </div>
  </div>

  <div class="actions-info">
    <button class="btn-icon-minimal" onclick={() => copyToClipboard(repo.clone_url)} title="Copy Clone URL">
      <span class="material-symbols-rounded">{copied ? 'check' : 'content_copy'}</span>
    </button>
    <span class="time">{timeAgo(repo.updated_at)}</span>
  </div>
</div>

<style>
  .repo-row {
    display: grid;
    grid-template-columns: 1fr 140px 220px 120px;
    align-items: center;
    padding: 1rem 1.25rem;
    background: transparent;
    transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1);
    gap: 1rem;
    position: relative;
    border: none;
  }

  .repo-row:hover {
    background: var(--md3-surface-3);
    transform: scale(1.01) translateX(8px);
    border-radius: 1rem;
    box-shadow: 0 8px 24px rgba(0,0,0,0.12);
    z-index: 10;
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
    background: rgba(0,0,0,0.2);
    margin: 0.1rem 0;
  }
  .lang-segment { height: 100%; }

  .desc-line { display: flex; align-items: center; gap: 0.5rem; min-width: 0; }
  .status-icon { font-size: 1rem; color: var(--text-muted); opacity: 0.6; flex-shrink: 0; }
  .repo-name { font-weight: 700; color: var(--primary-color); text-decoration: none; font-size: 1.1rem; display: inline-flex; align-items: center; gap: 0.5rem; }
  .repo-desc { font-size: 0.85rem; color: var(--text-muted); margin: 0; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }

  .badge-empty { font-size: 0.65rem; background: var(--md3-surface-2); padding: 0.1rem 0.5rem; border-radius: 2rem; font-weight: 800; letter-spacing: 0.05em; }
  .lang { display: flex; align-items: center; gap: 0.5rem; font-size: 0.85rem; font-weight: 600; }
  .lang-dot { width: 10px; height: 10px; border-radius: 50%; }

  .stats-info { display: flex; gap: 1.5rem; }
  .stat-item { display: flex; align-items: center; gap: 0.4rem; font-size: 0.85rem; }
  .stat-item.dim { opacity: 0.2; }
  .stat-item.active { color: var(--primary-color); font-weight: 700; }

  .actions-info { display: flex; align-items: center; justify-content: flex-end; gap: 1rem; }
  .time { font-size: 0.8rem; color: var(--text-muted); font-variant-numeric: tabular-nums; }

  .copy-btn { opacity: 0; transform: translateX(10px); transition: all 0.2s ease; background: var(--primary-color); color: var(--on-primary); border: none; width: 32px; height: 32px; border-radius: 50%; display: flex; align-items: center; justify-content: center; cursor: pointer; }
  .repo-row:hover .copy-btn { opacity: 1; transform: translateX(0); }

  @media (max-width: 900px) {
    .repo-row { grid-template-columns: 1fr 100px; }
    .tech-info, .stats-info { display: none; }
  }
</style>
