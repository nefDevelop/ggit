<script lang="ts">
  import { apiFetch } from "$lib/api";
  import { onMount, getContext } from "svelte";
  import Avatar from "$lib/components/Avatar.svelte";
  import MarkdownRenderer from "$lib/components/ui/MarkdownRenderer.svelte";
  import Button from "$lib/components/ui/Button.svelte";
  import Select from "$lib/components/ui/Select.svelte";

  let { index } = $props();
  const repo: any = getContext("repo");

  let issue = $state<any>(null);
  let comments = $state<any[]>([]);
  let allLabels = $state<any[]>([]);
  let allMilestones = $state<any[]>([]);
  let loading = $state(true);
  let error = $state("");

  let showLabelEditor = $state(false);
  let showMilestoneEditor = $state(false);

  let labelOptions = $derived(allLabels.map((l: any) => ({ value: l.id, label: l.name })));
  let milestoneOptions = $derived([
    { value: 0, label: 'No Milestone' },
    ...allMilestones.map(m => ({ value: m.id, label: m.title }))
  ]);

  let newCommentBody = $state("");
  let submitting = $state(false);

  async function fetchIssueDetail() {
    loading = true;
    error = "";
    try {
      const [issueRes, commentsRes, labelsRes, milestonesRes] = await Promise.all([
        apiFetch(
          `/api/v2/${repo.data.owner.username}/${repo.data.name}/issues/${index}`,
        ),
        apiFetch(
          `/api/v2/${repo.data.owner.username}/${repo.data.name}/issues/${index}/comments`,
        ),
        apiFetch(
            `/api/v2/${repo.data.owner.username}/${repo.data.name}/labels`,
        ),
        apiFetch(
            `/api/v2/${repo.data.owner.username}/${repo.data.name}/milestones`,
        ),
      ]);

      if (issueRes.ok && commentsRes.ok) {
        issue = await issueRes.json();
        comments = await commentsRes.json();
        if (labelsRes.ok) allLabels = await labelsRes.json();
        if (milestonesRes.ok) allMilestones = await milestonesRes.json();
      } else {
        error = "Issue not found.";
      }
    } catch (e) {
      error = "Failed to load issue details.";
    } finally {
      loading = false;
    }
  }

  async function postComment() {
    if (!newCommentBody.trim()) return;
    submitting = true;
    try {
      const res = await apiFetch(
        `/api/v2/${repo.data.owner.username}/${repo.data.name}/issues/${index}/comments`,
        {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({ body: newCommentBody }),
        },
      );
      if (res.ok) {
        const newComment = await res.json();
        comments = [...comments, newComment];
        newCommentBody = "";
      }
    } catch (e) {
      console.error(e);
    } finally {
      submitting = false;
    }
  }

  async function updateIssue(patch: any) {
    try {
      const res = await apiFetch(
        `/api/v2/${repo.data.owner.username}/${repo.data.name}/issues/${index}`,
        {
          method: "PATCH",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify(patch),
        },
      );
      if (res.ok) {
        const updated = await res.json();
        issue = updated;
      }
    } catch (e) {
      console.error(e);
    }
  }

  onMount(() => {
    if (repo.data) fetchIssueDetail();
  });

  function timeAgo(dateInput: any) {
    if (!dateInput) return "";
    const date = new Date(dateInput);
    const seconds = Math.floor((Date.now() - date.getTime()) / 1000);
    if (seconds < 60) return "just now";
    const minutes = Math.floor(seconds / 60);
    if (minutes < 60) return `${minutes}m ago`;
    const hours = Math.floor(minutes / 60);
    if (hours < 24) return `${hours}h ago`;
    return date.toLocaleDateString();
  }
</script>

