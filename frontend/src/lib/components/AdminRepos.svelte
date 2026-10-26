<script lang="ts">
  

import Database from "lucide-svelte/icons/database";
import PlusSquare from "lucide-svelte/icons/plus-square";
import Lock from "lucide-svelte/icons/lock";
import BookText from "lucide-svelte/icons/book-text";
import Star from "lucide-svelte/icons/star";
import GitFork from "lucide-svelte/icons/git-fork";
import Eye from "lucide-svelte/icons/eye";
import ShieldCheck from "lucide-svelte/icons/shield-check";
import Trash2 from "lucide-svelte/icons/trash-2";


  import { apiFetch } from "$lib/api";
  import { toast } from "$lib/toast.svelte";
  import { onMount } from "svelte";
  import Card from './ui/Card.svelte';
  import Button from './ui/Button.svelte';
  import Avatar from './Avatar.svelte';
  import { i18n } from '$lib/i18n.svelte';
  import { dialog } from '$lib/dialog.svelte';
  import TableSkeleton from './ui/TableSkeleton.svelte';
  import PageHeader from './ui/PageHeader.svelte';

  let repos = $state<any[]>([]);
  let loading = $state(true);
  let deletingIds = $state<Set<number>>(new Set());
  let page = $state(1);

  async function fetchRepos() {
    loading = true;
    try {
      const res = await apiFetch(`/api/v2/admin/repos?page=${page}`);
      if (res.ok) {
        const json = await res.json();
        repos = Array.isArray(json) ? json : [];
      }
    } catch (e) {
      console.error(e);
    } finally {
      loading = false;
    }
  }

  async function deleteRepo(id: number, name: string) {
    if (!await dialog.confirm(`Are you sure you want to delete repository "${name}"? This action cannot be undone.`)) return;

    deletingIds.add(id);
    try {
      const res = await apiFetch(`/api/v2/admin/repos/${id}`, { method: "DELETE" });
      if (res.ok) {
        repos = repos.filter((r) => r.id !== id);
        toast.success(`Repository "${name}" deleted successfully.`);
      }
    } catch (e) {
      console.error(e);
    } finally {
      deletingIds.delete(id);
    }
  }

  async function checkIntegrity(id: number, name: string) {
    if (!await dialog.confirm(`¿Seguro que quieres forzar una verificación de integridad (git fsck) en "${name}"?`)) return;
    
    try {
      const res = await apiFetch(`/api/v2/admin/repos/${id}/check-integrity`, { method: "POST" });
      if (res.ok) {
        dialog.alert(`Verificación de integridad iniciada para "${name}".`);
      }
    } catch (e) { console.error(e); }
  }

  onMount(fetchRepos);
</script>

