<script lang="ts">
  import { apiFetch } from "$lib/api";
  import { onMount, getContext } from 'svelte';
  import { dialog } from '$lib/dialog.svelte';
  import Button from './ui/Button.svelte';
  import Card from './ui/Card.svelte';
  import Badge from './ui/Badge.svelte';
  import Input from './ui/Input.svelte';
  import Select from './ui/Select.svelte';

  const repo: any = getContext('repo');
  let milestones = $state<any[]>([]);
  let loading = $state(true);
  let error = $state('');

  let showEditor = $state(false);
  let editingMilestone = $state<any>(null);
  let form = $state({
    title: '',
    description: '',
    due_on: '',
    state: 'open'
  });

  const stateOptions = [
    { value: 'open', label: 'Open' },
    { value: 'closed', label: 'Closed' }
  ];

  async function fetchMilestones() {
    loading = true;
    try {
      const res = await apiFetch(`/api/v2/${repo.data.owner.username}/${repo.data.name}/milestones`);
      if (res.ok) {
        milestones = await res.json();
      }
    } catch (e) {
      error = 'Failed to load milestones.';
    } finally {
      loading = false;
    }
  }

  async function saveMilestone() {
    if (!form.title.trim()) return;
    
    const body: any = { ...form };
    if (form.due_on) {
        body.due_on = new Date(form.due_on).toISOString();
    } else {
        delete body.due_on;
    }

    const method = editingMilestone ? 'PATCH' : 'POST';
    const url = editingMilestone 
      ? `/api/v2/${repo.data.owner.username}/${repo.data.name}/milestones/${editingMilestone.id}`
      : `/api/v2/${repo.data.owner.username}/${repo.data.name}/milestones`;

    try {
      const res = await apiFetch(url, {
        method,
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(body)
      });
      if (res.ok) {
        showEditor = false;
        editingMilestone = null;
        form = { title: '', description: '', due_on: '', state: 'open' };
        fetchMilestones();
      }
    } catch (e) {
      console.error(e);
    }
  }

  async function deleteMilestone(id: number) {
    if (!await dialog.confirm('Are you sure you want to delete this milestone?')) return;
    try {
      const res = await apiFetch(`/api/v2/${repo.data.owner.username}/${repo.data.name}/milestones/${id}`, {
        method: 'DELETE'
      });
      if (res.ok) fetchMilestones();
    } catch (e) {
      console.error(e);
    }
  }

  function handleEdit(ms: any) {
    editingMilestone = ms;
    form = {
      title: ms.title,
      description: ms.description || '',
      due_on: ms.due_on ? ms.due_on.split('T')[0] : '',
      state: ms.state,
      hill_state: ms.hill_state || 0
    };
    showEditor = true;
  }

  onMount(fetchMilestones);

  function getProgress(ms: any) {
    const total = ms.open_issues + ms.closed_issues;
    if (total === 0) return 0;
    return (ms.closed_issues / total) * 100;
  }
</script>

