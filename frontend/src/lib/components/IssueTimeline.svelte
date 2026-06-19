<script lang="ts">
  import Avatar from "$lib/components/Avatar.svelte";
  import MarkdownRenderer from "$lib/components/ui/MarkdownRenderer.svelte";
  import { apiFetch } from "$lib/api";

  let { issue = null, comments = [], repoOwner = "", repoName = "", onCommentPosted }: {
    issue: any;
    comments: any[];
    repoOwner: string;
    repoName: string;
    onCommentPosted?: (comment: any) => void;
  } = $props();

  let newCommentBody = $state("");
  let submitting = $state(false);

  async function postComment(index: number) {
    if (!newCommentBody.trim()) return;
    submitting = true;
    try {
      const res = await apiFetch(
        `/api/v2/${repoOwner}/${repoName}/issues/${index}/comments`,
        { method: "POST", headers: { "Content-Type": "application/json" }, body: JSON.stringify({ body: newCommentBody }) },
      );
      if (res.ok) {
        const newComment = await res.json();
        onCommentPosted?.(newComment);
        newCommentBody = "";
      }
    } catch (e) {
      console.error(e);
    } finally {
      submitting = false;
    }
  }

  function timeAgo(dateInput: any) {
    if (!dateInput) return "";
    const date = new Date(dateInput);
    const seconds = Math.floor((Date.now() - date.getTime()) / 1000);
    if (seconds < 60) return "just now";
    const minutes = Math.floor(seconds / 60);
    if (minutes < 60) return `${minutes}m ago`;
    const hours = Math.floor(minutes / 60);
    if (hours < 24) return `${hours}h ago`;
    const days = Math.floor(hours / 24);
    if (days < 30) return `${days}d ago`;
    return date.toLocaleDateString();
  }
</script>

<div class="timeline">
  <div class="timeline-item main-post">
    <Avatar src={issue?.user?.avatar_url} username={issue?.user?.username} size={48} class="item-avatar" />
    <div class="item-card-flat">
      <div class="item-header">
        <a href="/{issue?.user?.username}" class="user-link">{issue?.user?.username}</a> commented {timeAgo(issue?.created_at)}
      </div>
      <div class="item-body">
        <MarkdownRenderer content={issue?.body} repoContext={{owner: repoOwner, name: repoName}} />
      </div>
    </div>
  </div>

  <div class="timeline-divider"></div>

  {#each comments as comment}
    <div class="timeline-item">
      <Avatar src={comment.user?.avatar_url} username={comment.user?.username} size={48} class="item-avatar" />
      <div class="item-card-flat">
        <div class="item-header">
          <a href="/{comment.user?.username}" class="user-link">{comment.user?.username}</a> commented {timeAgo(comment.created_at)}
        </div>
        <div class="item-body">
          <MarkdownRenderer content={comment.body} repoContext={{owner: repoOwner, name: repoName}} />
        </div>
      </div>
    </div>
  {/each}

  <div class="timeline-item new-comment-item">
    <Avatar username={issue?.user?.username} size={48} class="item-avatar pulse" />
    <div class="item-card-flat">
      <textarea placeholder="Leave a comment..." bind:value={newCommentBody} disabled={submitting}></textarea>
      <div class="comment-actions">
        <button class="btn-action-primary" onclick={() => postComment(issue?.number)} disabled={submitting || !newCommentBody.trim()}>
          {submitting ? "Posting..." : "Post Comment"}
        </button>
      </div>
    </div>
  </div>
</div>

<style>
  .timeline { display: flex; flex-direction: column; gap: 0.25rem; }
  .timeline-item { display: flex; gap: 1rem; padding: 1rem 0; }
  .item-avatar { flex-shrink: 0; }
  .item-card-flat {
    flex: 1; background: var(--md-sys-color-surface-container-low);
    border-radius: 0.75rem; padding: 1rem 1.25rem;
    border: 1px solid var(--md-sys-color-outline-variant);
  }
  .item-header { font-size: 0.85rem; color: var(--text-muted); margin-bottom: 0.75rem; }
  .user-link { color: var(--primary-color); text-decoration: none; font-weight: 700; }
  .user-link:hover { text-decoration: underline; }
  .timeline-divider { height: 1px; background: var(--md-sys-color-outline-variant); margin: 0 0 0 60px; }
  .new-comment-item textarea {
    width: 100%; min-height: 100px; resize: vertical;
    background: rgba(0,0,0,0.2); border: 1px solid var(--md-sys-color-outline-variant);
    color: inherit; padding: 0.75rem; border-radius: 0.5rem;
    font-family: inherit; font-size: 0.95rem; outline: none; box-sizing: border-box;
  }
  .comment-actions { display: flex; justify-content: flex-end; gap: 0.5rem; margin-top: 0.75rem; }
  .btn-action-primary {
    background: var(--md-sys-color-primary); color: var(--md-sys-color-on-primary);
    border: none; padding: 0.5rem 1.25rem; border-radius: 0.5rem;
    font-weight: 700; cursor: pointer; font-size: 0.9rem;
  }
  .btn-action-primary:disabled { opacity: 0.5; cursor: not-allowed; }
</style>
