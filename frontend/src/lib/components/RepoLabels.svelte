<script lang="ts">
  import { apiFetch } from "$lib/api";
  import { onMount, getContext } from 'svelte';
  import { dialog } from '$lib/dialog.svelte';
  import Button from './ui/Button.svelte';
  import Card from './ui/Card.svelte';
  import Input from './ui/Input.svelte';

  const repo: any = getContext('repo');
  let labels = $state<any[]>([]);
  let loading = $state(true);
  let error = $state('');

  let showEditor = $state(false);
  let editingLabel = $state<any>(null);
  let newName = $state('');
  let newColor = $state('#444444');

  async function fetchLabels() {
    loading = true;
    try {
      const res = await apiFetch(`/api/v2/${repo.data.owner.username}/${repo.data.name}/labels`);
      if (res.ok) {
        labels = await res.json();
      }
    } catch (e) {
      error = 'Failed to load labels.';
    } finally {
      loading = false;
    }
  }

  async function saveLabel() {
    if (!newName.trim()) return;
    const body = { name: newName, color: newColor };
    const method = editingLabel ? 'PATCH' : 'POST';
    const url = editingLabel 
      ? `/api/v2/${repo.data.owner.username}/${repo.data.name}/labels/${editingLabel.id}`
      : `/api/v2/${repo.data.owner.username}/${repo.data.name}/labels`;

    try {
      const res = await apiFetch(url, {
        method,
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(body)
      });
      if (res.ok) {
        showEditor = false;
        editingLabel = null;
        newName = '';
        newColor = '#444444';
        fetchLabels();
      }
    } catch (e) {
      console.error(e);
    }
  }

  async function deleteLabel(id: number) {
    if (!await dialog.confirm('Are you sure you want to delete this label?')) return;
    try {
      const res = await apiFetch(`/api/v2/${repo.data.owner.username}/${repo.data.name}/labels/${id}`, {
        method: 'DELETE'
      });
      if (res.ok) fetchLabels();
    } catch (e) {
      console.error(e);
    }
  }

  function handleEdit(label: any) {
    editingLabel = label;
    newName = label.name;
    newColor = label.color;
    showEditor = true;
  }

  onMount(fetchLabels);
</script>

<div class="repo-labels">
  <header class="section-actions">
    <Button variant="primary" onclick={() => { editingLabel = null; newName = ''; newColor = '#444444'; showEditor = true; }} icon="add_circle">
      New Label
    </Button>
  </header>

  <div class="labels-content">

  {#if showEditor}
    <Card class="label-editor tonal" padding="lg">
      <h3 class="editor-title">{editingLabel ? 'Edit Label' : 'Create New Label'}</h3>
      <div class="editor-grid">
        <Input 
          label="Label Name" 
          bind:value={newName} 
          placeholder="e.g. bug, improvement..." 
        />
        
        <div class="color-field">
          <label class="color-label" for="native-color-picker">Color</label>
          <div class="color-picker-row">
            <input type="color" bind:value={newColor} class="native-picker" id="native-color-picker" />
            <Input bind:value={newColor} placeholder="#444444" />
          </div>
        </div>

        <div class="preview-zone">
          <label class="preview-label" for="label-preview-box">Preview</label>
          <div class="label-preview" id="label-preview-box" style="background: {newColor}">
            {newName || 'Label Preview'}
          </div>
        </div>
      </div>
      <div class="editor-footer">
        <Button variant="secondary" onclick={() => showEditor = false}>Cancel</Button>
        <Button variant="primary" onclick={saveLabel} disabled={!newName.trim()}>
          {editingLabel ? 'Update Label' : 'Create Label'}
        </Button>
      </div>
    </Card>
  {/if}

  <div class="labels-list">
    {#if loading}
      <div class="loading-state">Loading labels...</div>
    {:else if labels.length === 0}
      <div class="empty-state">
        <span class="material-symbols-rounded">error_outline</span>
        <h3>No labels found</h3>
        <p>There are no labels defined yet for this repository.</p>
      </div>
    {:else}
      {#each labels as label}
        <Card class="label-row" tonal={true} padding="none">
          <div class="label-info">
            <span class="label-tag" style="background: {label.color}">
              {label.name}
            </span>
          </div>
          <div class="label-actions">
            <Button variant="ghost" icon="edit" onclick={() => handleEdit(label)} />
            <Button variant="ghost" icon="delete" class="danger-btn" onclick={() => deleteLabel(label.id)} />
          </div>
        </Card>
      {/each}
    {/if}
    </div>
  </div>
</div>

<style>
  .repo-labels { display: flex; flex-direction: column; gap: 0; margin: -2rem; }
  .labels-content { padding: 2rem; display: flex; flex-direction: column; gap: 2rem; }
  .section-actions { display: flex; justify-content: flex-end; padding: 1.5rem 2rem; border-bottom: 1px solid var(--md-sys-color-outline-variant); }

  .editor-title { margin-top: 0; margin-bottom: 1.5rem; font-weight: 800; color: var(--primary-color); }
  .editor-grid { display: grid; grid-template-columns: 2fr 1.5fr 1fr; gap: 2rem; align-items: flex-end; margin-bottom: 2rem; }
  @media (max-width: 800px) { .editor-grid { grid-template-columns: 1fr; gap: 1.5rem; } }
  
  .color-field { display: flex; flex-direction: column; gap: 0.5rem; }
  .color-label, .preview-label { font-size: 0.9rem; font-weight: 600; color: var(--md-sys-color-on-surface-variant); margin-left: 0.25rem; }

  .color-picker-row { display: flex; gap: 0.5rem; align-items: flex-end; }
  .native-picker { width: 48px; height: 48px; padding: 0; border: 1px solid var(--border-color); background: transparent; cursor: pointer; border-radius: 8px; flex-shrink: 0; }

  .label-preview { 
    height: 48px; 
    display: flex; 
    align-items: center; 
    justify-content: center;
    padding: 0 1rem; 
    border-radius: 8px; 
    color: white; 
    font-weight: 800; 
    text-transform: uppercase; 
    font-size: 0.75rem; 
    letter-spacing: 0.05em;
    border: 1px solid rgba(255,255,255,0.2);
  }

  .editor-footer { display: flex; justify-content: flex-end; gap: 1rem; }

  .labels-list { display: flex; flex-direction: column; gap: 0.75rem; }
  
  :global(.label-row) { 
    display: flex; 
    align-items: center; 
    justify-content: space-between; 
    padding: 1rem 1.5rem !important; 
    border-radius: 1rem;
  }

  .label-tag { 
    font-size: 0.75rem; 
    font-weight: 800; 
    padding: 0.25rem 0.75rem; 
    border-radius: 4px; 
    color: white; 
    text-transform: uppercase;
    letter-spacing: 0.05em;
    border: 1px solid rgba(255,255,255,0.1);
  }

  .label-actions { display: flex; gap: 0.5rem; }
  :global(.danger-btn:hover) { color: var(--error) !important; }

  .empty-state { text-align: center; padding: 4rem; color: var(--text-muted); display: flex; flex-direction: column; align-items: center; gap: 1rem; opacity: 0.6; }
  .empty-state span { font-size: 3rem; }
</style>
