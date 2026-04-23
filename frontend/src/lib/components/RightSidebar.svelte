<script lang="ts">
  import { apiFetch } from "$lib/api";
  import { authState } from '$lib/auth.svelte';
  import { onMount } from 'svelte';
  import Avatar from '$lib/components/Avatar.svelte';

  let feeds = $state<any[]>([]);
  let loading = $state(true);

  async function loadFeeds() {
    try {
      const res = await apiFetch('/api/v2/user/dashboard/feeds');
      if (res.ok) {
        const data = await res.json();
        feeds = data.feeds;
      }
    } catch (e) {
      console.error("Failed to fetch activity feeds", e);
    } finally {
      loading = false;
    }
  }

  function formatAction(act: any) {
    try {
      switch (act.op_type) {
        case 1: 
          return `<span class="action-summary"><span class="material-symbols-rounded mini">add_box</span> <a href="/${act.repo_user_name}/${act.repo_name}">${act.repo_name}</a></span>`;
        case 5: // ActionCommitRepo
          const pushData = JSON.parse(act.content);
          return `
            <div class="action-summary">
              <span class="action-pill"><span class="material-symbols-rounded mini">history</span> ${pushData.Len}</span>
              <span class="action-pill"><span class="material-symbols-rounded mini">fork_right</span> ${act.ref_name}</span>
              <span class="repo-link-wrap">in <a href="/${act.repo_user_name}/${act.repo_name}" class="repo-link">${act.repo_name}</a></span>
            </div>
          `;
        case 11: // ActionMergePullRequest
          return `<span class="action-summary"><span class="material-symbols-rounded mini">call_merge</span> merged into <a href="/${act.repo_user_name}/${act.repo_name}">${act.repo_name}</a></span>`;
        default:
          return act.content || 'activity';
      }
    } catch (e) {
      return act.content;
    }
  }

  onMount(() => {
    loadFeeds();
  });
</script>

