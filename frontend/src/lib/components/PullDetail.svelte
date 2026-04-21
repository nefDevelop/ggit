<script lang="ts">
  import { apiFetch } from "$lib/api";
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
        apiFetch(`/api/v2/${repo.data.owner.username}/${repo.data.name}/pulls/${index}`),
        apiFetch(`/api/v2/${repo.data.owner.username}/${repo.data.name}/pulls/${index}/diff`)
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
         <strong>{pull.user?.username}</strong> wants to merge into <code>{pull.pull_request?.base_branch}</code> from <code>{pull.pull_request?.head_branch}</code>
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
              <div class="item-header"><strong>{pull.user?.username}</strong> commented</div>
              <div class="item-body markdown-body">
                {pull.body}
              </div>
           </div>
           <div class="coming-soon">
             <span class="material-symbols-rounded">construction</span>
             <p>Timeline details are being modernized...</p>
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
  
  .pull-detail-header { border-bottom: 1px solid var(--border-color); padding-bottom: 1.5rem; }
  .title-row { display: flex; align-items: center; gap: 1.5rem; }
  .title-row h1 { margin: 0; font-size: 1.75rem; font-weight: 800; }
  .index { color: var(--text-muted); font-weight: 400; }
  
  .status-pill { display: flex; align-items: center; gap: 0.5rem; padding: 0.4rem 1.25rem; border-radius: 2rem; background: var(--success); color: var(--on-primary); font-weight: 700; font-size: 0.9rem; }
  .status-pill.merged { background: var(--primary); }
  
  .subtitle { color: var(--text-muted); font-size: 0.9rem; margin-top: 0.5rem; }
  .subtitle code { background: var(--surface-variant); padding: 2px 6px; border-radius: 4px; font-family: 'JetBrains Mono', monospace; font-weight: 700; color: var(--primary-color); }

  .pull-tabs-bar { display: flex; gap: 0.5rem; border-bottom: 1px solid var(--border-color); margin-bottom: 1rem; }
  .nav-tab { background: transparent; border: none; padding: 0.85rem 1.25rem; color: var(--text-muted); font-weight: 700; cursor: pointer; display: flex; align-items: center; gap: 0.75rem; transition: all 0.2s; position: relative; }
  .nav-tab:hover { color: var(--text-color); background: var(--surface-variant); }
  .nav-tab.active { color: var(--primary-color); }
  .nav-tab.active::after { content: ''; position: absolute; bottom: -1px; left: 0; right: 0; height: 2px; background: var(--primary-color); }
  .nav-tab .count { font-size: 0.7rem; background: var(--surface-2); padding: 2px 6px; border-radius: 10px; opacity: 0.7; }

  .timeline-item-flat { 
    border: 1px solid var(--border-color); 
    background: var(--bg-color); 
    border-radius: 1rem; 
    overflow: hidden; 
    margin-bottom: 2rem;
  }
  .item-header { padding: 0.75rem 1.5rem; background: var(--surface-variant); border-bottom: 1px solid var(--border-color); font-size: 0.85rem; color: var(--text-muted); }
  .item-body { padding: 1.5rem; line-height: 1.6; }

  .coming-soon { text-align: center; padding: 4rem; color: var(--text-muted); opacity: 0.5; display: flex; flex-direction: column; align-items: center; gap: 1rem; }
  .coming-soon span { font-size: 3rem; }

  .skeleton-detail { height: 400px; background: var(--surface-variant); border-radius: 1rem; animation: pulse 1.5s infinite; }
  @keyframes pulse { 0%, 100% { opacity: 0.3; } 50% { opacity: 0.6; } }

  .error-msg { background: rgba(239, 68, 68, 0.1); color: #ef4444; padding: 2rem; border-radius: 1rem; text-align: center; font-weight: 700; }
</style>
