<script lang="ts">
  import { apiFetch } from "$lib/api";
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
  let page = $state(1);

  async function fetchRepos() {
    loading = true;
    try {
      const res = await apiFetch(`/api/v2/admin/repos?page=${page}`);
      if (res.ok) {
        repos = await res.json();
      }
    } catch (e) {
      console.error(e);
    } finally {
      loading = false;
    }
  }

  async function deleteRepo(id: number, name: string) {
    if (!await dialog.confirm(`Are you sure you want to delete repository "${name}"? This action cannot be undone.`)) return;

    try {
      const res = await apiFetch(`/api/v2/admin/repos/${id}`, { method: "DELETE" });
      if (res.ok) {
        repos = repos.filter((r) => r.id !== id);
      }
    } catch (e) {
      console.error(e);
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
    icon="database"
  >
    {#snippet actions()}
      <Button variant="primary" icon="add_box" href="/admin/repos/adoption">Adopción de Repos</Button>
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
                    <span class="material-symbols-rounded repo-icon">{repo.is_private ? 'lock' : 'book'}</span>
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
                    <span class="material-symbols-rounded">star</span>
                    {repo.num_stars}
                  </div>
                  <div class="stat">
                    <span class="material-symbols-rounded">fork_right</span>
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
                  <Button variant="ghost" size="sm" icon="visibility" href="/{repo.owner?.username}/{repo.name}" />
                  <Button variant="ghost" size="sm" icon="health_and_safety" onclick={() => checkIntegrity(repo.id, repo.name)} title="Verificar Integridad" />
                  <Button variant="ghost" size="sm" icon="delete" class="danger-btn" onclick={() => deleteRepo(repo.id, repo.name)} />
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
    padding: 0 1.5rem 4rem;
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
  .repo-icon { font-size: 1.5rem; color: var(--primary-color); opacity: 0.7; }
  
  .repo-link { text-decoration: none; display: flex; align-items: center; gap: 0.25rem; font-size: 1.2rem; }
  .repo-link .owner { color: var(--text-muted); font-weight: 600; }
  .repo-link .divider { color: var(--md-sys-color-outline); opacity: 0.4; }
  .repo-link .name { color: var(--text-color); font-weight: 900; }
  .repo-link:hover .name { color: var(--primary-color); text-decoration: underline; }

  .repo-description { margin: 0; font-size: 0.95rem; color: var(--text-muted); line-height: 1.5; max-width: 500px; }

  .repo-stats { display: flex; gap: 2rem; }
  .stat { display: flex; align-items: center; gap: 0.5rem; font-size: 1rem; font-weight: 800; color: var(--text-muted); }
  .stat span { font-size: 1.25rem; color: var(--primary-color); }

  .badge { padding: 0.5rem 1rem; border-radius: 1rem; font-size: 0.8rem; font-weight: 900; text-transform: uppercase; background: rgba(var(--primary-rgb), 0.12); color: var(--primary-color); border: 2px solid rgba(var(--primary-rgb), 0.2); }
  .badge.private { background: rgba(var(--error-rgb), 0.12); color: var(--error); border-color: rgba(var(--error-rgb), 0.2); }

  .action-row { display: flex; gap: 0.75rem; }
  :global(.danger-btn:hover) { color: var(--error) !important; background: rgba(var(--error-rgb), 0.15) !important; }
</style>
