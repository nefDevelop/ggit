<script lang="ts">
  import { apiFetch } from "$lib/api";
  import { onMount, getContext } from 'svelte';
  import { dialog } from '$lib/dialog.svelte';
  import Avatar from './Avatar.svelte';
  import Button from './ui/Button.svelte';
  import Input from './ui/Input.svelte';
  import Select from './ui/Select.svelte';
  import Card from './ui/Card.svelte';
  
  const repo: any = getContext('repo');
  let collaborators = $state<any[]>([]);
  let loading = $state(true);
  let newCollaborator = $state('');
  let submitting = $state(false);
  let message = $state('');

  const roleOptions = [
    { value: 1, label: 'Read' },
    { value: 2, label: 'Write' },
    { value: 3, label: 'Admin' }
  ];

  async function fetchCollaborators() {
    loading = true;
    try {
      const res = await apiFetch(`/api/v2/${repo.data.owner.username}/${repo.data.name}/settings/collaborators`);
      if (res.ok) {
        collaborators = await res.json();
      }
    } catch (e) {
      console.error(e);
    } finally {
      loading = false;
    }
  }

  async function addCollaborator() {
    if (!newCollaborator) return;
    submitting = true;
    message = '';
    try {
      const res = await apiFetch(`/api/v2/${repo.data.owner.username}/${repo.data.name}/settings/collaborators/${newCollaborator}`, {
        method: 'POST'
      });
      if (res.ok) {
        newCollaborator = '';
        fetchCollaborators();
      } else {
        const body = await res.json();
        message = body.message || 'User not found.';
      }
    } catch (e) {
      message = 'Failed to add collaborator.';
    } finally {
      submitting = false;
    }
  }

  async function updatePermission(username: string, mode: number) {
    try {
      const res = await apiFetch(`/api/v2/${repo.data.owner.username}/${repo.data.name}/settings/collaborators/${username}`, {
        method: 'PATCH',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ mode })
      });
      if (res.ok) {
        fetchCollaborators();
      }
    } catch (e) {
      console.error(e);
    }
  }

  async function removeCollaborator(username: string) {
    if (!await dialog.confirm(`Remove ${username} as collaborator?`)) return;
    try {
      const res = await apiFetch(`/api/v2/${repo.data.owner.username}/${repo.data.name}/settings/collaborators/${username}`, {
        method: 'DELETE'
      });
      if (res.ok) {
        fetchCollaborators();
      }
    } catch (e) {
      console.error(e);
    }
  }

  function getModeValue(col: any) {
    if (col.permissions.admin) return 3;
    if (col.permissions.push) return 2;
    return 1;
  }

  onMount(fetchCollaborators);
</script>

