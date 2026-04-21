<script lang="ts">
  import { apiFetch } from "$lib/api";
  import { onMount, getContext } from 'svelte';
  
  let { index } = $props();
  const repo: any = getContext('repo');
  
  let issue = $state<any>(null);
  let comments = $state<any[]>([]);
  let loading = $state(true);
  let error = $state('');
  
  let newCommentBody = $state('');
  let submitting = $state(false);

  async function fetchIssueDetail() {
    loading = true;
    error = '';
    try {
      const [issueRes, commentsRes] = await Promise.all([
        apiFetch(`/api/v2/${repo.data.owner.username}/${repo.data.name}/issues/${index}`),
        apiFetch(`/api/v2/${repo.data.owner.username}/${repo.data.name}/issues/${index}/comments`)
      ]);

      if (issueRes.ok && commentsRes.ok) {
        issue = await issueRes.json();
        comments = await commentsRes.json();
      } else {
        error = 'Issue not found.';
      }
    } catch (e) {
      error = 'Failed to load issue details.';
    } finally {
      loading = false;
    }
  }

  async function postComment() {
    if (!newCommentBody.trim()) return;
    submitting = true;
    try {
      const res = await apiFetch(`/api/v2/${repo.data.owner.username}/${repo.data.name}/issues/${index}/comments`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ body: newCommentBody })
      });
      if (res.ok) {
        const newComment = await res.json();
        comments = [...comments, newComment];
        newCommentBody = '';
      }
    } catch (e) {
      console.error(e);
    } finally {
      submitting = false;
    }
  }

  onMount(() => {
    if (repo.data) fetchIssueDetail();
  });

  function timeAgo(dateInput: any) {
    if (!dateInput) return '';
    const date = new Date(dateInput);
    const seconds = Math.floor((Date.now() - date.getTime()) / 1000);
    if (seconds < 60) return 'just now';
    const minutes = Math.floor(seconds / 60);
    if (minutes < 60) return `${minutes}m ago`;
    const hours = Math.floor(minutes / 60);
    if (hours < 24) return `${hours}h ago`;
    return date.toLocaleDateString();
  }
</script>

