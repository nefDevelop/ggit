<script lang="ts">
  import { apiFetch } from "$lib/api";
  import { onMount, getContext } from 'svelte';
  import { dialog } from '$lib/dialog.svelte';
  import Button from './ui/Button.svelte';
  import Card from './ui/Card.svelte';
  import Input from './ui/Input.svelte';
  import PageHeader from './ui/PageHeader.svelte';
  

import Plus from "lucide-svelte/icons/plus";
import Info from "lucide-svelte/icons/info";
import Edit from "lucide-svelte/icons/edit";
import Trash2 from "lucide-svelte/icons/trash-2";
import Tag from "lucide-svelte/icons/tag";



  const repo: any = getContext('repo');
  let labels = $state<any[]>([]);
  let loading = $state(true);
  let error = $state('');

  let showEditor = $state(false);
  let editingLabel = $state<any>(null);
  let saving = $state(false);
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
    if (saving || !newName.trim()) return;
    saving = true;
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
    } finally {
      saving = false;
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
  <PageHeader title="" icon="" compact noBorder>
    {#snippet actions()}
      {#if labels.length > 0}
        <Button variant="primary" onclick={() => { editingLabel = null; newName = ''; newColor = '#444444'; showEditor = true; }} icon={Plus}>
          New Label
        </Button>
      {/if}
    {/snippet}
  </PageHeader>

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
        <Button variant="ghost" size="sm" onclick={() => showEditor = false} disabled={saving}>Cancel</Button>
        <Button variant="primary" size="sm" onclick={saveLabel} disabled={saving || !newName.trim()} loading={saving}>
          {editingLabel ? 'Update Label' : 'Create Label'}
        </Button>
      </div>
    </Card>
  {/if}

    {#if loading}
      <div class="loading-state">
        <div class="spinner"></div>
      </div>
    {:else if labels.length === 0}
      <div class="empty-state">
        <div class="empty-icon-ring">
          <Tag size={40} />
        </div>
        <h2>No labels found</h2>
        <p>There are no labels defined yet for this repository.</p>
        <div style="margin-top: 1rem;">
          <Button variant="primary" onclick={() => { editingLabel = null; newName = ''; newColor = '#0061a4'; showEditor = true; }} icon={Plus}>
            Create your first label
          </Button>
        </div>
      </div>
    {:else}
      <div class="flat-list-container">
        {#each labels as label}
          <div class="list-item-row">
            <div class="list-item-icon">
              <Tag size={18} style="color: var(--text-muted); opacity: 0.5;" />
            </div>
            <div class="item-main-content" style="justify-content: center;">
              <div>
                <span class="badge-mini" style="background: {label.color}; font-size: 0.85rem; padding: 0.25rem 0.75rem;">
                  {label.name}
                </span>
              </div>
            </div>
            <div class="item-actions-right">
              <button class="btn-icon-minimal" onclick={() => handleEdit(label)} title="Edit"><Edit size={16} /></button>
              <button class="btn-icon-minimal delete" onclick={() => deleteLabel(label.id)} title="Delete"><Trash2 size={16} /></button>
            </div>
          </div>
        {/each}
      </div>
    {/if}
  </div>
</div>

<style>
  .repo-labels { display: flex; flex-direction: column; gap: 0; }
  .labels-content { padding: 0; display: flex; flex-direction: column; gap: 1.25rem; }

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

  .empty-state {
    text-align: center;
    padding: 6rem 2rem;
    color: var(--text-muted);
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 1.25rem;
  }

  .empty-state h2 {
    margin: 0;
    font-size: 1.75rem;
    font-weight: 950;
    color: var(--text-color);
  }

  .empty-state p {
    max-width: 400px;
    margin: 0;
    font-size: 1.1rem;
    line-height: 1.5;
  }

  .empty-icon-ring {
    width: 80px;
    height: 80px;
    border: 2px dashed var(--md-sys-color-outline-variant);
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;
    margin-bottom: 1rem;
    color: var(--primary-color);
    background: rgba(var(--primary-rgb), 0.05);
  }
  
  .loading-state { display: flex; justify-content: center; padding: 3rem; }
  .spinner { width: 32px; height: 32px; border: 3px solid var(--md-sys-color-outline-variant); border-top-color: var(--primary-color); border-radius: 50%; animation: spin 1s infinite linear; }
  @keyframes spin { to { transform: rotate(360deg); } }
</style>
