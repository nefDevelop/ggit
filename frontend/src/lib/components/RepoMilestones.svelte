<script lang="ts">
  import { PlusCircle, Calendar, TrendingUp, TrendingDown, Edit, Trash2, AlertCircle } from "lucide-svelte";
  import { apiFetch } from "$lib/api";
  import { onMount, getContext } from 'svelte';
  import { dialog } from '$lib/dialog.svelte';
  import Button from './ui/Button.svelte';
  import Card from './ui/Card.svelte';
  import Badge from './ui/Badge.svelte';
  import Input from './ui/Input.svelte';
  import Select from './ui/Select.svelte';
  import PageHeader from './ui/PageHeader.svelte';

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
    state: 'open',
    hill_state: 0
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
        form = { title: '', description: '', due_on: '', state: 'open', hill_state: 0 };
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
    <Button variant="primary" onclick={() => { editingMilestone = null; form = { title: '', description: '', due_on: '', state: 'open', hill_state: 0 }; showEditor = true; }} icon={PlusCircle}>
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
        <div class="full-width hill-state-editor-box tonal">
          <div class="hill-editor-header">
            <span class="label-text">Hill Chart Progress: {form.hill_state}%</span>
            <span class="hill-phase-label">{form.hill_state < 50 ? 'Uphill: Figuring it out' : 'Downhill: Making it happen'}</span>
          </div>
          
          <div class="interactive-hill-container">
            <svg 
              viewBox="0 0 100 45" 
              preserveAspectRatio="none"
              class="hill-chart-edit"
              role="slider"
              aria-label="Hill Chart Progress"
              aria-valuemin="0"
              aria-valuemax="100"
              aria-valuenow={form.hill_state}
              tabindex="0"
              onclick={(e) => {
                const svg = e.currentTarget;
                const rect = svg.getBoundingClientRect();
                const x = ((e.clientX - rect.left) / rect.width) * 100;
                form.hill_state = Math.round(Math.max(0, Math.min(100, x)));
              }}
              onkeydown={(e) => {
                if (e.key === 'ArrowRight' || e.key === 'ArrowUp') {
                  form.hill_state = Math.min(100, form.hill_state + 1);
                } else if (e.key === 'ArrowLeft' || e.key === 'ArrowDown') {
                  form.hill_state = Math.max(0, form.hill_state - 1);
                }
              }}
              onmousedown={(e) => {
                 const svg = e.currentTarget;
                 const onMouseMove = (moveEvent: MouseEvent) => {
                    const rect = svg.getBoundingClientRect();
                    const x = ((moveEvent.clientX - rect.left) / rect.width) * 100;
                    form.hill_state = Math.round(Math.max(0, Math.min(100, x)));
                 };
                 const onMouseUp = () => {
                    window.removeEventListener('mousemove', onMouseMove);
                    window.removeEventListener('mouseup', onMouseUp);
                 };
                 window.addEventListener('mousemove', onMouseMove);
                 window.addEventListener('mouseup', onMouseUp);
              }}
            >
              <path d="M 0 40 Q 50 -30 100 40" fill="none" stroke="var(--md-sys-color-outline-variant)" stroke-width="1.5" />
              <line x1="50" y1="5" x2="50" y2="40" stroke="var(--md-sys-color-outline-variant)" stroke-dasharray="3" stroke-width="0.5" />
              
              <!-- Background Clickable Area -->
              <rect x="0" y="0" width="100" height="45" fill="transparent" style="cursor: crosshair" />

              <circle 
                cx={form.hill_state} 
                cy={40 - 140 * (form.hill_state / 100) + 140 * Math.pow(form.hill_state / 100, 2)} 
                r="3.5" 
                fill="var(--primary-color)" 
                stroke="var(--bg-color)"
                stroke-width="1"
                filter="drop-shadow(0 0 6px var(--primary-color))"
                style="cursor: grab"
              />
            </svg>
            <div class="hill-labels-wide">
              <span>Uncertainty</span>
              <span>Certainty</span>
            </div>
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

  <div class="milestones-content">
    {#if loading}
      <div class="loading-state">
        <div class="spinner"></div>
      </div>
    {:else if milestones.length === 0}
      <div class="empty-state">
        <AlertCircle size={48} />
        <h3>No milestones found</h3>
        <p>There are no milestones in this repository.</p>
      </div>
    {:else}
      <div class="flat-list-container">
        {#each milestones as ms}
          <div class="list-item-row">
            <div class="item-main-content">
              <div class="ms-header">
                <div class="ms-title-row">
                  <a href={`/${repo.data.owner.username}/${repo.data.name}/issues?milestone=${ms.id}`} class="item-title-link">
                    {ms.title}
                  </a>
                  <span class="badge-mini" style="background: {ms.state === 'open' ? 'var(--success)' : 'var(--text-muted)'}">{ms.state}</span>
                </div>
                <div class="ms-meta">
                  {#if ms.due_on}
                    <span class="due-date">
                      <Calendar size={14} />
                      Due on {new Date(ms.due_on).toLocaleDateString()}
                    </span>
                  {/if}
                  <span class="updated">Last updated {new Date().toLocaleDateString()}</span>
                </div>
              </div>
              
              <div class="ms-description">
                {ms.description || 'No description provided.'}
              </div>

              <div class="hill-chart-display-box">
                <svg viewBox="0 0 100 45" preserveAspectRatio="none" class="hill-chart-large">
                  <path d="M 0 40 Q 50 -30 100 40" fill="none" stroke="var(--md-sys-color-outline-variant)" stroke-width="1.5" />
                  <line x1="50" y1="5" x2="50" y2="40" stroke="var(--md-sys-color-outline-variant)" stroke-dasharray="2" stroke-width="0.4" />
                  <circle 
                    cx={ms.hill_state} 
                    cy={40 - 140 * (ms.hill_state / 100) + 140 * Math.pow(ms.hill_state / 100, 2)} 
                    r="3" 
                    fill="var(--primary-color)" 
                    filter="drop-shadow(0 0 5px var(--primary-color))"
                  />
                </svg>
                <div class="hill-labels">
                  <div class="phase uphill" class:active={ms.hill_state < 50}>
                      <TrendingUp size={16} />
                      Figuring it out
                  </div>
                  <div class="phase downhill" class:active={ms.hill_state >= 50}>
                      <TrendingDown size={16} />
                      Making it happen
                  </div>
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

            <div class="item-actions-right">
              <button class="btn-icon-minimal" onclick={() => handleEdit(ms)} title="Edit"><Edit size={16} /></button>
              <button class="btn-icon-minimal delete" onclick={() => deleteMilestone(ms.id)} title="Delete"><Trash2 size={16} /></button>
            </div>
          </div>
        {/each}
      </div>
    {/if}
  </div>
</div>

<style>
  .repo-milestones { display: flex; flex-direction: column; gap: 0; }
  .milestones-list { display: flex; flex-direction: column; gap: 1.5rem; padding: 0; margin-top: 1.5rem; }
  .section-actions { display: flex; justify-content: flex-end; padding: 0.5rem 0; border-bottom: 1px solid var(--md-sys-color-outline-variant); }

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

  .ms-description { font-size: 0.95rem; opacity: 0.8; line-height: 1.5; }

  .ms-progress-zone { display: flex; flex-direction: column; gap: 0.5rem; }
  .progress-info { display: flex; justify-content: space-between; font-size: 0.8rem; font-weight: 700; color: var(--text-muted); }
  .progress-bar { height: 8px; background: rgba(255,255,255,0.05); border-radius: 10px; overflow: hidden; }
  .progress-fill { height: 100%; background: var(--primary-color); border-radius: 10px; transition: width 0.4s ease; }

  .ms-actions { display: flex; flex-direction: column; gap: 0.5rem; justify-content: center; }
  :global(.danger-btn:hover) { color: var(--error) !important; }

  .empty-state { text-align: center; padding: 4rem; color: var(--text-muted); display: flex; flex-direction: column; align-items: center; gap: 1rem; opacity: 0.6; }

  .hill-state-editor-box { padding: 1.5rem; border-radius: 1.25rem; display: flex; flex-direction: column; gap: 1rem; border: 1px solid var(--md-sys-color-outline-variant); }
  .hill-editor-header { display: flex; justify-content: space-between; align-items: baseline; }
  .label-text { font-size: 0.9rem; font-weight: 800; color: var(--text-color); }
  .hill-phase-label { font-size: 0.75rem; font-weight: 700; color: var(--primary-color); text-transform: uppercase; letter-spacing: 0.05em; }
  
  .interactive-hill-container { padding: 0.5rem 1rem; }
  .hill-chart-edit { width: 100%; height: 100px; overflow: visible; }
  .hill-labels-wide { display: flex; justify-content: space-between; font-size: 0.7rem; font-weight: 800; text-transform: uppercase; color: var(--text-muted); opacity: 0.6; margin-top: 0.25rem; }

  .hill-chart-display-box { margin-top: 0.5rem; padding: 1.5rem; border-radius: 1.25rem; border: 1px solid var(--md-sys-color-outline-variant); background: rgba(0,0,0,0.2) !important; }
  .hill-chart-large { width: 100%; height: 80px; overflow: visible; }
  .hill-labels { display: flex; justify-content: space-between; margin-top: 0.75rem; }
  .phase { display: flex; align-items: center; gap: 0.5rem; font-size: 0.75rem; font-weight: 800; text-transform: uppercase; color: var(--text-muted); opacity: 0.4; transition: all 0.3s; }
  .phase.active { opacity: 1; color: var(--primary-color); }
</style>
le>
