<script lang="ts">
  import { apiFetch } from "$lib/api";
  import { onMount, getContext } from "svelte";
  import Button from "./ui/Button.svelte";
  import Card from "./ui/Card.svelte";
  import { dialog } from "$lib/dialog.svelte";
  import { Trash2, Plus, Globe, Shield, RefreshCw, List, Cloud, Info } from "lucide-svelte";
  import PageHeader from "./ui/PageHeader.svelte";
  import Input from "./ui/Input.svelte";

  const repo: any = getContext("repo");
  let remotes = $state<string[]>([]);
  let loading = $state(true);
  let adding = $state(false);
  let deleting = $state<string | null>(null);

  let newName = $state("");
  let newURL = $state("");
  let message = $state("");
  let error = $state("");

  async function fetchRemotes() {
    loading = true;
    try {
      const res = await apiFetch(
        `/api/v2/${repo.data.owner.username}/${repo.data.name}/settings/remotes`
      );
      if (res.ok) remotes = await res.json();
    } catch (e) {
      console.error(e);
    } finally {
      loading = false;
    }
  }

  async function addRemote() {
    if (!newName || !newURL) return;
    adding = true;
    error = "";
    try {
      const res = await apiFetch(
        `/api/v2/${repo.data.owner.username}/${repo.data.name}/settings/remotes`,
        {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({ name: newName, url: newURL }),
        }
      );
      if (res.ok) {
        newName = "";
        newURL = "";
        message = "Remoto añadido correctamente";
        setTimeout(() => (message = ""), 3000);
        fetchRemotes();
      } else {
        const data = await res.json();
        error = data.message || "Error al añadir remoto";
      }
    } catch (e) {
      error = "Error de conexión";
    } finally {
      adding = false;
    }
  }

  async function deleteRemote(name: string) {
    if (!await dialog.confirm(`¿Estás seguro de eliminar el remoto "${name}"?`)) return;
    deleting = name;
    try {
      const res = await apiFetch(
        `/api/v2/${repo.data.owner.username}/${repo.data.name}/settings/remotes/${name}`,
        { method: "DELETE" }
      );
      if (res.ok) {
        fetchRemotes();
      }
    } catch (e) {
      console.error(e);
    } finally {
      deleting = null;
    }
  }

  onMount(fetchRemotes);
</script>

