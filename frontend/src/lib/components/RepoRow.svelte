<script lang="ts">
  import { getLanguageColor } from '$lib/utils/languages';

  let { repo } = $props();

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
</script>

<div class="repo-row">
  <div class="main-info">
    <div class="title-line">
      <a href="/{repo.full_name}" class="repo-name">
        {#if repo.private}
          <span class="material-symbols-rounded icon-sm" title="Private">lock</span>
        {/if}
        {#if repo.fork}
          <span class="material-symbols-rounded icon-sm">account_tree</span>
        {/if}
        {repo.name}
      </a>
      {#if repo.empty}
        <span class="badge-empty">Empty</span>
      {/if}
    </div>
    <p class="repo-desc">
      {repo.description || 'No description provided.'}
    </p>
  </div>

  <div class="tech-info">
    {#if repo.language}
      <div class="lang">
        <span class="lang-dot" style="background-color: {getLanguageColor(repo.language)}"></span>
        {repo.language}
      </div>
    {:else if repo.empty}
      <div class="setup-hint">Ready for code</div>
    {:else}
      <div class="lang-none">N/A</div>
    {/if}
  </div>

  <div class="stats-info">
    <div class="stat-item {repo.open_issues_count > 0 ? 'active' : 'dim'}" title="Open Issues">
      <span class="material-symbols-rounded">error</span>
      {repo.open_issues_count || 0}
    </div>
    <div class="stat-item dim" title="Pull Requests">
      <span class="material-symbols-rounded">merge</span>
      0
    </div>
  </div>

  <div class="activity-info">
    <span class="time">{timeAgo(repo.updated_at)}</span>
  </div>
</div>

<style>
  .repo-row {
    display: grid;
    grid-template-columns: 1fr 120px 140px 100px;
    align-items: center;
    padding: 0.85rem 1.25rem;
    border-bottom: 1px solid var(--border-color);
    transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1);
    gap: 1rem;
    position: relative;
    z-index: 1;
  }

  .repo-row:hover {
    background: var(--md3-surface-3);
    transform: scale(1.005) translateX(4px);
    z-index: 2;
    border-color: transparent;
    box-shadow: 0 4px 12px rgba(0,0,0,0.1);
    border-radius: 0.75rem;
  }

  .repo-row:last-child {
    border-bottom: none;
  }

  .main-info {
    display: flex;
    flex-direction: column;
    min-width: 0;
  }

  .repo-name {
    font-weight: 600;
    color: var(--primary-color);
    text-decoration: none;
    font-size: 1.05rem;
    display: inline-flex;
    align-items: center;
    gap: 0.4rem;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }

  .repo-name:hover {
    text-decoration: underline;
  }

  .repo-desc {
    font-size: 0.85rem;
    color: var(--text-muted);
    margin: 0.2rem 0 0 0;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
    font-style: italic;
  }

  .badge-empty {
    font-size: 0.7rem;
    background: var(--md3-surface-2);
    color: var(--text-muted);
    padding: 0.15rem 0.6rem;
    border-radius: 1rem;
    border: 1px solid var(--border-color);
    text-transform: uppercase;
    font-weight: 800;
    letter-spacing: 0.05em;
  }

  .setup-hint {
    font-size: 0.8rem;
    color: var(--primary-color);
    opacity: 0.8;
    font-weight: 600;
  }

  .lang-none {
    opacity: 0.4;
    font-size: 0.8rem;
  }

  .stat-item.dim {
    opacity: 0.25;
  }

  .stat-item.active {
    color: var(--primary-color);
    font-weight: 600;
  }

  .icon-sm {
    font-size: 1.1rem;
    opacity: 0.7;
  }

  .tech-info {
    font-size: 0.85rem;
    color: var(--text-muted);
  }

  .lang {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    font-weight: 500;
  }

  .lang-dot {
    width: 10px;
    height: 10px;
    border-radius: 50%;
    opacity: 0.9;
    box-shadow: 0 0 4px rgba(0,0,0,0.1);
  }

  .stats-info {
    display: flex;
    gap: 1.25rem;
    font-size: 0.85rem;
    color: var(--text-muted);
  }

  .stat-item {
    display: flex;
    align-items: center;
    gap: 0.4rem;
  }

  .stat-item span {
    font-size: 1.15rem;
  }

  .activity-info {
    text-align: right;
    font-size: 0.85rem;
    color: var(--text-muted);
    font-variant-numeric: tabular-nums;
  }

  @media (max-width: 800px) {
    .repo-row {
      grid-template-columns: 1fr 100px;
    }
    .tech-info, .stats-info {
      display: none;
    }
  }
</style>
