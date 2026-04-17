<script lang="ts">
  import { apiFetch } from "$lib/api";
  import { onMount } from 'svelte';
  import { page } from '$app/state';
  import { goto } from '$app/navigation';
  import DiffViewer from '$lib/components/DiffViewer.svelte';
  import CommitRow from '$lib/components/CommitRow.svelte';

  let username = $derived(page.params.username);
  let reponame = $derived(page.params.reponame);

  let branches = $state<any[]>([]);
  let baseBranch = $state('');
  let headBranch = $state('');
  let comparison = $state<any>(null);
  let loadingComp = $state(false);
  
  let title = $state('');
  let content = $state('');
  let submitting = $state(false);

  let activeTab = $state('diff'); // 'diff' or 'commits'

  async function fetchBranches() {
    try {
      const res = await apiFetch(`/api/v2/${username}/${reponame}/branches`);
      if (res.ok) {
          const data = await res.json();
          branches = [data.default, ...data.active, ...data.stale];
          baseBranch = data.default.name;
          // Try to find a branch that isn't the default to suggest as head
          const other = data.active.find((b: any) => b.name !== baseBranch);
          headBranch = other ? other.name : baseBranch;
      }
    } catch (e) {
      console.error(e);
    }
  }

  async function fetchComparison() {
    if (!baseBranch || !headBranch || baseBranch === headBranch) {
      comparison = null;
      return;
    }
    loadingComp = true;
    try {
      const res = await apiFetch(`/api/v2/${username}/${reponame}/compare/${baseBranch}...${headBranch}`);
      if (res.ok) {
        comparison = await res.json();
        // Auto-fill title if empty
        if (!title && comparison.num_commits === 1) {
          title = comparison.commits[0].summary;
        }
      } else {
        comparison = null;
      }
    } catch (e) {
      console.error(e);
      comparison = null;
    } finally {
      loadingComp = false;
    }
  }

  async function createPR() {
    if (!title) return;
    submitting = true;
    try {
      const res = await apiFetch(`/api/v2/${username}/${reponame}/pulls`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          title,
          content,
          base_branch: baseBranch,
          head_branch: headBranch
        })
      });
      if (res.ok) {
        const pr = await res.json();
        goto(`/${username}/${reponame}/pulls/${pr.index}`);
      }
    } catch (e) {
      console.error(e);
    } finally {
      submitting = false;
    }
  }

  $effect(() => {
    if (baseBranch && headBranch) fetchComparison();
  });

  onMount(fetchBranches);
</script>