<div class="repo-milestones">
  <header class="section-actions">
    <Button variant="primary" onclick={() => { editingMilestone = null; form = { title: '', description: '', due_on: '', state: 'open' }; showEditor = true; }} icon="add_circle">
      New Milestone
    </Button>
  </header>

  {#if showEditor}
    <Card class="milestone-editor tonal" padding="lg">
      <header class="editor-header">
        <h3 class="editor-title">{editingMilestone ? 'Edit Milestone' : 'Create New Milestone'}</h3>
        {#if editingMilestone}
          <div class="state-picker">
            <Select 
              label="State" 
              bind:value={form.state} 
              options={stateOptions} 
            />
          </div>
        {/if}
      </header>
      
      <div class="editor-form">
        <Input 
          label="Title" 
          bind:value={form.title} 
          placeholder="e.g. Version 1.0.0" 
        />
        <Input 
          label="Due Date (Optional)" 
          type="date" 
          bind:value={form.due_on} 
        />
        <div class="full-width">
          <Input 
            label="Description" 
            multiline 
            rows={4} 
            bind:value={form.description} 
            placeholder="What should be completed in this milestone?" 
          />
        </div>
        <div class="full-width hill-state-picker">
          <label for="hill-state">Hill Chart State: {form.hill_state}% ({form.hill_state < 50 ? 'Uphill / Figuring it out' : 'Downhill / Making it happen'})</label>
          <input 
            id="hill-state"
            type="range" 
            min="0" 
            max="100" 
            bind:value={form.hill_state} 
            class="hill-slider"
          />
          <div class="hill-markers">
            <span>Uncertainty</span>
            <span>Execution</span>
          </div>
        </div>
      </div>

      <div class="editor-footer">
        <Button variant="secondary" onclick={() => showEditor = false}>Cancel</Button>
        <Button variant="primary" onclick={saveMilestone} disabled={!form.title.trim()}>
          {editingMilestone ? 'Update Milestone' : 'Create Milestone'}
        </Button>
      </div>
    </Card>
  {/if}

  <div class="milestones-list">
    {#if loading}
      <div class="loading-state">Loading milestones...</div>
    {:else if milestones.length === 0}
      <div class="empty-state tonal">
        <span class="material-symbols-rounded">event_busy</span>
        <p>No milestones found for this repository.</p>
      </div>
    {:else}
      {#each milestones as ms}
        <Card class="milestone-card tonal" padding="lg">
          <div class="ms-main">
            <div class="ms-header">
              <div class="ms-title-row">
                <h4>{ms.title}</h4>
                <Badge variant={ms.state === 'open' ? 'success' : 'neutral'}>{ms.state}</Badge>
              </div>
              <div class="ms-meta">
                {#if ms.due_on}
                  <span class="due-date">
                    <span class="material-symbols-rounded mini">calendar_today</span>
                    Due on {new Date(ms.due_on).toLocaleDateString()}
                  </span>
                {/if}
                <span class="updated">Last updated {new Date().toLocaleDateString()}</span>
              </div>
            </div>
            
            <div class="ms-description">
              {ms.description || 'No description provided.'}
            </div>

            <div class="hill-chart-container">
              <svg viewBox="0 0 100 40" class="hill-chart">
                <path d="M 0 35 Q 50 -10 100 35" fill="none" stroke="var(--md-sys-color-outline-variant)" stroke-width="1" />
                <line x1="50" y1="5" x2="50" y2="35" stroke="var(--md-sys-color-outline-variant)" stroke-dasharray="2" stroke-width="0.5" />
                <circle 
                  cx={ms.hill_state} 
                  cy={35 - (Math.sin((ms.hill_state / 100) * Math.PI) * 30)} 
                  r="3" 
                  fill="var(--primary-color)" 
                  filter="drop-shadow(0 0 4px var(--primary-color))"
                />
              </svg>
              <div class="hill-labels">
                <span>Figuring it out</span>
                <span>Making it happen</span>
              </div>
            </div>

            <div class="ms-progress-zone">
              <div class="progress-info">
                <span>{getProgress(ms).toFixed(0)}% Complete</span>
                <span class="counts">{ms.closed_issues} closed · {ms.open_issues} open</span>
              </div>
              <div class="progress-bar">
                <div class="progress-fill" style="width: {getProgress(ms)}%"></div>
              </div>
            </div>
          </div>

          <div class="ms-actions">
            <Button variant="ghost" icon="edit" onclick={() => handleEdit(ms)} />
            <Button variant="ghost" icon="delete" class="danger-btn" onclick={() => deleteMilestone(ms.id)} />
          </div>
        </Card>
      {/each}
    {/if}
  </div>
</div>

<style>
  .repo-milestones { display: flex; flex-direction: column; gap: 2rem; }
  .section-actions { display: flex; justify-content: flex-end; }

  .editor-title { margin: 0; font-weight: 800; color: var(--primary-color); }
  .editor-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 2rem; gap: 2rem; }
  .state-picker { width: 180px; }

  .editor-form { display: grid; grid-template-columns: 1fr 1fr; gap: 1.5rem; margin-bottom: 2rem; }
  @media (max-width: 600px) { .editor-form { grid-template-columns: 1fr; } .editor-header { flex-direction: column; align-items: flex-start; } .state-picker { width: 100%; } }
  
  .full-width { grid-column: span 2; }
  @media (max-width: 600px) { .full-width { grid-column: span 1; } }

  .editor-footer { display: flex; justify-content: flex-end; gap: 1rem; }

  .milestones-list { display: flex; flex-direction: column; gap: 1.5rem; }
  
  :global(.milestone-card) { display: flex; gap: 2rem; border-radius: 1.5rem; }
  .ms-main { flex: 1; display: flex; flex-direction: column; gap: 1.25rem; }
  
  .ms-header { display: flex; flex-direction: column; gap: 0.25rem; }
  .ms-title-row { display: flex; align-items: center; gap: 1rem; }
  .ms-title-row h4 { margin: 0; font-size: 1.25rem; font-weight: 800; color: var(--primary-color); }
  
  .ms-meta { display: flex; gap: 1.5rem; font-size: 0.8rem; color: var(--text-muted); font-weight: 600; }
  .ms-meta .mini { font-size: 1rem; }

  .ms-description { font-size: 0.95rem; opacity: 0.8; line-height: 1.5; }

  .ms-progress-zone { display: flex; flex-direction: column; gap: 0.5rem; }
  .progress-info { display: flex; justify-content: space-between; font-size: 0.8rem; font-weight: 700; color: var(--text-muted); }
  .progress-bar { height: 8px; background: rgba(255,255,255,0.05); border-radius: 10px; overflow: hidden; }
  .progress-fill { height: 100%; background: var(--primary-color); border-radius: 10px; transition: width 0.4s ease; }

  .ms-actions { display: flex; flex-direction: column; gap: 0.5rem; justify-content: center; }
  :global(.danger-btn:hover) { color: var(--error) !important; }

  .empty-state { text-align: center; padding: 4rem; color: var(--text-muted); display: flex; flex-direction: column; align-items: center; gap: 1rem; opacity: 0.6; }
  .empty-state span { font-size: 3rem; }

  .hill-state-picker { display: flex; flex-direction: column; gap: 0.75rem; margin-top: 1rem; }
  .hill-state-picker label { font-size: 0.85rem; font-weight: 700; color: var(--text-muted); }
  .hill-slider { width: 100%; height: 6px; background: rgba(255,255,255,0.1); border-radius: 3px; outline: none; appearance: none; }
  .hill-slider::-webkit-slider-thumb { appearance: none; width: 16px; height: 16px; background: var(--primary-color); border-radius: 50%; cursor: pointer; border: 2px solid var(--bg-color); }
  .hill-markers { display: flex; justify-content: space-between; font-size: 0.7rem; font-weight: 800; text-transform: uppercase; color: var(--text-muted); opacity: 0.6; }

  .hill-chart-container { margin: 0.5rem 0; background: rgba(0,0,0,0.2); padding: 1.25rem; border-radius: 1rem; border: 1px solid var(--md-sys-color-outline-variant); }
  .hill-chart { width: 100%; height: 60px; overflow: visible; }
  .hill-labels { display: flex; justify-content: space-between; font-size: 0.7rem; font-weight: 700; text-transform: uppercase; color: var(--text-muted); margin-top: 0.5rem; opacity: 0.8; }
</style>