<div class="issue-detail-view">
  {#if loading}
    <div class="skeleton-detail island island-interactive"></div>
  {:else if error}
    <div class="error-msg island island-interactive">{error}</div>
  {:else if issue}
    <header class="issue-main-header">
      <div class="issue-title-row">
        <h1>{issue.title} <span class="index-tag">#{issue.number}</span></h1>
        <div class="status-pill" class:closed={issue.state === 'closed'}>
          <span class="material-symbols-rounded">
            {issue.state === 'closed' ? 'check_circle' : 'error'}
          </span>
          {issue.state === 'closed' ? 'Closed' : 'Open'}
        </div>
      </div>
      <div class="issue-subtitle">
        <strong>{issue.user?.username}</strong> opened this issue {timeAgo(issue.created_at)} · {issue.comments} comments
      </div>
    </header>

    <div class="issue-layout">
      <div class="timeline">
        <!-- Original post -->
        <div class="timeline-item main-post">
          <div class="item-avatar">
            <span class="material-symbols-rounded">person</span>
          </div>
          <div class="item-card-flat">
            <div class="item-header">
              <strong>{issue.user?.username}</strong> commented {timeAgo(issue.created_at)}
            </div>
            <div class="item-body markdown-body">
              {@html issue.body_html || issue.body}
            </div>
          </div>
        </div>

        <div class="timeline-divider"></div>

        <!-- Comments -->
        {#each comments as comment}
          <div class="timeline-item">
            <div class="item-avatar">
              <span class="material-symbols-rounded">chat_bubble</span>
            </div>
            <div class="item-card-flat">
              <div class="item-header">
                <strong>{comment.user?.username}</strong> commented {timeAgo(comment.created_at)}
              </div>
              <div class="item-body markdown-body">
                {@html comment.body_html || comment.body}
              </div>
            </div>
          </div>
        {/each}

        <!-- Add Comment -->
        <div class="timeline-item new-comment-item">
          <div class="item-avatar pulse">
            <span class="material-symbols-rounded">add_comment</span>
          </div>
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
                {submitting ? 'Posting...' : 'Post Comment'}
              </button>
            </div>
          </div>
        </div>
      </div>

      <aside class="sidebar">
        <div class="side-item">
          <h4 class="side-title">Assignee</h4>
          <div class="side-val">
            {#if issue.assignee}
               <div class="user-pill">
                 <span class="material-symbols-rounded">account_circle</span>
                 {issue.assignee.username}
               </div>
            {:else}
              <span class="muted">No one assigned</span>
            {/if}
          </div>
        </div>
        <div class="side-item">
            <h4 class="side-title">Labels</h4>
            <div class="side-val labels-grid">
              {#each issue.labels || [] as label}
                <span class="label-tag" style="background: {label.color}">
                  {label.name}
                </span>
              {/each}
              {#if !issue.labels?.length}
                <span class="muted">None yet</span>
              {/if}
            </div>
          </div>
      </aside>
    </div>
  {/if}
</div>

<style>
  .issue-detail-view { display: flex; flex-direction: column; gap: 2rem; }

  .issue-main-header { border-bottom: 1px solid var(--md-sys-color-outline-variant); padding-bottom: 1.5rem; }
  .issue-title-row { display: flex; align-items: center; gap: 1.5rem; flex-wrap: wrap; margin-bottom: 0.5rem; }
  .issue-title-row h1 { margin: 0; font-size: 2rem; font-weight: 800; }
  .index-tag { color: var(--text-muted); font-weight: 400; }
  
  .status-pill { display: flex; align-items: center; gap: 0.5rem; padding: 0.4rem 1.25rem; border-radius: 2rem; color: var(--on-primary); font-weight: 700; background: var(--success); font-size: 0.9rem; }
  .status-pill.closed { background: var(--error); }
  .issue-subtitle { color: var(--text-muted); font-size: 0.95rem; }

  .issue-layout { display: grid; grid-template-columns: 1fr 280px; gap: 3rem; }

  .timeline { display: flex; flex-direction: column; gap: 2rem; position: relative; }
  .timeline::before { content: ''; position: absolute; left: 1.5rem; top: 0; bottom: 0; width: 2px; background: rgba(255,255,255,0.05); z-index: 0; }

  .timeline-item { display: flex; gap: 1.5rem; position: relative; z-index: 1; }
  .item-avatar { width: 3rem; height: 3rem; background: var(--md3-surface-2); border-radius: 50%; display: flex; align-items: center; justify-content: center; color: var(--primary-color); flex-shrink: 0; border: 2px solid var(--md3-surface-2); box-shadow: 0 0 0 4px var(--bg-color); }
  .item-avatar.pulse { animation: borderPulse 2s infinite; }
  @keyframes borderPulse { 0%, 100% { border-color: transparent; } 50% { border-color: var(--primary-color); } }

  .item-card-flat { 
    flex: 1; 
    min-width: 0; 
    border: 1px solid var(--md-sys-color-outline-variant); 
    background: var(--bg-color); 
    border-radius: 1rem; 
    overflow: hidden; 
  }
  .item-header { padding: 0.75rem 1.25rem; background: var(--surface-variant); border-bottom: 1px solid var(--md-sys-color-outline-variant); font-size: 0.85rem; color: var(--text-muted); }
  .item-body { padding: 1.25rem; line-height: 1.6; }

  .timeline-divider { height: 1px; }

  .new-comment-item textarea { width: 100%; min-height: 120px; background: transparent; border: none; padding: 1.25rem; color: var(--text-color); font-family: inherit; font-size: 1rem; resize: vertical; border-bottom: 1px solid var(--md-sys-color-outline-variant); outline: none; }
  .comment-actions { padding: 0.75rem 1.25rem; display: flex; justify-content: flex-end; }
  
  .btn-action-primary { background: var(--primary-color); color: var(--on-primary); border: none; padding: 0.5rem 1.25rem; border-radius: 10px; font-weight: 700; cursor: pointer; transition: all 0.2s; }
  .btn-action-primary:hover:not(:disabled) { background: var(--primary-hover); transform: translateY(-1px); }
  .btn-action-primary:disabled { opacity: 0.5; cursor: not-allowed; }

  .sidebar { display: flex; flex-direction: column; gap: 2rem; }
  .side-item { border-bottom: 1px solid var(--md-sys-color-outline-variant); padding-bottom: 1.5rem; }
  .side-title { font-size: 0.75rem; font-weight: 800; text-transform: uppercase; letter-spacing: 0.05em; color: var(--text-muted); margin-bottom: 1rem; }
  .user-pill { display: flex; align-items: center; gap: 0.5rem; font-weight: 600; font-size: 0.9rem; }
  
  .labels-grid { display: flex; flex-wrap: wrap; gap: 0.5rem; }
  .label-tag { font-size: 0.75rem; font-weight: 700; padding: 0.1rem 0.6rem; border-radius: 2rem; color: white; border: 1px solid rgba(255,255,255,0.2); }

  .muted { color: var(--text-muted); font-size: 0.85rem; opacity: 0.5; }

  /* Markdown adjustment */
  :global(.markdown-body) { font-size: 1rem; color: rgba(255,255,255,0.85); }
  :global(.markdown-body p) { margin-bottom: 1rem; }

  .skeleton-detail { height: 400px; animation: pulse 1.5s infinite; }
  @keyframes pulse { 0%, 100% { opacity: 0.5; } 50% { opacity: 0.8; } }
</style>
