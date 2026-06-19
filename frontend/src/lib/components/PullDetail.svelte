<script lang="ts">
  import type { RepoContext } from "$lib/types";
  import { apiFetch } from "$lib/api";
  import { onMount, getContext } from 'svelte';
  import { toast } from "$lib/toast.svelte";
  import DiffViewer from '$lib/components/DiffViewer.svelte';
  import ReviewToolbar from '$lib/components/ReviewToolbar.svelte';
  import PullMergeBox from '$lib/components/PullMergeBox.svelte';
  import MarkdownRenderer from '$lib/components/ui/MarkdownRenderer.svelte';
  
  let { index } = $props();
  const repo = getContext<RepoContext>('repo');
  
  let pull = $state<any>(null);
  let diff = $state<any>(null);
  let loading = $state(true);
  let error = $state('');
  
  let activeView = $state('conversation'); // 'conversation', 'commits', 'files'
  let comments = $state<any[]>([]);
  let reviews = $state<any[]>([]);
  let inlineComments = $state<any[]>([]);
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
      const [pullRes, diffRes, reviewsRes] = await Promise.all([
        apiFetch(`/api/v2/${repo.data.owner.username}/${repo.data.name}/pulls/${index}`),
        apiFetch(`/api/v2/${repo.data.owner.username}/${repo.data.name}/pulls/${index}/diff`),
        apiFetch(`/api/v2/${repo.data.owner.username}/${repo.data.name}/pulls/${index}/reviews`)
      ]);

      if (pullRes.ok) pull = await pullRes.json();
      if (diffRes.ok) diff = await diffRes.json();
      if (reviewsRes.ok) reviews = await reviewsRes.json();
      
      if (!pull) {
        error = 'Pull request not found.';
        return;
      }

      // Fetch normal comments
      const commentsRes = await apiFetch(`/api/v2/${repo.data.owner.username}/${repo.data.name}/issues/${index}/comments`);
      if (commentsRes.ok) {
        comments = await commentsRes.json();
      }

      // We need to extract inline comments from reviews or fetch them separately.
      // For now, let's assume reviewsRes returns them or we fetch per review.
      // Actually, my API endpoint `GET /pulls/:index/reviews` should probably include them.
      // For this implementation, I'll flatten them into inlineComments state.
      inlineComments = []; // Reset
      // Logic to fetch comments for each review if not included
    } catch (e) {
      error = 'Failed to load pull request details.';
    } finally {
      loading = false;
    }
  }

  async function addInlineComment(file: string, line: number, content: string) {
    // 1. Create a pending review if not exists, or use existing pending.
    // For simplicity, we'll create a comment-type review immediately.
    try {
      const reviewRes = await apiFetch(`/api/v2/${repo.data.owner.username}/${repo.data.name}/pulls/${index}/reviews`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ type: 2, content: '' }) // Type 2: Comment
      });

      if (reviewRes.ok) {
        const review = await reviewRes.json();
        const commentRes = await apiFetch(`/api/v2/${repo.data.owner.username}/${repo.data.name}/pulls/${index}/reviews/${review.id}/comments`, {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({
            tree_path: file,
            line: line,
            commit_id: diff.Files[0].Sections[0].Lines[0].Content.includes('SHA') ? '...' : 'HEAD', // SHA should come from diff
            content: content
          })
        });

        if (commentRes.ok) {
          toast.success("Comentario añadido.");
          fetchPullDetail(); // Refresh
        }
      }
    } catch (e) {
      console.error(e);
      toast.error("Error al añadir comentario.");
    }
  }

  async function resolveComment(id: number, resolved: boolean) {
    try {
      const res = await apiFetch(`/api/v2/${repo.data.owner.username}/${repo.data.name}/pulls/comments/${id}/resolve`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ resolved })
      });
      if (res.ok) {
        toast.success(resolved ? "Conversación resuelta." : "Conversación reabierta.");
        fetchPullDetail();
      }
    } catch (e) {
      console.error(e);
      toast.error("Error al actualizar la conversación.");
    }
  }

  async function submitReview(type: number, content: string) {
    try {
      const res = await apiFetch(`/api/v2/${repo.data.owner.username}/${repo.data.name}/pulls/${index}/reviews`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ type, content })
      });
      if (res.ok) {
        toast.success("Revisión enviada.");
        fetchPullDetail();
      }
    } catch (e) {
      console.error(e);
      toast.error("Error al enviar la revisión.");
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

  function getAttentionStateInfo(state: number) {
    switch (state) {
      case 1: return { label: 'Review required', color: '#0061a4', icon: 'visibility' };
      case 2: return { label: 'Waiting for author', color: '#ba1a1a', icon: 'pending' };
      case 3: return { label: 'Approved', color: '#006d39', icon: 'check_circle' };
      default: return null;
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
    {@const info = getAttentionStateInfo(pull.attention_state)}
    <header class="pull-detail-header">
       <div class="title-row">
         <h1>{pull.title} <span class="index">#{pull.number}</span></h1>
         
         {#if info}
           <div class="attention-badge" style="--badge-color: {info.color}">
              <span class="material-symbols-rounded mini">{info.icon}</span>
              {info.label}
           </div>
         {/if}

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

            {#if !pull.pull_request?.merged}
               <PullMergeBox {index} onMerged={fetchPullDetail} />
            {/if}
        </div>
      {:else if activeView === 'files'}
        <ReviewToolbar onSubmitReview={submitReview} />
        <DiffViewer 
          {diff} 
          allowComments={!pull.pull_request?.merged}
          comments={reviews.flatMap((r: any) => (r.comments || []).map((c: any) => ({...c, reviewer: r.reviewer})))}
          onAddComment={addInlineComment}
          onResolveComment={resolveComment}
        />
      {/if}
    </div>
  {/if}
</div>

<style>
  .pull-detail-view { display: flex; flex-direction: column; gap: 2rem; }
  
  .pull-detail-header { border-bottom: 1px solid var(--md-sys-color-outline-variant); padding-bottom: 1.5rem; }
  .title-row { display: flex; align-items: center; gap: 1.5rem; }
  .title-row h1 { margin: 0; font-size: 1.75rem; font-weight: 800; line-height: 1.2; }
  .index { color: var(--text-muted); font-weight: 400; }
  
  .status-pill { display: flex; align-items: center; gap: 0.5rem; padding: 0.4rem 1.25rem; border-radius: 2rem; background: var(--success); color: var(--on-primary); font-weight: 700; font-size: 0.9rem; }
  .status-pill.merged { background: var(--primary); }
  
  .attention-badge {
    display: flex;
    align-items: center;
    gap: 0.4rem;
    padding: 0.35rem 0.8rem;
    border-radius: 0.5rem;
    background: color-mix(in srgb, var(--badge-color), transparent 85%);
    color: var(--badge-color);
    font-size: 0.85rem;
    font-weight: 700;
    border: 1px solid color-mix(in srgb, var(--badge-color), transparent 70%);
  }
  .attention-badge span { font-size: 1.1rem; }

  .subtitle { color: var(--text-muted); font-size: 0.9rem; margin-top: 0.5rem; line-height: 1.5; }
  .subtitle code { background: var(--surface-variant); padding: 2px 6px; border-radius: 4px; font-family: 'JetBrains Mono', monospace; font-weight: 700; color: var(--primary-color); }

  .pull-tabs-bar { display: flex; gap: 0.5rem; border-bottom: 1px solid var(--md-sys-color-outline-variant); margin-bottom: 1rem; }
  .tab-pill { background: transparent; border: none; padding: 0.85rem 1.25rem; color: var(--text-muted); font-weight: 700; cursor: pointer; display: flex; align-items: center; gap: 0.75rem; transition: all 0.2s; position: relative; }
  .tab-pill:hover { color: var(--text-color); background: var(--surface-variant); }
  .tab-pill.active { color: var(--primary-color); }
  .tab-pill.active::after { content: ''; position: absolute; bottom: -1px; left: 0; right: 0; height: 2px; background: var(--primary-color); }
  .tab-pill .count { font-size: 0.7rem; background: var(--surface-2); padding: 2px 6px; border-radius: 10px; opacity: 0.7; }

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