<div class="remotes-settings-wrapper">
  <Card padding="none" class="settings-premium-card">
    <PageHeader 
      title="Gestionar Remotos"
      subtitle="Configura destinos adicionales para sincronizar tu código con otros servicios."
      icon={RefreshCw}
    />

    <div class="settings-content">
      <div class="add-remote-section island">
        <header class="section-meta">
          <strong>Añadir Nuevo Remoto</strong>
          <p>Vincula este repositorio con otro servidor externo.</p>
        </header>

        <form class="form-grid" onsubmit={(e) => { e.preventDefault(); addRemote(); }}>
          <Input 
            label="Nombre" 
            bind:value={newName} 
            placeholder="ej: github, backup" 
          />
          <Input 
            label="URL del Repositorio" 
            bind:value={newURL} 
            placeholder="https://github.com/usuario/repo.git" 
            class="grow"
          />
          <div class="form-actions-inline">
            <Button 
              variant="primary" 
              type="submit"
              disabled={adding || !newName || !newURL}
              icon={Plus}
              loading={adding}
            >
              Añadir Remoto
            </Button>
          </div>
        </form>
        {#if error}<p class="error-text">{error}</p>{/if}
        {#if message}<p class="success-text">{message}</p>{/if}
      </div>

      <div class="remotes-list-section">
        <header class="list-header">
          <List size={24} />
          <strong>Remotos Configurados</strong>
        </header>

        {#if loading}
          <div class="skeleton-list"></div>
        {:else if remotes.length === 0}
          <div class="empty-state">
            <Globe size={48} />
            <h3>No hay remotos configurados</h3>
            <p>Este repositorio solo existe localmente en Ggit.</p>
          </div>
        {:else}
          <div class="remotes-grid">
            {#each remotes as remote}
              <div class="remote-item tonal">
                <div class="remote-main">
                  <div class="remote-icon-wrap">
                    <Cloud size={24} />
                  </div>
                  <div class="remote-details">
                    <span class="remote-name">{remote}</span>
                    <span class="remote-url-hint">Sincronización bidireccional</span>
                  </div>
                </div>
                <Button 
                  variant="ghost" 
                  size="sm"
                  class="danger-text"
                  icon={Trash2}
                  onclick={() => deleteRemote(remote)}
                  disabled={deleting === remote}
                />
              </div>
            {/each}
          </div>
        {/if}
      </div>
    </div>
  </Card>

  <div class="info-zone">
    <Card class="tonal warning-card" padding="lg">
      <div class="info-content">
        <Info size={32} class="warning-icon" />
        <div class="info-text">
          <strong>Nota sobre Autenticación</strong>
          <p>Actualmente, el servidor de Ggit debe tener acceso pre-configurado (vía SSH o tokens en el entorno) para interactuar con estos remotos externos.</p>
        </div>
      </div>
    </Card>
  </div>
</div>

<style>
  .settings-content { padding: 2.5rem; display: flex; flex-direction: column; gap: 3rem; }
  @media (max-width: 1300px) { .settings-content { padding: 1.5rem; } }
  
  .add-remote-section { padding: 2rem; border-radius: 2rem; display: flex; flex-direction: column; gap: 2rem; }
  .section-meta strong { font-size: 1.25rem; font-weight: 900; color: var(--text-color); }
  .section-meta p { margin: 0.25rem 0 0; color: var(--text-muted); font-size: 0.95rem; }

  .form-grid { display: grid; grid-template-columns: 1fr 2fr; gap: 1.5rem; align-items: flex-end; }
  .form-actions-inline { grid-column: span 2; display: flex; justify-content: flex-end; padding-top: 1rem; }

  @media (max-width: 1300px) {
    .form-grid { grid-template-columns: 1fr; }
    .form-actions-inline { grid-column: span 1; }
  }

  .remotes-list-section { display: flex; flex-direction: column; gap: 1.5rem; }
  .list-header { display: flex; align-items: center; gap: 0.75rem; color: var(--primary-color); }
  .list-header strong { font-weight: 950; font-size: 0.8rem; text-transform: uppercase; letter-spacing: 0.1em; }

  .remotes-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(350px, 1fr)); gap: 1.25rem; }
  @media (max-width: 1100px) { .remotes-grid { grid-template-columns: 1fr; } }

  .remote-item { display: flex; align-items: center; justify-content: space-between; padding: 1.25rem 1.5rem; border-radius: 1.5rem; border: 1px solid var(--md-sys-color-outline-variant); transition: all 0.2s; }
  .remote-item:hover { border-color: var(--primary-color); background: rgba(var(--primary-rgb), 0.05); }

  .remote-main { display: flex; align-items: center; gap: 1.25rem; }
  .remote-icon-wrap { width: 44px; height: 44px; border-radius: 12px; background: rgba(var(--primary-rgb), 0.1); color: var(--primary-color); display: flex; align-items: center; justify-content: center; }

  .remote-details { display: flex; flex-direction: column; }
  .remote-name { font-weight: 800; font-size: 1.1rem; color: var(--text-color); }
  .remote-url-hint { font-size: 0.75rem; color: var(--text-muted); opacity: 0.8; }

  .empty-state { padding: 4rem 2rem; text-align: center; color: var(--text-muted); display: flex; flex-direction: column; align-items: center; gap: 1rem; opacity: 0.6; }

  .error-text { color: var(--error); font-size: 0.85rem; margin-top: 1rem; font-weight: 600; }
  .success-text { color: var(--success); font-size: 0.85rem; margin-top: 1rem; font-weight: 600; }

  .info-zone { }
  .info-content { display: flex; gap: 1.5rem; align-items: flex-start; }
  :global(.warning-icon) { color: var(--primary-color); opacity: 0.8; }
  .info-text strong { display: block; margin-bottom: 0.25rem; font-size: 1.1rem; font-weight: 800; }
  .info-text p { margin: 0; font-size: 0.9rem; color: var(--text-muted); line-height: 1.5; }

  .skeleton-list { height: 200px; background: var(--md-sys-color-surface-variant); border-radius: 2rem; animation: pulse 1.5s infinite; }
  @keyframes pulse { 50% { opacity: 0.5; } }
</style>
