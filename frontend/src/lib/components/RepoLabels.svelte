<script lang="ts">
  import { apiFetch } from "$lib/api";
  import { onMount, getContext } from 'svelte';
  import Button from './ui/Button.svelte';
  import Card from './ui/Card.svelte';

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
    if (!confirm('Are you sure you want to delete this label?')) return;
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

  {#if showEditor}
    <Card class="label-editor tonal" padding="lg">
      <h3>{editingLabel ? 'Edit Label' : 'Create New Label'}</h3>
      <div class="editor-grid">
        <div class="field">
          <label for="name">Label Name</label>
          <input id="name" type="text" bind:value={newName} placeholder="e.g. bug, improvement..." />
        </div>
        <div class="field">
          <label for="color">Color</label>
          <div class="color-picker-row">
            <input id="color" type="color" bind:value={newColor} />
            <input type="text" bind:value={newColor} placeholder="#444444" />
          </div>
        </div>
        <div class="preview-zone">
          <label>Preview</label>
          <div class="label-preview" style="background: {newColor}">
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
      <div class="empty-state tonal">
        <span class="material-symbols-rounded">label_off</span>
        <p>No labels defined yet for this repository.</p>
      </div>
    {:else}
      {#each labels as label}
        <Card class="label-row tonal" padding="none">
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

<style>
  .repo-labels { display: flex; flex-direction: column; gap: 2rem; }
  .section-actions { display: flex; justify-content: flex-end; }

  .label-editor { margin-bottom: 2rem; }
  .editor-grid { display: grid; grid-template-columns: 2fr 1fr 1fr; gap: 2rem; align-items: flex-end; margin-bottom: 2rem; }
  
  .field { display: flex; flex-direction: column; gap: 0.5rem; }
  .field label { font-size: 0.75rem; font-weight: 800; text-transform: uppercase; color: var(--text-muted); }
  .field input { background: var(--surface-variant); border: 1px solid var(--border-color); padding: 0.75rem; border-radius: 8px; color: var(--text-color); }

  .color-picker-row { display: flex; gap: 0.5rem; }
  .color-picker-row input[type="color"] { width: 48px; padding: 0; border: none; background: transparent; cursor: pointer; }

  .label-preview { 
    height: 38px; 
    display: flex; 
    align-items: center; 
    padding: 0 1rem; 
    border-radius: 4px; 
    color: white; 
    font-weight: 700; 
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
