<script lang="ts">
  import { onMount, getContext } from 'svelte';
  import DiffViewer from '$lib/components/DiffViewer.svelte';
  
  let { index } = $props();
  const repo: any = getContext('repo');
  
  let pull = $state<any>(null);
  let diff = $state<any>(null);
  let loading = $state(true);
  let error = $state('');
  
  let activeView = $state('conversation'); // 'conversation', 'commits', 'files'

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
        fetch(`/api/v2/${repo.data.owner.username}/${repo.data.name}/pulls/${index}`),
        fetch(`/api/v2/${repo.data.owner.username}/${repo.data.name}/pulls/${index}/diff`)
      ]);

      if (pullRes.ok) pull = await pullRes.json();
      if (diffRes.ok) diff = await diffRes.json();
      
      if (!pull) error = 'Pull request not found.';
    } catch (e) {
      error = 'Failed to load pull request details.';
    } finally {
      loading = false;
    }
  }

  onMount(() => {
    if (repo.data) fetchPullDetail();
  });
</script>

<div class="pull-detail-view">
  {#if loading}
    <div class="skeleton-detail island"></div>
  {:else if error}
    <div class="error-msg island">{error}</div>
  {:else if pull}
    <header class="pull-header">
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
         <strong>{pull.user?.username}</strong> wants to merge commits into <code>{pull.pull_request?.base_branch}</code> from <code>{pull.pull_request?.head_branch}</code>
       </div>
    </header>

    <nav class="pull-tabs">
      <button class="pull-tab" class:active={activeView === 'conversation'} onclick={() => activeView = 'conversation'}>
        Conversation
      </button>
      <button class="pull-tab" class:active={activeView === 'files'} onclick={() => activeView = 'files'}>
        Files Changed <span class="count">{diff?.Files?.length || 0}</span>
      </button>
    </nav>

    <div class="pull-body">
      {#if activeView === 'conversation'}
        <div class="conversation-view">
           <div class="timeline-item island">
              <div class="item-header"><strong>{pull.user?.username}</strong> commented</div>
              <div class="item-body markdown-body">
                {pull.body}
              </div>
           </div>
           <p class="muted center">Timeline and comments coming soon...</p>
        </div>
      {:else if activeView === 'files'}
        <DiffViewer {diff} />
      {/if}
    </div>
  {/if}
</div>

<style>
  .pull-detail-view { display: flex; flex-direction: column; gap: 1.5rem; }
  .pull-header { border-bottom: 1px solid var(--border-color); padding-bottom: 1.5rem; }
  .title-row { display: flex; align-items: center; gap: 1rem; margin-bottom: 0.5rem; }
  .title-row h1 { margin: 0; font-size: 1.75rem; font-weight: 800; }
  .index { color: var(--text-muted); font-weight: 400; }
  
  .status-pill { display: flex; align-items: center; gap: 0.5rem; padding: 0.4rem 1rem; border-radius: 2rem; background: #10b981; color: white; font-weight: 700; }
  .status-pill.merged { background: #8b5cf6; }
  .subtitle { color: var(--text-muted); font-size: 0.9rem; }
  .subtitle code { background: rgba(255,255,255,0.05); padding: 0.1rem 0.3rem; border-radius: 4px; }

  .pull-tabs { display: flex; gap: 1rem; border-bottom: 1px solid var(--border-color); margin-bottom: 1rem; }
  .pull-tab { background: transparent; border: none; padding: 0.75rem 1rem; color: var(--text-muted); font-weight: 600; cursor: pointer; position: relative; }
  .pull-tab.active { color: var(--text-color); }
  .pull-tab.active::after { content: ''; position: absolute; bottom: 0; left: 0; right: 0; height: 2px; background: var(--primary-color); }
  .pull-tab .count { background: rgba(255,255,255,0.1); padding: 0.1rem 0.4rem; border-radius: 10px; font-size: 0.75rem; margin-left: 0.5rem; }

  .diff-file { margin-bottom: 2rem; overflow: hidden; }
  .diff-file-header { padding: 0.75rem 1rem; background: rgba(255,255,255,0.03); border-bottom: 1px solid var(--border-color); display: flex; align-items: center; gap: 0.5rem; }
  
  .diff-line { display: grid; grid-template-columns: 40px 40px 1fr; font-family: monospace; font-size: 0.85rem; line-height: 1.4; }
  .diff-line.add { background: rgba(16, 185, 129, 0.1); }
  .diff-line.del { background: rgba(239, 68, 68, 0.1); }
  .diff-line pre { margin: 0; padding: 0 0.5rem; }
  .ln { text-align: right; padding-right: 0.5rem; opacity: 0.3; border-right: 1px solid rgba(255,255,255,0.05); }

  .diff-section-header { padding: 0.25rem 1rem; background: rgba(255,255,255,0.01); color: var(--text-muted); font-size: 0.75rem; font-family: monospace; }
  
  .muted { color: var(--text-muted); opacity: 0.5; }
  .center { text-align: center; padding: 2rem; }

  .skeleton-detail { height: 400px; animation: pulse 1.5s infinite; }
  @keyframes pulse { 0%, 100% { opacity: 0.5; } 50% { opacity: 0.8; } }
</style>
