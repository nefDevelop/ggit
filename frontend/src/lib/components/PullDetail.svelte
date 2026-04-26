<script lang="ts">
  import { apiFetch } from "$lib/api";
  import { onMount, getContext } from 'svelte';
  import DiffViewer from '$lib/components/DiffViewer.svelte';
  import MarkdownRenderer from '$lib/components/ui/MarkdownRenderer.svelte';
  
  let { index } = $props();
  const repo: any = getContext('repo');
  
  let pull = $state<any>(null);
  let diff = $state<any>(null);
  let loading = $state(true);
  let error = $state('');
  
  let activeView = $state('conversation'); // 'conversation', 'commits', 'files'
  let comments = $state<any[]>([]);
  let newComment = $state('');
  let postingComment = $state(false);

  async function fetchPullDetail() {
    if (index === 'new') {
        loading = false;
        error = 'El flujo de creación de Pull Requests está siendo modernizado. Por ahora, usa la interfaz antigua.';
        return;
    }
    loading = true;
    error = '';
    try {
      const [pullRes, diffRes] = await Promise.all([
        apiFetch(`/api/v2/${repo.data.owner.username}/${repo.data.name}/pulls/${index}`),
        apiFetch(`/api/v2/${repo.data.owner.username}/${repo.data.name}/pulls/${index}/diff`)
      ]);

      if (pullRes.ok) pull = await pullRes.json();
      if (diffRes.ok) diff = await diffRes.json();
      
      if (!pull) {
        error = 'Pull request not found.';
        return;
      }

      // Fetch comments (as PR is an issue in Ggit)
      const commentsRes = await apiFetch(`/api/v2/${repo.data.owner.username}/${repo.data.name}/issues/${index}/comments`);
      if (commentsRes.ok) {
        comments = await commentsRes.json();
      }
    } catch (e) {
      error = 'Failed to load pull request details.';
    } finally {
      loading = false;
    }
  }

  async function postComment() {
    if (!newComment.trim()) return;
    postingComment = true;
    try {
      const res = await apiFetch(`/api/v2/${repo.data.owner.username}/${repo.data.name}/issues/${index}/comments`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ body: newComment })
      });
      if (res.ok) {
        const comment = await res.json();
        comments = [...comments, comment];
        newComment = '';
      }
    } catch (e) {
        console.error(e);
    } finally {
        postingComment = false;
    }
  }

  onMount(() => {
    if (repo.data) fetchPullDetail();
  });
</script>