<div class="admin-repos-container">
  <PageHeader 
    title={i18n.t('admin.repositories') || "Repositories"}
    subtitle="Supervisión y gestión global de repositorios del sistema"
    icon={Database}
  >
    {#snippet actions()}
      <Button variant="primary" icon={PlusSquare} href="/admin/repos/adoption">Adopción de Repos</Button>
    {/snippet}
  </PageHeader>

  <div class="admin-repos-body">

  <div class="modern-table-card">
    <table class="modern-table">
      <thead>
        <tr>
          <th>ID</th>
          <th>{i18n.t('repository')}</th>
          <th>Colaboradores</th>
          <th>Estado</th>
          <th>Acciones</th>
        </tr>
      </thead>
      <tbody>
        {#if loading}
          <TableSkeleton columns={5} rows={5} />
        {:else}
          {#each repos as repo}
            <tr class="hover-row">
              <td class="id-cell" data-label="ID">#{repo.id}</td>
              <td data-label={i18n.t('repository')}>
                <div class="repo-info-cell">
                  <div class="repo-main">
                    {#if repo.is_private}
                      <Lock size={18} class="repo-icon" />
                    {:else}
                      <BookText size={18} class="repo-icon" />
                    {/if}
                    <a href="/{repo.owner?.username}/{repo.name}" class="repo-link">
                      <span class="owner">{repo.owner?.username}</span>
                      <span class="divider">/</span>
                      <span class="name">{repo.name}</span>
                    </a>
                  </div>
                  {#if repo.description}
                    <p class="repo-description">{repo.description}</p>
                  {/if}
                </div>
              </td>
              <td data-label="Colaboradores">
                <div class="repo-stats">
                  <div class="stat">
                    <Star size={16} />
                    {repo.num_stars}
                  </div>
                  <div class="stat">
                    <GitFork size={16} />
                    {repo.num_forks}
                  </div>
                </div>
              </td>
              <td data-label="Estado">
                <span class="badge" class:private={repo.is_private}>
                  {repo.is_private ? 'Privado' : 'Público'}
                </span>
              </td>
              <td data-label="Acciones">
                <div class="action-row">
                  <a class="icon-btn" href="/{repo.owner?.username}/{repo.name}" title="Ver">
                    <Eye size={18} />
                  </a>
                  <button class="icon-btn" onclick={() => checkIntegrity(repo.id, repo.name)} title="Verificar Integridad">
                    <ShieldCheck size={18} />
                  </button>
                  <button 
                    class="icon-btn danger" 
                    onclick={() => deleteRepo(repo.id, repo.name)} 
                    title="Eliminar"
                    disabled={deletingIds.has(repo.id)}
                  >
                    {#if deletingIds.has(repo.id)}
                      <div class="spinner-small"></div>
                    {:else}
                      <Trash2 size={18} />
                    {/if}
                  </button>
                </div>
              </td>
            </tr>
          {/each}
        {/if}
      </tbody>
    </table>
  </div>
  </div>
</div>

<style>
  .admin-repos-body {
    padding: 0;
    animation: fadeIn 0.4s cubic-bezier(0.4, 0, 0.2, 1);
  }

  @keyframes fadeIn {
    from { opacity: 0; transform: translateY(4px); }
    to { opacity: 1; transform: translateY(0); }
  }
  

  @media (max-width: 1024px) {
    .repo-info-cell { text-align: right; }
    .repo-main { justify-content: flex-end; }
    .repo-stats { justify-content: flex-end; }
    .action-row { justify-content: flex-end; }
  }


  .id-cell { font-family: monospace; color: var(--text-muted); font-weight: 700; opacity: 0.5; font-size: 0.9rem; }

  .repo-info-cell { display: flex; flex-direction: column; gap: 0.6rem; }
  .repo-main { display: flex; align-items: center; gap: 1rem; }
  :global(.repo-icon) { color: var(--primary-color); opacity: 0.7; }
  
  .repo-link { text-decoration: none; display: flex; align-items: center; gap: 0.25rem; font-size: 1.2rem; }
  .repo-link .owner { color: var(--text-muted); font-weight: 600; }
  .repo-link .divider { color: var(--md-sys-color-outline); opacity: 0.4; }
  .repo-link .name { color: var(--text-color); font-weight: 900; }
  .repo-link:hover .name { color: var(--primary-color); text-decoration: underline; }

  .repo-description { margin: 0; font-size: 0.95rem; color: var(--text-muted); line-height: 1.5; max-width: 500px; }

  .repo-stats { display: flex; gap: 2rem; }
  .stat { display: flex; align-items: center; gap: 0.5rem; font-size: 1rem; font-weight: 800; color: var(--text-muted); }

  .badge { padding: 0.5rem 1rem; border-radius: 0.75rem; font-size: 0.8rem; font-weight: 900; text-transform: uppercase; background: rgba(var(--primary-rgb), 0.12); color: var(--primary-color); border: 2px solid rgba(var(--primary-rgb), 0.2); }
  .badge.private { background: rgba(var(--error-rgb), 0.12); color: var(--error); border-color: rgba(var(--error-rgb), 0.2); }

  .action-row { display: flex; gap: 0.75rem; }
  .icon-btn {
    background: transparent;
    border: none;
    color: var(--text-muted);
    cursor: pointer;
    padding: 0.5rem;
    display: flex;
    align-items: center;
    justify-content: center;
    border-radius: 0.5rem;
    transition: all 0.2s;
    text-decoration: none;
  }
  .icon-btn:hover {
    background: var(--md-sys-color-surface-variant);
    color: var(--text-color);
  }
  .icon-btn.danger:hover {
    background: rgba(var(--error-rgb), 0.1);
    color: var(--error);
  }
  .icon-btn:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }

  .spinner-small {
    width: 16px;
    height: 16px;
    border: 2px solid currentColor;
    border-top-color: transparent;
    border-radius: 50%;
    animation: spin 0.8s linear infinite;
  }

  @keyframes spin {
    to { transform: rotate(360deg); }
  }
</style>
