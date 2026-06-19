<script lang="ts">
  import type { RepoContext } from "$lib/types";
  import { apiFetch } from "$lib/api";
  import { onMount, getContext } from 'svelte';
  import { dialog } from '$lib/dialog.svelte';
  import Avatar from './Avatar.svelte';
  import Button from './ui/Button.svelte';
  import Input from './ui/Input.svelte';
  import Select from './ui/Select.svelte';
  import Card from './ui/Card.svelte';
  import PageHeader from './ui/PageHeader.svelte';
  import { i18n } from '$lib/i18n.svelte';
  
import Users from "lucide-svelte/icons/users";

  
  const repo = getContext<RepoContext>('repo');
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
    if (!await dialog.confirm(`¿Remover a ${username} como colaborador?`)) return;
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

<Card padding="none" class="settings-premium-card">
  <PageHeader 
    title="Colaboradores"
    subtitle="Gestiona quién tiene acceso a este repositorio y sus niveles de permiso."
    icon={Users}
  />

  <div class="settings-content">
    <div class="add-collab-box island">
      <div class="add-header-minimal">
        <span class="material-symbols-rounded">person_add</span>
        <strong>Añadir Colaborador</strong>
      </div>
      <form class="add-form" onsubmit={(e) => { e.preventDefault(); addCollaborator(); }}>
        <div class="input-row">
          <Input 
            bind:value={newCollaborator} 
            placeholder="Buscar por nombre de usuario..." 
            icon="search"
            error={message}
          />
          <Button 
            type="submit" 
            variant="primary" 
            disabled={submitting || !newCollaborator}
            loading={submitting}
          >
            Añadir
          </Button>
        </div>
      </form>
    </div>

    <div class="collab-list-header">
      <span class="count-tag">{collaborators.length} Colaboradores</span>
    </div>

    <div class="collab-stack">
      {#if loading}
        {#each Array(3) as _}
          <div class="skeleton-tile tonal"></div>
        {/each}
      {:else}
        {#each collaborators as col}
          <div class="collab-tile tonal hover-lift">
            <div class="collab-main">
              <Avatar src={col.avatar_url} username={col.username} size={44} class="collab-avatar" />
              <div class="collab-info">
                <span class="collab-name">{col.full_name || col.username}</span>
                <span class="collab-handle">@{col.username}</span>
              </div>
            </div>

            <div class="collab-actions">
              <div class="role-wrap">
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

        {#if collaborators.length === 0}
          <div class="empty-state">
            <span class="material-symbols-rounded">group_off</span>
            <p>Este repositorio no tiene colaboradores externos.</p>
          </div>
        {/if}
      {/if}
    </div>
  </div>
</Card>

<style>
  .settings-content { padding: 2.5rem; display: flex; flex-direction: column; gap: 2.5rem; }
  @media (max-width: 600px) { .settings-content { padding: 1.5rem; } }

  .add-collab-box { padding: 1.5rem; border-radius: 1.5rem; border: 1px solid var(--md-sys-color-outline-variant); }
  .add-header-minimal { display: flex; align-items: center; gap: 0.75rem; margin-bottom: 1.25rem; color: var(--primary-color); }
  .add-header-minimal span { font-size: 1.5rem; }
  .add-header-minimal strong { font-weight: 900; text-transform: uppercase; letter-spacing: 0.05em; font-size: 0.85rem; }

  .add-form .input-row { display: flex; gap: 1rem; align-items: flex-start; }
  :global(.add-form .input-row .input-container) { flex: 1; }

  .collab-list-header { border-bottom: 2px solid var(--md-sys-color-outline-variant); padding-bottom: 0.75rem; }
  .count-tag { font-size: 0.75rem; font-weight: 950; text-transform: uppercase; color: var(--primary-color); letter-spacing: 0.1em; opacity: 0.8; }

  .collab-stack { display: flex; flex-direction: column; gap: 1rem; }

  .collab-tile { display: flex; justify-content: space-between; align-items: center; padding: 1.25rem 1.75rem; border-radius: 1.5rem; border: 1px solid var(--md-sys-color-outline-variant); transition: all 0.3s; }
  .collab-tile:hover { transform: translateX(4px); border-color: var(--primary-color); background: rgba(var(--primary-rgb), 0.02); }

  .collab-main { display: flex; align-items: center; gap: 1.25rem; }
  :global(.collab-avatar) { border: 1px solid var(--md-sys-color-outline-variant); }
  
  .collab-info { display: flex; flex-direction: column; }
  .collab-name { font-size: 1.1rem; font-weight: 850; color: var(--text-color); }
  .collab-handle { font-size: 0.85rem; color: var(--text-muted); font-weight: 600; }

  .collab-actions { display: flex; align-items: center; gap: 1.5rem; }
  .role-wrap { width: 140px; }

  :global(.danger-btn:hover) { color: var(--error) !important; background: rgba(var(--error-rgb), 0.1) !important; }

  .empty-state { padding: 4rem 2rem; text-align: center; color: var(--text-muted); display: flex; flex-direction: column; align-items: center; gap: 1rem; }
  .empty-state span { font-size: 3rem; opacity: 0.3; }
  .empty-state p { font-weight: 700; }

  .skeleton-tile { height: 84px; border-radius: 1.5rem; animation: pulse 1.5s infinite; }
  @keyframes pulse { 50% { opacity: 0.5; } }

  @media (max-width: 600px) {
    .collab-tile { flex-direction: column; align-items: stretch; gap: 1.5rem; }
    .collab-actions { justify-content: space-between; border-top: 1px solid var(--md-sys-color-outline-variant); padding-top: 1rem; }
    .role-wrap { flex: 1; }
  }
</style>