<div class="issue-detail-view">
  {#if loading}
    <div class="skeleton-detail"></div>
  {:else if error}
    <div class="error-msg">{error}</div>
  {:else if issue}
    <header class="issue-main-header">
      <div class="issue-title-row">
        <h1>{issue.title} <span class="index-tag">#{issue.number}</span></h1>
        <div class="status-pill" class:closed={issue.state === "closed"}>
          <span class="material-symbols-rounded">
            {issue.state === "closed" ? "check_circle" : "error"}
          </span>
          {issue.state === "closed" ? "Closed" : "Open"}
        </div>
      </div>
      <div class="issue-subtitle">
        <a href="/{issue.user?.username}" class="user-link">{issue.user?.username}</a> opened this issue {timeAgo(
          issue.created_at,
        )} · {issue.comments} comments
      </div>
    </header>

    <div class="issue-layout">
      <div class="timeline">
        <!-- Original post -->
        <div class="timeline-item main-post">
          <Avatar
            src={issue.user?.avatar_url}
            username={issue.user?.username}
            size={48}
            class="item-avatar"
          />
          <div class="item-card-flat">
            <div class="item-header">
              <a href="/{issue.user?.username}" class="user-link">{issue.user?.username}</a> commented {timeAgo(
                issue.created_at,
              )}
            </div>
            <div class="item-body">
              <MarkdownRenderer 
                content={issue.body} 
                repoContext={{owner: repo.data.owner.username, name: repo.data.name}}
              />
            </div>
          </div>
        </div>

        <div class="timeline-divider"></div>

        <!-- Comments -->
        {#each comments as comment}
          <div class="timeline-item">
            <Avatar
              src={comment.user?.avatar_url}
              username={comment.user?.username}
              size={48}
              class="item-avatar"
            />
            <div class="item-card-flat">
              <div class="item-header">
                <a href="/{comment.user?.username}" class="user-link">{comment.user?.username}</a> commented {timeAgo(
                  comment.created_at,
                )}
              </div>
              <div class="item-body">
                <MarkdownRenderer 
                  content={comment.body} 
                  repoContext={{owner: repo.data.owner.username, name: repo.data.name}}
                />
              </div>
            </div>
          </div>
        {/each}

        <!-- Add Comment -->
        <div class="timeline-item new-comment-item">
          <Avatar
            username={issue.user?.username}
            size={48}
            class="item-avatar pulse"
          />
          <div class="item-card-flat">
            <textarea
              placeholder="Leave a comment..."
              bind:value={newCommentBody}
              disabled={submitting}
            ></textarea>
            <div class="comment-actions">
              <button
                class="btn-action-primary"
                onclick={postComment}
                disabled={submitting || !newCommentBody.trim()}
              >
                {submitting ? "Posting..." : "Post Comment"}
              </button>
            </div>
          </div>
        </div>
      </div>

      <aside class="sidebar">
        <div class="side-item">
          <div class="side-header">
            <h4 class="side-title">Assignee</h4>
            <Button variant="ghost" icon="settings" size="sm" />
          </div>
          <div class="side-val">
            {#if issue.assignee}
              <div class="user-pill">
                <Avatar
                  src={issue.assignee.avatar_url}
                  username={issue.assignee.username}
                  size={24}
                />
                {issue.assignee.username}
              </div>
            {:else}
              <span class="muted">No one assigned</span>
            {/if}
          </div>
        </div>

        <div class="side-item">
          <div class="side-header">
            <h4 class="side-title">Labels</h4>
            <Button 
              variant="ghost" 
              icon={showLabelEditor ? "close" : "settings"} 
              size="sm" 
              onclick={() => showLabelEditor = !showLabelEditor} 
            />
          </div>
          <div class="side-val">
            {#if showLabelEditor}
              <div class="editor-popover">
                <Select 
                  options={[{value: 0, label: 'Add Label...'}, ...labelOptions.filter(opt => !issue.labels?.some((l: any) => l.id === opt.value))]} 
                  value={0}
                  onchange={(id) => {
                    if (id === 0) return;
                    const newLabels = [...(issue.labels || []).map((l: any) => l.id), id];
                    updateIssue({ label_ids: newLabels });
                  }}
                />
              </div>
            {/if}
            <div class="labels-grid">
              {#each issue.labels || [] as label}
                <span class="label-tag" style="background: {label.color}">
                  {label.name}
                  {#if showLabelEditor}
                    <button class="remove-label" onclick={() => {
                      const newLabels = issue.labels.filter((l: any) => l.id !== label.id).map((l: any) => l.id);
                      updateIssue({ label_ids: newLabels });
                    }}>×</button>
                  {/if}
                </span>
              {/each}
              {#if !issue.labels?.length && !showLabelEditor}
                <span class="muted">None yet</span>
              {/if}
            </div>
          </div>
        </div>

        <div class="side-item">
          <div class="side-header">
            <h4 class="side-title">Milestone</h4>
            <Button 
              variant="ghost" 
              icon={showMilestoneEditor ? "close" : "settings"} 
              size="sm" 
              onclick={() => showMilestoneEditor = !showMilestoneEditor} 
            />
          </div>
          <div class="side-val">
            {#if showMilestoneEditor}
              <div class="editor-popover">
                <Select 
                  options={milestoneOptions} 
                  value={issue.milestone?.id || 0}
                  onchange={(id) => {
                    updateIssue({ milestone_id: id === 0 ? -1 : id }); // -1 usually means remove in some APIs, or 0. Check Gogs API.
                    showMilestoneEditor = false;
                  }}
                />
              </div>
            {:else if issue.milestone}
              <div class="milestone-pill">
                <span class="material-symbols-rounded mini">event</span>
                {issue.milestone.title}
              </div>
            {:else}
              <span class="muted">No milestone</span>
            {/if}
          </div>
        </div>
      </aside>
    </div>
  {/if}
</div>

<style>
  .issue-detail-view {
    display: flex;
    flex-direction: column;
    gap: 2rem;
  }

  .issue-main-header {
    border-bottom: 1px solid var(--md-sys-color-outline-variant);
    padding-bottom: 1.5rem;
  }
  .issue-title-row {
    display: flex;
    align-items: center;
    gap: 1.5rem;
    flex-wrap: wrap;
    margin-bottom: 0.5rem;
  }
  .issue-title-row h1 {
    margin: 0;
    font-size: 2rem;
    font-weight: 800;
    line-height: 1.2;
  }
  .index-tag {
    color: var(--text-muted);
    font-weight: 400;
  }

  .status-pill {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    padding: 0.4rem 1.25rem;
    border-radius: 2rem;
    color: var(--md-sys-color-on-primary);
    font-weight: 700;
    background: var(--md-sys-color-primary);
    font-size: 0.9rem;
    box-shadow: var(--shadow-sm);
  }
  .status-pill.closed {
    background: var(--tertiary);
    color: var(--on-tertiary);
  }
  .issue-subtitle {
    color: var(--text-muted);
    font-size: 0.95rem;
  }

  .issue-layout {
    display: grid;
    grid-template-columns: 1fr 280px;
    gap: 3rem;
  }
  @media (max-width: 1000px) {
    .issue-layout {
      grid-template-columns: 1fr;
      gap: 2rem;
    }
    .sidebar { order: -1; }
  }

  .timeline {
    display: flex;
    flex-direction: column;
    gap: 2rem;
    position: relative;
  }
  .timeline::before {
    content: "";
    position: absolute;
    left: 1.5rem;
    top: 0;
    bottom: 0;
    width: 2px;
    background: rgba(255, 255, 255, 0.05);
    z-index: 0;
  }

  .timeline-item {
    display: flex;
    gap: 1.5rem;
    position: relative;
    z-index: 1;
  }
  :global(.item-avatar) {
    width: 3rem;
    height: 3rem;
    background: var(--md-sys-color-secondary);
    color: var(--md-sys-color-on-secondary);
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;
    flex-shrink: 0;
    border: 2px solid var(--md-sys-color-outline-variant);
    box-shadow: 0 0 0 4px var(--bg-color);
  }
  :global(.item-avatar.pulse) {
    animation: borderPulse 2s infinite;
  }
  @keyframes borderPulse {
    0%,
    100% {
      border-color: transparent;
    }
    50% {
      border-color: var(--primary-color);
    }
  }

  .item-card-flat {
    flex: 1;
    min-width: 0;
    border: 1px solid var(--md-sys-color-outline-variant);
    background: var(--bg-color);
    border-radius: 1rem;
    overflow: hidden;
  }
  .item-header {
    padding: 0.75rem 1.25rem;
    background: var(--md-sys-color-secondary-container);
    border-bottom: 1px solid var(--md-sys-color-outline-variant);
    font-size: 0.85rem;
    color: var(--md-sys-color-on-secondary-container);
    font-weight: 600;
  }
  .item-body {
    padding: 1.25rem;
    line-height: 1.6;
  }

  .timeline-divider {
    height: 1px;
  }

  .new-comment-item textarea {
    width: 100%;
    min-height: 120px;
    background: transparent;
    border: none;
    padding: 1.25rem;
    color: var(--text-color);
    font-family: inherit;
    font-size: 1rem;
    resize: vertical;
    border-bottom: 1px solid var(--md-sys-color-outline-variant);
    outline: none;
  }
  .comment-actions {
    padding: 0.75rem 1.25rem;
    display: flex;
    justify-content: flex-end;
  }

  .btn-action-primary {
    background: var(--primary-color);
    color: var(--on-primary);
    border: none;
    padding: 0.5rem 1.25rem;
    border-radius: 10px;
    font-weight: 700;
    cursor: pointer;
    transition: all 0.2s;
  }
  .btn-action-primary:hover:not(:disabled) {
    background: var(--primary-hover);
    transform: translateY(-1px);
  }
  .btn-action-primary:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }

  .sidebar {
    display: flex;
    flex-direction: column;
    gap: 2rem;
  }
  .side-item {
    border-bottom: 1px solid var(--md-sys-color-outline-variant);
    padding-bottom: 1.5rem;
  }
  .side-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 0.75rem;
  }
  .side-title {
    font-size: 0.75rem;
    font-weight: 800;
    text-transform: uppercase;
    letter-spacing: 0.05em;
    color: var(--text-muted);
    margin: 0;
  }
  
  .editor-popover {
    margin-bottom: 1rem;
    animation: slideDown 0.2s ease-out;
  }
  @keyframes slideDown {
    from { opacity: 0; transform: translateY(-10px); }
    to { opacity: 1; transform: translateY(0); }
  }

  .user-pill, .milestone-pill {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    font-weight: 600;
    font-size: 0.9rem;
    background: var(--surface-variant);
    padding: 0.4rem 0.75rem;
    border-radius: 8px;
    width: fit-content;
  }
  .milestone-pill .mini {
    font-size: 1.1rem;
    opacity: 0.7;
  }

  .labels-grid {
    display: flex;
    flex-wrap: wrap;
    gap: 0.5rem;
  }
  .label-tag {
    font-size: 0.75rem;
    font-weight: 700;
    padding: 0.2rem 0.6rem;
    border-radius: 2rem;
    color: white;
    border: 1px solid rgba(255, 255, 255, 0.2);
    display: flex;
    align-items: center;
    gap: 0.4rem;
  }
  .remove-label {
    background: rgba(0, 0, 0, 0.2);
    border: none;
    color: white;
    width: 16px;
    height: 16px;
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;
    cursor: pointer;
    font-size: 12px;
    line-height: 1;
    padding: 0;
  }
  .remove-label:hover {
    background: rgba(255, 255, 255, 0.3);
  }

  .muted {
    color: var(--text-muted);
    font-size: 0.85rem;
    opacity: 0.5;
    padding-left: 0.25rem;
  }

  .user-link {
    color: var(--text-color);
    text-decoration: none;
    font-weight: 800;
  }
  .user-link:hover {
    color: var(--primary-color);
    text-decoration: underline;
  }

  /* Markdown adjustment */
  :global(.markdown-body) {
    font-size: 1rem;
    color: rgba(255, 255, 255, 0.85);
  }
  :global(.markdown-body p) {
    margin-bottom: 1rem;
  }

  .skeleton-detail {
    height: 400px;
    animation: pulse 1.5s infinite;
  }
  @keyframes pulse {
    0%,
    100% {
      opacity: 0.5;
    }
    50% {
      opacity: 0.8;
    }
  }
</style>
