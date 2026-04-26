<script lang="ts">
  import { onMount } from "svelte";
  import RepoRow from "./RepoRow.svelte";
  import { apiFetch } from "$lib/api";
  import { i18n } from '$lib/i18n.svelte';
  import { dialog } from '$lib/dialog.svelte';
  import Card from './ui/Card.svelte';
  import Button from './ui/Button.svelte';

  let repos = $state<any[]>([]);
  let loading = $state(true);

  async function fetchRepos() {
    loading = true;
    try {
      const res = await apiFetch("/api/v2/user/settings/repos");
      if (res.ok) repos = await res.json();
    } catch (e) {
      console.error(e);
    } finally {
      loading = false;
    }
  }

  async function leaveRepo(id: number, name: string) {
    if (!await dialog.confirm(`¿Estás seguro de que quieres dejar de colaborar en ${name}?`)) return;
    try {
      const res = await apiFetch(`/api/v2/user/settings/repos/${id}`, { method: "DELETE" });
      if (res.ok) await fetchRepos();
    } catch (e) { console.error(e); }
  }

  onMount(fetchRepos);
</script>

<Card class="settings-premium-card" padding="lg">
  <header class="section-header">
    <div class="title-with-icon">
      <span class="material-symbols-rounded hero-icon">book</span>
      <div>
        <h2>Mis Repositorios</h2>
        <p class="subtitle">Lista de todos tus proyectos personales y colaboraciones activas</p>
      </div>
    </div>
  </header>

  <div class="modular-stack">
    <div class="config-module">
      <div class="module-tag">
        <span class="material-symbols-rounded">inventory_2</span>
        Inventario de Código
      </div>
      <div class="island tonal-list">
        {#if loading}
          <div class="skeleton-repos"></div>
        {:else if repos.length === 0}
          <div class="empty-repos-box">
             <span class="material-symbols-rounded">folder_open</span>
             <p>Aún no tienes repositorios vinculados.</p>
             <Button variant="primary" href="/repo/create" icon="add">Crear Primer Repositorio</Button>
          </div>
        {:else}
          {#each repos as repo}
            <div class="repo-modular-row">
              <div class="repo-main-content">
                 <RepoRow {repo} />
              </div>
              {#if repo.owner?.username && repo.owner.username !== repo.context_user}
                <div class="repo-actions">
                  <Button 
                    variant="ghost" 
                    size="sm" 
                    icon="logout" 
                    class="danger-btn" 
                    onclick={() => leaveRepo(repo.id, repo.full_name)}
                  >Dejar Colaboración</Button>
                </div>
              {/if}
            </div>
          {/each}
        {/if}
      </div>
    </div>
  </div>
</Card>

<style>
  
  .section-header { margin-bottom: 4rem; }
  .title-with-icon { display: flex; align-items: center; gap: 2rem; }
  .hero-icon { font-size: 3rem; color: var(--primary-color); background: rgba(var(--primary-rgb), 0.1); padding: 1.25rem; border-radius: 2rem; }
  h2 { margin: 0; font-size: 2.25rem; font-weight: 950; letter-spacing: -0.04em; }
  .subtitle { margin: 0.5rem 0 0 0; color: var(--text-muted); font-size: 1.1rem; }

  @media (max-width: 600px) {
    .section-header { margin-bottom: 2rem; }
    h2 { font-size: 1.5rem; }
    .hero-icon { font-size: 1.5rem; padding: 0.6rem; border-radius: 1rem; }
    .title-with-icon { gap: 0.75rem; }
    .subtitle { font-size: 0.9rem; }
  }

  .modular-stack { display: flex; flex-direction: column; gap: 4rem; }
  @media (max-width: 600px) {
    .modular-stack { gap: 2rem; }
  }
  .config-module { display: flex; flex-direction: column; gap: 1.5rem; }

  .module-tag { font-size: 0.8rem; font-weight: 950; text-transform: uppercase; color: var(--primary-color); letter-spacing: 0.15em; display: flex; align-items: center; gap: 0.75rem; border-bottom: 2px solid var(--md-sys-color-outline-variant); padding-bottom: 0.75rem; width: fit-content; opacity: 0.8; }
  .module-tag span { font-size: 1.25rem; }

  .tonal-list { display: flex; flex-direction: column; background: rgba(0,0,0,0.15); border-radius: 2.25rem; overflow: hidden; border: 1px solid var(--md-sys-color-outline-variant); }
  @media (max-width: 600px) {
    .tonal-list { border-radius: 1.25rem; }
  }
  
  .repo-modular-row { display: flex; align-items: center; justify-content: space-between; padding: 2rem 2.5rem; border-bottom: 1px solid var(--md-sys-color-outline-variant); gap: 3rem; transition: background 0.3s; }
  @media (max-width: 600px) { .repo-modular-row { flex-direction: column; align-items: stretch; padding: 1rem 1.25rem; gap: 1rem; } }
  .repo-modular-row:last-child { border-bottom: none; }
  .repo-modular-row:hover { background: rgba(255,255,255,0.02); }

  .repo-main-content { flex: 1; }
  .repo-actions { flex-shrink: 0; }
  @media (max-width: 600px) {
    .repo-actions { display: flex; justify-content: flex-end; }
  }

  :global(.danger-btn:hover) { color: var(--error) !important; background: rgba(var(--error-rgb), 0.1) !important; }

  .empty-repos-box { padding: 5rem; text-align: center; color: var(--text-muted); display: flex; flex-direction: column; align-items: center; gap: 1.5rem; }
  .empty-repos-box span { font-size: 4rem; opacity: 0.2; }

  .skeleton-repos { height: 300px; background: var(--surface-variant); border-radius: 2rem; animation: pulse 1.5s infinite; }
</style>