<aside class="right-sidebar">
  <div class="sidebar-inner">
    <header class="sidebar-header">
      <h2 class="sidebar-label">Recent Activity</h2>
      {#if feeds.length > 0}
        <span class="count-badge">{feeds.length}</span>
      {/if}
    </header>

    <div class="activity-scroll-area">
      {#if loading}
        <div class="skeletons">
          {#each Array(6) as _}<div class="skeleton-item"></div>{/each}
        </div>
      {:else if feeds.length === 0}
        <div class="empty-activity">
          <span class="material-symbols-rounded">notifications_off</span>
          <p>No activity yet</p>
        </div>
      {:else}
        <div class="activity-timeline">
          {#each feeds as act}
            <div class="activity-entry">
              <div class="entry-marker">
                <div class="avatar-box">
                  <Avatar src={act.act_avatar} username={act.act_user_name} size={32} class="user-avatar" />
                  <div class="type-icon {act.op_type === 5 ? 'push' : 'create'}">
                    <span class="material-symbols-rounded">
                      {act.op_type === 5 ? 'commit' : 'add_box'}
                    </span>
                  </div>
                </div>
                <div class="marker-line"></div>
              </div>
              <div class="entry-content">
                <div class="entry-meta">
                  <span class="user-name">{act.act_user_name}</span>
                  <span class="timestamp">{new Date(act.created_at).toLocaleDateString(undefined, { month: 'short', day: 'numeric' })}</span>
                </div>
                <div class="entry-body">
                  {@html formatAction(act)}
                </div>
              </div>
            </div>
          {/each}
        </div>
      {/if}
    </div>
  </div>
</aside>

<style>
  .right-sidebar {
    width: 320px;
    height: 100vh;
    padding: 0.75rem 0.75rem 0.75rem 0; /* Align with main content margin */
    box-sizing: border-box;
    display: flex;
    flex-direction: column;
    background: transparent;
  }

  .sidebar-inner {
    flex: 1;
    display: flex;
    flex-direction: column;
    min-height: 0;
    position: relative;
    padding-top: 1rem;
  }

  .sidebar-header {
    padding: 1rem 1rem 1.5rem;
    display: flex;
    align-items: center;
    justify-content: space-between;
  }

  .sidebar-label {
    font-size: 0.75rem;
    font-weight: 900;
    text-transform: uppercase;
    letter-spacing: 0.1em;
    color: var(--text-muted);
    margin: 0;
  }

  .count-badge {
    background: var(--surface-variant);
    color: var(--text-muted);
    padding: 2px 8px;
    border-radius: 12px;
    font-size: 0.65rem;
    font-weight: 800;
  }

  .activity-scroll-area {
    flex: 1;
    overflow-y: auto;
    padding: 0 1rem 2rem;
  }

  .activity-timeline {
    display: flex;
    flex-direction: column;
  }

  .activity-entry {
    display: flex;
    gap: 1rem;
    padding: 0.75rem 0.5rem;
    border-radius: 1rem;
    transition: all 0.2s ease;
    margin: 0 -0.5rem;
  }

  .activity-entry:hover {
    background: rgba(255, 255, 255, 0.03);
    transform: translateX(-4px);
  }

  .entry-marker {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 4px;
    flex-shrink: 0;
  }

  .avatar-box {
    position: relative;
    width: 32px;
    height: 32px;
  }

  :global(.user-avatar) {
    width: 100%;
    height: 100%;
    border: 2px solid var(--bg-color);
  }

  .type-icon {
    position: absolute;
    right: -4px;
    bottom: -4px;
    width: 16px;
    height: 16px;
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;
    border: 2px solid var(--bg-color);
  }

  .type-icon span { font-size: 10px; color: white; }
  .type-icon.push { background: #10b981; }
  .type-icon.create { background: #3b82f6; }

  .marker-line {
    width: 2px;
    flex: 1;
    background: var(--md-sys-color-outline-variant);
    opacity: 0.2;
    border-radius: 2px;
  }

  .activity-entry:last-child .marker-line { display: none; }

  .entry-content {
    flex: 1;
    min-width: 0;
  }

  .entry-meta {
    display: flex;
    align-items: center;
    justify-content: space-between;
    margin-bottom: 2px;
  }

  .user-name {
    font-size: 0.85rem;
    font-weight: 850;
    color: var(--text-color);
  }

  .timestamp {
    font-size: 0.7rem;
    color: var(--text-muted);
    opacity: 0.6;
    font-weight: 600;
  }

  .entry-body {
    font-size: 0.82rem;
    line-height: 1.4;
    color: var(--text-muted);
  }

  .entry-body :global(.action-summary) { display: flex; align-items: center; gap: 0.35rem; flex-wrap: wrap; }
  .entry-body :global(.action-pill) { background: var(--surface-variant); padding: 0px 6px; border-radius: 4px; font-size: 0.7rem; font-weight: 700; color: var(--text-color); display: flex; align-items: center; gap: 4px; }
  .entry-body :global(.repo-link) { font-weight: 700; color: var(--primary-color); text-decoration: none; }
  .entry-body :global(.repo-link:hover) { text-decoration: underline; }
  .entry-body :global(.repo-link-wrap) { opacity: 0.7; }

  .skeleton-item {
    height: 60px;
    background: var(--surface-variant);
    border-radius: 0.75rem;
    margin-bottom: 1rem;
    animation: pulse 1.5s infinite;
  }

  @keyframes pulse { 0%, 100% { opacity: 0.5; } 50% { opacity: 0.8; } }

  @media (max-width: 1300px) {
    .right-sidebar { display: none; }
  }

  .empty-activity {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    padding: 3rem 1rem;
    opacity: 0.3;
  }
  .empty-activity span { font-size: 2.5rem; margin-bottom: 0.5rem; }
  .empty-activity p { font-size: 0.8rem; font-weight: 600; }
</style>