<div class="pr-creator">
  <header class="creator-header island tonal">
    <div class="compare-selectors">
      <div class="selector">
        <span class="label">Base</span>
        <select bind:value={baseBranch}>
          {#each branches as b}
            <option value={b.name}>{b.name}</option>
          {/each}
        </select>
      </div>
      <span class="material-symbols-rounded arrow">arrow_back</span>
      <div class="selector">
        <span class="label">Compare</span>
        <select bind:value={headBranch}>
          {#each branches as b}
            <option value={b.name}>{b.name}</option>
          {/each}
        </select>
      </div>
    </div>

    {#if baseBranch === headBranch}
      <div class="info-banner warning">
        <span class="material-symbols-rounded">info</span>
        Select different branches to compare.
      </div>
    {:else if comparison}
      <div class="info-banner success">
        <span class="material-symbols-rounded">check_circle</span>
        Able to merge. These branches can be automatically merged.
      </div>
    {/if}
  </header>

  {#if comparison}
    <div class="workspace-grid">
      <aside class="creation-form island">
        <h3>Create Pull Request</h3>
        <div class="field">
          <input type="text" bind:value={title} placeholder="Title" required />
        </div>
        <div class="field">
          <textarea bind:value={content} placeholder="Description (Markdown supported)"></textarea>
        </div>
        <button class="btn btn-primary lg" onclick={createPR} disabled={!title || submitting}>
          {submitting ? 'Creating...' : 'Create Pull Request'}
        </button>
      </aside>

      <main class="comparison-preview island tonal">
        <nav class="preview-tabs">
          <button class:active={activeTab === 'diff'} onclick={() => activeTab = 'diff'}>
            Changes <span class="badge">{comparison.diff?.Files?.length || 0}</span>
          </button>
          <button class:active={activeTab === 'commits'} onclick={() => activeTab = 'commits'}>
            Commits <span class="badge">{comparison.num_commits}</span>
          </button>
        </nav>

        <div class="preview-content">
          {#if loadingComp}
            <div class="loading-overlay">
              <div class="spinner"></div>
            </div>
          {:else}
            {#if activeTab === 'diff'}
               <DiffViewer diff={comparison.diff} />
            {:else}
               <div class="commit-list">
                 {#each comparison.commits as commit}
                   <CommitRow {commit} minimalist={true} />
                 {/each}
               </div>
            {/if}
          {/if}
        </div>
      </main>
    </div>
  {/if}
</div>

<style>
  .pr-creator { display: flex; flex-direction: column; gap: 2rem; animation: fadeIn 0.4s ease-out; }
  
  .creator-header { padding: 1.5rem 2rem; display: flex; align-items: center; gap: 3rem; }
  .compare-selectors { display: flex; align-items: center; gap: 1rem; }
  .selector { display: flex; flex-direction: column; gap: 0.25rem; }
  .selector .label { font-size: 0.7rem; font-weight: 800; text-transform: uppercase; color: var(--text-muted); opacity: 0.7; }
  .selector select { background: rgba(0,0,0,0.2); border: 1px solid var(--border-color); color: #fff; padding: 0.5rem 1rem; border-radius: 0.5rem; font-weight: 600; min-width: 150px; }
  .arrow { opacity: 0.5; margin-top: 1rem; }

  .info-banner { display: flex; align-items: center; gap: 0.75rem; padding: 0.75rem 1.5rem; border-radius: 0.75rem; font-size: 0.9rem; font-weight: 600; flex: 1; }
  .info-banner.warning { background: rgba(234, 179, 8, 0.1); color: #eab308; }
  .info-banner.success { background: rgba(34, 197, 94, 0.1); color: #22c55e; }

  .workspace-grid { display: grid; grid-template-columns: 400px 1fr; gap: 2rem; align-items: start; }
  
  .creation-form { padding: 2rem; display: flex; flex-direction: column; gap: 1.5rem; box-shadow: var(--shadow-lg); }
  .creation-form h3 { margin: 0; font-size: 1.25rem; }
  .field input { width: 100%; background: rgba(0,0,0,0.2); border: 1px solid var(--border-color); color: #fff; padding: 0.75rem 1rem; border-radius: 0.75rem; font-size: 1rem; }
  .field textarea { width: 100%; background: rgba(0,0,0,0.2); border: 1px solid var(--border-color); color: #fff; padding: 0.75rem 1rem; border-radius: 0.75rem; font-size: 0.95rem; height: 300px; resize: none; font-family: inherit; }
  .lg { padding: 1rem; font-weight: 700; width: 100%; }

  .comparison-preview { padding: 0; overflow: hidden; display: flex; flex-direction: column; min-height: 600px; }
  .preview-tabs { display: flex; background: rgba(0,0,0,0.1); border-bottom: 1px solid var(--border-color); }
  .preview-tabs button { flex: 1; padding: 1rem; border: none; background: transparent; color: var(--text-muted); cursor: pointer; font-weight: 700; display: flex; align-items: center; justify-content: center; gap: 0.5rem; transition: all 0.2s; }
  .preview-tabs button.active { color: var(--primary-color); background: rgba(var(--primary-rgb), 0.05); border-bottom: 2px solid var(--primary-color); }
  .badge { font-size: 0.75rem; background: rgba(255,255,255,0.05); padding: 0.1rem 0.5rem; border-radius: 1rem; }

  .preview-content { flex: 1; position: relative; }
  .loading-overlay { position: absolute; inset: 0; background: rgba(0,0,0,0.2); display: flex; align-items: center; justify-content: center; z-index: 10; backdrop-filter: blur(2px); }
  
  .commit-list { display: flex; flex-direction: column; }

  @keyframes fadeIn { from { opacity: 0; transform: translateY(10px); } to { opacity: 1; transform: translateY(0); } }
</style>
