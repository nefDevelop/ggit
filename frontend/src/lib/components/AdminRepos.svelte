<script lang="ts">
  import { apiFetch } from "$lib/api";
  import { onMount } from "svelte";
  import Card from './ui/Card.svelte';
  import Button from './ui/Button.svelte';
  import Avatar from './Avatar.svelte';
  import { i18n } from '$lib/i18n.svelte';

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
    if (!confirm(`Are you sure you want to delete repository "${name}"? This action cannot be undone.`)) return;

    try {
      const res = await apiFetch(`/api/v2/admin/repos/${id}`, { method: "DELETE" });
      if (res.ok) {
        repos = repos.filter((r) => r.id !== id);
      }
    } catch (e) {
      console.error(e);
    }
  }

  onMount(fetchRepos);
</script>

<Card class="settings-premium-card" padding="lg">
  <header class="section-header">
    <div class="title-with-icon">
      <span class="material-symbols-rounded hero-icon">database</span>
      <div>
        <h2>{i18n.t('admin.repositories')}</h2>
        <p class="subtitle">Supervisión y gestión global de repositorios del sistema</p>
      </div>
    </div>
  </header>

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
          {#each Array(5) as _}
            <tr class="skeleton-row">
              <td colspan="5"><div class="skeleton-line"></div></td>
            </tr>
          {/each}
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
                  <Button variant="ghost" size="sm" icon="delete" class="danger-btn" onclick={() => deleteRepo(repo.id, repo.name)} />
                </div>
              </td>
            </tr>
          {/each}
        {/if}
      </tbody>
    </table>
  </div>
</Card>

<style>
  :global(.settings-premium-card) { border-radius: 3rem; background: rgba(var(--surface-rgb), 0.5); backdrop-filter: blur(24px); border: 1px solid var(--md-sys-color-outline-variant); }
  
  .header-flex { display: flex; justify-content: space-between; align-items: center; gap: 2rem; }
  @media (max-width: 600px) {
    .header-flex { flex-direction: column; align-items: stretch; text-align: center; }
    .title-with-icon { flex-direction: column; gap: 1rem; }
  }

  .title-with-icon { display: flex; align-items: center; gap: 2rem; }
  .hero-icon { font-size: 3rem; color: var(--primary-color); background: rgba(var(--primary-rgb), 0.1); padding: 1.25rem; border-radius: 2rem; flex-shrink: 0; }

  h2 { margin: 0; font-size: 2.25rem; font-weight: 950; letter-spacing: -0.04em; }
  .subtitle { margin: 0.5rem 0 0 0; color: var(--text-muted); font-size: 1.1rem; }

  @media (max-width: 600px) {
    .section-header { margin-bottom: 2rem; }
    h2 { font-size: 1.5rem; }
    .hero-icon { font-size: 1.5rem; padding: 0.6rem; border-radius: 1rem; }
    .title-with-icon { gap: 0.75rem; }
    .subtitle { font-size: 0.9rem; }
  }

  .modern-table-card { border: 1px solid var(--md-sys-color-outline-variant); border-radius: 2.5rem; overflow: hidden; background: rgba(0,0,0,0.15); box-shadow: inset 0 2px 4px rgba(0,0,0,0.1); }

  .modern-table { width: 100%; border-collapse: collapse; text-align: left; }
  .modern-table th { padding: 1.5rem; background: rgba(var(--primary-rgb), 0.03); color: var(--text-muted); font-size: 0.8rem; text-transform: uppercase; letter-spacing: 0.15em; font-weight: 900; border-bottom: 1px solid var(--md-sys-color-outline-variant); }
  .modern-table td { padding: 1.75rem 1.5rem; border-bottom: 1px solid var(--md-sys-color-outline-variant); }
  .modern-table tr:last-child td { border-bottom: none; }

  @media (max-width: 1024px) {
    .modern-table-card { border-radius: 1.5rem; }
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

  .skeleton-line { height: 2rem; background: var(--surface-variant); border-radius: 0.75rem; animation: pulse 1.5s infinite; }
  @keyframes pulse { 0% { opacity: 0.3; } 50% { opacity: 0.6; } 100% { opacity: 0.3; } }
</style>
