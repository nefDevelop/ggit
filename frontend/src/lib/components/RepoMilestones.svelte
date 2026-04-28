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
      state: ms.state
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
</style>
