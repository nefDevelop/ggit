<script lang="ts">
  import { apiFetch } from "$lib/api";
  import { onMount } from "svelte";
  import Card from '$lib/components/ui/Card.svelte';
  import Button from '$lib/components/ui/Button.svelte';
  import { i18n } from '$lib/i18n.svelte';
  import { dialog } from '$lib/dialog.svelte';
  import TableSkeleton from '$lib/components/ui/TableSkeleton.svelte';

  let unmanaged = $state<any[]>([]);
  let loading = $state(true);
  let adopting = $state<string | null>(null);

  async function fetchUnmanaged() {
    loading = true;
    try {
      const res = await apiFetch('/api/v2/admin/repos/unmanaged');
      if (res.ok) {
        unmanaged = await res.json();
      }
    } catch (e) {
      console.error(e);
    } finally {
      loading = false;
    }
  }

  async function adoptRepo(owner: string, name: string) {
    if (!await dialog.confirm(`¿Seguro que quieres adoptar el repositorio "${owner}/${name}"? Se vinculará a la base de datos de Ggit.`)) return;

    adopting = `${owner}/${name}`;
    try {
      const res = await apiFetch('/api/v2/admin/repos/adopt', {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ owner_name: owner, repo_name: name })
      });
      if (res.ok) {
        unmanaged = unmanaged.filter((r) => !(r.owner_name === owner && r.repo_name === name));
        dialog.alert(`Repositorio "${owner}/${name}" adoptado con éxito.`);
      } else {
        const data = await res.json();
        dialog.alert(`Error al adoptar: ${data.error || 'Desconocido'}`);
      }
    } catch (e) {
      console.error(e);
      dialog.alert("Error de red al intentar adoptar el repositorio.");
    } finally {
      adopting = null;
    }
  }

  onMount(fetchUnmanaged);
</script>

<svelte:head>
  <title>Adopción de Repositorios · Admin · Ggit</title>
</svelte:head>

<div class="adoption-view">
  <div class="breadcrumb">
    <a href="/admin/repos">Repositorios</a>
    <span class="material-symbols-rounded">chevron_right</span>
    <span>Adopción</span>
  </div>

  <Card class="settings-premium-card" padding="lg">
    <header class="section-header">
      <div class="title-with-icon">
        <span class="material-symbols-rounded hero-icon">folder_zip</span>
        <div>
          <h2>Adopción de Repositorios</h2>
          <p class="subtitle">Se han detectado carpetas Git en el disco que no están registradas en la base de datos.</p>
        </div>
      </div>
    </header>

    <div class="modern-table-card">
      <table class="modern-table">
        <thead>
          <tr>
            <th>Propietario</th>
            <th>Nombre del Repo</th>
            <th>Ruta en Disco</th>
            <th>Acciones</th>
          </tr>
        </thead>
        <tbody>
          {#if loading}
            <TableSkeleton columns={4} rows={5} />
          {:else if unmanaged.length === 0}
            <tr>
              <td colspan="4" class="empty-state">
                <div class="empty-msg">
                  <span class="material-symbols-rounded">check_circle</span>
                  <p>No se encontraron repositorios huérfanos. ¡Todo está en orden!</p>
                </div>
              </td>
            </tr>
          {:else}
            {#each unmanaged as repo}
              <tr class="hover-row">
                <td data-label="Propietario">
                  <div class="owner-chip">
                    <span class="material-symbols-rounded">person</span>
                    {repo.owner_name}
                  </div>
                </td>
                <td data-label="Nombre del Repo">
                  <strong class="repo-name-text">{repo.repo_name}</strong>
                </td>
                <td data-label="Ruta en Disco">
                  <code class="path-code">{repo.path}</code>
                </td>
                <td data-label="Acciones">
                  <Button 
                    variant="primary" 
                    size="sm" 
                    icon="link"
                    onclick={() => adoptRepo(repo.owner_name, repo.repo_name)}
                    disabled={adopting !== null}
                  >
                    {adopting === `${repo.owner_name}/${repo.repo_name}` ? 'Adoptando...' : 'Adoptar'}
                  </Button>
                </td>
              </tr>
            {/each}
          {/if}
        </tbody>
      </table>
    </div>
  </Card>
</div>

<style>
  .adoption-view { display: flex; flex-direction: column; gap: 2rem; }
  
  .breadcrumb { display: flex; align-items: center; gap: 0.5rem; color: var(--text-muted); font-weight: 700; font-size: 0.9rem; margin-bottom: -1rem; }
  .breadcrumb a { color: var(--primary-color); text-decoration: none; }
  .breadcrumb a:hover { text-decoration: underline; }
  .breadcrumb span { font-size: 1.2rem; }

  .section-header { margin-bottom: 3rem; }
  .title-with-icon { display: flex; align-items: center; gap: 2rem; }
  .hero-icon { font-size: 3rem; color: var(--primary-color); background: rgba(var(--primary-rgb), 0.1); padding: 1.25rem; border-radius: 2rem; flex-shrink: 0; }

  h2 { margin: 0; font-size: 2.25rem; font-weight: 950; letter-spacing: -0.04em; }
  .subtitle { margin: 0.5rem 0 0 0; color: var(--text-muted); font-size: 1.1rem; }

  .owner-chip { display: flex; align-items: center; gap: 0.5rem; font-weight: 800; color: var(--text-color); background: var(--md-sys-color-surface-variant); padding: 0.4rem 0.8rem; border-radius: 0.75rem; width: fit-content; }
  .owner-chip span { font-size: 1.1rem; opacity: 0.6; }

  .repo-name-text { font-size: 1.1rem; color: var(--primary-color); font-weight: 900; }
  
  .path-code { font-family: 'JetBrains Mono', monospace; font-size: 0.85rem; color: var(--text-muted); background: rgba(0,0,0,0.2); padding: 0.4rem 0.75rem; border-radius: 0.5rem; border: 1px solid var(--md-sys-color-outline-variant); }

  .empty-state { padding: 5rem 0; text-align: center; }
  .empty-msg { display: flex; flex-direction: column; align-items: center; gap: 1.5rem; color: var(--text-muted); }
  .empty-msg span { font-size: 4rem; color: var(--primary-color); opacity: 0.4; }
  .empty-msg p { font-size: 1.2rem; font-weight: 700; margin: 0; }

  @media (max-width: 800px) {
    .path-code { display: none; }
  }
</style>