<div class="pull-detail-view">
  {#if loading}
    <div class="skeleton-detail"></div>
  {:else if error}
    <div class="error-msg">{error}</div>
  {:else if pull}
    <header class="pull-detail-header">
       <div class="title-row">
         <h1>{pull.title} <span class="index">#{pull.number}</span></h1>
         <div class="status-pill" class:merged={pull.pull_request?.merged}>
            <span class="material-symbols-rounded">
              {pull.pull_request?.merged ? 'merge' : 'merge_type'}
            </span>
            {pull.pull_request?.merged ? 'Merged' : 'Open'}
         </div>
       </div>
       <div class="subtitle">
         <a href="/{pull.user?.username}" class="user-link">{pull.user?.username}</a> wants to merge into <code>{pull.pull_request?.base_branch}</code> from <code>{pull.pull_request?.head_branch}</code>
       </div>
    </header>

    <nav class="pull-tabs-bar">
      <button class="tab-pill" class:active={activeView === 'conversation'} onclick={() => activeView = 'conversation'}>
        <span class="material-symbols-rounded mini">chat</span>
        {activeView === 'conversation' ? 'Conversation' : ''}
      </button>
      <button class="tab-pill" class:active={activeView === 'files'} onclick={() => activeView = 'files'}>
        <span class="material-symbols-rounded mini">diff</span>
        {activeView === 'files' ? 'Files' : ''}
        <span class="count">{diff?.Files?.length || 0}</span>
      </button>
    </nav>

    <div class="pull-content-body">
      {#if activeView === 'conversation'}
        <div class="conversation-timeline">
             <div class="timeline-item-flat">
                <div class="item-header"><a href="/{pull.user?.username}" class="user-link">{pull.user?.username}</a> commented</div>
                <div class="item-body">
                  <MarkdownRenderer 
                    content={pull.body} 
                    repoContext={{owner: repo.data.owner.username, name: repo.data.name}}
                  />
                </div>
             </div>

             {#each comments as comment}
               <div class="timeline-item-flat">
                 <div class="item-header">
                   <a href="/{comment.user?.username}" class="user-link">{comment.user?.username}</a> commented on {new Date(comment.created_at).toLocaleDateString()}
                 </div>
                 <div class="item-body">
                   <MarkdownRenderer 
                     content={comment.body} 
                     repoContext={{owner: repo.data.owner.username, name: repo.data.name}}
                   />
                 </div>
               </div>
             {/each}

            <div class="comment-composer island tonal">
              <textarea bind:value={newComment} placeholder="Leave a comment..."></textarea>
              <div class="composer-footer">
                <button class="btn btn-primary" onclick={postComment} disabled={!newComment.trim() || postingComment}>
                   {postingComment ? 'Posting...' : 'Comment'}
                </button>
              </div>
            </div>
        </div>
      {:else if activeView === 'files'}
        <DiffViewer {diff} />
      {/if}
    </div>
  {/if}
</div>

<style>
  .pull-detail-view { display: flex; flex-direction: column; gap: 2rem; }
  
  .pull-detail-header { border-bottom: 1px solid var(--md-sys-color-outline-variant); padding-bottom: 1.5rem; }
  .title-row { display: flex; align-items: center; gap: 1.5rem; }
  .title-row h1 { margin: 0; font-size: 1.75rem; font-weight: 800; }
  .index { color: var(--text-muted); font-weight: 400; }
  
  .status-pill { display: flex; align-items: center; gap: 0.5rem; padding: 0.4rem 1.25rem; border-radius: 2rem; background: var(--success); color: var(--on-primary); font-weight: 700; font-size: 0.9rem; }
  .status-pill.merged { background: var(--primary); }
  
  .subtitle { color: var(--text-muted); font-size: 0.9rem; margin-top: 0.5rem; }
  .subtitle code { background: var(--surface-variant); padding: 2px 6px; border-radius: 4px; font-family: 'JetBrains Mono', monospace; font-weight: 700; color: var(--primary-color); }

  .pull-tabs-bar { display: flex; gap: 0.5rem; border-bottom: 1px solid var(--md-sys-color-outline-variant); margin-bottom: 1rem; }
  .nav-tab { background: transparent; border: none; padding: 0.85rem 1.25rem; color: var(--text-muted); font-weight: 700; cursor: pointer; display: flex; align-items: center; gap: 0.75rem; transition: all 0.2s; position: relative; }
  .nav-tab:hover { color: var(--text-color); background: var(--surface-variant); }
  .nav-tab.active { color: var(--primary-color); }
  .nav-tab.active::after { content: ''; position: absolute; bottom: -1px; left: 0; right: 0; height: 2px; background: var(--primary-color); }
  .nav-tab .count { font-size: 0.7rem; background: var(--surface-2); padding: 2px 6px; border-radius: 10px; opacity: 0.7; }

  .timeline-item-flat { 
    border: 1px solid var(--md-sys-color-outline-variant); 
    background: var(--bg-color); 
    border-radius: 1rem; 
    overflow: hidden; 
    margin-bottom: 2rem;
  }
  .item-header { padding: 0.75rem 1.5rem; background: var(--surface-variant); border-bottom: 1px solid var(--md-sys-color-outline-variant); font-size: 0.85rem; color: var(--text-muted); }
  .item-body { padding: 1.5rem; line-height: 1.6; }

  .comment-composer { display: flex; flex-direction: column; gap: 1rem; border: 2px solid var(--md-sys-color-outline-variant); border-radius: 1rem; overflow: hidden; }
  .comment-composer textarea { background: transparent; border: none; padding: 1.5rem; color: var(--text-color); font-family: inherit; font-size: 1rem; resize: vertical; min-height: 120px; outline: none; }
  .composer-footer { padding: 0.75rem 1.5rem; background: var(--surface-variant); border-top: 1px solid var(--md-sys-color-outline-variant); display: flex; justify-content: flex-end; }

  .skeleton-detail { height: 400px; background: var(--surface-variant); border-radius: 1rem; animation: pulse 1.5s infinite; }
  @keyframes pulse { 0%, 100% { opacity: 0.3; } 50% { opacity: 0.6; } }

  .error-msg { background: rgba(239, 68, 68, 0.1); color: #ef4444; padding: 2rem; border-radius: 1rem; text-align: center; font-weight: 700; }

  .user-link {
    color: var(--text-color);
    text-decoration: none;
    font-weight: 800;
  }
  .user-link:hover {
    color: var(--primary-color);
    text-decoration: underline;
  }
</style>