<div class="collaboration-view">
  <Card class="add-collab-card tonal" padding="lg">
    <div class="add-header">
      <span class="material-symbols-rounded icon-hero">person_add</span>
      <div class="text">
        <h3>Add Collaborator</h3>
        <p>Invite users to contribute to this repository with specific permissions.</p>
      </div>
    </div>

    <form class="add-form" onsubmit={(e) => { e.preventDefault(); addCollaborator(); }}>
      <div class="input-row">
        <Input 
          bind:value={newCollaborator} 
          placeholder="Search by username..." 
          icon="search"
          error={message}
        />
        <Button 
          type="submit" 
          variant="primary" 
          disabled={submitting || !newCollaborator}
          loading={submitting}
        >
          Add
        </Button>
      </div>
    </form>
  </Card>

  <div class="collab-list-container">
    <div class="list-header">
      <span class="count">{collaborators.length} Collaborators</span>
    </div>

    <div class="collab-grid">
      {#each collaborators as col}
        <div class="collab-island hover-lift">
          <div class="collab-main">
            <Avatar src={col.avatar_url} username={col.username} size={48} />
            <div class="info">
              <span class="name">{col.full_name || col.username}</span>
              <span class="handle">@{col.username}</span>
            </div>
          </div>

          <div class="collab-actions">
            <div class="role-selector">
              <Select 
                options={roleOptions} 
                value={getModeValue(col)} 
                onchange={(val) => updatePermission(col.username, val)}
              />
            </div>
            <Button 
              variant="ghost" 
              icon="person_remove" 
              class="danger-btn" 
              onclick={() => removeCollaborator(col.username)} 
            />
          </div>
        </div>
      {/each}

      {#if !loading && collaborators.length === 0}
        <div class="empty-collabs">
          <span class="material-symbols-rounded">group_off</span>
          <p>No collaborators yet. Start by adding someone above.</p>
        </div>
      {/if}

      {#if loading}
        {#each Array(3) as _}
          <div class="skeleton-island"></div>
        {/each}
      {/if}
    </div>
  </div>
</div>

<style>
  .collaboration-view {
    display: flex;
    flex-direction: column;
    gap: 3rem;
    animation: fadeIn 0.4s ease-out;
  }

  .add-header {
    display: flex;
    align-items: center;
    gap: 1.5rem;
    margin-bottom: 2rem;
  }
  .icon-hero {
    font-size: 2.5rem;
    color: var(--primary-color);
    background: rgba(var(--primary-rgb), 0.1);
    padding: 1rem;
    border-radius: 1.25rem;
  }
  .add-header h3 {
    margin: 0;
    font-size: 1.5rem;
    font-weight: 850;
  }
  .add-header p {
    margin: 0.25rem 0 0 0;
    color: var(--text-muted);
    font-size: 0.95rem;
  }

  .add-form .input-row {
    display: flex;
    gap: 1rem;
    align-items: flex-start;
  }
  :global(.add-form .input-row .input-container) {
    flex: 1;
  }

  .list-header {
    margin-bottom: 1.5rem;
    padding-left: 0.5rem;
    border-bottom: 2px solid var(--md-sys-color-outline-variant);
    padding-bottom: 0.75rem;
  }
  .list-header .count {
    font-size: 0.75rem;
    font-weight: 900;
    text-transform: uppercase;
    letter-spacing: 0.1em;
    color: var(--primary-color);
    opacity: 0.8;
  }

  .collab-grid {
    display: flex;
    flex-direction: column;
    gap: 1rem;
  }

  .collab-island {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 1.25rem 1.75rem;
    background: var(--md3-surface-1);
    border: 1px solid var(--md-sys-color-outline-variant);
    border-radius: 1.5rem;
    transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  }
  .collab-island:hover {
    background: var(--surface-variant);
    border-color: var(--primary-color);
  }

  .collab-main {
    display: flex;
    align-items: center;
    gap: 1.25rem;
  }
  .info {
    display: flex;
    flex-direction: column;
  }
  .name {
    font-size: 1.1rem;
    font-weight: 800;
    color: var(--text-color);
  }
  .handle {
    font-size: 0.8rem;
    color: var(--text-muted);
    font-weight: 600;
  }

  .collab-actions {
    display: flex;
    align-items: center;
    gap: 1.5rem;
  }
  .role-selector {
    width: 140px;
  }

  :global(.danger-btn:hover) {
    color: var(--error) !important;
    background: rgba(var(--error-rgb), 0.1) !important;
  }

  .empty-collabs {
    padding: 5rem 0;
    text-align: center;
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 1rem;
    opacity: 0.5;
  }
  .empty-collabs span {
    font-size: 4rem;
    color: var(--primary-color);
  }
  .empty-collabs p {
    font-weight: 700;
    font-size: 1.1rem;
  }

  .skeleton-island {
    height: 80px;
    background: var(--surface-variant);
    border-radius: 1.5rem;
    animation: pulse 1.5s infinite;
  }

  @keyframes pulse {
    0%, 100% { opacity: 0.3; }
    50% { opacity: 0.6; }
  }
  @keyframes fadeIn {
    from { opacity: 0; transform: translateY(10px); }
    to { opacity: 1; transform: translateY(0); }
  }

  @media (max-width: 600px) {
    .collab-island {
      flex-direction: column;
      align-items: stretch;
      gap: 1.5rem;
    }
    .collab-actions {
      justify-content: space-between;
    }
    .role-selector {
      flex: 1;
    }
  }
</style>
